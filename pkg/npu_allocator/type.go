package npu_allocator

import (
	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/util"
)

type NpuAllocator interface {
	Allocate(available DeviceSet, required DeviceSet, size int) DeviceSet
}

type Device interface {
	// Index returns an index number of Device for sorting purpose.
	// Index must be injected from `furiosa-device-plugin`, and should not be modified by `libfuriosa-kubernetes`.
	Index() int

	// ID returns a unique ID of Device to identify the device.
	ID() string

	// TopologyHintKey returns unique key to retrieve TopologyHint using TopologyHintProvider.
	TopologyHintKey() TopologyHintKey

	// Equal check whether source Device is identical to the target Device.
	Equal(target Device) bool
}

type DeviceSet interface {
	Contains(target ...Device) bool
	Equal(target ...Device) bool
	Difference(target ...Device) DeviceSet
	Union(target ...Device) DeviceSet
	Insert(target ...Device)
	Devices() []Device
	Len() int
}

type deviceSet struct {
	btreeSet *util.BtreeSet[Device]
}

func NewDeviceSet(devices ...Device) DeviceSet {
	btreeSet := util.NewBtreeSetWithLessFunc(len(devices), func(a, b Device) bool {
		idx1, idx2 := a.Index(), b.Index()
		if idx1 == idx2 {
			id1, id2 := a.ID(), b.ID()
			return id1 < id2
		}

		return idx1 < idx2
	})

	for _, device := range devices {
		btreeSet.Insert(device)
	}

	return &deviceSet{btreeSet: btreeSet}
}

// Contains checks whether source DeviceSet contains target DeviceSet.
func (source *deviceSet) Contains(target ...Device) bool {
	if source.Len() == 0 || len(target) == 0 {
		return false
	}

	for _, targetDevice := range target {
		if !source.btreeSet.Has(targetDevice) {
			return false
		}
	}

	return true
}

// Equal check whether source DeviceSet and target DeviceSet is identical regardless of element order.
func (source *deviceSet) Equal(target ...Device) bool {
	if source.Len() == 0 && len(target) == 0 {
		return true
	}

	if source.Len() == 0 || len(target) == 0 {
		return false
	}

	if source.Len() != len(target) {
		return false
	}

	for _, targetDevice := range target {
		if !source.btreeSet.Has(targetDevice) {
			return false
		}
	}

	return true
}

// Difference returns a subset of the source DeviceSet that has no intersection with the target DeviceSet.
func (source *deviceSet) Difference(target ...Device) DeviceSet {
	targetDeviceSet := NewDeviceSet(target...)

	difference := NewDeviceSet()
	for _, sourceDevice := range source.Devices() {
		if !targetDeviceSet.Contains(sourceDevice) {
			difference.Insert(sourceDevice)
		}
	}

	return difference
}

// Union returns new DeviceSet containing elements of source and target DeviceSets
func (source *deviceSet) Union(target ...Device) DeviceSet {
	ds := NewDeviceSet(source.Devices()...)
	for _, targetDevice := range target {
		ds.Insert(targetDevice)
	}

	return ds
}

func (source *deviceSet) Insert(target ...Device) {
	for _, device := range target {
		source.btreeSet.Insert(device)
	}
}

func (source *deviceSet) Devices() []Device {
	return source.btreeSet.Items()
}

func (source *deviceSet) Len() int {
	return source.btreeSet.Len()
}

// TODO(@bg): impl Intersection()

// TopologyHintProvider takes two devices as argument return topology hint.
// The hint would be score, distance, preference of two devices.
type TopologyHintProvider func(device1, device2 Device) uint

// TopologyHintKey is named type of string, used for TopologyHintMatrix
type TopologyHintKey string

// TopologyHintMatrix provides score of device to device based on smi.Device smi.LinkType.
type TopologyHintMatrix map[TopologyHintKey]map[TopologyHintKey]uint

// TopologyScoreCalculator calculates sum of score of given topologyHintKeys based on smi.Device smi.LinkType.
type TopologyScoreCalculator func(keys []TopologyHintKey) uint

// NewTopologyHintMatrix generates TopologyHintMatrix using list of smi.Device.
func NewTopologyHintMatrix(smiDevices []smi.Device) (TopologyHintMatrix, error) {
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
			linkType, err := device1.DeviceToDeviceLinkType(device2)
			if err != nil {
				return nil, err
			}

			score := uint(linkType)

			pciBusID1, err := util.ParseBusIDFromBDF(deviceInfo1.BDF())
			if err != nil {
				return nil, err
			}

			pciBusID2, err := util.ParseBusIDFromBDF(deviceInfo2.BDF())
			if err != nil {
				return nil, err
			}

			key1, key2 := TopologyHintKey(pciBusID1), TopologyHintKey(pciBusID2)
			if key1 > key2 {
				key1, key2 = key2, key1
			}

			if _, ok := topologyHintMatrix[key1]; !ok {
				topologyHintMatrix[key1] = make(map[TopologyHintKey]uint)
			}

			topologyHintMatrix[key1][key2] = score
		}
	}

	return topologyHintMatrix, nil
}
