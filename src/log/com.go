package log

import (
	"bufio"
	"fmt"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// 通用日志格式
type ComLogFormatter struct{}

// 通用日志实例
var ComLoggerClient *logrus.Logger

// {功能} 通用日志格式方法
// {参数} 无
// {返回} 日志 错误
func (my *ComLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timeV := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf("[%s] %s", timeV, entry.Message)
	return []byte(msg + "\n"), nil
}

// {功能} 初始化日志
// {参数} 无
// {返回} 无
func initLog() {
	//日志实例化
	ComLoggerClient = logrus.New()
	file := ""
	dir, _ := os.Getwd()
	file = dir + "/logs/"
	//设置日志级别
	ComLoggerClient.SetLevel(logrus.InfoLevel)
	//设置Info日志切割
	logInfoWriter, _ := rotatelogs.New(
		file+"info.%Y-%m-%d.log",
		rotatelogs.WithMaxAge(20*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithRotationSize(5*1024*1024),
	)
	//设置error日志切割
	logErrorWriter, _ := rotatelogs.New(
		file+"error.%Y-%m-%d.log",
		rotatelogs.WithMaxAge(20*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithRotationSize(5*1024*1024),
	)
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logInfoWriter,
		logrus.ErrorLevel: logErrorWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &ComLogFormatter{})
	ComLoggerClient.AddHook(lfHook)

	if isDEV() {
		out := os.Stdout
		ComLoggerClient.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			PadLevelText:    true,
		})
		ComLoggerClient.SetOutput(out)
	} else {
		out, _ := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		writer := bufio.NewWriter(out)
		ComLoggerClient.SetOutput(writer)
	}
}

// {功能} 打印日志输出，仅控制台输出
// {参数} s 日志信息
// {返回} 无
func ComLoggerFmt(s ...interface{}) {
	timeV := time.Now().Format("2006-01-02 15:04:05")
	timeV = fmt.Sprintf("[%s] ", timeV)
	msg := append([]interface{}{timeV}, s...)
	msg = append(msg, "\n")
	fmt.Print(msg...)
}
