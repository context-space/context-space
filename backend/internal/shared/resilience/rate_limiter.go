package resilience

import (
	"sync"
	"time"

	"github.com/context-space/context-space/backend/internal/shared/utils"
	"golang.org/x/time/rate"
)

// RateLimiterConfig holds the configuration for a rate limiter
type RateLimiterConfig struct {
	// Name is the name of the rate limiter
	Name string
	// RPS is the maximum requests per second
	RPS int
	// Burst is the maximum burst size
	Burst int
	// KeyPrefix is a prefix for cache keys
	KeyPrefix string
}

// RateLimiter implements the rate limiter pattern
type RateLimiter struct {
	config     RateLimiterConfig
	limiters   map[string]*rate.Limiter
	mutex      sync.RWMutex
	cleanupTTL time.Duration
	lastClean  time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	// Set default values if not provided
	if config.RPS <= 0 {
		config.RPS = 10
	}

	if config.Burst <= 0 {
		config.Burst = config.RPS
	}

	return &RateLimiter{
		config:     config,
		limiters:   make(map[string]*rate.Limiter),
		mutex:      sync.RWMutex{},
		cleanupTTL: 1 * time.Hour,
		lastClean:  time.Now(),
	}
}

// getLimiter gets or creates a limiter for the given key
func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	// Clean old limiters if needed
	rl.cleanupIfNeeded()

	// Check if limiter exists
	rl.mutex.RLock()
	limiter, exists := rl.limiters[key]
	rl.mutex.RUnlock()

	if exists {
		return limiter
	}

	// Create new limiter
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	// Check again in case another goroutine created the limiter
	limiter, exists = rl.limiters[key]
	if exists {
		return limiter
	}

	// Create new limiter
	limiter = rate.NewLimiter(rate.Limit(rl.config.RPS), rl.config.Burst)
	rl.limiters[key] = limiter

	return limiter
}

// Allow checks if a request is allowed for the given key
func (rl *RateLimiter) Allow(key string) bool {
	// Create proper key with prefix
	fullKey := key
	if rl.config.KeyPrefix != "" {
		fullKey = utils.StringsBuilder(rl.config.KeyPrefix, ":", key)
	}

	// Get limiter
	limiter := rl.getLimiter(fullKey)

	// Check if request is allowed
	return limiter.Allow()
}

// AllowN checks if N requests are allowed for the given key
func (rl *RateLimiter) AllowN(key string, n int) bool {
	// Create proper key with prefix
	fullKey := key
	if rl.config.KeyPrefix != "" {
		fullKey = utils.StringsBuilder(rl.config.KeyPrefix, ":", key)
	}

	// Get limiter
	limiter := rl.getLimiter(fullKey)

	// Check if N requests are allowed
	return limiter.AllowN(time.Now(), n)
}

// cleanupIfNeeded removes old limiters that haven't been used
func (rl *RateLimiter) cleanupIfNeeded() {
	// Only clean once per hour
	if time.Since(rl.lastClean) < rl.cleanupTTL {
		return
	}

	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	// Update last clean time
	rl.lastClean = time.Now()

	// Create new map
	rl.limiters = make(map[string]*rate.Limiter)
}
