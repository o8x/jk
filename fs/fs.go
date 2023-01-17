package fs

import (
	"errors"
	"io/fs"
	"os"
)

func FileExist(name string) bool {
	_, err := os.Stat(name)
	return !errors.Is(err, fs.ErrNotExist)
}

func HasRDPermission(name string) bool {
	if !FileExist(name) {
		return false
	}

	f, err := os.OpenFile(name, os.O_RDONLY, 0755)
	if err == nil {
		f.Close()
		return true
	}

	return !errors.Is(err, fs.ErrPermission)
}

func HasWRPermission(name string) bool {
	if !FileExist(name) {
		return false
	}

	f, err := os.OpenFile(name, os.O_WRONLY, 0755)
	if err == nil {
		f.Close()
		return true
	}

	return !errors.Is(err, fs.ErrPermission)
}

func HasRWPermission(name string) bool {
	if !FileExist(name) {
		return false
	}

	f, err := os.OpenFile(name, os.O_RDONLY|os.O_WRONLY, 0755)
	if err == nil {
		f.Close()
		return true
	}

	return !errors.Is(err, fs.ErrPermission)
}
