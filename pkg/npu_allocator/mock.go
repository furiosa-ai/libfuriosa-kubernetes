package npu_allocator

import (
	"fmt"
	"strconv"

	"github.com/bradfitz/iter"
	"github.com/google/uuid"
)

var _ Device = (*mockDevice)(nil)

func NewMockDevice(index int, id string, topologyHintKey TopologyHintKey) Device {
	return &mockDevice{
		index:           index,
		id:              id,
		topologyHintKey: topologyHintKey,
	}
}

func buildMockDeviceSet(start, end int) DeviceSet {
	result := NewDeviceSet()
	for i := start; i <= end; i++ {
		result.Insert(buildMockDevice(i))
	}

	return result
}

func buildMockDevice(target int) Device {
	return &mockDevice{
		index:           target,
		id:              strconv.Itoa(target),
		topologyHintKey: TopologyHintKey(strconv.Itoa(target)),
	}
}

func getStaticHintKeys() []TopologyHintKey {
	return []TopologyHintKey{"0", "1", "2", "3", "4", "5", "6", "7"}
}

func buildStaticHintMatrixForTwoSocketBalancedConfig() TopologyHintMatrix {
	return TopologyHintMatrix{
		"0": {"0": 70, "1": 30, "2": 20, "3": 20, "4": 10, "5": 10, "6": 10, "7": 10},
		"1": {"1": 70, "2": 20, "3": 20, "4": 10, "5": 10, "6": 10, "7": 10},
		"2": {"2": 70, "3": 30, "4": 10, "5": 10, "6": 10, "7": 10},
		"3": {"3": 70, "4": 10, "5": 10, "6": 10, "7": 10},
		"4": {"4": 70, "5": 30, "6": 20, "7": 20},
		"5": {"5": 70, "6": 20, "7": 20},
		"6": {"6": 70, "7": 30},
		"7": {"7": 70},
	}
}

type mockDevice struct {
	index           int
	id              string
	topologyHintKey TopologyHintKey
}

func (m *mockDevice) Index() int {
	return m.index
}

func (m *mockDevice) ID() string {
	return m.id
}

func (m *mockDevice) TopologyHintKey() TopologyHintKey {
	return m.topologyHintKey
}

func (m *mockDevice) Equal(target Device) bool {
	if _, isMockDevicePtr := target.(*mockDevice); !isMockDevicePtr {
		return false
	}

	if m.id == target.ID() && m.topologyHintKey == target.TopologyHintKey() {
		return true
	}

	return false
}

func generateSameBoardMockDeviceSet(index int, cnt int, hintKey TopologyHintKey) DeviceSet {
	UUID, _ := uuid.NewUUID()

	devices := NewDeviceSet()
	for i := range iter.N(cnt) {
		devices.Insert(NewMockDevice(index, fmt.Sprintf("%s_%d", UUID.String(), i), hintKey))
	}

	return devices
}
