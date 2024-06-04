package main

import (
	"fmt"
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/manifest"
	"os"

	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/smi"
)

func main() {
	err := smi.Init()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	devices, err := smi.GetDevices()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	for _, device := range devices {
		rngdManifest, err := manifest.NewRngdManifest(device)
		if err != nil {
			return
		}

		nodes := rngdManifest.DeviceNodes()
		for _, node := range nodes {
			fmt.Printf("Node: %+v\n", node)
		}

		mounts := rngdManifest.MountPaths()
		for _, mount := range mounts {
			fmt.Printf("Mount: %+v\n", mount)
		}
	}

	_ = smi.Shutdown()
}
