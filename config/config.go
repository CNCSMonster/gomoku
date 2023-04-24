package config

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "github.com/jinzhu/gorm"
)

var (
	db *gorm.DB
)

func Connect() {
	// 从环境变量中获取
	account := os.Getenv("MYSQL_ROOT_ACCOUNT")
	password := os.Getenv("MYSQL_ROOT_PASSWORD")
	dsn := account + ":" + password + "@tcp(43.136.17.142:3306)/gomokudb?charset=utf8mb4&parseTime=True&loc=Local"
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}) //TOfix
	if err != nil {
		panic(err)
	}
	db = d
}

func GetDB() *gorm.DB {
	return db
}
