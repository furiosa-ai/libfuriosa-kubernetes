package furiosa_device

import (
	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
)

type PartitioningPolicy string

const (
	NonePolicy       PartitioningPolicy = "none"
	SingleCorePolicy PartitioningPolicy = "single-core"
	DualCorePolicy   PartitioningPolicy = "dual-core"
	QuadCorePolicy   PartitioningPolicy = "quad-core"
)

// CoreSize returns the number of cores per partition
func (strategy PartitioningPolicy) CoreSize() int {
	switch strategy {
	case SingleCorePolicy:
		return 1

	case DualCorePolicy:
		return 2

	case QuadCorePolicy:
		return 4

	default: // CoreSize should not be used for NonePolicy
		panic("unknown policy")
	}
}

type newDeviceFunc func(originDevice smi.Device, isDisabled bool) ([]FuriosaDevice, error)

func newDeviceFuncResolver(policy PartitioningPolicy) (ret newDeviceFunc) {
	// Note: config validation ensure that there is no exception other than listed strategies.
	switch policy {
	case NonePolicy:
		ret = func(originDevice smi.Device, isDisabled bool) ([]FuriosaDevice, error) {
			newExclusiveDevice, err := newExclusiveDevice(originDevice, isDisabled)
			if err != nil {
				return nil, err
			}

			return []FuriosaDevice{newExclusiveDevice}, nil
		}

	case SingleCorePolicy, DualCorePolicy, QuadCorePolicy:
		ret = func(originDevice smi.Device, isDisabled bool) ([]FuriosaDevice, error) {
			deviceInfo, err := originDevice.DeviceInfo()
			if err != nil {
				return nil, err
			}

			numOfCoresPerPartition := policy.CoreSize()
			totalCores := int(deviceInfo.CoreNum())
			newPartitionedDevices, err := newPartitionedDevices(originDevice, numOfCoresPerPartition, totalCores/numOfCoresPerPartition, isDisabled)
			if err != nil {
				return nil, err
			}

			return newPartitionedDevices, nil
		}

	default:
		panic("unknown policy")
	}

	return ret
}
