package device

import (
	"errors"
	"testing"
)

// TODO support renegade
func TestCollectDevices(t *testing.T) {
	tests := []struct {
		description            string
		inputPaths             []string
		expectedCoresLen       int
		expectedDeviceFilesLen int
	}{
		{
			description: "test warboy npu0",
			inputPaths: []string{
				Abs("../../test/device/testdata/test-0/dev/npu0"),
				Abs("../../test/device/testdata/test-0/dev/npu0pe0"),
				Abs("../../test/device/testdata/test-0/dev/npu0pe1"),
				Abs("../../test/device/testdata/test-0/dev/npu0pe0-1"),
			},
			expectedCoresLen:       2,
			expectedDeviceFilesLen: 4,
		},
		{
			description: "test warboy npu1",
			inputPaths: []string{
				Abs("../../test/device/testdata/test-0/dev/npu1"),
				Abs("../../test/device/testdata/test-0/dev/npu1pe0"),
				Abs("../../test/device/testdata/test-0/dev/npu1pe1"),
				Abs("../../test/device/testdata/test-0/dev/npu1pe0-1"),
			},
			expectedCoresLen:       2,
			expectedDeviceFilesLen: 4,
		},
	}

	for _, tc := range tests {
		actualCores, actualDeviceFiles := collectDevices(tc.inputPaths)
		if tc.expectedCoresLen != len(actualCores) {
			t.Errorf("expected %d but got %d", tc.expectedCoresLen, len(actualCores))
		}

		if tc.expectedDeviceFilesLen != len(actualDeviceFiles) {
			t.Errorf("expected %d but got %d", tc.expectedDeviceFilesLen, len(actualDeviceFiles))
		}
	}
}

