# Cron Plugin

Cron插件提供了强大的定时任务调度功能，支持多种任务类型、依赖管理、重试机制和事件监控。

## 功能特性

- ✅ 基于Cron表达式的任务调度
- ✅ 灵活的任务注册和管理
- ✅ 任务依赖关系管理
- ✅ 自动重试和错误处理
- ✅ 任务执行历史记录
- ✅ 事件监听和监控
- ✅ 并发控制和单例任务
- ✅ 动态任务管理
- ✅ 多时区支持
- ✅ 任务生命周期回调

## 快速开始

### 基础配置

```yaml
# config/config.yaml
cron:
  enabled: true
  timezone: "Asia/Shanghai"
  concurrency: 100
  timeout: 30m
  maxRetries: 3
  retryDelay: 5s
  logEnabled: true
  logLevel: "info"
  managementEnabled: true
  
  jobs:
    - id: "cleanup-temp-files"
      name: "清理临时文件"
      description: "每天凌晨清理临时文件"
      enabled: true
      cron: "0 2 * * *"  # 每天凌晨2点
      handler: "cleanup"
      timeout: 10m
      maxRetries: 2
      retryDelay: 5m
      tags: ["cleanup", "maintenance"]
      
    - id: "daily-report"
      name: "每日报告"
      description: "每天生成业务报告"
      enabled: true
      cron: "0 8 * * *"  # 每天早上8点
      handler: "report"
      timeout: 30m
      maxRetries: 3
      retryDelay: 10m
      dependsOn: ["cleanup-temp-files"]
      tags: ["report", "daily"]
```

### 基本使用

```go
package main

import (
    "gin-admin-pro/plugin/cron"
)

func main() {
    // 创建定时任务管理器
    config := cron.DefaultConfig()
    config.Timezone = "Asia/Shanghai"
    
    manager, err := cron.NewCronManager(config)
    if err != nil {
        panic(err)
    }
    
    // 初始化
    err = manager.Initialize()
    if err != nil {
        panic(err)
    }
    defer manager.Stop()
    
    // 启动管理器
    err = manager.Start()
    if err != nil {
        panic(err)
    }
    
    // 添加任务
    jobConfig := &cron.JobConfig{
        ID:          "my-task",
        Name:        "我的任务",
        Description: "这是一个示例任务",
        Enabled:     true,
        Cron:        "*/5 * * * *", // 每5分钟执行
        Handler:     "sample",
        Timeout:     time.Minute * 10,
        MaxRetries:  3,
        Tags:        []string{"example"},
    }
    
    err = manager.AddJob(jobConfig)
    if err != nil {
        panic(err)
    }
    
    // 保持程序运行
    select {}
}
```

## 任务开发

### 自定义任务处理器

```go
package main

import (
    "context"
    "gin-admin-pro/plugin/cron"
    "log"
)

type MyJobHandler struct{}

// Handle 执行任务
func (h *MyJobHandler) Handle(ctx context.Context, job cron.Job, params map[string]interface{}) error {
    log.Printf("Executing job: %s", job.GetID())
    
    // 业务逻辑
    err := doSomeWork(ctx)
    if err != nil {
        return err
    }
    
    log.Printf("Job %s completed successfully", job.GetID())
    return nil
}

// OnStart 任务开始回调
func (h *MyJobHandler) OnStart(ctx context.Context, job cron.Job) error {
    log.Printf("Job %s started", job.GetID())
    return nil
}

// OnComplete 任务完成回调
func (h *MyJobHandler) OnComplete(ctx context.Context, job cron.Job, err error) error {
    if err != nil {
        log.Printf("Job %s completed with error: %v", job.GetID(), err)
    } else {
        log.Printf("Job %s completed successfully", job.GetID())
    }
    return nil
}

// OnError 任务失败回调
func (h *MyJobHandler) OnError(ctx context.Context, job cron.Job, err error) error {
    log.Printf("Job %s failed: %v", job.GetID(), err)
    // 可以在这里添加告警逻辑
    return nil
}

// OnRetry 任务重试回调
func (h *MyJobHandler) OnRetry(ctx context.Context, job cron.Job, attempt int, err error) error {
    log.Printf("Job %s retry attempt %d, last error: %v", job.GetID(), attempt, err)
    return nil
}

func doSomeWork(ctx context.Context) error {
    // 模拟工作
    select {
    case <-time.After(time.Second * 2):
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

// 注册任务处理器
func main() {
    manager, _ := cron.NewCronManager(nil)
    
    // 注册自定义处理器
    myHandler := &MyJobHandler{}
    registry := manager.GetRegistry()
    registry.RegisterHandler("my-job", myHandler)
    
    // 使用自定义处理器
    jobConfig := &cron.JobConfig{
        ID:      "custom-job",
        Name:    "自定义任务",
        Enabled: true,
        Cron:    "0 */1 * * *", // 每小时执行
        Handler: "my-job",
    }
    
    manager.AddJob(jobConfig)
}
```

