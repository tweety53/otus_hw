package hw04_lru_cache //nolint:golint,stylecheck
import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		cap:   capacity,
		queue: NewList(),
		cmu:   sync.Mutex{},
		cache: make(map[Key]*listItem),
	}
}

type lruCache struct {
	cap   int
	queue List
	cmu   sync.Mutex
	cache map[Key]*listItem
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.cmu.Lock()
	defer l.cmu.Unlock()

	if el, ok := l.cache[key]; ok {
		l.queue.MoveToFront(el)
		el.Value = cacheItem{
			key:   key,
			value: value,
		}
		return true
	}

	el := cacheItem{
		key:   key,
		value: value,
	}
	newEl := l.queue.PushFront(el)
	l.cache[key] = newEl

	if l.queue.Len() > l.cap {
		last := l.queue.Back()
		l.queue.Remove(last)

		cacheItem, ok := last.Value.(cacheItem)
		if !ok {
			return false
		}

		delete(l.cache, cacheItem.key)
	}

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.cmu.Lock()
	defer l.cmu.Unlock()

	if el, ok := l.cache[key]; ok {
		l.queue.MoveToFront(el)
		cacheItem, ok := el.Value.(cacheItem)
		if !ok {
			return nil, false
		}
		return cacheItem.value, true
	}

	return nil, false
}

func (l *lruCache) Clear() {
	l.cmu.Lock()
	defer l.cmu.Unlock()
	l.queue = NewList()
	l.cache = make(map[Key]*listItem)
}
