package npu_allocator

import (
	furiosaSmi "github.com/furiosa-ai/libfuriosa-kubernetes/pkg/smi"

	"gonum.org/v1/gonum/stat/combin"
)

var _ NpuAllocator = (*scoreBasedOptimalNpuAllocator)(nil)

type topologyMatrix map[string]map[string]uint

type scoreBasedOptimalNpuAllocator struct {
	hintProvider TopologyHintProvider
}

func populateTopologyMatrix(devices []furiosaSmi.Device) (topologyMatrix, error) {
	topologyMatrix := make(topologyMatrix)
	deviceToDeviceInfo := make(map[furiosaSmi.Device]furiosaSmi.DeviceInfo)

	for _, device := range devices {
		deviceInfo, err := device.DeviceInfo()
		if err != nil {
			return nil, err
		}
		deviceToDeviceInfo[device] = deviceInfo
	}

	for device1, deviceInfo1 := range deviceToDeviceInfo {
		for device2, deviceInfo2 := range deviceToDeviceInfo {
			linkType, err := device1.GetDeviceToDeviceLinkType(device2)
			if err != nil {
				return nil, err
			}

			key1 := deviceInfo1.BDF()
			key2 := deviceInfo2.BDF()
			if key1 > key2 {
				key1, key2 = key2, key1
			}

			if _, ok := topologyMatrix[key1]; !ok {
				topologyMatrix[key1] = make(map[string]uint)
			}

			topologyMatrix[key1][key2] = uint(linkType)
		}
	}

	return topologyMatrix, nil
}

func NewScoreBasedOptimalNpuAllocator(devices []furiosaSmi.Device) (NpuAllocator, error) {
	topologyMatrix, err := populateTopologyMatrix(devices)
	if err != nil {
		return nil, err
	}

	return newScoreBasedOptimalNpuAllocator(
		func(device1, device2 Device) uint {
			if innerMap, exists := topologyMatrix[device1.TopologyHintKey()]; exists {
				if score, exists := innerMap[device2.TopologyHintKey()]; exists {
					return score
				}
			}

			return 0
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
	return n.hintProvider(device1, device2)
}
