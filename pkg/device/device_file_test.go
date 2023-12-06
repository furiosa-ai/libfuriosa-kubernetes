package device

import (
	"errors"
	"testing"
)

func TestNewDeviceFile(t *testing.T) {
	tests := []struct {
		description    string
		input          string
		expectedResult DeviceFile
		expectedError  error
	}{
		{
			description:    "test npu",
			input:          "/ASSUME_VALID_DEV_FS_PATH/npu0pe",
			expectedResult: nil,
			expectedError:  IncompatibleDriver,
		},
		{
			description: "test npu0",
			input:       "/ASSUME_VALID_DEV_FS_PATH/npu0",
			expectedResult: &deviceFile{
				index: 0,
				coreRange: coreRange{
					coreRangeType: CoreRangeTypeAll,
					start:         0,
					end:           0,
				},
				path:       "/ASSUME_VALID_DEV_FS_PATH/npu0",
				deviceMode: DeviceModeMultiCore,
			},
			expectedError: nil,
		},
		{
			description:    "test npu0pe",
			input:          "/ASSUME_VALID_DEV_FS_PATH/npu0pe",
			expectedResult: nil,
			expectedError:  IncompatibleDriver,
		},
		{
			description: "test npu0pe0",
			input:       "/ASSUME_VALID_DEV_FS_PATH/npu0pe0",
			expectedResult: &deviceFile{
				index: 0,
				coreRange: coreRange{
					coreRangeType: CoreRangeTypeRange,
					start:         0,
					end:           0,
				},
				path:       "/ASSUME_VALID_DEV_FS_PATH/npu0pe0",
				deviceMode: DeviceModeSingle,
			},
			expectedError: nil,
		},
		{
			description: "test npu0pe1",
			input:       "/ASSUME_VALID_DEV_FS_PATH/npu0pe1",
			expectedResult: &deviceFile{
				index: 0,
				coreRange: coreRange{
					coreRangeType: CoreRangeTypeRange,
					start:         1,
					end:           1,
				},
				path:       "/ASSUME_VALID_DEV_FS_PATH/npu0pe1",
				deviceMode: DeviceModeSingle,
			},
			expectedError: nil,
		},
		{
			description: "test npu0pe0-1",
			input:       "/ASSUME_VALID_DEV_FS_PATH/npu0pe0-1",
			expectedResult: &deviceFile{
				index: 0,
				coreRange: coreRange{
					coreRangeType: CoreRangeTypeRange,
					start:         0,
					end:           1,
				},
				path:       "/ASSUME_VALID_DEV_FS_PATH/npu0pe0-1",
				deviceMode: DeviceModeFusion,
			},
			expectedError: nil,
		},
		{
			description:    "test npu0pe0-",
			input:          "/ASSUME_VALID_DEV_FS_PATH/npu0pe",
			expectedResult: nil,
			expectedError:  IncompatibleDriver,
		},
	}

	for _, tc := range tests {
		actualResult, actualError := NewDeviceFile(tc.input)

		if actualError != nil || tc.expectedError != nil {
			if !errors.Is(actualError, tc.expectedError) {
				t.Errorf("expected %t but got %t", tc.expectedError, actualError)
			}

			continue
		}

		if actualResult.Filename() != tc.expectedResult.Filename() {
			t.Errorf("expected %s but got %s", tc.expectedResult.Filename(), actualResult.Filename())
			continue
		}

		if actualResult.Path() != tc.expectedResult.Path() {
			t.Errorf("expected %s but got %s", tc.expectedResult.Path(), actualResult.Path())
			continue
		}

		if actualResult.Mode() != tc.expectedResult.Mode() {
			t.Errorf("expected %s but got %s", tc.expectedResult.Mode(), actualResult.Mode())
			continue
		}

		if actualResult.DeviceIndex() != tc.expectedResult.DeviceIndex() {
			t.Errorf("expected %d but got %d", tc.expectedResult.DeviceIndex(), actualResult.DeviceIndex())
			continue
		}

		if actualResult.CoreRange().Type() != tc.expectedResult.CoreRange().Type() {
			t.Errorf("expected %s but got %s", tc.expectedResult.CoreRange().Type(), actualResult.CoreRange().Type())
			continue
		}

		if tc.expectedResult.CoreRange().Type() != CoreRangeTypeAll {
			expected := actualResult.CoreRange().Start()
			actual := tc.expectedResult.CoreRange().Start()
			if expected != actual {
				t.Errorf("expected %d but got %d", expected, actual)
				continue
			}

			expected = actualResult.CoreRange().End()
			actual = tc.expectedResult.CoreRange().End()
			if expected != actual {
				t.Errorf("expected %d but got %d", expected, actual)
				continue
			}
		}
	}
}
