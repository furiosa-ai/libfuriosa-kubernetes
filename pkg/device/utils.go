package device

import "path/filepath"

const (
	SepColon = ":"
)

func safeDerefUint8(ptr *uint8) uint8 {
	if ptr == nil {
		return 0
	}

	return *ptr
}

func Abs(input string) string {
	abs, _ := filepath.Abs(input)
	return abs
}
