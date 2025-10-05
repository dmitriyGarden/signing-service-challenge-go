# Service Architecture

## Overview
The service follows a layered structure that separates HTTP transport, domain rules, cryptography, and persistence. Request handlers translate HTTP payloads into domain calls, while the domain layer encapsulates signature device behavior, ensuring that signature counters remain consistent and the signing process stays reusable across algorithms and storage backends.

## Domain Layer
- `domain.Device` holds the persistent state: `ID string`, `Algorithm Algorithm`, `Label string`, `SignatureCounter uint64`, `LastSignature []byte`, `CreatedAt time.Time`, `UpdatedAt time.Time`. Helper methods guard incremental updates (`IncrementCounter`, `SetLabel`).
- `domain.Algorithm` enumerates supported algorithms (`AlgorithmRSA`, `AlgorithmECDSA`), validating input before device creation.
- `domain.DeviceService` orchestrates workflows. Dependencies:
  - `DeviceRepository`: CRUD access plus atomic counter updates (`Create`, `Get`, `List`, `Update`, `UpdateSignatureState`).
  - `SignerFactory`: returns a `Signer` for a device algorithm and keypair.
  - `KeyStore`: persists generated key pairs, enabling retrieval for signing.
- `domain.SignatureResult` returns the base64 signature and the secured payload string. `DeviceService.SignTransaction` ensures the counter increments only after the crypto call succeeds and persists the new `LastSignature`.

## Persistence Layer
- Default implementation: `persistence.InMemoryDeviceRepository` backed by `sync.RWMutex` and maps keyed by device ID. Counter updates use `sync.Mutex` to guarantee monotonic increments.
- `KeyStore` interface has `Store(deviceID string, publicKey, privateKey []byte)` and `Load(deviceID string) (publicKey, privateKey []byte, error)`. In-memory variant uses a separate map to isolate sensitive bytes from device metadata.
- Persistence interfaces allow swapping to SQL/NoSQL stores without changing domain logic.

## Crypto Layer
- `crypto.Signer` stays the abstraction. A factory (`SignerFactory` implementation) inspects `domain.Algorithm` and wires RSA/ECDSA signers by unmarshal­ing stored key material.
- RSA signer uses `crypto/rsa` with PKCS#1 v1.5 for signing; ECDSA signer uses `ecdsa.SignASN1`. Both return base64-encoded signatures.
- Key generation lives in `crypto/generation.go`; `domain.DeviceService.CreateDevice` invokes the corresponding generator via factory.

## HTTP Transport
- Router resides in `api/server.go`. Routes:
  - `GET /api/v0/health`
  - `POST /api/v0/devices`
  - `GET /api/v0/devices`
  - `GET /api/v0/devices/{id}`
  - `POST /api/v0/devices/{id}/sign`
- Handlers validate JSON payloads, invoke `DeviceService`, and return structured responses or errors (`ErrorResponse`). Validation failures respond with `422 Unprocessable Entity`.

## Cross-Cutting Concerns
- Errors use typed values (`domain.ErrNotFound`, `domain.ErrInvalidInput`) to drive HTTP status codes.
- Logging occurs in handlers wrapping domain calls; future instrumentation hooks can be added without leaking into domain code.
- Configuration collected in `main.go` (listen address, signer factory choice, repository backend) to keep `api` package free of global state.
