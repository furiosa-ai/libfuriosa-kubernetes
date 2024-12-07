package util

import (
	"github.com/google/btree"
)

type BtreeMapItem[K any, V any] struct {
	Key   K
	Value V
}

type BtreeMap[K any, V any] struct {
	tree *btree.BTreeG[BtreeMapItem[K, V]]
}

func NewBtreeMapWithLessFunc[K any, V any](degree int, less btree.LessFunc[BtreeMapItem[K, V]]) *BtreeMap[K, V] {
	if degree < 2 {
		degree = 2
	}

	return &BtreeMap[K, V]{
		tree: btree.NewG(degree, less),
	}
}

func (m *BtreeMap[K, V]) Get(key K) V {
	item, exists := m.tree.Get(BtreeMapItem[K, V]{Key: key})
	if !exists {
		return *new(V)
	}

	return item.Value
}

func (m *BtreeMap[K, V]) Has(key K) bool {
	return m.tree.Has(BtreeMapItem[K, V]{Key: key})
}

func (m *BtreeMap[K, V]) Keys() []K {
	result := make([]K, 0, m.tree.Len())
	m.tree.Ascend(func(item BtreeMapItem[K, V]) bool {
		result = append(result, item.Key)

		return true
	})

	return result
}

func (m *BtreeMap[K, V]) Insert(key K, value V) {
	m.tree.ReplaceOrInsert(BtreeMapItem[K, V]{Key: key, Value: value})
}

type BtreeSet[T any] struct {
	tree *btree.BTreeG[T]
}

func NewBtreeSetWithLessFunc[T any](degree int, less btree.LessFunc[T]) *BtreeSet[T] {
	if degree < 2 {
		degree = 2
	}

	return &BtreeSet[T]{
		tree: btree.NewG[T](degree, less),
	}
}

func (s *BtreeSet[T]) Len() int {
	return s.tree.Len()
}

func (s *BtreeSet[T]) Has(item T) bool {
	return s.tree.Has(item)
}

func (s *BtreeSet[T]) Items() []T {
	result := make([]T, 0, s.tree.Len())
	s.tree.Ascend(func(item T) bool {
		result = append(result, item)

		return true
	})

	return result
}

func (s *BtreeSet[T]) Insert(item T) {
	s.tree.ReplaceOrInsert(item)
}
