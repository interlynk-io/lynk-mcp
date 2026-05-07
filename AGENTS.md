# Repository Guidelines

## Project Structure & Module Organization

This is a Go module for the `lynk-mcp` command-line MCP server. The entry point lives in `cmd/lynk-mcp/main.go`. Internal packages are under `internal/`: `api` wraps Lynk API calls, `graphql` contains GraphQL client/query types, `mcp` defines MCP resources/tools/server wiring, `config` handles local configuration and token storage, and `retry` contains retry helpers. Tests live next to the package they cover, for example `internal/retry/retry_test.go`. Build artifacts are written to `build/`; `bin/lynk-mcp` is a checked-in binary and should not be used as the primary development output.

## Architecture Notes

This MCP server is a local bridge to the Lynk GraphQL API. For API-side context, use the sibling checkout at `../lynk-api` when available, or the upstream repository `github.com/interlynk-io/lynk-api`. Keep GraphQL query and type changes in sync with that API.

## Build, Test, and Development Commands

- `make deps`: downloads modules and runs `go mod tidy`.
- `make build`: builds `./cmd/lynk-mcp` into `build/lynk-mcp` for the current OS/architecture.
- `make run`: builds and starts `lynk-mcp serve`.
- `make configure`: builds and runs interactive token/config setup.
- `make check`: builds and runs `lynk-mcp verify`.
- `make test`: runs `go generate ./...` and `go test -cover -race ./...`.
- `make vet`: runs `go vet ./...`.
- `make ci`: runs dependency setup, generation, vetting, and tests.
- `make clean`: removes generated build and coverage artifacts.

## Coding Style & Naming Conventions

Use standard Go formatting: tabs for indentation, `gofmt` output, and idiomatic package names. Run `make fmt` before submitting changes. Keep packages focused and unexported by default; export identifiers only for cross-package APIs. Name tests `TestFunction_Behavior`, matching the existing pattern such as `TestDo_SuccessAfterRetries`. Preserve the repository’s Apache-2.0 file header on new Go source files.

## Testing Guidelines

Use Go’s standard `testing` package. Add tests beside the package under test using `*_test.go` files. Prefer deterministic tests with short timeouts and explicit context cancellation where retry or network-adjacent behavior is involved. Run `make test` before opening a PR; use targeted commands such as `go test ./internal/retry -run TestDo` while iterating.

## Commit & Pull Request Guidelines

Recent commit history uses short, imperative summaries with optional PR numbers, for example `Fix MCP connection timeout for new tokens with two-layer retry (#16)` and `Support for token (#12)`. Keep commits focused on one logical change. PRs should include a concise description, validation steps (`make test`, `make ci`, or targeted tests), related issues, and screenshots or terminal output only when user-visible CLI behavior changes.

## Security & Configuration Tips

Do not commit API tokens or local config. Runtime config is stored in `~/.lynk-mcp/config.yaml`, and tokens should come from the system keychain or `LYNK_API_TOKEN`. Use `lynk_test_*` or staging tokens for development.
