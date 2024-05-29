package manifest

import (
	"fmt"
	"reflect"
	"testing"

	furiosaSmi "github.com/furiosa-ai/libfuriosa-kubernetes/pkg/furiosa_smi_go"
)

func newTestWarboyDevice() furiosaSmi.Device {
	return furiosaSmi.GetStaticMockWarboyDevice(0)
}

func TestDeviceNodes(t *testing.T) {
	tests := []struct {
		description         string
		expectedDeviceNodes []*DeviceNode
	}{
		{
			description: "test DeviceNodes()",
			expectedDeviceNodes: []*DeviceNode{
				{
					ContainerPath: fmt.Sprintf(mgmtFileExp, "/dev/npu0"),
					HostPath:      fmt.Sprintf(mgmtFileExp, "/dev/npu0"),
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
					ContainerPath: "/dev/" + fmt.Sprintf(channelExp, "npu0", 0),
					HostPath:      "/dev/" + fmt.Sprintf(channelExp, "npu0", 0),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: "/dev/" + fmt.Sprintf(channelExp, "npu0", 1),
					HostPath:      "/dev/" + fmt.Sprintf(channelExp, "npu0", 1),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: "/dev/" + fmt.Sprintf(channelExp, "npu0", 2),
					HostPath:      "/dev/" + fmt.Sprintf(channelExp, "npu0", 2),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: "/dev/" + fmt.Sprintf(channelExp, "npu0", 3),
					HostPath:      "/dev/" + fmt.Sprintf(channelExp, "npu0", 3),
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

func TestMountPaths(t *testing.T) {
	tests := []struct {
		description        string
		expectedMountPaths []*Mount
	}{
		{
			description: "test MountPaths()",
			expectedMountPaths: []*Mount{
				{
					ContainerPath: sysClassRoot + fmt.Sprintf(mgmtFileExp, "npu0"),
					HostPath:      sysClassRoot + fmt.Sprintf(mgmtFileExp, "npu0"),
					Options:       []string{readOnlyOpt},
				},
				{
					ContainerPath: sysClassRoot + "npu0",
					HostPath:      sysClassRoot + "npu0",
					Options:       []string{readOnlyOpt},
				},
				{
					ContainerPath: sysClassRoot + "npu0pe0",
					HostPath:      sysClassRoot + "npu0pe0",
					Options:       []string{readOnlyOpt},
				},
				{
					ContainerPath: sysClassRoot + "npu0pe1",
					HostPath:      sysClassRoot + "npu0pe1",
					Options:       []string{readOnlyOpt},
				},
				{
					ContainerPath: sysClassRoot + "npu0pe0-1",
					HostPath:      sysClassRoot + "npu0pe0-1",
					Options:       []string{readOnlyOpt},
				},
				{
					ContainerPath: sysDevicesRoot + fmt.Sprintf(mgmtFileExp, "npu0"),
					HostPath:      sysDevicesRoot + fmt.Sprintf(mgmtFileExp, "npu0"),
					Options:       []string{readOnlyOpt},
				},
				{
					ContainerPath: sysDevicesRoot + "npu0",
					HostPath:      sysDevicesRoot + "npu0",
					Options:       []string{readOnlyOpt},
				},
				{
					ContainerPath: sysDevicesRoot + "npu0pe0",
					HostPath:      sysDevicesRoot + "npu0pe0",
					Options:       []string{readOnlyOpt},
				},
				{
					ContainerPath: sysDevicesRoot + "npu0pe1",
					HostPath:      sysDevicesRoot + "npu0pe1",
					Options:       []string{readOnlyOpt},
				},
				{
					ContainerPath: sysDevicesRoot + "npu0pe0-1",
					HostPath:      sysDevicesRoot + "npu0pe0-1",
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