### 自定义任务实现

```go
package main

import (
    "context"
    "gin-admin-pro/plugin/cron"
)

type CustomJob struct {
    *cron.BaseJob
    
    // 自定义字段
    customField string
}

func NewCustomJob(config *cron.JobConfig, customField string) *CustomJob {
    return &CustomJob{
        BaseJob:     cron.NewBaseJob(config),
        customField: customField,
    }
}

// Execute 执行任务
func (j *CustomJob) Execute(ctx context.Context, params map[string]interface{}) error {
    // 自定义执行逻辑
    log.Printf("Custom job executing with field: %s", j.customField)
    
    // 使用传入的参数
    if param, exists := params["customParam"]; exists {
        log.Printf("Custom param: %v", param)
    }
    
    return nil
}

// Cancel 取消任务
func (j *CustomJob) Cancel() error {
    j.setStatus(cron.JobStatusCancelled)
    log.Printf("Custom job %s cancelled", j.GetID())
    return nil
}
```

## 任务管理

### 动态任务管理

```go
// 添加任务
jobConfig := &cron.JobConfig{
    ID:      "dynamic-job",
    Name:    "动态任务",
    Enabled: true,
    Cron:    "*/10 * * * *",
    Handler: "sample",
}

err := manager.AddJob(jobConfig)

// 手动执行任务
err := manager.RunJob("dynamic-job", map[string]interface{}{
    "manual": true,
    "user":   "admin",
})

// 暂停任务
err := manager.PauseJob("dynamic-job")

// 恢复任务
err := manager.ResumeJob("dynamic-job")

// 移除任务
err := manager.RemoveJob("dynamic-job")
```

### 任务查询和监控

```go
// 获取单个任务
job, err := manager.GetJob("my-job")
if err == nil {
    log.Printf("Job status: %s", job.GetStatus())
    log.Printf("Job enabled: %v", job.IsEnabled())
    log.Printf("Job running: %v", job.IsRunning())
}

// 列出所有任务
jobs := manager.ListJobs()
for _, job := range jobs {
    log.Printf("Job: %s, Status: %s, Enabled: %v", 
        job.GetID(), job.GetStatus(), job.IsEnabled())
}

// 获取执行历史
executor := manager.GetExecutor()
history, err := executor.GetExecutionHistory("my-job", 10)
for _, execution := range history {
    log.Printf("Execution %s: %s, Duration: %v, Error: %s", 
        execution.ID, execution.Status, execution.Duration, execution.Error)
}
```

## 事件监听

```go
type LogEventHandler struct{}

func (h *LogEventHandler) Handle(event *cron.JobEvent) error {
    log.Printf("Job event: %s, Job: %s, Time: %v", 
        event.Type, event.JobID, event.Timestamp)
    
    switch event.Type {
    case "started":
        // 任务开始事件
    case "completed":
        // 任务完成事件
    case "failed":
        // 任务失败事件
        if data, ok := event.Data.(map[string]interface{}); ok {
            if errMsg, exists := data["error"]; exists {
                log.Printf("Job failed with error: %v", errMsg)
            }
        }
    }
    
    return nil
}

// 添加事件监听器
eventHandler := &LogEventHandler{}
manager.AddEventHandler(eventHandler)
```

## Cron表达式

### 基本格式

```
分 时 日 月 周
*  *  *  *  *
```

### 特殊字符

- `*`: 任意值
- `,`: 多个值 (1,3,5)
- `-`: 范围 (1-5)
- `/`: 步长 (*/5)
- `?`: 不指定值 (仅用于日和周)

### 常用表达式

```
*/5 * * * *    # 每5分钟
0 */1 * * *    # 每小时
0 2 * * *      # 每天凌晨2点
0 2 * * 1      # 每周一凌晨2点
0 1 1 * *      # 每月1号凌晨1点
0 0 1 1 *      # 每年1月1号午夜
```

## 高级功能

### 任务依赖

```go
jobConfig := &cron.JobConfig{
    ID:        "backup-database",
    Name:      "数据库备份",
    Enabled:   true,
    Cron:      "0 3 * * *",  # 凌晨3点
    Handler:   "backup",
    DependsOn: ["cleanup-temp-files"],  # 依赖清理任务完成
    MaxRetries: 2,
}
```

