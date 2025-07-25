package middleware

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/warjiang/page-spy-api/config"
)

// 用于签名JWT的密钥
var jwtSecret []byte

// Claims JWT声明结构
type Claims struct {
	jwt.RegisteredClaims
}

// InitJWTSecret 初始化JWT密钥
func InitJWTSecret(cfg *config.Config) {
	if cfg.AuthConfig == nil || cfg.AuthConfig.JwtSecret == "" {
		// 使用临时密钥，但不保存到配置文件
		jwtSecret = generateRandomKey(32)
		return
	}

	jwtSecret = []byte(cfg.AuthConfig.JwtSecret)
}

// 生成随机密钥
func generateRandomKey(length int) []byte {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		// 如果随机生成失败，使用时间戳作为备选方案
		now := time.Now().UnixNano()
		for i := 0; i < length; i++ {
			key[i] = byte((now >> (i % 8)) & 0xff)
		}
	}
	return key
}

// GetJWTExpirationHours 获取JWT过期时间
func GetJWTExpirationHours(cfg *config.Config) int {
	// 从配置中获取过期时间
	if cfg.AuthConfig != nil && cfg.AuthConfig.TokenExpiration > 0 {
		return cfg.AuthConfig.TokenExpiration
	}

	// 默认24小时
	return 24
}

// GenerateToken 生成JWT令牌
func GenerateToken(cfg *config.Config) (string, int, error) {
	// 确保JWT密钥已初始化
	if len(jwtSecret) == 0 {
		InitJWTSecret(cfg)
	}

	// 确定过期时间
	expirationHours := GetJWTExpirationHours(cfg)
	expirationTime := time.Now().Add(time.Hour * time.Duration(expirationHours))

	// 创建声明
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "page-spy",
		},
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", 0, err
	}

	return tokenString, expirationHours, nil
}

// ParseToken 解析和验证JWT令牌
func ParseToken(tokenString string) (*Claims, error) {
	// 确保JWT密钥已初始化
	if len(jwtSecret) == 0 {
		return nil, fmt.Errorf("JWT secret not initialized")
	}

	// 解析JWT令牌
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// 提取声明
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
