package furiosa_smi_go

import "github.com/furiosa-ai/libfuriosa-kubernetes/pkg/furiosa_smi_go/binding"

type VersionInfo interface {
	Arch() Arch
	Major() uint32
	Minor() uint32
	Patch() uint32
	Metadata() string
}

var _ VersionInfo = new(versionInfo)

type versionInfo struct {
	raw binding.FuriosaSmiDriverVersion
}

func newVersionInfo(raw binding.FuriosaSmiDriverVersion) VersionInfo {
	return &versionInfo{
		raw: raw,
	}
}

func (v versionInfo) Arch() Arch {
	return Arch(v.raw.Arch)
}

func (v versionInfo) Major() uint32 {
	return v.raw.Major
}

func (v versionInfo) Minor() uint32 {
	return v.raw.Minor
}

func (v versionInfo) Patch() uint32 {
	return v.raw.Patch
}

func (v versionInfo) Metadata() string {
	return string(v.raw.Metadata[:])
}
