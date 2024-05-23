package furiosa_smi_go

import (
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/furiosa_smi_go/binding"
)

type Arch uint32

const (
	ArchWarboy  = Arch(binding.FuriosaSmiArchWarboy)
	ArchRngd    = Arch(binding.FuriosaSmiArchRngd)
	ArchRngdMax = Arch(binding.FuriosaSmiArchRngdMax)
	ArchRngdS   = Arch(binding.FuriosaSmiArchRngdS)
)

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
	LinkTypeSoc          = LinkType(binding.FuriosaSmiDeviceToDeviceLinkTypeSoc)
)
