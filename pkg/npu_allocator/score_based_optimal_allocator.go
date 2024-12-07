package npu_allocator

import (
	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"gonum.org/v1/gonum/stat/combin"
)

var _ NpuAllocator = (*scoreBasedOptimalNpuAllocator)(nil)

type scoreBasedOptimalNpuAllocator struct {
	hintProvider TopologyHintProvider
}

func NewScoreBasedOptimalNpuAllocator(devices []smi.Device) (NpuAllocator, error) {
	topologyHintMatrix, err := NewTopologyHintMatrix(devices)
	if err != nil {
		return nil, err
	}

	hintProvider := func(device1, device2 Device) uint {
		key1, key2 := device1.TopologyHintKey(), device2.TopologyHintKey()
		if key1 > key2 {
			key1, key2 = key2, key1
		}

		if innerMap, innerMapExists := topologyHintMatrix[key1]; innerMapExists {
			if score, scoreExists := innerMap[key2]; scoreExists {
				return score
			}
		}

		return 0
	}

	return newScoreBasedOptimalNpuAllocator(hintProvider), nil
}

func NewMockScoreBasedOptimalNpuAllocator(mockHintProvider TopologyHintProvider) (NpuAllocator, error) {
	return newScoreBasedOptimalNpuAllocator(mockHintProvider), nil
}

func newScoreBasedOptimalNpuAllocator(hintProvider TopologyHintProvider) NpuAllocator {
	return &scoreBasedOptimalNpuAllocator{
		hintProvider: hintProvider,
	}
}

func (n *scoreBasedOptimalNpuAllocator) Allocate(available DeviceSet, required DeviceSet, request int) DeviceSet {
	subsetLen := request - required.Len()
	// length of required equals to request, it means allocating specific device sets
	if subsetLen == 0 {
		return required
	}

	// generate seed sets using differences
	difference := available.Difference(required.Devices()...)
	combinations := generateKDeviceSet(difference, subsetLen)

	// union subset and required to build full device set combination
	for idx, combination := range combinations {
		newDeviceSet := combination.Union(required.Devices()...)
		combinations[idx] = newDeviceSet
	}

	// score all survived device set
	// initialize with the first element to prevent edge case that score of all element in the filtered list is zero.
	var bestSet = combinations[0]
	var highestScore = n.scoreDeviceSet(bestSet)

	for _, set := range combinations {
		if score := n.scoreDeviceSet(set); score > highestScore {
			bestSet = set
			highestScore = score
		}
	}

	//pick the best one
	return bestSet
}

func generateKDeviceSet(ds DeviceSet, size int) (result []DeviceSet) {
	// NOTE(@bg): combin.Combinations internally uses binomial coefficient C(n, k) implementation,
	// which call panic() if k > n and either n and k is negative number.
	// https://github.com/gonum/gonum/blob/f74f45f5f3e9cc7c1d0f0af2ffd19ccf8972a87e/stat/combin/combin.go#L29
	if ds.Len() < 1 || size < 1 || size > ds.Len() {
		return result
	}

	devices := ds.Devices()
	for _, indices := range combin.Combinations(len(devices), size) {
		newDeviceSet := NewDeviceSet()
		for _, index := range indices {
			newDeviceSet.Insert(devices[index])
		}

		result = append(result, newDeviceSet)
	}

	return result
}

func (n *scoreBasedOptimalNpuAllocator) scoreDeviceSet(deviceSet DeviceSet) uint {
	total := uint(0)

	// calculate total score using distance of  two device
	for i, d1 := range deviceSet.Devices() {
		for j, d2 := range deviceSet.Devices() {
			if j > i {
				total += n.scoreDevicePair(d1, d2)
			}
		}
	}

	return total
}

func (n *scoreBasedOptimalNpuAllocator) scoreDevicePair(device1 Device, device2 Device) uint {
	return n.hintProvider(device1, device2)
}
