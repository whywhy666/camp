package src

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DBHelper *gorm.DB
var err error

func init() {
	dsn := "root:bytedancecamp@tcp(180.184.74.5:3306)/?charset=utf8mb4&parseTime=True&loc=Local"
	DBHelper, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		//fmt.Println(err)
		log.Fatal("DB初始化错误:", err)

		//出错之后发送信号
	}
	DBHelper = DBHelper.Debug()

	//链接池问题

}
