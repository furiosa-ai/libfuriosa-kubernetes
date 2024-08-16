package npu_allocator

import (
	"fmt"
	"regexp"

	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/smi"
	"gonum.org/v1/gonum/stat/combin"
)

var _ NpuAllocator = (*scoreBasedOptimalNpuAllocator)(nil)

type scoreBasedOptimalNpuAllocator struct {
	hintProvider TopologyHintProvider
}

// populateTopologyHintMatrixForScoreBasedAllocator generates TopologyHintMatrix using list of smi.Device.
func populateTopologyHintMatrixForScoreBasedAllocator(smiDevices []smi.Device) (TopologyHintMatrix, error) {
	topologyHintMatrix := make(TopologyHintMatrix)
	deviceToDeviceInfo := make(map[smi.Device]smi.DeviceInfo)

	for _, device := range smiDevices {
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

			busID1, err := parseBusIDFromBDF(deviceInfo1.BDF())
			if err != nil {
				return nil, err
			}

			busID2, err := parseBusIDFromBDF(deviceInfo2.BDF())
			if err != nil {
				return nil, err
			}

			key1 := TopologyHintKey(busID1)
			key2 := TopologyHintKey(busID2)
			if key1 > key2 {
				key1, key2 = key2, key1
			}

			if _, ok := topologyHintMatrix[key1]; !ok {
				topologyHintMatrix[key1] = make(map[TopologyHintKey]uint)
			}

			topologyHintMatrix[key1][key2] = uint(linkType)
		}
	}

	return topologyHintMatrix, nil
}

// parseBusIDFromBDF parses bdf and returns PCI bus ID.
func parseBusIDFromBDF(bdf string) (string, error) {
	bdfPattern := `^(?P<domain>[0-9a-fA-F]{1,4}):(?P<bus>[0-9a-fA-F]+):(?P<function>[0-9a-fA-F]+\.[0-9])$`
	subExpKeyBus := "bus"
	bdfRegExp := regexp.MustCompile(bdfPattern)

	matches := bdfRegExp.FindStringSubmatch(bdf)
	if matches == nil {
		return "", fmt.Errorf("couldn't parse the given string %s with bdf regex pattern: %s", bdf, bdfPattern)
	}

	subExpIndex := bdfRegExp.SubexpIndex(subExpKeyBus)
	if subExpIndex == -1 {
		return "", fmt.Errorf("couldn't parse bus id from the given bdf expression %s", bdf)
	}

	return matches[subExpIndex], nil
}

func NewScoreBasedOptimalNpuAllocator(devices []smi.Device) (NpuAllocator, error) {
	topologyHintMatrix, err := populateTopologyHintMatrixForScoreBasedAllocator(devices)
	if err != nil {
		return nil, err
	}

	hintProvider := func(device1, device2 Device) uint {
		if innerMap, innerMapExists := topologyHintMatrix[device1.GetTopologyHintKey()]; innerMapExists {
			if score, scoreExists := innerMap[device2.GetTopologyHintKey()]; scoreExists {
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
