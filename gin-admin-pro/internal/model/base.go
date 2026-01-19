package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型，包含所有实体的公共字段
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// AuditModel 审计模型，包含基础模型和审计字段
type AuditModel struct {
	BaseModel
	CreateBy uint   `json:"createBy"`
	UpdateBy uint   `json:"updateBy"`
	Remark   string `gorm:"size:500" json:"remark"`
}

// TreeModel 树形结构模型
type TreeModel struct {
	AuditModel
	ParentID  uint   `json:"parentId"`
	Level     int    `json:"level"`
	Sort      int    `json:"sort"`
	Name      string `gorm:"size:100;not null" json:"name"`
	Path      string `gorm:"size:500" json:"path"`
	Ancestors string `gorm:"size:500" json:"ancestors"`
}

// TableName 自定义表名前缀
const (
	TablePrefixSystem = "system_"
	TablePrefixInfra  = "infra_"
	TablePrefixBpm    = "bpm_"
	TablePrefixPay    = "pay_"
	TablePrefixMember = "member_"
	TablePrefixMall   = "mall_"
	TablePrefixReport = "report_"
)

// GetTableName 获取带前缀的表名
func GetTableName(prefix, name string) string {
	return prefix + name
}

// BeforeCreate GORM钩子 - 创建前
func (m *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	if m.UpdatedAt.IsZero() {
		m.UpdatedAt = time.Now()
	}
	return nil
}

// BeforeUpdate GORM钩子 - 更新前
func (m *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}

// IsDeleted 检查是否被软删除
func (m *BaseModel) IsDeleted() bool {
	return m.DeletedAt.Valid
}

// GetID 获取ID
func (m *BaseModel) GetID() uint {
	return m.ID
}

// SetCreatedBy 设置创建人
func (m *AuditModel) SetCreatedBy(userID uint) {
	m.CreateBy = userID
}

// SetUpdatedBy 设置更新人
func (m *AuditModel) SetUpdatedBy(userID uint) {
	m.UpdateBy = userID
}

// GetCreatedBy 获取创建人ID
func (m *AuditModel) GetCreatedBy() uint {
	return m.CreateBy
}

// GetUpdatedBy 获取更新人ID
func (m *AuditModel) GetUpdatedBy() uint {
	return m.UpdateBy
}

// GetRemark 获取备注
func (m *AuditModel) GetRemark() string {
	return m.Remark
}

// SetRemark 设置备注
func (m *AuditModel) SetRemark(remark string) {
	m.Remark = remark
}
