package cdi_spec

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/furiosa-ai/furiosa-smi-go/pkg/smi"
	"tags.cncf.io/container-device-interface/specs-go"
)

func NewPartitionedDeviceSpecRenderer(device smi.Device, coreStart int, coreEnd int) (Renderer, error) {
	deviceInfo, err := device.DeviceInfo()
	if err != nil {
		return nil, err
	}

	var deviceSpec DeviceSpec = nil
	switch deviceInfo.Arch() {
	case smi.ArchWarboy:
		return nil, fmt.Errorf("partitioned device is not supported for Warboy")
	case smi.ArchRngd:
		deviceSpec, err = newRngdDeviceSpec(device)
	}

	if err != nil {
		return nil, err
	}

	return &partitionedDeviceSpecRenderer{
		spec:      deviceSpec,
		coreStart: coreStart,
		coreEnd:   coreEnd,
	}, nil
}

type partitionedDeviceSpecRenderer struct {
	spec      DeviceSpec
	coreStart int
	coreEnd   int
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

func (p partitionedDeviceSpecRenderer) Render() *specs.Device {
	mutatedSpec := p.spec.DeviceSpec()
	mutatedSpec.ContainerEdits.DeviceNodes = filterPartitionedDeviceNodes(p.spec, p.coreStart, p.coreEnd)
	return mutatedSpec
}

// filterPartitionedDeviceNodes filters Device Nodes by following rules.
//   - npu{N}pe{X} will be dropped if X is not in the range between `peLowerBound` and `peUpperBound`.
//   - npu{N}pe{X}-{Y} will be dropped if X and Y are not in the range between `peLowerBound` and `peUpperBound`.
//
// e.g. If strategy is QuadCore and partition range is 0 to 3, below device files will be assigned.
//   - npu0pe0-3
//   - npu0pe0-1, npu0pe2-3
//   - npu0pe0, npu0pe1, npu0pe2, npu0pe3
func filterPartitionedDeviceNodes(original DeviceSpec, startCore int, endCore int) []*specs.DeviceNode {
	peLowerBound, peUpperBound := startCore, endCore

	var survivedDeviceNodes []*specs.DeviceNode
	for _, deviceNode := range original.deviceNodes() {
		path := deviceNode.Path
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

			// Note(@bg): we can ensure  startCore and endCore is always a number because of the regexp.
			peStartNum, _ := strconv.Atoi(startCore)
			peEndNum, _ := strconv.Atoi(endCore)

			if peLowerBound <= peStartNum && peEndNum <= peUpperBound {
				survivedDeviceNodes = append(survivedDeviceNodes, deviceNode)
			}
		} else {
			survivedDeviceNodes = append(survivedDeviceNodes, deviceNode)
		}
	}

	return survivedDeviceNodes
}
