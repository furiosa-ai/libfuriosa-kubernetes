package npu_allocator

import (
	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"gonum.org/v1/gonum/stat/combin"
	"sort"
)

var _ NpuAllocator = (*binPackingNpuAllocator)(nil)

type binPackingNpuAllocator struct {
	topologyHintMatrix TopologyHintMatrix
}

func NewBinPackingNpuAllocator(devices []smi.Device) (NpuAllocator, error) {
	topologyHintMatrix, err := populateTopologyHintMatrixForScoreBasedAllocator(devices)
	if err != nil {
		return nil, err
	}

	return &binPackingNpuAllocator{
		topologyHintMatrix: topologyHintMatrix,
	}, nil
}

func (b binPackingNpuAllocator) Allocate(available DeviceSet, required DeviceSet, size int) DeviceSet {
	// Step1: available DeviceSet을 TopologyHintKey를 기준으로 map을 만든다.
	availableByTopologyHintKey := make(map[TopologyHintKey]DeviceSet)
	for _, device := range available {
		topologyKey := device.GetTopologyHintKey()
		availableByTopologyHintKey[topologyKey] = append(availableByTopologyHintKey[topologyKey], device)
	}

	// Step2: required DeviceSet을 먼저 처리한다, 같은 물리 카드에서 할당을 우선시하기위해 requiredKey를 수집한다.
	collected := DeviceSet{}
	requiredKeys := make(map[TopologyHintKey]bool)
	for _, device := range required {
		collected = append(collected, device)

		topologyHintKey := device.GetTopologyHintKey()
		requiredKeys[topologyHintKey] = true

		deviceSetForKey := availableByTopologyHintKey[topologyHintKey]
		availableByTopologyHintKey[topologyHintKey] = deviceSetForKey.Difference(DeviceSet{device})
	}

	if len(collected) == size {
		return collected
	}

	// Step3: required key를 먼저 처리하여 이미 할당이된 물리카드에서 먼저 할당을 한다.
	for topologyHintKey := range requiredKeys {
		devices := availableByTopologyHintKey[topologyHintKey]
		for _, device := range devices {
			collected = append(collected, device)
			availableByTopologyHintKey[topologyHintKey] = availableByTopologyHintKey[topologyHintKey].Difference(DeviceSet{device})

			if len(collected) == size {
				return collected
			}
		}
	}

	// Step4: required key와 사용하지 않은 key의 조합을 생성하여 가장 높은 점수를 가지는 key부터 할당 한다.
	var availableKeys []TopologyHintKey
	for key := range availableByTopologyHintKey {
		if !requiredKeys[key] {
			availableKeys = append(availableKeys, key)
		}
	}

	requiredKeysSlice := make([]TopologyHintKey, 0, len(requiredKeys))
	for key := range requiredKeys {
		requiredKeysSlice = append(requiredKeysSlice, key)
	}

	var combinations [][]TopologyHintKey
	for _, availableKey := range availableKeys {
		copiedRequiredKeys := append([]TopologyHintKey{}, requiredKeysSlice...)
		// availableKey는 slice의 맨 마지막이다.
		combination := append(copiedRequiredKeys, availableKey)
		combinations = append(combinations, combination)
	}

	//calculateCombinationScore가지고 채점
	var scores []uint
	for i, combination := range combinations {
		scores[i] = calculateCombinationScore(combination, b.topologyHintMatrix)
	}

	//점수 높음순으로 정렬
	var indices []int
	for i := range indices {
		indices[i] = i
	}

	sort.Slice(indices, func(i, j int) bool {
		return scores[indices[i]] > scores[indices[j]]
	})

	//최종적으로 topology 점수가 가장높은 available key가 정렬이 되있으며, key를 추출해서 해당 카드부터 순차적으로 할당
	for _, idx := range indices {
		combination := combinations[idx]
		availableKey := combination[len(combination)-1]
		devices := availableByTopologyHintKey[availableKey]
		for _, device := range devices {
			collected = append(collected, device)

			if len(collected) == size {
				return collected
			}
		}
	}

	return collected
}

func calculateCombinationScore(keys []TopologyHintKey, topologyHintMatrix TopologyHintMatrix) uint {
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
