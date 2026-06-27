package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cachedItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	rw       sync.Mutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.rw.Lock()
	defer c.rw.Unlock()
	item := c.items[key]
	isExist := (item != nil)

	if !isExist {
		if c.queue.Len() == c.capacity {
			backItem := c.queue.Back()
			c.queue.Remove(backItem)
			oldItem := backItem.Value.(cachedItem)
			delete(c.items, oldItem.key)
		}

		newItem := c.queue.PushFront(cachedItem{key: key, value: value})
		c.items[key] = newItem
		return isExist
	}

	e := item.Value.(cachedItem)
	e.value = value
	item.Value = e
	c.queue.MoveToFront(item)
	return isExist
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.rw.Lock()
	defer c.rw.Unlock()
	item := c.items[key]
	isExist := (item != nil)
	if !isExist {
		return nil, isExist
	}

	c.queue.MoveToFront(item)

	val := item.Value.(cachedItem)
	return val.value, isExist
}

func (c *lruCache) Clear() {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
