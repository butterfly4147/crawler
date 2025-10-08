# Repository Guidelines

## Project Structure & Module Organization
- Source code is grouped by domain under directories such as `engine/` (task scheduling), `spider/` (site-specific crawlers), `parse/` (content parsing), `collect/` (data sinks), and `storage/` (persistence helpers). Cross-cutting concerns live in `auth/`, `limiter/`, `proxy/`, and `log/`.
- The CLI entrypoint is `main.go`, which calls the Cobra root in `cmd/`. `cmd/master` orchestrates control-plane commands, while `cmd/worker` starts worker agents. Shared protobuf definitions are in `proto/`; generated assets should be committed.
- Runtime configuration defaults sit in `config.toml` and `default.etcd/`. Documentation and design notes are under `docs/`, and release automation lives in `.github/`.

## Build, Test, and Development Commands
- `make build` compiles the `crawler` binary with version metadata embedded via `version/`.
- `make debug` builds with `-gcflags="all=-N -l"` for breakpoint-friendly binaries.
- `make lint` runs `golangci-lint ./...`; keep the tree clean before every PR.
- `make imports` enforces `goimports` formatting across the module.
- `make cover` runs `go test ./... -short` with coverage and prints a summary. Use this before pushing.
- For quick local runs, use `go run . master --config config.toml` or `go run . worker --config config.toml`.

## Coding Style & Naming Conventions
- Follow the Go toolchain defaults: `gofmt` (tabs, no trailing spaces) and `goimports` are required. Keep functions under ~40 lines, as encouraged in `README.md`.
- Packages, files, and directories use lowercase-without-spaces. Exported symbols need succinct names and doc comments whenever behavior is non-trivial.
- Log messages should originate from the `log/` package; avoid direct `fmt` prints in production paths.

## Testing Guidelines
- Use the standard `testing` package with table-driven tests where sensible. Name files `*_test.go` and functions `TestXxx`.
- When adding new behavior, include unit tests colocated with the package. Integration flows that touch etcd or external services belong under `engine/` or `master/` with build tags.
- Maintain or raise coverage reported by `make cover`; call out notable changes in PRs.

## Commit & Pull Request Guidelines
- Write informative, imperative commits. Prefer Conventional Commit style such as `feat(spider): add rate limit`. Squash noisy work-in-progress commits before review.
- PRs must describe the problem, summarize the solution, list verification steps (e.g., `make lint`, `make cover`), and link any relevant issues. Include configuration or schema changes in the description and attach screenshots or logs for CLI output.

## Configuration & Security Tips
- Keep secrets out of the repo. Use environment variables consumed by `auth/` and `proxy/` packages rather than committing credentials.
- Review `docker-compose.yml` and `kubernetes/` manifests when changing infrastructure-facing code to ensure deployment parity.
