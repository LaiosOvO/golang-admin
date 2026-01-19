package system

import (
	"gin-admin-pro/internal/model"
	"gin-admin-pro/internal/model/system"
	"time"

	"gorm.io/gorm"
)

// UserDAO 用户数据访问层
type UserDAO struct {
	db *gorm.DB
}

// NewUserDAO 创建用户DAO实例
func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

// UserPageReq 用户分页查询请求
type UserPageReq struct {
	model.PageReq
	Username   string   `json:"username"`
	Mobile     string   `json:"mobile"`
	Email      string   `json:"email"`
	Status     *int     `json:"status"`
	DeptID     *uint    `json:"deptId"`
	CreateTime []string `json:"createTime"`
}

// UserPageResp 用户分页查询响应
type UserPageResp struct {
	ID         uint       `json:"id"`
	Username   string     `json:"username"`
	Nickname   string     `json:"nickname"`
	DeptID     uint       `json:"deptId"`
	DeptName   string     `json:"deptName"`
	Email      string     `json:"email"`
	Mobile     string     `json:"mobile"`
	Avatar     string     `json:"avatar"`
	Status     int        `json:"status"`
	LoginIP    string     `json:"loginIp"`
	LoginDate  *time.Time `json:"loginDate"`
	CreateTime time.Time  `json:"createTime"`
	UpdateTime time.Time  `json:"updateTime"`
}

// UserDetailResp 用户详情响应
type UserDetailResp struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	DeptID   uint   `json:"deptId"`
	PostIDs  []uint `json:"postIds"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Avatar   string `json:"avatar"`
	Status   int    `json:"status"`
	Remark   string `json:"remark"`
}

// UserSimpleResp 用户简单信息响应
type UserSimpleResp struct {
	ID       uint   `json:"id"`
	Nickname string `json:"nickname"`
	DeptName string `json:"deptName"`
}

// CreateReq 创建用户请求
type CreateReq struct {
	Username string `json:"username" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	Password string `json:"password" binding:"required"`
	DeptID   uint   `json:"deptId" binding:"required"`
	PostIDs  []uint `json:"postIds"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Avatar   string `json:"avatar"`
	Status   int    `json:"status" binding:"required"`
	Remark   string `json:"remark"`
}

// UpdateReq 更新用户请求
type UpdateReq struct {
	ID       uint   `json:"id" binding:"required"`
	Nickname string `json:"nickname"`
	DeptID   *uint  `json:"deptId"`
	PostIDs  []uint `json:"postIds"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Avatar   string `json:"avatar"`
	Status   *int   `json:"status"`
	Remark   string `json:"remark"`
}

