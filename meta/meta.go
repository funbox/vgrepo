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
	*VMetadataRepository
	StorageURL  string
	StoragePath string
}

type VMetadataRepository struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Versions    []*VMetadataVersion `json:"versions"`
}

type VMetadataVersion struct {
	Version   string               `json:"version"`
	Providers []*VMetadataProvider `json:"providers"`
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

//func AreEqualVersions(first string, second string) bool {
//	return first == second
//}

//func NotEqualVersions(first string, second string) bool {
//	return !AreEqualVersions(first, second)
//}

func (m *VMetadata) IsEmptyMeta() bool {
	return len(m.Versions) == 0
}

func (m *VMetadata) BaseMeta() string {
	return path.Join(m.StoragePath, m.Name)
}

func (m *VMetadata) DirMeta() string {
	return path.Join(m.BaseMeta(), "metadata")
}

func (m *VMetadata) PathMeta() string {
	return path.Join(m.DirMeta(), getMetaFormat(m.Name))
}

func (m *VMetadata) CloneMeta() *VMetadata {
	clone := *m
	return &clone
}

func (m *VMetadata) URLMeta() string {
	return fmtc.Sprintf("%s/%s/%s/%s",
		strings.Trim(m.StorageURL, "/"),
		m.Name,
		"metadata",
		getMetaFormat(m.Name),
	)
}

func (m *VMetadata) HasMeta() bool {
	return fsutil.IsExist(m.PathMeta())
}

func (m *VMetadata) AnyVersion(version string, f func(string, string) bool) bool {
	for _, v := range m.Versions {
		if f(v.Version, version) {
			return true
		}
	}
	return false
}

func (m *VMetadata) FilterVersion(version string, f func(string, string) bool) *VMetadata {
	clone := m.CloneMeta()

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

func NewMetadata(storagePath, storageUrl, name, description string, versions []*VMetadataVersion) *VMetadata {
	m := &VMetadata{
		StoragePath:         storagePath,
		StorageURL:          storageUrl,
		VMetadataRepository: NewMetadataRepository(name, description, versions),
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

//func NewMetadataVersion(version string, providers []*VMetadataProvider) *VMetadataVersion {
//	m := &VMetadataVersion{
//		Version:   version,
//		Providers: providers,
//	}
//
//	return m
//}
//
//func NewMetadataProvider(name string, checksum string, checksumType string, url string) *VMetadataProvider {
//	m := &VMetadataProvider{
//		Name:         name,
//		Checksum:     checksum,
//		ChecksumType: checksumType,
//		URL:          url,
//	}
//
//	return m
//}

// ////////////////////////////////////////////////////////////////////////////////// //

func (m *VMetadata) ReadMeta() (*VMetadata, error) {
	f, err := m.LoadFromFile(m.PathMeta())

	return NewMetadata(m.StoragePath, m.StorageURL, f.Name, f.Description, f.Versions), err
}

func (m *VMetadata) LoadFromFile(metaPath string) (*VMetadataRepository, error) {
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

func (m *VMetadata) WriteMeta(metaPath string) error {
	return jsonutil.EncodeToFile(metaPath, m.VMetadataRepository)
}

// ////////////////////////////////////////////////////////////////////////////////// //
