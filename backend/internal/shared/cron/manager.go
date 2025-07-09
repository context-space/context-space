package cron

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/context-space/context-space/backend/internal/shared/infrastructure/cache"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	observability "github.com/context-space/cloud-observability"
)

const (
	// CronTaskLockKey is the distributed lock key format for cron tasks
	CronTaskLockKey = "cron_task_lock:%s"
	// CronTaskLockTimeout is the timeout for cron task locks
	CronTaskLockTimeout = 10 * time.Minute
)

// CronTask represents a single executable task
type CronTask struct {
	Name    string
	Handler func(ctx context.Context) error
}

// TaskGroup represents a group of tasks with the same schedule
type TaskGroup struct {
	Name     string
	Schedule string // cron expression
	Tasks    []CronTask
}

// TaskGroupConfig holds configuration for task groups
type TaskGroupConfig struct {
	// Task group executed on the 15th of each month
	Monthly15th TaskGroup
	// Task group executed on the 30th of each month
	Monthly30th TaskGroup
	// Task group executed every half hour
	HalfHourly TaskGroup
}

type taskResult struct {
	task CronTask
	err  error
}

// CronManager manages multiple cron task groups
type CronManager struct {
	scheduler   *cron.Cron
	taskGroups  map[string]*TaskGroup
	obs         *observability.ObservabilityProvider
	redisClient cache.Cache

	// Internal state
	mu        sync.RWMutex
	isRunning bool
}

// NewCronManager creates a new cron manager
func NewCronManager(obs *observability.ObservabilityProvider, redisClient cache.Cache) *CronManager {
	// Create cron scheduler with timezone support and middleware
	scheduler := cron.New(
		cron.WithLocation(time.UTC),
		cron.WithSeconds(),
		cron.WithChain(
			cron.SkipIfStillRunning(cron.DiscardLogger),
			cron.Recover(cron.DefaultLogger),
		),
	)

	return &CronManager{
		scheduler:   scheduler,
		taskGroups:  make(map[string]*TaskGroup),
		obs:         obs,
		redisClient: redisClient,
	}
}

// RegisterTaskGroup registers a task group with its schedule
func (m *CronManager) RegisterTaskGroup(ctx context.Context, group *TaskGroup) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if group == nil {
		return fmt.Errorf("task group cannot be nil")
	}

	if group.Name == "" {
		return fmt.Errorf("task group name cannot be empty")
	}

	if group.Schedule == "" {
		return fmt.Errorf("task group schedule cannot be empty")
	}

	if len(group.Tasks) == 0 {
		m.obs.Logger.Warn(ctx, "Registering task group with no tasks",
			zap.String("group", group.Name))
	}

	// Add the cron job for this task group
	_, err := m.scheduler.AddFunc(group.Schedule, func() {
		m.executeTaskGroup(ctx, group)
	})

	if err != nil {
		return fmt.Errorf("failed to schedule task group %s: %w", group.Name, err)
	}

	m.taskGroups[group.Name] = group

	m.obs.Logger.Info(ctx, "Task group registered",
		zap.String("group", group.Name),
		zap.String("schedule", group.Schedule),
		zap.Int("task_count", len(group.Tasks)),
	)

	return nil
}

// Start begins the cron scheduler
func (m *CronManager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.isRunning {
		return fmt.Errorf("cron manager is already running")
	}

	if len(m.taskGroups) == 0 {
		m.obs.Logger.Warn(ctx, "Starting cron manager with no registered task groups")
	}

	m.scheduler.Start()
	m.isRunning = true

	m.obs.Logger.Info(ctx, "Cron manager started",
		zap.Int("task_groups", len(m.taskGroups)),
	)

	return nil
}

// Stop gracefully stops the cron scheduler
func (m *CronManager) Stop(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isRunning {
		return nil
	}

	stopCtx := m.scheduler.Stop()

	// Wait for graceful shutdown or context timeout
	select {
	case <-stopCtx.Done():
		m.obs.Logger.Info(ctx, "Cron manager stopped gracefully")
	case <-ctx.Done():
		m.obs.Logger.Warn(ctx, "Cron manager forced shutdown due to context timeout")
	}

	m.isRunning = false
	return nil
}

