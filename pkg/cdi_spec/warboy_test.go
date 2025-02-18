package cdi_spec

import (
	"fmt"
	"sort"
	"tags.cncf.io/container-device-interface/specs-go"
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
		expectedDeviceNodes []*specs.DeviceNode
	}{
		{
			description: "test deviceNodes()",
			expectedDeviceNodes: []*specs.DeviceNode{
				{
					Path:        fmt.Sprintf(warboyMgmtFileExp, "/dev/npu0"),
					HostPath:    fmt.Sprintf(warboyMgmtFileExp, "/dev/npu0"),
					Permissions: readWriteOpt,
				}, {
					Path:        "/dev/npu0pe0",
					HostPath:    "/dev/npu0pe0",
					Permissions: readWriteOpt,
				}, {
					Path:        "/dev/npu0pe1",
					HostPath:    "/dev/npu0pe1",
					Permissions: readWriteOpt,
				}, {
					Path:        "/dev/npu0pe0-1",
					HostPath:    "/dev/npu0pe0-1",
					Permissions: readWriteOpt,
				}, {
					Path:        fmt.Sprintf(warboyChannelExp, "npu0", 0),
					HostPath:    fmt.Sprintf(warboyChannelExp, "npu0", 0),
					Permissions: readWriteOpt,
				}, {
					Path:        fmt.Sprintf(warboyChannelExp, "npu0", 1),
					HostPath:    fmt.Sprintf(warboyChannelExp, "npu0", 1),
					Permissions: readWriteOpt,
				}, {
					Path:        fmt.Sprintf(warboyChannelExp, "npu0", 2),
					HostPath:    fmt.Sprintf(warboyChannelExp, "npu0", 2),
					Permissions: readWriteOpt,
				}, {
					Path:        fmt.Sprintf(warboyChannelExp, "npu0", 3),
					HostPath:    fmt.Sprintf(warboyChannelExp, "npu0", 3),
					Permissions: readWriteOpt,
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			spec, _ := newWarboyDeviceSpec(newTestWarboyDevice())
			actualDeviceNodes := spec.deviceNodes()

			sort.Slice(tc.expectedDeviceNodes, func(i, j int) bool {
				return tc.expectedDeviceNodes[i].Path < tc.expectedDeviceNodes[j].Path
			})

			sort.Slice(actualDeviceNodes, func(i, j int) bool {
				return actualDeviceNodes[i].Path < actualDeviceNodes[j].Path
			})

			assert.Equal(t, tc.expectedDeviceNodes, actualDeviceNodes)
		})
	}
}
