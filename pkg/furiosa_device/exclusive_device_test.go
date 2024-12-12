package furiosa_device

import (
	"testing"

	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/npu_allocator"
	"github.com/stretchr/testify/assert"
	devicePluginAPIv1Beta1 "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

func TestDeviceID(t *testing.T) {
	tests := []struct {
		description    string
		mockDevice     smi.Device
		expectedResult string
	}{
		{
			description:    "test device id",
			mockDevice:     smi.GetStaticMockDevices(smi.ArchWarboy)[0],
			expectedResult: "A76AAD68-6855-40B1-9E86-D080852D1C80",
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			exclusiveDev, err := NewExclusiveDevice(tc.mockDevice, false)
			assert.NoError(t, err)

			actualResult := exclusiveDev.DeviceID()
			assert.Equal(t, tc.expectedResult, actualResult)
		})
	}
}

func TestPCIBusID(t *testing.T) {
	tests := []struct {
		description    string
		mockDevice     smi.Device
		expectedResult string
	}{
		{
			description:    "test pci bus id1",
			mockDevice:     smi.GetStaticMockDevices(smi.ArchWarboy)[0],
			expectedResult: "27",
		},
		{
			description:    "test pci bus id2",
			mockDevice:     smi.GetStaticMockDevices(smi.ArchWarboy)[1],
			expectedResult: "2a",
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			exclusiveDev, err := NewExclusiveDevice(tc.mockDevice, false)
			assert.NoError(t, err)

			actualResult := exclusiveDev.PCIBusID()
			assert.Equal(t, tc.expectedResult, actualResult)
		})
	}
}

func TestNUMANode(t *testing.T) {
	tests := []struct {
		description    string
		mockDevice     smi.Device
		expectedResult int
		expectError    bool
	}{
		{
			description:    "test numa node 1",
			mockDevice:     smi.GetStaticMockDevices(smi.ArchWarboy)[0],
			expectedResult: 0,
			expectError:    false,
		},
		{
			description:    "test numa node 2",
			mockDevice:     smi.GetStaticMockDevices(smi.ArchWarboy)[4],
			expectedResult: 1,
			expectError:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			exclusiveDev, err := NewExclusiveDevice(tc.mockDevice, false)
			assert.NoError(t, err)

			actualResult := exclusiveDev.NUMANode()
			assert.Equal(t, tc.expectedResult, actualResult)
		})
	}
}

func TestDeviceSpecs(t *testing.T) {
	tests := []struct {
		description    string
		mockDevice     smi.Device
		expectedResult []*devicePluginAPIv1Beta1.DeviceSpec
	}{
		{
			description: "test warboy exclusive device",
			mockDevice:  smi.GetStaticMockDevices(smi.ArchWarboy)[0],
			expectedResult: []*devicePluginAPIv1Beta1.DeviceSpec{
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
		//TODO(@bg): add testcases for rngd and other npu family later
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			exclusiveDev, err := NewExclusiveDevice(tc.mockDevice, false)
			assert.NoError(t, err)

			actualResult := exclusiveDev.DeviceSpecs()
			assert.Equal(t, tc.expectedResult, actualResult)
		})
	}
}

// This function tests the IsHealthy API only in terms of the deny list.
func TestIsHealthy(t *testing.T) {
	tests := []struct {
		description    string
		mockDevice     smi.Device
		isDisabled     bool
		expectedResult bool
	}{
		{
			description:    "test healthy device",
			mockDevice:     smi.GetStaticMockDevices(smi.ArchWarboy)[0],
			isDisabled:     false,
			expectedResult: true,
		},
		{
			description:    "test unhealthy device",
			mockDevice:     smi.GetStaticMockDevices(smi.ArchWarboy)[0],
			isDisabled:     true,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			exclusiveDev, err := NewExclusiveDevice(tc.mockDevice, tc.isDisabled)
			assert.NoError(t, err)

			actualResult, err := exclusiveDev.IsHealthy()
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedResult, actualResult)
		})
	}
}

