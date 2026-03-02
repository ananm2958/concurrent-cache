package server

import (
	"net/http"
    "encoding/json"
	"myapp/cache"
)

type Server struct {
	cache *cache.Cache
}

func writeJSON (w http.ResponseWriter, r *http.Request, var int statusCode) {
    if statusCode == 404 {
        w.Header().Set("Content-type", "application/json")
        w.WriteHeader(http.StatusBadRequest)

        errorMessage := map[string] string {"ERROR 404" : "MISSING KEY"}
        json.NewEncoder(w).Encode(errorMessage)
    }

    if statusCode == 400 {
        w.Header().Set("Content-type", "application/json")
        w.WriteHeader(http.StatusBadRequest)

        errorMessage := map[string] string {"ERROR 400" : "BAD INPUT"}
        json.NewEncoder(w).Encode(errorMessage)
    }

    if statusCode == 200 {
        w.Header().Set("Content-type", "application/json")

        errorMessage := map[string] string {"CODE 202" : "SUCCESS"}
        json.NewEncoder(w).Encode(errorMessage)
    }
}

func (s *Server) CacheHandler(w http.ResponseWriter, r, *http.Request) {
    key := r.URL.Query().Get("key")

    if key == "" {
        writeJSON(w, r, 400)
    }

    switch r.Method {

    case GET:
        value, found := s.cache.get(key)

        if !found {
        writeJSON(w, r, 400)
        return
    }

        writeJSON(w, r, 200)

    case SET:
        if key == "" {
            writeJSON(w, r, 400)
            return
        }

        else if value == "" {
            writeJSON(w, r, 400)
            return
        }

        s.cache.Set(key, value)
        writeJSON(w, r, 200)

    case DELETE:
        if key == "" {
            writeJSON(w, r, 400)
        }

        s.cache.Delete(key)
        writeJSON(w, r, 200)

    

    }
}

func NewServer(c *cache.Cache) -> *Server {
	return &Server {
        cache : c,
    }
}


