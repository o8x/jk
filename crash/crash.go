package crash

import (
	"runtime/debug"

	"github.com/o8x/jk/logger"
)

func Recover(message string) {
	if err := recover(); err != nil {
		logger.
			WithField("recover", err).
			WithField("stack", string(debug.Stack())).
			Error(message)
	}
}
