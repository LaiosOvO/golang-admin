package operlog

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
)

// OperLog 操作日志
type OperLog struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	Title         string    `gorm:"size:50" json:"title"`                      // 操作模块
	BusinessType  int       `gorm:"default:0" json:"businessType"`             // 业务类型（0其它 1新增 2修改 3删除）
	Method        string    `gorm:"size:100" json:"method"`                    // 请求方法
	RequestMethod string    `gorm:"size:10" json:"requestMethod"`              // 请求方式
	OperatorType  int       `gorm:"default:0" json:"operatorType"`             // 操作类别（0其它 1后台用户 2手机端用户）
	OperName      string    `gorm:"size:50" json:"operName"`                   // 操作人员
	DeptName      string    `gorm:"size:50" json:"deptName"`                   // 部门名称
	OperUrl       string    `gorm:"size:255" json:"operUrl"`                   // 请求URL
	OperIp        string    `gorm:"size:128" json:"operIp"`                    // 操作地址
	OperLocation  string    `gorm:"size:255" json:"operLocation"`              // 操作地点
	OperParam     string    `gorm:"size:2000" json:"operParam"`                // 请求参数
	JsonResult    string    `gorm:"size:2000" json:"jsonResult"`               // 返回参数
	Status        int       `gorm:"default:0" json:"status"`                   // 操作状态（0正常 1异常）
	ErrorMsg      string    `gorm:"size:2000" json:"errorMsg"`                 // 错误消息
	OperTime      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"operTime"` // 操作时间
	CostTime      int64     `json:"costTime"`                                  // 消耗时间（毫秒）
}

// TableName 设置表名
func (OperLog) TableName() string {
	return "system_oper_log"
}

// LogContext 日志上下文
type LogContext struct {
	RequestID   string                 `json:"requestId"`
	UserID      uint                   `json:"userId"`
	Username    string                 `json:"username"`
	DeptID      uint                   `json:"deptId"`
	DeptName    string                 `json:"deptName"`
	Method      string                 `json:"method"`
	Path        string                 `json:"path"`
	IP          string                 `json:"ip"`
	UserAgent   string                 `json:"userAgent"`
	RequestBody map[string]interface{} `json:"requestBody"`
	StartTime   time.Time              `json:"startTime"`
	EndTime     time.Time              `json:"endTime"`
	Status      int                    `json:"status"`
	Error       string                 `json:"error"`
	Response    interface{}            `json:"response"`
}

// Service 操作日志服务
type Service struct {
	db     *gorm.DB
	config *Config
}

// NewService 创建操作日志服务
func NewService(db *gorm.DB, config *Config) *Service {
	if config == nil {
		config = DefaultConfig()
	}
	return &Service{
		db:     db,
		config: config,
	}
}

// CreateOperLog 创建操作日志
func (s *Service) CreateOperLog(log *OperLog) error {
	if !s.config.Enabled {
		return nil
	}

	// 数据脱敏
	if s.config.EnableMask {
		log.OperParam = s.maskSensitiveData(log.OperParam)
		log.JsonResult = s.maskSensitiveData(log.JsonResult)
	}

	// 检查日志长度限制
	if len(log.OperParam) > s.config.MaxParamLength {
		log.OperParam = log.OperParam[:s.config.MaxParamLength] + "..."
	}
	if len(log.JsonResult) > s.config.MaxResultLength {
		log.JsonResult = log.JsonResult[:s.config.MaxResultLength] + "..."
	}

	return s.db.Create(log).Error
}

// CreateLogFromContext 从上下文创建日志
func (s *Service) CreateLogFromContext(ctx *LogContext) error {
	if !s.config.Enabled || s.db == nil {
		return nil
	}

	// 获取IP地址位置
	location := s.getLocationByIP(ctx.IP)

	// 计算耗时
	costTime := ctx.EndTime.Sub(ctx.StartTime).Milliseconds()

	// 获取操作类型
	businessType := s.getBusinessType(ctx.Path, ctx.Method)

	// 构建操作日志
	log := &OperLog{
		Title:         s.getModuleTitle(ctx.Path),
		BusinessType:  businessType,
		Method:        ctx.Method,
		RequestMethod: ctx.Method,
		OperatorType:  1, // 后台用户
		OperName:      ctx.Username,
		DeptName:      ctx.DeptName,
		OperUrl:       ctx.Path,
		OperIp:        ctx.IP,
		OperLocation:  location,
		Status:        ctx.Status,
		ErrorMsg:      ctx.Error,
		OperTime:      ctx.EndTime,
		CostTime:      costTime,
	}

	// 序列化请求参数
	if ctx.RequestBody != nil {
		paramBytes, _ := json.Marshal(ctx.RequestBody)
		log.OperParam = string(paramBytes)
	}

	// 序列化响应结果
	if ctx.Response != nil {
		resultBytes, _ := json.Marshal(ctx.Response)
		log.JsonResult = string(resultBytes)
	}

	return s.CreateOperLog(log)
}

// GetOperLogs 获取操作日志列表
func (s *Service) GetOperLogs(page, pageSize int, title, operName, businessType, status, operTime string) ([]OperLog, int64, error) {
	var logs []OperLog
	var total int64

	query := s.db.Model(&OperLog{})

	// 添加查询条件
	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if operName != "" {
		query = query.Where("oper_name LIKE ?", "%"+operName+"%")
	}
	if businessType != "" {
		query = query.Where("business_type = ?", businessType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if operTime != "" {
		// operTime格式：2023-01-01,2023-01-31
		times := strings.Split(operTime, ",")
		if len(times) == 2 {
			query = query.Where("oper_time BETWEEN ? AND ?", times[0]+" 00:00:00", times[1]+" 23:59:59")
		}
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Order("oper_time DESC").Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetOperLogByID 根据ID获取操作日志
func (s *Service) GetOperLogByID(id uint) (*OperLog, error) {
	var log OperLog
	err := s.db.First(&log, id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// DeleteOperLog 删除操作日志
func (s *Service) DeleteOperLog(id uint) error {
	return s.db.Delete(&OperLog{}, id).Error
}

// DeleteOperLogs 批量删除操作日志
func (s *Service) DeleteOperLogs(ids []uint) error {
	return s.db.Delete(&OperLog{}, ids).Error
}

// CleanOperLog 清空操作日志
func (s *Service) CleanOperLog() error {
	return s.db.Exec("DELETE FROM system_oper_log").Error
}

// ExportOperLog 导出操作日志
func (s *Service) ExportOperLog(ids []uint) ([]OperLog, error) {
	var logs []OperLog
	query := s.db.Model(&OperLog{})
	if len(ids) > 0 {
		query = query.Where("id IN ?", ids)
	}
	err := query.Find(&logs).Error
	return logs, err
}

// maskSensitiveData 脱敏敏感数据
func (s *Service) maskSensitiveData(data string) string {
	if data == "" {
		return data
	}

	// 定义需要脱敏的字段
	sensitiveFields := s.config.SensitiveFields
	for _, field := range sensitiveFields {
		// 使用正则表达式进行更精确的脱敏处理
		pattern := fmt.Sprintf(`"%s":\s*"[^"]*"`, field)
		replacement := fmt.Sprintf(`"%s":"***"`, field)
		data = regexp.MustCompile(pattern).ReplaceAllString(data, replacement)
	}

	return data
}

