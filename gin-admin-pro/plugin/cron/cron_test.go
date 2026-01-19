package cron

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.NotNil(t, config)
	assert.True(t, config.Enabled)
	assert.Equal(t, "UTC", config.Timezone)
	assert.Equal(t, 100, config.Concurrency)
	assert.Equal(t, 3, config.MaxRetries)
	assert.True(t, config.LogEnabled)
	assert.Equal(t, "info", config.LogLevel)
}

func TestDefaultJobConfig(t *testing.T) {
	config := DefaultJobConfig()

	assert.NotNil(t, config)
	assert.True(t, config.Enabled)
	assert.Equal(t, "0 */1 * * *", config.Cron)
	assert.Equal(t, "UTC", config.Timezone)
	assert.Equal(t, 3, config.MaxRetries)
	assert.Equal(t, 1, config.MaxInstances)
	assert.True(t, config.Singleton)
}

func TestNewCronManager(t *testing.T) {
	config := &Config{
		Enabled:     true,
		Timezone:    "Asia/Shanghai",
		Concurrency: 10,
	}

	manager, err := NewCronManager(config)
	assert.NoError(t, err)
	assert.NotNil(t, manager)

	retrievedConfig := manager.GetConfig()
	assert.Equal(t, true, retrievedConfig.Enabled)
	assert.Equal(t, "Asia/Shanghai", retrievedConfig.Timezone)
	assert.Equal(t, 10, retrievedConfig.Concurrency)
}

func TestCronManager_Initialize(t *testing.T) {
	manager, err := NewCronManager(nil)
	require.NoError(t, err)

	err = manager.Initialize()
	assert.NoError(t, err)

	// 检查默认处理器是否注册
	registry := manager.GetRegistry()
	handlers := registry.ListHandlers()

	assert.Contains(t, handlers, "sample")
	assert.Contains(t, handlers, "cleanup")
	assert.Contains(t, handlers, "report")
}

func TestCronManager_AddJob(t *testing.T) {
	manager, err := NewCronManager(nil)
	require.NoError(t, err)

	err = manager.Initialize()
	require.NoError(t, err)

	// 添加任务
	jobConfig := &JobConfig{
		ID:          "test-job",
		Name:        "Test Job",
		Description: "Test job description",
		Enabled:     true,
		Cron:        "*/5 * * * *", // 每5分钟执行
		Handler:     "sample",
		Timeout:     time.Second * 30,
		MaxRetries:  3,
	}

	err = manager.AddJob(jobConfig)
	assert.NoError(t, err)

	// 获取任务
	job, err := manager.GetJob("test-job")
	assert.NoError(t, err)
	assert.Equal(t, "test-job", job.GetID())
	assert.Equal(t, "Test Job", job.GetName())
	assert.Equal(t, "Test job description", job.GetDescription())
	assert.True(t, job.IsEnabled())
}

func TestCronManager_StartStop(t *testing.T) {
	manager, err := NewCronManager(&Config{Enabled: true})
	require.NoError(t, err)

	err = manager.Initialize()
	require.NoError(t, err)

	// 添加任务
	jobConfig := &JobConfig{
		ID:      "test-job",
		Name:    "Test Job",
		Enabled: true,
		Cron:    "*/1 * * * *", // 每分钟执行
		Handler: "sample",
	}

	err = manager.AddJob(jobConfig)
	assert.NoError(t, err)

	// 启动管理器
	err = manager.Start()
	assert.NoError(t, err)
	assert.True(t, manager.IsStarted())

	// 等待一段时间
	time.Sleep(time.Millisecond * 100)

	// 停止管理器
	err = manager.Stop()
	assert.NoError(t, err)
	assert.False(t, manager.IsStarted())
}

func TestCronManager_RunJob(t *testing.T) {
	manager, err := NewCronManager(nil)
	require.NoError(t, err)

	err = manager.Initialize()
	require.NoError(t, err)

	// 添加任务
	jobConfig := &JobConfig{
		ID:      "test-job",
		Name:    "Test Job",
		Enabled: true,
		Cron:    "*/5 * * * *",
		Handler: "sample",
	}

	err = manager.AddJob(jobConfig)
	assert.NoError(t, err)

	// 手动运行任务
	err = manager.RunJob("test-job", nil)
	assert.NoError(t, err)

	// 等待任务完成
	time.Sleep(time.Millisecond * 500)
}

