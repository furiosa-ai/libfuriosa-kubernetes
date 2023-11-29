package device

import (
	"errors"
	"fmt"
)

var (
	DeviceNotfound     = errors.New("device not found")
	UnknownArch        = errors.New("unknown arch")
	IncompatibleDriver = errors.New("incompatible device driver")
	UnexpectedValue    = errors.New("unexpected value")
	//TODO: DeviceBusy
	//TODO: IoErro
	//TODO: PermissionDenied
	//TODO: HwmonError
	//TODO: PerformanceCounterError
	//TODO: ParseError
)

func NewUnrecognizedFileError(file string) error {
	return fmt.Errorf("%w : "+file+" cannot be recognized", IncompatibleDriver)
}

func NewUnknownArchError(arch string, rev string, msg string) error {
	return fmt.Errorf("%w : arch: %s, rev: %s, msg: %s ", UnknownArch, arch, rev, msg)
}

func NewDeviceNotFoundError(device string) error {
	return fmt.Errorf("%w : device %s not found", DeviceNotfound, device)
}

func NewUnexpectedValue(msg string) error {
	return fmt.Errorf("%w : msg: %s", UnexpectedValue, msg)
}
