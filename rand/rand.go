package rand

import (
	"math/rand"
	"time"
)

var rnd *rand.Rand

func init() {
	rnd = rand.New(rand.NewSource(time.Now().UnixMilli()))
}

func Intn(n int) int {
	return rnd.Intn(n)
}

func Int() int {
	return rnd.Int()
}

func Int63() int64 {
	return rnd.Int63()
}

func Int63n(n int64) int64 {
	return rnd.Int63n(n)
}
