package domain

import "fmt"

// ValidationError represents invalid user-supplied data.
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	if e.Field == "" {
		return fmt.Sprintf("validation failed: %s", e.Message)
	}
	return fmt.Sprintf("validation failed on '%s': %s", e.Field, e.Message)
}

// NotFoundError wraps missing resource scenarios.
type NotFoundError struct {
	Resource string
	ID       string
}

func (e NotFoundError) Error() string {
	if e.Resource == "" {
		return "resource not found"
	}
	if e.ID == "" {
		return fmt.Sprintf("%s not found", e.Resource)
	}
	return fmt.Sprintf("%s '%s' not found", e.Resource, e.ID)
}

// ConflictError signals state conflicts such as duplicate IDs or counter regressions.
type ConflictError struct {
	Reason string
}

func (e ConflictError) Error() string {
	if e.Reason == "" {
		return "conflict"
	}
	return e.Reason
}

// InternalError indicates server-side issues; wraps root cause but hides specifics.
type InternalError struct {
	Reason string
}

func (e InternalError) Error() string {
	if e.Reason == "" {
		return "internal error"
	}
	return e.Reason
}

var (
	ErrInvalidAlgorithm   = ValidationError{Field: "algorithm", Message: "unsupported algorithm"}
	ErrInvalidDeviceID    = ValidationError{Field: "id", Message: "device ID must be a valid UUID"}
	ErrDeviceNotFound     = NotFoundError{Resource: "device"}
	ErrDeviceExists       = ConflictError{Reason: "device already exists"}
	ErrKeyMaterialMissing = InternalError{Reason: "key material missing"}
)
