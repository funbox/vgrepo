package meta

import (
	"sort"
	"strings"

	"pkg.re/essentialkaos/ek.v8/fmtc"
	"pkg.re/essentialkaos/ek.v8/fsutil"
	"pkg.re/essentialkaos/ek.v8/jsonutil"
	"pkg.re/essentialkaos/ek.v8/path"
	"pkg.re/essentialkaos/ek/version"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type VMetadata struct {
	*VMetadataRepository        // metadata of the repository
	StorageURL           string // repository URL and port
	StoragePath          string // repository path to store images and metadata
}

type VMetadataVersionsList []*VMetadataVersion

type VMetadataRepository struct {
	Name        string                `json:"name"`        // name of the repository
	Description string                `json:"description"` // description of the repository
	Versions    VMetadataVersionsList `json:"versions"`    // list of available versions
}

type VMetadataProvidersList []*VMetadataProvider

type VMetadataVersion struct {
	Version   string                 `json:"version"`   // version of the image
	Providers VMetadataProvidersList `json:"providers"` // list of available providers
}

type VMetadataProvider struct {
	Name         string `json:"Name"`          // name of provider (e.g. virtualbox)
	URL          string `json:"url"`           // url to downloadable image
	Checksum     string `json:"checksum"`      // checksum string
	ChecksumType string `json:"checksum_type"` // checksum type of calculated checksum string
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

// Len implements interface method for Sort
func (s VMetadataVersionsList) Len() int {
	return len(s)
}

// Swap implements interface method for Sort
func (s VMetadataVersionsList) Swap(i, j int) {
	*s[i], *s[j] = *s[j], *s[i]
}

// Less implements interface method for Sort
func (s VMetadataVersionsList) Less(i, j int) bool {
	first, _ := version.Parse(s[i].Version)
	second, _ := version.Parse(s[j].Version)

	return first.Less(second)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// equalVersions returns true if versions are equal
func equalVersions(first string, second string) bool {
	return first == second
}

// notEqualVersions returns true if versions are not equal
func notEqualVersions(first string, second string) bool {
	return !equalVersions(first, second)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// AnyVersions returns true if version if present on the list
func (m *VMetadata) AnyVersions(version string, f func(string, string) bool) bool {
	for _, v := range m.Versions {
		if f(v.Version, version) {
			return true
		}
	}
	return false
}

// IsVersionExist returns true if version is already exist in the metadata
func (m *VMetadata) IsVersionExist(version string) bool {
	return m.AnyVersions(version, equalVersions)
}

func (m *VMetadata) CountVersions() int {
	return len(m.Versions)
}

// OldestVersion returns the oldest version from the list
func (m *VMetadata) OldestVersion() string {
	if !m.IsEmptyMeta() {
		return m.Versions[0].Version
	} else {
		return ""
	}
}

// LatestVersion returns the latest version from the list
func (m *VMetadata) LatestVersion() string {
	if !m.IsEmptyMeta() {
		return m.Versions[m.CountVersions()-1].Version
	} else {
		return ""
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// SortVersions sorts list of versions in the metadata
func (m *VMetadata) SortVersions() {
	sort.Sort(VMetadataVersionsList(m.Versions))
}

// AddVersion adds version to the metadata list
func (m *VMetadata) AddVersion(version *VMetadataVersion) error {
	if m.IsVersionExist(version.Version) {
		return fmtc.Errorf(
			"Cannot add version to metadata: version %s is already exist",
			version.Version,
		)
	}

	m.Versions = append(m.Versions, version)
	m.SortVersions()

	return nil
}

// FilterVersion filters list of versions in the metadata by given function
func (m *VMetadata) FilterVersion(version string, f func(string, string) bool) {
	versionsList := make(VMetadataVersionsList, 0)

	for _, v := range m.Versions {
		if f(v.Version, version) {
			versionsList = append(versionsList, v)
		}
	}

	m.Versions = versionsList
}

// RemoveVersion removes version from the list or do nothing
func (m *VMetadata) RemoveVersion(version string) {
	m.FilterVersion(version, notEqualVersions)
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

// NewMetadataRepository returns new VMetadataRepository struct
func NewMetadataRepository(name string, description string, versions VMetadataVersionsList) *VMetadataRepository {
	m := &VMetadataRepository{
		Name:        name,
		Description: description,
		Versions:    versions,
	}

	return m
}

// NewMetadataVersion returns new VMetadataVersion struct
func NewMetadataVersion(version string, providers VMetadataProvidersList) *VMetadataVersion {
	m := &VMetadataVersion{
		Version:   version,
		Providers: providers,
	}

	return m
}

// NewMetadataProvider returns new VMetadataProvider struct
func NewMetadataProvider(name string, checksum string, checksumType string, url string) *VMetadataProvider {
	m := &VMetadataProvider{
		Name:         name,
		Checksum:     checksum,
		ChecksumType: checksumType,
		URL:          url,
	}

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
