package binding

/*
#include "../furiosa_smi.h"

static inline char* get_bdf_at(FuriosaSmiDisabledDevices* devices, int index) {
    return devices->bdfs[index];
}

static inline uint32_t get_disabled_devices_count(FuriosaSmiDisabledDevices* devices) {
    return devices->count;
}
*/
import "C"
import (
	"unsafe"
)

type DisabledDevices struct {
	Count uint32
	Bdfs  [64]FuriosaSmiBdf
}

func ConvertDisabledDevices(cDevices *FuriosaSmiDisabledDevices) *DisabledDevices {
	out := &DisabledDevices{}
	out.Count = uint32(cDevices.count)

	out.Count = uint32(C.get_disabled_devices_count((*C.FuriosaSmiDisabledDevices)(cDevices)))

	for i := 0; i < int(out.Count); i++ {
		cBdfPtr := unsafe.Pointer(C.get_bdf_at((*C.FuriosaSmiDisabledDevices)(cDevices), C.int(i)))
		out.Bdfs[i] = *(*FuriosaSmiBdf)(cBdfPtr)
	}

	return out
}
