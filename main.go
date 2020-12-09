package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-kit/kit/log/level"
	"github.com/hostinger/litespeed_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	// Version set during build
	Version string
	// Date set during build
	Date string
	// Revision set during build
	Revision string
)

func main() {
	var (
		exporter = "litespeed_exporter"

		metricsPath              = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
		listenAddress            = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9777").String()
		litespeedScrapePattern   = kingpin.Flag("litespeed.scrape-pattern", "Pattern of files to scrape LiteSpeed metrics from.").Default("/tmp/lshttpd/.rtreport*").String()
		litespeedExcludedMetrics = kingpin.Flag("litespeed.exclude-metrics", "Comma-separated list of metrics to exclude. Available options: ["+collector.LitespeedMetrics.String()+"]").Default("").String()
		litespeedReqRatesByHost  = kingpin.Flag("litespeed.req-rates-by-host", "Export Request Rates by host.").Bool()
		litespeedMetricsByCore   = kingpin.Flag("litespeed.metrics-by-core", "Export metrics by core filename.").Bool()
		litespeedExcludeExtapp   = kingpin.Flag("litespeed.exclude-extapp", "Exclude EXTAPP metrics altogether.").Bool()
	)

	promlogConfig := &promlog.Config{}
	logger := promlog.New(promlogConfig)

	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.HelpFlag.Short('h')
	kingpin.Version(fmt.Sprintf("%s v%s (%s %s)", exporter, Version, Date, Revision))
	kingpin.Parse()

	excludedMetricFlags := strings.Split(*litespeedExcludedMetrics, ",")

	collector := collector.NewLitespeedCollector(
		collector.LitespeedCollectorOpts{
			FilePattern:     *litespeedScrapePattern,
			ReqRatesByHost:  *litespeedReqRatesByHost,
			MetricsByCore:   *litespeedMetricsByCore,
			ExcludeExtapp:   *litespeedExcludeExtapp,
			ExcludedMetrics: collector.ParseFlagsToMap(excludedMetricFlags),
		},
		logger,
	)

	prometheus.MustRegister(collector)

	level.Info(logger).Log("build", version.Info())
	level.Info(logger).Log("address", *listenAddress)

	http.Handle(*metricsPath, promhttp.Handler())

	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		level.Error(logger).Log("msg", "Could not start HTTP server", "err", err)
		os.Exit(1)
	}
}
