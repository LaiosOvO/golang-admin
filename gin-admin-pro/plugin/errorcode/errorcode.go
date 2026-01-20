package errorcode

import (
	"fmt"
	"gorm.io/gorm"
)

// ErrorCode 错误码
type ErrorCode struct {
	ID       int    `gorm:"primarykey" json:"id"`
	Type     string `gorm:"size:100;not null" json:"type"`    // 错误类型
	Code     int    `gorm:"not null;uniqueIndex" json:"code"` // 错误码
	Name     string `gorm:"size:100;not null" json:"name"`    // 错误名称
	Message  string `gorm:"size:500;not null" json:"message"` // 错误消息
	Solution string `gorm:"size:1000" json:"solution"`        // 解决方案
	Status   int    `gorm:"default:1" json:"status"`          // 状态 0-禁用 1-启用
	Remark   string `gorm:"size:500" json:"remark"`           // 备注
	CreateBy uint   `json:"createBy"`
	UpdateBy uint   `json:"updateBy"`
	CreateAt uint   `json:"createAt"`
	UpdateAt uint   `json:"updateAt"`
}

// TableName 设置表名
func (ErrorCode) TableName() string {
	return "system_error_code"
}

// Service 错误码服务
type Service struct {
	db     *gorm.DB
	config *Config
}

// NewService 创建错误码服务
func NewService(db *gorm.DB, config *Config) *Service {
	if config == nil {
		config = DefaultConfig()
	}
	return &Service{
		db:     db,
		config: config,
	}
}

// GetErrorCodes 获取错误码列表
func (s *Service) GetErrorCodes(page, pageSize int, errorCodeType, name, status string) ([]ErrorCode, int64, error) {
	var errorCodes []ErrorCode
	var total int64

	query := s.db.Model(&ErrorCode{})

	// 添加查询条件
	if errorCodeType != "" {
		query = query.Where("type LIKE ?", "%"+errorCodeType+"%")
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Order("code ASC").Find(&errorCodes).Error
	if err != nil {
		return nil, 0, err
	}

	return errorCodes, total, nil
}

// GetErrorCodeByID 根据ID获取错误码
func (s *Service) GetErrorCodeByID(id int) (*ErrorCode, error) {
	var errorCode ErrorCode
	err := s.db.First(&errorCode, id).Error
	if err != nil {
		return nil, err
	}
	return &errorCode, nil
}

// GetErrorCodeByCode 根据错误码获取错误信息
func (s *Service) GetErrorCodeByCode(code int) (*ErrorCode, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	var errorCode ErrorCode
	err := s.db.Where("code = ? AND status = 1", code).First(&errorCode).Error
	if err != nil {
		return nil, err
	}
	return &errorCode, nil
}

// CreateErrorCode 创建错误码
func (s *Service) CreateErrorCode(errorCode *ErrorCode) error {
	// 检查错误码是否已存在
	var count int64
	err := s.db.Model(&ErrorCode{}).Where("code = ?", errorCode.Code).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("错误码 %d 已存在", errorCode.Code)
	}

	return s.db.Create(errorCode).Error
}

// UpdateErrorCode 更新错误码
func (s *Service) UpdateErrorCode(errorCode *ErrorCode) error {
	// 检查错误码是否被其他记录使用
	var count int64
	err := s.db.Model(&ErrorCode{}).Where("code = ? AND id != ?", errorCode.Code, errorCode.ID).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("错误码 %d 已存在", errorCode.Code)
	}

	return s.db.Save(errorCode).Error
}

// DeleteErrorCode 删除错误码
func (s *Service) DeleteErrorCode(id int) error {
	return s.db.Delete(&ErrorCode{}, id).Error
}

// GetErrorMessage 获取错误消息
func (s *Service) GetErrorMessage(code int) string {
	if s.db == nil {
		return s.getDefaultErrorMessage(code)
	}
	errorCode, err := s.GetErrorCodeByCode(code)
	if err != nil {
		return s.getDefaultErrorMessage(code)
	}
	return errorCode.Message
}

// GetErrorMessageWithParams 获取带参数的错误消息
func (s *Service) GetErrorMessageWithParams(code int, params ...interface{}) string {
	message := s.GetErrorMessage(code)
	if len(params) > 0 {
		return fmt.Sprintf(message, params...)
	}
	return message
}

// GetErrorSolution 获取解决方案
func (s *Service) GetErrorSolution(code int) string {
	errorCode, err := s.GetErrorCodeByCode(code)
	if err != nil {
		return ""
	}
	return errorCode.Solution
}

// ValidateErrorCode 验证错误码
func (s *Service) ValidateErrorCode(code int) error {
	_, err := s.GetErrorCodeByCode(code)
	return err
}

// GetErrorCodesByType 根据类型获取错误码列表
func (s *Service) GetErrorCodesByType(errorType string) ([]ErrorCode, error) {
	var errorCodes []ErrorCode
	err := s.db.Where("type = ? AND status = 1", errorType).Order("code ASC").Find(&errorCodes).Error
	return errorCodes, err
}

