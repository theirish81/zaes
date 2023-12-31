package lib

import (
	"os"
	"path"
	"strings"
)

// checkFileExistenceAndType checks whether the provided path leads to a file/dir, and it is the expected type as
// expressed by the "dir" parameter
func checkFileExistenceAndType(file string, dir bool) bool {
	info, err := os.Stat(file)
	return err != nil || (dir && info.IsDir()) || (!dir && !info.IsDir())
}

// checkExists will return true if the provided file path leads to an existing file
func checkExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

// zipFilePath given a parent directory and a file name, it generates a ZIP file name for it
func zipFilePath(parentDir string, fileName string) string {
	return path.Join(parentDir, fileName+".zip")
}

// zaesFilePath given a parent directory and a file name, it generates a ZAES file name for it
func zaesFilePath(parentDir string, fileName string) string {
	return path.Join(parentDir, fileName+".zaes")
}

// cleanupDirPath checks whether the provided path ends with a path separator and removes it
func cleanupDirPath(path string) string {
	if strings.HasSuffix(path, string(os.PathSeparator)) {
		return path[0 : len(path)-1]
	}
	return path
}
