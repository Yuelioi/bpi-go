# Slice 07 — Probe, contract sync, and safety gates

## Outcome

Complete. The repository has deterministic source synchronization, parity documentation, privacy auditing, and a Go-native read-only Probe.

## Delivered

- `bpi-sourcegen -check` verifies the locked source manifest, source lock, and byte-for-byte copied contract tree against the fixed Rust commit.
- `bpi-probe audit` verifies commit/schema locks, exact mappings and risks, contract cases, fixtures, hashes, and sanitization.
- `bpi-probe api-doc` deterministically generates `docs/api-index.md`.
- `bpi-probe batch-run` supports anonymous, normal, and VIP read-only profiles and emits body-free summaries containing only bounded metadata and hashes.
- Live network execution requires both `BPI_PROBE=1` and `--read-only`; login-session, mutating, and spending contracts cannot run through this command.

## Acceptance

- [x] Default invocation performs no network request.
- [x] Read-only gates cannot execute higher-risk classes.
- [x] The 642-file snapshot lock and 206 mappings are deterministic.
- [x] Sanitizer tests reject representative credentials and sensitive response material.
