package resilience

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/context-space/context-space/backend/internal/shared/apierrors"
)

// CircuitBreakerState represents the state of a circuit breaker
type CircuitBreakerState string

const (
	// CircuitBreakerStateClosed indicates the circuit breaker is closed and requests are allowed
	CircuitBreakerStateClosed CircuitBreakerState = "CLOSED"
	// CircuitBreakerStateOpen indicates the circuit breaker is open and requests are not allowed
	CircuitBreakerStateOpen CircuitBreakerState = "OPEN"
	// CircuitBreakerStateHalfOpen indicates the circuit breaker is allowing a limited number of requests to check if the service is recovered
	CircuitBreakerStateHalfOpen CircuitBreakerState = "HALF_OPEN"
)

// CircuitBreakerConfig holds the configuration for a circuit breaker
type CircuitBreakerConfig struct {
	// Name is the name of the circuit breaker
	Name string
	// FailureThreshold is the number of consecutive failures before opening the circuit
	FailureThreshold int
	// ResetTimeout is the time to wait before transitioning from open to half-open state
	ResetTimeout time.Duration
	// HalfOpenMaxCalls is the maximum number of calls allowed in half-open state
	HalfOpenMaxCalls int
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	config          CircuitBreakerConfig
	state           CircuitBreakerState
	failureCount    int
	lastStateChange time.Time
	halfOpenCalls   int
	mutex           sync.RWMutex
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	// Set default values if not provided
	if config.FailureThreshold <= 0 {
		config.FailureThreshold = 5
	}

	if config.ResetTimeout <= 0 {
		config.ResetTimeout = 60 * time.Second
	}

	if config.HalfOpenMaxCalls <= 0 {
		config.HalfOpenMaxCalls = 1
	}

	return &CircuitBreaker{
		config:          config,
		state:           CircuitBreakerStateClosed,
		failureCount:    0,
		lastStateChange: time.Now(),
		halfOpenCalls:   0,
		mutex:           sync.RWMutex{},
	}
}

// Execute executes a function with circuit breaker protection
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
	// Check if circuit breaker is open
	if !cb.AllowRequest() {
		err := apierrors.NewCircuitOpenError(
			fmt.Sprintf("Circuit breaker '%s' is open", cb.config.Name),
			fmt.Errorf("circuit breaker open"),
		)

		return err
	}

	// Execute the function
	err := fn()

	// Update circuit breaker state based on the result
	if err != nil {
		cb.OnFailure()
		return err
	}

	// On success
	cb.OnSuccess()
	return nil
}

// AllowRequest checks if a request is allowed
func (cb *CircuitBreaker) AllowRequest() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	// If closed, always allow
	if cb.state == CircuitBreakerStateClosed {
		return true
	}

	// If open, check if reset timeout has passed
	if cb.state == CircuitBreakerStateOpen {
		if time.Since(cb.lastStateChange) > cb.config.ResetTimeout {
			// Transition to half-open state
			cb.mutex.RUnlock()
			cb.transitionState(CircuitBreakerStateHalfOpen)
			cb.mutex.RLock()
			return true
		}
		return false
	}

	// If half-open, check if we've exceeded max calls
	if cb.state == CircuitBreakerStateHalfOpen {
		return cb.halfOpenCalls < cb.config.HalfOpenMaxCalls
	}

	// Default to allowing the request
	return true
}

// OnSuccess records a successful call
func (cb *CircuitBreaker) OnSuccess() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	// Reset failure count
	cb.failureCount = 0

	// If half-open, transition to closed
	if cb.state == CircuitBreakerStateHalfOpen {
		cb.transitionState(CircuitBreakerStateClosed)
	}
}

// OnFailure records a failed call
func (cb *CircuitBreaker) OnFailure() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	// Increment failure count
	cb.failureCount++

	// Update half-open calls counter
	if cb.state == CircuitBreakerStateHalfOpen {
		cb.halfOpenCalls++
	}

	// Check if threshold is reached in closed state
	if cb.state == CircuitBreakerStateClosed && cb.failureCount >= cb.config.FailureThreshold {
		cb.transitionState(CircuitBreakerStateOpen)
		return
	}

	// Check if we should transition from half-open to open
	if cb.state == CircuitBreakerStateHalfOpen && cb.halfOpenCalls >= cb.config.HalfOpenMaxCalls {
		cb.transitionState(CircuitBreakerStateOpen)
		return
	}
}

// transitionState changes the circuit breaker state
func (cb *CircuitBreaker) transitionState(newState CircuitBreakerState) {
	cb.state = newState
	cb.lastStateChange = time.Now()

	// Reset counters on state change
	if newState == CircuitBreakerStateHalfOpen {
		cb.halfOpenCalls = 0
	} else if newState == CircuitBreakerStateClosed {
		cb.failureCount = 0
		cb.halfOpenCalls = 0
	}
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() CircuitBreakerState {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

// Reset resets the circuit breaker to the closed state
func (cb *CircuitBreaker) Reset() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.transitionState(CircuitBreakerStateClosed)
}
