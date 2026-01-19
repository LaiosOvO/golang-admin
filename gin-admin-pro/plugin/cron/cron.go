package cron

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// CronManager 定时任务管理器
type CronManager struct {
	config   *Config
	cron     *cron.Cron
	registry JobRegistry
	executor JobExecutor

	jobs     map[string]Job
	handlers map[string]JobHandler
	mu       sync.RWMutex

	// 执行管理
	runningJobs map[string]context.CancelFunc
	runningMu   sync.RWMutex

	// 事件管理
	eventHandlers []EventHandler

	// 状态
	started bool
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

// NewCronManager 创建定时任务管理器
func NewCronManager(config *Config) (*CronManager, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// 创建cron实例
	c := cron.New(
		cron.WithLocation(time.UTC),
		cron.WithChain(cron.Recover(cron.DefaultLogger)),
	)

	ctx, cancel := context.WithCancel(context.Background())

	manager := &CronManager{
		config:        config,
		cron:          c,
		jobs:          make(map[string]Job),
		handlers:      make(map[string]JobHandler),
		runningJobs:   make(map[string]context.CancelFunc),
		eventHandlers: []EventHandler{},
		ctx:           ctx,
		cancel:        cancel,
	}

	// 创建默认的注册器和执行器
	manager.registry = NewDefaultRegistry(manager)
	manager.executor = NewDefaultExecutor(manager)

	return manager, nil
}

// Initialize 初始化定时任务管理器
func (cm *CronManager) Initialize() error {
	if !cm.config.Enabled {
		log.Println("Cron manager is disabled")
		return nil
	}

	cm.mu.Lock()
	defer cm.mu.Unlock()

	// 注册默认任务处理器
	if err := cm.RegisterDefaultHandlers(); err != nil {
		return fmt.Errorf("register default handlers failed: %w", err)
	}

	// 加载配置中的任务
	for _, jobConfig := range cm.config.Jobs {
		if err := cm.AddJobFromConfig(jobConfig); err != nil {
			log.Printf("Failed to add job %s: %v", jobConfig.ID, err)
			continue
		}
	}

	log.Printf("Cron manager initialized with %d jobs", len(cm.jobs))
	return nil
}

// RegisterDefaultHandlers 注册默认任务处理器
func (cm *CronManager) RegisterDefaultHandlers() error {
	// 注册示例任务处理器
	sampleHandler := &SampleJobHandler{}
	if err := cm.registry.RegisterHandler("sample", sampleHandler); err != nil {
		return err
	}

	// 注册清理任务处理器
	cleanupHandler := &CleanupJobHandler{}
	if err := cm.registry.RegisterHandler("cleanup", cleanupHandler); err != nil {
		return err
	}

	// 注册报告任务处理器
	reportHandler := &ReportJobHandler{}
	if err := cm.registry.RegisterHandler("report", reportHandler); err != nil {
		return err
	}

	return nil
}

// Start 启动定时任务管理器
func (cm *CronManager) Start() error {
	if !cm.config.Enabled {
		log.Println("Cron manager is disabled, not starting")
		return nil
	}

	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.started {
		return fmt.Errorf("cron manager is already started")
	}

	// 启动cron调度器
	cm.cron.Start()
	cm.started = true

	log.Println("Cron manager started successfully")
	return nil
}

// Stop 停止定时任务管理器
func (cm *CronManager) Stop() error {
	if !cm.config.Enabled {
		return nil
	}

	cm.mu.Lock()
	defer cm.mu.Unlock()

	if !cm.started {
		return nil
	}

	// 取消所有运行中的任务
	cm.runningMu.Lock()
	for jobID, cancel := range cm.runningJobs {
		cancel()
		delete(cm.runningJobs, jobID)
	}
	cm.runningMu.Unlock()

	// 停止cron调度器
	ctx := cm.cron.Stop()
	select {
	case <-ctx.Done():
		log.Println("Cron manager stopped gracefully")
	case <-time.After(time.Second * 10):
		log.Println("Cron manager stop timeout, forcing stop")
	}

	// 取消主上下文
	cm.cancel()

	cm.started = false
	return nil
}

// AddJobFromConfig 从配置添加任务
func (cm *CronManager) AddJobFromConfig(config *JobConfig) error {
	if config == nil {
		return fmt.Errorf("job config is nil")
	}

	if config.ID == "" {
		return fmt.Errorf("job ID is required")
	}

	if config.Handler == "" {
		return fmt.Errorf("job handler is required")
	}

	// 检查任务是否已存在
	if _, exists := cm.jobs[config.ID]; exists {
		return fmt.Errorf("job %s already exists", config.ID)
	}

	// 创建任务
	job := NewDefaultJob(config)

	// 添加到管理器
	cm.jobs[config.ID] = job

	// 如果启用了任务，添加到调度器
	if config.Enabled && cm.started {
		return cm.scheduleJob(job)
	}

	return nil
}

// scheduleJob 调度任务
func (cm *CronManager) scheduleJob(job Job) error {
	config := job.GetConfig()

	// 添加任务到cron
	entryID, err := cm.cron.AddFunc(config.Cron, func() {
		cm.executeJob(job)
	})

	if err != nil {
		return fmt.Errorf("failed to schedule job %s: %w", job.GetID(), err)
	}

	// 存储entry ID（可以在扩展中添加到任务元数据）
	if config.Metadata == nil {
		config.Metadata = make(map[string]string)
	}
	config.Metadata["cronEntryID"] = fmt.Sprintf("%d", entryID)

	log.Printf("Job %s scheduled with cron: %s", job.GetID(), config.Cron)
	return nil
}

// executeJob 执行任务
func (cm *CronManager) executeJob(job Job) {
	config := job.GetConfig()

	// 检查任务是否启用
	if !job.IsEnabled() {
		return
	}

	// 检查单例任务
	if config.Singleton && job.IsRunning() {
		log.Printf("Job %s is singleton and already running, skipping", job.GetID())
		return
	}

	// 检查依赖
	if !cm.checkDependencies(job) {
		log.Printf("Job %s dependencies not satisfied, skipping", job.GetID())
		return
	}

	// 创建执行上下文
	ctx, cancel := context.WithTimeout(cm.ctx, config.Timeout)

	// 记录运行中的任务
	cm.runningMu.Lock()
	cm.runningJobs[job.GetID()] = cancel
	cm.runningMu.Unlock()

	// 启动goroutine执行任务
	cm.wg.Add(1)
	go func() {
		defer cm.wg.Done()
		defer func() {
			cm.runningMu.Lock()
			delete(cm.runningJobs, job.GetID())
			cm.runningMu.Unlock()
		}()

		// 发送开始事件
		cm.emitEvent(&JobEvent{
			Type:      "started",
			JobID:     job.GetID(),
			JobName:   job.GetName(),
			Timestamp: time.Now(),
		})

		// 执行任务
		err := cm.executor.Execute(ctx, job)

		// 发送完成事件
		eventType := "completed"
		if err != nil {
			eventType = "failed"
		}

		cm.emitEvent(&JobEvent{
			Type:      eventType,
			JobID:     job.GetID(),
			JobName:   job.GetName(),
			Timestamp: time.Now(),
			Data:      map[string]interface{}{"error": err},
		})
	}()
}

// checkDependencies 检查任务依赖
func (cm *CronManager) checkDependencies(job Job) bool {
	dependsOn := job.GetDependsOn()
	if len(dependsOn) == 0 {
		return true
	}

	for _, depID := range dependsOn {
		if depJob, exists := cm.jobs[depID]; exists {
			if !depJob.IsEnabled() || depJob.GetStatus() != JobStatusCompleted {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

// AddJob 动态添加任务
func (cm *CronManager) AddJob(config *JobConfig) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if err := cm.AddJobFromConfig(config); err != nil {
		return err
	}

	// 如果管理器已启动且任务启用，立即调度
	if cm.started && config.Enabled {
		job := cm.jobs[config.ID]
		return cm.scheduleJob(job)
	}

	return nil
}

// RemoveJob 移除任务
func (cm *CronManager) RemoveJob(jobID string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	job, exists := cm.jobs[jobID]
	if !exists {
		return fmt.Errorf("job %s not found", jobID)
	}

	// 取消正在运行的任务
	cm.runningMu.Lock()
	if cancel, exists := cm.runningJobs[jobID]; exists {
		cancel()
		delete(cm.runningJobs, jobID)
	}
	cm.runningMu.Unlock()

	// 从cron中移除（需要entry ID）
	if cronEntryID, exists := job.GetMetadata()["cronEntryID"]; exists {
		if entryID, err := fmt.Sscanf(cronEntryID, "%d", new(int64)); err == nil && entryID == 1 {
			// 从cron中移除任务（cron v3不支持直接移除特定任务，需要重建）
			// 这里采用重新创建cron的方式
			cm.rebuildCron()
		}
	}

	delete(cm.jobs, jobID)
	log.Printf("Job %s removed", jobID)

	return nil
}

// rebuildCron 重建cron调度器
func (cm *CronManager) rebuildCron() {
	if !cm.started {
		return
	}

	// 停止当前cron
	ctx := cm.cron.Stop()
	<-ctx.Done()

	// 重新创建cron
	cm.cron = cron.New(
		cron.WithLocation(time.UTC),
		cron.WithChain(cron.Recover(cron.DefaultLogger)),
	)

	// 重新调度所有启用的任务
	for _, job := range cm.jobs {
		if job.IsEnabled() {
			if err := cm.scheduleJob(job); err != nil {
				log.Printf("Failed to reschedule job %s: %v", job.GetID(), err)
			}
		}
	}

	// 重新启动
	cm.cron.Start()
}

// GetJob 获取任务
func (cm *CronManager) GetJob(jobID string) (Job, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	job, exists := cm.jobs[jobID]
	if !exists {
		return nil, fmt.Errorf("job %s not found", jobID)
	}

	return job, nil
}

// ListJobs 列出所有任务
func (cm *CronManager) ListJobs() []Job {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	jobs := make([]Job, 0, len(cm.jobs))
	for _, job := range cm.jobs {
		jobs = append(jobs, job)
	}

	return jobs
}

// RunJob 手动运行任务
func (cm *CronManager) RunJob(jobID string, params map[string]interface{}) error {
	job, err := cm.GetJob(jobID)
	if err != nil {
		return err
	}

	if !job.IsEnabled() {
		return fmt.Errorf("job %s is disabled", jobID)
	}

	// 直接执行任务
	go cm.executeJob(job)

	return nil
}

// PauseJob 暂停任务
func (cm *CronManager) PauseJob(jobID string) error {
	job, err := cm.GetJob(jobID)
	if err != nil {
		return err
	}

	// 取消正在运行的任务
	cm.runningMu.Lock()
	if cancel, exists := cm.runningJobs[jobID]; exists {
		cancel()
		delete(cm.runningJobs, jobID)
	}
	cm.runningMu.Unlock()

	// 设置任务为禁用状态
	config := job.GetConfig()
	config.Enabled = false

	// 重建调度器
	cm.rebuildCron()

	return nil
}

// ResumeJob 恢复任务
func (cm *CronManager) ResumeJob(jobID string) error {
	job, err := cm.GetJob(jobID)
	if err != nil {
		return err
	}

	// 启用任务
	config := job.GetConfig()
	config.Enabled = true

	// 重新调度任务
	return cm.scheduleJob(job)
}

// GetRegistry 获取任务注册器
func (cm *CronManager) GetRegistry() JobRegistry {
	return cm.registry
}

// GetExecutor 获取任务执行器
func (cm *CronManager) GetExecutor() JobExecutor {
	return cm.executor
}

// AddEventHandler 添加事件处理器
func (cm *CronManager) AddEventHandler(handler EventHandler) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.eventHandlers = append(cm.eventHandlers, handler)
}

// emitEvent 发送事件
func (cm *CronManager) emitEvent(event *JobEvent) {
	for _, handler := range cm.eventHandlers {
		if err := handler.Handle(event); err != nil {
			log.Printf("Event handler error: %v", err)
		}
	}
}

// IsStarted 检查管理器是否已启动
func (cm *CronManager) IsStarted() bool {
	return cm.started
}

// GetConfig 获取配置
func (cm *CronManager) GetConfig() *Config {
	return cm.config
}
