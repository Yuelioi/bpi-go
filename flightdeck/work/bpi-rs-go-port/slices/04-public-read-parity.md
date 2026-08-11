# Slice 04 — Promoted public-read parity

## Outcome

Complete. All 115 promoted `public-read` Rust contracts map to supported, vertically verified Go methods.

## Delivered

- Completed public discovery, media catalogue, interaction, live HTTP, and danmaku capabilities across all affected domains.
- Added typed params and stable response models, with narrowly scoped raw JSON only for demonstrably unstable subtrees.
- Verified method, path, query or body, signing, headers, credential scope, success fixtures, and error fixtures through public domain methods.
- Mapped every public contract exactly once in `parity/implemented.json`.
- Kept all copied source contracts byte-faithful to the locked Rust tree.

## Acceptance

- [x] All 115 public-read contracts are supported.
- [x] The parity audit rejects duplicate, missing, stale, extra, or risk-mismatched mappings.
- [x] Offline checks need no account file, network access, or sibling Rust checkout.
- [x] The public API remains domain-oriented and exposes no test-only request-builder seam.

## Evidence

- [Locked source inventory](../../../../parity/source.json)
- [Implemented mappings](../../../../parity/implemented.json)
- [Contract lock](../../../../parity/contracts.lock.json)
- [Generated API index](../../../../docs/api-index.md)
