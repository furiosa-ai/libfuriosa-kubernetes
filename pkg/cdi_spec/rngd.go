package cdi_spec

import (
	"fmt"
	"github.com/bradfitz/iter"
	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"tags.cncf.io/container-device-interface/specs-go"
)

const (
	rngdDevFsRoot        = "/dev/rngd/"
	rngdMgmtFileExp      = "%smgmt"
	rngdMaxChannel       = 8
	rngdMaxRemoteChannel = 8
	rngdChannelExp       = "%sch%d"
	rngdRemoteChannelExp = "%sch%dr"
	rngdDmaRemappingExp  = "%sdmar"
	rngdBar0Exp          = "%sbar0"
	rngdBar2Exp          = "%sbar2"
	rngdBar4Exp          = "%sbar4"
)

type rngdDeviceSpec struct {
	device      smi.Device
	deviceInfo  smi.DeviceInfo
	deviceFiles []smi.DeviceFile
}

func newRngdDeviceSpec(device smi.Device) (CDISpec, error) {
	deviceInfo, err := device.DeviceInfo()
	if err != nil {
		return nil, err
	}

	deviceFiles, err := device.DeviceFiles()
	if err != nil {
		return nil, err
	}

	return &rngdDeviceSpec{
		device:      device,
		deviceInfo:  deviceInfo,
		deviceFiles: deviceFiles,
	}, nil
}

func (w *rngdDeviceSpec) containerEdits() *specs.ContainerEdits {
	return &specs.ContainerEdits{
		Env:            nil,
		DeviceNodes:    w.deviceNodes(),
		Hooks:          nil,
		Mounts:         w.mounts(),
		IntelRdt:       nil,
		AdditionalGIDs: nil,
	}
}

func (w *rngdDeviceSpec) DeviceSpec() *specs.Device {
	containerEdits := w.containerEdits()

	return &specs.Device{
		Name:           w.deviceInfo.Name(),
		ContainerEdits: *containerEdits,
	}
}

func (w *rngdDeviceSpec) deviceNodes() []*specs.DeviceNode {
	var deviceNodes []*specs.DeviceNode
	devName := w.deviceInfo.Name()

	// mount npu mgmt deviceFile under "/dev/rngd"
	deviceNodes = append(deviceNodes, &specs.DeviceNode{
		Path:        rngdDevFsRoot + fmt.Sprintf(rngdMgmtFileExp, devName),
		HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdMgmtFileExp, devName),
		Permissions: readWriteOpt,
	})

	// mount devFiles such as "/dev/rngd/npu0pe0", "/dev/rngd/npu0pe0-1"
	for _, deviceFile := range w.deviceFiles {
		deviceNodes = append(deviceNodes, &specs.DeviceNode{
			Path:        deviceFile.Path(),
			HostPath:    deviceFile.Path(),
			Permissions: readWriteOpt,
		})
	}

	// mount channel fd for dma such as "/dev/rngd/npu0ch0" ~ "/dev/rngd/npu0ch7"
	for idx := range iter.N(rngdMaxChannel) {
		deviceNodes = append(deviceNodes, &specs.DeviceNode{
			Path:        rngdDevFsRoot + fmt.Sprintf(rngdChannelExp, devName, idx),
			HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdChannelExp, devName, idx),
			Permissions: readWriteOpt,
		})
	}

	// mount remote channel fd for dma such as "/dev/rngd/npu0ch0r" ~ "/dev/rngd/npu0ch7r"
	for idx := range iter.N(rngdMaxRemoteChannel) {
		deviceNodes = append(deviceNodes, &specs.DeviceNode{
			Path:        rngdDevFsRoot + fmt.Sprintf(rngdRemoteChannelExp, devName, idx),
			HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdRemoteChannelExp, devName, idx),
			Permissions: readWriteOpt,
		})
	}

	// mount dma remapping fd
	deviceNodes = append(deviceNodes, &specs.DeviceNode{
		Path:        rngdDevFsRoot + fmt.Sprintf(rngdDmaRemappingExp, devName),
		HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdDmaRemappingExp, devName),
		Permissions: readWriteOpt,
	})

	// mount bar0
	deviceNodes = append(deviceNodes, &specs.DeviceNode{
		Path:        rngdDevFsRoot + fmt.Sprintf(rngdBar0Exp, devName),
		HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdBar0Exp, devName),
		Permissions: readWriteOpt,
	})

	// mount bar2
	deviceNodes = append(deviceNodes, &specs.DeviceNode{
		Path:        rngdDevFsRoot + fmt.Sprintf(rngdBar2Exp, devName),
		HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdBar2Exp, devName),
		Permissions: readWriteOpt,
	})

	// mount bar4
	deviceNodes = append(deviceNodes, &specs.DeviceNode{
		Path:        rngdDevFsRoot + fmt.Sprintf(rngdBar4Exp, devName),
		HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdBar4Exp, devName),
		Permissions: readWriteOpt,
	})

	return deviceNodes
}

func (w *rngdDeviceSpec) mounts() []*specs.Mount {
	return nil
}
