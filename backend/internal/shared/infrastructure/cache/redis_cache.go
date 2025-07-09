package cache

import (
	"context"
	"fmt"
	"time"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/shared/config"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	AccessTokenLockKey     = "access_token_lock:%s:%s"
	AccessTokenLockTimeout = 1 * time.Second
)

// RedisClient wraps a Redis client with observability
type RedisClient struct {
	client          *redis.Client
	obs             *observability.ObservabilityProvider
	traceOperations bool
}

// NewRedisClient creates a new Redis client
func NewRedisClient(cfg *config.Config, observabilityProvider *observability.ObservabilityProvider) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisClient{
		client:          client,
		obs:             observabilityProvider,
		traceOperations: cfg.Logging.Level == "debug",
	}, nil
}

// Set stores a key-value pair with expiration
func (c *RedisClient) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	var span trace.Span
	if c.traceOperations {
		ctx, span = c.obs.Tracer.Start(ctx, "redis.Set")
		span.SetAttributes(attribute.String("key", key))
		span.SetAttributes(attribute.String("expiration", expiration.String()))
		defer span.End()
	}

	return c.client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value by key
func (c *RedisClient) Get(ctx context.Context, key string) (string, error) {
	var span trace.Span
	if c.traceOperations {
		ctx, span = c.obs.Tracer.Start(ctx, "redis.Get")
		span.SetAttributes(attribute.String("key", key))
		defer span.End()
	}

	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key not found: %s", key)
	} else if err != nil {
		return "", err
	}

	return val, nil
}

// Delete removes a key
func (c *RedisClient) Delete(ctx context.Context, key string) error {
	var span trace.Span
	if c.traceOperations {
		ctx, span = c.obs.Tracer.Start(ctx, "redis.Delete")
		span.SetAttributes(attribute.String("key", key))
		defer span.End()
	}

	return c.client.Del(ctx, key).Err()
}

// Close closes the Redis connection
func (c *RedisClient) Close() error {
	return c.client.Close()
}

func (c *RedisClient) AcquireLock(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	var span trace.Span
	if c.traceOperations {
		ctx, span = c.obs.Tracer.Start(ctx, "redis.AcquireLock")
		span.SetAttributes(attribute.String("key", key))
		span.SetAttributes(attribute.String("expiration", expiration.String()))
		defer span.End()
	}

	return c.client.SetNX(ctx, key, "1", expiration).Result()
}

func (c *RedisClient) ReleaseLock(ctx context.Context, key string) error {
	var span trace.Span
	if c.traceOperations {
		ctx, span = c.obs.Tracer.Start(ctx, "redis.ReleaseLock")
		span.SetAttributes(attribute.String("key", key))
		defer span.End()
	}

	return c.client.Del(ctx, key).Err()
}
