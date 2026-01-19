package system

import (
	"gin-admin-pro/internal/model"
	"gin-admin-pro/internal/model/system"
	"strings"

	"gorm.io/gorm"
)

// MenuDAO 菜单数据访问层
type MenuDAO struct {
	db *gorm.DB
}

// NewMenuDAO 创建菜单DAO实例
func NewMenuDAO(db *gorm.DB) *MenuDAO {
	return &MenuDAO{db: db}
}

// MenuListReq 菜单列表请求
type MenuListReq struct {
	Name   string `json:"name"`
	Status *int   `json:"status"`
	Type   *int   `json:"type"`
}

// MenuResp 菜单响应
type MenuResp struct {
	ID            uint       `json:"id"`
	Name          string     `json:"name"`
	ParentID      uint       `json:"parentId"`
	Level         int        `json:"level"`
	Sort          int        `json:"sort"`
	Path          string     `json:"path"`
	Component     string     `json:"component"`
	ComponentName string     `json:"componentName"`
	Icon          string     `json:"icon"`
	Type          int        `json:"type"`
	Perms         string     `json:"perms"`
	Status        int        `json:"status"`
	Visible       int        `json:"visible"`
	KeepAlive     int        `json:"keepAlive"`
	AlwaysShow    int        `json:"alwaysShow"`
	Children      []MenuResp `json:"children"`
	CreateTime    int64      `json:"createTime"`
}

// MenuDetailResp 菜单详情响应
type MenuDetailResp struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	ParentID      uint   `json:"parentId"`
	Level         int    `json:"level"`
	Sort          int    `json:"sort"`
	Path          string `json:"path"`
	Component     string `json:"component"`
	ComponentName string `json:"componentName"`
	Icon          string `json:"icon"`
	Type          int    `json:"type"`
	Perms         string `json:"perms"`
	Status        int    `json:"status"`
	Visible       int    `json:"visible"`
	KeepAlive     int    `json:"keepAlive"`
	AlwaysShow    int    `json:"alwaysShow"`
	Remark        string `json:"remark"`
}

// MenuSimpleResp 菜单简单响应（用于角色授权）
type MenuSimpleResp struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type int    `json:"type"`
}

// CreateMenuReq 创建菜单请求
type CreateMenuReq struct {
	ParentID      uint   `json:"parentId" binding:"required"`
	Name          string `json:"name" binding:"required"`
	Sort          int    `json:"sort" binding:"required"`
	Path          string `json:"path" binding:"required"`
	Component     string `json:"component"`
	ComponentName string `json:"componentName"`
	Icon          string `json:"icon"`
	Type          int    `json:"type" binding:"required"`
	Perms         string `json:"perms"`
	Status        int    `json:"status" binding:"required"`
	Visible       int    `json:"visible"`
	KeepAlive     int    `json:"keepAlive"`
	AlwaysShow    int    `json:"alwaysShow"`
	Remark        string `json:"remark"`
}

// UpdateMenuReq 更新菜单请求
type UpdateMenuReq struct {
	ID            uint   `json:"id" binding:"required"`
	ParentID      *uint  `json:"parentId"`
	Name          string `json:"name"`
	Sort          *int   `json:"sort"`
	Path          string `json:"path"`
	Component     string `json:"component"`
	ComponentName string `json:"componentName"`
	Icon          string `json:"icon"`
	Type          *int   `json:"type"`
	Perms         string `json:"perms"`
	Status        *int   `json:"status"`
	Visible       *int   `json:"visible"`
	KeepAlive     *int   `json:"keepAlive"`
	AlwaysShow    *int   `json:"alwaysShow"`
	Remark        string `json:"remark"`
}

// GetList 获取菜单列表（树形结构）
func (dao *MenuDAO) GetList(req *MenuListReq) ([]MenuResp, error) {
	query := dao.db.Model(&system.Menu{})

	// 名称模糊查询
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}

	// 状态查询
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// 类型查询
	if req.Type != nil {
		query = query.Where("type = ?", *req.Type)
	}

	var menus []system.Menu
	if err := query.Order("sort ASC, id ASC").Find(&menus).Error; err != nil {
		return nil, err
	}

	return dao.buildMenuTree(menus, 0), nil
}

// GetAllSimpleList 获取所有菜单简单列表（用于角色授权）
func (dao *MenuDAO) GetAllSimpleList() ([]MenuSimpleResp, error) {
	var menus []system.Menu
	if err := dao.db.Model(&system.Menu{}).
		Select("id, name, type").
		Where("status = ?", 1).
		Order("sort ASC, id ASC").
		Find(&menus).Error; err != nil {
		return nil, err
	}

	resp := make([]MenuSimpleResp, len(menus))
	for i, menu := range menus {
		resp[i] = MenuSimpleResp{
			ID:   menu.ID,
			Name: menu.Name,
			Type: menu.Type,
		}
	}

	return resp, nil
}