// SearchErrorCodes 搜索错误码
func (s *Service) SearchErrorCodes(keyword string) ([]ErrorCode, error) {
	var errorCodes []ErrorCode
	err := s.db.Where("(name LIKE ? OR message LIKE ? OR code = ?) AND status = 1",
		"%"+keyword+"%", "%"+keyword+"%", keyword).Find(&errorCodes).Error
	return errorCodes, err
}

// ExportErrorCodes 导出错误码
func (s *Service) ExportErrorCodes(ids []int, errorType string) ([]ErrorCode, error) {
	var errorCodes []ErrorCode
	query := s.db.Model(&ErrorCode{})

	if len(ids) > 0 {
		query = query.Where("id IN ?", ids)
	}
	if errorType != "" {
		query = query.Where("type = ?", errorType)
	}

	err := query.Find(&errorCodes).Error
	return errorCodes, err
}

// ImportErrorCodes 导入错误码
func (s *Service) ImportErrorCodes(errorCodes []ErrorCode) (int, int, error) {
	var successCount, errorCount int

	for _, errorCode := range errorCodes {
		// 检查是否已存在
		var existing ErrorCode
		err := s.db.Where("code = ?", errorCode.Code).First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			// 不存在，创建新记录
			err = s.db.Create(&errorCode).Error
			if err != nil {
				errorCount++
				continue
			}
			successCount++
		} else if err == nil {
			// 已存在，更新记录
			errorCode.ID = existing.ID
			err = s.db.Save(&errorCode).Error
			if err != nil {
				errorCount++
				continue
			}
			successCount++
		} else {
			errorCount++
		}
	}

	return successCount, errorCount, nil
}

// GetErrorCodeTypes 获取所有错误类型
func (s *Service) GetErrorCodeTypes() ([]string, error) {
	var types []string
	err := s.db.Model(&ErrorCode{}).Distinct("type").Where("status = 1").Pluck("type", &types).Error
	return types, err
}

// getDefaultValue 获取默认错误消息
func (s *Service) getDefaultErrorMessage(code int) string {
	// 检查预定义错误码
	switch code {
	case CodeUnknownError:
		return "未知错误，请联系管理员"
	case CodeParamError:
		return "请求参数错误：%s"
	case CodeDataNotFound:
		return "数据不存在"
	case CodeDataExists:
		return "数据已存在"
	case CodeOperationFailed:
		return "操作失败"
	case CodePermissionDenied:
		return "权限不足"
	case CodeTokenExpired:
		return "Token已过期"
	case CodeTokenInvalid:
		return "Token无效"
	case CodeRateLimitExceeded:
		return "请求频率超限"
	case CodeUserNotFound:
		return "用户不存在"
	case CodeUserExists:
		return "用户已存在"
	case CodeUserDisabled:
		return "用户已被禁用"
	case CodePasswordError:
		return "密码错误"
	case CodeAccountLocked:
		return "账户已被锁定"
	case CodeLoginRequired:
		return "请先登录"
	case CodeRoleNotFound:
		return "角色不存在"
	case CodeRoleExists:
		return "角色已存在"
	case CodeRoleInUse:
		return "角色正在使用中"
	case CodePermissionNotFound:
		return "权限不存在"
	case CodePermissionExists:
		return "权限已存在"
	case CodeBusinessError:
		return "业务错误"
	case CodeDataInvalid:
		return "数据无效"
	case CodeConfigError:
		return "配置错误"
	case CodeServiceUnavailable:
		return "服务不可用"
	}

	// 根据错误码范围返回默认消息
	switch {
	case code >= 200 && code < 300:
		return "操作成功"
	case code >= 400 && code < 500:
		return "请求参数错误"
	case code >= 500 && code < 600:
		return "服务器内部错误"
	default:
		return fmt.Sprintf("未知错误 %d", code)
	}
}

// predefined error codes
const (
	// 通用错误码 (1000-1999)
	CodeSuccess           = 0    // 成功
	CodeUnknownError      = 1001 // 未知错误
	CodeParamError        = 1002 // 参数错误
	CodeDataNotFound      = 1003 // 数据不存在
	CodeDataExists        = 1004 // 数据已存在
	CodeOperationFailed   = 1005 // 操作失败
	CodePermissionDenied  = 1006 // 权限不足
	CodeTokenExpired      = 1007 // Token已过期
	CodeTokenInvalid      = 1008 // Token无效
	CodeRateLimitExceeded = 1009 // 请求频率超限

	// 用户相关错误码 (2000-2999)
	CodeUserNotFound  = 2001 // 用户不存在
	CodeUserExists    = 2002 // 用户已存在
	CodeUserDisabled  = 2003 // 用户已禁用
	CodePasswordError = 2004 // 密码错误
	CodeAccountLocked = 2005 // 账户已锁定
	CodeLoginRequired = 2006 // 需要登录

	// 角色相关错误码 (3000-3999)
	CodeRoleNotFound = 3001 // 角色不存在
	CodeRoleExists   = 3002 // 角色已存在
	CodeRoleInUse    = 3003 // 角色正在使用

	// 权限相关错误码 (4000-4999)
	CodePermissionNotFound = 4001 // 权限不存在
	CodePermissionExists   = 4002 // 权限已存在

	// 业务相关错误码 (5000-9999)
	CodeBusinessError      = 5001 // 业务逻辑错误
	CodeDataInvalid        = 5002 // 数据无效
	CodeConfigError        = 5003 // 配置错误
	CodeServiceUnavailable = 5004 // 服务不可用
)

