package system

import (
	"gin-admin-pro/internal/dao/system"
	"gin-admin-pro/internal/pkg/response"
	deptservice "gin-admin-pro/internal/service/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeptController 部门控制器
type DeptController struct {
	deptService *deptservice.DeptService
}

// NewDeptController 创建部门控制器实例
func NewDeptController(deptDAO *system.DeptDAO) *DeptController {
	return &DeptController{
		deptService: deptservice.NewDeptService(deptDAO),
	}
}

// List 获取部门列表
// @Summary 获取部门列表
// @Description 获取部门树形结构列表
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param name query string false "部门名称"
// @Param status query int false "状态：0-禁用 1-启用"
// @Success 200 {object} response.Response{data=[]system.DeptResp}
// @Failure 400 {object} response.Response
// @Router /api/v1/system/dept/list [get]
func (ctrl *DeptController) List(c *gin.Context) {
	var req system.DeptListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	depts, err := ctrl.deptService.GetList(&req)
	if err != nil {
		response.Error(c, "获取部门列表失败")
		return
	}

	response.Success(c, depts)
}

// Get 获取部门详情
// @Summary 获取部门详情
// @Description 根据ID获取部门详情
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param id query int true "部门ID"
// @Success 200 {object} response.Response{data=system.DeptDetailResp}
// @Failure 400 {object} response.Response
// @Router /api/v1/system/dept/get [get]
func (ctrl *DeptController) Get(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "部门ID格式错误")
		return
	}

	dept, err := ctrl.deptService.GetByID(uint(id))
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, dept)
}

// Create 创建部门
// @Summary 创建部门
// @Description 创建新部门
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param request body system.CreateDeptReq true "创建部门请求"
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Failure 400 {object} response.Response
// @Router /api/v1/system/dept/create [post]
func (ctrl *DeptController) Create(c *gin.Context) {
	var req system.CreateDeptReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	// 从上下文获取当前用户ID
	userID := c.GetUint("user_id")

	id, err := ctrl.deptService.Create(&req, userID)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"id": id,
	})
}

// Update 更新部门
// @Summary 更新部门
// @Description 更新部门信息
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param request body system.UpdateDeptReq true "更新部门请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/system/dept/update [put]
func (ctrl *DeptController) Update(c *gin.Context) {
	var req system.UpdateDeptReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	// 从上下文获取当前用户ID
	userID := c.GetUint("user_id")

	err := ctrl.deptService.Update(&req, userID)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete 删除部门
// @Summary 删除部门
// @Description 删除部门（如果存在子部门或用户则无法删除）
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param id query int true "部门ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/system/dept/delete [delete]
func (ctrl *DeptController) Delete(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "部门ID格式错误")
		return
	}

	err = ctrl.deptService.Delete(uint(id))
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// ListAllSimple 获取部门简单列表
// @Summary 获取部门简单列表
// @Description 获取所有启用的部门简单列表（用于用户选择等场景）
// @Tags 部门管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]system.DeptSimpleResp}
// @Failure 400 {object} response.Response
// @Router /api/v1/system/dept/list-all-simple [get]
func (ctrl *DeptController) ListAllSimple(c *gin.Context) {
	depts, err := ctrl.deptService.GetAllSimpleList()
	if err != nil {
		response.Error(c, "获取部门列表失败")
		return
	}

	response.Success(c, depts)
}

// GetUsers 获取部门用户
// @Summary 获取部门用户
// @Description 获取指定部门下的用户列表
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param deptId query int true "部门ID"
// @Success 200 {object} response.Response{data=[]system.User}
// @Failure 400 {object} response.Response
// @Router /api/v1/system/dept/users [get]
func (ctrl *DeptController) GetUsers(c *gin.Context) {
	deptIdStr := c.Query("deptId")
	deptID, err := strconv.ParseUint(deptIdStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "部门ID格式错误")
		return
	}

	users, err := ctrl.deptService.GetUsersByDept(uint(deptID))
	if err != nil {
		response.Error(c, "获取部门用户失败")
		return
	}

	response.Success(c, users)
}
