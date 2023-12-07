package ini

import (
	"path/filepath"

	"github.com/go-ini/ini"
)

var cfg *ini.File

var DefaultIniConfigFile = "app.ini"

func init() {
	DefaultIniConfigFile, _ = filepath.Abs("app.ini")
	cfg, _ = ini.Load(DefaultIniConfigFile)
}

func Get() *ini.File {
	return cfg
}
