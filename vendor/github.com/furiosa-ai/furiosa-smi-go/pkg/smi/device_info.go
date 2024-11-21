package smi

import "github.com/furiosa-ai/furiosa-smi-go/pkg/smi/binding"

// DeviceInfo represents a device information.
type DeviceInfo interface {
	// Index returns an index number of the device based on hardware topology.
	Index() uint32
	// Arch returns an architecture of device.
	Arch() Arch
	// CoreNum returns the number of PE cores.
	CoreNum() uint32
	// NumaNode returns a numa node of device.
	NumaNode() uint32
	// Name returns a name of device.
	Name() string
	// Serial returns a serial of device.
	Serial() string
	// UUID returns an uuid of device.
	UUID() string
	// BDF returns a bdf of device.
	BDF() string
	// Major returns a major part of pci device.
	Major() uint16
	// Minor returns a minor part of pci device.
	Minor() uint16
	// FirmwareVersion returns a firmware version of device.
	FirmwareVersion() VersionInfo
	// PertVersion returns a PERT version of device.
	PertVersion() VersionInfo
}

var _ DeviceInfo = new(deviceInfo)

type deviceInfo struct {
	raw binding.FuriosaSmiDeviceInfo
}

func newDeviceInfo(raw binding.FuriosaSmiDeviceInfo) DeviceInfo {
	return &deviceInfo{
		raw: raw,
	}
}

func (d *deviceInfo) Index() uint32 {
	return d.raw.Index
}

func (d *deviceInfo) Arch() Arch {
	return Arch(d.raw.Arch)
}

func (d *deviceInfo) CoreNum() uint32 {
	return d.raw.CoreNum
}

func (d *deviceInfo) NumaNode() uint32 {
	return d.raw.NumaNode
}

func (d *deviceInfo) Name() string {
	return byteBufferToString(d.raw.Name[:])
}

func (d *deviceInfo) Serial() string {
	return byteBufferToString(d.raw.Serial[:])
}

func (d *deviceInfo) UUID() string {
	return byteBufferToString(d.raw.Uuid[:])
}

func (d *deviceInfo) BDF() string {
	return byteBufferToString(d.raw.Bdf[:])
}

func (d *deviceInfo) Major() uint16 {
	return d.raw.Major
}

func (d *deviceInfo) Minor() uint16 {
	return d.raw.Minor
}

func (d *deviceInfo) FirmwareVersion() VersionInfo {
	return newVersionInfo(d.raw.FirmwareVersion)
}

func (d *deviceInfo) PertVersion() VersionInfo {
	return newVersionInfo(d.raw.PertVersion)
}
