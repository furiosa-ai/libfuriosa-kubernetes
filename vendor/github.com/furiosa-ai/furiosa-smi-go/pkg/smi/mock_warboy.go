package smi

import (
	"fmt"
)

var _ Device = new(staticWarboyMockDevice)

type staticWarboyMockDevice struct {
	arch    Arch
	nodeIdx int
}

func (m *staticWarboyMockDevice) DeviceInfo() (DeviceInfo, error) {
	return &staticWarboyMockDeviceInfo{
		nodeIdx: m.nodeIdx,
	}, nil
}

func (m *staticWarboyMockDevice) DeviceFiles() ([]DeviceFile, error) {
	return []DeviceFile{
		&staticMockDeviceFile{
			cores: []uint32{0},
			path:  fmt.Sprintf("/dev/npu%dpe0", m.nodeIdx),
		},
		&staticMockDeviceFile{
			cores: []uint32{1},
			path:  fmt.Sprintf("/dev/npu%dpe1", m.nodeIdx),
		},
		&staticMockDeviceFile{
			cores: []uint32{0, 1},
			path:  fmt.Sprintf("/dev/npu%dpe0-1", m.nodeIdx),
		},
	}, nil
}

func (m *staticWarboyMockDevice) CoreStatus() (map[uint32]CoreStatus, error) {
	return map[uint32]CoreStatus{0: CoreStatusAvailable, 1: CoreStatusAvailable}, nil
}

func (m *staticWarboyMockDevice) DeviceErrorInfo() (DeviceErrorInfo, error) {
	return &staticMockDeviceErrorInfo{}, nil
}

func (m *staticWarboyMockDevice) Liveness() (bool, error) {
	return true, nil
}

func (m *staticWarboyMockDevice) CoreUtilization() (CoreUtilization, error) {
	return &staticMockCoreUtilization{
		pe: []PeUtilization{
			&staticMockPeUtilization{core: 0, timeWindow: 1000, usage: 50},
		},
	}, nil
}

func (m *staticWarboyMockDevice) MemoryUtilization() (MemoryUtilization, error) {
	return &staticMockMemoryUtilization{}, nil
}

func (m *staticWarboyMockDevice) PowerConsumption() (float64, error) {
	return float64(100), nil
}

func (m *staticWarboyMockDevice) DeviceTemperature() (DeviceTemperature, error) {
	return &staticMockDeviceTemperature{}, nil
}

func (m *staticWarboyMockDevice) DeviceToDeviceLinkType(target Device) (LinkType, error) {
	selfNodeIdx := m.nodeIdx
	targetNodeIdx := target.(*staticWarboyMockDevice).nodeIdx

	if selfNodeIdx > targetNodeIdx {
		selfNodeIdx, targetNodeIdx = targetNodeIdx, selfNodeIdx
	}

	ret := linkTypeHintMap[selfNodeIdx][targetNodeIdx]
	return ret, nil
}

func (m *staticWarboyMockDevice) P2PAccessible(_ Device) (bool, error) {
	return false, nil
}

func (m *staticWarboyMockDevice) DevicePerformanceCounter() (DevicePerformanceCounter, error) {
	return &staticMockDevicePerformanceCounter{}, nil
}

type staticWarboyMockDeviceInfo struct {
	nodeIdx int
}

func (m *staticWarboyMockDeviceInfo) Index() uint32 {
	return uint32(m.nodeIdx)
}

var _ DeviceInfo = new(staticWarboyMockDeviceInfo)

func (m *staticWarboyMockDeviceInfo) Arch() Arch {
	return ArchWarboy
}

func (m *staticWarboyMockDeviceInfo) CoreNum() uint32 {
	return 2
}

func (m *staticWarboyMockDeviceInfo) NumaNode() uint32 {
	return uint32(staticMockHintMap[m.nodeIdx].numaNode)
}

func (m *staticWarboyMockDeviceInfo) Name() string {
	return fmt.Sprintf("/dev/npu%d", m.nodeIdx)
}

func (m *staticWarboyMockDeviceInfo) Serial() string {
	return staticMockHintMap[m.nodeIdx].serial
}

func (m *staticWarboyMockDeviceInfo) UUID() string {
	return staticMockHintMap[m.nodeIdx].uuid
}

func (m *staticWarboyMockDeviceInfo) BDF() string {
	return staticMockHintMap[m.nodeIdx].bdf
}

func (m *staticWarboyMockDeviceInfo) Major() uint16 {
	return staticMockHintMap[m.nodeIdx].major
}

func (m *staticWarboyMockDeviceInfo) Minor() uint16 {
	return staticMockHintMap[m.nodeIdx].minor
}

// FirmwareVersion e.g. version: 1.6.0, c1bebfd
func (m *staticWarboyMockDeviceInfo) FirmwareVersion() VersionInfo {
	return newStaticMockVersionInfo(1, 6, 0, "c1bebfd")
}

func (m *staticWarboyMockDeviceInfo) PertVersion() VersionInfo {
	return newStaticMockVersionInfo(0, 0, 0, "")
}
