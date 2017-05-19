package meta

// ////////////////////////////////////////////////////////////////////////////////// //

type VMetadataProvidersList []*VMetadataProvider

type VMetadataProvider struct {
	Name         string `json:"Name"`          // name of provider (e.g. virtualbox)
	URL          string `json:"url"`           // url to downloadable image
	Checksum     string `json:"checksum"`      // checksum string
	ChecksumType string `json:"checksum_type"` // checksum type of calculated checksum string
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
