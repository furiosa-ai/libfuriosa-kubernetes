package furiosa_device

import (
	"testing"

	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"github.com/stretchr/testify/assert"
)

func TestBuildFuriosaDevices(t *testing.T) {
	tests := []struct {
		description           string
		policy                PartitioningPolicy
		expectExclusiveDevice bool
	}{
		{
			description:           "test generic policy",
			policy:                NonePolicy,
			expectExclusiveDevice: true,
		},
		{
			description:           "test single core policy",
			policy:                SingleCorePolicy,
			expectExclusiveDevice: false,
		},
		{
			description:           "test dual core policy",
			policy:                DualCorePolicy,
			expectExclusiveDevice: false,
		},
		{
			description:           "test quad core policy",
			policy:                QuadCorePolicy,
			expectExclusiveDevice: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			devices := smi.GetStaticMockDevices(smi.ArchRngd)

			actualDevices, err := NewFuriosaDevices(devices, nil, tc.policy)
			assert.NoError(t, err)

			for _, actualDevice := range actualDevices {
				if tc.expectExclusiveDevice {
					assert.IsType(t, new(exclusiveDevice), actualDevice)
				} else {
					assert.IsType(t, new(partitionedDevice), actualDevice)
				}
			}
		})
	}
}
