package server

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"concurrent-cache/src/cache"
	"concurrent-cache/src/metrics"
	"concurrent-cache/src/persistence"
)

type Server struct {
	cache   *cache.Cache
	metrics *metrics.Metrics
	aof     *persistence.AOF
	mux     *http.ServeMux
}

func New(c *cache.Cache, m *metrics.Metrics, aof *persistence.AOF) *Server {
	return &Server{cache: c, metrics: m, aof: aof}
}
func (s *Server) Handler() http.Handler {
	if s.mux == nil {
		s.mux = http.NewServeMux()
		s.mux.HandleFunc("/cache", s.cacheHandler)
		s.mux.HandleFunc("/metrics", s.metricsHandler)
	}
	return s.mux
}
func (s *Server) Start(port int) error {
	return (&http.Server{
		Addr: ":" + strconv.Itoa(port), Handler: s.Handler(),
		ReadHeaderTimeout: 5 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second,
		ConnState: func(_ net.Conn, state http.ConnState) { s.metrics.RecordConnectionState(state) },
	}).ListenAndServe()
}
