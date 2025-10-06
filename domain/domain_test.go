package domain_test

import (
	"testing"
	"time"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/google/uuid"
)

func TestBuildSecuredPayload(t *testing.T) {
	reference := []byte{0x01, 0x02}
	payload := domain.BuildSecuredPayload(42, "payload", reference)
	expected := "42_payload_AQI="
	if payload != expected {
		t.Fatalf("expected %q, got %q", expected, payload)
	}
}

func TestParseAlgorithm_Success(t *testing.T) {
	algorithm, err := domain.ParseAlgorithm("  rsa ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if algorithm != domain.AlgorithmRSA {
		t.Fatalf("expected RSA, got %s", algorithm)
	}
}

func TestParseAlgorithm_Invalid(t *testing.T) {
	_, err := domain.ParseAlgorithm("sha256")
	if err == nil {
		t.Fatal("expected error for unsupported algorithm")
	}
	if err != domain.ErrInvalidAlgorithm {
		t.Fatalf("expected ErrInvalidAlgorithm, got %v", err)
	}
}

func TestValidateAlgorithm(t *testing.T) {
	if err := domain.ValidateAlgorithm(domain.AlgorithmECDSA); err != nil {
		t.Fatalf("unexpected error for ECDSA: %v", err)
	}
	err := domain.ValidateAlgorithm(domain.Algorithm("HMAC"))
	if err != domain.ErrInvalidAlgorithm {
		t.Fatalf("expected ErrInvalidAlgorithm, got %v", err)
	}
}

func TestDeviceClone(t *testing.T) {
	now := time.Now().UTC()
	device := domain.Device{
		ID:        uuid.New(),
		Algorithm: domain.AlgorithmRSA,
		Label:     "original",
		CreatedAt: now,
		UpdatedAt: now,
	}
	clone := device.Clone()
	clone.Label = "modified"
	if device.Label != "original" {
		t.Fatalf("expected original label to remain, got %q", device.Label)
	}
	if clone.Label != "modified" {
		t.Fatalf("expected clone label to change, got %q", clone.Label)
	}
}

func TestDeviceWithLabel(t *testing.T) {
	base := time.Now().UTC()
	device := domain.Device{ID: uuid.New(), Algorithm: domain.AlgorithmECDSA, Label: "old", CreatedAt: base, UpdatedAt: base}
	newTime := base.Add(10 * time.Minute)
	updated := device.WithLabel("new-label", newTime)
	if updated.Label != "new-label" {
		t.Fatalf("expected label to update, got %q", updated.Label)
	}
	if !updated.UpdatedAt.Equal(newTime) {
		t.Fatalf("expected UpdatedAt %v, got %v", newTime, updated.UpdatedAt)
	}
	if updated.CreatedAt != device.CreatedAt {
		t.Fatalf("expected CreatedAt to remain %v, got %v", device.CreatedAt, updated.CreatedAt)
	}
}

func TestValidationError_Error(t *testing.T) {
	err := domain.ValidationError{Field: "field", Message: "missing"}
	expected := "validation failed on 'field': missing"
	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}
}
