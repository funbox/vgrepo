package storage

import (
	"pkg.re/essentialkaos/ek.v9/fsutil"

	"github.com/gongled/vgrepo/prefs"
	"github.com/gongled/vgrepo/repository"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// VStorage provides struct for storage entity
type VStorage struct {
	*prefs.Preferences
	repositories repository.VRepositoryList
}

// ////////////////////////////////////////////////////////////////////////////////// //

// listDirs returns list of directories by given path
func listDirs(dir string) []string {
	return fsutil.List(dir, true, fsutil.ListingFilter{Perms: "DRX"})
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Repositories provides a list of repositories
func (s *VStorage) Repositories() repository.VRepositoryList {
	return s.repositories
}

// ////////////////////////////////////////////////////////////////////////////////// //

// NewStorage returns VStorage struct by given preferences
func NewStorage(settings *prefs.Preferences) *VStorage {
	repositories := make(repository.VRepositoryList, 0)

	for _, r := range listDirs(settings.PackagesPath()) {
		repositories = append(
			repositories, repository.NewRepository(
				settings,
				r,
			),
		)
	}

	return &VStorage{
		settings,
		repositories,
	}
}
