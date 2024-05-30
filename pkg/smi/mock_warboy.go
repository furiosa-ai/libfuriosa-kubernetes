package smi

import (
	"fmt"
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
	0: {bdf: "0000:27:00.0", numaNode: 0, serial: "WBYB0236FH505KRE0", uuid: "A76AAD68-6855-40B1-9E86-D080852D1C80", major: 234, minor: 0},
	1: {bdf: "0000:2a:00.0", numaNode: 0, serial: "WBYB0236FH505KRE1", uuid: "A76AAD68-6855-40B1-9E86-D080852D1C81", major: 235, minor: 0},
	2: {bdf: "0000:51:00.0", numaNode: 0, serial: "WBYB0236FH505KRE2", uuid: "A76AAD68-6855-40B1-9E86-D080852D1C82", major: 236, minor: 0},
	3: {bdf: "0000:57:00.0", numaNode: 0, serial: "WBYB0236FH505KRE3", uuid: "A76AAD68-6855-40B1-9E86-D080852D1C83", major: 237, minor: 0},
	4: {bdf: "0000:9e:00.0", numaNode: 1, serial: "WBYB0236FH505KRE4", uuid: "A76AAD68-6855-40B1-9E86-D080852D1C84", major: 238, minor: 0},
	5: {bdf: "0000:a4:00.0", numaNode: 1, serial: "WBYB0236FH505KRE5", uuid: "A76AAD68-6855-40B1-9E86-D080852D1C85", major: 239, minor: 0},
	6: {bdf: "0000:c7:00.0", numaNode: 1, serial: "WBYB0236FH505KRE6", uuid: "A76AAD68-6855-40B1-9E86-D080852D1C86", major: 240, minor: 0},
	7: {bdf: "0000:ca:00.0", numaNode: 1, serial: "WBYB0236FH505KRE7", uuid: "A76AAD68-6855-40B1-9E86-D080852D1C87", major: 241, minor: 0},
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

func GetStaticMockWarboyDevices() (mockDevices []Device) {
	for i := range iter.N(8) {
		mockDevices = append(mockDevices, GetStaticMockWarboyDevice(i))
	}

	return
}

func GetStaticMockWarboyDevice(nodeIdx int) Device {
	return &staticMockDevice{
		nodeIdx: nodeIdx,
	}
}

var _ Device = new(staticMockDevice)

type staticMockDevice struct {
	nodeIdx int
}

func (m staticMockDevice) DeviceInfo() (DeviceInfo, error) {
	return &staticMockDeviceInfo{
		nodeIdx: m.nodeIdx,
	}, nil
}

func (m staticMockDevice) DeviceFiles() ([]DeviceFile, error) {
	return []DeviceFile{
		&staticMockDeviceFile{
			cores: []uint32{0},
			path:  fmt.Sprintf("/dev/npu%d", m.nodeIdx),
		},
		&staticMockDeviceFile{
			cores: []uint32{0},
			path:  fmt.Sprintf("/dev/npu%dpe0", m.nodeIdx),
		},
		&staticMockDeviceFile{
			cores: []uint32{1},
			path:  fmt.Sprintf("/dev/npu%dpe1", m.nodeIdx),
		},
		&staticMockDeviceFile{
			cores: []uint32{1},
			path:  fmt.Sprintf("/dev/npu%dpe0-1", m.nodeIdx),
		},
	}, nil
}

func (m staticMockDevice) CoreStatus() (map[uint32]CoreStatus, error) {
	return map[uint32]CoreStatus{0: CoreStatusAvailable, 1: CoreStatusAvailable}, nil
}

func (m staticMockDevice) DeviceErrorInfo() (DeviceErrorInfo, error) {
	return &staticMockDeviceErrorInfo{}, nil
}

func (m staticMockDevice) Liveness() (bool, error) {
	return true, nil
}

func (m staticMockDevice) DeviceUtilization() (DeviceUtilization, error) {
	return &staticMockDeviceUtilization{
		pe: []PeUtilization{
			&staticMockPeUtilization{cores: []uint32{0}, timeWindow: 1000, usage: 50},
		},
		mem: &staticMockMemoryUtilization{},
	}, nil
}

func (m staticMockDevice) PowerConsumption() (uint32, error) {
	return 100, nil
}

func (m staticMockDevice) DeviceTemperature() (DeviceTemperature, error) {
	return &staticMockDeviceTemperature{}, nil
}

func (m staticMockDevice) GetDeviceToDeviceLinkType(target Device) (LinkType, error) {
	selfNodeIdx := m.nodeIdx
	targetNodeIdx := target.(*staticMockDevice).nodeIdx

	if selfNodeIdx > targetNodeIdx {
		selfNodeIdx, targetNodeIdx = targetNodeIdx, selfNodeIdx
	}

	ret := linkTypeHintMap[selfNodeIdx][targetNodeIdx]
	return ret, nil
}

type staticMockDeviceInfo struct {
	nodeIdx int
}

var _ DeviceInfo = new(staticMockDeviceInfo)

func (m staticMockDeviceInfo) Arch() Arch {
	return ArchWarboy
}

func (m staticMockDeviceInfo) CoreNum() uint32 {
	return 2
}

func (m staticMockDeviceInfo) NumaNode() uint32 {
	return uint32(staticMockHintMap[m.nodeIdx].numaNode)
}

func (m staticMockDeviceInfo) Name() string {
	return fmt.Sprintf("/dev/npu%d", m.nodeIdx)
}

func (m staticMockDeviceInfo) Serial() string {
	return staticMockHintMap[m.nodeIdx].serial
}

func (m staticMockDeviceInfo) UUID() string {
	return staticMockHintMap[m.nodeIdx].uuid
}

func (m staticMockDeviceInfo) BDF() string {
	return staticMockHintMap[m.nodeIdx].bdf
}

func (m staticMockDeviceInfo) Major() uint16 {
	return staticMockHintMap[m.nodeIdx].major
}

func (m staticMockDeviceInfo) Minor() uint16 {
	return staticMockHintMap[m.nodeIdx].minor
}

func (m staticMockDeviceInfo) FirmwareVersion() VersionInfo {
	return &staticMockFirmwareVersionInfo{}
}

func (m staticMockDeviceInfo) DriverVersion() VersionInfo {
	return &staticMockDriverVersionInfo{}
}

type staticMockDeviceFile struct {
	cores []uint32
	path  string
}

var _ DeviceFile = new(staticMockDeviceFile)

func (m staticMockDeviceFile) Cores() []uint32 {
	return m.cores
}

func (m staticMockDeviceFile) Path() string {
	return m.path
}

type staticMockDeviceErrorInfo struct{}

var _ DeviceErrorInfo = new(staticMockDeviceErrorInfo)

func (m staticMockDeviceErrorInfo) AxiPostErrorCount() uint32 {
	return 0
}

func (m staticMockDeviceErrorInfo) AxiFetchErrorCount() uint32 {
	return 0
}

func (m staticMockDeviceErrorInfo) AxiDiscardErrorCount() uint32 {
	return 0
}

func (m staticMockDeviceErrorInfo) AxiDoorbellErrorCount() uint32 {
	return 0
}

func (m staticMockDeviceErrorInfo) PciePostErrorCount() uint32 {
	return 0
}

func (m staticMockDeviceErrorInfo) PcieFetchErrorCount() uint32 {
	return 0
}

func (m staticMockDeviceErrorInfo) PcieDiscardErrorCount() uint32 {
	return 0
}

func (m staticMockDeviceErrorInfo) PcieDoorbellErrorCount() uint32 {
	return 0
}

func (m staticMockDeviceErrorInfo) DeviceErrorCount() uint32 {
	return 0
}

type staticMockPeUtilization struct {
	cores      []uint32
	timeWindow uint32
	usage      uint32
}

var _ PeUtilization = new(staticMockPeUtilization)

func (m staticMockPeUtilization) Cores() []uint32 {
	return m.cores
}

func (m staticMockPeUtilization) TimeWindowMill() uint32 {
	return m.timeWindow
}

func (m staticMockPeUtilization) PeUsagePercentage() uint32 {
	return m.usage
}

type staticMockMemoryUtilization struct{}

var _ MemoryUtilization = new(staticMockMemoryUtilization)

func (m staticMockMemoryUtilization) TotalBytes() uint64 {
	return 0
}

func (m staticMockMemoryUtilization) InUseBytes() uint64 {
	return 0
}

type staticMockDeviceUtilization struct {
	pe  []PeUtilization
	mem MemoryUtilization
}

var _ DeviceUtilization = new(staticMockDeviceUtilization)

func (m staticMockDeviceUtilization) PeUtilization() []PeUtilization {
	return m.pe
}

func (m staticMockDeviceUtilization) MemoryUtilization() MemoryUtilization {
	return m.mem
}

type staticMockDeviceTemperature struct{}

var _ DeviceTemperature = new(staticMockDeviceTemperature)

func (m staticMockDeviceTemperature) SocPeak() int32 {
	return 0
}

func (m staticMockDeviceTemperature) Ambient() int32 {
	return 0
}

// version: 1.9.2, 3def9c2
type staticMockDriverVersionInfo struct{}

var _ VersionInfo = new(staticMockDriverVersionInfo)

func (m staticMockDriverVersionInfo) Arch() Arch {
	return ArchWarboy
}

func (m staticMockDriverVersionInfo) Major() uint32 {
	return 1
}

func (m staticMockDriverVersionInfo) Minor() uint32 {
	return 9
}

func (m staticMockDriverVersionInfo) Patch() uint32 {
	return 2
}

func (m staticMockDriverVersionInfo) Metadata() string {
	return "3def9c2"
}

// version: 1.6.0, c1bebfd
type staticMockFirmwareVersionInfo struct {
}

var _ VersionInfo = new(staticMockFirmwareVersionInfo)

func (m staticMockFirmwareVersionInfo) Arch() Arch {
	return ArchWarboy
}

func (m staticMockFirmwareVersionInfo) Major() uint32 {
	return 1
}

func (m staticMockFirmwareVersionInfo) Minor() uint32 {
	return 6
}

func (m staticMockFirmwareVersionInfo) Patch() uint32 {
	return 0
}

func (m staticMockFirmwareVersionInfo) Metadata() string {
	return "c1bebfd"
}
