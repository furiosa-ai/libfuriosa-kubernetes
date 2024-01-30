package device

import (
	"fmt"
	"strconv"
)

type mockWarboyDevice struct {
	deviceIndex     uint8
	numaNode        int
	busname         string
	pciDev          string
	deviceSN        string
	firmwareVersion string
	driverVersion   string
	deviceUUID      string
	cores           []uint8
	devFiles        []DeviceFile
}

var _ Device = new(mockWarboyDevice)

func NewMockWarboyDevice(deviceIndex uint8, numaNode int, busname, pciDev, deviceSN, firmwareVersion, driverVersion, deviceUUID string) Device {
	return &mockWarboyDevice{
		deviceIndex:     deviceIndex,
		numaNode:        numaNode,
		cores:           []uint8{0, 1},
		busname:         busname,
		pciDev:          pciDev,
		deviceSN:        deviceSN,
		firmwareVersion: firmwareVersion,
		driverVersion:   driverVersion,
		deviceUUID:      deviceUUID,
		devFiles: []DeviceFile{
			&deviceFile{
				index: deviceIndex,
				coreRange: coreRange{
					coreRangeType: CoreRangeTypeAll,
					start:         0,
					end:           0,
				},
				path:       fmt.Sprintf("/dev/npu%d", deviceIndex),
				deviceMode: DeviceModeMultiCore,
			},
			&deviceFile{
				index: deviceIndex,
				coreRange: coreRange{
					coreRangeType: CoreRangeTypeRange,
					start:         0,
					end:           0,
				},
				path:       fmt.Sprintf("/dev/npu%dpe0", deviceIndex),
				deviceMode: DeviceModeSingle,
			},
			&deviceFile{
				index: deviceIndex,
				coreRange: coreRange{
					coreRangeType: CoreRangeTypeRange,
					start:         1,
					end:           1,
				},
				path:       fmt.Sprintf("/dev/npu%dpe1", deviceIndex),
				deviceMode: DeviceModeSingle,
			},
			&deviceFile{
				index: deviceIndex,
				coreRange: coreRange{
					coreRangeType: CoreRangeTypeRange,
					start:         0,
					end:           1,
				},
				path:       fmt.Sprintf("/dev/npu%dpe0-1", deviceIndex),
				deviceMode: DeviceModeFusion,
			},
		},
	}
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
	return m.busname, nil
}

func (m *mockWarboyDevice) PCIDev() (string, error) {
	return m.pciDev, nil
}

func (m *mockWarboyDevice) DeviceSn() (string, error) {
	return m.deviceSN, nil
}

func (m *mockWarboyDevice) DeviceUUID() (string, error) {
	return m.deviceUUID, nil
}

func (m *mockWarboyDevice) FirmwareVersion() (string, error) {
	return m.firmwareVersion, nil
}

func (m *mockWarboyDevice) DriverVersion() (string, error) {
	return m.driverVersion, nil
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
