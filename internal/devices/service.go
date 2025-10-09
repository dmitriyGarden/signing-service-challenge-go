package devices

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/google/uuid"
)

// Service encapsulates domain rules for managing devices and signatures.
type Service struct {
	repo           Repository
	keyStore       KeyStore
	keyGenerator   KeyGenerator
	signerFactory  SignerFactory
	signatureStore SignatureStore
	clock          func() time.Time
	signMX         sync.RWMutex // guards Append operations to keep signature counters monotonic
}

// NewService constructs a Service with injectable dependencies for testing.
func NewService(repo Repository, keyStore KeyStore, generator KeyGenerator, signerFactory SignerFactory, signatureStore SignatureStore) *Service {
	return &Service{
		repo:           repo,
		keyStore:       keyStore,
		keyGenerator:   generator,
		signerFactory:  signerFactory,
		signatureStore: signatureStore,
		clock:          time.Now,
	}
}

// WithClock allows overriding the clock function (mostly for tests).
func (s *Service) WithClock(clock func() time.Time) {
	if clock != nil {
		s.clock = clock
	}
}

// CreateDeviceInput captures user-provided data to create a new device.
type CreateDeviceInput struct {
	ID        uuid.UUID
	Algorithm domain.Algorithm
	Label     string
}

// CreateDeviceResult bundles the persisted device with its generated key material.
type CreateDeviceResult struct {
	Device domain.Device
}

// CreateDevice provisions a new signature device.
func (s *Service) CreateDevice(ctx context.Context, input CreateDeviceInput) (*CreateDeviceResult, error) {
	if s == nil {
		return nil, errors.New("device service is nil")
	}

	label := strings.TrimSpace(input.Label)

	keys, err := s.keyGenerator.Generate(input.Algorithm)
	if err != nil {
		return nil, fmt.Errorf("generate key pair: %w", err)
	}

	now := s.clock().UTC()
	device := domain.Device{
		ID:        input.ID,
		Algorithm: input.Algorithm,
		Label:     label,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, device); err != nil {
		return nil, err
	}

	if err := s.keyStore.Store(ctx, device.ID, keys); err != nil {
		_ = s.repo.Delete(ctx, device.ID)
		return nil, fmt.Errorf("store key material: %w", err)
	}

	return &CreateDeviceResult{
		Device: device,
	}, nil
}

// SignTransactionInput carries the parameters required to sign a payload.
type SignTransactionInput struct {
	DeviceID uuid.UUID
	Data     string
}

// SignatureResult represents the outcome of a signing operation.
type SignatureResult struct {
	Signature    string
	SignedData   string
	CounterValue uint64
}

// SignTransaction creates a signature for the given payload while keeping counters consistent.
func (s *Service) SignTransaction(ctx context.Context, input SignTransactionInput) (*SignatureResult, error) {
	if s == nil {
		return nil, errors.New("device service is nil")
	}

	if strings.TrimSpace(input.Data) == "" {
		return nil, domain.ValidationError{Field: "data", Message: "data is required"}
	}

	device, err := s.repo.Get(ctx, input.DeviceID)
	if err != nil {
		return nil, err
	}

	material, err := s.keyStore.Load(ctx, device.ID)
	if err != nil {
		return nil, fmt.Errorf("load key material: %w", err)
	}

	signer, err := s.signerFactory.SignerFor(device, material)
	if err != nil {
		return nil, fmt.Errorf("resolve signer: %w", err)
	}
	s.signMX.Lock()
	defer s.signMX.Unlock()
	prevRecord, found, err := s.signatureStore.Last(ctx, device.ID)
	if err != nil {
		return nil, err
	}

	var reference []byte
	var counter uint64
	if !found {
		reference = device.ID[:]
		counter = 0
	} else {
		decoded, decodeErr := base64.StdEncoding.DecodeString(prevRecord.Signature)
		if decodeErr != nil {
			return nil, fmt.Errorf("decode previous signature: %w", decodeErr)
		}
		reference = decoded
		counter = prevRecord.Counter
	}

	signedData := domain.BuildSecuredPayload(counter+1, input.Data, reference)

	signatureBytes, err := signer.Sign([]byte(signedData))
	if err != nil {
		return nil, fmt.Errorf("sign payload: %w", err)
	}

	encodedSignature := base64.StdEncoding.EncodeToString(signatureBytes)
	record := SignatureRecord{
		Signature:  encodedSignature,
		SignedData: signedData,
		CreatedAt:  s.clock().UTC(),
	}

	storedRecord, err := s.signatureStore.Append(ctx, device.ID, record)
	if err != nil {
		return nil, fmt.Errorf("append signature record: %w", err)
	}

	return &SignatureResult{
		Signature:    storedRecord.Signature,
		SignedData:   storedRecord.SignedData,
		CounterValue: storedRecord.Counter,
	}, nil
}

