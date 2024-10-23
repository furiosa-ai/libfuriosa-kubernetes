package npu_allocator

import (
	"reflect"
	"testing"

	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
)

func TestGenerateNonDuplicatedDeviceSet(t *testing.T) {
	tests := []struct {
		description string
		devices     DeviceSet
		size        int
		expected    []DeviceSet
	}{
		{
			description: "empty input",
			devices:     nil,
			size:        4,
			expected:    []DeviceSet{},
		},
		{
			description: "size 0",
			devices:     nil,
			size:        0,
			expected:    []DeviceSet{},
		},
		{
			description: "size 1",
			devices:     buildMockDeviceSet(0, 3),
			size:        1,
			expected: func() []DeviceSet {
				deviceSets := make([]DeviceSet, 0)
				for i := 0; i <= 3; i++ {
					deviceSets = append(deviceSets, []Device{buildMockDevice(i)})
				}

				return deviceSets
			}(),
		},
		{
			description: "size greater than input slice length",
			devices:     buildMockDeviceSet(0, 7),
			size:        10,
			expected:    []DeviceSet{},
		},
		{
			description: "size equal to input slice length",
			devices:     buildMockDeviceSet(0, 7),
			size:        8,
			expected: []DeviceSet{
				buildMockDeviceSet(0, 7),
			},
		},
		{
			description: "generate combinations of two from eight",
			devices:     buildMockDeviceSet(0, 7),
			size:        2,
			expected: func() []DeviceSet {
				deviceSets := make([]DeviceSet, 0)
				for i := 0; i <= 7; i++ {
					for j := i + 1; j <= 7; j++ {
						deviceSets = append(deviceSets, []Device{
							buildMockDevice(i),
							buildMockDevice(j),
						})
					}
				}

				return deviceSets
			}(),
		},
	}

	for _, tc := range tests {
		actual := generateKDeviceSet(tc.devices, tc.size)
		if len(actual) != len(tc.expected) {
			t.Errorf("two slices are not identical")
			continue
		}

		for idx, inner := range tc.expected {
			inner2 := actual[idx]

			if !inner.Equal(inner2) {
				t.Errorf("two slices are not identical")
			}
		}
	}
}

func mockTopologyHintProvider(hints TopologyHintMatrix) TopologyHintProvider {
	return func(device1, device2 Device) uint {
		topologyHintKey1 := device1.TopologyHintKey()
		topologyHintKey2 := device2.TopologyHintKey()

		if topologyHintKey1 > topologyHintKey2 {
			topologyHintKey1, topologyHintKey2 = topologyHintKey2, topologyHintKey1
		}

		if hint, ok := hints[topologyHintKey1][topologyHintKey2]; ok {
			return hint
		}

		return 0
	}
}

func buildStaticHintMatrixForTwoSocketBalancedConfig() TopologyHintMatrix {
	return TopologyHintMatrix{
		"0": {"0": 70, "1": 30, "2": 20, "3": 20, "4": 10, "5": 10, "6": 10, "7": 10},
		"1": {"1": 70, "2": 20, "3": 20, "4": 10, "5": 10, "6": 10, "7": 10},
		"2": {"2": 70, "3": 30, "4": 10, "5": 10, "6": 10, "7": 10},
		"3": {"3": 70, "4": 10, "5": 10, "6": 10, "7": 10},
		"4": {"4": 70, "5": 30, "6": 20, "7": 20},
		"5": {"5": 70, "6": 20, "7": 20},
		"6": {"6": 70, "7": 30},
		"7": {"7": 70},
	}
}

