package main

import (
	"fmt"
	"os"

	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/cdi_spec_gen"
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/furiosa_device"
)

func main() {
	err := smi.Init()
	if err != nil {
		return
	}

	devices, err := smi.ListDevices()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	furiosaDevices, err := furiosa_device.NewFuriosaDevices(devices, nil, furiosa_device.NonePolicy)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	medianIdx := len(furiosaDevices) / 2

	cdiSpec := cdi_spec_gen.NewSpec(cdi_spec_gen.WithDevices(furiosaDevices...),
		cdi_spec_gen.WithAggregatedDevice(),
		cdi_spec_gen.WithSpecDirs(cdi_spec_gen.DefaultStaticDir),
		cdi_spec_gen.WithSpecFileName("furiosa.yaml"),
		cdi_spec_gen.WithFilePermissions(cdi_spec_gen.DefaultPermissions),
		cdi_spec_gen.WithGroupDevice("group1", false, furiosaDevices...),
		cdi_spec_gen.WithGroupDevice("group2", true, furiosaDevices...),
		cdi_spec_gen.WithGroupDevice("group3", false, furiosaDevices[medianIdx:]...),
		cdi_spec_gen.WithGroupDevice("group4", true, furiosaDevices[medianIdx:]...),
	)

	err = cdiSpec.Write()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
}
