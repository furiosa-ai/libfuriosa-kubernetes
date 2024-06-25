package smi

import (
	"errors"

	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/smi/binding"
)

func ToError(code binding.FuriosaSmiReturnCode) (ret error) {
	switch code {
	case binding.FuriosaSmiReturnCodeOk:
		ret = nil
	case binding.FuriosaSmiReturnCodeInvalidArgumentError:
		ret = errors.New("invalid argument error")
	case binding.FuriosaSmiReturnCodeNullPointerError:
		ret = errors.New("null pointer error")
	case binding.FuriosaSmiReturnCodeMaxBufferSizeExceedError:
		ret = errors.New("max buffer size exceed error")
	case binding.FuriosaSmiReturnCodeDeviceNotFoundError:
		ret = errors.New("device not found error")
	case binding.FuriosaSmiReturnCodeDeviceBusyError:
		ret = errors.New("device busy error")
	case binding.FuriosaSmiReturnCodeIoError:
		ret = errors.New("io error")
	case binding.FuriosaSmiReturnCodePermissionDeniedError:
		ret = errors.New("permission denied error")
	case binding.FuriosaSmiReturnCodeUnknownArchError:
		ret = errors.New("unknown arch error")
	case binding.FuriosaSmiReturnCodeIncompatibleDriverError:
		ret = errors.New("incompatible driver error")
	case binding.FuriosaSmiReturnCodeUnexpectedValueError:
		ret = errors.New("unexpected value error")
	case binding.FuriosaSmiReturnCodeParseError:
		ret = errors.New("parse error")
	case binding.FuriosaSmiReturnCodeUnknownError:
		ret = errors.New("unknown error")
	case binding.FuriosaSmiReturnCodeInternalError:
		ret = errors.New("internal error")
	case binding.FuriosaSmiReturnCodeUninitializedError:
		ret = errors.New("uninitialized error")
	case binding.FuriosaSmiReturnCodeContextError:
		ret = errors.New("context error")
	}
	return
}
