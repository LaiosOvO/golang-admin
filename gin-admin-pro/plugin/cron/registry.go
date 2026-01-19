package cron

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// DefaultRegistry 默认任务注册器
type DefaultRegistry struct {
	manager  *CronManager
	handlers map[string]JobHandler
	mu       sync.RWMutex
}

// NewDefaultRegistry 创建默认任务注册器
func NewDefaultRegistry(manager *CronManager) *DefaultRegistry {
	return &DefaultRegistry{
		manager:  manager,
		handlers: make(map[string]JobHandler),
	}
}

// RegisterHandler 注册任务处理器
func (r *DefaultRegistry) RegisterHandler(name string, handler JobHandler) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.handlers[name]; exists {
		return fmt.Errorf("handler %s already exists", name)
	}

	r.handlers[name] = handler
	log.Printf("Handler %s registered successfully", name)
	return nil
}

// GetHandler 获取任务处理器
func (r *DefaultRegistry) GetHandler(name string) (JobHandler, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	handler, exists := r.handlers[name]
	if !exists {
		return nil, fmt.Errorf("handler %s not found", name)
	}

	return handler, nil
}

// ListHandlers 列出所有处理器
func (r *DefaultRegistry) ListHandlers() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	handlers := make([]string, 0, len(r.handlers))
	for name := range r.handlers {
		handlers = append(handlers, name)
	}

	return handlers
}

// UnregisterHandler 注销任务处理器
func (r *DefaultRegistry) UnregisterHandler(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.handlers[name]; !exists {
		return fmt.Errorf("handler %s not found", name)
	}

	delete(r.handlers, name)
	log.Printf("Handler %s unregistered successfully", name)
	return nil
}

// DefaultExecutor 默认任务执行器
type DefaultExecutor struct {
	manager *CronManager
	history map[string][]*JobExecution
	mu      sync.RWMutex
}

// NewDefaultExecutor 创建默认任务执行器
func NewDefaultExecutor(manager *CronManager) *DefaultExecutor {
	return &DefaultExecutor{
		manager: manager,
		history: make(map[string][]*JobExecution),
	}
}

// Execute 执行任务
func (e *DefaultExecutor) Execute(ctx context.Context, job Job) error {
	config := job.GetConfig()

	// 获取任务处理器
	handler, err := e.manager.GetRegistry().GetHandler(config.Handler)
	if err != nil {
		return fmt.Errorf("get handler failed: %w", err)
	}

	// 创建执行记录
	execution := &JobExecution{
		ID:          fmt.Sprintf("%s-%d", job.GetID(), time.Now().Unix()),
		JobID:       job.GetID(),
		JobName:     job.GetName(),
		Status:      JobStatusRunning,
		StartTime:   time.Now(),
		Attempt:     1,
		MaxAttempts: config.MaxRetries + 1,
		Params:      config.Params,
		Trigger:     "schedule",
		Metadata:    make(map[string]string),
	}

	// 记录执行历史
	e.recordExecution(execution)

	// 调用任务开始回调
	if err := handler.OnStart(ctx, job); err != nil {
		log.Printf("Job %s start callback failed: %v", job.GetID(), err)
	}

	// 执行任务，包含重试逻辑
	var lastErr error
	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		execution.Attempt = attempt + 1

		if attempt > 0 {
			// 等待重试延迟
			select {
			case <-time.After(config.RetryDelay):
			case <-ctx.Done():
				return ctx.Err()
			}

			// 调用重试回调
			if err := handler.OnRetry(ctx, job, attempt, lastErr); err != nil {
				log.Printf("Job %s retry callback failed: %v", job.GetID(), err)
			}
		}

		// 执行任务
		err := handler.Handle(ctx, job, config.Params)
		if err == nil {
			// 任务成功
			execution.Status = JobStatusCompleted
			now := time.Now()
			execution.EndTime = &now
			execution.Duration = now.Sub(execution.StartTime)

			e.recordExecution(execution)

			// 调用完成回调
			if err := handler.OnComplete(ctx, job, nil); err != nil {
				log.Printf("Job %s complete callback failed: %v", job.GetID(), err)
			}

			return nil
		}

		lastErr = err
		log.Printf("Job %s attempt %d failed: %v", job.GetID(), attempt+1, err)
	}

	// 任务失败
	execution.Status = JobStatusFailed
	execution.Error = lastErr.Error()
	now := time.Now()
	execution.EndTime = &now
	execution.Duration = now.Sub(execution.StartTime)

	e.recordExecution(execution)

	// 调用错误回调
	if err := handler.OnError(ctx, job, lastErr); err != nil {
		log.Printf("Job %s error callback failed: %v", job.GetID(), err)
	}

	return lastErr
}

// GetExecutionHistory 获取执行历史
func (e *DefaultExecutor) GetExecutionHistory(jobID string, limit int) ([]*JobExecution, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	history, exists := e.history[jobID]
	if !exists {
		return []*JobExecution{}, nil
	}

	if limit > 0 && len(history) > limit {
		return history[len(history)-limit:], nil
	}

	return history, nil
}

// GetRunningJobs 获取正在执行的任务
func (e *DefaultExecutor) GetRunningJobs() []Job {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var runningJobs []Job
	for _, executions := range e.history {
		for _, execution := range executions {
			if execution.Status == JobStatusRunning {
				// 这里应该从管理器获取对应的Job实例
				// 简化实现，返回空列表
				continue
			}
		}
	}

	return runningJobs
}

// CancelExecution 取消任务执行
func (e *DefaultExecutor) CancelExecution(jobID string) error {
	// 简化实现，实际应该通过context取消
	log.Printf("Cancel execution for job %s", jobID)
	return nil
}

