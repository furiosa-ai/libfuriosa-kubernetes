package npu_allocator

var _ Device = (*mockDevice)(nil)

func NewMockDevice(index int, id string, topologyHintKey TopologyHintKey) Device {
	return &mockDevice{
		index:           index,
		id:              id,
		topologyHintKey: topologyHintKey,
	}
}

type mockDevice struct {
	index           int
	id              string
	topologyHintKey TopologyHintKey
}

func (m *mockDevice) GetIndex() int {
	return m.index
}

func (m *mockDevice) GetID() string {
	return m.id
}

func (m *mockDevice) GetTopologyHintKey() TopologyHintKey {
	return m.topologyHintKey
}

func (m *mockDevice) Equal(target Device) bool {
	if _, isMockDevicePtr := target.(*mockDevice); !isMockDevicePtr {
		return false
	}

	if m.id == target.GetID() && m.topologyHintKey == target.GetTopologyHintKey() {
		return true
	}

	return false
}
