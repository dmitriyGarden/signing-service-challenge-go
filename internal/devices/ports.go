package devices

import (
	"context"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/google/uuid"
)

// Repository abstracts persistence capabilities for devices.
type Repository interface {
	Create(ctx context.Context, device domain.Device) error
	Get(ctx context.Context, id uuid.UUID) (domain.Device, error)
	List(ctx context.Context) ([]domain.Device, error)
	Update(ctx context.Context, device domain.Device) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// KeyStore manages storing and retrieving key material independently of devices.
type KeyStore interface {
	Store(ctx context.Context, deviceID uuid.UUID, material domain.KeyMaterial) error
	Load(ctx context.Context, deviceID uuid.UUID) (domain.KeyMaterial, error)
	Delete(ctx context.Context, deviceID uuid.UUID) error
}

// Signer describes the minimal behaviour required from crypto signers.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

// SignerFactory resolves a signer implementation for a given device and key material.
type SignerFactory interface {
	SignerFor(device domain.Device, material domain.KeyMaterial) (Signer, error)
}

// KeyGenerator can produce key pairs for the configured algorithms.
type KeyGenerator interface {
	Generate(algorithm domain.Algorithm) (domain.KeyMaterial, error)
}

// SignatureStore persists signature records per device.
type SignatureStore interface {
	Append(ctx context.Context, deviceID uuid.UUID, record SignatureRecord) (SignatureRecord, error)
	List(ctx context.Context, deviceID uuid.UUID) ([]SignatureRecord, error)
	Get(ctx context.Context, deviceID uuid.UUID, counter uint64) (SignatureRecord, error)
	Last(ctx context.Context, deviceID uuid.UUID) (SignatureRecord, bool, error)
	GetCounters(ctx context.Context, deviceIDs []uuid.UUID) (map[uuid.UUID]uint64, error)
	Delete(ctx context.Context, deviceID uuid.UUID) error
}
