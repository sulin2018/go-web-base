package log

import (
	"path"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/sulin2018/go-web-base/src/app/config"
	"github.com/sulin2018/go-web-base/src/utils"
)

func InitLogrus() {
	logrus.Trace("init logrus")

	err := utils.MkDir(config.AppConfig.LogFilePath) // create dir
	if err != nil {
		logrus.Panicln(err)
	}

	if config.AppConfig.AppRunMode == "dev" {
		logrus.SetLevel(logrus.TraceLevel) // trace debug info warn error fatal panic
		logrus.SetFormatter(&logrus.TextFormatter{
			// DisableColors: true,
			// ForceQuote:      true,
			FullTimestamp:   true,
			TimestampFormat: utils.TIMEFORMAT,
			// CallerPrettyfier: callerPrettyfier,
		})
	} else {
		logrus.SetLevel(logrus.InfoLevel) // info warn error fatal panic
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: utils.TIMEFORMAT,
		})
	}

	// logrus.SetReportCaller(true)

	logrus.AddHook(newRotateHook(config.AppConfig.LogFilePath, config.AppConfig.AppName, 7*24*time.Hour, 24*time.Hour))

	logrus.Trace("init logrus complate")
	// logrus.WithFields(logrus.Fields{"more": "Init logrus success"}).Info("Info")
}

func newRotateHook(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) *lfshook.LfsHook {
	baseLogPath := path.Join(logPath, logFileName)
	logrus.Println(baseLogPath)
	writer, err := rotatelogs.New(
		baseLogPath+".%Y-%m-%d.log",
		rotatelogs.WithLinkName(baseLogPath),      // link filename to new file
		rotatelogs.WithMaxAge(maxAge),             // 7 days
		rotatelogs.WithRotationTime(rotationTime), // file split
	)
	if err != nil {
		logrus.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}

	return lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.TextFormatter{DisableColors: true, TimestampFormat: "2006-01-02 15:04:05.000"})
}
