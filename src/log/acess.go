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

// 访问日志实例
var AccessLoggerClient *logrus.Logger

// 访问日志格式
type AccessFormatter struct{}

// {功能} 访问日志格式方法
// {参数} 无
// {返回} 日志 错误
func (my *AccessFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timeV := time.Now().Format("2006-01-02 15:04:05")
	var (
		ip     string = "0.0.0.0"
		method string = "GET"
		url    string = "/"
	)
	if entry.Data["ip"] != nil {
		ip = entry.Data["ip"].(string)
	}
	if entry.Data["method"] != nil {
		method = entry.Data["method"].(string)
	}
	if entry.Data["url"] != nil {
		url = entry.Data["url"].(string)
	}
	msg := fmt.Sprintf("[%s] [%s] [%s] [%s] %s", timeV, ip, method, url, entry.Message)
	return []byte(msg + "\n"), nil
}

// {功能} 初始化访问日志
// {参数} 无
// {返回} 无
func initAcessLog() {
	//日志实例化
	AccessLoggerClient = logrus.New()
	file := ""           // 日志文件地址
	dir, _ := os.Getwd() // 当前运行目录
	file = dir + "/logs/"
	AccessLoggerClient.SetLevel(logrus.InfoLevel) //设置日志级别
	// 设置日志切割
	logInfoWriter, _ := rotatelogs.New(
		file+"access.%Y-%m-%d.log",
		rotatelogs.WithMaxAge(20*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithRotationSize(5*1024*1024),
	)
	// 设置日志输入
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel: logInfoWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &AccessFormatter{}) // 自定义日志hook
	AccessLoggerClient.AddHook(lfHook)                      // 添加hook
	if isDEV() {
		AccessLoggerClient.SetFormatter(&AccessFormatter{})
		// 添加终端输出
		out := os.Stdout
		AccessLoggerClient.SetOutput(out)
	} else {
		// 生产环境不在终端显示日志
		out, _ := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		writer := bufio.NewWriter(out)
		AccessLoggerClient.SetOutput(writer)
	}
}
