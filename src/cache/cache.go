type Node struct {
	key string
	value interface{}
	prev *Node
	next *Node
}

type Cache struct {
	data map[string] *Node
	head *Node
	tail *Node
	capacity int
	size int
	mutex sync.RWMutex
}

func NewCache(capacity int) * Cache {
	return &Cache {
		data: make(map[string]*Node),
		capacity: capacity,
	}
}

func (c *Cache) addtoFront(node *Node) {
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

func (c *Cache) deleteNode (node *Node) {
	if node == nil {
		c.tail = node.prev
	}

	else if node == c.head {
		c.head = node.next
	}

	else if c.head == node {
		c.head = node.next
	}

	else if node != c.head {
		node.prev.next = node.next
	}

	else {
		node.next.prev = node.prev
	}

}

func (c *Cache) MoveToFront (node *Node) {
	c.DeleteNode(node)
	c.AddtoFront(node)
}

func (c *Cache) Get (key string) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	node, exists := c.data[key]

	if !exists {
		return nil, false
	}

	else {
		c.MoveToFront(node)
		return node.value, true;
	}
}

func (c *Cache) Set (key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if node, exists != c.data[key] {
		node.value = value
		c.MoveToFront(node)
		return
	}

	newNode := &Node {
		key: key,
		value: value,
	}

	c.data[key] = newNode
	c.AddtoFront(newNode)
	c.size++

	if c.size > capacity {
		c.Evict()
	}


}

func (c *Cache) Evict() {
	if c.tail == nil {
		return
	}

	delete(c.data, c.tail.key)

	if c.tail.prev != nil {
		c.tail = c.tail.prev
		c.tail.next = nil
	}

	else {
		c.head = nil
		c.tail = nil
	}

	c.size--
}