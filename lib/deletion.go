package lib

import (
	"os"
	"path/filepath"
)

// WipeDir will securely wipe an entire directory. Provide path in the arguments
func WipeDir(path string) error {
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				return WipeFile(path)
			}
			return nil

		})
	if err == nil {
		return os.RemoveAll(path)
	}
	return err
}

// WipeFile will securely wipe a file. Provide path as argument
func WipeFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	size := fileInfo.Size()
	zeroBytes := make([]byte, size)
	copy(zeroBytes[:], "0")
	if _, err = file.Write(zeroBytes); err != nil {
		return err
	}
	_ = file.Close()
	return os.Remove(path)
}
