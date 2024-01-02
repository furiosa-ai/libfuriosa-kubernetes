package device

import (
	"errors"
	"testing"
)

func TestListDevices(t *testing.T) {
	tests := []struct {
		description       string
		inputSysFs        string
		inputDevFs        string
		expectedDeviceLen int
		expectedError     error
	}{
		{
			description:       "test positive case",
			inputDevFs:        Abs("../../testdata/device/testdata/test-0/dev"),
			inputSysFs:        Abs("../../testdata/device/testdata/test-0/sys"),
			expectedDeviceLen: 2,
			expectedError:     nil,
		},
		{
			description:       "test negative case",
			inputDevFs:        "WRONG_DEV_FS",
			inputSysFs:        "WRONG_DEV_FS",
			expectedDeviceLen: 0,
			expectedError:     nil,
		},
	}

	for _, tc := range tests {
		mockDeviceLister := newMockDeviceLister(tc.inputDevFs, tc.inputSysFs)
		actualResult, actualErr := mockDeviceLister.ListDevices()

		if tc.expectedError != nil || actualErr != nil {
			if errors.Is(actualErr, tc.expectedError) {
				t.Errorf("expected %s but got %s", tc.expectedError, actualErr)
				continue
			}
		}

		if tc.expectedDeviceLen != len(actualResult) {
			t.Errorf("expected %d but got %d", tc.expectedDeviceLen, len(actualResult))
			continue
		}
	}
}

func newMockDeviceLister(devFs string, sysFs string) DeviceLister {
	return newDeviceLister(devFs, sysFs, func(dev DevFile) bool {
		return true
	})
}
