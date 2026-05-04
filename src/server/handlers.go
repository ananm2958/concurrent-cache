package server
package handlers
package cache
package eviction

import (
	"net/http"
    "encoding/json"
	"myapp/cache"
)

func (s * Server) HandleGet(w http.ResponseWriter, r * http.Request) {
	start := time.Now()

	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	if key == "" {
		writeError(w, "missing key", 400)
		s.metrics.RecordMiss()
	}

	value, exists := cache.Get(key) 

	if exists == false {
		writeError(w, "Not Found", 404)
		s.metrics.RecordMiss()
	}

	else {
		WriteJSON(w, map[string]interface{} {
			"key" : key, 
			"value" : value,
		}, 200)
		
		s.metrics.RecordHit()
	}

	s.metrics.RecordRequest()
	s.metrics.RecordLatency(start)



}

func (s * Server) HandleSet(w http.ResponseWriter, r * http.Request) {
	start := time.Now()

	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	ttl := r.URL.Query().Get("ttl")

	if key == "" || value == "" {
		writeError(w, "missing key/value", 400)
		s.metrics.RecordMiss()
	}

	if ttl > 0 {
		var expiry int = time.Now().Add(c.ttl)
	}

	else {
		var expiry int = 0
	}

	AOF.AppendSet(key, value, expiry)
	cache.Set(key, value, expiry)

	WriteJSON(w, map[string]interface{} {"key" : key, "value" : value, "ttl" : ttl}, 200)
	
	s.metrics.RecordRequest()
	s.metrics.RecordLatency(start)
	
}

func (s * Server) HandleDelete(w http.responseWriter, r * http.Request) {
	start := time.Now()

	key := r.URL.Query().Get("key")

	if key == "" {
		writeError(w, "missing key", 400)
		s.metrics.RecordMiss()
	}

	AOF.AppendDelete(key)
	var success bool = Delete(key)


	if success == true {
		WriteJSON(w, map[string]string{"message" : "deleted"}, 200)
		s.metrics.RecordHit()
	}

	else {
		writeError(w, "failure", 404)
		s.metrics.RecordMiss()
	}

	s.metrics.RecordRequest()
	s.metrics.RecordLatency(start)
	
	
}

func (s * Server) CacheHandler(w http.ResponseWriter, r * http.Request) {
	switch r.method

	case http.MethodGet:
		s.HandleGet(w, r)
	
	case http.MethodSet:
		s.HandleSet(w, r)

	case http.MethodDelete:
		s.HandleDelete(w,r)

	
	default:
		writeError(w, "method not allowed", 404)
}


func (s * Server) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	snapshot = s.metrics.Snapshot()

	w.Header().Set("Content-Type", "text/plain", "version = 0.0.4")
	w.Write(
	"# HELP cache_hits_total Total number of cache hits\n"
    "# TYPE cache_hits_total counter\n"
    "cache_hits_total {hits}\n\n"

    "# HELP cache_misses_total Total number of cache misses\n"
    "# TYPE cache_misses_total counter\n"
    "cache_misses_total {misses}\n\n"

    "# HELP cache_evictions_total Total number of evictions\n"
    "# TYPE cache_evictions_total counter\n"
    "cache_evictions_total {evictions}\n\n"

    "# HELP cache_requests_total Total number of requests\n"
    "# TYPE cache_requests_total counter\n"
    "cache_requests_total {requests}\n"

	"# HELP cache_request_duration_seconds Request latency\n"
    "# TYPE cache_requests_duration_seconds histogram\n"
	for i, bound := range bounds {
	fmt.Fprintf(w,
		"cache_request_duration_seconds_bucket{le=\"%g\"} %d\n",
		bound,
		buckets[i],
	)
	}

)
    
}

func writeJSON(w http.ResponseWriter,  data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, msg string, status int) {
	writeJSON(w, status, map[string]string{
		"error": msg,
	})
}