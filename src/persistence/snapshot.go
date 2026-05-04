package snapshot

import (
	"fmt"
	"time"
	"os"
)

type Entry struct {
	Key    string        `json:"key"`
	Value  interface{}   `json:"value"`
	Expiry time.Time     `json:"expiry"`
}

func (s * Server) TakeSnap() {
	f, err := os.OpenFile("snapshot.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		panic(err)
	}

	entries := make([]Entry, 0, len(s.cache.data))

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for key, node := range cache.data {
		entries = append(entries, Entry{
			Key:    key,
			Value:  node.value,
			Expiry: node.expiry,
		})
	}

	data, err1 := json.Marshal(entries)

	if err != nil {
		panic(err)
	}

}