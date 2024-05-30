package smi

import (
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/smi/binding"
)

func Init() error {
	if ret := binding.FuriosaSmiInit(); ret != binding.FuriosaSmiReturnCodeOk {
		return ToError(ret)
	}

	return nil
}

func Shutdown() error {
	if ret := binding.FuriosaSmiShutdown(); ret != binding.FuriosaSmiReturnCodeOk {
		return ToError(ret)
	}

	return nil
}
