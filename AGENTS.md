# Repository Guidelines

## Project Structure & Module Organization
Core setup lives in `main.go`, which invokes `internal/app` to compose dependencies. HTTP handlers are grouped in `api/` with the versioned transport under `api/v0/devices`. Domain contracts and the main device service sit in `internal/devices` (including signature history types), while persistence adapters (`persistence/`) provide in-memory repositories and signature stores. Cryptographic helpers and signer factories are placed under `crypto/`. Integration documentation lives in `docs/`, and request collections can live alongside the transport (e.g. `api/http/api.http`).

Runtime configuration is centralised in `internal/config`, which currently reads environment variables such as `LISTEN_ADDRESS` (default `:8080`).

## Build, Test, and Development Commands
- `make build` – compile all packages.
- `make run` – run the HTTP API on `:8080`.
- `make test` – execute the entire suite (unit + integration).
- `make race` – run the suite with the race detector enabled.
- `make integration` – execute only the integration tests in `api/tests`.
- `make mocks` – regenerate gomock doubles under `pkg/mocks` after interface changes.
- `make tidy` – sync module dependencies.
- `make clean` – remove local Go caches and generated mocks.

## Coding Style & Naming Conventions
Stick to `gofmt`/`goimports` formatting (tabs, 120-char guidance). Exported types use PascalCase; locals camelCase. File names mirror concepts (`service.go`, `handler.go`, `memory_store.go`). Keep interfaces concise and named after their roles (e.g., `Repository`, `KeyStore`).

## Testing Guidelines
- Unit tests rely on gomock (see `internal/devices/service_test.go`, `api/v0/devices/handler_test.go`). Use expectations to verify storage interactions (`SignatureStore.Append`, etc.).
- Integration tests (`api/tests/integration_test.go`) cover REST scenarios: lifecycle, validation errors, missing devices, signature history retrieval, and concurrent signing against the in-memory stack.
- Concurrency is exercised by `TestConcurrentSigningIntegration` in `api/tests/integration_test.go`, ensuring signature counters remain monotonic under parallel requests.
- Use `make test` for the full suite; avoid skipping tests unless justified.

## Commit & Pull Request Guidelines
Write imperative commit messages focused on intent (`Add concurrent signing test`). Reference tasks where possible (`[#12]`). Pull requests should:
- Summarize changes and affected modules.
- Include test evidence (`make test` output snippet).
- Highlight API or schema updates, with sample payloads when applicable.

## Security & Configuration Tips
Keep generated keys ephemeral; do not commit `.pem` files. Configuration currently relies on defaults in `main.go`; add environment-based overrides in `internal/app` if needed. Future persistence backends should encrypt keys at rest and document rotation procedures.
