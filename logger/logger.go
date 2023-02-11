package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var std = logrus.New()

func SetLevel(level logrus.Level) {
	std.SetLevel(level)
}

func ResetLevel() {
	SetLevel(logrus.InfoLevel)
}

func Get() *logrus.Logger {
	return std
}

func Init(level, file string) {
	var output io.Writer = &lumberjack.Logger{
		Filename:  file,
		MaxSize:   128,
		MaxAge:    90,
		Compress:  true,
		LocalTime: true,
	}

	logLevel := ParseLevel(level, logrus.InfoLevel)
	std = logrus.New()
	std.SetLevel(logLevel)

	if logLevel == logrus.DebugLevel {
		output = os.Stdout
	}

	std.SetFormatter(&logrus.TextFormatter{})
	std.SetOutput(output)
}

func ParseLevel(level string, def logrus.Level) logrus.Level {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return def
	}

	return logLevel
}
