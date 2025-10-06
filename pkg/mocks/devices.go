package mocks

import (
	"context"
	"reflect"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/devices"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

// MockRepository implements devices.Repository for testing.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder records invocations on MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository constructs a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT exposes the recorder for expectation setup.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Create mocks the base method.
func (m *MockRepository) Create(ctx context.Context, device domain.Device) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, device)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create records an expected call.
func (mr *MockRepositoryMockRecorder) Create(ctx, device interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository)(nil).Create), ctx, device)
}

// Get mocks the base method.
func (m *MockRepository) Get(ctx context.Context, id uuid.UUID) (domain.Device, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(domain.Device)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get records an expected call.
func (mr *MockRepositoryMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRepository)(nil).Get), ctx, id)
}

// List mocks the base method.
func (m *MockRepository) List(ctx context.Context) ([]domain.Device, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx)
	ret0, _ := ret[0].([]domain.Device)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List records an expected call.
func (mr *MockRepositoryMockRecorder) List(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockRepository)(nil).List), ctx)
}

// Update mocks the base method.
func (m *MockRepository) Update(ctx context.Context, device domain.Device) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, device)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update records an expected call.
func (mr *MockRepositoryMockRecorder) Update(ctx, device interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepository)(nil).Update), ctx, device)
}

// Delete mocks the base method.
func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete records an expected call.
func (mr *MockRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepository)(nil).Delete), ctx, id)
}

// MockKeyStore implements devices.KeyStore for testing.
type MockKeyStore struct {
	ctrl     *gomock.Controller
	recorder *MockKeyStoreMockRecorder
}

// MockKeyStoreMockRecorder records invocations on MockKeyStore.
type MockKeyStoreMockRecorder struct {
	mock *MockKeyStore
}

// NewMockKeyStore constructs a new mock instance.
func NewMockKeyStore(ctrl *gomock.Controller) *MockKeyStore {
	mock := &MockKeyStore{ctrl: ctrl}
	mock.recorder = &MockKeyStoreMockRecorder{mock}
	return mock
}

// EXPECT exposes the recorder for expectation setup.
func (m *MockKeyStore) EXPECT() *MockKeyStoreMockRecorder {
	return m.recorder
}

// Store mocks the base method.
func (m *MockKeyStore) Store(ctx context.Context, deviceID uuid.UUID, material domain.KeyMaterial) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", ctx, deviceID, material)
	ret0, _ := ret[0].(error)
	return ret0
}

// Store records an expected call.
func (mr *MockKeyStoreMockRecorder) Store(ctx, deviceID, material interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockKeyStore)(nil).Store), ctx, deviceID, material)
}

// Load mocks the base method.
func (m *MockKeyStore) Load(ctx context.Context, deviceID uuid.UUID) (domain.KeyMaterial, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load", ctx, deviceID)
	ret0, _ := ret[0].(domain.KeyMaterial)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Load records an expected call.
func (mr *MockKeyStoreMockRecorder) Load(ctx, deviceID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockKeyStore)(nil).Load), ctx, deviceID)
}

// Delete mocks the base method.
func (m *MockKeyStore) Delete(ctx context.Context, deviceID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, deviceID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete records an expected call.
func (mr *MockKeyStoreMockRecorder) Delete(ctx, deviceID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockKeyStore)(nil).Delete), ctx, deviceID)
}

// MockSigner implements devices.Signer for testing.
type MockSigner struct {
	ctrl     *gomock.Controller
	recorder *MockSignerMockRecorder
}

// MockSignerMockRecorder records invocations on MockSigner.
type MockSignerMockRecorder struct {
	mock *MockSigner
}

// NewMockSigner constructs a new mock instance.
func NewMockSigner(ctrl *gomock.Controller) *MockSigner {
	mock := &MockSigner{ctrl: ctrl}
	mock.recorder = &MockSignerMockRecorder{mock}
	return mock
}

// EXPECT exposes the recorder for expectation setup.
func (m *MockSigner) EXPECT() *MockSignerMockRecorder {
	return m.recorder
}

