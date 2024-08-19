package smi

import "github.com/furiosa-ai/furiosa-smi-go/pkg/smi/binding"

type Arch uint32

const (
	ArchWarboy  = Arch(binding.FuriosaSmiArchWarboy)
	ArchRngd    = Arch(binding.FuriosaSmiArchRngd)
	ArchRngdMax = Arch(binding.FuriosaSmiArchRngdMax)
	ArchRngdS   = Arch(binding.FuriosaSmiArchRngdS)
)

func (a Arch) ToString() string {
	switch a {
	case ArchWarboy:
		return "warboy"
	case ArchRngd:
		return "rngd"
	case ArchRngdMax:
		return "rngd-max"
	case ArchRngdS:
		return "rngd-s"
	default:
		return "unknown"
	}
}

type CoreStatus uint32

const (
	CoreStatusAvailable = CoreStatus(binding.FuriosaSmiCoreStatusAvailable)
	CoreStatusOccupied  = CoreStatus(binding.FuriosaSmiCoreStatusOccupied)
)

type LinkType uint32

const (
	LinkTypeUnknown      = LinkType(binding.FuriosaSmiDeviceToDeviceLinkTypeUnknown)
	LinkTypeInterconnect = LinkType(binding.FuriosaSmiDeviceToDeviceLinkTypeInterconnect)
	LinkTypeCpu          = LinkType(binding.FuriosaSmiDeviceToDeviceLinkTypeCpu)
	LinkTypeHostBridge   = LinkType(binding.FuriosaSmiDeviceToDeviceLinkTypeBridge)
	LinkTypeNoc          = LinkType(binding.FuriosaSmiDeviceToDeviceLinkTypeNoc)
)
