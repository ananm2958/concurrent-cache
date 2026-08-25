package metrics

import (
	"net/http"
	"testing"
	"time"
)

func TestSnapshotReportsRequestsPerSecond(t *testing.T) {
	m := New()
	m.RecordRequest()
	m.RecordRequest()
	if got := m.Snapshot().RequestsPerSecond; got != 2 {
		t.Fatalf("requests per second = %d, want 2", got)
	}

	expired := New()
	tick := time.Now().UnixNano()/int64(requestRateBucketWidth) - requestRateBucketCount
	bucket := &expired.requestRate[tick%requestRateBucketCount]

	bucket.mutex.Lock()
	bucket.tick = tick
	bucket.count = 1
	bucket.mutex.Unlock()

	if got := expired.Snapshot().RequestsPerSecond; got != 0 {
		t.Fatalf("expired requests per second = %d, want 0", got)
	}
}

func TestSnapshotReportsActiveConnections(t *testing.T) {
	m := New()
	m.RecordConnectionState(http.StateNew)
	m.RecordConnectionState(http.StateNew)
	m.RecordConnectionState(http.StateIdle)
	m.RecordConnectionState(http.StateClosed)

	snap := m.Snapshot()
	if snap.ConnectionsTotal != 2 || snap.ActiveConnections != 1 {
		t.Fatalf("connections = total %d active %d, want total 2 active 1", snap.ConnectionsTotal, snap.ActiveConnections)
	}
}
