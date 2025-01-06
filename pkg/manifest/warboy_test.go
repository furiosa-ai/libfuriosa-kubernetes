package manifest

import (
	"fmt"
	"sort"
	"testing"

	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"github.com/stretchr/testify/assert"
)

func newTestWarboyDevice() smi.Device {
	return smi.GetStaticMockDevice(smi.ArchWarboy, 0)
}

func TestWarboyDeviceNodes(t *testing.T) {
	tests := []struct {
		description         string
		expectedDeviceNodes []*DeviceNode
	}{
		{
			description: "test DeviceNodes()",
			expectedDeviceNodes: []*DeviceNode{
				{
					ContainerPath: fmt.Sprintf(warboyMgmtFileExp, "/dev/npu0"),
					HostPath:      fmt.Sprintf(warboyMgmtFileExp, "/dev/npu0"),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: "/dev/npu0pe0",
					HostPath:      "/dev/npu0pe0",
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: "/dev/npu0pe1",
					HostPath:      "/dev/npu0pe1",
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: "/dev/npu0pe0-1",
					HostPath:      "/dev/npu0pe0-1",
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: fmt.Sprintf(warboyChannelExp, "npu0", 0),
					HostPath:      fmt.Sprintf(warboyChannelExp, "npu0", 0),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: fmt.Sprintf(warboyChannelExp, "npu0", 1),
					HostPath:      fmt.Sprintf(warboyChannelExp, "npu0", 1),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: fmt.Sprintf(warboyChannelExp, "npu0", 2),
					HostPath:      fmt.Sprintf(warboyChannelExp, "npu0", 2),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: fmt.Sprintf(warboyChannelExp, "npu0", 3),
					HostPath:      fmt.Sprintf(warboyChannelExp, "npu0", 3),
					Permissions:   readWriteOpt,
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			manifest, _ := NewWarboyManifest(newTestWarboyDevice())
			actualDeviceNodes := manifest.DeviceNodes()

			sort.Slice(tc.expectedDeviceNodes, func(i, j int) bool {
				return tc.expectedDeviceNodes[i].ContainerPath < tc.expectedDeviceNodes[j].ContainerPath
			})

			sort.Slice(actualDeviceNodes, func(i, j int) bool {
				return actualDeviceNodes[i].ContainerPath < actualDeviceNodes[j].ContainerPath
			})

			assert.Equal(t, tc.expectedDeviceNodes, actualDeviceNodes)
		})
	}
}

func TestWarboyMountPaths(t *testing.T) {
	tests := []struct {
		description        string
		expectedMountPaths []*Mount
	}{
		{
			description: "test MountPaths()",
			expectedMountPaths: []*Mount{
				{
					ContainerPath: warboySysClassRoot + fmt.Sprintf(warboyMgmtFileExp, "npu0"),
					HostPath:      warboySysClassRoot + fmt.Sprintf(warboyMgmtFileExp, "npu0"),
					Options:       []string{readOnlyOpt},
				},
				{
					ContainerPath: warboySysClassRoot + "npu0pe0",
					HostPath:      warboySysClassRoot + "npu0pe0",
					Options:       []string{readOnlyOpt},
				},
				{
					ContainerPath: warboySysClassRoot + "npu0pe1",
					HostPath:      warboySysClassRoot + "npu0pe1",
					Options:       []string{readOnlyOpt},
				},
				{
					ContainerPath: warboySysClassRoot + "npu0pe0-1",
					HostPath:      warboySysClassRoot + "npu0pe0-1",
					Options:       []string{readOnlyOpt},
				},
				{
					ContainerPath: warboySysDevicesRoot + fmt.Sprintf(warboyMgmtFileExp, "npu0"),
					HostPath:      warboySysDevicesRoot + fmt.Sprintf(warboyMgmtFileExp, "npu0"),
					Options:       []string{readOnlyOpt},
				},
				{
					ContainerPath: warboySysDevicesRoot + "npu0pe0",
					HostPath:      warboySysDevicesRoot + "npu0pe0",
					Options:       []string{readOnlyOpt},
				},
				{
					ContainerPath: warboySysDevicesRoot + "npu0pe1",
					HostPath:      warboySysDevicesRoot + "npu0pe1",
					Options:       []string{readOnlyOpt},
				},
				{
					ContainerPath: warboySysDevicesRoot + "npu0pe0-1",
					HostPath:      warboySysDevicesRoot + "npu0pe0-1",
					Options:       []string{readOnlyOpt},
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			manifest, _ := NewWarboyManifest(newTestWarboyDevice())
			actualMountPaths := manifest.MountPaths()

			sort.Slice(tc.expectedMountPaths, func(i, j int) bool {
				return tc.expectedMountPaths[i].ContainerPath < tc.expectedMountPaths[j].ContainerPath
			})

			sort.Slice(actualMountPaths, func(i, j int) bool {
				return actualMountPaths[i].ContainerPath < actualMountPaths[j].ContainerPath
			})

			assert.Equal(t, tc.expectedMountPaths, actualMountPaths)
		})
	}
}
