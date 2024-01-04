package device

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

const (
	NpuExp = "npu%d"
)

// CoreStatus Enum for NPU device status.
type DeviceStatus string

const (
	DeviceStatusAvailable DeviceStatus = "DeviceStatusAvailable"
	DeviceStatusOccupied  DeviceStatus = "DeviceStatusOccupied"
)

// CoreStatus Enum for NPU core status.
type CoreStatus string

const (
	CoreStatusAvailable   CoreStatus = "CoreStatusAvailable"
	CoreStatusOccupied    CoreStatus = "CoreStatusOccupied"
	CoreStatusUnavailable CoreStatus = "CoreStatusUnavailable"
)

// DeviceMode Enum for NPU's operating mode.
type DeviceMode string

const (
	DeviceModeSingle    DeviceMode = "DeviceModeSingle"
	DeviceModeFusion    DeviceMode = "DeviceModeFusion"
	DeviceModeMultiCore DeviceMode = "DeviceModeMultiCore"
)

type CoreRangeType string

const (
	CoreRangeTypeAll   CoreRangeType = "CoreRangeTypeAll"
	CoreRangeTypeRange CoreRangeType = "CoreRangeTypeRange"
)

type CoreRange interface {
	Type() CoreRangeType
	Start() uint8
	End() uint8
	Contains(idx uint8) bool
	//TODO has_intersection
}

var _ CoreRange = new(coreRange)

type coreRange struct {
	coreRangeType CoreRangeType
	start         uint8
	end           uint8
}

func (c coreRange) Type() CoreRangeType {
	return c.coreRangeType
}

func (c coreRange) Start() uint8 {
	return c.start
}

func (c coreRange) End() uint8 {
	return c.end
}

func (c coreRange) Contains(idx uint8) bool {
	if c.coreRangeType == CoreRangeTypeAll {
		return true
	}

	return idx >= c.start && idx <= c.end
}

type Device interface {
	// Name returns the name of the device (e.g., npu0).
	Name() string
	// DeviceIndex returns the device index (e.g., 1 for npu1pe0).
	DeviceIndex() uint8
	// Arch returns `Arch` of the device(e.g., `Warboy`).
	Arch() Arch
	// Alive returns a liveness state of the device.
	Alive() (bool, error)
	// AtrErr returns error states of the device.
	AtrErr() (map[string]uint32, error)
	// Busname returns PCI bus number of the device.
	Busname() (string, error)
	// PCIDev returns PCI device ID of the device.
	PCIDev() (string, error)
	// DeviceSn returns serial number of the device.
	DeviceSn() (string, error)
	// DeviceUUID returns UUID of the device.
	DeviceUUID() (string, error)
	// FirmwareVersion retrieves firmware revision from the device.
	FirmwareVersion() (string, error)
	// DriverVersion retrieves driver version for the device.
	DriverVersion() (string, error)
	// HeartBeat returns uptime of the device.
	HeartBeat() (uint32, error)
	// NumaNode retrieve NUMA node ID associated with the NPU's PCI lane
	NumaNode() (uint8, error)
	// CoreNum counts the number of cores.
	CoreNum() uint8
	// Cores list the core indices.
	Cores() []uint8
	// DevFiles list device files under this device.
	DevFiles() []DeviceFile
	// GetStatusCore examine a specific core of the device, whether it is available or not.
	GetStatusCore(core uint8) (CoreStatus, error)
	// GetStatusAll Examine each core of the device, whether it is available or not.
	GetStatusAll() (map[uint8]CoreStatus, error)
	//TODO: performance_counter
	//TODO: clock_frequency
	//TODO: get_hwmon_fetcher
}

var _ Device = new(device)

type device struct {
	deviceIndex uint8
	devRoot     string
	sysRoot     string
	arch        Arch
	meta        map[string]string
	numaNode    int
	//TODO: hwmon_fetcher
	cores    []uint8
	devFiles []DeviceFile
}

