package error

import "github.com/o8x/jk/v2/logger"

func Hide(err error) {
	if err != nil {
		logger.WithError(err).Error("hide error")
	}
}
