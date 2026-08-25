package persistence

import (
	"encoding/json"
	"fmt"
	"os"

	"concurrent-cache/src/cache"
)

// SaveSnapshot atomically replaces path with a JSON snapshot of the cache.
func SaveSnapshot(path string, c *cache.Cache) error {
	data, err := json.Marshal(c.Entries())
	if err != nil { return fmt.Errorf("marshal snapshot: %w", err) }
	return os.WriteFile(path, data, 0644)
}
