package system

import (
	"gin-admin-pro/internal/dao/system"
	"gin-admin-pro/internal/model"
	"gin-admin-pro/internal/pkg/response"
	"gin-admin-pro/internal/pkg/token"
	userservice "gin-admin-pro/internal/service/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService *userservice.UserService
}

// NewUserController 创建用户控制器实例
func NewUserController(userDAO *system.UserDAO, tokenSvc *token.TokenService) *UserController {
	return &UserController{
		userService: userservice.NewUserService(userDAO, tokenSvc),
	}
}

// Page 获取用户分页列表
// @Summary 获取用户分页列表
// @Description 分页查询用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param pageNo query int true "页码"
// @Param pageSize query int true "每页数量"
// @Param username query string false "用户名"
// @Param mobile query string false "手机号"
// @Param email query string false "邮箱"
// @Param status query int false "状态：0-禁用 1-启用"
// @Param deptId query int false "部门ID"
// @Param createTime query []string false "创建时间范围"
// @Success 200 {object} response.Response{data=model.PageResp}
// @Failure 400 {object} response.Response
// @Router /api/v1/system/user/page [get]
func (ctrl *UserController) Page(c *gin.Context) {
	var req system.UserPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	page, err := ctrl.userService.GetPage(&req)
	if err != nil {
		response.Error(c, "查询用户列表失败")
		return
	}

	response.Success(c, page)
}

// Get 获取用户详情
// @Summary 获取用户详情
// @Description 根据ID获取用户详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id query uint true "用户ID"
// @Success 200 {object} response.Response{data=system.UserDetailResp}
// @Failure 400 {object} response.Response
// @Router /api/v1/system/user/get [get]
func (ctrl *UserController) Get(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		response.BadRequest(c, "用户ID不能为空")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "用户ID格式错误")
		return
	}

	user, err := ctrl.userService.GetByID(uint(id))
	if err != nil {
		if err == userservice.ErrUserNotFound {
			response.NotFound(c, "用户不存在")
			return
		}
		response.Error(c, "获取用户详情失败")
		return
	}

	response.Success(c, user)
}

// Create 创建用户
// @Summary 创建用户
// @Description 创建新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body system.CreateReq true "创建用户请求"
// @Success 200 {object} response.Response{data=uint}
// @Failure 400 {object} response.Response
// @Router /api/v1/system/user/create [post]
func (ctrl *UserController) Create(c *gin.Context) {
	var req system.CreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	// 获取当前操作用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "未获取到用户信息")
		return
	}

	id, err := ctrl.userService.Create(&req, userID.(uint))
	if err != nil {
		switch err {
		case userservice.ErrUsernameExists:
			response.BadRequest(c, "用户名已存在")
			return
		case userservice.ErrEmailExists:
			response.BadRequest(c, "邮箱已存在")
			return
		case userservice.ErrMobileExists:
			response.BadRequest(c, "手机号已存在")
			return
		default:
			response.Error(c, "创建用户失败")
			return
		}
	}

	response.Success(c, id)
}

// Update 更新用户
// @Summary 更新用户
// @Description 更新用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body system.UpdateReq true "更新用户请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/system/user/update [put]
func (ctrl *UserController) Update(c *gin.Context) {
	var req system.UpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	// 获取当前操作用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "未获取到用户信息")
		return
	}

	err := ctrl.userService.Update(&req, userID.(uint))
	if err != nil {
		switch err {
		case userservice.ErrUserNotFound:
			response.NotFound(c, "用户不存在")
			return
		case userservice.ErrEmailExists:
			response.BadRequest(c, "邮箱已存在")
			return
		case userservice.ErrMobileExists:
			response.BadRequest(c, "手机号已存在")
			return
		default:
			response.Error(c, "更新用户失败")
			return
		}
	}

	response.Success(c, nil)
}

// Delete 删除用户
// @Summary 删除用户
// @Description 根据ID删除用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id query uint true "用户ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/system/user/delete [delete]
func (ctrl *UserController) Delete(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		response.BadRequest(c, "用户ID不能为空")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "用户ID格式错误")
		return
	}

	err = ctrl.userService.Delete(uint(id))
	if err != nil {
		if err == userservice.ErrUserNotFound {
			response.NotFound(c, "用户不存在")
			return
		}
		response.Error(c, "删除用户失败")
		return
	}

	response.Success(c, nil)
}

// DeleteBatch 批量删除用户
// @Summary 批量删除用户
// @Description 批量删除用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body model.DeleteBatchReq true "批量删除请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/system/user/delete-list [delete]
func (ctrl *UserController) DeleteBatch(c *gin.Context) {
	var req model.DeleteBatchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	if len(req.IDs) == 0 {
		response.BadRequest(c, "请选择要删除的用户")
		return
	}

	err := ctrl.userService.DeleteBatch(req.IDs)
	if err != nil {
		if err == userservice.ErrUserNotFound {
			response.NotFound(c, "用户不存在")
			return
		}
		response.Error(c, "批量删除用户失败")
		return
	}

	response.Success(c, nil)
}

