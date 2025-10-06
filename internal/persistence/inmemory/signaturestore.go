package inmemory

import (
	"context"
	"fmt"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/devices"
	"github.com/google/uuid"
)

// SignatureStore keeps signature history per device in memory.
type SignatureStore struct {
	mu      sync.RWMutex
	records map[uuid.UUID][]devices.SignatureRecord
}

var _ devices.SignatureStore = (*SignatureStore)(nil)

// NewSignatureStore creates an empty signature store.
func NewSignatureStore() *SignatureStore {
	return &SignatureStore{
		records: make(map[uuid.UUID][]devices.SignatureRecord),
	}
}

// Append adds a signature record for the given device and assigns a counter.
func (s *SignatureStore) Append(_ context.Context, deviceID uuid.UUID, record devices.SignatureRecord) (devices.SignatureRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cloned := record.Clone()
	cloned.Counter = uint64(len(s.records[deviceID]) + 1)
	s.records[deviceID] = append(s.records[deviceID], cloned)

	return cloned.Clone(), nil
}

// List returns all signature records for a device.
func (s *SignatureStore) List(_ context.Context, deviceID uuid.UUID) ([]devices.SignatureRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	records := s.records[deviceID]
	result := make([]devices.SignatureRecord, len(records))
	for i, record := range records {
		result[i] = record.Clone()
	}
	return result, nil
}

// Get retrieves a specific signature by counter.
func (s *SignatureStore) Get(_ context.Context, deviceID uuid.UUID, counter uint64) (devices.SignatureRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if counter == 0 {
		return devices.SignatureRecord{}, domain.NotFoundError{Resource: "signature", ID: fmt.Sprintf("%s#%d", deviceID.String(), counter)}
	}

	records := s.records[deviceID]
	if int(counter) > len(records) {
		return devices.SignatureRecord{}, domain.NotFoundError{Resource: "signature", ID: fmt.Sprintf("%s#%d", deviceID.String(), counter)}
	}

	return records[counter-1].Clone(), nil
}

func (s *SignatureStore) GetCounters(_ context.Context, deviceIDs []uuid.UUID) (map[uuid.UUID]uint64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make(map[uuid.UUID]uint64, len(deviceIDs))
	for _, deviceID := range deviceIDs {
		res[deviceID] = uint64(len(s.records[deviceID]))
	}
	return res, nil
}

// Last returns the most recent signature record for a device.
func (s *SignatureStore) Last(_ context.Context, deviceID uuid.UUID) (devices.SignatureRecord, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	records := s.records[deviceID]
	if len(records) == 0 {
		return devices.SignatureRecord{}, false, nil
	}

	return records[len(records)-1].Clone(), true, nil
}

func (s *SignatureStore) Delete(_ context.Context, deviceID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.records, deviceID)
	return nil
}
