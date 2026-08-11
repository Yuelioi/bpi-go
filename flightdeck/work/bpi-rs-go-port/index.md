# bpi-rs Go port

**Status:** Complete

## Goal

Deliver an idiomatic Go SDK that reproduces the mature `bpi-rs` 0.2 module-client behaviour and promoted contract surface without preserving the discarded legacy Go implementation.

## Result

The initial parity scope is complete against `bpi-rs` commit `36cb1104befee33b4281c59a4365d42dfdde45a4` (Cargo version 0.2.4):

- All 27 promoted domain clients have Go counterparts.
- All 206 promoted contracts map one-to-one to supported Go methods; there are no missing, duplicate, stale, extra, or unsupported mappings.
- Risk coverage is exact: 115 public reads, 32 authenticated reads, 52 private reads, and 7 login-session contracts.
- The committed evidence contains 432 unique response fixtures and 642 byte-locked contract files.
- Public methods are context-first, return typed business payloads, validate parameters before transport, and preserve recoverable decode bodies without logging them.
- Credentials, Cookie refresh, WBI and Bili-ticket signing, request policies, special response formats, and login-session primitives are implemented with per-client state and offline verification.
- `bpi-probe` provides deterministic parity, snapshot, documentation, and sanitization audits plus a doubly gated read-only live Probe.
- README, generated API index, migration, Probe, contribution, security, release, changelog, CI, and dependency-update configuration are present.

## Verification

Passed on 2026-08-11 from the Go repository without loading an account file or making a live Bilibili request:

- `gofmt` drift check and `go mod tidy`
- `go vet ./...`
- `go run honnef.co/go/tools/cmd/staticcheck@latest ./...`
- `go test -count=1 ./...`
- `go test -race -count=1 ./...`
- `go run golang.org/x/vuln/cmd/govulncheck@latest ./...` — no vulnerabilities found
- `go run ./cmd/bpi-probe audit` — 206 contracts, 206 mappings, 432 fixtures, 642 locked files
- `go run ./cmd/bpi-probe api-doc --check`
- `go run ./cmd/bpi-sourcegen -rust ../bpi-rs -expect-commit 36cb1104befee33b4281c59a4365d42dfdde45a4 -expect-domains 27 -expect-contracts 206 -check`

The source-generator check confirms byte-for-byte contract-tree parity with the sibling Rust checkout. Ordinary tests and audits rely only on committed Go-repository evidence and do not require that sibling checkout.

## Safety boundary

No live Probe was run during completion. The live runner cannot execute login-session, mutating, or spending contracts, and read-only network execution requires both `BPI_PROBE=1` and `--read-only`. A later simplification removed the TOML loader; authenticated Probe profiles now receive raw Cookie headers through explicit profile-specific environment variables. Probe output remains ignored.

## Repository state

The obsolete local implementation was removed after review, Git was initialized on `main`, and no commit was created. The completed repository remains uncommitted so the owner can review and choose the initial commit history.

## Follow-up

[Plan item 9](plan.md#9-post-parity-capability-review) is intentionally outside this completed Work. Unpromoted legacy, mutating, or spending capabilities require a separately authorized review rather than automatic migration.

## References

- [Fixed decisions](context.md)
- [Completion plan](plan.md)
- [Generated Go API index](../../../docs/api-index.md)
- [Parity workflow](../../../parity/README.md)
- [bpi-rs README](../../../../bpi-rs/README.md)
- [bpi-rs promoted contracts](../../../../bpi-rs/tests/contracts)
