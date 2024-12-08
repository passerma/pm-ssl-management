package log

import "os"

func isDEV() bool {
	return os.Getenv("env") == "dev"
}

// {功能} log 初始化
// {参数} 无
// {返回} 无
func init() {
	initLog()
	initAcessLog()
}
