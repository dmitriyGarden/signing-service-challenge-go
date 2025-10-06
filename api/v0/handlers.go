package v0

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/api/v0/devices"
	"github.com/go-chi/chi/v5"
)

// Handler wires version specific routes.
type Handler struct {
	service devices.Service
}

// NewHandler creates an API v0 handler with the given service.
func NewHandler(srv devices.Service) *Handler {
	return &Handler{
		service: srv,
	}
}

// Register mounts the versioned device routes and health endpoint.
func (h *Handler) Register(r chi.Router) {
	handler := devices.New(h.service)
	handler.Register(r)
	r.Get("/health", health)
}
