package furiosa_device

import (
	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/cdi_spec"
	"tags.cncf.io/container-device-interface/specs-go"
)

var _ FuriosaDevice = (*exclusiveDevice)(nil)

type exclusiveDevice struct {
	index      int
	origin     smi.Device
	manifest   cdi_spec.Renderer
	deviceID   string
	pciBusID   string
	numaNode   int
	isDisabled bool
}

func newExclusiveDevice(originDevice smi.Device, isDisabled bool) (FuriosaDevice, error) {
	deviceID, pciBusID, numaNode, originIndex, err := parseDeviceInfo(originDevice)
	if err != nil {
		return nil, err
	}

	newExclusiveDeviceManifest, err := cdi_spec.NewExclusiveDeviceSpecRenderer(originDevice)
	if err != nil {
		return nil, err
	}

	return &exclusiveDevice{
		index:      originIndex,
		origin:     originDevice,
		manifest:   newExclusiveDeviceManifest,
		deviceID:   deviceID,
		pciBusID:   pciBusID,
		numaNode:   int(numaNode),
		isDisabled: isDisabled,
	}, nil
}

func (f *exclusiveDevice) DeviceID() string {
	return f.deviceID
}

func (f *exclusiveDevice) PCIBusID() string {
	return f.pciBusID
}

func (f *exclusiveDevice) NUMANode() int {
	return f.numaNode
}

func (f *exclusiveDevice) IsHealthy() (bool, error) {
	//TODO(@bg): use more sophisticated way
	if f.isDisabled {
		return false, nil
	}
	liveness, err := f.origin.Liveness()
	if err != nil {
		return liveness, err
	}
	return liveness, nil
}

func (f *exclusiveDevice) DeviceSpec() *specs.Device {
	renderer, _ := cdi_spec.NewExclusiveDeviceSpecRenderer(f.origin)
	return renderer.Render()
}

func (f *exclusiveDevice) Index() int {
	return f.index
}
