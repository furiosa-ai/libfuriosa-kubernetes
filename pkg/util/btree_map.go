package util

import "github.com/google/btree"

type BTreeMapItem[K btree.Ordered, V any] struct {
	key   K
	value V
}

type BtreeMap[K btree.Ordered, V any] struct {
	tree *btree.BTreeG[BTreeMapItem[K, V]]
}

func NewBtreeMap[K btree.Ordered, V any](degree int) *BtreeMap[K, V] {
	if degree < 2 {
		degree = 2
	}

	return &BtreeMap[K, V]{
		tree: btree.NewG(degree, func(a, b BTreeMapItem[K, V]) bool {
			return a.key < b.key
		}),
	}
}

func (m *BtreeMap[K, V]) Get(key K) V {
	item, exists := m.tree.Get(BTreeMapItem[K, V]{key: key})
	if !exists {
		return *new(V)
	}

	return item.value
}

func (m *BtreeMap[K, V]) Keys() []K {
	result := make([]K, 0, m.tree.Len())
	m.tree.Ascend(func(item BTreeMapItem[K, V]) bool {
		result = append(result, item.key)

		return true
	})

	return result
}

func (m *BtreeMap[K, V]) Insert(key K, value V) {
	m.tree.ReplaceOrInsert(BTreeMapItem[K, V]{key: key, value: value})
}

type BTreeSet[T btree.Ordered] BtreeMap[T, struct{}]

func NewBtreeSet[T btree.Ordered](degree int) *BTreeSet[T] {
	if degree < 2 {
		degree = 2
	}

	btreeMap := NewBtreeMap[T, struct{}](degree)
	return (*BTreeSet[T])(btreeMap)
}

func (s *BTreeSet[T]) Len() int {
	return s.tree.Len()
}

func (s *BTreeSet[T]) Has(key T) bool {
	return s.tree.Has(BTreeMapItem[T, struct{}]{key: key})
}

func (s *BTreeSet[T]) Keys() []T {
	result := make([]T, 0, s.tree.Len())
	s.tree.Ascend(func(item BTreeMapItem[T, struct{}]) bool {
		result = append(result, item.key)

		return true
	})

	return result
}

func (s *BTreeSet[T]) Insert(key T) {
	s.tree.ReplaceOrInsert(BTreeMapItem[T, struct{}]{key: key})
}