// TODO(@bg): add hint matrix and test for up to four socket configuration
// TODO(@bg): add hint matrix and test for non-optimal configuration
func TestAllocation(t *testing.T) {
	tests := []struct {
		description string
		available   DeviceSet
		required    DeviceSet
		request     int
		hints       TopologyHintMatrix
		expected    DeviceSet
	}{
		{
			description: "[topology hint 0x0011] request eight devices from total eight devices",
			available:   buildMockDeviceSet(0, 7),
			required:    nil,
			request:     8,
			hints:       buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected:    buildMockDeviceSet(0, 7),
		},
		{
			description: "[topology hint 0x0011] request six devices from total eight devices",
			available:   buildMockDeviceSet(0, 7),
			required:    nil,
			request:     6,
			hints:       buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected:    buildMockDeviceSet(0, 5),
		},
		{
			description: "[topology hint 0x0011] request four devices from total eight devices",
			available:   buildMockDeviceSet(0, 7),
			required:    nil,
			request:     4,
			hints:       buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected:    buildMockDeviceSet(0, 3),
		},
		{
			description: "[topology hint 0x0001] request four devices from filtered devices",
			available:   buildMockDeviceSet(0, 3),
			required:    nil,
			request:     4,
			hints:       buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected:    buildMockDeviceSet(0, 3),
		},
		{
			description: "[topology hint 0x0001] request two devices from filtered devices",
			available:   buildMockDeviceSet(0, 3),
			required:    nil,
			request:     2,
			hints:       buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected:    buildMockDeviceSet(0, 1),
		},
		{
			description: "[topology hint 0x0010] request four devices from filtered devices",
			available:   buildMockDeviceSet(4, 7),
			required:    nil,
			request:     4,
			hints:       buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected:    buildMockDeviceSet(4, 7),
		},
		{
			description: "[topology hint 0x0010] request two devices from filtered devices",
			available:   buildMockDeviceSet(4, 7),
			required:    nil,
			request:     2,
			hints:       buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected:    buildMockDeviceSet(4, 5),
		},
		{
			description: "[topology hint 0x0011] request four devices from five devices",
			available: DeviceSet{
				buildMockDevice(0),
				buildMockDevice(1),
				buildMockDevice(3),
				buildMockDevice(4),
				buildMockDevice(7),
			},
			required: nil,
			request:  4,
			hints:    buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: DeviceSet{
				buildMockDevice(0),
				buildMockDevice(1),
				buildMockDevice(3),
				buildMockDevice(4),
			},
		},
		{
			description: "[topology hint 0x0011] request four devices from five devices, require specific device",
			available: DeviceSet{
				buildMockDevice(0),
				buildMockDevice(1),
				buildMockDevice(3),
				buildMockDevice(4),
				buildMockDevice(7),
			},
			required: DeviceSet{
				buildMockDevice(7),
			},
			request: 4,
			hints:   buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: DeviceSet{
				buildMockDevice(0),
				buildMockDevice(1),
				buildMockDevice(3),
				buildMockDevice(7),
			},
		},
		{
			description: "[topology hint 0x0010] request four devices from eight devices, require specific device set",
			available:   buildMockDeviceSet(0, 7),
			required:    buildMockDeviceSet(4, 5),
			request:     4,
			hints:       buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: DeviceSet{
				buildMockDevice(4),
				buildMockDevice(5),
				buildMockDevice(6),
				buildMockDevice(7),
			},
		},
		{
			description: "[no topology hint] request four devices from six devices",
			available: DeviceSet{
				buildMockDevice(1),
				buildMockDevice(3),
				buildMockDevice(4),
				buildMockDevice(5),
				buildMockDevice(7),
			},
			required: nil,
			request:  4,
			hints:    buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: DeviceSet{
				buildMockDevice(1),
				buildMockDevice(4),
				buildMockDevice(5),
				buildMockDevice(7),
			},
		},
		{
			description: "[no topology hint] request four devices from 6 devices, require specific device set",
			available: DeviceSet{
				buildMockDevice(1),
				buildMockDevice(3),
				buildMockDevice(4),
				buildMockDevice(5),
				buildMockDevice(7),
			},
			required: DeviceSet{
				buildMockDevice(1),
				buildMockDevice(3),
			},
			request: 4,
			hints:   buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: DeviceSet{
				buildMockDevice(1),
				buildMockDevice(3),
				buildMockDevice(4),
				buildMockDevice(5),
			},
		},
		{
			description: "[no topology hint] request two devices from five devices",
			available: DeviceSet{
				buildMockDevice(0),
				buildMockDevice(2),
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(6),
			},
			required: nil,
			request:  2,
			hints:    buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: DeviceSet{
				buildMockDevice(2),
				buildMockDevice(3),
			},
		},
		{
			description: "[no topology hint] request two devices from five devices, require specific device set",
			available: DeviceSet{
				buildMockDevice(0),
				buildMockDevice(2),
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(6),
			},
			required: DeviceSet{
				buildMockDevice(5),
			},
			request: 2,
			hints:   buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: DeviceSet{
				buildMockDevice(5),
				buildMockDevice(6),
			},
		},
		{
			description: "[no topology hint] request two devices from eight devices, require specific device set",
			available:   buildMockDeviceSet(0, 7),
			required: DeviceSet{
				buildMockDevice(4),
			},
			request: 2,
			hints:   buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: DeviceSet{
				buildMockDevice(4),
				buildMockDevice(5),
			},
		},
		{
			description: "[no topology hint] request one device from eight devices",
			available:   buildMockDeviceSet(0, 7),
			required:    nil,
			request:     1,
			hints:       buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: DeviceSet{
				buildMockDevice(0),
			},
		},
		{
			description: "[no topology hint] request one device from eight devices, require specific device set",
			available:   buildMockDeviceSet(0, 7),
			required: DeviceSet{
				buildMockDevice(3),
			},
			request: 1,
			hints:   buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: DeviceSet{
				buildMockDevice(3),
			},
		},
		{
			description: "[no topology hint] request one device from three devices",
			available: DeviceSet{
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(7),
			},
			required: nil,
			request:  1,
			hints:    buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: DeviceSet{
				buildMockDevice(3),
			},
		},
		{
			description: "[no topology hint] request one device from three devices, require specific device",
			available: DeviceSet{
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(7),
			},
			required: DeviceSet{
				buildMockDevice(7),
			},
			request: 1,
			hints:   buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: DeviceSet{
				buildMockDevice(7),
			},
		},
		{
			description: "[reboot] allocate reserved resources",
			available:   buildMockDeviceSet(0, 7),
			required:    buildMockDeviceSet(0, 3),
			request:     4,
			hints:       buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected:    buildMockDeviceSet(0, 3),
		},
	}

	for _, tc := range tests {
		allocator, _ := NewMockScoreBasedOptimalNpuAllocator(mockTopologyHintProvider(tc.hints))
		actualResult := allocator.Allocate(tc.available, tc.required, tc.request)

		if len(actualResult) != len(tc.expected) {
			t.Errorf("expected %v but got %v", tc.expected, actualResult)
		}

		actualResult.Sort()
		tc.expected.Sort()

		for idx, actual := range actualResult {
			if actual.ID() != tc.expected[idx].ID() || actual.TopologyHintKey() != tc.expected[idx].TopologyHintKey() {
				t.Errorf("expected %v but got %v", actual.(*mockDevice), tc.expected[idx].(*mockDevice))
				break
			}
		}
	}
}

