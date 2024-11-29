package npu_allocator

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"github.com/stretchr/testify/assert"
)

// TestValidHintKeysCombinationsGeneratorWith8HintKeys tests `generateValidHintKeysCombinations(...)` generates correct results.
func TestValidHintKeysCombinationsGeneratorWith8HintKeys(t *testing.T) {
	hintKeys := []TopologyHintKey{"0", "1", "2", "3", "4", "5", "6", "7"}

	tests := []struct {
		description             string
		unusedHintKeys          []TopologyHintKey
		deviceCountByHintKeyMap map[TopologyHintKey]int
		remainingDevicesSize    int
		expected                [][]TopologyHintKey
	}{
		{
			description:    "No hintKeys are in use, every devices in each hintKeys are available with size 8, remaining 8",
			unusedHintKeys: hintKeys[:],
			deviceCountByHintKeyMap: map[TopologyHintKey]int{
				hintKeys[0]: 8,
				hintKeys[1]: 8,
				hintKeys[2]: 8,
				hintKeys[3]: 8,
				hintKeys[4]: 8,
				hintKeys[5]: 8,
				hintKeys[6]: 8,
				hintKeys[7]: 8,
			},
			remainingDevicesSize: 8,
			expected: [][]TopologyHintKey{
				{hintKeys[0]},
				{hintKeys[1]},
				{hintKeys[2]},
				{hintKeys[3]},
				{hintKeys[4]},
				{hintKeys[5]},
				{hintKeys[6]},
				{hintKeys[7]},
			},
		},
		{
			description:    "No hintKeys are in use, every devices in each hintKeys are available with size 8, remaining 10",
			unusedHintKeys: hintKeys[:],
			deviceCountByHintKeyMap: map[TopologyHintKey]int{
				hintKeys[0]: 8,
				hintKeys[1]: 8,
				hintKeys[2]: 8,
				hintKeys[3]: 8,
				hintKeys[4]: 8,
				hintKeys[5]: 8,
				hintKeys[6]: 8,
				hintKeys[7]: 8,
			},
			remainingDevicesSize: 10,
			expected: func() [][]TopologyHintKey {
				expectedHintKeys := make([][]TopologyHintKey, 0)
				for i := 0; i < 8; i++ {
					for j := i + 1; j < 8; j++ {
						expectedHintKeys = append(expectedHintKeys, []TopologyHintKey{hintKeys[i], hintKeys[j]})
					}
				}

				return expectedHintKeys
			}(),
		},
		{
			description:    "No hintKeys are in use, each hintKeys are partially available from 1 to 4, remaining 9",
			unusedHintKeys: hintKeys[:],
			deviceCountByHintKeyMap: map[TopologyHintKey]int{
				hintKeys[0]: 1,
				hintKeys[1]: 1,
				hintKeys[2]: 2,
				hintKeys[3]: 2,
				hintKeys[4]: 3,
				hintKeys[5]: 3,
				hintKeys[6]: 4,
				hintKeys[7]: 4,
			},
			remainingDevicesSize: 9,
			expected: [][]TopologyHintKey{
				{hintKeys[2], hintKeys[4], hintKeys[6]},
				{hintKeys[3], hintKeys[4], hintKeys[6]},
				{hintKeys[2], hintKeys[5], hintKeys[6]},
				{hintKeys[3], hintKeys[5], hintKeys[6]},
				{hintKeys[4], hintKeys[5], hintKeys[6]},
				{hintKeys[2], hintKeys[4], hintKeys[7]},
				{hintKeys[3], hintKeys[4], hintKeys[7]},
				{hintKeys[2], hintKeys[5], hintKeys[7]},
				{hintKeys[3], hintKeys[5], hintKeys[7]},
				{hintKeys[4], hintKeys[5], hintKeys[7]},
				{hintKeys[0], hintKeys[6], hintKeys[7]},
				{hintKeys[1], hintKeys[6], hintKeys[7]},
				{hintKeys[2], hintKeys[6], hintKeys[7]},
				{hintKeys[3], hintKeys[6], hintKeys[7]},
				{hintKeys[4], hintKeys[6], hintKeys[7]},
				{hintKeys[5], hintKeys[6], hintKeys[7]},
			},
		},
		{
			description:    "No hintKeys are in use, each hintKeys are partially available, remaining 5, do not pick hintKey with deviceCount 0",
			unusedHintKeys: hintKeys[:],
			deviceCountByHintKeyMap: map[TopologyHintKey]int{
				hintKeys[0]: 0,
				hintKeys[1]: 0,
				hintKeys[2]: 1,
				hintKeys[3]: 1,
				hintKeys[4]: 1,
				hintKeys[5]: 2,
				hintKeys[6]: 2,
				hintKeys[7]: 2,
			},
			remainingDevicesSize: 5,
			expected: [][]TopologyHintKey{
				{hintKeys[2], hintKeys[5], hintKeys[6]},
				{hintKeys[3], hintKeys[5], hintKeys[6]},
				{hintKeys[4], hintKeys[5], hintKeys[6]},
				{hintKeys[2], hintKeys[5], hintKeys[7]},
				{hintKeys[3], hintKeys[5], hintKeys[7]},
				{hintKeys[4], hintKeys[5], hintKeys[7]},
				{hintKeys[2], hintKeys[6], hintKeys[7]},
				{hintKeys[3], hintKeys[6], hintKeys[7]},
				{hintKeys[4], hintKeys[6], hintKeys[7]},
				{hintKeys[5], hintKeys[6], hintKeys[7]},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			actual := generateValidHintKeysCombinations(tc.unusedHintKeys, tc.deviceCountByHintKeyMap, tc.remainingDevicesSize)

			assert.Equal(t, len(tc.expected), len(actual))
			assert.ElementsMatch(t, tc.expected, actual)
		})
	}
}

