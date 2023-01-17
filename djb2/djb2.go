package djb2

import "fmt"

func Make(s string) string {
	return fmt.Sprintf("%d", Sum(s))
}

// Sum http://www.cse.yorku.ca/~oz/hash.html
// https://github.com/singularian/mdhash
func Sum(s string) uint64 {
	var hash uint64 = 5381
	for _, c := range s {
		hash = ((hash << 5) + hash) + uint64(c)
	}

	return hash
}

func Check(s string, hash string) bool {
	return Make(s) == hash
}
