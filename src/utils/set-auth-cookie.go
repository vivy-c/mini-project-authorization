package utils

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func SetAuthCookie(echoContext echo.Context, token string) {
	authCookie := http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	}
	echoContext.SetCookie(&authCookie)
}
