package oss

import "io"

// OSSInterface 对象存储接口
type OSSInterface interface {
	// UploadFile 上传文件
	UploadFile(key string, reader io.Reader, size int64, contentType string) (string, error)

	// DeleteFile 删除文件
	DeleteFile(key string) error

	// GetFileURL 获取文件访问URL
	GetFileURL(key string) string

	// IsExists 检查文件是否存在
	IsExists(key string) (bool, error)
}

// UploadResult 上传结果
type UploadResult struct {
	Key      string `json:"key"`      // 文件key
	URL      string `json:"url"`      // 访问URL
	FileName string `json:"fileName"` // 文件名
	Size     int64  `json:"size"`     // 文件大小
}

// FileConfig 文件配置
type FileConfig struct {
	MaxSize   int64    `yaml:"maxSize"`   // 最大文件大小（字节）
	AllowExts []string `yaml:"allowExts"` // 允许的文件扩展名
	Provider  string   `yaml:"provider"`  // 存储提供商：local/minio/aliyun
	Path      string   `yaml:"path"`      // 存储路径
}

// DefaultFileConfig 默认文件配置
func DefaultFileConfig() FileConfig {
	return FileConfig{
		MaxSize:   10 * 1024 * 1024, // 10MB
		AllowExts: []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt", ".zip", ".rar"},
		Provider:  "local",
		Path:      "./uploads",
	}
}
