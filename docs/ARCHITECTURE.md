# Service Architecture

## Overview
The service follows a layered structure that separates HTTP transport, domain rules, cryptography, and persistence. Request handlers translate HTTP payloads into domain calls, while the domain layer encapsulates signature device behavior, ensuring that signature counters remain consistent and the signing process stays reusable across algorithms and storage backends.

## Domain Layer
- `domain.Device` holds canonical device state (`ID`, `Algorithm`, `Label`, timestamps) and exposes immutable update helpers. Signature counters are derived dynamically via the signature store.
- `domain.Algorithm`, `domain.ValidateAlgorithm`, and `domain.ParseAlgorithm` centralise validation for supported algorithms (`RSA`, `ECDSA`).
- `domain.BuildSecuredPayload` composes the `<counter>_<payload>_<reference>` string used for signing, ensuring consistent behaviour across service implementations.
- Error types (`ValidationError`, `NotFoundError`, `ConflictError`, `InternalError`) convey failure semantics without binding to transport concerns.
- `internal/devices.SignatureRecord` captures stored signature metadata (counter, signature, signed payload, timestamp) for retrieval endpoints.

## Persistence Layer
- `internal/devices.Repository` and `internal/devices.KeyStore` describe the storage ports. The default in-memory implementations (`persistence.InMemoryDeviceRepository`, `persistence.InMemoryKeyStore`) satisfy them with `sync.RWMutex`-guarded maps.
- `internal/devices.SignatureStore` abstracts signature history. `persistence.InMemorySignatureStore` implements it with append-only slices and counter lookup maps.
- Repository methods return typed domain errors for duplicates and missing IDs, while `SignatureStore` guarantees sequential counters.

## Crypto Layer
- `crypto.DefaultKeyGenerator` implements `internal/devices.KeyGenerator`, emitting PEM-encoded `domain.KeyMaterial` for RSA and ECDSA pairs.
- `crypto.SignerFactory` implements `internal/devices.SignerFactory`, decoding private keys and returning algorithm-specific signers (`RSASigner`, `ECDSASigner`).
- Signers normalise on SHA-256 hashing and output raw signature bytes for the service to base64-encode.

## Application Layer
- `internal/devices.Service` orchestrates device workflows (create, list, update label, delete, sign). It validates input, coordinates persistence, and ensures counters advance monotonically before persisting signatures.
- `internal/devices.LoggingService` decorates the core service with optional structured logging hooks.
- `internal/app.NewServer` is the composition root: it wires repositories, keystore, crypto providers, services, logging decorator, and HTTP handlers.
- `internal/config` centralises environment-driven settings (e.g. `LISTEN_ADDRESS`) that are loaded before the server bootstraps.

## HTTP Transport
- `api/server.go` configures the HTTP mux, registering the health endpoint and delegating device routes to `api/v0/devices.Handler`.
- `api/v0/devices.Handler` owns JSON validation, error translation, and response envelopes for `/api/v0/devices` CRUD operations and the `/sign` action.
- Additional endpoints (`GET /api/v0/devices/{id}/signatures`, `GET /api/v0/devices/{id}/signatures/{counter}`) expose signature history backed by the domain service.
- Typed domain errors are mapped to `422` (validation), `404` (missing devices), `409` (conflicts), or `500` (unexpected issues), while successful responses follow a `{ "data": ... }` convention.

## Cross-Cutting Concerns
- Logging remains opt-in via the service decorator, keeping the core logic oblivious to `log.Printf` or future tracing frameworks.
- Composition centralised in `internal/app` simplifies testing (dependency injection) and upcoming backend swaps.
- Error typing keeps HTTP, CLI, or gRPC frontends consistent—new transports can read the same error taxonomy without string matching.
- Integration confidence comes from gomock-driven unit tests and the `api/tests` suite, which exercises the router with in-memory adapters (including a parallel signing stress case).
