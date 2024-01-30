package main

import (
	"fmt"
	"os"

	furiosaDevice "github.com/furiosa-ai/libfuriosa-kubernetes/pkg/device"
	"github.com/jedib0t/go-pretty/v6/table"
)

func main() {
	devices, err := furiosaDevice.NewDeviceLister().ListDevices()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	topology, err := furiosaDevice.NewTopology(devices)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	header := table.Row{"#"}
	for _, device := range devices {
		header = append(header, device.Name())
	}
	t.AppendHeader(header)

	for _, device1 := range devices {
		row := table.Row{device1.Name()}
		key1, _ := device1.Busname()
		for _, device2 := range devices {
			key2, _ := device2.Busname()
			linkType := topology.GetLinkType(key1, key2)
			row = append(row, linkTypeToString(linkType))
		}
		t.AppendRow(row)
	}

	t.Render()
}

func linkTypeToString(linkType furiosaDevice.LinkType) string {
	switch linkType {
	case furiosaDevice.LinkTypeInterconnect:
		return "Interconnect"
	case furiosaDevice.LinkTypeCPU:
		return "CPU"
	case furiosaDevice.LinkTypeHostBridge:
		return "Host Bridge"
	case furiosaDevice.LinkTypeSoc:
		return "SoC"
	default:
		return "Unknown"
	}
}
