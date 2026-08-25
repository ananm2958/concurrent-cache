package metrics

import (
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	requestRateBucketCount = 100
	requestRateBucketWidth = 10 * time.Millisecond
)

type requestRateBucket struct {
	mutex sync.Mutex
	tick  int64
	count int64
}

type Metrics struct {
	hits              int64
	misses            int64
	evictions         int64
	requests          int64
	connectionsTotal  int64
	activeConnections int64
	requestRate       [requestRateBucketCount]requestRateBucket

	buckets      []int64
	bucketBounds []float64

	totalLatencyNanos int64
	latencyCount      int64
}

type Snapshot struct {
	Hits, Misses, Evictions, Requests   int64
	ConnectionsTotal, ActiveConnections int64
	RequestsPerSecond                   int64
	Buckets                             []int64
	Bounds                              []float64
	TotalLatency                        float64
	LatencyCount                        int64
}

func New() *Metrics {
	bounds := []float64{0.001, 0.01, 0.1, 1}
	return &Metrics{
		bucketBounds: bounds,
		buckets:      make([]int64, len(bounds)),
	}
}

func (m *Metrics) Snapshot() Snapshot {
	buckets := make([]int64, len(m.buckets))
	for i := range m.buckets {
		buckets[i] = atomic.LoadInt64(&m.buckets[i])
	}

	return Snapshot{
		Hits:              atomic.LoadInt64(&m.hits),
		Misses:            atomic.LoadInt64(&m.misses),
		Evictions:         atomic.LoadInt64(&m.evictions),
		Requests:          atomic.LoadInt64(&m.requests),
		ConnectionsTotal:  atomic.LoadInt64(&m.connectionsTotal),
		ActiveConnections: atomic.LoadInt64(&m.activeConnections),
		RequestsPerSecond: m.requestsInLastSecond(),
		Buckets:           buckets,
		Bounds:            append([]float64(nil), m.bucketBounds...),
		TotalLatency:      float64(atomic.LoadInt64(&m.totalLatencyNanos)) / float64(time.Second),
		LatencyCount:      atomic.LoadInt64(&m.latencyCount),
	}
}

func (m *Metrics) RecordRequest() {
	atomic.AddInt64(&m.requests, 1)

	tick := time.Now().UnixNano() / int64(requestRateBucketWidth)
	bucket := &m.requestRate[tick%requestRateBucketCount]

	bucket.mutex.Lock()
	if bucket.tick != tick {
		bucket.tick = tick
		bucket.count = 0
	}
	bucket.count++
	bucket.mutex.Unlock()
}

func (m *Metrics) RecordConnectionState(state http.ConnState) {
	switch state {
	case http.StateNew:
		atomic.AddInt64(&m.connectionsTotal, 1)
		atomic.AddInt64(&m.activeConnections, 1)

	case http.StateClosed, http.StateHijacked:
		for {
			active := atomic.LoadInt64(&m.activeConnections)
			if active == 0 ||
				atomic.CompareAndSwapInt64(&m.activeConnections, active, active-1) {
				break
			}
		}
	}
}

func (m *Metrics) RecordHit() {
	atomic.AddInt64(&m.hits, 1)
}

func (m *Metrics) RecordMiss() {
	atomic.AddInt64(&m.misses, 1)
}

func (m *Metrics) RecordEviction() {
	atomic.AddInt64(&m.evictions, 1)
}

func (m *Metrics) RecordLatency(start time.Time) {
	duration := time.Since(start)

	atomic.AddInt64(&m.totalLatencyNanos, duration.Nanoseconds())
	atomic.AddInt64(&m.latencyCount, 1)

	for i, bound := range m.bucketBounds {
		if duration.Seconds() <= bound {
			atomic.AddInt64(&m.buckets[i], 1)
			break
		}
	}
}

func (m *Metrics) requestsInLastSecond() int64 {
	nowTick := time.Now().UnixNano() / int64(requestRateBucketWidth)
	var requests int64

	for i := range m.requestRate {
		bucket := &m.requestRate[i]

		bucket.mutex.Lock()
		if nowTick-bucket.tick >= 0 &&
			nowTick-bucket.tick < requestRateBucketCount {
			requests += bucket.count
		}
		bucket.mutex.Unlock()
	}

	return requests
}
