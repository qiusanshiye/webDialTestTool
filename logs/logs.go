package logs

import (
	"github.com/alecthomas/log4go"
)

type TLoggerConf struct {
	LogPath  string `ini:"file"`
	LogLevel int    `ini:"level"`
}

const (
	logFilterName = "file"
)

var (
	logger log4go.Logger
)

func Init(lc *TLoggerConf) {
	logger = make(log4go.Logger)
	if lc != nil && lc.LogPath != "stdout" {
		initFileLogger(lc)
	} else {
		initConsoleLogger()
	}
	logger.Info("init logger...Done")
}

func initFileLogger(lc *TLoggerConf) {
	flw := log4go.NewFileLogWriter(lc.LogPath, false)
	flw.SetFormat("[%D %T] [%L] (%S) %M")
	flw.SetRotate(true)
	flw.SetRotateSize(0)
	flw.SetRotateLines(0)
	flw.SetRotateDaily(true)
	logger.AddFilter(logFilterName, log4go.Level(lc.LogLevel), flw)
}

func initConsoleLogger() {
	w := log4go.NewConsoleLogWriter()
	w.SetFormat("[%D %T] [%L] (%S) %M")
	logger.AddFilter(logFilterName, log4go.DEBUG, w)
}

func Uninit() {
	logger[logFilterName].Close()
	logger.Info("uninit logger...Done")
}

func Debug(format string, args ...interface{}) {
	if len(args) > 0 {
		logger.Debug(format, args...)
	} else {
		logger.Debug(format)
	}
}

func Info(format string, args ...interface{}) {
	if len(args) > 0 {
		logger.Info(format, args...)
	} else {
		logger.Info(format)
	}
}

func Warn(format string, args ...interface{}) {
	if len(args) > 0 {
		logger.Warn(format, args...)
	} else {
		logger.Warn(format)
	}
}

func Error(format string, args ...interface{}) {
	if len(args) > 0 {
		logger.Error(format, args...)
	} else {
		logger.Error(format)
	}
}

func Fatal(format string, args ...interface{}) {
	if len(args) > 0 {
		logger.Critical(format, args...)
	} else {
		logger.Critical(format)
	}
}
