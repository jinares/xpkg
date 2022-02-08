package xlog

import (
	"github.com/sirupsen/logrus"
	"os"
)

//Std Std log
func Std() *logrus.Logger {
	std := logrus.New()
	std.SetOutput(os.Stderr)
	std.SetLevel(logrus.TraceLevel)
	std.SetFormatter(&TextLoggerFormatter{})
	std.SetReportCaller(true)
	return std
}
