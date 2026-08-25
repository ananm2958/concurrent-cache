package cache

import (
	"sync"
	"time"
)

type Node struct {
	key    string
	value  interface{}
	prev   *Node
	next   *Node
	expiry time.Time
}

type Cache struct {
	data     map[string]*Node
	head     *Node
	tail     *Node
	capacity int
	size     int
	ttl      time.Duration
	mutex    sync.RWMutex
}

// Entry is a serializable view of a cache value.
type Entry struct {
	Key string
	Value interface{}
	Expiry time.Time
}

func NewCache(capacity int, ttl time.Duration) *Cache {
	return &Cache{
		data:     make(map[string]*Node),
		capacity: capacity,
		ttl:      ttl,
	}
}


func (c *Cache) addToFront(node *Node) {
	node.prev = nil
	node.next = c.head

	if c.head != nil {
		c.head.prev = node
	}

	c.head = node

	if c.tail == nil {
		c.tail = node
	}
}

func (c *Cache) deleteNode(node *Node) {
	if node == nil {
		return
	}

	if node == c.head {
		c.head = node.next
		if c.head != nil {
			c.head.prev = nil
		}
	} else if node == c.tail {
		c.tail = node.prev
		if c.tail != nil {
			c.tail.next = nil
		}
	} else {
		if node.prev != nil {
			node.prev.next = node.next
		}
		if node.next != nil {
			node.next.prev = node.prev
		}
	}
}

func (c *Cache) moveToFront(node *Node) {
	c.deleteNode(node)
	c.addToFront(node)
}



func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	node, exists := c.data[key]
	if !exists {
		return nil, false
	}

	// TTL check
	if c.isExpired(node) {
		c.removeNodeCompletely(node)
		return nil, false
	}

	c.moveToFront(node)
	return node.value, true
}

func (c *Cache) Set(key string, value interface{}) {
	c.SetWithExpiry(key, value, time.Now().Add(c.ttl))
}

// SetWithExpiry stores a value with an explicit expiry. A zero expiry means it
// does not expire.
func (c *Cache) SetWithExpiry(key string, value interface{}, expiry time.Time) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if node, exists := c.data[key]; exists {
		node.value = value
		node.expiry = expiry
		c.moveToFront(node)
		return
	}

	newNode := &Node{
		key:    key,
		value:  value,
		expiry: expiry,
	}

	c.data[key] = newNode
	c.addToFront(newNode)
	c.size++

	c.evictIfNeeded()
}

// Entries returns a consistent copy of all non-expired cache entries.
func (c *Cache) Entries() []Entry {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	entries := make([]Entry, 0, c.size)
	for _, node := range c.data {
		if c.isExpired(node) {
			c.removeNodeCompletely(node)
			continue
		}
		entries = append(entries, Entry{Key: node.key, Value: node.value, Expiry: node.expiry})
	}
	return entries
}

func (c *Cache) Delete(key string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	node, exists := c.data[key]
	if !exists {
		return false
	}

	c.removeNodeCompletely(node)
	return true
}


func (c *Cache) removeNodeCompletely(node *Node) {
	c.deleteNode(node)
	delete(c.data, node.key)
	c.size--
}
