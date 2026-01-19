package system

import (
	"gin-admin-pro/internal/dao/system"
	"gin-admin-pro/internal/model"
	"gin-admin-pro/internal/pkg/response"
	roleservice "gin-admin-pro/internal/service/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RoleController 角色控制器
type RoleController struct {
	roleService *roleservice.RoleService
}

// NewRoleController 创建角色控制器实例
func NewRoleController(roleDAO *system.RoleDAO) *RoleController {
	return &RoleController{
		roleService: roleservice.NewRoleService(roleDAO),
	}
}

// Page 获取角色分页列表
// @Summary 获取角色分页列表
// @Description 分页查询角色列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param pageNo query int true "页码"
// @Param pageSize query int true "每页数量"
// @Param name query string false "角色名称"
// @Param code query string false "角色代码"
// @Param status query int false "状态：0-禁用 1-启用"
// @Param createTime query []string false "创建时间范围"
// @Success 200 {object} response.Response{data=model.PageResp}
// @Failure 400 {object} response.Response
// @Router /api/v1/system/role/page [get]
func (ctrl *RoleController) Page(c *gin.Context) {
	var req system.RolePageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	page, err := ctrl.roleService.GetPage(&req)
	if err != nil {
		response.Error(c, "查询失败："+err.Error())
		return
	}

	response.Success(c, page)
}

// Get 获取角色详情
// @Summary 获取角色详情
// @Description 根据ID获取角色详情信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id query int true "角色ID"
// @Success 200 {object} response.Response{data=system.RoleDetailResp}
// @Failure 400 {object} response.Response
// @Router /api/v1/system/role/get [get]
func (ctrl *RoleController) Get(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		response.BadRequest(c, "角色ID不能为空")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "角色ID格式错误")
		return
	}

	role, err := ctrl.roleService.GetByID(uint(id))
	if err != nil {
		if err == roleservice.ErrRoleNotFound {
			response.NotFound(c, "角色不存在")
			return
		}
		response.Error(c, "查询失败："+err.Error())
		return
	}

	response.Success(c, role)
}

// ListAllSimple 获取角色精简列表
// @Summary 获取角色精简列表
// @Description 获取所有启用的角色精简列表，用于下拉选择
// @Tags 角色管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]system.RoleSimpleResp}
// @Failure 400 {object} response.Response
// @Router /api/v1/system/role/list-all-simple [get]
func (ctrl *RoleController) ListAllSimple(c *gin.Context) {
	roles, err := ctrl.roleService.GetAllSimple()
	if err != nil {
		response.Error(c, "查询失败："+err.Error())
		return
	}

	response.Success(c, roles)
}

// Create 创建角色
// @Summary 创建角色
// @Description 创建新的角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param request body system.RoleCreateReq true "角色信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/system/role/create [post]
func (ctrl *RoleController) Create(c *gin.Context) {
	var req system.RoleCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	// 设置默认值
	if req.Status == 0 {
		req.Status = 1 // 默认启用
	}
	if req.DataScope == 0 {
		req.DataScope = 1 // 默认全部数据权限
	}
	if req.Type == 0 {
		req.Type = 2 // 默认自定义角色
	}

	// 获取当前操作用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "未获取到用户信息")
		return
	}
	createBy := userID.(uint)

	err := ctrl.roleService.Create(&req, createBy)
	if err != nil {
		if err == roleservice.ErrRoleCodeExists {
			response.BadRequest(c, "角色代码已存在")
			return
		}
		if err == roleservice.ErrInvalidDataScope {
			response.BadRequest(c, "无效的数据权限范围")
			return
		}
		response.Error(c, "创建失败："+err.Error())
		return
	}

	response.Success(c, nil)
}

// Update 更新角色
// @Summary 更新角色
// @Description 更新角色信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param request body system.RoleUpdateReq true "角色信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/system/role/update [put]
func (ctrl *RoleController) Update(c *gin.Context) {
	var req system.RoleUpdateReq
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
	updateBy := userID.(uint)

	err := ctrl.roleService.Update(&req, updateBy)
	if err != nil {
		if err == roleservice.ErrRoleNotFound {
			response.NotFound(c, "角色不存在")
			return
		}
		if err == roleservice.ErrRoleCodeExists {
			response.BadRequest(c, "角色代码已存在")
			return
		}
		if err == roleservice.ErrInvalidDataScope {
			response.BadRequest(c, "无效的数据权限范围")
			return
		}
		response.Error(c, "更新失败："+err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete 删除角色
// @Summary 删除角色
// @Description 根据ID删除角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id query int true "角色ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/system/role/delete [delete]
func (ctrl *RoleController) Delete(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		response.BadRequest(c, "角色ID不能为空")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "角色ID格式错误")
		return
	}

	err = ctrl.roleService.Delete(uint(id))
	if err != nil {
		if err == roleservice.ErrRoleNotFound {
			response.NotFound(c, "角色不存在")
			return
		}
		if err == roleservice.ErrRoleIsBuiltin {
			response.BadRequest(c, "内置角色，无法删除")
			return
		}
		if err == roleservice.ErrRoleHasUsers {
			response.BadRequest(c, "角色下存在用户，无法删除")
			return
		}
		response.Error(c, "删除失败："+err.Error())
		return
	}

	response.Success(c, nil)
}

// AssignMenuPermissions 分配菜单权限
// @Summary 分配菜单权限
// @Description 为角色分配菜单权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param roleId path int true "角色ID"
// @Param request body model.DeleteBatchReq true "菜单ID列表"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/system/role/assign-menu/{roleId} [put]
func (ctrl *RoleController) AssignMenuPermissions(c *gin.Context) {
	roleIdStr := c.Param("roleId")
	if roleIdStr == "" {
		response.BadRequest(c, "角色ID不能为空")
		return
	}

	roleId, err := strconv.ParseUint(roleIdStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "角色ID格式错误")
		return
	}

	var req model.DeleteBatchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	err = ctrl.roleService.AssignMenuPermissions(uint(roleId), req.IDs)
	if err != nil {
		if err == roleservice.ErrRoleNotFound {
			response.NotFound(c, "角色不存在")
			return
		}
		response.Error(c, "分配失败："+err.Error())
		return
	}

	response.Success(c, nil)
}

// GetMenuPermissions 获取菜单权限
// @Summary 获取菜单权限
// @Description 获取角色的菜单权限列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param roleId path int true "角色ID"
// @Success 200 {object} response.Response{data=[]uint}
// @Failure 400 {object} response.Response
// @Router /api/v1/system/role/menu-permissions/{roleId} [get]
func (ctrl *RoleController) GetMenuPermissions(c *gin.Context) {
	roleIdStr := c.Param("roleId")
	if roleIdStr == "" {
		response.BadRequest(c, "角色ID不能为空")
		return
	}

	roleId, err := strconv.ParseUint(roleIdStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "角色ID格式错误")
		return
	}

	menuIDs, err := ctrl.roleService.GetMenuPermissions(uint(roleId))
	if err != nil {
		if err == roleservice.ErrRoleNotFound {
			response.NotFound(c, "角色不存在")
			return
		}
		response.Error(c, "查询失败："+err.Error())
		return
	}

	response.Success(c, menuIDs)
}
