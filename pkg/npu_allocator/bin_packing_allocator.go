package npu_allocator

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/smi"
)

var _ NpuAllocator = (*binPackingNpuAllocator)(nil)

type binPackingNpuAllocator struct {
	hintProvider TopologyHintProvider
}

func populateTopologyHintMatrixForBinPackingAllocator(smiDevices []smi.Device) (TopologyHintMatrix, error) {
	topologyHintMatrix := make(TopologyHintMatrix)
	deviceToDeviceInfo := make(map[smi.Device]smi.DeviceInfo)

	for _, device := range smiDevices {
		deviceInfo, err := device.DeviceInfo()
		if err != nil {
			return nil, err
		}

		deviceToDeviceInfo[device] = deviceInfo
	}

	// parseBusIDFromBDF parses bdf and returns PCI bus ID.
	parseBusIDFromBDF := func(bdf string) (string, error) {
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

	for device1, deviceInfo1 := range deviceToDeviceInfo {
		for device2, deviceInfo2 := range deviceToDeviceInfo {
			linkType, err := device1.GetDeviceToDeviceLinkType(device2)
			if err != nil {
				return nil, err
			}

			score := uint(linkType)

			pciBusID1, err := parseBusIDFromBDF(deviceInfo1.BDF())
			if err != nil {
				return nil, err
			}

			pciBusID2, err := parseBusIDFromBDF(deviceInfo2.BDF())
			if err != nil {
				return nil, err
			}

			key1, key2 := TopologyHintKey(pciBusID1), TopologyHintKey(pciBusID2)
			if key1 > key2 {
				key1, key2 = key2, key1
			}

			if _, ok := topologyHintMatrix[key1]; !ok {
				topologyHintMatrix[key1] = make(map[TopologyHintKey]uint)
			}

			topologyHintMatrix[key1][key2] = score
		}
	}

	return topologyHintMatrix, nil
}

func NewBinPackingNpuAllocator(smiDevices []smi.Device) (NpuAllocator, error) {
	topologyHintMatrix, err := populateTopologyHintMatrixForBinPackingAllocator(smiDevices)
	if err != nil {
		return nil, err
	}

	hintProvider := func(device1, device2 Device) uint {
		key1, key2 := device1.GetTopologyHintKey(), device2.GetTopologyHintKey()
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

	return &binPackingNpuAllocator{hintProvider: hintProvider}, nil
}

func (b *binPackingNpuAllocator) Allocate(available DeviceSet, required DeviceSet, request int) DeviceSet {
	subsetLen := request - len(required)
	// If subsetLen is zero, it means pre-allocated devices already satisfies device request quantity.
	if subsetLen <= 0 {
		return required
	}

	// difference contains devices in `available` set, excluding `required` set.
	difference := available.Difference(required)
	remainingDevicesByHintMap := make(map[TopologyHintKey]DeviceSet)
	for _, device := range difference {
		topologyHintKey := device.GetTopologyHintKey()
		if _, ok := remainingDevicesByHintMap[topologyHintKey]; !ok {
			remainingDevicesByHintMap[topologyHintKey] = make(DeviceSet, 0)
		}

		remainingDevicesByHintMap[topologyHintKey] = append(remainingDevicesByHintMap[topologyHintKey], device)
	}

	// finalizedDevices contains finalized allocated devices.
	finalizedDevices := make(DeviceSet, 0, request)
	finalizedDevices = finalizedDevices.Union(required)

	for subsetLen > 0 {
		selectedDevices := b.selectBestScoredDevices(subsetLen, finalizedDevices, remainingDevicesByHintMap)
		subsetLen -= len(selectedDevices)
		finalizedDevices = finalizedDevices.Union(selectedDevices)
	}

	return finalizedDevices
}

// scoreDeviceSet returns total sum of scores for each pair of devices.
func (b *binPackingNpuAllocator) scoreDeviceSet(deviceSet DeviceSet) uint {
	var total uint = 0
	for i := 0; i < len(deviceSet); i++ {
		for j := i + 1; j < len(deviceSet); j++ {
			total += b.scoreDevicePair(deviceSet[i], deviceSet[j])
		}
	}

	return total
}

// scoreDevicePair returns score based on distance between two devices.
// Higher score means lower distance.
func (b *binPackingNpuAllocator) scoreDevicePair(device1 Device, device2 Device) uint {
	return b.hintProvider(device1, device2)
}

func (b *binPackingNpuAllocator) selectBestScoredDevices(subsetLen int, finalizedDevices DeviceSet, remainingDevicesByHintMap map[TopologyHintKey]DeviceSet) DeviceSet {
	var highestScore uint = 0
	var selectedHintKey TopologyHintKey = ""

	wg := new(sync.WaitGroup)
	lock := new(sync.Mutex)

	for topologyHintKey, devices := range remainingDevicesByHintMap {
		wg.Add(1)
		go func() {
			defer wg.Done()

			partialDevices := devices
			if len(devices) > subsetLen {
				partialDevices = devices[:subsetLen]
			}

			score := b.scoreDeviceSet(partialDevices.Union(finalizedDevices))

			lock.Lock()
			if selectedHintKey == "" || highestScore < score {
				highestScore = score
				selectedHintKey = topologyHintKey
			}
			lock.Unlock()
		}()
	}

	wg.Wait()

	selectedDevices := remainingDevicesByHintMap[selectedHintKey]
	if len(selectedDevices) > subsetLen {
		selectedDevices = selectedDevices[:subsetLen]
		remainingDevicesByHintMap[selectedHintKey] = remainingDevicesByHintMap[selectedHintKey][subsetLen:]
	} else {
		delete(remainingDevicesByHintMap, selectedHintKey)
	}

	return finalizedDevices.Union(selectedDevices)
}
