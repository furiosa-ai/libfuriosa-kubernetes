package smi

import (
	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi/binding"
)

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

type CoreStatuses interface {
	// PeStatus returns a core status of the device.
	PeStatus() []PeStatus
}

var _ CoreStatuses = new(coreStatuses)

type coreStatuses struct {
	raw binding.FuriosaSmiCoreStatuses
}

func newCoreStatuses(raw binding.FuriosaSmiCoreStatuses) CoreStatuses {
	return &coreStatuses{
		raw: raw,
	}
}

func (c *coreStatuses) PeStatus() (ret []PeStatus) {
	for i := uint32(0); i < c.raw.Count; i++ {
		ret = append(ret, newPeStatus(c.raw.CoreStatus[i]))
	}

	return
}

// PeStatus represents a device core status.
type PeStatus interface {
	// Core returns a core index.
	Core() uint32
	// Status returns a core status.
	Status() CoreStatus
}

var _ PeStatus = new(peStatus)

type peStatus struct {
	raw binding.FuriosaSmiPeStatus
}

func newPeStatus(raw binding.FuriosaSmiPeStatus) PeStatus {
	return &peStatus{
		raw: raw,
	}
}

func (p *peStatus) Core() uint32 {
	return p.raw.Core
}

func (p *peStatus) Status() CoreStatus {
	return CoreStatus(p.raw.Status)
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

type PeFrequency interface {
	Core() uint32
	Frequency() uint32
}

var _ PeFrequency = new(peFrequency)

type peFrequency struct {
	raw binding.FuriosaSmiPeFrequency
}

func newPeFrequency(raw binding.FuriosaSmiPeFrequency) PeFrequency {
	return &peFrequency{
		raw: raw,
	}
}

func (p *peFrequency) Core() uint32 {
	return p.raw.Core
}

func (p *peFrequency) Frequency() uint32 {
	return p.raw.Frequency
}

type CoreFrequency interface {
	PeFrequency() []PeFrequency
}

var _ CoreFrequency = new(coreFrequency)

type coreFrequency struct {
	raw binding.FuriosaSmiCoreFrequency
}

func newCoreFrequency(raw binding.FuriosaSmiCoreFrequency) CoreFrequency {
	return &coreFrequency{
		raw: raw,
	}
}

func (c *coreFrequency) PeFrequency() (ret []PeFrequency) {
	for i := uint32(0); i < c.raw.PeCount; i++ {
		ret = append(ret, newPeFrequency(c.raw.Pe[i]))
	}

	return
}

type MemoryFrequency interface {
	Frequency() uint32
}

var _ MemoryFrequency = new(memoryFrequency)

type memoryFrequency struct {
	raw binding.FuriosaSmiMemoryFrequency
}

func newMemoryFrequency(raw binding.FuriosaSmiMemoryFrequency) MemoryFrequency {
	return &memoryFrequency{
		raw: raw,
	}
}

func (m *memoryFrequency) Frequency() uint32 {
	return m.raw.Frequency
}

// GovernorProfile Represents a governor profile
type GovernorProfile uint32

const (
	// GovernorProfileOnDemand governor profile
	GovernorProfileOnDemand = GovernorProfile(binding.FuriosaSmiGovernorProfileOnDemand)
	// GovernorProfilePerformance governor profile
	GovernorProfilePerformance = GovernorProfile(binding.FuriosaSmiGovernorProfilePerformance)
	// GovernorProfilePowerSave governor profile
	GovernorProfilePowerSave = GovernorProfile(binding.FuriosaSmiGovernorProfilePowerSave)
)

func (p GovernorProfile) String() string {
	switch p {
	case GovernorProfileOnDemand:
		return "OnDemand"

	case GovernorProfilePerformance:
		return "Performance"

	case GovernorProfilePowerSave:
		return "PowerSave"

	default: // should not reach here!
		return "Unknown"
	}
}