// TestBinPackingNpuAllocator tests NpuAllocator.Allocate
func TestBinPackingNpuAllocator(t *testing.T) {
	staticHintMatrix := buildStaticHintMatrixForTwoSocketBalancedConfig()

	sut, _ := NewMockBinPackingNpuAllocator(staticHintMatrix)
	generateMockDevices := func(devicesPerBoard int) DeviceSet {
		mockDevices := make(DeviceSet, 0)
		for idx, hintKey := range getStaticHintKeys() {
			mockDevices = mockDevices.Union(generateSameBoardMockDeviceSet(idx, devicesPerBoard, hintKey))
		}

		return mockDevices
	}

	t.Run("2 devices share single board", func(t *testing.T) {
		// nodeIdx of each mock devices will be:
		// [
		//	  0, 0,
		//	  1, 1,
		//	  2, 2,
		//	  3, 3,
		//	  4, 4,
		//	  5, 5,
		//	  6, 6,
		//	  7, 7,
		// ]
		mockDevices := generateMockDevices(2)

		tests := []struct {
			description      string
			available        DeviceSet
			required         DeviceSet
			request          int
			verificationFunc func(DeviceSet) error
		}{
			{
				description: "all devices are available, no required devices exist, 4 requested",
				available:   mockDevices[:],
				required:    DeviceSet{},
				request:     4,
				verificationFunc: func(deviceSet DeviceSet) error {
					hintKeyCntMap := make(map[TopologyHintKey]int)
					for _, device := range deviceSet {
						hintKeyCntMap[device.TopologyHintKey()] += 1
					}

					if len(hintKeyCntMap) != 2 {
						return fmt.Errorf("expected 2 hintKeys, got %d", len(hintKeyCntMap))
					}

					if len(deviceSet) != 4 {
						return fmt.Errorf("expected 4 devices, got %d", len(deviceSet))
					}

					return nil
				},
			},
			{
				description: "all devices are available, no required devices exist, 7 requested",
				available:   mockDevices[:],
				required:    DeviceSet{},
				request:     7,
				verificationFunc: func(deviceSet DeviceSet) error {
					hintKeyCntMap := make(map[TopologyHintKey]int)
					for _, device := range deviceSet {
						hintKeyCntMap[device.TopologyHintKey()] += 1
					}

					if len(hintKeyCntMap) != 4 {
						return fmt.Errorf("expected 4 hintKeys, got %d", len(hintKeyCntMap))
					}

					if len(deviceSet) != 7 {
						return fmt.Errorf("expected 7 devices, got %d", len(deviceSet))
					}

					return nil
				},
			},
			{
				description: "15 devices are available, no required devices exist, 8 requested",
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
						hintKeyCntMap[device.TopologyHintKey()] += 1
					}

					for _, idx := range []int{8, 9, 10, 11, 12, 13, 14, 15} {
						expectedHintKey := mockDevices[idx].TopologyHintKey()
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
			{
				description: "all devices are available, 4 required devices exist, 8 requested",
				available:   mockDevices[:],
				required: DeviceSet{
					mockDevices[0],
					mockDevices[1],
					mockDevices[2],
					mockDevices[4],
				},
				request: 7,
				verificationFunc: func(deviceSet DeviceSet) error {
					hintKeyCntMap := make(map[TopologyHintKey]int)
					for _, device := range deviceSet {
						hintKeyCntMap[device.TopologyHintKey()] += 1
					}

					for _, idx := range []int{0, 1, 2, 4} {
						expectedHintKey := mockDevices[idx].TopologyHintKey()
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
				description: "all devices are available, 4 required devices exist, 4 requested",
				available:   mockDevices[:],
				required: DeviceSet{
					mockDevices[4],
					mockDevices[5],
					mockDevices[6],
					mockDevices[7],
				},
				request: 4,
				verificationFunc: func(deviceSet DeviceSet) error {
					idCheckMap := make(map[string]struct{})
					for _, device := range deviceSet {
						idCheckMap[device.ID()] = struct{}{}
					}

					// because length of required devices and value of request are same,
					// selected devices must be same as given required devices.
					for _, idx := range []int{4, 5, 6, 7} {
						expectedDevice := mockDevices[idx]
						if _, ok := idCheckMap[expectedDevice.ID()]; !ok {
							return fmt.Errorf("expected device %s is not in the selected list %v", expectedDevice, idCheckMap)
						}
					}

					return nil
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.description, func(t *testing.T) {
				allocatedDevices := sut.Allocate(tc.available, tc.required, tc.request)

				assert.Equal(t, tc.request, len(allocatedDevices))
				assert.NoError(t, tc.verificationFunc(allocatedDevices))
			})
		}
	})

	t.Run("4 devices share single board", func(t *testing.T) {
		// nodeIdx of each mock devices will be:
		// [
		//	  0, 0, 0, 0,
		//	  1, 1, 1, 1,
		//	  2, 2, 2, 2,
		//	  3, 3, 3, 3,
		//	  4, 4, 4, 4,
		//	  5, 5, 5, 5,
		//	  6, 6, 6, 6,
		//	  7, 7, 7, 7,
		// ]
		mockDevices := generateMockDevices(4)

		tests := []struct {
			description      string
			available        DeviceSet
			required         DeviceSet
			request          int
			verificationFunc func(DeviceSet) error
		}{
			{
				description: "all devices are available, no required devices exist, 4 requested",
				available:   mockDevices[:],
				required:    DeviceSet{},
				request:     4,
				verificationFunc: func(deviceSet DeviceSet) error {
					hintKeyCntMap := make(map[TopologyHintKey]int)
					for _, device := range deviceSet {
						hintKeyCntMap[device.TopologyHintKey()] += 1
					}

					if len(hintKeyCntMap) != 1 {
						return fmt.Errorf("expected 1 hintKeys, got %d", len(hintKeyCntMap))
					}

					if len(deviceSet) != 4 {
						return fmt.Errorf("expected 4 devices, got %d", len(deviceSet))
					}

					return nil
				},
			},
			{
				description: "all devices are available, no required devices exist, 7 requested",
				available:   mockDevices[:],
				required:    DeviceSet{},
				request:     7,
				verificationFunc: func(deviceSet DeviceSet) error {
					hintKeyCntMap := make(map[TopologyHintKey]int)
					for _, device := range deviceSet {
						hintKeyCntMap[device.TopologyHintKey()] += 1
					}

					if len(hintKeyCntMap) != 2 {
						return fmt.Errorf("expected 4 hintKeys, got %d", len(hintKeyCntMap))
					}

					if len(deviceSet) != 7 {
						return fmt.Errorf("expected 7 devices, got %d", len(deviceSet))
					}

					return nil
				},
			},
			{
				description: "all devices are available, no required devices exist, 8 requested",
				available:   mockDevices[:],
				required:    DeviceSet{},
				request:     8,
				verificationFunc: func(deviceSet DeviceSet) error {
					hintKeyCntMap := make(map[TopologyHintKey]int)
					for _, device := range deviceSet {
						hintKeyCntMap[device.TopologyHintKey()] += 1
					}

					if len(hintKeyCntMap) != 2 {
						return fmt.Errorf("expected 1 hintKeys, got %d", len(hintKeyCntMap))
					}

					hintKeys := make([]TopologyHintKey, 0)
					for hintKey := range hintKeyCntMap {
						hintKeys = append(hintKeys, hintKey)
					}

					scoreSum := sut.(*binPackingNpuAllocator).topologyScoreCalculator(hintKeys)
					if scoreSum != uint(smi.LinkTypeHostBridge) {
						return fmt.Errorf("2 hintKeys must be picked as `LinkTypeHostBridge` but it is not %d", scoreSum)
					}

					return nil
				},
			},
			{
				description: "all devices are available, no required devices exist, 10 requested",
				available:   mockDevices[:],
				required:    DeviceSet{},
				request:     10,
				verificationFunc: func(deviceSet DeviceSet) error {
					hintKeyCntMap := make(map[TopologyHintKey]int)
					for _, device := range deviceSet {
						hintKeyCntMap[device.TopologyHintKey()] += 1
					}

					if len(hintKeyCntMap) != 3 {
						return fmt.Errorf("expected 4 hintKeys, got %d", len(hintKeyCntMap))
					}

					hintKeys := make([]TopologyHintKey, 0)
					for hintKey := range hintKeyCntMap {
						hintKeys = append(hintKeys, hintKey)
					}

					scoreSum := sut.(*binPackingNpuAllocator).topologyScoreCalculator(hintKeys)
					if scoreSum != (uint(smi.LinkTypeHostBridge) + uint(smi.LinkTypeCpu)*2) {
						return fmt.Errorf("2 `LinkTypeCpu` and 1 `LinkTypeHostBridge` expected, but it is not: %d", scoreSum)
					}

					return nil
				},
			},
			{
				description: "21 devices are available, 2 required devices exist, 6 requested",
				// devices in the same line means they share same topologyHintKey.
				available: DeviceSet{
					mockDevices[0], mockDevices[1], mockDevices[3], // 3ea
					mockDevices[4], mockDevices[5], mockDevices[6], mockDevices[7], // 4ea
					mockDevices[8],                   // 1ea
					mockDevices[12], mockDevices[13], // 2ea
					mockDevices[16], mockDevices[17], mockDevices[18], mockDevices[19], // 4ea
					mockDevices[20], mockDevices[21], mockDevices[22], // 3ea
					// 0ea
					mockDevices[28], mockDevices[29], mockDevices[30], mockDevices[31], // 4ea
				},
				required: DeviceSet{
					mockDevices[8],
					mockDevices[28],
				},
				request: 6,
				verificationFunc: func(deviceSet DeviceSet) error {
					hintKeyCntMap := make(map[TopologyHintKey]int)
					for _, device := range deviceSet {
						hintKeyCntMap[device.TopologyHintKey()] += 1
					}

					if len(hintKeyCntMap) != 3 {
						return fmt.Errorf("expected 3 hintKeys, got %d", len(hintKeyCntMap))
					}

					if _, ok := hintKeyCntMap[mockDevices[8].TopologyHintKey()]; !ok {
						return fmt.Errorf("hintKey for mockDevices[8] not found in selected hintKeys")
					}

					if _, ok := hintKeyCntMap[mockDevices[28].TopologyHintKey()]; !ok {
						return fmt.Errorf("hintKey for mockDevices[28] not found in selected hintKeys")
					}

					hintKeys := make([]TopologyHintKey, 0)
					for hintKey := range hintKeyCntMap {
						hintKeys = append(hintKeys, hintKey)
					}

					scoreSum := sut.(*binPackingNpuAllocator).topologyScoreCalculator(hintKeys)
					if scoreSum != (uint(smi.LinkTypeHostBridge) + uint(smi.LinkTypeInterconnect)*2) {
						return fmt.Errorf("2 `LinkTypeInterconnect` and 1 `LinkTypeHostBridge` expected, but it is not: %d", scoreSum)
					}

					return nil
				},
			},
			{
				description: "all devices are available, 4 required devices exist, 4 requested",
				available:   mockDevices[:],
				required: DeviceSet{
					mockDevices[8],
					mockDevices[15],
					mockDevices[22],
					mockDevices[27],
				},
				request: 4,
				verificationFunc: func(deviceSet DeviceSet) error {
					idCheckMap := make(map[string]struct{})
					for _, device := range deviceSet {
						idCheckMap[device.ID()] = struct{}{}
					}

					// because length of required devices and value of request are same,
					// selected devices must be same as given required devices.
					for _, idx := range []int{8, 15, 22, 27} {
						expectedDevice := mockDevices[idx]
						if _, ok := idCheckMap[expectedDevice.ID()]; !ok {
							return fmt.Errorf("expected device %s is not in the selected list %v", expectedDevice, idCheckMap)
						}
					}

					return nil
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.description, func(t *testing.T) {
				allocatedDevices := sut.Allocate(tc.available, tc.required, tc.request)

				assert.Equal(t, tc.request, len(allocatedDevices))
				assert.NoError(t, tc.verificationFunc(allocatedDevices))
			})
		}
	})

	t.Run("8 devices share single board", func(t *testing.T) {
		// nodeIdx of each mock devices will be:
		// [
		//	  0, 0, 0, 0, 0, 0, 0, 0,
		//	  1, 1, 1, 1, 1, 1, 1, 1,
		//	  2, 2, 2, 2, 2, 2, 2, 2,
		//	  3, 3, 3, 3, 3, 3, 3, 3,
		//	  4, 4, 4, 4, 4, 4, 4, 4,
		//	  5, 5, 5, 5, 5, 5, 5, 5,
		//	  6, 6, 6, 6, 6, 6, 6, 6,
		//	  7, 7, 7, 7, 7, 7, 7, 7,
		// ]
		mockDevices := generateMockDevices(8)

		tests := []struct {
			description      string
			available        DeviceSet
			required         DeviceSet
			request          int
			verificationFunc func(DeviceSet) error
		}{
			{
				description: "all devices are available, no required devices exist, 8 requested",
				available:   mockDevices[:],
				required:    DeviceSet{},
				request:     8,
				verificationFunc: func(deviceSet DeviceSet) error {
					hintKeyCntMap := make(map[TopologyHintKey]int)
					for _, device := range deviceSet {
						hintKeyCntMap[device.TopologyHintKey()] += 1
					}

					if len(hintKeyCntMap) != 1 {
						return fmt.Errorf("expected 1 hintKeys, got %d", len(hintKeyCntMap))
					}

					if len(deviceSet) != 8 {
						return fmt.Errorf("expected 8 devices, got %d", len(deviceSet))
					}

					return nil
				},
			},
			{
				description: "all devices are available, no required devices exist, 16 requested",
				available:   mockDevices[:],
				required:    DeviceSet{},
				request:     16,
				verificationFunc: func(deviceSet DeviceSet) error {
					hintKeyCntMap := make(map[TopologyHintKey]int)
					for _, device := range deviceSet {
						hintKeyCntMap[device.TopologyHintKey()] += 1
					}

					if len(hintKeyCntMap) != 2 {
						return fmt.Errorf("expected 1 hintKeys, got %d", len(hintKeyCntMap))
					}

					hintKeys := make([]TopologyHintKey, 0)
					for hintKey := range hintKeyCntMap {
						hintKeys = append(hintKeys, hintKey)
					}

					scoreSum := sut.(*binPackingNpuAllocator).topologyScoreCalculator(hintKeys)
					if scoreSum != uint(smi.LinkTypeHostBridge) {
						return fmt.Errorf("2 hintKeys must be picked as `LinkTypeHostBridge` but it is not %d", scoreSum)
					}

					return nil
				},
			},
			{
				description: "all devices are available, no required devices exist, 20 requested",
				available:   mockDevices[:],
				required:    DeviceSet{},
				request:     20,
				verificationFunc: func(deviceSet DeviceSet) error {
					hintKeyCntMap := make(map[TopologyHintKey]int)
					for _, device := range deviceSet {
						hintKeyCntMap[device.TopologyHintKey()] += 1
					}

					if len(hintKeyCntMap) != 3 {
						return fmt.Errorf("expected 4 hintKeys, got %d", len(hintKeyCntMap))
					}

					hintKeys := make([]TopologyHintKey, 0)
					for hintKey := range hintKeyCntMap {
						hintKeys = append(hintKeys, hintKey)
					}

					scoreSum := sut.(*binPackingNpuAllocator).topologyScoreCalculator(hintKeys)
					if scoreSum != (uint(smi.LinkTypeHostBridge) + uint(smi.LinkTypeCpu)*2) {
						return fmt.Errorf("2 `LinkTypeCpu` and 1 `LinkTypeHostBridge` expected, but it is not: %d", scoreSum)
					}

					return nil
				},
			},
			{
				description: "58 devices are available, 4 required devices exist, 20 requested",
				available: DeviceSet{
					mockDevices[0], mockDevices[1], // 2ea
					mockDevices[8], mockDevices[9], mockDevices[10], mockDevices[11], mockDevices[12], mockDevices[13], mockDevices[14], mockDevices[15], // 8ea
					mockDevices[16], mockDevices[17], mockDevices[18], mockDevices[19], mockDevices[20], mockDevices[21], mockDevices[22], mockDevices[23], // 8ea
					mockDevices[24], mockDevices[25], mockDevices[26], mockDevices[27], mockDevices[28], mockDevices[29], mockDevices[30], mockDevices[31], // 8ea
					mockDevices[32], mockDevices[33], mockDevices[34], mockDevices[35], mockDevices[36], mockDevices[37], mockDevices[38], mockDevices[39], // 8ea
					mockDevices[40], mockDevices[41], mockDevices[42], mockDevices[43], mockDevices[44], mockDevices[45], mockDevices[46], mockDevices[47], // 8ea
					mockDevices[48], mockDevices[49], mockDevices[50], mockDevices[51], mockDevices[52], mockDevices[53], mockDevices[54], mockDevices[55], // 8ea
					mockDevices[56], mockDevices[57], mockDevices[58], mockDevices[59], mockDevices[60], mockDevices[61], mockDevices[62], mockDevices[63], // 8ea
				},
				required: DeviceSet{
					mockDevices[8], mockDevices[9], mockDevices[10], mockDevices[11],
				},
				request: 20,
				verificationFunc: func(deviceSet DeviceSet) error {
					hintKeyCntMap := make(map[TopologyHintKey]int)
					for _, device := range deviceSet {
						hintKeyCntMap[device.TopologyHintKey()] += 1
					}

					if len(hintKeyCntMap) != 3 {
						return fmt.Errorf("expected 3 hintKeys, got %d", len(hintKeyCntMap))
					}

					if _, ok := hintKeyCntMap[mockDevices[8].TopologyHintKey()]; !ok {
						return fmt.Errorf("hintKey for mockDevices[8] not found in selected hintKeys")
					}

					if _, ok := hintKeyCntMap[mockDevices[16].TopologyHintKey()]; !ok {
						return fmt.Errorf("hintKey for mockDevices[8] not found in selected hintKeys")
					}

					if _, ok := hintKeyCntMap[mockDevices[24].TopologyHintKey()]; !ok {
						return fmt.Errorf("hintKey for mockDevices[8] not found in selected hintKeys")
					}

					if hintKeyCntMap[mockDevices[8].TopologyHintKey()] != 8 {
						return fmt.Errorf("total 8 devices must be selected from hintKey %s, but it is not: %v", mockDevices[8].TopologyHintKey(), hintKeyCntMap)
					}

					if (hintKeyCntMap[mockDevices[16].TopologyHintKey()] + hintKeyCntMap[mockDevices[24].TopologyHintKey()]) != 12 {
						return fmt.Errorf("total 12 devices must be selected from hintKeys %s and %s, but it is not: %v", mockDevices[16].TopologyHintKey(), mockDevices[24].TopologyHintKey(), hintKeyCntMap)
					}

					return nil
				},
			},
			{
				description: "all devices are available, 8 required devices exist, 8 requested",
				available:   mockDevices[:],
				required: DeviceSet{
					mockDevices[1],
					mockDevices[4],
					mockDevices[6],
					mockDevices[11],
					mockDevices[15],
					mockDevices[16],
					mockDevices[22],
					mockDevices[27],
				},
				request: 8,
				verificationFunc: func(deviceSet DeviceSet) error {
					idCheckMap := make(map[string]struct{})
					for _, device := range deviceSet {
						idCheckMap[device.ID()] = struct{}{}
					}

					// because length of required devices and value of request are same,
					// selected devices must be same as given required devices.
					for _, idx := range []int{1, 4, 6, 11, 15, 16, 22, 27} {
						expectedDevice := mockDevices[idx]
						if _, ok := idCheckMap[expectedDevice.ID()]; !ok {
							return fmt.Errorf("expected device %s is not in the selected list %v", expectedDevice, idCheckMap)
						}
					}

					return nil
				},
			},
			{
				description: "single board is available, no required devices, 4 requested",
				available: func() DeviceSet {
					// Simulates devices are given with non-ordered through `Allocate(...)` call.
					partialMockDevices := make(DeviceSet, 8)
					copy(partialMockDevices, mockDevices[0:8])

					rand.Shuffle(len(partialMockDevices), func(i, j int) {
						partialMockDevices[i], partialMockDevices[j] = partialMockDevices[j], partialMockDevices[i]
					})

					return partialMockDevices
				}(),
				required: DeviceSet{},
				request:  4,
				verificationFunc: func(deviceSet DeviceSet) error {
					idCheckMap := make(map[string]struct{})
					for _, device := range deviceSet {
						idCheckMap[device.ID()] = struct{}{}
					}

					// Must always return first 4 PE cores.
					for _, idx := range []int{0, 1, 2, 3} {
						expectedDevice := mockDevices[idx]
						if _, ok := idCheckMap[expectedDevice.ID()]; !ok {
							return fmt.Errorf("expected device %s is not in the selected list %v", expectedDevice, idCheckMap)
						}
					}

					return nil
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.description, func(t *testing.T) {
				allocatedDevices := sut.Allocate(tc.available, tc.required, tc.request)

				assert.Equal(t, tc.request, len(allocatedDevices))
				assert.NoError(t, tc.verificationFunc(allocatedDevices))
			})
		}
	})
}
