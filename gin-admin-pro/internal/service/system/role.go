package system

import (
	"errors"
	"gin-admin-pro/internal/dao/system"
	"gin-admin-pro/internal/model"
	sysmodel "gin-admin-pro/internal/model/system"
	"gorm.io/gorm"
)

var (
	ErrRoleNotFound     = errors.New("角色不存在")
	ErrRoleCodeExists   = errors.New("角色代码已存在")
	ErrRoleHasUsers     = errors.New("角色下存在用户，无法删除")
	ErrRoleIsBuiltin    = errors.New("内置角色，无法删除")
	ErrInvalidDataScope = errors.New("无效的数据权限范围")
)

// RoleService 角色服务层
type RoleService struct {
	roleDAO *system.RoleDAO
}

// NewRoleService 创建角色服务实例
func NewRoleService(roleDAO *system.RoleDAO) *RoleService {
	return &RoleService{
		roleDAO: roleDAO,
	}
}

// GetPage 获取角色分页列表
func (rs *RoleService) GetPage(req *system.RolePageReq) (*model.PageResp, error) {
	roles, total, err := rs.roleDAO.GetPage(req)
	if err != nil {
		return nil, err
	}

	return &model.PageResp{
		List:  roles,
		Total: total,
	}, nil
}

// GetByID 根据ID获取角色详情
func (rs *RoleService) GetByID(id uint) (*system.RoleDetailResp, error) {
	role, err := rs.roleDAO.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}
	return role, nil
}

// GetAllSimple 获取所有角色精简列表
func (rs *RoleService) GetAllSimple() ([]*system.RoleSimpleResp, error) {
	return rs.roleDAO.GetAllSimple()
}

// Create 创建角色
func (rs *RoleService) Create(req *system.RoleCreateReq, createBy uint) error {
	// 检查角色代码是否已存在
	exists, err := rs.roleDAO.CheckCodeExists(req.Code, nil)
	if err != nil {
		return err
	}
	if exists {
		return ErrRoleCodeExists
	}

	// 验证数据权限范围
	if !rs.isValidDataScope(req.DataScope) {
		return ErrInvalidDataScope
	}

	return rs.roleDAO.Create(req, createBy)
}

// Update 更新角色
func (rs *RoleService) Update(req *system.RoleUpdateReq, updateBy uint) error {
	// 检查角色是否存在
	_, err := rs.roleDAO.GetByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRoleNotFound
		}
		return err
	}

	// 检查角色代码是否已存在（排除当前角色）
	exists, err := rs.roleDAO.CheckCodeExists(req.Code, &req.ID)
	if err != nil {
		return err
	}
	if exists {
		return ErrRoleCodeExists
	}

	// 验证数据权限范围
	if !rs.isValidDataScope(req.DataScope) {
		return ErrInvalidDataScope
	}

	return rs.roleDAO.Update(req, updateBy)
}

// Delete 删除角色
func (rs *RoleService) Delete(id uint) error {
	// 检查角色是否存在
	role, err := rs.roleDAO.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRoleNotFound
		}
		return err
	}

	// 内置角色不能删除
	if role.Type == 1 {
		return ErrRoleIsBuiltin
	}

	// TODO: 检查角色下是否有关联用户（需要实现角色-用户关联查询）
	// 这里暂时允许删除，后续可以加上用户关联检查

	return rs.roleDAO.Delete(id)
}

// AssignMenuPermissions 分配菜单权限
func (rs *RoleService) AssignMenuPermissions(roleID uint, menuIDs []uint) error {
	// 检查角色是否存在
	_, err := rs.roleDAO.GetByID(roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRoleNotFound
		}
		return err
	}

	return rs.roleDAO.AssignMenuPermissions(roleID, menuIDs)
}

// GetMenuPermissions 获取角色的菜单权限
func (rs *RoleService) GetMenuPermissions(roleID uint) ([]uint, error) {
	// 检查角色是否存在
	_, err := rs.roleDAO.GetByID(roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}

	return rs.roleDAO.GetMenuIDsByRoleID(roleID)
}

// isValidDataScope 验证数据权限范围是否有效
func (rs *RoleService) isValidDataScope(dataScope int) bool {
	validScopes := []int{1, 2, 3, 4, 5} // 1-全部 2-自定义 3-本部门 4-本部门及以下 5-仅本人
	for _, scope := range validScopes {
		if scope == dataScope {
			return true
		}
	}
	return false
}

// DataScopeDesc 数据权限范围描述
func (rs *RoleService) DataScopeDesc() map[int]string {
	return map[int]string{
		1: "全部数据权限",
		2: "自定义数据权限",
		3: "本部门数据权限",
		4: "本部门及以下数据权限",
		5: "仅本人数据权限",
	}
}

// RoleTypeDesc 角色类型描述
func (rs *RoleService) RoleTypeDesc() map[int]string {
	return map[int]string{
		1: "内置角色",
		2: "自定义角色",
	}
}

// StatusDesc 状态描述
func (rs *RoleService) StatusDesc() map[int]string {
	return map[int]string{
		0: "禁用",
		1: "启用",
	}
}

// GetRolesByUserID 获取用户的角色列表
func (rs *RoleService) GetRolesByUserID(userID uint) ([]*sysmodel.Role, error) {
	return rs.roleDAO.GetRolesByUserID(userID)
}
