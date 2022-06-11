package dao

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB
var sqlDb *sql.DB

// 数据库初始化连接
// 调用dao包就会默认执行init函数
func init() {
	user := "root"
	password := "123456"
	host := "localhost"
	dbName := "douyin"

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user, password, host, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}
	sqlDb, err = db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(100)
}

func CloseDB() {
	err := sqlDb.Close()
	if err != nil {
		fmt.Println("closeDb err")
		return
	}
}
