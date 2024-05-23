package furiosa_smi_go

// TODO(@bg): support configuration for mock devices(e.g. number of devices and arch)
func GetMockDevices() []Device {
	return []Device{
		&mockDevice{},
		&mockDevice{},
		&mockDevice{},
		&mockDevice{},
		&mockDevice{},
		&mockDevice{},
		&mockDevice{},
		&mockDevice{},
	}
}

var _ Device = new(mockDevice)

type mockDevice struct {
	//binding.FuriosaSmiDeviceInfo{}
	//binding.FuriosaSmiDeviceFiles{}
	//binding.FuriosaSmiCoreStatuses{}
	//binding.FuriosaSmiDeviceErrorInfo{}
	//liveness
	//binding.FuriosaSmiDeviceUtilization{}
	//binding.FuriosaSmiDevicePowerConsumption{}
	//binding.FuriosaSmiDeviceTemperature{}
}

func (m mockDevice) DeviceInfo() (DeviceInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockDevice) DeviceFiles() ([]DeviceFile, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockDevice) CoreStatus() (map[uint32]CoreStatus, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockDevice) DeviceErrorInfo() (DeviceErrorInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockDevice) Liveness() (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockDevice) DeviceUtilization() (DeviceUtilization, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockDevice) PowerConsumption() (uint32, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockDevice) DeviceTemperature() (DeviceTemperature, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockDevice) GetDeviceToDeviceLinkType(target Device) (LinkType, error) {
	//TODO implement me
	panic("implement me")
}
