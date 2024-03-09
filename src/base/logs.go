package base

import (
	"github.com/sirupsen/logrus"
	"time"
)

type ExtraFieldsHook struct {
	logrus.Hook
	AppName string
}

func (hook ExtraFieldsHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook ExtraFieldsHook) Fire(entry *logrus.Entry) error {
	entry.Data["timestamp"] = time.Now().Format(time.StampMilli)
	entry.Data["app_name"] = hook.AppName

	return nil
}

var Logger = CreateLogger(nil)

func CreateLogger(config *BackendConfig) *logrus.Logger {
	var level logrus.Level
	var hook ExtraFieldsHook
	if config != nil {
		var err error
		level, err = logrus.ParseLevel(config.Logs.Level)
		if err != nil {
			panic(err)
		}
		hook = ExtraFieldsHook{
			AppName: config.Logs.AppName,
		}
	} else {
		level = logrus.InfoLevel
		hook = ExtraFieldsHook{
			AppName: "stealthy-backend",
		}
	}
	formatter := logrus.JSONFormatter{}
	formatter.DisableTimestamp = true

	logger := logrus.New()
	logger.SetFormatter(&formatter)
	logger.SetLevel(level)
	logger.AddHook(&hook)
	return logger
}
