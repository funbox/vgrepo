package meta

import (
	"fmt"
	"sort"

	"pkg.re/essentialkaos/ek.v9/sortutil"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type VMetadataVersionsList []*VMetadataVersion

type VMetadataVersion struct {
	Version   string                 `json:"version"`   // version of the image
	Providers VMetadataProvidersList `json:"providers"` // list of available providers
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
	return sortutil.VersionCompare(s[i].Version, s[j].Version)
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
func (m *VMetadata) LatestVersion() *VMetadataVersion {
	if !m.IsEmptyMeta() {
		return m.Versions[m.CountVersions()-1]
	} else {
		return nil
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
		return fmt.Errorf(
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

// NewMetadataVersion returns new VMetadataVersion struct
func NewMetadataVersion(version string, providers VMetadataProvidersList) *VMetadataVersion {
	m := &VMetadataVersion{
		Version:   version,
		Providers: providers,
	}

	return m
}
