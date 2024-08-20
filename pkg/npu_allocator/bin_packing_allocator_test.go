package npu_allocator

import (
	"fmt"
	"testing"

	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/smi"
	"github.com/stretchr/testify/assert"
)

func generateSameBoardMockDeviceSet(start int, end int, hintKey TopologyHintKey) DeviceSet {
	devices := make(DeviceSet, 0, end-start+1)
	for i := start; i < end; i++ {
		devices = append(devices, NewMockDevice(fmt.Sprintf("%s-%02d", hintKey, i), hintKey))
	}

	return devices
}

// TestSelectBestScoredDevices tests binPackingNpuAllocator.selectBestScoredNewDevices().
//   - It only tests single trial, not the final version of total allocated devices.
func TestSelectBestScoredDevices(t *testing.T) {
	mockHintMatrix := TopologyHintMatrix{
		"0": {"0": 70, "1": 30, "2": 20, "3": 20, "4": 10, "5": 10, "6": 10, "7": 10},
		"1": {"1": 70, "2": 20, "3": 20, "4": 10, "5": 10, "6": 10, "7": 10},
		"2": {"2": 70, "3": 30, "4": 10, "5": 10, "6": 10, "7": 10},
		"3": {"3": 70, "4": 10, "5": 10, "6": 10, "7": 10},
		"4": {"4": 70, "5": 30, "6": 20, "7": 20},
		"5": {"5": 70, "6": 20, "7": 20},
		"6": {"6": 70, "7": 30},
		"7": {"7": 70},
	}

	mockHintProvider := func(device1, device2 Device) uint {
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
			expectedSelectedLength: 4, // new selection 4
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
			expectedSelectedLength: 7, // new selection 7
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
			expectedSelectedLength: 7, // new selection 7
		},
	}

	t.Parallel()
	for _, tc := range tests {
		t.Run(tc.description, func(subT *testing.T) {
			selectedDevices := sut.selectBestScoredNewDevices(tc.maxSelectLength, tc.previouslyAllocatedDevices, tc.remainingDevicesByHintMap)

			assert.Equal(subT, tc.expectedSelectedLength, len(selectedDevices))
			for _, device := range selectedDevices {
				assert.Contains(subT, tc.expectedIn, device.GetTopologyHintKey())
			}
		})
	}
}

// TestBinPackingNpuAllocator_Warboy tests NpuAllocator.Allocate for Warboy arch with single core strategy.
func TestBinPackingNpuAllocator_Warboy(t *testing.T) {
	mockSMIDevices := smi.GetStaticMockDevices(smi.ArchWarboy) // we have total 8 mock devices
	sut, _ := NewBinPackingNpuAllocator(mockSMIDevices)

	// assume we have warboy with single core strategy.
	// therefore, each smi device will have 2 cores, which will generate 2 mock devices by each iteration.
	// after this iteration, we will have total 16 mock devices.
	// nodeIdx of each mock devices will be [0, 0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7].
	mockDevices := make(DeviceSet, 0)
	for _, smiDevice := range mockSMIDevices {
		deviceInfo, _ := smiDevice.DeviceInfo()
		pciBusID, _ := parseBusIDFromBDF(deviceInfo.BDF())
		mockDevices = mockDevices.Union(generateSameBoardMockDeviceSet(0, 2, TopologyHintKey(pciBusID)))
	}

	tests := []struct {
		description      string
		available        DeviceSet
		required         DeviceSet
		request          int
		verificationFunc func(DeviceSet) error
	}{
		{
			description: "total 4 devices must be allocated with 2 hintKey groups",
			available:   mockDevices[:],
			required:    DeviceSet{},
			request:     4,
			verificationFunc: func(deviceSet DeviceSet) error {
				hintKeyCntMap := make(map[TopologyHintKey]int)
				for _, device := range deviceSet {
					hintKeyCntMap[device.GetTopologyHintKey()] += 1
				}

				if len(hintKeyCntMap) != 2 {
					return fmt.Errorf("expected 2 hintKeys, got %d", len(hintKeyCntMap))
				}

				return nil
			},
		},
		{
			description: "total 8 devices must be allocated with 4 hintKey groups",
			available:   mockDevices[:],
			required:    DeviceSet{},
			request:     8,
			verificationFunc: func(deviceSet DeviceSet) error {
				hintKeyCntMap := make(map[TopologyHintKey]int)
				for _, device := range deviceSet {
					hintKeyCntMap[device.GetTopologyHintKey()] += 1
				}

				if len(hintKeyCntMap) != 4 {
					return fmt.Errorf("expected 4 hintKeys, got %d", len(hintKeyCntMap))
				}

				return nil
			},
		},
		{
			description: "total 11 devices must be allocated with 6 hintKey groups",
			available:   mockDevices[:],
			required:    DeviceSet{},
			request:     11,
			verificationFunc: func(deviceSet DeviceSet) error {
				hintKeyCntMap := make(map[TopologyHintKey]int)
				for _, device := range deviceSet {
					hintKeyCntMap[device.GetTopologyHintKey()] += 1
				}

				if len(hintKeyCntMap) != 6 {
					return fmt.Errorf("expected 6 hintKeys, got %d", len(hintKeyCntMap))
				}

				return nil
			},
		},
		{
			description: "total 8 devices must be allocated with 2 new hintKey groups and 2 existing hintKey groups",
			available:   mockDevices[:],
			required: DeviceSet{
				mockDevices[0],
				mockDevices[1],
				mockDevices[2],
				mockDevices[3],
			},
			request: 8,
			verificationFunc: func(deviceSet DeviceSet) error {
				hintKeyCntMap := make(map[TopologyHintKey]int)
				for _, device := range deviceSet {
					hintKeyCntMap[device.GetTopologyHintKey()] += 1
				}

				for i := range []int{0, 1, 2, 3} {
					if _, ok := hintKeyCntMap[mockDevices[i].GetTopologyHintKey()]; !ok {
						return fmt.Errorf("expected to find hintKey %v in %v", mockDevices[i].GetTopologyHintKey(), hintKeyCntMap)
					}
				}

				if len(hintKeyCntMap) != 4 {
					return fmt.Errorf("expected 6 hintKeys, got %d", len(hintKeyCntMap))
				}

				return nil
			},
		},
		{
			description: "total 8 devices must be allocated with 4 existing hintKey groups",
			available:   mockDevices[:],
			required: DeviceSet{
				mockDevices[0],
				mockDevices[2],
				mockDevices[4],
				mockDevices[6],
			},
			request: 8,
			verificationFunc: func(deviceSet DeviceSet) error {
				hintKeyCntMap := make(map[TopologyHintKey]int)
				for _, device := range deviceSet {
					hintKeyCntMap[device.GetTopologyHintKey()] += 1
				}

				for i := range []int{0, 2, 4, 6} {
					if _, ok := hintKeyCntMap[mockDevices[i].GetTopologyHintKey()]; !ok {
						return fmt.Errorf("expected to find hintKey %v in %v", mockDevices[i].GetTopologyHintKey(), hintKeyCntMap)
					}
				}

				if len(hintKeyCntMap) != 4 {
					return fmt.Errorf("expected 6 hintKeys, got %d", len(hintKeyCntMap))
				}

				return nil
			},
		},
		{
			description: "total 8 devices must be allocated with 1 new hintKey groups and 3 existing hintKey groups",
			available:   mockDevices[:],
			required: DeviceSet{
				mockDevices[0],
				mockDevices[1],
				mockDevices[2],
				mockDevices[4],
			},
			request: 8,
			verificationFunc: func(deviceSet DeviceSet) error {
				hintKeyCntMap := make(map[TopologyHintKey]int)
				for _, device := range deviceSet {
					hintKeyCntMap[device.GetTopologyHintKey()] += 1
				}

				for i := range []int{0, 1, 2, 4} {
					if _, ok := hintKeyCntMap[mockDevices[i].GetTopologyHintKey()]; !ok {
						return fmt.Errorf("expected to find hintKey %v in %v", mockDevices[i].GetTopologyHintKey(), hintKeyCntMap)
					}
				}

				if len(hintKeyCntMap) != 4 {
					return fmt.Errorf("expected 6 hintKeys, got %d", len(hintKeyCntMap))
				}

				return nil
			},
		},
	}

	t.Parallel()
	for _, tc := range tests {
		t.Run(tc.description, func(subT *testing.T) {
			allocatedDevices := sut.Allocate(tc.available, tc.required, tc.request)
			assert.Equal(subT, tc.request, len(allocatedDevices))
			assert.NoError(subT, tc.verificationFunc(allocatedDevices))
		})
	}
}

