package cron

import (
	"context"
	"testing"
	"time"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/credentialmanagement/infrastructure/persistence"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/cache"
	credentialmanagement_mocks "github.com/context-space/context-space/backend/internal/shared/testing/mocks/credentialmanagement"
	shared_mocks "github.com/context-space/context-space/backend/internal/shared/testing/mocks/shared"
	"github.com/robfig/cron/v3"
)

func TestRefreshToken(t *testing.T) {
	// Create mock dependencies
	mockOAuthProvider := &credentialmanagement_mocks.MockOAuthProvider{}
	mockCredentialRepo := &credentialmanagement_mocks.MockCredentialRepository{}
	mockOAuthCredentialRepo := &credentialmanagement_mocks.MockOAuthCredentialRepository{}
	mockVaultService := &credentialmanagement_mocks.MockVaultService{}
	mockRedisClient := &shared_mocks.MockCache{}

	// Create mock observability components
	logConfig := &observability.LogConfig{
		Level:       observability.DebugLevel,
		Format:      observability.ConsoleFormat,
		OutputPaths: []string{"stdout"},
		Development: true,
	}
	logger, _ := observability.NewLogger(logConfig)

	mockObs := &observability.ObservabilityProvider{
		Logger:  logger,
		Tracer:  observability.NewTracer("test-tracer"),
		Metrics: &observability.Metrics{},
	}

	// Create token refresh service with all required dependencies
	tokenRefreshService := persistence.NewTokenRefreshService(
		mockRedisClient,
		mockOAuthProvider,
		mockCredentialRepo,
		mockOAuthCredentialRepo,
		mockVaultService,
		mockObs,
	)

	// Create task builder
	taskBuilder := NewCronTaskBuilder(tokenRefreshService, mockObs)

	// Test single task execution
	t.Run("SingleTask", func(t *testing.T) {
		task := taskBuilder.BuildCronTask()
		err := task.Handler(context.Background())
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	// Test all task groups execution
	t.Run("AllTaskGroups", func(t *testing.T) {
		taskGroups := taskBuilder.CreateAllTaskGroups()
		for _, taskGroup := range taskGroups {
			// Execute each task in the group
			for _, task := range taskGroup.Tasks {
				err := task.Handler(context.Background())
				if err != nil {
					t.Errorf("Expected no error for task %s, got %v", task.Name, err)
				}
			}
		}
	})

	// Test task group creation
	t.Run("TaskGroupCreation", func(t *testing.T) {
		halfHourlyGroup := taskBuilder.CreateHalfHourlyTaskGroup()
		if halfHourlyGroup.Name != "half_hourly_refresh" {
			t.Errorf("Expected group name 'half_hourly_refresh', got %s", halfHourlyGroup.Name)
		}
		if halfHourlyGroup.Schedule != "0 */3 * * *" {
			t.Errorf("Expected schedule '0 */3 * * *', got %s", halfHourlyGroup.Schedule)
		}
		if len(halfHourlyGroup.Tasks) != 1 {
			t.Errorf("Expected 1 task, got %d", len(halfHourlyGroup.Tasks))
		}
	})
}

func TestCronExpression(t *testing.T) {
	// Test if cron expression is correct
	halfHourlyGroup := &TaskGroup{
		Name:     "test_half_hourly",
		Schedule: "0 */30 * * * *", // Execute every 30 minutes
		Tasks:    []CronTask{},
	}

	// Create cron scheduler
	scheduler := cron.New(
		cron.WithLocation(time.UTC),
		cron.WithSeconds(), // Enable second-level precision
	)

	// Verify if cron expression is valid
	_, err := scheduler.AddFunc(halfHourlyGroup.Schedule, func() {
		t.Log("Cron job executed")
	})
	if err != nil {
		t.Fatalf("Invalid cron expression: %v", err)
	}

	t.Logf("Cron expression '%s' is valid", halfHourlyGroup.Schedule)
}

func TestCronManagerExecution(t *testing.T) {
	// Create mock observability
	logConfig := &observability.LogConfig{
		Level:       observability.DebugLevel,
		Format:      observability.ConsoleFormat,
		OutputPaths: []string{"stdout"},
		Development: true,
	}
	logger, _ := observability.NewLogger(logConfig)

	mockObs := &observability.ObservabilityProvider{
		Logger:  logger,
		Tracer:  observability.NewTracer("test-tracer"),
		Metrics: &observability.Metrics{},
	}

	// Create mock Redis client
	mockRedisClient := &cache.RedisClient{}

	// Create cron manager
	cronManager := NewCronManager(mockObs, mockRedisClient)

	// Create a test task group that executes every minute (for testing)
	testTaskGroup := &TaskGroup{
		Name:     "test_task_group",
		Schedule: "0 */1 * * * *", // Execute every minute
		Tasks: []CronTask{
			{
				Name: "test_task",
				Handler: func(ctx context.Context) error {
					t.Log("Test task executed successfully")
					return nil
				},
			},
		},
	}

	// Register task group
	err := cronManager.RegisterTaskGroup(context.Background(), testTaskGroup)
	if err != nil {
		t.Fatalf("Failed to register task group: %v", err)
	}

	// Start cron manager
	err = cronManager.Start(context.Background())
	if err != nil {
		t.Fatalf("Failed to start cron manager: %v", err)
	}

	// Wait for a while to let tasks execute
	time.Sleep(5 * time.Minute)

	// Stop cron manager
	err = cronManager.Stop(context.Background())
	if err != nil {
		t.Fatalf("Failed to stop cron manager: %v", err)
	}

	// Verify if task group is registered
	taskGroups := cronManager.ListTaskGroups()
	if len(taskGroups) != 1 {
		t.Errorf("Expected 1 task group, got %d", len(taskGroups))
	}
	if taskGroups[0] != "test_task_group" {
		t.Errorf("Expected task group 'test_task_group', got %s", taskGroups[0])
	}
}
