package system

import (
	"errors"
	"gin-admin-pro/internal/dao/system"
	"gin-admin-pro/internal/model"
	sysmodel "gin-admin-pro/internal/model/system"
	"gin-admin-pro/internal/pkg/token"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound       = errors.New("用户不存在")
	ErrUsernameExists     = errors.New("用户名已存在")
	ErrEmailExists        = errors.New("邮箱已存在")
	ErrMobileExists       = errors.New("手机号已存在")
	ErrPasswordIncorrect  = errors.New("密码错误")
	ErrUserDisabled       = errors.New("用户已被禁用")
	ErrInvalidCredentials = errors.New("用户名或密码错误")
)

// UserService 用户服务层
type UserService struct {
	userDAO  *system.UserDAO
	tokenSvc *token.TokenService
}

// NewUserService 创建用户服务实例
func NewUserService(userDAO *system.UserDAO, tokenSvc *token.TokenService) *UserService {
	return &UserService{
		userDAO:  userDAO,
		tokenSvc: tokenSvc,
	}
}

// LoginReq 登录请求
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResp 登录响应
type LoginResp struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresTime  int64  `json:"expiresTime"`
}

// UserInfoResp 用户信息响应
type UserInfoResp struct {
	User        *UserInfo      `json:"user"`
	Roles       []string       `json:"roles"`
	Permissions []string       `json:"permissions"`
	Menus       []MenuTreeNode `json:"menus"`
}

// UserInfo 用户基本信息
type UserInfo struct {
	ID         uint   `json:"id"`
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	Email      string `json:"email"`
	Mobile     string `json:"mobile"`
	Avatar     string `json:"avatar"`
	Status     int    `json:"status"`
	CreateTime string `json:"createTime"`
}

// MenuTreeNode 菜单树节点
type MenuTreeNode struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	Path      string         `json:"path"`
	Component string         `json:"component"`
	Icon      string         `json:"icon"`
	Sort      int            `json:"sort"`
	ParentID  uint           `json:"parentId"`
	Type      int            `json:"type"`
	Status    int            `json:"status"`
	Visible   bool           `json:"visible"`
	Children  []MenuTreeNode `json:"children,omitempty"`
}

// GetPage 获取用户分页列表
func (s *UserService) GetPage(req *system.UserPageReq) (*model.PageResp, error) {
	users, total, err := s.userDAO.GetPage(req)
	if err != nil {
		return nil, err
	}

	return &model.PageResp{
		List:  users,
		Total: total,
	}, nil
}

// GetByID 根据ID获取用户详情
func (s *UserService) GetByID(id uint) (*system.UserDetailResp, error) {
	return s.userDAO.GetByID(id)
}

// Create 创建用户
func (s *UserService) Create(req *system.CreateReq, operatorID uint) (uint, error) {
	// 验证用户名唯一性
	exists, err := s.userDAO.CheckUsernameExists(req.Username, nil)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, ErrUsernameExists
	}

	// 验证邮箱唯一性
	if req.Email != "" {
		exists, err = s.userDAO.CheckEmailExists(req.Email, nil)
		if err != nil {
			return 0, err
		}
		if exists {
			return 0, ErrEmailExists
		}
	}

	// 验证手机号唯一性
	if req.Mobile != "" {
		exists, err = s.userDAO.CheckMobileExists(req.Mobile, nil)
		if err != nil {
			return 0, err
		}
		if exists {
			return 0, ErrMobileExists
		}
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	createReq := &system.CreateReq{
		Username: req.Username,
		Nickname: req.Nickname,
		Password: string(hashedPassword),
		DeptID:   req.DeptID,
		PostIDs:  req.PostIDs,
		Email:    req.Email,
		Mobile:   req.Mobile,
		Avatar:   req.Avatar,
		Status:   req.Status,
		Remark:   req.Remark,
	}

	return s.userDAO.Create(createReq, operatorID)
}

// Update 更新用户
func (s *UserService) Update(req *system.UpdateReq, operatorID uint) error {
	// 检查用户是否存在
	user, err := s.userDAO.GetByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// 验证邮箱唯一性
	if req.Email != "" && req.Email != user.Email {
		exists, err := s.userDAO.CheckEmailExists(req.Email, &req.ID)
		if err != nil {
			return err
		}
		if exists {
			return ErrEmailExists
		}
	}

	// 验证手机号唯一性
	if req.Mobile != "" && req.Mobile != user.Mobile {
		exists, err := s.userDAO.CheckMobileExists(req.Mobile, &req.ID)
		if err != nil {
			return err
		}
		if exists {
			return ErrMobileExists
		}
	}

	return s.userDAO.Update(req, operatorID)
}