func TestMounts(t *testing.T) {
	tests := []struct {
		description    string
		mockDevice     smi.Device
		expectedResult []*devicePluginAPIv1Beta1.Mount
	}{
		{
			description: "test warboy mount",
			mockDevice:  smi.GetStaticMockDevices(smi.ArchWarboy)[0],
			expectedResult: []*devicePluginAPIv1Beta1.Mount{
				{
					ContainerPath: "/sys/class/npu_mgmt/npu0_mgmt",
					HostPath:      "/sys/class/npu_mgmt/npu0_mgmt",
					ReadOnly:      true,
				},
				{
					ContainerPath: "/sys/class/npu_mgmt/npu0pe0",
					HostPath:      "/sys/class/npu_mgmt/npu0pe0",
					ReadOnly:      true,
				},
				{
					ContainerPath: "/sys/class/npu_mgmt/npu0pe1",
					HostPath:      "/sys/class/npu_mgmt/npu0pe1",
					ReadOnly:      true,
				},
				{
					ContainerPath: "/sys/class/npu_mgmt/npu0pe0-1",
					HostPath:      "/sys/class/npu_mgmt/npu0pe0-1",
					ReadOnly:      true,
				},
				{
					ContainerPath: "/sys/devices/virtual/npu_mgmt/npu0_mgmt",
					HostPath:      "/sys/devices/virtual/npu_mgmt/npu0_mgmt",
					ReadOnly:      true,
				},
				{
					ContainerPath: "/sys/devices/virtual/npu_mgmt/npu0pe0",
					HostPath:      "/sys/devices/virtual/npu_mgmt/npu0pe0",
					ReadOnly:      true,
				},
				{
					ContainerPath: "/sys/devices/virtual/npu_mgmt/npu0pe1",
					HostPath:      "/sys/devices/virtual/npu_mgmt/npu0pe1",
					ReadOnly:      true,
				},
				{
					ContainerPath: "/sys/devices/virtual/npu_mgmt/npu0pe0-1",
					HostPath:      "/sys/devices/virtual/npu_mgmt/npu0pe0-1",
					ReadOnly:      true,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			exclusiveDev, err := NewExclusiveDevice(tc.mockDevice, false)
			assert.NoError(t, err)

			actualResult := exclusiveDev.Mounts()
			assert.Equal(t, tc.expectedResult, actualResult)
		})
	}
}

func TestID(t *testing.T) {
	tests := []struct {
		description    string
		mockDevice     smi.Device
		expectedResult string
	}{
		{
			description:    "test id",
			mockDevice:     smi.GetStaticMockDevices(smi.ArchWarboy)[0],
			expectedResult: "A76AAD68-6855-40B1-9E86-D080852D1C80",
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			exclusiveDev, err := NewExclusiveDevice(tc.mockDevice, false)
			assert.NoError(t, err)

			actualResult := exclusiveDev.ID()
			assert.Equal(t, tc.expectedResult, actualResult)
		})
	}
}

func TestTopologyHintKey(t *testing.T) {
	tests := []struct {
		description    string
		mockDevice     smi.Device
		expectedResult npu_allocator.TopologyHintKey
	}{
		{
			description:    "test topology hint",
			mockDevice:     smi.GetStaticMockDevices(smi.ArchWarboy)[0],
			expectedResult: "27",
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			exclusiveDev, err := NewExclusiveDevice(tc.mockDevice, false)
			assert.NoError(t, err)

			actualResult := exclusiveDev.TopologyHintKey()
			assert.Equal(t, tc.expectedResult, actualResult)
		})
	}
}

func TestEqual(t *testing.T) {
	tests := []struct {
		description      string
		mockSourceDevice smi.Device
		mockTargetDevice smi.Device
		expected         bool
	}{
		{
			description:      "expect source and target are identical",
			mockSourceDevice: smi.GetStaticMockDevices(smi.ArchWarboy)[0],
			mockTargetDevice: smi.GetStaticMockDevices(smi.ArchWarboy)[0],
			expected:         true,
		},
		{
			description:      "expect source and target are not identical",
			mockSourceDevice: smi.GetStaticMockDevices(smi.ArchWarboy)[0],
			mockTargetDevice: smi.GetStaticMockDevices(smi.ArchWarboy)[1],
			expected:         false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			source, err := NewExclusiveDevice(tc.mockSourceDevice, false)
			assert.NoError(t, err)

			target, err := NewExclusiveDevice(tc.mockTargetDevice, false)
			assert.NoError(t, err)

			actual := source.Equal(target)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
