package npu_allocator

import (
	"sync"

	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/smi"
)

var _ NpuAllocator = (*binPackingNpuAllocator)(nil)

type binPackingNpuAllocator struct {
	hintProvider TopologyHintProvider
}

func NewBinPackingNpuAllocator(smiDevices []smi.Device) (NpuAllocator, error) {
	topologyHintMatrix, err := populateTopologyHintMatrixFromSMIDevices(smiDevices)
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
	differenceByHintMap := make(map[TopologyHintKey]DeviceSet)
	for _, device := range difference {
		topologyHintKey := device.GetTopologyHintKey()
		if _, ok := differenceByHintMap[topologyHintKey]; !ok {
			differenceByHintMap[topologyHintKey] = make(DeviceSet, 0)
		}

		differenceByHintMap[topologyHintKey] = append(differenceByHintMap[topologyHintKey], device)
	}

	// allocatedDevices contains finalized allocated devices.
	allocatedDevices := make(DeviceSet, 0, request)
	allocatedDevices = allocatedDevices.Union(required)

	for subsetLen > 0 {
		selectedDevices := b.selectBestScoredDevices(subsetLen, allocatedDevices, differenceByHintMap)
		subsetLen -= len(selectedDevices)
		allocatedDevices = allocatedDevices.Union(selectedDevices)
	}

	return allocatedDevices
}

func (b *binPackingNpuAllocator) selectBestScoredDevices(subsetLen int, allocatedDevices DeviceSet, remainingDevicesByHintMap map[TopologyHintKey]DeviceSet) DeviceSet {
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

			score := b.scoreDeviceSet(partialDevices.Union(allocatedDevices))

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

	return allocatedDevices.Union(selectedDevices)
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
