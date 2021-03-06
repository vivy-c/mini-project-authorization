package middlewares

import (
	"login-api/src/configs"
	"login-api/src/utils"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func VerifyAuthentication() echo.MiddlewareFunc {
	cfg, _ := configs.LoadServerConfig(".")
	return middleware.JWTWithConfig(middleware.JWTConfig{
		ErrorHandlerWithContext: func(err error, c echo.Context) error {
			status := http.StatusUnauthorized
			return utils.CreateResponse(c, status, http.StatusText(status), err.Error())
		},
		SigningKey:  []byte(cfg.JWTsecret),
		ContextKey:  "token",
		Claims:      jwt.MapClaims{},
		TokenLookup: "cookie:token",
	})
}
