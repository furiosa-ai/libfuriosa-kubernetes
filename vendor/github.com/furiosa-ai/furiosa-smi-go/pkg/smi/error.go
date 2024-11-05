package smi

import (
	"errors"

	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi/binding"
)

func ToError(code binding.FuriosaSmiReturnCode) error {
	switch code {
	case binding.FuriosaSmiReturnCodeOk:
		return nil

	case binding.FuriosaSmiReturnCodeInvalidArgumentError:
		return errors.New("invalid argument error")

	case binding.FuriosaSmiReturnCodeNullPointerError:
		return errors.New("null pointer error")

	case binding.FuriosaSmiReturnCodeMaxBufferSizeExceedError:
		return errors.New("max buffer size exceed error")

	case binding.FuriosaSmiReturnCodeDeviceNotFoundError:
		return errors.New("device not found error")

	case binding.FuriosaSmiReturnCodeDeviceBusyError:
		return errors.New("device busy error")

	case binding.FuriosaSmiReturnCodeIoError:
		return errors.New("io error")

	case binding.FuriosaSmiReturnCodePermissionDeniedError:
		return errors.New("permission denied error")

	case binding.FuriosaSmiReturnCodeUnknownArchError:
		return errors.New("unknown arch error")

	case binding.FuriosaSmiReturnCodeIncompatibleDriverError:
		return errors.New("incompatible driver error")

	case binding.FuriosaSmiReturnCodeUnexpectedValueError:
		return errors.New("unexpected value error")

	case binding.FuriosaSmiReturnCodeParseError:
		return errors.New("parse error")

	case binding.FuriosaSmiReturnCodeUnknownError:
		return errors.New("unknown error")

	case binding.FuriosaSmiReturnCodeInternalError:
		return errors.New("internal error")

	case binding.FuriosaSmiReturnCodeUninitializedError:
		return errors.New("uninitialized error")

	case binding.FuriosaSmiReturnCodeContextError:
		return errors.New("context error")

	case binding.FuriosaSmiReturnCodeNotSupportedError:
		return errors.New("not supported error")
	}

	return nil
}
