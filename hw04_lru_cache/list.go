package hw04lrucache

import "sync"

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	length    int
	frontItem *ListItem
	backItem  *ListItem
	lock      sync.Mutex
}

func (l *list) Len() int {
	return l.length
}
func (l *list) Front() *ListItem {
	return l.frontItem
}
func (l *list) Back() *ListItem {
	return l.backItem
}
func (l *list) PushFront(v interface{}) *ListItem {
	l.lock.Lock()
	newItem := ListItem{}
	newItem.Value = v
	if l.frontItem != nil {
		newItem.Next = l.frontItem
		l.frontItem.Prev = &newItem
		l.frontItem = &newItem
	} else {
		l.frontItem = &newItem
		l.backItem = &newItem
	}
	l.length += 1
	l.lock.Unlock()
	return &newItem
}
func (l *list) PushBack(v interface{}) *ListItem {
	l.lock.Lock()
	newItem := ListItem{}
	newItem.Value = v
	if l.backItem != nil {
		newItem.Prev = l.backItem
		l.backItem.Next = &newItem
		l.backItem = &newItem
	} else {
		l.backItem = &newItem
		l.frontItem = &newItem

	}
	l.length += 1
	l.lock.Unlock()
	return &newItem
}

// Метод вызывается только от существующих в списке элементов.
func (l *list) Remove(i *ListItem) {
	l.lock.Lock()
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.backItem = i.Prev
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.frontItem = i.Next
	}
	l.length -= 1
	l.lock.Unlock()
}

// Метод вызывается только от существующих в списке элементов.
func (l *list) MoveToFront(i *ListItem) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if i.Prev == nil {
		return
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	if l.frontItem != nil {
		l.frontItem.Prev = i
	}
	i.Prev = nil
	i.Next = l.frontItem
	l.frontItem = i
}

func NewList() List {
	return new(list)
}
