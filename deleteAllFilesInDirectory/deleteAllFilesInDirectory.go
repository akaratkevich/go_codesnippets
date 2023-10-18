package deleteAllFilesInDirectory

import (
	"os"
	"path/filepath"
)

func DeleteAllFilesInDirectory(directory string) error {
	d, err := os.Open(directory)
	if err != nil {
		return err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		err = os.RemoveAll(filepath.Join(directory, name))
		if err != nil {
			return err
		}
	}

	return nil
}
