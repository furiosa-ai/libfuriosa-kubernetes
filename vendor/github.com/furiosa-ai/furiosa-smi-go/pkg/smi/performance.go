package smi

import (
	"time"

	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi/binding"
)

// PeUtilization represents a PE utilization.
type PeUtilization interface {
	// Core returns a PE core index.
	Core() uint32
	// TimeWindowMill returns time window for utilization.
	TimeWindowMill() uint32
	// PeUsagePercentage returns PE usage percentage.
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

// CoreUtilization represents a core utilization.
type CoreUtilization interface {
	// PeUtilization returns the list of utilizations for each PE cores.
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

// DeviceTemperature represents a temperature information of the device.
type DeviceTemperature interface {
	// SocPeak returns the highest temperature observed from SoC sensors.
	SocPeak() float64
	// Ambient returns the temperature observed from sensors attached to the board.
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

// DevicePerformanceCounter represents a device performance counter.
type DevicePerformanceCounter interface {
	// PerformanceCounter returns a list of performance counters.
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

// PerformanceCounter represents a performance counter.
type PerformanceCounter interface {
	// Timestamp returns timestamp.
	Timestamp() time.Time
	// Core returns a core index.
	Core() uint32
	// CycleCount returns total cycle count in 64-bit unsigned int.
	CycleCount() uint64
	// TaskExecutionCycle returns cycle count used for task execution in 64-bit unsigned int.
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

func (p *performanceCounter) Core() uint32 {
	return p.raw.Core
}

func (p *performanceCounter) CycleCount() uint64 {
	return p.raw.CycleCount
}

func (p *performanceCounter) TaskExecutionCycle() uint64 {
	return p.raw.TaskExecutionCycle
}

func newGovernorProfile(profile binding.FuriosaSmiGovernorProfile) GovernorProfile {
	switch profile {
	case binding.FuriosaSmiGovernorProfileOnDemand:
		return GovernorProfileOnDemand

	case binding.FuriosaSmiGovernorProfilePerformance:
		return GovernorProfilePerformance

	case binding.FuriosaSmiGovernorProfilePowerSave:
		return GovernorProfilePowerSave

	default:
		return GovernorProfileOnDemand
	}
}
