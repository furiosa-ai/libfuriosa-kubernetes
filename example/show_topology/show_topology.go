package main

import (
	"fmt"
	"os"
	"path/filepath"

	furiosaSmi "github.com/furiosa-ai/libfuriosa-kubernetes/pkg/furiosa_smi_go"
	"github.com/jedib0t/go-pretty/v6/table"
)

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

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	header := table.Row{"#"}
	for _, device := range devices {
		info, err := device.DeviceInfo()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}

		header = append(header, filepath.Base(info.Name()))
	}
	t.AppendHeader(header)

	for _, device1 := range devices {
		info1, err := device1.DeviceInfo()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}

		row := table.Row{filepath.Base(info1.Name())}
		for _, device2 := range devices {
			linkType, err := device1.GetDeviceToDeviceLinkType(device2)
			if err != nil {
				fmt.Printf("%s\n", err.Error())
				os.Exit(1)
			}

			row = append(row, linkTypeToString(linkType))
		}
		t.AppendRow(row)
	}

	t.Render()

	_ = furiosaSmi.Shutdown()
}

func linkTypeToString(linkType furiosaSmi.LinkType) string {
	switch linkType {
	case furiosaSmi.LinkTypeInterconnect:
		return "Interconnect"
	case furiosaSmi.LinkTypeCpu:
		return "CPU"
	case furiosaSmi.LinkTypeHostBridge:
		return "Host Bridge"
	case furiosaSmi.LinkTypeNoc:
		return "NoC"
	default:
		return "Unknown"
	}
}
