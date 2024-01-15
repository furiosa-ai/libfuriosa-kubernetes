package manifest

import (
	"fmt"

	"github.com/bradfitz/iter"
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/device"
)

const (
	sysClassRoot         = "/sys/class/npu_mgmt/"
	sysDevicesRoot       = "/sys/devices/virtual/npu_mgmt/"
	devRoot              = "/dev/"
	mgmtFileExp          = "%s_mgmt"
	readOnlyOpt          = "ro"
	readWriteOpt         = "rw"
	channelExp           = devRoot + "%sch%d"
	warboyMaxChannel int = 4
)

var _ Manifest = (*warboyMenifest)(nil)

type warboyMenifest struct {
	device device.Device
}

func NewWarboyManifest(origin device.Device) Manifest {
	return &warboyMenifest{
		device: origin,
	}
}

// EnvVars Note: older version of device plugin sets `NPU_DEVNAME`, `NPU_NPUNAME`, `NPU_PENAME`.
// However, those env variables are now deprecated and replaced with device-api.
func (w warboyMenifest) EnvVars() map[string]string {
	return nil
}

// Annotations Note: order version of device plugin set the annotation `alpha.furiosa.ai/npu.devname`.
// This annotation is used for CRI Runtime injection, however the annotation was not consumed.
func (w warboyMenifest) Annotations() map[string]string {
	return nil

}

func (w warboyMenifest) DeviceNodes() []DeviceNode {
	var deviceNodes []DeviceNode

	deviceNodes = append(deviceNodes, DeviceNode{
		ContainerPath: devRoot + fmt.Sprintf(mgmtFileExp, w.device.Name()),
		HostPath:      devRoot + fmt.Sprintf(mgmtFileExp, w.device.Name()),
		permissions:   readWriteOpt,
	})

	for _, file := range w.device.DevFiles() {
		deviceNodes = append(deviceNodes, DeviceNode{
			ContainerPath: file.Path(),
			HostPath:      file.Path(),
			permissions:   readWriteOpt,
		})
	}

	for idx := range iter.N(warboyMaxChannel) {
		deviceNodes = append(deviceNodes, DeviceNode{
			ContainerPath: fmt.Sprintf(channelExp, w.device.Name(), idx),
			HostPath:      fmt.Sprintf(channelExp, w.device.Name(), idx),
			permissions:   readWriteOpt,
		})
	}

	return deviceNodes
}

func (w warboyMenifest) MountPaths() []Mount {
	var mounts []Mount
	devName := w.device.Name()

	// mount "/sys/class/npu_mgmt/npu{x}_mgmt" path
	mounts = append(mounts, Mount{
		ContainerPath: sysClassRoot + fmt.Sprintf(mgmtFileExp, devName),
		HostPath:      sysClassRoot + fmt.Sprintf(mgmtFileExp, devName),
		Options:       []string{readOnlyOpt},
	})

	// mount all dev files under "/sys/class/npu_mgmt/"
	for _, file := range w.device.DevFiles() {
		mounts = append(mounts, Mount{
			ContainerPath: sysClassRoot + file.Filename(),
			HostPath:      sysClassRoot + file.Filename(),
			Options:       []string{readOnlyOpt},
		})
	}

	// mount "/sys/devices/virtual/npu_mgmt/npu{x}_mgmt" path
	mounts = append(mounts, Mount{
		ContainerPath: sysDevicesRoot + fmt.Sprintf(mgmtFileExp, devName),
		HostPath:      sysDevicesRoot + fmt.Sprintf(mgmtFileExp, devName),
		Options:       []string{readOnlyOpt},
	})

	// mount all dev files under "/sys/devices/virtual/npu_mgmt/"
	for _, file := range w.device.DevFiles() {
		mounts = append(mounts, Mount{
			ContainerPath: sysDevicesRoot + file.Filename(),
			HostPath:      sysDevicesRoot + file.Filename(),
			Options:       []string{readOnlyOpt},
		})
	}

	return mounts
}
