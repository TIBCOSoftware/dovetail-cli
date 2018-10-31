package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func CreateDirIfNotExist(subdir ...string) string {
	dir := path.Join(subdir...)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Printf("createDirIfNotExist err %v", err)
			panic(err)
		}
	}

	return dir
}

func CopyFile(src string, dest string) error {
	content, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	return copyContent(content, dest)
}
func copyContent(content []byte, dest string) error {
	ft, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer ft.Close()
	ft.Write(content)
	return nil
}

func CreateTargetDirs(targetPath string) string {
	os.RemoveAll(targetPath)
	target := CreateDirIfNotExist(targetPath)
	return target
}
