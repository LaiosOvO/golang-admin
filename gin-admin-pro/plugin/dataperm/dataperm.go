package dataperm

import (
	"fmt"

	"gorm.io/gorm"
)

// 数据权限范围常量
const (
	DataScopeAll       = 1 // 全部数据权限
	DataScopeCustom    = 2 // 自定义数据权限
	DataScopeDept      = 3 // 本部门数据权限
	DataScopeDeptChild = 4 // 本部门及以下数据权限
	DataScopeSelf      = 5 // 仅本人数据权限
)

// DataPermission 数据权限接口
type DataPermission interface {
	// GetDataScope 获取用户的数据权限范围
	GetDataScope(userID uint) (int, error)
	// GetDataScopeDeptIDs 获取用户数据权限的部门ID列表
	GetDataScopeDeptIDs(userID uint) ([]uint, error)
	// BuildDataScopeSQL 构建数据权限SQL条件
	BuildDataScopeSQL(db *gorm.DB, userID uint, deptAlias, userAlias string) *gorm.DB
}

// Service 数据权限服务
type Service struct {
	db *gorm.DB
}

// NewService 创建数据权限服务
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// GetDataScope 获取用户的数据权限范围
func (s *Service) GetDataScope(userID uint) (int, error) {
	type Result struct {
		DataScope int
	}

	var result Result
	err := s.db.Raw(`
		SELECT COALESCE(MIN(r.data_scope), ?) as data_scope
	 FROM user_role ur
	 JOIN role r ON ur.role_id = r.id
	 WHERE ur.user_id = ? AND r.status = 1
	`, DataScopeSelf, userID).Scan(&result).Error

	if err != nil {
		return DataScopeSelf, fmt.Errorf("failed to get user data scope: %w", err)
	}

	return result.DataScope, nil
}

// GetDataScopeDeptIDs 获取用户数据权限的部门ID列表
func (s *Service) GetDataScopeDeptIDs(userID uint) ([]uint, error) {
	// 获取用户的部门信息
	var userDept struct {
		DeptID     uint
		DeptPath   string
		DeptParent uint
	}

	err := s.db.Table("user").
		Select("dept_id").
		Where("id = ?", userID).
		Scan(&userDept).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get user dept: %w", err)
	}

	// 获取用户的数据权限范围
	dataScope, err := s.GetDataScope(userID)
	if err != nil {
		return nil, err
	}

	switch dataScope {
	case DataScopeAll:
		// 全部数据权限，返回空列表表示不限制
		return nil, nil
	case DataScopeDept:
		// 本部门数据权限
		if userDept.DeptID == 0 {
			return []uint{}, nil
		}
		return []uint{userDept.DeptID}, nil
	case DataScopeDeptChild:
		// 本部门及以下数据权限
		if userDept.DeptID == 0 {
			return []uint{}, nil
		}
		return s.getDeptChildIDs(userDept.DeptID)
	case DataScopeSelf:
		// 仅本人数据权限，在查询时需要特殊处理
		return nil, nil
	case DataScopeCustom:
		// 自定义数据权限，查询用户可访问的部门列表
		return s.getCustomDeptIDs(userID)
	default:
		return []uint{}, nil
	}
}

// getDeptChildIDs 获取部门及其所有子部门ID
func (s *Service) getDeptChildIDs(deptID uint) ([]uint, error) {
	var deptIDs []uint

	// 查询当前部门及所有子部门
	err := s.db.Raw(`
		WITH RECURSIVE dept_tree AS (
			SELECT id FROM dept WHERE id = ? AND status = 1
			UNION ALL
			SELECT d.id FROM dept d
			INNER JOIN dept_tree dt ON d.parent_id = dt.id
			WHERE d.status = 1
		)
		SELECT id FROM dept_tree
	`, deptID).Scan(&deptIDs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get dept child ids: %w", err)
	}

	return deptIDs, nil
}

// getCustomDeptIDs 获取用户自定义数据权限的部门ID列表
func (s *Service) getCustomDeptIDs(userID uint) ([]uint, error) {
	var deptIDs []uint

	// 查询用户可访问的部门（通过角色关联的部门）
	err := s.db.Raw(`
		SELECT DISTINCT rd.dept_id
		FROM user_role ur
		JOIN role r ON ur.role_id = r.id
		JOIN role_dept rd ON r.id = rd.role_id
		JOIN dept d ON rd.dept_id = d.id
		WHERE ur.user_id = ? AND r.status = 1 AND d.status = 1
	`, userID).Scan(&deptIDs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get custom dept ids: %w", err)
	}

	return deptIDs, nil
}

