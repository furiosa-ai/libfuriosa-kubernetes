package smi

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi/binding"
)

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
	case binding.FuriosaSmiGovernorProfilePerformance:
		return GovernorProfilePerformance

	case binding.FuriosaSmiGovernorProfilePowerSave:
		return GovernorProfilePowerSave

	default:
		return GovernorProfilePerformance
	}
}

type performanceCounterMap struct {
	mu   sync.RWMutex
	data map[binding.FuriosaSmiDeviceHandle]PerformanceCounterInfo
}

func newPerformanceCounterMap() performanceCounterMap {
	return performanceCounterMap{
		data: make(map[binding.FuriosaSmiDeviceHandle]PerformanceCounterInfo),
	}
}

func (pcm *performanceCounterMap) get(dev Device) (PerformanceCounterInfo, bool) {
	pcm.mu.RLock()
	defer pcm.mu.RUnlock()

	info, exists := pcm.data[dev.(*device).handle]
	return info, exists
}

func (pcm *performanceCounterMap) set(dev Device, info PerformanceCounterInfo) {
	pcm.mu.Lock()
	defer pcm.mu.Unlock()
	pcm.data[dev.(*device).handle] = info
}

type PerformanceCounterInfo struct {
	beforeCounter DevicePerformanceCounter
	afterCounter  DevicePerformanceCounter
}

type ObserverOpt struct {
	devices  []Device
	interval uint32
}

func NewOptForObserver() (ObserverOpt, error) {
	devices, err := ListDevices()

	if err != nil {
		return ObserverOpt{}, err
	}

	return ObserverOpt{
		devices:  devices,
		interval: 500,
	}, nil
}

func (o *ObserverOpt) SetDevices(devices []Device) {
	o.devices = devices
}

func (o *ObserverOpt) SetInterval(interval uint32) {
	o.interval = interval
}

type observer struct {
	once                  sync.Once
	performanceCounterMap performanceCounterMap
	stopCh                chan struct{}
}

var _ Observer = new(observer)

type Observer interface {
	// GetCoreUtilization returns the core utilization for the given device.
	GetCoreUtilization(device Device) ([]CoreUtilization, error)
	// Destroy stops the observer when it calls explicitly or observer is destroyed by GC.
	Destroy()
}

func newObserverWithOpt(opt ObserverOpt) (Observer, error) {
	devices := opt.devices
	interval := opt.interval

	o := &observer{
		performanceCounterMap: newPerformanceCounterMap(),
		stopCh:                make(chan struct{}),
	}

	o.start(devices, time.Duration(interval)*time.Millisecond)

	runtime.SetFinalizer(o, func(o Observer) {
		o.Destroy()
	})
	return o, nil
}

func (o *observer) isDestroyed() bool {
	select {
	case <-o.stopCh:
		return true
	default:
		return false
	}
}

func (o *observer) GetCoreUtilization(device Device) ([]CoreUtilization, error) {
	if o.isDestroyed() {
		return nil, fmt.Errorf("observer is already destroyed")
	}

	if device == nil {
		return nil, fmt.Errorf("device is nil")
	}

	utilization, err := o.CalculateUtilization(device)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate utilization: %w", err)
	}

	return utilization, nil
}

func (o *observer) start(devices []Device, interval time.Duration) {
	o.updateUtilization(devices)
	time.Sleep(100 * time.Millisecond)
	o.updateUtilization(devices)

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				o.updateUtilization(devices)

			case <-o.stopCh:
				return
			}
		}
	}()
}

func (o *observer) Destroy() {
	o.once.Do(func() {
		close(o.stopCh)
	})
}

