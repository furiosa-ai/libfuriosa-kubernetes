package smi

import (
	"fmt"
	"math"

	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi/binding"
)

type PcieInfo interface {
	// DeviceInfo returns PCIe device information.
	DeviceInfo() PcieDeviceInfo
	// LinkInfo returns PCIe link information.
	LinkInfo() PcieLinkInfo
	// SriovInfo returns SR-IOV information.
	SriovInfo() SriovInfo
	// RootComplexInfo returns PCIe Root Complex information.
	RootComplexInfo() PcieRootComplexInfo
	// SwitchInfo returns PCIe switch information if available.
	SwitchInfo() PcieSwitchInfo
}

var _ PcieInfo = new(pcieInfo)

type pcieInfo struct {
	pcieDeviceInfo      PcieDeviceInfo
	pcieLinkInfo        PcieLinkInfo
	sriovInfo           SriovInfo
	pcieRootComplexInfo PcieRootComplexInfo
	pcieSwitchInfo      PcieSwitchInfo
}

func newPcieInfo(pcieDeviceInfo PcieDeviceInfo,
	pcieLinkInfo PcieLinkInfo,
	sriovInfo SriovInfo,
	pcieRootComplexInfo PcieRootComplexInfo,
	pcieSwitchInfo PcieSwitchInfo) PcieInfo {
	return &pcieInfo{
		pcieDeviceInfo:      pcieDeviceInfo,
		pcieLinkInfo:        pcieLinkInfo,
		sriovInfo:           sriovInfo,
		pcieRootComplexInfo: pcieRootComplexInfo,
		pcieSwitchInfo:      pcieSwitchInfo,
	}
}

func (p *pcieInfo) DeviceInfo() PcieDeviceInfo {
	return p.pcieDeviceInfo
}

func (p *pcieInfo) LinkInfo() PcieLinkInfo {
	return p.pcieLinkInfo
}

func (p *pcieInfo) SriovInfo() SriovInfo {
	return p.sriovInfo
}

func (p *pcieInfo) RootComplexInfo() PcieRootComplexInfo {
	return p.pcieRootComplexInfo
}

func (p *pcieInfo) SwitchInfo() PcieSwitchInfo {
	if p.pcieSwitchInfo == nil {
		return nil
	}
	return p.pcieSwitchInfo
}

type PcieDeviceInfo interface {
	// DeviceId returns device id.
	DeviceId() uint16
	// VendorId returns vendor id.
	VendorId() uint16
	// SubsystemId returns subsystem device id.
	SubsystemId() uint16
	// RevisionId returns revision id.
	RevisionId() uint8
	// ClassId returns class id.
	ClassId() uint8
	// SubClassId returns subclass id.
	SubClassId() uint8
}

var _ PcieDeviceInfo = new(pcieDeviceInfo)

type pcieDeviceInfo struct {
	raw binding.FuriosaSmiPcieDeviceInfo
}

func newPcieDeviceInfo(raw binding.FuriosaSmiPcieDeviceInfo) PcieDeviceInfo {
	return &pcieDeviceInfo{
		raw: raw,
	}
}

func (p *pcieDeviceInfo) DeviceId() uint16 {
	return p.raw.DeviceId
}

func (p *pcieDeviceInfo) VendorId() uint16 {
	return p.raw.SubsystemVendorId
}

func (p *pcieDeviceInfo) SubsystemId() uint16 {
	return p.raw.SubsystemDeviceId
}

func (p *pcieDeviceInfo) RevisionId() uint8 {
	return p.raw.RevisionId
}

func (p *pcieDeviceInfo) ClassId() uint8 {
	return p.raw.ClassId
}

func (p *pcieDeviceInfo) SubClassId() uint8 {
	return p.raw.SubClassId
}

type PcieLinkInfo interface {
	// PcieGenStatus returns PCIe generation status.
	PcieGenStatus() uint8
	// PcieWidthStatus returns link width status.
	LinkWidthStatus() uint32
	// PcieSpeedStatus returns link speed status in GT/s.
	LinkSpeedStatus() float64
	// MaxLinkGenCapability returns maximum link generation capability.
	MaxLinkWidthCapability() uint32
	// MaxLinkSpeedCapability returns maximum link speed capability in GT/s.
	MaxLinkSpeedCapability() float64
}

