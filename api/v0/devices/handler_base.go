package devices

import (
	"context"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	appdevices "github.com/fiskaly/coding-challenges/signing-service-challenge/internal/devices"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const devicesPathPrefix = "/api/v0/devices/"

var (
	_ Service = (*appdevices.Service)(nil)
	_ Service = (*appdevices.LoggingService)(nil)
)

// Service captures the contract used by HTTP handlers.
type Service interface {
	CreateDevice(ctx context.Context, input appdevices.CreateDeviceInput) (*appdevices.CreateDeviceResult, error)
	ListDevices(ctx context.Context) ([]domain.Device, error)
	GetDevice(ctx context.Context, id uuid.UUID) (domain.Device, error)
	UpdateDeviceLabel(ctx context.Context, id uuid.UUID, label string) (domain.Device, error)
	DeleteDevice(ctx context.Context, id uuid.UUID) error
	SignTransaction(ctx context.Context, input appdevices.SignTransactionInput) (*appdevices.SignatureResult, error)
	GetCounters(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]uint64, error)
	ListSignatures(ctx context.Context, deviceID uuid.UUID) ([]appdevices.SignatureRecord, error)
	GetSignature(ctx context.Context, deviceID uuid.UUID, counter uint64) (appdevices.SignatureRecord, error)
}

// Handler manages device-related HTTP endpoints.
type Handler struct {
	service Service
}

// New constructs a device handler.
func New(service Service) *Handler {
	return &Handler{service: service}
}

// Register wires handler routes into the provided mux.
func (h *Handler) Register(r chi.Router) {
	r.Route("/devices", h.registerDevices)
}

func (h *Handler) registerDevices(r chi.Router) {
	r.Post("/", h.createDevice)
	r.Get("/", h.listDevices)
	r.Get("/{device_id}", h.getDevice)
	r.Put("/{device_id}", h.updateDevice)
	r.Delete("/{device_id}", h.deleteDevice)

	r.Post("/{device_id}/sign", h.signTransaction)
	r.Get("/{device_id}/signatures", h.listSignatures)
	r.Get("/{device_id}/signatures/{counter}", h.getSignature)
}

func (h *Handler) deviceID(r *http.Request) (uuid.UUID, error) {
	chi.URLParam(r, "device_id")
	deviceID, err := uuid.Parse(chi.URLParam(r, "device_id"))
	if err != nil {
		return uuid.Nil, domain.ErrInvalidDeviceID
	}
	return deviceID, nil
}