// executeTaskGroup executes all tasks in a task group
func (m *CronManager) executeTaskGroup(ctx context.Context, group *TaskGroup) {
	// Add top-level panic recovery
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			m.obs.Logger.Error(ctx, "executeTaskGroup recovered from panic",
				zap.String("group", group.Name),
				zap.Any("recover", r),
				zap.String("stack", string(buf[:n])))
		}
	}()

	ctx, span := m.obs.Tracer.Start(ctx, "CronManager.executeTaskGroup")
	defer span.End()

	startTime := time.Now()

	m.obs.Logger.Info(ctx, "Task group execution started",
		zap.String("group", group.Name),
	)

	const maxConcurrency = 5 // Limit concurrency
	semaphore := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup
	results := make(chan taskResult, len(group.Tasks))

	taskList := make([]string, 0, len(group.Tasks))

	for _, task := range group.Tasks {
		taskCopy := task // Avoid closure references
		lockKey := fmt.Sprintf(CronTaskLockKey, taskCopy.Name)

		// Try to acquire distributed lock
		lock, err := m.redisClient.AcquireLock(ctx, lockKey, CronTaskLockTimeout)
		if err != nil {
			m.obs.Logger.Error(ctx, "Failed to acquire lock",
				zap.String("task", taskCopy.Name), zap.Error(err))
			continue
		}
		if !lock {
			m.obs.Logger.Info(ctx, "Skipped task due to lock held by another instance",
				zap.String("task", taskCopy.Name))
			continue
		}

		wg.Add(1)
		go func(t CronTask, lockKey string) {
			defer wg.Done()

			defer func() {
				if r := recover(); r != nil {
					buf := make([]byte, 4096)
					n := runtime.Stack(buf, false)
					m.obs.Logger.Error(ctx, "Recovered from panic in task",
						zap.String("task", t.Name), zap.Any("recover", r),
						zap.String("stack", string(buf[:n])))
				}

				if err := m.redisClient.ReleaseLock(ctx, lockKey); err != nil {
					m.obs.Logger.Error(ctx, "Failed to release lock",
						zap.String("task", t.Name), zap.Error(err))
				}

			}()

			// Acquire concurrency semaphore
			select {
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			case <-ctx.Done():
				m.obs.Logger.Warn(ctx, "Context cancelled before acquiring semaphore",
					zap.String("task", t.Name))
				return
			case <-time.After(10 * time.Minute):
				m.obs.Logger.Error(ctx, "Timeout acquiring semaphore",
					zap.String("task", t.Name))
				return
			}

			taskCtx, span := m.obs.Tracer.Start(ctx, "CronTask."+t.Name)
			defer span.End()

			err := m.executeTask(taskCtx, t)

			select {
			case results <- taskResult{task: t, err: err}:
			case <-ctx.Done():
				m.obs.Logger.Warn(ctx, "Context cancelled before sending result",
					zap.String("task", t.Name))
			case <-time.After(30 * time.Second):
				m.obs.Logger.Error(ctx, "Timeout sending task result",
					zap.String("task", t.Name))
			}
		}(taskCopy, lockKey)

		taskList = append(taskList, taskCopy.Name)
	}

	wg.Wait()
	close(results)

	// Collect results
	var successCount, failCount int64
	for result := range results {
		if result.err != nil {
			failCount++
			m.obs.Logger.Error(ctx, "Task execution failed",
				zap.String("task", result.task.Name), zap.Error(result.err))
		} else {
			successCount++
		}
	}

	duration := time.Since(startTime)

	m.obs.Logger.Info(ctx, "Task group execution completed",
		zap.String("group", group.Name),
		zap.Strings("task_list", taskList),
		zap.Int64("duration", duration.Milliseconds()),
		zap.Int64("successful", successCount),
		zap.Int64("failed", failCount),
	)
}

// executeTask executes a single task with timeout protection
func (m *CronManager) executeTask(ctx context.Context, task CronTask) error {
	// Create timeout context for task execution (10 minutes)
	taskCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	// Execute the task
	return task.Handler(taskCtx)
}

// IsRunning returns whether the cron manager is currently running
func (m *CronManager) IsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.isRunning
}

// ListTaskGroups returns all registered task groups
func (m *CronManager) ListTaskGroups() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	groups := make([]string, 0, len(m.taskGroups))
	for name := range m.taskGroups {
		groups = append(groups, name)
	}
	return groups
}
