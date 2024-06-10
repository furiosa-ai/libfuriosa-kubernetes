package npu_allocator

import (
	furiosaSmi "github.com/furiosa-ai/libfuriosa-kubernetes/pkg/smi"
	"reflect"
	"strconv"
	"testing"
)

func buildMockDeviceSet(start, end int) DeviceSet {
	result := DeviceSet{}
	for i := start; i <= end; i++ {
		result = append(result, &mockDevice{
			id:              strconv.Itoa(i),
			topologyHintKey: strconv.Itoa(i),
		})
	}

	return result
}

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
			expected: []DeviceSet{
				[]Device{
					&mockDevice{
						id:              "0",
						topologyHintKey: "0",
					},
				},
				[]Device{
					&mockDevice{
						id:              "1",
						topologyHintKey: "1",
					},
				},
				[]Device{
					&mockDevice{
						id:              "2",
						topologyHintKey: "2",
					},
				},
				[]Device{
					&mockDevice{
						id:              "3",
						topologyHintKey: "3",
					},
				},
			},
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
			expected: []DeviceSet{
				[]Device{
					&mockDevice{
						id:              "0",
						topologyHintKey: "0",
					},
					&mockDevice{
						id:              "1",
						topologyHintKey: "1",
					},
				},
				[]Device{
					&mockDevice{
						id:              "0",
						topologyHintKey: "0",
					},
					&mockDevice{
						id:              "2",
						topologyHintKey: "2",
					},
				},
				[]Device{
					&mockDevice{
						id:              "0",
						topologyHintKey: "0",
					},
					&mockDevice{
						id:              "3",
						topologyHintKey: "3",
					},
				},
				[]Device{
					&mockDevice{
						id:              "0",
						topologyHintKey: "0",
					},
					&mockDevice{
						id:              "4",
						topologyHintKey: "4",
					},
				},
				[]Device{
					&mockDevice{
						id:              "0",
						topologyHintKey: "0",
					},
					&mockDevice{
						id:              "5",
						topologyHintKey: "5",
					},
				},
				[]Device{
					&mockDevice{
						id:              "0",
						topologyHintKey: "0",
					},
					&mockDevice{
						id:              "6",
						topologyHintKey: "6",
					},
				},
				[]Device{
					&mockDevice{
						id:              "0",
						topologyHintKey: "0",
					},
					&mockDevice{
						id:              "7",
						topologyHintKey: "7",
					},
				},
				[]Device{
					&mockDevice{
						id:              "1",
						topologyHintKey: "1",
					},
					&mockDevice{
						id:              "2",
						topologyHintKey: "2",
					},
				},
				[]Device{
					&mockDevice{
						id:              "1",
						topologyHintKey: "1",
					},
					&mockDevice{
						id:              "3",
						topologyHintKey: "3",
					},
				},
				[]Device{
					&mockDevice{
						id:              "1",
						topologyHintKey: "1",
					},
					&mockDevice{
						id:              "4",
						topologyHintKey: "4",
					},
				},
				[]Device{
					&mockDevice{
						id:              "1",
						topologyHintKey: "1",
					},
					&mockDevice{
						id:              "5",
						topologyHintKey: "5",
					},
				},
				[]Device{
					&mockDevice{
						id:              "1",
						topologyHintKey: "1",
					},
					&mockDevice{
						id:              "6",
						topologyHintKey: "6",
					},
				},
				[]Device{
					&mockDevice{
						id:              "1",
						topologyHintKey: "1",
					},
					&mockDevice{
						id:              "7",
						topologyHintKey: "7",
					},
				},
				[]Device{
					&mockDevice{
						id:              "2",
						topologyHintKey: "2",
					},
					&mockDevice{
						id:              "3",
						topologyHintKey: "3",
					},
				},
				[]Device{
					&mockDevice{
						id:              "2",
						topologyHintKey: "2",
					},
					&mockDevice{
						id:              "4",
						topologyHintKey: "4",
					},
				},
				[]Device{
					&mockDevice{
						id:              "2",
						topologyHintKey: "2",
					},
					&mockDevice{
						id:              "5",
						topologyHintKey: "5",
					},
				},
				[]Device{
					&mockDevice{
						id:              "2",
						topologyHintKey: "2",
					},
					&mockDevice{
						id:              "6",
						topologyHintKey: "6",
					},
				},
				[]Device{
					&mockDevice{
						id:              "2",
						topologyHintKey: "2",
					},
					&mockDevice{
						id:              "7",
						topologyHintKey: "7",
					},
				},
				[]Device{
					&mockDevice{
						id:              "3",
						topologyHintKey: "3",
					},
					&mockDevice{
						id:              "4",
						topologyHintKey: "4",
					},
				},
				[]Device{
					&mockDevice{
						id:              "3",
						topologyHintKey: "3",
					},
					&mockDevice{
						id:              "5",
						topologyHintKey: "5",
					},
				},
				[]Device{
					&mockDevice{
						id:              "3",
						topologyHintKey: "3",
					},
					&mockDevice{
						id:              "6",
						topologyHintKey: "6",
					},
				},
				[]Device{
					&mockDevice{
						id:              "3",
						topologyHintKey: "3",
					},
					&mockDevice{
						id:              "7",
						topologyHintKey: "7",
					},
				},
				[]Device{
					&mockDevice{
						id:              "4",
						topologyHintKey: "4",
					},
					&mockDevice{
						id:              "5",
						topologyHintKey: "5",
					},
				},
				[]Device{
					&mockDevice{
						id:              "4",
						topologyHintKey: "4",
					},
					&mockDevice{
						id:              "6",
						topologyHintKey: "6",
					},
				},
				[]Device{
					&mockDevice{
						id:              "4",
						topologyHintKey: "4",
					},
					&mockDevice{
						id:              "7",
						topologyHintKey: "7",
					},
				},
				[]Device{
					&mockDevice{
						id:              "5",
						topologyHintKey: "5",
					},
					&mockDevice{
						id:              "6",
						topologyHintKey: "6",
					},
				},
				[]Device{
					&mockDevice{
						id:              "5",
						topologyHintKey: "5",
					},
					&mockDevice{
						id:              "7",
						topologyHintKey: "7",
					},
				},
				[]Device{
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

func mockTopologyHintProvider(hints map[string]map[string]uint) TopologyHintProvider {
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

func buildStaticHintMatrixForTwoSocketBalancedConfig() map[string]map[string]uint {
	return map[string]map[string]uint{
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
		hints       map[string]map[string]uint
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
				&mockDevice{
					id:              "0",
					topologyHintKey: "0",
				},
				&mockDevice{
					id:              "1",
					topologyHintKey: "1",
				},
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
				&mockDevice{
					id:              "4",
					topologyHintKey: "4",
				},
				&mockDevice{
					id:              "7",
					topologyHintKey: "7",
				},
			},
			required: nil,
			request:  4,
			hints:    buildStaticHintMatrixForTwoSocketBalancedConfig(),
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
					id:              "3",
					topologyHintKey: "3",
				},
				&mockDevice{
					id:              "4",
					topologyHintKey: "4",
				},
			},
		},
		{
			description: "[topology hint 0x0011] request four devices from five devices, require specific device",
			available: DeviceSet{
				&mockDevice{
					id:              "0",
					topologyHintKey: "0",
				},
				&mockDevice{
					id:              "1",
					topologyHintKey: "1",
				},
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
				&mockDevice{
					id:              "4",
					topologyHintKey: "4",
				},
				&mockDevice{
					id:              "7",
					topologyHintKey: "7",
				},
			},
			required: DeviceSet{
				&mockDevice{
					id:              "7",
					topologyHintKey: "7",
				},
			},
			request: 4,
			hints:   buildStaticHintMatrixForTwoSocketBalancedConfig(),
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
					id:              "3",
					topologyHintKey: "3",
				},
				&mockDevice{
					id:              "7",
					topologyHintKey: "7",
				},
			},
		},
		{
			description: "[topology hint 0x0010] request four devices from eight devices, require specific device set",
			available:   buildMockDeviceSet(0, 7),
			required:    buildMockDeviceSet(4, 5),
			request:     4,
			hints:       buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: DeviceSet{
				&mockDevice{
					id:              "4",
					topologyHintKey: "4",
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
			description: "[no topology hint] request four devices from six devices",
			available: DeviceSet{
				&mockDevice{
					id:              "1",
					topologyHintKey: "1",
				},
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
				&mockDevice{
					id:              "4",
					topologyHintKey: "4",
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
			required: nil,
			request:  4,
			hints:    buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: DeviceSet{
				&mockDevice{
					id:              "1",
					topologyHintKey: "1",
				},
				&mockDevice{
					id:              "4",
					topologyHintKey: "4",
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
		},
		{
			description: "[no topology hint] request four devices from 6 devices, require specific device set",
			available: DeviceSet{
				&mockDevice{
					id:              "1",
					topologyHintKey: "1",
				},
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
				&mockDevice{
					id:              "4",
					topologyHintKey: "4",
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
			required: DeviceSet{
				&mockDevice{
					id:              "1",
					topologyHintKey: "1",
				},
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
			},
			request: 4,
			hints:   buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: DeviceSet{
				&mockDevice{
					id:              "1",
					topologyHintKey: "1",
				},
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
				&mockDevice{
					id:              "4",
					topologyHintKey: "4",
				},
				&mockDevice{
					id:              "5",
					topologyHintKey: "5",
				},
			},
		},
		{
			description: "[no topology hint] request two devices from five devices",
			available: DeviceSet{
				&mockDevice{
					id:              "0",
					topologyHintKey: "0",
				},
				&mockDevice{
					id:              "2",
					topologyHintKey: "2",
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
			required: nil,
			request:  2,
			hints:    buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: DeviceSet{
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
		{
			description: "[no topology hint] request two devices from five devices, require specific device set",
			available: DeviceSet{
				&mockDevice{
					id:              "0",
					topologyHintKey: "0",
				},
				&mockDevice{
					id:              "2",
					topologyHintKey: "2",
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
			required: DeviceSet{
				&mockDevice{
					id:              "5",
					topologyHintKey: "5",
				},
			},
			request: 2,
			hints:   buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: DeviceSet{
				&mockDevice{
					id:              "5",
					topologyHintKey: "5",
				},
				&mockDevice{
					id:              "6",
					topologyHintKey: "6",
				},
			},
		},
		{
			description: "[no topology hint] request two devices from eight devices, require specific device set",
			available:   buildMockDeviceSet(0, 7),
			required: DeviceSet{
				&mockDevice{
					id:              "4",
					topologyHintKey: "4",
				},
			},
			request: 2,
			hints:   buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: DeviceSet{
				&mockDevice{
					id:              "4",
					topologyHintKey: "4",
				},
				&mockDevice{
					id:              "5",
					topologyHintKey: "5",
				},
			},
		},
		{
			description: "[no topology hint] request one device from eight devices",
			available:   buildMockDeviceSet(0, 7),
			required:    nil,
			request:     1,
			hints:       buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: DeviceSet{
				&mockDevice{
					id:              "0",
					topologyHintKey: "0",
				},
			},
		},
		{
			description: "[no topology hint] request one device from eight devices, require specific device set",
			available:   buildMockDeviceSet(0, 7),
			required: DeviceSet{
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
			},
			request: 1,
			hints:   buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: DeviceSet{
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
			},
		},
		{
			description: "[no topology hint] request one device from three devices",
			available: DeviceSet{
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
			required: nil,
			request:  1,
			hints:    buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: DeviceSet{
				&mockDevice{
					id:              "3",
					topologyHintKey: "3",
				},
			},
		},
		{
			description: "[no topology hint] request one device from three devices, require specific device",
			available: DeviceSet{
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
			required: DeviceSet{
				&mockDevice{
					id:              "7",
					topologyHintKey: "7",
				},
			},
			request: 1,
			hints:   buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: DeviceSet{
				&mockDevice{
					id:              "7",
					topologyHintKey: "7",
				},
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
		input       []furiosaSmi.Device
		expected    topologyMatrix
	}{
		{
			description: "test 8 npu configuration",
			input:       furiosaSmi.GetStaticMockDevices(furiosaSmi.ArchWarboy),
			expected: topologyMatrix{
				"0000:27:00.0": {
					"0000:27:00.0": 70,
					"0000:2a:00.0": 30,
					"0000:51:00.0": 20,
					"0000:57:00.0": 20,
					"0000:9e:00.0": 10,
					"0000:a4:00.0": 10,
					"0000:c7:00.0": 10,
					"0000:ca:00.0": 10,
				},
				"0000:2a:00.0": {
					"0000:2a:00.0": 70,
					"0000:51:00.0": 20,
					"0000:57:00.0": 20,
					"0000:9e:00.0": 10,
					"0000:a4:00.0": 10,
					"0000:c7:00.0": 10,
					"0000:ca:00.0": 10,
				},
				"0000:51:00.0": {
					"0000:51:00.0": 70,
					"0000:57:00.0": 30,
					"0000:9e:00.0": 10,
					"0000:a4:00.0": 10,
					"0000:c7:00.0": 10,
					"0000:ca:00.0": 10,
				},
				"0000:57:00.0": {
					"0000:57:00.0": 70,
					"0000:9e:00.0": 10,
					"0000:a4:00.0": 10,
					"0000:c7:00.0": 10,
					"0000:ca:00.0": 10,
				},
				"0000:9e:00.0": {
					"0000:9e:00.0": 70,
					"0000:a4:00.0": 30,
					"0000:c7:00.0": 20,
					"0000:ca:00.0": 20,
				},
				"0000:a4:00.0": {
					"0000:a4:00.0": 70,
					"0000:c7:00.0": 20,
					"0000:ca:00.0": 20,
				},
				"0000:c7:00.0": {
					"0000:c7:00.0": 70,
					"0000:ca:00.0": 30,
				},
				"0000:ca:00.0": {
					"0000:ca:00.0": 70,
				},
			},
		},
	}

	for _, tc := range tests {
		actual, _ := populateTopologyMatrix(tc.input)

		if !reflect.DeepEqual(actual, tc.expected) {
			t.Errorf("expected %v but got %v", tc.expected, actual)
		}
	}
}
