package main

import (
	"fmt"
	"os"

	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
)

func main() {
	devices, err := smi.ListDevices()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("found %d device(s)\n", len(devices))

	for _, device := range devices {
		deviceInfo, err := device.DeviceInfo()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Printf("Device Arch: %v\n", deviceInfo.Arch())
		fmt.Printf("Device CoreNum: %d\n", deviceInfo.CoreNum())
		fmt.Printf("Device NumaNode: %d\n", deviceInfo.NumaNode())
		fmt.Printf("Device Name: %s\n", deviceInfo.Name())
		fmt.Printf("Device Serial: %s\n", deviceInfo.Serial())
		fmt.Printf("Device UUID: %s\n", deviceInfo.UUID())
		fmt.Printf("Device BDF: %s\n", deviceInfo.BDF())
		fmt.Printf("Device Major: %d\n", deviceInfo.Major())
		fmt.Printf("Device Minor: %d\n", deviceInfo.Minor())
		fmt.Printf("Device FirmwareVersion")
		fmt.Printf("  Arch: %v\n", deviceInfo.FirmwareVersion().Arch())
		fmt.Printf("  Major: %d\n", deviceInfo.FirmwareVersion().Major())
		fmt.Printf("  Minor: %d\n", deviceInfo.FirmwareVersion().Minor())
		fmt.Printf("  Patch: %d\n", deviceInfo.FirmwareVersion().Patch())
		fmt.Printf("  Meta: %s\n", deviceInfo.FirmwareVersion().Metadata())
		fmt.Printf("Device DriverVersion")
		fmt.Printf("  Arch: %v\n", deviceInfo.DriverVersion().Arch())
		fmt.Printf("  Major: %d\n", deviceInfo.DriverVersion().Major())
		fmt.Printf("  Minor: %d\n", deviceInfo.DriverVersion().Minor())
		fmt.Printf("  Patch: %d\n", deviceInfo.DriverVersion().Patch())
		fmt.Printf("  Meta: %s\n", deviceInfo.DriverVersion().Metadata())

		liveness, err := device.Liveness()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Printf("Liveness: %v\n", liveness)

		coreStatus, err := device.CoreStatus()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Printf("Core Status:\n")
		for core, status := range coreStatus {
			fmt.Printf("  Core %d: %v\n", core, status)
		}

		deviceFiles, err := device.DeviceFiles()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Printf("Device Files:\n")
		for _, deviceFile := range deviceFiles {
			fmt.Printf("  Cores: %v\n", deviceFile.Cores())
			fmt.Printf("  Path: %s\n", deviceFile.Path())
		}

		//print DeviceErrorInfo nicely
		deviceErrorInfo, err := device.DeviceErrorInfo()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Printf("Device Error Info:\n")
		fmt.Printf("  AxiPostErrorCount: %d\n", deviceErrorInfo.AxiPostErrorCount())
		fmt.Printf("  AxiFetchErrorCount: %d\n", deviceErrorInfo.AxiFetchErrorCount())
		fmt.Printf("  AxiDiscardErrorCount: %d\n", deviceErrorInfo.AxiDiscardErrorCount())
		fmt.Printf("  AxiDoorbellErrorCount: %d\n", deviceErrorInfo.AxiDoorbellErrorCount())
		fmt.Printf("  PciePostErrorCount: %d\n", deviceErrorInfo.PciePostErrorCount())
		fmt.Printf("  PcieFetchErrorCount: %d\n", deviceErrorInfo.PcieFetchErrorCount())
		fmt.Printf("  PcieDiscardErrorCount: %d\n", deviceErrorInfo.PcieDiscardErrorCount())
		fmt.Printf("  PcieDoorbellErrorCount: %d\n", deviceErrorInfo.PcieDoorbellErrorCount())
		fmt.Printf("  DeviceErrorCount: %d\n", deviceErrorInfo.DeviceErrorCount())

		//printDeviceUtilization
		utilization, err := device.DeviceUtilization()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Printf("Device Utilization:\n")
		for _, peUtilization := range utilization.PeUtilization() {
			fmt.Printf("  PE Utilization:\n")
			fmt.Printf("    Cores: %v\n", peUtilization.Cores())
			fmt.Printf("    Time Window Mill: %d\n", peUtilization.TimeWindowMill())
			fmt.Printf("    PE Usage Percentage: %f\n", peUtilization.PeUsagePercentage())
		}
		fmt.Printf("  Memory Utilization:\n")
		fmt.Printf("    Total Bytes: %d\n", utilization.MemoryUtilization().TotalBytes())
		fmt.Printf("    In Use Bytes: %d\n", utilization.MemoryUtilization().InUseBytes())

		temperature, err := device.DeviceTemperature()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Printf("Device Temperature:\n")
		fmt.Printf("  Soc Peak: %f\n", temperature.SocPeak())
		fmt.Printf("  Ambient: %f\n", temperature.Ambient())

		powerConsumption, err := device.PowerConsumption()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Printf("Power Consumption: %f\n", powerConsumption)

	}
}
