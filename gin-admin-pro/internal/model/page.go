package model

import "strings"

// PageReq 分页请求基础结构
type PageReq struct {
	PageNo    int    `form:"pageNo" json:"pageNo" binding:"min=1"`
	PageSize  int    `form:"pageSize" json:"pageSize" binding:"min=1,max=200"`
	SortField string `form:"sortField" json:"sortField"`
	SortOrder string `form:"sortOrder" json:"sortOrder" binding:"omitempty,oneof=asc desc"`
}

// GetSort 获取排序信息
func (p *PageReq) GetSort() string {
	if p.SortField == "" {
		return "id desc"
	}

	sortOrder := "asc"
	if p.SortOrder == "desc" {
		sortOrder = "desc"
	}

	// 转换为下划线命名
	field := strings.ReplaceAll(p.SortField, ".", "_")
	return field + " " + sortOrder
}

// GetOffset 获取偏移量
func (p *PageReq) GetOffset() int {
	if p.PageNo <= 0 {
		p.PageNo = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	return (p.PageNo - 1) * p.PageSize
}

// PageResp 分页响应基础结构
type PageResp struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

// DeleteBatchReq 批量删除请求
type DeleteBatchReq struct {
	IDs []uint `json:"ids" binding:"required"`
}

// RefreshTokenReq 刷新token请求
type RefreshTokenReq struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
