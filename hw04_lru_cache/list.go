package hw04lrucache

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
	newItem := ListItem{Value: v}
	newItem.Next = l.frontItem
	if l.frontItem != nil {
		l.frontItem.Prev = &newItem
	} else {
		l.backItem = &newItem
	}
	l.frontItem = &newItem
	l.length++
	return &newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := ListItem{Value: v}
	newItem.Prev = l.backItem
	if l.backItem != nil {
		l.backItem.Next = &newItem
	} else {
		l.frontItem = &newItem
	}
	l.backItem = &newItem
	l.length++
	return &newItem
}

// Метод вызывается только от существующих в списке элементов.
func (l *list) Remove(i *ListItem) {
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
	i.Prev = nil
	i.Next = nil
	l.length--
}

// Метод вызывается только от существующих в списке элементов.
func (l *list) MoveToFront(i *ListItem) {
	if i == nil || i == l.frontItem {
		return
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if i == l.backItem {
		l.backItem = i.Prev
	}

	// attach at front
	i.Prev = nil
	i.Next = l.frontItem
	if l.frontItem != nil {
		l.frontItem.Prev = i
	}
	l.frontItem = i
	if l.backItem == nil {
		l.backItem = i
	}
}

func NewList() List {
	return new(list)
}
