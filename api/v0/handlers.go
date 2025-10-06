package v0

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/api/v0/devices"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service devices.Service
}

func NewHandler(srv devices.Service) *Handler {
	return &Handler{
		service: srv,
	}
}

func (h *Handler) Register(r chi.Router) {
	handler := devices.New(h.service)
	handler.Register(r)
	r.Get("/health", health)
}
