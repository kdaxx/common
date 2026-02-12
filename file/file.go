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

	// 若目录不存在，创建该路径下的所有目录
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
