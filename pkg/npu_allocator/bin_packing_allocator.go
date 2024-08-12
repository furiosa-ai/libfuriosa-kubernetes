package npu_allocator

import (
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/smi"
	"math"
)

var _ NpuAllocator = (*binPackingNpuAllocator)(nil)

type binPackingNpuAllocator struct{}

func NewBinPackingNpuAllocator(devices []smi.Device) (NpuAllocator, error) {
	return &binPackingNpuAllocator{}, nil
}

func (b *binPackingNpuAllocator) Allocate(available DeviceSet, required DeviceSet, request int) DeviceSet {
	remainingCnt := request - len(required)
	// If remainingCnt is zero, it means pre-allocated devices already satisfies device request quantity.
	if remainingCnt == 0 {
		return required
	}

	// candidates contains devices in `available` set, excluding `required` set.
	candidates := available.Difference(required)
	candidatesByHintMap := make(map[string]DeviceSet) // construct map by TopologyHintKey and DeviceSet
	for _, candidate := range candidates {
		topologyHintKey := candidate.TopologyHintKey()
		if _, ok := candidatesByHintMap[topologyHintKey]; !ok {
			candidatesByHintMap[topologyHintKey] = make(DeviceSet, 0)
		}

		candidatesByHintMap[topologyHintKey] = append(candidatesByHintMap[topologyHintKey], candidate)
	}

	// finalizedDevices contains finalized allocated devices.
	finalizedDevices := make(DeviceSet, 0, request)
	finalizedDevices = append(finalizedDevices, required...)

	for remainingCnt > 0 {
		// Step 1: Use Best Fit Bin Packing algorithm to select candidates.
		minDiff := math.MaxInt32
		selectedTopologyHintKey := ""
		for key, partialCandidates := range candidatesByHintMap {
			diff := len(partialCandidates) - remainingCnt
			if diff >= 0 && diff < minDiff {
				minDiff = diff
				selectedTopologyHintKey = key
			}
		}

		// If Step 1 failed to find proper candidates, `selectedTopologyHintKey` will be empty.
		if selectedTopologyHintKey != "" {
			finalizedDevices = append(finalizedDevices, candidatesByHintMap[selectedTopologyHintKey][:remainingCnt]...)
			candidatesByHintMap[selectedTopologyHintKey] = candidatesByHintMap[selectedTopologyHintKey][remainingCnt:]
			remainingCnt = 0
			break
		}

		// Step 2: Find candidates which have the largest length.
		maxLen := 0
		selectedTopologyHintKey = ""
		for key, partialCandidates := range candidatesByHintMap {
			if selectedTopologyHintKey == "" || len(partialCandidates) > maxLen {
				maxLen = len(partialCandidates)
				selectedTopologyHintKey = key
			}
		}

		finalizedDevices = append(finalizedDevices, candidatesByHintMap[selectedTopologyHintKey]...)
		remainingCnt -= len(candidatesByHintMap[selectedTopologyHintKey])
		delete(candidatesByHintMap, selectedTopologyHintKey)
	}

	return finalizedDevices
}
