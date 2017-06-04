package repository

// ////////////////////////////////////////////////////////////////////////////////// //

// VPackage provides a struct with an input parameters
type VPackage struct {
	name     string // name of given package
	version  string // version of given package
	provider string // provider name of given package
}

// VPackageList is a list of packages (VPackage structs)
type VPackageList []*VPackage

// PKG_EXTENSION_TYPE is an extension name with a leading dot
const PKG_EXTENSION_TYPE = ".box"

// PKG_CHECKSUM_TYPE is a type of supported hash function
const PKG_CHECKSUM_TYPE  = "sha256"

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

// NewPackage returns new VPackage struct by given parameters
func NewPackage(name string, version string, provider string) *VPackage {
	b := &VPackage{
		name:     name,
		version:  version,
		provider: provider,
	}

	return b
}
