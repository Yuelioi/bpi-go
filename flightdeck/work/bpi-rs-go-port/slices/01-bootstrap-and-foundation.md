# Slice 01 — Bootstrap and core foundation

## Outcome

Create a clean, compilable Go SDK foundation that can support every promoted domain without inheriting the legacy local implementation shape.

## Current

Complete. The legacy Go files and local account file were deleted after their absolute targets were verified inside the workspace. Git was initialized on `main` with no commit. The repository now contains a standard-library-only root `bpi` package, a deterministic Rust source inventory, an isolated client/session implementation, bounded HTTP execution, sanitized logging, typed errors, envelope decoding, fixture adapters, and offline tests.

## Next

Continue with [Slice 02](02-session-signing-and-policy.md).

## Steps

- [x] Generate `parity/source.json` or equivalent from Cargo features, `src/<domain>/client.rs`, and `tests/contracts/**/contract.json`.
- [x] Verify the inventory reports 27 domain clients, 206 promoted contracts, unique contract names, and the expected risk totals.
- [x] Record the Rust commit in a contract/source lock file.
- [x] Resolve the absolute legacy deletion targets and confirm Flightdeck files are outside them.
- [x] Replace the old package layout while preserving `flightdeck/` and repository instructions.
- [x] Rebuild `go.mod` around the root `bpi` package, starting standard-library-only.
- [x] Add package documentation and the minimal `Client`, `Option`, domain-client accessor, and configuration types.
- [x] Add the internal request executor using `http.Client` and the `http.RoundTripper` seam.
- [x] Add envelope decoding and typed Go errors with response-body recovery/redaction.
- [x] Add fixture adapters and table tests for construction, isolation, context cancellation, headers, aliases, API errors, missing data, and model mismatch.
- [x] Run formatting, `go vet ./...`, `go test ./...`, and `go test -race ./...`.
- [x] Update the Work handoff and check this Slice complete only when all acceptance checks pass.

## Verification

- Source generation completed deterministically for 27 domains and 206 promoted contracts at Rust commit `36cb1104befee33b4281c59a4365d42dfdde45a4`.
- `go vet ./...` passed.
- `go test -count=1 ./...` passed.
- `go test -race -count=1 ./...` passed.
- Transport-error, request-log, account, and response-decode representations were verified not to leak test secrets.

## Interface sketch to validate

```go
type Client struct {
    // private HTTP, session, policy, cache, and logger state
}

func NewClient(opts ...Option) (*Client, error)

func WithHTTPClient(client *http.Client) Option
func WithCookie(cookie string) Option
func WithAccount(account Account) Option
func WithLogger(logger *slog.Logger) Option

func (c *Client) Video() VideoClient

type VideoClient struct {
    client *Client
}

func (v VideoClient) View(ctx context.Context, params video.ViewParams) (*video.View, error)
```

The sketch is a design target, not permission to expose internal transport or request-builder details. Adjust names when Go documentation or import-cycle checks reveal a clearer interface, but preserve the call shape and invariants in `context.md`.

## Foundation acceptance

- `NewClient()` is deterministic, offline, and free of global side effects.
- Invalid options fail construction without leaving partially initialized shared state.
- Two clients remain isolated under concurrent account, Cookie, logger, and WBI-cache use.
- Domain clients are cheap wrappers and do not duplicate HTTP/session implementation.
- Tests substitute `http.RoundTripper`; production and test callers cross the same client/domain interface.
- Context cancellation reaches the adapter and interrupts body reads.
- Default request headers match Bilibili expectations, but credentials are scoped to Bilibili hosts.
- Logs contain operation, method, sanitized URL, status, API code, and duration without credential-bearing query/header values.
- Envelope decoding handles Rust fixture aliases and distinguishes API failure, missing required data, optional data, and model mismatch.
- Response decode errors retain an explicitly accessible copy of the raw body but redact it from `Error`, formatting, slog, and serialization.
- The repository passes offline formatting, vet, tests, and race tests from a clean checkout.

## Evidence to consult during execution

- [Root client](../../../../../bpi-rs/src/client.rs)
- [Request helpers](../../../../../bpi-rs/src/request.rs)
- [Envelope semantics](../../../../../bpi-rs/src/response.rs)
- [Error semantics](../../../../../bpi-rs/src/err/error.rs)
- [Transport metadata and redaction](../../../../../bpi-rs/src/transport)
- [Envelope fixtures](../../../../../bpi-rs/tests/fixtures/envelope)
