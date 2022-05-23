package database

import (
	"login-api/src/configs"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	config, _ := configs.LoadServerConfig(".")
	db, err := gorm.Open(mysql.Open(config.ConnectionString), &gorm.Config{})

	if err != nil {
		panic("connection to database error!")
	}
	return db
}
