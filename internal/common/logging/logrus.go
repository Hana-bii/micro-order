package logging

import (
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

func Init() {
	SetFormatter(logrus.StandardLogger())
	logrus.SetLevel(logrus.DebugLevel)
}

func SetFormatter(logger *logrus.Logger) {
	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyMsg:   "message",
		},
	})
	// 需要在本地环境变量设置LOCAL_MODE
	if isLocal, _ := strconv.ParseBool(os.Getenv("LOCAL_MODE")); isLocal {
		logger.SetFormatter(&prefixed.TextFormatter{
			FullTimestamp:   true,
			ForceFormatting: true,
		})
	}
}
