package repo

import (
	"os"
	"path"
	"strings"

	"pkg.re/essentialkaos/ek.v8/fmtc"
	"pkg.re/essentialkaos/ek.v8/fsutil"

	"github.com/gongled/vgrepo/meta"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type VRepository struct {
	*meta.VMetadata
}

// ////////////////////////////////////////////////////////////////////////////////// //

func getVersion(version string) string {
	ver := "latest"

	if version != "" {
		ver = version
	}

	return ver
}

func getBoxFormat(name string, version string) string {
	return fmtc.Sprintf("%s-%s.box", name, getVersion(version))
}

func (r *VRepository) BaseRepo() string {
	return path.Join(r.StoragePath, r.Name)
}

func (r *VRepository) DirRepo() string {
	return path.Join(r.BaseRepo(), "boxes")
}

func (r *VRepository) PathRepo(version string) string {
	return path.Join(r.DirRepo(), getBoxFormat(r.Name, version))
}

func (r *VRepository) URLRepo(version string) string {
	return fmtc.Sprintf("%s/%s/%s/%s",
		strings.Trim(r.StorageURL, "/"),
		r.Name,
		"boxes",
		getBoxFormat(r.Name, version),
	)
}

// ////////////////////////////////////////////////////////////////////////////////// //

func (r *VRepository) HasBox(version string) bool {
	return fsutil.IsExist(r.PathRepo(version))
}

func (r *VRepository) IsExist(version string) bool {
	return r.HasMeta() && r.HasBox(version)
}

func (r *VRepository) AddBox(src string) error {
	var err error

	//if !fsutil.IsExist(src) {
	//	err := fmtc.Errorf("Unable to read file %s\n", src)
	//	return err
	//}

	//if r.Meta.IsEmptyMeta() {
	//	r.Meta.Description = metadata.Description
	//}

	fmtc.Printf("Reading metadata...\n")
	r.Name = "qweqwe"
	fmtc.Println(r.Name)

	//if fsutil.IsExist(ver) {
	//	err := fmtc.Errorf("File %s is already exist", ver)
	//	return err
	//}
	//
	//if !fsutil.IsDir(basedir) {
	//	err := os.MkdirAll(basedir, 0755)
	//	return err
	//}

	//err = fsutil2.CopyFile(src, dst)

	return err
}

func (r *VRepository) DeleteBox(version string) error {
	ver := getVersion(version)

	fmtc.Println(ver)

	return nil
}

func (r *VRepository) Destroy() error {
	err := os.RemoveAll(r.BaseRepo())

	return err
}

// ////////////////////////////////////////////////////////////////////////////////// //

func NewRepository(storagePath string, storageUrl string, name string) *VRepository {
	r := meta.NewMetadata(
		storagePath,
		storageUrl,
		name,
		"",
		make([]*meta.VMetadataVersion, 0))

	m := &VRepository{r}

	if m.HasMeta() {
		// TODO: remove returning value *VMetadataRepository
		m.VMetadata, _ = m.ReadMeta()
	}

	return m
}

// ////////////////////////////////////////////////////////////////////////////////// //
