package npu_allocator

import (
	"sort"
)

// TopologyHintProvider takes two devices as argument return topology hint.
// The hint would be score, distance, preference of two devices.
type TopologyHintProvider func(device1, device2 Device) uint

// TopologyHintKey is named type of string, used for TopologyHintMatrix
type TopologyHintKey string

// TopologyHintMatrix provides score of device to device based on smi.Device smi.LinkType.
type TopologyHintMatrix map[TopologyHintKey]map[TopologyHintKey]uint

type NpuAllocator interface {
	Allocate(available DeviceSet, required DeviceSet, size int) DeviceSet
}

type Device interface {
	// GetID returns a unique ID of Device to identify the device.
	GetID() string

	// GetTopologyHintKey returns unique key to retrieve TopologyHint using TopologyHintProvider.
	GetTopologyHintKey() TopologyHintKey

	// Equal check whether source Device is identical to the target Device.
	Equal(target Device) bool
}

type DeviceSet []Device

// Contains checks whether source DeviceSet contains target DeviceSet.
func (source DeviceSet) Contains(target DeviceSet) bool {
	if len(source) == 0 || len(target) == 0 {
		return false
	}

	visited := map[string]bool{}
	for _, device := range source {
		visited[device.GetID()] = true
	}

	for _, device := range target {
		if _, ok := visited[device.GetID()]; !ok {
			return false
		}
	}

	return true
}

// Sort sorts source DeviceSet.
func (source DeviceSet) Sort() {
	sort.Slice(source, func(i, j int) bool {
		return source[i].GetID() < source[j].GetID()
	})
}

// Equal check whether source DeviceSet and target DeviceSet is identical regardless of element order.
func (source DeviceSet) Equal(target DeviceSet) bool {
	if len(source) != len(target) {
		return false
	}

	visited := make(map[string]TopologyHintKey)
	for _, device := range source {
		visited[device.GetID()] = device.GetTopologyHintKey()
	}

	for _, device := range target {
		if visited[device.GetID()] != device.GetTopologyHintKey() {
			return false
		}
	}

	return true
}

// Difference returns a subset of the source DeviceSet that has no intersection with the target DeviceSet.
func (source DeviceSet) Difference(target DeviceSet) (difference DeviceSet) {
	for _, device := range source {
		if !target.Contains(DeviceSet{device}) {
			difference = append(difference, device)
		}
	}

	return difference
}

// Union returns new DeviceSet containing elements of source and target DeviceSets
func (source DeviceSet) Union(target DeviceSet) (union DeviceSet) {
	union = append(union, source...)
	visited := map[string]bool{}
	for _, device := range source {
		visited[device.GetID()] = true
	}

	for _, device := range target {
		if _, ok := visited[device.GetID()]; !ok {
			union = append(union, device)
		}
	}

	return union
}

// TODO(@bg): impl Intersection()
