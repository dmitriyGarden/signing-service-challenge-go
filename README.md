# Signing Service Challenge

## Overview
This service offers an HTTP API for managing signature devices and signing payloads using RSA or ECDSA keys. The system is structured into modular layers (API handlers, device service, crypto providers, persistence) and is accompanied by unit and integration tests.

> Development note: an AI coding assistant (Codex) collaborated on implementing features, documentation, and tests in this repository.

## Quick Start
- Install Go 1.20+
- Run all tests: `make test`
- Optional: `make race` to exercise the suite with the race detector.
- Start the API locally: `make run` (listens on `:8080`).
- Regenerate gomock doubles after interface changes: `make mocks`.

### Configuration
- `LISTEN_ADDRESS` – override the default `:8080` listen address for the HTTP server.

## API Highlights
- `POST /api/v0/devices` — create a device (`algorithm` must be `rsa` or `ecdsa`)
- `GET /api/v0/devices` — list devices
- `GET /api/v0/devices/{id}` — fetch a device
- `PATCH /api/v0/devices/{id}` — update label
- `DELETE /api/v0/devices/{id}` — delete device
- `POST /api/v0/devices/{id}/sign` — sign payload; response includes signature and secured data
- `GET /api/v0/devices/{id}/signatures` — retrieve signature history for a device
- `GET /api/v0/devices/{id}/signatures/{counter}` — fetch a specific signature by counter value

Device retrieval endpoints embed the current signature counter and last signature reference, computed from the signature history.

Refer to `api/tests/integration_test.go` for sample request/response bodies.

## Testing Strategy
- Unit tests: `internal/devices/service_test.go` (device service business logic with gomock) and `api/v0/devices/handler_test.go` (transport validation and error mapping).
- Integration tests: `api/tests/integration_test.go` spin up the HTTP router with in-memory adapters to exercise lifecycle and signing flows end-to-end.
- Concurrency checks: `internal/devices/concurrency_test.go` and the integration concurrency case ensure signature counters stay monotonic under parallel requests.

Run `make test` to execute the full suite.

## Architecture Notes
- Dependency wiring lives in `internal/app/app.go`.
- Crypto implementations and key generation reside in `crypto/`.
- In-memory persistence resides in `persistence/` and satisfies service ports defined in `internal/devices/ports.go`.

For a deeper breakdown, see `docs/ARCHITECTURE.md`.
