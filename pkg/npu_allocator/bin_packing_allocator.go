package npu_allocator

import (
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/smi"
	"math"
	"sort"
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
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].ID() < candidates[j].ID()
	})

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
		selectedTopologyHintKey := getTopologyHintKeyUsingBestFitBinPacking(remainingCnt, &candidatesByHintMap)
		if selectedTopologyHintKey != "" {
			finalizedDevices = append(finalizedDevices, candidatesByHintMap[selectedTopologyHintKey][:remainingCnt]...)
			candidatesByHintMap[selectedTopologyHintKey] = candidatesByHintMap[selectedTopologyHintKey][remainingCnt:]
			remainingCnt = 0
			break
		}

		// Step 2: Find candidates which have the largest length.
		selectedTopologyHintKey = getLargestLengthCandidatesTopologyHintKey(&candidatesByHintMap)
		finalizedDevices = append(finalizedDevices, candidatesByHintMap[selectedTopologyHintKey]...)
		remainingCnt -= len(candidatesByHintMap[selectedTopologyHintKey])
		delete(candidatesByHintMap, selectedTopologyHintKey)
	}

	return finalizedDevices
}

// getTopologyHintKeyUsingBestFitBinPacking uses Best Fit Bin Packing algorithm to select candidates key
func getTopologyHintKeyUsingBestFitBinPacking(remainingCnt int, candidatesByHintMap *map[string]DeviceSet) string {
	minDiff := math.MaxInt32
	topologyHintKey := ""
	for key, candidates := range *candidatesByHintMap {
		diff := len(candidates) - remainingCnt
		if diff >= 0 && diff < minDiff {
			minDiff = diff
			topologyHintKey = key
		}
	}

	return topologyHintKey
}

// getLargestLengthCandidatesTopologyHintKey selects candidates key which has the largest length
func getLargestLengthCandidatesTopologyHintKey(candidatesByHintMap *map[string]DeviceSet) string {
	maxLen := 0
	topologyHintKey := ""
	for key, candidates := range *candidatesByHintMap {
		if topologyHintKey == "" || len(candidates) > maxLen {
			maxLen = len(candidates)
			topologyHintKey = key
		}
	}

	return topologyHintKey
}
