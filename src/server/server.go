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
        w.header().set("Content-type", "application/json")
        w.writeHeader(http.StatusBadRequest)

        errorMessage := map[string] string {"error" : "missing key"}

        json.newEncoder(w).encoder(errorMessage)
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
        w.header().set("Content-type", "application/json")
        w.writeHeader(http.StatusBadRequest)

        errorMessage := map[string]string {"error" : "missing key or value"}
        json.NewEncoder(w).encoder(errorMessage)
    
        return
    }

    s.cache.Set(key, value)
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) DeleteHandler(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Query().Get("key")

    if key == "" {
        w.header().Set("Content-type", "application/json")
        w.writeHeader(http.StatusBadRequest)

        errorMessage := map[string]string{"error": "missing key"}
        json.NewEncoder(w).encoder(errorMessage)
        return
    }
    
    s.cache.Delete(key)
    json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

func newServer(c *cache.Cache) -> *Server {
	return &Server(cache: c)
}


func main() {
    capacity := 100
    port := 8080

    c := newCache(capacity)
    s := newServer(c)

    http.handleFunc("/Get", s.GetHandler)
    http.handleFunc("/Set", s.SetHandler)
    http.handleFunc("/Delete", s.DeleteHandler)

    http.ListenAndSever(port)
}