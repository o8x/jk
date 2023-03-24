package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var std = New(logrus.DebugLevel, os.Stdout)

func Get() *logrus.Logger {
	return std
}

func UseLogger(l *logrus.Logger) {
	std = l
}

func UseDefault() {
	UseLogger(NewFile("info", "log/jk.log"))
}

func NewFile(level string, out string) *logrus.Logger {
	l := ParseLevel(level, logrus.InfoLevel)
	return New(l, NewRotater(out))
}

func New(level logrus.Level, out io.Writer) *logrus.Logger {
	l := logrus.New()
	l.SetLevel(level)
	l.SetFormatter(&logrus.TextFormatter{})
	l.SetOutput(out)
	return l
}

func NewRotater(file string) io.Writer {
	return &lumberjack.Logger{
		Filename:  file,
		MaxSize:   128,
		MaxAge:    90,
		Compress:  true,
		LocalTime: true,
	}
}

func ParseLevel(level string, def logrus.Level) logrus.Level {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return def
	}

	return logLevel
}
