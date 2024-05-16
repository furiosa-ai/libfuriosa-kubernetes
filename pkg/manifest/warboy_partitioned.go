package manifest

import (
	"fmt"
	"github.com/bradfitz/iter"
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/device"
	"tags.cncf.io/container-device-interface/pkg/cdi"
)

var _ Manifest = (*partitionedWarboyManifest)(nil)

type partitionedWarboyManifest struct {
	device    device.Device
	coreStart uint8
	coreEnd   uint8
}

func NewPartitionedWarboyManifest(origin device.Device, coreStart uint8, coreEnd uint8) Manifest {
	return &partitionedWarboyManifest{
		device:    origin,
		coreStart: coreStart,
		coreEnd:   coreEnd,
	}
}

func (p partitionedWarboyManifest) EnvVars() map[string]string {
	return nil
}

func (p partitionedWarboyManifest) Annotations() map[string]string {
	return nil
}

func (p partitionedWarboyManifest) DeviceNodes() []*DeviceNode {
	var deviceNodes []*DeviceNode

	deviceNodes = append(deviceNodes, &DeviceNode{
		ContainerPath: devRoot + fmt.Sprintf(mgmtFileExp, p.device.Name()),
		HostPath:      devRoot + fmt.Sprintf(mgmtFileExp, p.device.Name()),
		Permissions:   readWriteOpt,
	})

	for _, file := range collectDevFiles(p.device.DevFiles(), p.coreStart, p.coreEnd) {
		deviceNodes = append(deviceNodes, &DeviceNode{
			ContainerPath: file.Path(),
			HostPath:      file.Path(),
			Permissions:   readWriteOpt,
		})
	}

	for idx := range iter.N(warboyMaxChannel) {
		deviceNodes = append(deviceNodes, &DeviceNode{
			ContainerPath: fmt.Sprintf(channelExp, p.device.Name(), idx),
			HostPath:      fmt.Sprintf(channelExp, p.device.Name(), idx),
			Permissions:   readWriteOpt,
		})
	}

	return deviceNodes
}

func (p partitionedWarboyManifest) MountPaths() []*Mount {
	var mounts []*Mount
	devName := p.device.Name()

	filteredDeviceFiles := collectDevFiles(p.device.DevFiles(), p.coreStart, p.coreEnd)

	// mount "/sys/class/npu_mgmt/npu{x}_mgmt" path
	mounts = append(mounts, &Mount{
		ContainerPath: sysClassRoot + fmt.Sprintf(mgmtFileExp, devName),
		HostPath:      sysClassRoot + fmt.Sprintf(mgmtFileExp, devName),
		Options:       []string{readOnlyOpt},
	})

	// mount all dev files under "/sys/class/npu_mgmt/"
	for _, file := range filteredDeviceFiles {
		mounts = append(mounts, &Mount{
			ContainerPath: sysClassRoot + file.Filename(),
			HostPath:      sysClassRoot + file.Filename(),
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
	for _, file := range filteredDeviceFiles {
		mounts = append(mounts, &Mount{
			ContainerPath: sysDevicesRoot + file.Filename(),
			HostPath:      sysDevicesRoot + file.Filename(),
			Options:       []string{readOnlyOpt},
		})
	}

	return mounts
}

func (p partitionedWarboyManifest) ToCDIContainerEdits() *cdi.ContainerEdits {
	return toCDIContainerEdits(p)
}