// UpdateDeviceLabel updates the display label of an existing device.
func (s *Service) UpdateDeviceLabel(ctx context.Context, id uuid.UUID, label string) (domain.Device, error) {
	if s == nil {
		return domain.Device{}, errors.New("device service is nil")
	}

	device, err := s.repo.Get(ctx, id)
	if err != nil {
		return domain.Device{}, err
	}

	updated := device.WithLabel(strings.TrimSpace(label), s.clock().UTC())
	if err := s.repo.Update(ctx, updated); err != nil {
		return domain.Device{}, err
	}

	return updated, nil
}

// GetDevice fetches a device by its identifier.
func (s *Service) GetDevice(ctx context.Context, id uuid.UUID) (domain.Device, error) {
	if s == nil {
		return domain.Device{}, errors.New("device service is nil")
	}

	device, err := s.repo.Get(ctx, id)
	if err != nil {
		return domain.Device{}, err
	}

	return device, nil
}

// GetCounters reports the current signature counter for each requested device ID.
func (s *Service) GetCounters(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]uint64, error) {
	if s == nil {
		return nil, errors.New("device service is nil")
	}
	s.signMX.RLock()
	defer s.signMX.RUnlock()
	counters, err := s.signatureStore.GetCounters(ctx, ids)
	if err != nil {
		return nil, err
	}

	return counters, nil
}

// ListDevices returns all devices known to the service.
func (s *Service) ListDevices(ctx context.Context) ([]domain.Device, error) {
	if s == nil {
		return nil, errors.New("device service is nil")
	}

	return s.repo.List(ctx)
}

