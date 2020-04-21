package log

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	os.Chmod(path, os.ModePerm)
	return nil
}

func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

func remove(file string) error {
	// 删除文件
	err := os.Remove(file)
	if err != nil {
		return err
	}
	return nil
}

func walkDir(dir string, suffix int64, f func(string) error) (err error) {
	err = filepath.Walk(dir, func(fname string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			//忽略目录
			return nil
		}
		if fi.Size()/(1024*1024) > suffix {
			fmt.Println(fi.Size() / (1024 * 1024))
			return remove(fname)
		}
		return nil
	})
	return err
}
