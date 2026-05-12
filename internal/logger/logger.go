package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func Init() {
	levelStr := strings.ToLower(os.Getenv("LOG_LEVEL"))
	if levelStr == "" {
		levelStr = "info"
	}

	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		level = logrus.InfoLevel
	}

	Log.SetLevel(level)

	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	Log.Info("logger initialized")
	Log.Infof("log level: %s", level.String())
}
