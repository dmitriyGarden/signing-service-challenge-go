package devices

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	appdevices "github.com/fiskaly/coding-challenges/signing-service-challenge/internal/devices"
	"github.com/google/uuid"
)

// createDevice provisions a new device and returns its metadata.
func (h *Handler) createDevice(w http.ResponseWriter, r *http.Request) {
	var request createDeviceRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, []string{"invalid request payload"})
		return
	}
	errs := request.Validate()
	if len(errs) > 0 {
		writeErrorsResponse(w, http.StatusBadRequest, errs)
		return
	}

	algorithm, err := domain.ParseAlgorithm(request.Algorithm)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	uid, err := uuid.Parse(request.ID)
	if err != nil {
		writeDomainError(w, domain.ErrInvalidDeviceID)
		return
	}

	result, err := h.service.CreateDevice(r.Context(), appdevices.CreateDeviceInput{
		ID:        uid,
		Algorithm: algorithm,
		Label:     request.Label,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeAPIResponse(w, http.StatusCreated, devicePayload{
		ID:        result.Device.ID.String(),
		Algorithm: string(result.Device.Algorithm),
		Label:     result.Device.Label,
		Counter:   0,
	})
}

// listDevices returns all registered devices.
func (h *Handler) listDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := h.service.ListDevices(r.Context())
	if err != nil {
		writeDomainError(w, err)
		return
	}
	payloads, err := h.makeDevicesPayload(r.Context(), devices)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeAPIResponse(w, http.StatusOK, payloads)
}

// getDevice fetches a device by ID.
func (h *Handler) getDevice(w http.ResponseWriter, r *http.Request) {
	id, err := h.deviceID(r)
	if err != nil {
		writeDomainError(w, err)
	}
	device, err := h.service.GetDevice(r.Context(), id)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	payloads, err := h.makeDevicesPayload(r.Context(), []domain.Device{device})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeAPIResponse(w, http.StatusOK, payloads[0])
}

// updateDevice updates a device's label.
func (h *Handler) updateDevice(w http.ResponseWriter, r *http.Request) {
	id, err := h.deviceID(r)
	if err != nil {
		writeDomainError(w, err)
	}
	var request updateDeviceRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, []string{"invalid request payload"})
		return
	}
	errs := request.Validate()
	if len(errs) > 0 {
		writeErrorsResponse(w, http.StatusBadRequest, errs)
		return
	}
	updated, err := h.service.UpdateDeviceLabel(r.Context(), id, request.Label)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	payloads, err := h.makeDevicesPayload(r.Context(), []domain.Device{updated})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeAPIResponse(w, http.StatusOK, payloads[0])
}

// deleteDevice removes a device and associated state.
func (h *Handler) deleteDevice(w http.ResponseWriter, r *http.Request) {
	id, err := h.deviceID(r)
	if err != nil {
		writeDomainError(w, err)
	}
	if err := h.service.DeleteDevice(r.Context(), id); err != nil {
		writeDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// makeDevicesPayload enriches devices with their current counters.
func (h *Handler) makeDevicesPayload(ctx context.Context, devices []domain.Device) ([]devicePayload, error) {
	ids := make([]uuid.UUID, len(devices))
	for i, device := range devices {
		ids[i] = device.ID
	}
	counters, err := h.service.GetCounters(ctx, ids)
	if err != nil {
		return nil, err
	}
	res := make([]devicePayload, len(devices))
	for i, device := range devices {
		res[i] = devicePayload{
			ID:        device.ID.String(),
			Algorithm: string(device.Algorithm),
			Label:     device.Label,
			Counter:   counters[device.ID],
		}
	}
	return res, nil
}
