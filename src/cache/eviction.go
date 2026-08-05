package cache
import "time"


func (c *Cache) evictIfNeeded() {
	for c.size > c.capacity {
		c.evict()
	}

}

func (c *Cache) evict() {
	if c.tail == nil {
		return
	}

	node := c.tail
	c.removeNodeCompletely(node)
}



func (c *Cache) isExpired(node *Node) bool {
	if node.expiry.IsZero() {
		return false
	}
	return time.Now().After(node.expiry)
}


func (c *Cache) RemoveExpired() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, node := range c.data {
		if c.isExpired(node) {
			c.removeNodeCompletely(node)
			
		}
	}
}
