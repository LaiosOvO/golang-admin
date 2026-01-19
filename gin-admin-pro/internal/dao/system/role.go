package system

import (
	"gin-admin-pro/internal/model"
	"gin-admin-pro/internal/model/system"
	"time"

	"gorm.io/gorm"
)

// RoleDAO 角色数据访问层
type RoleDAO struct {
	db *gorm.DB
}

// NewRoleDAO 创建角色DAO实例
func NewRoleDAO(db *gorm.DB) *RoleDAO {
	return &RoleDAO{db: db}
}

// RolePageReq 角色分页查询请求
type RolePageReq struct {
	model.PageReq
	Name       string   `json:"name"`
	Code       string   `json:"code"`
	Status     *int     `json:"status"`
	CreateTime []string `json:"createTime"`
}

// RolePageResp 角色分页查询响应
type RolePageResp struct {
	ID         uint      `json:"id"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	Sort       int       `json:"sort"`
	DataScope  int       `json:"dataScope"`
	Status     int       `json:"status"`
	Type       int       `json:"type"`
	Remark     string    `json:"remark"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}

// RoleDetailResp 角色详情响应
type RoleDetailResp struct {
	ID         uint      `json:"id"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	Sort       int       `json:"sort"`
	DataScope  int       `json:"dataScope"`
	Status     int       `json:"status"`
	Type       int       `json:"type"`
	Remark     string    `json:"remark"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
	CreateBy   uint      `json:"createBy"`
	UpdateBy   uint      `json:"updateBy"`
}

// RoleCreateReq 创建角色请求
type RoleCreateReq struct {
	Code      string `json:"code" binding:"required,max=100"`
	Name      string `json:"name" binding:"required,max=30"`
	Sort      int    `json:"sort"`
	DataScope int    `json:"dataScope"`
	Status    int    `json:"status"`
	Type      int    `json:"type"`
	Remark    string `json:"remark" binding:"max=500"`
}

// RoleUpdateReq 更新角色请求
type RoleUpdateReq struct {
	ID        uint   `json:"id" binding:"required"`
	Code      string `json:"code" binding:"required,max=100"`
	Name      string `json:"name" binding:"required,max=30"`
	Sort      int    `json:"sort"`
	DataScope int    `json:"dataScope"`
	Status    int    `json:"status"`
	Type      int    `json:"type"`
	Remark    string `json:"remark" binding:"max=500"`
}

