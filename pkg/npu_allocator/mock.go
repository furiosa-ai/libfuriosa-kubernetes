package npu_allocator

var _ Device = (*mockDevice)(nil)

func NewMockDevice(id string, topologyHintKey TopologyHintKey) Device {
	return &mockDevice{
		id:              id,
		topologyHintKey: topologyHintKey,
	}
}

type mockDevice struct {
	id              string
	topologyHintKey TopologyHintKey
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
