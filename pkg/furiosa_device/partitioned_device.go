package furiosa_device

import (
	"fmt"
	"strconv"

	"github.com/bradfitz/iter"
	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/manifest"
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/npu_allocator"
	devicePluginAPIv1Beta1 "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

var _ FuriosaDevice = (*partitionedDevice)(nil)

// e.g. If UUID is a3e78042-9cc7-4344-9541-d2d3ffd28106 and Partition is 0-1,
// DeviceID should be "a3e78042-9cc7-4344-9541-d2d3ffd28106_cores_0-1".
const deviceIdDelimiter = "_cores_"

// Partition contains the index of PE cores.
// If partition has PE cores from 0 to 4, `start` will be 0 and `end` wil be 4.
type Partition struct {
	start int
	end   int
}

func (p Partition) String() string {
	if p.start == p.end {
		return strconv.Itoa(p.start)
	}

	return fmt.Sprintf("%d-%d", p.start, p.end)
}

type partitionedDevice struct {
	index      int
	origin     smi.Device
	manifest   manifest.Manifest
	uuid       string
	partition  Partition
	pciBusID   string
	numaNode   int
	isDisabled bool
}

// generateIndexForPartitionedDevice generated final index value for Partitioned Device
func generateIndexForPartitionedDevice(originalIndex, partitionIndex, partitionsLength int) int {
	return originalIndex*partitionsLength + partitionIndex
}

// NewPartitionedDevices returns list of FuriosaDevice based on given config.ResourceUnitStrategy.
func NewPartitionedDevices(originDevice smi.Device, numOfCoresPerPartition int, numOfPartitions int, isDisabled bool) ([]FuriosaDevice, error) {
	arch, uuid, pciBusID, numaNode, originIndex, err := parseDeviceInfo(originDevice)
	if err != nil {
		return nil, err
	}

	// This block checks architecture and gets manifest of it.
	// If architecture is invalid, it returns error.
	var originalManifest manifest.Manifest
	switch arch {
	case smi.ArchWarboy:
		if originalManifest, err = manifest.NewWarboyManifest(originDevice); err != nil {
			return nil, err
		}

	case smi.ArchRngd:
		if originalManifest, err = manifest.NewRngdManifest(originDevice); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported architecture: %s", arch.ToString())
	}

	partitionedDevices := make([]FuriosaDevice, 0)
	for partitionIndex := range iter.N(numOfPartitions) {
		partition := Partition{
			start: partitionIndex * numOfCoresPerPartition,
			end:   (partitionIndex+1)*numOfCoresPerPartition - 1,
		}

		partitionedManifest, err := NewPartitionedDeviceManifest(arch, originalManifest, partition)
		if err != nil {
			return nil, err
		}

		partitionedDevices = append(partitionedDevices, &partitionedDevice{
			index:      generateIndexForPartitionedDevice(originIndex, partitionIndex, numOfPartitions),
			origin:     originDevice,
			manifest:   partitionedManifest,
			uuid:       uuid,
			partition:  partition,
			pciBusID:   pciBusID,
			numaNode:   int(numaNode),
			isDisabled: isDisabled,
		})
	}

	return partitionedDevices, nil
}

func (p *partitionedDevice) DeviceID() string {
	// e.g. If UUID is a3e78042-9cc7-4344-9541-d2d3ffd28106 and Partition is 0-1,
	// DeviceID should be "a3e78042-9cc7-4344-9541-d2d3ffd28106_cores_0-1".
	return fmt.Sprintf("%s%s%s", p.uuid, deviceIdDelimiter, p.partition.String())
}

func (p *partitionedDevice) PCIBusID() string {
	return p.pciBusID
}

func (p *partitionedDevice) NUMANode() int {
	return p.numaNode
}

func (p *partitionedDevice) IsHealthy() (bool, error) {
	if p.isDisabled {
		return false, nil
	}

	liveness, err := p.origin.Liveness()
	if err != nil {
		return liveness, err
	}

	return liveness, nil
}

func (p *partitionedDevice) IsExclusiveDevice() bool {
	return false
}

func (p *partitionedDevice) EnvVars() map[string]string {
	return p.manifest.EnvVars()
}

func (p *partitionedDevice) Annotations() map[string]string {
	return p.manifest.Annotations()
}

func (p *partitionedDevice) DeviceSpecs() []*devicePluginAPIv1Beta1.DeviceSpec {
	deviceSpecs := make([]*devicePluginAPIv1Beta1.DeviceSpec, 0)
	for _, deviceNode := range p.manifest.DeviceNodes() {
		deviceSpecs = append(deviceSpecs, &devicePluginAPIv1Beta1.DeviceSpec{
			ContainerPath: deviceNode.ContainerPath,
			HostPath:      deviceNode.HostPath,
			Permissions:   deviceNode.Permissions,
		})
	}

	return deviceSpecs
}

func (p *partitionedDevice) Mounts() []*devicePluginAPIv1Beta1.Mount {
	mounts := make([]*devicePluginAPIv1Beta1.Mount, 0)
	for _, mount := range p.manifest.MountPaths() {
		readOnly := false
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

func (p *partitionedDevice) CDIDevices() []*devicePluginAPIv1Beta1.CDIDevice {
	// TODO: implement it when CDI is ready.
	return nil
}

func (p *partitionedDevice) Index() int {
	return p.index
}

func (p *partitionedDevice) ID() string {
	return p.DeviceID()
}

func (p *partitionedDevice) TopologyHintKey() npu_allocator.TopologyHintKey {
	return npu_allocator.TopologyHintKey(p.PCIBusID())
}

func (p *partitionedDevice) Equal(target npu_allocator.Device) bool {
	converted, isPartitionedDevice := target.(*partitionedDevice)
	if !isPartitionedDevice {
		return false
	}

	if p.DeviceID() != converted.DeviceID() {
		return false
	}

	return true
}
