package npu_allocator

import (
	"fmt"
	"testing"

	"github.com/bradfitz/iter"
	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func generateSameBoardMockDeviceSet(cnt int, hintKey TopologyHintKey) DeviceSet {
	devices := make(DeviceSet, 0, cnt)
	for range iter.N(cnt) {
		UUID, _ := uuid.NewUUID()
		devices = append(devices, NewMockDevice(UUID.String(), hintKey))
	}

	return devices
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
		pciBusID, _ := util.ParseBusIDFromBDF(deviceInfo.BDF())
		mockDevices = mockDevices.Union(generateSameBoardMockDeviceSet(2, TopologyHintKey(pciBusID)))
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

				for _, idx := range []int{0, 1, 2, 3} {
					expectedHintKey := mockDevices[idx].GetTopologyHintKey()
					if _, ok := hintKeyCntMap[expectedHintKey]; !ok {
						return fmt.Errorf("expected hintKey %s is not in the selected list %v", expectedHintKey, hintKeyCntMap)
					}
				}

				if len(hintKeyCntMap) != 4 {
					return fmt.Errorf("expected 4 hintKeys, got %d", len(hintKeyCntMap))
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

				for _, idx := range []int{0, 2, 4, 6} {
					expectedHintKey := mockDevices[idx].GetTopologyHintKey()
					if _, ok := hintKeyCntMap[expectedHintKey]; !ok {
						return fmt.Errorf("expected hintKey %s is not in the selected list %v", expectedHintKey, hintKeyCntMap)
					}
				}

				if len(hintKeyCntMap) != 4 {
					return fmt.Errorf("expected 4 hintKeys, got %d", len(hintKeyCntMap))
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

				for _, idx := range []int{0, 1, 2, 4} {
					expectedHintKey := mockDevices[idx].GetTopologyHintKey()
					if _, ok := hintKeyCntMap[expectedHintKey]; !ok {
						return fmt.Errorf("expected hintKey %s is not in the selected list %v", expectedHintKey, hintKeyCntMap)
					}
				}

				if len(hintKeyCntMap) != 4 {
					return fmt.Errorf("expected 4 hintKeys, got %d", len(hintKeyCntMap))
				}

				return nil
			},
		},
		{
			description: "total 8 devices must be allocated with 4 new hintKey groups at the same NUMA node",
			available: DeviceSet{
				mockDevices[0],
				mockDevices[2],
				mockDevices[3],
				mockDevices[4],
				mockDevices[5],
				mockDevices[6],
				mockDevices[7],
				mockDevices[8],
				mockDevices[9],
				mockDevices[10],
				mockDevices[11],
				mockDevices[12],
				mockDevices[13],
				mockDevices[14],
				mockDevices[15],
			},
			required: DeviceSet{},
			request:  8,
			verificationFunc: func(deviceSet DeviceSet) error {
				hintKeyCntMap := make(map[TopologyHintKey]int)
				for _, device := range deviceSet {
					hintKeyCntMap[device.GetTopologyHintKey()] += 1
				}

				for _, idx := range []int{8, 9, 10, 11, 12, 13, 14, 15} {
					expectedHintKey := mockDevices[idx].GetTopologyHintKey()
					if _, ok := hintKeyCntMap[expectedHintKey]; !ok {
						return fmt.Errorf("expected hintKey %s is not in the selected list %v", expectedHintKey, hintKeyCntMap)
					}
				}

				if len(hintKeyCntMap) != 4 {
					return fmt.Errorf("expected 4 hintKeys, got %d in %v", len(hintKeyCntMap), hintKeyCntMap)
				}

				return nil
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(subT *testing.T) {
			allocatedDevices := sut.Allocate(tc.available, tc.required, tc.request)
			for _, device := range allocatedDevices {
				fmt.Printf("%s ", device.GetTopologyHintKey())
			}
			fmt.Println()

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
