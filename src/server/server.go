package server

import (
	"net/http"
    "encoding/json"
	"myapp/cache"
)

type Server struct {
	cache *cache.Cache
}

func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Query().Get("key")

    if key == "" {
        http.Error(w, "missing key", http.StatusBadRequest)
        return
    }

    value, ok := s.cache.Get(key)

    if !ok {
        http.Error(w, "key not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"key": key, "value": value})
}

func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Query().Get("key")
    value := r.URL.Query().Get("value")

    if key == "" || value == "" {
        http.Error(w, "missing key or value", http.StatusBadRequest)
        return
    }

    s.cache.Set(key, value)
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) DeleteHandler(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Query().Get("key")

    if key == "" {
        http.Error(w, "missing key", http.StatusBadRequest)
        return
    }
    
    s.cache.Delete(key)
    json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

func newServer(c *cache.Cache) *Server {
	return &Server(cache: c)
}


func main() {
    c := newCache(capacity)
    newServer(c)

    http.HandleFunc("/get", server.GetHandler)
    http.HandleFunc("/set", server.SetHandler)
    http.HandleFunc("/delete", server.DeleteHandler)

    http.ListenAndServe(":8080", nil)
}