package devices_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api/v0/devices"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	appdevices "github.com/fiskaly/coding-challenges/signing-service-challenge/internal/devices"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func newRouter(service devices.Service) http.Handler {
	r := chi.NewRouter()
	h := devices.New(service)
	h.Register(r)
	return r
}

func decodeResponse[T any](t *testing.T, body io.Reader) T {
	t.Helper()
	var wrapper struct {
		Data T `json:"data"`
	}
	if err := json.NewDecoder(body).Decode(&wrapper); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return wrapper.Data
}

func TestCreateDevice_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockDevicesService(ctrl)
	router := newRouter(svc)

	deviceID := uuid.New()
	svc.EXPECT().CreateDevice(gomock.Any(), appdevices.CreateDeviceInput{
		ID:        deviceID,
		Algorithm: domain.AlgorithmRSA,
		Label:     "Terminal",
	}).Return(&appdevices.CreateDeviceResult{Device: domain.Device{ID: deviceID, Algorithm: domain.AlgorithmRSA, Label: "Terminal"}}, nil)

	payload := map[string]string{
		"id":        deviceID.String(),
		"algorithm": "RSA",
		"label":     "Terminal",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/devices/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	type createResp struct {
		ID        string `json:"id"`
		Algorithm string `json:"algorithm"`
		Label     string `json:"label"`
		Counter   uint64 `json:"counter"`
	}
	resp := decodeResponse[createResp](t, w.Body)

	if resp.ID != deviceID.String() {
		t.Fatalf("expected id %s, got %s", deviceID, resp.ID)
	}
	if resp.Counter != 0 {
		t.Fatalf("expected counter 0, got %d", resp.Counter)
	}
}

func TestCreateDevice_InvalidPayload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockDevicesService(ctrl)
	router := newRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/devices/", bytes.NewReader([]byte("not-json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestListDevices_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockDevicesService(ctrl)
	router := newRouter(svc)

	deviceID := uuid.New()
	devicesList := []domain.Device{{
		ID:        deviceID,
		Algorithm: domain.AlgorithmECDSA,
		Label:     "POS",
	}}
	svc.EXPECT().ListDevices(gomock.Any()).Return(devicesList, nil)
	svc.EXPECT().GetCounters(gomock.Any(), []uuid.UUID{deviceID}).Return(map[uuid.UUID]uint64{deviceID: 5}, nil)

	req := httptest.NewRequest(http.MethodGet, "/devices/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	type listResp struct {
		ID        string `json:"id"`
		Algorithm string `json:"algorithm"`
		Label     string `json:"label"`
		Counter   uint64 `json:"counter"`
	}
	payload := decodeResponse[[]listResp](t, w.Body)
	if len(payload) != 1 {
		t.Fatalf("expected 1 device, got %d", len(payload))
	}
	if payload[0].Counter != 5 {
		t.Fatalf("expected counter 5, got %d", payload[0].Counter)
	}
}

func TestSignTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockDevicesService(ctrl)
	router := newRouter(svc)

	deviceID := uuid.New()
	svc.EXPECT().SignTransaction(gomock.Any(), appdevices.SignTransactionInput{DeviceID: deviceID, Data: "payload"}).Return(&appdevices.SignatureResult{
		Signature:    "sig",
		SignedData:   "signed",
		CounterValue: 1,
	}, nil)

	body, _ := json.Marshal(map[string]string{"data": "payload"})
	req := httptest.NewRequest(http.MethodPost, "/devices/"+deviceID.String()+"/sign", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	type signResp struct {
		Signature  string `json:"signature"`
		SignedData string `json:"signed_data"`
	}
	resp := decodeResponse[signResp](t, w.Body)
	if resp.Signature != "sig" {
		t.Fatalf("unexpected signature %q", resp.Signature)
	}
}

func TestGetSignature_InvalidCounter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockDevicesService(ctrl)
	router := newRouter(svc)

	deviceID := uuid.New()

	req := httptest.NewRequest(http.MethodGet, "/devices/"+deviceID.String()+"/signatures/not-a-number", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d", w.Code)
	}
}

func TestListSignatures_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockDevicesService(ctrl)
	router := newRouter(svc)

	deviceID := uuid.New()
	records := []appdevices.SignatureRecord{{
		Counter:    2,
		Signature:  "sig",
		SignedData: "payload",
		CreatedAt:  time.Unix(0, 0).UTC(),
	}}
	svc.EXPECT().ListSignatures(gomock.Any(), deviceID).Return(records, nil)

	req := httptest.NewRequest(http.MethodGet, "/devices/"+deviceID.String()+"/signatures", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	type signatureResp struct {
		Counter    uint64    `json:"counter"`
		Signature  string    `json:"signature"`
		SignedData string    `json:"signed_data"`
		CreatedAt  time.Time `json:"created_at"`
	}
	payload := decodeResponse[[]signatureResp](t, w.Body)
	if len(payload) != 1 || payload[0].Counter != 2 {
		t.Fatalf("unexpected payload: %+v", payload)
	}
}
