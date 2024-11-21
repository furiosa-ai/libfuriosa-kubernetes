package smi

import "github.com/furiosa-ai/furiosa-smi-go/pkg/smi/binding"

// Arch represents NPU architecture.
type Arch uint32

const (
	// ArchWarboy represents Warboy architecture.
	ArchWarboy = Arch(binding.FuriosaSmiArchWarboy)
	// ArchRngd represents RNGD architecture.
	ArchRngd = Arch(binding.FuriosaSmiArchRngd)
	// ArchRngdMax represents RNGD-Max architecture.
	ArchRngdMax = Arch(binding.FuriosaSmiArchRngdMax)
	// ArchRngdS represents RNGD-S architecture.
	ArchRngdS = Arch(binding.FuriosaSmiArchRngdS)
)

// ToString converts given arch into the string representation.
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

// CoreStatus represents a device core status
type CoreStatus uint32

const (
	// CoreStatusAvailable represents core is available.
	CoreStatusAvailable = CoreStatus(binding.FuriosaSmiCoreStatusAvailable)
	// CoreStatusOccupied represents core is occupied.
	CoreStatusOccupied = CoreStatus(binding.FuriosaSmiCoreStatusOccupied)
)

// LinkType represents a topology link type between 2 NPU devices.
type LinkType uint32

const (
	// LinkTypeUnknown means unknown link type.
	LinkTypeUnknown = LinkType(binding.FuriosaSmiDeviceToDeviceLinkTypeUnknown)
	// LinkTypeInterconnect represents link type under same machine.
	LinkTypeInterconnect = LinkType(binding.FuriosaSmiDeviceToDeviceLinkTypeInterconnect)
	// LinkTypeCpu represents link type under same cpu.
	LinkTypeCpu = LinkType(binding.FuriosaSmiDeviceToDeviceLinkTypeCpu)
	// LinkTypeHostBridge represents link type under same switch.
	LinkTypeHostBridge = LinkType(binding.FuriosaSmiDeviceToDeviceLinkTypeBridge)
	// LinkTypeNoc represents link type under same socket.
	LinkTypeNoc = LinkType(binding.FuriosaSmiDeviceToDeviceLinkTypeNoc)
)
