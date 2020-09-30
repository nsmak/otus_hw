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
	front *listItem
	back  *listItem
	len   int
}

// Len - длина списка.
func (l *list) Len() int {
	return l.len
}

// Front - первый элемент списка.
func (l *list) Front() *listItem {
	return l.front
}

// Back последний элемент списка.
func (l *list) Back() *listItem {
	return l.back
}

// PushFront добавляет значение в начало списка.
func (l *list) PushFront(v interface{}) *listItem {
	newItem := &listItem{Value: v}

	defer func() {
		l.front = newItem
		l.len++
	}()

	if l.len == 0 {
		l.back = newItem
		return newItem
	}

	newItem.Next = l.front
	l.front.Prev = newItem
	return newItem
}

// PushBack добавляет значение в конец списка.
func (l *list) PushBack(v interface{}) *listItem {
	newItem := &listItem{Value: v}

	defer func() {
		l.back = newItem
		l.len++
	}()

	if l.len == 0 {
		l.front = newItem
		return newItem
	}
	newItem.Prev = l.back
	l.back.Next = newItem
	return newItem
}

// Remove удаляет элемент из списка.
func (l *list) Remove(i *listItem) {
	if i.Next == nil {
		l.back = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	if i.Prev == nil {
		l.front = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	l.len--
}

// MoveToFront перемещает элемент в начало списка.
func (l *list) MoveToFront(i *listItem) {
	if i.Prev == nil {
		return
	}

	l.Remove(i)
	_ = l.PushFront(i.Value)
}

// NewList возвращает новый инстанс списка.
func NewList() List {
	return &list{}
}
