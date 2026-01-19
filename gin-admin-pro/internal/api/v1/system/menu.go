package system

import (
	"gin-admin-pro/internal/dao/system"
	"gin-admin-pro/internal/pkg/response"
	menuservice "gin-admin-pro/internal/service/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MenuController 菜单控制器
type MenuController struct {
	menuService *menuservice.MenuService
}

// NewMenuController 创建菜单控制器实例
func NewMenuController(menuDAO *system.MenuDAO) *MenuController {
	return &MenuController{
		menuService: menuservice.NewMenuService(menuDAO),
	}
}

// List 获取菜单列表
// @Summary 获取菜单列表
// @Description 获取菜单树形结构列表
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param name query string false "菜单名称"
// @Param status query int false "状态：0-禁用 1-启用"
// @Param type query int false "菜单类型：1-目录 2-菜单 3-按钮"
// @Success 200 {object} response.Response{data=[]system.MenuResp}
// @Failure 400 {object} response.Response
// @Router /api/v1/system/menu/list [get]
func (ctrl *MenuController) List(c *gin.Context) {
	var req system.MenuListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	menus, err := ctrl.menuService.GetList(&req)
	if err != nil {
		response.Error(c, "获取菜单列表失败")
		return
	}

	response.Success(c, menus)
}

// Get 获取菜单详情
// @Summary 获取菜单详情
// @Description 根据ID获取菜单详情
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param id query int true "菜单ID"
// @Success 200 {object} response.Response{data=system.MenuDetailResp}
// @Failure 400 {object} response.Response
// @Router /api/v1/system/menu/get [get]
func (ctrl *MenuController) Get(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "菜单ID格式错误")
		return
	}

	menu, err := ctrl.menuService.GetByID(uint(id))
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, menu)
}

// Create 创建菜单
// @Summary 创建菜单
// @Description 创建新菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param request body system.CreateMenuReq true "创建菜单请求"
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Failure 400 {object} response.Response
// @Router /api/v1/system/menu/create [post]
func (ctrl *MenuController) Create(c *gin.Context) {
	var req system.CreateMenuReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	// 从上下文获取当前用户ID
	userID := c.GetUint("user_id")

	id, err := ctrl.menuService.Create(&req, userID)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"id": id,
	})
}

// Update 更新菜单
// @Summary 更新菜单
// @Description 更新菜单信息
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param request body system.UpdateMenuReq true "更新菜单请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/system/menu/update [put]
func (ctrl *MenuController) Update(c *gin.Context) {
	var req system.UpdateMenuReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	// 从上下文获取当前用户ID
	userID := c.GetUint("user_id")

	err := ctrl.menuService.Update(&req, userID)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete 删除菜单
// @Summary 删除菜单
// @Description 删除菜单（如果存在子菜单则无法删除）
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param id query int true "菜单ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/system/menu/delete [delete]
func (ctrl *MenuController) Delete(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "菜单ID格式错误")
		return
	}

	err = ctrl.menuService.Delete(uint(id))
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// ListAllSimple 获取菜单简单列表
// @Summary 获取菜单简单列表
// @Description 获取所有启用的菜单简单列表（用于角色授权）
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]system.MenuSimpleResp}
// @Failure 400 {object} response.Response
// @Router /api/v1/system/menu/list-all-simple [get]
func (ctrl *MenuController) ListAllSimple(c *gin.Context) {
	menus, err := ctrl.menuService.GetAllSimpleList()
	if err != nil {
		response.Error(c, "获取菜单列表失败")
		return
	}

	response.Success(c, menus)
}

// ListUserPermissions 获取用户菜单权限
// @Summary 获取用户菜单权限
// @Description 获取当前用户的菜单权限列表
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]system.MenuResp}
// @Failure 400 {object} response.Response
// @Router /api/v1/system/permission/list-user-permissions [get]
func (ctrl *MenuController) ListUserPermissions(c *gin.Context) {
	// 从上下文获取当前用户ID
	userID := c.GetUint("user_id")

	menus, err := ctrl.menuService.GetUserMenus(userID)
	if err != nil {
		response.Error(c, "获取用户菜单权限失败")
		return
	}

	response.Success(c, menus)
}
