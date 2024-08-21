package npu_allocator

// scoreDeviceSet calculates total score of given devices using given hintProvider.
func scoreDeviceSet(hintProvider TopologyHintProvider, devices DeviceSet) uint {
	scoreDevicePair := func(device1 Device, device2 Device) uint {
		return hintProvider(device1, device2)
	}

	var total uint = 0
	for i := 0; i < len(devices); i++ {
		for j := i + 1; j < len(devices); j++ {
			total += scoreDevicePair(devices[i], devices[j])
		}
	}

	return total
}
