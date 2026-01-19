package migration

import (
	"fmt"
	"gin-admin-pro/internal/model"
	"gin-admin-pro/internal/model/system"
	"gorm.io/gorm"
	"log"
)

// Migrator 数据库迁移器
type Migrator struct {
	db *gorm.DB
}

// NewMigrator 创建迁移器
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db: db}
}

// AutoMigrate 自动迁移所有表
func (m *Migrator) AutoMigrate() error {
	log.Println("开始数据库迁移...")

	// 系统管理模块表
	models := []interface{}{
		// 用户相关
		&system.User{},
		&system.Role{},
		&system.Menu{},
		&system.Dept{},
		&system.Post{},

		// 关联表
		&system.UserRole{},
		&system.RoleMenu{},
		&system.UserPost{},
	}

	// 执行迁移
	for _, model := range models {
		if err := m.db.AutoMigrate(model); err != nil {
			return fmt.Errorf("迁移表 %T 失败: %w", model, err)
		}
		log.Printf("成功迁移表: %T", model)
	}

	// 创建索引
	if err := m.createIndexes(); err != nil {
		return fmt.Errorf("创建索引失败: %w", err)
	}

	// 插入初始数据
	if err := m.insertInitialData(); err != nil {
		return fmt.Errorf("插入初始数据失败: %w", err)
	}

	log.Println("数据库迁移完成!")
	return nil
}

// createIndexes 创建索引
func (m *Migrator) createIndexes() error {
	log.Println("创建索引...")

	// 用户表索引
	indexes := []struct {
		table   string
		index   string
		columns []string
		comment string
	}{
		{"system_user", "idx_username", []string{"username"}, "用户名索引"},
		{"system_user", "idx_mobile", []string{"mobile"}, "手机号索引"},
		{"system_user", "idx_email", []string{"email"}, "邮箱索引"},
		{"system_user", "idx_dept_id", []string{"dept_id"}, "部门ID索引"},
		{"system_user", "idx_status", []string{"status"}, "状态索引"},

		{"system_role", "idx_code", []string{"code"}, "角色编码索引"},
		{"system_role", "idx_status", []string{"status"}, "状态索引"},

		{"system_menu", "idx_parent_id", []string{"parent_id"}, "父级ID索引"},
		{"system_menu", "idx_type", []string{"type"}, "菜单类型索引"},
		{"system_menu", "idx_status", []string{"status"}, "状态索引"},

		{"system_dept", "idx_parent_id", []string{"parent_id"}, "父级ID索引"},
		{"system_dept", "idx_status", []string{"status"}, "状态索引"},

		{"system_post", "idx_code", []string{"code"}, "岗位编码索引"},
		{"system_post", "idx_status", []string{"status"}, "状态索引"},
	}

	for _, idx := range indexes {
		// 检查索引是否存在
		var count int64
		err := m.db.Raw(`
			SELECT COUNT(*) FROM information_schema.statistics 
			WHERE table_schema = DATABASE() 
			AND table_name = ? 
			AND index_name = ?
		`, idx.table, idx.index).Scan(&count).Error

		if err != nil {
			log.Printf("检查索引 %s 失败: %v", idx.index, err)
			continue
		}

		if count == 0 {
			// 创建索引
			columns := ""
			for i, col := range idx.columns {
				if i > 0 {
					columns += ", "
				}
				columns += col
			}

			err := m.db.Exec(fmt.Sprintf(
				"CREATE INDEX %s ON %s (%s) COMMENT '%s'",
				idx.index, idx.table, columns, idx.comment,
			)).Error

			if err != nil {
				log.Printf("创建索引 %s 失败: %v", idx.index, err)
			} else {
				log.Printf("成功创建索引: %s", idx.index)
			}
		}
	}

	return nil
}

// insertInitialData 插入初始数据
func (m *Migrator) insertInitialData() error {
	log.Println("插入初始数据...")

	// 插入超级管理员角色
	var count int64
	if err := m.db.Model(&system.Role{}).Where("code = ?", "super_admin").Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		superAdmin := &system.Role{
			Code:      "super_admin",
			Name:      "超级管理员",
			Sort:      0,
			DataScope: 1, // 全部数据权限
			Status:    1,
			Type:      1, // 内置角色
			Remark:    "系统内置超级管理员角色，拥有所有权限",
		}
		if err := m.db.Create(superAdmin).Error; err != nil {
			return err
		}
		log.Println("插入超级管理员角色成功")
	}

	// 插入管理员角色
	if err := m.db.Model(&system.Role{}).Where("code = ?", "admin").Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		admin := &system.Role{
			Code:      "admin",
			Name:      "管理员",
			Sort:      1,
			DataScope: 1, // 全部数据权限
			Status:    1,
			Type:      1, // 内置角色
			Remark:    "系统内置管理员角色",
		}
		if err := m.db.Create(admin).Error; err != nil {
			return err
		}
		log.Println("插入管理员角色成功")
	}

	// 插入普通用户角色
	if err := m.db.Model(&system.Role{}).Where("code = ?", "common").Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		common := &system.Role{
			Code:      "common",
			Name:      "普通用户",
			Sort:      2,
			DataScope: 4, // 仅本人数据权限
			Status:    1,
			Type:      1, // 内置角色
			Remark:    "系统内置普通用户角色",
		}
		if err := m.db.Create(common).Error; err != nil {
			return err
		}
		log.Println("插入普通用户角色成功")
	}

	// 插入根部门
	if err := m.db.Model(&system.Dept{}).Where("parent_id = ?", 0).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		rootDept := &system.Dept{
			TreeModel: model.TreeModel{
				AuditModel: model.AuditModel{},
				ParentID:   0,
				Level:      0,
				Sort:       0,
				Name:       "总公司",
				Ancestors:  "0",
				Path:       "0",
			},
			LeaderUserId: 0,
			Phone:        "",
			Email:        "",
		}
		if err := m.db.Create(rootDept).Error; err != nil {
			return err
		}
		log.Println("插入根部门成功")
	}

	log.Println("初始数据插入完成")
	return nil
}

// DropAllTables 删除所有表（慎用）
func (m *Migrator) DropAllTables() error {
	log.Println("警告：正在删除所有表...")

	tables := []string{
		"user_post",
		"role_menu",
		"user_role",
		"system_post",
		"system_dept",
		"system_menu",
		"system_role",
		"system_user",
	}

	for _, table := range tables {
		if m.db.Migrator().HasTable(table) {
			if err := m.db.Migrator().DropTable(table); err != nil {
				return fmt.Errorf("删除表 %s 失败: %w", table, err)
			}
			log.Printf("删除表: %s", table)
		}
	}

	log.Println("所有表删除完成")
	return nil
}

// ResetDatabase 重置数据库（删除所有表并重新迁移）
func (m *Migrator) ResetDatabase() error {
	log.Println("开始重置数据库...")

	// 删除所有表
	if err := m.DropAllTables(); err != nil {
		return err
	}

	// 重新迁移
	if err := m.AutoMigrate(); err != nil {
		return err
	}

	log.Println("数据库重置完成")
	return nil
}
