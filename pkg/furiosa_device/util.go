package furiosa_device

import (
	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/util"
)

func parseDeviceInfo(originDevice smi.Device) (arch smi.Arch, deviceID, pciBusID string, numaNode uint, originIndex int, err error) {
	info, err := originDevice.DeviceInfo()
	if err != nil {
		return 0, "", "", 0, 0, err
	}

	arch = info.Arch()
	deviceID = info.UUID()
	pciBusID, err = util.ParseBusIDFromBDF(info.BDF())
	numaNode = uint(info.NumaNode())
	originIndex = int(info.Index())

	return arch, deviceID, pciBusID, numaNode, originIndex, err
}

func contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
