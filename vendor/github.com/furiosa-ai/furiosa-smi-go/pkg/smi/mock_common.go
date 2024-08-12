package smi

import "github.com/bradfitz/iter"

type mockHint struct {
	bdf      string
	numaNode int
	serial   string
	uuid     string
	major    uint16
	minor    uint16
}

var staticMockHintMap = map[int]mockHint{
	0: {bdf: "0000:27:00.0", numaNode: 0, serial: "TEST0236FH505KRE0", uuid: "A76AAD68-6855-40B1-9E86-D080852D1C80", major: 234, minor: 0},
	1: {bdf: "0000:2a:00.0", numaNode: 0, serial: "TEST0236FH505KRE1", uuid: "A76AAD68-6855-40B1-9E86-D080852D1C81", major: 235, minor: 0},
	2: {bdf: "0000:51:00.0", numaNode: 0, serial: "TEST0236FH505KRE2", uuid: "A76AAD68-6855-40B1-9E86-D080852D1C82", major: 236, minor: 0},
	3: {bdf: "0000:57:00.0", numaNode: 0, serial: "TEST0236FH505KRE3", uuid: "A76AAD68-6855-40B1-9E86-D080852D1C83", major: 237, minor: 0},
	4: {bdf: "0000:9e:00.0", numaNode: 1, serial: "TEST0236FH505KRE4", uuid: "A76AAD68-6855-40B1-9E86-D080852D1C84", major: 238, minor: 0},
	5: {bdf: "0000:a4:00.0", numaNode: 1, serial: "TEST0236FH505KRE5", uuid: "A76AAD68-6855-40B1-9E86-D080852D1C85", major: 239, minor: 0},
	6: {bdf: "0000:c7:00.0", numaNode: 1, serial: "TEST0236FH505KRE6", uuid: "A76AAD68-6855-40B1-9E86-D080852D1C86", major: 240, minor: 0},
	7: {bdf: "0000:ca:00.0", numaNode: 1, serial: "TEST0236FH505KRE7", uuid: "A76AAD68-6855-40B1-9E86-D080852D1C87", major: 241, minor: 0},
}

// linkTypeHintMap provides hint matrix for optimized 2socket server like the below topology.
// Machine
// ├── Package (CPU)
// │   ├── Host Bridge
// │   │   └····· PCI Bridge
// │   │       ├── NPU0
// │   │       └── NPU1
// │   └── Host Bridge
// │       └····· PCI Bridge
// │           ├── NPU2
// │           └── NPU3
// └── Package (CPU)
//
//	├── Host Bridge
//	│   └····· PCI Bridge
//	│       ├── NPU4
//	│       └── NPU5
//	└── Host Bridge
//	    └····· PCI Bridge
//	        ├── NPU6
//	        └── NPU7
var linkTypeHintMap = map[int]map[int]LinkType{
	0: {0: LinkTypeNoc, 1: LinkTypeHostBridge, 2: LinkTypeCpu, 3: LinkTypeCpu, 4: LinkTypeInterconnect, 5: LinkTypeInterconnect, 6: LinkTypeInterconnect, 7: LinkTypeInterconnect},
	1: {1: LinkTypeNoc, 2: LinkTypeCpu, 3: LinkTypeCpu, 4: LinkTypeInterconnect, 5: LinkTypeInterconnect, 6: LinkTypeInterconnect, 7: LinkTypeInterconnect},
	2: {2: LinkTypeNoc, 3: LinkTypeHostBridge, 4: LinkTypeInterconnect, 5: LinkTypeInterconnect, 6: LinkTypeInterconnect, 7: LinkTypeInterconnect},
	3: {3: LinkTypeNoc, 4: LinkTypeInterconnect, 5: LinkTypeInterconnect, 6: LinkTypeInterconnect, 7: LinkTypeInterconnect},
	4: {4: LinkTypeNoc, 5: LinkTypeHostBridge, 6: LinkTypeCpu, 7: LinkTypeCpu},
	5: {5: LinkTypeNoc, 6: LinkTypeCpu, 7: LinkTypeCpu},
	6: {6: LinkTypeNoc, 7: LinkTypeHostBridge},
	7: {7: LinkTypeNoc},
}

func GetStaticMockDevices(arch Arch) (mockDevices []Device) {
	for i := range iter.N(8) {
		mockDevices = append(mockDevices, GetStaticMockDevice(arch, i))
	}

	return
}

