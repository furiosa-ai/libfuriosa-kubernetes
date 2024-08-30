package smi

import (
	"runtime"

	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi/binding"
)

type furiosaSmiObserverInstance = *binding.FuriosaSmiObserver

func ListDevices() ([]Device, error) {
	var outDeviceHandle binding.FuriosaSmiDeviceHandles
	if ret := binding.FuriosaSmiGetDeviceHandles(&outDeviceHandle); ret != binding.FuriosaSmiReturnCodeOk {
		return nil, ToError(ret)
	}

	var outObserverInstance = new(furiosaSmiObserverInstance)
	if ret := binding.FuriosaSmiCreateObserver(outObserverInstance); ret != binding.FuriosaSmiReturnCodeOk {
		return nil, ToError(ret)
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

type Device interface {
	DeviceInfo() (DeviceInfo, error)
	DeviceFiles() ([]DeviceFile, error)
	CoreStatus() (map[uint32]CoreStatus, error)
	DeviceErrorInfo() (DeviceErrorInfo, error)
	Liveness() (bool, error)
	DeviceUtilization() (DeviceUtilization, error)
	PowerConsumption() (float64, error)
	DeviceTemperature() (DeviceTemperature, error)
	GetDeviceToDeviceLinkType(target Device) (LinkType, error)
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
		return nil, ToError(ret)
	}

	return newDeviceInfo(out), nil
}

func (d *device) DeviceFiles() ([]DeviceFile, error) {
	var out binding.FuriosaSmiDeviceFiles

	if ret := binding.FuriosaSmiGetDeviceFiles(d.handle, &out); ret != binding.FuriosaSmiReturnCodeOk {
		return nil, ToError(ret)
	}

	var deviceFiles []DeviceFile
	for i := 0; i < int(out.Count); i++ {
		deviceFiles = append(deviceFiles, newDeviceFile(out.DeviceFiles[i]))
	}

	return deviceFiles, nil
}

func (d *device) CoreStatus() (map[uint32]CoreStatus, error) {
	var out binding.FuriosaSmiCoreStatuses

	if ret := binding.FuriosaSmiGetDeviceCoreStatus(d.handle, &out); ret != binding.FuriosaSmiReturnCodeOk {
		return nil, ToError(ret)
	}

	coreStatusMap := make(map[uint32]CoreStatus)
	for i := 0; i < int(out.Count); i++ {
		coreStatusMap[uint32(i)] = CoreStatus(out.CoreStatus[i])
	}

	return coreStatusMap, nil
}

func (d *device) DeviceErrorInfo() (DeviceErrorInfo, error) {
	var out binding.FuriosaSmiDeviceErrorInfo

	if ret := binding.FuriosaSmiGetDeviceErrorInfo(d.handle, &out); ret != binding.FuriosaSmiReturnCodeOk {
		return nil, ToError(ret)
	}

	return newDeviceErrorInfo(out), nil
}

func (d *device) Liveness() (bool, error) {
	var out bool

	if ret := binding.FuriosaSmiGetDeviceLiveness(d.handle, &out); ret != binding.FuriosaSmiReturnCodeOk {
		return false, ToError(ret)
	}

	return out, nil
}

func (d *device) DeviceUtilization() (DeviceUtilization, error) {
	var out binding.FuriosaSmiDeviceUtilization

	if ret := binding.FuriosaSmiGetDeviceUtilization(*d.observerInstance, d.handle, &out); ret != binding.FuriosaSmiReturnCodeOk {
		return nil, ToError(ret)
	}

	return newDeviceUtilization(out), nil
}

func (d *device) PowerConsumption() (float64, error) {
	var out binding.FuriosaSmiDevicePowerConsumption

	if ret := binding.FuriosaSmiGetDevicePowerConsumption(d.handle, &out); ret != binding.FuriosaSmiReturnCodeOk {
		return 0, ToError(ret)
	}

	return out.RmsTotal, nil
}

func (d *device) DeviceTemperature() (DeviceTemperature, error) {
	var out binding.FuriosaSmiDeviceTemperature

	if ret := binding.FuriosaSmiGetDeviceTemperature(d.handle, &out); ret != binding.FuriosaSmiReturnCodeOk {
		return nil, ToError(ret)
	}

	return newDeviceTemperature(out), nil
}

func (d *device) GetDeviceToDeviceLinkType(target Device) (LinkType, error) {
	var linkType binding.FuriosaSmiDeviceToDeviceLinkType

	if ret := binding.FuriosaSmiGetDeviceToDeviceLinkType(d.handle, target.(*device).handle, &linkType); ret != binding.FuriosaSmiReturnCodeOk {
		return LinkTypeUnknown, ToError(ret)
	}

	return LinkType(linkType), nil
}
