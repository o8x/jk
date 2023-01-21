package logger

import (
	"strings"

	"github.com/sirupsen/logrus"
)

func Info(format string, args ...any) {
	std.Logf(logrus.InfoLevel, format, args...)
}

func Fatal(format string, args ...any) {
	std.Fatal([]any{format, args}...)
}

func Warn(format string, args ...any) {
	std.Logf(logrus.WarnLevel, format, args...)
}

func Error(format string, args ...any) {
	std.Logf(logrus.ErrorLevel, format, args...)
}

func Group(group string) *logrus.Entry {
	return WithField("group", strings.ToUpper(group))
}

func WithField(key string, value any) *logrus.Entry {
	return std.WithField(key, value)
}

func WithError(err error) *logrus.Entry {
	return std.WithError(err)
}
