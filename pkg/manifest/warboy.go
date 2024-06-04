package manifest

import (
	"fmt"

	"github.com/bradfitz/iter"
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/smi"
	"path/filepath"
)

const (
	warboyMgmtFileExp        = "%s_mgmt"
	warboySysClassRoot       = "/sys/class/npu_mgmt/"
	warboySysDevicesRoot     = "/sys/devices/virtual/npu_mgmt/"
	warboyMaxChannel     int = 4
	warboyChannelExp         = "%sch%d"
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
		ContainerPath: fmt.Sprintf(warboyMgmtFileExp, devName),
		HostPath:      fmt.Sprintf(warboyMgmtFileExp, devName),
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
			ContainerPath: fmt.Sprintf(warboyChannelExp, devName, idx),
			HostPath:      fmt.Sprintf(warboyChannelExp, devName, idx),
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
		ContainerPath: warboySysClassRoot + fmt.Sprintf(warboyMgmtFileExp, devName),
		HostPath:      warboySysClassRoot + fmt.Sprintf(warboyMgmtFileExp, devName),
		Options:       []string{readOnlyOpt},
	})

	// mount all dev files under "/sys/class/npu_mgmt/"
	for _, file := range w.deviceFiles {
		fileName := filepath.Base(file.Path())
		mounts = append(mounts, &Mount{
			ContainerPath: warboySysClassRoot + fileName,
			HostPath:      warboySysClassRoot + fileName,
			Options:       []string{readOnlyOpt},
		})
	}

	// mount "/sys/devices/virtual/npu_mgmt/npu{x}_mgmt" path
	mounts = append(mounts, &Mount{
		ContainerPath: warboySysDevicesRoot + fmt.Sprintf(warboyMgmtFileExp, devName),
		HostPath:      warboySysDevicesRoot + fmt.Sprintf(warboyMgmtFileExp, devName),
		Options:       []string{readOnlyOpt},
	})

	// mount all dev files under "/sys/devices/virtual/npu_mgmt/"
	for _, file := range w.deviceFiles {
		fileName := filepath.Base(file.Path())
		mounts = append(mounts, &Mount{
			ContainerPath: warboySysDevicesRoot + fileName,
			HostPath:      warboySysDevicesRoot + fileName,
			Options:       []string{readOnlyOpt},
		})
	}

	return mounts
}
