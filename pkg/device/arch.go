package device // Arch Enum for the NPU architecture.

import "errors"

type Arch string

const (
	ArchWarboy   Arch = "Warboy"
	ArchRenegade Arch = "Renegade"
)

func archFromStr(archStr string, revStr string) (arch Arch, error error) {
	switch archStr + revStr {
	case "WarboyB0":
		arch = ArchWarboy
		error = nil
	case "Renegade":
		arch = ArchRenegade
		error = nil
	default:
		arch = ""
		error = errors.New("unknown arch")
	}

	return
}