// DeleteDevice removes a device and its key material.
func (s *Service) DeleteDevice(ctx context.Context, id uuid.UUID) error {
	if s == nil {
		return errors.New("device service is nil")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	if err := s.keyStore.Delete(ctx, id); err != nil {
		return err
	}
	s.signMX.Lock()
	defer s.signMX.Unlock()
	if err := s.signatureStore.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

// ListSignatures retrieves all signature records for a device.
func (s *Service) ListSignatures(ctx context.Context, deviceID uuid.UUID) ([]SignatureRecord, error) {
	if s == nil {
		return nil, errors.New("device service is nil")
	}

	return s.signatureStore.List(ctx, deviceID)
}

// GetSignature fetches a specific signature record by counter.
func (s *Service) GetSignature(ctx context.Context, deviceID uuid.UUID, counter uint64) (SignatureRecord, error) {
	if s == nil {
		return SignatureRecord{}, errors.New("device service is nil")
	}

	return s.signatureStore.Get(ctx, deviceID, counter)
}

// LoggingService decorates a Service with structured logging callbacks.
type LoggingService struct {
	inner  *Service
	logger func(event string, fields map[string]interface{})
}

// NewLoggingService returns a logging decorator for the provided service.
func NewLoggingService(inner *Service, logger func(event string, fields map[string]interface{})) *LoggingService {
	return &LoggingService{inner: inner, logger: logger}
}

func (l *LoggingService) log(event string, fields map[string]interface{}) {
	if l.logger != nil {
		l.logger(event, fields)
	}
}

// CreateDevice proxies to the wrapped service while emitting log hooks.
func (l *LoggingService) CreateDevice(ctx context.Context, input CreateDeviceInput) (*CreateDeviceResult, error) {
	l.log("device.create", map[string]interface{}{"id": input.ID, "algorithm": input.Algorithm})
	result, err := l.inner.CreateDevice(ctx, input)
	if err != nil {
		l.log("device.create.error", map[string]interface{}{"id": input.ID, "error": err.Error()})
	}
	return result, err
}

// SignTransaction proxies signing calls and adds log events.
func (l *LoggingService) SignTransaction(ctx context.Context, input SignTransactionInput) (*SignatureResult, error) {
	l.log("device.sign", map[string]interface{}{"id": input.DeviceID})
	result, err := l.inner.SignTransaction(ctx, input)
	if err != nil {
		l.log("device.sign.error", map[string]interface{}{"id": input.DeviceID, "error": err.Error()})
	}
	return result, err
}

// UpdateDeviceLabel wraps the service update call with logging.
func (l *LoggingService) UpdateDeviceLabel(ctx context.Context, id uuid.UUID, label string) (domain.Device, error) {
	l.log("device.update", map[string]interface{}{"id": id})
	device, err := l.inner.UpdateDeviceLabel(ctx, id, label)
	if err != nil {
		l.log("device.update.error", map[string]interface{}{"id": id, "error": err.Error()})
	}
	return device, err
}

// GetCounters records failures when retrieving signature counters.
func (l *LoggingService) GetCounters(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]uint64, error) {
	counters, err := l.inner.GetCounters(ctx, ids)
	if err != nil {
		l.log("device.get.error", map[string]interface{}{"ids": ids, "error": err.Error()})
	}
	return counters, err
}

// GetDevice logs retrieval attempts for individual devices.
func (l *LoggingService) GetDevice(ctx context.Context, id uuid.UUID) (domain.Device, error) {
	device, err := l.inner.GetDevice(ctx, id)
	if err != nil {
		l.log("device.get.error", map[string]interface{}{"id": id, "error": err.Error()})
	}
	return device, err
}

// ListDevices wraps ListDevices with error logging.
func (l *LoggingService) ListDevices(ctx context.Context) ([]domain.Device, error) {
	devices, err := l.inner.ListDevices(ctx)
	if err != nil {
		l.log("device.list.error", map[string]interface{}{"error": err.Error()})
	}
	return devices, err
}

// DeleteDevice logs delete attempts and errors.
func (l *LoggingService) DeleteDevice(ctx context.Context, id uuid.UUID) error {
	l.log("device.delete", map[string]interface{}{"id": id})
	err := l.inner.DeleteDevice(ctx, id)
	if err != nil {
		l.log("device.delete.error", map[string]interface{}{"id": id, "error": err.Error()})
	}
	return err
}

// ListSignatures logs list operations and failures.
func (l *LoggingService) ListSignatures(ctx context.Context, deviceID uuid.UUID) ([]SignatureRecord, error) {
	records, err := l.inner.ListSignatures(ctx, deviceID)
	if err != nil {
		l.log("signature.list.error", map[string]interface{}{"device_id": deviceID, "error": err.Error()})
	}
	return records, err
}

// GetSignature wraps the service call to fetch a specific signature.
func (l *LoggingService) GetSignature(ctx context.Context, deviceID uuid.UUID, counter uint64) (SignatureRecord, error) {
	record, err := l.inner.GetSignature(ctx, deviceID, counter)
	if err != nil {
		l.log("signature.get.error", map[string]interface{}{"device_id": deviceID, "counter": counter, "error": err.Error()})
	}
	return record, err
}
