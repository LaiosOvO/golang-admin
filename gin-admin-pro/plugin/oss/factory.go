package oss

import (
	"gin-admin-pro/internal/pkg/config"
)

// StorageFactory 存储工厂
type StorageFactory struct{}

// NewStorageFactory 创建存储工厂
func NewStorageFactory() *StorageFactory {
	return &StorageFactory{}
}

// CreateStorage 创建存储实例
func (f *StorageFactory) CreateStorage(provider string) (OSSInterface, error) {
	cfg := config.GetConfig()

	if provider == "" {
		provider = cfg.Upload.Path
		if provider == "" {
			provider = "local"
		}
	}

	switch provider {
	case "local":
		return NewLocalStorage()
	case "minio":
		return NewMinIOStorage()
	default:
		// 暂时默认使用本地存储
		return NewLocalStorage()
	}
}

// GetDefaultStorage 获取默认存储实例
func GetDefaultStorage() (OSSInterface, error) {
	factory := NewStorageFactory()
	return factory.CreateStorage("local")
}
