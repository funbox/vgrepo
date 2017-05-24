package prefs

import (
	"pkg.re/essentialkaos/ek.v9/path"
)

type Preferences struct {
	storagePath string // path to packages and metadata
	storageURL  string // URL and port to packages and metadata
}

// ////////////////////////////////////////////////////////////////////////////////// //

// StoragePath provides path to store packages and metadata
func (p *Preferences) StoragePath() string {
	return p.storagePath
}

// StorageURL shares URL to store and download packages with metadata
func (p *Preferences) StorageURL() string {
	return p.storageURL
}

// PackagesPath provides path to packages directory
func (p *Preferences) PackagesPath() string {
	return path.Join(p.StoragePath(), "packages")
}

// MetadataPath provides path to metadata directory
func (p *Preferences) MetadataPath() string {
	return path.Join(p.StoragePath(), "metadata")
}

// ////////////////////////////////////////////////////////////////////////////////// //

func NewPreferences(storagePath string, storageURL string) *Preferences {
	return &Preferences{
		storagePath,
		storageURL,
	}
}