func GetStaticMockDevice(arch Arch, nodeIdx int) Device {
	switch arch {
	case ArchWarboy:
		return &staticWarboyMockDevice{
			arch:    ArchWarboy,
			nodeIdx: nodeIdx,
		}
	case ArchRngd:
		return &staticRngdMockDevice{
			arch:    ArchRngd,
			nodeIdx: nodeIdx,
		}
	}

	panic("unknown arch")
}

type staticMockDeviceFile struct {
	cores []uint32
	path  string
}

var _ DeviceFile = new(staticMockDeviceFile)

func (m *staticMockDeviceFile) Cores() []uint32 {
	return m.cores
}

func (m *staticMockDeviceFile) Path() string {
	return m.path
}

type staticMockDeviceErrorInfo struct{}

var _ DeviceErrorInfo = new(staticMockDeviceErrorInfo)

func (m *staticMockDeviceErrorInfo) AxiPostErrorCount() uint32 {
	return 0
}

func (m *staticMockDeviceErrorInfo) AxiFetchErrorCount() uint32 {
	return 0
}

func (m *staticMockDeviceErrorInfo) AxiDiscardErrorCount() uint32 {
	return 0
}

func (m *staticMockDeviceErrorInfo) AxiDoorbellErrorCount() uint32 {
	return 0
}

func (m *staticMockDeviceErrorInfo) PciePostErrorCount() uint32 {
	return 0
}

func (m *staticMockDeviceErrorInfo) PcieFetchErrorCount() uint32 {
	return 0
}

func (m *staticMockDeviceErrorInfo) PcieDiscardErrorCount() uint32 {
	return 0
}

func (m *staticMockDeviceErrorInfo) PcieDoorbellErrorCount() uint32 {
	return 0
}

func (m *staticMockDeviceErrorInfo) DeviceErrorCount() uint32 {
	return 0
}

type staticMockPeUtilization struct {
	core       uint32
	timeWindow uint32
	usage      float64
}

var _ PeUtilization = new(staticMockPeUtilization)

func (m *staticMockPeUtilization) Core() uint32 {
	return m.core
}

func (m *staticMockPeUtilization) TimeWindowMill() uint32 {
	return m.timeWindow
}

func (m *staticMockPeUtilization) PeUsagePercentage() float64 {
	return m.usage
}

type staticMockMemoryUtilization struct{}

var _ MemoryUtilization = new(staticMockMemoryUtilization)

func (m *staticMockMemoryUtilization) TotalBytes() uint64 {
	return 0
}

func (m *staticMockMemoryUtilization) InUseBytes() uint64 {
	return 0
}

type staticMockDeviceUtilization struct {
	pe  []PeUtilization
	mem MemoryUtilization
}

var _ DeviceUtilization = new(staticMockDeviceUtilization)

func (m *staticMockDeviceUtilization) PeUtilization() []PeUtilization {
	return m.pe
}

func (m *staticMockDeviceUtilization) MemoryUtilization() MemoryUtilization {
	return m.mem
}

type staticMockDeviceTemperature struct{}

var _ DeviceTemperature = new(staticMockDeviceTemperature)

func (m *staticMockDeviceTemperature) SocPeak() float64 {
	return 20
}

func (m *staticMockDeviceTemperature) Ambient() float64 {
	return 10
}

func newStaticMockVersionInfo(arch Arch, major, minor, patch uint32, metadata string) VersionInfo {
	return &staticMockVersionInfo{
		arch:     arch,
		major:    major,
		minor:    minor,
		patch:    patch,
		metadata: metadata,
	}
}

type staticMockVersionInfo struct {
	arch     Arch
	major    uint32
	minor    uint32
	patch    uint32
	metadata string
}

var _ VersionInfo = new(staticMockVersionInfo)

func (m *staticMockVersionInfo) Arch() Arch {
	return m.arch
}

func (m *staticMockVersionInfo) Major() uint32 {
	return m.major
}

func (m *staticMockVersionInfo) Minor() uint32 {
	return m.minor
}

func (m *staticMockVersionInfo) Patch() uint32 {
	return m.patch
}

func (m *staticMockVersionInfo) Metadata() string {
	return m.metadata
}

func getDeviceToDeviceLinkType(src, dst Device) (LinkType, error) {
	selfNodeIdx := src.(*staticRngdMockDevice).nodeIdx
	targetNodeIdx := dst.(*staticRngdMockDevice).nodeIdx

	if selfNodeIdx > targetNodeIdx {
		selfNodeIdx, targetNodeIdx = targetNodeIdx, selfNodeIdx
	}

	ret := linkTypeHintMap[selfNodeIdx][targetNodeIdx]
	return ret, nil
}
