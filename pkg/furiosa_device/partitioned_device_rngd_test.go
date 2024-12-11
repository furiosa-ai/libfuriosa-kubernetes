package furiosa_device

/*
import (
	"fmt"
	"slices"
	"testing"

	//"github.com/furiosa-ai/furiosa-device-plugin/internal/config"
	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/npu_allocator"
	"github.com/stretchr/testify/assert"
	devicePluginAPIv1Beta1 "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

const (
	totalCoresOfRNGD = 8
)

func TestFinalIndexGeneration_RNGD_PartitionedDevice(t *testing.T) {
	rngdMockDevices := smi.GetStaticMockDevices(smi.ArchRngd)

	tests := []struct {
		description                  string
		strategy                     config.ResourceUnitStrategy
		expectedIndexes              []int
		expectedIndexToDeviceUUIDMap map[int]string // key: index, value: uuid
	}{
		{
			description: "Single Core Strategy",
			strategy:    config.SingleCoreStrategy,
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
			strategy:    config.DualCoreStrategy,
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
			strategy:    config.QuadCoreStrategy,
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
			deviceMgr, _ := NewDeviceManager(smi.ArchRngd, rngdMockDevices, tc.strategy, nil, false)

			furiosaDeviceMap := deviceMgr.(*deviceManager).furiosaDevices
			furiosaDevices := make([]FuriosaDevice, 0, len(furiosaDeviceMap))
			for _, device := range furiosaDeviceMap {
				furiosaDevices = append(furiosaDevices, device)
			}

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
		strategy        config.ResourceUnitStrategy
		expectedResults []string
	}{
		{
			description: "should return a list of RNGD Device ID for single core strategy",
			mockDevice:  rngdMockDevice,
			strategy:    config.SingleCoreStrategy,
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
			description: "should return a list of RNGD Device ID for dual core strategy",
			mockDevice:  rngdMockDevice,
			strategy:    config.DualCoreStrategy,
			expectedResults: []string{
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "0-1"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "2-3"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "4-5"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "6-7"),
			},
		},
		{
			description: "should return a list of RNGD Device ID for quad core strategy",
			mockDevice:  rngdMockDevice,
			strategy:    config.QuadCoreStrategy,
			expectedResults: []string{
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "0-3"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "4-7"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			numOfCoresPerPartition := tc.strategy.CoreSize()
			partitionedDevices, err := NewPartitionedDevices(tc.mockDevice, numOfCoresPerPartition, totalCoresOfRNGD/numOfCoresPerPartition, false)
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
		strategy       config.ResourceUnitStrategy
		expectedResult string
	}{
		{
			description:    "returned devices must have same PCI Bus IDs - RNGD 0",
			mockDevice:     rngdMockDevice0,
			strategy:       config.SingleCoreStrategy,
			expectedResult: rngdMockDevice0PciBusId,
		},
		{
			description:    "returned devices must have same PCI Bus IDs - RNGD 1",
			mockDevice:     rngdMockDevice1,
			strategy:       config.SingleCoreStrategy,
			expectedResult: rngdMockDevice1PciBusId,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			numOfCoresPerPartition := tc.strategy.CoreSize()
			partitionedDevices, err := NewPartitionedDevices(tc.mockDevice, numOfCoresPerPartition, totalCoresOfRNGD/numOfCoresPerPartition, false)
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
		strategy       config.ResourceUnitStrategy
		expectedResult int
	}{
		{
			description:    "returned devices must have same NUMA node - RNGD 0",
			mockDevice:     rngdMockDevice0,
			strategy:       config.SingleCoreStrategy,
			expectedResult: rngdMockDevice0NUMANode,
		},
		{
			description:    "returned devices must have same NUMA node - RNGD 1",
			mockDevice:     rngdMockDevice1,
			strategy:       config.SingleCoreStrategy,
			expectedResult: rngdMockDevice1NUMANode,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			numOfCoresPerPartition := tc.strategy.CoreSize()
			partitionedDevices, err := NewPartitionedDevices(tc.mockDevice, numOfCoresPerPartition, totalCoresOfRNGD/numOfCoresPerPartition, true)
			assert.NoError(t, err)

			expectedNUMANode := tc.expectedResult
			for _, device := range partitionedDevices {
				actualNUMANode := device.NUMANode()
				assert.Equal(t, expectedNUMANode, actualNUMANode)
			}
		})
	}
}

func TestDeviceSpecs_RNGD_PartitionedDevice(t *testing.T) {
	rngdMockDevice := smi.GetStaticMockDevices(smi.ArchRngd)[0]

	rngdExpectedResultCandidatesForSingleCoreStrategy := func() [][]*devicePluginAPIv1Beta1.DeviceSpec {
		candidates := make([][]*devicePluginAPIv1Beta1.DeviceSpec, 0, 8)
		for i := 0; i < 8; i++ {
			candidates = append(candidates, []*devicePluginAPIv1Beta1.DeviceSpec{
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

	rngdExpectedResultCandidatesForDualCoreStrategy := func() [][]*devicePluginAPIv1Beta1.DeviceSpec {
		candidates := make([][]*devicePluginAPIv1Beta1.DeviceSpec, 0, 8)
		for i := 0; i < 8; i += 2 {
			candidates = append(candidates, []*devicePluginAPIv1Beta1.DeviceSpec{
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

	rngdExpectedResultCandidatesForQuadCoreStrategy := func() [][]*devicePluginAPIv1Beta1.DeviceSpec {
		candidates := make([][]*devicePluginAPIv1Beta1.DeviceSpec, 0, 8)
		for i := 0; i < 8; i += 4 {
			candidates = append(candidates, []*devicePluginAPIv1Beta1.DeviceSpec{
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
		strategy                 config.ResourceUnitStrategy
		expectedResultCandidates [][]*devicePluginAPIv1Beta1.DeviceSpec
	}{
		{
			description:              "[SingleCoreStrategy] each RNGD mock device must contains all DeviceSpecs",
			mockDevice:               rngdMockDevice,
			strategy:                 config.SingleCoreStrategy,
			expectedResultCandidates: rngdExpectedResultCandidatesForSingleCoreStrategy,
		},
		{
			description:              "[DualCoreStrategy] each RNGD mock device must contains all DeviceSpecs",
			mockDevice:               rngdMockDevice,
			strategy:                 config.DualCoreStrategy,
			expectedResultCandidates: rngdExpectedResultCandidatesForDualCoreStrategy,
		},
		{
			description:              "[QuadCoreStrategy] each RNGD mock device must contains all DeviceSpecs",
			mockDevice:               rngdMockDevice,
			strategy:                 config.QuadCoreStrategy,
			expectedResultCandidates: rngdExpectedResultCandidatesForQuadCoreStrategy,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			numOfCoresPerPartition := tc.strategy.CoreSize()
			partitionedDevices, err := NewPartitionedDevices(tc.mockDevice, numOfCoresPerPartition, totalCoresOfRNGD/numOfCoresPerPartition, false)
			assert.NoError(t, err)

			assert.Len(t, partitionedDevices, len(tc.expectedResultCandidates))

			for i, device := range partitionedDevices {
				actualResult := device.DeviceSpecs()

				assert.ElementsMatch(t, tc.expectedResultCandidates[i], actualResult)
			}
		})
	}
}

func TestIsHealthy_RNGD_PartitionedDevice(t *testing.T) {
	tests := []struct {
		description     string
		mockDevice      smi.Device
		strategy        config.ResourceUnitStrategy
		isDisabled      bool
		expectedResults bool
	}{
		{
			description:     "Enabled device must be healthy - RNGD",
			mockDevice:      smi.GetStaticMockDevices(smi.ArchRngd)[0],
			strategy:        config.SingleCoreStrategy,
			isDisabled:      true,
			expectedResults: false,
		},
		{
			description:     "Disabled device must be unhealthy - RNGD",
			mockDevice:      smi.GetStaticMockDevices(smi.ArchRngd)[0],
			strategy:        config.SingleCoreStrategy,
			isDisabled:      true,
			expectedResults: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			numOfCoresPerPartition := tc.strategy.CoreSize()
			partitionedDevices, err := NewPartitionedDevices(tc.mockDevice, numOfCoresPerPartition, totalCoresOfRNGD/numOfCoresPerPartition, tc.isDisabled)
			assert.NoError(t, err)

			for _, device := range partitionedDevices {
				actualResult, err := device.IsHealthy()
				assert.NoError(t, err)

				assert.Equal(t, tc.expectedResults, actualResult)
			}
		})
	}
}

func TestMounts_RNGD_PartitionedDevice(t *testing.T) {
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
		strategy        config.ResourceUnitStrategy
		expectedResults []*devicePluginAPIv1Beta1.Mount
	}{
		{
			description:     "each RNGD mock device must contains all Mounts",
			mockDevice:      rngdMockDevice,
			strategy:        config.SingleCoreStrategy,
			expectedResults: rngdMockDeviceMounts,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			numOfCoresPerPartition := tc.strategy.CoreSize()
			partitionedDevices, err := NewPartitionedDevices(tc.mockDevice, numOfCoresPerPartition, totalCoresOfRNGD/numOfCoresPerPartition, false)
			assert.NoError(t, err)

			for _, device := range partitionedDevices {
				actualResults := device.Mounts()
				assert.Equal(t, tc.expectedResults, actualResults)
			}
		})
	}
}

func TestID_RNGD_PartitionedDevice(t *testing.T) {
	rngdMockDevice := smi.GetStaticMockDevices(smi.ArchRngd)[0]
	rngdMockDeviceUUID := "A76AAD68-6855-40B1-9E86-D080852D1C80"

	tests := []struct {
		description     string
		mockDevice      smi.Device
		strategy        config.ResourceUnitStrategy
		expectedResults []string
	}{
		{
			description: "should return a list of RNGD Device ID for single core strategy",
			mockDevice:  rngdMockDevice,
			strategy:    config.SingleCoreStrategy,
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
			description: "should return a list of RNGD Device ID for dual core strategy",
			mockDevice:  rngdMockDevice,
			strategy:    config.DualCoreStrategy,
			expectedResults: []string{
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "0-1"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "2-3"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "4-5"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "6-7"),
			},
		},
		{
			description: "should return a list of RNGD Device ID for quad core strategy",
			mockDevice:  rngdMockDevice,
			strategy:    config.QuadCoreStrategy,
			expectedResults: []string{
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "0-3"),
				fmt.Sprintf("%s%s%s", rngdMockDeviceUUID, deviceIdDelimiter, "4-7"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			numOfCoresPerPartition := tc.strategy.CoreSize()
			partitionedDevices, err := NewPartitionedDevices(tc.mockDevice, numOfCoresPerPartition, totalCoresOfRNGD/numOfCoresPerPartition, false)
			assert.NoError(t, err)

			assert.Lenf(t, partitionedDevices, len(tc.expectedResults), "length of expectedResults and partitioned devices are not equal for strategy %s: expected: %d, got: %d", tc.strategy, len(tc.expectedResults), len(partitionedDevices))

			for i, device := range partitionedDevices {
				expectedId := tc.expectedResults[i]
				actualId := device.ID()

				assert.Equal(t, expectedId, actualId)
			}
		})
	}
}

func TestTopologyHintKey_RNGD_PartitionedDevice(t *testing.T) {
	rngdMockDevice0 := smi.GetStaticMockDevices(smi.ArchRngd)[0]
	rngdMockDevice0PciBusId := "27"

	rngdMockDevice1 := smi.GetStaticMockDevices(smi.ArchRngd)[1]
	rngdMockDevice1PciBusId := "2a"

	tests := []struct {
		description    string
		mockDevice     smi.Device
		strategy       config.ResourceUnitStrategy
		expectedResult npu_allocator.TopologyHintKey
	}{
		{
			description:    "returned devices must have same TopologyHintKeys - RNGD 0",
			mockDevice:     rngdMockDevice0,
			strategy:       config.SingleCoreStrategy,
			expectedResult: npu_allocator.TopologyHintKey(rngdMockDevice0PciBusId),
		},
		{
			description:    "returned devices must have same TopologyHintKeys - RNGD 1",
			mockDevice:     rngdMockDevice1,
			strategy:       config.SingleCoreStrategy,
			expectedResult: npu_allocator.TopologyHintKey(rngdMockDevice1PciBusId),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			numOfCoresPerPartition := tc.strategy.CoreSize()
			partitionedDevices, err := NewPartitionedDevices(tc.mockDevice, numOfCoresPerPartition, totalCoresOfRNGD/numOfCoresPerPartition, false)
			assert.NoError(t, err)

			for _, device := range partitionedDevices {
				actualResult := device.TopologyHintKey()

				assert.Equal(t, tc.expectedResult, actualResult)
			}
		})
	}
}

func TestEqual_RNGD_PartitionedDevice(t *testing.T) {
	tests := []struct {
		description      string
		mockSourceDevice smi.Device
		mockTargetDevice smi.Device
		strategy         config.ResourceUnitStrategy
		expected         bool
	}{
		{
			description:      "expect source and target are identical",
			mockSourceDevice: smi.GetStaticMockDevices(smi.ArchRngd)[0],
			mockTargetDevice: smi.GetStaticMockDevices(smi.ArchRngd)[0],
			strategy:         config.SingleCoreStrategy,
			expected:         true,
		},
		{
			description:      "expect source and target are not identical",
			mockSourceDevice: smi.GetStaticMockDevices(smi.ArchRngd)[0],
			mockTargetDevice: smi.GetStaticMockDevices(smi.ArchRngd)[1],
			strategy:         config.SingleCoreStrategy,
			expected:         false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			numOfCoresPerPartition := tc.strategy.CoreSize()

			sourcePartitionedDevices, err := NewPartitionedDevices(tc.mockSourceDevice, numOfCoresPerPartition, totalCoresOfRNGD/numOfCoresPerPartition, false)
			assert.NoError(t, err)

			targetPartitionedDevices, err := NewPartitionedDevices(tc.mockTargetDevice, numOfCoresPerPartition, totalCoresOfRNGD/numOfCoresPerPartition, false)
			assert.NoError(t, err)

			assert.Len(t, sourcePartitionedDevices, len(targetPartitionedDevices))

			for i := range sourcePartitionedDevices {
				sourceDevice := sourcePartitionedDevices[i]
				targetDevice := targetPartitionedDevices[i]

				actualResult := sourceDevice.Equal(targetDevice)

				assert.Equal(t, tc.expected, actualResult)
			}
		})
	}
}
*/