// GetByID 根据ID获取菜单详情
func (dao *MenuDAO) GetByID(id uint) (*MenuDetailResp, error) {
	var menu system.Menu
	if err := dao.db.First(&menu, id).Error; err != nil {
		return nil, err
	}

	return &MenuDetailResp{
		ID:            menu.ID,
		Name:          menu.Name,
		ParentID:      menu.ParentID,
		Level:         menu.Level,
		Sort:          menu.Sort,
		Path:          menu.Path,
		Component:     menu.Component,
		ComponentName: menu.ComponentName,
		Icon:          menu.Icon,
		Type:          menu.Type,
		Perms:         menu.Perms,
		Status:        menu.Status,
		Visible:       menu.Visible,
		KeepAlive:     menu.KeepAlive,
		AlwaysShow:    menu.AlwaysShow,
		Remark:        menu.Remark,
	}, nil
}

// Create 创建菜单
func (dao *MenuDAO) Create(req *CreateMenuReq, createBy uint) (uint, error) {
	// 获取父菜单信息
	var parent system.Menu
	level := 1
	ancestors := "0"

	if req.ParentID != 0 {
		if err := dao.db.First(&parent, req.ParentID).Error; err != nil {
			return 0, err
		}
		level = parent.Level + 1
		ancestors = parent.Ancestors + "," + string(rune(parent.ID))
	}

	menu := system.Menu{
		TreeModel: model.TreeModel{
			ParentID:  req.ParentID,
			Level:     level,
			Sort:      req.Sort,
			Name:      req.Name,
			Path:      req.Path,
			Ancestors: ancestors,
			AuditModel: model.AuditModel{
				CreateBy: createBy,
				Remark:   req.Remark,
			},
		},
		Type:          req.Type,
		Icon:          req.Icon,
		Component:     req.Component,
		ComponentName: req.ComponentName,
		Perms:         req.Perms,
		Status:        req.Status,
		Visible:       req.Visible,
		KeepAlive:     req.KeepAlive,
		AlwaysShow:    req.AlwaysShow,
	}

	if err := dao.db.Create(&menu).Error; err != nil {
		return 0, err
	}

	// 更新路径
	menu.Path = dao.buildPath(menu.ID, req.ParentID)
	dao.db.Model(&menu).Update("path", menu.Path)

	return menu.ID, nil
}