// UpdatePassword 重置用户密码
// @Summary 重置用户密码
// @Description 重置用户密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body system.UpdatePasswordReq true "重置密码请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/system/user/update-password [put]
func (ctrl *UserController) UpdatePassword(c *gin.Context) {
	var req system.UpdatePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	err := ctrl.userService.UpdatePassword(&req, 0) // 管理员重置密码不需要记录操作人
	if err != nil {
		if err == userservice.ErrUserNotFound {
			response.NotFound(c, "用户不存在")
			return
		}
		response.Error(c, "重置密码失败")
		return
	}

	response.Success(c, nil)
}

// UpdateStatus 修改用户状态
// @Summary 修改用户状态
// @Description 启用或禁用用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body system.UpdateStatusReq true "修改状态请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/system/user/update-status [put]
func (ctrl *UserController) UpdateStatus(c *gin.Context) {
	var req system.UpdateStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	// 获取当前操作用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "未获取到用户信息")
		return
	}

	err := ctrl.userService.UpdateStatus(&req, userID.(uint))
	if err != nil {
		if err == userservice.ErrUserNotFound {
			response.NotFound(c, "用户不存在")
			return
		}
		response.Error(c, "修改用户状态失败")
		return
	}

	response.Success(c, nil)
}

// SimpleList 获取用户简单列表
// @Summary 获取用户简单列表
// @Description 获取用户简单信息列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param deptId query uint false "部门ID"
// @Success 200 {object} response.Response{data=[]system.UserSimpleResp}
// @Failure 400 {object} response.Response
// @Router /api/v1/system/user/simple-list [get]
func (ctrl *UserController) SimpleList(c *gin.Context) {
	deptIDStr := c.Query("deptId")
	var deptID *uint

	if deptIDStr != "" {
		id, err := strconv.ParseUint(deptIDStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "部门ID格式错误")
			return
		}
		deptIDUint := uint(id)
		deptID = &deptIDUint
	}

	users, err := ctrl.userService.GetSimpleList(deptID)
	if err != nil {
		response.Error(c, "获取用户列表失败")
		return
	}

	response.Success(c, users)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body system.LoginReq true "登录请求"
// @Success 200 {object} response.Response{data=system.LoginResp}
// @Failure 400 {object} response.Response
// @Router /api/v1/system/auth/login [post]
func (ctrl *UserController) Login(c *gin.Context) {
	var req userservice.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	// 获取客户端IP
	clientIP := c.ClientIP()

	loginResp, err := ctrl.userService.Login(&req, clientIP)
	if err != nil {
		switch err {
		case userservice.ErrInvalidCredentials:
			response.BadRequest(c, "用户名或密码错误")
			return
		case userservice.ErrUserDisabled:
			response.Forbidden(c, "用户已被禁用")
			return
		default:
			response.Error(c, "登录失败")
			return
		}
	}

	response.Success(c, loginResp)
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出接口
// @Tags 认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/system/auth/logout [post]
func (ctrl *UserController) Logout(c *gin.Context) {
	// 从请求头获取token
	token := c.GetHeader("Authorization")
	if token == "" {
		response.BadRequest(c, "token不能为空")
		return
	}

	// 移除Bearer前缀
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	err := ctrl.userService.Logout(token)
	if err != nil {
		response.Error(c, "登出失败")
		return
	}

	response.Success(c, nil)
}

// RefreshToken 刷新token
// @Summary 刷新token
// @Description 刷新访问token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body model.RefreshTokenReq true "刷新token请求"
// @Success 200 {object} response.Response{data=system.LoginResp}
// @Failure 400 {object} response.Response
// @Router /api/v1/system/auth/refresh-token [post]
func (ctrl *UserController) RefreshToken(c *gin.Context) {
	var req model.RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	loginResp, err := ctrl.userService.RefreshToken(req.RefreshToken)
	if err != nil {
		response.BadRequest(c, "刷新token失败")
		return
	}

	response.Success(c, loginResp)
}

// GetPermissionInfo 获取用户权限信息
// @Summary 获取用户权限信息
// @Description 获取当前用户的详细信息、角色、权限和菜单
// @Tags 认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=system.UserInfoResp}
// @Failure 400 {object} response.Response
// @Router /api/v1/system/auth/get-permission-info [get]
func (ctrl *UserController) GetPermissionInfo(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "未获取到用户信息")
		return
	}

	userInfo, err := ctrl.userService.GetUserInfo(userID.(uint))
	if err != nil {
		response.Error(c, "获取用户信息失败")
		return
	}

	response.Success(c, userInfo)
}
