package smi

import (
	"fmt"
)

var _ Device = new(staticRngdPlusMockDevice)

type staticRngdPlusMockDevice struct {
	arch    Arch
	nodeIdx int
}

func (m *staticRngdPlusMockDevice) DeviceInfo() (DeviceInfo, error) {
	return &staticRngdPlusMockDeviceInfo{
		nodeIdx: m.nodeIdx,
	}, nil
}

func (m *staticRngdPlusMockDevice) DeviceFiles() ([]DeviceFile, error) {
	return []DeviceFile{
		&staticMockDeviceFile{
			cores: []uint32{0},
			path:  fmt.Sprintf("/dev/rngd/npu%dpe0", m.nodeIdx),
		},
		&staticMockDeviceFile{
			cores: []uint32{1},
			path:  fmt.Sprintf("/dev/rngd/npu%dpe1", m.nodeIdx),
		},
		&staticMockDeviceFile{
			cores: []uint32{0, 1},
			path:  fmt.Sprintf("/dev/rngd/npu%dpe0-1", m.nodeIdx),
		},
		&staticMockDeviceFile{
			cores: []uint32{2},
			path:  fmt.Sprintf("/dev/rngd/npu%dpe2", m.nodeIdx),
		},
		&staticMockDeviceFile{
			cores: []uint32{3},
			path:  fmt.Sprintf("/dev/rngd/npu%dpe3", m.nodeIdx),
		},
		&staticMockDeviceFile{
			cores: []uint32{2, 3},
			path:  fmt.Sprintf("/dev/rngd/npu%dpe2-3", m.nodeIdx),
		},
		&staticMockDeviceFile{
			cores: []uint32{0, 1, 2, 3},
			path:  fmt.Sprintf("/dev/rngd/npu%dpe0-3", m.nodeIdx),
		},
		&staticMockDeviceFile{
			cores: []uint32{4},
			path:  fmt.Sprintf("/dev/rngd/npu%dpe4", m.nodeIdx),
		},
		&staticMockDeviceFile{
			cores: []uint32{5},
			path:  fmt.Sprintf("/dev/rngd/npu%dpe5", m.nodeIdx),
		},
		&staticMockDeviceFile{
			cores: []uint32{4, 5},
			path:  fmt.Sprintf("/dev/rngd/npu%dpe4-5", m.nodeIdx),
		},
		&staticMockDeviceFile{
			cores: []uint32{6},
			path:  fmt.Sprintf("/dev/rngd/npu%dpe6", m.nodeIdx),
		},
		&staticMockDeviceFile{
			cores: []uint32{7},
			path:  fmt.Sprintf("/dev/rngd/npu%dpe7", m.nodeIdx),
		},
		&staticMockDeviceFile{
			cores: []uint32{6, 7},
			path:  fmt.Sprintf("/dev/rngd/npu%dpe6-7", m.nodeIdx),
		},
		&staticMockDeviceFile{
			cores: []uint32{4, 5, 6, 7},
			path:  fmt.Sprintf("/dev/rngd/npu%dpe4-7", m.nodeIdx),
		},
	}, nil
}

func (m *staticRngdPlusMockDevice) CoreStatus() (CoreStatuses, error) {
	return staticMockCoreStatuses{
		coreStatus: []PeStatus{
			&staticMockPeStatus{core: 0, status: CoreStatusAvailable},
			&staticMockPeStatus{core: 1, status: CoreStatusAvailable},
			&staticMockPeStatus{core: 2, status: CoreStatusAvailable},
			&staticMockPeStatus{core: 3, status: CoreStatusAvailable},
			&staticMockPeStatus{core: 4, status: CoreStatusAvailable},
			&staticMockPeStatus{core: 5, status: CoreStatusAvailable},
			&staticMockPeStatus{core: 6, status: CoreStatusAvailable},
			&staticMockPeStatus{core: 7, status: CoreStatusAvailable},
		},
	}, nil
}

func (m *staticRngdPlusMockDevice) Liveness() (bool, error) {
	return true, nil
}

func (m *staticRngdPlusMockDevice) CoreFrequency() (CoreFrequency, error) {
	return &staticMockCoreFrequency{
		pe: []PeFrequency{
			&staticMockPeFrequency{core: 0, frequency: 500},
			&staticMockPeFrequency{core: 1, frequency: 500},
			&staticMockPeFrequency{core: 2, frequency: 500},
			&staticMockPeFrequency{core: 3, frequency: 500},
			&staticMockPeFrequency{core: 4, frequency: 500},
			&staticMockPeFrequency{core: 5, frequency: 500},
			&staticMockPeFrequency{core: 6, frequency: 500},
			&staticMockPeFrequency{core: 7, frequency: 500},
		},
	}, nil
}

