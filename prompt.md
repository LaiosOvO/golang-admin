ralph "
# Gin-Vue-Admin 克隆项目 - 参考若依Vue Pro实现

## 重要说明
- 这是一个大型企业级后台管理系统开发项目
- 严格按照阶段顺序执行，每完成一个小功能就提交代码
- 代码提交信息**不要包含**任何AI模型信息（如GPT、Claude等字样）
- 每个步骤都要在 \`docs/tech/step/\` 下记录实现逻辑
- **只有在所有阶段都完成后，才输出 <promise>GIN_ADMIN_PROJECT_COMPLETE</promise>**
- 在此之前**不要**输出任何 <promise> 标记

---

## 项目目标

基于 gin-vue-admin 架构，参考若依Vue Pro (ruoyi-vue-pro) 的实现规范，构建一个完整的企业级后台管理系统。

### 技术栈要求

#### 后端框架
- **语言**: Golang
- **Web框架**: Gin
- **ORM**: GORM

#### 数据库
- **关系型数据库**:
  - MySQL（主数据库）
  - PostgreSQL（包含 PostGIS 插件用于地理信息，pgvector 插件用于向量检索）
- **非关系型数据库**:
  - MongoDB（文档存储）
  - Redis（缓存、Session）
  - Elasticsearch（全文搜索）
- **向量数据库**:
  - Milvus（向量检索）

#### 中间件
- **消息队列**: Kafka
- **定时任务**: Cron 或其他调度框架

#### AI 集成
- 参考 xiaozhi-server（Golang版本）集成 AI 对话功能
- 支持多种 AI 模型接入

#### 参考项目
1. **架构参考**: https://github.com/flipped-aurora/gin-vue-admin
2. **接口规范参考**: /Volumes/T7/workspace/company/studio/code/admin/ruoyi-vue-pro
3. **AI 集成参考**: xiaozhi-server (Golang)

---

## 阶段 1：项目初始化与技术调研

### 任务清单
- [ ] 克隆并研究 gin-vue-admin 项目结构
- [ ] 分析 ruoyi-vue-pro 的接口规范和响应格式
- [ ] 研究 xiaozhi-server 的 AI 集成方案
- [ ] 设计项目目录结构
- [ ] 编写技术选型文档

### 产出文档
创建以下文档（全部中文）：
- \`docs/tech/01_技术调研报告.md\`
  - gin-vue-admin 架构分析
  - ruoyi-vue-pro 接口规范总结
  - xiaozhi-server AI 集成方案
  - 各数据库选型理由
  
- \`docs/tech/02_项目架构设计.md\`
  - 整体架构图
  - 目录结构说明
  - 模块划分
  - 数据流设计

- \`docs/tech/03_接口规范文档.md\`
  - 参考 ruoyi-vue-pro 的接口格式
  - 统一响应结构
  - 错误码定义
  - 分页规范

- \`docs/tech/04_数据库设计文档.md\`
  - MySQL 表设计
  - PostgreSQL GIS 数据设计
  - MongoDB 文档结构
  - Redis 缓存策略
  - Milvus 向量索引设计

### 项目目录结构（初始化）
\`\`\`
gin-admin-pro/
├── cmd/                    # 程序入口
│   └── server/
│       └── main.go
├── config/                 # 配置文件
│   ├── config.yaml         # 主配置
│   ├── config.dev.yaml     # 开发环境
│   └── config.prod.yaml    # 生产环境
├── internal/               # 内部代码
│   ├── api/                # API 控制器
│   ├── service/            # 业务逻辑
│   ├── dao/                # 数据访问层
│   ├── model/              # 数据模型
│   ├── middleware/         # 中间件
│   ├── router/             # 路由
│   └── pkg/                # 内部工具包
├── plugin/                 # 插件（参考 ruoyi starter）
│   ├── redis/
│   ├── mysql/
│   ├── postgresql/
│   ├── mongodb/
│   ├── elasticsearch/
│   ├── milvus/
│   ├── kafka/
│   ├── oss/                # 文件上传
│   └── ai/                 # AI 集成
├── docs/                   # 文档目录
│   ├── tech/               # 技术文档
│   │   ├── step/           # 实现步骤记录
│   │   └── api/            # API 文档
│   └── sql/                # 数据库脚本
├── scripts/                # 脚本
├── test/                   # 测试
├── docker/                 # Docker 配置
│   └── docker-compose.yaml
├── go.mod
├── go.sum
├── Makefile
└── README.md
\`\`\`

### Git 提交规范
提交信息格式：
\`\`\`
<type>: <description>

类型(type):
- init: 项目初始化
- feat: 新功能
- fix: 修复bug
- docs: 文档更新
- refactor: 重构
- test: 测试相关
- chore: 构建/工具链

示例:
init: 项目初始化，创建基础目录结构
feat: 实现配置文件加载功能
docs: 添加配置加载实现文档
\`\`\`

**不要在提交信息中包含**: GPT、Claude、AI生成、模型等字样

### 验收标准
- 项目目录创建完成
- 技术文档齐全
- 明确了实现路线

**完成后提交代码**: \`init: 项目初始化，创建基础架构和文档\`

**在 docs/tech/step/ 创建**: \`01_项目初始化.md\` 记录本阶段实现思路

---

## 阶段 2：配置文件加载（步骤1）

### 任务清单
- [ ] 使用 Viper 实现配置文件加载
- [ ] 支持多环境配置（dev/test/prod）
- [ ] 实现配置热加载（可选）
- [ ] 定义配置结构体

### 实现要求
创建 \`internal/pkg/config/\` 目录：
\`\`\`go
// config/config.yaml 示例
server:
  port: 8080
  mode: debug

database:
  mysql:
    host: localhost
    port: 3306
    database: gin_admin
    username: root
    password: password
  postgresql:
    host: localhost
    port: 5432
    extensions:
      - postgis
      - vector
  mongodb:
    uri: mongodb://localhost:27017
  redis:
    host: localhost
    port: 6379
  elasticsearch:
    url: http://localhost:9200
  milvus:
    host: localhost
    port: 19530

kafka:
  brokers:
    - localhost:9092

ai:
  enabled: true
  provider: openai
  api_key: sk-xxx
\`\`\`

### 文件清单
- \`internal/pkg/config/config.go\` - 配置加载逻辑
- \`internal/pkg/config/types.go\` - 配置结构体定义
- \`config/config.yaml\` - 配置文件模板

### 验收标准
- 可以成功加载配置文件
- 支持环境变量覆盖
- 有配置验证逻辑
- 单元测试通过

**完成后提交**: \`feat: 实现配置文件加载功能\`

**记录文档**: \`docs/tech/step/02_配置文件加载.md\`
内容包括：
- 使用的库（Viper）
- 配置结构设计
- 多环境支持方案
- 代码示例
- 测试方法

---

## 阶段 3：数据库配置与连接（步骤2-3）

### 任务清单

#### MySQL 配置
- [ ] 实现 MySQL 连接池
- [ ] 配置 GORM
- [ ] 实现数据库健康检查
- [ ] 创建 plugin/mysql/mysql.go

#### PostgreSQL 配置
- [ ] 实现 PostgreSQL 连接
- [ ] 启用 PostGIS 插件
- [ ] 启用 pgvector 插件
- [ ] 创建 plugin/postgresql/postgresql.go

#### MongoDB 配置
- [ ] 实现 MongoDB 连接
- [ ] 配置连接池
- [ ] 创建 plugin/mongodb/mongodb.go

#### Redis 配置
- [ ] 实现 Redis 连接
- [ ] 支持单机/集群模式
- [ ] 创建 plugin/redis/redis.go

#### Elasticsearch 配置
- [ ] 实现 ES 客户端
- [ ] 创建 plugin/elasticsearch/elasticsearch.go

#### Milvus 配置
- [ ] 实现 Milvus 连接
- [ ] 创建 plugin/milvus/milvus.go

### Plugin 目录结构（参考 ruoyi starter）
\`\`\`
plugin/
├── mysql/
│   ├── mysql.go          # MySQL 插件
│   ├── config.go         # 配置
│   └── README.md
├── postgresql/
│   ├── postgresql.go
│   ├── postgis.go        # PostGIS 扩展
│   ├── pgvector.go       # pgvector 扩展
│   └── README.md
├── mongodb/
│   ├── mongodb.go
│   └── README.md
├── redis/
│   ├── redis.go
│   ├── cache.go          # 缓存封装
│   └── README.md
├── elasticsearch/
│   ├── elasticsearch.go
│   └── README.md
└── milvus/
    ├── milvus.go
    └── README.md
\`\`\`

### 验收标准
- 所有数据库可以成功连接
- 有连接池管理
- 有错误处理和重连机制
- 每个 plugin 有独立的 README

**分多次提交**:
- \`feat: 实现MySQL数据库插件\`
- \`feat: 实现PostgreSQL插件及GIS/Vector扩展\`
- \`feat: 实现MongoDB插件\`
- \`feat: 实现Redis插件\`
- \`feat: 实现Elasticsearch插件\`
- \`feat: 实现Milvus向量数据库插件\`

**记录文档**: 
- \`docs/tech/step/03_MySQL数据库集成.md\`
- \`docs/tech/step/04_PostgreSQL及扩展集成.md\`
- \`docs/tech/step/05_NoSQL数据库集成.md\`
- \`docs/tech/step/06_向量数据库集成.md\`

---

## 阶段 4：ORM 集成与模型定义（步骤3）

### 任务清单
- [ ] 配置 GORM（已在阶段3部分完成）
- [ ] 定义基础模型结构
- [ ] 实现软删除
- [ ] 实现审计字段（创建时间、更新时间、创建人等）
- [ ] 数据库迁移工具

### 基础模型定义
\`\`\`go
// internal/model/base.go
type BaseModel struct {
    ID        uint           \`gorm:\"primarykey\" json:\"id\"\`
    CreatedAt time.Time      \`json:\"createdAt\"\`
    UpdatedAt time.Time      \`json:\"updatedAt\"\`
    DeletedAt gorm.DeletedAt \`gorm:\"index\" json:\"-\"\`
}

type AuditModel struct {
    BaseModel
    CreateBy uint   \`json:\"createBy\"\`
    UpdateBy uint   \`json:\"updateBy\"\`
    Remark   string \`json:\"remark\"\`
}
\`\`\`

### 参考 ruoyi-vue-pro 的表设计
需要实现的核心表：
- system_user（用户表）
- system_role（角色表）
- system_menu（菜单表）
- system_dept（部门表）
- system_post（岗位表）
- system_user_role（用户角色关联）
- system_role_menu（角色菜单关联）

### 验收标准
- GORM 配置完成
- 基础模型定义完成
- 可以执行数据库迁移
- 有表结构文档

**完成后提交**: \`feat: 实现GORM集成及基础模型定义\`

**记录文档**: \`docs/tech/step/07_ORM集成与模型设计.md\`

---

## 阶段 5：Gin 框架配置（步骤4）

### 任务清单
- [ ] 初始化 Gin 引擎
- [ ] 配置路由分组
- [ ] 设置日志
- [ ] 配置跨域
- [ ] 优雅关闭

### 实现要求
\`\`\`go
// internal/router/router.go
func InitRouter() *gin.Engine {
    r := gin.New()
    r.Use(gin.Logger())
    r.Use(gin.Recovery())
    r.Use(Cors()) // 跨域中间件
    
    // API 分组
    api := r.Group(\"/api\")
    {
        v1 := api.Group(\"/v1\")
        {
            // 系统管理
            system := v1.Group(\"/system\")
            {
                // 用户管理
                system.GET(\"/user/page\", handler.GetUserPage)
                // ... 其他路由
            }
        }
    }
    
    return r
}
\`\`\`

### 验收标准
- Gin 服务可以启动
- 路由分组清晰
- 有健康检查接口
- 日志格式规范

**完成后提交**: \`feat: 配置Gin框架及路由\`

**记录文档**: \`docs/tech/step/08_Gin框架配置.md\`

---

## 阶段 6：拦截器/中间件配置（步骤5）

### 任务清单
- [ ] 日志中间件
- [ ] CORS 中间件
- [ ] 限流中间件
- [ ] 认证中间件（JWT）
- [ ] 权限中间件
- [ ] 操作日志中间件
- [ ] 异常处理中间件

### 中间件目录
\`\`\`
internal/middleware/
├── logger.go           # 日志中间件
├── cors.go             # 跨域
├── rate_limit.go       # 限流
├── auth.go             # JWT认证
├── permission.go       # 权限校验
├── operation_log.go    # 操作日志
└── recovery.go         # 异常恢复
\`\`\`

### 验收标准
- 所有中间件实现完成
- 中间件可以正确执行
- 有中间件使用示例

**分多次提交**:
- \`feat: 实现日志和CORS中间件\`
- \`feat: 实现限流和异常处理中间件\`
- \`feat: 实现认证和权限中间件\`
- \`feat: 实现操作日志中间件\`

**记录文档**: \`docs/tech/step/09_中间件实现.md\`

---

## 阶段 7：JWT 功能实现（步骤6）

### 任务清单
- [ ] 生成 JWT Token
- [ ] 验证 Token
- [ ] 刷新 Token
- [ ] Token 存储（Redis）
- [ ] 单点登录支持

### 实现要求
\`\`\`go
// internal/pkg/jwt/jwt.go
type Claims struct {
    UserID   uint   \`json:\"userId\"\`
    Username string \`json:\"username\"\`
    jwt.StandardClaims
}

func GenerateToken(userId uint, username string) (string, error)
func ParseToken(tokenString string) (*Claims, error)
func RefreshToken(tokenString string) (string, error)
\`\`\`

### 验收标准
- 可以生成和验证 Token
- Token 过期自动刷新
- 支持黑名单（用户登出）

**完成后提交**: \`feat: 实现JWT认证功能\`

**记录文档**: \`docs/tech/step/10_JWT认证实现.md\`

---

## 阶段 8：用户管理模块实现（参考 ruoyi-vue-pro）

### 任务清单
- [ ] 创建用户表和模型
- [ ] 实现用户 CRUD
- [ ] 用户登录/登出
- [ ] 用户信息查询
- [ ] 用户密码加密
- [ ] 用户状态管理

### 接口规范（完全参考 ruoyi-vue-pro）

#### 1. 用户分页查询
\`\`\`
GET /api/v1/system/user/page
参数: pageNo, pageSize, username, mobile, status, createTime
响应格式:
{
  \"code\": 0,
  \"data\": {
    \"list\": [...],
    \"total\": 100
  },
  \"msg\": \"操作成功\"
}
\`\`\`

#### 2. 用户详情
\`\`\`
GET /api/v1/system/user/get?id=1
\`\`\`

#### 3. 创建用户
\`\`\`
POST /api/v1/system/user/create
Body: {\"username\": \"admin\", \"password\": \"123456\", ...}
\`\`\`

#### 4. 更新用户
\`\`\`
PUT /api/v1/system/user/update
\`\`\`

#### 5. 删除用户
\`\`\`
DELETE /api/v1/system/user/delete?id=1
\`\`\`

#### 6. 用户登录
\`\`\`
POST /api/v1/system/auth/login
Body: {\"username\": \"admin\", \"password\": \"123456\"}
响应: {\"code\": 0, \"data\": {\"token\": \"xxx\", \"expiresTime\": 1234567890}}
\`\`\`

### 文件结构
\`\`\`
internal/
├── api/v1/system/
│   └── user.go           # 用户控制器
├── service/system/
│   └── user.go           # 用户服务
├── dao/system/
│   └── user.go           # 用户数据访问
└── model/system/
    └── user.go           # 用户模型
\`\`\`

### 验收标准
- 所有接口地址与 ruoyi-vue-pro 一致
- 响应格式与 ruoyi-vue-pro 一致
- 密码使用 bcrypt 加密
- 有接口测试用例

**完成后提交**: \`feat: 实现用户管理模块\`

**记录文档**: \`docs/tech/step/11_用户管理实现.md\`
- 详细说明与 ruoyi-vue-pro 的对应关系
- 接口列表
- 实现思路

---

## 阶段 9：角色管理模块实现

### 任务清单
- [ ] 创建角色表和模型
- [ ] 实现角色 CRUD
- [ ] 角色菜单权限分配
- [ ] 角色数据权限配置

### 接口规范（参考 ruoyi-vue-pro）
\`\`\`
GET    /api/v1/system/role/page
GET    /api/v1/system/role/get?id=1
POST   /api/v1/system/role/create
PUT    /api/v1/system/role/update
DELETE /api/v1/system/role/delete?id=1
GET    /api/v1/system/role/list-all-simple  # 获取角色精简列表
\`\`\`

### 验收标准
- 接口与 ruoyi-vue-pro 一致
- 支持角色菜单权限配置
- 支持数据权限范围配置

**完成后提交**: \`feat: 实现角色管理模块\`

**记录文档**: \`docs/tech/step/12_角色管理实现.md\`

---

## 阶段 10：菜单管理模块实现

### 任务清单
- [ ] 创建菜单表（树形结构）
- [ ] 实现菜单 CRUD
- [ ] 菜单树形结构返回
- [ ] 用户菜单权限查询

### 接口规范
\`\`\`
GET    /api/v1/system/menu/list
GET    /api/v1/system/menu/get?id=1
POST   /api/v1/system/menu/create
PUT    /api/v1/system/menu/update
DELETE /api/v1/system/menu/delete?id=1
GET    /api/v1/system/permission/list-user-permissions  # 获取用户菜单
\`\`\`

### 验收标准
- 支持树形菜单结构
- 菜单权限控制正确
- 接口格式与 ruoyi 一致

**完成后提交**: \`feat: 实现菜单管理模块\`

**记录文档**: \`docs/tech/step/13_菜单管理实现.md\`

---

## 阶段 11：部门管理模块实现

### 任务清单
- [ ] 创建部门表（树形结构）
- [ ] 实现部门 CRUD
- [ ] 部门树形结构

### 接口规范
\`\`\`
GET    /api/v1/system/dept/list
GET    /api/v1/system/dept/get?id=1
POST   /api/v1/system/dept/create
PUT    /api/v1/system/dept/update
DELETE /api/v1/system/dept/delete?id=1
\`\`\`

**完成后提交**: \`feat: 实现部门管理模块\`

**记录文档**: \`docs/tech/step/14_部门管理实现.md\`

---

## 阶段 12：文件上传功能实现

### 任务清单
- [ ] 实现本地文件上传
- [ ] 实现 MinIO/OSS 文件上传
- [ ] 文件类型验证
- [ ] 文件大小限制
- [ ] 文件访问URL生成
- [ ] 创建 plugin/oss/

### Plugin 目录
\`\`\`
plugin/oss/
├── oss.go              # OSS 接口定义
├── local.go            # 本地存储
├── minio.go            # MinIO
├── aliyun.go           # 阿里云OSS（可选）
└── README.md
\`\`\`

### 接口规范
\`\`\`
POST /api/v1/infra/file/upload
响应: {\"code\": 0, \"data\": {\"url\": \"http://...\", \"fileName\": \"xxx\"}}
\`\`\`

### 验收标准
- 支持多种存储方式
- 文件上传成功
- 可以访问上传的文件

**完成后提交**: \`feat: 实现文件上传功能及OSS插件\`

**记录文档**: \`docs/tech/step/15_文件上传实现.md\`

---

## 阶段 13：Kafka 消息队列集成

### 任务清单
- [ ] 实现 Kafka 生产者
- [ ] 实现 Kafka 消费者
- [ ] 创建 plugin/kafka/
- [ ] 实现消息发送和接收示例

### Plugin 目录
\`\`\`
plugin/kafka/
├── producer.go
├── consumer.go
├── config.go
└── README.md
\`\`\`

**完成后提交**: \`feat: 实现Kafka消息队列插件\`

**记录文档**: \`docs/tech/step/16_Kafka集成.md\`

---

## 阶段 14：定时任务实现

### 任务清单
- [ ] 集成 Cron 库
- [ ] 实现任务调度
- [ ] 任务管理接口
- [ ] 创建 plugin/cron/

### Plugin 目录
\`\`\`
plugin/cron/
├── cron.go
├── job.go
└── README.md
\`\`\`

**完成后提交**: \`feat: 实现定时任务调度功能\`

**记录文档**: \`docs/tech/step/17_定时任务实现.md\`

---

## 阶段 15：AI 集成（参考 xiaozhi-server）

### 任务清单
- [ ] 研究 xiaozhi-server 的 AI 集成方案
- [ ] 实现 AI 对话接口
- [ ] 支持多 AI 模型切换（OpenAI、Claude等）
- [ ] 实现对话历史管理
- [ ] 创建 plugin/ai/

### Plugin 目录
\`\`\`
plugin/ai/
├── ai.go               # AI 接口定义
├── openai.go           # OpenAI 实现
├── claude.go           # Claude 实现
├── conversation.go     # 对话管理
└── README.md
\`\`\`

### 接口规范
\`\`\`
POST /api/v1/ai/chat
Body: {\"message\": \"你好\", \"conversationId\": \"xxx\"}
响应: {\"code\": 0, \"data\": {\"reply\": \"你好！\", \"conversationId\": \"xxx\"}}
\`\`\`

### 验收标准
- 可以调用 AI 接口
- 支持对话上下文
- 有错误处理

**完成后提交**: \`feat: 实现AI对话集成\`

**记录文档**: \`docs/tech/step/18_AI集成实现.md\`
- 参考 xiaozhi-server 的方案
- AI 模型配置
- 对话流程设计

---

## 阶段 16：Ruoyi Starter 功能移植

### 任务清单
分析 ruoyi-vue-pro 的所有的 starter 模块，将以下功能实现到 plugin/ 中：

- [ ] spring-boot-starter-biz-data-permission（数据权限）
- [ ] spring-boot-starter-biz-dict（字典管理）
- [ ] spring-boot-starter-biz-operatelog（操作日志）
- [ ] spring-boot-starter-biz-error-code（错误码）
- [ ] spring-boot-starter-security（安全）
等

### 需要创建的 Plugin
\`\`\`
plugin/
├── dataperm/          # 数据权限
├── dict/              # 字典管理
├── operlog/           # 操作日志
├── errorcode/         # 错误码
└── security/          # 安全（已在JWT中部分实现）
\`\`\`

### 每个功能独立提交
- \`feat: 实现数据权限插件\`
- \`feat: 实现字典管理插件\`
- \`feat: 实现操作日志插件\`
- \`feat: 实现错误码管理插件\`

**记录文档**: 每个插件一个文档，如：
- \`docs/tech/step/19_数据权限实现.md\`
- \`docs/tech/step/20_字典管理实现.md\`
- 等等...

---

## 阶段 17：Docker 部署配置

### 任务清单
- [ ] 编写 Dockerfile
- [ ] 编写 docker-compose.yaml（包含所有服务）
- [ ] 配置数据库初始化脚本
- [ ] 编写部署文档

### docker-compose.yaml 需要包含
\`\`\`yaml
services:
  mysql:
    image: mysql:8.0
  postgresql:
    image: postgis/postgis:15-3.3
  mongodb:
    image: mongo:6
  redis:
    image: redis:7
  elasticsearch:
    image: elasticsearch:8.8.0
  milvus:
    image: milvusdb/milvus:latest
  kafka:
    image: confluentinc/cp-kafka:latest
  zookeeper:
    image: confluentinc/cp-zookeeper:latest
  app:
    build: .
\`\`\`

**完成后提交**: \`feat: 添加Docker部署配置\`

**记录文档**: \`docs/tech/step/21_Docker部署.md\`

---

## 阶段 18：集成测试

### 任务清单
- [ ] 编写接口测试用例
- [ ] 测试所有核心功能
- [ ] 性能测试
- [ ] 压力测试

**完成后提交**: \`test: 添加集成测试用例\`

**记录文档**: \`docs/tech/step/22_集成测试.md\`

---

## 阶段 19：文档完善

### 任务清单
- [ ] 完善 API 文档（Swagger）
- [ ] 编写部署文档
- [ ] 编写开发文档
- [ ] 编写运维文档

### 产出文档
- \`docs/api/swagger.yaml\` - API 文档
- \`docs/deploy.md\` - 部署指南
- \`docs/development.md\` - 开发指南
- \`docs/operation.md\` - 运维手册

**完成后提交**: \`docs: 完善项目文档\`

---

## 阶段 20：最终验收

### 任务清单
- [ ] 验证所有功能可用
- [ ] 检查所有接口与 ruoyi-vue-pro 的一致性
- [ ] 验证 Docker 部署
- [ ] 代码质量检查
- [ ] 生成最终项目报告

### 验收标准
- 所有核心功能实现
- 接口规范与 ruoyi-vue-pro 一致
- 所有数据库正常运行
- Docker 可以一键部署
- 文档齐全

**只有在所有验收通过后，才输出: <promise>GIN_ADMIN_PROJECT_COMPLETE</promise>**

**最终文档**: \`docs/tech/23_项目总结.md\`

---

## 执行规则

1. **严格按阶段顺序执行**: 1 → 2 → 3 → ... → 20
2. **每完成一个小功能就提交**:
   - 不要等整个阶段完成才提交
   - 每个提交对应一个具体功能
3. **提交信息规范**:
   - 使用约定的格式
   - **禁止**包含 AI、GPT、Claude 等字样
4. **文档同步更新**:
   - 每完成一个功能，立即在 \`docs/tech/step/\` 写文档
   - 文档要详细记录实现思路、代码示例、遇到的问题
5. **接口完全对齐**:
   - 所有接口地址必须与 ruoyi-vue-pro 一致
   - 响应格式必须一致
6. **Plugin 独立性**:
   - 每个 plugin 可以独立使用
   - 有独立的 README
7. **代码质量**:
   - 有适当的注释
   - 有错误处理
   - 有单元测试

---

## 进度跟踪

创建 \`PROGRESS.md\`，实时更新：

\`\`\`markdown
# 开发进度

## 已完成阶段
- [x] 阶段 1: 项目初始化 ✅
- [ ] 阶段 2: 配置文件加载
- [ ] 阶段 3: 数据库集成
...

## 当前状态
- 正在执行: 阶段 X
- 已提交次数: X 次
- 已完成功能: XXX

## Plugin 完成情况
- [x] plugin/mysql ✅
- [ ] plugin/postgresql
- [ ] plugin/mongodb
...

## 核心模块完成情况
- [ ] 用户管理
- [ ] 角色管理
- [ ] 菜单管理
...
\`\`\`

---

**现在开始执行阶段 1：项目初始化与技术调研**

请先：
1. 研究 gin-vue-admin 项目结构
2. 分析 ruoyi-vue-pro 接口规范
3. 查看 xiaozhi-server AI 集成方案
4. 生成技术文档

" --max-iterations 50 --completion-promise "GIN_ADMIN_PROJECT_COMPLETE"
