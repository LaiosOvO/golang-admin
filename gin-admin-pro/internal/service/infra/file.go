package infra

import (
	"errors"
	"fmt"
	"gin-admin-pro/plugin/oss"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

// FileService 文件服务
type FileService struct {
	storage oss.OSSInterface
}

// NewFileService 创建文件服务实例
func NewFileService(storage oss.OSSInterface) *FileService {
	return &FileService{
		storage: storage,
	}
}

// UploadResult 上传结果
type UploadResult struct {
	URL      string `json:"url"`      // 文件访问URL
	FileName string `json:"fileName"` // 文件名
	Size     int64  `json:"size"`     // 文件大小
	Type     string `json:"type"`     // 文件类型
}

// UploadFile 上传单个文件
func (s *FileService) UploadFile(file *multipart.FileHeader) (*UploadResult, error) {
	if file == nil {
		return nil, errors.New("请选择要上传的文件")
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		return nil, errors.New("打开文件失败")
	}
	defer src.Close()

	// 生成文件名
	fileName := s.generateFileName(file.Filename)

	// 上传文件
	fileURL, err := s.storage.UploadFile(fileName, src, file.Size, file.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	return &UploadResult{
		URL:      fileURL,
		FileName: file.Filename,
		Size:     file.Size,
		Type:     filepath.Ext(file.Filename),
	}, nil
}

// UploadMultipleFiles 上传多个文件
func (s *FileService) UploadMultipleFiles(files []*multipart.FileHeader) ([]*UploadResult, error) {
	if len(files) == 0 {
		return nil, errors.New("请选择要上传的文件")
	}

	var results []*UploadResult
	var errors []string

	for i, file := range files {
		if file == nil {
			continue
		}

		result, err := s.UploadFile(file)
		if err != nil {
			errors = append(errors, fmt.Sprintf("文件%d上传失败: %s", i+1, err.Error()))
			continue
		}

		results = append(results, result)
	}

	if len(errors) > 0 {
		return results, fmt.Errorf(strings.Join(errors, "; "))
	}

	return results, nil
}

// DeleteFile 删除文件
func (s *FileService) DeleteFile(fileURL string) error {
	// 从URL中提取文件key
	key := s.extractKeyFromURL(fileURL)
	if key == "" {
		return errors.New("无效的文件URL")
	}

	return s.storage.DeleteFile(key)
}

// generateFileName 生成文件名
func (s *FileService) generateFileName(originalName string) string {
	ext := strings.ToLower(filepath.Ext(originalName))
	nameWithoutExt := strings.TrimSuffix(originalName, ext)

	// 清理文件名中的特殊字符
	nameWithoutExt = strings.ReplaceAll(nameWithoutExt, " ", "_")
	nameWithoutExt = strings.ReplaceAll(nameWithoutExt, "/", "_")
	nameWithoutExt = strings.ReplaceAll(nameWithoutExt, "\\", "_")

	// 生成时间戳
	timestamp := time.Now().Format("20060102150405")

	return fmt.Sprintf("%s_%s%s", nameWithoutExt, timestamp, ext)
}

// extractKeyFromURL 从URL中提取文件key
func (s *FileService) extractKeyFromURL(fileURL string) string {
	// 移除URL前缀
	if strings.HasPrefix(fileURL, "/uploads/") {
		return strings.TrimPrefix(fileURL, "/uploads/")
	}

	// 如果是完整URL，提取路径部分
	if strings.Contains(fileURL, "://") {
		parts := strings.Split(fileURL, "/")
		if len(parts) > 3 {
			return strings.Join(parts[3:], "/")
		}
	}

	return fileURL
}

// ValidateFile 验证文件
func (s *FileService) ValidateFile(file *multipart.FileHeader) error {
	if file == nil {
		return errors.New("请选择要上传的文件")
	}

	// 检查文件大小
	if file.Size <= 0 {
		return errors.New("文件大小不能为0")
	}

	return nil
}
