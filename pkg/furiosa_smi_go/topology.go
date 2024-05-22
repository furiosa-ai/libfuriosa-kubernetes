package furiosa_smi_go

import (
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/furiosa_smi_go/binding"
)

type Topology interface {
	GetLinkType(dev1 Device, dev2 Device) (LinkType, error)
}

var _ Topology = new(topology)

type topology struct{}

func (t topology) GetLinkType(dev1 Device, dev2 Device) (LinkType, error) {
	var linkType binding.FuriosaSmiDeviceToDeviceLinkType
	if ret := binding.FuriosaSmiGetDeviceToDeviceLinkType(dev1.(*device).handle, dev2.(*device).handle, &linkType); ret != binding.FuriosaSmiReturnCodeOk {
		return LinkTypeUnknown, ToError(ret)
	}

	return LinkType(linkType), nil
}
