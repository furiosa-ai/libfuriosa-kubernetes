package util

import (
	"cmp"
	"github.com/google/btree"
)

type BTreeMapItem[K any, V any] struct {
	Key   K
	Value V
}

type BtreeMap[K any, V any] struct {
	tree *btree.BTreeG[BTreeMapItem[K, V]]
}

func NewBtreeMap[K cmp.Ordered, V any](degree int) *BtreeMap[K, V] {
	if degree < 2 {
		degree = 2
	}

	return &BtreeMap[K, V]{
		tree: btree.NewG(degree, func(a, b BTreeMapItem[K, V]) bool {
			return a.Key < b.Key
		}),
	}
}

func NewBtreeMapWithLessFunc[K any, V any](degree int, less btree.LessFunc[BTreeMapItem[K, V]]) *BtreeMap[K, V] {
	if degree < 2 {
		degree = 2
	}

	return &BtreeMap[K, V]{
		tree: btree.NewG(degree, less),
	}
}

func (m *BtreeMap[K, V]) Get(key K) V {
	item, exists := m.tree.Get(BTreeMapItem[K, V]{Key: key})
	if !exists {
		return *new(V)
	}

	return item.Value
}

func (m *BtreeMap[K, V]) Keys() []K {
	result := make([]K, 0, m.tree.Len())
	m.tree.Ascend(func(item BTreeMapItem[K, V]) bool {
		result = append(result, item.Key)

		return true
	})

	return result
}

func (m *BtreeMap[K, V]) Insert(key K, value V) {
	m.tree.ReplaceOrInsert(BTreeMapItem[K, V]{Key: key, Value: value})
}

type BTreeSet[T any] BtreeMap[T, struct{}]

func NewBtreeSet[T cmp.Ordered](degree int) *BTreeSet[T] {
	btreeMap := NewBtreeMap[T, struct{}](degree)
	return (*BTreeSet[T])(btreeMap)
}

func NewBtreeSetWithLessFunc[T any](degree int, less btree.LessFunc[BTreeMapItem[T, struct{}]]) *BTreeSet[T] {
	btreeMap := NewBtreeMapWithLessFunc[T, struct{}](degree, less)
	return (*BTreeSet[T])(btreeMap)
}

func (s *BTreeSet[T]) Len() int {
	return s.tree.Len()
}

func (s *BTreeSet[T]) Has(key T) bool {
	return s.tree.Has(BTreeMapItem[T, struct{}]{Key: key})
}

func (s *BTreeSet[T]) Keys() []T {
	result := make([]T, 0, s.tree.Len())
	s.tree.Ascend(func(item BTreeMapItem[T, struct{}]) bool {
		result = append(result, item.Key)

		return true
	})

	return result
}

func (s *BTreeSet[T]) Insert(key T) {
	s.tree.ReplaceOrInsert(BTreeMapItem[T, struct{}]{Key: key})
}
