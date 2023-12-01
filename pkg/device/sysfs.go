package device

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	NumaNodeFilePath = "%s/bus/pci/devices/%s/numa_node"
	MgmtFilePath     = "class/npu_mgmt/npu%d_mgmt/%s"
)

type MgmtFile interface {
	Filename() string
	IsStatic() bool
}

type StaticMgmtFile string

const (
	Busname      StaticMgmtFile = "busname"
	Dev          StaticMgmtFile = "dev"
	DeviceSn     StaticMgmtFile = "device_sn"
	DeviceType   StaticMgmtFile = "device_type"
	DeviceUuid   StaticMgmtFile = "device_uuid"
	PlatformType StaticMgmtFile = "platform_type"
	SocRev       StaticMgmtFile = "soc_rev"
	SocUid       StaticMgmtFile = "soc_uid"
)

func (s StaticMgmtFile) Filename() string {
	//@bg We don't need more complexity here, since StaticMgmtFile is custom type of string.
	return string(s)
}

func (s StaticMgmtFile) IsStatic() bool {
	return true
}

func ListStaticMgmtFiles() []MgmtFile {
	return []MgmtFile{Busname, Dev, DeviceSn, DeviceType, DeviceUuid, PlatformType, SocRev, SocUid}
}

var _ MgmtFile = new(StaticMgmtFile)

type DynamicMgmtFile string

const (
	Alive         DynamicMgmtFile = "alive"
	AtrError      DynamicMgmtFile = "atr_error"
	FwVersion     DynamicMgmtFile = "fw_version"
	Heartbeat     DynamicMgmtFile = "heartbeat"
	NeClkFreqInfo DynamicMgmtFile = "ne_clk_freq_info"
	Version       DynamicMgmtFile = "version"
)

var _ MgmtFile = new(DynamicMgmtFile)

func (d DynamicMgmtFile) Filename() string {
	// We don't need more complexity here since DynamicMgmtFile is custom type of string.
	return string(d)
}

func (d DynamicMgmtFile) IsStatic() bool {
	return false
}

func ReadMgmtFile(sysFs string, mgmtFile string, deviceIndex uint8) (string, error) {
	sysFsPath, err := filepath.Abs(sysFs)
	if err != nil {
		return "", err
	}

	mgmtFileAbsPath := BuildMgmtFilePath(sysFsPath, mgmtFile, deviceIndex)
	bytes, err := os.ReadFile(mgmtFileAbsPath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(bytes)), nil
}

func BuildMgmtFilePath(baseDir string, file string, idx uint8) string {
	return filepath.Join(baseDir, fmt.Sprintf(MgmtFilePath, idx, file))
}

func ReadNumaNode(sysFs string, bdfIdentifier string) (int, error) {
	// BDF(bus, device, function) identifier is unique id of device.
	numaNodeFilePath := fmt.Sprintf(NumaNodeFilePath, sysFs, bdfIdentifier)
	bytes, err := os.ReadFile(numaNodeFilePath)
	if err != nil {
		return -1, err
	}

	numaNode, err := strconv.Atoi(strings.TrimSpace(string(bytes)))
	if err != nil {
		return -1, err
	}

	return numaNode, nil
}