func (m *staticRngdPlusMockDevice) MemoryFrequency() (MemoryFrequency, error) {
	return &staticMockMemoryFrequency{frequency: 6000}, nil
}

func (m *staticRngdPlusMockDevice) PowerConsumption() (float64, error) {
	return float64(100), nil
}

func (m *staticRngdPlusMockDevice) DeviceTemperature() (DeviceTemperature, error) {
	return &staticMockDeviceTemperature{}, nil
}

func (m *staticRngdPlusMockDevice) DeviceToDeviceLinkType(target Device) (LinkType, error) {
	return getDeviceToDeviceLinkTypeRngdPlus(m, target)
}

func (m *staticRngdPlusMockDevice) P2PAccessible(_ Device) (bool, error) {
	return true, nil
}

func (m *staticRngdPlusMockDevice) DevicePerformanceCounter() (DevicePerformanceCounter, error) {
	return &staticMockDevicePerformanceCounter{}, nil
}

func (m *staticRngdPlusMockDevice) GovernorProfile() (GovernorProfile, error) {
	return GovernorProfilePerformance, nil
}

func (m *staticRngdPlusMockDevice) SetGovernorProfile(profile GovernorProfile) error {
	return nil
}

func (m *staticRngdPlusMockDevice) EnableDevice() error {
	return nil
}

func (m *staticRngdPlusMockDevice) DisableDevice() error {
	return nil
}

type staticRngdPlusMockDeviceInfo struct {
	nodeIdx int
}

var _ DeviceInfo = new(staticRngdPlusMockDeviceInfo)

func (m *staticRngdPlusMockDeviceInfo) Index() uint32 {
	return uint32(m.nodeIdx)
}

func (m *staticRngdPlusMockDeviceInfo) Arch() Arch {
	return ArchRngdPlus
}

func (m *staticRngdPlusMockDeviceInfo) CoreNum() uint32 {
	return 8
}

func (m *staticRngdPlusMockDeviceInfo) NumaNode() int32 {
	return int32(staticMockHintMap[m.nodeIdx].numaNode)
}

func (m *staticRngdPlusMockDeviceInfo) Name() string {
	return fmt.Sprintf("npu%d", m.nodeIdx)
}

func (m *staticRngdPlusMockDeviceInfo) Serial() string {
	return staticMockHintMap[m.nodeIdx].serial
}

func (m *staticRngdPlusMockDeviceInfo) UUID() string {
	return staticMockHintMap[m.nodeIdx].uuid
}

func (m *staticRngdPlusMockDeviceInfo) BDF() string {
	return staticMockHintMap[m.nodeIdx].bdf
}

func (m *staticRngdPlusMockDeviceInfo) Major() uint16 {
	return staticMockHintMap[m.nodeIdx].major
}

func (m *staticRngdPlusMockDeviceInfo) Minor() uint16 {
	return staticMockHintMap[m.nodeIdx].minor
}

func (m *staticRngdPlusMockDeviceInfo) FirmwareVersion() VersionInfo {
	return newStaticMockVersionInfo(1, 6, 0, "c1bebfd", "dev0")
}

func (m *staticRngdPlusMockDevice) PcieInfo() (PcieInfo, error) {
	return &staticMockPcieInfo{}, nil
}

func (staticRngdPlusMockDevice) CoreUtilization(observer Observer) ([]CoreUtilization, error) {
	return nil, nil
}

func (staticRngdPlusMockDevice) ThrottleReason() (ThrottleReason, error) {
	return 0, nil
}

func (staticRngdPlusMockDevice) MemoryUtilization() (MemoryUtilization, error) {
	return nil, nil
}

func getDeviceToDeviceLinkTypeRngdPlus(src, dst Device) (LinkType, error) {
	selfNodeIdx := src.(*staticRngdPlusMockDevice).nodeIdx
	targetNodeIdx := dst.(*staticRngdPlusMockDevice).nodeIdx

	if selfNodeIdx > targetNodeIdx {
		selfNodeIdx, targetNodeIdx = targetNodeIdx, selfNodeIdx
	}

	ret := linkTypeHintMap[selfNodeIdx][targetNodeIdx]
	return ret, nil
}