func TestPopulateTopologyMatrix(t *testing.T) {
	tests := []struct {
		description string
		input       []smi.Device
		expected    TopologyHintMatrix
	}{
		{
			description: "test 8 npu configuration",
			input:       smi.GetStaticMockDevices(smi.ArchWarboy),
			expected: TopologyHintMatrix{
				"27": {"27": 70, "2a": 30, "51": 20, "57": 20, "9e": 10, "a4": 10, "c7": 10, "ca": 10},
				"2a": {"2a": 70, "51": 20, "57": 20, "9e": 10, "a4": 10, "c7": 10, "ca": 10},
				"51": {"51": 70, "57": 30, "9e": 10, "a4": 10, "c7": 10, "ca": 10},
				"57": {"57": 70, "9e": 10, "a4": 10, "c7": 10, "ca": 10},
				"9e": {"9e": 70, "a4": 30, "c7": 20, "ca": 20},
				"a4": {"a4": 70, "c7": 20, "ca": 20},
				"c7": {"c7": 70, "ca": 30},
				"ca": {"ca": 70},
			},
		},
	}

	for _, tc := range tests {
		actual, _ := NewTopologyHintMatrix(tc.input)

		if !reflect.DeepEqual(actual, tc.expected) {
			t.Errorf("expected %v but got %v", tc.expected, actual)
		}
	}
}
