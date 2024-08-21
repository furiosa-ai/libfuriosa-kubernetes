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
	topologyHintMatrix, err := populateTopologyHintMatrix(smiDevices)
	if err != nil {
		return nil, err
	}

	return newBinPackingNpuAllocator(getGenericHintProvider(topologyHintMatrix)), nil
}

func NewMockBinPackingNpuAllocator(mockHintProvider TopologyHintProvider) (NpuAllocator, error) {
	return newBinPackingNpuAllocator(mockHintProvider), nil
}

func newBinPackingNpuAllocator(hintProvider TopologyHintProvider) NpuAllocator {
	return &binPackingNpuAllocator{hintProvider: hintProvider}
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
		// select best scored devices at every iteration, until subsetLen reaches 0.
		selectedDevices := b.selectBestScoredNewDevices(subsetLen, allocatedDevices, differenceByHintMap)
		subsetLen -= len(selectedDevices)
		allocatedDevices = allocatedDevices.Union(selectedDevices)
	}

	return allocatedDevices
}

// selectBestScoredNewDevices selects devices which get the largest score with previouslyAllocatedDevices.
// It only returns newly selected devices, not including previously allocated devices.
func (b *binPackingNpuAllocator) selectBestScoredNewDevices(
	maxSelectLength int,
	previouslyAllocatedDevices DeviceSet,
	remainingDevicesByHintMap map[TopologyHintKey]DeviceSet,
) DeviceSet {
	var highestScoreAvg = 0.0
	var selectedHintKey TopologyHintKey = ""

	wg := new(sync.WaitGroup)
	lock := new(sync.Mutex)

	for topologyHintKey, devices := range remainingDevicesByHintMap {
		// FIXME: Starting from go 1.22, below single line of codes can be removed.
		// Please see https://github.com/furiosa-ai/cloud-native-toolkit/issues/1#issuecomment-2301123405
		hintKey, devs := topologyHintKey, devices

		wg.Add(1)
		go func() {
			defer wg.Done()

			partialDevices := devs
			if len(partialDevices) > maxSelectLength {
				partialDevices = devs[:maxSelectLength]
			}

			scoringTargetDevices := previouslyAllocatedDevices.Union(partialDevices)
			scoreSum := scoreDeviceSet(b.hintProvider, scoringTargetDevices)
			scoreAvg := float64(scoreSum) / float64(len(scoringTargetDevices))

			lock.Lock()
			if selectedHintKey == "" || highestScoreAvg < scoreAvg {
				highestScoreAvg = scoreAvg
				selectedHintKey = hintKey
			}
			lock.Unlock()
		}()
	}

	wg.Wait()

	selectedDevices := remainingDevicesByHintMap[selectedHintKey]
	if len(selectedDevices) > maxSelectLength {
		// if length of selected device is longer than maxSelectLength, cut it.
		selectedDevices = selectedDevices[:maxSelectLength]
		remainingDevicesByHintMap[selectedHintKey] = remainingDevicesByHintMap[selectedHintKey][maxSelectLength:]
	} else {
		delete(remainingDevicesByHintMap, selectedHintKey)
	}

	return selectedDevices
}
