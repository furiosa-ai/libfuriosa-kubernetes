package smi

import (
	"runtime"

	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi/binding"
)

type furiosaSmiObserverInstance = *binding.FuriosaSmiObserver

// ListDevices lists all Furiosa NPU devices in the system.
func ListDevices() ([]Device, error) {
	var outDeviceHandle binding.FuriosaSmiDeviceHandles
	if ret := binding.FuriosaSmiGetDeviceHandles(&outDeviceHandle); ret != binding.FuriosaSmiReturnCodeOk {
		return nil, toError(ret)
	}

	var outObserverInstance = new(furiosaSmiObserverInstance)
	if ret := binding.FuriosaSmiCreateObserver(outObserverInstance); ret != binding.FuriosaSmiReturnCodeOk {
		return nil, toError(ret)
	}

	defer runtime.SetFinalizer(outObserverInstance, func(observerInstance *furiosaSmiObserverInstance) {
		_ = binding.FuriosaSmiDestroyObserver(observerInstance)
	})

	var devices []Device
	for i := 0; i < int(outDeviceHandle.Count); i++ {
		devices = append(devices, newDevice(outDeviceHandle.DeviceHandles[i], outObserverInstance))
	}

	return devices, nil
}

// Device represents the abstraction for a single Furiosa NPU device.
type Device interface {
	// DeviceInfo returns `DeviceInfo` which contains information about NPU device. (e.g. arch, serial, ...)
	DeviceInfo() (DeviceInfo, error)
	// DeviceFiles list device files under this device.
	DeviceFiles() ([]DeviceFile, error)
	// CoreStatus examine each core of the device, whether it is occupied or available.
	CoreStatus() (CoreStatuses, error)
	// Liveness returns a liveness state of the device.
	Liveness() (bool, error)
	// CoreFrequency returns a core frequency (MHz) of the device.
	CoreFrequency() (CoreFrequency, error)
	// MemoryFrequency returns a memory frequency (MHz) of the device.
	MemoryFrequency() (MemoryFrequency, error)
	// CoreUtilization returns a core utilization of the device.
	CoreUtilization() (CoreUtilization, error)
	// PowerConsumption returns a power consumption of the device.
	PowerConsumption() (float64, error)
	// DeviceTemperature returns a temperature of the device.
	DeviceTemperature() (DeviceTemperature, error)
	// DeviceToDeviceLinkType returns a device link type between two devices.
	DeviceToDeviceLinkType(target Device) (LinkType, error)
	// P2PAccessible returns whether two devices are p2p accessible each other or not.
	P2PAccessible(target Device) (bool, error)
	// DevicePerformanceCounter returns a performance counter of the device.
	DevicePerformanceCounter() (DevicePerformanceCounter, error)
	// GovernorProfile returns a governor profile of the device.
	GovernorProfile() (GovernorProfile, error)
	// SetGovernorProfile set a governor profile of the device.
	SetGovernorProfile(governorProfile GovernorProfile) error
}

var _ Device = new(device)

type device struct {
	observerInstance *furiosaSmiObserverInstance
	handle           binding.FuriosaSmiDeviceHandle
}

func newDevice(handle binding.FuriosaSmiDeviceHandle, observerInstance *furiosaSmiObserverInstance) Device {
	return &device{
		observerInstance: observerInstance,
		handle:           handle,
	}
}

func (d *device) DeviceInfo() (DeviceInfo, error) {
	var out binding.FuriosaSmiDeviceInfo
	if ret := binding.FuriosaSmiGetDeviceInfo(d.handle, &out); ret != binding.FuriosaSmiReturnCodeOk {
		return nil, toError(ret)
	}

	return newDeviceInfo(out), nil
}

func (d *device) DeviceFiles() ([]DeviceFile, error) {
	var out binding.FuriosaSmiDeviceFiles

	if ret := binding.FuriosaSmiGetDeviceFiles(d.handle, &out); ret != binding.FuriosaSmiReturnCodeOk {
		return nil, toError(ret)
	}

	var deviceFiles []DeviceFile
	for i := 0; i < int(out.Count); i++ {
		deviceFiles = append(deviceFiles, newDeviceFile(out.DeviceFiles[i]))
	}

	return deviceFiles, nil
}

