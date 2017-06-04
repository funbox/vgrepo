package metadata

import (
	"fmt"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// VMetadataProvidersList provides a list of VMetadataProvider structs
type VMetadataProvidersList []*VMetadataProvider

// VMetadataProvider struct contains field of provider
type VMetadataProvider struct {
	Name         string `json:"name"`          // name of provider (e.g. virtualbox)
	URL          string `json:"url"`           // url to downloadable image
	Checksum     string `json:"checksum"`      // checksum string
	ChecksumType string `json:"checksum_type"` // checksum type of calculated checksum string
}

// ////////////////////////////////////////////////////////////////////////////////// //

// CountProviders returns count of providers
func (v *VMetadataVersion) CountProviders() int {
	return len(v.Providers)
}

// IsProviderExist returns true if provider exist
func (v *VMetadataVersion) IsProviderExist(provider *VMetadataProvider) bool {
	return v.AnyProviders(provider, isEqualProviders)
}

// isEqualProviders returns true if providers are equal
func isEqualProviders(first *VMetadataProvider, second *VMetadataProvider) bool {
	return first.Name == second.Name
}

// notEqualProviders returns true if providers are not equal
func notEqualProviders(first *VMetadataProvider, second *VMetadataProvider) bool {
	return !isEqualProviders(first, second)
}

// AnyProviders returns true if provider if present on the list
func (v *VMetadataVersion) AnyProviders(provider *VMetadataProvider, f func(*VMetadataProvider, *VMetadataProvider) bool) bool {
	for _, p := range v.Providers {
		if f(p, provider) {
			return true
		}
	}
	return false
}

// FindProvider returns provider by name
func (v *VMetadataVersion) FindProvider(name string) *VMetadataProvider {
	for _, p := range v.Providers {
		if p.Name == name {
			return p
		}
	}
	return nil
}

// FilterProvider filters list of providers in the metadata by given function
func (v *VMetadataVersion) FilterProvider(provider *VMetadataProvider, f func(*VMetadataProvider, *VMetadataProvider) bool) {
	providersList := make(VMetadataProvidersList, 0)
	for _, p := range v.Providers {
		if f(p, provider) {
			providersList = append(providersList, p)
		}
	}
	v.Providers = providersList
}

// AddProvider adds version to the metadata list
func (v *VMetadataVersion) AddProvider(provider *VMetadataProvider) error {
	if v.IsProviderExist(provider) {
		return fmt.Errorf(
			"provider %s is already exist",
			provider.Name,
		)
	}
	v.Providers = append(v.Providers, provider)

	return nil
}

// RemoveProvider removes provider from the list or do nothing
func (v *VMetadataVersion) RemoveProvider(provider *VMetadataProvider) {
	v.FilterProvider(provider, notEqualProviders)
}

// ////////////////////////////////////////////////////////////////////////////////// //

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
