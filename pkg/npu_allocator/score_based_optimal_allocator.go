package npu_allocator

import (
	topology "github.com/furiosa-ai/libfuriosa-kubernetes/pkg/device"
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

func newScoreBasedOptimalNpuAllocator(hintProvider TopologyHintProvider) NpuAllocator {
	return &scoreBasedOptimalNpuAllocator{
		hintProvider: hintProvider,
	}
}

func (n *scoreBasedOptimalNpuAllocator) Allocate(available DeviceSet, required DeviceSet, request int) DeviceSet {
	// generate all possible device set
	combinations := generateNonDuplicatedDeviceSet(available, request)

	// filter device set doesn't contain required devices
	var filtered []DeviceSet
	if len(required) == 0 {
		filtered = combinations
	} else {
		for _, combination := range combinations {
			if combination.Contains(required) {
				filtered = append(filtered, combination)
			}
		}
	}

	// no device set survived
	if len(filtered) == 0 {
		return nil
	}

	// score all survived device set
	// initialize with the first element to prevent edge case that score of all element in the filtered list is zero.
	var bestSet = filtered[0]
	var highestScore = n.scoreDeviceSet(bestSet)

	for _, set := range filtered {
		score := n.scoreDeviceSet(set)
		if score > highestScore {
			bestSet = set
			highestScore = score
		}
	}

	//pick the best one
	return bestSet
}

func generateNonDuplicatedDeviceSet(devices DeviceSet, size int) (result []DeviceSet) {
	if len(devices) == 0 || size == 0 || size > len(devices) {
		return result
	}

	devices.Sort()

	// use iterative approach for memory efficient
	total := len(devices)
	indices := make([]int, size)

	// initialize indices
	for i := range indices {
		indices[i] = i
	}

	for {
		newDeviceSet := make(DeviceSet, size)
		// generate combination at the current position
		for i, index := range indices {
			newDeviceSet[i] = devices[index]
		}
		result = append(result, newDeviceSet)

		// initialize pivot with the last index
		pivot := size - 1

		// move(decrease) pivot if it reached max at the current position.
		for pivot >= 0 && indices[pivot] == pivot+total-size {
			pivot--
		}

		// exit loop we have visited all combination
		if pivot < 0 {
			break
		}

		// visit next element at the current position
		indices[pivot]++

		// initialize indices next to pivot to `pivot +1, +1, +1` if pivot moved
		for i := pivot + 1; i < size; i++ {
			indices[i] = indices[i-1] + 1
		}

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
