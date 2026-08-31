# Context

## Baselines

- The current Go parity lock names bpi-rs commit `36cb1104befee33b4281c59a4365d42dfdde45a4`, Cargo version 0.2.4.
- The requested source release is bpi-rs tag `v0.3.0` at commit `94bcf43`.
- The sibling checkout at `../bpi-rs` is clean on `main` at `4b8d491`, one documentation-only commit after the tag.
- The Go repository is clean on `main` before this Work begins.

## Upstream release claims

bpi-rs 0.3.0 declares breaking response-model corrections for FLAC DASH audio and signed Bangumi fields, fixes null backup URL and Bangumi `durl` decoding, preserves nonzero API errors instead of masking them with payload decode errors, and broadens historical Dynamic and Bangumi response compatibility. It also adds representative response matrices and Windows/feature-isolation fixes.

## Constraints

- Keep the established Go package layout and context-first public methods unless an upstream behavior requires an intentional API change.
- Preserve the response-body recovery and redaction guarantees already provided by the Go client.
- Treat tag `v0.3.0`, not the later documentation-only `main` commit, as the parity source of truth.
- Verification must remain offline unless live Probe execution is separately authorized.
- Do not commit, tag, publish, or alter the sibling bpi-rs checkout.

## Verified scope decisions

- bpi-rs 0.3.0 changes no promoted contract, endpoint, risk class, domain client, or contract fixture; only the locked source commit and Cargo version change in generated source evidence.
- Go slices already accept missing and JSON `null` list values, so the Rust-specific nullable backup URL deserializer requires regression coverage but no new Go decoding hook.
- Go models already use zero values for missing display fields and retain Bangumi activity/payment payloads as raw JSON, so Rust's added serde defaults and signed pendant field do not require equivalent Go public type changes.
- `Envelope` continues decoding compatible nonzero payloads for explicit unchecked extraction, but suppresses only incompatible payload decoding when a nonzero code provides the authoritative API error.
