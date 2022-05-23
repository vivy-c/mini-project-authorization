package routes

import (
	"login-api/src/controllers"
	"login-api/src/middlewares"

	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	route := echo.New()
	api := route.Group("/api")

	// user authenticated
	useronly := api.Group("")
	useronly.Use(middlewares.VerifyAuthentication())

	// auths
	api.POST("/login", controllers.LoginHandler)
	api.POST("/register", controllers.RegisterHandler)

	// users
	useronly.GET("/users", controllers.GetAllUsersHandler)
	useronly.GET("/users/d/:id", controllers.GetUserDetailByIDHandler)
	useronly.DELETE("/users/d/:id", controllers.RemoveUserByIDHandler)
	useronly.PUT("/users/p/:id", controllers.AmendProfileByIDHandler)
	useronly.PUT("/users/s/:id", controllers.AmendPasswordByIDHandler)

	return route
}
