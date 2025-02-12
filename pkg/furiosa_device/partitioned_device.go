package furiosa_device

import (
	"fmt"
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/cdi_spec"
	"strconv"
	"tags.cncf.io/container-device-interface/specs-go"

	"github.com/bradfitz/iter"
	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
)

var _ FuriosaDevice = (*partitionedDevice)(nil)

// e.g. If UUID is a3e78042-9cc7-4344-9541-d2d3ffd28106 and Partition is 0-1,
// DeviceID should be "a3e78042-9cc7-4344-9541-d2d3ffd28106_cores_0-1".
const deviceIdDelimiter = "_cores_"

// Partition contains the index of PE cores.
// If partition has PE cores from 0 to 4, `Start` will be 0 and `End` wil be 4.
type Partition struct {
	Start int
	End   int
}

func (p Partition) String() string {
	if p.Start == p.End {
		return strconv.Itoa(p.Start)
	}

	return fmt.Sprintf("%d-%d", p.Start, p.End)
}

type partitionedDevice struct {
	index      int
	origin     smi.Device
	renderer   cdi_spec.Renderer
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

// newPartitionedDevices returns list of FuriosaDevice based on given config.ResourceUnitStrategy.
func newPartitionedDevices(originDevice smi.Device, numOfCoresPerPartition int, numOfPartitions int, isDisabled bool) ([]FuriosaDevice, error) {
	uuid, pciBusID, numaNode, originIndex, err := parseDeviceInfo(originDevice)
	if err != nil {
		return nil, err
	}

	partitionedDevices := make([]FuriosaDevice, 0)
	for partitionIndex := range iter.N(numOfPartitions) {
		partition := Partition{
			Start: partitionIndex * numOfCoresPerPartition,
			End:   (partitionIndex+1)*numOfCoresPerPartition - 1,
		}

		partitionedManifest, err := cdi_spec.NewPartitionedDeviceSpecRenderer(originDevice, partition.Start, partition.End)
		if err != nil {
			return nil, err
		}

		partitionedDevices = append(partitionedDevices, &partitionedDevice{
			index:      generateIndexForPartitionedDevice(originIndex, partitionIndex, numOfPartitions),
			origin:     originDevice,
			renderer:   partitionedManifest,
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

func (p *partitionedDevice) CDISpec() (*specs.Device, error) {
	renderer, err := cdi_spec.NewPartitionedDeviceSpecRenderer(p.origin, p.partition.Start, p.partition.End)
	if err != nil {
		return nil, err
	}
	return renderer.Render(), nil
}

func (p *partitionedDevice) Index() int {
	return p.index
}
