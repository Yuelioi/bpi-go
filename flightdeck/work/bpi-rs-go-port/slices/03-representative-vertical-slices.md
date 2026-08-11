# Slice 03 — Representative vertical domain slices

## Outcome

Prove that the shared foundation supports stable Go domain methods across public JSON, authenticated state, WBI signing, optional payload, aliases, and raw-byte responses before scaling to all 206 promoted contracts.

## Current

Complete. Eleven promoted contracts are implemented through public Go domain methods and exercise all representative response and request-policy shapes.

## Next

Continue with [Slice 04](04-public-read-parity.md), beginning with small independent public-read domains and retaining the same contract-first vertical workflow.

## Steps

- [x] Port `activity.info` and `activity.list`, including pagination defaults and fixture tests.
- [x] Port `video.view`, `video.detail`, `video.pagelist`, and `video.desc` using typed AID/BVID params.
- [x] Port `login.nav` with explicit authenticated/anonymous envelope semantics.
- [x] Port WBI-signed `video.play_url` with request-signing verification.
- [x] Port optional-payload `comment.read.hot` and `result`-alias `bangumi.timeline`.
- [x] Port raw binary `danmaku.web.seg` without bypassing Client policy.
- [x] Add a Go parity mapping and audit that connects each supported contract to one public method/model.
- [x] Verify all representative methods through their public call surface with offline contract fixtures.
- [x] Run formatting, vet, normal tests, race tests, and examples.
- [x] Save Work and check Plan item 3 complete only when every representative shape passes.

## Acceptance

- Public call style is consistent across activity, video, login, optional, and raw-response methods.
- Endpoint tests assert method, URL, query/body, credential policy, and WBI requirements through captured HTTP requests.
- Response tests use copied sanitized promoted fixtures and cover success plus API/model errors.
- Models type stable fields and isolate only genuinely unstable subtrees as `json.RawMessage`.
- No public seam exists solely for tests; `http.RoundTripper` remains the adapter.
- Every implemented method has exactly one parity mapping tied to the locked Rust contract.
- Full offline formatting, vet, normal, and race suites pass.

## Evidence

- `parity/implemented.json` contains 11 unique mappings locked to the Rust source commit.
- Copied promoted contracts and sanitized fixtures live below `testdata/contracts/`.
- Domain tests assert observable HTTP requests through `http.RoundTripper` and decode the promoted fixtures.
- `go run ./cmd/bpi-sourcegen` reproduces `parity/source.json` byte-for-byte.
- `gofmt`, `go vet ./...`, `go test -count=1 ./...`, and `go test -race -count=1 ./...` pass offline.

- [Rust Activity client](../../../../../bpi-rs/src/activity/client.rs)
- [Rust Activity params/models](../../../../../bpi-rs/src/activity)
- [Rust Video client](../../../../../bpi-rs/src/video/client.rs)
- [Rust Login client](../../../../../bpi-rs/src/login/client.rs)
- [Promoted contracts](../../../../../bpi-rs/tests/contracts)