// Update 更新菜单
func (dao *MenuDAO) Update(req *UpdateMenuReq, updateBy uint) error {
	// 开始事务
	tx := dao.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 获取原菜单信息
	var menu system.Menu
	if err := tx.First(&menu, req.ID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新数据
	updateData := map[string]interface{}{
		"update_by": updateBy,
	}

	if req.ParentID != nil {
		updateData["parent_id"] = *req.ParentID

		// 重新计算层级和祖先
		if *req.ParentID != 0 {
			var parent system.Menu
			if err := tx.First(&parent, *req.ParentID).Error; err != nil {
				tx.Rollback()
				return err
			}
			updateData["level"] = parent.Level + 1
			updateData["ancestors"] = parent.Ancestors + "," + string(rune(parent.ID))
		} else {
			updateData["level"] = 1
			updateData["ancestors"] = "0"
		}
	}

	if req.Name != "" {
		updateData["name"] = req.Name
	}
	if req.Sort != nil {
		updateData["sort"] = *req.Sort
	}
	if req.Path != "" {
		updateData["path"] = req.Path
	}
	if req.Component != "" {
		updateData["component"] = req.Component
	}
	if req.ComponentName != "" {
		updateData["component_name"] = req.ComponentName
	}
	if req.Icon != "" {
		updateData["icon"] = req.Icon
	}
	if req.Type != nil {
		updateData["type"] = *req.Type
	}
	if req.Perms != "" {
		updateData["perms"] = req.Perms
	}
	if req.Status != nil {
		updateData["status"] = *req.Status
	}
	if req.Visible != nil {
		updateData["visible"] = *req.Visible
	}
	if req.KeepAlive != nil {
		updateData["keep_alive"] = *req.KeepAlive
	}
	if req.AlwaysShow != nil {
		updateData["always_show"] = *req.AlwaysShow
	}
	if req.Remark != "" {
		updateData["remark"] = req.Remark
	}

	if err := tx.Model(&system.Menu{}).Where("id = ?", req.ID).Updates(updateData).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 如果父级变更，需要更新子菜单的层级和祖先
	if req.ParentID != nil && *req.ParentID != menu.ParentID {
		if err := dao.updateChildrenTree(tx, req.ID); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// Delete 删除菜单
func (dao *MenuDAO) Delete(id uint) error {
	// 检查是否有子菜单
	var count int64
	if err := dao.db.Model(&system.Menu{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return gorm.ErrRecordNotFound // 使用已存在菜单下有子菜单
	}

	// 删除角色菜单关联
	if err := dao.db.Where("menu_id = ?", id).Delete(&system.RoleMenu{}).Error; err != nil {
		return err
	}

	return dao.db.Delete(&system.Menu{}, id).Error
}

// GetUserMenus 获取用户菜单列表
func (dao *MenuDAO) GetUserMenus(userID uint) ([]MenuResp, error) {
	query := `
		SELECT DISTINCT m.* FROM system_menu m
		INNER JOIN system_role_menu rm ON m.id = rm.menu_id
		INNER JOIN system_user_role ur ON rm.role_id = ur.role_id
		INNER JOIN system_role r ON rm.role_id = r.id
		WHERE ur.user_id = ? 
		AND m.status = 1 
		AND r.status = 1
		AND m.visible = 1
		ORDER BY m.sort ASC, m.id ASC
	`

	var menus []system.Menu
	if err := dao.db.Raw(query, userID).Find(&menus).Error; err != nil {
		return nil, err
	}

	return dao.buildMenuTree(menus, 0), nil
}

// buildMenuTree 构建菜单树
func (dao *MenuDAO) buildMenuTree(menus []system.Menu, parentID uint) []MenuResp {
	var tree []MenuResp

	for _, menu := range menus {
		if menu.ParentID == parentID {
			node := MenuResp{
				ID:            menu.ID,
				Name:          menu.Name,
				ParentID:      menu.ParentID,
				Level:         menu.Level,
				Sort:          menu.Sort,
				Path:          menu.Path,
				Component:     menu.Component,
				ComponentName: menu.ComponentName,
				Icon:          menu.Icon,
				Type:          menu.Type,
				Perms:         menu.Perms,
				Status:        menu.Status,
				Visible:       menu.Visible,
				KeepAlive:     menu.KeepAlive,
				AlwaysShow:    menu.AlwaysShow,
				CreateTime:    menu.CreatedAt.Unix(),
			}

			// 递归构建子菜单
			children := dao.buildMenuTree(menus, menu.ID)
			if len(children) > 0 {
				node.Children = children
			}

			tree = append(tree, node)
		}
	}

	return tree
}

// buildPath 构建菜单路径
func (dao *MenuDAO) buildPath(id, parentID uint) string {
	if parentID == 0 {
		return "/" + string(rune(id))
	}

	var parent system.Menu
	if err := dao.db.First(&parent, parentID).Error; err != nil {
		return "/" + string(rune(id))
	}

	return parent.Path + "/" + string(rune(id))
}

// updateChildrenTree 更新子菜单树
func (dao *MenuDAO) updateChildrenTree(tx *gorm.DB, parentID uint) error {
	var children []system.Menu
	if err := tx.Find(&children, "parent_id = ?", parentID).Error; err != nil {
		return err
	}

	for _, child := range children {
		// 重新计算层级和祖先
		var parent system.Menu
		if err := tx.First(&parent, parentID).Error; err != nil {
			return err
		}

		updateData := map[string]interface{}{
			"level":     parent.Level + 1,
			"ancestors": parent.Ancestors + "," + string(rune(parent.ID)),
		}

		if err := tx.Model(&child).Updates(updateData).Error; err != nil {
			return err
		}

		// 递归更新子菜单
		if err := dao.updateChildrenTree(tx, child.ID); err != nil {
			return err
		}
	}

	return nil
}

// CheckNameExists 检查菜单名称是否存在（同级下唯一）
func (dao *MenuDAO) CheckNameExists(name string, parentID uint, excludeID *uint) (bool, error) {
	query := dao.db.Model(&system.Menu{}).
		Where("name = ? AND parent_id = ?", name, parentID)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetMaxSort 获取同级下的最大排序值
func (dao *MenuDAO) GetMaxSort(parentID uint) (int, error) {
	var maxSort int
	err := dao.db.Model(&system.Menu{}).
		Where("parent_id = ?", parentID).
		Select("COALESCE(MAX(sort), 0)").
		Scan(&maxSort).Error

	return maxSort, err
}

// GetParentChain 获取父级菜单链
func (dao *MenuDAO) GetParentChain(menuID uint) ([]system.Menu, error) {
	var menu system.Menu
	if err := dao.db.First(&menu, menuID).Error; err != nil {
		return nil, err
	}

	if menu.Ancestors == "" || menu.Ancestors == "0" {
		return []system.Menu{}, nil
	}

	ancestorIDs := strings.Split(menu.Ancestors, ",")
	var menus []system.Menu

	for _, idStr := range ancestorIDs {
		if idStr == "0" {
			continue
		}

		var ancestor system.Menu
		if err := dao.db.First(&ancestor, idStr).Error; err != nil {
			return nil, err
		}
		menus = append(menus, ancestor)
	}

	return menus, nil
}
