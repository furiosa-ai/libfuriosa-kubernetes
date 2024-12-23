package cdi_spec

import (
	"testing"

	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"github.com/stretchr/testify/assert"

	"tags.cncf.io/container-device-interface/specs-go"
)

func TestFilterPartitionedDeviceNodes(t *testing.T) {
	rngd := smi.GetStaticMockDevices(smi.ArchRngd)[0]
	rngdSpec, _ := newRngdDeviceSpec(rngd)

	readWriteOpt := "rw"

	tests := []struct {
		description  string
		spec         DeviceSpec
		startCore    int
		endCore      int
		mustContains []*specs.DeviceNode
	}{
		{
			description: "RNGD, single-core strategy",
			spec:        rngdSpec,
			startCore:   1,
			endCore:     1,
			mustContains: []*specs.DeviceNode{
				{
					Path:        "/dev/rngd/npu0pe1",
					HostPath:    "/dev/rngd/npu0pe1",
					Permissions: readWriteOpt,
				},
			},
		},
		{
			description: "RNGD, dual-core strategy",
			spec:        rngdSpec,
			startCore:   2,
			endCore:     3,
			mustContains: []*specs.DeviceNode{
				{
					Path:        "/dev/rngd/npu0pe2-3",
					HostPath:    "/dev/rngd/npu0pe2-3",
					Permissions: readWriteOpt,
				},
				{
					Path:        "/dev/rngd/npu0pe2",
					HostPath:    "/dev/rngd/npu0pe2",
					Permissions: readWriteOpt,
				},
				{
					Path:        "/dev/rngd/npu0pe3",
					HostPath:    "/dev/rngd/npu0pe3",
					Permissions: readWriteOpt,
				},
			},
		},
		{
			description: "RNGD, quad-core strategy",
			spec:        rngdSpec,
			startCore:   4,
			endCore:     7,
			mustContains: []*specs.DeviceNode{
				{
					Path:        "/dev/rngd/npu0pe4-7",
					HostPath:    "/dev/rngd/npu0pe4-7",
					Permissions: readWriteOpt,
				},
				{
					Path:        "/dev/rngd/npu0pe4-5",
					HostPath:    "/dev/rngd/npu0pe4-5",
					Permissions: readWriteOpt,
				},
				{
					Path:        "/dev/rngd/npu0pe6-7",
					HostPath:    "/dev/rngd/npu0pe6-7",
					Permissions: readWriteOpt,
				},
				{
					Path:        "/dev/rngd/npu0pe4",
					HostPath:    "/dev/rngd/npu0pe4",
					Permissions: readWriteOpt,
				},
				{
					Path:        "/dev/rngd/npu0pe5",
					HostPath:    "/dev/rngd/npu0pe5",
					Permissions: readWriteOpt,
				},
				{
					Path:        "/dev/rngd/npu0pe6",
					HostPath:    "/dev/rngd/npu0pe6",
					Permissions: readWriteOpt,
				},
				{
					Path:        "/dev/rngd/npu0pe7",
					HostPath:    "/dev/rngd/npu0pe7",
					Permissions: readWriteOpt,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			actual := filterPartitionedDeviceNodes(tc.spec, tc.startCore, tc.endCore)
			assert.Subset(t, actual, tc.mustContains)
		})
	}
}
