package device

import "path/filepath"

const (
	SepColon = ":"
)

func Abs(input string) string {
	abs, _ := filepath.Abs(input)
	return abs
}
