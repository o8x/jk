package crash

import (
	"runtime/debug"

	"github.com/sirupsen/logrus"

	"github.com/o8x/jk/logger"
)

type RecoverFunc func(*Crash)

func Logger(log *logrus.Logger) func(c *Crash) {
	return func(c *Crash) {
		c.Logger = log
	}
}

type Crash struct {
	Logger  *logrus.Logger `json:"logger"`
	Message string         `json:"message"`
}

func Recover(message string, fns ...RecoverFunc) {
	c := &Crash{
		Logger:  logger.Get(),
		Message: message,
	}

	for _, fn := range fns {
		fn(c)
	}

	if err := recover(); err != nil {
		c.Logger.
			WithField("recover", err).
			WithField("stack", string(debug.Stack())).
			Error(message)
	}
}
