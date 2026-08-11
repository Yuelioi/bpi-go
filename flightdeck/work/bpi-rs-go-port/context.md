# Context

## Product intent

This repository will become the Go counterpart of `bpi-rs`. The existing local Go source is not a migration base and is not a compatibility constraint. Reuse from it is allowed only after independent comparison with the Rust source and promoted contract evidence.

The intended user experience is the Rust 0.2 module-client style expressed with Go conventions:

```go
client, err := bpi.NewClient(bpi.WithCookie(cookie))
if err != nil {
    return err
}

bvid, err := ids.ParseBVID("BV1xx411c7mD")
if err != nil {
    return err
}

view, err := client.Video().View(ctx, video.ViewByBVID(bvid))
```

## Source of truth

The baseline source is sibling repository `bpi-rs` at commit `36cb1104befee33b4281c59a4365d42dfdde45a4` and Cargo version `0.2.4`.

When sources disagree, use this precedence:

1. Sanitized promoted contracts and response fixtures under `bpi-rs/tests/contracts/` for observed request and response behaviour.
2. Domain-client methods under `bpi-rs/src/<domain>/client.rs` for the supported public capability surface.
3. Typed params, IDs, response, session, signing, and transport modules for invariants and implementation semantics.
4. `docs/api-index.md`, README, migration, development, and risk documentation for intended usage and safety policy.
5. Legacy free functions, ignored live tests, placeholders, and unpromoted source only as investigation leads.

The initial parity target is:

- 27 domain clients.
- 206 promoted contracts.
- Risk distribution: 115 public reads, 32 authenticated reads, 52 private reads, and 7 login-session contracts.
- Direct business-payload returns for normal domain methods.
- Explicit account and Cookie configuration.
- Recoverable response-model decode failures without logging raw bodies.

## Fixed Go design decisions

- The module path remains `github.com/Yuelioi/bpi-go`; the root package name is `bpi`.
- There is no global client, global account, implicit configuration-file read, or automatic live request.
- `NewClient(opts ...Option) (*Client, error)` creates an isolated, concurrency-safe client.
- Every network operation accepts `context.Context` as its first argument.
- `Client.Video()`, `Client.User()`, and peers return lightweight domain clients defined in the root package. Domain params and models live in top-level packages such as `video` and `user`; this preserves `client.Video().View(...)` without Go import cycles.
- Domain methods return decoded payloads. The generic envelope remains available only for custom requests and recovery.
- Parameter constructors prevent invalid identifier combinations. Numeric IDs reject zero; string IDs validate their documented form.
- Standard-library `http.RoundTripper` is the transport test seam. Do not add a second public transport abstraction unless a genuinely different production adapter appears.
- Use `log/slog` for optional structured logging. The library is quiet by default and all request metadata is sanitized before logging.
- Prefer the standard library. Add a dependency only for a protocol the standard library cannot reasonably implement, and record the reason in the relevant Slice.
- Superseded after the initial port: the Go module no longer exposes account-profile loading or depends on a TOML parser. Callers supply `Account` values or raw Cookie request headers.
- Shared query, form, streaming multipart, CSRF, and WBI construction helpers remain package-internal; callers use domain methods or the existing `http.Request` escape hatch.
- Public errors follow Go conventions: typed errors, `Unwrap`, `errors.Is`/`errors.As`, stable sentinel errors where appropriate, and semantic helper functions for login, VIP, permission, and risk-control conditions.
- Parameter validation errors are implemented in an internal dependency-neutral package and re-exported as a root-package type alias, so domain packages avoid import cycles while callers can still use `errors.As` with `*bpi.ParameterError`.
- Domain parameter structs own their public `EncodeQuery` validation seam. Root domain clients consume that seam and keep URL/request execution internal.
- IP request parameters use `net/netip`; scoped IPv6 addresses are rejected because the Rust source accepts plain `IpAddr` values only.
- Stable JSON embedded inside a response string is decoded by the domain model's `UnmarshalJSON`, preserving the wire string while routing malformed embedded data through the root typed response-decode error.
- Domain-specific series IDs remain in package `video` because the Rust cross-domain typed-ID set does not define one. Collection pagination accepts both promoted `page_num`/`page_size` and `num`/`size` response aliases in one deep model.
- Default `go test ./...` is offline and side-effect free. Network probes and mutating flows require explicit environment gates.
- Public implementation and response models are handwritten. Generation may maintain parity manifests and documentation, but must not generate an opaque public SDK surface.

## Contract handling

The Go repository must contain a reproducible snapshot of the Rust promoted contracts so CI does not depend on a sibling checkout. A lock file records the source commit. Contract files remain evidence and are not silently rewritten to fit Go models.

Rust-specific metadata such as `rust_model` may remain in the source snapshot. A Go parity manifest maps each contract name to its Go domain method and model. CI fails for duplicate, missing, or stale mappings.

## Safety and privacy

- Never commit `account.toml`, Cookie strings, `SESSDATA`, `bili_jct`, `buvid3`, access tokens, raw private responses, or unsanitized Probe output.
- Error strings, debug formatting, serialization, and logs must not contain raw response bodies or credentials.
- Read-only live probes require one explicit opt-in gate.
- Authenticated and private live probes require explicit profile-specific Cookie environment variables.
- Mutating operations require an action gate and explicit target identifiers.
- Spending operations require a second independent confirmation gate.

## Non-goals for initial parity

- Preserving the current local Go package layout or method names.
- Translating Rust syntax or file structure one-to-one.
- Promoting Rust placeholders, legacy free functions, or ignored tests without contract evidence.
- Adding mutating or spending operations to the default examples or offline suite.
- Hiding unstable fields behind broad `map[string]any` models when a stable typed subset is known; raw JSON is reserved for explicitly unstable subtrees.

## Tooling constraint

The current Codex skill catalog does not expose a skill named `matt`. Until its exact package or path is supplied, migration design uses the available `codebase-design` guidance: deep modules, small interfaces, the standard HTTP test seam, and verification through public interfaces.
