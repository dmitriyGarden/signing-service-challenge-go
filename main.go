package main

import (
	"log"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/app"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/config"
)

func main() {
	cfg := config.Load()

	server := app.NewServer(cfg.ListenAddress)

	if err := server.Run(); err != nil {
		log.Fatalf("could not start server on %s: %v", cfg.ListenAddress, err)
	}
}
