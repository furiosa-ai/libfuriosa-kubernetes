package device

import "sort"

const (
	defaultDevFsPath = "/dev"
	defaultSysFsPath = "/sys"
)

type DeviceLister interface {
	// ListDevices lists all Furiosa NPU devices in the system.
	ListDevices() ([]Device, error)
}

var _ DeviceLister = new(deviceLister)

type deviceLister struct {
	devFs           string
	sysFs           string
	deviceValidator deviceValidateFunc
}

// TODO(@bg): pass logger
func NewDeviceLister() DeviceLister {
	return newDeviceLister(defaultDevFsPath, defaultSysFsPath, defaultDeviceValidator)
}

func newDeviceLister(devFs string, sysFs string, deviceValidator deviceValidateFunc) DeviceLister {
	return &deviceLister{
		devFs:           devFs,
		sysFs:           sysFs,
		deviceValidator: deviceValidator,
	}
}

func (d *deviceLister) ListDevices() ([]Device, error) {
	var devices []Device

	devFiles, err := ListDevFs(d.devFs)
	if err != nil {
		return nil, err
	}

	//TODO(@bg): run filterDevFiles with up to 4~5goroutines
	//https://pkg.go.dev/golang.org/x/sync/errgroup
	for deviceIdx, paths := range filterDevFiles(devFiles, d.deviceValidator) {
		device, err := NewDevice(deviceIdx, paths, d.devFs, d.sysFs)
		if err != nil {
			//NOTE:(@bg): this function returns only successfully parsed device as same as rust implementation.
			//https://github.com/furiosa-ai/device-api/blob/5037256c6d00903e8955cdbcd39f4a5b59290623/device-api/src/blocking.rs#L70-L74
			//TODO(@bg): log errors
			continue
		}

		devices = append(devices, device)
	}

	sort.SliceStable(devices, func(i, j int) bool {
		return devices[i].DeviceIndex() < devices[j].DeviceIndex()
	})

	return devices, nil
}
