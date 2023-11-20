package version

import (
	"encoding/base64"
	"fmt"
	"runtime"
	"strings"

	"github.com/o8x/jk/v2/json"
)

var date string
var copyright string
var version string
var changelog string
var hash string

type Version struct {
	Copyright  string `json:"copyright"`
	Version    string `json:"version"`
	GoVersion  string `json:"go_version"`
	Changelog  string `json:"changelog"`
	CommitHash string `json:"hash"`
	Date       string `json:"date"`
	Compiler   string `json:"compiler"`
	Platform   string `json:"platform"`
}

func (v Version) String() string {
	return json.PrettifyString(v)
}

func GetBuildFlags() string {
	return strings.Join([]string{
		"-ldflags", `"`,
		fmt.Sprintf("-X github.com/o8x/jk/v2/version.date=%s \\\n", date),
		fmt.Sprintf("-X github.com/o8x/jk/v2/version.copyright=%s \\\n", copyright),
		fmt.Sprintf("-X github.com/o8x/jk/v2/version.version=%s \\\n", version),
		fmt.Sprintf("-X github.com/o8x/jk/v2/version.changelog=%s \\\n", changelog),
		fmt.Sprintf("-X github.com/o8x/jk/v2/version.hash=%s", hash),
		`"`,
	}, " ")
}

func Get() *Version {
	s, err := base64.StdEncoding.DecodeString(copyright)
	if err == nil {
		copyright = string(s)
	}

	if s, err = base64.StdEncoding.DecodeString(changelog); err == nil {
		changelog = string(s)
	}

	return &Version{
		Copyright:  copyright,
		Version:    version,
		Changelog:  changelog,
		CommitHash: hash,
		Date:       date,
		GoVersion:  runtime.Version(),
		Platform:   fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		Compiler:   runtime.Compiler,
	}
}