// TestBinPackingNpuAllocator_RNGD tests NpuAllocator.Allocate for RNGD arch with single core strategy.
func TestBinPackingNpuAllocator_RNGD(t *testing.T) {
	// FIXME: due to bug exists in smi package, it is not possible to write tests for RNGD right now.
	//  Please fix it after https://github.com/furiosa-ai/libfuriosa-kubernetes/pull/45 PR is merged.
	//  ---
	//  panic: interface conversion: smi.Device is smi.staticRngdMockDevice, not *smi.staticRngdMockDevice [recovered]
	//        panic: interface conversion: smi.Device is smi.staticRngdMockDevice, not *smi.staticRngdMockDevice

	//mockSMIDevices := smi.GetStaticMockDevices(smi.ArchRngd) // we have total 8 mock devices
	//sut, _ := NewBinPackingNpuAllocator(mockSMIDevices)
	//
	//// assume we have warboy with single core strategy.
	//// therefore, each smi device will have 2 cores, which will generate 2 mock devices by each iteration.
	//// after this iteration, we will have total 16 mock devices.
	//// nodeIdx of each mock devices will be like below.
	//// [
	////   0, 0, 0, 0, 0, 0, 0, 0,
	////   1, 1, 1, 1, 1, 1, 1, 1,
	////   2, 2, 2, 2, 2, 2, 2, 2,
	////   3, 3, 3, 3, 3, 3, 3, 3,
	////   4, 4, 4, 4, 4, 4, 4, 4,
	////   5, 5, 5, 5, 5, 5, 5, 5,
	////   6, 6, 6, 6, 6, 6, 6, 6,
	////   7, 7, 7, 7, 7, 7, 7, 7,
	//// ].
	//mockDevices := make(DeviceSet, 0)
	//for _, smiDevice := range mockSMIDevices {
	//	deviceInfo, _ := smiDevice.DeviceInfo()
	//	pciBusID, _ := parseBusIDFromBDF(deviceInfo.BDF())
	//	mockDevices = mockDevices.Union(generateSameBoardMockDeviceSet(0, 8, TopologyHintKey(pciBusID)))
	//}
	//
	//tests := []struct {
	//	description      string
	//	available        DeviceSet
	//	required         DeviceSet
	//	request          int
	//	verificationFunc func(DeviceSet) error
	//}{}
	//
	//t.Parallel()
	//for _, tc := range tests {
	//	t.Run(tc.description, func(subT *testing.T) {
	//		allocatedDevices := sut.Allocate(tc.available, tc.required, tc.request)
	//		assert.Equal(subT, tc.request, len(allocatedDevices))
	//		assert.NoError(subT, tc.verificationFunc(allocatedDevices))
	//	})
	//}
}
