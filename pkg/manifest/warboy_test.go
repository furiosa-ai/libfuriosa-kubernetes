package manifest

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/smi"
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
					ContainerPath: "/dev/npu0",
					HostPath:      "/dev/npu0",
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
					ContainerPath: "/dev/" + fmt.Sprintf(warboyChannelExp, "npu0", 0),
					HostPath:      "/dev/" + fmt.Sprintf(warboyChannelExp, "npu0", 0),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: "/dev/" + fmt.Sprintf(warboyChannelExp, "npu0", 1),
					HostPath:      "/dev/" + fmt.Sprintf(warboyChannelExp, "npu0", 1),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: "/dev/" + fmt.Sprintf(warboyChannelExp, "npu0", 2),
					HostPath:      "/dev/" + fmt.Sprintf(warboyChannelExp, "npu0", 2),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: "/dev/" + fmt.Sprintf(warboyChannelExp, "npu0", 3),
					HostPath:      "/dev/" + fmt.Sprintf(warboyChannelExp, "npu0", 3),
					Permissions:   readWriteOpt,
				},
			},
		},
	}
	for _, tc := range tests {
		manifest, _ := NewWarboyManifest(newTestWarboyDevice())
		actualDeviceNodes := manifest.DeviceNodes()
		if !reflect.DeepEqual(actualDeviceNodes, tc.expectedDeviceNodes) {
			t.Errorf("expected %v but got %v", tc.expectedDeviceNodes, actualDeviceNodes)
		}
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
					ContainerPath: warboySysClassRoot + "npu0",
					HostPath:      warboySysClassRoot + "npu0",
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
					ContainerPath: warboySysDevicesRoot + "npu0",
					HostPath:      warboySysDevicesRoot + "npu0",
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
		manifest, _ := NewWarboyManifest(newTestWarboyDevice())
		actualMountPaths := manifest.MountPaths()
		if !reflect.DeepEqual(actualMountPaths, tc.expectedMountPaths) {
			t.Errorf("expected %v but got %v", tc.expectedMountPaths, actualMountPaths)
		}
	}
}
