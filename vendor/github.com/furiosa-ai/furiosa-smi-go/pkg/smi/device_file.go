package smi

import "github.com/furiosa-ai/furiosa-smi-go/pkg/smi/binding"

// DeviceFile represents a device file.
type DeviceFile interface {
	// Cores returns a list of core for device file.
	Cores() []uint32
	// Path returns a device file path.
	Path() string
}

var _ DeviceFile = new(deviceFile)

type deviceFile struct {
	raw binding.FuriosaSmiDeviceFile
}

func newDeviceFile(raw binding.FuriosaSmiDeviceFile) DeviceFile {
	return &deviceFile{
		raw: raw,
	}
}

func (d *deviceFile) Cores() []uint32 {
	var cores []uint32

	for i := d.raw.CoreStart; i <= d.raw.CoreEnd; i++ {
		cores = append(cores, i)
	}

	return cores
}

func (d *deviceFile) Path() string {
	return byteBufferToString(d.raw.Path[:])
}
