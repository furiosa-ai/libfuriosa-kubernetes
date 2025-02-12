package furiosa_device

import (
	"fmt"
	"slices"
	"testing"

	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
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

const (
	totalCoresOfRNGD = 8
)

func TestFinalIndexGeneration_RNGD_PartitionedDevice(t *testing.T) {
	rngdMockDevices := smi.GetStaticMockDevices(smi.ArchRngd)

	tests := []struct {
		description                  string
		strategy                     PartitioningPolicy
		expectedIndexes              []int
		expectedIndexToDeviceUUIDMap map[int]string // key: index, value: uuid
	}{
		{
			description: "Single Core Strategy",
			strategy:    SingleCorePolicy,
			expectedIndexes: func() []int {
				indexes := make([]int, 64)
				for i := range indexes {
					indexes[i] = i
				}

				return indexes
			}(),
			expectedIndexToDeviceUUIDMap: func() map[int]string {
				mapping := make(map[int]string)
				for i := 0; i < 64; i++ {
					deviceInfo, _ := rngdMockDevices[i/8].DeviceInfo()
					mapping[i] = deviceInfo.UUID()
				}

				return mapping
			}(),
		},
		{
			description: "Dual Core Strategy",
			strategy:    DualCorePolicy,
			expectedIndexes: func() []int {
				indexes := make([]int, 32)
				for i := range indexes {
					indexes[i] = i
				}

				return indexes
			}(),
			expectedIndexToDeviceUUIDMap: func() map[int]string {
				mapping := make(map[int]string)
				for i := 0; i < 32; i++ {
					deviceInfo, _ := rngdMockDevices[i/4].DeviceInfo()
					mapping[i] = deviceInfo.UUID()
				}

				return mapping
			}(),
		},
		{
			description: "Quad Core Strategy",
			strategy:    QuadCorePolicy,
			expectedIndexes: func() []int {
				indexes := make([]int, 16)
				for i := range indexes {
					indexes[i] = i
				}

				return indexes
			}(),
			expectedIndexToDeviceUUIDMap: func() map[int]string {
				mapping := make(map[int]string)
				for i := 0; i < 16; i++ {
					deviceInfo, _ := rngdMockDevices[i/2].DeviceInfo()
					mapping[i] = deviceInfo.UUID()
				}

				return mapping
			}(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			furiosaDevices, _ := NewFuriosaDevices(rngdMockDevices, nil, tc.strategy)

			slices.SortFunc(furiosaDevices, func(dev1, dev2 FuriosaDevice) int {
				return dev1.Index() - dev2.Index()
			})

			finalIndexes := make([]int, 0, len(furiosaDevices))
			for _, device := range furiosaDevices {
				finalIndexes = append(finalIndexes, device.Index())
			}

			assert.Equal(t, tc.expectedIndexes, finalIndexes)

			finalIndexToDeviceUUIDMap := make(map[int]string)
			for _, furiosaDevice := range furiosaDevices {
				finalIndexToDeviceUUIDMap[furiosaDevice.Index()] = furiosaDevice.(*partitionedDevice).uuid
			}

			assert.Equal(t, tc.expectedIndexToDeviceUUIDMap, finalIndexToDeviceUUIDMap)
		})
	}
}

func TestDeviceIDs_RNGD_PartitionedDevice(t *testing.T) {
	rngdMockDevice := smi.GetStaticMockDevices(smi.ArchRngd)[0]
	rngdMockDeviceUUID := "A76AAD68-6855-40B1-9E86-D080852D1C80"

	tests := []struct {
		description     string
		mockDevice      smi.Device
		strategy        PartitioningPolicy
		expectedResults []string
	}{
		{
			description: "should return a list of RNGD Device DeviceID for single core strategy",
			mockDevice:  rngdMockDevice,
			strategy:    SingleCorePolicy,
			expectedResults: []string{
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "0"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "1"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "2"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "3"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "4"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "5"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "6"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "7"),
			},
		},
		{
			description: "should return a list of RNGD Device DeviceID for dual core strategy",
			mockDevice:  rngdMockDevice,
			strategy:    DualCorePolicy,
			expectedResults: []string{
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "0-1"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "2-3"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "4-5"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "6-7"),
			},
		},
		{
			description: "should return a list of RNGD Device DeviceID for quad core strategy",
			mockDevice:  rngdMockDevice,
			strategy:    QuadCorePolicy,
			expectedResults: []string{
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "0-3"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "4-7"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			numOfCoresPerPartition := tc.strategy.CoreSize()
			partitionedDevices, err := newPartitionedDevices(tc.mockDevice, numOfCoresPerPartition, totalCoresOfRNGD/numOfCoresPerPartition, false)
			assert.NoError(t, err)

			assert.Lenf(t, partitionedDevices, len(tc.expectedResults), "length of expectedResults and partitioned devices are not equal for strategy %s: expected: %d, got: %d", tc.strategy, len(tc.expectedResults), len(partitionedDevices))

			for i, device := range partitionedDevices {
				expectedDeviceId := tc.expectedResults[i]
				actualDeviceId := device.DeviceID()

				assert.Equal(t, expectedDeviceId, actualDeviceId)
			}
		})
	}
}

