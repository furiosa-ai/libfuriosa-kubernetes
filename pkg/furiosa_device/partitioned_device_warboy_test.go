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
	totalCoresOfWarboy = 2
)

func TestFinalIndexGeneration_Warboy_PartitionedDevice(t *testing.T) {
	warboyMockDevices := smi.GetStaticMockDevices(smi.ArchWarboy)

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
				indexes := make([]int, 16)
				for i := range indexes {
					indexes[i] = i
				}

				return indexes
			}(),
			expectedIndexToDeviceUUIDMap: func() map[int]string {
				mapping := make(map[int]string)
				for i := 0; i < 16; i++ {
					deviceInfo, _ := warboyMockDevices[i/2].DeviceInfo()
					mapping[i] = deviceInfo.UUID()
				}

				return mapping
			}(),
		},
		{
			description: "Dual Core Strategy",
			strategy:    config.DualCoreStrategy,
			expectedIndexes: func() []int {
				indexes := make([]int, 8)
				for i := range indexes {
					indexes[i] = i
				}

				return indexes
			}(),
			expectedIndexToDeviceUUIDMap: func() map[int]string {
				mapping := make(map[int]string)
				for i := 0; i < 8; i++ {
					deviceInfo, _ := warboyMockDevices[i].DeviceInfo()
					mapping[i] = deviceInfo.UUID()
				}

				return mapping
			}(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			deviceMgr, _ := NewDeviceManager(smi.ArchWarboy, warboyMockDevices, tc.strategy, nil, false)

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

func TestDeviceIDs_Warboy_PartitionedDevice(t *testing.T) {
	warboyMockDevice := smi.GetStaticMockDevices(smi.ArchWarboy)[0]
	warboyMockDeviceUUID := "A76AAD68-6855-40B1-9E86-D080852D1C80"

	tests := []struct {
		description     string
		mockDevice      smi.Device
		strategy        config.ResourceUnitStrategy
		expectedResults []string
	}{
		{
			description: "should return a list of Warboy Device ID for single core strategy",
			mockDevice:  warboyMockDevice,
			strategy:    config.SingleCoreStrategy,
			expectedResults: []string{
				fmt.Sprintf("%s%s%s", warboyMockDeviceUUID, deviceIdDelimiter, "0"),
				fmt.Sprintf("%s%s%s", warboyMockDeviceUUID, deviceIdDelimiter, "1"),
			},
		},
		{
			description: "should return a list of Warboy Device ID for dual core strategy",
			mockDevice:  warboyMockDevice,
			strategy:    config.DualCoreStrategy,
			expectedResults: []string{
				fmt.Sprintf("%s%s%s", warboyMockDeviceUUID, deviceIdDelimiter, "0-1"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			numOfCoresPerPartition := tc.strategy.CoreSize()
			partitionedDevices, err := NewPartitionedDevices(tc.mockDevice, numOfCoresPerPartition, totalCoresOfWarboy/numOfCoresPerPartition, false)
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

func TestPCIBusIDs_Warboy_PartitionedDevice(t *testing.T) {
	warboyMockDevice0 := smi.GetStaticMockDevices(smi.ArchWarboy)[0]
	warboyMockDevice0PciBusId := "27"

	warboyMockDevice1 := smi.GetStaticMockDevices(smi.ArchWarboy)[1]
	warboyMockDevice1PciBusId := "2a"

	tests := []struct {
		description    string
		mockDevice     smi.Device
		strategy       config.ResourceUnitStrategy
		expectedResult string
	}{
		{
			description:    "returned devices must have same PCI Bus IDs - WARBOY 0",
			mockDevice:     warboyMockDevice0,
			strategy:       config.SingleCoreStrategy,
			expectedResult: warboyMockDevice0PciBusId,
		},
		{
			description:    "returned devices must have same PCI Bus IDs - WARBOY 1",
			mockDevice:     warboyMockDevice1,
			strategy:       config.SingleCoreStrategy,
			expectedResult: warboyMockDevice1PciBusId,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			numOfCoresPerPartition := tc.strategy.CoreSize()
			partitionedDevices, err := NewPartitionedDevices(tc.mockDevice, numOfCoresPerPartition, totalCoresOfWarboy/numOfCoresPerPartition, false)
			assert.NoError(t, err)

			expectedPCIBusID := tc.expectedResult
			for _, device := range partitionedDevices {
				actualPCIBusID := device.PCIBusID()

				assert.Equal(t, expectedPCIBusID, actualPCIBusID)
			}
		})
	}
}

func TestNUMANode_Warboy_PartitionedDevice(t *testing.T) {
	warboyMockDevice0 := smi.GetStaticMockDevices(smi.ArchWarboy)[0]
	warboyMockDevice0NUMANode := 0

	warboyMockDevice1 := smi.GetStaticMockDevices(smi.ArchWarboy)[4]
	warboyMockDevice1NUMANode := 1

	tests := []struct {
		description    string
		mockDevice     smi.Device
		strategy       config.ResourceUnitStrategy
		expectedResult int
	}{
		{
			description:    "returned devices must have same NUMA node - WARBOY 0",
			mockDevice:     warboyMockDevice0,
			strategy:       config.SingleCoreStrategy,
			expectedResult: warboyMockDevice0NUMANode,
		},
		{
			description:    "returned devices must have same NUMA node - WARBOY 1",
			mockDevice:     warboyMockDevice1,
			strategy:       config.SingleCoreStrategy,
			expectedResult: warboyMockDevice1NUMANode,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			numOfCoresPerPartition := tc.strategy.CoreSize()
			partitionedDevices, err := NewPartitionedDevices(tc.mockDevice, numOfCoresPerPartition, totalCoresOfWarboy/numOfCoresPerPartition, false)
			assert.NoError(t, err)

			expectedNUMANode := tc.expectedResult
			for _, device := range partitionedDevices {
				actualNUMANode := device.NUMANode()

				assert.Equal(t, expectedNUMANode, actualNUMANode)
			}
		})
	}
}

func TestDeviceSpecs_Warboy_PartitionedDevice(t *testing.T) {
	warboyMockDevice := smi.GetStaticMockDevices(smi.ArchWarboy)[0]

	tests := []struct {
		description              string
		mockDevice               smi.Device
		strategy                 config.ResourceUnitStrategy
		expectedResultCandidates [][]*devicePluginAPIv1Beta1.DeviceSpec
	}{
		{
			description: "[SingleCoreStrategy] each Warboy mock device must contains all DeviceSpecs",
			mockDevice:  warboyMockDevice,
			strategy:    config.SingleCoreStrategy,
			expectedResultCandidates: [][]*devicePluginAPIv1Beta1.DeviceSpec{
				{
					{
						ContainerPath: "/dev/npu0_mgmt",
						HostPath:      "/dev/npu0_mgmt",
						Permissions:   "rw",
					},
					{
						ContainerPath: "/dev/npu0pe0",
						HostPath:      "/dev/npu0pe0",
						Permissions:   "rw",
					},
					{
						ContainerPath: "/dev/npu0ch0",
						HostPath:      "/dev/npu0ch0",
						Permissions:   "rw",
					},
					{
						ContainerPath: "/dev/npu0ch1",
						HostPath:      "/dev/npu0ch1",
						Permissions:   "rw",
					},
					{
						ContainerPath: "/dev/npu0ch2",
						HostPath:      "/dev/npu0ch2",
						Permissions:   "rw",
					},
					{
						ContainerPath: "/dev/npu0ch3",
						HostPath:      "/dev/npu0ch3",
						Permissions:   "rw",
					},
				},
				{
					{
						ContainerPath: "/dev/npu0_mgmt",
						HostPath:      "/dev/npu0_mgmt",
						Permissions:   "rw",
					},
					{
						ContainerPath: "/dev/npu0pe1",
						HostPath:      "/dev/npu0pe1",
						Permissions:   "rw",
					},
					{
						ContainerPath: "/dev/npu0ch0",
						HostPath:      "/dev/npu0ch0",
						Permissions:   "rw",
					},
					{
						ContainerPath: "/dev/npu0ch1",
						HostPath:      "/dev/npu0ch1",
						Permissions:   "rw",
					},
					{
						ContainerPath: "/dev/npu0ch2",
						HostPath:      "/dev/npu0ch2",
						Permissions:   "rw",
					},
					{
						ContainerPath: "/dev/npu0ch3",
						HostPath:      "/dev/npu0ch3",
						Permissions:   "rw",
					},
				},
			},
		},
		{
			description: "[DualCoreStrategy] each Warboy mock device must contains all DeviceSpecs",
			mockDevice:  warboyMockDevice,
			strategy:    config.DualCoreStrategy,
			expectedResultCandidates: [][]*devicePluginAPIv1Beta1.DeviceSpec{
				{
					{
						ContainerPath: "/dev/npu0_mgmt",
						HostPath:      "/dev/npu0_mgmt",
						Permissions:   "rw",
					},
					{
						ContainerPath: "/dev/npu0pe0",
						HostPath:      "/dev/npu0pe0",
						Permissions:   "rw",
					},
					{
						ContainerPath: "/dev/npu0pe1",
						HostPath:      "/dev/npu0pe1",
						Permissions:   "rw",
					},
					{
						ContainerPath: "/dev/npu0pe0-1",
						HostPath:      "/dev/npu0pe0-1",
						Permissions:   "rw",
					},
					{
						ContainerPath: "/dev/npu0ch0",
						HostPath:      "/dev/npu0ch0",
						Permissions:   "rw",
					},
					{
						ContainerPath: "/dev/npu0ch1",
						HostPath:      "/dev/npu0ch1",
						Permissions:   "rw",
					},
					{
						ContainerPath: "/dev/npu0ch2",
						HostPath:      "/dev/npu0ch2",
						Permissions:   "rw",
					},
					{
						ContainerPath: "/dev/npu0ch3",
						HostPath:      "/dev/npu0ch3",
						Permissions:   "rw",
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			numOfCoresPerPartition := tc.strategy.CoreSize()
			partitionedDevices, err := NewPartitionedDevices(tc.mockDevice, numOfCoresPerPartition, totalCoresOfWarboy/numOfCoresPerPartition, false)
			assert.NoError(t, err)

			assert.Lenf(t, partitionedDevices, len(tc.expectedResultCandidates), "%s: expected %d partitioned devices, got %d", tc.description, len(tc.expectedResultCandidates), len(partitionedDevices))

			for i, device := range partitionedDevices {
				actualResult := device.DeviceSpecs()

				assert.ElementsMatch(t, tc.expectedResultCandidates[i], actualResult)
			}
		})
	}
}

func TestIsHealthy_Warboy_PartitionedDevice(t *testing.T) {
	tests := []struct {
		description     string
		mockDevice      smi.Device
		strategy        config.ResourceUnitStrategy
		isDisabled      bool
		expectedResults bool
	}{
		{
			description:     "Enabled device must be healthy - WARBOY",
			mockDevice:      smi.GetStaticMockDevices(smi.ArchWarboy)[0],
			strategy:        config.SingleCoreStrategy,
			isDisabled:      false,
			expectedResults: true,
		},
		{
			description:     "Disabled device must be unhealthy - WARBOY",
			mockDevice:      smi.GetStaticMockDevices(smi.ArchWarboy)[0],
			strategy:        config.SingleCoreStrategy,
			isDisabled:      true,
			expectedResults: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			numOfCoresPerPartition := tc.strategy.CoreSize()
			partitionedDevices, err := NewPartitionedDevices(tc.mockDevice, numOfCoresPerPartition, totalCoresOfWarboy/numOfCoresPerPartition, tc.isDisabled)
			assert.NoError(t, err)

			for _, device := range partitionedDevices {
				actualResult, err := device.IsHealthy()
				assert.NoError(t, err)

				assert.Equal(t, tc.expectedResults, actualResult)
			}
		})
	}
}

func TestID_Warboy_PartitionedDevice(t *testing.T) {
	warboyMockDevice := smi.GetStaticMockDevices(smi.ArchWarboy)[0]
	warboyMockDeviceUUID := "A76AAD68-6855-40B1-9E86-D080852D1C80"

	tests := []struct {
		description     string
		mockDevice      smi.Device
		strategy        config.ResourceUnitStrategy
		expectedResults []string
	}{
		{
			description: "should return a list of Warboy Device ID for single core strategy",
			mockDevice:  warboyMockDevice,
			strategy:    config.SingleCoreStrategy,
			expectedResults: []string{
				fmt.Sprintf("%s%s%s", warboyMockDeviceUUID, deviceIdDelimiter, "0"),
				fmt.Sprintf("%s%s%s", warboyMockDeviceUUID, deviceIdDelimiter, "1"),
			},
		},
		{
			description: "should return a list of Warboy Device ID for dual core strategy",
			mockDevice:  warboyMockDevice,
			strategy:    config.DualCoreStrategy,
			expectedResults: []string{
				fmt.Sprintf("%s%s%s", warboyMockDeviceUUID, deviceIdDelimiter, "0-1"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			numOfCoresPerPartition := tc.strategy.CoreSize()
			partitionedDevices, err := NewPartitionedDevices(tc.mockDevice, numOfCoresPerPartition, totalCoresOfWarboy/numOfCoresPerPartition, false)
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

func TestTopologyHintKey_Warboy_PartitionedDevice(t *testing.T) {
	warboyMockDevice0 := smi.GetStaticMockDevices(smi.ArchWarboy)[0]
	warboyMockDevice0PciBusId := "27"

	warboyMockDevice1 := smi.GetStaticMockDevices(smi.ArchWarboy)[1]
	warboyMockDevice1PciBusId := "2a"

	tests := []struct {
		description    string
		mockDevice     smi.Device
		strategy       config.ResourceUnitStrategy
		expectedResult npu_allocator.TopologyHintKey
	}{
		{
			description:    "returned devices must have same TopologyHintKeys - WARBOY 0",
			mockDevice:     warboyMockDevice0,
			strategy:       config.SingleCoreStrategy,
			expectedResult: npu_allocator.TopologyHintKey(warboyMockDevice0PciBusId),
		},
		{
			description:    "returned devices must have same TopologyHintKeys - WARBOY 1",
			mockDevice:     warboyMockDevice1,
			strategy:       config.SingleCoreStrategy,
			expectedResult: npu_allocator.TopologyHintKey(warboyMockDevice1PciBusId),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			numOfCoresPerPartition := tc.strategy.CoreSize()
			partitionedDevices, err := NewPartitionedDevices(tc.mockDevice, numOfCoresPerPartition, totalCoresOfWarboy/numOfCoresPerPartition, false)
			assert.NoError(t, err)

			for _, device := range partitionedDevices {
				actualResult := device.TopologyHintKey()

				assert.Equal(t, tc.expectedResult, actualResult)
			}
		})
	}
}

func TestEqual_Warboy_PartitionedDevice(t *testing.T) {
	tests := []struct {
		description      string
		mockSourceDevice smi.Device
		mockTargetDevice smi.Device
		strategy         config.ResourceUnitStrategy
		expected         bool
	}{
		{
			description:      "expect source and target are identical",
			mockSourceDevice: smi.GetStaticMockDevices(smi.ArchWarboy)[0],
			mockTargetDevice: smi.GetStaticMockDevices(smi.ArchWarboy)[0],
			strategy:         config.SingleCoreStrategy,
			expected:         true,
		},
		{
			description:      "expect source and target are not identical",
			mockSourceDevice: smi.GetStaticMockDevices(smi.ArchWarboy)[0],
			mockTargetDevice: smi.GetStaticMockDevices(smi.ArchWarboy)[1],
			strategy:         config.SingleCoreStrategy,
			expected:         false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			numOfCoresPerPartition := tc.strategy.CoreSize()

			sourcePartitionedDevices, err := NewPartitionedDevices(tc.mockSourceDevice, numOfCoresPerPartition, totalCoresOfWarboy/numOfCoresPerPartition, false)
			assert.NoError(t, err)

			targetPartitionedDevices, err := NewPartitionedDevices(tc.mockTargetDevice, numOfCoresPerPartition, totalCoresOfWarboy/numOfCoresPerPartition, false)
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
