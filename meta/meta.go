package meta

import (
	"strings"

	"pkg.re/essentialkaos/ek.v8/fmtc"
	"pkg.re/essentialkaos/ek.v8/fsutil"
	"pkg.re/essentialkaos/ek.v8/jsonutil"
	"pkg.re/essentialkaos/ek.v8/path"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type VMetadata struct {
	*VMetadataRepository        // metadata of the repository
	StorageURL           string // repository URL and port
	StoragePath          string // repository path to store images and metadata
}

// ////////////////////////////////////////////////////////////////////////////////// //

// nameMeta returns name of the metadata file with extension
func (m *VMetadata) nameMeta() string {
	return fmtc.Sprintf("%s.json", m.Name)
}

// baseMeta returns name of metadata files
func (m *VMetadata) baseMeta() string {
	return path.Join(m.StoragePath, m.Name)
}

// DirMeta returns directory string to metadata file
func (m *VMetadata) DirMeta() string {
	return path.Join(m.baseMeta(), "metadata")
}

// PathMeta returns full path to metadata file
func (m *VMetadata) PathMeta() string {
	return path.Join(m.DirMeta(), m.nameMeta())
}

// URLMeta returns direct link to metadata file
func (m *VMetadata) URLMeta() string {
	return fmtc.Sprintf("%s/%s/%s/%s",
		strings.Trim(m.StorageURL, "/"),
		m.Name,
		"metadata",
		m.nameMeta(),
	)
}

// HasMeta returns true if metadata file is present on disk
func (m *VMetadata) HasMeta() bool {
	return fsutil.IsExist(m.PathMeta())
}

// IsEmptyMeta returns true if versions list is empty
func (m *VMetadata) IsEmptyMeta() bool {
	return m.CountVersions() == 0
}

// ////////////////////////////////////////////////////////////////////////////////// //

// NewMetadata returns new VMetadata struct
func NewMetadata(storagePath string, storageUrl string, repository *VMetadataRepository) *VMetadata {
	m := &VMetadata{
		StoragePath:         storagePath,
		StorageURL:          storageUrl,
		VMetadataRepository: repository,
	}

	m.SortVersions()

	return m
}

// ////////////////////////////////////////////////////////////////////////////////// //

// loadFromFile returns new VMetadataRepository struct which was read from the metadata file
func (m *VMetadata) loadFromFile(metaPath string) (*VMetadataRepository, error) {
	if !fsutil.IsExist(metaPath) {
		return nil, fmtc.Errorf("Metadata %s does not exist", metaPath)
	}

	info := &VMetadataRepository{}

	err := jsonutil.DecodeFile(metaPath, info)

	if err != nil {
		return nil, err
	}

	return info, err
}


// dumpToFile dumps VMetadata struct on the metadata file on the disk
func (m *VMetadata) dumpToFile(metaPath string) error {
	m.SortVersions()

	return jsonutil.EncodeToFile(metaPath, m.VMetadataRepository)
}

// ReadMeta returns new VMetadata struct which was read from the metadata file
func (m *VMetadata) ReadMeta() (*VMetadata, error) {
	md, err := m.loadFromFile(m.PathMeta())

	return NewMetadata(
		m.StoragePath,
		m.StorageURL,
		md,
	), err
}

// WriteMeta dumps VMetadata struct on the linked metadata file on the disk
func (m *VMetadata) WriteMeta() error {
	return m.dumpToFile(m.PathMeta())
}

// ////////////////////////////////////////////////////////////////////////////////// //
