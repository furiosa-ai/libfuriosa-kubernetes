package hwloc

import "testing"

func TestGetCommonAncestorObjType(t *testing.T) {
	tests := []struct {
		description     string
		xmlTopologyPath string
		bdf1            string
		bdf2            string
		expectedResult  HwlocObjType
		expectedError   bool
	}{
		{
			description:     "expect bridge",
			xmlTopologyPath: "./test.xml",
			bdf1:            "0000:5e:00.0",
			bdf2:            "0000:3b:00.0",
			expectedResult:  HwlocObjTypePackage,
			expectedError:   false,
		},
		{
			description:     "expect machine",
			xmlTopologyPath: "./test.xml",
			bdf1:            "0000:5e:00.0",
			bdf2:            "0000:af:00.0",
			expectedResult:  HwlocObjTypeMachine,
			expectedError:   false,
		},
		{
			description:     "expect unknown",
			xmlTopologyPath: "./test.xml",
			bdf1:            "0000:00:00.0",
			bdf2:            "0000:00:00.0",
			expectedResult:  HwlocObjTypeUnknown,
			expectedError:   true,
		},
	}

	for _, tc := range tests {
		mock := NewMockHwloc(tc.xmlTopologyPath)
		err := mock.TopologyInit()
		if err != nil {
			t.Errorf("unexpected error: %t", err)
			continue
		}

		err = mock.SetIoTypeFilter()
		if err != nil {
			t.Errorf("unexpected error: %t", err)
			continue
		}

		err = mock.TopologyLoad()
		if err != nil {
			t.Errorf("unexpected error: %t", err)
			continue
		}

		objType, err := mock.GetCommonAncestorObjType(tc.bdf1, tc.bdf2)
		if err != nil != tc.expectedError {
			t.Errorf("unexpected error: %t", err)
			continue
		}

		if objType != tc.expectedResult {
			t.Errorf("expectedResult %v but got %v", tc.expectedResult, objType)
			continue
		}

		mock.TopologyDestroy()
	}

}
