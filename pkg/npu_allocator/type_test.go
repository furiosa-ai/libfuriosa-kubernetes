package npu_allocator

import (
	"testing"

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
			source:      NewDeviceSet(),
			target:      NewDeviceSet(),
			expected:    false,
		},
		{
			description: "compare source and empty target",
			source: NewDeviceSet(
				buildMockDevice(3),
				buildMockDevice(5),
			),
			target:   NewDeviceSet(),
			expected: false,
		},
		{
			description: "compare empty source and target",
			source:      NewDeviceSet(),
			target: NewDeviceSet(
				buildMockDevice(3),
				buildMockDevice(5),
			),
			expected: false,
		},
		{
			description: "compare source and subset",
			source: NewDeviceSet(
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(6),
				buildMockDevice(7),
			),
			target: NewDeviceSet(
				buildMockDevice(3),
				buildMockDevice(7),
			),
			expected: true,
		},
		{
			description: "compare source and non subset",
			source: NewDeviceSet(
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(6),
				buildMockDevice(7),
			),
			target: NewDeviceSet(
				buildMockDevice(1),
				buildMockDevice(4),
			),
			expected: false,
		},
		{
			description: "compare subset and target",
			source: NewDeviceSet(
				buildMockDevice(3),
				buildMockDevice(5),
			),
			target: NewDeviceSet(
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(6),
				buildMockDevice(7),
			),
			expected: false,
		},
		{
			description: "compare identical DeviceSets",
			source: NewDeviceSet(
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(6),
				buildMockDevice(7),
			),
			target: NewDeviceSet(
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(6),
				buildMockDevice(7),
			),
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			actual := tc.source.Contains(tc.target.Devices()...)

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
			description: "sort empty device set",
			source:      NewDeviceSet(),
			expected:    NewDeviceSet(),
		},
		{
			description: "sort unsorted device set",
			source: NewDeviceSet(
				buildMockDevice(7),
				buildMockDevice(5),
				buildMockDevice(6),
				buildMockDevice(3),
			),
			expected: NewDeviceSet(
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(6),
				buildMockDevice(7),
			),
		},
		{
			description: "sort sorted device set",
			source: NewDeviceSet(
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(6),
				buildMockDevice(7),
			),
			expected: NewDeviceSet(
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(6),
				buildMockDevice(7),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			source := tc.source.Devices()
			expected := tc.expected.Devices()

			for idx, sourceDevice := range tc.source.Devices() {
				assert.Equal(t, expected[idx].ID(), sourceDevice.ID())
				assert.Equal(t, expected[idx].TopologyHintKey(), source[idx].TopologyHintKey())
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
			description: "compare empty DeviceSets",
			source:      NewDeviceSet(),
			target:      NewDeviceSet(),
			expected:    true,
		},
		{
			description: "compare un-identical DeviceSets",
			source: NewDeviceSet(
				buildMockDevice(0),
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(7),
			),
			target: NewDeviceSet(
				buildMockDevice(1),
				buildMockDevice(2),
				buildMockDevice(4),
				buildMockDevice(6),
			),
			expected: false,
		},
		{
			description: "compare un-identical DeviceSets with intersection",
			source: NewDeviceSet(
				buildMockDevice(0),
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(7),
			),
			target: NewDeviceSet(
				buildMockDevice(1),
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(6),
			),
			expected: false,
		},
		{
			description: "compare identical DeviceSets in different order",
			source: NewDeviceSet(
				buildMockDevice(0),
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(7),
			),
			target: NewDeviceSet(
				buildMockDevice(7),
				buildMockDevice(3),
				buildMockDevice(0),
				buildMockDevice(5),
			),
			expected: true,
		},
		{
			description: "compare identical DeviceSets in same order",
			source: NewDeviceSet(
				buildMockDevice(0),
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(7),
			),
			target: NewDeviceSet(
				buildMockDevice(0),
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(7),
			),
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			actual := tc.source.Equal(tc.target.Devices()...)

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
			source:      NewDeviceSet(),
			target:      NewDeviceSet(),
			expected:    NewDeviceSet(),
		},
		{
			description: "diff empty DeviceSets",
			source:      NewDeviceSet(),
			target:      NewDeviceSet(),
			expected:    NewDeviceSet(),
		},
		{
			description: "diff source and empty target",
			source: NewDeviceSet(
				buildMockDevice(0),
				buildMockDevice(1),
			),
			target: NewDeviceSet(),
			expected: NewDeviceSet(
				buildMockDevice(0),
				buildMockDevice(1),
			),
		},
		{
			description: "diff source and target without intersection",
			source: NewDeviceSet(
				buildMockDevice(0),
				buildMockDevice(1),
			),
			target: NewDeviceSet(
				buildMockDevice(2),
				buildMockDevice(3),
			),
			expected: NewDeviceSet(
				buildMockDevice(0),
				buildMockDevice(1),
			),
		},
		{
			description: "diff empty source and target",
			source:      NewDeviceSet(),
			target: NewDeviceSet(
				buildMockDevice(0),
				buildMockDevice(1),
			),
			expected: NewDeviceSet(),
		},
		{
			description: "diff source and target with intersection",
			source: NewDeviceSet(
				buildMockDevice(0),
				buildMockDevice(1),
				buildMockDevice(2),
				buildMockDevice(3),
			),
			target: NewDeviceSet(
				buildMockDevice(2),
				buildMockDevice(3),
			),
			expected: NewDeviceSet(
				buildMockDevice(0),
				buildMockDevice(1),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			target := tc.target.Devices()
			expected := tc.expected.Devices()

			actual := tc.source.Difference(target...).Devices()
			assert.Equal(t, expected, actual)
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
			source:      NewDeviceSet(),
			target:      NewDeviceSet(),
			expected:    NewDeviceSet(),
		},
		{
			description: "Union empty DeviceSets",
			source:      NewDeviceSet(),
			target:      NewDeviceSet(),
			expected:    NewDeviceSet(),
		},
		{
			description: "Union empty target",
			source: NewDeviceSet(
				buildMockDevice(0),
				buildMockDevice(1),
			),
			target: NewDeviceSet(),
			expected: NewDeviceSet(
				buildMockDevice(0),
				buildMockDevice(1),
			),
		},
		{
			description: "Union empty source",
			source:      NewDeviceSet(),
			target: NewDeviceSet(
				buildMockDevice(0),
				buildMockDevice(1),
			),
			expected: NewDeviceSet(
				buildMockDevice(0),
				buildMockDevice(1),
			),
		},
		{
			description: "Union source and target",
			source: NewDeviceSet(
				buildMockDevice(0),
				buildMockDevice(1),
			),
			target: NewDeviceSet(
				buildMockDevice(2),
				buildMockDevice(3),
			),
			expected: NewDeviceSet(
				buildMockDevice(0),
				buildMockDevice(1),
				buildMockDevice(2),
				buildMockDevice(3),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			target := tc.target.Devices()
			expected := tc.expected.Devices()

			actual := tc.source.Union(target...).Devices()
			assert.Equal(t, expected, actual)
		})
	}
}
