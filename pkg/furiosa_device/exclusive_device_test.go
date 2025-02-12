package furiosa_device

import (
	"testing"

	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"github.com/stretchr/testify/assert"
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
			exclusiveDev, err := newExclusiveDevice(tc.mockDevice, false)
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
			exclusiveDev, err := newExclusiveDevice(tc.mockDevice, false)
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
			exclusiveDev, err := newExclusiveDevice(tc.mockDevice, false)
			assert.NoError(t, err)

			actualResult := exclusiveDev.NUMANode()
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
			exclusiveDev, err := newExclusiveDevice(tc.mockDevice, tc.isDisabled)
			assert.NoError(t, err)

			actualResult, err := exclusiveDev.IsHealthy()
			assert.NoError(t, err)

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
			exclusiveDev, err := newExclusiveDevice(tc.mockDevice, false)
			assert.NoError(t, err)

			actualResult := exclusiveDev.DeviceID()
			assert.Equal(t, tc.expectedResult, actualResult)
		})
	}
}
