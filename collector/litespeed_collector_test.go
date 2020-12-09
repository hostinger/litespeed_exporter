package collector

import (
	"os"
	"path"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func assertMetricsEqual(t *testing.T, c prometheus.Collector, expected string) {
	exp, err := os.Open(path.Join("..", "testdata", expected))
	if err != nil {
		t.Fatalf("Error opening expected result file %q: %v", expected, err)
	}
	if err := testutil.CollectAndCompare(c, exp); err != nil {
		t.Fatal("Metrics not equal:", err)
	}
}

func TestScrapeReportsHandlesZeroMatchingFiles(t *testing.T) {
	c := NewLitespeedCollector(
		LitespeedCollectorOpts{
			FilePattern:     path.Join("..", "testdata", "non-existing-pattern"),
			ReqRatesByHost:  false,
			MetricsByCore:   true,
			ExcludeExtapp:   false,
			ExcludedMetrics: ParseFlagsToMap([]string{}),
		},
		log.NewNopLogger(),
	)
	r, err := c.scrapeReports(c.options.FilePattern)

	assert.Nil(t, err)
	assert.Equal(t, map[string]litespeedReport{}, r)
}

func TestScrapeReportsHandlesMalformedReportFile(t *testing.T) {
	c := NewLitespeedCollector(
		LitespeedCollectorOpts{
			FilePattern:     path.Join("..", "testdata", "malformed_report"),
			ReqRatesByHost:  false,
			MetricsByCore:   true,
			ExcludeExtapp:   false,
			ExcludedMetrics: ParseFlagsToMap([]string{}),
		},
		log.NewNopLogger(),
	)
	r, err := c.scrapeReports(c.options.FilePattern)

	assert.Len(t, r, 0)
	assert.Nil(t, err)
}

func TestScrapeReportsHandlesMatchingFiles(t *testing.T) {
	c := NewLitespeedCollector(
		LitespeedCollectorOpts{
			FilePattern:     path.Join("..", "testdata", ".rtreport*"),
			ReqRatesByHost:  true,
			MetricsByCore:   true,
			ExcludeExtapp:   false,
			ExcludedMetrics: ParseFlagsToMap([]string{reqRatePubCacheHitsPerSecField}),
		},
		log.NewNopLogger(),
	)
	r, err := c.scrapeReports(c.options.FilePattern)

	assert.Len(t, r, 3)
	assert.Nil(t, err)
}
func TestMetricIsTrackedReturnsExpected(t *testing.T) {
	c := NewLitespeedCollector(
		LitespeedCollectorOpts{
			FilePattern:     path.Join("..", "testdata", "none"),
			ReqRatesByHost:  false,
			MetricsByCore:   true,
			ExcludeExtapp:   false,
			ExcludedMetrics: ParseFlagsToMap([]string{extappReqPerSecField}),
		},
		log.NewNopLogger(),
	)

	tests := []struct {
		inputFlag string
		want      bool
	}{
		{reqRateReqPerSecField, true},
		{extappReqPerSecField, false},
	}

	for _, tc := range tests {
		r := c.metricIsTracked(tc.inputFlag)
		assert.Equal(t, tc.want, r)
	}
}

func TestCollectHandlesInvalidReportValueTypes(t *testing.T) {
	c := NewLitespeedCollector(
		LitespeedCollectorOpts{
			FilePattern:     path.Join("..", "testdata", "invalid_value_types_report"),
			ReqRatesByHost:  false,
			MetricsByCore:   true,
			ExcludeExtapp:   false,
			ExcludedMetrics: ParseFlagsToMap([]string{}),
		},
		log.NewNopLogger(),
	)

	assertMetricsEqual(t, c, "invalid_value_types.metrics")
}

func TestCollectSkipsExcludedMetricsWhenAllExcluded(t *testing.T) {
	var ef []string
	for flag := range LitespeedMetrics {
		ef = append(ef, flag)
	}

	c := NewLitespeedCollector(
		LitespeedCollectorOpts{
			FilePattern:     path.Join("..", "testdata", ".rtreport"),
			ReqRatesByHost:  false,
			MetricsByCore:   true,
			ExcludeExtapp:   false,
			ExcludedMetrics: ParseFlagsToMap(ef),
		},
		log.NewNopLogger(),
	)

	assertMetricsEqual(t, c, "all_filtered.metrics")
}

func TestCollectSkipsExcludedMetricsWhenSomeExcluded(t *testing.T) {
	var ef []string
	for flag := range LitespeedMetrics {
		if flag != bpsInField && flag != reqRateTotReqsField && flag != extappReqPerSecField {
			ef = append(ef, flag)
		}
	}

	c := NewLitespeedCollector(
		LitespeedCollectorOpts{
			FilePattern:     path.Join("..", "testdata", ".rtreport"),
			ReqRatesByHost:  false,
			MetricsByCore:   true,
			ExcludeExtapp:   false,
			ExcludedMetrics: ParseFlagsToMap(ef),
		},
		log.NewNopLogger(),
	)

	assertMetricsEqual(t, c, "some_filtered.metrics")
}

func TestCollectSkipsExcludedMetricsExportsReqRateByHost(t *testing.T) {
	var ef []string
	for flag := range LitespeedMetrics {
		if flag != bpsInField && flag != reqRateTotReqsField && flag != extappReqPerSecField {
			ef = append(ef, flag)
		}
	}

	c := NewLitespeedCollector(
		LitespeedCollectorOpts{
			FilePattern:     path.Join("..", "testdata", ".rtreport"),
			ReqRatesByHost:  true,
			MetricsByCore:   true,
			ExcludeExtapp:   false,
			ExcludedMetrics: ParseFlagsToMap(ef),
		},
		log.NewNopLogger(),
	)

	assertMetricsEqual(t, c, "some_filtered_by_host.metrics")
}

func TestCollectSkipsExcludedMetricsExportsSummedUpReports(t *testing.T) {
	var ef []string
	for flag := range LitespeedMetrics {
		if flag != bpsInField && flag != reqRateTotReqsField && flag != extappReqPerSecField {
			ef = append(ef, flag)
		}
	}

	c := NewLitespeedCollector(
		LitespeedCollectorOpts{
			FilePattern:     path.Join("..", "testdata", ".rtreport*"),
			ReqRatesByHost:  false,
			MetricsByCore:   false,
			ExcludeExtapp:   false,
			ExcludedMetrics: ParseFlagsToMap(ef),
		},
		log.NewNopLogger(),
	)

	assertMetricsEqual(t, c, "some_filtered_summed_up.metrics")
}

func TestCollectSkipsExtappMetrics(t *testing.T) {
	var ef []string
	for flag := range LitespeedMetrics {
		if flag != bpsInField && flag != reqRateTotReqsField && flag != extappReqPerSecField {
			ef = append(ef, flag)
		}
	}

	c := NewLitespeedCollector(
		LitespeedCollectorOpts{
			FilePattern:     path.Join("..", "testdata", ".rtreport*"),
			ReqRatesByHost:  false,
			MetricsByCore:   false,
			ExcludeExtapp:   true,
			ExcludedMetrics: ParseFlagsToMap(ef),
		},
		log.NewNopLogger(),
	)

	assertMetricsEqual(t, c, "extapp_excluded.metrics")
}
