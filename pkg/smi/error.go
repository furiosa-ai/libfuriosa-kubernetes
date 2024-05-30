package smi

import (
	"errors"

	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/smi/binding"
)

func ToError(code binding.FuriosaSmiReturnCode) (ret error) {
	switch code {
	case binding.FuriosaSmiReturnCodeOk:
		ret = nil
	case binding.FuriosaSmiReturnCodeInitializeError:
		ret = errors.New("initialize error")
	case binding.FuriosaSmiReturnCodeUninitializedError:
		ret = errors.New("uninitialized error")
	case binding.FuriosaSmiReturnCodeInvalidArgumentError:
		ret = errors.New("invalid argument error")
	case binding.FuriosaSmiReturnCodeNullPointerError:
		ret = errors.New("null pointer error")
	case binding.FuriosaSmiReturnCodeMaxBufferSizeExceedError:
		ret = errors.New("max buffer size exceed error")
	case binding.FuriosaSmiReturnCodeDeviceFileNotFoundError:
		ret = errors.New("device file not found error")
	case binding.FuriosaSmiReturnCodeDeviceFileFormatError:
		ret = errors.New("device file format error")
	case binding.FuriosaSmiReturnCodeDeviceNotInUseError:
		ret = errors.New("device not in use error")
	case binding.FuriosaSmiReturnCodeDeviceNodeError:
		ret = errors.New("device node error")
	case binding.FuriosaSmiReturnCodeParseError:
		ret = errors.New("parse error")
	case binding.FuriosaSmiReturnCodeUnknownError:
		ret = errors.New("unknown error")
	}

	return
}