func TestPCIBusIDs_RNGD_PartitionedDevice(t *testing.T) {
	rngdMockDevice0 := smi.GetStaticMockDevices(smi.ArchRngd)[0]
	rngdMockDevice0PciBusId := "27"

	rngdMockDevice1 := smi.GetStaticMockDevices(smi.ArchRngd)[1]
	rngdMockDevice1PciBusId := "2a"

	tests := []struct {
		description    string
		mockDevice     smi.Device
		strategy       PartitioningPolicy
		expectedResult string
	}{
		{
			description:    "returned devices must have same PCI Bus IDs - RNGD 0",
			mockDevice:     rngdMockDevice0,
			strategy:       SingleCorePolicy,
			expectedResult: rngdMockDevice0PciBusId,
		},
		{
			description:    "returned devices must have same PCI Bus IDs - RNGD 1",
			mockDevice:     rngdMockDevice1,
			strategy:       SingleCorePolicy,
			expectedResult: rngdMockDevice1PciBusId,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			numOfCoresPerPartition := tc.strategy.CoreSize()
			partitionedDevices, err := newPartitionedDevices(tc.mockDevice, numOfCoresPerPartition, totalCoresOfRNGD/numOfCoresPerPartition, false)
			assert.NoError(t, err)

			expectedPCIBusID := tc.expectedResult
			for _, device := range partitionedDevices {
				actualPCIBusID := device.PCIBusID()
				assert.Equal(t, expectedPCIBusID, actualPCIBusID)
			}
		})
	}
}

func TestNUMANode_RNGD_PartitionedDevice(t *testing.T) {
	rngdMockDevice0 := smi.GetStaticMockDevices(smi.ArchRngd)[0]
	rngdMockDevice0NUMANode := 0

	rngdMockDevice1 := smi.GetStaticMockDevices(smi.ArchRngd)[4]
	rngdMockDevice1NUMANode := 1

	tests := []struct {
		description    string
		mockDevice     smi.Device
		strategy       PartitioningPolicy
		expectedResult int
	}{
		{
			description:    "returned devices must have same NUMA node - RNGD 0",
			mockDevice:     rngdMockDevice0,
			strategy:       SingleCorePolicy,
			expectedResult: rngdMockDevice0NUMANode,
		},
		{
			description:    "returned devices must have same NUMA node - RNGD 1",
			mockDevice:     rngdMockDevice1,
			strategy:       SingleCorePolicy,
			expectedResult: rngdMockDevice1NUMANode,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			numOfCoresPerPartition := tc.strategy.CoreSize()
			partitionedDevices, err := newPartitionedDevices(tc.mockDevice, numOfCoresPerPartition, totalCoresOfRNGD/numOfCoresPerPartition, true)
			assert.NoError(t, err)

			expectedNUMANode := tc.expectedResult
			for _, device := range partitionedDevices {
				actualNUMANode := device.NUMANode()
				assert.Equal(t, expectedNUMANode, actualNUMANode)
			}
		})
	}
}

