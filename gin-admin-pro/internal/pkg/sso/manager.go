package sso

import (
	"context"
	"fmt"

	"gin-admin-pro/internal/service"
)

// SSOManager 单点登录管理器
type SSOManager struct {
	tokenService interface {
		RevokeAllUserTokens(userID uint) error
		RevokeToken(tokenString string) error
	}
}

// NewSSOManager 创建SSO管理器
func NewSSOManager() *SSOManager {
	return &SSOManager{
		tokenService: service.Services.TokenService,
	}
}

// EnableSingleSignOn 启用单点登录
// 当用户在新设备登录时，撤销之前所有的Token
func (s *SSOManager) EnableSingleSignOn(userID uint) error {
	// 撤销用户所有现有的Token
	return s.tokenService.RevokeAllUserTokens(userID)
}

// EnableMultiDeviceLogin 允许多设备登录
// 用户可以在多个设备上同时登录
func (s *SSOManager) EnableMultiDeviceLogin(userID uint) error {
	// 不撤销现有Token，允许新的Token生成
	return nil
}

// IsUserOnline 检查用户是否在线
func (s *SSOManager) IsUserOnline(userID uint) (bool, error) {
	ctx := context.Background()
	userTokensKey := fmt.Sprintf("jwt:user_tokens:%d", userID)

	exists, err := service.Services.RedisClient.Exists(ctx, userTokensKey)
	if err != nil {
		return false, fmt.Errorf("检查用户在线状态失败: %w", err)
	}

	return exists > 0, nil
}

// GetActiveDevices 获取用户活跃设备数量
func (s *SSOManager) GetActiveDevices(userID uint) (int, error) {
	ctx := context.Background()
	userTokensKey := fmt.Sprintf("jwt:user_tokens:%d", userID)

	tokens, err := service.Services.RedisClient.SMembers(ctx, userTokensKey)
	if err != nil {
		return 0, fmt.Errorf("获取用户设备信息失败: %w", err)
	}

	// 每个访问Token代表一个设备
	deviceCount := 0
	for _, token := range tokens {
		// 检查是否是访问Token（访问Token存储在 jwt:access: 前缀下）
		accessKey := fmt.Sprintf("jwt:access:%s", token)
		exists, _ := service.Services.RedisClient.Exists(ctx, accessKey)
		if exists > 0 {
			deviceCount++
		}
	}

	return deviceCount, nil
}

// RevokeDevice 撤销指定设备的Token
func (s *SSOManager) RevokeDevice(userID uint, deviceToken string) error {
	return s.tokenService.RevokeToken(deviceToken)
}

// RevokeAllDevices 撤销用户所有设备的Token
func (s *SSOManager) RevokeAllDevices(userID uint) error {
	return s.tokenService.RevokeAllUserTokens(userID)
}

// LoginOptions 登录选项
type LoginOptions struct {
	EnableSSO  bool   `json:"enableSSO"`  // 是否启用单点登录
	MaxDevices int    `json:"maxDevices"` // 最大设备数量（0表示无限制）
	DeviceName string `json:"deviceName"` // 设备名称
	DeviceInfo string `json:"deviceInfo"` // 设备信息
}

// LoginWithOptions 使用选项登录
func (s *SSOManager) LoginWithOptions(userID uint, username string, opts LoginOptions) error {
	// 如果启用单点登录，先撤销现有Token
	if opts.EnableSSO {
		if err := s.EnableSingleSignOn(userID); err != nil {
			return fmt.Errorf("启用单点登录失败: %w", err)
		}
	}

	// 如果设置了最大设备数限制
	if opts.MaxDevices > 0 {
		activeDevices, err := s.GetActiveDevices(userID)
		if err != nil {
			return fmt.Errorf("获取活跃设备数失败: %w", err)
		}

		// 如果达到最大设备数，撤销最早的Token
		if activeDevices >= opts.MaxDevices {
			ctx := context.Background()
			userTokensKey := fmt.Sprintf("jwt:user_tokens:%d", userID)
			tokens, err := service.Services.RedisClient.SMembers(ctx, userTokensKey)
			if err != nil {
				return fmt.Errorf("获取用户Token失败: %w", err)
			}

			// 撤销第一个访问Token
			for _, token := range tokens {
				accessKey := fmt.Sprintf("jwt:access:%s", token)
				exists, _ := service.Services.RedisClient.Exists(ctx, accessKey)
				if exists > 0 {
					s.tokenService.RevokeToken(token)
					break // 只撤销一个
				}
			}
		}
	}

	return nil
}
