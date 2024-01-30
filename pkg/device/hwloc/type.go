package hwloc

type HwlocObjType string

const (
	HwlocObjTypeMachine  HwlocObjType = "Machine"
	HwlocObjTypePackage  HwlocObjType = "Package"
	HwlocObjTypeCore     HwlocObjType = "Core"
	HwlocObjTypePU       HwlocObjType = "PU"
	HwlocObjTypeL1Cache  HwlocObjType = "L1Cache"
	HwlocObjTypeL2Cache  HwlocObjType = "L2Cache"
	HwlocObjTypeL3Cache  HwlocObjType = "L3Cache"
	HwlocObjTypeL4Cache  HwlocObjType = "L4Cache"
	HwlocObjTypeL5Cache  HwlocObjType = "L5Cache"
	HwlocObjTypeL1iCache HwlocObjType = "L1iCache"
	HwlocObjTypeL2iCache HwlocObjType = "L2iCache"
	HwlocObjTypeL3iCache HwlocObjType = "L3iCache"
	HwlocObjTypeGroup    HwlocObjType = "Group"
	HwlocObjTypeNUMANode HwlocObjType = "NUMANode"
	HwlocObjTypeBridge   HwlocObjType = "Bridge"
	HwlocObjTypePCIDev   HwlocObjType = "PCIDev"
	HwlocObjTypeOSDev    HwlocObjType = "OSDev"
	HwlocObjTypeMisc     HwlocObjType = "Misc"
	HwlocObjTypeUnknown  HwlocObjType = "Unknown"
)
