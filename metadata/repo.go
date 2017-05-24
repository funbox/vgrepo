package metadata

// ////////////////////////////////////////////////////////////////////////////////// //

type VMetadataRepository struct {
	Name        string                `json:"name"`        // name of the repository
	Description string                `json:"description"` // description of the repository
	Versions    VMetadataVersionsList `json:"versions"`    // list of available versions
}

// ////////////////////////////////////////////////////////////////////////////////// //

// NewMetadataRepository returns new VMetadataRepository struct
func NewMetadataRepository(name string, description string, versions VMetadataVersionsList) *VMetadataRepository {
	m := &VMetadataRepository{
		Name:        name,
		Description: description,
		Versions:    versions,
	}

	return m
}
