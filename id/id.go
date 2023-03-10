package id

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"github.com/o8x/jk/v2/fs"
	"github.com/o8x/jk/v2/logger"
	"github.com/o8x/jk/v2/uniqid"
)

var FileName = "/var/.jk-u"

func init() {
	dir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	FileName = path.Join(dir, ".jk-u")
	logger.WithField("name", FileName).Info("with jk-u file")
}

func Generate() (string, error) {
	if fs.FileExist(FileName) {
		c, err := os.ReadFile(FileName)
		if err != nil {
			return "", err
		}

		if len(c) > 4 && string(c[:4]) == "jk-u" {
			return Read()
		}
	}

	number := uniqid.String()
	err := os.WriteFile(FileName, []byte(fmt.Sprintf("jk-u@%s", number)), 0755)
	return number, err
}

func Cleanup() {
	_ = os.Remove(FileName)
}

func Read() (string, error) {
	c, err := os.ReadFile(FileName)
	if err != nil {
		return "", err
	}

	return string(bytes.Split(c, []byte("@"))[1]), nil
}
