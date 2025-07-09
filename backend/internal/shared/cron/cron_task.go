package cron

import (
	"context"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	"go.uber.org/zap"
)

// CronTaskBuilder constructs token refresh tasks for different providers
type CronTaskBuilder struct {
	tokenRefreshService domain.TokenRefresh
	obs                 *observability.ObservabilityProvider
}

// NewCronTaskBuilder creates a new token refresh task builder
func NewCronTaskBuilder(
	tokenRefreshService domain.TokenRefresh,
	obs *observability.ObservabilityProvider,
) *CronTaskBuilder {
	return &CronTaskBuilder{
		tokenRefreshService: tokenRefreshService,
		obs:                 obs,
	}
}

// BuildCronTask creates a cron task for a specific provider
func (b *CronTaskBuilder) BuildCronTask() CronTask {
	return CronTask{
		Name: "refresh_access_tokens",
		Handler: func(ctx context.Context) error {
			return b.refreshAccessTokens(ctx)
		},
	}
}

// CreateHalfHourlyTaskGroup creates the task group for half-hourly execution
func (b *CronTaskBuilder) CreateHalfHourlyTaskGroup() *TaskGroup {
	return &TaskGroup{
		Name:     "half_hourly_refresh",
		Schedule: "0 */30 * * * *", // Execute every 30 minutes (6-field cron expression)
		Tasks: []CronTask{
			b.BuildCronTask(),
		},
	}
}

// CreateAllTaskGroups creates all predefined task groups
func (b *CronTaskBuilder) CreateAllTaskGroups() []*TaskGroup {
	return []*TaskGroup{
		b.CreateHalfHourlyTaskGroup(),
	}
}

// refreshTokensForProvider refreshes tokens for a specific provider
func (b *CronTaskBuilder) refreshAccessTokens(ctx context.Context) error {
	ctx, span := b.obs.Tracer.Start(ctx, "TokenRefreshTaskBuilder.refreshTokensForProvider")
	defer span.End()

	b.obs.Logger.Info(ctx, "Starting token refresh for all providers")

	err := b.tokenRefreshService.RefreshAccessTokens(ctx)
	if err != nil {
		b.obs.Logger.Error(ctx, "Failed to refresh access tokens", zap.Error(err))
		return err
	}

	return nil
}
