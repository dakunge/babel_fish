package util

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logc"
)

var jwtKey = []byte("my_secret_key")

type MyCustomClaims struct {
	Username string `json:"username"`
	UserID   uint   `json:"user_id"`
	jwt.StandardClaims
}

func GenerateJWT(ctx context.Context, username string, userID uint, duration time.Duration) (string, error) {
	// 设置自定义声明
	claims := MyCustomClaims{
		Username: username,
		UserID:   userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
	}

	// 创建 JWT 令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并获取完整的编码后的字符串
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		logc.Errorf(ctx, "jwt sign err %v", err)
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(tokenString string) (*MyCustomClaims, error) {
	// 解析 JWT 令牌
	claims := MyCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return &claims, nil
}
