package device

import "testing"

func TestGetCommonAncestorObjType(t *testing.T) {
	tests := []struct {
		description     string
		xmlTopologyPath string
		dev1            Device
		dev2            Device
		expected        LinkType
	}{
		{
			description:     "expect CPU ",
			xmlTopologyPath: "./hwloc/test.xml",
			dev1:            NewMockWarboyDevice(0, 0, "0000:5e:00.0", "", "", "", "", ""),
			dev2:            NewMockWarboyDevice(1, 0, "0000:3b:00.0", "", "", "", "", ""),
			expected:        LinkTypeCPU,
		},
		{
			description:     "expect CrossCPU ",
			xmlTopologyPath: "./hwloc/test.xml",
			dev1:            NewMockWarboyDevice(0, 0, "0000:5e:00.0", "", "", "", "", ""),
			dev2:            NewMockWarboyDevice(1, 1, "0000:af:00.0", "", "", "", "", ""),
			expected:        LinkTypeCrossCPU,
		},
		//TODO(@bg): add test case for LinkTypeHostBridge, it requires to replace ./hwloc/test.xml for proper test data.
	}

	for _, tc := range tests {

		mockTopology, err := NewMockTopology([]Device{tc.dev1, tc.dev2}, tc.xmlTopologyPath)
		if err != nil {
			t.Errorf("unexpected error %t", err)
		}

		actual := mockTopology.GetLinkType(tc.dev1, tc.dev2)
		if tc.expected != actual {
			t.Errorf("expected %v but got %v", tc.expected, actual)
		}
	}
}
