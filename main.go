package main

import (
	"login-api/src/database"
	"login-api/src/routes"
)

func init() {
	database.Migrate()
}

func main() {
	api := routes.New()
	api.Logger.Fatal(api.Start(":8000"))
}