// BuildDataScopeSQL 构建数据权限SQL条件
func (s *Service) BuildDataScopeSQL(db *gorm.DB, userID uint, deptAlias, userAlias string) *gorm.DB {
	dataScope, err := s.GetDataScope(userID)
	if err != nil {
		return db.Where("1 = 0") // 如果获取权限失败，不返回任何数据
	}

	switch dataScope {
	case DataScopeAll:
		// 全部数据权限，不添加任何条件
		return db
	case DataScopeDept:
		// 本部门数据权限
		if deptAlias == "" {
			deptAlias = "dept_id"
		}
		return s.buildDeptScopeSQL(db, userID, deptAlias, false)
	case DataScopeDeptChild:
		// 本部门及以下数据权限
		if deptAlias == "" {
			deptAlias = "dept_id"
		}
		return s.buildDeptScopeSQL(db, userID, deptAlias, true)
	case DataScopeSelf:
		// 仅本人数据权限
		if userAlias == "" {
			userAlias = "id"
		}
		return db.Where(fmt.Sprintf("%s = ?", userAlias), userID)
	case DataScopeCustom:
		// 自定义数据权限
		deptIDs, err := s.getCustomDeptIDs(userID)
		if err != nil || len(deptIDs) == 0 {
			return db.Where("1 = 0")
		}
		if deptAlias == "" {
			deptAlias = "dept_id"
		}
		return db.Where(fmt.Sprintf("%s IN ?", deptAlias), deptIDs)
	default:
		return db.Where("1 = 0")
	}
}

// buildDeptScopeSQL 构建部门范围SQL
func (s *Service) buildDeptScopeSQL(db *gorm.DB, userID uint, deptAlias string, includeChild bool) *gorm.DB {
	// 获取用户部门ID
	var userDept struct {
		DeptID uint
	}

	err := s.db.Table("user").
		Select("dept_id").
		Where("id = ?", userID).
		Scan(&userDept).Error

	if err != nil || userDept.DeptID == 0 {
		return db.Where("1 = 0")
	}

	if includeChild {
		// 本部门及以下
		deptIDs, err := s.getDeptChildIDs(userDept.DeptID)
		if err != nil || len(deptIDs) == 0 {
			return db.Where("1 = 0")
		}
		return db.Where(fmt.Sprintf("%s IN ?", deptAlias), deptIDs)
	} else {
		// 仅本部门
		return db.Where(fmt.Sprintf("%s = ?", deptAlias), userDept.DeptID)
	}
}

// CheckDataPermission 检查用户是否有访问指定数据的权限
func (s *Service) CheckDataPermission(userID uint, dataDeptID uint, dataType string) bool {
	dataScope, err := s.GetDataScope(userID)
	if err != nil {
		return false
	}

	switch dataScope {
	case DataScopeAll:
		return true
	case DataScopeSelf:
		// 对于仅本人权限，需要检查数据是否属于当前用户
		return dataType == "user" && dataDeptID == userID
	case DataScopeDept, DataScopeDeptChild:
		// 获取用户可访问的部门列表
		deptIDs, err := s.GetDataScopeDeptIDs(userID)
		if err != nil {
			return false
		}
		for _, deptID := range deptIDs {
			if deptID == dataDeptID {
				return true
			}
		}
		return false
	case DataScopeCustom:
		// 自定义权限，检查数据部门是否在用户可访问的部门列表中
		deptIDs, err := s.getCustomDeptIDs(userID)
		if err != nil {
			return false
		}
		for _, deptID := range deptIDs {
			if deptID == dataDeptID {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// Scope 数据权限作用域，用于GORM查询
func Scope(userID uint, deptAlias, userAlias string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		service := NewService(db)
		return service.BuildDataScopeSQL(db, userID, deptAlias, userAlias)
	}
}

// HasDeptDataPermission 检查部门数据权限
func HasDeptDataPermission(userID uint, deptID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		service := NewService(db)
		if service.CheckDataPermission(userID, deptID, "dept") {
			return db
		}
		return db.Where("1 = 0")
	}
}

// HasUserDataPermission 检查用户数据权限
func HasUserDataPermission(userID uint, targetUserID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		service := NewService(db)
		if service.CheckDataPermission(userID, targetUserID, "user") {
			return db
		}
		return db.Where("1 = 0")
	}
}

// GetDataScopeName 获取数据权限范围名称
func GetDataScopeName(dataScope int) string {
	switch dataScope {
	case DataScopeAll:
		return "全部数据权限"
	case DataScopeCustom:
		return "自定义数据权限"
	case DataScopeDept:
		return "本部门数据权限"
	case DataScopeDeptChild:
		return "本部门及以下数据权限"
	case DataScopeSelf:
		return "仅本人数据权限"
	default:
		return "未知权限"
	}
}

// GetDataScopeOptions 获取数据权限范围选项
func GetDataScopeOptions() []map[string]interface{} {
	return []map[string]interface{}{
		{"value": DataScopeAll, "label": GetDataScopeName(DataScopeAll)},
		{"value": DataScopeCustom, "label": GetDataScopeName(DataScopeCustom)},
		{"value": DataScopeDept, "label": GetDataScopeName(DataScopeDept)},
		{"value": DataScopeDeptChild, "label": GetDataScopeName(DataScopeDeptChild)},
		{"value": DataScopeSelf, "label": GetDataScopeName(DataScopeSelf)},
	}
}

// RoleDept 角色部门关联表（用于自定义数据权限）
type RoleDept struct {
	ID     uint `gorm:"primarykey" json:"id"`
	RoleID uint `gorm:"not null;index" json:"roleId"`
	DeptID uint `gorm:"not null;index" json:"deptId"`
}
