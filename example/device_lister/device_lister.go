package main

import (
	"fmt"
	"os"

	furiosaSmi "github.com/furiosa-ai/libfuriosa-kubernetes/pkg/furiosa_smi_go"
)

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -L/usr/local/lib -lfuriosa_smi
*/
import "C"

func main() {
	err := furiosaSmi.Init()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	devices, err := furiosaSmi.GetDevices()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("found %d\n", len(devices))

	for _, device := range devices {
		fmt.Printf("%v\n", device)
		coreStatus, err := device.CoreStatus()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Printf("%v\n", coreStatus)

		deviceFiles, err := device.DeviceFiles()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		for _, deviceFile := range deviceFiles {
			fmt.Printf("%v\n", deviceFile)
		}
	}

	_ = furiosaSmi.Shutdown()
}
