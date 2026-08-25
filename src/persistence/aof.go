package persistence

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type AOF struct { path string; mutex sync.Mutex }
type operation struct { Type string `json:"type"`; Key string `json:"key"`; Value interface{} `json:"value,omitempty"`; Expiry time.Time `json:"expiry,omitempty"` }
func NewAOF(path string) *AOF { return &AOF{path: path} }

func (a *AOF) AppendSet(key string, value interface{}, expiry time.Time) error {
	a.mutex.Lock(); defer a.mutex.Unlock()
	f, err := os.OpenFile(a.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return fmt.Errorf("open AOF: %w", err)
	}

	defer f.Close()
	return json.NewEncoder(f).Encode(operation{Type: "set", Key: key, Value: value, Expiry: expiry})
}


func (a *AOF) AppendDelete(key string) error {
	a.mutex.Lock(); defer a.mutex.Unlock()
	f, err := os.OpenFile(a.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return fmt.Errorf("open AOF: %w", err)
	}
	
	defer f.Close()
	return json.NewEncoder(f).Encode(operation{Type: "delete", Key: key})
}
