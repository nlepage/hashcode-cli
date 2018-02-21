package upload

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mholt/archiver"
	"github.com/nlepage/hashcode-cli/config"
)

func archiveSource() (string, error) {
	sourceDir, err := filepath.Abs(config.SourceDir())
	if err != nil {
		return "", err
	}

	f, err := ioutil.TempFile(os.TempDir(), "hashcode-source-")
	if err != nil {
		return "", err
	}

	fName := f.Name() + ".tar.gz"

	fmt.Printf("Archiving source dir %s to %s...\n", sourceDir, fName)

	return fName, archiver.TarGz.Make(fName, []string{sourceDir})
}
