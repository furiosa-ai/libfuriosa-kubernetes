package smi

import "github.com/furiosa-ai/furiosa-smi-go/pkg/smi/binding"

type DeviceErrorInfo interface {
	AxiPostErrorCount() uint32
	AxiFetchErrorCount() uint32
	AxiDiscardErrorCount() uint32
	AxiDoorbellErrorCount() uint32
	PciePostErrorCount() uint32
	PcieFetchErrorCount() uint32
	PcieDiscardErrorCount() uint32
	PcieDoorbellErrorCount() uint32
	DeviceErrorCount() uint32
}

var _ DeviceErrorInfo = new(deviceErrorInfo)

type deviceErrorInfo struct {
	raw binding.FuriosaSmiDeviceErrorInfo
}

func newDeviceErrorInfo(raw binding.FuriosaSmiDeviceErrorInfo) DeviceErrorInfo {
	return &deviceErrorInfo{
		raw: raw,
	}
}

func (d *deviceErrorInfo) AxiPostErrorCount() uint32 {
	return d.raw.AxiPostErrorCount
}

func (d *deviceErrorInfo) AxiFetchErrorCount() uint32 {
	return d.raw.AxiFetchErrorCount
}

func (d *deviceErrorInfo) AxiDiscardErrorCount() uint32 {
	return d.raw.AxiDiscardErrorCount
}

func (d *deviceErrorInfo) AxiDoorbellErrorCount() uint32 {
	return d.raw.AxiDoorbellErrorCount
}

func (d *deviceErrorInfo) PciePostErrorCount() uint32 {
	return d.raw.PciePostErrorCount
}

func (d *deviceErrorInfo) PcieFetchErrorCount() uint32 {
	return d.raw.PcieFetchErrorCount
}

func (d *deviceErrorInfo) PcieDiscardErrorCount() uint32 {
	return d.raw.PcieDiscardErrorCount
}

func (d *deviceErrorInfo) PcieDoorbellErrorCount() uint32 {
	return d.raw.PcieDoorbellErrorCount
}

func (d *deviceErrorInfo) DeviceErrorCount() uint32 {
	return d.raw.DeviceErrorCount
}
