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
			description:     "expect HostBridge ",
			xmlTopologyPath: "./hwloc/test.xml",
			dev1:            NewMockWarboyDevice(0, 0, "0000:27:00.0", "", "", "", "", ""),
			dev2:            NewMockWarboyDevice(1, 0, "0000:2a:00.0", "", "", "", "", ""),
			expected:        LinkTypeHostBridge,
		},
		{
			description:     "expect CPU ",
			xmlTopologyPath: "./hwloc/test.xml",
			dev1:            NewMockWarboyDevice(0, 0, "0000:27:00.0", "", "", "", "", ""),
			dev2:            NewMockWarboyDevice(1, 0, "0000:51:00.0", "", "", "", "", ""),
			expected:        LinkTypeCPU,
		},
		{
			description:     "expect CrossCPU ",
			xmlTopologyPath: "./hwloc/test.xml",
			dev1:            NewMockWarboyDevice(0, 0, "0000:27:00.0", "", "", "", "", ""),
			dev2:            NewMockWarboyDevice(1, 1, "0000:c7:00.0", "", "", "", "", ""),
			expected:        LinkTypeCrossCPU,
		},
		{
			description:     "expect CrossCPU ",
			xmlTopologyPath: "./hwloc/test.xml",
			dev1:            NewMockWarboyDevice(0, 0, "0000:27:00.0", "", "", "", "", ""),
			dev2:            NewMockWarboyDevice(0, 0, "0000:27:00.0", "", "", "", "", ""),
			expected:        LinkTypeUnknown,
		},
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
