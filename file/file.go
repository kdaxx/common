package file

import (
	"os"
	"path/filepath"
)

func MkdirWithFilePath(path string) error {
	dir, _ := filepath.Split(path)

	return Mkdir(dir)
}

func Mkdir(dir string) error {
	_, err := os.Stat(dir)

	// create all dir if it's not exist
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
