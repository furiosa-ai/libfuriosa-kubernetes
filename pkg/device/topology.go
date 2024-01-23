package device

import "C"
import (
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/device/hwloc"
)

// LinkType represents distance between two devices
type LinkType uint

const (
	// LinkTypeUnknown unknown
	LinkTypeUnknown LinkType = iota
	// LinkTypeCrossCPU two devices are connected across different cpus.
	LinkTypeCrossCPU
	// LinkTypeCPU two devices are connected under the same cpu, it may mean:
	// devices are directly attached to the cpu pcie lane without PCIE switch.
	// devices are attached to different PCIE switches under the same cpu.
	LinkTypeCPU
	// LinkTypeHostBridge two devices are connected under the same PCIE host bridge.
	// Note that this does not guarantee devices are attached to the same PCIE switch.
	// More switches could exist under the host bridge switch.
	LinkTypeHostBridge
)

type Topology interface {
	// GetLinkType queries distance of two devices.
	GetLinkType(device1 Device, device2 Device) LinkType
}

var _ Topology = new(topology)

type topology struct {
	hwlocClient    hwloc.Hwloc
	topologyMatrix map[string]map[string]LinkType
}

func NewTopology(devices []Device) (Topology, error) {
	return newTopology(devices, hwloc.NewHwloc())
}

func NewMockTopology(devices []Device, xmlTopologyPath string) (Topology, error) {
	return newTopology(devices, hwloc.NewMockHwloc(xmlTopologyPath))
}

func (t *topology) getLinkType(dev1BDF string, dev2BDF string) (LinkType, error) {
	commonAncestorObjType, err := t.hwlocClient.GetCommonAncestorObjType(dev1BDF, dev2BDF)
	if err != nil {
		return LinkTypeUnknown, err
	}

	switch commonAncestorObjType {
	case hwloc.HwlocObjTypeMachine:
		return LinkTypeCrossCPU, nil
	case hwloc.HwlocObjTypePackage:
		return LinkTypeCPU, nil
	case hwloc.HwlocObjTypeBridge:
		return LinkTypeHostBridge, nil
	}

	return LinkTypeUnknown, nil
}

func (t *topology) populateTopologyMatrix(Devices []Device) error {
	for _, dev1 := range Devices {
		for _, dev2 := range Devices {
			key1, err := dev1.Busname()
			if err != nil {
				return err
			}

			key2, err := dev2.Busname()
			if err != nil {
				return err
			}

			if key1 == key2 {
				continue
			}

			// order keys to reduce memory usage of two-dimensional map
			if key1 > key2 {
				key1, key2 = key2, key1
			}

			linkType, err := t.getLinkType(key1, key2)
			if err != nil {
				return err
			}

			if _, ok := t.topologyMatrix[key1]; !ok {
				t.topologyMatrix[key1] = make(map[string]LinkType)
			}

			t.topologyMatrix[key1][key2] = linkType
		}
	}
	return nil
}

func newTopology(devices []Device, hwlocClient hwloc.Hwloc) (Topology, error) {
	newTopo := &topology{
		hwlocClient:    hwlocClient,
		topologyMatrix: map[string]map[string]LinkType{},
	}

	err := hwlocClient.TopologyInit()
	defer hwlocClient.TopologyDestroy()
	if err != nil {
		return nil, err
	}

	err = hwlocClient.SetIoTypeFilter()
	if err != nil {
		return nil, err
	}

	err = hwlocClient.TopologyLoad()
	if err != nil {
		return nil, err
	}

	err = newTopo.populateTopologyMatrix(devices)
	if err != nil {
		return nil, err
	}

	return newTopo, nil
}

// GetLinkType takes two devices as parameters and return link type between them.
func (t *topology) GetLinkType(device1 Device, device2 Device) LinkType {
	//NOTE(@bg): ignoring errors is not good pattern, in this case, all devices are validated when topology matrix was populated.
	key1, _ := device1.Busname()
	key2, _ := device2.Busname()

	if key1 > key2 {
		key1, key2 = key2, key1
	}

	linkType, ok := t.topologyMatrix[key1][key2]
	if !ok {
		return LinkTypeUnknown
	}

	return linkType
}
