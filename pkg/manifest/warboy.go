package manifest

import (
	"fmt"

	"github.com/bradfitz/iter"
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/smi"
	"path/filepath"
)

const (
	sysClassRoot         = "/sys/class/npu_mgmt/"
	sysDevicesRoot       = "/sys/devices/virtual/npu_mgmt/"
	mgmtFileExp          = "%s_mgmt"
	readOnlyOpt          = "ro"
	readWriteOpt         = "rw"
	channelExp           = "%sch%d"
	warboyMaxChannel int = 4
)

var _ Manifest = (*warboyManifest)(nil)

type warboyManifest struct {
	device      smi.Device
	deviceInfo  smi.DeviceInfo
	deviceFiles []smi.DeviceFile
}

func NewWarboyManifest(device smi.Device) (Manifest, error) {
	deviceInfo, err := device.DeviceInfo()
	if err != nil {
		return nil, err
	}

	deviceFiles, err := device.DeviceFiles()
	if err != nil {
		return nil, err
	}

	return &warboyManifest{
		device:      device,
		deviceInfo:  deviceInfo,
		deviceFiles: deviceFiles,
	}, nil
}

// EnvVars Note: older version of device plugin sets `NPU_DEVNAME`, `NPU_NPUNAME`, `NPU_PENAME`.
// However, those env variables are now deprecated and replaced with device-api.
func (w warboyManifest) EnvVars() map[string]string {
	return nil
}

// Annotations Note: order version of device plugin set the annotation `alpha.furiosa.ai/npu.devname`.
// This annotation is used for CRI Runtime injection, however the annotation was not consumed.
func (w warboyManifest) Annotations() map[string]string {
	return nil
}

func (w warboyManifest) DeviceNodes() []*DeviceNode {
	var deviceNodes []*DeviceNode
	devName := w.deviceInfo.Name()

	// mount npu mgmt file under "/dev"
	deviceNodes = append(deviceNodes, &DeviceNode{
		ContainerPath: fmt.Sprintf(mgmtFileExp, devName),
		HostPath:      fmt.Sprintf(mgmtFileExp, devName),
		Permissions:   readWriteOpt,
	})

	// mount devFiles such as "/dev/npu0", "/dev/npu0pe0"
	for _, file := range w.deviceFiles {
		deviceNodes = append(deviceNodes, &DeviceNode{
			ContainerPath: file.Path(),
			HostPath:      file.Path(),
			Permissions:   readWriteOpt,
		})
	}

	// mount channel fd for dma such as "/dev/npu0ch0" ~ "/dev/npu0ch3"
	for idx := range iter.N(warboyMaxChannel) {
		deviceNodes = append(deviceNodes, &DeviceNode{
			ContainerPath: fmt.Sprintf(channelExp, devName, idx),
			HostPath:      fmt.Sprintf(channelExp, devName, idx),
			Permissions:   readWriteOpt,
		})
	}

	return deviceNodes
}

func (w warboyManifest) MountPaths() []*Mount {
	var mounts []*Mount
	devName := filepath.Base(w.deviceInfo.Name())

	// mount "/sys/class/npu_mgmt/npu{x}_mgmt" path
	mounts = append(mounts, &Mount{
		ContainerPath: sysClassRoot + fmt.Sprintf(mgmtFileExp, devName),
		HostPath:      sysClassRoot + fmt.Sprintf(mgmtFileExp, devName),
		Options:       []string{readOnlyOpt},
	})

	// mount all dev files under "/sys/class/npu_mgmt/"
	for _, file := range w.deviceFiles {
		fileName := filepath.Base(file.Path())
		mounts = append(mounts, &Mount{
			ContainerPath: sysClassRoot + fileName,
			HostPath:      sysClassRoot + fileName,
			Options:       []string{readOnlyOpt},
		})
	}

	// mount "/sys/devices/virtual/npu_mgmt/npu{x}_mgmt" path
	mounts = append(mounts, &Mount{
		ContainerPath: sysDevicesRoot + fmt.Sprintf(mgmtFileExp, devName),
		HostPath:      sysDevicesRoot + fmt.Sprintf(mgmtFileExp, devName),
		Options:       []string{readOnlyOpt},
	})

	// mount all dev files under "/sys/devices/virtual/npu_mgmt/"
	for _, file := range w.deviceFiles {
		fileName := filepath.Base(file.Path())
		mounts = append(mounts, &Mount{
			ContainerPath: sysDevicesRoot + fileName,
			HostPath:      sysDevicesRoot + fileName,
			Options:       []string{readOnlyOpt},
		})
	}

	return mounts
}
