# ORM集成与模型设计

## 概述

本文档记录了阶段4：ORM集成与模型设计的实现过程，包括GORM配置、基础模型定义、系统核心模型设计和数据库迁移工具的开发。

## 实现内容

### 1. GORM配置检查

通过检查现有的MySQL插件配置，确认了GORM已经正确配置：

- ✅ 单数表名配置 (`SingularTable: true`)
- ✅ 日志级别配置
- ✅ 连接池配置
- ✅ 自动迁移方法 (`AutoMigrate`)
- ✅ 表检查方法 (`HasTable`)

### 2. 基础模型设计

#### 2.1 BaseModel

位置：`internal/model/base.go`

```go
type BaseModel struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `json:"createdAt"`
    UpdatedAt time.Time      `json:"updatedAt"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

**功能特性：**
- 主键ID (uint类型，自增)
- 创建时间 (自动设置)
- 更新时间 (自动更新)
- 软删除标记 (GORM DeletedAt类型)

**GORM钩子：**
- `BeforeCreate`: 自动设置创建时间和更新时间
- `BeforeUpdate`: 自动更新更新时间

#### 2.2 AuditModel

继承自BaseModel，添加审计字段：

```go
type AuditModel struct {
    BaseModel
    CreateBy uint   `json:"createBy"`
    UpdateBy uint   `json:"updateBy"`
    Remark   string `gorm:"size:500" json:"remark"`
}
```

**功能特性：**
- 创建人ID
- 更新人ID
- 备注信息 (最大500字符)

**辅助方法：**
- `SetCreatedBy()`: 设置创建人
- `SetUpdatedBy()`: 设置更新人
- `GetCreatedBy()`: 获取创建人
- `GetUpdatedBy()`: 获取更新人

#### 2.3 TreeModel

继承自AuditModel，支持树形结构：

```go
type TreeModel struct {
    AuditModel
    ParentID  uint   `json:"parentId"`
    Level     int    `json:"level"`
    Sort      int    `json:"sort"`
    Name      string `gorm:"size:100;not null" json:"name"`
    Path      string `gorm:"size:500" json:"path"`
    Ancestors string `gorm:"size:500" json:"ancestors"`
}
```

**功能特性：**
- 父级ID
- 层级深度
- 排序字段
- 名称字段
- 路径字段
- 祖先路径

### 3. 系统核心模型

位置：`internal/model/system/system.go`

#### 3.1 用户表 (system_user)

```go
type User struct {
    AuditModel
    Username    string `gorm:"size:30;not null;uniqueIndex" json:"username"`
    Nickname    string `gorm:"size:30" json:"nickname"`
    Password    string `gorm:"size:100;not null" json:"-"`
    Mobile      string `gorm:"size:11" json:"mobile"`
    Email       string `gorm:"size:50" json:"email"`
    Avatar      string `gorm:"size:512" json:"avatar"`
    Status      int    `gorm:"default:1" json:"status"`
    LoginIP     string `gorm:"size:50" json:"loginIP"`
    LoginDate   *gorm.DeletedAt `json:"loginDate"`
    DeptID      uint   `json:"deptId"`
    Dept        *Dept  `gorm:"foreignKey:DeptID" json:"dept,omitempty"`
    PostIDs     string `gorm:"size:255" json:"postIds"`
    Posts       []Post `gorm:"many2many:user_post;" json:"posts,omitempty"`
    Roles       []Role `gorm:"many2many:user_role;" json:"roles,omitempty"`
}
```

**字段说明：**
- `Username`: 用户名 (30字符，唯一索引)
- `Password`: 密码 (JSON响应时隐藏)
- `Status`: 状态 (0-禁用 1-启用)
- `DeptID`: 所属部门ID
- `PostIDs`: 岗位ID列表 (逗号分隔)
- `Posts`: 多对多关联岗位
- `Roles`: 多对多关联角色

#### 3.2 角色表 (system_role)

