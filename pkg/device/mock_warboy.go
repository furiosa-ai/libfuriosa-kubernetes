package device

import (
	"fmt"
	"strconv"
)

type mockWarboyDevice struct {
	deviceIndex uint8
	numaNode    int
	cores       []uint8
	devFiles    []DeviceFile
}

var _ Device = new(mockWarboyDevice)

func NewMockWarboyDevice(deviceIndex uint8, numaNode int) (Device, error) {
	return &mockWarboyDevice{
		deviceIndex: deviceIndex,
		numaNode:    numaNode,
		cores:       []uint8{0, 1},
		devFiles: []DeviceFile{
			&deviceFile{
				index: deviceIndex,
				coreRange: coreRange{
					coreRangeType: CoreRangeTypeAll,
					start:         0,
					end:           0,
				},
				path:       "/dev/npu0",
				deviceMode: DeviceModeMultiCore,
			},
			&deviceFile{
				index: deviceIndex,
				coreRange: coreRange{
					coreRangeType: CoreRangeTypeRange,
					start:         0,
					end:           0,
				},
				path:       "/dev/npu0pe0",
				deviceMode: DeviceModeSingle,
			},
			&deviceFile{
				index: deviceIndex,
				coreRange: coreRange{
					coreRangeType: CoreRangeTypeRange,
					start:         1,
					end:           1,
				},
				path:       "/dev/npu0pe1",
				deviceMode: DeviceModeSingle,
			},
			&deviceFile{
				index: deviceIndex,
				coreRange: coreRange{
					coreRangeType: CoreRangeTypeRange,
					start:         0,
					end:           1,
				},
				path:       "/dev/npu0pe0-1",
				deviceMode: DeviceModeFusion,
			},
		},
	}, nil
}

func (m *mockWarboyDevice) Name() string {
	return fmt.Sprintf(NpuExp, m.deviceIndex)
}

func (m *mockWarboyDevice) DeviceIndex() uint8 {
	return m.deviceIndex
}

func (m *mockWarboyDevice) Arch() Arch {
	return ArchWarboy
}

func (m *mockWarboyDevice) Alive() (bool, error) {
	return true, nil
}

func (m *mockWarboyDevice) AtrErr() (map[string]uint32, error) {
	return nil, nil
}

func (m *mockWarboyDevice) Busname() (string, error) {
	return "0000:6d:00.0", nil
}

func (m *mockWarboyDevice) PCIDev() (string, error) {
	return "234:0", nil
}

func (m *mockWarboyDevice) DeviceSn() (string, error) {
	return "WBYB0236FH505KREO", nil
}

func (m *mockWarboyDevice) DeviceUUID() (string, error) {
	return "A76AAD68-6855-40B1-9E86-D080852D1C84", nil
}

func (m *mockWarboyDevice) FirmwareVersion() (string, error) {
	return "1.6.0, c1bebfd", nil
}

func (m *mockWarboyDevice) DriverVersion() (string, error) {
	return "1.9.2, 3def9c2", nil
}

func (m *mockWarboyDevice) HeartBeat() (uint32, error) {
	return 1001, nil
}

func (m *mockWarboyDevice) NumaNode() (uint8, error) {
	if m.numaNode < 0 {
		return 0, NewUnexpectedValue(strconv.Itoa(m.numaNode))
	}

	return uint8(m.numaNode), nil
}

func (m *mockWarboyDevice) CoreNum() uint8 {
	return uint8(len(m.cores))
}

func (m *mockWarboyDevice) Cores() []uint8 {
	return m.cores
}

func (m *mockWarboyDevice) DevFiles() []DeviceFile {
	return m.devFiles
}

func (m *mockWarboyDevice) GetStatusCore(_ uint8) (CoreStatus, error) {
	return CoreStatusAvailable, nil
}

func (m *mockWarboyDevice) GetStatusAll() (map[uint8]CoreStatus, error) {
	return map[uint8]CoreStatus{0: CoreStatusAvailable, 1: CoreStatusAvailable}, nil
}