func NewDevice(deviceIndex uint8, paths []string, devFs string, sysFs string) (Device, error) {
	if !IsFuriosaDevice(deviceIndex, sysFs) {
		return nil, NewDeviceNotFoundError(fmt.Sprintf(NpuExp, deviceIndex))
	}

	meta := map[string]string{}
	for _, file := range ListStaticMgmtFiles() {
		value, err := ReadMgmtFile(sysFs, file.Filename(), deviceIndex)
		if err != nil {
			return nil, err
		}

		meta[file.Filename()] = value
	}

	//build device arch
	deviceType := meta[DeviceType.Filename()]
	socRev := meta[SocRev.Filename()]
	arch, err := archFromStr(deviceType, socRev)
	if err != nil {
		return nil, NewUnknownArchError(deviceType, socRev, err.Error())
	}

	//parse numa node
	busname := meta[Busname.Filename()]
	numaNode, err := ReadNumaNode(sysFs, busname)
	if err != nil {
		return nil, NewUnexpectedValue(err.Error())
	}

	// do collect_devices here
	cores, devFiles := collectDevices(paths)

	// TODO(@bg): we don't need hwmon fetcher support at this moment: Nov, 2023.
	return &device{
		deviceIndex: deviceIndex,
		devRoot:     devFs,
		sysRoot:     sysFs,
		arch:        arch,
		meta:        meta,
		numaNode:    numaNode,
		cores:       cores,
		devFiles:    devFiles,
	}, nil
}

func collectDevices(paths []string) ([]uint8, []DeviceFile) {
	coreMap := map[uint8]uint8{}
	cores := []uint8{}
	devFiles := []DeviceFile{}

	for _, path := range paths {
		devFile, err := NewDeviceFile(path)
		if err != nil {
			//NOTE(@bg): in rust impl, we have assumption here that there will be no error.
			continue
		}
		devFiles = append(devFiles, devFile)

		_, coreIndices, err := ParseIndices(filepath.Base(path))
		if err != nil {
			continue
		}

		for _, coreIdx := range coreIndices {
			coreMap[coreIdx] = coreIdx
		}
	}

	for coreIdx := range coreMap {
		cores = append(cores, coreIdx)
	}

	sort.SliceStable(devFiles, func(i, j int) bool {
		if devFiles[i].DeviceIndex() != devFiles[j].DeviceIndex() {
			return devFiles[i].DeviceIndex() < devFiles[j].DeviceIndex()
		}

		typeI, typeJ := devFiles[i].CoreRange().Type(), devFiles[j].CoreRange().Type()

		if typeI != typeJ {
			return typeI == CoreRangeTypeAll
		}

		if typeI == CoreRangeTypeRange {
			effectiveRangeI := devFiles[i].CoreRange().End() - devFiles[i].CoreRange().Start()
			effectiveRangeJ := devFiles[j].CoreRange().End() - devFiles[j].CoreRange().Start()
			return effectiveRangeI < effectiveRangeJ
		}

		return false
	})

	return cores, devFiles
}

func (d *device) Name() string {
	return fmt.Sprintf(NpuExp, d.deviceIndex)
}

func (d *device) DeviceIndex() uint8 {
	return d.deviceIndex
}

func (d *device) Arch() Arch {
	return d.arch
}

func (d *device) Alive() (bool, error) {
	aliveStr, err := ReadMgmtFile(d.sysRoot, Alive.Filename(), d.deviceIndex)
	if err != nil {
		return false, NewUnexpectedValue(err.Error())
	}

	alive, err := strconv.ParseBool(aliveStr)
	if err != nil {
		return false, NewUnexpectedValue(err.Error())
	}

	return alive, nil
}

func (d *device) AtrErr() (map[string]uint32, error) {
	atrErrMap := map[string]uint32{}
	atrErrStr, err := ReadMgmtFile(d.sysRoot, AtrError.Filename(), d.deviceIndex)
	if err != nil {
		return nil, NewUnexpectedValue(err.Error())
	}

	atrErrStr = strings.TrimSpace(atrErrStr)
	scanner := bufio.NewScanner(strings.NewReader(atrErrStr))
	for scanner.Scan() {
		if rawKey, rawValue, success := strings.Cut(strings.TrimSpace(scanner.Text()), SepColon); success {
			key := strings.ToLower(strings.TrimSpace(rawKey))
			value, parseErr := strconv.ParseUint(rawValue, 10, 32)
			if parseErr != nil {
				continue
			}

			atrErrMap[key] = uint32(value)
		}
	}

	return atrErrMap, nil
}