```go
type Role struct {
    AuditModel
    Code      string `gorm:"size:100;not null;uniqueIndex" json:"code"`
    Name      string `gorm:"size:30;not null" json:"name"`
    Sort      int    `gorm:"default:0" json:"sort"`
    DataScope int    `gorm:"default:1" json:"dataScope"`
    Status    int    `gorm:"default:1" json:"status"`
    Type      int    `gorm:"default:1" json:"type"`
    Remark    string `gorm:"size:500" json:"remark"`
    Users     []User `gorm:"many2many:user_role;" json:"users,omitempty"`
    Menus     []Menu `gorm:"many2many:role_menu;" json:"menus,omitempty"`
}
```

**数据权限范围：**
- 1: 全部数据权限
- 2: 自定义数据权限
- 3: 本部门数据权限
- 4: 本部门及以下数据权限
- 5: 仅本人数据权限

**角色类型：**
- 1: 内置角色
- 2: 自定义角色

#### 3.3 菜单表 (system_menu)

```go
type Menu struct {
    TreeModel
    Type        int    `gorm:"not null" json:"type"`
    Icon        string `gorm:"size:100" json:"icon"`
    Component   string `gorm:"size:255" json:"component"`
    ComponentName string `gorm:"size:255" json:"componentName"`
    Perms       string `gorm:"size:100" json:"perms"`
    Status      int    `gorm:"default:1" json:"status"`
    Visible     int    `gorm:"default:1" json:"visible"`
    KeepAlive   int    `gorm:"default:1" json:"keepAlive"`
    AlwaysShow  int    `gorm:"default:1" json:"alwaysShow"`
    Parent      *Menu  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
    Children    []Menu `gorm:"foreignKey:ParentID" json:"children,omitempty"`
    Roles       []Role `gorm:"many2many:role_menu;" json:"roles,omitempty"`
}
```

**菜单类型：**
- 1: 目录
- 2: 菜单
- 3: 按钮

#### 3.4 部门表 (system_dept)

```go
type Dept struct {
    TreeModel
    LeaderUserId uint   `json:"leaderUserId"`
    Phone        string `gorm:"size:11" json:"phone"`
    Email        string `gorm:"size:50" json:"email"`
    Status       int    `gorm:"default:1" json:"status"`
    Parent       *Dept  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
    Children     []Dept `gorm:"foreignKey:ParentID" json:"children,omitempty"`
    Leader       *User  `gorm:"foreignKey:LeaderUserId" json:"leader,omitempty"`
    Users        []User `gorm:"foreignKey:DeptID" json:"users,omitempty"`
}
```

#### 3.5 岗位表 (system_post)

```go
type Post struct {
    AuditModel
    Code   string `gorm:"size:64;not null;uniqueIndex" json:"code"`
    Name   string `gorm:"size:50;not null" json:"name"`
    Sort   int    `gorm:"default:0" json:"sort"`
    Status int    `gorm:"default:1" json:"status"`
    Remark string `gorm:"size:500" json:"remark"`
    Users  []User `gorm:"many2many:user_post;" json:"users,omitempty"`
}
```

#### 3.6 关联表

多对多关联表：
- `user_role`: 用户角色关联
- `role_menu`: 角色菜单关联
- `user_post`: 用户岗位关联

### 4. 数据库迁移工具

位置：`internal/migration/migrator.go`

#### 4.1 迁移器结构

```go
type Migrator struct {
    db *gorm.DB
}
```

#### 4.2 主要功能

##### 自动迁移

```go
func (m *Migrator) AutoMigrate() error
```

**功能：**
- 自动迁移所有系统表
- 创建必要索引
- 插入初始数据

**迁移的表：**
- system_user (用户表)
- system_role (角色表)
- system_menu (菜单表)
- system_dept (部门表)
- system_post (岗位表)
- user_role (用户角色关联表)
- role_menu (角色菜单关联表)
- user_post (用户岗位关联表)

##### 索引创建

