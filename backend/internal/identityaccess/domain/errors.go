package domain

import "fmt"

// ValidationError represents a validation error
type ValidationError struct {
	Message string
	Err     error
}

// Error returns the error message
func (e *ValidationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Err.Error())
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *ValidationError) Unwrap() error {
	return e.Err
}

// ConflictError represents a conflict error
type ConflictError struct {
	Message string
	Err     error
}

// Error returns the error message
func (e *ConflictError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Err.Error())
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *ConflictError) Unwrap() error {
	return e.Err
}

// NotFoundError represents a not found error
type NotFoundError struct {
	Message string
	Err     error
}

// Error returns the error message
func (e *NotFoundError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Err.Error())
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *NotFoundError) Unwrap() error {
	return e.Err
}

// UnauthorizedError represents an unauthorized error
type UnauthorizedError struct {
	Message string
	Err     error
}

// Error returns the error message
func (e *UnauthorizedError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Err.Error())
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *UnauthorizedError) Unwrap() error {
	return e.Err
}