// TODO support renegade
func TestNewDevice(t *testing.T) {
	tests := []struct {
		description    string
		inputDevIdx    uint8
		inputPaths     []string
		inputDevFs     string
		inputSysFs     string
		expectedResult Device
		expectedError  error
	}{
		{
			description: "test warboy npu0",
			inputDevIdx: 0,
			inputPaths: []string{
				Abs("../../test/device/testdata/test-0/dev/npu0"),
				Abs("../../test/device/testdata/test-0/dev/npu0pe0"),
				Abs("../../test/device/testdata/test-0/dev/npu0pe1"),
				Abs("../../test/device/testdata/test-0/dev/npu0pe0-1"),
			},
			inputDevFs: "../../test/device/testdata/test-0/dev",
			inputSysFs: "../../test/device/testdata/test-0/sys",
			expectedResult: &device{
				deviceIndex: 0,
				devRoot:     "../../test/device/testdata/test-0/dev",
				sysRoot:     "../../test/device/testdata/test-0/sys",
				arch:        ArchWarboy,
				meta: map[string]string{
					"busname":       "0000:6d:00.0",
					"dev":           "234:0",
					"device_sn":     "WBYB0236FH505KREO",
					"device_type":   "Warboy",
					"device_uuid":   "A76AAD68-6855-40B1-9E86-D080852D1C84",
					"platform_type": "FuriosaAI",
					"soc_rev":       "B0",
					"soc_uid":       "A76AAD68-8224DC84",
				},
				numaNode: 0,
				cores:    []uint8{0, 1},
				devFiles: []DeviceFile{
					&deviceFile{
						index: 0,
						coreRange: coreRange{
							coreRangeType: CoreRangeTypeAll,
							start:         nil,
							end:           nil,
						},
						path:       Abs("../../test/device/testdata/test-0/dev/npu0"),
						deviceMode: DeviceModeMultiCore,
					},
					&deviceFile{
						index: 0,
						coreRange: coreRange{
							coreRangeType: CoreRangeTypeRange,
							start:         newUint8Pointer(0),
							end:           newUint8Pointer(0),
						},
						path:       Abs("../../test/device/testdata/test-0/dev/npu0pe0"),
						deviceMode: DeviceModeSingle,
					},
					&deviceFile{
						index: 0,
						coreRange: coreRange{
							coreRangeType: CoreRangeTypeRange,
							start:         newUint8Pointer(1),
							end:           newUint8Pointer(1),
						},
						path:       Abs("../../test/device/testdata/test-0/dev/npu0pe1"),
						deviceMode: DeviceModeSingle,
					},
					&deviceFile{
						index: 0,
						coreRange: coreRange{
							coreRangeType: CoreRangeTypeRange,
							start:         newUint8Pointer(0),
							end:           newUint8Pointer(1),
						},
						path:       Abs("../../test/device/testdata/test-0/dev/npu0pe0-1"),
						deviceMode: DeviceModeFusion,
					},
				},
			},
			expectedError: nil,
		},
		{
			description: "test warboy npu1",
			inputDevIdx: 1,
			inputPaths: []string{
				Abs("../../test/device/testdata/test-0/dev/npu1"),
				Abs("../../test/device/testdata/test-0/dev/npu1pe0"),
				Abs("../../test/device/testdata/test-0/dev/npu1pe1"),
				Abs("../../test/device/testdata/test-0/dev/npu1pe0-1"),
			},
			inputDevFs: "../../test/device/testdata/test-0/dev",
			inputSysFs: "../../test/device/testdata/test-0/sys",
			expectedResult: &device{
				deviceIndex: 1,
				devRoot:     "../../test/device/testdata/test-0/dev",
				sysRoot:     "../../test/device/testdata/test-0/sys",
				arch:        ArchWarboy,
				meta: map[string]string{
					"busname":       "0000:ff:00.0",
					"dev":           "510:0",
					"device_sn":     "WBYB0236FH543KREO",
					"device_type":   "Warboy",
					"device_uuid":   "A76AAD68-96A2-4B6A-A879-F91B8224DC84",
					"platform_type": "FuriosaAI",
					"soc_rev":       "B0",
					"soc_uid":       "A76AAD68-852D1C84",
				},
				numaNode: -1,
				cores:    []uint8{0, 1},
				devFiles: []DeviceFile{
					&deviceFile{
						index: 1,
						coreRange: coreRange{
							coreRangeType: CoreRangeTypeAll,
							start:         nil,
							end:           nil,
						},
						path:       Abs("../../test/device/testdata/test-0/dev/npu1"),
						deviceMode: DeviceModeMultiCore,
					},
					&deviceFile{
						index: 1,
						coreRange: coreRange{
							coreRangeType: CoreRangeTypeRange,
							start:         newUint8Pointer(0),
							end:           newUint8Pointer(0),
						},
						path:       Abs("../../test/device/testdata/test-0/dev/npu1pe0"),
						deviceMode: DeviceModeSingle,
					},
					&deviceFile{
						index: 1,
						coreRange: coreRange{
							coreRangeType: CoreRangeTypeRange,
							start:         newUint8Pointer(1),
							end:           newUint8Pointer(1),
						},
						path:       Abs("../../test/device/testdata/test-0/dev/npu1pe1"),
						deviceMode: DeviceModeSingle,
					},
					&deviceFile{
						index: 1,
						coreRange: coreRange{
							coreRangeType: CoreRangeTypeRange,
							start:         newUint8Pointer(0),
							end:           newUint8Pointer(1),
						},
						path:       Abs("../../test/device/testdata/test-0/dev/npu1pe0-1"),
						deviceMode: DeviceModeFusion,
					},
				},
			},
			expectedError: nil,
		},
		{
			description:    "test wrong npu",
			inputDevIdx:    40,
			inputPaths:     nil,
			inputDevFs:     "WRONG_DEV_FS_PATH",
			inputSysFs:     "WRONG_SYS_FS_PATH",
			expectedResult: nil,
			expectedError:  DeviceNotfound,
		},
	}

	for _, tc := range tests {
		actualResult, actualErr := NewDevice(tc.inputDevIdx, tc.inputPaths, tc.inputDevFs, tc.inputSysFs)

		if tc.expectedError != nil || actualErr != nil {
			if !errors.Is(actualErr, tc.expectedError) {
				t.Errorf("expected %s but got %s", tc.expectedError, actualErr)
				continue
			}
		}

		if tc.expectedResult == nil {
			continue
		}

		if tc.expectedResult.Name() != actualResult.Name() {
			t.Errorf("expected %s but got %s", tc.expectedResult.Name(), actualResult.Name())
			continue
		}

		if tc.expectedResult.DeviceIndex() != actualResult.DeviceIndex() {
			t.Errorf("expected %d but got %d", tc.expectedResult.DeviceIndex(), actualResult.DeviceIndex())
			continue
		}

		if tc.expectedResult.Arch() != actualResult.Arch() {
			t.Errorf("expected %s but got %s", tc.expectedResult.Arch(), actualResult.Arch())
			continue
		}

		expectedAliveValue, expectedAliveError := tc.expectedResult.Alive()
		actualAliveValue, actualAliveError := actualResult.Alive()
		if expectedAliveError != nil || actualAliveError != nil {
			if !errors.Is(expectedAliveError, actualAliveError) {
				t.Errorf("expected %s but got %s", expectedAliveError, actualAliveError)
				continue
			}
		}

		if expectedAliveValue != actualAliveValue {
			t.Errorf("expected %t but got %t", expectedAliveValue, actualAliveValue)
			continue
		}

		expectedAtrErrValue, expectedAtrErrError := tc.expectedResult.AtrErr()
		actualAtrErrValue, actualAtrErrError := actualResult.AtrErr()
		if expectedAtrErrError != nil || actualAtrErrError != nil {
			if !errors.Is(expectedAtrErrError, actualAtrErrError) {
				t.Errorf("expected %s but got %s", expectedAtrErrError, actualAtrErrError)
				continue
			}
		}

		for key, expectedValue := range expectedAtrErrValue {
			if expectedValue != actualAtrErrValue[key] {
				t.Errorf("expected value %d for key %s but got %d", expectedValue, key, actualAtrErrValue[key])
				continue
			}
		}

		expectedBusnameValue, expectedBusnameError := tc.expectedResult.Busname()
		actualBusnameValue, actualBusnameError := actualResult.Busname()
		if expectedBusnameError != nil || actualBusnameError != nil {
			if !errors.Is(expectedBusnameError, actualBusnameError) {
				t.Errorf("expected %s but got %s", expectedBusnameError, actualBusnameError)
				continue
			}
		}

		if expectedBusnameValue != actualBusnameValue {
			t.Errorf("expected %s but got %s", expectedBusnameValue, actualBusnameValue)
			continue
		}

		expectedPCIDevValue, expectedPCIDevError := tc.expectedResult.PCIDev()
		actualPCIDevValue, actualPCIDevError := actualResult.PCIDev()
		if expectedPCIDevError != nil || actualPCIDevError != nil {
			if !errors.Is(expectedPCIDevError, actualPCIDevError) {
				t.Errorf("expected %s but got %s", expectedPCIDevError, actualPCIDevError)
				continue
			}
		}

		if expectedPCIDevValue != actualPCIDevValue {
			t.Errorf("expected %s but got %s", expectedPCIDevValue, actualPCIDevValue)
			continue
		}

		expectedDeviceSnValue, expectedDeviceSnError := tc.expectedResult.DeviceSn()
		actualDeviceSnValue, actualDeviceSnError := actualResult.DeviceSn()
		if expectedDeviceSnError != nil || actualDeviceSnError != nil {
			if !errors.Is(expectedDeviceSnError, actualDeviceSnError) {
				t.Errorf("expected %s but got %s", expectedDeviceSnError, actualDeviceSnError)
				continue
			}
		}

		if expectedDeviceSnValue != actualDeviceSnValue {
			t.Errorf("expected %s but got %s", expectedDeviceSnValue, actualDeviceSnValue)
			continue
		}

		expectedDeviceUUIDValue, expectedDeviceUUIDError := tc.expectedResult.DeviceUUID()
		actualDeviceUUIDValue, actualDeviceUUIDError := actualResult.DeviceUUID()
		if expectedDeviceUUIDError != nil || actualDeviceUUIDError != nil {
			if !errors.Is(expectedDeviceUUIDError, actualDeviceUUIDError) {
				t.Errorf("expected %s but got %s", expectedDeviceUUIDError, actualDeviceUUIDError)
				continue
			}
		}

		if expectedDeviceUUIDValue != actualDeviceUUIDValue {
			t.Errorf("expected %s but got %s", expectedDeviceUUIDValue, actualDeviceUUIDValue)
			continue
		}

		expectedFirmwareVersionValue, expectedFirmwareVersionError := tc.expectedResult.FirmwareVersion()
		actualFirmwareVersionValue, actualFirmwareVersionError := actualResult.FirmwareVersion()
		if expectedFirmwareVersionError != nil || actualFirmwareVersionError != nil {
			if !errors.Is(expectedFirmwareVersionError, actualFirmwareVersionError) {
				t.Errorf("expected %s but got %s", expectedFirmwareVersionError, actualFirmwareVersionError)
				continue
			}
		}

		if expectedFirmwareVersionValue != actualFirmwareVersionValue {
			t.Errorf("expected %s but got %s", expectedFirmwareVersionValue, actualFirmwareVersionValue)
			continue
		}

		expectedDriverVersionValue, expectedDriverVersionError := tc.expectedResult.DriverVersion()
		actualDriverVersionValue, actualDriverVersionError := actualResult.DriverVersion()
		if expectedDriverVersionError != nil || actualDriverVersionError != nil {
			if !errors.Is(expectedDriverVersionError, actualDriverVersionError) {
				t.Errorf("expected %s but got %s", expectedDriverVersionError, actualDriverVersionError)
				continue
			}
		}

		if expectedDriverVersionValue != actualDriverVersionValue {
			t.Errorf("expected %s but got %s", expectedDriverVersionValue, actualDriverVersionValue)
			continue
		}

		expectedHeartBeatValue, expectedHeartBeatError := tc.expectedResult.HeartBeat()
		actualHeartBeatValue, actualHeartBeatError := actualResult.HeartBeat()
		if expectedHeartBeatError != nil || actualHeartBeatError != nil {
			if !errors.Is(expectedHeartBeatError, actualHeartBeatError) {
				t.Errorf("expected %s but got %s", expectedHeartBeatError, actualHeartBeatError)
				continue
			}
		}

		if expectedHeartBeatValue != actualHeartBeatValue {
			t.Errorf("expected %d but got %d", expectedHeartBeatValue, actualHeartBeatValue)
			continue
		}

		expectedNumaNodeValue, expectedNumaNodeError := tc.expectedResult.NumaNode()
		actualNumaNodeValue, actualNumaNodeError := actualResult.NumaNode()
		if expectedNumaNodeError != nil || actualNumaNodeError != nil {
			if !errors.Is(expectedNumaNodeError, UnexpectedValue) || !errors.Is(actualNumaNodeError, UnexpectedValue) {
				t.Errorf("expected %s but got %s", expectedNumaNodeError, actualNumaNodeError)
				continue
			}
		}

		if expectedNumaNodeValue != actualNumaNodeValue {
			t.Errorf("expected %d but got %d", expectedNumaNodeValue, actualNumaNodeValue)
			continue
		}

		if tc.expectedResult.CoreNum() != actualResult.CoreNum() {
			t.Errorf("expected %d but got %d", tc.expectedResult.CoreNum(), actualResult.CoreNum())
			continue
		}

		if len(tc.expectedResult.Cores()) != len(actualResult.Cores()) {
			t.Errorf("expected %d but got %d", len(tc.expectedResult.Cores()), len(actualResult.Cores()))
			continue
		}

		for _, expectedElement := range tc.expectedResult.Cores() {
			found := false
			for _, actualElement := range actualResult.Cores() {
				if expectedElement == actualElement {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("expected %d but not found", expectedElement)
				break
			}
		}

		if len(tc.expectedResult.DevFiles()) != len(actualResult.DevFiles()) {
			t.Errorf("expected %d but got %d", len(tc.expectedResult.DevFiles()), len(actualResult.DevFiles()))
			continue
		}

		for _, expectedDeviceFile := range tc.expectedResult.DevFiles() {
			found := false
			for _, actualDeviceFile := range actualResult.DevFiles() {
				if expectedDeviceFile.Filename() != actualDeviceFile.Filename() {
					continue
				}

				if expectedDeviceFile.Path() != actualDeviceFile.Path() {
					continue
				}

				if expectedDeviceFile.Mode() != actualDeviceFile.Mode() {
					continue
				}

				if expectedDeviceFile.DeviceIndex() != actualDeviceFile.DeviceIndex() {
					continue
				}

				if expectedDeviceFile.CoreRange().Type() != actualDeviceFile.CoreRange().Type() {
					continue
				}

				expected := safeDerefUint8(expectedDeviceFile.CoreRange().Start())
				actual := safeDerefUint8(actualDeviceFile.CoreRange().Start())
				if expected != actual {
					t.Errorf("expected %d but got %d", expected, actual)
					continue
				}

				expected = safeDerefUint8(expectedDeviceFile.CoreRange().End())
				actual = safeDerefUint8(actualDeviceFile.CoreRange().End())
				if expected != actual {
					t.Errorf("expected %d but got %d", expected, actual)
					continue
				}

				found = true
				break
			}

			if !found {
				t.Errorf("expected %s but not found", expectedDeviceFile.Filename())
				break
			}
		}

		//TODO: we need distinguish warboy and renegade here
		for _, coreIdx := range []uint8{0, 1} {
			expectedGetStatusCoreValue, expectedGetStatusCoreError := tc.expectedResult.GetStatusCore(coreIdx)
			actualGetStatusCoreValue, actualGetStatusCoreError := actualResult.GetStatusCore(coreIdx)
			if expectedGetStatusCoreError != nil || actualGetStatusCoreError != nil {
				if !errors.Is(expectedGetStatusCoreError, actualGetStatusCoreError) {
					t.Errorf("expected %s but got %s", expectedGetStatusCoreError, actualGetStatusCoreError)
					continue
				}
			}

			if expectedGetStatusCoreValue != actualGetStatusCoreValue {
				t.Errorf("expected %s but got %s", expectedDriverVersionValue, actualDriverVersionValue)
				continue
			}
		}

		expectedGetStatusAllValue, expectedGetStatusAllError := tc.expectedResult.GetStatusAll()
		actualGetStatusAllValue, actualGetStatusAllError := actualResult.GetStatusAll()
		if expectedGetStatusAllError != nil || actualGetStatusAllError != nil {
			if !errors.Is(expectedGetStatusAllError, actualGetStatusAllError) {
				t.Errorf("expected %s but got %s", expectedGetStatusAllError, actualGetStatusAllError)
				continue
			}
		}

		for key, expectedValue := range expectedGetStatusAllValue {
			if expectedValue != actualGetStatusAllValue[key] {
				t.Errorf("expected value %s for key %d but got %s", expectedValue, key, actualGetStatusAllValue[key])
				continue
			}
		}

	}
}
