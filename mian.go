package main

import (
	"pm-ssl-management/src/conf"
	"pm-ssl-management/src/cron"
	"pm-ssl-management/src/log"
	"pm-ssl-management/src/model"
	"pm-ssl-management/src/route"
)

func main() {
	log.ComLoggerFmt("启动服务: ", conf.GetConf("name"))

	model.Init()

	cron.Init()

	route.Init()
}
