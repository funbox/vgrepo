package index

import (
	"fmt"
	"os"
	"io/ioutil"
	"text/template"

	"pkg.re/essentialkaos/ek.v9/fsutil"

	"github.com/gongled/vgrepo/storage"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func ExportIndex(index *storage.VStorage, templateFile string, outputFile string) error {
	if templateFile == "" {
		return fmt.Errorf("Can't use given template")
	}

	if fsutil.IsExist(outputFile) {
		err := os.Remove(outputFile)

		if err != nil {
			return err
		}
	}

	fd, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	defer fd.Close()

	tpl, err := ioutil.ReadFile(templateFile)

	if err != nil {
		return err
	}

	t := template.New("template")
	t, err = t.Parse(string(tpl[:]))

	return t.Execute(fd, index)
}

// ////////////////////////////////////////////////////////////////////////////////// //