// GetPredefinedErrorCodes 获取预定义错误码列表
func GetPredefinedErrorCodes() []ErrorCode {
	return []ErrorCode{
		// 通用错误码
		{Type: "common", Code: CodeSuccess, Name: "操作成功", Message: "操作成功", Status: 1},
		{Type: "common", Code: CodeUnknownError, Name: "未知错误", Message: "未知错误，请联系管理员", Status: 1},
		{Type: "common", Code: CodeParamError, Name: "参数错误", Message: "请求参数错误：%s", Status: 1},
		{Type: "common", Code: CodeDataNotFound, Name: "数据不存在", Message: "数据不存在：%s", Status: 1},
		{Type: "common", Code: CodeDataExists, Name: "数据已存在", Message: "数据已存在：%s", Status: 1},
		{Type: "common", Code: CodeOperationFailed, Name: "操作失败", Message: "操作失败：%s", Status: 1},
		{Type: "common", Code: CodePermissionDenied, Name: "权限不足", Message: "权限不足，无法执行此操作", Status: 1},
		{Type: "common", Code: CodeTokenExpired, Name: "Token已过期", Message: "登录已过期，请重新登录", Status: 1},
		{Type: "common", Code: CodeTokenInvalid, Name: "Token无效", Message: "登录信息无效，请重新登录", Status: 1},
		{Type: "common", Code: CodeRateLimitExceeded, Name: "请求频率超限", Message: "请求过于频繁，请稍后再试", Status: 1},

		// 用户相关错误码
		{Type: "user", Code: CodeUserNotFound, Name: "用户不存在", Message: "用户不存在", Status: 1},
		{Type: "user", Code: CodeUserExists, Name: "用户已存在", Message: "用户已存在：%s", Status: 1},
		{Type: "user", Code: CodeUserDisabled, Name: "用户已禁用", Message: "用户已被禁用", Status: 1},
		{Type: "user", Code: CodePasswordError, Name: "密码错误", Message: "密码错误，请重新输入", Status: 1},
		{Type: "user", Code: CodeAccountLocked, Name: "账户已锁定", Message: "账户已锁定，请联系管理员", Status: 1},
		{Type: "user", Code: CodeLoginRequired, Name: "需要登录", Message: "请先登录后再进行操作", Status: 1},

		// 角色相关错误码
		{Type: "role", Code: CodeRoleNotFound, Name: "角色不存在", Message: "角色不存在", Status: 1},
		{Type: "role", Code: CodeRoleExists, Name: "角色已存在", Message: "角色已存在：%s", Status: 1},
		{Type: "role", Code: CodeRoleInUse, Name: "角色正在使用", Message: "角色正在使用中，无法删除", Status: 1},

		// 权限相关错误码
		{Type: "permission", Code: CodePermissionNotFound, Name: "权限不存在", Message: "权限不存在", Status: 1},
		{Type: "permission", Code: CodePermissionExists, Name: "权限已存在", Message: "权限已存在：%s", Status: 1},

		// 业务相关错误码
		{Type: "business", Code: CodeBusinessError, Name: "业务逻辑错误", Message: "业务逻辑错误：%s", Status: 1},
		{Type: "business", Code: CodeDataInvalid, Name: "数据无效", Message: "数据无效：%s", Status: 1},
		{Type: "business", Code: CodeConfigError, Name: "配置错误", Message: "配置错误：%s", Status: 1},
		{Type: "business", Code: CodeServiceUnavailable, Name: "服务不可用", Message: "服务暂时不可用，请稍后再试", Status: 1},
	}
}

// GetErrorTypeName 获取错误类型名称
func GetErrorTypeName(errorType string) string {
	typeNames := map[string]string{
		"common":     "通用错误",
		"user":       "用户错误",
		"role":       "角色错误",
		"permission": "权限错误",
		"business":   "业务错误",
		"system":     "系统错误",
	}

	if name, exists := typeNames[errorType]; exists {
		return name
	}
	return errorType
}

// GetErrorTypeOptions 获取错误类型选项
func GetErrorTypeOptions() []map[string]interface{} {
	types := []string{"common", "user", "role", "permission", "business", "system"}
	var options []map[string]interface{}

	for _, t := range types {
		options = append(options, map[string]interface{}{
			"value": t,
			"label": GetErrorTypeName(t),
		})
	}

	return options
}
