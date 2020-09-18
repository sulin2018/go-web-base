package log

import (
	"path"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	logs "github.com/sirupsen/logrus"
	"github.com/sulin2018/go-web-base/src/app/config"
	"github.com/sulin2018/go-web-base/src/utils"
)

func InitLogrus() {
	logs.Trace("init logrus")

	err := utils.MkDir("logs") // create dir
	if err != nil {
		logs.Panicln(err)
	}

	if config.AppConfig.AppRunMode == "dev" {
		logs.SetLevel(logs.TraceLevel) // trace debug info warn error fatal panic
		logs.SetFormatter(&logs.TextFormatter{
			// DisableColors: true,
			// ForceQuote:      true,
			FullTimestamp:   true,
			TimestampFormat: utils.TIMEFORMAT,
			// CallerPrettyfier: callerPrettyfier,
		})
	} else {
		logs.SetLevel(logs.InfoLevel) // info warn error fatal panic
		logs.SetFormatter(&logs.JSONFormatter{
			TimestampFormat: utils.TIMEFORMAT,
		})
	}

	// logs.SetReportCaller(true)

	logs.AddHook(newRotateHook(config.AppConfig.LogFilePath, config.AppConfig.AppName, 7*24*time.Hour, 24*time.Hour))

	logs.Trace("init logrus complate")
	// logs.WithFields(logs.Fields{"more": "Init logrus success"}).Info("Info")
}

func newRotateHook(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) *lfshook.LfsHook {
	baseLogPath := path.Join(logPath, logFileName)
	logs.Println(baseLogPath)
	writer, err := rotatelogs.New(
		baseLogPath+".%Y-%m-%d.log",
		rotatelogs.WithLinkName(baseLogPath),      // link filename to new file
		rotatelogs.WithMaxAge(maxAge),             // 7 days
		rotatelogs.WithRotationTime(rotationTime), // file split
	)
	if err != nil {
		logs.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}

	return lfshook.NewHook(lfshook.WriterMap{
		logs.DebugLevel: writer,
		logs.InfoLevel:  writer,
		logs.WarnLevel:  writer,
		logs.ErrorLevel: writer,
		logs.FatalLevel: writer,
		logs.PanicLevel: writer,
	}, &logs.TextFormatter{DisableColors: true, TimestampFormat: "2006-01-02 15:04:05.000"})
}