// getLocationByIP 根据IP获取地址位置
func (s *Service) getLocationByIP(ip string) string {
	// 这里应该调用IP地址库获取地理位置
	// 为了简化，这里返回默认值
	if ip == "127.0.0.1" || ip == "::1" {
		return "本地访问"
	}
	return "未知位置"
}

// getBusinessType 根据路径和方法获取业务类型
func (s *Service) getBusinessType(path, method string) int {
	path = strings.ToLower(path)

	// 特殊路径处理
	if strings.Contains(path, "/auth/login") || strings.Contains(path, "/auth/logout") {
		return 0 // 其它
	}

	// 根据路径前缀和HTTP方法判断业务类型
	if strings.Contains(path, "/create") || (method == "POST" && !strings.Contains(path, "/auth")) {
		return 1 // 新增
	}
	if strings.Contains(path, "/update") || method == "PUT" {
		return 2 // 修改
	}
	if strings.Contains(path, "/delete") || method == "DELETE" {
		return 3 // 删除
	}

	return 0 // 其它
}

// getModuleTitle 根据路径获取模块标题
func (s *Service) getModuleTitle(path string) string {
	path = strings.ToLower(path)

	// 根据路径映射模块标题
	if strings.Contains(path, "/user") {
		return "用户管理"
	}
	if strings.Contains(path, "/role") {
		return "角色管理"
	}
	if strings.Contains(path, "/menu") {
		return "菜单管理"
	}
	if strings.Contains(path, "/dept") {
		return "部门管理"
	}
	if strings.Contains(path, "/dict") {
		return "字典管理"
	}
	if strings.Contains(path, "/file") {
		return "文件管理"
	}

	return "系统操作"
}

