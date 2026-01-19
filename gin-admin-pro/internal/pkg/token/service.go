package token

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"gin-admin-pro/internal/pkg/config"
	"gin-admin-pro/internal/pkg/jwt"
	"gin-admin-pro/plugin/redis"
)

// TokenService Token服务
type TokenService struct {
	redisClient *redis.Client
	config      *config.Config
}

// NewTokenService 创建Token服务
func NewTokenService(redisClient *redis.Client) *TokenService {
	return &TokenService{
		redisClient: redisClient,
		config:      config.GetConfig(),
	}
}

// GenerateTokens 生成访问Token和刷新Token
func (s *TokenService) GenerateTokens(userID uint, username string) (*TokenPair, error) {
	// 生成访问Token
	accessToken, err := jwt.GenerateToken(userID, username)
	if err != nil {
		return nil, fmt.Errorf("生成访问Token失败: %w", err)
	}

	// 生成刷新Token
	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("生成刷新Token失败: %w", err)
	}

	// 计算过期时间
	accessExpire := time.Duration(s.config.JWT.AccessTokenExpire) * 24 * time.Hour
	refreshExpire := time.Duration(s.config.JWT.RefreshTokenExpire) * 24 * time.Hour

	ctx := context.Background()

	// 存储访问Token到Redis
	accessKey := s.getAccessTokenKey(accessToken)
	if err := s.redisClient.Set(ctx, accessKey, userID, accessExpire); err != nil {
		return nil, fmt.Errorf("存储访问Token失败: %w", err)
	}

	// 存储刷新Token到Redis
	refreshKey := s.getRefreshTokenKey(refreshToken)
	if err := s.redisClient.Set(ctx, refreshKey, userID, refreshExpire); err != nil {
		return nil, fmt.Errorf("存储刷新Token失败: %w", err)
	}

	// 存储用户的Token映射（用于单点登录）
	userTokensKey := s.getUserTokensKey(fmt.Sprintf("%d", userID))
	s.redisClient.SAdd(ctx, userTokensKey, accessToken, refreshToken)
	s.redisClient.Expire(ctx, userTokensKey, refreshExpire)

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(accessExpire.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// ValidateToken 验证Token
func (s *TokenService) ValidateToken(tokenString string) (*TokenInfo, error) {
	ctx := context.Background()

	// 检查Token是否在Redis中存在
	accessKey := s.getAccessTokenKey(tokenString)
	exists, err := s.redisClient.Exists(ctx, accessKey)
	if err != nil {
		return nil, fmt.Errorf("检查Token失败: %w", err)
	}

	if exists == 0 {
		return nil, fmt.Errorf("Token不存在或已过期")
	}

	// 解析JWT Token
	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("解析Token失败: %w", err)
	}

	return &TokenInfo{
		UserID:    claims.UserID,
		Username:  claims.Username,
		Subject:   claims.Subject,
		ExpiresAt: claims.ExpiresAt.Time,
		IssuedAt:  claims.IssuedAt.Time,
	}, nil
}

// RefreshToken 刷新Token
func (s *TokenService) RefreshToken(refreshToken string) (*TokenPair, error) {
	ctx := context.Background()

	// 检查刷新Token是否存在
	refreshKey := s.getRefreshTokenKey(refreshToken)
	exists, err := s.redisClient.Exists(ctx, refreshKey)
	if err != nil {
		return nil, fmt.Errorf("检查刷新Token失败: %w", err)
	}

	if exists == 0 {
		return nil, fmt.Errorf("刷新Token不存在或已过期")
	}

	// 获取用户ID（暂时不使用，后续优化）
	_, err = s.redisClient.Get(ctx, refreshKey)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// TODO: 从数据库获取用户名，这里使用临时方式
	username := "user"

	// 撤销旧的Token
	s.RevokeToken(refreshToken)

	// 生成新的Token对（TODO: 修复用户ID转换）
	return s.GenerateTokens(0, username) // TODO: 修复用户ID转换
}

// RevokeToken 撤销Token
func (s *TokenService) RevokeToken(tokenString string) error {
	ctx := context.Background()

	// 检查是否是访问Token
	accessKey := s.getAccessTokenKey(tokenString)
	if exists, _ := s.redisClient.Exists(ctx, accessKey); exists > 0 {
		// 获取用户ID
		userIDStr, err := s.redisClient.Get(ctx, accessKey)
		if err == nil {
			s.removeUserToken(ctx, userIDStr, tokenString)
		}
		return s.redisClient.Del(ctx, accessKey)
	}

	// 检查是否是刷新Token
	refreshKey := s.getRefreshTokenKey(tokenString)
	if exists, _ := s.redisClient.Exists(ctx, refreshKey); exists > 0 {
		// 获取用户ID
		userIDStr, err := s.redisClient.Get(ctx, refreshKey)
		if err == nil {
			s.removeUserToken(ctx, userIDStr, tokenString)
		}
		return s.redisClient.Del(ctx, refreshKey)
	}

	return fmt.Errorf("Token不存在")
}

// RevokeAllUserTokens 撤销用户所有Token（用于单点登录）
func (s *TokenService) RevokeAllUserTokens(userID uint) error {
	ctx := context.Background()

	userTokensKey := s.getUserTokensKey(fmt.Sprintf("%d", userID))
	tokens, err := s.redisClient.SMembers(ctx, userTokensKey)
	if err != nil {
		return fmt.Errorf("获取用户Token失败: %w", err)
	}

	// 删除所有Token
	for _, token := range tokens {
		accessKey := s.getAccessTokenKey(token)
		refreshKey := s.getRefreshTokenKey(token)

		s.redisClient.Del(ctx, accessKey)
		s.redisClient.Del(ctx, refreshKey)
	}

	// 删除用户Token集合
	return s.redisClient.Del(ctx, userTokensKey)
}

// generateRefreshToken 生成随机刷新Token
func (s *TokenService) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// getAccessTokenKey 获取访问Token的Redis键
func (s *TokenService) getAccessTokenKey(token string) string {
	return fmt.Sprintf("jwt:access:%s", token)
}

// getRefreshTokenKey 获取刷新Token的Redis键
func (s *TokenService) getRefreshTokenKey(token string) string {
	return fmt.Sprintf("jwt:refresh:%s", token)
}

// getUserTokensKey 获取用户Token集合的Redis键
func (s *TokenService) getUserTokensKey(userID string) string {
	return fmt.Sprintf("jwt:user_tokens:%s", userID)
}

// removeUserToken 从用户Token集合中移除Token
func (s *TokenService) removeUserToken(ctx context.Context, userIDStr, token string) {
	userTokensKey := s.getUserTokensKey(userIDStr)
	s.redisClient.SRem(ctx, userTokensKey, token)
}

// TokenPair Token对
type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
	TokenType    string `json:"tokenType"`
}

// TokenInfo Token信息
type TokenInfo struct {
	UserID    uint      `json:"userId"`
	Username  string    `json:"username"`
	Subject   string    `json:"subject"`
	ExpiresAt time.Time `json:"expiresAt"`
	IssuedAt  time.Time `json:"issuedAt"`
}
