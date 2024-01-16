package manifest

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/device"
)

func newTestWarboyDevice() device.Device {
	newDevice, err := device.NewDevice(0,
		[]string{
			device.Abs("../../testdata/device/testdata/test-0/dev/npu0"),
			device.Abs("../../testdata/device/testdata/test-0/dev/npu0pe0"),
			device.Abs("../../testdata/device/testdata/test-0/dev/npu0pe1"),
			device.Abs("../../testdata/device/testdata/test-0/dev/npu0pe0-1"),
		},
		"../../testdata/device/testdata/test-0/dev",
		"../../testdata/device/testdata/test-0/sys")
	if err != nil {
		return nil
	}

	return newDevice
}

func TestDeviceNodes(t *testing.T) {
	tests := []struct {
		description         string
		expectedDeviceNodes []DeviceNode
	}{
		{
			description: "test DeviceNodes()",
			expectedDeviceNodes: []DeviceNode{
				{
					ContainerPath: devRoot + fmt.Sprintf(mgmtFileExp, "npu0"),
					HostPath:      devRoot + fmt.Sprintf(mgmtFileExp, "npu0"),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: device.Abs("../../testdata/device/testdata/test-0/dev/npu0"),
					HostPath:      device.Abs("../../testdata/device/testdata/test-0/dev/npu0"),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: device.Abs("../../testdata/device/testdata/test-0/dev/npu0pe0"),
					HostPath:      device.Abs("../../testdata/device/testdata/test-0/dev/npu0pe0"),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: device.Abs("../../testdata/device/testdata/test-0/dev/npu0pe1"),
					HostPath:      device.Abs("../../testdata/device/testdata/test-0/dev/npu0pe1"),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: device.Abs("../../testdata/device/testdata/test-0/dev/npu0pe0-1"),
					HostPath:      device.Abs("../../testdata/device/testdata/test-0/dev/npu0pe0-1"),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: fmt.Sprintf(channelExp, "npu0", 0),
					HostPath:      fmt.Sprintf(channelExp, "npu0", 0),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: fmt.Sprintf(channelExp, "npu0", 1),
					HostPath:      fmt.Sprintf(channelExp, "npu0", 1),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: fmt.Sprintf(channelExp, "npu0", 2),
					HostPath:      fmt.Sprintf(channelExp, "npu0", 2),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: fmt.Sprintf(channelExp, "npu0", 3),
					HostPath:      fmt.Sprintf(channelExp, "npu0", 3),
					Permissions:   readWriteOpt,
				},
			},
		},
	}
	for _, tc := range tests {
		manifest := NewWarboyManifest(newTestWarboyDevice())
		actualDeviceNodes := manifest.DeviceNodes()
		if !reflect.DeepEqual(actualDeviceNodes, tc.expectedDeviceNodes) {
			t.Errorf("expected %v but got %v", tc.expectedDeviceNodes, actualDeviceNodes)
		}

	}

}

func TestMountPaths(t *testing.T) {
	tests := []struct {
		description        string
		expectedMountPaths []Mount
	}{
		{
			description: "test MountPaths()",
			expectedMountPaths: []Mount{
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
		manifest := NewWarboyManifest(newTestWarboyDevice())
		actualMountPaths := manifest.MountPaths()
		if !reflect.DeepEqual(actualMountPaths, tc.expectedMountPaths) {
			t.Errorf("expected %v but got %v", tc.expectedMountPaths, actualMountPaths)
		}

	}

}
