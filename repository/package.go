package repository

// ////////////////////////////////////////////////////////////////////////////////// //

type VPackage struct {
	name     string // name of given package
	version  string // version of given package
	provider string // provider name of given package
}

type VPackageList []*VPackage

const (
	PKG_EXTENSION_TYPE = ".box"   // extension name with a leading dot
	PKG_CHECKSUM_TYPE  = "sha256" // type of supported hash function
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Name provides name of the package
func (b *VPackage) Name() string {
	return b.name
}

// Version provides version of the package
func (b *VPackage) Version() string {
	return b.version
}

// Provider provides provider name of the package
func (b *VPackage) Provider() string {
	return b.provider
}

// ////////////////////////////////////////////////////////////////////////////////// //

func NewPackage(name string, version string, provider string) *VPackage {
	b := &VPackage{
		name:     name,
		version:  version,
		provider: provider,
	}

	return b
}