/*func TestDeviceSpecs_RNGD_PartitionedDevice(t *testing.T) {
	rngdMockDevice := smi.GetStaticMockDevices(smi.ArchRngd)[0]

	rngdExpectedResultCandidatesForSingleCoreStrategy := func() [][]*devicePluginAPIv1Beta1.CDISpec {
		candidates := make([][]*devicePluginAPIv1Beta1.CDISpec, 0, 8)
		for i := 0; i < 8; i++ {
			candidates = append(candidates, []*devicePluginAPIv1Beta1.CDISpec{
				{
					ContainerPath: "/dev/rngd/npu0mgmt",
					HostPath:      "/dev/rngd/npu0mgmt",
					Permissions:   "rw",
				},
				{
					ContainerPath: fmt.Sprintf("/dev/rngd/npu0pe%d", i),
					HostPath:      fmt.Sprintf("/dev/rngd/npu0pe%d", i),
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch0",
					HostPath:      "/dev/rngd/npu0ch0",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch1",
					HostPath:      "/dev/rngd/npu0ch1",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch2",
					HostPath:      "/dev/rngd/npu0ch2",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch3",
					HostPath:      "/dev/rngd/npu0ch3",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch4",
					HostPath:      "/dev/rngd/npu0ch4",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch5",
					HostPath:      "/dev/rngd/npu0ch5",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch6",
					HostPath:      "/dev/rngd/npu0ch6",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch7",
					HostPath:      "/dev/rngd/npu0ch7",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch0r",
					HostPath:      "/dev/rngd/npu0ch0r",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch1r",
					HostPath:      "/dev/rngd/npu0ch1r",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch2r",
					HostPath:      "/dev/rngd/npu0ch2r",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch3r",
					HostPath:      "/dev/rngd/npu0ch3r",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch4r",
					HostPath:      "/dev/rngd/npu0ch4r",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch5r",
					HostPath:      "/dev/rngd/npu0ch5r",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch6r",
					HostPath:      "/dev/rngd/npu0ch6r",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch7r",
					HostPath:      "/dev/rngd/npu0ch7r",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0dmar",
					HostPath:      "/dev/rngd/npu0dmar",
					Permissions:   "rw",
				},
			})
		}

		return candidates
	}()

	rngdExpectedResultCandidatesForDualCoreStrategy := func() [][]*devicePluginAPIv1Beta1.CDISpec {
		candidates := make([][]*devicePluginAPIv1Beta1.CDISpec, 0, 8)
		for i := 0; i < 8; i += 2 {
			candidates = append(candidates, []*devicePluginAPIv1Beta1.CDISpec{
				{
					ContainerPath: "/dev/rngd/npu0mgmt",
					HostPath:      "/dev/rngd/npu0mgmt",
					Permissions:   "rw",
				},
				{
					ContainerPath: fmt.Sprintf("/dev/rngd/npu0pe%d-%d", i, i+1),
					HostPath:      fmt.Sprintf("/dev/rngd/npu0pe%d-%d", i, i+1),
					Permissions:   "rw",
				},
				{
					ContainerPath: fmt.Sprintf("/dev/rngd/npu0pe%d", i),
					HostPath:      fmt.Sprintf("/dev/rngd/npu0pe%d", i),
					Permissions:   "rw",
				},
				{
					ContainerPath: fmt.Sprintf("/dev/rngd/npu0pe%d", i+1),
					HostPath:      fmt.Sprintf("/dev/rngd/npu0pe%d", i+1),
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch0",
					HostPath:      "/dev/rngd/npu0ch0",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch1",
					HostPath:      "/dev/rngd/npu0ch1",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch2",
					HostPath:      "/dev/rngd/npu0ch2",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch3",
					HostPath:      "/dev/rngd/npu0ch3",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch4",
					HostPath:      "/dev/rngd/npu0ch4",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch5",
					HostPath:      "/dev/rngd/npu0ch5",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch6",
					HostPath:      "/dev/rngd/npu0ch6",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch7",
					HostPath:      "/dev/rngd/npu0ch7",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch0r",
					HostPath:      "/dev/rngd/npu0ch0r",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch1r",
					HostPath:      "/dev/rngd/npu0ch1r",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch2r",
					HostPath:      "/dev/rngd/npu0ch2r",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch3r",
					HostPath:      "/dev/rngd/npu0ch3r",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch4r",
					HostPath:      "/dev/rngd/npu0ch4r",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch5r",
					HostPath:      "/dev/rngd/npu0ch5r",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch6r",
					HostPath:      "/dev/rngd/npu0ch6r",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch7r",
					HostPath:      "/dev/rngd/npu0ch7r",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0dmar",
					HostPath:      "/dev/rngd/npu0dmar",
					Permissions:   "rw",
				},
			})
		}

		return candidates
	}()

	rngdExpectedResultCandidatesForQuadCoreStrategy := func() [][]*devicePluginAPIv1Beta1.CDISpec {
		candidates := make([][]*devicePluginAPIv1Beta1.CDISpec, 0, 8)
		for i := 0; i < 8; i += 4 {
			candidates = append(candidates, []*devicePluginAPIv1Beta1.CDISpec{
				{
					ContainerPath: "/dev/rngd/npu0mgmt",
					HostPath:      "/dev/rngd/npu0mgmt",
					Permissions:   "rw",
				},
				{
					ContainerPath: fmt.Sprintf("/dev/rngd/npu0pe%d-%d", i, i+3),
					HostPath:      fmt.Sprintf("/dev/rngd/npu0pe%d-%d", i, i+3),
					Permissions:   "rw",
				},
				{
					ContainerPath: fmt.Sprintf("/dev/rngd/npu0pe%d-%d", i, i+1),
					HostPath:      fmt.Sprintf("/dev/rngd/npu0pe%d-%d", i, i+1),
					Permissions:   "rw",
				},
				{
					ContainerPath: fmt.Sprintf("/dev/rngd/npu0pe%d-%d", i+2, i+3),
					HostPath:      fmt.Sprintf("/dev/rngd/npu0pe%d-%d", i+2, i+3),
					Permissions:   "rw",
				},
				{
					ContainerPath: fmt.Sprintf("/dev/rngd/npu0pe%d", i),
					HostPath:      fmt.Sprintf("/dev/rngd/npu0pe%d", i),
					Permissions:   "rw",
				},
				{
					ContainerPath: fmt.Sprintf("/dev/rngd/npu0pe%d", i+1),
					HostPath:      fmt.Sprintf("/dev/rngd/npu0pe%d", i+1),
					Permissions:   "rw",
				},
				{
					ContainerPath: fmt.Sprintf("/dev/rngd/npu0pe%d", i+2),
					HostPath:      fmt.Sprintf("/dev/rngd/npu0pe%d", i+2),
					Permissions:   "rw",
				},
				{
					ContainerPath: fmt.Sprintf("/dev/rngd/npu0pe%d", i+3),
					HostPath:      fmt.Sprintf("/dev/rngd/npu0pe%d", i+3),
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch0",
					HostPath:      "/dev/rngd/npu0ch0",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch1",
					HostPath:      "/dev/rngd/npu0ch1",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch2",
					HostPath:      "/dev/rngd/npu0ch2",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch3",
					HostPath:      "/dev/rngd/npu0ch3",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch4",
					HostPath:      "/dev/rngd/npu0ch4",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch5",
					HostPath:      "/dev/rngd/npu0ch5",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch6",
					HostPath:      "/dev/rngd/npu0ch6",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch7",
					HostPath:      "/dev/rngd/npu0ch7",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch0r",
					HostPath:      "/dev/rngd/npu0ch0r",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch1r",
					HostPath:      "/dev/rngd/npu0ch1r",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch2r",
					HostPath:      "/dev/rngd/npu0ch2r",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch3r",
					HostPath:      "/dev/rngd/npu0ch3r",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch4r",
					HostPath:      "/dev/rngd/npu0ch4r",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch5r",
					HostPath:      "/dev/rngd/npu0ch5r",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch6r",
					HostPath:      "/dev/rngd/npu0ch6r",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0ch7r",
					HostPath:      "/dev/rngd/npu0ch7r",
					Permissions:   "rw",
				},
				{
					ContainerPath: "/dev/rngd/npu0dmar",
					HostPath:      "/dev/rngd/npu0dmar",
					Permissions:   "rw",
				},
			})
		}

		return candidates
	}()

	tests := []struct {
		description              string
		mockDevice               smi.Device
		strategy                 PartitioningPolicy
		expectedResultCandidates [][]*devicePluginAPIv1Beta1.CDISpec
	}{
		{
			description:              "[SingleCoreStrategy] each RNGD mock device must contains all DeviceSpecs",
			mockDevice:               rngdMockDevice,
			strategy:                 SingleCorePolicy,
			expectedResultCandidates: rngdExpectedResultCandidatesForSingleCoreStrategy,
		},
		{
			description:              "[DualCoreStrategy] each RNGD mock device must contains all DeviceSpecs",
			mockDevice:               rngdMockDevice,
			strategy:                 DualCorePolicy,
			expectedResultCandidates: rngdExpectedResultCandidatesForDualCoreStrategy,
		},
		{
			description:              "[QuadCoreStrategy] each RNGD mock device must contains all DeviceSpecs",
			mockDevice:               rngdMockDevice,
			strategy:                 QuadCorePolicy,
			expectedResultCandidates: rngdExpectedResultCandidatesForQuadCoreStrategy,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			numOfCoresPerPartition := tc.strategy.CoreSize()
			partitionedDevices, err := newPartitionedDevices(tc.mockDevice, numOfCoresPerPartition, totalCoresOfRNGD/numOfCoresPerPartition, false)
			assert.NoError(t, err)

			assert.Len(t, partitionedDevices, len(tc.expectedResultCandidates))

			for i, device := range partitionedDevices {
				actualResult := device.DeviceSpecs()

				assert.ElementsMatch(t, tc.expectedResultCandidates[i], actualResult)
			}
		})
	}
}*/

