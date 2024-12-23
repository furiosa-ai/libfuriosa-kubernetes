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
	warboySysDevicesRoot       = "/sys/devices/virtual/npu_mgmt/"
	warboyMaxChannel       int = 4
	warboyChannelExp           = "/dev/%sch%d"

	warboySinglePeExp = "pe%d"
	warboyDualPeExp   = "pe%d-%d"
	warboyBar0Exp     = "%s_bar0"
	warboyBar2Exp     = "%s_bar2"
	warboyBar4Exp     = "%s_bar4"
)

type warboyDeviceSpec struct {
	device      smi.Device
	deviceInfo  smi.DeviceInfo
	deviceFiles []smi.DeviceFile
}

func newWarboyDeviceSpec(device smi.Device) (DeviceSpec, error) {
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
		Mounts:         nil,
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
	var mounts []*specs.Mount
	devName := w.deviceInfo.Name()

	// mount "/sys/devices/virtual/npu_mgmt/npu{x}" path
	mounts = append(mounts, &specs.Mount{
		HostPath:      warboySysDevicesRoot + devName,
		ContainerPath: warboySysDevicesRoot + devName,
		Options:       []string{readOnlyOpt, bindOpt},
	})

	// mount /sys/devices/virtual/npu_mgmt/npu{x}pe0 path
	// mount /sys/devices/virtual/npu_mgmt/npu{x}pe1 path
	for idx := range iter.N(2) {
		mounts = append(mounts, &specs.Mount{
			ContainerPath: warboySysDevicesRoot + devName + fmt.Sprintf(warboySinglePeExp, idx),
			HostPath:      warboySysDevicesRoot + devName + fmt.Sprintf(warboySinglePeExp, idx),
			Options:       []string{readOnlyOpt, bindOpt},
		})
	}

	// mount /sys/devices/virtual/npu_mgmt/npu{x}pe0-1 path
	mounts = append(mounts, &specs.Mount{
		ContainerPath: warboySysDevicesRoot + devName + fmt.Sprintf(warboyDualPeExp, 0, 1),
		HostPath:      warboySysDevicesRoot + devName + fmt.Sprintf(warboyDualPeExp, 0, 1),
		Options:       []string{readOnlyOpt, bindOpt},
	})

	// mount /sys/devices/virtual/npu_mgmt/npu{x}_mgmt path
	mounts = append(mounts, &specs.Mount{
		HostPath:      warboySysDevicesRoot + fmt.Sprintf(warboyMgmtFileExp, devName),
		ContainerPath: warboySysDevicesRoot + fmt.Sprintf(warboyMgmtFileExp, devName),
		Options:       []string{readOnlyOpt, bindOpt},
	})

	// mount /sys/devices/virtual/npu_mgmt/npu{x}_bar0 path
	mounts = append(mounts, &specs.Mount{
		ContainerPath: warboySysDevicesRoot + fmt.Sprintf(warboyBar0Exp, devName),
		HostPath:      warboySysDevicesRoot + fmt.Sprintf(warboyBar0Exp, devName),
		Options:       []string{readOnlyOpt, bindOpt},
	})

	// mount /sys/devices/virtual/npu_mgmt/npu{x}_bar2 path
	mounts = append(mounts, &specs.Mount{
		ContainerPath: warboySysDevicesRoot + fmt.Sprintf(warboyBar2Exp, devName),
		HostPath:      warboySysDevicesRoot + fmt.Sprintf(warboyBar2Exp, devName),
		Options:       []string{readOnlyOpt, bindOpt},
	})

	// mount /sys/devices/virtual/npu_mgmt/npu{x}_bar4 path
	mounts = append(mounts, &specs.Mount{
		ContainerPath: warboySysDevicesRoot + fmt.Sprintf(warboyBar4Exp, devName),
		HostPath:      warboySysDevicesRoot + fmt.Sprintf(warboyBar4Exp, devName),
		Options:       []string{readOnlyOpt, bindOpt},
	})

	return mounts
}
