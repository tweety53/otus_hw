package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int
	Front() *listItem
	Back() *listItem
	PushFront(v interface{}) *listItem
	PushBack(v interface{}) *listItem
	Remove(i *listItem)
	MoveToFront(i *listItem)
}

type listItem struct {
	Value interface{}
	Next  *listItem
	Prev  *listItem
}

type list struct {
	len   int
	first *listItem
	last  *listItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *listItem {
	return l.first
}

func (l *list) Back() *listItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *listItem {
	l.len++
	firstElem := l.Front()
	newElem := &listItem{
		Value: extractValue(v),
		Next:  firstElem,
		Prev:  nil,
	}

	if firstElem != nil {
		firstElem.Prev = newElem
	}
	l.first = newElem
	if l.Back() == nil {
		l.last = newElem
	}

	return newElem
}

func (l *list) PushBack(v interface{}) *listItem {
	l.len++
	lastElem := l.Back()
	newElem := &listItem{
		Value: extractValue(v),
		Next:  nil,
		Prev:  lastElem,
	}

	if lastElem != nil {
		lastElem.Next = newElem
	}
	l.last = newElem
	if l.Front() == nil {
		l.first = newElem
	}

	return newElem
}

func (l *list) Remove(i *listItem) {
	l.len--
	if i == l.first {
		l.first = i.Next
	} else if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	if i == l.last {
		l.last = i.Prev
	} else if i.Next != nil {
		i.Next.Prev = i.Prev
	}
}

func (l *list) MoveToFront(i *listItem) {
	if i == l.Front() {
		return
	}

	l.Remove(i)
	l.PushFront(i)
}

func NewList() List {
	return &list{}
}

func extractValue(v interface{}) interface{} {
	switch val := v.(type) {
	case *listItem:
		return val.Value
	default:
		return val
	}
}
