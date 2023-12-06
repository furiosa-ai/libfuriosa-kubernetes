package main

import (
	"fmt"
	"os"

	furiosaDevice "github.com/furiosa-ai/libfuriosa-kubernetes/pkg/device"
)

func main() {
	devices, err := furiosaDevice.NewDeviceLister().ListDevices()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("found %d\n", len(devices))

	for _, device := range devices {
		fmt.Printf("%v\n", device)
		coreStatus, err := device.GetStatusAll()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Printf("%v\n", coreStatus)
		for _, deviceFile := range device.DevFiles() {
			fmt.Printf("%v\n", deviceFile)
		}
	}
}
