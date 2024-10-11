package npu_allocator

import (
	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"gonum.org/v1/gonum/stat/combin"
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

	// step4: 남은 부분을 계산한다.
	devicesToAllocate := size - len(collected)

	unusedKeys := []TopologyHintKey{}
	for key := range availableByTopologyHintKey {
		if !requiredKeys[key] {
			unusedKeys = append(unusedKeys, key)
		}
	}

	keyDeviceCounts := make(map[TopologyHintKey]int)
	for key, devices := range availableByTopologyHintKey {
		keyDeviceCounts[key] = len(devices)
	}

	//step5: 1개의 key만 요소로 가지는 조합을 생성해서 최대 unusedKeys의 길이많큼 요소를 가지는 조합을 생성한다.
	// 1, 2, 3, 4란 key가 있다고 하면 아래와 같이 조합을 생성한다.
	// (1), (2), (3), (4)
	// (1, 2), (1, 3) (1, 4), (2, 3), (2, 4), (3, 4)
	// (1, 2, 3), (1, 2, 4), (1, 3, 4), (2, 3, 4)
	// (1, 2, 3, 4)
	// 요소의 수가 작은것부터 시작해서 조합에서 size를 만족할수 있다면 조합을 결과에 넣는다.
	// 결과적으로 validCombinations에 있는 집합의 크기는 모두 동일하다.
	var validCombinations [][]TopologyHintKey
	for k := 1; k <= len(unusedKeys); k++ {
		indicesCombinations := combin.Combinations(len(unusedKeys), k)
		combinationFound := false // 조합이 성공적으로 만들어졌는지 체크하는 변수

		for _, indices := range indicesCombinations {
			combinationKeys := []TopologyHintKey{}
			totalDevices := 0

			// 현재 조합에 있는 key들을 계산하여 totalDevices 값을 계산
			for _, idx := range indices {
				key := unusedKeys[idx]
				combinationKeys = append(combinationKeys, key)
				totalDevices += keyDeviceCounts[key]
			}

			// 조건을 만족하는 경우 validCombinations에 추가
			if totalDevices >= devicesToAllocate {
				validCombinations = append(validCombinations, combinationKeys)
				combinationFound = true
			}
		}

		// 조건을 만족하는 최소 조합의 크기에서 종료
		if combinationFound {
			break
		}
	}

	// step6: 만약 required keys가 존재한다면 위에서 만들어진 각각의 조합들에 required key를 추가해주어야 한다.
	requiredKeysSlice := make([]TopologyHintKey, 0, len(requiredKeys))
	for key := range requiredKeys {
		requiredKeysSlice = append(requiredKeysSlice, key)
	}

	for i := range validCombinations {
		validCombinations[i] = append(validCombinations[i], requiredKeysSlice...)
	}

	// step7: 채점을 하고 가장 높은 점수를 가지는 조합을 찾는다.
	highestScore := uint(0)
	bestCombination := []TopologyHintKey{}

	for _, combination := range validCombinations {
		score := calculateCombinationScore(combination, b.topologyHintMatrix)
		if score > highestScore {
			highestScore = score
			bestCombination = combination
		}
	}

	// step8: collected에 추가하고 반환한다.
	for _, key := range bestCombination {
		devices := availableByTopologyHintKey[key]
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
	// key가 1개짜리 조합인 경우에는 채점을 하는 모든 조합이 1개짜리이다. 이경우엔 조합을 채점하는 의미가 없기떄문에 0으로 처리한다.
	// 이부분은 아래 combin.Combination에서 n이 k보다 작으면 패닉이 나기때문에 간단한 처리를 해준다.
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
