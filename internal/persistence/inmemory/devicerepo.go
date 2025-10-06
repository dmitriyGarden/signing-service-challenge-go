package inmemory

import (
	"context"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/devices"
	"github.com/google/uuid"
)

// DeviceRepository stores devices in process memory with mutex protection.
type DeviceRepository struct {
	mu      sync.RWMutex
	devices map[uuid.UUID]domain.Device
}

var _ devices.Repository = (*DeviceRepository)(nil)

// NewDeviceRepository constructs a repository ready for use.
func NewDeviceRepository() *DeviceRepository {
	return &DeviceRepository{
		devices: make(map[uuid.UUID]domain.Device),
	}
}

// Create inserts a new device if it does not already exist.
func (r *DeviceRepository) Create(_ context.Context, device domain.Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.devices[device.ID]; exists {
		return domain.ErrDeviceExists
	}

	r.devices[device.ID] = device.Clone()
	return nil
}

// Get retrieves a device by its identifier.
func (r *DeviceRepository) Get(_ context.Context, id uuid.UUID) (domain.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	device, exists := r.devices[id]
	if !exists {
		return domain.Device{}, domain.NotFoundError{Resource: "device", ID: id.String()}
	}

	return device.Clone(), nil
}

// List returns all stored devices in unspecified order.
func (r *DeviceRepository) List(_ context.Context) ([]domain.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]domain.Device, 0, len(r.devices))
	for _, device := range r.devices {
		result = append(result, device.Clone())
	}

	return result, nil
}

// Update replaces the stored device state.
func (r *DeviceRepository) Update(_ context.Context, device domain.Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.devices[device.ID]; !exists {
		return domain.NotFoundError{Resource: "device", ID: device.ID.String()}
	}

	r.devices[device.ID] = device.Clone()
	return nil
}

// Delete removes a device from storage.
func (r *DeviceRepository) Delete(_ context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.devices, id)
	return nil
}
