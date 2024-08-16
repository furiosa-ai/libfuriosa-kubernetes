package npu_allocator

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

	return &binPackingNpuAllocator{
		hintProvider: newTopologyHintProviderForBinPackingAllocator(topologyHintMatrix),
	}, nil
}

func newTopologyHintProviderForBinPackingAllocator(topologyHintMatrix TopologyHintMatrix) TopologyHintProvider {
	return func(device1, device2 Device) uint {
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
}

func (b *binPackingNpuAllocator) Allocate(available DeviceSet, required DeviceSet, request int) DeviceSet {
	return nil
}
