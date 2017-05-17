package prefs

import (
	"pkg.re/essentialkaos/ek.v8/fmtc"
	"pkg.re/essentialkaos/ek.v8/fsutil"
	"pkg.re/essentialkaos/ek.v8/knf"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	STORAGE_URL  = "storage:url"
	STORAGE_PATH = "storage:path"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type Preferences struct {
	StorageURL  string
	StoragePath string
}

// ////////////////////////////////////////////////////////////////////////////////// //

func (p *Preferences) Validate() []error {
	var errs []error

	if !fsutil.IsExist(p.StoragePath) {
		errs = append(errs, fmtc.Errorf("Storage path %s does not exist\n", p.StoragePath))
	}

	if !fsutil.IsWritable(p.StoragePath) {
		errs = append(errs, fmtc.Errorf("Storage path %s is not writable\n", p.StoragePath))
	}

	return errs
}

func New(configPath string) (*Preferences, error) {
	cnf, err := knf.Read(configPath)

	if err != nil {
		return nil, fmtc.Errorf("Unable to read settings")
	}

	p := &Preferences{
		StorageURL:  cnf.GetS(STORAGE_URL),
		StoragePath: cnf.GetS(STORAGE_PATH),
	}

	return p, err
}

// ////////////////////////////////////////////////////////////////////////////////// //
