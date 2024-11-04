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

func (m *staticRngdMockDevice) CoreStatus() (map[uint32]CoreStatus, error) {
	return map[uint32]CoreStatus{
		0: CoreStatusAvailable,
		1: CoreStatusAvailable,
		2: CoreStatusAvailable,
		3: CoreStatusAvailable,
		4: CoreStatusAvailable,
		5: CoreStatusAvailable,
		6: CoreStatusAvailable,
		7: CoreStatusAvailable,
	}, nil
}

func (m *staticRngdMockDevice) DeviceErrorInfo() (DeviceErrorInfo, error) {
	return &staticMockDeviceErrorInfo{}, nil
}

func (m *staticRngdMockDevice) Liveness() (bool, error) {
	return true, nil
}

func (m *staticRngdMockDevice) CoreUtilization() (CoreUtilization, error) {
	return &staticMockCoreUtilization{
		pe: []PeUtilization{
			&staticMockPeUtilization{core: 0, timeWindow: 1000, usage: 50},
			&staticMockPeUtilization{core: 1, timeWindow: 1000, usage: 50},
			&staticMockPeUtilization{core: 2, timeWindow: 1000, usage: 50},
			&staticMockPeUtilization{core: 3, timeWindow: 1000, usage: 50},
			&staticMockPeUtilization{core: 4, timeWindow: 1000, usage: 50},
			&staticMockPeUtilization{core: 5, timeWindow: 1000, usage: 50},
			&staticMockPeUtilization{core: 6, timeWindow: 1000, usage: 50},
			&staticMockPeUtilization{core: 7, timeWindow: 1000, usage: 50},
		},
	}, nil
}

func (m *staticRngdMockDevice) MemoryUtilization() (MemoryUtilization, error) {
	return &staticMockMemoryUtilization{}, nil
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
	return fmt.Sprintf("/dev/rngd/npu%d", m.nodeIdx)
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
