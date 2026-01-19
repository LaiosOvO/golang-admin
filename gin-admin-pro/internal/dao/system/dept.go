package system

import (
	"errors"
	"gin-admin-pro/internal/model"
	"gin-admin-pro/internal/model/system"
	"strings"

	"gorm.io/gorm"
)

// DeptDAO 部门数据访问层
type DeptDAO struct {
	db *gorm.DB
}

// NewDeptDAO 创建部门DAO实例
func NewDeptDAO(db *gorm.DB) *DeptDAO {
	return &DeptDAO{db: db}
}

// DeptListReq 部门列表请求
type DeptListReq struct {
	Name   string `json:"name"`
	Status *int   `json:"status"`
}

// DeptResp 部门响应
type DeptResp struct {
	ID         uint       `json:"id"`
	Name       string     `json:"name"`
	ParentID   uint       `json:"parentId"`
	Level      int        `json:"level"`
	Sort       int        `json:"sort"`
	Leader     uint       `json:"leader"`
	LeaderName string     `json:"leaderName"`
	Phone      string     `json:"phone"`
	Email      string     `json:"email"`
	Status     int        `json:"status"`
	Children   []DeptResp `json:"children"`
	CreateTime int64      `json:"createTime"`
}

// DeptDetailResp 部门详情响应
type DeptDetailResp struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	ParentID uint   `json:"parentId"`
	Level    int    `json:"level"`
	Sort     int    `json:"sort"`
	Leader   uint   `json:"leader"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Status   int    `json:"status"`
	Remark   string `json:"remark"`
}

// DeptSimpleResp 部门简单响应
type DeptSimpleResp struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// CreateDeptReq 创建部门请求
type CreateDeptReq struct {
	ParentID uint   `json:"parentId" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Sort     int    `json:"sort" binding:"required"`
	Leader   uint   `json:"leader"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Status   int    `json:"status" binding:"required"`
	Remark   string `json:"remark"`
}

// UpdateDeptReq 更新部门请求
type UpdateDeptReq struct {
	ID       uint   `json:"id" binding:"required"`
	ParentID *uint  `json:"parentId"`
	Name     string `json:"name"`
	Sort     *int   `json:"sort"`
	Leader   *uint  `json:"leader"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Status   *int   `json:"status"`
	Remark   string `json:"remark"`
}

// GetList 获取部门列表（树形结构）
func (dao *DeptDAO) GetList(req *DeptListReq) ([]DeptResp, error) {
	query := dao.db.Model(&system.Dept{}).Preload("Leader")

	// 名称模糊查询
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}

	// 状态查询
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	var depts []system.Dept
	if err := query.Order("sort ASC, id ASC").Find(&depts).Error; err != nil {
		return nil, err
	}

	return dao.buildDeptTree(depts, 0), nil
}

// GetAllSimpleList 获取所有部门简单列表
func (dao *DeptDAO) GetAllSimpleList() ([]DeptSimpleResp, error) {
	var depts []system.Dept
	if err := dao.db.Model(&system.Dept{}).
		Select("id, name").
		Where("status = ?", 1).
		Order("sort ASC, id ASC").
		Find(&depts).Error; err != nil {
		return nil, err
	}

	resp := make([]DeptSimpleResp, len(depts))
	for i, dept := range depts {
		resp[i] = DeptSimpleResp{
			ID:   dept.ID,
			Name: dept.Name,
		}
	}

	return resp, nil
}

// GetByID 根据ID获取部门详情
func (dao *DeptDAO) GetByID(id uint) (*DeptDetailResp, error) {
	var dept system.Dept
	if err := dao.db.Preload("Leader").First(&dept, id).Error; err != nil {
		return nil, err
	}

	return &DeptDetailResp{
		ID:       dept.ID,
		Name:     dept.Name,
		ParentID: dept.ParentID,
		Level:    dept.Level,
		Sort:     dept.Sort,
		Leader:   dept.LeaderUserId,
		Phone:    dept.Phone,
		Email:    dept.Email,
		Status:   dept.Status,
		Remark:   dept.Remark,
	}, nil
}

// Create 创建部门
func (dao *DeptDAO) Create(req *CreateDeptReq, createBy uint) (uint, error) {
	// 获取父部门信息
	var parent system.Dept
	level := 1
	ancestors := "0"

	if req.ParentID != 0 {
		if err := dao.db.First(&parent, req.ParentID).Error; err != nil {
			return 0, err
		}
		level = parent.Level + 1
		ancestors = parent.Ancestors + "," + string(rune(parent.ID))
	}

	dept := system.Dept{
		TreeModel: model.TreeModel{
			ParentID:  req.ParentID,
			Level:     level,
			Sort:      req.Sort,
			Name:      req.Name,
			Ancestors: ancestors,
			AuditModel: model.AuditModel{
				CreateBy: createBy,
				Remark:   req.Remark,
			},
		},
		LeaderUserId: req.Leader,
		Phone:        req.Phone,
		Email:        req.Email,
		Status:       req.Status,
	}

	if err := dao.db.Create(&dept).Error; err != nil {
		return 0, err
	}

	// 更新路径
	dept.Path = dao.buildPath(dept.ID, req.ParentID)
	dao.db.Model(&dept).Update("path", dept.Path)

	return dept.ID, nil
}

