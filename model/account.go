package model

import (
	"cncsmonster/gomoku/config"
	"errors"

	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Name     string `gorm:"column:name"`
	Password string `gorm:"column:password"`
}

var accountDB *gorm.DB

// 该函数用来连接数据库，
func init() {
	// 连接数据库,然后从数据库中获取当前最新id
	config.Connect()
	accountDB = config.GetDB()
	accountDB.AutoMigrate(&Account{})
}

func CreateAccount(name, password string) (*Account, error) {

	findAccount := FindAccountByName(name)
	if findAccount != nil {
		return nil, errors.New("Account already exists")
	}
	var out Account
	result := accountDB.Create(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	out.Name = name
	out.Password = password
	accountDB.Save(&out)
	return &out, nil
}

// 通过名字查找账号,如果不存在返回nil
func FindAccountByName(name string) *Account {
	var findAccount Account
	result := accountDB.Where("name=?", name).Find(&findAccount)
	if result.RowsAffected == 0 {
		return nil
	}
	return &findAccount
}