// UpdatePasswordReq 更新密码请求
type UpdatePasswordReq struct {
	ID       uint   `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateStatusReq 更新状态请求
type UpdateStatusReq struct {
	ID     uint `json:"id" binding:"required"`
	Status int  `json:"status" binding:"required"`
}

// GetPage 获取用户分页列表
func (dao *UserDAO) GetPage(req *UserPageReq) ([]UserPageResp, int64, error) {
	var users []system.User
	var total int64

	query := dao.db.Model(&system.User{}).
		Preload("Dept")

	// 用户名模糊查询
	if req.Username != "" {
		query = query.Where("username LIKE ?", "%"+req.Username+"%")
	}

	// 手机号模糊查询
	if req.Mobile != "" {
		query = query.Where("mobile LIKE ?", "%"+req.Mobile+"%")
	}

	// 邮箱模糊查询
	if req.Email != "" {
		query = query.Where("email LIKE ?", "%"+req.Email+"%")
	}

	// 状态查询
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// 部门查询
	if req.DeptID != nil {
		query = query.Where("dept_id = ?", *req.DeptID)
	}

	// 创建时间范围查询
	if len(req.CreateTime) == 2 {
		query = query.Where("created_at BETWEEN ? AND ?", req.CreateTime[0], req.CreateTime[1])
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.PageNo - 1) * req.PageSize
	if err := query.
		Offset(offset).
		Limit(req.PageSize).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	// 转换响应格式
	resp := make([]UserPageResp, len(users))
	for i, user := range users {
		resp[i] = UserPageResp{
			ID:         user.ID,
			Username:   user.Username,
			Nickname:   user.Nickname,
			DeptID:     user.DeptID,
			Email:      user.Email,
			Mobile:     user.Mobile,
			Avatar:     user.Avatar,
			Status:     user.Status,
			LoginIP:    user.LoginIP,
			LoginDate:  &time.Time{},
			CreateTime: user.CreatedAt,
			UpdateTime: user.UpdatedAt,
		}

		if user.Dept != nil {
			resp[i].DeptName = user.Dept.Name
		}

		if user.LoginDate.Valid {
			resp[i].LoginDate = &user.LoginDate.Time
		}
	}

	return resp, total, nil
}

// GetByID 根据ID获取用户详情
func (dao *UserDAO) GetByID(id uint) (*UserDetailResp, error) {
	var user system.User
	if err := dao.db.Preload("Posts").First(&user, id).Error; err != nil {
		return nil, err
	}

	postIDs := make([]uint, len(user.Posts))
	for i, post := range user.Posts {
		postIDs[i] = post.ID
	}

	return &UserDetailResp{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		DeptID:   user.DeptID,
		PostIDs:  postIDs,
		Email:    user.Email,
		Mobile:   user.Mobile,
		Avatar:   user.Avatar,
		Status:   user.Status,
		Remark:   user.Remark,
	}, nil
}

// GetByUsername 根据用户名获取用户
func (dao *UserDAO) GetByUsername(username string) (*system.User, error) {
	var user system.User
	err := dao.db.Preload("Roles").First(&user, "username = ?", username).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (dao *UserDAO) Create(req *CreateReq, createBy uint) (uint, error) {
	user := system.User{
		Username: req.Username,
		Nickname: req.Nickname,
		Password: req.Password,
		DeptID:   req.DeptID,
		Email:    req.Email,
		Mobile:   req.Mobile,
		Avatar:   req.Avatar,
		Status:   req.Status,
		AuditModel: model.AuditModel{
			CreateBy: createBy,
			Remark:   req.Remark,
		},
	}

	// 开始事务
	tx := dao.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建用户
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// 关联岗位
	if len(req.PostIDs) > 0 {
		var posts []system.Post
		if err := tx.Find(&posts, req.PostIDs).Error; err != nil {
			tx.Rollback()
			return 0, err
		}

		if err := tx.Model(&user).Association("Posts").Append(&posts); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}

// Update 更新用户
func (dao *UserDAO) Update(req *UpdateReq, updateBy uint) error {
	// 开始事务
	tx := dao.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新用户基本信息
	updateData := map[string]interface{}{
		"update_by": updateBy,
	}

	if req.Nickname != "" {
		updateData["nickname"] = req.Nickname
	}
	if req.DeptID != nil {
		updateData["dept_id"] = *req.DeptID
	}
	if req.Email != "" {
		updateData["email"] = req.Email
	}
	if req.Mobile != "" {
		updateData["mobile"] = req.Mobile
	}
	if req.Avatar != "" {
		updateData["avatar"] = req.Avatar
	}
	if req.Status != nil {
		updateData["status"] = *req.Status
	}
	if req.Remark != "" {
		updateData["remark"] = req.Remark
	}

	if err := tx.Model(&system.User{}).Where("id = ?", req.ID).Updates(updateData).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新岗位关联
	if req.PostIDs != nil {
		var user system.User
		if err := tx.First(&user, req.ID).Error; err != nil {
			tx.Rollback()
			return err
		}

		// 清除现有关联
		if err := tx.Model(&user).Association("Posts").Clear(); err != nil {
			tx.Rollback()
			return err
		}

		// 添加新关联
		if len(req.PostIDs) > 0 {
			var posts []system.Post
			if err := tx.Find(&posts, req.PostIDs).Error; err != nil {
				tx.Rollback()
				return err
			}

			if err := tx.Model(&user).Association("Posts").Append(&posts); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// 提交事务
	return tx.Commit().Error
}

// Delete 删除用户
func (dao *UserDAO) Delete(id uint) error {
	return dao.db.Delete(&system.User{}, id).Error
}

// DeleteBatch 批量删除用户
func (dao *UserDAO) DeleteBatch(ids []uint) error {
	return dao.db.Delete(&system.User{}, ids).Error
}

// UpdatePassword 更新用户密码
func (dao *UserDAO) UpdatePassword(req *UpdatePasswordReq) error {
	return dao.db.Model(&system.User{}).
		Where("id = ?", req.ID).
		Update("password", req.Password).Error
}

// UpdateStatus 更新用户状态
func (dao *UserDAO) UpdateStatus(req *UpdateStatusReq, updateBy uint) error {
	return dao.db.Model(&system.User{}).
		Where("id = ?", req.ID).
		Updates(map[string]interface{}{
			"status":    req.Status,
			"update_by": updateBy,
		}).Error
}

// UpdateLoginInfo 更新登录信息
func (dao *UserDAO) UpdateLoginInfo(userID uint, loginIP string) error {
	now := time.Now()
	return dao.db.Model(&system.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"login_ip":   loginIP,
			"login_date": &now,
		}).Error
}

// GetSimpleList 获取用户简单列表
func (dao *UserDAO) GetSimpleList(deptID *uint) ([]UserSimpleResp, error) {
	query := dao.db.Model(&system.User{}).
		Preload("Dept").
		Where("status = ?", 1) // 只查询启用用户

	if deptID != nil {
		query = query.Where("dept_id = ?", *deptID)
	}

	var users []system.User
	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	resp := make([]UserSimpleResp, len(users))
	for i, user := range users {
		resp[i] = UserSimpleResp{
			ID:       user.ID,
			Nickname: user.Nickname,
		}

		if user.Dept != nil {
			resp[i].DeptName = user.Dept.Name
		}
	}

	return resp, nil
}

// CheckUsernameExists 检查用户名是否存在
func (dao *UserDAO) CheckUsernameExists(username string, excludeID *uint) (bool, error) {
	query := dao.db.Model(&system.User{}).Where("username = ?", username)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// CheckEmailExists 检查邮箱是否存在
func (dao *UserDAO) CheckEmailExists(email string, excludeID *uint) (bool, error) {
	if email == "" {
		return false, nil
	}

	query := dao.db.Model(&system.User{}).Where("email = ?", email)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// CheckMobileExists 检查手机号是否存在
func (dao *UserDAO) CheckMobileExists(mobile string, excludeID *uint) (bool, error) {
	if mobile == "" {
		return false, nil
	}

	query := dao.db.Model(&system.User{}).Where("mobile = ?", mobile)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
