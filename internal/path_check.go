package internal

import (
	"fmt"
	"os"
	"syscall"
)

/*
check path for writing file process.

	check that exist, its a directory and has correct permissions for writing.
	return error if path is invalid.
*/

func IsWritable(path string) (bool, error) {
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
