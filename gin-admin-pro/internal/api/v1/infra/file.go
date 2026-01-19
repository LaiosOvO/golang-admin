package infra

import (
	"gin-admin-pro/internal/pkg/response"
	fileservice "gin-admin-pro/internal/service/infra"
	"gin-admin-pro/plugin/oss"

	"github.com/gin-gonic/gin"
)

// FileController 文件控制器
type FileController struct {
	fileService *fileservice.FileService
}

// NewFileController 创建文件控制器实例
func NewFileController(storage oss.OSSInterface) *FileController {
	return &FileController{
		fileService: fileservice.NewFileService(storage),
	}
}

// Upload 上传单个文件
// @Summary 上传文件
// @Description 上传单个文件到对象存储
// @Tags 文件管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Success 200 {object} response.Response{data=UploadResult}
// @Failure 400 {object} response.Response
// @Router /api/v1/infra/file/upload [post]
func (ctrl *FileController) Upload(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "获取上传文件失败："+err.Error())
		return
	}

	// 验证文件
	if err := ctrl.fileService.ValidateFile(file); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 上传文件
	result, err := ctrl.fileService.UploadFile(file)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, result)
}

// UploadMultiple 上传多个文件
// @Summary 上传多个文件
// @Description 上传多个文件到对象存储
// @Tags 文件管理
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "文件（可多选）"
// @Success 200 {object} response.Response{data=[]UploadResult}
// @Failure 400 {object} response.Response
// @Router /api/v1/infra/file/upload-multiple [post]
func (ctrl *FileController) UploadMultiple(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		response.BadRequest(c, "获取上传文件失败："+err.Error())
		return
	}

	// 获取所有上传的文件
	files := form.File["files"]
	if len(files) == 0 {
		response.BadRequest(c, "请选择要上传的文件")
		return
	}

	// 上传文件
	results, err := ctrl.fileService.UploadMultipleFiles(files)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, results)
}

// Delete 删除文件
// @Summary 删除文件
// @Description 删除指定URL的文件
// @Tags 文件管理
// @Accept json
// @Produce json
// @Param request body DeleteFileRequest true "删除文件请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/infra/file/delete [delete]
func (ctrl *FileController) Delete(c *gin.Context) {
	var req DeleteFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	if req.URL == "" {
		response.BadRequest(c, "文件URL不能为空")
		return
	}

	// 删除文件
	err := ctrl.fileService.DeleteFile(req.URL)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteFileRequest 删除文件请求
type DeleteFileRequest struct {
	URL string `json:"url" binding:"required"`
}