func TestCronManager_PauseResumeJob(t *testing.T) {
	manager, err := NewCronManager(nil)
	require.NoError(t, err)

	err = manager.Initialize()
	require.NoError(t, err)

	// 添加任务
	jobConfig := &JobConfig{
		ID:      "test-job",
		Name:    "Test Job",
		Enabled: true,
		Cron:    "*/1 * * * *",
		Handler: "sample",
	}

	err = manager.AddJob(jobConfig)
	assert.NoError(t, err)

	// 启动管理器
	err = manager.Start()
	require.NoError(t, err)

	// 暂停任务
	err = manager.PauseJob("test-job")
	assert.NoError(t, err)

	job, _ := manager.GetJob("test-job")
	assert.False(t, job.IsEnabled())

	// 恢复任务
	err = manager.ResumeJob("test-job")
	assert.NoError(t, err)

	job, _ = manager.GetJob("test-job")
	assert.True(t, job.IsEnabled())

	// 停止管理器
	err = manager.Stop()
	assert.NoError(t, err)
}

func TestCronManager_ListJobs(t *testing.T) {
	manager, err := NewCronManager(nil)
	require.NoError(t, err)

	err = manager.Initialize()
	require.NoError(t, err)

	// 添加多个任务
	jobs := []*JobConfig{
		{ID: "job1", Name: "Job 1", Enabled: true, Cron: "*/5 * * * *", Handler: "sample"},
		{ID: "job2", Name: "Job 2", Enabled: true, Cron: "*/10 * * * *", Handler: "cleanup"},
		{ID: "job3", Name: "Job 3", Enabled: false, Cron: "*/15 * * * *", Handler: "report"},
	}

	for _, jobConfig := range jobs {
		err := manager.AddJob(jobConfig)
		assert.NoError(t, err)
	}

	// 列出所有任务
	allJobs := manager.ListJobs()
	assert.Len(t, allJobs, 3)

	// 获取任务
	job, err := manager.GetJob("job1")
	assert.NoError(t, err)
	assert.Equal(t, "job1", job.GetID())
}

func TestCronManager_RemoveJob(t *testing.T) {
	manager, err := NewCronManager(nil)
	require.NoError(t, err)

	err = manager.Initialize()
	require.NoError(t, err)

	// 添加任务
	jobConfig := &JobConfig{
		ID:      "test-job",
		Name:    "Test Job",
		Enabled: true,
		Cron:    "*/5 * * * *",
		Handler: "sample",
	}

	err = manager.AddJob(jobConfig)
	assert.NoError(t, err)

	// 确认任务存在
	_, err = manager.GetJob("test-job")
	assert.NoError(t, err)

	// 移除任务
	err = manager.RemoveJob("test-job")
	assert.NoError(t, err)

	// 确认任务已移除
	_, err = manager.GetJob("test-job")
	assert.Error(t, err)
}

func TestDefaultRegistry(t *testing.T) {
	registry := &DefaultRegistry{
		handlers: make(map[string]JobHandler),
	}

	// 注册处理器
	handler := &SampleJobHandler{}
	err := registry.RegisterHandler("test", handler)
	assert.NoError(t, err)

	// 重复注册应该失败
	err = registry.RegisterHandler("test", handler)
	assert.Error(t, err)

	// 获取处理器
	retrievedHandler, err := registry.GetHandler("test")
	assert.NoError(t, err)
	assert.Equal(t, handler, retrievedHandler)

	// 获取不存在的处理器
	_, err = registry.GetHandler("nonexistent")
	assert.Error(t, err)

	// 列出处理器
	handlers := registry.ListHandlers()
	assert.Contains(t, handlers, "test")

	// 注销处理器
	err = registry.UnregisterHandler("test")
	assert.NoError(t, err)

	// 注销不存在的处理器
	err = registry.UnregisterHandler("nonexistent")
	assert.Error(t, err)
}

func TestDefaultExecutor(t *testing.T) {
	manager, err := NewCronManager(nil)
	require.NoError(t, err)

	// Initialize the manager to register handlers
	err = manager.Initialize()
	require.NoError(t, err)

	executor := NewDefaultExecutor(manager)
	assert.NotNil(t, executor)

	// 创建任务
	jobConfig := &JobConfig{
		ID:      "test-job",
		Name:    "Test Job",
		Enabled: true,
		Cron:    "*/5 * * * *",
		Handler: "sample",
		Timeout: time.Second * 10,
	}

	job := NewDefaultJob(jobConfig)

	// 执行任务
	ctx := context.Background()
	err = executor.Execute(ctx, job)
	assert.NoError(t, err)

	// 获取执行历史
	history, err := executor.GetExecutionHistory("test-job", 10)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(history), 1)
	// Check the last execution is completed
	lastExecution := history[len(history)-1]
	assert.Equal(t, JobStatusCompleted, lastExecution.Status)
}