// Delete 删除用户
func (s *UserService) Delete(id uint) error {
	// 检查用户是否存在
	_, err := s.userDAO.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	return s.userDAO.Delete(id)
}

// DeleteBatch 批量删除用户
func (s *UserService) DeleteBatch(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	// 验证所有用户是否存在
	for _, id := range ids {
		_, err := s.userDAO.GetByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrUserNotFound
			}
			return err
		}
	}

	return s.userDAO.DeleteBatch(ids)
}

// UpdatePassword 更新用户密码
func (s *UserService) UpdatePassword(req *system.UpdatePasswordReq, operatorID uint) error {
	// 检查用户是否存在
	_, err := s.userDAO.GetByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	updateReq := &system.UpdatePasswordReq{
		ID:       req.ID,
		Password: string(hashedPassword),
	}

	return s.userDAO.UpdatePassword(updateReq)
}

// UpdateStatus 更新用户状态
func (s *UserService) UpdateStatus(req *system.UpdateStatusReq, operatorID uint) error {
	// 检查用户是否存在
	_, err := s.userDAO.GetByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	return s.userDAO.UpdateStatus(req, operatorID)
}

// Login 用户登录
func (s *UserService) Login(req *LoginReq, clientIP string) (*LoginResp, error) {
	// 获取用户信息
	user, err := s.userDAO.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, ErrUserDisabled
	}

	// 生成JWT Token
	tokenPair, err := s.tokenSvc.GenerateTokens(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	// 更新登录信息
	if err := s.userDAO.UpdateLoginInfo(user.ID, clientIP); err != nil {
		// 记录日志但不影响登录流程
		// TODO: 添加日志记录
	}

	return &LoginResp{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresTime:  tokenPair.ExpiresIn,
	}, nil
}

// Logout 用户登出
func (s *UserService) Logout(accessToken string) error {
	// 将token加入黑名单
	return s.tokenSvc.RevokeToken(accessToken)
}

// RefreshToken 刷新token
func (s *UserService) RefreshToken(refreshToken string) (*LoginResp, error) {
	tokenPair, err := s.tokenSvc.RefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	return &LoginResp{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresTime:  tokenPair.ExpiresIn,
	}, nil
}

// GetUserInfo 获取当前用户信息
func (s *UserService) GetUserInfo(userID uint) (*UserInfoResp, error) {
	// 获取用户信息
	user, err := s.userDAO.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// 获取用户角色
	userModel, err := s.userDAO.GetByUsername(user.Username)
	if err != nil {
		return nil, err
	}

	roles := make([]string, len(userModel.Roles))
	permissions := make([]string, 0)

	for i, role := range userModel.Roles {
		roles[i] = role.Code
		// 收集权限代码
		for _, menu := range role.Menus {
			if menu.Perms != "" {
				permissions = append(permissions, menu.Perms)
			}
		}
	}

	// 构建菜单树
	menus := s.buildMenuTree(userModel.Roles)

	userInfo := &UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Email:    user.Email,
		Mobile:   user.Mobile,
		Avatar:   user.Avatar,
		Status:   user.Status,
	}

	return &UserInfoResp{
		User:        userInfo,
		Roles:       roles,
		Permissions: permissions,
		Menus:       menus,
	}, nil
}

// GetSimpleList 获取用户简单列表
func (s *UserService) GetSimpleList(deptID *uint) ([]system.UserSimpleResp, error) {
	return s.userDAO.GetSimpleList(deptID)
}

// buildMenuTree 构建菜单树
func (s *UserService) buildMenuTree(roles []sysmodel.Role) []MenuTreeNode {
	// 收集所有菜单
	menuMap := make(map[uint]*MenuTreeNode)
	rootMenus := make([]*MenuTreeNode, 0)

	for _, role := range roles {
		for _, menu := range role.Menus {
			if menu.Status != 1 { // 只处理启用的菜单
				continue
			}

			node := &MenuTreeNode{
				ID:        menu.ID,
				Name:      menu.Name,
				Path:      menu.Path,
				Component: menu.Component,
				Icon:      menu.Icon,
				Sort:      menu.Sort,
				ParentID:  menu.ParentID,
				Type:      menu.Type,
				Status:    menu.Status,
				Visible:   menu.Visible == 1,
			}

			menuMap[menu.ID] = node

			if menu.ParentID == 0 {
				rootMenus = append(rootMenus, node)
			}
		}
	}

	// 构建树形结构
	for _, menu := range menuMap {
		if menu.ParentID != 0 {
			if parent, exists := menuMap[menu.ParentID]; exists {
				parent.Children = append(parent.Children, *menu)
			}
		}
	}

	// 转换为返回格式
	result := make([]MenuTreeNode, len(rootMenus))
	for i, menu := range rootMenus {
		result[i] = *menu
	}

	return result
}