func (d *device) CoreStatus() (CoreStatuses, error) {
	var out binding.FuriosaSmiCoreStatuses

	if ret := binding.FuriosaSmiGetDeviceCoreStatus(d.handle, &out); ret != binding.FuriosaSmiReturnCodeOk {
		return nil, toError(ret)
	}

	return newCoreStatuses(out), nil
}

func (d *device) Liveness() (bool, error) {
	var out bool

	if ret := binding.FuriosaSmiGetDeviceLiveness(d.handle, &out); ret != binding.FuriosaSmiReturnCodeOk {
		return false, toError(ret)
	}

	return out, nil
}

func (d *device) CoreFrequency() (CoreFrequency, error) {
	var out binding.FuriosaSmiCoreFrequency

	if ret := binding.FuriosaSmiGetCoreFrequency(d.handle, &out); ret != binding.FuriosaSmiReturnCodeOk {
		return nil, toError(ret)
	}

	return newCoreFrequency(out), nil
}

func (d *device) MemoryFrequency() (MemoryFrequency, error) {
	var out binding.FuriosaSmiMemoryFrequency

	if ret := binding.FuriosaSmiGetMemoryFrequency(d.handle, &out); ret != binding.FuriosaSmiReturnCodeOk {
		return nil, toError(ret)
	}

	return newMemoryFrequency(out), nil
}

func (d *device) CoreUtilization() (CoreUtilization, error) {
	var out binding.FuriosaSmiCoreUtilization

	if ret := binding.FuriosaSmiGetCoreUtilization(*d.observerInstance, d.handle, &out); ret != binding.FuriosaSmiReturnCodeOk {
		return nil, toError(ret)
	}

	return newCoreUtilization(out), nil
}

func (d *device) PowerConsumption() (float64, error) {
	var out binding.FuriosaSmiDevicePowerConsumption

	if ret := binding.FuriosaSmiGetDevicePowerConsumption(d.handle, &out); ret != binding.FuriosaSmiReturnCodeOk {
		return 0, toError(ret)
	}

	return out.RmsTotal, nil
}

func (d *device) DeviceTemperature() (DeviceTemperature, error) {
	var out binding.FuriosaSmiDeviceTemperature

	if ret := binding.FuriosaSmiGetDeviceTemperature(d.handle, &out); ret != binding.FuriosaSmiReturnCodeOk {
		return nil, toError(ret)
	}

	return newDeviceTemperature(out), nil
}

func (d *device) DeviceToDeviceLinkType(target Device) (LinkType, error) {
	var linkType binding.FuriosaSmiDeviceToDeviceLinkType

	if ret := binding.FuriosaSmiGetDeviceToDeviceLinkType(d.handle, target.(*device).handle, &linkType); ret != binding.FuriosaSmiReturnCodeOk {
		return LinkTypeUnknown, toError(ret)
	}

	return LinkType(linkType), nil
}

func (d *device) P2PAccessible(target Device) (bool, error) {
	var out bool

	if ret := binding.FuriosaSmiGetP2pAccessible(d.handle, target.(*device).handle, &out); ret != binding.FuriosaSmiReturnCodeOk {
		return false, toError(ret)
	}

	return out, nil
}

func (d *device) DevicePerformanceCounter() (DevicePerformanceCounter, error) {
	var out binding.FuriosaSmiDevicePerformanceCounter

	if ret := binding.FuriosaSmiGetDevicePerformanceCounter(d.handle, &out); ret != binding.FuriosaSmiReturnCodeOk {
		return nil, toError(ret)
	}

	return newDevicePerformanceCounter(out), nil
}

func (d *device) GovernorProfile() (GovernorProfile, error) {
	var out binding.FuriosaSmiGovernorProfile

	if ret := binding.FuriosaSmiGetGovernorProfile(d.handle, &out); ret != binding.FuriosaSmiReturnCodeOk {
		return 0, toError(ret)
	}

	return newGovernorProfile(out), nil
}

func (d *device) SetGovernorProfile(profile GovernorProfile) error {
	rawProfile := binding.FuriosaSmiGovernorProfile(profile)
	if ret := binding.FuriosaSmiSetGovernorProfile(d.handle, rawProfile); ret != binding.FuriosaSmiReturnCodeOk {
		return toError(ret)
	}

	return nil
}
