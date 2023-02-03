package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
)

func Sha1(text string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(text)))
}

func Sha256(text string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(text)))
}

func Sha224(text string) string {
	return fmt.Sprintf("%x", sha256.Sum224([]byte(text)))
}

func Md5(text string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(text)))
}
