package v0

import (
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api/v0/utils"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

func health(response http.ResponseWriter, _ *http.Request) {
	utils.WriteAPIResponse(response, http.StatusOK, HealthResponse{
		Status:  "pass",
		Version: "v0",
	})
}
