package model

import (
	"pm-ssl-management/src/log"
	"pm-ssl-management/src/util"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 数据库链接实例
var DdClient *gorm.DB

func initTable() {
	DdClient.AutoMigrate(&Certificate{})
}

func initDb() {
	filePath := util.GetWdFile("data.db")

	// 设置mysql日志
	loggerLevel := logger.Silent
	// if util.IsDEV() {
	// 	loggerLevel = logger.Info
	// }

	db, err := gorm.Open(sqlite.Open(filePath), &gorm.Config{
		Logger: logger.Default.LogMode(loggerLevel),
	})

	if err != nil {
		log.ComLoggerFmt("[sqlite] 启动失败: ", err.Error())
		panic("[sqlite] 启动失败: " + err.Error())
	} else {
		log.ComLoggerFmt("[sqlite] 启动成功")
	}

	DdClient = db
}

func Init() {
	initDb()
	initTable()
}
