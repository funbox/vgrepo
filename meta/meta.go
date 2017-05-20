package meta

import (
	"fmt"
	"strings"

	"pkg.re/essentialkaos/ek.v9/fsutil"
	"pkg.re/essentialkaos/ek.v9/jsonutil"
	"pkg.re/essentialkaos/ek.v9/path"
	"github.com/gongled/vgrepo/prefs"
	"os"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type VMetadata struct {
	*prefs.Preferences          // settings
	*VMetadataRepository        // metadata of the repository
}

// ////////////////////////////////////////////////////////////////////////////////// //

// nameMeta returns name of the metadata file with extension
func (m *VMetadata) nameMeta() string {
	return fmt.Sprintf("%s.json", m.Name)
}

// baseMeta returns name of metadata files
func (m *VMetadata) baseMeta() string {
	return path.Join(m.StoragePath(), "metadata")
}

// DirMeta returns directory string to metadata file
func (m *VMetadata) DirMeta() string {
	return path.Join(m.baseMeta(), m.Name)
}

// PathMeta returns full path to metadata file
func (m *VMetadata) PathMeta() string {
	return path.Join(m.DirMeta(), m.nameMeta())
}

// URLMeta returns direct link to metadata file
func (m *VMetadata) URLMeta() string {
	return fmt.Sprintf("%s/%s/%s/%s",
		strings.Trim(m.StorageURL(), "/"),
		"metadata",
		m.Name,
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
func NewMetadata(settings *prefs.Preferences, repository *VMetadataRepository) *VMetadata {
	m := &VMetadata{
		settings,
		repository,
	}

	m.SortVersions()

	return m
}

// ////////////////////////////////////////////////////////////////////////////////// //

// loadFromFile returns new VMetadataRepository struct which was read from the metadata file
func (m *VMetadata) loadFromFile(metaPath string) (*VMetadataRepository, error) {
	if !fsutil.IsExist(metaPath) {
		return nil, fmt.Errorf("Metadata %s does not exist", metaPath)
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
	var err error

	m.SortVersions()

	err = os.MkdirAll(path.Dir(metaPath), 0755)

	if err != nil {
		return err
	}

	return jsonutil.EncodeToFile(metaPath, m.VMetadataRepository)
}

// ReadMeta returns new VMetadata struct which was read from the metadata file
func (m *VMetadata) ReadMeta() (*VMetadata, error) {
	md, err := m.loadFromFile(m.PathMeta())

	return NewMetadata(
		m.Preferences,
		md,
	), err
}

// WriteMeta dumps VMetadata struct on the linked metadata file on the disk
func (m *VMetadata) WriteMeta() error {
	return m.dumpToFile(m.PathMeta())
}

// DeleteMeta deletes file with metadata from the disk
func (m *VMetadata) DeleteMeta() error {
	return os.RemoveAll(m.DirMeta())
}

// ////////////////////////////////////////////////////////////////////////////////// //
