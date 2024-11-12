package manifest

import (
	"fmt"
	"testing"

	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"github.com/stretchr/testify/assert"
)

func newTestRngdDevice() smi.Device {
	return smi.GetStaticMockDevice(smi.ArchRngd, 0)
}

func TestRngdDeviceNodes(t *testing.T) {
	tests := []struct {
		description         string
		expectedDeviceNodes []*DeviceNode
	}{
		{
			description: "test DeviceNodes()",
			expectedDeviceNodes: []*DeviceNode{
				{
					ContainerPath: fmt.Sprintf(devFsRngdMgmtFileExp, "npu0"),
					HostPath:      fmt.Sprintf(devFsRngdMgmtFileExp, "npu0"),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: "/dev/rngd/npu0pe0",
					HostPath:      "/dev/rngd/npu0pe0",
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: "/dev/rngd/npu0pe1",
					HostPath:      "/dev/rngd/npu0pe1",
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: "/dev/rngd/npu0pe0-1",
					HostPath:      "/dev/rngd/npu0pe0-1",
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: "/dev/rngd/npu0pe2",
					HostPath:      "/dev/rngd/npu0pe2",
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: "/dev/rngd/npu0pe3",
					HostPath:      "/dev/rngd/npu0pe3",
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: "/dev/rngd/npu0pe2-3",
					HostPath:      "/dev/rngd/npu0pe2-3",
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: "/dev/rngd/npu0pe0-3",
					HostPath:      "/dev/rngd/npu0pe0-3",
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: "/dev/rngd/npu0pe4",
					HostPath:      "/dev/rngd/npu0pe4",
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: "/dev/rngd/npu0pe5",
					HostPath:      "/dev/rngd/npu0pe5",
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: "/dev/rngd/npu0pe4-5",
					HostPath:      "/dev/rngd/npu0pe4-5",
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: "/dev/rngd/npu0pe6",
					HostPath:      "/dev/rngd/npu0pe6",
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: "/dev/rngd/npu0pe7",
					HostPath:      "/dev/rngd/npu0pe7",
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: "/dev/rngd/npu0pe6-7",
					HostPath:      "/dev/rngd/npu0pe6-7",
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: "/dev/rngd/npu0pe4-7",
					HostPath:      "/dev/rngd/npu0pe4-7",
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: fmt.Sprintf(rngdChannelExp, "npu0", 0),
					HostPath:      fmt.Sprintf(rngdChannelExp, "npu0", 0),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: fmt.Sprintf(rngdChannelExp, "npu0", 1),
					HostPath:      fmt.Sprintf(rngdChannelExp, "npu0", 1),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: fmt.Sprintf(rngdChannelExp, "npu0", 2),
					HostPath:      fmt.Sprintf(rngdChannelExp, "npu0", 2),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: fmt.Sprintf(rngdChannelExp, "npu0", 3),
					HostPath:      fmt.Sprintf(rngdChannelExp, "npu0", 3),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: fmt.Sprintf(rngdChannelExp, "npu0", 4),
					HostPath:      fmt.Sprintf(rngdChannelExp, "npu0", 4),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: fmt.Sprintf(rngdChannelExp, "npu0", 5),
					HostPath:      fmt.Sprintf(rngdChannelExp, "npu0", 5),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: fmt.Sprintf(rngdChannelExp, "npu0", 6),
					HostPath:      fmt.Sprintf(rngdChannelExp, "npu0", 6),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: fmt.Sprintf(rngdChannelExp, "npu0", 7),
					HostPath:      fmt.Sprintf(rngdChannelExp, "npu0", 7),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: fmt.Sprintf(rngdRemoteChannelExp, "npu0", 0),
					HostPath:      fmt.Sprintf(rngdRemoteChannelExp, "npu0", 0),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: fmt.Sprintf(rngdRemoteChannelExp, "npu0", 1),
					HostPath:      fmt.Sprintf(rngdRemoteChannelExp, "npu0", 1),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: fmt.Sprintf(rngdRemoteChannelExp, "npu0", 2),
					HostPath:      fmt.Sprintf(rngdRemoteChannelExp, "npu0", 2),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: fmt.Sprintf(rngdRemoteChannelExp, "npu0", 3),
					HostPath:      fmt.Sprintf(rngdRemoteChannelExp, "npu0", 3),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: fmt.Sprintf(rngdRemoteChannelExp, "npu0", 4),
					HostPath:      fmt.Sprintf(rngdRemoteChannelExp, "npu0", 4),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: fmt.Sprintf(rngdRemoteChannelExp, "npu0", 5),
					HostPath:      fmt.Sprintf(rngdRemoteChannelExp, "npu0", 5),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: fmt.Sprintf(rngdRemoteChannelExp, "npu0", 6),
					HostPath:      fmt.Sprintf(rngdRemoteChannelExp, "npu0", 6),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: fmt.Sprintf(rngdRemoteChannelExp, "npu0", 7),
					HostPath:      fmt.Sprintf(rngdRemoteChannelExp, "npu0", 7),
					Permissions:   readWriteOpt,
				}, {
					ContainerPath: fmt.Sprintf(rngdDmaRemappingExp, "npu0"),
					HostPath:      fmt.Sprintf(rngdDmaRemappingExp, "npu0"),
					Permissions:   readWriteOpt,
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			manifest, _ := NewRngdManifest(newTestRngdDevice())
			actualDeviceNodes := manifest.DeviceNodes()

			assert.Equal(t, tc.expectedDeviceNodes, actualDeviceNodes)
		})
	}

}

func TestRngdMountPaths(t *testing.T) {
	tests := []struct {
		description        string
		expectedMountPaths []*Mount
	}{
		{
			description: "test MountPaths()",
			expectedMountPaths: []*Mount{
				{
					ContainerPath: rngdSysClassRoot + fmt.Sprintf(sysFsRngdMgmtFileExp, "npu0"),
					HostPath:      rngdSysClassRoot + fmt.Sprintf(sysFsRngdMgmtFileExp, "npu0"),
					Options:       []string{readOnlyOpt},
				},
				{
					ContainerPath: rngdSysDevicesRoot + fmt.Sprintf(sysFsRngdMgmtFileExp, "npu0"),
					HostPath:      rngdSysDevicesRoot + fmt.Sprintf(sysFsRngdMgmtFileExp, "npu0"),
					Options:       []string{readOnlyOpt},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			manifest, _ := NewRngdManifest(newTestRngdDevice())
			actualMountPaths := manifest.MountPaths()

			assert.Equal(t, tc.expectedMountPaths, actualMountPaths)
		})
	}
}