自动创建以下索引：
- 用户表：username, mobile, email, dept_id, status
- 角色表：code, status
- 菜单表：parent_id, type, status
- 部门表：parent_id, status
- 岗位表：code, status

##### 初始数据

自动插入：
- 超级管理员角色 (super_admin)
- 管理员角色 (admin)
- 普通用户角色 (common)
- 根部门 (总公司)

#### 4.3 命令行工具

位置：`cmd/migrate/main.go`

**使用方法：**

```bash
# 自动迁移
go run cmd/migrate/main.go -action=migrate -env=dev

# 重置数据库（删除所有表并重新创建）
go run cmd/migrate/main.go -action=reset -env=dev

# 删除所有表
go run cmd/migrate/main.go -action=drop -env=dev
```

**参数说明：**
- `-action`: 迁移动作 (migrate/reset/drop)
- `-env`: 环境配置 (dev/test/prod)

### 5. 表名前缀规范

定义了统一的表名前缀：

```go
const (
    TablePrefixSystem = "system_"  // 系统管理模块
    TablePrefixInfra  = "infra_"   // 基础设施模块
    TablePrefixBpm    = "bpm_"     // 工作流模块
    TablePrefixPay    = "pay_"     // 支付模块
    TablePrefixMember = "member_"  // 会员模块
    TablePrefixMall   = "mall_"    // 商城模块
    TablePrefixReport = "report_"  // 报表模块
)
```

### 6. 单元测试

位置：`internal/migration/migrator_test.go`

**测试覆盖：**
- ✅ 自动迁移功能
- ✅ 基础模型方法
- ✅ 审计模型方法
- ✅ 系统模型字段验证
- ✅ 表名生成

**运行测试：**

```bash
go test ./internal/migration/...
```

## 与Ruoyi-Vue-Pro的对应关系

| 表名 | 说明 | 对应关系 |
|------|------|----------|
| system_user | 用户表 | ✅ 完全对应 |
| system_role | 角色表 | ✅ 完全对应 |
| system_menu | 菜单表 | ✅ 完全对应 |
| system_dept | 部门表 | ✅ 完全对应 |
| system_post | 岗位表 | ✅ 完全对应 |
| user_role | 用户角色关联表 | ✅ 完全对应 |
| role_menu | 角色菜单关联表 | ✅ 完全对应 |
| user_post | 用户岗位关联表 | ✅ 完全对应 |

## 字段设计规范

### 命名规范
- 使用小写字母和下划线
- 主键统一使用 `id`
- 外键使用 `表名_id` 格式
- 时间字段使用 `created_at`, `updated_at`
- 状态字段使用 `status`

### 数据类型
- ID: uint (自增主键)
- 文本内容: varchar (根据内容长度指定)
- 状态: int (0-禁用 1-启用)
- 时间: time.Time
- 外键: uint

### 索引规范
- 主键自动索引
- 外键索引
- 唯一字段唯一索引
- 查询频繁字段普通索引

## 下一步计划

1. ✅ 完成ORM集成和模型定义
2. ✅ 完成数据库迁移工具
3. ⏳ 开始阶段5：Gin框架配置
4. ⏳ 实现用户管理API接口
5. ⏳ 集成JWT认证功能

## 技术要点总结

### 优势
1. **统一的模型设计**: 基础模型提供通用的CRUD功能
2. **完善的审计功能**: 自动记录创建人、更新人、操作时间
3. **灵活的软删除**: 支持数据恢复和查询过滤
4. **树形结构支持**: 适用于菜单、部门等层级数据
5. **自动化迁移**: 简化数据库初始化和更新流程

### 注意事项
1. **外键约束**: 使用GORM外键关联，保证数据一致性
2. **索引优化**: 根据查询需求创建合适的索引
3. **数据迁移**: 提供安全的数据库迁移和重置功能
4. **测试覆盖**: 保证模型和迁移功能的正确性

---

**阶段4完成时间**: 2025-01-20  
**完成状态**: ✅ 已完成  
**下一步**: 阶段5 - Gin框架配置