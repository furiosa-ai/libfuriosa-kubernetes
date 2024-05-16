package cdi

import (
	"github.com/furiosa-ai/libfuriosa-kubernetes/pkg/manifest"
	"path/filepath"
	"tags.cncf.io/container-device-interface/pkg/cdi"
)

// "/var/run/cdi" is directory for generated CDI specs that runtimes watches

type Sync interface {
	Sync() error
	ListSpecs() ([]*cdi.Spec, error)
	CreateSpec(spec *cdi.Spec, key string) error
	GetSpec(specName string) (*cdi.Spec, error)
	DeleteSpec(specName string) error
}

var _ Sync = (*sync)(nil)

type sync struct {
	cdiRoot  string
	registry cdi.Registry
	cache    cdi.Cache
}

func NewSync(devices []manifest.Manifest) Sync {
	return newSync(cdi.DefaultDynamicDir, devices)
}

func newSync(cdiRoot string, devices []manifest.Manifest) Sync {
	return &sync{
		cdiRoot:  cdiRoot,
		registry: cdi.GetRegistry(cdi.WithSpecDirs(cdiRoot)),
	}
}

func (s sync) ListSpecs() ([]*cdi.Spec, error) {
	return nil, nil
}

func (s sync) Sync() error {
	// list all specs

	// validate specs?

	// generate specs from current states

	// compares and updates

	return nil
}

func (s sync) CreateSpec(spec *cdi.Spec, key string) error {
	s.cache.ListDevices()
	s.cache.GetDevice("")
	s.cache.RemoveSpec()
	s.cache.ListClasses()

	specName, err := cdi.GenerateNameForTransientSpec(spec, key)
	if err != nil {
		return err
	}

	specPath := filepath.Join(s.cdiRoot, specName+".json")
	err = s.registry.SpecDB().WriteSpec(spec, specPath)
	if err != nil {
		return err
	}

}

func (s sync) GetSpec(specName string) (*cdi.Spec, error) {
	//TODO implement me
	panic("implement me")
}

func (s sync) DeleteSpec(specName string) error {
	//TODO implement me
	panic("implement me")
}