func TestDefaultJob(t *testing.T) {
	config := &JobConfig{
		ID:          "test-job",
		Name:        "Test Job",
		Description: "Test job description",
		Enabled:     true,
		Cron:        "*/5 * * * *",
		Handler:     "sample",
		Tags:        []string{"test", "sample"},
		Metadata:    map[string]string{"env": "test"},
	}

	job := NewDefaultJob(config)

	// 基础方法测试
	assert.Equal(t, "test-job", job.GetID())
	assert.Equal(t, "Test Job", job.GetName())
	assert.Equal(t, "Test job description", job.GetDescription())
	assert.Equal(t, config, job.GetConfig())
	assert.True(t, job.IsEnabled())
	assert.Equal(t, []string{"test", "sample"}, job.GetTags())
	assert.Equal(t, "test", job.GetMetadata()["env"])

	// 状态测试
	assert.Equal(t, JobStatusPending, job.GetStatus())
	assert.False(t, job.IsRunning())

	// 执行控制测试
	ctx := context.Background()
	err := job.Execute(ctx, nil)
	assert.NoError(t, err)

	err = job.Pause()
	assert.NoError(t, err)
	assert.Equal(t, JobStatusPaused, job.GetStatus())

	err = job.Resume()
	assert.NoError(t, err)

	err = job.Cancel()
	assert.NoError(t, err)
	assert.Equal(t, JobStatusCancelled, job.GetStatus())
}

func TestSampleJobHandler(t *testing.T) {
	handler := &SampleJobHandler{}

	config := &JobConfig{
		ID:      "test-job",
		Name:    "Test Job",
		Enabled: true,
		Handler: "sample",
	}

	job := NewDefaultJob(config)
	ctx := context.Background()

	// 测试生命周期方法
	err := handler.OnStart(ctx, job)
	assert.NoError(t, err)

	err = handler.Handle(ctx, job, nil)
	assert.NoError(t, err)

	err = handler.OnComplete(ctx, job, nil)
	assert.NoError(t, err)

	err = handler.OnError(ctx, job, nil)
	assert.NoError(t, err)

	err = handler.OnRetry(ctx, job, 1, nil)
	assert.NoError(t, err)
}

func TestJobExecution(t *testing.T) {
	execution := &JobExecution{
		ID:          "test-execution",
		JobID:       "test-job",
		JobName:     "Test Job",
		Status:      JobStatusRunning,
		StartTime:   time.Now(),
		Attempt:     1,
		MaxAttempts: 3,
		Params:      map[string]interface{}{"key": "value"},
		Trigger:     "manual",
		TriggeredBy: "admin",
		Metadata:    map[string]string{"env": "test"},
	}

	assert.Equal(t, "test-execution", execution.ID)
	assert.Equal(t, "test-job", execution.JobID)
	assert.Equal(t, "Test Job", execution.JobName)
	assert.Equal(t, JobStatusRunning, execution.Status)
	assert.Equal(t, "manual", execution.Trigger)
	assert.Equal(t, "admin", execution.TriggeredBy)
	assert.Equal(t, "value", execution.Params["key"])
	assert.Equal(t, "test", execution.Metadata["env"])
}

// Integration tests (would require more complex setup)
func TestCronManager_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This test would require a more complex setup
	// with actual cron execution and verification
	// For now, just test basic functionality

	manager, err := NewCronManager(&Config{
		Enabled:  true,
		Timezone: "UTC",
	})
	require.NoError(t, err)

	err = manager.Initialize()
	require.NoError(t, err)

	// Add a job that runs every second for testing
	jobConfig := &JobConfig{
		ID:      "integration-test",
		Name:    "Integration Test Job",
		Enabled: true,
		Cron:    "*/1 * * * * *", // Note: This format is not standard, should be "* * * * *"
		Handler: "sample",
		Timeout: time.Second * 5,
	}

	// Fix the cron expression
	jobConfig.Cron = "* * * * *" // Every minute

	err = manager.AddJob(jobConfig)
	assert.NoError(t, err)

	err = manager.Start()
	assert.NoError(t, err)

	// Wait a bit for potential execution
	time.Sleep(time.Millisecond * 100)

	err = manager.Stop()
	assert.NoError(t, err)
}
