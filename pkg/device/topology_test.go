package device

import "testing"

func TestGetCommonAncestorObjType(t *testing.T) {
	tests := []struct {
		description     string
		xmlTopologyPath string
		dev1            Device
		dev2            Device
		query1          string
		query2          string
		expected        LinkType
	}{
		{
			description:     "expect HostBridge ",
			xmlTopologyPath: "./hwloc/test.xml",
			dev1:            NewMockWarboyDevice(0, 0, "0000:27:00.0", "", "", "", "", ""),
			dev2:            NewMockWarboyDevice(1, 0, "0000:2a:00.0", "", "", "", "", ""),
			query1:          "0000:27:00.0",
			query2:          "0000:2a:00.0",
			expected:        LinkTypeHostBridge,
		},
		{
			description:     "expect CPU ",
			xmlTopologyPath: "./hwloc/test.xml",
			dev1:            NewMockWarboyDevice(0, 0, "0000:27:00.0", "", "", "", "", ""),
			dev2:            NewMockWarboyDevice(1, 0, "0000:51:00.0", "", "", "", "", ""),
			query1:          "0000:27:00.0",
			query2:          "0000:51:00.0",
			expected:        LinkTypeCPU,
		},
		{
			description:     "expect CrossCPU ",
			xmlTopologyPath: "./hwloc/test.xml",
			dev1:            NewMockWarboyDevice(0, 0, "0000:27:00.0", "", "", "", "", ""),
			dev2:            NewMockWarboyDevice(1, 1, "0000:c7:00.0", "", "", "", "", ""),
			query1:          "0000:27:00.0",
			query2:          "0000:c7:00.0",
			expected:        LinkTypeInterconnect,
		},
		{
			description:     "expect Unknown ",
			xmlTopologyPath: "./hwloc/test.xml",
			dev1:            NewMockWarboyDevice(0, 0, "0000:27:00.0", "", "", "", "", ""),
			dev2:            NewMockWarboyDevice(0, 0, "0000:27:00.0", "", "", "", "", ""),
			query1:          "0000:27:00.0",
			query2:          "0000:27:00.0",
			expected:        LinkTypeSoc,
		},
		{
			description:     "expect Unknown2 ",
			xmlTopologyPath: "./hwloc/test.xml",
			dev1:            NewMockWarboyDevice(0, 0, "0000:27:00.0", "", "", "", "", ""),
			dev2:            NewMockWarboyDevice(0, 0, "0000:27:00.0", "", "", "", "", ""),
			query1:          "wrong_query1",
			query2:          "wrong_query2",
			expected:        LinkTypeUnknown,
		},
	}

	for _, tc := range tests {
		mockTopology, err := NewMockTopology([]Device{tc.dev1, tc.dev2}, tc.xmlTopologyPath)
		if err != nil {
			t.Errorf("unexpected error %t", err)
			continue
		}

		actual := mockTopology.GetLinkType(tc.query1, tc.query2)
		if tc.expected != actual {
			t.Errorf("expected %v but got %v", tc.expected, actual)
		}
	}
}
