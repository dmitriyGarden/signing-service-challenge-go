package domain

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Algorithm enumerates supported signing algorithms.
type Algorithm string

// Supported signing algorithms.
const (
	AlgorithmRSA   Algorithm = "RSA"
	AlgorithmECDSA Algorithm = "ECDSA"
)

// Device represents a signature device managed by the service.
type Device struct {
	ID        uuid.UUID `json:"id"`
	Algorithm Algorithm `json:"algorithm"`
	Label     string    `json:"label"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Clone provides a deep copy to avoid leaking internal state.
// Clone returns a shallow copy of the device to avoid exposing internal references.
func (d Device) Clone() Device {
	clone := d
	return clone
}

// WithLabel returns a new copy of the device with an updated label and timestamp.
func (d Device) WithLabel(label string, at time.Time) Device {
	clone := d.Clone()
	clone.Label = label
	clone.UpdatedAt = at
	return clone
}

// KeyMaterial holds the serialized public and private keys for a device.
type KeyMaterial struct {
	Public  []byte // PEM encoded public key material.
	Private []byte // PEM encoded private key material.
}

// BuildSecuredPayload composes the string to be signed following domain rules.
func BuildSecuredPayload(counter uint64, data string, reference []byte) string {
	encoded := base64.StdEncoding.EncodeToString(reference)
	return fmt.Sprintf("%d_%s_%s", counter, data, encoded)
}

// ParseAlgorithm converts an external string into a supported Algorithm.
func ParseAlgorithm(value string) (Algorithm, error) {
	normalized := strings.ToUpper(strings.TrimSpace(value))
	algorithm := Algorithm(normalized)
	if err := ValidateAlgorithm(algorithm); err != nil {
		return "", err
	}
	return algorithm, nil
}

// ValidateAlgorithm ensures the provided algorithm is supported.
func ValidateAlgorithm(algorithm Algorithm) error {
	switch algorithm {
	case AlgorithmRSA, AlgorithmECDSA:
		return nil
	default:
		return ErrInvalidAlgorithm
	}
}
