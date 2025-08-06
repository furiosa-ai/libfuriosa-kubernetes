package smi

import "fmt"

var _ Device = new(staticRngdMockDevice)

type staticRngdMockDevice struct {
	arch    Arch
	nodeIdx int
}

func (m *staticRngdMockDevice) DeviceInfo() (DeviceInfo, error) {
	return &staticRngdMockDeviceInfo{
		nodeIdx: m.nodeIdx,
	}, nil
}

func (m *staticRngdMockDevice) DeviceFiles() ([]DeviceFile, error) {
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

func (m *staticRngdMockDevice) CoreStatus() (CoreStatuses, error) {
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

func (m *staticRngdMockDevice) Liveness() (bool, error) {
	return true, nil
}

func (m *staticRngdMockDevice) CoreFrequency() (CoreFrequency, error) {
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

func (m *staticRngdMockDevice) MemoryFrequency() (MemoryFrequency, error) {
	return &staticMockMemoryFrequency{frequency: 6000}, nil
}

func (m *staticRngdMockDevice) PowerConsumption() (float64, error) {
	return float64(100), nil
}

func (m *staticRngdMockDevice) DeviceTemperature() (DeviceTemperature, error) {
	return &staticMockDeviceTemperature{}, nil
}

func (m *staticRngdMockDevice) DeviceToDeviceLinkType(target Device) (LinkType, error) {
	return getDeviceToDeviceLinkType(m, target)
}

func (m *staticRngdMockDevice) P2PAccessible(_ Device) (bool, error) {
	return true, nil

}

func (m *staticRngdMockDevice) DevicePerformanceCounter() (DevicePerformanceCounter, error) {
	return &staticMockDevicePerformanceCounter{}, nil
}

func (m *staticRngdMockDevice) GovernorProfile() (GovernorProfile, error) {
	return GovernorProfilePerformance, nil
}

func (m *staticRngdMockDevice) SetGovernorProfile(profile GovernorProfile) error {
	return nil
}

func (m *staticRngdMockDevice) EnableDevice() error {
	return nil
}

func (m *staticRngdMockDevice) DisableDevice() error {
	return nil
}

type staticRngdMockDeviceInfo struct {
	nodeIdx int
}

var _ DeviceInfo = new(staticRngdMockDeviceInfo)

func (m *staticRngdMockDeviceInfo) Index() uint32 {
	return uint32(m.nodeIdx)
}

func (m *staticRngdMockDeviceInfo) Arch() Arch {
	return ArchRngd
}

func (m *staticRngdMockDeviceInfo) CoreNum() uint32 {
	return 8
}

func (m *staticRngdMockDeviceInfo) NumaNode() uint32 {
	return uint32(staticMockHintMap[m.nodeIdx].numaNode)
}

func (m *staticRngdMockDeviceInfo) Name() string {
	return fmt.Sprintf("npu%d", m.nodeIdx)
}

func (m *staticRngdMockDeviceInfo) Serial() string {
	return staticMockHintMap[m.nodeIdx].serial
}

func (m *staticRngdMockDeviceInfo) UUID() string {
	return staticMockHintMap[m.nodeIdx].uuid
}

func (m *staticRngdMockDeviceInfo) BDF() string {
	return staticMockHintMap[m.nodeIdx].bdf
}

func (m *staticRngdMockDeviceInfo) Major() uint16 {
	return staticMockHintMap[m.nodeIdx].major
}

func (m *staticRngdMockDeviceInfo) Minor() uint16 {
	return staticMockHintMap[m.nodeIdx].minor
}

// FirmwareVersion e.g. version: 1.6.0, c1bebfd
func (m *staticRngdMockDeviceInfo) FirmwareVersion() VersionInfo {
	return newStaticMockVersionInfo(1, 6, 0, "c1bebfd")
}

func (m *staticRngdMockDeviceInfo) PertVersion() VersionInfo {
	return newStaticMockVersionInfo(0, 0, 0, "")
}

func (m *staticRngdMockDevice) PcieInfo() (PcieInfo, error) {
	return &staticMockPcieInfo{}, nil
}
