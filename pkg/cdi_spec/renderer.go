package cdi_spec

import (
	"tags.cncf.io/container-device-interface/specs-go"
)

type Renderer interface {
	Render() *specs.Device
}

type CDISpec interface {
	DeviceSpec() *specs.Device
	containerEdits() *specs.ContainerEdits
	deviceNodes() []*specs.DeviceNode
	mounts() []*specs.Mount
}
