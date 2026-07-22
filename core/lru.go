package core

// Node represents a single element in our doubly-linked list
type Node struct {
	key   string
	value string
	prev  *Node
	next  *Node
}

// LRUCache combines a hash map for O(1) lookups and a doubly-linked list for O(1) evictions
type LRUCache struct {
	capacity int
	cache    map[string]*Node
	head     *Node
	tail     *Node
}

// NewLRU initializes a new LRU cache with the given capacity
func NewLRU(capacity int) *LRUCache {
	l := &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*Node),
		head:     &Node{},
		tail:     &Node{},
	}
	// Wire the dummy head and tail together
	l.head.next = l.tail
	l.tail.prev = l.head
	return l
}

// removeNode detaches a node from the linked list
func (l *LRUCache) removeNode(node *Node) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

// moveToHead moves an accessed or newly created node right behind the head
func (l *LRUCache) moveToHead(node *Node) {
	l.removeNode(node)
	node.prev = l.head
	node.next = l.head.next
	l.head.next.prev = node
	l.head.next = node
}

// Get retrieves an item and marks it as recently used
func (l *LRUCache) Get(key string) (string, bool) {
	if node, found := l.cache[key]; found {
		l.moveToHead(node) // Mark as most recently used
		return node.value, true
	}
	return "", false
}

// Put inserts a new item or updates an existing one, triggering eviction if at capacity
func (l *LRUCache) Put(key string, value string) {
	if node, found := l.cache[key]; found {
		node.value = value
		l.moveToHead(node)
		return
	}

	// If at capacity, evict the LEAST recently used item (the one right before the tail)
	if len(l.cache) >= l.capacity {
		tailNode := l.tail.prev
		l.removeNode(tailNode)
		delete(l.cache, tailNode.key)
	}

	// Create new node and insert at head
	newNode := &Node{key: key, value: value}
	l.cache[key] = newNode
	
	newNode.prev = l.head
	newNode.next = l.head.next
	l.head.next.prev = newNode
	l.head.next = newNode
}
