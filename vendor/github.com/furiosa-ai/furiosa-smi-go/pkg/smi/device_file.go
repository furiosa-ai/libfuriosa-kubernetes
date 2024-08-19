package smi

import "github.com/furiosa-ai/furiosa-smi-go/pkg/smi/binding"

type DeviceFile interface {
	Cores() []uint32
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
