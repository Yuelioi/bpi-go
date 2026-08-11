# Plan

> Historical note: the completed initial plan included an optional account-profile loader. That helper and its TOML dependency were later removed; credential-source loading now belongs to callers.

This is the ordered completion rollup for the Go reimplementation. A plan item is complete only when its listed acceptance checks pass; compiling a skeleton is not parity.

## Rollup

- [x] 0. Initialize durable Work and lock the Rust source baseline.
- [x] 1. [Bootstrap the clean Go module and core foundation](slices/01-bootstrap-and-foundation.md).
- [x] 2. [Implement session, credentials, signing, and request-policy modules](slices/02-session-signing-and-policy.md).
- [x] 3. [Prove the architecture with representative vertical slices](slices/03-representative-vertical-slices.md).
- [x] 4. [Port all promoted public-read JSON capabilities](slices/04-public-read-parity.md).
- [x] 5. [Port authenticated-read, private-read, and login-session capabilities](slices/05-authenticated-private-login-parity.md).
- [x] 6. [Port promoted non-JSON and special transport flows](slices/06-special-transports.md).
- [x] 7. [Build the Go-native Probe, contract sync, and safety gates](slices/07-probe-contract-safety.md).
- [x] 8. [Complete parity documentation, examples, CI, and release readiness](slices/08-docs-ci-release.md).
- [ ] 9. Evaluate unpromoted legacy, mutating, and spending capabilities as a separate post-parity scope (intentionally deferred; not part of this Work's definition of complete).

## 0. Baseline and scope

Completed in this initialization:

- Locked source commit `36cb1104befee33b4281c59a4365d42dfdde45a4`.
- Defined promoted domain clients and 206 promoted contracts as the initial parity target.
- Declared the old local Go implementation out of scope for compatibility.
- Established Flightdeck recovery documents and the first executable Slice.

## 1. Clean module and core foundation

Deliverables:

- Replace the legacy package layout with a root `bpi` package and domain model/param packages.
- Create a machine-readable source-surface and contract parity manifest.
- Implement isolated `Client` construction with validated functional options.
- Implement the internal HTTP request executor, default headers, context propagation, safe logging metadata, and body-size policy.
- Implement generic Bilibili envelope decoding with `data`/`result`, `code`/`errno`, and message aliases.
- Implement typed HTTP, API, parameter, authentication, missing-data, and response-decode errors.
- Preserve response bytes only on model decode failure and explicitly redact them from formatting and logs.
- Establish `RoundTripper` adapters and fixture helpers for offline tests.

Acceptance:

- Two clients have isolated Cookie, cache, logger, and HTTP state.
- Constructor execution does not read files, install global loggers, or perform network I/O.
- Cancellation and deadlines propagate to the HTTP request.
- Sensitive query values and headers never appear in logs.
- Envelope fixtures from Rust decode equivalently, including aliases, API errors, missing data, and recoverable model mismatch.
- `go test ./...`, `go test -race ./...`, `go vet ./...`, and formatting checks pass offline.

## 2. Session, credentials, signing, and policy

Deliverables:

- `Account` and Cookie parsing with validation, cloning, clearing, CSRF access, and redacted debug representation.
- Explicit account-profile file loading as an optional helper outside client construction.
- Cookie scoping limited to Bilibili hosts; external custom URLs do not receive credentials.
- Typed IDs matching Rust invariants: `AID`, `AudioID`, `CID`, `MID`, `RoomID`, `MediaID`, `SeasonID`, `EpisodeID`, `NoteID`, `CVID`, `BVID`, and `DynamicID`.
- Deterministic WBI signing, URL-key extraction, per-client key caching, clock injection for tests, and cache refresh behaviour.
- Bili-ticket HMAC signing and remaining promoted signing helpers.
- Shared request policies for JSON envelope, optional payload, raw bytes, form, multipart, CSRF, and WBI-signed requests.

Acceptance:

- Rust signing vectors pass unchanged in Go.
- Cookie parsing fuzz tests do not panic and reject malformed segments predictably.
- Account and decode-error secrets do not appear in `%v`, `%+v`, slog output, or JSON serialization.
- WBI cache is concurrency-safe, isolated per client, and deterministic under an injected clock.
- Requests to non-Bilibili hosts carry no session Cookie unless explicitly supplied by the caller.

## 3. Representative vertical slices

Port these first, in order:

1. `activity.info` and `activity.list` — simple public JSON and pagination defaults.
2. `video.view`, `video.detail`, `video.pagelist`, and `video.desc` — typed video identifiers and nested models.
3. `login.nav` — authenticated payload and explicit session behaviour.
4. `video.play_url` — WBI signing and playback parameters.
5. One optional-payload endpoint and one endpoint whose success payload uses a response alias.
6. One raw-byte or XML endpoint to validate that JSON is not baked into the domain-client interface.

For every method, deliver params, models, domain method, request contract assertion, fixture decode assertion, error fixture assertion, documentation, and parity mapping in the same change.

Acceptance:

- The public call style remains consistent across all representative shapes.
- Tests exercise domain methods through the same interface callers use; they do not assert private request-builder state.
- No new public seam is introduced solely for tests.
- The first slices demonstrate anonymous, authenticated, WBI-signed, optional, and non-JSON behaviour.

## 4. Promoted public-read parity

Port public JSON capabilities in capability waves while keeping each method vertically complete:

1. Core discovery: `video`, `search`, `user`, `video_ranking`, `web_widget`, `clientinfo`, and `opus`.
2. Media catalogues: `activity`, `article`, `audio`, `bangumi`, `cheese`, and public `manga`.
3. Interaction reads: public `comment`, `dynamic`, `fav`, `historytoview`, `note`, and `misc`.
4. Public HTTP portions of `live` and `danmaku`.

Acceptance:

- Every `public-read` promoted contract maps to exactly one supported Go method or an explicitly justified unsupported entry.
- Each mapped method has parameter, request, fixture, and error verification.
- Stable contract fields are typed; deliberately unstable nested payloads are narrowly isolated as raw JSON.
- API-index generation reports no missing public-read implementation.

## 5. Authenticated, private, and login-session parity

Port in increasing sensitivity:

1. Authenticated read-only status and personalization.
2. User-owned collections, histories, and account state.
3. Creative center, message, wallet, moderation, and other private payloads.
4. Login QR generation/polling, session refresh, exit, and associated response-Cookie handling.

Acceptance:

- All 32 authenticated-read, 52 private-read, and 7 login-session contract classifications are represented in the parity report.
- Default tests use only sanitized committed fixtures and never load local credentials.
- Account profiles are explicit and private fields are minimized in committed models and fixtures.
- Login flows expose primitives and state, while caller applications retain control of QR rendering, polling cadence, timeout, and persistence.

## 6. Special transports and formats

Implement protocol-specific modules without weakening the ordinary JSON interface:

- XML and compressed danmaku responses.
- Binary/protobuf danmaku segments where supported by the Rust source and contracts.
- Live WebSocket framing, heartbeat, decompression, and event decoding only when the source capability is complete enough to specify.
- Streaming media/download responses with explicit ownership and close semantics.
- Multipart upload and publish preparation.
- Multi-step flows represented by Probe contracts.

Acceptance:

- Streaming methods document ownership, cancellation, buffering, and close requirements.
- Compression and frame parsers have deterministic fixture tests and fuzz coverage for malformed input.
- Placeholder Rust files do not become advertised Go capabilities without evidence and acceptance tests.

## 7. Probe, contract sync, and risk gates

Deliverables:

- A contract snapshot sync command with source commit lock and drift report.
- A Go API/parity index generated from contracts plus Go declarations.
- A read-only Probe CLI supporting anonymous, normal, and VIP profiles.
- Sanitization audit that rejects credential keys, account identifiers, private content, and unsanitized raw output.
- Risk gates matching the Rust classifications.

Acceptance:

- Default Probe invocation performs no network request without opt-in.
- Mutating and spending classifications cannot run under the read-only gate.
- Contract sync is deterministic and CI detects source revision or mapping drift.
- Sanitization tests fail on representative credential and private-data leaks.

## 8. Documentation, CI, and release readiness

Deliverables:

- README quick start, anonymous/authenticated examples, custom request recovery, and explicit safety notes.
- Generated API index with Go symbols, contracts, risk class, and support status.
- Migration notes describing semantic correspondence with `bpi-rs`, not compatibility with the discarded Go skeleton.
- CI for formatting, vet, offline tests, race tests, examples, dependency audit, and contract parity.
- Contribution, security, release, and changelog documents.

Acceptance:

- All documentation examples compile; live examples require explicit execution gates.
- The promoted parity report has no unexplained gaps.
- A clean checkout passes the full offline validation command set without sibling repositories or account files.
- Public package documentation states concurrency, context, credential, logging, and response-body guarantees.

## 9. Post-parity capability review

After promoted parity, inventory unpromoted Rust free functions and legacy modules. Do not port them automatically.

For each candidate, require live evidence, risk classification, a stable Go interface, sanitized fixtures where permitted, and explicit safety gates. Mutating and spending capabilities remain outside default examples and default tests.

## Definition of complete

The initial Work is complete when:

- All 27 promoted domain clients have a Go counterpart.
- All 206 promoted contracts are implemented or carry a reviewed, documented unsupported decision.
- All supported contracts have offline request and response verification.
- Default construction and testing have no hidden file, global-state, credential, network, or side-effect behaviour.
- The SDK passes formatting, vet, tests, race tests, examples, parity audit, and privacy audit from a clean checkout.
- README and generated API index accurately describe the shipped interface.