// recordExecution 记录执行历史
func (e *DefaultExecutor) recordExecution(execution *JobExecution) {
	e.mu.Lock()
	defer e.mu.Unlock()

	jobID := execution.JobID
	if _, exists := e.history[jobID]; !exists {
		e.history[jobID] = []*JobExecution{}
	}

	e.history[jobID] = append(e.history[jobID], execution)

	// 限制历史记录数量
	maxHistory := 100
	if len(e.history[jobID]) > maxHistory {
		e.history[jobID] = e.history[jobID][len(e.history[jobID])-maxHistory:]
	}
}

// DefaultJob 默认任务实现
type DefaultJob struct {
	*BaseJob
}

// NewDefaultJob 创建默认任务
func NewDefaultJob(config *JobConfig) *DefaultJob {
	return &DefaultJob{
		BaseJob: NewBaseJob(config),
	}
}

// Execute 执行任务
func (j *DefaultJob) Execute(ctx context.Context, params map[string]interface{}) error {
	// 默认实现什么都不做
	return nil
}

// Cancel 取消任务
func (j *DefaultJob) Cancel() error {
	j.setStatus(JobStatusCancelled)
	return nil
}

// Pause 暂停任务
func (j *DefaultJob) Pause() error {
	j.setStatus(JobStatusPaused)
	return nil
}

// Resume 恢复任务
func (j *DefaultJob) Resume() error {
	if j.GetStatus() == JobStatusPaused {
		j.setStatus(JobStatusPending)
	}
	return nil
}

// CheckDependencies 检查依赖
func (j *DefaultJob) CheckDependencies() bool {
	// 简化实现，返回true
	return true
}

// SampleJobHandler 示例任务处理器
type SampleJobHandler struct{}

// Handle 执行任务
func (h *SampleJobHandler) Handle(ctx context.Context, job Job, params map[string]interface{}) error {
	log.Printf("Sample job %s is running with params: %v", job.GetID(), params)

	// 模拟任务执行
	select {
	case <-time.After(time.Second * 2):
		log.Printf("Sample job %s completed", job.GetID())
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// OnStart 任务开始
func (h *SampleJobHandler) OnStart(ctx context.Context, job Job) error {
	log.Printf("Sample job %s started", job.GetID())
	return nil
}

// OnComplete 任务完成
func (h *SampleJobHandler) OnComplete(ctx context.Context, job Job, err error) error {
	if err != nil {
		log.Printf("Sample job %s completed with error: %v", job.GetID(), err)
	} else {
		log.Printf("Sample job %s completed successfully", job.GetID())
	}
	return nil
}

// OnError 任务失败
func (h *SampleJobHandler) OnError(ctx context.Context, job Job, err error) error {
	log.Printf("Sample job %s failed: %v", job.GetID(), err)
	return nil
}

// OnRetry 任务重试
func (h *SampleJobHandler) OnRetry(ctx context.Context, job Job, attempt int, err error) error {
	log.Printf("Sample job %s retry attempt %d, last error: %v", job.GetID(), attempt, err)
	return nil
}

// CleanupJobHandler 清理任务处理器
type CleanupJobHandler struct{}

// Handle 执行任务
func (h *CleanupJobHandler) Handle(ctx context.Context, job Job, params map[string]interface{}) error {
	log.Printf("Cleanup job %s is running", job.GetID())

	// 模拟清理操作
	select {
	case <-time.After(time.Second * 1):
		log.Printf("Cleanup job %s completed", job.GetID())
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// OnStart 任务开始
func (h *CleanupJobHandler) OnStart(ctx context.Context, job Job) error {
	log.Printf("Cleanup job %s started", job.GetID())
	return nil
}

// OnComplete 任务完成
func (h *CleanupJobHandler) OnComplete(ctx context.Context, job Job, err error) error {
	if err != nil {
		log.Printf("Cleanup job %s completed with error: %v", job.GetID(), err)
	} else {
		log.Printf("Cleanup job %s completed successfully", job.GetID())
	}
	return nil
}

// OnError 任务失败
func (h *CleanupJobHandler) OnError(ctx context.Context, job Job, err error) error {
	log.Printf("Cleanup job %s failed: %v", job.GetID(), err)
	return nil
}

// OnRetry 任务重试
func (h *CleanupJobHandler) OnRetry(ctx context.Context, job Job, attempt int, err error) error {
	log.Printf("Cleanup job %s retry attempt %d, last error: %v", job.GetID(), attempt, err)
	return nil
}

// ReportJobHandler 报告任务处理器
type ReportJobHandler struct{}

// Handle 执行任务
func (h *ReportJobHandler) Handle(ctx context.Context, job Job, params map[string]interface{}) error {
	log.Printf("Report job %s is running", job.GetID())

	// 模拟报告生成
	select {
	case <-time.After(time.Second * 3):
		log.Printf("Report job %s completed", job.GetID())
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// OnStart 任务开始
func (h *ReportJobHandler) OnStart(ctx context.Context, job Job) error {
	log.Printf("Report job %s started", job.GetID())
	return nil
}

// OnComplete 任务完成
func (h *ReportJobHandler) OnComplete(ctx context.Context, job Job, err error) error {
	if err != nil {
		log.Printf("Report job %s completed with error: %v", job.GetID(), err)
	} else {
		log.Printf("Report job %s completed successfully", job.GetID())
	}
	return nil
}

// OnError 任务失败
func (h *ReportJobHandler) OnError(ctx context.Context, job Job, err error) error {
	log.Printf("Report job %s failed: %v", job.GetID(), err)
	return nil
}

// OnRetry 任务重试
func (h *ReportJobHandler) OnRetry(ctx context.Context, job Job, attempt int, err error) error {
	log.Printf("Report job %s retry attempt %d, last error: %v", job.GetID(), attempt, err)
	return nil
}
