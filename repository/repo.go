package repository

import (
	"fmt"
	"os"
	"strings"

	"pkg.re/essentialkaos/ek.v9/fsutil"
	"pkg.re/essentialkaos/ek.v9/hash"
	"pkg.re/essentialkaos/ek.v9/path"

	"github.com/gongled/vgrepo/metadata"
	"github.com/gongled/vgrepo/prefs"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// REPO_DEFAULT_DIR_PERM sets default permissions to new directories
const REPO_DEFAULT_DIR_PERM = 0755

// REPO_DEFAULT_PKG_PERM sets default permissions to new packages
const REPO_DEFAULT_PKG_PERM = 0644

// ////////////////////////////////////////////////////////////////////////////////// //

// VRepository provides a struct for repository entity
type VRepository struct {
	*metadata.VMetadata
}

// VRepositoryList is a list of VRepository structs
type VRepositoryList []*VRepository

// ////////////////////////////////////////////////////////////////////////////////// //

// boxNameFormat provides name with an extension
func (r *VRepository) boxNameFormat(pkg *VPackage) string {
	return fmt.Sprintf("%s%s", pkg.Name(), PKG_EXTENSION_TYPE)
}

// repoNameDir provides the directory with a name
func (r *VRepository) repoNameDir(pkg *VPackage) string {
	return path.Join(r.PackagesPath(), pkg.Name())
}

// repoVersionDir provides the directory with a version
func (r *VRepository) repoVersionDir(pkg *VPackage) string {
	return path.Join(r.repoNameDir(pkg), pkg.Version())
}

// repoProviderDir provides the directory with a provider name
func (r *VRepository) repoProviderDir(pkg *VPackage) string {
	return path.Join(r.repoVersionDir(pkg), pkg.Provider())
}

// ////////////////////////////////////////////////////////////////////////////////// //

// RepoPath provides the full path by given package parameters
func (r *VRepository) RepoPath(pkg *VPackage) string {
	return path.Join(r.repoProviderDir(pkg), r.boxNameFormat(pkg))
}

// RepoURL provides URL by given package parameters.
func (r *VRepository) RepoURL(pkg *VPackage) string {
	return fmt.Sprintf("%s/packages/%s/%s/%s/%s",
		strings.Trim(r.StorageURL(), "/"),
		pkg.Name(),
		pkg.Version(),
		pkg.Provider(),
		r.boxNameFormat(pkg),
	)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// copyPackage copies file from one source to another
func copyPackage(src string, dst string) error {
	var err error

	basedir := path.Dir(dst)

	err = os.MkdirAll(basedir, REPO_DEFAULT_DIR_PERM)

	if err != nil {
		return err
	}

	err = fsutil.CopyFile(src, dst, REPO_DEFAULT_PKG_PERM)

	if err != nil {
		return err
	}

	return err
}

// AddPackage adds package and import image by given source and package parameters.
func (r *VRepository) AddPackage(src string, pkg *VPackage) error {
	dst := r.RepoPath(pkg)

	// create empty provider list
	providerList := make(metadata.VMetadataProvidersList, 0)

	// add a new provider with calculated URL and SHA256 checksum.
	//
	// Note: in case you want to change checksum type, make sure that it is
	// well supported by Vagrant client. See the list of current supported values on
	// https://www.vagrantup.com/docs/vagrantfile/machine_settings.html
	providerList = append(providerList,
		metadata.NewMetadataProvider(
			pkg.Provider(),
			hash.FileHash(src),
			PKG_CHECKSUM_TYPE,
			r.RepoURL(pkg),
		),
	)

	version := metadata.NewMetadataVersion(pkg.Version(), providerList)

	err := r.AddVersion(version)

	if err != nil {
		return err
	}

	err = copyPackage(src, dst)

	if err != nil {
		return err
	}

	return r.WriteMeta()
}

// RemovePackage removes package or repository by given package parameters
func (r *VRepository) RemovePackage(pkg *VPackage) error {
	version := r.FindVersion(pkg.Version())

	if version == nil {
		return fmt.Errorf("unable to find version %s", pkg.Version())
	}

	provider := version.FindProvider(pkg.Provider())

	if provider == nil {
		return fmt.Errorf("unable to find provider %s", pkg.Provider())
	}

	// Path for deletion
	deletedPath := r.RepoPath(pkg)

	// Delete the provider from the list if the length of the list larger than one.
	// In other case it make sense to remove the whole version.
	if version.CountProviders() > 1 {
		deletedPath = r.repoProviderDir(pkg)
		version.RemoveProvider(provider)
	} else {
		deletedPath = r.repoVersionDir(pkg)
		r.RemoveVersion(version)
	}

	// In case metadata is empty (e.g. version list is empty) we have to
	// purge the repository.
	if r.IsEmptyMeta() {
		deletedPath = r.repoNameDir(pkg)
		r.DeleteMeta()
	} else {
		r.WriteMeta()
	}

	// Remove package or repository
	return os.RemoveAll(deletedPath)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// NewRepository returns new VRepository struct by given parameters or exits metadata
func NewRepository(settings *prefs.Preferences, name string) *VRepository {
	m := metadata.NewMetadata(
		settings,
		metadata.NewMetadataRepository(
			name,
			"",
			make(metadata.VMetadataVersionsList, 0),
		),
	)

	r := &VRepository{m}

	// Trying to create struct by given metadata JSON file
	if m.HasMeta() {
		r.VMetadata, _ = r.ReadMeta()
	}

	return r
}

// ////////////////////////////////////////////////////////////////////////////////// //
