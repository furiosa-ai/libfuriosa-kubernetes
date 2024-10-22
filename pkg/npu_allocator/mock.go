package npu_allocator

import "strconv"

var _ Device = (*mockDevice)(nil)

func NewMockDevice(index int, id string, topologyHintKey TopologyHintKey) Device {
	return &mockDevice{
		index:           index,
		id:              id,
		topologyHintKey: topologyHintKey,
	}
}

func buildMockDeviceSet(start, end int) DeviceSet {
	result := DeviceSet{}
	for i := start; i <= end; i++ {
		result = append(result, buildMockDevice(i))
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
