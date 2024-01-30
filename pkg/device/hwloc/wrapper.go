package hwloc

/*
#include <hwloc.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

const (
	HwlocSuccess = 0
)

func topologyInit(topologyContext *C.hwloc_topology_t) error {
	err := C.hwloc_topology_init(topologyContext)
	if int(err) != HwlocSuccess {
		return fmt.Errorf("couldn't initialize Hwloc topology handle, got non-zero error code: %d", int(err))
	}

	return nil
}

func setIoTypeFilter(topologyContext C.hwloc_topology_t) error {
	err := C.hwloc_topology_set_io_types_filter(topologyContext, C.HWLOC_TYPE_FILTER_KEEP_IMPORTANT)
	if int(err) != HwlocSuccess {
		return fmt.Errorf("couldn't set filter for io devices, got non-zero error code: %d", int(err))
	}

	return nil
}

func topologyLoad(topologyContext C.hwloc_topology_t) error {
	err := C.hwloc_topology_load(topologyContext)
	if int(err) != HwlocSuccess {
		return fmt.Errorf("couldn't build topology, got non-zero error code: %d", int(err))
	}

	return nil
}

func topologySetXML(topologyContext C.hwloc_topology_t, xmlTopologyPath string) error {
	cStr := C.CString(xmlTopologyPath)
	defer C.free(unsafe.Pointer(cStr))

	err := C.hwloc_topology_set_xml(topologyContext, cStr)
	if int(err) != HwlocSuccess {
		return fmt.Errorf("couldn't set xml topology, got non-zero error code: %d", int(err))
	}

	return nil
}

func getPciDevByBusIDstring(topologyContext C.hwloc_topology_t, bdf string) C.hwloc_obj_t {
	cStr := C.CString(bdf)
	defer C.free(unsafe.Pointer(cStr))
	return C.hwloc_get_pcidev_by_busidstring(topologyContext, cStr)
}

func getCommonAncestorObj(topologyContext C.hwloc_topology_t, obj1 C.hwloc_obj_t, obj2 C.hwloc_obj_t) C.hwloc_obj_t {
	return C.hwloc_get_common_ancestor_obj(topologyContext, obj1, obj2)
}

func topologyDestroy(topologyContext C.hwloc_topology_t) {
	C.hwloc_topology_destroy(topologyContext)
}
