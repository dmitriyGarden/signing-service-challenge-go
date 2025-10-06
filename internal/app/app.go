package app

import (
	"log"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	v0 "github.com/fiskaly/coding-challenges/signing-service-challenge/api/v0"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/devices"
	inmemory2 "github.com/fiskaly/coding-challenges/signing-service-challenge/internal/persistence/inmemory"
	crypto2 "github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/crypto"
)

// NewServer wires together application dependencies and returns a configured HTTP server.
func NewServer(listenAddress string) *api.Server {
	repository := inmemory2.NewDeviceRepository()
	keyStore := inmemory2.NewKeyStore()
	keyGenerator := crypto2.NewDefaultKeyGenerator()
	signerFactory := crypto2.NewSignerFactory()
	signatureStore := inmemory2.NewSignatureStore()

	coreService := devices.NewService(repository, keyStore, keyGenerator, signerFactory, signatureStore)
	loggingService := devices.NewLoggingService(coreService, func(event string, fields map[string]interface{}) {
		log.Printf("event=%s fields=%v", event, fields)
	})

	apiV0Handler := v0.NewHandler(loggingService)

	return api.NewServer(listenAddress, map[string]api.DeviceHandler{
		"/api/v0": apiV0Handler,
	})
}
