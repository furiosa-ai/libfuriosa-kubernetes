package furiosa_smi_go

import (
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/furiosa_smi_go/binding"
)

func GetDevices() ([]Device, error) {
	out := binding.FuriosaSmiDeviceHandles{
		Count:         0,
		DeviceHandles: [64]binding.FuriosaSmiDeviceHandle{},
	}

	if ret := binding.FuriosaSmiGetDeviceHandles(&out); ret != binding.FuriosaSmiReturnCodeOk {
		return nil, ToError(ret)
	}

	devices := make([]Device, out.Count)
	for i := 0; i < int(out.Count); i++ {
		devices = append(devices, newDevice(out.DeviceHandles[i]))
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
	PowerConsumption() (uint32, error)
	DeviceTemperature() (DeviceTemperature, error)
	GetDeviceToDeviceLinkType(target Device) (LinkType, error)
}

var _ Device = new(device)

type device struct {
	handle binding.FuriosaSmiDeviceHandle
}

func newDevice(handle binding.FuriosaSmiDeviceHandle) Device {
	return &device{handle: handle}
}

func (d *device) DeviceInfo() (DeviceInfo, error) {
	out := binding.FuriosaSmiDeviceInfo{}
	if ret := binding.FuriosaSmiGetDeviceInfo(d.handle, &out); ret != binding.FuriosaSmiReturnCodeOk {
		return nil, ToError(ret)
	}

	return newDeviceInfo(out), nil
}

func (d *device) DeviceFiles() ([]DeviceFile, error) {
	out := binding.FuriosaSmiDeviceFiles{}

	if ret := binding.FuriosaSmiGetDeviceFiles(d.handle, &out); ret != binding.FuriosaSmiReturnCodeOk {
		return nil, ToError(ret)
	}

	deviceFiles := make([]DeviceFile, out.Count)
	for i := 0; i < int(out.Count); i++ {
		deviceFiles = append(deviceFiles, newDeviceFile(out.DeviceFiles[i]))
	}

	return deviceFiles, nil
}

func (d *device) CoreStatus() (map[uint32]CoreStatus, error) {

	out := binding.FuriosaSmiCoreStatuses{}

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
	out := binding.FuriosaSmiDeviceErrorInfo{}

	if ret := binding.FuriosaSmiGetDeviceErrorInfo(d.handle, &out); ret != binding.FuriosaSmiReturnCodeOk {
		return nil, ToError(ret)
	}

	return newDeviceErrorInfo(out), nil
}

func (d *device) Liveness() (bool, error) {
	out := false

	if ret := binding.FuriosaSmiGetDeviceLiveness(d.handle, &out); ret != binding.FuriosaSmiReturnCodeOk {
		return false, ToError(ret)
	}

	return out, nil
}

func (d *device) DeviceUtilization() (DeviceUtilization, error) {
	out := binding.FuriosaSmiDeviceUtilization{}

	if ret := binding.FuriosaSmiGetDeviceUtilization(d.handle, &out); ret != binding.FuriosaSmiReturnCodeOk {
		return nil, ToError(ret)
	}

	return newDeviceUtilization(out), nil
}

func (d *device) PowerConsumption() (uint32, error) {
	out := binding.FuriosaSmiDevicePowerConsumption{}

	if ret := binding.FuriosaSmiGetDevicePowerConsumption(d.handle, &out); ret != binding.FuriosaSmiReturnCodeOk {
		return 0, ToError(ret)
	}

	return uint32(out.RmsTotal), nil
}

func (d *device) DeviceTemperature() (DeviceTemperature, error) {
	out := binding.FuriosaSmiDeviceTemperature{}

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
