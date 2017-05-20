package storage

import (
	"pkg.re/essentialkaos/ek.v9/fsutil"
	"pkg.re/essentialkaos/ek.v9/path"

	"github.com/gongled/vgrepo/prefs"
	"github.com/gongled/vgrepo/repo"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type VStorage struct {
	*prefs.Preferences
	repositories repo.VRepositoryList
}

// ////////////////////////////////////////////////////////////////////////////////// //

func listDirs(dir string) []string {
	return fsutil.List(dir, false)
}

// ////////////////////////////////////////////////////////////////////////////////// //

func packagesPath(settings *prefs.Preferences) string {
	return path.Join(settings.StoragePath(), "packages")
}

func (s *VStorage) Repositories() repo.VRepositoryList {
	return s.repositories
}

func NewStorage(settings *prefs.Preferences) *VStorage {
	repositories := make(repo.VRepositoryList, 0)

	for _, r := range listDirs(packagesPath(settings)) {
		repositories = append(
			repositories, repo.NewRepository(
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
