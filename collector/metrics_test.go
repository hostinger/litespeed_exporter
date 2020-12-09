package collector

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestStringReturnsAvailableMetricsString(t *testing.T) {
	s := LitespeedMetrics.String()
	e := "AVAILCONN, AVAILSSL, BPS_IN, BPS_OUT, EXTAPP_CMAXCONN, EXTAPP_EMAXCONN, EXTAPP_IDLE_CONN, EXTAPP_INUSE_CONN, EXTAPP_POOL_SIZE, EXTAPP_REQ_PER_SEC, EXTAPP_TOT_REQS, EXTAPP_WAITQUE_DEPTH, IDLECONN, MAXCONN, MAXSSL_CONN, PLAINCONN, REQ_RATE_PRIVATE_CACHE_HITS_PER_SEC, REQ_RATE_PUB_CACHE_HITS_PER_SEC, REQ_RATE_REQ_PER_SEC, REQ_RATE_REQ_PROCESSING, REQ_RATE_STATIC_HITS_PER_SEC, REQ_RATE_TOTAL_PRIVATE_CACHE_HITS, REQ_RATE_TOTAL_PUB_CACHE_HITS, REQ_RATE_TOTAL_STATIC_HITS, REQ_RATE_TOT_REQS, SSLCONN, SSL_BPS_IN, SSL_BPS_OUT"

	assert.Equal(t, e, s)
}

func TestCreatesNewGenericMetric(t *testing.T) {
	m := newGenericMetric(bpsInField, "test metric", prometheus.GaugeValue)

	assert.Equal(t, prometheus.GaugeValue, m.Type)
	assert.Equal(t, "Desc{fqName: \"litespeed_bps_in\", help: \"test metric\", constLabels: {}, variableLabels: [core]}", m.Desc.String())
}

func TestCreatesNewReqReportMetric(t *testing.T) {
	m := newReqRateMetric(reqRateTotReqsField, "test metric", prometheus.GaugeValue)

	assert.Equal(t, prometheus.GaugeValue, m.Type)
	assert.Equal(t, "Desc{fqName: \"litespeed_req_rate_tot_reqs\", help: \"test metric\", constLabels: {}, variableLabels: [core hostname]}", m.Desc.String())
}

func TestCreatesNewExtappMetric(t *testing.T) {
	m := newExtappMetric(extappCmaxconnField, "test metric", prometheus.GaugeValue)

	assert.Equal(t, prometheus.GaugeValue, m.Type)
	assert.Equal(t, "Desc{fqName: \"litespeed_extapp_cmaxconn\", help: \"test metric\", constLabels: {}, variableLabels: [core service hostname handler]}", m.Desc.String())
}
