# 数据权限插件

数据权限插件基于角色的数据范围设置，实现对数据的细粒度访问控制。

## 功能特性

- 支持5种数据权限范围：
  - 全部数据权限
  - 自定义数据权限
  - 本部门数据权限
  - 本部门及以下数据权限
  - 仅本人数据权限
- 支持部门树形结构的权限继承
- 支持角色级别的权限设置
- 提供GORM查询Scope扩展
- 支持权限检查和验证

## 使用方法

### 1. 初始化服务

```go
import "gin-admin-pro/plugin/dataperm"

// 创建数据权限服务
dataPermService := dataperm.NewService(db)
```

### 2. 在查询中应用数据权限

```go
// 获取用户列表（应用数据权限）
var users []system.User
err = db.Model(&system.User{}).
    Scopes(dataperm.Scope(userID, "dept_id", "id")).
    Find(&users).Error
```

### 3. 检查数据权限

```go
// 检查用户是否有权限访问指定部门数据
hasPermission := dataPermService.CheckDataPermission(userID, deptID, "dept")

// 使用GORM Scope检查权限
err = db.Model(&system.Dept{}).
    Scopes(dataperm.HasDeptDataPermission(userID, deptID)).
    First(&dept).Error
```

### 4. 获取用户权限信息

```go
// 获取用户的数据权限范围
dataScope, err := dataPermService.GetDataScope(userID)

// 获取用户可访问的部门ID列表
deptIDs, err := dataPermService.GetDataScopeDeptIDs(userID)
```

## 权限范围说明

| 值 | 名称 | 说明 |
|----|------|------|
| 1 | 全部数据权限 | 可以查看所有数据 |
| 2 | 自定义数据权限 | 只能查看指定部门的数据 |
| 3 | 本部门数据权限 | 只能查看本部门的数据 |
| 4 | 本部门及以下数据权限 | 可以查看本部门及子部门的数据 |
| 5 | 仅本人数据权限 | 只能查看自己的数据 |

## 数据库表结构

### role_dept 表

自定义数据权限的角色部门关联表：

```sql
CREATE TABLE role_dept (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    role_id BIGINT NOT NULL COMMENT '角色ID',
    dept_id BIGINT NOT NULL COMMENT '部门ID',
    INDEX idx_role_id (role_id),
    INDEX idx_dept_id (dept_id)
) COMMENT='角色部门关联表';
```

## 配置说明

### 角色数据权限配置

在角色管理中，可以设置角色的数据权限范围：

1. **全部数据权限**：用户可以看到所有数据
2. **自定义数据权限**：需要为角色分配可访问的部门
3. **本部门数据权限**：用户只能看到自己部门的数据
4. **本部门及以下数据权限**：用户可以看到自己部门及下级部门的数据
5. **仅本人数据权限**：用户只能看到自己的数据

### 部门树结构

部门支持树形结构，数据权限会根据部门层级自动继承：

- 本部门权限：只查看当前部门
- 本部门及以下权限：查看当前部门和所有子部门

## API集成

数据权限通常与中间件配合使用：

```go
// 在中间件中应用数据权限
func DataPermissionMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := getCurrentUserID(c)
        
        // 将数据权限服务存储到上下文中
        c.Set("dataPermService", dataperm.NewService(db))
        
        c.Next()
    }
}
```

## 注意事项

1. **性能优化**：数据权限查询会使用子查询，建议为相关字段添加索引
2. **权限继承**：部门权限会自动继承子部门的数据
3. **缓存策略**：可以考虑缓存用户的权限信息以提高性能
4. **权限验证**：在修改数据时也需要验证数据权限

## 示例代码

### 用户列表查询

```go
func GetUserList(c *gin.Context) {
    userID := getCurrentUserID(c)
    
    var users []system.User
    query := db.Model(&system.User{})
    
    // 应用数据权限
    query = query.Scopes(dataperm.Scope(userID, "dept_id", "id"))
    
    // 其他查询条件
    if username := c.Query("username"); username != "" {
        query = query.Where("username LIKE ?", "%"+username+"%")
    }
    
    err := query.Find(&users).Error
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{"data": users})
}
```

### 部门数据检查

```go
func UpdateDept(c *gin.Context) {
    userID := getCurrentUserID(c)
    deptID := c.Param("id")
    
    // 检查权限
    err := db.Model(&system.Dept{}).
        Scopes(dataperm.HasDeptDataPermission(userID, deptID)).
        First(&dept).Error
        
    if err != nil {
        c.JSON(403, gin.H{"error": "无权限访问该部门"})
        return
    }
    
    // 执行更新操作
    // ...
}
```

## 测试

```go
func TestDataPermission(t *testing.T) {
    service := dataperm.NewService(testDB)
    
    // 测试获取数据权限范围
    dataScope, err := service.GetDataScope(1)
    assert.NoError(t, err)
    assert.Equal(t, dataperm.DataScopeDept, dataScope)
    
    // 测试构建数据权限SQL
    query := testDB.Model(&system.User{})
    query = service.BuildDataScopeSQL(query, 1, "dept_id", "id")
    
    // 验证生成的SQL条件
    // ...
}
```