package npu_allocator

import (
	"sync"

	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"gonum.org/v1/gonum/stat/combin"
)

var _ NpuAllocator = (*binPackingNpuAllocator)(nil)

type binPackingNpuAllocator struct {
	topologyScoreCalculator TopologyScoreCalculator
}

func NewBinPackingNpuAllocator(devices []smi.Device) (NpuAllocator, error) {
	topologyHintMatrix, err := populateTopologyHintMatrix(devices)
	if err != nil {
		return nil, err
	}

	return &binPackingNpuAllocator{
		topologyScoreCalculator: getTopologyScoreCalculatorBinPackingAllocator(topologyHintMatrix),
	}, nil
}

func (b *binPackingNpuAllocator) Allocate(available DeviceSet, required DeviceSet, size int) DeviceSet {
	// Step 1: available DeviceSet 을 TopologyHintKey 를 기준으로 map 을 만든다.
	availableDevicesByHintKeyMap := make(map[TopologyHintKey]DeviceSet)
	for _, device := range available {
		hintKey := device.GetTopologyHintKey()
		if _, ok := availableDevicesByHintKeyMap[hintKey]; !ok {
			availableDevicesByHintKeyMap[hintKey] = make(DeviceSet, 0)
		}

		availableDevicesByHintKeyMap[hintKey] = append(availableDevicesByHintKeyMap[hintKey], device)
	}

	// Step 2: required DeviceSet 을 먼저 처리한다. 같은 물리 카드에서 할당을 우선시하기 위해 requiredKey 를 수집한다.
	collectedDevices := make(DeviceSet, 0, size)
	requiredHintKeySet := make(map[TopologyHintKey]struct{})
	for _, device := range required {
		collectedDevices = append(collectedDevices, device)

		hintKey := device.GetTopologyHintKey()
		requiredHintKeySet[hintKey] = struct{}{}

		availableDevicesByHintKeyMap[hintKey] = availableDevicesByHintKeyMap[hintKey].Difference(DeviceSet{device})
	}

	if len(collectedDevices) >= size {
		return collectedDevices[:size]
	}

	// Step 3: required key 를 먼저 처리하여 이미 할당이 된 물리카드에서 먼저 할당을 한다.
	for hintKey := range requiredHintKeySet {
		devices := availableDevicesByHintKeyMap[hintKey]
		for _, device := range devices {
			collectedDevices = append(collectedDevices, device)
			availableDevicesByHintKeyMap[hintKey] = availableDevicesByHintKeyMap[hintKey].Difference(DeviceSet{device})
		}

		if len(collectedDevices) >= size {
			return collectedDevices[:size]
		}
	}

	// Step 4: 남은 부분을 계산한다.
	remainingDevicesSize := size - len(collectedDevices)

	unusedHintKeys := make([]TopologyHintKey, 0)
	deviceCountByHintKeyMap := make(map[TopologyHintKey]int)
	for hintKey, devices := range availableDevicesByHintKeyMap {
		if _, ok := requiredHintKeySet[hintKey]; !ok {
			unusedHintKeys = append(unusedHintKeys, hintKey)
		}

		deviceCountByHintKeyMap[hintKey] = len(devices)
	}

	// Step 5: 1개의 key 만 요소로 가지는 조합을 생성해서 최대 unusedHintKeys 의 길이만큼 요소를 가지는 조합을 생성한다.
	validHintKeysList := generateValidHintKeysCombinations(unusedHintKeys, deviceCountByHintKeyMap, remainingDevicesSize)

	// Step 6: 만약 required keys 가 존재한다면 위에서 만들어진 각각의 조합들에 required key 를 추가해 주어야 한다.
	requiredHintKeys := make([]TopologyHintKey, 0, len(requiredHintKeySet))
	for hintKey := range requiredHintKeySet {
		requiredHintKeys = append(requiredHintKeys, hintKey)
	}

	for i := range validHintKeysList {
		validHintKeysList[i] = append(validHintKeysList[i], requiredHintKeys...)
	}

	// Step 7: 채점을 하고 가장 높은 점수를 가지는 조합을 찾는다.
	highestScore := uint(0)
	bestHintKeys := make([]TopologyHintKey, 0)
	{
		wg := sync.WaitGroup{}
		lock := sync.Mutex{}

		for _, hintKeys := range validHintKeysList {
			go func() {
				wg.Add(1)
				wg.Done()

				score := b.topologyScoreCalculator(hintKeys)

				lock.Lock()
				defer lock.Unlock()

				if score > highestScore {
					highestScore = score
					bestHintKeys = hintKeys
				}
			}()
		}

		wg.Wait()
	}

	// Step 8: collectedDevices 에 추가하고 반환한다.
	for _, hintKey := range bestHintKeys {
		devices := availableDevicesByHintKeyMap[hintKey]
		collectedDevices = append(collectedDevices, devices...)
	}

	return collectedDevices[:size]
}

// getTopologyScoreCalculatorBinPackingAllocator calculates total sum of given TopologyHintKey list.
func getTopologyScoreCalculatorBinPackingAllocator(topologyHintMatrix TopologyHintMatrix) TopologyScoreCalculator {
	return func(keys []TopologyHintKey) uint {
		// key 가 1개짜리 조합인 경우에는 채점을 하는 모든 조합이 1개짜리 이다. 이 경우엔 조합을 채점하는 의미가 없기 때문에 0 으로 처리한다.
		// 이 부분은 아래 combin.Combinations 에서 n이 k보다 작으면 패닉이 나기 때문에 간단한 처리를 해준다.
		if len(keys) <= 1 {
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

func generateValidHintKeysCombinations(unusedHintKeys []TopologyHintKey, deviceCountByHintKeyMap map[TopologyHintKey]int, remainingDevicesSize int) [][]TopologyHintKey {
	// 1, 2, 3, 4란 key 가 있다고 하면 아래와 같이 조합을 생성한다.
	// (1), (2), (3), (4)
	// (1, 2), (1, 3) (1, 4), (2, 3), (2, 4), (3, 4)
	// (1, 2, 3), (1, 2, 4), (1, 3, 4), (2, 3, 4)
	// (1, 2, 3, 4)
	// 요소의 수가 작은것부터 시작해서 조합에서 size를 만족할수 있다면 조합을 결과에 넣는다.
	// 결과적으로 validHintKeysList 에 있는 집합의 크기는 모두 동일하다.

	validHintKeysList := make([][]TopologyHintKey, 0)
	for k := 1; k <= len(unusedHintKeys); k++ {
		wg := sync.WaitGroup{}
		lock := sync.Mutex{}

		indicesCombinations := combin.Combinations(len(unusedHintKeys), k)
		for _, indices := range indicesCombinations {
			go func() {
				wg.Add(1)
				defer wg.Done()

				hintKeys := make([]TopologyHintKey, 0)
				totalDevices := 0

				// 현재 조합에 있는 key 들을 계산하여 totalDevices 값을 계산
				for _, idx := range indices {
					hintKey := unusedHintKeys[idx]
					hintKeys = append(hintKeys, hintKey)

					totalDevices += deviceCountByHintKeyMap[hintKey]
				}

				// 조건을 만족하는 경우 validHintKeysList 에 추가
				lock.Lock()
				defer lock.Unlock()

				if totalDevices >= remainingDevicesSize {
					validHintKeysList = append(validHintKeysList, hintKeys)
				}
			}()
		}

		wg.Wait()

		// 조건을 만족하는 최소 조합의 크기에서 종료
		if len(validHintKeysList) > 0 {
			break
		}
	}

	return validHintKeysList
}
