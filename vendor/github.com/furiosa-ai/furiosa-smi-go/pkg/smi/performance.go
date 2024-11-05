package smi

import (
	"time"

	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi/binding"
)

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

type CoreUtilization interface {
	PeUtilization() []PeUtilization
}

var _ CoreUtilization = new(coreUtilization)

type coreUtilization struct {
	raw binding.FuriosaSmiCoreUtilization
}

func newCoreUtilization(raw binding.FuriosaSmiCoreUtilization) CoreUtilization {
	return &coreUtilization{
		raw: raw,
	}
}

func (d *coreUtilization) PeUtilization() (ret []PeUtilization) {
	for i := uint32(0); i < d.raw.PeCount; i++ {
		ret = append(ret, newPeUtilization(d.raw.Pe[i]))
	}

	return
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

type DevicePerformanceCounter interface {
	PerformanceCounter() []PerformanceCounter
}

var _ DevicePerformanceCounter = new(devicePerformanceCounter)

type devicePerformanceCounter struct {
	raw binding.FuriosaSmiDevicePerformanceCounter
}

func newDevicePerformanceCounter(raw binding.FuriosaSmiDevicePerformanceCounter) DevicePerformanceCounter {
	return &devicePerformanceCounter{
		raw: raw,
	}
}

func (d *devicePerformanceCounter) PerformanceCounter() []PerformanceCounter {
	var ret []PerformanceCounter

	for i := uint32(0); i < d.raw.PeCount; i++ {
		ret = append(ret, newPerformanceCounter(d.raw.PePerformanceCounters[i]))
	}

	return ret
}

type PerformanceCounter interface {
	Timestamp() time.Time
	CycleCount() uint64
	TaskExecutionCycle() uint64
}

var _ PerformanceCounter = new(performanceCounter)

type performanceCounter struct {
	raw binding.FuriosaSmiPePerformanceCounter
}

func newPerformanceCounter(raw binding.FuriosaSmiPePerformanceCounter) PerformanceCounter {
	return &performanceCounter{
		raw: raw,
	}
}

func (p *performanceCounter) Timestamp() time.Time {
	return time.Unix(p.raw.Timestamp, 0)
}

func (p *performanceCounter) CycleCount() uint64 {
	return p.raw.CycleCount
}

func (p *performanceCounter) TaskExecutionCycle() uint64 {
	return p.raw.TaskExecutionCycle
}
