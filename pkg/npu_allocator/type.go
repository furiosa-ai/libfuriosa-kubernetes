package npu_allocator

import (
	"sort"

	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/smi"
)

// TopologyHintProvider takes two devices as argument return topology hint.
// The hint would be score, distance, preference of two devices.
type TopologyHintProvider func(device1, device2 Device) uint

// TopologyHintMatrix provides score of each device to device based on device link type.
type TopologyHintMatrix map[string]map[string]uint

// populateTopologyMatrix generates TopologyHintMatrix using provided list of smi.Device.
func populateTopologyMatrix(smiDevices []smi.Device) (TopologyHintMatrix, error) {
	topologyHintMatrix := make(TopologyHintMatrix)
	deviceToDeviceInfo := make(map[smi.Device]smi.DeviceInfo)

	for _, device := range smiDevices {
		deviceInfo, err := device.DeviceInfo()
		if err != nil {
			return nil, err
		}
		deviceToDeviceInfo[device] = deviceInfo
	}

	for device1, deviceInfo1 := range deviceToDeviceInfo {
		for device2, deviceInfo2 := range deviceToDeviceInfo {
			linkType, err := device1.GetDeviceToDeviceLinkType(device2)
			if err != nil {
				return nil, err
			}

			key1 := deviceInfo1.BDF()
			key2 := deviceInfo2.BDF()
			if key1 > key2 {
				key1, key2 = key2, key1
			}

			if _, ok := topologyHintMatrix[key1]; !ok {
				topologyHintMatrix[key1] = make(map[string]uint)
			}

			topologyHintMatrix[key1][key2] = uint(linkType)
		}
	}

	return topologyHintMatrix, nil
}

type NpuAllocator interface {
	Allocate(available DeviceSet, required DeviceSet, size int) DeviceSet
}

type Device interface {
	// ID returns a unique ID of Device to identify the device.
	ID() string
	// TopologyHintKey returns unique key to retrieve TopologyHint using TopologyHintProvider.
	TopologyHintKey() string
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
		visited[device.ID()] = true
	}

	for _, device := range target {
		if _, ok := visited[device.ID()]; !ok {
			return false
		}
	}

	return true
}

// Sort sorts source DeviceSet.
func (source DeviceSet) Sort() {
	sort.Slice(source, func(i, j int) bool {
		return source[i].ID() < source[j].ID()
	})
}

// Equal check whether source DeviceSet and target DeviceSet is identical regardless of element order.
func (source DeviceSet) Equal(target DeviceSet) bool {
	if len(source) != len(target) {
		return false
	}

	visited := make(map[string]string)
	for _, device := range source {
		visited[device.ID()] = device.TopologyHintKey()
	}

	for _, device := range target {
		if visited[device.ID()] != device.TopologyHintKey() {
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
		visited[device.ID()] = true
	}

	for _, device := range target {
		if _, ok := visited[device.ID()]; !ok {
			union = append(union, device)
		}
	}

	return union
}

// TODO(@bg): impl Intersection()
