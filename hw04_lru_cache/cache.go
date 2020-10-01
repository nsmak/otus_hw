package hw04_lru_cache //nolint:golint,stylecheck
import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*listItem
	mu       *sync.Mutex
}

// Set добавляет значение в кэш по ключу.
func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if !ok {
		newItem := c.queue.PushFront(cacheItem{
			key:   key,
			value: value,
		})
		c.items[key] = newItem

		if c.queue.Len() > c.capacity {
			lastItem := c.queue.Back()
			c.queue.Remove(lastItem)
			delete(c.items, lastItem.value.(cacheItem).key)
		}
		return false
	}

	item.value = cacheItem{
		key:   key,
		value: value,
	}
	c.queue.MoveToFront(item)
	return true
}

// Get возвращает значение из кэша по ключу.
func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if !ok {
		return nil, false
	}
	c.queue.MoveToFront(item)
	return item.value.(cacheItem).value, true
}

// Clear очищает кэш.
func (c *lruCache) Clear() {
	c.queue = &list{}
	c.items = make(map[Key]*listItem)
}

type cacheItem struct {
	key   Key
	value interface{}
}

// NewCache возвращает новый инстанс кэша.
func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		items:    make(map[Key]*listItem),
		queue:    &list{},
		mu:       &sync.Mutex{},
	}
}
