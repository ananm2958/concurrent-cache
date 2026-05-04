package server
package cache
package eviction

import (
	"net/http"
    "encoding/json"
	"myapp/cache"
)

type Server struct {
	cache *cache.Cache
    metrics *metrics.Metrics
}



func newServer(c *Cache.cache) * Server {
    return &server {
        cache : c,
    }
}

func RegisterRoutes() {
    mux := http.newServeMux()
    mux.http.HandleCacheFunc("/cache", CacheHandler)
    mux.http.HandleMetricsFunc("/metrics", MetricsHandler)

   s.router = mux
}

func Start(port int) {
    RegisterRoutes()

    server := &http.Server{
        Addr: ":" + strconv.Itoa(port)
        Handler: mux,
        ReadTimeout: 5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout: 60 * time.Second, 
    }

    server.ListenAndServe()
}