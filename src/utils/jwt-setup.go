package utils

import (
	"log"
	"login-api/src/configs"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	jwt.StandardClaims
}

func GenerateJwt(userId string) (token string, err error) {
	claims := JwtCustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			Id:        userId,
		},
	}

	secret, _ := configs.LoadServerConfig(".")
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = rawToken.SignedString([]byte(secret.JWTsecret))
	return
}

func SetJwtCookie(ec echo.Context, token string) {
	authCookie := http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	}
	ec.SetCookie(&authCookie)
}

func ExtractClaims(tokenStr string) (jwt.MapClaims, bool) {
	hmacSecretString := "configs.GetJwtSecret().SecretKey"
	hmacSecret := []byte(hmacSecretString)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return hmacSecret, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}