func (o *observer) updateUtilization(devices []Device) {
	for _, device := range devices {
		performanceCounter, err := device.DevicePerformanceCounter()

		pc, exists := o.performanceCounterMap.get(device)

		if exists {
			o.performanceCounterMap.set(device, PerformanceCounterInfo{
				beforeCounter: pc.afterCounter,
				afterCounter:  performanceCounter})
		} else {
			o.performanceCounterMap.set(device, PerformanceCounterInfo{
				beforeCounter: performanceCounter,
				afterCounter:  performanceCounter})
		}

		if err != nil {
			continue
		}
	}
}

type CoreUtilization struct {
	Core              uint32
	TimeWindowMil     uint32
	PeUsagePercentage float64
}

func (o *observer) CalculateUtilization(device Device) ([]CoreUtilization, error) {
	performanceCounterInfo, exists := o.performanceCounterMap.get(device)

	if !exists {
		return nil, fmt.Errorf("no performance counter info found for device %v", device)
	}

	beforeCounter := performanceCounterInfo.beforeCounter
	afterCounter := performanceCounterInfo.afterCounter

	utilizationResult := make([]CoreUtilization, 0)

	afterPerfCounter := afterCounter.PerformanceCounter()
	beforePerfCounter := beforeCounter.PerformanceCounter()

	for i, beforePeCounter := range beforePerfCounter {
		afterPeCounter := afterPerfCounter[i]

		if afterPeCounter.CycleCount() < beforePeCounter.CycleCount() {
			utilization := CoreUtilization{
				Core:              beforePeCounter.Core(),
				TimeWindowMil:     0,
				PeUsagePercentage: 0.0,
			}

			utilizationResult = append(utilizationResult, utilization)
			continue
		}

		taskExecutionCycleDiff := afterPeCounter.TaskExecutionCycle() - beforePeCounter.TaskExecutionCycle()
		cycleCountDiff := afterPeCounter.CycleCount() - beforePeCounter.CycleCount()

		peUsagePercentage := safeUsizeDivide(taskExecutionCycleDiff, cycleCountDiff) * 100.0

		utilization := CoreUtilization{
			Core:              beforePeCounter.Core(),
			TimeWindowMil:     uint32(afterPeCounter.Timestamp().Sub(beforePeCounter.Timestamp()).Milliseconds()),
			PeUsagePercentage: peUsagePercentage,
		}

		utilizationResult = append(utilizationResult, utilization)
	}
	return utilizationResult, nil
}

func safeUsizeDivide(fst, snd uint64) float64 {
	if snd == 0 {
		return 0.0
	}
	return float64(fst) / float64(snd)
}

// A type for representing a throttle reason
type ThrottleReason uint32

const (
	// Throttling not active
	ThrottleReasonNone = ThrottleReason(binding.FuriosaSmiThrottleReasonNone)
	// Throttling in idle or unused state
	ThrottleReasonIdle = ThrottleReason(binding.FuriosaSmiThrottleReasonIdle)
	// Throttling triggered by high temperature
	ThrottleReasonThermalSlowdown = ThrottleReason(binding.FuriosaSmiThrottleReasonThermalSlowdown)
	// Throttling due to host-defined power limit
	ThrottleReasonAppPowerCap = ThrottleReason(binding.FuriosaSmiThrottleReasonAppPowerCap)
	// Throttling due to host-defined clock limit
	ThrottleReasonAppClockCap = ThrottleReason(binding.FuriosaSmiThrottleReasonAppClockCap)
	// Throttling from device-internal clock limit
	ThrottleReasonHwClockCap = ThrottleReason(binding.FuriosaSmiThrottleReasonHwClockCap)
	// Throttling from internal bus/NoC bandwidth limit
	ThrottleReasonHwBusLimit = ThrottleReason(binding.FuriosaSmiThrottleReasonHwBusLimit)
	// Throttling from device-enforced power limit
	ThrottleReasonHwPowerCap = ThrottleReason(binding.FuriosaSmiThrottleReasonHwPowerCap)
	// Throttling due to other undefined reasons
	ThrottleReasonOtherReason = ThrottleReason(binding.FuriosaSmiThrottleReasonOtherReason)
)