// Update 更新部门
func (dao *DeptDAO) Update(req *UpdateDeptReq, updateBy uint) error {
	// 开始事务
	tx := dao.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 获取原部门信息
	var dept system.Dept
	if err := tx.First(&dept, req.ID).Error; err != nil {
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
			var parent system.Dept
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
	if req.Leader != nil {
		updateData["leader_user_id"] = *req.Leader
	}
	if req.Phone != "" {
		updateData["phone"] = req.Phone
	}
	if req.Email != "" {
		updateData["email"] = req.Email
	}
	if req.Status != nil {
		updateData["status"] = *req.Status
	}
	if req.Remark != "" {
		updateData["remark"] = req.Remark
	}

	if err := tx.Model(&system.Dept{}).Where("id = ?", req.ID).Updates(updateData).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 如果父级变更，需要更新子部门的层级和祖先
	if req.ParentID != nil && *req.ParentID != dept.ParentID {
		if err := dao.updateChildrenTree(tx, req.ID); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// Delete 删除部门
func (dao *DeptDAO) Delete(id uint) error {
	// 检查是否有子部门
	var count int64
	if err := dao.db.Model(&system.Dept{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return gorm.ErrRecordNotFound // 使用已存在部门下有子部门
	}

	// 检查部门下是否有用户
	var userCount int64
	if err := dao.db.Model(&system.User{}).Where("dept_id = ?", id).Count(&userCount).Error; err != nil {
		return err
	}
	if userCount > 0 {
		return errors.New("部门下存在用户，无法删除")
	}

	return dao.db.Delete(&system.Dept{}, id).Error
}

// GetUsersByDept 获取部门下的用户列表
func (dao *DeptDAO) GetUsersByDept(deptID uint) ([]system.User, error) {
	var users []system.User
	err := dao.db.Model(&system.User{}).
		Preload("Dept").
		Where("dept_id = ? AND status = ?", deptID, 1).
		Find(&users).Error

	return users, err
}

// buildDeptTree 构建部门树
func (dao *DeptDAO) buildDeptTree(depts []system.Dept, parentID uint) []DeptResp {
	var tree []DeptResp

	for _, dept := range depts {
		if dept.ParentID == parentID {
			node := DeptResp{
				ID:         dept.ID,
				Name:       dept.Name,
				ParentID:   dept.ParentID,
				Level:      dept.Level,
				Sort:       dept.Sort,
				Leader:     dept.LeaderUserId,
				Phone:      dept.Phone,
				Email:      dept.Email,
				Status:     dept.Status,
				CreateTime: dept.CreatedAt.Unix(),
			}

			// 设置负责人姓名
			if dept.Leader != nil {
				node.LeaderName = dept.Leader.Nickname
			}

			// 递归构建子部门
			children := dao.buildDeptTree(depts, dept.ID)
			if len(children) > 0 {
				node.Children = children
			}

			tree = append(tree, node)
		}
	}

	return tree
}

// buildPath 构建部门路径
func (dao *DeptDAO) buildPath(id, parentID uint) string {
	if parentID == 0 {
		return "/" + string(rune(id))
	}

	var parent system.Dept
	if err := dao.db.First(&parent, parentID).Error; err != nil {
		return "/" + string(rune(id))
	}

	return parent.Path + "/" + string(rune(id))
}

// updateChildrenTree 更新子部门树
func (dao *DeptDAO) updateChildrenTree(tx *gorm.DB, parentID uint) error {
	var children []system.Dept
	if err := tx.Find(&children, "parent_id = ?", parentID).Error; err != nil {
		return err
	}

	for _, child := range children {
		// 重新计算层级和祖先
		var parent system.Dept
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

		// 递归更新子部门
		if err := dao.updateChildrenTree(tx, child.ID); err != nil {
			return err
		}
	}

	return nil
}

// CheckNameExists 检查部门名称是否存在（同级下唯一）
func (dao *DeptDAO) CheckNameExists(name string, parentID uint, excludeID *uint) (bool, error) {
	query := dao.db.Model(&system.Dept{}).
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
func (dao *DeptDAO) GetMaxSort(parentID uint) (int, error) {
	var maxSort int
	err := dao.db.Model(&system.Dept{}).
		Where("parent_id = ?", parentID).
		Select("COALESCE(MAX(sort), 0)").
		Scan(&maxSort).Error

	return maxSort, err
}

// GetParentChain 获取父级部门链
func (dao *DeptDAO) GetParentChain(deptID uint) ([]system.Dept, error) {
	var dept system.Dept
	if err := dao.db.First(&dept, deptID).Error; err != nil {
		return nil, err
	}

	if dept.Ancestors == "" || dept.Ancestors == "0" {
		return []system.Dept{}, nil
	}

	ancestorIDs := strings.Split(dept.Ancestors, ",")
	var depts []system.Dept

	for _, idStr := range ancestorIDs {
		if idStr == "0" {
			continue
		}

		var ancestor system.Dept
		if err := dao.db.First(&ancestor, idStr).Error; err != nil {
			return nil, err
		}
		depts = append(depts, ancestor)
	}

	return depts, nil
}
