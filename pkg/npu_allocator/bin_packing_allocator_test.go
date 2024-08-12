package npu_allocator

import (
	"fmt"
	"github.com/bradfitz/iter"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func generateSameBoardMockDeviceSet(amount int, topologyHintKey string) DeviceSet {
	id, _ := uuid.NewUUID()
	devices := make(DeviceSet, 0, amount)
	for i := range iter.N(amount) {
		devices = append(devices, &mockDevice{
			id:              fmt.Sprintf("%s_%2d", id.String(), i),
			topologyHintKey: topologyHintKey,
		})
	}

	return devices
}

// TestGetTopologyHintKeyUsingBestFitBinPacking tests whether `getTopologyHintKeyUsingBestFitBinPacking()` is picking the TopologyHintKey correctly or not.
func TestGetTopologyHintKeyUsingBestFitBinPacking(t *testing.T) {
	tests := []struct {
		description      string
		devicesByHintMap map[string]DeviceSet
		remainingCnt     int
		expectedIn       []string
	}{
		{
			description: "must pick '02' which is the target for best fit bin packing algorithm",
			devicesByHintMap: func() map[string]DeviceSet {
				deviceMap := make(map[string]DeviceSet)
				deviceMap["00"] = generateSameBoardMockDeviceSet(8, "00")
				deviceMap["01"] = generateSameBoardMockDeviceSet(7, "01")
				deviceMap["02"] = generateSameBoardMockDeviceSet(6, "02")
				deviceMap["03"] = generateSameBoardMockDeviceSet(4, "03")
				deviceMap["04"] = generateSameBoardMockDeviceSet(3, "04")

				return deviceMap
			}(),
			remainingCnt: 5,
			expectedIn:   []string{"02"},
		},
		{
			description: "must pick '04' which is the target for best fit bin packing algorithm",
			devicesByHintMap: func() map[string]DeviceSet {
				deviceMap := make(map[string]DeviceSet)
				deviceMap["00"] = generateSameBoardMockDeviceSet(8, "00")
				deviceMap["01"] = generateSameBoardMockDeviceSet(7, "01")
				deviceMap["02"] = generateSameBoardMockDeviceSet(6, "02")
				deviceMap["03"] = generateSameBoardMockDeviceSet(3, "03")
				deviceMap["04"] = generateSameBoardMockDeviceSet(5, "04")

				return deviceMap
			}(),
			remainingCnt: 5,
			expectedIn:   []string{"04"},
		},
		{
			description: "must not pick any key because none of them has sufficient size for requested remaining cnt",
			devicesByHintMap: func() map[string]DeviceSet {
				deviceMap := make(map[string]DeviceSet)
				deviceMap["00"] = generateSameBoardMockDeviceSet(5, "00")
				deviceMap["01"] = generateSameBoardMockDeviceSet(4, "01")
				deviceMap["02"] = generateSameBoardMockDeviceSet(3, "02")
				deviceMap["03"] = generateSameBoardMockDeviceSet(2, "03")
				deviceMap["04"] = generateSameBoardMockDeviceSet(1, "04")

				return deviceMap
			}(),
			remainingCnt: 7,
			expectedIn:   []string{""},
		},
	}

	t.Parallel()
	for _, tc := range tests {
		t.Run(tc.description, func(subT *testing.T) {
			actual := getTopologyHintKeyUsingBestFitBinPacking(tc.remainingCnt, &tc.devicesByHintMap)
			assert.Contains(subT, tc.expectedIn, actual)
		})
	}
}

// TestGetLargestLengthCandidatesTopologyHintKey tests whether `getLargestLengthCandidatesTopologyHintKey()` is picking the TopologyHintKey correctly or not.
func TestGetLargestLengthCandidatesTopologyHintKey(t *testing.T) {
	tests := []struct {
		description      string
		devicesByHintMap map[string]DeviceSet
		expectedIn       []string
	}{
		{
			description: "must pick '02' which is the target for best fit bin packing algorithm",
			devicesByHintMap: func() map[string]DeviceSet {
				deviceMap := make(map[string]DeviceSet)
				deviceMap["00"] = generateSameBoardMockDeviceSet(8, "00")
				deviceMap["01"] = generateSameBoardMockDeviceSet(7, "01")
				deviceMap["02"] = generateSameBoardMockDeviceSet(6, "02")
				deviceMap["03"] = generateSameBoardMockDeviceSet(4, "03")
				deviceMap["04"] = generateSameBoardMockDeviceSet(3, "04")

				return deviceMap
			}(),
			expectedIn: []string{"00"},
		},
		{
			description: "must pick '04' which is the target for best fit bin packing algorithm",
			devicesByHintMap: func() map[string]DeviceSet {
				deviceMap := make(map[string]DeviceSet)
				deviceMap["00"] = generateSameBoardMockDeviceSet(1, "00")
				deviceMap["01"] = generateSameBoardMockDeviceSet(5, "01")
				deviceMap["02"] = generateSameBoardMockDeviceSet(2, "02")
				deviceMap["03"] = generateSameBoardMockDeviceSet(4, "03")
				deviceMap["04"] = generateSameBoardMockDeviceSet(3, "04")

				return deviceMap
			}(),
			expectedIn: []string{"01"},
		},
		{
			description: "must not pick any key because none of them has sufficient size for requested remaining cnt",
			devicesByHintMap: func() map[string]DeviceSet {
				deviceMap := make(map[string]DeviceSet)
				deviceMap["00"] = generateSameBoardMockDeviceSet(3, "00")
				deviceMap["01"] = generateSameBoardMockDeviceSet(3, "01")
				deviceMap["02"] = generateSameBoardMockDeviceSet(4, "02")
				deviceMap["03"] = generateSameBoardMockDeviceSet(4, "03")
				deviceMap["04"] = generateSameBoardMockDeviceSet(4, "04")

				return deviceMap
			}(),
			expectedIn: []string{"02", "03", "04"},
		},
	}

	t.Parallel()
	for _, tc := range tests {
		t.Run(tc.description, func(subT *testing.T) {
			actual := getLargestLengthCandidatesTopologyHintKey(&tc.devicesByHintMap)
			assert.Contains(subT, tc.expectedIn, actual)
		})
	}
}

func TestAllocationOfBinPackingAllocator(t *testing.T) {
	mockDevices := []DeviceSet{
		generateSameBoardMockDeviceSet(1, "00"),
		generateSameBoardMockDeviceSet(2, "01"),
		generateSameBoardMockDeviceSet(3, "02"),
		generateSameBoardMockDeviceSet(4, "03"),
		generateSameBoardMockDeviceSet(5, "04"),
		generateSameBoardMockDeviceSet(6, "05"),
		generateSameBoardMockDeviceSet(7, "06"),
		generateSameBoardMockDeviceSet(8, "07"),
	}

	allocator, _ := NewBinPackingNpuAllocator(nil)

	tests := []struct {
		description string
		available   DeviceSet
		required    DeviceSet
		request     int
		expectedIn  DeviceSet
	}{
		{
			description: "pre-allocated devices already satisfies device request quantity. no additional allocation needed",
			available: func() DeviceSet {
				devices := make(DeviceSet, 0)
				devices = append(devices, mockDevices[0]...)
				devices = append(devices, mockDevices[1]...)
				devices = append(devices, mockDevices[2]...)
				devices = append(devices, mockDevices[3]...)
				devices = append(devices, mockDevices[4]...)
				devices = append(devices, mockDevices[5]...)
				devices = append(devices, mockDevices[6]...)
				devices = append(devices, mockDevices[7]...)

				return devices
			}(),
			required: func() DeviceSet {
				devices := make(DeviceSet, 0)
				devices = append(devices, mockDevices[0]...)
				devices = append(devices, mockDevices[2]...)

				return devices
			}(),
			request: 4,
			expectedIn: func() DeviceSet {
				devices := make(DeviceSet, 0)
				devices = append(devices, mockDevices[0]...)
				devices = append(devices, mockDevices[2]...)

				return devices
			}(),
		},
		{
			description: "must pick '01' which is the target for best fit bin packing algorithm",
			available: func() DeviceSet {
				devices := make(DeviceSet, 0)
				devices = append(devices, mockDevices[0]...)
				devices = append(devices, mockDevices[1]...)
				devices = append(devices, mockDevices[2]...)
				devices = append(devices, mockDevices[3]...)
				devices = append(devices, mockDevices[4]...)
				devices = append(devices, mockDevices[5]...)
				devices = append(devices, mockDevices[6]...)
				devices = append(devices, mockDevices[7]...)

				return devices
			}(),
			required: func() DeviceSet {
				devices := make(DeviceSet, 0)
				devices = append(devices, mockDevices[0]...)
				devices = append(devices, mockDevices[2]...)

				return devices
			}(),
			request: 5,
			expectedIn: func() DeviceSet {
				devices := make(DeviceSet, 0)
				devices = append(devices, mockDevices[0]...)
				devices = append(devices, mockDevices[2]...)
				devices = append(devices, mockDevices[1]...)

				return devices
			}(),
		},
		{
			description: "must pick '07' which is the target for best fit bin packing algorithm",
			available: func() DeviceSet {
				devices := make(DeviceSet, 0)
				devices = append(devices, mockDevices[0]...)
				devices = append(devices, mockDevices[1]...)
				devices = append(devices, mockDevices[2]...)
				devices = append(devices, mockDevices[3]...)
				devices = append(devices, mockDevices[4]...)
				devices = append(devices, mockDevices[5]...)
				devices = append(devices, mockDevices[6]...)
				devices = append(devices, mockDevices[7]...)

				return devices
			}(),
			required: func() DeviceSet {
				devices := make(DeviceSet, 0)
				devices = append(devices, mockDevices[0]...)
				devices = append(devices, mockDevices[2]...)

				return devices
			}(),
			request: 12,
			expectedIn: func() DeviceSet {
				devices := make(DeviceSet, 0)
				devices = append(devices, mockDevices[0]...)
				devices = append(devices, mockDevices[2]...)
				devices = append(devices, mockDevices[7]...)

				return devices
			}(),
		},
		{
			description: "must pick '03', '07' which is the target for best fit bin packing algorithm",
			available: func() DeviceSet {
				devices := make(DeviceSet, 0)
				devices = append(devices, mockDevices[0]...)
				devices = append(devices, mockDevices[1]...)
				devices = append(devices, mockDevices[2]...)
				devices = append(devices, mockDevices[3]...)
				devices = append(devices, mockDevices[4]...)
				devices = append(devices, mockDevices[5]...)
				devices = append(devices, mockDevices[6]...)
				devices = append(devices, mockDevices[7]...)

				return devices
			}(),
			required: func() DeviceSet {
				devices := make(DeviceSet, 0)
				devices = append(devices, mockDevices[0]...)
				devices = append(devices, mockDevices[2]...)

				return devices
			}(),
			request: 15,
			expectedIn: func() DeviceSet {
				devices := make(DeviceSet, 0)
				devices = append(devices, mockDevices[0]...)
				devices = append(devices, mockDevices[2]...)
				devices = append(devices, mockDevices[7]...)
				devices = append(devices, mockDevices[3]...)

				return devices
			}(),
		},
	}

	t.Parallel()
	for _, tc := range tests {
		t.Run(tc.description, func(subT *testing.T) {
			actual := allocator.Allocate(tc.available, tc.required, tc.request)
			for _, device := range actual {
				assert.Contains(subT, tc.expectedIn, device)
			}
		})
	}
}
