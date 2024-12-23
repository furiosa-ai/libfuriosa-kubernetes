package furiosa_device

import (
	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"tags.cncf.io/container-device-interface/specs-go"
)

type FuriosaDevice interface {
	Index() int
	DeviceID() string
	PCIBusID() string
	NUMANode() int
	IsHealthy() (bool, error)
	DeviceSpec() *specs.Device
}

func NewFuriosaDevices(devices []smi.Device, blockedList []string, policy PartitioningPolicy) ([]FuriosaDevice, error) {
	var furiosaDevices []FuriosaDevice
	var newDevFunc = newDeviceFuncResolver(policy)
	for _, origin := range devices {
		info, err := origin.DeviceInfo()
		if err != nil {
			return nil, err
		}

		isDisabled := contains(blockedList, info.UUID())
		newDevices, err := newDevFunc(origin, isDisabled)
		if err != nil {
			return nil, err
		}

		furiosaDevices = append(furiosaDevices, newDevices...)
	}
	return furiosaDevices, nil
}
