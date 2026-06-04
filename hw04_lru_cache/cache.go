package hw04lrucache

import (
	"fmt"
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
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
	fmt.Printf("[LRU Set] key=%q exist=%v value=%v\n", key, isExist, value)

	if !isExist {
		if c.queue.Len() == c.capacity {
			backItem := c.queue.Back()
			fmt.Printf("[LRU Set] evict value=%v\n", backItem.Value)
			c.queue.Remove(backItem)
		}
		newItem := c.queue.PushFront(value)
		c.items[key] = newItem
	} else {
		item.Value = value
		c.queue.MoveToFront(item)
		return true

	}
	return isExist

}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.rw.Lock()
	defer c.rw.Unlock()
	var item *ListItem = c.items[key]
	var isExist bool = (item != nil)
	if !isExist {
		return nil, isExist
	}
	c.queue.MoveToFront(item)
	isExist = true
	fmt.Printf("[LRU Get] key=%q found=%v value=%v\n", key, isExist, item.Value)
	return item.Value, isExist
}

func (c *lruCache) Clear() {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
