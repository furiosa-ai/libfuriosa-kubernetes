package device

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	deviceFilePattern  = `^npu(?P<device_id>\d+)((?:pe)(?P<start_core>\d+)(-(?P<end_core>\d+))?)?$`
	subExpKeyDeviceId  = "device_id"
	subExpKeyStartCore = "start_core"
	subExpKeyEndCore   = "end_core"
	defaultPlatform    = "FuriosaAI"
	nonDefaultPlatform = "VITIS"
)

var (
	deviceFileRegExp = regexp.MustCompile(deviceFilePattern)
)

type deviceValidateFunc func(dev DevFile) bool

type DevFile struct {
	fileAbsPath string
	fileType    fs.FileMode
}

// ListDevFs builds slice of DevFile from dev fs
func ListDevFs(devFs string) ([]DevFile, error) {
	devFsPath, err := filepath.Abs(devFs)
	if err != nil {
		return nil, err
	}
	files, err := os.ReadDir(devFs)
	if err != nil {
		return nil, err
	}

	var devFiles []DevFile
	for _, file := range files {
		devFiles = append(devFiles, DevFile{
			fileAbsPath: filepath.Join(devFsPath, file.Name()),
			fileType:    file.Type(),
		})
	}
	return devFiles, nil
}

// ParseIndices parse device id and core range from the given file.
func ParseIndices(filename string) (uint8, []uint8, error) {
	if !deviceFileRegExp.MatchString(filename) {
		return 0, nil, NewUnrecognizedFileError(filename)
	}

	matches := deviceFileRegExp.FindStringSubmatch(filename)
	subExps := deviceFileRegExp.SubexpNames()

	namedMatches := map[string]string{}
	for i, match := range matches {
		subExp := subExps[i]
		if subExp == "" {
			continue
		}
		namedMatches[subExp] = match
	}

	// parse device_id
	deviceIdStr := namedMatches[subExpKeyDeviceId]
	deviceId, err := strconv.ParseUint(deviceIdStr, 10, 8)
	if err != nil {
		return 0, nil, NewUnrecognizedFileError(filename)
	}

	// parse start_core
	coreStartStr := namedMatches[subExpKeyStartCore]
	coreStart, err := strconv.ParseUint(coreStartStr, 10, 8)
	if err != nil {
		return uint8(deviceId), nil, nil
	}

	// parse end_core
	endCoreStr := namedMatches[subExpKeyEndCore]
	endCore, err := strconv.ParseUint(endCoreStr, 10, 8)
	if err != nil {
		return uint8(deviceId), []uint8{uint8(coreStart)}, nil
	}

	var cores []uint8
	for i := coreStart; i <= endCore; i++ {
		cores = append(cores, uint8(i))
	}

	return uint8(deviceId), cores, nil
}

// defaultDeviceValidator is default validator function for DevFile
func defaultDeviceValidator(dev DevFile) bool {
	return dev.fileType&fs.ModeCharDevice != 0
}

func filterDevFiles(devFiles []DevFile, validator deviceValidateFunc) map[uint8][]string {
	filtered := map[uint8][]string{}

	for _, devFile := range devFiles {
		if validator(devFile) {

			//Base returns the last element of path
			deviceIdx, _, err := ParseIndices(filepath.Base(devFile.fileAbsPath))
			if err != nil {
				continue
			}

			// in rust, we canonicalize file path here, but in go, we always have abs path here.
			filtered[deviceIdx] = append(filtered[deviceIdx], devFile.fileAbsPath)
		}
	}

	return filtered
}

func isFuiosaPlatform(contents string) bool {
	examine := strings.TrimSpace(contents)

	return examine == defaultPlatform || examine == nonDefaultPlatform
}

func IsFuriosaDevice(idx uint8, sysFs string) bool {
	platformType, err := ReadMgmtFile(sysFs, PlatformType.Filename(), idx)
	if err != nil {
		return false
	}

	return isFuiosaPlatform(platformType)
}
