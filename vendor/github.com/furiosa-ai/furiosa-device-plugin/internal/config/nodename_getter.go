package config

import "os"

type NodeNameGetter interface {
	GetNodename() string
}

func NewNodeNameGetter() NodeNameGetter {
	return &nodeNameGetterImpl{isCached: false, cache: ""}
}

type nodeNameGetterImpl struct {
	isCached bool
	cache    string
}

func (n *nodeNameGetterImpl) GetNodename() string {
	if n.isCached {
		return n.cache
	}
	nodeNameFromEnv := os.Getenv("NODE_NAME")
	n.isCached = true
	n.cache = nodeNameFromEnv
	return nodeNameFromEnv
}