var _ PcieLinkInfo = new(pcieLinkInfo)

type pcieLinkInfo struct {
	raw binding.FuriosaSmiPcieLinkInfo
}

func newPcieLinkInfo(raw binding.FuriosaSmiPcieLinkInfo) PcieLinkInfo {
	return &pcieLinkInfo{
		raw: raw,
	}
}

func (p *pcieLinkInfo) PcieGenStatus() uint8 {
	return p.raw.PcieGenStatus
}

func (p *pcieLinkInfo) LinkWidthStatus() uint32 {
	return p.raw.LinkWidthStatus
}

func (p *pcieLinkInfo) LinkSpeedStatus() float64 {
	return p.raw.LinkSpeedStatus
}

func (p *pcieLinkInfo) MaxLinkWidthCapability() uint32 {
	return p.raw.MaxLinkWidthCapability
}

func (p *pcieLinkInfo) MaxLinkSpeedCapability() float64 {
	return p.raw.MaxLinkSpeedCapability
}

type SriovInfo interface {
	// SriovTotalVfs returns the total number of VFs
	SriovTotalVfs() uint32
	// SriovEnabledVfs returns the number of enabled VFs
	SriovEnabledVfs() uint32
}

var _ SriovInfo = new(sriovInfo)

type sriovInfo struct {
	raw binding.FuriosaSmiSriovInfo
}

func newSriovInfo(raw binding.FuriosaSmiSriovInfo) SriovInfo {
	return &sriovInfo{
		raw: raw,
	}
}

func (s *sriovInfo) SriovTotalVfs() uint32 {
	return s.raw.SriovTotalVfs
}

func (s *sriovInfo) SriovEnabledVfs() uint32 {
	return s.raw.SriovEnabledVfs
}

type PcieRootComplexInfo interface {
	// Domain returns domain information.
	Domain() uint16
	// Bus returns bus information.
	Bus() uint8
	// String returns a string representation of the PCIe root complex information in BDF format.
	String() string
}

var _ PcieRootComplexInfo = new(pcieRootComplexInfo)

type pcieRootComplexInfo struct {
	raw binding.FuriosaSmiPcieRootComplexInfo
}

func newPcieRootComplexInfo(raw binding.FuriosaSmiPcieRootComplexInfo) PcieRootComplexInfo {
	return &pcieRootComplexInfo{
		raw: raw,
	}
}

func (p *pcieRootComplexInfo) Domain() uint16 {
	return p.raw.Domain
}

func (p *pcieRootComplexInfo) Bus() uint8 {
	return p.raw.Bus
}

func (p *pcieRootComplexInfo) String() string {
	return fmt.Sprintf("%04x:%02x", p.Domain(), p.Bus())
}

type PcieSwitchInfo interface {
	// Domain returns domain information.
	Domain() uint16
	// Bus returns bus information.
	Bus() uint8
	// Device returns device information.
	Device() uint8
	// Function returns function information.
	Function() uint8
	// String returns a string representation of the PCIe switch information in BDF format.
	String() string
}

var _ PcieSwitchInfo = new(pcieSwitchInfo)

type pcieSwitchInfo struct {
	raw binding.FuriosaSmiPcieSwitchInfo
}

func newPcieSwitchInfo(raw binding.FuriosaSmiPcieSwitchInfo) PcieSwitchInfo {
	if raw.Domain == math.MaxUint16 && raw.Bus == math.MaxUint8 && raw.Device == math.MaxUint8 && raw.Function == math.MaxUint8 {
		return nil
	} else {
		return &pcieSwitchInfo{
			raw: raw,
		}
	}
}

func (p *pcieSwitchInfo) Domain() uint16 {
	return p.raw.Domain
}

func (p *pcieSwitchInfo) Bus() uint8 {
	return p.raw.Bus
}

func (p *pcieSwitchInfo) Device() uint8 {
	return p.raw.Device
}

func (p *pcieSwitchInfo) Function() uint8 {
	return p.raw.Function
}

func (p *pcieSwitchInfo) String() string {
	return fmt.Sprintf("%04x:%02x:%02x.%d",
		p.Domain(),
		p.Bus(),
		p.Device(),
		p.Function(),
	)
}
