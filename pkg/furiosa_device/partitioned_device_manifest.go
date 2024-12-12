package furiosa_device

import (
	"regexp"
	"strconv"

	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/manifest"
)

type partitionedDeviceManifest struct {
	arch        smi.Arch
	original    manifest.Manifest
	partition   Partition
	deviceNodes []*manifest.DeviceNode
	mounts      []*manifest.Mount
}

func NewPartitionedDeviceManifest(arch smi.Arch, original manifest.Manifest, partition Partition) (manifest.Manifest, error) {
	deviceNodes, err := filterPartitionedDeviceNodes(original, partition)
	if err != nil {
		return nil, err
	}

	return &partitionedDeviceManifest{
		arch:        arch,
		original:    original,
		partition:   partition,
		deviceNodes: deviceNodes,

		// do not filter any mount paths right now as right now, we don't need to filter any mount paths.
		// see: https://github.com/furiosa-ai/furiosa-device-plugin/pull/30#discussion_r1819763238
		mounts: original.MountPaths(),
	}, nil
}

func (p *partitionedDeviceManifest) EnvVars() map[string]string {
	return p.original.EnvVars()
}

func (p *partitionedDeviceManifest) Annotations() map[string]string {
	return p.original.Annotations()
}

func (p *partitionedDeviceManifest) DeviceNodes() []*manifest.DeviceNode {
	return p.deviceNodes
}

func (p *partitionedDeviceManifest) MountPaths() []*manifest.Mount {
	return p.mounts
}

const (
	deviceIdExp  = "device_id"
	startCoreExp = "start_core"
	endCoreExp   = "end_core"

	regexpPattern = `^\S+npu(?P<` + deviceIdExp + `>\d+)((?:pe)(?P<` + startCoreExp + `>\d+)(-(?P<` + endCoreExp + `>\d+))?)?$`
)

var (
	deviceNodePeRegex = regexp.MustCompile(regexpPattern)
	deviceNodeSubExps = deviceNodePeRegex.SubexpNames()
)

// filterPartitionedDeviceNodes filters Device Nodes by following rules.
//   - npu{N}pe{X} will be dropped if X is not in the range between `peLowerBound` and `peUpperBound`.
//   - npu{N}pe{X}-{Y} will be dropped if X and Y are not in the range between `peLowerBound` and `peUpperBound`.
//
// e.g. If strategy is QuadCore and partition range is 0 to 3, below device files will be assigned.
//   - npu0pe0-3
//   - npu0pe0-1, npu0pe2-3
//   - npu0pe0, npu0pe1, npu0pe2, npu0pe3
func filterPartitionedDeviceNodes(original manifest.Manifest, partition Partition) ([]*manifest.DeviceNode, error) {
	peLowerBound, peUpperBound := partition.start, partition.end

	var survivedDeviceNodes []*manifest.DeviceNode
	for _, deviceNode := range original.DeviceNodes() {
		path := deviceNode.ContainerPath
		matches := deviceNodePeRegex.FindStringSubmatch(path)
		namedMatches := map[string]string{}
		for i, match := range matches {
			subExp := deviceNodeSubExps[i]
			if subExp == "" {
				continue
			}

			namedMatches[subExp] = match
		}

		if len(namedMatches) > 1 {
			startCore := namedMatches[startCoreExp]
			endCore := namedMatches[endCoreExp]
			if endCore == "" {
				endCore = startCore
			}

			peStartNum, err := strconv.Atoi(startCore)
			if err != nil {
				return nil, err
			}

			peEndNum, err := strconv.Atoi(endCore)
			if err != nil {
				return nil, err
			}

			if peLowerBound <= peStartNum && peEndNum <= peUpperBound {
				survivedDeviceNodes = append(survivedDeviceNodes, deviceNode)
			}
		} else {
			survivedDeviceNodes = append(survivedDeviceNodes, deviceNode)
		}
	}

	return survivedDeviceNodes, nil
}
