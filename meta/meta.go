package meta

import (
	"pkg.re/essentialkaos/ek.v8/fmtc"
	"pkg.re/essentialkaos/ek.v8/fsutil"
	"pkg.re/essentialkaos/ek.v8/jsonutil"
	"pkg.re/essentialkaos/ek.v8/path"

	"github.com/gongled/vgrepo/prefs"
	"strings"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	M_PATH_SUFFIX = "metadata"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type VMetadata struct {
	*VMetadataRepository
	Preferences *prefs.Preferences
}

type VMetadataRepository struct {
	Name        string             	`json:"name"`
	Description string             	`json:"description"`
	Versions    []*VMetadataVersion `json:"versions"`
}

type VMetadataVersion struct {
	Version   string              	`json:"version"`
	Providers []*VMetadataProvider 	`json:"providers"`
}

type VMetadataProvider struct {
	Name         string `json:"Name"`
	URL          string `json:"url"`
	Checksum     string `json:"checksum"`
	ChecksumType string `json:"checksum_type"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

func getMetaFormat(name string) string {
	return fmtc.Sprintf("%s.json", name)
}

func AreEqualVersions(first string, second string) bool {
	return first == second
}

func NotEqualVersions(first string, second string) bool {
	return !AreEqualVersions(first, second)
}

func (m *VMetadata) IsEmpty() bool {
	return len(m.Versions) == 0
}

func (m *VMetadata) Base() string {
	return path.Join(m.Preferences.StoragePath, m.Name)
}

func (m *VMetadata) Dir() string {
	return path.Join(m.Base(), M_PATH_SUFFIX)
}

func (m *VMetadata) Path() string {
	return path.Join(m.Dir(), getMetaFormat(m.Name))
}

func (m *VMetadata) Clone() *VMetadata {
	clone := *m
	return &clone
}

func (m *VMetadata) URL() string {
	return fmtc.Sprintf("%s/%s/%s/%s",
		strings.Trim(m.Preferences.StorageURL, "/"),
		m.Name,
		M_PATH_SUFFIX,
		getMetaFormat(m.Name),
	)
}

func (m *VMetadata) HasMeta() bool {
	return fsutil.IsExist(m.Path())
}

func (m *VMetadata) AnyVersion(version string, f func (string, string) bool) bool {
	for _, v := range m.Versions {
		if f(v.Version, version) {
			return true
		}
	}
	return false
}

func (m *VMetadata) FilterVersion(version string, f func(string, string) bool) *VMetadata {
	clone := m.Clone()

	versionsList := make([]*VMetadataVersion, 0)

	for _, v := range clone.Versions {
		if f(v.Version, version) {
			versionsList = append(versionsList, v)
		}
	}

	 clone.Versions = versionsList

	return clone
}

// ////////////////////////////////////////////////////////////////////////////////// //

func NewMetadata(prefs *prefs.Preferences, name string) *VMetadata {
	repo := NewMetadataRepository(name, "", make([]*VMetadataVersion, 0))

	m := &VMetadata{
		repo,
		prefs,
	}

	if m.HasMeta() {
		m, _ = m.Read()
	}

	return m
}

func NewMetadataRepository(name string, description string, versions []*VMetadataVersion) *VMetadataRepository {
	m := &VMetadataRepository{
		Name:        name,
		Description: description,
		Versions:    versions,
	}

	return m
}

func NewMetadataVersion(version string, providers []*VMetadataProvider) *VMetadataVersion {
	m := &VMetadataVersion{
		Version: version,
		Providers: providers,
	}

	return m
}

func NewMetadataProvider(name string, checksum string, checksumType string, url string) *VMetadataProvider {
	m := &VMetadataProvider{
		Name: name,
		Checksum: checksum,
		ChecksumType: checksumType,
		URL: url,
	}

	return m
}

// ////////////////////////////////////////////////////////////////////////////////// //

func (m *VMetadata) Read() (*VMetadata, error) {
	return m.LoadFromFile(m.Path())
}

func (m *VMetadata) LoadFromFile(metaPath string) (*VMetadata, error) {
	if !fsutil.IsExist(metaPath) {
		return nil, fmtc.Errorf("Metadata %s does not exist", metaPath)
	}

	info := &VMetadataRepository{}

	err := jsonutil.DecodeFile(metaPath, info)

	if err != nil {
		return nil, err
	}

	m.VMetadataRepository = info

	return m, err
}

func (m *VMetadata) Write(metaPath string) error {
	err := jsonutil.EncodeToFile(metaPath, m)

	return err
}

// ////////////////////////////////////////////////////////////////////////////////// //
