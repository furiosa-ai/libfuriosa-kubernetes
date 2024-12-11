package furiosa_device

import (
	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/npu_allocator"
	devicePluginAPIv1Beta1 "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

const (
	readOnlyOpt = "ro"
)

type DeviceInfo interface {
	DeviceID() string
	PCIBusID() string
	NUMANode() int
	IsHealthy() (bool, error)
	IsExclusiveDevice() bool
}

type Manifest interface {
	EnvVars() map[string]string
	Annotations() map[string]string
	DeviceSpecs() []*devicePluginAPIv1Beta1.DeviceSpec
	Mounts() []*devicePluginAPIv1Beta1.Mount
	CDIDevices() []*devicePluginAPIv1Beta1.CDIDevice
}

type FuriosaDevice interface {
	DeviceInfo
	Manifest
	npu_allocator.Device
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
