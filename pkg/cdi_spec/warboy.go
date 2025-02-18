package cdi_spec

import (
	"fmt"
	"tags.cncf.io/container-device-interface/specs-go"

	"github.com/bradfitz/iter"
	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
)

const (
	devFsWarboyMgmtFileExp     = "/dev/%s_mgmt"
	warboyMgmtFileExp          = "%s_mgmt"
	warboyMaxChannel       int = 4
	warboyChannelExp           = "/dev/%sch%d"
)

type warboyDeviceSpec struct {
	device      smi.Device
	deviceInfo  smi.DeviceInfo
	deviceFiles []smi.DeviceFile
}

func newWarboyDeviceSpec(device smi.Device) (CDISpec, error) {
	deviceInfo, err := device.DeviceInfo()
	if err != nil {
		return nil, err
	}

	deviceFiles, err := device.DeviceFiles()
	if err != nil {
		return nil, err
	}

	return &warboyDeviceSpec{
		device:      device,
		deviceInfo:  deviceInfo,
		deviceFiles: deviceFiles,
	}, nil
}

func (w *warboyDeviceSpec) containerEdits() *specs.ContainerEdits {
	return &specs.ContainerEdits{
		Env:            nil,
		DeviceNodes:    w.deviceNodes(),
		Hooks:          nil,
		Mounts:         w.mounts(),
		IntelRdt:       nil,
		AdditionalGIDs: nil,
	}
}

func (w *warboyDeviceSpec) DeviceSpec() *specs.Device {
	containerEdits := w.containerEdits()

	return &specs.Device{
		Name:           w.deviceInfo.Name(),
		ContainerEdits: *containerEdits,
	}
}

func (w *warboyDeviceSpec) deviceNodes() []*specs.DeviceNode {
	var deviceNodes []*specs.DeviceNode
	devName := w.deviceInfo.Name()

	// mount npu mgmt file under "/dev"
	deviceNodes = append(deviceNodes, &specs.DeviceNode{
		Path:        fmt.Sprintf(devFsWarboyMgmtFileExp, devName),
		HostPath:    fmt.Sprintf(devFsWarboyMgmtFileExp, devName),
		Permissions: readWriteOpt,
	})

	// mount devFiles such as "/dev/npu0", "/dev/npu0pe0"
	for _, file := range w.deviceFiles {
		deviceNodes = append(deviceNodes, &specs.DeviceNode{
			Path:        file.Path(),
			HostPath:    file.Path(),
			Permissions: readWriteOpt,
		})
	}

	// mount channel fd for dma such as "/dev/npu0ch0" ~ "/dev/npu0ch3"
	for idx := range iter.N(warboyMaxChannel) {
		deviceNodes = append(deviceNodes, &specs.DeviceNode{
			Path:        fmt.Sprintf(warboyChannelExp, devName, idx),
			HostPath:    fmt.Sprintf(warboyChannelExp, devName, idx),
			Permissions: readWriteOpt,
		})
	}

	return deviceNodes
}

func (w *warboyDeviceSpec) mounts() []*specs.Mount {
	return nil
}
