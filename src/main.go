package main
package cache
package eviction 
package persistence
package handlers
package server

import (
    "fmt"       // Package for printing text to console or writing responses
    "net/http"  // Package for creating HTTP servers and handling requests
)


func startCleanup(c *cache.Cache, interval time.Duration) {
    go func() {
        ticker := time.NewTicker(interval)
        defer ticker.Stop()

        for range ticker.C {
            c.RemoveExpired()
        }
    }()
}

func main() {
    capacity int := 10000
    ttl time.Duration = 60 * time.Second
    cache := newCache(capacity, ttl)
   
    server := newServer(cache)

    port int := 8080

    startCleanup(cache, 10 *time.Second )
    Start(8080)



}