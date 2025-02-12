package cdi_spec

import (
	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"tags.cncf.io/container-device-interface/specs-go"
)

func NewExclusiveDeviceSpecRenderer(device smi.Device) (Renderer, error) {
	deviceInfo, err := device.DeviceInfo()
	if err != nil {
		return nil, err
	}

	var deviceSpec CDISpec = nil
	switch deviceInfo.Arch() {
	case smi.ArchWarboy:
		deviceSpec, err = newWarboyDeviceSpec(device)
	case smi.ArchRngd:
		deviceSpec, err = newRngdDeviceSpec(device)
	}

	if err != nil {
		return nil, err
	}

	return &exclusiveDeviceSpecRenderer{
		spec: deviceSpec,
	}, nil
}

var _ Renderer = (*exclusiveDeviceSpecRenderer)(nil)

type exclusiveDeviceSpecRenderer struct {
	spec CDISpec
}

func (e *exclusiveDeviceSpecRenderer) Render() *specs.Device {
	return e.spec.DeviceSpec()
}
