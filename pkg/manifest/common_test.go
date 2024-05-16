package manifest

import (
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/device"
	"reflect"
	"tags.cncf.io/container-device-interface/pkg/cdi"
	"testing"
)

type dummyCoreRange struct {
	coreRangeType device.CoreRangeType
	start         uint8
	end           uint8
}

var _ device.CoreRange = new(dummyCoreRange)

func (d dummyCoreRange) Type() device.CoreRangeType {
	return d.coreRangeType
}

func (d dummyCoreRange) Start() uint8 {
	return d.start
}

func (d dummyCoreRange) End() uint8 {
	return d.end
}

func (d dummyCoreRange) Contains(_ uint8) bool {
	return false
}

type dummyDeviceFile struct {
	coreRange dummyCoreRange
}

var _ device.DeviceFile = new(dummyDeviceFile)

func (d dummyDeviceFile) Path() string {
	return ""
}

func (d dummyDeviceFile) Filename() string {
	return ""
}

func (d dummyDeviceFile) DeviceIndex() uint8 {
	return 0
}

func (d dummyDeviceFile) CoreRange() device.CoreRange {
	return d.coreRange
}

func (d dummyDeviceFile) Mode() device.DeviceMode {
	panic("implement me")
}

func genWarboyDeviceFiles() []device.DeviceFile {
	return []device.DeviceFile{
		&dummyDeviceFile{
			coreRange: dummyCoreRange{
				coreRangeType: device.CoreRangeTypeAll,
				start:         0,
				end:           0,
			},
		},
		&dummyDeviceFile{
			coreRange: dummyCoreRange{
				coreRangeType: device.CoreRangeTypeRange,
				start:         0,
				end:           0,
			},
		},
		&dummyDeviceFile{
			coreRange: dummyCoreRange{
				coreRangeType: device.CoreRangeTypeRange,
				start:         1,
				end:           1,
			},
		},
		&dummyDeviceFile{
			coreRange: dummyCoreRange{
				coreRangeType: device.CoreRangeTypeRange,
				start:         0,
				end:           1,
			},
		},
	}
}

func genRngdDeviceFiles() []device.DeviceFile {
	return []device.DeviceFile{
		&dummyDeviceFile{
			coreRange: dummyCoreRange{
				coreRangeType: device.CoreRangeTypeAll,
				start:         0,
				end:           0,
			},
		},
		&dummyDeviceFile{
			coreRange: dummyCoreRange{
				coreRangeType: device.CoreRangeTypeRange,
				start:         0,
				end:           0,
			},
		},
		&dummyDeviceFile{
			coreRange: dummyCoreRange{
				coreRangeType: device.CoreRangeTypeRange,
				start:         1,
				end:           1,
			},
		},
		&dummyDeviceFile{
			coreRange: dummyCoreRange{
				coreRangeType: device.CoreRangeTypeRange,
				start:         2,
				end:           2,
			},
		},
		&dummyDeviceFile{
			coreRange: dummyCoreRange{
				coreRangeType: device.CoreRangeTypeRange,
				start:         3,
				end:           3,
			},
		},
		&dummyDeviceFile{
			coreRange: dummyCoreRange{
				coreRangeType: device.CoreRangeTypeRange,
				start:         0,
				end:           1,
			},
		},
		&dummyDeviceFile{
			coreRange: dummyCoreRange{
				coreRangeType: device.CoreRangeTypeRange,
				start:         2,
				end:           3,
			},
		},
		&dummyDeviceFile{
			coreRange: dummyCoreRange{
				coreRangeType: device.CoreRangeTypeRange,
				start:         0,
				end:           3,
			},
		},
		// skip rngd pe 4~7
	}
}

