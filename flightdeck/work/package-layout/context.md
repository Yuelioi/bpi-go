# Context

## Stable interface

- Module path stays `github.com/Yuelioi/bpi-go`; the root package remains `bpi`.
- Existing construction and domain calls remain source-compatible: `bpi.NewClient(options...).Activity().Info(...)` and equivalent accessors for all 27 domains.
- Root account, option, response, envelope, and error names remain available with identical error identity and `errors.Is`/`errors.As` behaviour.
- Promoted request behaviour, models, risk classes, fixture evidence, and parity mappings must not change during this structural refactor.

## Target package graph

```text
bpi facade ───────► client
    │                 ▲
    └────► domains ───┘
```

- The root `bpi` module owns only construction compatibility, public aliases/wrappers, domain accessors, package documentation, examples, and repository-wide parity tests.
- `client/` owns HTTP execution, Cookie/account session state, options, envelopes, errors, request policies, WBI signing/cache, and their focused tests. Credential-source loading remains the caller's responsibility.
- Each public domain module owns its params, models, endpoints, domain `Client`, and contract tests.
- `http.RoundTripper` remains the real adapter seam. The package split must not add a second public mock/transport abstraction.

## Design constraints

- Prefer type aliases and narrow forwarding methods in the root facade so public error/type identity is preserved.
- Keep the root `Client` state private rather than embedding and exposing the low-level client.
- Domain constructors may accept `*client.Client`; the ordinary root facade remains the documented interface.
- Tests should cross the same public domain interface as callers. Shared fixture/query helpers may live in `internal/testutil`.
- Do not introduce a `domains/` catch-all directory: existing top-level domain import paths are already public and should remain stable.
- No Git commit, push, tag, live Probe, or credential read is part of this Work.

## Acceptance

- Root contains only a small, coherent set of Go facade/documentation/parity files.
- Every domain implementation and domain contract test is collocated with that domain.
- Existing README examples and exported root call shape compile unchanged.
- The 206/206 parity audit, generated API index, normal tests, race tests, vet, static analysis, and privacy audit pass.
