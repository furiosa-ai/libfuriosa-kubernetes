package device

import "C"
import (
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/device/hwloc"
)

// LinkType represents distance between two devices in integer.
// The value will be used to calculate total score of device set to select the best NPU pairs.
type LinkType uint

const (
	// LinkTypeUnknown unknown
	LinkTypeUnknown LinkType = 0
	// LinkTypeInterconnect two devices are connected across different cpus through interconnect.
	LinkTypeInterconnect LinkType = 10
	// LinkTypeCPU two devices are connected under the same cpu, it may mean:
	// devices are directly attached to the cpu pcie lane without PCIE switch.
	// devices are attached to different PCIE switches under the same cpu.
	LinkTypeCPU LinkType = 20
	// LinkTypeHostBridge two devices are connected under the same PCIE host bridge.
	// Note that this does not guarantee devices are attached to the same PCIE switch.
	// More switches could exist under the host bridge switch.
	LinkTypeHostBridge LinkType = 30

	// NOTE(@bg): Score 40 and 50 is reserved for LinkTypeMultiSwitch and LinkTypeSingleSwitch.
	// NOTE(@bg): Score 60 is reserved for LinkTypeBoard

	// LinkTypeSoc two devices are on the same Soc chip.
	LinkTypeSoc LinkType = 70
)

type Topology interface {
	// GetLinkType queries distance of two devices.
	GetLinkType(dev1BDF string, dev2BDF string) LinkType
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
	if dev1BDF == dev2BDF {
		return LinkTypeSoc, nil
	}

	commonAncestorObjType, err := t.hwlocClient.GetCommonAncestorObjType(dev1BDF, dev2BDF)
	if err != nil {
		return LinkTypeUnknown, err
	}

	switch commonAncestorObjType {
	case hwloc.HwlocObjTypeMachine:
		return LinkTypeInterconnect, nil
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
func (t *topology) GetLinkType(dev1BDF string, dev2BDF string) LinkType {
	if dev1BDF > dev2BDF {
		dev1BDF, dev2BDF = dev2BDF, dev1BDF
	}

	linkType, ok := t.topologyMatrix[dev1BDF][dev2BDF]
	if !ok {
		return LinkTypeUnknown
	}

	return linkType
}
