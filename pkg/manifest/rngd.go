package manifest

import (
	"fmt"

	"github.com/bradfitz/iter"
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/smi"
	"path/filepath"
)

const (
	devFsRngdMgmtFileExp = "/dev/rngd/%smgmt"
	sysFsRngdMgmtFileExp = "rngd!%smgmt"
	rngdSysClassRoot     = "/sys/class/rngd_mgmt/"
	rngdSysDevicesRoot   = "/sys/devices/virtual/rngd_mgmt/"
	rngdMaxChannel       = 8
	rngdChannelExp       = "/dev/rngd/%sch%d"
)

var _ Manifest = (*rngdManifest)(nil)

type rngdManifest struct {
	device      smi.Device
	deviceInfo  smi.DeviceInfo
	deviceFiles []smi.DeviceFile
}

func NewRngdManifest(device smi.Device) (Manifest, error) {
	deviceInfo, err := device.DeviceInfo()
	if err != nil {
		return nil, err
	}

	deviceFiles, err := device.DeviceFiles()
	if err != nil {
		return nil, err
	}

	return &rngdManifest{
		device:      device,
		deviceInfo:  deviceInfo,
		deviceFiles: deviceFiles,
	}, nil
}

// EnvVars Note: older version of device plugin sets `NPU_DEVNAME`, `NPU_NPUNAME`, `NPU_PENAME`.
// However, those env variables are now deprecated and replaced with device-api.
func (w rngdManifest) EnvVars() map[string]string {
	return nil
}

// Annotations Note: order version of device plugin set the annotation `alpha.furiosa.ai/npu.devname`.
// This annotation is used for CRI Runtime injection, however the annotation was not consumed.
func (w rngdManifest) Annotations() map[string]string {
	return nil
}

func (w rngdManifest) DeviceNodes() []*DeviceNode {
	var deviceNodes []*DeviceNode
	devName := filepath.Base(w.deviceInfo.Name())

	// mount npu mgmt file under "/dev/rngd"
	deviceNodes = append(deviceNodes, &DeviceNode{
		ContainerPath: fmt.Sprintf(devFsRngdMgmtFileExp, devName),
		HostPath:      fmt.Sprintf(devFsRngdMgmtFileExp, devName),
		Permissions:   readWriteOpt,
	})

	// mount devFiles such as "/dev/rngd/npu0pe0", "/dev/rngd/npu0pe0-1"
	for _, file := range w.deviceFiles {
		deviceNodes = append(deviceNodes, &DeviceNode{
			ContainerPath: file.Path(),
			HostPath:      file.Path(),
			Permissions:   readWriteOpt,
		})
	}

	// mount channel fd for dma such as "/dev/rngd/npu0ch0" ~ "/dev/rngd/npu0ch7"
	for idx := range iter.N(rngdMaxChannel) {
		deviceNodes = append(deviceNodes, &DeviceNode{
			ContainerPath: fmt.Sprintf(rngdChannelExp, devName, idx),
			HostPath:      fmt.Sprintf(rngdChannelExp, devName, idx),
			Permissions:   readWriteOpt,
		})
	}

	return deviceNodes
}

func (w rngdManifest) MountPaths() []*Mount {
	var mounts []*Mount
	devName := filepath.Base(w.deviceInfo.Name())

	// mount "/sys/class/rngd_mgmt/rngd!npu{x}_mgmt" path
	mounts = append(mounts, &Mount{
		ContainerPath: rngdSysClassRoot + fmt.Sprintf(sysFsRngdMgmtFileExp, devName),
		HostPath:      rngdSysClassRoot + fmt.Sprintf(sysFsRngdMgmtFileExp, devName),
		Options:       []string{readOnlyOpt},
	})

	// mount "/sys/devices/virtual/rngd_mgmt/rngd!npu{x}_mgmt" path
	mounts = append(mounts, &Mount{
		ContainerPath: rngdSysDevicesRoot + fmt.Sprintf(sysFsRngdMgmtFileExp, devName),
		HostPath:      rngdSysDevicesRoot + fmt.Sprintf(sysFsRngdMgmtFileExp, devName),
		Options:       []string{readOnlyOpt},
	})

	return mounts
}
