package x

import (
	"strconv"
)

func ParseInt(s string, def int) int {
	return int(ParseInt64(s, int64(def)))
}
func ParseUint(s string, def uint) uint {
	return uint(ParseInt64(s, int64(def)))
}

func ParseInt64(s string, def int64) int64 {
	l, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return def
	}

	return l
}
