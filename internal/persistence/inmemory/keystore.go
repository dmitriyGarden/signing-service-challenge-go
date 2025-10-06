package inmemory

import (
	"context"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/devices"
	"github.com/google/uuid"
)

// InMemoryKeyStore keeps key material in RAM keyed by device ID.
type KeyStore struct {
	mu    sync.RWMutex
	store map[uuid.UUID]domain.KeyMaterial
}

var _ devices.KeyStore = (*KeyStore)(nil)

// NewKeyStore instantiates the in-memory key store.
func NewKeyStore() *KeyStore {
	return &KeyStore{
		store: make(map[uuid.UUID]domain.KeyMaterial),
	}
}

// Store persists key material for a device.
func (k *KeyStore) Store(_ context.Context, deviceID uuid.UUID, material domain.KeyMaterial) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	copyMaterial := domain.KeyMaterial{
		Public:  append([]byte(nil), material.Public...),
		Private: append([]byte(nil), material.Private...),
	}

	k.store[deviceID] = copyMaterial
	return nil
}

// Load fetches key material for a device.
func (k *KeyStore) Load(_ context.Context, deviceID uuid.UUID) (domain.KeyMaterial, error) {
	k.mu.RLock()
	defer k.mu.RUnlock()

	material, exists := k.store[deviceID]
	if !exists {
		return domain.KeyMaterial{}, domain.ErrKeyMaterialMissing
	}

	return domain.KeyMaterial{
		Public:  append([]byte(nil), material.Public...),
		Private: append([]byte(nil), material.Private...),
	}, nil
}

// Delete removes key material from the store.
func (k *KeyStore) Delete(_ context.Context, deviceID uuid.UUID) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	delete(k.store, deviceID)
	return nil
}