### 单例任务

```go
jobConfig := &cron.JobConfig{
    ID:          "long-running-job",
    Name:        "长时间运行任务",
    Enabled:     true,
    Cron:        "0 */2 * * *",
    Handler:     "processor",
    Singleton:   true,           # 防止并发执行
    MaxInstances: 1,             # 最大实例数
    Timeout:     time.Hour,      # 超时时间
}
```

### 参数传递

```go
jobConfig := &cron.JobConfig{
    ID:      "param-job",
    Name:    "带参数的任务",
    Enabled: true,
    Cron:    "0 */1 * * *",
    Handler: "email",
    Params: map[string]interface{}{
        "to":       ["admin@example.com", "ops@example.com"],
        "subject":  "系统报告",
        "template": "daily-report",
        "retry":    true,
    },
}
```

### 标签和元数据

```go
jobConfig := &cron.JobConfig{
    ID:          "monitored-job",
    Name:        "监控任务",
    Enabled:     true,
    Cron:        "*/5 * * * *",
    Handler:     "monitor",
    Tags:        ["monitoring", "system", "critical"],
    Metadata: map[string]string{
        "owner":      "ops-team",
        "severity":   "high",
        "alert-when": "failed",
        "run-on":    "prod",
    },
}
```

## 配置参数详解

### 基础配置

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `enabled` | bool | true | 是否启用定时任务 |
| `timezone` | string | "UTC" | 时区设置 |
| `concurrency` | int | 100 | 最大并发数 |
| `timeout` | duration | 30m | 默认任务超时时间 |
| `maxRetries` | int | 3 | 默认重试次数 |
| `retryDelay` | duration | 5s | 重试间隔 |
| `logEnabled` | bool | true | 是否启用日志 |
| `logLevel` | string | "info" | 日志级别 |
| `managementEnabled` | bool | true | 是否启用管理接口 |

### 任务配置

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `id` | string | 必填 | 任务唯一标识 |
| `name` | string | 必填 | 任务名称 |
| `description` | string | "" | 任务描述 |
| `enabled` | bool | true | 是否启用 |
| `cron` | string | 必填 | Cron表达式 |
| `handler` | string | 必填 | 处理器名称 |
| `timeout` | duration | 30m | 任务超时时间 |
| `maxRetries` | int | 3 | 重试次数 |
| `retryDelay` | duration | 5s | 重试间隔 |
| `dependsOn` | []string | [] | 依赖任务ID列表 |
| `singleton` | bool | true | 是否为单例任务 |
| `maxInstances` | int | 1 | 最大并发实例数 |
| `params` | map | {} | 任务参数 |
| `tags` | []string | [] | 任务标签 |
| `metadata` | map | {} | 任务元数据 |

## 最佳实践

### 1. 任务设计

```go
// ✅ 好的设计
jobConfig := &cron.JobConfig{
    ID:          "user-cleanup",
    Name:        "用户数据清理",
    Description: "清理30天前的非活跃用户数据",
    Enabled:     true,
    Cron:        "0 2 * * *",  # 非高峰期执行
    Handler:     "cleanup",
    Timeout:     time.Hour,
    MaxRetries:  2,
    Singleton:   true,        # 防止重复执行
    Tags:        ["cleanup", "maintenance"],
    Metadata: map[string]string{
        "retention": "30d",
        "backup":    "true",
    },
}

// ❌ 避免的设计
jobConfig := &cron.JobConfig{
    ID:       "bad-job",
    Cron:      "* * * * *",  # 频率太高
    Timeout:   time.Second,   # 超时太短
    MaxRetries: 0,           # 不重试
}
```

### 2. 错误处理

```go
type RobustJobHandler struct{}

func (h *RobustJobHandler) Handle(ctx context.Context, job cron.Job, params map[string]interface{}) error {
    // 记录开始时间
    start := time.Now()
    log.Printf("Job %s started", job.GetID())
    
    defer func() {
        duration := time.Since(start)
        log.Printf("Job %s completed in %v", job.GetID(), duration)
        
        // 发送监控指标
        metrics.RecordJobDuration(job.GetID(), duration)
    }()
    
    // 执行业务逻辑
    err := h.doWork(ctx, params)
    if err != nil {
        // 记录详细错误
        log.Printf("Job %s failed: %v", job.GetID(), err)
        
        // 发送告警
        alert.SendJobFailureAlert(job.GetID(), err)
        
        return err
    }
    
    return nil
}

func (h *RobustJobHandler) OnError(ctx context.Context, job cron.Job, err error) error {
    // 增强错误处理
    if errors.Is(err, context.DeadlineExceeded) {
        log.Printf("Job %s timed out", job.GetID())
    }
    
    // 检查是否需要人工干预
    if h.requiresManualIntervention(err) {
        h.createIncident(job, err)
    }
    
    return nil
}
```

