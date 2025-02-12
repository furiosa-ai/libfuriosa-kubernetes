package cdi_spec_gen

import (
	"fmt"
	"path/filepath"
	"regexp"
	"tags.cncf.io/container-device-interface/specs-go"
)

const (
	basePattern = "npu(\\d+)"
)

var (
	baseRegex = regexp.MustCompile(basePattern)
)

func mutateContainerPath(origin *specs.Device, newDeviceIdx int) *specs.Device {
	mutatedSpec := origin

	// mutate device nodes
	for idx, deviceNode := range mutatedSpec.ContainerEdits.DeviceNodes {
		dir := filepath.Dir(deviceNode.Path)
		base := filepath.Base(deviceNode.Path)

		newBase := baseRegex.ReplaceAllString(base, fmt.Sprintf("npu%d", newDeviceIdx))
		mutatedSpec.ContainerEdits.DeviceNodes[idx].Path = filepath.Join(dir, newBase)
	}

	// mutate mounts
	for idx, mount := range origin.ContainerEdits.Mounts {
		dir := filepath.Dir(mount.ContainerPath)
		base := filepath.Base(mount.ContainerPath)

		newBase := baseRegex.ReplaceAllString(base, fmt.Sprintf("npu%d", newDeviceIdx))
		mutatedSpec.ContainerEdits.Mounts[idx].ContainerPath = filepath.Join(dir, newBase)
	}

	return mutatedSpec
}
