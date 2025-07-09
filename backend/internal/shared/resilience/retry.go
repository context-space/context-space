package resilience

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"
)

// RetryFunction is a function that can be retried
type RetryFunction func() error

// RetryPolicy defines when to retry and when to give up
type RetryPolicy interface {
	// ShouldRetry returns true if the operation should be retried after the given error
	ShouldRetry(err error, attempt int) bool

	// NextBackoff returns the backoff duration for the next retry
	NextBackoff(attempt int) time.Duration

	// MaxAttempts returns the maximum number of attempts
	MaxAttempts() int
}

// ExponentialBackoffPolicy implements RetryPolicy with exponential backoff
type ExponentialBackoffPolicy struct {
	maxAttempts     int
	initialBackoff  time.Duration
	maxBackoff      time.Duration
	backoffFactor   float64
	retryableErrors []func(error) bool
}

// NewExponentialBackoffPolicy creates a new exponential backoff policy
func NewExponentialBackoffPolicy(
	maxAttempts int,
	initialBackoff time.Duration,
	maxBackoff time.Duration,
	backoffFactor float64,
	retryableErrors []func(error) bool,
) *ExponentialBackoffPolicy {
	return &ExponentialBackoffPolicy{
		maxAttempts:     maxAttempts,
		initialBackoff:  initialBackoff,
		maxBackoff:      maxBackoff,
		backoffFactor:   backoffFactor,
		retryableErrors: retryableErrors,
	}
}

// ShouldRetry returns true if the operation should be retried
func (p *ExponentialBackoffPolicy) ShouldRetry(err error, attempt int) bool {
	// Check if we've exceeded the maximum number of attempts
	if attempt >= p.maxAttempts {
		return false
	}

	// Check if the error is retryable
	for _, isRetryable := range p.retryableErrors {
		if isRetryable(err) {
			return true
		}
	}

	return false
}

// NextBackoff returns the backoff duration for the next retry
func (p *ExponentialBackoffPolicy) NextBackoff(attempt int) time.Duration {
	// Calculate exponential backoff
	backoff := float64(p.initialBackoff) * math.Pow(p.backoffFactor, float64(attempt))

	// Add jitter (randomness) to avoid thundering herd problem
	jitter := (0.5 + 0.5*math.Sin(float64(attempt)*math.Pi/3)) // value between 0 and 1
	backoff = backoff * (0.5 + jitter)

	// Ensure we don't exceed max backoff
	if backoff > float64(p.maxBackoff) {
		backoff = float64(p.maxBackoff)
	}

	return time.Duration(backoff)
}

// MaxAttempts returns the maximum number of attempts
func (p *ExponentialBackoffPolicy) MaxAttempts() int {
	return p.maxAttempts
}

// Retry retries a function according to the retry policy
func Retry(ctx context.Context, fn RetryFunction, policy RetryPolicy) error {
	var lastErr error

	for attempt := 0; attempt < policy.MaxAttempts(); attempt++ {
		// Execute the function
		err := fn()
		if err == nil {
			// Success, return nil
			return nil
		}

		lastErr = err

		// Check if we should retry
		if !policy.ShouldRetry(err, attempt) {
			// We shouldn't retry, return the error
			return fmt.Errorf("giving up after %d attempts: %w", attempt+1, err)
		}

		// Calculate backoff duration
		backoff := policy.NextBackoff(attempt)

		// Wait for the backoff duration or until context is canceled
		select {
		case <-ctx.Done():
			// Context was canceled, return the context error
			return fmt.Errorf("retry canceled by context: %w", ctx.Err())
		case <-time.After(backoff):
			// Backoff completed, continue to next attempt
		}
	}

	// If we get here, we've exceeded the maximum number of attempts
	return fmt.Errorf("exceeded maximum retry attempts: %w", lastErr)
}

// IsNetworkError returns true if the error is likely a network error
func IsNetworkError(err error) bool {
	// Check for common network error strings
	// This is a simplified implementation, in a real system you would check for specific error types
	errStr := err.Error()
	networkErrorPhrases := []string{
		"connection refused",
		"connection reset",
		"connection closed",
		"no route to host",
		"network is unreachable",
		"timeout",
		"connection timed out",
		"dial tcp",
	}

	for _, phrase := range networkErrorPhrases {
		if strings.Contains(errStr, phrase) {
			return true
		}
	}

	return false
}

// IsTemporaryError checks if the error is a temporary error that might succeed on retry
func IsTemporaryError(err error) bool {
	// Check if the error implements a Temporary() method
	type temporary interface {
		Temporary() bool
	}

	te, ok := err.(temporary)
	return ok && te.Temporary()
}

// IsTimeoutError checks if the error is a timeout error
func IsTimeoutError(err error) bool {
	// Check if the error implements a Timeout() method
	type timeout interface {
		Timeout() bool
	}

	te, ok := err.(timeout)
	return ok && te.Timeout()
}

// IsServerError checks if the error is a 5xx server error
func IsServerError(err error) bool {
	// This is a simplified implementation, in a real system you would check HTTP status codes
	errStr := err.Error()
	serverErrorPhrases := []string{
		"500",
		"501",
		"502",
		"503",
		"504",
		"505",
		"internal server error",
		"bad gateway",
		"service unavailable",
		"gateway timeout",
	}

	for _, phrase := range serverErrorPhrases {
		if strings.Contains(errStr, phrase) {
			return true
		}
	}

	return false
}

// DefaultRetryableErrors returns a list of default error checkers for retryable errors
func DefaultRetryableErrors() []func(error) bool {
	return []func(error) bool{
		IsNetworkError,
		IsTemporaryError,
		IsTimeoutError,
		IsServerError,
	}
}
