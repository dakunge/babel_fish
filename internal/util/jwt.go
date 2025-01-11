package util

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("my_secret_key")

type MyCustomClaims struct {
	Username string `json:"username"`
	UserID   uint   `json:"user_id"`
	jwt.StandardClaims
}

func GenerateJWT(ctx context.Context, secretKey string, userID uint, iat, seconds int64) (string, error) {
	// 设置自定义声明
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userID"] = userID
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
