package npu_allocator

var _ NpuAllocator = (*binPackingNpuAllocator)(nil)

type binPackingNpuAllocator struct{}

func NewBinPackingNpuAllocator() (NpuAllocator, error) {
	//TODO implement me
	panic("implement me")
}

func (b binPackingNpuAllocator) Allocate(available DeviceSet, required DeviceSet, size int) DeviceSet {
	//TODO implement me
	panic("implement me")
}
