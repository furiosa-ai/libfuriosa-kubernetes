package npu_allocator

import (
	topology "github.com/furiosa-ai/libfuriosa-kubernetes/pkg/device"
	"gonum.org/v1/gonum/stat/combin"
)

var _ NpuAllocator = (*scoreBasedOptimalNpuAllocator)(nil)

type scoreBasedOptimalNpuAllocator struct {
	hintProvider TopologyHintProvider
}

func NewScoreBasedOptimalNpuAllocator(device []topology.Device) (NpuAllocator, error) {
	newTopology, err := topology.NewTopology(device)
	if err != nil {
		return nil, err
	}

	return newScoreBasedOptimalNpuAllocator(func(topologyHintKey1, topologyHintKey2 string) uint {
		return uint(newTopology.GetLinkType(topologyHintKey1, topologyHintKey2))
	}), nil
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
	subsetLen := request - len(required)
	// length of required equals to request, it means allocating specific device sets
	if subsetLen == 0 {
		return required
	}

	// generate seed sets using differences
	difference := available.Difference(required)
	combinations := generateKDeviceSet(difference, subsetLen)

	// union subset and required to build full device set combination
	for idx, combination := range combinations {
		newDeviceSet := combination.Union(required)
		newDeviceSet.Sort()
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

func generateKDeviceSet(devices DeviceSet, size int) (result []DeviceSet) {
	// NOTE(@bg): combin.Combinations internally uses binomial coefficient C(n, k) implementation,
	// which call panic() if k > n and either n and k is negative number.
	// https://github.com/gonum/gonum/blob/f74f45f5f3e9cc7c1d0f0af2ffd19ccf8972a87e/stat/combin/combin.go#L29
	if len(devices) < 1 || size < 1 || size > len(devices) {
		return result
	}

	for _, indices := range combin.Combinations(len(devices), size) {
		newDeviceSet := DeviceSet{}
		for _, index := range indices {
			newDeviceSet = append(newDeviceSet, devices[index])
		}
		result = append(result, newDeviceSet)
	}

	return result
}

func (n *scoreBasedOptimalNpuAllocator) scoreDeviceSet(deviceSet DeviceSet) uint {
	total := uint(0)

	// calculate total score using distance of  two device
	for i, d1 := range deviceSet {
		for j, d2 := range deviceSet {
			if j > i {
				total += n.scoreDevicePair(d1, d2)
			}
		}
	}

	return total
}

func (n *scoreBasedOptimalNpuAllocator) scoreDevicePair(device1 Device, device2 Device) uint {
	return n.hintProvider(device1.TopologyHintKey(), device2.TopologyHintKey())
}
