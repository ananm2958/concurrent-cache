type Metrics struct {
    hits int64
    misses int64
    evictions int64
    requests int64

    buckets      []int64
	bucketBounds []float64

	totalLatency float64

    mutex sync.RLock
}

func Snapshot() {
	metrics.mutex.RLock()

	hits int = metrics.hits
	misses int = metrics.misses
	evictions int = metrics.evictions
	requests int = metrics.requests
	buckets := append([]int64(nil), m.buckets)
	bounds := append([]float64(nil), m.bucketBounds)

	metrics.mutex.RUnlock()
}

func RecordRequest() {
	metrics.Snapshot()
	
	metrics.mutex.Lock()
	defer metrics.mutex.Unlock()

	metrics.requests++


func RecordHit() {
	metrics.RecordRequest()
	metrics.mutex.Lock()
	defer metrics.mutex.Unlock()

	metrics.hits++

}

func RecordMiss() {
	metrics.RecordRequest()
	metrics.mutex.Lock()
	defer metrics.mutex.Unlock()

	metrics.misses++
}

func RecordEviction() {
	metrics.RecordRequest()
	metrics.mutex.Lock()
	defer metrics.mutex.Unlock()

	metrics.evictions++
}

}

func RecordLatency(start int) {
	metrics.Snapshot()
	metrics.RecordRequest()

	duration := time.Now() - start

	metrics.mutex.Lock()
	defer metrics.mutex.Unlock()

	metrics.totalLatency += duration

	for i := 0; i < bucketBounds; i++ {
		if duration <= bucketBounds[i] {
			buckets[i]++
			break	
		}
	}

}
