package smi

import (
	"fmt"

	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi/binding"
)

// VersionInfo represents a version information.
type VersionInfo interface {
	fmt.Stringer // added for `String() string` method

	// Major returns a major part of version.
	Major() uint32
	// Minor returns a minor part of version.
	Minor() uint32
	// Patch returns a patch part of version.
	Patch() uint32
	// Metadata returns a metadata of version.
	Metadata() string
	// Prerelease returns a prerelease of version.
	Prerelease() string
}

var _ VersionInfo = new(versionInfo)

type versionInfo struct {
	raw binding.FuriosaSmiVersion
}

func newVersionInfo(raw binding.FuriosaSmiVersion) VersionInfo {
	return &versionInfo{raw: raw}
}

func (v *versionInfo) Major() uint32 {
	return v.raw.Major
}

func (v *versionInfo) Minor() uint32 {
	return v.raw.Minor
}

func (v *versionInfo) Patch() uint32 {
	return v.raw.Patch
}

func (v *versionInfo) Metadata() string {
	return byteBufferToString(v.raw.Metadata[:])
}

func (v *versionInfo) Prerelease() string {
	return byteBufferToString(v.raw.Prerelease[:])
}

func (v *versionInfo) String() string {
	prerelease := v.Prerelease()

	if prerelease == "" {
		return fmt.Sprintf("%d.%d.%d, %s", v.Major(), v.Minor(), v.Patch(), v.Metadata())
	}

	return fmt.Sprintf("%d.%d.%d(%s), %s", v.Major(), v.Minor(), v.Patch(), prerelease, v.Metadata())
}
