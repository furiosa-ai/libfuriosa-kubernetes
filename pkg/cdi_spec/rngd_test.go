package cdi_spec

import (
	"fmt"
	"tags.cncf.io/container-device-interface/specs-go"
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
		expectedDeviceNodes []*specs.DeviceNode
	}{
		{
			description: "test deviceNodes()",
			expectedDeviceNodes: []*specs.DeviceNode{
				{
					Path:        rngdDevFsRoot + fmt.Sprintf(rngdMgmtFileExp, "npu0"),
					HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdMgmtFileExp, "npu0"),
					Permissions: readWriteOpt,
				}, {
					Path:        "/dev/rngd/npu0pe0",
					HostPath:    "/dev/rngd/npu0pe0",
					Permissions: readWriteOpt,
				}, {
					Path:        "/dev/rngd/npu0pe1",
					HostPath:    "/dev/rngd/npu0pe1",
					Permissions: readWriteOpt,
				}, {
					Path:        "/dev/rngd/npu0pe0-1",
					HostPath:    "/dev/rngd/npu0pe0-1",
					Permissions: readWriteOpt,
				}, {
					Path:        "/dev/rngd/npu0pe2",
					HostPath:    "/dev/rngd/npu0pe2",
					Permissions: readWriteOpt,
				}, {
					Path:        "/dev/rngd/npu0pe3",
					HostPath:    "/dev/rngd/npu0pe3",
					Permissions: readWriteOpt,
				}, {
					Path:        "/dev/rngd/npu0pe2-3",
					HostPath:    "/dev/rngd/npu0pe2-3",
					Permissions: readWriteOpt,
				}, {
					Path:        "/dev/rngd/npu0pe0-3",
					HostPath:    "/dev/rngd/npu0pe0-3",
					Permissions: readWriteOpt,
				}, {
					Path:        "/dev/rngd/npu0pe4",
					HostPath:    "/dev/rngd/npu0pe4",
					Permissions: readWriteOpt,
				}, {
					Path:        "/dev/rngd/npu0pe5",
					HostPath:    "/dev/rngd/npu0pe5",
					Permissions: readWriteOpt,
				}, {
					Path:        "/dev/rngd/npu0pe4-5",
					HostPath:    "/dev/rngd/npu0pe4-5",
					Permissions: readWriteOpt,
				}, {
					Path:        "/dev/rngd/npu0pe6",
					HostPath:    "/dev/rngd/npu0pe6",
					Permissions: readWriteOpt,
				}, {
					Path:        "/dev/rngd/npu0pe7",
					HostPath:    "/dev/rngd/npu0pe7",
					Permissions: readWriteOpt,
				}, {
					Path:        "/dev/rngd/npu0pe6-7",
					HostPath:    "/dev/rngd/npu0pe6-7",
					Permissions: readWriteOpt,
				}, {
					Path:        "/dev/rngd/npu0pe4-7",
					HostPath:    "/dev/rngd/npu0pe4-7",
					Permissions: readWriteOpt,
				}, {
					Path:        rngdDevFsRoot + fmt.Sprintf(rngdChannelExp, "npu0", 0),
					HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdChannelExp, "npu0", 0),
					Permissions: readWriteOpt,
				}, {
					Path:        rngdDevFsRoot + fmt.Sprintf(rngdChannelExp, "npu0", 1),
					HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdChannelExp, "npu0", 1),
					Permissions: readWriteOpt,
				}, {
					Path:        rngdDevFsRoot + fmt.Sprintf(rngdChannelExp, "npu0", 2),
					HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdChannelExp, "npu0", 2),
					Permissions: readWriteOpt,
				}, {
					Path:        rngdDevFsRoot + fmt.Sprintf(rngdChannelExp, "npu0", 3),
					HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdChannelExp, "npu0", 3),
					Permissions: readWriteOpt,
				}, {
					Path:        rngdDevFsRoot + fmt.Sprintf(rngdChannelExp, "npu0", 4),
					HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdChannelExp, "npu0", 4),
					Permissions: readWriteOpt,
				}, {
					Path:        rngdDevFsRoot + fmt.Sprintf(rngdChannelExp, "npu0", 5),
					HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdChannelExp, "npu0", 5),
					Permissions: readWriteOpt,
				}, {
					Path:        rngdDevFsRoot + fmt.Sprintf(rngdChannelExp, "npu0", 6),
					HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdChannelExp, "npu0", 6),
					Permissions: readWriteOpt,
				}, {
					Path:        rngdDevFsRoot + fmt.Sprintf(rngdChannelExp, "npu0", 7),
					HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdChannelExp, "npu0", 7),
					Permissions: readWriteOpt,
				}, {
					Path:        rngdDevFsRoot + fmt.Sprintf(rngdRemoteChannelExp, "npu0", 0),
					HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdRemoteChannelExp, "npu0", 0),
					Permissions: readWriteOpt,
				}, {
					Path:        rngdDevFsRoot + fmt.Sprintf(rngdRemoteChannelExp, "npu0", 1),
					HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdRemoteChannelExp, "npu0", 1),
					Permissions: readWriteOpt,
				}, {
					Path:        rngdDevFsRoot + fmt.Sprintf(rngdRemoteChannelExp, "npu0", 2),
					HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdRemoteChannelExp, "npu0", 2),
					Permissions: readWriteOpt,
				}, {
					Path:        rngdDevFsRoot + fmt.Sprintf(rngdRemoteChannelExp, "npu0", 3),
					HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdRemoteChannelExp, "npu0", 3),
					Permissions: readWriteOpt,
				}, {
					Path:        rngdDevFsRoot + fmt.Sprintf(rngdRemoteChannelExp, "npu0", 4),
					HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdRemoteChannelExp, "npu0", 4),
					Permissions: readWriteOpt,
				}, {
					Path:        rngdDevFsRoot + fmt.Sprintf(rngdRemoteChannelExp, "npu0", 5),
					HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdRemoteChannelExp, "npu0", 5),
					Permissions: readWriteOpt,
				}, {
					Path:        rngdDevFsRoot + fmt.Sprintf(rngdRemoteChannelExp, "npu0", 6),
					HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdRemoteChannelExp, "npu0", 6),
					Permissions: readWriteOpt,
				}, {
					Path:        rngdDevFsRoot + fmt.Sprintf(rngdRemoteChannelExp, "npu0", 7),
					HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdRemoteChannelExp, "npu0", 7),
					Permissions: readWriteOpt,
				}, {
					Path:        rngdDevFsRoot + fmt.Sprintf(rngdDmaRemappingExp, "npu0"),
					HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdDmaRemappingExp, "npu0"),
					Permissions: readWriteOpt,
				}, {
					Path:        rngdDevFsRoot + fmt.Sprintf(rngdBar0Exp, "npu0"),
					HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdBar0Exp, "npu0"),
					Permissions: readWriteOpt,
				}, {
					Path:        rngdDevFsRoot + fmt.Sprintf(rngdBar2Exp, "npu0"),
					HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdBar2Exp, "npu0"),
					Permissions: readWriteOpt,
				}, {
					Path:        rngdDevFsRoot + fmt.Sprintf(rngdBar4Exp, "npu0"),
					HostPath:    rngdDevFsRoot + fmt.Sprintf(rngdBar4Exp, "npu0"),
					Permissions: readWriteOpt,
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			spec, _ := newRngdDeviceSpec(newTestRngdDevice())
			actualDeviceNodes := spec.deviceNodes()

			assert.Equal(t, tc.expectedDeviceNodes, actualDeviceNodes)
		})
	}

}

