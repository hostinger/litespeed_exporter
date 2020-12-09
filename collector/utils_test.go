package collector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSumOrAppendHandlesProperly(t *testing.T) {
	tests := []struct {
		kv map[string]float64
		k  string
		v  float64
		e  map[string]float64
	}{
		{map[string]float64{}, "test", 0.1, map[string]float64{"test": 0.1}},
		{map[string]float64{"test2": 1.2}, "test2", 0.8, map[string]float64{"test2": 2.0}},
	}

	for _, tc := range tests {
		sumOrAppend(tc.kv, tc.k, tc.v)
		assert.Equal(t, tc.e, tc.kv)
	}
}

func TestParseMetricValueParsesValidMetricProperly(t *testing.T) {
	tests := []struct {
		f string
		v string
		e float64
	}{
		{reqRateReqPerSecField, "0.1", 0.1},
		{extappReqPerSecField, "1.2", 1.2},
		{reqRatePubCacheHitsPerSecField, "2.34", 2.34},
		{reqRatePrivateCacheHitsPerSecField, "4.5", 4.5},
		{reqRateStaticHitsPerSecField, "5.6", 5.6},
		{bpsInField, "50", 50.0},
		{plainconnField, "123", 123.0},
	}

	for _, tc := range tests {
		value, err := parseMetricValue(tc.f, tc.v)
		assert.Equal(t, tc.e, value)
		assert.Nil(t, err, "Unexpected error while parsing metric value")
	}
}

func TestParseMetricValueHandlesInvalidValueType(t *testing.T) {
	value, err := parseMetricValue(versionField, "test")
	assert.Equal(t, 0.0, value)
	assert.Error(t, err)
}

func TestParseKeyValPairReturnsExpected(t *testing.T) {
	tests := []struct {
		s  string
		ek string
		ev string
	}{
		{"TEST: true", "TEST", "true"},
		{"TEST 2: true too I guess", "TEST 2", "true too I guess"},
		{"REQ_PER_SEC: 123.45", "REQ_PER_SEC", "123.45"},
	}

	for _, tc := range tests {
		k, v := parseKeyValPair(tc.s, ": ")
		assert.Equal(t, tc.ek, k)
		assert.Equal(t, tc.ev, v)
	}
}

func TestParseKeyValLineReturnsExpected(t *testing.T) {
	tests := []struct {
		s string
		e map[string]string
	}{
		{"", map[string]string{}},
		{"TEST: true", map[string]string{"TEST": "true"}},
		{"TEST: true, NOT_TEST: false, RANDOM: 123.45", map[string]string{"TEST": "true", "NOT_TEST": "false", "RANDOM": "123.45"}},
	}

	for _, tc := range tests {
		m := parseKeyValLineToMap(tc.s)
		assert.Equal(t, tc.e, m)
	}
}

func TestParseFlagsToMapReturnsExpected(t *testing.T) {
	tests := []struct {
		s []string
		e map[string]bool
	}{
		{[]string{}, map[string]bool{}},
		{[]string{"TEST_1", "TEST_2"}, map[string]bool{"TEST_1": true, "TEST_2": true}},
	}

	for _, tc := range tests {
		m := ParseFlagsToMap(tc.s)
		assert.Equal(t, tc.e, m)
	}
}