func TestIsHealthy_RNGD_PartitionedDevice(t *testing.T) {
	tests := []struct {
		description     string
		mockDevice      smi.Device
		strategy        PartitioningPolicy
		isDisabled      bool
		expectedResults bool
	}{
		{
			description:     "Enabled device must be healthy - RNGD",
			mockDevice:      smi.GetStaticMockDevices(smi.ArchRngd)[0],
			strategy:        SingleCorePolicy,
			isDisabled:      true,
			expectedResults: false,
		},
		{
			description:     "Disabled device must be unhealthy - RNGD",
			mockDevice:      smi.GetStaticMockDevices(smi.ArchRngd)[0],
			strategy:        SingleCorePolicy,
			isDisabled:      true,
			expectedResults: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			numOfCoresPerPartition := tc.strategy.CoreSize()
			partitionedDevices, err := newPartitionedDevices(tc.mockDevice, numOfCoresPerPartition, totalCoresOfRNGD/numOfCoresPerPartition, tc.isDisabled)
			assert.NoError(t, err)

			for _, device := range partitionedDevices {
				actualResult, err := device.IsHealthy()
				assert.NoError(t, err)

				assert.Equal(t, tc.expectedResults, actualResult)
			}
		})
	}
}

/*func TestMounts_RNGD_PartitionedDevice(t *testing.T) {
	rngdMockDevice := smi.GetStaticMockDevices(smi.ArchRngd)[0]
	rngdMockDeviceMounts := []*devicePluginAPIv1Beta1.Mount{
		{
			ContainerPath: "/sys/class/rngd_mgmt/rngd!npu0mgmt",
			HostPath:      "/sys/class/rngd_mgmt/rngd!npu0mgmt",
			ReadOnly:      true,
		},
		{
			ContainerPath: "/sys/devices/virtual/rngd_mgmt/rngd!npu0mgmt",
			HostPath:      "/sys/devices/virtual/rngd_mgmt/rngd!npu0mgmt",
			ReadOnly:      true,
		},
	}

	tests := []struct {
		description     string
		mockDevice      smi.Device
		strategy        PartitioningPolicy
		expectedResults []*devicePluginAPIv1Beta1.Mount
	}{
		{
			description:     "each RNGD mock device must contains all Mounts",
			mockDevice:      rngdMockDevice,
			strategy:        SingleCorePolicy,
			expectedResults: rngdMockDeviceMounts,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			numOfCoresPerPartition := tc.strategy.CoreSize()
			partitionedDevices, err := newPartitionedDevices(tc.mockDevice, numOfCoresPerPartition, totalCoresOfRNGD/numOfCoresPerPartition, false)
			assert.NoError(t, err)

			for _, device := range partitionedDevices {
				actualResults := device.Mounts()
				assert.Equal(t, tc.expectedResults, actualResults)
			}
		})
	}
}*/

