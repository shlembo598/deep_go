package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Entry struct {
	key         int
	value       int
	left, right *Entry
}

type OrderedMap struct {
	Root *Entry
	size int
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{}
}

func (m *OrderedMap) Insert(key, value int) {
	current := &m.Root

	for {
		if *current == nil {
			*current = &Entry{key: key, value: value}
			m.size++
			return
		}

		if (**current).key == key {
			return
		}

		if key > (**current).key {
			current = &(*current).right
		} else {
			current = &(*current).left
		}
	}
}

func (m *OrderedMap) Erase(key int) {
	m.Root = m.eraseHandler(m.Root, key)
}

func (m *OrderedMap) eraseHandler(entry *Entry, key int) *Entry {
	if entry == nil {
		return nil
	}

	if key < entry.key {
		entry.left = m.eraseHandler(entry.left, key)
	} else if key > entry.key {
		entry.right = m.eraseHandler(entry.right, key)
	} else {
		m.size--

		if entry.left == nil && entry.right == nil {
			return nil
		}

		if entry.left == nil {
			return entry.right
		}
		if entry.right == nil {
			return entry.left
		}

		predecessor := entry.left
		for predecessor.right != nil {
			predecessor = predecessor.right
		}
		entry.key = predecessor.key
		entry.value = predecessor.value
		entry.left = m.eraseHandler(entry.left, predecessor.key)

		return entry
	}

	return entry
}

func (m *OrderedMap) Contains(key int) bool {
	current := m.Root

	for {
		if current == nil {
			return false
		}

		if current.key == key {
			return true
		}

		if key > current.key {
			current = current.right
		} else {
			current = current.left
		}
	}
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	m.forEachHandler(m.Root, action)
}

func (m *OrderedMap) forEachHandler(node *Entry, action func(int, int)) {
	if node == nil {
		return
	}

	m.forEachHandler(node.left, action)
	action(node.key, node.value)
	m.forEachHandler(node.right, action)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
