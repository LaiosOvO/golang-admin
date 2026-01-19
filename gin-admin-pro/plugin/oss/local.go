package oss

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gin-admin-pro/internal/pkg/config"
)

// LocalStorage 本地文件存储
type LocalStorage struct {
	config FileConfig
}

// NewLocalStorage 创建本地存储实例
func NewLocalStorage() (*LocalStorage, error) {
	cfg := config.GetConfig()
	fileConfig := DefaultFileConfig()

	// 从配置文件读取文件配置
	if cfg.Upload.MaxSize > 0 {
		fileConfig.MaxSize = int64(cfg.Upload.MaxSize)
	}
	if len(cfg.Upload.AllowedTypes) > 0 {
		fileConfig.AllowExts = cfg.Upload.AllowedTypes
	}
	if cfg.Upload.Path != "" {
		fileConfig.Path = cfg.Upload.Path
	}

	// 创建存储目录
	if err := os.MkdirAll(fileConfig.Path, 0755); err != nil {
		return nil, fmt.Errorf("创建存储目录失败: %v", err)
	}

	return &LocalStorage{
		config: fileConfig,
	}, nil
}

// UploadFile 上传文件
func (ls *LocalStorage) UploadFile(key string, reader io.Reader, size int64, contentType string) (string, error) {
	// 检查文件大小
	if size > ls.config.MaxSize {
		return "", fmt.Errorf("文件大小超过限制，最大允许 %d MB", ls.config.MaxSize/(1024*1024))
	}

	// 生成文件路径
	datePath := time.Now().Format("2006/01/02")
	fullDir := filepath.Join(ls.config.Path, datePath)

	// 创建目录
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return "", fmt.Errorf("创建文件目录失败: %v", err)
	}

	// 完整文件路径
	fullPath := filepath.Join(fullDir, key)

	// 创建文件
	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	// 复制文件内容
	written, err := io.Copy(file, reader)
	if err != nil {
		os.Remove(fullPath) // 删除已创建的文件
		return "", fmt.Errorf("保存文件失败: %v", err)
	}

	if written != size {
		os.Remove(fullPath) // 删除已创建的文件
		return "", fmt.Errorf("文件保存不完整")
	}

	return ls.GetFileURL(key), nil
}

// DeleteFile 删除文件
func (ls *LocalStorage) DeleteFile(key string) error {
	fullPath := ls.getFullPath(key)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil // 文件不存在，视为删除成功
	}

	return os.Remove(fullPath)
}

// GetFileURL 获取文件访问URL
func (ls *LocalStorage) GetFileURL(key string) string {
	datePath := time.Now().Format("2006/01/02")
	return fmt.Sprintf("/uploads/%s/%s", datePath, key)
}

// IsExists 检查文件是否存在
func (ls *LocalStorage) IsExists(key string) (bool, error) {
	fullPath := ls.getFullPath(key)
	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// getFullPath 获取文件的完整路径
func (ls *LocalStorage) getFullPath(key string) string {
	// 从key中提取日期路径和文件名
	parts := strings.Split(key, "/")
	if len(parts) >= 4 {
		// 格式: uploads/2025/01/20/filename.ext
		datePath := strings.Join(parts[1:4], "/")
		filename := strings.Join(parts[4:], "/")
		return filepath.Join(ls.config.Path, datePath, filename)
	}

	// 简单情况，直接在根目录
	return filepath.Join(ls.config.Path, key)
}

// ValidateFile 验证文件
func (ls *LocalStorage) ValidateFile(file *multipart.FileHeader) error {
	// 检查文件大小
	if file.Size > ls.config.MaxSize {
		return fmt.Errorf("文件大小超过限制，最大允许 %d MB", ls.config.MaxSize/(1024*1024))
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := false
	for _, allowedExt := range ls.config.AllowExts {
		if ext == allowedExt {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("不支持的文件类型: %s", ext)
	}

	return nil
}

// GenerateFileName 生成文件名
func (ls *LocalStorage) GenerateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	nameWithoutExt := strings.TrimSuffix(originalName, ext)

	// 清理文件名中的特殊字符
	nameWithoutExt = strings.ReplaceAll(nameWithoutExt, " ", "_")
	nameWithoutExt = strings.ReplaceAll(nameWithoutExt, "/", "_")
	nameWithoutExt = strings.ReplaceAll(nameWithoutExt, "\\", "_")

	// 生成时间戳
	timestamp := time.Now().Format("20060102150405")

	return fmt.Sprintf("%s_%s%s", nameWithoutExt, timestamp, ext)
}

// GetConfig 获取配置
func (ls *LocalStorage) GetConfig() FileConfig {
	return ls.config
}
