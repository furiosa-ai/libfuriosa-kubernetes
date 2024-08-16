package npu_allocator

import (
	"fmt"
	"math"

	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/smi"
)

var _ NpuAllocator = (*binPackingNpuAllocator)(nil)

type binPackingNpuAllocator struct{}

func NewBinPackingNpuAllocator(_ []smi.Device) (NpuAllocator, error) {
	return &binPackingNpuAllocator{}, nil
}

func (b *binPackingNpuAllocator) Allocate(available DeviceSet, required DeviceSet, request int) DeviceSet {
	fmt.Printf("available: %d, required: %d, request: %d\n", len(available), len(required), request)

	subsetLen := request - len(required)
	// If subsetLen is zero, it means pre-allocated devices already satisfies device request quantity.
	if subsetLen == 0 {
		return required
	}

	// difference contains devices in `available` set, excluding `required` set.
	difference := available.Difference(required)

	differenceByHintMap := make(map[TopologyHintKey]DeviceSet) // construct map by GetTopologyHintKey and DeviceSet
	for _, device := range difference {
		topologyHintKey := device.GetTopologyHintKey()
		if _, ok := differenceByHintMap[topologyHintKey]; !ok {
			differenceByHintMap[topologyHintKey] = make(DeviceSet, 0)
		}

		differenceByHintMap[topologyHintKey] = append(differenceByHintMap[topologyHintKey], device)
	}

	// finalizedDevices contains finalized allocated devices.
	finalizedDevices := make(DeviceSet, 0, request)
	finalizedDevices = finalizedDevices.Union(required)

	for subsetLen > 0 {
		// Step 1: Use Best Fit Bin Packing algorithm to select difference.
		selectedTopologyHintKey := getTopologyHintKeyUsingBestFitBinPacking(subsetLen, differenceByHintMap)
		if selectedTopologyHintKey != "" {
			finalizedDevices = finalizedDevices.Union(differenceByHintMap[selectedTopologyHintKey][:subsetLen])
			differenceByHintMap[selectedTopologyHintKey] = differenceByHintMap[selectedTopologyHintKey][subsetLen:]
			break
		}

		// Step 2: Find difference which have the largest length.
		selectedTopologyHintKey = getLargestLengthDifferenceTopologyHintKey(differenceByHintMap)
		finalizedDevices = finalizedDevices.Union(differenceByHintMap[selectedTopologyHintKey])
		subsetLen -= len(differenceByHintMap[selectedTopologyHintKey])
		delete(differenceByHintMap, selectedTopologyHintKey)
	}

	return finalizedDevices
}

// getTopologyHintKeyUsingBestFitBinPacking uses Best Fit Bin Packing algorithm to select difference key
func getTopologyHintKeyUsingBestFitBinPacking(subsetLen int, differenceByHintMap map[TopologyHintKey]DeviceSet) TopologyHintKey {
	minDiff := math.MaxInt32
	var topologyHintKey TopologyHintKey = ""
	for key, difference := range differenceByHintMap {
		diff := len(difference) - subsetLen
		if diff >= 0 && diff < minDiff {
			minDiff = diff
			topologyHintKey = key
		}
	}

	return topologyHintKey
}

// getLargestLengthDifferenceTopologyHintKey selects difference key which has the largest length
func getLargestLengthDifferenceTopologyHintKey(differenceByHintMap map[TopologyHintKey]DeviceSet) TopologyHintKey {
	maxLen := 0
	var topologyHintKey TopologyHintKey = ""
	for key, difference := range differenceByHintMap {
		if topologyHintKey == "" || len(difference) > maxLen {
			maxLen = len(difference)
			topologyHintKey = key
		}
	}

	return topologyHintKey
}