func TestCollectDevFiles(t *testing.T) {
	tests := []struct {
		description string
		deviceFiles []device.DeviceFile
		coreStart   uint8
		coreEnd     uint8
		expected    []device.DeviceFile
	}{
		{
			description: "test partitioned warboy PE 0",
			deviceFiles: genWarboyDeviceFiles(),
			coreStart:   0,
			coreEnd:     0,
			expected: []device.DeviceFile{
				&dummyDeviceFile{
					coreRange: dummyCoreRange{
						coreRangeType: device.CoreRangeTypeRange,
						start:         0,
						end:           0,
					},
				},
			},
		},
		{
			description: "test partitioned warboy PE 1",
			deviceFiles: genWarboyDeviceFiles(),
			coreStart:   1,
			coreEnd:     1,
			expected: []device.DeviceFile{
				&dummyDeviceFile{
					coreRange: dummyCoreRange{
						coreRangeType: device.CoreRangeTypeRange,
						start:         1,
						end:           1,
					},
				},
			},
		},
		{
			description: "test partitioned rngd PE 0",
			deviceFiles: genRngdDeviceFiles(),
			coreStart:   0,
			coreEnd:     0,
			expected: []device.DeviceFile{
				&dummyDeviceFile{
					coreRange: dummyCoreRange{
						coreRangeType: device.CoreRangeTypeRange,
						start:         0,
						end:           0,
					},
				},
			},
		},
		{
			description: "test partitioned rngd PE 1",
			deviceFiles: genRngdDeviceFiles(),
			coreStart:   1,
			coreEnd:     1,
			expected: []device.DeviceFile{
				&dummyDeviceFile{
					coreRange: dummyCoreRange{
						coreRangeType: device.CoreRangeTypeRange,
						start:         1,
						end:           1,
					},
				},
			},
		},
		{
			description: "test partitioned rngd PE 2",
			deviceFiles: genRngdDeviceFiles(),
			coreStart:   2,
			coreEnd:     2,
			expected: []device.DeviceFile{
				&dummyDeviceFile{
					coreRange: dummyCoreRange{
						coreRangeType: device.CoreRangeTypeRange,
						start:         2,
						end:           2,
					},
				},
			},
		},
		{
			description: "test partitioned rngd PE 3",
			deviceFiles: genRngdDeviceFiles(),
			coreStart:   3,
			coreEnd:     3,
			expected: []device.DeviceFile{
				&dummyDeviceFile{
					coreRange: dummyCoreRange{
						coreRangeType: device.CoreRangeTypeRange,
						start:         3,
						end:           3,
					},
				},
			},
		},
		{
			description: "test partitioned rngd PE 0-1",
			deviceFiles: genRngdDeviceFiles(),
			coreStart:   0,
			coreEnd:     1,
			expected: []device.DeviceFile{
				&dummyDeviceFile{
					coreRange: dummyCoreRange{
						coreRangeType: device.CoreRangeTypeRange,
						start:         0,
						end:           0,
					},
				},
				&dummyDeviceFile{
					coreRange: dummyCoreRange{
						coreRangeType: device.CoreRangeTypeRange,
						start:         1,
						end:           1,
					},
				},
				&dummyDeviceFile{
					coreRange: dummyCoreRange{
						coreRangeType: device.CoreRangeTypeRange,
						start:         0,
						end:           1,
					},
				},
			},
		},
		{
			description: "test partitioned rngd PE 2-3",
			deviceFiles: genRngdDeviceFiles(),
			coreStart:   2,
			coreEnd:     3,
			expected: []device.DeviceFile{
				&dummyDeviceFile{
					coreRange: dummyCoreRange{
						coreRangeType: device.CoreRangeTypeRange,
						start:         2,
						end:           2,
					},
				},
				&dummyDeviceFile{
					coreRange: dummyCoreRange{
						coreRangeType: device.CoreRangeTypeRange,
						start:         3,
						end:           3,
					},
				},
				&dummyDeviceFile{
					coreRange: dummyCoreRange{
						coreRangeType: device.CoreRangeTypeRange,
						start:         2,
						end:           3,
					},
				},
			},
		},
		{
			description: "test partitioned rngd PE 0-3",
			deviceFiles: genRngdDeviceFiles(),
			coreStart:   0,
			coreEnd:     3,
			expected: []device.DeviceFile{
				&dummyDeviceFile{
					coreRange: dummyCoreRange{
						coreRangeType: device.CoreRangeTypeRange,
						start:         0,
						end:           0,
					},
				},
				&dummyDeviceFile{
					coreRange: dummyCoreRange{
						coreRangeType: device.CoreRangeTypeRange,
						start:         1,
						end:           1,
					},
				},
				&dummyDeviceFile{
					coreRange: dummyCoreRange{
						coreRangeType: device.CoreRangeTypeRange,
						start:         2,
						end:           2,
					},
				},
				&dummyDeviceFile{
					coreRange: dummyCoreRange{
						coreRangeType: device.CoreRangeTypeRange,
						start:         3,
						end:           3,
					},
				},
				&dummyDeviceFile{
					coreRange: dummyCoreRange{
						coreRangeType: device.CoreRangeTypeRange,
						start:         0,
						end:           1,
					},
				},
				&dummyDeviceFile{
					coreRange: dummyCoreRange{
						coreRangeType: device.CoreRangeTypeRange,
						start:         2,
						end:           3,
					},
				},
				&dummyDeviceFile{
					coreRange: dummyCoreRange{
						coreRangeType: device.CoreRangeTypeRange,
						start:         0,
						end:           3,
					},
				},
			},
		},
		// skip rngd pe 4~7
	}
	for _, tc := range tests {
		result := collectDevFiles(tc.deviceFiles, tc.coreStart, tc.coreEnd)

		if !reflect.DeepEqual(tc.expected, result) {
			t.Errorf("%s: expected %v, got %v", tc.description, tc.expected, result)
		}
	}
}

func TestToCDIContainerEdits(t *testing.T) {
	tests := []struct {
		description string
		manifest    Manifest
		expected    *cdi.ContainerEdits
	}{
		{
			//tbd
		},
	}
	for _, tc := range tests {
		actual := toCDIContainerEdits(tc.manifest)

		if !reflect.DeepEqual(tc.expected, actual) {
			t.Errorf("%s: expected %v, got %v", tc.description, tc.expected, actual)
		}
	}
}
