package npu_allocator

import "sync"

var _ NpuAllocator = (*binPackingNpuAllocator)(nil)

type binPackingNpuAllocator struct {
	hintProvider TopologyHintProvider
}

func populateTopologyHintMatrixForBinPackingAllocator(devices DeviceSet) (TopologyHintMatrix, error) {
	topologyHintMatrix := make(TopologyHintMatrix)

	for _, device1 := range devices {
		for _, device2 := range devices {
			score := device1.CalculateScoreToOtherDevice(device2)

			key1, key2 := device1.GetTopologyHintKey(), device2.GetTopologyHintKey()
			if key1 > key2 {
				key1, key2 = key2, key1
			}

			if _, exists := topologyHintMatrix[key1]; !exists {
				topologyHintMatrix[key1] = make(map[TopologyHintKey]uint)
			}

			topologyHintMatrix[key1][key2] = score
		}
	}

	return topologyHintMatrix, nil
}

func NewBinPackingNpuAllocator(devices DeviceSet) (NpuAllocator, error) {
	topologyHintMatrix, err := populateTopologyHintMatrixForBinPackingAllocator(devices)
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
