package collector

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexOfReqRateReturnsExpected(t *testing.T) {
	lr := litespeedReport{
		ReqRates: []requestRateReport{
			{Hostname: ""},
			{Hostname: "Test 1"},
			{Hostname: "Test 2"},
		},
	}

	tests := []struct {
		hostname string
		want     int
	}{
		{"", 0},
		{"Test 1", 1},
		{"Test 2", 2},
		{"Non existing", -1},
	}

	for _, tc := range tests {
		f := func(r requestRateReport) bool {
			return strings.HasPrefix(r.Hostname, tc.hostname)
		}
		assert.Equal(t, tc.want, lr.indexOfReqRate(f), "Unexpected filtered reports.")
	}
}

func TestIndexOfExtappReturnsExpected(t *testing.T) {
	lr := litespeedReport{
		ExtApps: []externalAppReport{
			{Service: "Proxy", Hostname: "localhost", Handler: "lsphp.10000"},
			{Service: "LSAPI", Hostname: "localhost", Handler: "lsphp.10001"},
			{Service: "LSAPI", Hostname: "test.com", Handler: "lsphp.10002"},
		},
	}

	tests := []struct {
		service  string
		hostname string
		want     int
	}{
		{"Proxy", "localhost", 0},
		{"LSAPI", "localhost", 1},
		{"LSAPI", "test.com", 2},
		{"Proxy", "test.com", -1},
	}

	for _, tc := range tests {
		f := func(r externalAppReport) bool {
			return r.Service == tc.service && r.Hostname == tc.hostname
		}
		assert.Equal(t, tc.want, lr.indexOfExtApp(f), "Unexpected filtered reports.")
	}
}

func TestAddsReportsProperly(t *testing.T) {
	a := litespeedReport{
		GeneralInfo: generalInfoReport{Version: "Test version", Uptime: "00:10:10", KeyValues: map[string]float64{bpsInField: 11}},
		ReqRates: []requestRateReport{
			{Hostname: "", KeyValues: map[string]float64{reqRateReqPerSecField: 12.34, reqRateTotReqsField: 5}},
		},
		ExtApps: []externalAppReport{
			{Service: "Proxy", Hostname: "", Handler: "lsphp.10000", KeyValues: map[string]float64{}},
		},
	}

	b := litespeedReport{
		GeneralInfo: generalInfoReport{Version: "Test version 2", Uptime: "00:20:20", KeyValues: map[string]float64{bpsInField: 123, bpsOutField: 321}},
		ReqRates: []requestRateReport{
			{Hostname: "", KeyValues: map[string]float64{reqRateReqPerSecField: 10, reqRateTotReqsField: 6}},
			{Hostname: "Test 1", KeyValues: map[string]float64{reqRateTotalPubCacheHitsField: 200}},
		},
		ExtApps: []externalAppReport{
			{Service: "Proxy", Hostname: "", Handler: "lsphp.10000", KeyValues: map[string]float64{extappReqPerSecField: 7.8, extappTotReqsField: 456}},
			{Service: "LSAPI", Hostname: "localhost", Handler: "lsphp.10001", KeyValues: map[string]float64{extappCmaxconnField: 1000, extappEmaxconnField: 1678}},
		},
	}

	a.Add(b)

	assert.Equal(t, b.GeneralInfo.Version, a.GeneralInfo.Version)
	assert.Equal(t, b.GeneralInfo.Uptime, a.GeneralInfo.Uptime)
	assert.Len(t, a.GeneralInfo.KeyValues, 2)
	assert.Equal(t, 134.0, a.GeneralInfo.KeyValues[bpsInField])
	assert.Equal(t, 321.0, a.GeneralInfo.KeyValues[bpsOutField])
	assert.Len(t, a.ReqRates, 2)
	assert.Equal(t, "", a.ReqRates[0].Hostname)
	assert.Equal(t, 22.34, a.ReqRates[0].KeyValues[reqRateReqPerSecField])
	assert.Equal(t, 11.0, a.ReqRates[0].KeyValues[reqRateTotReqsField])
	assert.Equal(t, b.ReqRates[1], a.ReqRates[1])
	assert.Len(t, a.ExtApps, 2)
	assert.Equal(t, "Proxy", a.ExtApps[0].Service)
	assert.Equal(t, "", a.ExtApps[0].Hostname)
	assert.Equal(t, "lsphp.10000", a.ExtApps[0].Handler)
	assert.Len(t, a.ExtApps[0].KeyValues, 2)
	assert.Equal(t, 7.8, a.ExtApps[0].KeyValues[extappReqPerSecField])
	assert.Equal(t, 456.0, a.ExtApps[0].KeyValues[extappTotReqsField])
	assert.Equal(t, b.ExtApps[1], a.ExtApps[1])
}

