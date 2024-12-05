package util

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/bradfitz/iter"
	"github.com/stretchr/testify/assert"
)

func TestBtreeMap(t *testing.T) {
	numbers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	for i := range iter.N(3) {
		t.Run(fmt.Sprintf("[Trial %d] inject random order sequence", i+1), func(t *testing.T) {
			assign := make([]int, 10)
			copy(assign, numbers)

			rand.Shuffle(len(assign), func(i, j int) {
				assign[i], assign[j] = assign[j], assign[i]
			})

			expected := make([]int, 10)
			copy(expected, numbers)

			sut := NewBtreeMapWithLessFunc[int, struct{}](10, func(a, b BtreeMapItem[int, struct{}]) bool {
				return a.Key < b.Key
			})
			for idx := range assign {
				sut.Insert(assign[idx], struct{}{})
			}

			actual := sut.Keys()

			assert.Equal(t, expected, actual)
		})
	}
}

func TestBtreeSet(t *testing.T) {
	numbers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	for i := range iter.N(3) {
		t.Run(fmt.Sprintf("[Trial %d] inject random order sequence", i+1), func(t *testing.T) {
			assign := make([]int, 10)
			copy(assign, numbers)

			rand.Shuffle(len(assign), func(i, j int) {
				assign[i], assign[j] = assign[j], assign[i]
			})

			expected := make([]int, 10)
			copy(expected, numbers)

			sut := NewBtreeSetWithLessFunc[int](10, func(a, b int) bool {
				return a < b
			})
			for idx := range assign {
				sut.Insert(assign[idx])
			}

			actual := sut.Items()

			assert.Equal(t, expected, actual)
		})
	}
}
