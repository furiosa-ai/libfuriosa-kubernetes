package device // Arch Enum for the NPU architecture.

import "errors"

type Arch string

const (
	ArchWarboy Arch = "Warboy"
	ArchRngd   Arch = "Rngd"
)

func archFromStr(archStr string, revStr string) (arch Arch, error error) {
	switch archStr + revStr {
	case "WarboyB0":
		arch = ArchWarboy
		error = nil
	case "Rngd":
		arch = ArchRngd
		error = nil
	default:
		arch = ""
		error = errors.New("unknown arch")
	}

	return
}
