package mocks

import (
	"context"
	"reflect"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api/v0/devices"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	appdevices "github.com/fiskaly/coding-challenges/signing-service-challenge/internal/devices"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

// MockDevicesService mocks the API layer service contract.
type MockDevicesService struct {
	ctrl     *gomock.Controller
	recorder *MockDevicesServiceMockRecorder
}

var _ devices.Service = (*MockDevicesService)(nil)

// MockDevicesServiceMockRecorder records invocations on MockDevicesService.
type MockDevicesServiceMockRecorder struct {
	mock *MockDevicesService
}

// NewMockDevicesService constructs a new mock instance.
func NewMockDevicesService(ctrl *gomock.Controller) *MockDevicesService {
	mock := &MockDevicesService{ctrl: ctrl}
	mock.recorder = &MockDevicesServiceMockRecorder{mock}
	return mock
}

// EXPECT exposes the recorder for expectation setup.
func (m *MockDevicesService) EXPECT() *MockDevicesServiceMockRecorder {
	return m.recorder
}

// CreateDevice mocks the base method.
func (m *MockDevicesService) CreateDevice(ctx context.Context, input appdevices.CreateDeviceInput) (*appdevices.CreateDeviceResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDevice", ctx, input)
	ret0, _ := ret[0].(*appdevices.CreateDeviceResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDevice records an expected call.
func (mr *MockDevicesServiceMockRecorder) CreateDevice(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDevice", reflect.TypeOf((*MockDevicesService)(nil).CreateDevice), ctx, input)
}

// ListDevices mocks the base method.
func (m *MockDevicesService) ListDevices(ctx context.Context) ([]domain.Device, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListDevices", ctx)
	ret0, _ := ret[0].([]domain.Device)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListDevices records an expected call.
func (mr *MockDevicesServiceMockRecorder) ListDevices(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListDevices", reflect.TypeOf((*MockDevicesService)(nil).ListDevices), ctx)
}

// GetDevice mocks the base method.
func (m *MockDevicesService) GetDevice(ctx context.Context, id uuid.UUID) (domain.Device, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDevice", ctx, id)
	ret0, _ := ret[0].(domain.Device)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDevice records an expected call.
func (mr *MockDevicesServiceMockRecorder) GetDevice(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDevice", reflect.TypeOf((*MockDevicesService)(nil).GetDevice), ctx, id)
}

// UpdateDeviceLabel mocks the base method.
func (m *MockDevicesService) UpdateDeviceLabel(ctx context.Context, id uuid.UUID, label string) (domain.Device, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDeviceLabel", ctx, id, label)
	ret0, _ := ret[0].(domain.Device)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateDeviceLabel records an expected call.
func (mr *MockDevicesServiceMockRecorder) UpdateDeviceLabel(ctx, id, label interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDeviceLabel", reflect.TypeOf((*MockDevicesService)(nil).UpdateDeviceLabel), ctx, id, label)
}

// DeleteDevice mocks the base method.
func (m *MockDevicesService) DeleteDevice(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteDevice", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteDevice records an expected call.
func (mr *MockDevicesServiceMockRecorder) DeleteDevice(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDevice", reflect.TypeOf((*MockDevicesService)(nil).DeleteDevice), ctx, id)
}

// SignTransaction mocks the base method.
func (m *MockDevicesService) SignTransaction(ctx context.Context, input appdevices.SignTransactionInput) (*appdevices.SignatureResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignTransaction", ctx, input)
	ret0, _ := ret[0].(*appdevices.SignatureResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignTransaction records an expected call.
func (mr *MockDevicesServiceMockRecorder) SignTransaction(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignTransaction", reflect.TypeOf((*MockDevicesService)(nil).SignTransaction), ctx, input)
}

// GetCounters mocks the base method.
func (m *MockDevicesService) GetCounters(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCounters", ctx, ids)
	ret0, _ := ret[0].(map[uuid.UUID]uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCounters records an expected call.
func (mr *MockDevicesServiceMockRecorder) GetCounters(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCounters", reflect.TypeOf((*MockDevicesService)(nil).GetCounters), ctx, ids)
}

// ListSignatures mocks the base method.
func (m *MockDevicesService) ListSignatures(ctx context.Context, deviceID uuid.UUID) ([]appdevices.SignatureRecord, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSignatures", ctx, deviceID)
	ret0, _ := ret[0].([]appdevices.SignatureRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSignatures records an expected call.
func (mr *MockDevicesServiceMockRecorder) ListSignatures(ctx, deviceID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSignatures", reflect.TypeOf((*MockDevicesService)(nil).ListSignatures), ctx, deviceID)
}

// GetSignature mocks the base method.
func (m *MockDevicesService) GetSignature(ctx context.Context, deviceID uuid.UUID, counter uint64) (appdevices.SignatureRecord, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSignature", ctx, deviceID, counter)
	ret0, _ := ret[0].(appdevices.SignatureRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSignature records an expected call.
func (mr *MockDevicesServiceMockRecorder) GetSignature(ctx, deviceID, counter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSignature", reflect.TypeOf((*MockDevicesService)(nil).GetSignature), ctx, deviceID, counter)
}
