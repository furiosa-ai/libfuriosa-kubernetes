package device

import (
	"errors"
	"io/fs"
	"path/filepath"
	"slices"
	"testing"
)

func TestListDevFs(t *testing.T) {
	tests := []struct {
		description    string
		inputDevFs     string
		expectedResult []DevFile
		expectedError  error
	}{
		{
			description: "test positive case",
			inputDevFs:  "../../test/device/testdata/test-0/dev",
			expectedResult: []DevFile{
				{
					fileAbsPath: "/npu0",
				},
				{
					fileAbsPath: "/npu0_mgmt",
				},
				{
					fileAbsPath: "/npu0pe0",
				},
				{
					fileAbsPath: "/npu0pe0-1",
				},
				{
					fileAbsPath: "/npu0pe1",
				},
				{
					fileAbsPath: "/npu1",
				},
				{
					fileAbsPath: "/npu1_mgmt",
				},
				{
					fileAbsPath: "/npu1pe0",
				},
				{
					fileAbsPath: "/npu1pe0-1",
				},
				{
					fileAbsPath: "/npu1pe1",
				},
			},
			expectedError: nil,
		},
		{
			description:    "test negative case",
			inputDevFs:     "WRONG_DEV_PATH",
			expectedResult: nil,
			expectedError:  fs.ErrNotExist,
		},
	}

	for _, tc := range tests {
		actualResult, actualErr := ListDevFs(tc.inputDevFs)
		if !errors.Is(actualErr, tc.expectedError) {
			t.Errorf("expected %v but got %v", tc.expectedError, actualErr)
			continue
		}

		count := 0
		for _, actual := range actualResult {
			for _, expected := range tc.expectedResult {
				if filepath.Base(actual.fileAbsPath) == filepath.Base(expected.fileAbsPath) {
					count++
					break
				}
			}
		}

		if len(actualResult) != count {
			t.Errorf("expected %v but got %v", len(tc.expectedResult), count)
			continue
		}

	}
}

func TestParseIndices(t *testing.T) {
	tests := []struct {
		description       string
		inputDeviceFile   string
		expectedDeviceId  uint8
		expectedCoreRange []uint8
		expectedError     error
	}{
		{
			description:       "try parse wrong device",
			inputDeviceFile:   "gpu",
			expectedDeviceId:  0,
			expectedCoreRange: nil,
			expectedError:     IncompatibleDriver,
		},
		{
			description:       "try parse full device",
			inputDeviceFile:   "npu0",
			expectedDeviceId:  0,
			expectedCoreRange: []uint8{},
			expectedError:     nil,
		},
		{
			description:       "try parse incomplete full device",
			inputDeviceFile:   "npu",
			expectedDeviceId:  0,
			expectedCoreRange: nil,
			expectedError:     IncompatibleDriver,
		},
		{
			description:       "try parse single pe",
			inputDeviceFile:   "npu2pe0",
			expectedDeviceId:  2,
			expectedCoreRange: []uint8{0},
			expectedError:     nil,
		},
		{
			description:       "try parse incomplete single pe",
			inputDeviceFile:   "npu0pe",
			expectedDeviceId:  0,
			expectedCoreRange: nil,
			expectedError:     IncompatibleDriver,
		},
		{
			description:       "try parse range of pe",
			inputDeviceFile:   "npu3pe2-4",
			expectedDeviceId:  3,
			expectedCoreRange: []uint8{2, 3, 4},
			expectedError:     nil,
		},
		{
			description:       "try parse incomplete range of pe",
			inputDeviceFile:   "npu4pe2-",
			expectedDeviceId:  0,
			expectedCoreRange: nil,
			expectedError:     IncompatibleDriver,
		},
		{
			description:       "try parse oversize values",
			inputDeviceFile:   "npu259",
			expectedDeviceId:  0,
			expectedCoreRange: nil,
			expectedError:     IncompatibleDriver,
		},
	}

	for _, tc := range tests {
		actualDeviceId, actualCoreRange, actualErr := ParseIndices(tc.inputDeviceFile)

		if !errors.Is(actualErr, tc.expectedError) {
			t.Errorf("expected %v but got %v", tc.expectedError, actualErr)
			continue
		}

		if actualDeviceId != tc.expectedDeviceId {
			t.Errorf("expected %d but got %d", tc.expectedDeviceId, actualDeviceId)
			continue
		}

		slices.Sort(actualCoreRange)
		slices.Sort(tc.expectedCoreRange)

		if !slices.Equal(actualCoreRange, tc.expectedCoreRange) {
			t.Errorf("expected %v but got %v", tc.expectedCoreRange, actualCoreRange)
			continue
		}
	}
}

