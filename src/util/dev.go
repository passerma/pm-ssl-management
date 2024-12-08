package util

import "os"

// {功能} 判断是否为开发环境
// {参数} 无
// {返回} 是否为开发环境
func IsDEV() bool {
	return os.Getenv("env") == "dev"
}
