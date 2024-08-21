package npu_allocator

import (
	"fmt"
	"regexp"
)

const (
	bdfPattern   = `^(?P<domain>[0-9a-fA-F]{1,4}):(?P<bus>[0-9a-fA-F]+):(?P<function>[0-9a-fA-F]+\.[0-9])$`
	subExpKeyBus = "bus"
)

var (
	bdfRegExp = regexp.MustCompile(bdfPattern)
)

// ParseBusIDFromBDF parses bdf and returns PCI bus ID.
func ParseBusIDFromBDF(bdf string) (string, error) {
	matches := bdfRegExp.FindStringSubmatch(bdf)
	if matches == nil {
		return "", fmt.Errorf("couldn't parse the given string %s with bdf regex pattern: %s", bdf, bdfPattern)
	}

	subExpIndex := bdfRegExp.SubexpIndex(subExpKeyBus)
	if subExpIndex == -1 {
		return "", fmt.Errorf("couldn't parse bus id from the given bdf expression %s", bdf)
	}

	return matches[subExpIndex], nil
}
