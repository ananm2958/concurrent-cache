package cache

import (
	"testing"
	"time"
)

func TestCacheEvictsLeastRecentlyUsed(t *testing.T) {
	c := NewCache(2, time.Minute)
	c.Set("a", 1); c.Set("b", 2)
	if _, ok := c.Get("a"); !ok { t.Fatal("expected a") }
	c.Set("c", 3)
	if _, ok := c.Get("b"); ok { t.Fatal("expected b to be evicted") }
}

func TestCacheExpires(t *testing.T) {
	c := NewCache(2, time.Minute)
	c.SetWithExpiry("a", 1, time.Now().Add(-time.Second))
	if _, ok := c.Get("a"); ok { t.Fatal("expected expired key to be absent") }
}
