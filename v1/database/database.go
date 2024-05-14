package database

import (
	"backend/v1/config"
	"backend/v1/model"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error

	dbConfig := config.Config{}
	config.LoadConfig(&dbConfig)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbConfig.DBUser, dbConfig.DBPass, dbConfig.DBHost, dbConfig.DBPort, dbConfig.DBName)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Kết nối database thất bại: " + err.Error())
	}

	err = DB.AutoMigrate(&model.User{}, model.Contact{})
	if err != nil {
		panic("Migration thất bại: " + err.Error())
	}

	fmt.Println("Kết nối database thành công ^-^")
}
