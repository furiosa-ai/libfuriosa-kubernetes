package npu_allocator

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/bradfitz/iter"
	"github.com/stretchr/testify/assert"
)

func TestDeviceSetContains(t *testing.T) {
	tests := []struct {
		description string
		source      DeviceSet
		target      DeviceSet
		expected    bool
	}{
		{
			description: "compare empty source and empty target",
			source:      nil,
			target:      nil,
			expected:    false,
		},
		{
			description: "compare source and empty target",
			source: DeviceSet{
				buildMockDevice(3),
				buildMockDevice(5),
			},
			target:   nil,
			expected: false,
		},
		{
			description: "compare empty source and target",
			source:      nil,
			target: DeviceSet{
				buildMockDevice(3),
				buildMockDevice(5),
			},
			expected: false,
		},
		{
			description: "compare source and subset",
			source: DeviceSet{
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(6),
				buildMockDevice(7),
			},
			target: DeviceSet{
				buildMockDevice(3),
				buildMockDevice(7),
			},
			expected: true,
		},
		{
			description: "compare source and non subset",
			source: DeviceSet{
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(6),
				buildMockDevice(7),
			},
			target: DeviceSet{
				buildMockDevice(1),
				buildMockDevice(4),
			},
			expected: false,
		},
		{
			description: "compare subset and target",
			source: DeviceSet{
				buildMockDevice(3),
				buildMockDevice(5),
			},
			target: DeviceSet{
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(6),
				buildMockDevice(7),
			},
			expected: false,
		},
		{
			description: "compare identical DeviceSets",
			source: DeviceSet{
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(6),
				buildMockDevice(7),
			},
			target: DeviceSet{
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(6),
				buildMockDevice(7),
			},
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			actual := tc.source.Contains(tc.target)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestDeviceSetSort(t *testing.T) {
	tests := []struct {
		description string
		source      DeviceSet
		expected    DeviceSet
	}{
		{
			description: "sort nil device set",
			source:      nil,
			expected:    nil,
		},
		{
			description: "sort empty device set",
			source:      DeviceSet{},
			expected:    DeviceSet{},
		},
		{
			description: "sort unsorted device set",
			source: DeviceSet{
				buildMockDevice(7),
				buildMockDevice(5),
				buildMockDevice(6),
				buildMockDevice(3),
			},
			expected: DeviceSet{
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(6),
				buildMockDevice(7),
			},
		},
		{
			description: "sort sorted device set",
			source: DeviceSet{
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(6),
				buildMockDevice(7),
			},
			expected: DeviceSet{
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(6),
				buildMockDevice(7),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			tc.source.Sort()
			for idx, sourceDevice := range tc.source {

				assert.Equal(t, tc.expected[idx].ID(), sourceDevice.ID())
				assert.Equal(t, tc.expected[idx].TopologyHintKey(), tc.source[idx].TopologyHintKey())
			}
		})
	}
}

func TestDeviceSetEqual(t *testing.T) {
	tests := []struct {
		description string
		source      DeviceSet
		target      DeviceSet
		expected    bool
	}{
		{
			description: "compare nil DeviceSets",
			source:      nil,
			target:      nil,
			expected:    true,
		},
		{
			description: "compare empty DeviceSets",
			source:      DeviceSet{},
			target:      DeviceSet{},
			expected:    true,
		},
		{
			description: "compare un-identical DeviceSets",
			source: DeviceSet{
				buildMockDevice(0),
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(7),
			},
			target: DeviceSet{
				buildMockDevice(1),
				buildMockDevice(2),
				buildMockDevice(4),
				buildMockDevice(6),
			},
			expected: false,
		},
		{
			description: "compare un-identical DeviceSets with intersection",
			source: DeviceSet{
				buildMockDevice(0),
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(7),
			},
			target: DeviceSet{
				buildMockDevice(1),
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(6),
			},
			expected: false,
		},
		{
			description: "compare identical DeviceSets in different order",
			source: DeviceSet{
				buildMockDevice(0),
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(7),
			},
			target: DeviceSet{
				buildMockDevice(7),
				buildMockDevice(3),
				buildMockDevice(0),
				buildMockDevice(5),
			},
			expected: true,
		},
		{
			description: "compare identical DeviceSets in same order",
			source: DeviceSet{
				buildMockDevice(0),
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(7),
			},
			target: DeviceSet{
				buildMockDevice(0),
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(7),
			},
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			actual := tc.source.Equal(tc.target)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestDeviceSetDifference(t *testing.T) {
	tests := []struct {
		description string
		source      DeviceSet
		target      DeviceSet
		expected    DeviceSet
	}{
		{
			description: "diff nil DeviceSets",
			source:      nil,
			target:      nil,
			expected:    DeviceSet{},
		},
		{
			description: "diff empty DeviceSets",
			source:      DeviceSet{},
			target:      DeviceSet{},
			expected:    DeviceSet{},
		},
		{
			description: "diff source and empty target",
			source: DeviceSet{
				buildMockDevice(0),
				buildMockDevice(1),
			},
			target: nil,
			expected: DeviceSet{
				buildMockDevice(0),
				buildMockDevice(1),
			},
		},
		{
			description: "diff source and target without intersection",
			source: DeviceSet{
				buildMockDevice(0),
				buildMockDevice(1),
			},
			target: DeviceSet{
				buildMockDevice(2),
				buildMockDevice(3),
			},
			expected: DeviceSet{
				buildMockDevice(0),
				buildMockDevice(1),
			},
		},
		{
			description: "diff empty source and target",
			source:      nil,
			target: DeviceSet{
				buildMockDevice(0),
				buildMockDevice(1),
			},
			expected: DeviceSet{},
		},
		{
			description: "diff source and target with intersection",
			source: DeviceSet{
				buildMockDevice(0),
				buildMockDevice(1),
				buildMockDevice(2),
				buildMockDevice(3),
			},
			target: DeviceSet{
				buildMockDevice(2),
				buildMockDevice(3),
			},
			expected: DeviceSet{
				buildMockDevice(0),
				buildMockDevice(1),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			actual := tc.source.Difference(tc.target)
			if !actual.Equal(tc.expected) {
				t.Errorf("expected %v but got %v", tc.expected, actual)
			}
		})
	}
}

func TestDeviceSetUnion(t *testing.T) {
	tests := []struct {
		description string
		source      DeviceSet
		target      DeviceSet
		expected    DeviceSet
	}{
		{
			description: "Union nil DeviceSets",
			source:      nil,
			target:      nil,
			expected:    DeviceSet{},
		},
		{
			description: "Union empty DeviceSets",
			source:      DeviceSet{},
			target:      DeviceSet{},
			expected:    DeviceSet{},
		},
		{
			description: "Union empty target",
			source: DeviceSet{
				buildMockDevice(0),
				buildMockDevice(1),
			},
			target: DeviceSet{},
			expected: DeviceSet{
				buildMockDevice(0),
				buildMockDevice(1),
			},
		},
		{
			description: "Union empty source",
			source:      DeviceSet{},
			target: DeviceSet{
				buildMockDevice(0),
				buildMockDevice(1),
			},
			expected: DeviceSet{
				buildMockDevice(0),
				buildMockDevice(1),
			},
		},
		{
			description: "Union source and target",
			source: DeviceSet{
				buildMockDevice(0),
				buildMockDevice(1),
			},
			target: DeviceSet{
				buildMockDevice(2),
				buildMockDevice(3),
			},
			expected: DeviceSet{
				buildMockDevice(0),
				buildMockDevice(1),
				buildMockDevice(2),
				buildMockDevice(3),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			actual := tc.source.Union(tc.target)
			if !actual.Equal(tc.expected) {
				t.Errorf("expected %v but got %v", tc.expected, actual)
			}
		})
	}
}

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

			sut := NewBtreeMap[int, struct{}](10)
			for idx := range assign {
				sut.ReplaceOrInsert(assign[idx], struct{}{})
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

			sut := NewBtreeSet[int](10)
			for idx := range assign {
				sut.ReplaceOrInsert(assign[idx])
			}

			actual := sut.Keys()

			assert.Equal(t, expected, actual)
		})
	}
}