// BusName returns PCI bus number of the device.
func (d *device) Busname() (ret string, err error) {
	ret = d.meta[Busname.Filename()]
	if ret == "" {
		err = NewUnexpectedValue("could not retrieve busname")
	}

	return
}

// PCIDev returns PCI device ID of the device.
func (d *device) PCIDev() (ret string, err error) {
	ret = d.meta[Dev.Filename()]
	if ret == "" {
		err = NewUnexpectedValue("could not retrieve pci device id")
	}

	return
}

// DeviceSn returns serial number of the device.
func (d *device) DeviceSn() (ret string, err error) {
	ret = d.meta[DeviceSn.Filename()]
	if ret == "" {
		err = NewUnexpectedValue("could not retrieve device sn")
	}

	return
}

// DeviceUUID returns UUID of the device
func (d *device) DeviceUUID() (ret string, err error) {
	ret = d.meta[DeviceUuid.Filename()]
	if ret == "" {
		err = NewUnexpectedValue("could not retrieve device uuid")
	}

	return
}

func (d *device) FirmwareVersion() (ret string, err error) {
	firmwareVersion, err := ReadMgmtFile(d.sysRoot, FwVersion.Filename(), d.deviceIndex)
	if err != nil {
		return "", NewUnexpectedValue(err.Error())
	}

	if firmwareVersion == "" {
		return "", NewUnexpectedValue("could not retrieve firmware version")
	}

	return firmwareVersion, nil
}

func (d *device) DriverVersion() (string, error) {
	driverVersion, err := ReadMgmtFile(d.sysRoot, Version.Filename(), d.deviceIndex)
	if err != nil {
		return "", NewUnexpectedValue(err.Error())
	}

	if driverVersion == "" {
		return "", NewUnexpectedValue("could not retrieve driver version")
	}

	return driverVersion, nil
}

func (d *device) HeartBeat() (uint32, error) {
	heartBeatStr, err := ReadMgmtFile(d.sysRoot, Heartbeat.Filename(), d.deviceIndex)
	if err != nil {
		return 0, NewUnexpectedValue(err.Error())
	}
	heartBeat, err := strconv.ParseInt(heartBeatStr, 10, 32)
	if err != nil {
		return 0, NewUnexpectedValue(err.Error())
	}

	return uint32(heartBeat), nil
}

func (d *device) NumaNode() (ret uint8, err error) {
	ret = uint8(d.numaNode)
	if d.numaNode < 0 {
		ret = 0
		err = NewUnexpectedValue(strconv.Itoa(d.numaNode))
	}

	return
}

func (d *device) CoreNum() uint8 {
	return uint8(len(d.cores))
}

func (d *device) Cores() []uint8 {
	return d.cores
}

func (d *device) DevFiles() []DeviceFile {
	return d.devFiles
}

func getDeviceStatus(filePath string) (DeviceStatus, error) {
	//NOTE(@bg): the perm argument is not used unless a file is being created.
	file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		if errors.Is(err, syscall.EBUSY) {
			return DeviceStatusOccupied, nil
		} else {
			return "", NewUnexpectedValue(err.Error())
		}
	}
	defer file.Close()
	return DeviceStatusAvailable, nil
}

// GetStatusCore examine a specific core of the device, whether it is available or not.
func (d *device) GetStatusCore(core uint8) (CoreStatus, error) {
	for _, file := range d.devFiles {
		if file.Mode() != DeviceModeSingle {
			continue
		}

		if file.CoreRange().Contains(core) {
			//Note(@bg): error is ignored here for two reasons.
			//1) file is already validated when a deviceFile was created.
			//2) to stay aligned with rust impl.
			status, _ := getDeviceStatus(file.Path())
			if status == DeviceStatusOccupied {
				return CoreStatusOccupied, nil
			}
		}
	}
	return CoreStatusAvailable, nil
}

// GetStatusAll examine each core of the device, whether it is available or not.
func (d *device) GetStatusAll() (map[uint8]CoreStatus, error) {
	statusMap := map[uint8]CoreStatus{}

	for _, core := range d.cores {
		//Note(@bg): error is never returned as same as GetStatusCore
		statusMap[core], _ = d.GetStatusCore(core)
	}

	return statusMap, nil
}
