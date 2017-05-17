package repo

import (
	"path"
	"pkg.re/essentialkaos/ek/fmtc"

	"github.com/gongled/vgrepo/prefs"
	"github.com/gongled/vgrepo/meta"
	"pkg.re/essentialkaos/ek/fsutil"
	"strings"
	//"os"
	//fsutil2 "pkg.re/essentialkaos/ek.v8/fsutil"
	//"pkg.re/essentialkaos/ek/hash"
	//"pkg.re/essentialkaos/ek/hash"
	"os"
	//"pkg.re/essentialkaos/ek/hash"
)

const (
	R_PATH_SUFFIX string = "boxes"
	R_HASH_FUNC string   = "sha256"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type VRepository struct {
	Meta *meta.VMetadata
	Preferences *prefs.Preferences
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

func (r *VRepository) Base() string {
	return path.Join(r.Preferences.StoragePath, r.Meta.Name)
}

func (r *VRepository) Dir() string {
	return path.Join(r.Base(), R_PATH_SUFFIX)
}

func (r *VRepository) Path(version string) string {
	return path.Join(r.Dir(), getBoxFormat(r.Meta.Name, version))
}

func (r *VRepository) URL(version string) string {
	return fmtc.Sprintf("%s/%s/%s/%s",
		strings.Trim(r.Preferences.StorageURL, "/"),
		r.Meta.Name,
		R_PATH_SUFFIX,
		getBoxFormat(r.Meta.Name, version),
	)
}

// ////////////////////////////////////////////////////////////////////////////////// //

func (r *VRepository) HasBox(version string) bool {
	return fsutil.IsExist(r.Path(version))
}

func (r *VRepository) IsExist(version string) bool {
	return r.Meta.HasMeta() && r.HasBox(version)
}

func (r *VRepository) AddBox(src string) error {
	var err error

	//if !fsutil.IsExist(src) {
	//	err := fmtc.Errorf("Unable to read file %s\n", src)
	//	return err
	//}

	//if r.Meta.IsEmpty() {
	//	r.Meta.Description = metadata.Description
	//}

	fmtc.Printf("Reading metadata...\n")
	r.Meta.Name = "Hello, world"
	fmtc.Println(r.Meta.Name)

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
	err := os.RemoveAll(r.Base())

	return err
}

// ////////////////////////////////////////////////////////////////////////////////// //

func New(prefs *prefs.Preferences, name string) (*VRepository) {
	r := &VRepository{
		meta.NewMetadata(prefs, name),
		prefs,
	}

	return r
}

// ////////////////////////////////////////////////////////////////////////////////// //
