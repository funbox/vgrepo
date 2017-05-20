package repo

import (
	"fmt"

	"pkg.re/essentialkaos/ek.v9/path"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const DEFAULT_PROVIDER = "virtualbox"

type VPackage struct {
	Name     string
	Version  string
	Provider string
}

type VPackageList []*VPackage

// ////////////////////////////////////////////////////////////////////////////////// //

func (b *VPackage) nameBoxFormat() string {
	return fmt.Sprintf("%s.box", b.Name)
}

func (b *VPackage) DirBoxFormat() string {
	return path.Join(b.Name, b.Version, b.Provider)
}

func (b *VPackage) PathBoxFormat() string {
	return path.Join(b.DirBoxFormat(), b.nameBoxFormat())
}

func (b *VPackage) URLBoxFormat() string {
	return fmt.Sprintf(
		"packages/%s/%s/%s/%s",
		b.Name,
		b.Version,
		b.Provider,
		b.nameBoxFormat(),
	)
}

// ////////////////////////////////////////////////////////////////////////////////// //

func NewPackage(name string, version string, provider string) *VPackage {
	if provider == "" {
		provider = DEFAULT_PROVIDER
	}

	b := &VPackage{
		Name:     name,
		Version:  version,
		Provider: provider,
	}

	return b
}
