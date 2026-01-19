package jwt

import (
	"errors"
	"time"

	"gin-admin-pro/internal/pkg/config"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 声明
type Claims struct {
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID uint, username string) (string, error) {
	cfg := config.GetConfig()

	// 创建过期时间
	accessExpire := time.Duration(cfg.JWT.AccessTokenExpire) * 24 * time.Hour
	refreshExpire := time.Duration(cfg.JWT.RefreshTokenExpire) * 24 * time.Hour

	// 创建访问 Token 声明
	accessClaims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "gin-admin-pro",
			Subject:   "access-token",
		},
	}

	// 生成访问 Token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", err
	}

	// 创建刷新 Token 声明
	refreshClaims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "gin-admin-pro",
			Subject:   "refresh-token",
		},
	}

	// 生成刷新 Token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	_, err = refreshToken.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", err
	}

	// 返回访问 Token（前端主要使用）
	return accessTokenString, nil
}

// ParseToken 解析 Token
func ParseToken(tokenString string) (*Claims, error) {
	cfg := config.GetConfig()

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的 token")
}

// RefreshToken 刷新 Token
func RefreshToken(refreshTokenString string) (string, error) {
	cfg := config.GetConfig()

	token, err := jwt.ParseWithClaims(refreshTokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", errors.New("无效的刷新 token")
	}

	// 检查是否是刷新 Token
	if claims.Subject != "refresh-token" {
		return "", errors.New("不是刷新 token")
	}

	// 生成新的访问 Token
	return GenerateToken(claims.UserID, claims.Username)
}

// ValidateToken 验证 Token 是否有效
func ValidateToken(tokenString string) bool {
	_, err := ParseToken(tokenString)
	return err == nil
}
