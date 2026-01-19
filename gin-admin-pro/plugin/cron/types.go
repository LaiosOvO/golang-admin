package cron

import (
	"context"
	"sync"
	"time"
)

// JobStatus 任务状态
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
	JobStatusCancelled JobStatus = "cancelled"
	JobStatusPaused    JobStatus = "paused"
)

// Job 任务接口
type Job interface {
	// 基础信息
	GetID() string
	GetName() string
	GetDescription() string
	GetConfig() *JobConfig

	// 状态管理
	GetStatus() JobStatus
	IsEnabled() bool
	IsRunning() bool

	// 执行控制
	Execute(ctx context.Context, params map[string]interface{}) error
	Cancel() error
	Pause() error
	Resume() error

	// 依赖管理
	GetDependsOn() []string
	CheckDependencies() bool

	// 元数据
	GetTags() []string
	GetMetadata() map[string]string
}

// JobHandler 任务处理器接口
type JobHandler interface {
	// 执行任务
	Handle(ctx context.Context, job Job, params map[string]interface{}) error

	// 任务开始
	OnStart(ctx context.Context, job Job) error

	// 任务完成
	OnComplete(ctx context.Context, job Job, err error) error

	// 任务失败
	OnError(ctx context.Context, job Job, err error) error

	// 任务重试
	OnRetry(ctx context.Context, job Job, attempt int, err error) error
}

// JobRegistry 任务注册器接口
type JobRegistry interface {
	// 注册任务处理器
	RegisterHandler(name string, handler JobHandler) error

	// 获取任务处理器
	GetHandler(name string) (JobHandler, error)

	// 列出所有处理器
	ListHandlers() []string

	// 注销任务处理器
	UnregisterHandler(name string) error
}

// JobExecutor 任务执行器接口
type JobExecutor interface {
	// 执行任务
	Execute(ctx context.Context, job Job) error

	// 获取执行历史
	GetExecutionHistory(jobID string, limit int) ([]*JobExecution, error)

	// 获取正在执行的任务
	GetRunningJobs() []Job

	// 取消任务执行
	CancelExecution(jobID string) error
}

// JobExecution 任务执行记录
type JobExecution struct {
	ID          string                 `json:"id"`
	JobID       string                 `json:"jobId"`
	JobName     string                 `json:"jobName"`
	Status      JobStatus              `json:"status"`
	StartTime   time.Time              `json:"startTime"`
	EndTime     *time.Time             `json:"endTime,omitempty"`
	Duration    time.Duration          `json:"duration"`
	Attempt     int                    `json:"attempt"`
	MaxAttempts int                    `json:"maxAttempts"`
	Params      map[string]interface{} `json:"params"`
	Result      map[string]interface{} `json:"result,omitempty"`
	Error       string                 `json:"error,omitempty"`
	Trigger     string                 `json:"trigger"` // schedule, manual, api
	TriggeredBy string                 `json:"triggeredBy,omitempty"`
	Metadata    map[string]string      `json:"metadata"`
}

// JobMetrics 任务指标
type JobMetrics struct {
	JobID           string        `json:"jobId"`
	JobName         string        `json:"jobName"`
	TotalExecutions int64         `json:"totalExecutions"`
	SuccessCount    int64         `json:"successCount"`
	FailureCount    int64         `json:"failureCount"`
	AverageDuration time.Duration `json:"averageDuration"`
	LastExecution   time.Time     `json:"lastExecution"`
	NextExecution   time.Time     `json:"nextExecution,omitempty"`
}

// JobEvent 任务事件
type JobEvent struct {
	Type      string      `json:"type"` // started, completed, failed, cancelled, paused, resumed
	JobID     string      `json:"jobId"`
	JobName   string      `json:"jobName"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data,omitempty"`
}

// EventHandler 事件处理器接口
type EventHandler interface {
	Handle(event *JobEvent) error
}

// BaseJob 基础任务实现
type BaseJob struct {
	config    *JobConfig
	status    JobStatus
	startTime time.Time
	endTime   *time.Time
	attempt   int
	mu        sync.RWMutex
}

// NewBaseJob 创建基础任务
func NewBaseJob(config *JobConfig) *BaseJob {
	return &BaseJob{
		config: config,
		status: JobStatusPending,
	}
}

// GetID 获取任务ID
func (j *BaseJob) GetID() string {
	return j.config.ID
}

// GetName 获取任务名称
func (j *BaseJob) GetName() string {
	return j.config.Name
}

// GetDescription 获取任务描述
func (j *BaseJob) GetDescription() string {
	return j.config.Description
}

// GetConfig 获取任务配置
func (j *BaseJob) GetConfig() *JobConfig {
	return j.config
}

// GetStatus 获取任务状态
func (j *BaseJob) GetStatus() JobStatus {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.status
}

// IsEnabled 检查任务是否启用
func (j *BaseJob) IsEnabled() bool {
	return j.config.Enabled
}

// IsRunning 检查任务是否正在运行
func (j *BaseJob) IsRunning() bool {
	return j.GetStatus() == JobStatusRunning
}

// GetDependsOn 获取依赖任务
func (j *BaseJob) GetDependsOn() []string {
	return j.config.DependsOn
}

// GetTags 获取标签
func (j *BaseJob) GetTags() []string {
	return j.config.Tags
}

// GetMetadata 获取元数据
func (j *BaseJob) GetMetadata() map[string]string {
	return j.config.Metadata
}

// setStatus 设置任务状态
func (j *BaseJob) setStatus(status JobStatus) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.status = status

	if status == JobStatusRunning {
		j.startTime = time.Now()
		j.endTime = nil
	} else if status == JobStatusCompleted || status == JobStatusFailed || status == JobStatusCancelled {
		now := time.Now()
		j.endTime = &now
	}
}

// setAttempt 设置执行次数
func (j *BaseJob) setAttempt(attempt int) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.attempt = attempt
}

// getAttempt 获取执行次数
func (j *BaseJob) getAttempt() int {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.attempt
}
