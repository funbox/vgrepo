package repo

import (
	"fmt"
	"os"
	"strings"

	"pkg.re/essentialkaos/ek.v9/fsutil"
	"pkg.re/essentialkaos/ek.v9/hash"
	"pkg.re/essentialkaos/ek.v9/path"

	"github.com/gongled/vgrepo/meta"
	"github.com/gongled/vgrepo/prefs"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type VRepository struct {
	*meta.VMetadata
}

type VRepositoryList []*VRepository

// ////////////////////////////////////////////////////////////////////////////////// //

func (r *VRepository) DirRepo() string {
	return path.Join(r.StoragePath(), "packages")
}

func (r *VRepository) PathRepo(pkg *VPackage) string {
	return path.Join(r.DirRepo(), pkg.PathBoxFormat())
}

func (r *VRepository) URLRepo(pkg *VPackage) string {
	return fmt.Sprintf("%s/%s",
		strings.Trim(r.StorageURL(), "/"),
		pkg.URLBoxFormat(),
	)
}

// ////////////////////////////////////////////////////////////////////////////////// //

func (r *VRepository) CreateRepo() error {
	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

func (r *VRepository) copyPackage(src string, dst string) error {
	var err error

	basedir := path.Dir(dst)

	err = os.MkdirAll(basedir, 0755)

	if err != nil {
		return err
	}

	err = fsutil.CopyFile(src, dst, 0644)

	if err != nil {
		return err
	}

	return err
}

func (r *VRepository) AddPackage(src string, pkg *VPackage) error {
	r.CreateRepo()

	dst := r.PathRepo(pkg)
	err := r.copyPackage(src, dst)

	if err != nil {
		return err
	}

	providerList := make(meta.VMetadataProvidersList, 0)

	providerList = append(providerList,
		meta.NewMetadataProvider(
			pkg.Provider,
			hash.FileHash(dst),
			"sha256",
			r.URLRepo(pkg),
		),
	)

	version := meta.NewMetadataVersion(pkg.Version, providerList)

	r.AddVersion(version)
	r.WriteMeta()

	return err
}

func (r *VRepository) RemovePackage(pkg *VPackage) error {
	err := os.RemoveAll(r.PathRepo(pkg))

	if err != nil {
		return err
	}

	r.RemoveVersion(pkg.Version)

	if r.IsEmptyMeta() {
		r.DeleteMeta()
	}

	return err
}

// ////////////////////////////////////////////////////////////////////////////////// //

func NewRepository(settings *prefs.Preferences, name string) *VRepository {
	m := meta.NewMetadata(
		settings,
		meta.NewMetadataRepository(
			name,
			"",
			make(meta.VMetadataVersionsList, 0),
		),
	)

	r := &VRepository{m}

	if m.HasMeta() {
		r.VMetadata, _ = r.ReadMeta()
	}

	return r
}

// ////////////////////////////////////////////////////////////////////////////////// //
