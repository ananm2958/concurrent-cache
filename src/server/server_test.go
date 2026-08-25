package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"concurrent-cache/src/cache"
	"concurrent-cache/src/metrics"
)

func TestCacheHTTPFlow(t *testing.T) {
	s := New(cache.NewCache(10, time.Minute), metrics.New(), nil)
	set := httptest.NewRequest(http.MethodPost, "/cache", strings.NewReader(`{"key":"name","value":"codex"}`))
	set.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	s.Handler().ServeHTTP(response, set)
	if response.Code != http.StatusOK {
		t.Fatalf("set status = %d", response.Code)
	}
	response = httptest.NewRecorder()
	s.Handler().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/cache?key=name", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("get status = %d", response.Code)
	}
}

func TestMetricsExposeRequestRateAndConnectionGauges(t *testing.T) {
	m := metrics.New()
	s := New(cache.NewCache(10, time.Minute), m, nil)
	m.RecordConnectionState(http.StateNew)
	cacheResponse := httptest.NewRecorder()
	s.Handler().ServeHTTP(cacheResponse, httptest.NewRequest(http.MethodGet, "/cache?key=missing", nil))

	response := httptest.NewRecorder()
	s.Handler().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	body := response.Body.String()
	for _, metric := range []string{
		"cache_requests_per_second 1",
		"cache_connections_total 1",
		"cache_concurrent_connections 1",
		"cache_request_duration_seconds_count 1",
	} {
		if !strings.Contains(body, metric) {
			t.Errorf("metrics response missing %q:\n%s", metric, body)
		}
	}
}
