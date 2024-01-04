package device

import (
	"path/filepath"
)

type DeviceFile interface {
	// Path returns path to the device file.
	Path() string
	// Filename returns the file name (e.g., npu0pe0 for /dev/npu0pe0).
	Filename() string
	// DeviceIndex returns the device index (e.g., 1 for npu1pe0).
	DeviceIndex() uint8
	// CoreRange returns the device index (e.g., 1 for npu1pe0).
	CoreRange() CoreRange
	// Mode return the mode of this device file.
	Mode() DeviceMode
	//TODO: has_intersection
}

var _ DeviceFile = new(deviceFile)

type deviceFile struct {
	index      uint8
	coreRange  coreRange
	path       string
	deviceMode DeviceMode
}

func NewDeviceFile(path string) (DeviceFile, error) {
	filename := filepath.Base(path)
	deviceIndex, coreIndices, err := ParseIndices(filename)
	if err != nil {
		return nil, err
	}

	devFile := deviceFile{
		index: deviceIndex,
		path:  path,
	}

	switch len(coreIndices) {
	case 0:
		devFile.deviceMode = DeviceModeMultiCore
		devFile.coreRange.coreRangeType = CoreRangeTypeAll
		devFile.coreRange.start = 0
		devFile.coreRange.end = 0
	case 1:
		devFile.deviceMode = DeviceModeSingle
		devFile.coreRange.coreRangeType = CoreRangeTypeRange
		devFile.coreRange.start = coreIndices[0]
		devFile.coreRange.end = coreIndices[0]
	default:
		devFile.deviceMode = DeviceModeFusion
		devFile.coreRange.coreRangeType = CoreRangeTypeRange
		devFile.coreRange.start = coreIndices[0]
		devFile.coreRange.end = coreIndices[len(coreIndices)-1]
	}

	if devFile.coreRange.start > devFile.coreRange.end {
		return nil, NewUnrecognizedFileError(path)
	}

	return &devFile, nil
}

func (d *deviceFile) Path() string {
	return d.path
}

func (d *deviceFile) Filename() string {
	return filepath.Base(d.path)
}

func (d *deviceFile) DeviceIndex() uint8 {
	return d.index
}

func (d *deviceFile) CoreRange() CoreRange {
	return &d.coreRange
}

func (d *deviceFile) Mode() DeviceMode {
	return d.deviceMode
}
