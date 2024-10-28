package smi

import "github.com/furiosa-ai/furiosa-smi-go/pkg/smi/binding"

type PeUtilization interface {
	Core() uint32
	TimeWindowMill() uint32
	PeUsagePercentage() float64
}

var _ PeUtilization = new(peUtilization)

type peUtilization struct {
	raw binding.FuriosaSmiPeUtilization
}

func newPeUtilization(raw binding.FuriosaSmiPeUtilization) PeUtilization {
	return &peUtilization{
		raw: raw,
	}
}

func (p *peUtilization) Core() uint32 {
	return p.raw.Core
}

func (p *peUtilization) TimeWindowMill() uint32 {
	return p.raw.TimeWindowMil
}

func (p *peUtilization) PeUsagePercentage() float64 {
	return p.raw.PeUsagePercentage
}

type MemoryUtilization interface {
	TotalBytes() uint64
	InUseBytes() uint64
}

var _ MemoryUtilization = new(memoryUtilization)

func newMemoryUtilization(raw binding.FuriosaSmiMemoryUtilization) MemoryUtilization {
	return &memoryUtilization{
		raw: raw,
	}
}

type memoryUtilization struct {
	raw binding.FuriosaSmiMemoryUtilization
}

func (m *memoryUtilization) TotalBytes() uint64 {
	return m.raw.TotalBytes
}

func (m *memoryUtilization) InUseBytes() uint64 {
	return m.raw.InUseBytes
}

type DeviceUtilization interface {
	PeUtilization() []PeUtilization
	MemoryUtilization() MemoryUtilization
}

var _ DeviceUtilization = new(deviceUtilization)

type deviceUtilization struct {
	raw binding.FuriosaSmiDeviceUtilization
}

func newDeviceUtilization(raw binding.FuriosaSmiDeviceUtilization) DeviceUtilization {
	return &deviceUtilization{
		raw: raw,
	}
}

func (d *deviceUtilization) PeUtilization() (ret []PeUtilization) {
	for i := uint32(0); i < d.raw.PeCount; i++ {
		ret = append(ret, newPeUtilization(d.raw.Pe[i]))
	}

	return
}

func (d *deviceUtilization) MemoryUtilization() MemoryUtilization {
	return newMemoryUtilization(d.raw.Memory)
}

type DeviceTemperature interface {
	SocPeak() float64
	Ambient() float64
}

var _ DeviceTemperature = new(deviceTemperature)

type deviceTemperature struct {
	raw binding.FuriosaSmiDeviceTemperature
}

func newDeviceTemperature(raw binding.FuriosaSmiDeviceTemperature) DeviceTemperature {
	return &deviceTemperature{
		raw: raw,
	}
}

func (d *deviceTemperature) SocPeak() float64 {
	return d.raw.SocPeak
}

func (d *deviceTemperature) Ambient() float64 {
	return d.raw.Ambient
}
