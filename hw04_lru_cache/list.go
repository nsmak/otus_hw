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
	value interface{}
	next  *listItem
	prev  *listItem
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
	newItem := &listItem{value: v}

	defer func() {
		l.front = newItem
		l.len++
	}()

	if l.len == 0 {
		l.back = newItem
		return newItem
	}

	newItem.next = l.front
	l.front.prev = newItem
	return newItem
}

// PushBack добавляет значение в конец списка.
func (l *list) PushBack(v interface{}) *listItem {
	newItem := &listItem{value: v}

	defer func() {
		l.back = newItem
		l.len++
	}()

	if l.len == 0 {
		l.front = newItem
		return newItem
	}
	newItem.prev = l.back
	l.back.next = newItem
	return newItem
}

// Remove удаляет элемент из списка.
func (l *list) Remove(i *listItem) {
	if i.next == nil {
		l.back = i.prev
	} else {
		i.next.prev = i.prev
	}

	if i.prev == nil {
		l.front = i.next
	} else {
		i.prev.next = i.next
	}

	l.len--
}

// MoveToFront перемещает элемент в начало списка.
func (l *list) MoveToFront(i *listItem) {
	if i.prev == nil {
		return
	}

	l.Remove(i)
	_ = l.PushFront(i.value)
}

// NewList возвращает новый инстанс списка.
func NewList() List {
	return &list{}
}
