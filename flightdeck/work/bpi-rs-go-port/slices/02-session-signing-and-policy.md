# Slice 02 — Session, signing, and request policy

> Current note: the account-profile/TOML loader recorded below was removed after the initial port. The supported interface is now explicit `Account` or raw Cookie input, and the module has no third-party dependency.

## Outcome

Complete the shared typed identifiers, credential helpers, signing algorithms, caches, and request policies needed by promoted domain clients.

## Current

Complete. The root client now has explicit structured account-profile loading, full Rust-parity typed identifiers, deterministic WBI and Bili-ticket signing, per-client single-flight WBI key caches, response-Cookie refresh, and small package-internal request policies for query, form, streaming multipart, CSRF, WBI, optional payload, and raw response flows.

## Next

Continue with [Slice 03](03-representative-vertical-slices.md), starting with `activity.info` and `activity.list` through the public domain-client surface.

## Steps

- [x] Implement numeric and string identifier types with constructors, parsing, formatting, JSON/text encoding, zero-value validation, and Rust parity tests.
- [x] Add explicit account-profile loading as an optional helper that never runs during `NewClient`.
- [x] Add deterministic WBI mixin/signing functions and Rust test vectors.
- [x] Add per-client WBI key cache, clock injection, navigation-key decoding, and single-flight refresh behaviour.
- [x] Add Bili-ticket HMAC helpers and Rust test vectors.
- [x] Inventory and implement only the remaining signing helpers used by promoted domain-client methods.
- [x] Add shared form, multipart, raw-byte, optional-payload, CSRF, and WBI request policies behind the client interface.
- [x] Verify Cookie response updates and external-host credential scoping under concurrency.
- [x] Run formatting, vet, normal tests, race tests, and security-focused fuzz seeds.
- [x] Save the Work and check Plan item 2 complete only when all acceptance checks pass.

## Acceptance

- Every Rust typed-ID vector has an equivalent Go test and invalid zero/blank forms fail before network execution.
- WBI and Bili-ticket outputs match fixed Rust vectors exactly.
- WBI keys refresh once per cache bucket under concurrent requests and caches remain isolated between clients.
- No signing secret or generated signature appears in logs.
- Account-profile parsing is explicit, validates completeness, and never causes constructor file I/O.
- Request policies reuse the Client implementation instead of exposing a new public transport seam.
- The full offline validation suite passes with the race detector.

## Evidence

- [Rust typed IDs](../../../../../bpi-rs/src/ids.rs)
- [Rust WBI signing](../../../../../bpi-rs/src/sign/wbi.rs)
- [Rust WBI client integration](../../../../../bpi-rs/src/sign/wbi_client.rs)
- [Rust Bili-ticket signing](../../../../../bpi-rs/src/sign/bili_ticket.rs)

## Verification

- All Rust typed-ID examples and fixed WBI/Bili-ticket signing vectors pass in Go.
- The WBI navigation cache performs one fetch per client/hour bucket under concurrent use and remains isolated between clients.
- Account profile loading supports only explicit `[normal]` and `[vip]` sections and never runs during construction.
- Cookie parser and WBI fuzz seeds execute in the default offline suite without panics.
- Response Cookie updates, expiration, caller-supplied external Cookies, and external response-Cookie isolation are covered by tests.
- Form, streaming multipart, CSRF, WBI, required/optional envelope, and raw response policies use the existing Client/HTTP seam.
- `go vet ./...`, `go test -count=1 ./...`, and `go test -race -count=1 ./...` passed on 2026-08-11.

## Decisions

- `github.com/pelletier/go-toml/v2` v2.4.3 is the only dependency added; it is narrowly scoped to the explicit optional account-profile loader.
- `misc/sign/appkey.rs` and `misc/sign/v_voucher.rs` are placeholders and are not used by promoted domain-client methods, so they are not advertised or ported.
- Request construction stays package-internal and composable. The public standard-library `http.Request`/`Client.Do` seam remains the custom-request escape hatch.