### 3. 监控和告警

```go
type MonitoringHandler struct{}

func (h *MonitoringHandler) Handle(event *cron.JobEvent) error {
    // 发送指标到监控系统
    metrics.Counter("job_events_total").WithLabelValues(
        event.Type, event.JobID,
    ).Inc()
    
    // 关键事件告警
    if event.Type == "failed" {
        h.sendAlert(event)
    }
    
    return nil
}

func (h *MonitoringHandler) sendAlert(event *cron.JobEvent) {
    // 检查告警规则
    if h.shouldAlert(event.JobID) {
        alert.Send(&alert.Alert{
            Level:   "critical",
            Title:   fmt.Sprintf("Job Failed: %s", event.JobName),
            Message: fmt.Sprintf("Job %s failed at %v", event.JobID, event.Timestamp),
            Tags:    []string{"cron", "job-failure"},
        })
    }
}
```

## 故障排除

### 常见问题

1. **任务不执行**
   - 检查Cron表达式是否正确
   - 确认任务已启用
   - 检查时区设置

2. **任务重复执行**
   - 设置Singleton=true
   - 检查MaxInstances配置

3. **任务超时**
   - 增加Timeout配置
   - 优化任务逻辑

4. **依赖任务不执行**
   - 检查依赖任务是否成功完成
   - 确认依赖任务ID正确

### 调试技巧

```go
// 启用详细日志
config := cron.DefaultConfig()
config.LogLevel = "debug"
config.LogEnabled = true

// 添加调试事件处理器
debugHandler := &DebugEventHandler{}
manager.AddEventHandler(debugHandler)

// 手动执行任务进行调试
err := manager.RunJob("problematic-job", nil)
```

## 性能优化

### 1. 并发控制

```yaml
cron:
  concurrency: 50  # 根据服务器配置调整
```

### 2. 任务分组

```go
// 将相关任务分组，避免同时执行大量IO密集型任务
jobConfig := &cron.JobConfig{
    ID:      "io-intensive-job",
    Cron:    "0 3 * * *",  # 统一在低峰期执行
    Handler: "data-processor",
}
```

### 3. 资源监控

```go
type ResourceAwareHandler struct{}

func (h *ResourceAwareHandler) Handle(ctx context.Context, job cron.Job, params map[string]interface{}) error {
    // 检查系统资源
    if h.isSystemBusy() {
        return fmt.Errorf("system too busy, deferring job")
    }
    
    return h.doWork(ctx, params)
}
```

## 部署建议

### 生产环境配置

```yaml
cron:
  enabled: true
  timezone: "UTC"  # 统一使用UTC
  concurrency: 20  # 保守的并发数
  timeout: 1h
  maxRetries: 3
  retryDelay: 10m
  logLevel: "info"
  
  jobs:
    # 核心业务任务
    - id: "data-sync"
      enabled: true
      cron: "0 */6 * * *"  # 每6小时同步
      handler: "sync"
      timeout: 2h
      maxRetries: 5
      
    # 维护任务
    - id: "cleanup"
      enabled: true
      cron: "0 2 * * *"  # 凌晨执行
      handler: "cleanup"
      timeout: 3h
```

### 监控告警

- 任务执行成功率 > 95%
- 平均执行时间 < 预期时间
- 失败任务立即告警
- 长时间运行任务告警

## 扩展开发

### 自定义执行器

```go
type CustomExecutor struct {
    // 自定义字段
}

func (e *CustomExecutor) Execute(ctx context.Context, job cron.Job) error {
    // 自定义执行逻辑
    return nil
}
```

### 集成外部系统

```go
type DatabaseJobHandler struct {
    db *sql.DB
}

func (h *DatabaseJobHandler) Handle(ctx context.Context, job cron.Job, params map[string]interface{}) error {
    // 使用数据库连接
    _, err := h.db.ExecContext(ctx, "CALL cleanup_procedure()")
    return err
}
```

## 总结

Cron插件提供了完整的定时任务解决方案：

1. **强大的调度**: 基于标准Cron表达式
2. **灵活的管理**: 动态添加、删除、暂停、恢复任务
3. **可靠的执行**: 重试机制、超时控制、错误处理
4. **完善的监控**: 事件系统、执行历史、指标收集
5. **易于扩展**: 插件化架构，支持自定义处理器和执行器

该插件为企业级应用提供了可靠、高性能的定时任务管理功能。