package furiosa_device

import (
	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/manifest"
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/npu_allocator"
	devicePluginAPIv1Beta1 "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

var _ FuriosaDevice = (*exclusiveDevice)(nil)

type exclusiveDevice struct {
	index      int
	origin     smi.Device
	manifest   manifest.Manifest
	deviceID   string
	pciBusID   string
	numaNode   int
	isDisabled bool
}

func NewExclusiveDevice(originDevice smi.Device, isDisabled bool) (FuriosaDevice, error) {
	arch, deviceID, pciBusID, numaNode, originIndex, err := parseDeviceInfo(originDevice)
	if err != nil {
		return nil, err
	}

	var newExclusiveDeviceManifest manifest.Manifest
	var manifestErr error

	switch arch {
	case smi.ArchWarboy:
		newExclusiveDeviceManifest, manifestErr = manifest.NewWarboyManifest(originDevice)
	case smi.ArchRngd:
		newExclusiveDeviceManifest, manifestErr = manifest.NewRngdManifest(originDevice)
	}

	if manifestErr != nil {
		return nil, manifestErr

	}

	return &exclusiveDevice{
		index:      originIndex,
		origin:     originDevice,
		manifest:   newExclusiveDeviceManifest,
		deviceID:   deviceID,
		pciBusID:   pciBusID,
		numaNode:   int(numaNode),
		isDisabled: isDisabled,
	}, nil
}

func (f *exclusiveDevice) DeviceID() string {
	return f.deviceID
}

func (f *exclusiveDevice) PCIBusID() string {
	return f.pciBusID
}

func (f *exclusiveDevice) NUMANode() int {
	return f.numaNode
}

func (f *exclusiveDevice) IsHealthy() (bool, error) {
	//TODO(@bg): use more sophisticated way
	if f.isDisabled {
		return false, nil
	}
	liveness, err := f.origin.Liveness()
	if err != nil {
		return liveness, err
	}
	return liveness, nil
}

func (f *exclusiveDevice) IsExclusiveDevice() bool {
	return true
}

func (f *exclusiveDevice) EnvVars() map[string]string {
	return f.manifest.EnvVars()
}

func (f *exclusiveDevice) Annotations() map[string]string {
	return f.manifest.Annotations()
}

func buildDeviceSpec(node *manifest.DeviceNode) *devicePluginAPIv1Beta1.DeviceSpec {
	return &devicePluginAPIv1Beta1.DeviceSpec{
		ContainerPath: node.ContainerPath,
		HostPath:      node.HostPath,
		Permissions:   node.Permissions,
	}
}

func (f *exclusiveDevice) DeviceSpecs() []*devicePluginAPIv1Beta1.DeviceSpec {
	var deviceSpecs []*devicePluginAPIv1Beta1.DeviceSpec

	for _, deviceNode := range f.manifest.DeviceNodes() {
		deviceSpecs = append(deviceSpecs, buildDeviceSpec(deviceNode))
	}

	return deviceSpecs
}

func (f *exclusiveDevice) Mounts() []*devicePluginAPIv1Beta1.Mount {
	var mounts []*devicePluginAPIv1Beta1.Mount

	for _, mount := range f.manifest.MountPaths() {
		var readOnly = false
		// NOTE(@bg): available options are "nodev", "bind", "noexec" and file permission("ro", "rw", ...).
		// However, device-plugin only consume file permission.
		for _, opt := range mount.Options {
			if opt == readOnlyOpt {
				readOnly = true
				break
			}
		}

		mounts = append(mounts, &devicePluginAPIv1Beta1.Mount{
			ContainerPath: mount.ContainerPath,
			HostPath:      mount.HostPath,
			ReadOnly:      readOnly,
		})
	}

	return mounts
}

func (f *exclusiveDevice) CDIDevices() []*devicePluginAPIv1Beta1.CDIDevice {
	//TODO(@bg): CDI will be supported once libfuriosa-kubernetes is ready for CDI and DRA.
	return nil
}

func (f *exclusiveDevice) Index() int {
	return f.index
}

func (f *exclusiveDevice) ID() string {
	return f.DeviceID()
}

func (f *exclusiveDevice) TopologyHintKey() npu_allocator.TopologyHintKey {
	return npu_allocator.TopologyHintKey(f.PCIBusID())
}

func (f *exclusiveDevice) Equal(target npu_allocator.Device) bool {
	converted, isExclusiveDevice := target.(*exclusiveDevice)
	if !isExclusiveDevice {
		return false
	}

	if f.DeviceID() != converted.DeviceID() {
		return false
	}

	return true
}
