package devices

import (
	"time"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/google/uuid"
)

type createDeviceRequest struct {
	ID        string `json:"id"`
	Algorithm string `json:"algorithm"`
	Label     string `json:"label"`
}

func (c *createDeviceRequest) Validate() []error {
	errs := make([]error, 0)
	_, err := domain.ParseAlgorithm(c.Algorithm)
	if err != nil {
		errs = append(errs, err)
	}

	_, err = uuid.Parse(c.ID)
	if err != nil {
		errs = append(errs, domain.ErrInvalidDeviceID)
	}
	return errs
}

type devicePayload struct {
	ID        string `json:"id"`
	Algorithm string `json:"algorithm"`
	Label     string `json:"label"`
	Counter   uint64 `json:"counter"`
}

type updateDeviceRequest struct {
	Label string `json:"label"`
}

func (c *updateDeviceRequest) Validate() []error {
	return nil
}

type signRequest struct {
	Data string `json:"data"`
}

func (c *signRequest) Validate() []error {
	errs := make([]error, 0)

	return errs
}

type signResponse struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
}

type signaturePayload struct {
	Counter    uint64    `json:"counter"`
	Signature  string    `json:"signature"`
	SignedData string    `json:"signed_data"`
	CreatedAt  time.Time `json:"created_at"`
}
