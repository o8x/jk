package uniqid

import (
	"fmt"
	"time"

	"github.com/o8x/jk/v2/djb2"
)

func Number() uint64 {
	return djb2.Sum(fmt.Sprintf("%d", time.Now().UnixNano()))
}

func String() string {
	return djb2.Make(fmt.Sprintf("%d", time.Now().UnixNano()))
}
