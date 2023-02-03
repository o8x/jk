package size

import "fmt"

func Beautify(size int64, length int) string {
	f := float64(size)

	if f < 1000 {
		return fmt.Sprintf("%d", size)
	}

	if f >= 1000 && f < 1000000 {
		format := fmt.Sprintf("%%.%dfKB", length)

		return fmt.Sprintf(format, f/1000)
	}

	if f >= 1000000 && f < 1000000000 {
		format := fmt.Sprintf("%%.%dfMB", length)
		return fmt.Sprintf(format, f/1000000)
	}

	format := fmt.Sprintf("%%.%dfGB", length)
	return fmt.Sprintf(format, f/1000000000)
}
