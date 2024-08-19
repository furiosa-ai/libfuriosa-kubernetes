package npu_allocator

import (
	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/util"
)

// TopologyHintProvider takes two devices as argument return topology hint.
// The hint would be score, distance, preference of two devices.
type TopologyHintProvider func(device1, device2 Device) uint

// TopologyHintKey is named type of string, used for TopologyHintMatrix
type TopologyHintKey string

// TopologyHintMatrix provides score of device to device based on smi.Device smi.LinkType.
type TopologyHintMatrix map[TopologyHintKey]map[TopologyHintKey]uint

// populateTopologyHintMatrix generates TopologyHintMatrix using list of smi.Device.
func populateTopologyHintMatrix(smiDevices []smi.Device) (TopologyHintMatrix, error) {
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

func getGenericHintProvider(topologyHintMatrix TopologyHintMatrix) TopologyHintProvider {
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
