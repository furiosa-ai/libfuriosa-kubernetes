package cdi_spec_gen

import (
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/furiosa_device"
	"os"
	"sort"
	"tags.cncf.io/container-device-interface/pkg/cdi"
	"tags.cncf.io/container-device-interface/specs-go"
)

const (
	DefaultStaticDir     = cdi.DefaultStaticDir
	DefaultDynamicDir    = cdi.DefaultDynamicDir
	DefaultSpecFileName  = "furiosa.yaml"
	DefaultPermissions   = 0644
	version              = "0.6.0"
	vendor               = "furiosa.ai"
	class                = "npu"
	aggregatedDeviceName = "all"
)

type Spec interface {
	Raw() *specs.Spec
	Write() error
}

type spec struct {
	root     string
	filename string
	raw      *specs.Spec
}

var _ Spec = (*spec)(nil)

func (s *spec) Raw() *specs.Spec {
	return s.raw
}

func (s *spec) Write() error {
	cache, err := cdi.NewCache(cdi.WithAutoRefresh(false),
		cdi.WithSpecDirs(s.root),
	)

	if err != nil {
		return err
	}

	err = cache.WriteSpec(s.raw, s.filename)
	if err != nil {
		return err
	}

	err = os.Chmod(s.root+"/"+s.filename, DefaultPermissions)
	if err != nil {
		return err
	}

	return nil
}

func NewSpec(opts ...Option) (Spec, error) {
	gen := &specGenerator{
		root:                 DefaultStaticDir,
		filename:             DefaultSpecFileName,
		devices:              nil,
		permissions:          DefaultPermissions,
		groupDevices:         make(map[string]groupDevice),
		withAggregatedDevice: false,
	}

	for _, opt := range opts {
		opt(gen)
	}

	return gen.Build()
}

type groupDevice struct {
	groupDeviceName string
	tenantDevices   []furiosa_device.FuriosaDevice
}

type specGenerator struct {
	root                 string
	filename             string
	permissions          int
	devices              []furiosa_device.FuriosaDevice
	groupDevices         map[string]groupDevice
	withAggregatedDevice bool
}

func mergeDeviceSpec(specName string, devices []furiosa_device.FuriosaDevice) (*specs.Device, error) {
	if len(devices) == 0 {
		return nil, nil
	}

	// sort devices to ascending order by index
	sort.Slice(devices, func(i, j int) bool {
		return devices[i].Index() < devices[j].Index()
	})

	var deviceSpecs []specs.Device

	for _, device := range devices {
		target, err := device.CDISpec()
		if err != nil {
			return nil, err
		}

		deviceSpecs = append(deviceSpecs, *target)
	}

	var mergedDeviceNodes []*specs.DeviceNode
	var mergedMounts []*specs.Mount

	for _, deviceSpec := range deviceSpecs {
		mergedDeviceNodes = append(mergedDeviceNodes, deviceSpec.ContainerEdits.DeviceNodes...)
		mergedMounts = append(mergedMounts, deviceSpec.ContainerEdits.Mounts...)
	}

	// FIXME(@bg): this is the cheapest workaround to merge device specs
	aggregatedDevice := specs.Device{
		Name:        specName,
		Annotations: nil,
		ContainerEdits: specs.ContainerEdits{
			Env:            nil,
			DeviceNodes:    mergedDeviceNodes,
			Hooks:          nil,
			Mounts:         mergedMounts,
			IntelRdt:       nil,
			AdditionalGIDs: nil,
		},
	}

	return &aggregatedDevice, nil
}

func (b *specGenerator) Build() (Spec, error) {
	var deviceSpecs []specs.Device

	// handle native devices
	for _, device := range b.devices {
		deviceSpec, err := device.CDISpec()
		if err != nil {
			return nil, err
		}
		deviceSpecs = append(deviceSpecs, *deviceSpec)
	}

	// handle aggregated device
	if b.withAggregatedDevice {
		merged, err := mergeDeviceSpec(aggregatedDeviceName, b.devices)
		if err != nil {
			return nil, err
		}
		aggregatedDevice := merged
		if aggregatedDevice != nil {
			deviceSpecs = append(deviceSpecs, *aggregatedDevice)
		}
	}

	// handle group devices
	for _, group := range b.groupDevices {
		merged, err := mergeDeviceSpec(group.groupDeviceName, group.tenantDevices)
		if err != nil {
			return nil, err
		}

		if merged != nil {
			deviceSpecs = append(deviceSpecs, *merged)
		}
	}

	// TODO: validate class, vendor, device names with parser ex) parser.ValidateClassName()

	return &spec{
		root:     b.root,
		filename: b.filename,
		raw: &specs.Spec{
			Version:     version,
			Kind:        vendor + "/" + class,
			Annotations: nil,
			Devices:     deviceSpecs,
		},
	}, nil
}

type Option func(*specGenerator)

func WithSpecDirs(specDirs string) Option {
	return func(b *specGenerator) {
		b.root = specDirs
	}
}

func WithDevices(devices ...furiosa_device.FuriosaDevice) Option {
	return func(b *specGenerator) {
		b.devices = devices
	}
}

func WithAggregatedDevice() Option {
	return func(b *specGenerator) {
		b.withAggregatedDevice = true
	}
}

func WithSpecFileName(specFileName string) Option {
	return func(b *specGenerator) {
		b.filename = specFileName
	}
}

func WithFilePermissions(permissions int) Option {
	return func(b *specGenerator) {
		b.permissions = permissions
	}
}

func WithGroupDevice(groupDeviceName string, tenantDevices ...furiosa_device.FuriosaDevice) Option {
	return func(b *specGenerator) {
		b.groupDevices[groupDeviceName] = groupDevice{
			groupDeviceName: groupDeviceName,
			tenantDevices:   tenantDevices,
		}
	}
}
