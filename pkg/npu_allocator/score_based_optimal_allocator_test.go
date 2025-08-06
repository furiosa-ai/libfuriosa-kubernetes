package npu_allocator

import (
	"testing"

	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"github.com/stretchr/testify/assert"
)

func TestGenerateNonDuplicatedDeviceSet(t *testing.T) {
	tests := []struct {
		description string
		deviceSet   DeviceSet
		size        int
		expected    []DeviceSet
	}{
		{
			description: "empty input",
			deviceSet:   NewDeviceSet(),
			size:        4,
			expected:    []DeviceSet{},
		},
		{
			description: "size 0",
			deviceSet:   NewDeviceSet(),
			size:        0,
			expected:    []DeviceSet{},
		},
		{
			description: "size 1",
			deviceSet:   buildMockDeviceSet(0, 3),
			size:        1,
			expected: func() []DeviceSet {
				deviceSets := make([]DeviceSet, 0)
				for i := 0; i <= 3; i++ {
					deviceSets = append(deviceSets, NewDeviceSet(buildMockDevice(i)))
				}

				return deviceSets
			}(),
		},
		{
			description: "size greater than input slice length",
			deviceSet:   buildMockDeviceSet(0, 7),
			size:        10,
			expected:    []DeviceSet{},
		},
		{
			description: "size equal to input slice length",
			deviceSet:   buildMockDeviceSet(0, 7),
			size:        8,
			expected: []DeviceSet{
				buildMockDeviceSet(0, 7),
			},
		},
		{
			description: "generate combinations of two from eight",
			deviceSet:   buildMockDeviceSet(0, 7),
			size:        2,
			expected: func() []DeviceSet {
				deviceSets := make([]DeviceSet, 0)
				for i := 0; i <= 7; i++ {
					for j := i + 1; j <= 7; j++ {
						deviceSets = append(deviceSets, NewDeviceSet(buildMockDevice(i), buildMockDevice(j)))
					}
				}

				return deviceSets
			}(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			actual := generateKDeviceSet(tc.deviceSet, tc.size)

			assert.Len(t, actual, len(tc.expected))

			for idx, inner := range tc.expected {
				inner2 := actual[idx]

				assert.Equal(t, inner.Devices(), inner2.Devices())
			}
		})
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
			required:    NewDeviceSet(),
			request:     8,
			hints:       buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected:    buildMockDeviceSet(0, 7),
		},
		{
			description: "[topology hint 0x0011] request six devices from total eight devices",
			available:   buildMockDeviceSet(0, 7),
			required:    NewDeviceSet(),
			request:     6,
			hints:       buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected:    buildMockDeviceSet(0, 5),
		},
		{
			description: "[topology hint 0x0011] request four devices from total eight devices",
			available:   buildMockDeviceSet(0, 7),
			required:    NewDeviceSet(),
			request:     4,
			hints:       buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected:    buildMockDeviceSet(0, 3),
		},
		{
			description: "[topology hint 0x0001] request four devices from filtered devices",
			available:   buildMockDeviceSet(0, 3),
			required:    NewDeviceSet(),
			request:     4,
			hints:       buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected:    buildMockDeviceSet(0, 3),
		},
		{
			description: "[topology hint 0x0001] request two devices from filtered devices",
			available:   buildMockDeviceSet(0, 3),
			required:    NewDeviceSet(),
			request:     2,
			hints:       buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected:    buildMockDeviceSet(0, 1),
		},
		{
			description: "[topology hint 0x0010] request four devices from filtered devices",
			available:   buildMockDeviceSet(4, 7),
			required:    NewDeviceSet(),
			request:     4,
			hints:       buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected:    buildMockDeviceSet(4, 7),
		},
		{
			description: "[topology hint 0x0010] request two devices from filtered devices",
			available:   buildMockDeviceSet(4, 7),
			required:    NewDeviceSet(),
			request:     2,
			hints:       buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected:    buildMockDeviceSet(4, 5),
		},
		{
			description: "[topology hint 0x0011] request four devices from five devices",
			available: NewDeviceSet(
				buildMockDevice(0),
				buildMockDevice(1),
				buildMockDevice(3),
				buildMockDevice(4),
				buildMockDevice(7),
			),
			required: NewDeviceSet(),
			request:  4,
			hints:    buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: NewDeviceSet(
				buildMockDevice(0),
				buildMockDevice(1),
				buildMockDevice(3),
				buildMockDevice(4),
			),
		},
		{
			description: "[topology hint 0x0011] request four devices from five devices, require specific device",
			available: NewDeviceSet(
				buildMockDevice(0),
				buildMockDevice(1),
				buildMockDevice(3),
				buildMockDevice(4),
				buildMockDevice(7),
			),
			required: NewDeviceSet(buildMockDevice(7)),
			request:  4,
			hints:    buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: NewDeviceSet(
				buildMockDevice(0),
				buildMockDevice(1),
				buildMockDevice(3),
				buildMockDevice(7),
			),
		},
		{
			description: "[topology hint 0x0010] request four devices from eight devices, require specific device set",
			available:   buildMockDeviceSet(0, 7),
			required:    buildMockDeviceSet(4, 5),
			request:     4,
			hints:       buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: NewDeviceSet(
				buildMockDevice(4),
				buildMockDevice(5),
				buildMockDevice(6),
				buildMockDevice(7),
			),
		},
		{
			description: "[no topology hint] request four devices from six devices",
			available: NewDeviceSet(
				buildMockDevice(1),
				buildMockDevice(3),
				buildMockDevice(4),
				buildMockDevice(5),
				buildMockDevice(7),
			),
			required: NewDeviceSet(),
			request:  4,
			hints:    buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: NewDeviceSet(
				buildMockDevice(1),
				buildMockDevice(4),
				buildMockDevice(5),
				buildMockDevice(7),
			),
		},
		{
			description: "[no topology hint] request four devices from 6 devices, require specific device set",
			available: NewDeviceSet(
				buildMockDevice(1),
				buildMockDevice(3),
				buildMockDevice(4),
				buildMockDevice(5),
				buildMockDevice(7),
			),
			required: NewDeviceSet(
				buildMockDevice(1),
				buildMockDevice(3),
			),
			request: 4,
			hints:   buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: NewDeviceSet(
				buildMockDevice(1),
				buildMockDevice(3),
				buildMockDevice(4),
				buildMockDevice(5),
			),
		},
		{
			description: "[no topology hint] request two devices from five devices",
			available: NewDeviceSet(
				buildMockDevice(0),
				buildMockDevice(2),
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(6),
			),
			required: NewDeviceSet(),
			request:  2,
			hints:    buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: NewDeviceSet(
				buildMockDevice(2),
				buildMockDevice(3),
			),
		},
		{
			description: "[no topology hint] request two devices from five devices, require specific device set",
			available: NewDeviceSet(
				buildMockDevice(0),
				buildMockDevice(2),
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(6),
			),
			required: NewDeviceSet(buildMockDevice(5)),
			request:  2,
			hints:    buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: NewDeviceSet(
				buildMockDevice(5),
				buildMockDevice(6),
			),
		},
		{
			description: "[no topology hint] request two devices from eight devices, require specific device set",
			available:   buildMockDeviceSet(0, 7),
			required:    NewDeviceSet(buildMockDevice(4)),
			request:     2,
			hints:       buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: NewDeviceSet(
				buildMockDevice(4),
				buildMockDevice(5),
			),
		},
		{
			description: "[no topology hint] request one device from eight devices",
			available:   buildMockDeviceSet(0, 7),
			required:    NewDeviceSet(),
			request:     1,
			hints:       buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected:    NewDeviceSet(buildMockDevice(0)),
		},
		{
			description: "[no topology hint] request one device from eight devices, require specific device set",
			available:   buildMockDeviceSet(0, 7),
			required:    NewDeviceSet(buildMockDevice(3)),
			request:     1,
			hints:       buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected:    NewDeviceSet(buildMockDevice(3)),
		},
		{
			description: "[no topology hint] request one device from three devices",
			available: NewDeviceSet(
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(7),
			),
			required: NewDeviceSet(),
			request:  1,
			hints:    buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: NewDeviceSet(buildMockDevice(3)),
		},
		{
			description: "[no topology hint] request one device from three devices, require specific device",
			available: NewDeviceSet(
				buildMockDevice(3),
				buildMockDevice(5),
				buildMockDevice(7),
			),
			required: NewDeviceSet(buildMockDevice(7)),
			request:  1,
			hints:    buildStaticHintMatrixForTwoSocketBalancedConfig(),
			expected: NewDeviceSet(buildMockDevice(7)),
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
		t.Run(tc.description, func(t *testing.T) {
			allocator, _ := NewMockScoreBasedOptimalNpuAllocator(mockTopologyHintProvider(tc.hints))
			actualResult := allocator.Allocate(tc.available, tc.required, tc.request)

			assert.Equal(t, tc.expected.Len(), actualResult.Len())

			expected := tc.expected.Devices()
			for idx, actual := range actualResult.Devices() {
				assert.Equalf(t, expected[idx].ID(), actual.ID(), "expected %v but got %v", actual.(*mockDevice), expected[idx].(*mockDevice))
				assert.Equalf(t, expected[idx].TopologyHintKey(), actual.TopologyHintKey(), "expected %v but got %v", actual.(*mockDevice), expected[idx].(*mockDevice))
			}
		})
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
			input:       smi.GetStaticMockDevices(smi.ArchRngd),
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
		t.Run(tc.description, func(t *testing.T) {
			actual, _ := NewTopologyHintMatrix(tc.input)

			assert.Equal(t, tc.expected, actual)
		})
	}
}
