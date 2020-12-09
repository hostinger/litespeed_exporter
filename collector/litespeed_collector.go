package collector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

// LitespeedCollectorOpts carries the options used in LitespeedCollector
type LitespeedCollectorOpts struct {
	FilePattern     string
	ReqRatesByHost  bool
	MetricsByCore   bool
	ExcludeExtapp   bool
	ExcludedMetrics map[string]bool
}

// LitespeedCollector collects LiteSpeed stats from the given files and exports them as Prometheus metrics
type LitespeedCollector struct {
	mutex                        sync.RWMutex
	options                      LitespeedCollectorOpts
	totalScrapes, scrapeFailures prometheus.Counter
	logger                       log.Logger
}

// NewLitespeedCollector returns constructed collector
func NewLitespeedCollector(opts LitespeedCollectorOpts, logger log.Logger) *LitespeedCollector {
	return &LitespeedCollector{
		options: opts,
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "exporter_scrapes_total",
			Help:      "Current total LiteSpeed scrapes.",
		}),
		scrapeFailures: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "exporter_scrape_failures_total",
			Help:      "Number of errors while scraping files.",
		}),
		logger: logger,
	}
}

func (c *LitespeedCollector) metricIsTracked(flag string) bool {
	_, ok := c.options.ExcludedMetrics[flag]
	return !ok
}

// Describe describes all the metrics that can be exported by the LiteSpeed exporter
func (c *LitespeedCollector) Describe(ch chan<- *prometheus.Desc) {
	for flag, metric := range LitespeedMetrics {
		if c.metricIsTracked(flag) {
			ch <- metric.Desc
		}
	}
	ch <- litespeedVersion
	ch <- litespeedUp
	ch <- c.totalScrapes.Desc()
	ch <- c.scrapeFailures.Desc()
}

// Collect fetches the stats from target files and delivers them as Prometheus metrics
func (c *LitespeedCollector) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	up := c.collectReports(ch)

	ch <- prometheus.MustNewConstMetric(litespeedUp, prometheus.GaugeValue, up)
	ch <- c.totalScrapes
	ch <- c.scrapeFailures
}

func (c *LitespeedCollector) collectReports(ch chan<- prometheus.Metric) (up float64) {
	c.totalScrapes.Inc()

	reports, err := c.scrapeReports(c.options.FilePattern)
	if err != nil {
		c.scrapeFailures.Inc()
		return 0
	}

	versionScraped := false

	for core, report := range reports {
		if !versionScraped {
			ch <- prometheus.MustNewConstMetric(litespeedVersion, prometheus.GaugeValue, 1, report.GeneralInfo.Version)
			versionScraped = true
		}

		c.collectGeneralInfoMetrics(core, report.GeneralInfo, ch)
		c.collectReqRateMetrics(core, report.ReqRates, ch)
		c.collectExtAppMetrics(core, report.ExtApps, ch)
	}

	return 1
}

func (c *LitespeedCollector) collectGeneralInfoMetrics(core string, generalInfo generalInfoReport, ch chan<- prometheus.Metric) {
	for flag, value := range generalInfo.KeyValues {
		if metric, ok := LitespeedMetrics[flag]; ok {
			ch <- prometheus.MustNewConstMetric(metric.Desc, metric.Type, value, core)
		}
	}
}

func (c *LitespeedCollector) collectReqRateMetrics(core string, reports []requestRateReport, ch chan<- prometheus.Metric) {
	for _, rrReport := range reports {
		for flag, value := range rrReport.KeyValues {
			if metric, ok := LitespeedMetrics[flag]; ok {
				ch <- prometheus.MustNewConstMetric(metric.Desc, metric.Type, value, core, rrReport.Hostname)
			}
		}
	}
}

func (c *LitespeedCollector) collectExtAppMetrics(core string, reports []externalAppReport, ch chan<- prometheus.Metric) {
	for _, eaReport := range reports {
		for flag, value := range eaReport.KeyValues {
			if metric, ok := LitespeedMetrics[flag]; ok {
				ch <- prometheus.MustNewConstMetric(metric.Desc, metric.Type, value, core, eaReport.Service, eaReport.Hostname, eaReport.Handler)
			}
		}
	}
}