// GetBusinessTypeName 获取业务类型名称
func GetBusinessTypeName(businessType int) string {
	switch businessType {
	case 0:
		return "其它"
	case 1:
		return "新增"
	case 2:
		return "修改"
	case 3:
		return "删除"
	case 4:
		return "授权"
	case 5:
		return "导出"
	case 6:
		return "导入"
	default:
		return "未知"
	}
}

// GetBusinessTypeOptions 获取业务类型选项
func GetBusinessTypeOptions() []map[string]interface{} {
	return []map[string]interface{}{
		{"value": 0, "label": GetBusinessTypeName(0)},
		{"value": 1, "label": GetBusinessTypeName(1)},
		{"value": 2, "label": GetBusinessTypeName(2)},
		{"value": 3, "label": GetBusinessTypeName(3)},
		{"value": 4, "label": GetBusinessTypeName(4)},
		{"value": 5, "label": GetBusinessTypeName(5)},
		{"value": 6, "label": GetBusinessTypeName(6)},
	}
}

// GetOperatorTypeName 获取操作类别名称
func GetOperatorTypeName(operatorType int) string {
	switch operatorType {
	case 0:
		return "其它"
	case 1:
		return "后台用户"
	case 2:
		return "手机端用户"
	default:
		return "未知"
	}
}

// GetOperatorTypeOptions 获取操作类别选项
func GetOperatorTypeOptions() []map[string]interface{} {
	return []map[string]interface{}{
		{"value": 0, "label": GetOperatorTypeName(0)},
		{"value": 1, "label": GetOperatorTypeName(1)},
		{"value": 2, "label": GetOperatorTypeName(2)},
	}
}

// GetStatusName 获取操作状态名称
func GetStatusName(status int) string {
	switch status {
	case 0:
		return "正常"
	case 1:
		return "异常"
	default:
		return "未知"
	}
}

// GetStatusOptions 获取操作状态选项
func GetStatusOptions() []map[string]interface{} {
	return []map[string]interface{}{
		{"value": 0, "label": GetStatusName(0)},
		{"value": 1, "label": GetStatusName(1)},
	}
}

// GetOperLogStatistics 获取操作日志统计信息
func (s *Service) GetOperLogStatistics() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 今日操作总数
	var todayTotal int64
	err := s.db.Model(&OperLog{}).
		Where("DATE(oper_time) = CURDATE()").
		Count(&todayTotal).Error
	if err != nil {
		return nil, err
	}
	stats["todayTotal"] = todayTotal

	// 今日异常数
	var todayError int64
	err = s.db.Model(&OperLog{}).
		Where("DATE(oper_time) = CURDATE() AND status = 1").
		Count(&todayError).Error
	if err != nil {
		return nil, err
	}
	stats["todayError"] = todayError

	// 本周操作总数
	var weekTotal int64
	err = s.db.Model(&OperLog{}).
		Where("YEARWEEK(oper_time) = YEARWEEK(NOW())").
		Count(&weekTotal).Error
	if err != nil {
		return nil, err
	}
	stats["weekTotal"] = weekTotal

	// 本月操作总数
	var monthTotal int64
	err = s.db.Model(&OperLog{}).
		Where("DATE_FORMAT(oper_time, '%Y-%m') = DATE_FORMAT(NOW(), '%Y-%m')").
		Count(&monthTotal).Error
	if err != nil {
		return nil, err
	}
	stats["monthTotal"] = monthTotal

	return stats, nil
}
