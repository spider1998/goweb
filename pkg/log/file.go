package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	LineFeed = "\r\n"
)

//WriteLog return error
func writeLog(path, msg string) {
	r, _ := regexp.Compile(`\d{4}-\d{2}-\d{2}`)
	times := r.FindString(msg)

	var (
		err error
		f   *os.File
	)
	path = path + strings.Split(msg, ":")[0]
	if !IsExist(path) {
		if err = CreateDir(path); err != nil {
			logf(os.Stderr, priorityFatal, err.Error())
		}
	}

	f, err = os.OpenFile(path+"/"+times+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	_, err = io.WriteString(f, LineFeed+msg)

	defer f.Close()
	return
}

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

// 获取指定路径下以及所有子目录下的所有文件，并且删除大于 suffix 的文件
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
