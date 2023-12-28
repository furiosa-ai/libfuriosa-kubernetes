package device

import (
	"errors"
	"io/fs"
	"strings"
	"testing"
)

func TestReadMgmtFile(t *testing.T) {
	tests := []struct {
		description      string
		inputSysFs       string
		inputMgmtFile    MgmtFile
		inputDeviceIndex uint8
		expectedResult   string
		expectedError    error
	}{
		{
			description:      "test busname",
			inputSysFs:       "../../testdata/device/testdata/test-0/sys",
			inputMgmtFile:    Busname,
			inputDeviceIndex: 0,
			expectedResult:   "0000:6d:00.0",
			expectedError:    nil,
		},
		{
			description:      "test dev",
			inputSysFs:       "../../testdata/device/testdata/test-0/sys",
			inputMgmtFile:    Dev,
			inputDeviceIndex: 0,
			expectedResult:   "234:0",
			expectedError:    nil,
		},
		{
			description:      "test device sn",
			inputSysFs:       "../../testdata/device/testdata/test-0/sys",
			inputMgmtFile:    DeviceSn,
			inputDeviceIndex: 0,
			expectedResult:   "WBYB0236FH505KREO",
			expectedError:    nil,
		},
		{
			description:      "test device type",
			inputSysFs:       "../../testdata/device/testdata/test-0/sys",
			inputMgmtFile:    DeviceType,
			inputDeviceIndex: 0,
			expectedResult:   "Warboy",
			expectedError:    nil,
		},
		{
			description:      "test device uuid",
			inputSysFs:       "../../testdata/device/testdata/test-0/sys",
			inputMgmtFile:    DeviceUuid,
			inputDeviceIndex: 0,
			expectedResult:   "A76AAD68-6855-40B1-9E86-D080852D1C84",
			expectedError:    nil,
		},
		{
			description:      "test platform type",
			inputSysFs:       "../../testdata/device/testdata/test-0/sys",
			inputMgmtFile:    PlatformType,
			inputDeviceIndex: 0,
			expectedResult:   "FuriosaAI",
			expectedError:    nil,
		},
		{
			description:      "test soc rev",
			inputSysFs:       "../../testdata/device/testdata/test-0/sys",
			inputMgmtFile:    SocRev,
			inputDeviceIndex: 0,
			expectedResult:   "B0",
			expectedError:    nil,
		},
		{
			description:      "test soc uid",
			inputSysFs:       "../../testdata/device/testdata/test-0/sys",
			inputMgmtFile:    SocUid,
			inputDeviceIndex: 0,
			expectedResult:   "A76AAD68-8224DC84",
			expectedError:    nil,
		},
		{
			description:      "test alive",
			inputSysFs:       "../../testdata/device/testdata/test-0/sys",
			inputMgmtFile:    Alive,
			inputDeviceIndex: 0,
			expectedResult:   "1",
			expectedError:    nil,
		},
		{
			description:      "test atr error",
			inputSysFs:       "../../testdata/device/testdata/test-0/sys",
			inputMgmtFile:    AtrError,
			inputDeviceIndex: 0,
			expectedResult:   "AXI Post Error: 0\nAXI Fetch Error: 0\nAXI Discard Error: 0\nAXI Doorbell done: 0\nPCIe Post Error: 0\nPCIe Fetch Error: 0\nPCIe Discard Error: 0\nPCIe Doorbell done: 0\nDevice Error: 0",
			expectedError:    nil,
		},
		{
			description:      "test fw version",
			inputSysFs:       "../../testdata/device/testdata/test-0/sys",
			inputMgmtFile:    FwVersion,
			inputDeviceIndex: 0,
			expectedResult:   "1.6.0, c1bebfd",
			expectedError:    nil,
		},
		{
			description:      "test heartbeat",
			inputSysFs:       "../../testdata/device/testdata/test-0/sys",
			inputMgmtFile:    Heartbeat,
			inputDeviceIndex: 0,
			expectedResult:   "14649",
			expectedError:    nil,
		},
		{
			description:      "test NeClkFreqInfo",
			inputSysFs:       "../../testdata/device/testdata/test-0/sys",
			inputMgmtFile:    NeClkFreqInfo,
			inputDeviceIndex: 0,
			expectedResult:   "ne tensor (MHz): 2000\nne memory (MHz): 1000\nne operation (MHz): 1000\ndram (MT/s): 4266",
			expectedError:    nil,
		},
		{
			description:      "test version",
			inputSysFs:       "../../testdata/device/testdata/test-0/sys",
			inputMgmtFile:    Version,
			inputDeviceIndex: 0,
			expectedResult:   "1.9.2, 3def9c2",
			expectedError:    nil,
		},
		{
			description:      "test with wrong sys path",
			inputSysFs:       "WRONG_PATH",
			inputMgmtFile:    Version,
			inputDeviceIndex: 0,
			expectedResult:   "",
			expectedError:    fs.ErrNotExist,
		},
	}

	for _, tc := range tests {
		actualResult, actualErr := ReadMgmtFile(tc.inputSysFs, tc.inputMgmtFile.Filename(), tc.inputDeviceIndex)

		if tc.expectedError != nil || actualErr != nil {
			if !errors.Is(actualErr, tc.expectedError) {
				t.Errorf("expected %v but got %v", tc.expectedError, actualErr)
				continue
			}
		}

		tc.expectedResult = strings.Trim(tc.expectedResult, " ")
		tc.expectedResult = strings.Trim(tc.expectedResult, "\n")
		actualResult = strings.Trim(actualResult, " ")
		actualResult = strings.Trim(actualResult, "\n")
		if tc.expectedResult != actualResult {
			t.Errorf("expected %s but got %s", tc.expectedResult, actualResult)
		}
	}

}

func TestReadNumaNode(t *testing.T) {
	tests := []struct {
		description        string
		inputSysFs         string
		inputBdfIdentifier string
		expectedResult     int
		expectedError      error
	}{
		{
			description:        "expect numa node 0",
			inputSysFs:         "../../testdata/device/testdata/test-0/sys",
			inputBdfIdentifier: "0000:6d:00.0",
			expectedResult:     0,
			expectedError:      nil,
		},
		{
			description:        "expect numa node -1",
			inputSysFs:         "../../testdata/device/testdata/test-0/sys",
			inputBdfIdentifier: "0000:ff:00.0",
			expectedResult:     -1,
			expectedError:      nil,
		},
		{
			description:        "expect error : wrong identifier",
			inputSysFs:         "../../testdata/device/testdata/test-0/sys",
			inputBdfIdentifier: "",
			expectedResult:     -1,
			expectedError:      fs.ErrNotExist,
		},
		{
			description:        "expect error : wrong sys path",
			inputSysFs:         "WRONG_PATH",
			inputBdfIdentifier: "",
			expectedResult:     -1,
			expectedError:      fs.ErrNotExist,
		},
	}

	for _, tc := range tests {
		actualResult, actualErr := ReadNumaNode(tc.inputSysFs, tc.inputBdfIdentifier)
		if actualErr != nil || tc.expectedError != nil {
			if !errors.Is(actualErr, tc.expectedError) {
				t.Errorf("expected %v but got %v", tc.expectedError, actualErr)
				continue
			}
		}

		if actualResult != tc.expectedResult {
			t.Errorf("expected %d but got %d", tc.expectedResult, actualResult)
		}
	}

}
