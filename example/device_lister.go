package main

import (
	"fmt"
	"os"

	furiosaDevice "github.com/furiosa-ai/libfuriosa-kubernetes/pkg/device"
)

func main() {
	devices, err := furiosaDevice.NewDeviceLister().ListDevices()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	println("found ", len(devices))

	for _, device := range devices {
		println(fmt.Sprintf("%v", device))
		for _, deviceFile := range device.DevFiles() {
			println(fmt.Sprintf("%v", deviceFile))
		}
	}
}
