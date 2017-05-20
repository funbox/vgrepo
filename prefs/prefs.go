package prefs

type Preferences struct {
	storagePath string
	storageURL  string
}

// ////////////////////////////////////////////////////////////////////////////////// //

func (p *Preferences) StoragePath() string {
	return p.storagePath
}

func (p *Preferences) StorageURL() string {
	return p.storageURL
}

// ////////////////////////////////////////////////////////////////////////////////// //

func NewPreferences(storagePath string, storageURL string) *Preferences {
	return &Preferences{
		storagePath,
		storageURL,
	}
}
