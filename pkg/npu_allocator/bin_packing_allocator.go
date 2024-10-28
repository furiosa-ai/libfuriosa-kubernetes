package npu_allocator

import (
	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"gonum.org/v1/gonum/stat/combin"
)

var _ NpuAllocator = (*binPackingNpuAllocator)(nil)

type binPackingNpuAllocator struct {
	topologyScoreCalculator TopologyScoreCalculator
}

func NewBinPackingNpuAllocator(devices []smi.Device) (NpuAllocator, error) {
	topologyHintMatrix, err := NewTopologyHintMatrix(devices)
	if err != nil {
		return nil, err
	}

	return newBinPackingNpuAllocator(generateTopologyScoreCalculator(topologyHintMatrix)), nil
}

func NewMockBinPackingNpuAllocator(topologyHintMatrix TopologyHintMatrix) (NpuAllocator, error) {
	return newBinPackingNpuAllocator(generateTopologyScoreCalculator(topologyHintMatrix)), nil
}

// generateTopologyScoreCalculator returns calculator that calculates total sum of given TopologyHintKey list.
func generateTopologyScoreCalculator(topologyHintMatrix TopologyHintMatrix) TopologyScoreCalculator {
	return func(keys []TopologyHintKey) uint {
		// If there is only one key in keys, scoring combinations has no meaning.
		// This also prevents a panic from combin.Combinations when n is less than k.
		if len(keys) == 1 {
			return 0
		}

		totalScore := uint(0)

		indices := len(keys)
		combinations := combin.Combinations(indices, 2)

		for _, keyPair := range combinations {
			i, j := keyPair[0], keyPair[1]
			key1, key2 := keys[i], keys[j]
			if key1 > key2 {
				key1, key2 = key2, key1
			}

			if innerMap, exists := topologyHintMatrix[key1]; exists {
				if score, scoreExists := innerMap[key2]; scoreExists {
					totalScore += score
				}
			}
		}

		return totalScore
	}
}

func newBinPackingNpuAllocator(topologyScoreCalculator TopologyScoreCalculator) NpuAllocator {
	return &binPackingNpuAllocator{topologyScoreCalculator: topologyScoreCalculator}
}

func (b *binPackingNpuAllocator) Allocate(available DeviceSet, required DeviceSet, size int) DeviceSet {
	// If length of `required` already satisfies given `size`, just return it.
	if len(required) == size {
		return required
	}

	// Step 1: available DeviceSet 을 TopologyHintKey 를 기준으로 map 을 만든다.
	availableDevicesByHintKeyMap := make(map[TopologyHintKey]DeviceSet)
	for _, device := range available {
		hintKey := device.TopologyHintKey()
		if _, ok := availableDevicesByHintKeyMap[hintKey]; !ok {
			availableDevicesByHintKeyMap[hintKey] = make(DeviceSet, 0)
		}

		availableDevicesByHintKeyMap[hintKey] = append(availableDevicesByHintKeyMap[hintKey], device)
	}

	// Step 2: Process the required DeviceSet first. Collect required keys to prioritize allocations from the same physical card.
	collectedDevices := make(DeviceSet, 0, size)
	requiredHintKeySet := make(map[TopologyHintKey]struct{})
	for _, device := range required {
		collectedDevices = append(collectedDevices, device)

		hintKey := device.TopologyHintKey()
		requiredHintKeySet[hintKey] = struct{}{}

		availableDevicesByHintKeyMap[hintKey] = availableDevicesByHintKeyMap[hintKey].Difference(DeviceSet{device})
	}

	if len(collectedDevices) == size {
		return collectedDevices
	}

	// Step 3: Consume required keys first to mitigate fragmentation.
	for hintKey := range requiredHintKeySet {
		devices := availableDevicesByHintKeyMap[hintKey]
		for _, device := range devices {
			collectedDevices = append(collectedDevices, device)
			availableDevicesByHintKeyMap[hintKey] = availableDevicesByHintKeyMap[hintKey].Difference(DeviceSet{device})

			if len(collectedDevices) == size {
				return collectedDevices
			}
		}
	}

	// Step 4: Calculate device count to be allocated.
	remainingDevicesSize := size - len(collectedDevices)

	unusedHintKeys := make([]TopologyHintKey, 0)
	deviceCountByHintKeyMap := make(map[TopologyHintKey]int)
	for hintKey, devices := range availableDevicesByHintKeyMap {
		if _, ok := requiredHintKeySet[hintKey]; !ok {
			unusedHintKeys = append(unusedHintKeys, hintKey)
		}

		deviceCountByHintKeyMap[hintKey] = len(devices)
	}

	// Step 5: Generate combinations using unused hint keys, with size ranging form up to the number of unused hint keys.
	validCombinationsOfHintKeys := generateValidHintKeysCombinations(unusedHintKeys, deviceCountByHintKeyMap, remainingDevicesSize)

	// Step 6: If required keys exists, add them to all combinations to ensure correct scoring.
	requiredHintKeys := make([]TopologyHintKey, 0, len(requiredHintKeySet))
	for hintKey := range requiredHintKeySet {
		requiredHintKeys = append(requiredHintKeys, hintKey)
	}

	for i := range validCombinationsOfHintKeys {
		validCombinationsOfHintKeys[i] = append(validCombinationsOfHintKeys[i], requiredHintKeys...)
	}

	// Step 7: Score each combination and find the one with the highest score.
	var highestScore *uint = nil
	var bestHintKeys []TopologyHintKey
	for _, hintKeys := range validCombinationsOfHintKeys {
		score := b.topologyScoreCalculator(hintKeys)

		if highestScore == nil || score > *highestScore {
			highestScore = &score
			bestHintKeys = hintKeys
		}
	}

	// Step 8: Add to collectedDevices and return.
BestHintKeysLoop:
	for _, hintKey := range bestHintKeys {
		devices := availableDevicesByHintKeyMap[hintKey]
		for _, device := range devices {
			collectedDevices = append(collectedDevices, device)
			if len(collectedDevices) == size {
				break BestHintKeysLoop
			}
		}
	}

	return collectedDevices
}

func generateValidHintKeysCombinations(unusedHintKeys []TopologyHintKey, deviceCountByHintKeyMap map[TopologyHintKey]int, remainingDevicesSize int) [][]TopologyHintKey {
	// Given keys like 1, 2, 3, and 4, generate combinations as follows:
	// (1), (2), (3), (4)
	// (1, 2), (1, 3), (1, 4), (2, 3), (2, 4), (3, 4)
	// (1, 2, 3), (1, 2, 4), (1, 3, 4), (2, 3, 4)
	// (1, 2, 3, 4)
	// Start with smaller sets and add combinations to the result as soon as they satisfy the required size.
	// All sets in validCombinationsOfHintKeys will be of equal size.

	validCombinationsOfHintKeys := make([][]TopologyHintKey, 0)
	for k := 1; k <= len(unusedHintKeys); k++ {
		indicesCombinations := combin.Combinations(len(unusedHintKeys), k)
		for _, indices := range indicesCombinations {
			hintKeys := make([]TopologyHintKey, 0)
			totalDevices := 0

			for _, idx := range indices {
				hintKey := unusedHintKeys[idx]
				hintKeys = append(hintKeys, hintKey)

				totalDevices += deviceCountByHintKeyMap[hintKey]
			}

			if totalDevices >= remainingDevicesSize {
				validCombinationsOfHintKeys = append(validCombinationsOfHintKeys, hintKeys)
			}
		}

		if len(validCombinationsOfHintKeys) > 0 {
			break
		}
	}

	return validCombinationsOfHintKeys
}
