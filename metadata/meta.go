package metadata

import (
	"fmt"
	"os"
	"strings"

	"pkg.re/essentialkaos/ek.v9/fsutil"
	"pkg.re/essentialkaos/ek.v9/jsonutil"
	"pkg.re/essentialkaos/ek.v9/path"
	"pkg.re/essentialkaos/ek.v9/req"

	"github.com/gongled/vgrepo/prefs"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// MD_EXT_TYPE provides extension name of metadata file with a leading dot
const MD_EXT_TYPE = ".json"

// ////////////////////////////////////////////////////////////////////////////////// //

// VMetadata provides struct with the preferences and repository
type VMetadata struct {
	*prefs.Preferences   // settings
	*VMetadataRepository // metadata of the repository
}

// ////////////////////////////////////////////////////////////////////////////////// //

// metaName returns name of the metadata file with extension
func (m *VMetadata) metaName() string {
	return m.Name + MD_EXT_TYPE
}

// MetaDir returns directory string to metadata file
func (m *VMetadata) MetaDir() string {
	return path.Join(m.MetadataPath(), m.Name)
}

// MetaPath returns full path to metadata file
func (m *VMetadata) MetaPath() string {
	return path.Join(m.MetaDir(), m.metaName())
}

// MetaURL returns direct link to metadata file
func (m *VMetadata) MetaURL() string {
	return fmt.Sprintf("%s/%s/%s/%s",
		strings.Trim(m.StorageURL(), "/"),
		"metadata",
		m.Name,
		m.metaName(),
	)
}

// HasMeta returns true if metadata file is present on disk
func (m *VMetadata) HasMeta() bool {
	return fsutil.IsExist(m.MetaPath())
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
		return nil, fmt.Errorf("metadata %s does not exist", metaPath)
	}

	info := &VMetadataRepository{}

	err := jsonutil.DecodeFile(metaPath, info)

	if err != nil {
		return nil, err
	}

	return info, err
}

// loadFromURL returns new VMetadataRepository struct which was fetched from the HTTP/HTTPS resource
func (m *VMetadata) loadFromURL(metaURL string) (*VMetadataRepository, error) {
	resp, err := req.Request{URL: metaURL}.Get()

	if err != nil {
		return nil, fmt.Errorf("cannot fetch metadata from URL %s", metaURL)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(
			"cannot fetch metadata: server return status %s",
			resp.StatusCode,
		)
	}

	info := &VMetadataRepository{}

	err = resp.JSON(info)

	if err != nil {
		return nil, fmt.Errorf("cannot decode server response: %v", err)
	}

	return info, nil
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

// ReadMeta returns new VMetadata struct which was read from the metadata file or URL
func (m *VMetadata) ReadMeta() (*VMetadata, error) {
	md, err := m.loadFromFile(m.MetaPath())

	if err != nil {
		md, err = m.loadFromURL(m.MetaURL())
	}

	return NewMetadata(
		m.Preferences,
		md,
	), err
}

// WriteMeta dumps VMetadata struct on the linked metadata file on the disk
func (m *VMetadata) WriteMeta() error {
	return m.dumpToFile(m.MetaPath())
}

// DeleteMeta deletes file with metadata from the disk
func (m *VMetadata) DeleteMeta() error {
	return os.RemoveAll(m.MetaDir())
}

// ////////////////////////////////////////////////////////////////////////////////// //
