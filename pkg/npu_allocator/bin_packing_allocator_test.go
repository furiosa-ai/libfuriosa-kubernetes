package npu_allocator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mockHintMatrix = TopologyHintMatrix{
		"0": {"0": 70, "1": 30, "2": 20, "3": 20, "4": 10, "5": 10, "6": 10, "7": 10},
		"1": {"1": 70, "2": 20, "3": 20, "4": 10, "5": 10, "6": 10, "7": 10},
		"2": {"2": 70, "3": 30, "4": 10, "5": 10, "6": 10, "7": 10},
		"3": {"3": 70, "4": 10, "5": 10, "6": 10, "7": 10},
		"4": {"4": 70, "5": 30, "6": 20, "7": 20},
		"5": {"5": 70, "6": 20, "7": 20},
		"6": {"6": 70, "7": 30},
		"7": {"7": 70},
	}

	mockHintProvider = func(device1, device2 Device) uint {
		key1, key2 := device1.GetTopologyHintKey(), device2.GetTopologyHintKey()
		if key1 > key2 {
			key1, key2 = key2, key1
		}

		if innerMap, innerMapExists := mockHintMatrix[key1]; innerMapExists {
			if score, scoreExists := innerMap[key2]; scoreExists {
				return score
			}
		}

		return 0
	}
)

func generateSameBoardMockDeviceSet(start int, end int, hintKey TopologyHintKey) DeviceSet {
	devices := make(DeviceSet, 0, end-start+1)
	for i := start; i < end; i++ {
		devices = append(devices, NewMockDevice(fmt.Sprintf("%s-%2d", hintKey, i), hintKey))
	}

	return devices
}