func TestSumsMultipleReportsProperly(t *testing.T) {
	a := litespeedReport{
		GeneralInfo: generalInfoReport{Version: "Test version", Uptime: "00:10:10", KeyValues: map[string]float64{bpsInField: 11}},
		ReqRates: []requestRateReport{
			{Hostname: "", KeyValues: map[string]float64{reqRateReqPerSecField: 12.34, reqRateTotReqsField: 5}},
		},
		ExtApps: []externalAppReport{
			{Service: "Proxy", Hostname: "", Handler: "lsphp.10000", KeyValues: map[string]float64{}},
		},
	}

	b := litespeedReport{
		GeneralInfo: generalInfoReport{Version: "Test version 2", Uptime: "00:20:20", KeyValues: map[string]float64{bpsInField: 123, bpsOutField: 321}},
		ReqRates: []requestRateReport{
			{Hostname: "", KeyValues: map[string]float64{reqRateReqPerSecField: 10, reqRateTotReqsField: 6}},
			{Hostname: "Test 1", KeyValues: map[string]float64{reqRateTotalPubCacheHitsField: 200}},
		},
		ExtApps: []externalAppReport{
			{Service: "Proxy", Hostname: "", Handler: "lsphp.10000", KeyValues: map[string]float64{extappReqPerSecField: 7.8, extappTotReqsField: 456}},
			{Service: "LSAPI", Hostname: "localhost", Handler: "lsphp.10001", KeyValues: map[string]float64{extappCmaxconnField: 1000, extappEmaxconnField: 1678}},
		},
	}

	c := litespeedReport{
		GeneralInfo: generalInfoReport{Version: "Test version 3", Uptime: "00:30:30", KeyValues: map[string]float64{bpsInField: 6, bpsOutField: 9}},
		ReqRates: []requestRateReport{
			{Hostname: "Test 1", KeyValues: map[string]float64{reqRateTotalPubCacheHitsField: 100}},
		},
		ExtApps: []externalAppReport{
			{Service: "Proxy", Hostname: "", Handler: "lsphp.10000", KeyValues: map[string]float64{extappReqPerSecField: 0.2, extappTotReqsField: 4}},
		},
	}

	r := sumReports(map[string]litespeedReport{"report.1": a, "report.2": b, "report.3": c})

	assert.Equal(t, c.GeneralInfo.Version, r.GeneralInfo.Version)
	assert.Equal(t, c.GeneralInfo.Uptime, r.GeneralInfo.Uptime)
	assert.Equal(t, 140.0, r.GeneralInfo.KeyValues[bpsInField])
	assert.Equal(t, 330.0, r.GeneralInfo.KeyValues[bpsOutField])
	assert.Equal(t, 300.0, r.ReqRates[1].KeyValues[reqRateTotalPubCacheHitsField])
	assert.Equal(t, 8.0, r.ExtApps[0].KeyValues[extappReqPerSecField])
	assert.Equal(t, 460.0, r.ExtApps[0].KeyValues[extappTotReqsField])
}
