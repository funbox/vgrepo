package repo

import (
	//"os"
	"strings"

	"pkg.re/essentialkaos/ek.v9/fmtc"
	"pkg.re/essentialkaos/ek.v9/hash"
	"pkg.re/essentialkaos/ek.v9/path"

	"github.com/gongled/vgrepo/meta"
	"pkg.re/essentialkaos/ek/fsutil"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type VRepository struct {
	*meta.VMetadata
}

type VRepositoryList []*VRepository

// ////////////////////////////////////////////////////////////////////////////////// //

func (r *VRepository) BaseRepo() string {
	return r.StoragePath
}

func (r *VRepository) DirRepo() string {
	return path.Join(r.BaseRepo(), "packages")
}

func (r *VRepository) PathRepo(pkg *VPackage) string {
	return path.Join(r.DirRepo(), pkg.PathBoxFormat())
}

func (r *VRepository) URLRepo(pkg *VPackage) string {
	return fmtc.Sprintf("%s/%s",
		strings.Trim(r.StorageURL, "/"),
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

	return err
}

func (r *VRepository) AddPackage(src string, pkg *VPackage) error {
	r.CreateRepo()

	// dst := r.PathRepo(pkg)
	dst := src

	if !fsutil.IsExist(src) {
		return fmtc.Errorf("File %s does not exist", src)
	}

	var err error
	//err = r.copyPackage(src, dst)
	//
	//if err != nil {
	//	return fmtc.Errorf(
	//		"Unable to copy package from %s to %s",
	//		src,
	//		dst,
	//	)
	//}

	providerList := make(meta.VMetadataProvidersList, 0)

	provider := meta.NewMetadataProvider(
		pkg.Provider,
		hash.FileHash(dst),
		"sha256",
		r.URLRepo(pkg),
	)

	providerList = append(providerList, provider)

	version := meta.NewMetadataVersion(pkg.Version, providerList)

	r.AddVersion(version)
	// r.WriteMeta()

	for _, d := range r.Versions {
		fmtc.Println(d.Version)
		for _, q := range d.Providers {
			fmtc.Println(q.Name, q.Checksum, q.ChecksumType, q.URL)
		}
	}

	return err
}

func (r *VRepository) ListPackages() *VPackageList {
	return nil
}

func (r *VRepository) RemovePackage(pkg *VPackage) error {
	// err := os.RemoveAll(r.PathRepo(pkg))
	fmtc.Println(r.PathRepo(pkg))

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

func NewRepository(storagePath string, storageUrl string, name string) *VRepository {
	m := meta.NewMetadata(
		storagePath,
		storageUrl,
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
