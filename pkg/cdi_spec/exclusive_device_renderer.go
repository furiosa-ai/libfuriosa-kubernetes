package cdi_spec

import (
	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"tags.cncf.io/container-device-interface/specs-go"
)

func NewExclusiveDeviceSpecRenderer(device smi.Device) (Renderer, error) {
	deviceSpec, err := newRngdDeviceSpec(device)
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
