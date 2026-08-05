package metrics

import (
	"sync"
	"time"
)

type Metrics struct {
    hits int64
    misses int64
    evictions int64
    requests int64

    buckets      []int64
	bucketBounds []float64

	totalLatency float64

	mutex sync.RWMutex
}

type Snapshot struct { Hits, Misses, Evictions, Requests int64; Buckets []int64; Bounds []float64; TotalLatency float64 }

func New() *Metrics {
	bounds := []float64{0.001, 0.01, 0.1, 1}
	return &Metrics{bucketBounds: bounds, buckets: make([]int64, len(bounds))}
}

func (m *Metrics) Snapshot() Snapshot {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	hits := m.hits
	misses := m.misses
	evictions := m.evictions
	requests := m.requests
	buckets := append([]int64(nil), m.buckets)
	bounds := append([]float64(nil), m.bucketBounds)
	return Snapshot{hits, misses, evictions, requests, buckets, bounds, m.totalLatency}
}

func (m *Metrics) RecordRequest() { m.mutex.Lock(); defer m.mutex.Unlock(); m.requests++ }

func (m *Metrics) RecordHit() { m.mutex.Lock(); defer m.mutex.Unlock(); m.hits++ }
func (m *Metrics) RecordMiss() { m.mutex.Lock(); defer m.mutex.Unlock(); m.misses++ }

func (m *Metrics) RecordEviction() { m.mutex.Lock(); defer m.mutex.Unlock(); m.evictions++ }

func (m *Metrics) RecordLatency(start time.Time) {
	duration := time.Since(start).Seconds()
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.totalLatency += duration
	for i, bound := range m.bucketBounds {
		if duration <= bound {
			m.buckets[i]++
			break	
		}
	}

}
