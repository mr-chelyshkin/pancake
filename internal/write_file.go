package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
)

/*
write file to local.

	check path for writing file process.
	(check that exist, its a directory and has correct permissions for writing)
*/

func WriteFile(filePath string, dataBytes []byte) error {
	fileDir, _ := filepath.Split(filePath)
	if fileDir == "" {
		fileDir = "./"
	}

	if ok, err := __isWritable__(fileDir); !ok {
		return err
	}

	file, err := os.Create(filePath)
	if file == nil || err != nil {
		return err
	}
	defer file.Close()

	if err := ioutil.WriteFile(file.Name(), dataBytes, 0644); err != nil {
		return err
	}
	return nil
}

// --- >
func __isWritable__(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf(fmt.Sprintf("%s doesn't exist", path))
	}

	if !info.IsDir() {
		return false, fmt.Errorf(fmt.Sprintf("%s isn't a directory", path))
	}

	if info.Mode().Perm() & (1 << (uint(7))) == 0 {
		return false, fmt.Errorf(fmt.Sprintf("%s write permission bit is not set on this file for user", path))
	}

	var stat syscall.Stat_t
	if err = syscall.Stat(path, &stat); err != nil {
		return false, fmt.Errorf(fmt.Sprintf("%s unable to get stat", path))
	}

	return true, nil
}
