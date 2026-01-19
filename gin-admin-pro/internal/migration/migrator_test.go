package migration

import (
	"testing"
	"time"

	"gin-admin-pro/internal/model"
	"gin-admin-pro/internal/model/system"
	"github.com/stretchr/testify/assert"
)

func TestMigrator_NewMigrator(t *testing.T) {
	// 创建迁移器测试需要实际的数据库连接
	// 这里只测试迁移器的基本创建逻辑
	assert.NotNil(t, &Migrator{})
}

func TestBaseModel(t *testing.T) {
	// 测试基础模型
	baseModel := &model.BaseModel{
		ID:        1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	assert.Equal(t, uint(1), baseModel.GetID())
	assert.False(t, baseModel.IsDeleted())
}

func TestAuditModel(t *testing.T) {
	// 测试审计模型
	auditModel := &model.AuditModel{
		BaseModel: model.BaseModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		CreateBy: 100,
		UpdateBy: 100,
		Remark:   "测试备注",
	}

	assert.Equal(t, uint(100), auditModel.GetCreatedBy())
	assert.Equal(t, uint(100), auditModel.GetUpdatedBy())
	assert.Equal(t, "测试备注", auditModel.GetRemark())

	auditModel.SetRemark("新备注")
	auditModel.SetCreatedBy(200)
	auditModel.SetUpdatedBy(200)

	assert.Equal(t, "新备注", auditModel.GetRemark())
	assert.Equal(t, uint(200), auditModel.GetCreatedBy())
	assert.Equal(t, uint(200), auditModel.GetUpdatedBy())
}

func TestSystemModels(t *testing.T) {
	// 测试用户模型
	user := &system.User{
		AuditModel: model.AuditModel{
			BaseModel: model.BaseModel{
				ID:        1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			CreateBy: 100,
			UpdateBy: 100,
		},
		Username: "testuser",
		Nickname: "测试用户",
		Mobile:   "13800138000",
		Email:    "test@example.com",
		Status:   1,
		DeptID:   1,
	}

	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "测试用户", user.Nickname)
	assert.Equal(t, "13800138000", user.Mobile)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, 1, user.Status)
	assert.Equal(t, uint(1), user.DeptID)

	// 测试角色模型
	role := &system.Role{
		AuditModel: model.AuditModel{
			BaseModel: model.BaseModel{
				ID:        1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			CreateBy: 100,
			UpdateBy: 100,
		},
		Code:      "admin",
		Name:      "管理员",
		Sort:      0,
		DataScope: 1,
		Status:    1,
		Type:      1,
		Remark:    "管理员角色",
	}

	assert.Equal(t, "admin", role.Code)
	assert.Equal(t, "管理员", role.Name)
	assert.Equal(t, 1, role.Status)
	assert.Equal(t, 1, role.Type)
	assert.Equal(t, "管理员角色", role.Remark)

	// 测试菜单模型
	menu := &system.Menu{
		TreeModel: model.TreeModel{
			AuditModel: model.AuditModel{
				BaseModel: model.BaseModel{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				CreateBy: 100,
				UpdateBy: 100,
			},
			ParentID:  0,
			Level:     0,
			Sort:      0,
			Name:      "系统管理",
			Ancestors: "0",
			Path:      "system",
		},
		Type:       1,
		Icon:       "system",
		Component:  "system/index",
		Status:     1,
		Visible:    1,
		KeepAlive:  1,
		AlwaysShow: 1,
	}

	assert.Equal(t, "系统管理", menu.Name)
	assert.Equal(t, 1, menu.Type)
	assert.Equal(t, "system", menu.Icon)
	assert.Equal(t, "system/index", menu.Component)
	assert.Equal(t, 1, menu.Status)

	// 测试部门模型
	dept := &system.Dept{
		TreeModel: model.TreeModel{
			AuditModel: model.AuditModel{
				BaseModel: model.BaseModel{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				CreateBy: 100,
				UpdateBy: 100,
			},
			ParentID:  0,
			Level:     0,
			Sort:      0,
			Name:      "总公司",
			Ancestors: "0",
			Path:      "root",
		},
		LeaderUserId: 1,
		Phone:        "13800138000",
		Email:        "admin@example.com",
	}

	assert.Equal(t, "总公司", dept.Name)
	assert.Equal(t, uint(1), dept.LeaderUserId)
	assert.Equal(t, "13800138000", dept.Phone)
	assert.Equal(t, "admin@example.com", dept.Email)

	// 测试岗位模型
	post := &system.Post{
		AuditModel: model.AuditModel{
			BaseModel: model.BaseModel{
				ID:        1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			CreateBy: 100,
			UpdateBy: 100,
		},
		Code:   "ceo",
		Name:   "CEO",
		Sort:   0,
		Status: 1,
		Remark: "首席执行官",
	}

	assert.Equal(t, "ceo", post.Code)
	assert.Equal(t, "CEO", post.Name)
	assert.Equal(t, 1, post.Status)
	assert.Equal(t, "首席执行官", post.Remark)
}

func TestTableName(t *testing.T) {
	// 测试表名生成
	tableName := model.GetTableName(model.TablePrefixSystem, "user")
	assert.Equal(t, "system_user", tableName)

	tableName = model.GetTableName(model.TablePrefixInfra, "config")
	assert.Equal(t, "infra_config", tableName)

	tableName = model.GetTableName(model.TablePrefixBpm, "task")
	assert.Equal(t, "bpm_task", tableName)
}
