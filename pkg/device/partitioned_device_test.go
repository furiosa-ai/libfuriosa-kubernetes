package device

import (
	"fmt"
	"testing"

	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/manifest"
	"github.com/stretchr/testify/assert"
)

func TestGenerateIndexForPartitionedDevice(t *testing.T) {
	var tests []struct {
		description      string
		originalIndex    int
		partitionIndex   int
		partitionsLength int
		expected         int
	}

	// first element: partitionIndex, second element: finalIndex
	boardMatrix := [][]int{
		{0, 1, 2, 3, 4, 5, 6, 7},
		{8, 9, 10, 11, 12, 13, 14, 15},
		{16, 17, 18, 19, 20, 21, 22, 23},
		{24, 25, 26, 27, 28, 29, 30, 31},
	}

	for i := 0; i < 4; i++ {
		for j := 0; j < 8; j++ {
			tests = append(tests, struct {
				description      string
				originalIndex    int
				partitionIndex   int
				partitionsLength int
				expected         int
			}{
				description:      fmt.Sprintf("Original Board %d, Partition %d", i, j),
				originalIndex:    i,
				partitionIndex:   j,
				partitionsLength: 8,
				expected:         boardMatrix[i][j],
			})
		}
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			generatedFinalIndex := generateIndexForPartitionedDevice(tc.originalIndex, tc.partitionIndex, tc.partitionsLength)
			assert.Equal(t, tc.expected, generatedFinalIndex)
		})
	}
}

func TestFilterPartitionedDeviceNodes(t *testing.T) {
	warboy := smi.GetStaticMockDevices(smi.ArchWarboy)[0]
	rngd := smi.GetStaticMockDevices(smi.ArchRngd)[0]

	warboyManifest, _ := manifest.NewWarboyManifest(warboy)
	rngdManifest, _ := manifest.NewRngdManifest(rngd)

	readWriteOpt := "rw"

	tests := []struct {
		description  string
		manifest     manifest.Manifest
		partition    Partition
		mustContains []*manifest.DeviceNode
	}{
		{
			description: "Warboy, single-core strategy",
			manifest:    warboyManifest,
			partition:   Partition{start: 0, end: 0},
			mustContains: []*manifest.DeviceNode{
				{
					ContainerPath: "/dev/npu0pe0",
					HostPath:      "/dev/npu0pe0",
					Permissions:   readWriteOpt,
				},
			},
		},
		{
			description: "Warboy, dual-core strategy",
			manifest:    warboyManifest,
			partition:   Partition{start: 0, end: 1},
			mustContains: []*manifest.DeviceNode{
				{
					ContainerPath: "/dev/npu0pe0-1",
					HostPath:      "/dev/npu0pe0-1",
					Permissions:   readWriteOpt,
				},
				{
					ContainerPath: "/dev/npu0pe0",
					HostPath:      "/dev/npu0pe0",
					Permissions:   readWriteOpt,
				},
				{
					ContainerPath: "/dev/npu0pe1",
					HostPath:      "/dev/npu0pe1",
					Permissions:   readWriteOpt,
				},
			},
		},
		{
			description: "RNGD, single-core strategy",
			manifest:    rngdManifest,
			partition:   Partition{start: 1, end: 1},
			mustContains: []*manifest.DeviceNode{
				{
					ContainerPath: "/dev/rngd/npu0pe1",
					HostPath:      "/dev/rngd/npu0pe1",
					Permissions:   readWriteOpt,
				},
			},
		},
		{
			description: "RNGD, dual-core strategy",
			manifest:    rngdManifest,
			partition:   Partition{start: 2, end: 3},
			mustContains: []*manifest.DeviceNode{
				{
					ContainerPath: "/dev/rngd/npu0pe2-3",
					HostPath:      "/dev/rngd/npu0pe2-3",
					Permissions:   readWriteOpt,
				},
				{
					ContainerPath: "/dev/rngd/npu0pe2",
					HostPath:      "/dev/rngd/npu0pe2",
					Permissions:   readWriteOpt,
				},
				{
					ContainerPath: "/dev/rngd/npu0pe3",
					HostPath:      "/dev/rngd/npu0pe3",
					Permissions:   readWriteOpt,
				},
			},
		},
		{
			description: "RNGD, quad-core strategy",
			manifest:    rngdManifest,
			partition:   Partition{start: 4, end: 7},
			mustContains: []*manifest.DeviceNode{
				{
					ContainerPath: "/dev/rngd/npu0pe4-7",
					HostPath:      "/dev/rngd/npu0pe4-7",
					Permissions:   readWriteOpt,
				},
				{
					ContainerPath: "/dev/rngd/npu0pe4-5",
					HostPath:      "/dev/rngd/npu0pe4-5",
					Permissions:   readWriteOpt,
				},
				{
					ContainerPath: "/dev/rngd/npu0pe6-7",
					HostPath:      "/dev/rngd/npu0pe6-7",
					Permissions:   readWriteOpt,
				},
				{
					ContainerPath: "/dev/rngd/npu0pe4",
					HostPath:      "/dev/rngd/npu0pe4",
					Permissions:   readWriteOpt,
				},
				{
					ContainerPath: "/dev/rngd/npu0pe5",
					HostPath:      "/dev/rngd/npu0pe5",
					Permissions:   readWriteOpt,
				},
				{
					ContainerPath: "/dev/rngd/npu0pe6",
					HostPath:      "/dev/rngd/npu0pe6",
					Permissions:   readWriteOpt,
				},
				{
					ContainerPath: "/dev/rngd/npu0pe7",
					HostPath:      "/dev/rngd/npu0pe7",
					Permissions:   readWriteOpt,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			actual, err := filterPartitionedDeviceNodes(tc.manifest, tc.partition)
			assert.NoError(t, err)

			assert.Subset(t, actual, tc.mustContains)
		})
	}
}
