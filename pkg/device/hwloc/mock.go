package hwloc

type mockTopologyCtx struct {
	origin          topologyCtx
	xmlTopologyPath string
}

var _ Hwloc = new(mockTopologyCtx)

func NewMockHwloc(xmlTopologyPath string) Hwloc {
	return &mockTopologyCtx{
		origin: topologyCtx{
			topologyContext: nil,
		},
		xmlTopologyPath: xmlTopologyPath,
	}
}

func (m *mockTopologyCtx) TopologyInit() error {
	return m.origin.TopologyInit()
}

func (m *mockTopologyCtx) SetIoTypeFilter() error {
	return m.origin.SetIoTypeFilter()
}

func (m *mockTopologyCtx) TopologyLoad() error {
	// import topology from designated xml file
	if err := topologySetXML(m.origin.topologyContext, m.xmlTopologyPath); err != nil {
		return err
	}
	return m.origin.TopologyLoad()
}

func (m *mockTopologyCtx) GetCommonAncestorObjType(dev1BDF string, dev2BDF string) (HwlocObjType, error) {
	return m.origin.GetCommonAncestorObjType(dev1BDF, dev2BDF)
}

func (m *mockTopologyCtx) TopologyDestroy() {
	m.origin.TopologyDestroy()
}
