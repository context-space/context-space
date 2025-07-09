package domain

import (
	"encoding/json"
	"time"
)

// InvocationStatus represents the status of an invocation
type InvocationStatus string

const (
	// InvocationStatusPending represents a pending invocation
	InvocationStatusPending InvocationStatus = "pending"
	// InvocationStatusSuccess represents a successful invocation
	InvocationStatusSuccess InvocationStatus = "success"
	// InvocationStatusFailed represents a failed invocation
	InvocationStatusFailed InvocationStatus = "failed"
)

// Invocation represents an invocation of an operation on a provider
type Invocation struct {
	ID                  string
	UserID              string
	ProviderIdentifier  string
	OperationIdentifier string
	Status              InvocationStatus
	Duration            int64
	StartedAt           *time.Time
	CompletedAt         *time.Time
	ErrorMessage        string
	Parameters          map[string]interface{}
	ResponseData        json.RawMessage
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time
}

// NewInvocation creates a new invocation
func NewInvocation(id, userID, providerIdentifier, operationIdentifier string, parameters map[string]interface{}) *Invocation {
	return &Invocation{
		ID:                  id,
		UserID:              userID,
		ProviderIdentifier:  providerIdentifier,
		OperationIdentifier: operationIdentifier,
		Status:              InvocationStatusPending,
		Parameters:          parameters,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
}

// SetStarted marks the invocation as started
func (i *Invocation) SetStarted() {
	now := time.Now()
	i.Status = InvocationStatusPending
	i.StartedAt = &now
	i.UpdatedAt = now
}

// SetSuccess marks the invocation as successful
func (i *Invocation) SetSuccess(responseData json.RawMessage) {
	now := time.Now()
	i.Status = InvocationStatusSuccess
	i.ResponseData = responseData
	i.CompletedAt = &now
	i.Duration = i.CalculateDuration()
	i.UpdatedAt = now
}

// SetFailed marks the invocation as failed
func (i *Invocation) SetFailed(errorMessage string) {
	now := time.Now()
	i.Status = InvocationStatusFailed
	i.ErrorMessage = errorMessage
	i.CompletedAt = &now
	i.Duration = i.CalculateDuration()
	i.UpdatedAt = now
}

// CalculateDuration calculates the duration between start and completion
func (i *Invocation) CalculateDuration() int64 {
	if i.StartedAt == nil || i.CompletedAt == nil {
		return 0
	}
	return i.CompletedAt.UnixNano()/int64(time.Millisecond) - i.StartedAt.UnixNano()/int64(time.Millisecond)
}

// IsCompleted returns true if the invocation is completed
func (i *Invocation) IsCompleted() bool {
	return i.Status == InvocationStatusSuccess || i.Status == InvocationStatusFailed
}

// IsSuccessful returns true if the invocation was successful
func (i *Invocation) IsSuccessful() bool {
	return i.Status == InvocationStatusSuccess
}

// IsFailed returns true if the invocation failed
func (i *Invocation) IsFailed() bool {
	return i.Status == InvocationStatusFailed
}
