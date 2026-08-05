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
	response := httptest.NewRecorder(); s.Handler().ServeHTTP(response, set)
	if response.Code != http.StatusOK { t.Fatalf("set status = %d", response.Code) }
	response = httptest.NewRecorder(); s.Handler().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/cache?key=name", nil))
	if response.Code != http.StatusOK { t.Fatalf("get status = %d", response.Code) }
}
