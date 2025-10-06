package devices

import (
	"errors"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api/v0/utils"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

func writeDomainError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case domain.ValidationError:
		writeErrorResponse(w, http.StatusUnprocessableEntity, []string{e.Error()})
	case domain.NotFoundError:
		writeErrorResponse(w, http.StatusNotFound, []string{e.Error()})
	case domain.ConflictError:
		writeErrorResponse(w, http.StatusConflict, []string{e.Error()})
	case domain.InternalError:
		writeErrorResponse(w, http.StatusInternalServerError, []string{e.Error()})
	default:
		if errors.Is(err, domain.ErrInvalidAlgorithm) {
			writeErrorResponse(w, http.StatusUnprocessableEntity, []string{err.Error()})
			return
		}
		if errors.Is(err, domain.ErrDeviceExists) {
			writeErrorResponse(w, http.StatusConflict, []string{err.Error()})
			return
		}
		if errors.Is(err, domain.ErrKeyMaterialMissing) {
			writeErrorResponse(w, http.StatusInternalServerError, []string{err.Error()})
		}
		writeInternalError(w)
	}
}

func writeErrorResponse(w http.ResponseWriter, code int, errors []string) {
	utils.WriteErrorResponse(w, code, errors)
}

func writeErrorsResponse(w http.ResponseWriter, code int, errors []error) {
	res := make([]string, 0, len(errors))
	for _, err := range errors {
		if err != nil {
			res = append(res, err.Error())
		}
	}
	writeErrorResponse(w, code, res)
}

func writeAPIResponse(w http.ResponseWriter, code int, data interface{}) {
	utils.WriteAPIResponse(w, code, data)
}

func writeInternalError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