func TestRngdDeviceMounts(t *testing.T) {
	tests := []struct {
		description    string
		expectedMounts []*specs.Mount
	}{
		{
			description: "test rngd mounts()",
			expectedMounts: []*specs.Mount{
				{
					HostPath:      rngdSysDevicesRoot + rngdSysPrefix + fmt.Sprintf(rngdMgmtFileExp, "npu0"),
					ContainerPath: rngdSysDevicesRoot + rngdSysPrefix + fmt.Sprintf(rngdMgmtFileExp, "npu0"),
					Options:       []string{readOnlyOpt, bindOpt},
				},
				{
					HostPath:      rngdSysDevicesRoot + rngdSysPrefix + fmt.Sprintf(rngdBar0Exp, "npu0"),
					ContainerPath: rngdSysDevicesRoot + rngdSysPrefix + fmt.Sprintf(rngdBar0Exp, "npu0"),
					Options:       []string{readOnlyOpt, bindOpt},
				},
				{
					HostPath:      rngdSysDevicesRoot + rngdSysPrefix + fmt.Sprintf(rngdBar2Exp, "npu0"),
					ContainerPath: rngdSysDevicesRoot + rngdSysPrefix + fmt.Sprintf(rngdBar2Exp, "npu0"),
					Options:       []string{readOnlyOpt, bindOpt},
				},
				{
					HostPath:      rngdSysDevicesRoot + rngdSysPrefix + fmt.Sprintf(rngdBar4Exp, "npu0"),
					ContainerPath: rngdSysDevicesRoot + rngdSysPrefix + fmt.Sprintf(rngdBar4Exp, "npu0"),
					Options:       []string{readOnlyOpt, bindOpt},
				},
				{
					HostPath:      rngdSysDevicesRoot + rngdSysPrefix + fmt.Sprintf(rngdDmaRemappingExp, "npu0"),
					ContainerPath: rngdSysDevicesRoot + rngdSysPrefix + fmt.Sprintf(rngdDmaRemappingExp, "npu0"),
					Options:       []string{readOnlyOpt, bindOpt},
				},
				{
					HostPath:      rngdSysDevicesRoot + rngdSysPrefix + "npu0pe0",
					ContainerPath: rngdSysDevicesRoot + rngdSysPrefix + "npu0pe0",
					Options:       []string{readOnlyOpt, bindOpt},
				},
				{
					HostPath:      rngdSysDevicesRoot + rngdSysPrefix + "npu0pe1",
					ContainerPath: rngdSysDevicesRoot + rngdSysPrefix + "npu0pe1",
					Options:       []string{readOnlyOpt, bindOpt},
				},
				{
					HostPath:      rngdSysDevicesRoot + rngdSysPrefix + "npu0pe0-1",
					ContainerPath: rngdSysDevicesRoot + rngdSysPrefix + "npu0pe0-1",
					Options:       []string{readOnlyOpt, bindOpt},
				},
				{
					HostPath:      rngdSysDevicesRoot + rngdSysPrefix + "npu0pe2",
					ContainerPath: rngdSysDevicesRoot + rngdSysPrefix + "npu0pe2",
					Options:       []string{readOnlyOpt, bindOpt},
				},
				{
					HostPath:      rngdSysDevicesRoot + rngdSysPrefix + "npu0pe3",
					ContainerPath: rngdSysDevicesRoot + rngdSysPrefix + "npu0pe3",
					Options:       []string{readOnlyOpt, bindOpt},
				},
				{
					HostPath:      rngdSysDevicesRoot + rngdSysPrefix + "npu0pe2-3",
					ContainerPath: rngdSysDevicesRoot + rngdSysPrefix + "npu0pe2-3",
					Options:       []string{readOnlyOpt, bindOpt},
				},
				{
					HostPath:      rngdSysDevicesRoot + rngdSysPrefix + "npu0pe0-3",
					ContainerPath: rngdSysDevicesRoot + rngdSysPrefix + "npu0pe0-3",
					Options:       []string{readOnlyOpt, bindOpt},
				},
				{
					HostPath:      rngdSysDevicesRoot + rngdSysPrefix + "npu0pe4",
					ContainerPath: rngdSysDevicesRoot + rngdSysPrefix + "npu0pe4",
					Options:       []string{readOnlyOpt, bindOpt},
				},
				{
					HostPath:      rngdSysDevicesRoot + rngdSysPrefix + "npu0pe5",
					ContainerPath: rngdSysDevicesRoot + rngdSysPrefix + "npu0pe5",
					Options:       []string{readOnlyOpt, bindOpt},
				},
				{
					HostPath:      rngdSysDevicesRoot + rngdSysPrefix + "npu0pe4-5",
					ContainerPath: rngdSysDevicesRoot + rngdSysPrefix + "npu0pe4-5",
					Options:       []string{readOnlyOpt, bindOpt},
				},
				{
					HostPath:      rngdSysDevicesRoot + rngdSysPrefix + "npu0pe6",
					ContainerPath: rngdSysDevicesRoot + rngdSysPrefix + "npu0pe6",
					Options:       []string{readOnlyOpt, bindOpt},
				},
				{
					HostPath:      rngdSysDevicesRoot + rngdSysPrefix + "npu0pe7",
					ContainerPath: rngdSysDevicesRoot + rngdSysPrefix + "npu0pe7",
					Options:       []string{readOnlyOpt, bindOpt},
				},
				{
					HostPath:      rngdSysDevicesRoot + rngdSysPrefix + "npu0pe6-7",
					ContainerPath: rngdSysDevicesRoot + rngdSysPrefix + "npu0pe6-7",
					Options:       []string{readOnlyOpt, bindOpt},
				},
				{
					HostPath:      rngdSysDevicesRoot + rngdSysPrefix + "npu0pe4-7",
					ContainerPath: rngdSysDevicesRoot + rngdSysPrefix + "npu0pe4-7",
					Options:       []string{readOnlyOpt, bindOpt},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			spec, _ := newRngdDeviceSpec(newTestRngdDevice())
			actualMounts := spec.mounts()

			assert.Equal(t, tc.expectedMounts, actualMounts)
		})
	}
}