func TestID_RNGD_PartitionedDevice(t *testing.T) {
	rngdMockDevice := smi.GetStaticMockDevices(smi.ArchRngd)[0]
	rngdMockDeviceUUID := "A76AAD68-6855-40B1-9E86-D080852D1C80"

	tests := []struct {
		description     string
		mockDevice      smi.Device
		strategy        PartitioningPolicy
		expectedResults []string
	}{
		{
			description: "should return a list of RNGD Device DeviceID for single core strategy",
			mockDevice:  rngdMockDevice,
			strategy:    SingleCorePolicy,
			expectedResults: []string{
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "0"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "1"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "2"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "3"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "4"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "5"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "6"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "7"),
			},
		},
		{
			description: "should return a list of RNGD Device DeviceID for dual core strategy",
			mockDevice:  rngdMockDevice,
			strategy:    DualCorePolicy,
			expectedResults: []string{
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "0-1"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "2-3"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "4-5"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "6-7"),
			},
		},
		{
			description: "should return a list of RNGD Device DeviceID for quad core strategy",
			mockDevice:  rngdMockDevice,
			strategy:    QuadCorePolicy,
			expectedResults: []string{
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "0-3"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "4-7"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			numOfCoresPerPartition := tc.strategy.CoreSize()
			partitionedDevices, err := newPartitionedDevices(tc.mockDevice, numOfCoresPerPartition, totalCoresOfRNGD/numOfCoresPerPartition, false)
			assert.NoError(t, err)

			assert.Lenf(t, partitionedDevices, len(tc.expectedResults), "length of expectedResults and partitioned devices are not equal for strategy %s: expected: %d, got: %d", tc.strategy, len(tc.expectedResults), len(partitionedDevices))

			for i, device := range partitionedDevices {
				expectedId := tc.expectedResults[i]
				actualId := device.DeviceID()

				assert.Equal(t, expectedId, actualId)
			}
		})
	}
}
