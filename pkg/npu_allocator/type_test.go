package npu_allocator

import (
	"testing"
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
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
				&mockDevice{
					id:              "5",
					topologyHintKey: "5",
				},
			},
			target:   nil,
			expected: false,
		},
		{
			description: "compare empty source and target",
			source:      nil,
			target: DeviceSet{
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
				&mockDevice{
					id:              "5",
					topologyHintKey: "5",
				},
			},
			expected: false,
		},
		{
			description: "compare source and subset",
			source: DeviceSet{
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
				&mockDevice{
					id:              "5",
					topologyHintKey: "5",
				},
				&mockDevice{
					id:              "6",
					topologyHintKey: "6",
				},
				&mockDevice{
					id:              "7",
					topologyHintKey: "7",
				},
			},
			target: DeviceSet{
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
				&mockDevice{
					id:              "7",
					topologyHintKey: "7",
				},
			},
			expected: true,
		},
		{
			description: "compare source and non subset",
			source: DeviceSet{
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
				&mockDevice{
					id:              "5",
					topologyHintKey: "5",
				},
				&mockDevice{
					id:              "6",
					topologyHintKey: "6",
				},
				&mockDevice{
					id:              "7",
					topologyHintKey: "7",
				},
			},
			target: DeviceSet{
				&mockDevice{
					id:              "1",
					topologyHintKey: "1",
				},
				&mockDevice{
					id:              "4",
					topologyHintKey: "4",
				},
			},
			expected: false,
		},
		{
			description: "compare subset and target",
			source: DeviceSet{
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
				&mockDevice{
					id:              "5",
					topologyHintKey: "5",
				},
			},
			target: DeviceSet{
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
				&mockDevice{
					id:              "5",
					topologyHintKey: "5",
				},
				&mockDevice{
					id:              "6",
					topologyHintKey: "6",
				},
				&mockDevice{
					id:              "7",
					topologyHintKey: "7",
				},
			},
			expected: false,
		},
		{
			description: "compare identical DeviceSets",
			source: DeviceSet{
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
				&mockDevice{
					id:              "5",
					topologyHintKey: "5",
				},
				&mockDevice{
					id:              "6",
					topologyHintKey: "6",
				},
				&mockDevice{
					id:              "7",
					topologyHintKey: "7",
				},
			},
			target: DeviceSet{
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
				&mockDevice{
					id:              "5",
					topologyHintKey: "5",
				},
				&mockDevice{
					id:              "6",
					topologyHintKey: "6",
				},
				&mockDevice{
					id:              "7",
					topologyHintKey: "7",
				},
			},
			expected: true,
		},
	}

	for _, tc := range tests {
		actual := tc.source.Contains(tc.target)
		if tc.expected != actual {
			t.Errorf("expected %v but got %v", tc.expected, actual)
		}
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
				&mockDevice{
					id:              "7",
					topologyHintKey: "7",
				},
				&mockDevice{
					id:              "5",
					topologyHintKey: "5",
				},
				&mockDevice{
					id:              "6",
					topologyHintKey: "6",
				},
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
			},
			expected: DeviceSet{
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
				&mockDevice{
					id:              "5",
					topologyHintKey: "5",
				},
				&mockDevice{
					id:              "6",
					topologyHintKey: "6",
				},
				&mockDevice{
					id:              "7",
					topologyHintKey: "7",
				},
			},
		},
		{
			description: "sort sorted device set",
			source: DeviceSet{
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
				&mockDevice{
					id:              "5",
					topologyHintKey: "5",
				},
				&mockDevice{
					id:              "6",
					topologyHintKey: "6",
				},
				&mockDevice{
					id:              "7",
					topologyHintKey: "7",
				},
			},
			expected: DeviceSet{
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
				&mockDevice{
					id:              "5",
					topologyHintKey: "5",
				},
				&mockDevice{
					id:              "6",
					topologyHintKey: "6",
				},
				&mockDevice{
					id:              "7",
					topologyHintKey: "7",
				},
			},
		},
	}

	for _, tc := range tests {
		tc.source.Sort()
		for idx, sourceDevice := range tc.source {
			if sourceDevice.GetID() != tc.expected[idx].GetID() || sourceDevice.GetTopologyHintKey() != tc.expected[idx].GetTopologyHintKey() {
				t.Errorf("expected %v but got %v", sourceDevice.(*mockDevice), tc.expected[idx].(*mockDevice))
				break
			}
		}
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
				&mockDevice{
					id:              "0",
					topologyHintKey: "0",
				},
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
				&mockDevice{
					id:              "5",
					topologyHintKey: "5",
				},
				&mockDevice{
					id:              "7",
					topologyHintKey: "7",
				},
			},
			target: DeviceSet{
				&mockDevice{
					id:              "1",
					topologyHintKey: "1",
				},
				&mockDevice{
					id:              "2",
					topologyHintKey: "2",
				},
				&mockDevice{
					id:              "4",
					topologyHintKey: "4",
				},
				&mockDevice{
					id:              "6",
					topologyHintKey: "6",
				},
			},
			expected: false,
		},
		{
			description: "compare un-identical DeviceSets with intersection",
			source: DeviceSet{
				&mockDevice{
					id:              "0",
					topologyHintKey: "0",
				},
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
				&mockDevice{
					id:              "5",
					topologyHintKey: "5",
				},
				&mockDevice{
					id:              "7",
					topologyHintKey: "7",
				},
			},
			target: DeviceSet{
				&mockDevice{
					id:              "1",
					topologyHintKey: "1",
				},
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
				&mockDevice{
					id:              "5",
					topologyHintKey: "5",
				},
				&mockDevice{
					id:              "6",
					topologyHintKey: "6",
				},
			},
			expected: false,
		},
		{
			description: "compare identical DeviceSets in different order",
			source: DeviceSet{
				&mockDevice{
					id:              "0",
					topologyHintKey: "0",
				},
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
				&mockDevice{
					id:              "5",
					topologyHintKey: "5",
				},
				&mockDevice{
					id:              "7",
					topologyHintKey: "7",
				},
			},
			target: DeviceSet{
				&mockDevice{
					id:              "7",
					topologyHintKey: "7",
				},

				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
				&mockDevice{
					id:              "0",
					topologyHintKey: "0",
				},
				&mockDevice{
					id:              "5",
					topologyHintKey: "5",
				},
			},
			expected: true,
		},
		{
			description: "compare identical DeviceSets in same order",
			source: DeviceSet{
				&mockDevice{
					id:              "0",
					topologyHintKey: "0",
				},
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
				&mockDevice{
					id:              "5",
					topologyHintKey: "5",
				},
				&mockDevice{
					id:              "7",
					topologyHintKey: "7",
				},
			},
			target: DeviceSet{
				&mockDevice{
					id:              "0",
					topologyHintKey: "0",
				},
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
				&mockDevice{
					id:              "5",
					topologyHintKey: "5",
				},
				&mockDevice{
					id:              "7",
					topologyHintKey: "7",
				},
			},
			expected: true,
		},
	}

	for _, tc := range tests {
		actual := tc.source.Equal(tc.target)
		if actual != tc.expected {
			t.Errorf("expected %v, but got %v", tc.expected, actual)
		}
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
				&mockDevice{
					id:              "0",
					topologyHintKey: "0",
				},
				&mockDevice{
					id:              "1",
					topologyHintKey: "1",
				},
			},
			target: nil,
			expected: DeviceSet{
				&mockDevice{
					id:              "0",
					topologyHintKey: "0",
				},
				&mockDevice{
					id:              "1",
					topologyHintKey: "1",
				},
			},
		},
		{
			description: "diff source and target without intersection",
			source: DeviceSet{
				&mockDevice{
					id:              "0",
					topologyHintKey: "0",
				},
				&mockDevice{
					id:              "1",
					topologyHintKey: "1",
				},
			},
			target: DeviceSet{
				&mockDevice{
					id:              "2",
					topologyHintKey: "2",
				},
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
			},
			expected: DeviceSet{
				&mockDevice{
					id:              "0",
					topologyHintKey: "0",
				},
				&mockDevice{
					id:              "1",
					topologyHintKey: "1",
				},
			},
		},
		{
			description: "diff empty source and target",
			source:      nil,
			target: DeviceSet{
				&mockDevice{
					id:              "0",
					topologyHintKey: "0",
				},
				&mockDevice{
					id:              "1",
					topologyHintKey: "1",
				},
			},
			expected: DeviceSet{},
		},
		{
			description: "diff source and target with intersection",
			source: DeviceSet{
				&mockDevice{
					id:              "0",
					topologyHintKey: "0",
				},
				&mockDevice{
					id:              "1",
					topologyHintKey: "1",
				},
				&mockDevice{
					id:              "2",
					topologyHintKey: "2",
				},
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
			},
			target: DeviceSet{
				&mockDevice{
					id:              "2",
					topologyHintKey: "2",
				},
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
			},
			expected: DeviceSet{
				&mockDevice{
					id:              "0",
					topologyHintKey: "0",
				},
				&mockDevice{
					id:              "1",
					topologyHintKey: "1",
				},
			},
		},
	}

	for _, tc := range tests {
		actual := tc.source.Difference(tc.target)
		if !actual.Equal(tc.expected) {
			t.Errorf("expected %v but got %v", tc.expected, actual)
		}
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
				&mockDevice{
					id:              "0",
					topologyHintKey: "0",
				},
				&mockDevice{
					id:              "1",
					topologyHintKey: "1",
				},
			},
			target: DeviceSet{},
			expected: DeviceSet{
				&mockDevice{
					id:              "0",
					topologyHintKey: "0",
				},
				&mockDevice{
					id:              "1",
					topologyHintKey: "1",
				},
			},
		},
		{
			description: "Union empty source",
			source:      DeviceSet{},
			target: DeviceSet{
				&mockDevice{
					id:              "0",
					topologyHintKey: "0",
				},
				&mockDevice{
					id:              "1",
					topologyHintKey: "1",
				},
			},
			expected: DeviceSet{
				&mockDevice{
					id:              "0",
					topologyHintKey: "0",
				},
				&mockDevice{
					id:              "1",
					topologyHintKey: "1",
				},
			},
		},
		{
			description: "Union source and target",
			source: DeviceSet{
				&mockDevice{
					id:              "0",
					topologyHintKey: "0",
				},
				&mockDevice{
					id:              "1",
					topologyHintKey: "1",
				},
			},
			target: DeviceSet{
				&mockDevice{
					id:              "2",
					topologyHintKey: "2",
				},
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
			},
			expected: DeviceSet{
				&mockDevice{
					id:              "0",
					topologyHintKey: "0",
				},
				&mockDevice{
					id:              "1",
					topologyHintKey: "1",
				},
				&mockDevice{
					id:              "2",
					topologyHintKey: "2",
				},
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
			},
		},
	}

	for _, tc := range tests {
		actual := tc.source.Union(tc.target)
		if !actual.Equal(tc.expected) {
			t.Errorf("expected %v but got %v", tc.expected, actual)
		}
	}
}
