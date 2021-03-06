package hw04_lru_cache //nolint:golint,stylecheck

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().value)
		require.Equal(t, 70, l.Back().value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.next {
			elems = append(elems, i.value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("moving test", func(t *testing.T) {
		l := NewList()

		secondItem := l.PushFront(10) // [10]
		firstItem := l.PushFront(20)  // [20, 10]
		thirdItem := l.PushBack(30)   // [20, 10, 30]

		l.MoveToFront(thirdItem)
		require.Equal(t, secondItem, l.Back())

		l.MoveToFront(firstItem)
		require.NotEqual(t, thirdItem, l.Front())
		require.Equal(t, secondItem, l.Back())
	})

	t.Run("actions with item not from the list", func(t *testing.T) {
		l := NewList()

		firstItem := l.PushFront(10)
		_ = l.PushBack(11)

		notInListItem := &listItem{value: 22}
		l.MoveToFront(notInListItem)
		require.Equal(t, firstItem, l.Front())

		l.Remove(notInListItem)
		require.Equal(t, 2, l.Len())
	})
}