// TestSelectBestScoredDevices tests binPackingNpuAllocator.selectBestScoredDevices().
//   - It only tests single trial, not the final version of total allocated devices.
func TestSelectBestScoredDevices(t *testing.T) {
	mockBinPackingAllocator, _ := NewMockBinPackingNpuAllocator(mockHintProvider)
	sut := mockBinPackingAllocator.(*binPackingNpuAllocator)

	tests := []struct {
		description                string
		maxSelectLength            int
		previouslyAllocatedDevices DeviceSet
		remainingDevicesByHintMap  map[TopologyHintKey]DeviceSet
		expectedIn                 []TopologyHintKey
		expectedSelectedLength     int
	}{
		{
			description:                "any devices can be allocated because all devices in remainingDevicesByHintMap have equal size and previouslyAllocatedDevices is empty",
			maxSelectLength:            16,
			previouslyAllocatedDevices: DeviceSet{},
			remainingDevicesByHintMap: map[TopologyHintKey]DeviceSet{
				"0": generateSameBoardMockDeviceSet(0, 8, "0"),
				"1": generateSameBoardMockDeviceSet(0, 8, "1"),
				"2": generateSameBoardMockDeviceSet(0, 8, "2"),
				"3": generateSameBoardMockDeviceSet(0, 8, "3"),
				"4": generateSameBoardMockDeviceSet(0, 8, "4"),
				"5": generateSameBoardMockDeviceSet(0, 8, "5"),
				"6": generateSameBoardMockDeviceSet(0, 8, "6"),
				"7": generateSameBoardMockDeviceSet(0, 8, "7"),
			},
			expectedIn:             []TopologyHintKey{"0", "1", "2", "3", "4", "5", "6", "7"},
			expectedSelectedLength: 8, // new selection 8
		},
		{
			description:                "largest one should be picked because previouslyAllocatedDevices is empty",
			maxSelectLength:            16,
			previouslyAllocatedDevices: DeviceSet{},
			remainingDevicesByHintMap: map[TopologyHintKey]DeviceSet{
				"0": generateSameBoardMockDeviceSet(0, 1, "0"),
				"1": generateSameBoardMockDeviceSet(0, 6, "1"),
				"2": generateSameBoardMockDeviceSet(0, 2, "2"),
				"3": generateSameBoardMockDeviceSet(0, 4, "3"),
				"4": generateSameBoardMockDeviceSet(0, 7, "4"),
				"5": generateSameBoardMockDeviceSet(0, 8, "5"),
				"6": generateSameBoardMockDeviceSet(0, 5, "6"),
				"7": generateSameBoardMockDeviceSet(0, 3, "7"),
			},
			expectedIn:             []TopologyHintKey{"5"},
			expectedSelectedLength: 8, // new selection 8
		},
		{
			description:     "hintKey '1' must be picked",
			maxSelectLength: 15,
			previouslyAllocatedDevices: func() DeviceSet {
				devices := make(DeviceSet, 0)
				devices = devices.Union(generateSameBoardMockDeviceSet(0, 1, "1"))

				return devices
			}(),
			remainingDevicesByHintMap: map[TopologyHintKey]DeviceSet{
				"0": generateSameBoardMockDeviceSet(0, 4, "0"),
				"1": generateSameBoardMockDeviceSet(1, 5, "1"),
				"2": generateSameBoardMockDeviceSet(0, 4, "2"),
				"3": generateSameBoardMockDeviceSet(0, 4, "3"),
				"4": generateSameBoardMockDeviceSet(0, 4, "4"),
				"5": generateSameBoardMockDeviceSet(0, 4, "5"),
				"6": generateSameBoardMockDeviceSet(0, 4, "6"),
				"7": generateSameBoardMockDeviceSet(0, 4, "7"),
			},
			expectedIn:             []TopologyHintKey{"1"},
			expectedSelectedLength: 5, // existing 1, new selection 4
		},
		{
			description:     "7 devices from hintKey '3' must be picked",
			maxSelectLength: 7,
			previouslyAllocatedDevices: func() DeviceSet {
				devices := make(DeviceSet, 0)
				devices = devices.Union(generateSameBoardMockDeviceSet(0, 1, "1"))
				devices = devices.Union(generateSameBoardMockDeviceSet(0, 8, "2"))

				return devices
			}(),
			remainingDevicesByHintMap: map[TopologyHintKey]DeviceSet{
				"0": generateSameBoardMockDeviceSet(0, 4, "0"),
				"1": generateSameBoardMockDeviceSet(1, 8, "1"),
				// no available slots for hintKey "2"
				"3": generateSameBoardMockDeviceSet(0, 8, "3"),
				"4": generateSameBoardMockDeviceSet(0, 8, "4"),
				"5": generateSameBoardMockDeviceSet(0, 8, "5"),
				"6": generateSameBoardMockDeviceSet(0, 8, "6"),
				"7": generateSameBoardMockDeviceSet(0, 8, "7"),
			},
			expectedIn:             []TopologyHintKey{"1", "2", "3"},
			expectedSelectedLength: 16, // existing 9, new selection 7
		},
		{
			description:     "7 devices from hintKey '0' must be picked",
			maxSelectLength: 7,
			previouslyAllocatedDevices: func() DeviceSet {
				devices := make(DeviceSet, 0)
				devices = devices.Union(generateSameBoardMockDeviceSet(0, 1, "0"))
				devices = devices.Union(generateSameBoardMockDeviceSet(0, 8, "1"))

				return devices
			}(),
			remainingDevicesByHintMap: map[TopologyHintKey]DeviceSet{
				"0": generateSameBoardMockDeviceSet(1, 8, "0"),
				// no available slots for hintKey "1"
				"2": generateSameBoardMockDeviceSet(0, 8, "2"),
				"3": generateSameBoardMockDeviceSet(0, 8, "3"),
				"4": generateSameBoardMockDeviceSet(0, 8, "4"),
				"5": generateSameBoardMockDeviceSet(0, 8, "5"),
				"6": generateSameBoardMockDeviceSet(0, 8, "6"),
				"7": generateSameBoardMockDeviceSet(0, 8, "7"),
			},
			expectedIn:             []TopologyHintKey{"0", "1"},
			expectedSelectedLength: 16, // existing 9, new selection 7
		},
	}

	t.Parallel()
	for _, tc := range tests {
		t.Run(tc.description, func(subT *testing.T) {
			selectedDevices := sut.selectBestScoredDevices(tc.maxSelectLength, tc.previouslyAllocatedDevices, tc.remainingDevicesByHintMap)

			assert.Equal(subT, tc.expectedSelectedLength, len(selectedDevices))
			for _, device := range selectedDevices {
				assert.Contains(subT, tc.expectedIn, device.GetTopologyHintKey())
			}
		})
	}
}
