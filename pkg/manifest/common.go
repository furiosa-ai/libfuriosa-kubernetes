package manifest

import (
	"fmt"
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/device"
	"os"
	"tags.cncf.io/container-device-interface/pkg/cdi"
	"tags.cncf.io/container-device-interface/specs-go"
)

const charDeviceType = "c"
const sysfsMountType = "sysfs"

func toCDIContainerEdits(manifest Manifest) *cdi.ContainerEdits {
	edit := cdi.ContainerEdits{ContainerEdits: &specs.ContainerEdits{}}
	envVars := manifest.EnvVars()

	for k, v := range envVars {
		edit.ContainerEdits.Env = append(edit.ContainerEdits.Env, fmt.Sprintf("%s=%s", k, v))
	}

	// skit annotations in ContainerEdits scope
	for _, devNode := range manifest.DeviceNodes() {
		edit.ContainerEdits.DeviceNodes = append(edit.ContainerEdits.DeviceNodes, &specs.DeviceNode{
			Path:        devNode.ContainerPath,
			HostPath:    devNode.HostPath,
			Type:        charDeviceType,
			Major:       0,
			Minor:       0,
			FileMode:    func() *os.FileMode { mode := os.ModeCharDevice; return &mode }(),
			Permissions: devNode.Permissions,
			UID:         nil,
			GID:         nil,
		})
	}

	for _, mounts := range manifest.MountPaths() {
		edit.ContainerEdits.Mounts = append(edit.ContainerEdits.Mounts, &specs.Mount{
			HostPath:      mounts.HostPath,
			ContainerPath: mounts.ContainerPath,
			Options:       mounts.Options,
			Type:          sysfsMountType,
		})
	}

	return &edit
}

// collectDevFiles collects device files corresponding to the given core range.
func collectDevFiles(deviceFiles []device.DeviceFile, coreStart uint8, coreEnd uint8) []device.DeviceFile {
	var result []device.DeviceFile

	for _, deviceFile := range deviceFiles {
		if deviceFile.CoreRange().Type() == device.CoreRangeTypeRange {
			if deviceFile.CoreRange().Start() >= coreStart && deviceFile.CoreRange().End() <= coreEnd {
				result = append(result, deviceFile)
			}
		}
	}

	return result
}
