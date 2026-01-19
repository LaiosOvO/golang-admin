package oss

import (
	"fmt"
	"gin-admin-pro/internal/pkg/config"
	"io"
)

// MinIOStorage MinIO对象存储（暂未实现）
type MinIOStorage struct {
	config FileConfig
}

// NewMinIOStorage 创建MinIO存储实例
func NewMinIOStorage() (*MinIOStorage, error) {
	cfg := config.GetConfig()
	fileConfig := DefaultFileConfig()
	fileConfig.Provider = "minio"

	// 从配置文件读取文件配置
	if cfg.Upload.MaxSize > 0 {
		fileConfig.MaxSize = int64(cfg.Upload.MaxSize)
	}
	if len(cfg.Upload.AllowedTypes) > 0 {
		fileConfig.AllowExts = cfg.Upload.AllowedTypes
	}

	return &MinIOStorage{
		config: fileConfig,
	}, nil
}

// UploadFile 上传文件
func (ms *MinIOStorage) UploadFile(key string, reader io.Reader, size int64, contentType string) (string, error) {
	return "", fmt.Errorf("MinIO存储暂未实现，请使用本地存储")
}

// DeleteFile 删除文件
func (ms *MinIOStorage) DeleteFile(key string) error {
	return fmt.Errorf("MinIO存储暂未实现，请使用本地存储")
}

// GetFileURL 获取文件访问URL
func (ms *MinIOStorage) GetFileURL(key string) string {
	return "/minio-not-implemented"
}

// IsExists 检查文件是否存在
func (ms *MinIOStorage) IsExists(key string) (bool, error) {
	return false, fmt.Errorf("MinIO存储暂未实现，请使用本地存储")
}

// GetConfig 获取配置
func (ms *MinIOStorage) GetConfig() FileConfig {
	return ms.config
}
