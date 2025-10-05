# Repository Guidelines

## Project Structure & Module Organization
The service code lives in `main.go`, which wires the HTTP server from `api/`. API handlers are split into `api/server.go` for router setup, `api/device.go` for device endpoints, and `api/health.go` for readiness checks. Domain logic for devices, counters, and signing contracts sits in `domain/device.go`. Cryptographic primitives and key encoding helpers are under `crypto/`, while storage abstractions reside in `persistence/`. Keep integration assets and future fixtures inside `persistence/fixtures` (create the folder if needed) so runtime data stays separate from Go sources.

## Build, Test, and Development Commands
Use `go build ./...` to ensure all packages compile. Run `go test ./...` to execute unit tests once they exist. Start the local server with `go run ./main.go` and point clients at `http://localhost:8080`. Add `GOLOG=debug` when running locally to surface verbose logs while developing new handlers.

## Coding Style & Naming Conventions
Follow standard Go formatting (`gofmt` or `goimports` before committing). Stick with tabs for indentation and keep lines under 120 characters. Use PascalCase for exported types and functions, camelCase for locals, and uppercase abbreviations only when common (`UUID`, `RSA`). Name files after the concept they encapsulate (`device_service.go`, `memory_store.go`). Prefer concise interfaces that expose domain verbs (`Signer`, `DeviceRepository`).

## Testing Guidelines
Favor table-driven tests placed next to the code under test, e.g., `domain/device_test.go`. Name tests `TestFunction_Scenario` for clarity. Cover happy paths, signature counter edge cases, and error propagation from persistence and crypto layers. For HTTP handlers, use `httptest` with JSON fixtures to verify response codes and payloads. Aim for coverage on all domain mutations before merging.

## Commit & Pull Request Guidelines
Write commits in the imperative mood (`Add RSA signer cache`), scoped to a single logical change; the history currently uses short "Initial commit" messages—continue in that concise style. Reference issue IDs when available (`[#12]`). Pull requests should summarize intent, list test evidence (`go test ./...` output), and mention any follow-up tasks. Include API or schema changes in the description and attach example requests or responses when they help reviewers.

## Security & Configuration Tips
Key material should remain ephemeral in development; avoid committing generated keys or `.pem` files. Store sensitive runtime configuration in environment variables, not source. When adding persistence backends, ensure keys and signatures are encrypted at rest and document rotation procedures in `README.md`.