func TestFilterDevFiles(t *testing.T) {
	tests := []struct {
		description string
		input       []DevFile
		expected    map[uint8][]string
	}{
		{
			description: "try parse npu0 dev files",
			input: []DevFile{
				{
					fileAbsPath: "npu0",
				},
				{
					fileAbsPath: "npu0_mgmt",
				},
				{
					fileAbsPath: "npu0pe0",
				},
				{
					fileAbsPath: "npu0pe0-1",
				},
				{
					fileAbsPath: "npu0pe1",
				},
			},
			expected: map[uint8][]string{
				0: {
					"npu0",
					"npu0pe0",
					"npu0pe0-1",
					"npu0pe1",
				},
			},
		},
		{
			description: "try parse npu0, npu1 dev files",
			input: []DevFile{
				{
					fileAbsPath: "npu0",
				},
				{
					fileAbsPath: "npu0_mgmt",
				},
				{
					fileAbsPath: "npu0pe0",
				},
				{
					fileAbsPath: "npu0pe0-1",
				},
				{
					fileAbsPath: "npu0pe1",
				},
				{
					fileAbsPath: "npu1",
				},
				{
					fileAbsPath: "npu1_mgmt",
				},
				{
					fileAbsPath: "npu1pe0",
				},
				{
					fileAbsPath: "npu1pe0-1",
				},
				{
					fileAbsPath: "npu1pe1",
				},
			},
			expected: map[uint8][]string{
				0: {
					"npu0",
					"npu0pe0",
					"npu0pe0-1",
					"npu0pe1",
				},
				1: {
					"npu1",
					"npu1pe0",
					"npu1pe0-1",
					"npu1pe1",
				},
			},
		},
		{
			description: "try parse npu0 dev files and wrong dev files",
			input: []DevFile{
				{
					fileAbsPath: "npu0",
				},
				{
					fileAbsPath: "npu0_mgmt",
				},
				{
					fileAbsPath: "npu0pe0",
				},
				{
					fileAbsPath: "npu0pe0-1",
				},
				{
					fileAbsPath: "npu0pe1",
				},
				{
					fileAbsPath: "gpu1",
				},
				{
					fileAbsPath: "gpu1_mgmt",
				},
				{
					fileAbsPath: "gpu1pe0",
				},
				{
					fileAbsPath: "gpu1pe0-1",
				},
				{
					fileAbsPath: "gpu1pe1",
				},
			},
			expected: map[uint8][]string{
				0: {
					"npu0",
					"npu0pe0",
					"npu0pe0-1",
					"npu0pe1",
				},
			},
		},
	}
	for _, tc := range tests {
		actual := filterDevFiles(tc.input, func(dev DevFile) bool {
			return true
		})
		for key, value := range actual {
			expected := tc.expected[key]
			slices.Sort(value)
			slices.Sort(expected)

			if !slices.Equal(value, expected) {
				t.Errorf("expected %v but got %v", expected, value)
				continue
			}
		}
	}
}

func TestIsFuriosaDevice(t *testing.T) {
	tests := []struct {
		description string
		inputIndex  uint8
		inputSysFs  string
		expected    bool
	}{
		{
			description: "test npu0",
			inputIndex:  0,
			inputSysFs:  "../../test/device/testdata/test-0/sys",
			expected:    true,
		},
		{
			description: "test npu1",
			inputIndex:  1,
			inputSysFs:  "../../test/device/testdata/test-0/sys",
			expected:    true,
		},
		{
			description: "test npu2",
			inputIndex:  2,
			inputSysFs:  "../../test/device/testdata/test-0/sys",
			expected:    false,
		},
		{
			description: "test wrong path",
			inputIndex:  2,
			inputSysFs:  "WRONG_PATH",
			expected:    false,
		},
	}
	for _, tc := range tests {
		actual := IsFuriosaDevice(tc.inputIndex, tc.inputSysFs)
		if actual != tc.expected {
			t.Errorf("expected %t but got %t", tc.expected, actual)
			continue
		}
	}
}
