package queue

import (
	"container/list"
	"sync"
)

type baseQueue struct {
	sync.RWMutex
	list *list.List
}

func newBase() *baseQueue {
	return &baseQueue{list: list.New()}
}

func (l *baseQueue) Push(v interface{}) {
	l.Lock()
	l.list.PushFront(v)
	l.Unlock()
}

func (l *baseQueue) PushBatch(vs []interface{}) {
	l.Lock()
	for _, item := range vs {
		l.list.PushFront(item)
	}
	l.Unlock()
}

func (l *baseQueue) Pop() interface{} {
	l.Lock()
	if elem := l.list.Back(); elem != nil {
		item := l.list.Remove(elem)
		l.Unlock()
		return item
	}
	l.Unlock()
	return nil
}

func (l *baseQueue) PopBatch(max int) []interface{} {
	l.Lock()
	count := l.list.Len()
	if count == 0 {
		l.Unlock()
		return []interface{}{}
	}
	if count > max {
		count = max
	}
	items := make([]interface{}, 0, count)
	for i := 0; i < count; i++ {
		item := l.list.Remove(l.list.Back())
		items = append(items, item)
	}
	l.Unlock()
	return items
}

func (l *baseQueue) Len() int {
	l.RLock()
	defer l.RUnlock()
	return l.list.Len()
}

// Queue is definition of simple queue
type Queue interface {
	Pop() interface{}
	PopBatch(int) []interface{}
	Push(interface{}) bool
	PushBatch([]interface{}) bool
	Len() int
}

var _ Queue = (*Limited)(nil)

// Limited is list with size limited
type Limited struct {
	maxSize int
	list    *baseQueue
}

// New creates limited list with size
func New(maxSize int) *Limited {
	return &Limited{list: newBase(), maxSize: maxSize}
}

// Pop pops element from back
func (ll *Limited) Pop() interface{} {
	return ll.list.Pop()
}

// PopBatch pops some elements from back
func (ll *Limited) PopBatch(max int) []interface{} {
	return ll.list.PopBatch(max)
}

// Push pushes element at front
func (ll *Limited) Push(v interface{}) bool {
	if ll.list.Len() >= ll.maxSize {
		return false
	}
	ll.list.Push(v)
	return true
}

// PushBatch pushes some elements at front
func (ll *Limited) PushBatch(vs []interface{}) bool {
	// --- May exceed maxSize here, for example:
	// ll.list.Len() is 3 and the ll.maxSize is 5, but the len(vs) is 3
	if ll.list.Len() >= ll.maxSize {
		return false
	}
	ll.list.PushBatch(vs)
	return true
}

// Len returns length of list
func (ll *Limited) Len() int {
	return ll.list.Len()
}
