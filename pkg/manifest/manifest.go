package manifest

type Manifest interface {
	// EnvVars returns map consisting of env variable name and value
	EnvVars() map[string]string
	// Annotations returns set of annotation for CRI Runtime injection
	// https://github.com/kubernetes/kubernetes/pull/58172
	Annotations() map[string]string
	// DeviceNodes returns set of DeviceNode for dev file mount
	DeviceNodes() []*DeviceNode
	// MountPaths returns set of Mount for extra file and directory mount
	MountPaths() []*Mount
}

// Mount is subset of oci-runtime Mount spec
type Mount struct {
	// mount path in container environment
	ContainerPath string
	// origin host path
	HostPath string
	// mount options such as "nosuid", "nodev", "bind", "noexec" and file permission("ro", "rw", ...)
	Options []string
}

// DeviceNode is subset struct of oci-runtime DeviceNode spec
type DeviceNode struct {
	// mount path in container environment
	ContainerPath string
	// origin host path
	HostPath string
	// Cgroups permissions of the device, candidates are one or more of
	// * r - allows container to read from the specified device.
	// * w - allows container to write to the specified device.
	// * m - allows container to create device files that do not yet exist.
	Permissions string
}
