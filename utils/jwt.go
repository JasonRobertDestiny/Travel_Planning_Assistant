package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 一些常见的JWT错误
var (
	ErrInvalidToken = errors.New("令牌无效")
	ErrExpiredToken = errors.New("令牌已过期")
)

// Claims 定义claims结构体用于JWT
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID int64, username string) (string, error) {
	// 获取配置
	config := GetConfig()

	// 设置过期时间
	expirationTime := time.Now().Add(config.GetJWTExpirationTime())

	// 创建claims
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "traveler-api",
			Subject:   username,
		},
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名token
	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken 验证JWT令牌
func ValidateToken(tokenString string) (*Claims, error) {
	// 获取配置
	config := GetConfig()

	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken 刷新JWT令牌
func RefreshToken(tokenString string) (string, error) {
	// 首先验证当前token
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// 检查过期时间是否在允许的刷新时间范围内
	if claims.ExpiresAt == nil {
		return "", errors.New("invalid token: missing expiration time")
	}

	expirationTime := claims.ExpiresAt.Time

	if time.Until(expirationTime) > 0 {
		// Token还没过期，不需要刷新
		return tokenString, nil
	}

	if time.Since(expirationTime) > time.Hour*24*7 {
		// Token过期超过一周，不允许刷新
		return "", errors.New("token expired too long ago to refresh")
	}

	// 创建新token
	return GenerateToken(claims.UserID, claims.Username)
}
