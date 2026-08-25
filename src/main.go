package main

import (
	"log"
	"time"

	"concurrent-cache/src/cache"
	"concurrent-cache/src/metrics"
	"concurrent-cache/src/persistence"
	"concurrent-cache/src/server"
)

func main() {
	c := cache.NewCache(10000, 60*time.Second)
	if err := persistence.LoadSnapshot("snapshot.json", c); err != nil { log.Printf("load snapshot: %v", err) }
	if err := persistence.ReplayAOF("appendonly.aof", c); err != nil { log.Printf("replay AOF: %v", err) }
	go func() { ticker:=time.NewTicker(10*time.Second); defer ticker.Stop(); for range ticker.C { c.RemoveExpired() } }()
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			if err := persistence.SaveSnapshot("snapshot.json", c); err != nil { log.Printf("save snapshot: %v", err) }
		}
	}()
	s := server.New(c, metrics.New(), persistence.NewAOF("appendonly.aof"))
	log.Fatal(s.Start(8080))
}
