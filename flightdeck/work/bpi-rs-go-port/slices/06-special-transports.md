# Slice 06 — Special transports and formats

## Outcome

Complete for every promoted contract evidenced by the locked Rust baseline.

## Delivered

- Added raw-byte, XML, compressed danmaku, playback, multipart/form, CSRF, heartbeat, and multi-step request handling where promoted contracts require them.
- Kept protocol-specific handling behind deep internal modules so ordinary JSON domain calls retain one consistent interface.
- Documented response ownership and bounded-body behaviour; context cancellation propagates to every transport.
- Did not advertise incomplete Rust placeholders such as unsupported live WebSocket or unpromoted mutation/spending surfaces.

## Acceptance

- [x] Every promoted non-JSON or special request shape has deterministic fixtures and public-surface tests.
- [x] Parsers and decoders handle malformed inputs without panics; fuzz seeds run in the default suite.
- [x] No placeholder capability was promoted without contract evidence.