// RoleSimpleResp 角色精简响应
type RoleSimpleResp struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// GetPage 获取角色分页列表
func (r *RoleDAO) GetPage(req *RolePageReq) ([]*RolePageResp, int64, error) {
	var roles []*system.Role
	var total int64

	db := r.db.Model(&system.Role{})

	// 查询条件
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Code != "" {
		db = db.Where("code LIKE ?", "%"+req.Code+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if len(req.CreateTime) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", req.CreateTime[0], req.CreateTime[1])
	}

	// 计算总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.PageNo - 1) * req.PageSize
	if err := db.Offset(offset).Limit(req.PageSize).Order("sort ASC, created_at DESC").Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	// 转换响应格式
	resps := make([]*RolePageResp, len(roles))
	for i, role := range roles {
		resps[i] = &RolePageResp{
			ID:         role.ID,
			Code:       role.Code,
			Name:       role.Name,
			Sort:       role.Sort,
			DataScope:  role.DataScope,
			Status:     role.Status,
			Type:       role.Type,
			Remark:     role.Remark,
			CreateTime: role.CreatedAt,
			UpdateTime: role.UpdatedAt,
		}
	}

	return resps, total, nil
}

// GetByID 根据ID获取角色
func (r *RoleDAO) GetByID(id uint) (*RoleDetailResp, error) {
	var role system.Role
	if err := r.db.First(&role, id).Error; err != nil {
		return nil, err
	}

	return &RoleDetailResp{
		ID:         role.ID,
		Code:       role.Code,
		Name:       role.Name,
		Sort:       role.Sort,
		DataScope:  role.DataScope,
		Status:     role.Status,
		Type:       role.Type,
		Remark:     role.Remark,
		CreateTime: role.CreatedAt,
		UpdateTime: role.UpdatedAt,
		CreateBy:   role.CreateBy,
		UpdateBy:   role.UpdateBy,
	}, nil
}

// GetAllSimple 获取所有角色精简列表
func (r *RoleDAO) GetAllSimple() ([]*RoleSimpleResp, error) {
	var roles []*system.Role
	if err := r.db.Where("status = ?", 1).Order("sort ASC").Find(&roles).Error; err != nil {
		return nil, err
	}

	resps := make([]*RoleSimpleResp, len(roles))
	for i, role := range roles {
		resps[i] = &RoleSimpleResp{
			ID:   role.ID,
			Code: role.Code,
			Name: role.Name,
		}
	}

	return resps, nil
}

// Create 创建角色
func (r *RoleDAO) Create(req *RoleCreateReq, createBy uint) error {
	role := &system.Role{
		Code:      req.Code,
		Name:      req.Name,
		Sort:      req.Sort,
		DataScope: req.DataScope,
		Status:    req.Status,
		Type:      req.Type,
		Remark:    req.Remark,
	}
	role.CreateBy = createBy

	return r.db.Create(role).Error
}

// Update 更新角色
func (r *RoleDAO) Update(req *RoleUpdateReq, updateBy uint) error {
	role := &system.Role{
		Code:      req.Code,
		Name:      req.Name,
		Sort:      req.Sort,
		DataScope: req.DataScope,
		Status:    req.Status,
		Type:      req.Type,
		Remark:    req.Remark,
	}
	role.UpdateBy = updateBy

	return r.db.Model(&system.Role{}).Where("id = ?", req.ID).Updates(role).Error
}

// Delete 删除角色
func (r *RoleDAO) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除角色菜单关联
		if err := tx.Where("role_id = ?", id).Delete(&system.RoleMenu{}).Error; err != nil {
			return err
		}

		// 删除用户角色关联
		if err := tx.Where("role_id = ?", id).Delete(&system.UserRole{}).Error; err != nil {
			return err
		}

		// 删除角色
		return tx.Delete(&system.Role{}, id).Error
	})
}

// GetByCode 根据代码获取角色
func (r *RoleDAO) GetByCode(code string) (*system.Role, error) {
	var role system.Role
	if err := r.db.Where("code = ?", code).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// CheckCodeExists 检查角色代码是否存在（排除指定ID）
func (r *RoleDAO) CheckCodeExists(code string, excludeID *uint) (bool, error) {
	var count int64
	db := r.db.Model(&system.Role{}).Where("code = ?", code)
	if excludeID != nil {
		db = db.Where("id != ?", *excludeID)
	}
	if err := db.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// AssignMenuPermissions 分配菜单权限
func (r *RoleDAO) AssignMenuPermissions(roleID uint, menuIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 先删除原有权限
		if err := tx.Where("role_id = ?", roleID).Delete(&system.RoleMenu{}).Error; err != nil {
			return err
		}

		// 添加新权限
		if len(menuIDs) > 0 {
			roleMenus := make([]*system.RoleMenu, len(menuIDs))
			for i, menuID := range menuIDs {
				roleMenus[i] = &system.RoleMenu{
					RoleID: roleID,
					MenuID: menuID,
				}
			}
			if err := tx.Create(&roleMenus).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetMenuIDsByRoleID 获取角色的菜单ID列表
func (r *RoleDAO) GetMenuIDsByRoleID(roleID uint) ([]uint, error) {
	var menuIDs []uint
	if err := r.db.Model(&system.RoleMenu{}).Where("role_id = ?", roleID).Pluck("menu_id", &menuIDs).Error; err != nil {
		return nil, err
	}
	return menuIDs, nil
}

// GetRolesByUserID 获取用户的角色列表
func (r *RoleDAO) GetRolesByUserID(userID uint) ([]*system.Role, error) {
	var roles []*system.Role
	if err := r.db.Joins("JOIN user_role ON user_role.role_id = role.id").
		Where("user_role.user_id = ? AND role.status = ?", userID, 1).
		Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}
