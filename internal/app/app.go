package app

import (
	"log"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	v0 "github.com/fiskaly/coding-challenges/signing-service-challenge/api/v0"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/devices"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/persistence/inmemory"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/crypto"
)

// NewServer wires together application dependencies and returns a configured HTTP server.
func NewServer(listenAddress string) *api.Server {
	repository := inmemory.NewDeviceRepository()
	keyStore := inmemory.NewKeyStore()
	keyGenerator := crypto.NewDefaultKeyGenerator()
	signerFactory := crypto.NewSignerFactory()
	signatureStore := inmemory.NewSignatureStore()

	coreService := devices.NewService(repository, keyStore, keyGenerator, signerFactory, signatureStore)
	loggingService := devices.NewLoggingService(coreService, func(event string, fields map[string]interface{}) {
		log.Printf("event=%s fields=%v", event, fields)
	})

	apiV0Handler := v0.NewHandler(loggingService)

	return api.NewServer(listenAddress, map[string]api.DeviceHandler{
		"/api/v0": apiV0Handler,
	})
}
