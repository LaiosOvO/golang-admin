package system

import (
	"gin-admin-pro/internal/model"
	"gorm.io/gorm"
)

// User 用户表
type User struct {
	model.AuditModel
	Username  string          `gorm:"size:30;not null;uniqueIndex" json:"username"`
	Nickname  string          `gorm:"size:30" json:"nickname"`
	Password  string          `gorm:"size:100;not null" json:"-"`
	Mobile    string          `gorm:"size:11" json:"mobile"`
	Email     string          `gorm:"size:50" json:"email"`
	Avatar    string          `gorm:"size:512" json:"avatar"`
	Status    int             `gorm:"default:1" json:"status"` // 0-禁用 1-启用
	LoginIP   string          `gorm:"size:50" json:"loginIP"`
	LoginDate *gorm.DeletedAt `json:"loginDate"`
	DeptID    uint            `json:"deptId"`
	Dept      *Dept           `gorm:"foreignKey:DeptID" json:"dept,omitempty"`
	PostIDs   string          `gorm:"size:255" json:"postIds"` // 岗位ID列表，逗号分隔
	Posts     []Post          `gorm:"many2many:user_post;" json:"posts,omitempty"`
	Roles     []Role          `gorm:"many2many:user_role;" json:"roles,omitempty"`
}

// Role 角色表
type Role struct {
	model.AuditModel
	Code      string `gorm:"size:100;not null;uniqueIndex" json:"code"`
	Name      string `gorm:"size:30;not null" json:"name"`
	Sort      int    `gorm:"default:0" json:"sort"`
	DataScope int    `gorm:"default:1" json:"dataScope"` // 数据范围 1-全部数据权限 2-自定义数据权限 3-本部门数据权限 4-本部门及以下数据权限 5-仅本人数据权限
	Status    int    `gorm:"default:1" json:"status"`    // 0-禁用 1-启用
	Type      int    `gorm:"default:1" json:"type"`      // 角色类型 1-内置角色 2-自定义角色
	Remark    string `gorm:"size:500" json:"remark"`
	Users     []User `gorm:"many2many:user_role;" json:"users,omitempty"`
	Menus     []Menu `gorm:"many2many:role_menu;" json:"menus,omitempty"`
}

// Menu 菜单表
type Menu struct {
	model.TreeModel
	Type          int    `gorm:"not null" json:"type"`          // 菜单类型 1-目录 2-菜单 3-按钮
	Icon          string `gorm:"size:100" json:"icon"`          // 菜单图标
	Component     string `gorm:"size:255" json:"component"`     // 组件路径
	ComponentName string `gorm:"size:255" json:"componentName"` // 组件名
	Perms         string `gorm:"size:100" json:"perms"`         // 权限标识
	Status        int    `gorm:"default:1" json:"status"`       // 0-禁用 1-启用
	Visible       int    `gorm:"default:1" json:"visible"`      // 0-隐藏 1-显示
	KeepAlive     int    `gorm:"default:1" json:"keepAlive"`    // 0-关闭 1-开启
	AlwaysShow    int    `gorm:"default:1" json:"alwaysShow"`   // 0-关闭 1-开启
	Parent        *Menu  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children      []Menu `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Roles         []Role `gorm:"many2many:role_menu;" json:"roles,omitempty"`
}

// Dept 部门表
type Dept struct {
	model.TreeModel
	LeaderUserId uint   `json:"leaderUserId"` // 负责人用户ID
	Phone        string `gorm:"size:11" json:"phone"`
	Email        string `gorm:"size:50" json:"email"`
	Status       int    `gorm:"default:1" json:"status"` // 0-禁用 1-启用
	Parent       *Dept  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children     []Dept `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Leader       *User  `gorm:"foreignKey:LeaderUserId" json:"leader,omitempty"`
	Users        []User `gorm:"foreignKey:DeptID" json:"users,omitempty"`
}

// Post 岗位表
type Post struct {
	model.AuditModel
	Code   string `gorm:"size:64;not null;uniqueIndex" json:"code"`
	Name   string `gorm:"size:50;not null" json:"name"`
	Sort   int    `gorm:"default:0" json:"sort"`
	Status int    `gorm:"default:1" json:"status"` // 0-禁用 1-启用
	Remark string `gorm:"size:500" json:"remark"`
	Users  []User `gorm:"many2many:user_post;" json:"users,omitempty"`
}

// UserRole 用户角色关联表
type UserRole struct {
	ID     uint `gorm:"primarykey" json:"id"`
	UserID uint `gorm:"not null;index" json:"userId"`
	RoleID uint `gorm:"not null;index" json:"roleId"`
	User   User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role   Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// RoleMenu 角色菜单关联表
type RoleMenu struct {
	ID     uint `gorm:"primarykey" json:"id"`
	RoleID uint `gorm:"not null;index" json:"roleId"`
	MenuID uint `gorm:"not null;index" json:"menuId"`
	Role   Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Menu   Menu `gorm:"foreignKey:MenuID" json:"menu,omitempty"`
}

// UserPost 用户岗位关联表
type UserPost struct {
	ID     uint `gorm:"primarykey" json:"id"`
	UserID uint `gorm:"not null;index" json:"userId"`
	PostID uint `gorm:"not null;index" json:"postId"`
	User   User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Post   Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
}
