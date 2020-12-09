package collector

import (
	"sort"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "litespeed"
)

var (
	// LitespeedMetrics includes all available LiteSpeed metrics
	LitespeedMetrics = metrics{
		bpsInField:                         newGenericMetric(bpsInField, "BPS_IN metric.", prometheus.GaugeValue),
		bpsOutField:                        newGenericMetric(bpsOutField, "BPS_OUT metric.", prometheus.GaugeValue),
		sslBpsInField:                      newGenericMetric(sslBpsInField, "SSL_BPS_IN metric.", prometheus.GaugeValue),
		sslBpsOutField:                     newGenericMetric(sslBpsOutField, "SSL_BPS_OUT metric.", prometheus.GaugeValue),
		maxconnField:                       newGenericMetric(maxconnField, "MAXCONN metric.", prometheus.GaugeValue),
		maxsslConnField:                    newGenericMetric(maxsslConnField, "MAXSSL_CONN metric.", prometheus.GaugeValue),
		plainconnField:                     newGenericMetric(plainconnField, "PLAINCONN metric.", prometheus.GaugeValue),
		availconnField:                     newGenericMetric(availconnField, "AVAILCONN metric.", prometheus.GaugeValue),
		idleconnField:                      newGenericMetric(idleconnField, "IDLECONN metric.", prometheus.GaugeValue),
		sslconnField:                       newGenericMetric(sslconnField, "SSLCONN metric.", prometheus.GaugeValue),
		availsslField:                      newGenericMetric(availsslField, "AVAILSSL metric.", prometheus.GaugeValue),
		reqRateReqProcessingField:          newReqRateMetric(reqRateReqProcessingField, "REQ_RATE_REQ_PROCESSING metric.", prometheus.GaugeValue),
		reqRateReqPerSecField:              newReqRateMetric(reqRateReqPerSecField, "REQ_RATE_REQ_PER_SEC metric.", prometheus.GaugeValue),
		reqRateTotReqsField:                newReqRateMetric(reqRateTotReqsField, "REQ_RATE_TOT_REQS metric.", prometheus.GaugeValue),
		reqRatePubCacheHitsPerSecField:     newReqRateMetric(reqRatePubCacheHitsPerSecField, "REQ_RATE_PUB_CACHE_HITS_PER_SEC metric.", prometheus.GaugeValue),
		reqRateTotalPubCacheHitsField:      newReqRateMetric(reqRateTotalPubCacheHitsField, "REQ_RATE_TOTAL_PUB_CACHE_HITS metric.", prometheus.GaugeValue),
		reqRatePrivateCacheHitsPerSecField: newReqRateMetric(reqRatePrivateCacheHitsPerSecField, "REQ_RATE_PRIVATE_CACHE_HITS_PER_SEC metric.", prometheus.GaugeValue),
		reqRateTotalPrivateCacheHitsField:  newReqRateMetric(reqRateTotalPrivateCacheHitsField, "REQ_RATE_TOTAL_PRIVATE_CACHE_HITS metric.", prometheus.GaugeValue),
		reqRateStaticHitsPerSecField:       newReqRateMetric(reqRateStaticHitsPerSecField, "REQ_RATE_STATIC_HITS_PER_SEC metric.", prometheus.GaugeValue),
		reqRateTotalStaticHitsField:        newReqRateMetric(reqRateTotalStaticHitsField, "REQ_RATE_TOTAL_STATIC_HITS metric.", prometheus.GaugeValue),
		extappCmaxconnField:                newExtappMetric(extappCmaxconnField, "EXTAPP_CMAXCONN metric.", prometheus.GaugeValue),
		extappEmaxconnField:                newExtappMetric(extappEmaxconnField, "EXTAPP_EMAXCONN metric.", prometheus.GaugeValue),
		extappPoolSizeField:                newExtappMetric(extappPoolSizeField, "EXTAPP_POOL_SIZE metric.", prometheus.GaugeValue),
		extappInuseConnField:               newExtappMetric(extappInuseConnField, "EXTAPP_INUSE_CONN metric.", prometheus.GaugeValue),
		extappIdleConnField:                newExtappMetric(extappIdleConnField, "EXTAPP_IDLE_CONN metric.", prometheus.GaugeValue),
		extappWaitqueDepthField:            newExtappMetric(extappWaitqueDepthField, "EXTAPP_WAITQUE_DEPTH metric.", prometheus.GaugeValue),
		extappReqPerSecField:               newExtappMetric(extappReqPerSecField, "EXTAPP_REQ_PER_SEC metric.", prometheus.GaugeValue),
		extappTotReqsField:                 newExtappMetric(extappTotReqsField, "EXTAPP_TOT_REQS metric.", prometheus.GaugeValue),
	}
	litespeedVersion = prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "version"), "A metric with a constant '1' value labeled by the LiteSpeed version.", []string{"version"}, nil)
	litespeedUp      = prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "up"), "Was the last scrape of LiteSpeed successful.", nil, nil)
)

type metricInfo struct {
	Desc *prometheus.Desc
	Type prometheus.ValueType
}

type metrics map[string]metricInfo

func (m metrics) String() string {
	s := []string{}
	for k := range m {
		s = append(s, k)
	}
	sort.Strings(s)
	return strings.Join(s, ", ")
}

func newGenericMetric(metricName string, docString string, t prometheus.ValueType) metricInfo {
	return metricInfo{
		Desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", strings.ToLower(metricName)),
			docString,
			[]string{"core"},
			nil,
		),
		Type: t,
	}
}

func newReqRateMetric(metricName string, docString string, t prometheus.ValueType) metricInfo {
	return metricInfo{
		Desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", strings.ToLower(metricName)),
			docString,
			[]string{"core", "hostname"},
			nil,
		),
		Type: t,
	}
}

func newExtappMetric(metricName string, docString string, t prometheus.ValueType) metricInfo {
	return metricInfo{
		Desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", strings.ToLower(metricName)),
			docString,
			[]string{"core", "service", "hostname", "handler"},
			nil,
		),
		Type: t,
	}
}
