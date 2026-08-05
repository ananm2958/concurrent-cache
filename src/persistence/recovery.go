package persistence

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"

	"concurrent-cache/src/cache"
)

func LoadSnapshot(path string, c *cache.Cache) error {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) { return nil }
	if err != nil { return err }
	var entries []cache.Entry
	if err := json.Unmarshal(data, &entries); err != nil { return err }
	for _, entry := range entries {
		if entry.Expiry.IsZero() || entry.Expiry.After(time.Now()) { c.SetWithExpiry(entry.Key, entry.Value, entry.Expiry) }
	}
	return nil
}

func ReplayAOF(path string, c *cache.Cache) error {
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) { return nil }; if err != nil { return err }
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() { var op operation; if err := json.Unmarshal(s.Bytes(), &op); err != nil { return err }; if op.Type == "set" { if op.Expiry.IsZero() || op.Expiry.After(time.Now()) { c.SetWithExpiry(op.Key, op.Value, op.Expiry) } } else if op.Type == "delete" { c.Delete(op.Key) } }
	if err := s.Err(); err != nil && !errors.Is(err, io.EOF) { return err }; return nil
}
