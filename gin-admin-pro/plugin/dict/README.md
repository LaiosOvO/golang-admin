# 字典管理插件

字典管理插件提供系统字典数据的统一管理功能，支持字典类型和字典数据的增删改查，并提供缓存机制提高访问性能。

## 功能特性

- 字典类型管理：增删改查字典类型
- 字典数据管理：增删改查字典数据
- 数据缓存：支持Redis缓存，提高访问性能
- 自动刷新：支持定时刷新缓存数据
- 数据验证：提供字典数据有效性验证
- 导出功能：支持字典数据导出
- 简单列表：提供下拉选择用的简化数据

## 使用方法

### 1. 初始化服务

```go
import "gin-admin-pro/plugin/dict"

// 创建字典插件
dictPlugin := dict.NewPlugin(db, nil)

// 初始化插件
err := dictPlugin.Init()
if err != nil {
    log.Fatal("Failed to init dict plugin:", err)
}

// 获取字典服务
dictService := dictPlugin.GetService()
```

### 2. 字典类型管理

```go
// 创建字典类型
dictType := &dict.DictType{
    Name:   "用户状态",
    Type:   "user_status",
    Status: 1,
    Remark: "用户状态字典",
}
err := dictService.CreateDictType(dictType)

// 查询字典类型列表
dictTypes, total, err := dictService.GetDictTypes(1, 20, "", "user_status", 1)

// 获取字典类型详情
dictType, err := dictService.GetDictTypeByID(1)

// 更新字典类型
dictType.Remark = "更新后的备注"
err := dictService.UpdateDictType(dictType)

// 删除字典类型
err := dictService.DeleteDictType(1)
```

### 3. 字典数据管理

```go
// 创建字典数据
dictData := &dict.DictData{
    DictSort: 1,
    Label:    "启用",
    Value:    "1",
    DictType: "user_status",
    Status:   1,
}
err := dictService.CreateDictData(dictData)

// 查询字典数据列表
dictDataList, total, err := dictService.GetDictDataList(1, 20, "user_status", "", 1)

// 根据类型获取所有启用的字典数据
dictDataList, err := dictService.GetDictDataByType("user_status")

// 获取简单列表（用于下拉选择）
simpleList, err := dictService.GetDictDataSimple("user_status")
```

### 4. 字典数据转换

```go
// 根据标签获取值
value, err := dictService.GetDictValueByLabel("user_status", "启用")

// 根据值获取标签
label, err := dictService.GetDictLabelByValue("user_status", "1")
```

### 5. 缓存管理

```go
// 刷新字典缓存
err := dictService.RefreshDictCache()
```

## 配置说明

### 基础配置

```yaml
dict:
  enabled: true           # 是否启用字典管理
  enableCache: true       # 是否启用缓存
  cacheExpire: 3600       # 缓存过期时间（秒）
  cachePrefix: "dict:"    # 缓存前缀
  autoRefresh: false      # 是否启用自动刷新
  refreshInterval: 300    # 刷新间隔（秒）
  defaultPageSize: 20     # 默认分页大小
  maxPageSize: 100        # 最大分页大小
```

### 自定义配置

```go
customConfig := &dict.Config{
    Enabled:         true,
    EnableCache:     true,
    CacheExpire:     7200,    // 2小时
    CachePrefix:     "myapp:dict:",
    AutoRefresh:     true,
    RefreshInterval: 600,     // 10分钟
    DefaultPageSize: 10,
    MaxPageSize:     50,
}

dictPlugin := dict.NewPlugin(db, customConfig)
```

## 数据库表结构

### system_dict_type 表（字典类型表）

| 字段名 | 类型 | 说明 | 约束 |
|--------|------|------|------|
| id | bigint | 主键 | PRIMARY KEY |
| name | varchar(100) | 字典名称 | NOT NULL |
| type | varchar(100) | 字典类型 | NOT NULL, UNIQUE |
| status | tinyint | 状态 | DEFAULT 1 |
| remark | varchar(500) | 备注 | |
| create_by | bigint | 创建人 | |
| update_by | bigint | 更新人 | |
| create_at | datetime | 创建时间 | |
| update_at | datetime | 更新时间 | |

### system_dict_data 表（字典数据表）

| 字段名 | 类型 | 说明 | 约束 |
|--------|------|------|------|
| id | bigint | 主键 | PRIMARY KEY |
| dict_sort | int | 排序 | DEFAULT 0 |
| label | varchar(100) | 字典标签 | NOT NULL |
| value | varchar(100) | 字典值 | NOT NULL |
| dict_type | varchar(100) | 字典类型 | NOT NULL, INDEX |
| status | tinyint | 状态 | DEFAULT 1 |
| remark | varchar(500) | 备注 | |
| create_by | bigint | 创建人 | |
| update_by | bigint | 更新人 | |
| create_at | datetime | 创建时间 | |
| update_at | datetime | 更新时间 | |

## API接口规范

### 字典类型接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/system/dict-type/page | 分页查询字典类型 |
| GET | /api/v1/system/dict-type/get | 获取字典类型详情 |
| POST | /api/v1/system/dict-type/create | 创建字典类型 |
| PUT | /api/v1/system/dict-type/update | 更新字典类型 |
| DELETE | /api/v1/system/dict-type/delete | 删除字典类型 |
| GET | /api/v1/system/dict-type/list-all-simple | 获取字典类型简单列表 |

