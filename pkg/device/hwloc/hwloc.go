package hwloc

/*
#cgo LDFLAGS: -lhwloc
#include <hwloc.h>

// NOTE: we need field accessor since the field name type is reserved keyword in go.
// It is guaranteed that obj is always not NULL.
hwloc_obj_type_t getObjType(hwloc_obj_t obj)  {
    return obj->type;
}
*/
import "C"
import "fmt"

type Hwloc interface {
	TopologyInit() error
	SetIoTypeFilter() error
	TopologyLoad() error
	GetCommonAncestorObjType(dev1BDF string, dev2BDF string) (HwlocObjType, error)
	TopologyDestroy()
}

var _ Hwloc = new(topologyCtx)

type topologyCtx struct {
	topologyContext C.hwloc_topology_t
}

func NewHwloc() Hwloc {
	return &topologyCtx{topologyContext: nil}
}

func (t *topologyCtx) TopologyInit() error {
	return topologyInit(&t.topologyContext)
}

func (t *topologyCtx) SetIoTypeFilter() error {
	return setIoTypeFilter(t.topologyContext)
}

func (t *topologyCtx) TopologyLoad() error {
	return topologyLoad(t.topologyContext)
}

func (t *topologyCtx) GetCommonAncestorObjType(dev1BDF string, dev2BDF string) (HwlocObjType, error) {
	obj1 := getPciDevByBusIDstring(t.topologyContext, dev1BDF)
	if obj1 == nil {
		return HwlocObjTypeUnknown, fmt.Errorf("couldn't find object for %s", dev1BDF)
	}

	obj2 := getPciDevByBusIDstring(t.topologyContext, dev2BDF)
	if obj2 == nil {
		return HwlocObjTypeUnknown, fmt.Errorf("couldn't find object for %s", dev2BDF)
	}

	ancestor := getCommonAncestorObj(t.topologyContext, obj1, obj2)
	if ancestor == nil {
		return HwlocObjTypeUnknown, fmt.Errorf("couldn't find common ancestor for %s and %s", dev1BDF, dev2BDF)
	}

	return getHwlocObjType(ancestor), nil
}

func (t *topologyCtx) TopologyDestroy() {
	topologyDestroy(t.topologyContext)
}

func getHwlocObjType(obj C.hwloc_obj_t) HwlocObjType {
	if obj == nil {
		return HwlocObjTypeUnknown
	}

	switch C.getObjType(obj) {
	case C.HWLOC_OBJ_MACHINE:
		return HwlocObjTypeMachine
	case C.HWLOC_OBJ_PACKAGE:
		return HwlocObjTypePackage
	case C.HWLOC_OBJ_CORE:
		return HwlocObjTypeCore
	case C.HWLOC_OBJ_PU:
		return HwlocObjTypePU
	case C.HWLOC_OBJ_L1CACHE:
		return HwlocObjTypeL1Cache
	case C.HWLOC_OBJ_L2CACHE:
		return HwlocObjTypeL2Cache
	case C.HWLOC_OBJ_L3CACHE:
		return HwlocObjTypeL3Cache
	case C.HWLOC_OBJ_L4CACHE:
		return HwlocObjTypeL4Cache
	case C.HWLOC_OBJ_L5CACHE:
		return HwlocObjTypeL5Cache
	case C.HWLOC_OBJ_L1ICACHE:
		return HwlocObjTypeL1iCache
	case C.HWLOC_OBJ_L2ICACHE:
		return HwlocObjTypeL2iCache
	case C.HWLOC_OBJ_L3ICACHE:
		return HwlocObjTypeL3iCache
	case C.HWLOC_OBJ_GROUP:
		return HwlocObjTypeGroup
	case C.HWLOC_OBJ_NUMANODE:
		return HwlocObjTypeNUMANode
	case C.HWLOC_OBJ_BRIDGE:
		return HwlocObjTypeBridge
	case C.HWLOC_OBJ_PCI_DEVICE:
		return HwlocObjTypePCIDev
	case C.HWLOC_OBJ_OS_DEVICE:
		return HwlocObjTypeOSDev
	case C.HWLOC_OBJ_MISC:
		return HwlocObjTypeMisc
	}
	return HwlocObjTypeUnknown
}
