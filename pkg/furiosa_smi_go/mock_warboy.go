package furiosa_smi_go

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

var mockHintMap = map[int]mockHint{
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

func GetMockWarboyDevices() (mockDevices []Device) {
	for i := range iter.N(8) {
		mockDevices = append(mockDevices, GetMockWarboyDevice(i))
	}

	return
}
func GetMockWarboyDevice(nodeIdx int) Device {
	return &mockDevice{
		nodeIdx: nodeIdx,
	}
}

var _ Device = new(mockDevice)

type mockDevice struct {
	nodeIdx int
}

func (m mockDevice) DeviceInfo() (DeviceInfo, error) {
	return &mockDeviceInfo{
		nodeIdx: m.nodeIdx,
	}, nil
}

func (m mockDevice) DeviceFiles() ([]DeviceFile, error) {
	return []DeviceFile{
		&mockDeviceFile{
			cores: []uint32{0},
			path:  fmt.Sprintf("/dev/npu%d", m.nodeIdx),
		},
		&mockDeviceFile{
			cores: []uint32{0},
			path:  fmt.Sprintf("/dev/npu%dpe0", m.nodeIdx),
		},
		&mockDeviceFile{
			cores: []uint32{1},
			path:  fmt.Sprintf("/dev/npu%dpe1", m.nodeIdx),
		},
		&mockDeviceFile{
			cores: []uint32{1},
			path:  fmt.Sprintf("/dev/npu%dpe0-1", m.nodeIdx),
		},
	}, nil
}

func (m mockDevice) CoreStatus() (map[uint32]CoreStatus, error) {
	return map[uint32]CoreStatus{0: CoreStatusAvailable, 1: CoreStatusAvailable}, nil
}

func (m mockDevice) DeviceErrorInfo() (DeviceErrorInfo, error) {
	return &mockDeviceErrorInfo{}, nil
}

func (m mockDevice) Liveness() (bool, error) {
	return true, nil
}

func (m mockDevice) DeviceUtilization() (DeviceUtilization, error) {
	return &mockDeviceUtilization{
		pe: []PeUtilization{
			&mockPeUtilization{cores: []uint32{0}, timeWindow: 1000, usage: 50},
		},
		mem: &mockMemoryUtilization{},
	}, nil
}

func (m mockDevice) PowerConsumption() (uint32, error) {
	return 100, nil
}

func (m mockDevice) DeviceTemperature() (DeviceTemperature, error) {
	return &mockDeviceTemperature{}, nil
}

func (m mockDevice) GetDeviceToDeviceLinkType(target Device) (LinkType, error) {
	selfNodeIdx := m.nodeIdx
	targetNodeIdx := target.(*mockDevice).nodeIdx

	if selfNodeIdx > targetNodeIdx {
		selfNodeIdx, targetNodeIdx = targetNodeIdx, selfNodeIdx
	}

	ret := linkTypeHintMap[selfNodeIdx][targetNodeIdx]
	return ret, nil
}

type mockDeviceInfo struct {
	nodeIdx int
}

var _ DeviceInfo = new(mockDeviceInfo)

func (m mockDeviceInfo) Arch() Arch {
	return ArchWarboy
}

func (m mockDeviceInfo) CoreNum() uint32 {
	return 2
}

func (m mockDeviceInfo) NumaNode() uint32 {
	return 0
}

func (m mockDeviceInfo) Name() string {
	return fmt.Sprintf("/dev/npu%d", m.nodeIdx)
}

func (m mockDeviceInfo) Serial() string {
	return mockHintMap[m.nodeIdx].serial
}

func (m mockDeviceInfo) UUID() string {
	return mockHintMap[m.nodeIdx].uuid
}

func (m mockDeviceInfo) BDF() string {
	return mockHintMap[m.nodeIdx].bdf
}

func (m mockDeviceInfo) Major() uint16 {
	return mockHintMap[m.nodeIdx].major
}

func (m mockDeviceInfo) Minor() uint16 {
	return mockHintMap[m.nodeIdx].minor
}

func (m mockDeviceInfo) FirmwareVersion() VersionInfo {
	return &mockFirmwareVersionInfo{}
}

func (m mockDeviceInfo) DriverVersion() VersionInfo {
	return &mockDriverVersionInfo{}
}

type mockDeviceFile struct {
	cores []uint32
	path  string
}

var _ DeviceFile = new(mockDeviceFile)

func (m mockDeviceFile) Cores() []uint32 {
	return m.cores
}

func (m mockDeviceFile) Path() string {
	return m.path
}

type mockDeviceErrorInfo struct{}

var _ DeviceErrorInfo = new(mockDeviceErrorInfo)

func (m mockDeviceErrorInfo) AxiPostErrorCount() uint32 {
	return 0
}

func (m mockDeviceErrorInfo) AxiFetchErrorCount() uint32 {
	return 0
}

func (m mockDeviceErrorInfo) AxiDiscardErrorCount() uint32 {
	return 0
}

func (m mockDeviceErrorInfo) AxiDoorbellErrorCount() uint32 {
	return 0
}

func (m mockDeviceErrorInfo) PciePostErrorCount() uint32 {
	return 0
}

func (m mockDeviceErrorInfo) PcieFetchErrorCount() uint32 {
	return 0
}

func (m mockDeviceErrorInfo) PcieDiscardErrorCount() uint32 {
	return 0
}

func (m mockDeviceErrorInfo) PcieDoorbellErrorCount() uint32 {
	return 0
}

func (m mockDeviceErrorInfo) DeviceErrorCount() uint32 {
	return 0
}

type mockPeUtilization struct {
	cores      []uint32
	timeWindow uint32
	usage      uint32
}

var _ PeUtilization = new(mockPeUtilization)

func (m mockPeUtilization) Cores() []uint32 {
	return m.cores
}

func (m mockPeUtilization) TimeWindowMill() uint32 {
	return m.timeWindow
}

func (m mockPeUtilization) PeUsagePercentage() uint32 {
	return m.usage
}

type mockMemoryUtilization struct{}

var _ MemoryUtilization = new(mockMemoryUtilization)

func (m mockMemoryUtilization) TotalBytes() uint64 {
	return 0
}

func (m mockMemoryUtilization) InUseBytes() uint64 {
	return 0
}

type mockDeviceUtilization struct {
	pe  []PeUtilization
	mem MemoryUtilization
}

var _ DeviceUtilization = new(mockDeviceUtilization)

func (m mockDeviceUtilization) PeUtilization() []PeUtilization {
	return m.pe
}

func (m mockDeviceUtilization) MemoryUtilization() MemoryUtilization {
	return m.mem
}

type mockDeviceTemperature struct{}

var _ DeviceTemperature = new(mockDeviceTemperature)

func (m mockDeviceTemperature) SocPeak() int32 {
	return 0
}

func (m mockDeviceTemperature) Ambient() int32 {
	return 0
}

// version: 1.9.2, 3def9c2
type mockDriverVersionInfo struct{}

var _ VersionInfo = new(mockDriverVersionInfo)

func (m mockDriverVersionInfo) Arch() Arch {
	return ArchWarboy
}

func (m mockDriverVersionInfo) Major() uint32 {
	return 1
}

func (m mockDriverVersionInfo) Minor() uint32 {
	return 9
}

func (m mockDriverVersionInfo) Patch() uint32 {
	return 2
}

func (m mockDriverVersionInfo) Metadata() string {
	return "3def9c2"
}

// version: 1.6.0, c1bebfd
type mockFirmwareVersionInfo struct {
}

var _ VersionInfo = new(mockFirmwareVersionInfo)

func (m mockFirmwareVersionInfo) Arch() Arch {
	return ArchWarboy
}

func (m mockFirmwareVersionInfo) Major() uint32 {
	return 1
}

func (m mockFirmwareVersionInfo) Minor() uint32 {
	return 6
}

func (m mockFirmwareVersionInfo) Patch() uint32 {
	return 0
}

func (m mockFirmwareVersionInfo) Metadata() string {
	return "c1bebfd"
}
