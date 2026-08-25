package server

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"concurrent-cache/src/cache"
	"concurrent-cache/src/metrics"
)

// TestMaxRequestsPerSecond sends requests for a fixed interval and reports the
// observed throughput. It has no fixed RPS target because that varies by host.
func TestMaxRequestsPerSecond(t *testing.T) {
	workers := loadTestInt("CACHE_LOAD_WORKERS", 32)
	duration := time.Duration(loadTestInt("CACHE_LOAD_DURATION_MS", 750)) * time.Millisecond
	if workers < 1 || duration <= 0 {
		t.Fatal("CACHE_LOAD_WORKERS and CACHE_LOAD_DURATION_MS must be positive")
	}

	ts, m := newLoadTestServer(t)
	client := &http.Client{Transport: &http.Transport{MaxIdleConns: workers, MaxIdleConnsPerHost: workers}}
	t.Cleanup(func() { client.CloseIdleConnections() })

	deadline := time.Now().Add(duration)
	var successful int64
	var wg sync.WaitGroup
	started := time.Now()
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for time.Now().Before(deadline) {
				response, err := client.Get(ts.URL + "/cache?key=load-test")
				if err != nil {
					continue
				}
				_, _ = io.Copy(io.Discard, response.Body)
				response.Body.Close()
				if response.StatusCode == http.StatusOK {
					atomic.AddInt64(&successful, 1)
				}
			}
		}()
	}
	wg.Wait()
	elapsed := time.Since(started)
	requests := atomic.LoadInt64(&successful)
	if requests == 0 {
		t.Fatal("load test completed without a successful request")
	}
	snapshot := m.Snapshot()
	if snapshot.Requests != requests {
		t.Fatalf("recorded requests = %d, want %d", snapshot.Requests, requests)
	}
	t.Logf("max request-rate sample: %.0f successful requests/second (%d requests in %s; rolling metric=%d)", float64(requests)/elapsed.Seconds(), requests, elapsed.Round(time.Millisecond), snapshot.RequestsPerSecond)
}

// TestMaxConcurrentConnections holds TCP connections open and verifies that
// the server's active-connection gauge reaches the requested load.
func TestMaxConcurrentConnections(t *testing.T) {
	connections := loadTestInt("CACHE_LOAD_CONNECTIONS", 64)
	if connections < 1 {
		t.Fatal("CACHE_LOAD_CONNECTIONS must be positive")
	}

	ts, m := newLoadTestServer(t)
	open := make([]net.Conn, 0, connections)
	for i := 0; i < connections; i++ {
		connection, err := net.DialTimeout("tcp", ts.Listener.Addr().String(), time.Second)
		if err != nil {
			for _, openConnection := range open {
				_ = openConnection.Close()
			}
			t.Fatalf("open connection %d/%d: %v", i+1, connections, err)
		}
		open = append(open, connection)
	}
	t.Cleanup(func() {
		for _, connection := range open {
			_ = connection.Close()
		}
	})

	var snapshot metrics.Snapshot
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		snapshot = m.Snapshot()
		if snapshot.ActiveConnections >= int64(connections) {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if snapshot.ActiveConnections < int64(connections) {
		t.Fatalf("active connections = %d, want at least %d", snapshot.ActiveConnections, connections)
	}
	t.Logf("max concurrent-connection sample: %d active connections (%d total accepted)", snapshot.ActiveConnections, snapshot.ConnectionsTotal)
}

func newLoadTestServer(t *testing.T) (*httptest.Server, *metrics.Metrics) {
	t.Helper()
	m := metrics.New()
	c := cache.NewCache(10, time.Minute)
	c.Set("load-test", "ready")
	s := New(c, m, nil)
	ts := httptest.NewUnstartedServer(s.Handler())
	ts.Config.ConnState = func(_ net.Conn, state http.ConnState) { m.RecordConnectionState(state) }
	ts.Start()
	t.Cleanup(ts.Close)
	return ts, m
}

func loadTestInt(name string, defaultValue int) int {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("%s must be an integer: %v", name, err))
	}
	return parsed
}
