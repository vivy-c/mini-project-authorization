package database

import (
	userModel "login-api/src/models/user"
)

func Migrate() {
	db := Connect()
	db.AutoMigrate(&userModel.User{})
}

func Demigrate() {
	db := Connect()
	db.Migrator().DropTable(&userModel.User{})
}