### 字典数据接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/system/dict-data/page | 分页查询字典数据 |
| GET | /api/v1/system/dict-data/get | 获取字典数据详情 |
| POST | /api/v1/system/dict-data/create | 创建字典数据 |
| PUT | /api/v1/system/dict-data/update | 更新字典数据 |
| DELETE | /api/v1/system/dict-data/delete | 删除字典数据 |
| GET | /api/v1/system/dict-data/type | 根据类型获取字典数据 |
| GET | /api/v1/system/dict-data/refresh | 刷新字典缓存 |

## 中间件集成

### 字典缓存中间件

```go
func DictCacheMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 将字典服务存储到上下文中
        dictService := dictPlugin.GetService()
        c.Set("dictService", dictService)
        c.Next()
    }
}
```

### 数据转换中间件

```go
func DictTransformMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 在响应中自动转换字典值
        c.Next()
        
        // 获取响应数据
        if data, exists := c.Get("responseData"); exists {
            if list, ok := data.([]map[string]interface{}); ok {
                // 转换字典值
                for i, item := range list {
                    if status, ok := item["status"].(int); ok {
                        if label, err := dictService.GetDictLabelByValue("system_status", fmt.Sprintf("%d", status)); err == nil {
                            item["statusLabel"] = label
                        }
                    }
                }
            }
        }
    }
}
```

## 缓存策略

### Redis缓存结构

```
dict:type:{dict_type}           # 字典类型信息
dict:data:{dict_type}           # 字典数据列表
dict:label:{dict_type}:{value}  # 根据值获取标签
dict:value:{dict_type}:{label}  # 根据标签获取值
```

### 缓存更新策略

1. **创建/更新/删除时自动刷新**：修改字典数据时自动清除相关缓存
2. **定时刷新**：配置自动刷新时，定时从数据库重新加载数据
3. **手动刷新**：提供API接口手动刷新缓存

### 缓存键命名规范

```
{prefix}:{type}:{identifier}
```

- `prefix`: 配置的缓存前缀，默认为 "dict:"
- `type`: 缓存类型（type/data/label/value）
- `identifier`: 具体标识符

## 默认字典数据

插件初始化时会自动创建以下默认字典：

### 系统状态字典（system_status）
- 启用 (1)
- 禁用 (0)

### 用户性别字典（user_gender）
- 男 (1)
- 女 (2)
- 保密 (3)

### 数据权限字典（data_scope）
- 全部数据权限 (1)
- 自定义数据权限 (2)
- 本部门数据权限 (3)
- 本部门及以下数据权限 (4)
- 仅本人数据权限 (5)

## 使用示例

### 在用户管理中的应用

```go
// 获取用户状态选项
statusOptions, err := dictService.GetDictDataSimple("system_status")

// 转换用户状态显示
func getUserStatusLabel(status int) string {
    label, _ := dictService.GetDictLabelByValue("system_status", fmt.Sprintf("%d", status))
    return label
}

// 验证状态值有效性
func validateUserStatus(status int) bool {
    err := dictService.ValidateDictType("system_status")
    return err == nil
}
```

### 在表单验证中的应用

```go
// 验证字典值
func validateDictValue(dictType, value string) bool {
    _, err := dictService.GetDictLabelByValue(dictType, value)
    return err == nil
}

// 获取字典选项用于下拉选择
func getDictOptions(dictType string) ([]map[string]string, error) {
    dataList, err := dictService.GetDictDataSimple(dictType)
    if err != nil {
        return nil, err
    }
    
    var options []map[string]string
    for _, data := range dataList {
        options = append(options, map[string]string{
            "label": data.Label,
            "value": data.Value,
        })
    }
    
    return options, nil
}
```

## 注意事项

1. **性能优化**：字典数据应启用缓存，减少数据库查询
2. **数据一致性**：修改字典数据后记得刷新缓存
3. **规范命名**：字典类型使用英文和下划线，避免特殊字符
4. **状态管理**：禁用的字典类型和数据不会在查询中返回
5. **排序规则**：字典数据按排序字段和ID排序显示

## 测试

```go
func TestDictService(t *testing.T) {
    service := NewService(testDB)
    
    // 测试创建字典类型
    dictType := &DictType{
        Name:   "测试类型",
        Type:   "test_type",
        Status: 1,
    }
    err := service.CreateDictType(dictType)
    assert.NoError(t, err)
    
    // 测试创建字典数据
    dictData := &DictData{
        DictSort: 1,
        Label:    "测试标签",
        Value:    "test_value",
        DictType: "test_type",
        Status:   1,
    }
    err = service.CreateDictData(dictData)
    assert.NoError(t, err)
    
    // 测试查询
    dataList, err := service.GetDictDataByType("test_type")
    assert.NoError(t, err)
    assert.Len(t, dataList, 1)
    
    // 测试转换
    label, err := service.GetDictLabelByValue("test_type", "test_value")
    assert.NoError(t, err)
    assert.Equal(t, "测试标签", label)
}
```