func (c *LitespeedCollector) scrapeFile(fileName string) (report *litespeedReport, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer func() {
		file.Close()
		if r := recover(); r != nil {
			err = fmt.Errorf("failed scraping file: %s", r)
			report = nil
		}
	}()

	idRegex := regexp.MustCompile(`^\w*`)
	ibRegex := regexp.MustCompile(`\[([^\[\]]*)\]`)

	report = &litespeedReport{
		GeneralInfo: generalInfoReport{KeyValues: make(map[string]float64)},
		ReqRates:    []requestRateReport{},
		ExtApps:     []externalAppReport{},
	}
	reader := bufio.NewReader(file)
	var line string

	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimRight(line, "\n")

		identifier := idRegex.FindString(line)
		if identifier != "" {
			switch identifier {
			case versionField:
				_, v := parseKeyValPair(line, ": ")
				report.GeneralInfo.Version = v
			case uptimeField:
				_, v := parseKeyValPair(line, ": ")
				report.GeneralInfo.Uptime = v
			case bpsInField, maxconnField:
				m := parseKeyValLineToMap(line)
				for k, v := range m {
					if !c.metricIsTracked(k) {
						continue
					}

					vf, err := parseMetricValue(k, v)
					if err != nil {
						level.Error(c.logger).Log("msg", "Can't parse field value", "value", v, "err", err)
						c.scrapeFailures.Inc()
					} else {
						report.GeneralInfo.KeyValues[k] = vf
					}
				}
			case reqRateField:
				parts := strings.SplitN(line, ": ", 2)
				matches := ibRegex.FindStringSubmatch(line)

				if !c.options.ReqRatesByHost && matches[1] != "" {
					continue
				}

				m := parseKeyValLineToMap(parts[1])
				rr := requestRateReport{
					Hostname:  matches[1],
					KeyValues: make(map[string]float64),
				}
				for k, v := range m {
					prefixedFlag := reqRateField + "_" + k
					if !c.metricIsTracked(prefixedFlag) {
						continue
					}

					vf, err := parseMetricValue(prefixedFlag, v)
					if err != nil {
						level.Error(c.logger).Log("msg", "Can't parse field value", "value", v, "err", err)
						c.scrapeFailures.Inc()
					} else {
						rr.KeyValues[prefixedFlag] = vf
					}
				}
				report.ReqRates = append(report.ReqRates, rr)
			case extappField:
				if c.options.ExcludeExtapp {
					break
				}

				parts := strings.SplitN(line, ": ", 2)
				m := parseKeyValLineToMap(parts[1])
				matches := ibRegex.FindAllStringSubmatch(line, -1)
				er := externalAppReport{
					Service:   matches[0][1],
					Hostname:  matches[1][1],
					Handler:   matches[2][1],
					KeyValues: make(map[string]float64),
				}
				for k, v := range m {
					prefixedFlag := extappField + "_" + k
					if !c.metricIsTracked(prefixedFlag) {
						continue
					}

					vf, err := parseMetricValue(prefixedFlag, v)
					if err != nil {
						level.Error(c.logger).Log("msg", "Can't parse field value", "value", v, "err", err)
						c.scrapeFailures.Inc()
					} else {
						er.KeyValues[prefixedFlag] = vf
					}
				}
				report.ExtApps = append(report.ExtApps, er)
			}
		}
	}

	if err != io.EOF {
		return nil, err
	}

	return report, nil
}

func (c *LitespeedCollector) scrapeReports(filePattern string) (map[string]litespeedReport, error) {
	matches, err := filepath.Glob(filePattern)
	if err != nil {
		return nil, err
	}

	reports := make(map[string]litespeedReport)
	for _, match := range matches {
		report, err := c.scrapeFile(match)
		if err == nil {
			reports[match] = *report
		}
	}

	if !c.options.MetricsByCore {
		return map[string]litespeedReport{"": *sumReports(reports)}, nil
	}

	return reports, nil
}
