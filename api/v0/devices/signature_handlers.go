package devices

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	appdevices "github.com/fiskaly/coding-challenges/signing-service-challenge/internal/devices"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) signTransaction(w http.ResponseWriter, r *http.Request) {
	id, err := h.deviceID(r)
	if err != nil {
		writeDomainError(w, err)
	}
	var request signRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, []string{"invalid request payload"})
		return
	}

	result, err := h.service.SignTransaction(r.Context(), appdevices.SignTransactionInput{
		DeviceID: id,
		Data:     request.Data,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeAPIResponse(w, http.StatusOK, signResponse{
		Signature:  result.Signature,
		SignedData: result.SignedData,
	})
}

func (h *Handler) listSignatures(w http.ResponseWriter, r *http.Request) {
	deviceID, err := h.deviceID(r)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	records, err := h.service.ListSignatures(r.Context(), deviceID)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	payloads := make([]signaturePayload, 0, len(records))
	for _, record := range records {
		payloads = append(payloads, signaturePayload{
			Counter:    record.Counter,
			Signature:  record.Signature,
			SignedData: record.SignedData,
			CreatedAt:  record.CreatedAt,
		})
	}

	writeAPIResponse(w, http.StatusOK, payloads)
}

func (h *Handler) getSignature(w http.ResponseWriter, r *http.Request) {
	deviceID, err := h.deviceID(r)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	counter := chi.URLParam(r, "counter")
	cnt, err := strconv.ParseUint(counter, 10, 64)
	if err != nil {
		writeDomainError(w, domain.ValidationError{
			Field:   "counter",
			Message: "counter must be a positive integer",
		})
		return
	}
	record, err := h.service.GetSignature(r.Context(), deviceID, cnt)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	payload := signaturePayload{
		Counter:    record.Counter,
		Signature:  record.Signature,
		SignedData: record.SignedData,
		CreatedAt:  record.CreatedAt,
	}

	writeAPIResponse(w, http.StatusOK, payload)
}