// Sign mocks the base method.
func (m *MockSigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sign", dataToBeSigned)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sign records an expected call.
func (mr *MockSignerMockRecorder) Sign(dataToBeSigned interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sign", reflect.TypeOf((*MockSigner)(nil).Sign), dataToBeSigned)
}

// MockSignerFactory implements devices.SignerFactory for testing.
type MockSignerFactory struct {
	ctrl     *gomock.Controller
	recorder *MockSignerFactoryMockRecorder
}

// MockSignerFactoryMockRecorder records invocations on MockSignerFactory.
type MockSignerFactoryMockRecorder struct {
	mock *MockSignerFactory
}

// NewMockSignerFactory constructs a new mock instance.
func NewMockSignerFactory(ctrl *gomock.Controller) *MockSignerFactory {
	mock := &MockSignerFactory{ctrl: ctrl}
	mock.recorder = &MockSignerFactoryMockRecorder{mock}
	return mock
}

// EXPECT exposes the recorder for expectation setup.
func (m *MockSignerFactory) EXPECT() *MockSignerFactoryMockRecorder {
	return m.recorder
}

// SignerFor mocks the base method.
func (m *MockSignerFactory) SignerFor(device domain.Device, material domain.KeyMaterial) (devices.Signer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignerFor", device, material)
	ret0, _ := ret[0].(devices.Signer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignerFor records an expected call.
func (mr *MockSignerFactoryMockRecorder) SignerFor(device, material interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignerFor", reflect.TypeOf((*MockSignerFactory)(nil).SignerFor), device, material)
}

// MockKeyGenerator implements devices.KeyGenerator for testing.
type MockKeyGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockKeyGeneratorMockRecorder
}

// MockKeyGeneratorMockRecorder records invocations on MockKeyGenerator.
type MockKeyGeneratorMockRecorder struct {
	mock *MockKeyGenerator
}

// NewMockKeyGenerator constructs a new mock instance.
func NewMockKeyGenerator(ctrl *gomock.Controller) *MockKeyGenerator {
	mock := &MockKeyGenerator{ctrl: ctrl}
	mock.recorder = &MockKeyGeneratorMockRecorder{mock}
	return mock
}

// EXPECT exposes the recorder for expectation setup.
func (m *MockKeyGenerator) EXPECT() *MockKeyGeneratorMockRecorder {
	return m.recorder
}

// Generate mocks the base method.
func (m *MockKeyGenerator) Generate(algorithm domain.Algorithm) (domain.KeyMaterial, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate", algorithm)
	ret0, _ := ret[0].(domain.KeyMaterial)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Generate records an expected call.
func (mr *MockKeyGeneratorMockRecorder) Generate(algorithm interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockKeyGenerator)(nil).Generate), algorithm)
}

// MockSignatureStore implements devices.SignatureStore for testing.
type MockSignatureStore struct {
	ctrl     *gomock.Controller
	recorder *MockSignatureStoreMockRecorder
}

// MockSignatureStoreMockRecorder records invocations on MockSignatureStore.
type MockSignatureStoreMockRecorder struct {
	mock *MockSignatureStore
}

// NewMockSignatureStore constructs a new mock instance.
func NewMockSignatureStore(ctrl *gomock.Controller) *MockSignatureStore {
	mock := &MockSignatureStore{ctrl: ctrl}
	mock.recorder = &MockSignatureStoreMockRecorder{mock}
	return mock
}

// EXPECT exposes the recorder for expectation setup.
func (m *MockSignatureStore) EXPECT() *MockSignatureStoreMockRecorder {
	return m.recorder
}

// Append mocks the base method.
func (m *MockSignatureStore) Append(ctx context.Context, deviceID uuid.UUID, record devices.SignatureRecord) (devices.SignatureRecord, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Append", ctx, deviceID, record)
	ret0, _ := ret[0].(devices.SignatureRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Append records an expected call.
func (mr *MockSignatureStoreMockRecorder) Append(ctx, deviceID, record interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Append", reflect.TypeOf((*MockSignatureStore)(nil).Append), ctx, deviceID, record)
}

// List mocks the base method.
func (m *MockSignatureStore) List(ctx context.Context, deviceID uuid.UUID) ([]devices.SignatureRecord, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, deviceID)
	ret0, _ := ret[0].([]devices.SignatureRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List records an expected call.
func (mr *MockSignatureStoreMockRecorder) List(ctx, deviceID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockSignatureStore)(nil).List), ctx, deviceID)
}

// Get mocks the base method.
func (m *MockSignatureStore) Get(ctx context.Context, deviceID uuid.UUID, counter uint64) (devices.SignatureRecord, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, deviceID, counter)
	ret0, _ := ret[0].(devices.SignatureRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get records an expected call.
func (mr *MockSignatureStoreMockRecorder) Get(ctx, deviceID, counter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockSignatureStore)(nil).Get), ctx, deviceID, counter)
}

// Last mocks the base method.
func (m *MockSignatureStore) Last(ctx context.Context, deviceID uuid.UUID) (devices.SignatureRecord, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Last", ctx, deviceID)
	ret0, _ := ret[0].(devices.SignatureRecord)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Last records an expected call.
func (mr *MockSignatureStoreMockRecorder) Last(ctx, deviceID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Last", reflect.TypeOf((*MockSignatureStore)(nil).Last), ctx, deviceID)
}

// GetCounters mocks the base method.
func (m *MockSignatureStore) GetCounters(ctx context.Context, deviceIDs []uuid.UUID) (map[uuid.UUID]uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCounters", ctx, deviceIDs)
	ret0, _ := ret[0].(map[uuid.UUID]uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCounters records an expected call.
func (mr *MockSignatureStoreMockRecorder) GetCounters(ctx, deviceIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCounters", reflect.TypeOf((*MockSignatureStore)(nil).GetCounters), ctx, deviceIDs)
}

// Delete mocks the base method.
func (m *MockSignatureStore) Delete(ctx context.Context, deviceID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, deviceID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete records an expected call.
func (mr *MockSignatureStoreMockRecorder) Delete(ctx, deviceID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSignatureStore)(nil).Delete), ctx, deviceID)
}
