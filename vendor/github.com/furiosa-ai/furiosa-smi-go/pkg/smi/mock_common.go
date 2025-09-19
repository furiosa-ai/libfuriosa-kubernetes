package smi

import (
	"fmt"
	"time"

	"github.com/bradfitz/iter"
)

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

// GetStaticMockDevices returns mock devices with given Arch for test purpose.
func GetStaticMockDevices(arch Arch) (mockDevices []Device) {
	for i := range iter.N(8) {
		mockDevices = append(mockDevices, GetStaticMockDevice(arch, i))
	}

	return
}

// GetStaticMockDevice returns single mock device with given Arch and idx number for test purpose.
func GetStaticMockDevice(arch Arch, nodeIdx int) Device {
	switch arch {
	case ArchRngd:
		return &staticRngdMockDevice{
			arch:    ArchRngd,
			nodeIdx: nodeIdx,
		}
		// TODO: add more arch
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

type staticMockCoreFrequency struct {
	pe []PeFrequency
}

var _ CoreFrequency = new(staticMockCoreFrequency)

func (m *staticMockCoreFrequency) PeFrequency() []PeFrequency {
	return m.pe
}

type staticMockPeFrequency struct {
	core      uint32
	frequency uint32
}

var _ PeFrequency = new(staticMockPeFrequency)

func (m *staticMockPeFrequency) Core() uint32 {
	return m.core
}

func (m *staticMockPeFrequency) Frequency() uint32 {
	return m.frequency
}

type staticMockCoreStatuses struct {
	coreStatus []PeStatus
}

var _ CoreStatuses = new(staticMockCoreStatuses)

func (m staticMockCoreStatuses) PeStatus() []PeStatus {
	return m.coreStatus
}

type staticMockPeStatus struct {
	core   uint32
	status CoreStatus
}

var _ PeStatus = new(staticMockPeStatus)

func (m *staticMockPeStatus) Core() uint32 {
	return m.core
}

func (m *staticMockPeStatus) Status() CoreStatus {
	return m.status
}

type staticMockMemoryFrequency struct {
	frequency uint32
}

var _ MemoryFrequency = new(staticMockMemoryFrequency)

func (m *staticMockMemoryFrequency) Frequency() uint32 {
	return m.frequency
}

type staticMockDeviceTemperature struct{}

var _ DeviceTemperature = new(staticMockDeviceTemperature)

func (m *staticMockDeviceTemperature) SocPeak() float64 {
	return 20
}

func (m *staticMockDeviceTemperature) Ambient() float64 {
	return 10
}

func newStaticMockVersionInfo(major, minor, patch uint32, metadata string) VersionInfo {
	return &staticMockVersionInfo{
		major:    major,
		minor:    minor,
		patch:    patch,
		metadata: metadata,
	}
}

type staticMockVersionInfo struct {
	major    uint32
	minor    uint32
	patch    uint32
	metadata string
}

var _ VersionInfo = new(staticMockVersionInfo)

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

func (m *staticMockVersionInfo) String() string {
	return fmt.Sprintf("%d.%d.%d+%s", m.Major(), m.Minor(), m.Patch(), m.Metadata())
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

type staticMockDevicePerformanceCounter struct{}

var _ DevicePerformanceCounter = new(staticMockDevicePerformanceCounter)

func (s staticMockDevicePerformanceCounter) PerformanceCounter() []PerformanceCounter {
	return []PerformanceCounter{
		&staticMockPerformanceCounter{},
	}
}

var _ DevicePerformanceCounter = new(staticMockDevicePerformanceCounter)

type staticMockPerformanceCounter struct{}

var _ PerformanceCounter = new(staticMockPerformanceCounter)

func (s staticMockPerformanceCounter) Timestamp() time.Time {
	return time.Now()
}

func (s staticMockPerformanceCounter) Core() uint32 {
	return 0
}

func (s staticMockPerformanceCounter) CycleCount() uint64 {
	return 0
}

func (s staticMockPerformanceCounter) TaskExecutionCycle() uint64 {
	return 0
}

type staticMockPcieInfo struct{}

var _ PcieInfo = new(staticMockPcieInfo)

func (s *staticMockPcieInfo) DeviceInfo() PcieDeviceInfo {
	return &staticMockPcieDeviceInfo{}
}

func (s *staticMockPcieInfo) LinkInfo() PcieLinkInfo {
	return &staticMockPcieLinkInfo{}
}

func (s *staticMockPcieInfo) SriovInfo() SriovInfo {
	return &staticMockSriovInfo{}
}

func (s *staticMockPcieInfo) RootComplexInfo() PcieRootComplexInfo {
	return &staticMockPcieRootComplexInfo{}
}

func (s *staticMockPcieInfo) SwitchInfo() PcieSwitchInfo {
	return &staticMockPcieSwitchInfo{}
}

type staticMockPcieDeviceInfo struct{}

var _ PcieDeviceInfo = new(staticMockPcieDeviceInfo)

func (s *staticMockPcieDeviceInfo) DeviceId() uint16 {
	return 0x0001
}

func (s *staticMockPcieDeviceInfo) VendorId() uint16 {
	return 0x1ed2
}

func (s *staticMockPcieDeviceInfo) SubsystemId() uint16 {
	return 0x0001
}

func (s *staticMockPcieDeviceInfo) RevisionId() uint8 {
	return 0x0001
}

func (s *staticMockPcieDeviceInfo) ClassId() uint8 {
	return 0x12
}

func (s *staticMockPcieDeviceInfo) SubClassId() uint8 {
	return 0x00
}

type staticMockPcieLinkInfo struct{}

var _ PcieLinkInfo = new(staticMockPcieLinkInfo)

func (s *staticMockPcieLinkInfo) PcieGenStatus() uint8 {
	return 5
}

func (s *staticMockPcieLinkInfo) LinkWidthStatus() uint32 {
	return 16
}

func (s *staticMockPcieLinkInfo) LinkSpeedStatus() float64 {
	return 32.0
}

func (s *staticMockPcieLinkInfo) MaxLinkWidthCapability() uint32 {
	return 16
}

func (s *staticMockPcieLinkInfo) MaxLinkSpeedCapability() float64 {
	return 32.0
}

type staticMockSriovInfo struct{}

var _ SriovInfo = new(staticMockSriovInfo)

func (s *staticMockSriovInfo) SriovTotalVfs() uint32 {
	return 4
}

func (s *staticMockSriovInfo) SriovEnabledVfs() uint32 {
	return 0
}

type staticMockPcieRootComplexInfo struct{}

var _ PcieRootComplexInfo = new(staticMockPcieRootComplexInfo)

func (s *staticMockPcieRootComplexInfo) Domain() uint16 {
	return 0x0000
}

func (s *staticMockPcieRootComplexInfo) Bus() uint8 {
	return 0x42
}

func (s *staticMockPcieRootComplexInfo) String() string {
	return fmt.Sprintf("%04x:%02x", s.Domain(), s.Bus())
}

type staticMockPcieSwitchInfo struct{}

var _ PcieSwitchInfo = new(staticMockPcieSwitchInfo)

func (s *staticMockPcieSwitchInfo) Domain() uint16 {
	return 0x0000
}

func (s *staticMockPcieSwitchInfo) Bus() uint8 {
	return 0x41
}

func (s *staticMockPcieSwitchInfo) Device() uint8 {
	return 0x00
}

func (s *staticMockPcieSwitchInfo) Function() uint8 {
	return 0x00
}

func (s *staticMockPcieSwitchInfo) String() string {
	return fmt.Sprintf("%04x:%02x:%02x.%d", s.Domain(), s.Bus(), s.Device(), s.Function())
}
