package npu_allocator

import (
	"container/list"
	"sync"

	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
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
	// If 'request - len(required)' is zero, it means pre-allocated devices already satisfies device request quantity.
	if request-len(required) <= 0 {
		return required
	}

	deviceIdToDeviceMap := make(map[string]Device)
	for _, device := range available {
		deviceIdToDeviceMap[device.GetID()] = device
	}

	var highestScore = scoreDeviceSet(b.hintProvider, required)
	var bestDevices = required

	wg := new(sync.WaitGroup)
	lock := new(sync.Mutex)

	allocatedDevicesQueue := list.New()
	allocatedDevicesQueue.PushBack(required)

	for allocatedDevicesQueue.Len() > 0 {
		allocatedDevices := allocatedDevicesQueue.Remove(allocatedDevicesQueue.Front()).(DeviceSet)
		if len(allocatedDevices) == request {
			wg.Add(1)
			go func() {
				defer wg.Done()
				scoreSum := scoreDeviceSet(b.hintProvider, allocatedDevices)

				lock.Lock()
				defer lock.Unlock()

				if scoreSum > highestScore {
					highestScore = scoreSum
					bestDevices = allocatedDevices
				}
			}()
		} else {
			selectedDevices := b.selectBestScoredNewDevices(allocatedDevices, deviceIdToDeviceMap)
			for _, selectedDevice := range selectedDevices {
				nextAllocatedDevices := make(DeviceSet, len(allocatedDevices), len(allocatedDevices)+1)
				copy(nextAllocatedDevices, allocatedDevices)
				nextAllocatedDevices = append(nextAllocatedDevices, selectedDevice)

				allocatedDevicesQueue.PushBack(nextAllocatedDevices)
			}
		}
	}

	wg.Wait()

	return bestDevices
}

// selectBestScoredNewDevices selects devices which have the largest scoreSum with alreadyAllocatedDevices.
// It only returns newly selected devices, not including previously allocated devices.
// If some devices have same scoreSum, more than 1 device can be returned.
func (b *binPackingNpuAllocator) selectBestScoredNewDevices(
	alreadyAllocatedDevices DeviceSet,
	deviceIdToDeviceMap map[string]Device,
) DeviceSet {
	var highestScoreSum uint = 0
	var selectedDevices DeviceSet = nil

	allocatedDeviceIdMap := make(map[string]Device)
	for _, device := range alreadyAllocatedDevices {
		allocatedDeviceIdMap[device.GetID()] = device
	}

	wg := new(sync.WaitGroup)
	lock := new(sync.Mutex)

	for deviceId, device := range deviceIdToDeviceMap {
		// If deviceId belongs to the already allocated device, skip it.
		if _, ok := allocatedDeviceIdMap[deviceId]; ok {
			continue
		}

		// FIXME: Starting from go 1.22, below single line of codes can be removed.
		// Please see https://github.com/furiosa-ai/cloud-native-toolkit/issues/1#issuecomment-2301123405
		targetDevice := device

		wg.Add(1)
		go func() {
			defer wg.Done()

			scoringTargetDevices := make(DeviceSet, len(alreadyAllocatedDevices), len(alreadyAllocatedDevices)+1)
			copy(scoringTargetDevices, alreadyAllocatedDevices)
			scoringTargetDevices = append(scoringTargetDevices, targetDevice)

			scoreSum := scoreDeviceSet(b.hintProvider, scoringTargetDevices)

			lock.Lock()
			defer lock.Unlock()

			if selectedDevices == nil || scoreSum > highestScoreSum {
				highestScoreSum = scoreSum

				selectedDevices = make(DeviceSet, 0, 1)
				selectedDevices = append(selectedDevices, targetDevice)
			} else if scoreSum == highestScoreSum {
				selectedDevices = append(selectedDevices, targetDevice)
			}
		}()
	}

	wg.Wait()

	return selectedDevices
}
