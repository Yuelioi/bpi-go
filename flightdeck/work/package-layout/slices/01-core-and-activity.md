# Slice 01 — Shared client module and activity proof

## Outcome

Establish the final dependency direction and prove it through the complete `activity` module before scaling the migration.

## Current

Complete. The shared client seam and activity package established the dependency direction used by every remaining domain.

## Next

None.

## Steps

- [x] Move account, options, transport, session, envelopes, errors, policies, and WBI implementation into `client/`.
- [x] Expose only the low-level request/signing operations required by domain implementations.
- [x] Add root aliases/wrappers that preserve current type and error identity.
- [x] Move the activity client and tests beside activity params/models.
- [x] Verify root and direct domain construction, public examples, and activity fixtures.
- [x] Run normal tests, vet, formatting, and static analysis before scaling.

## Acceptance

- No root/domain import cycle.
- `bpi.NewClient().Activity()` returns the collocated activity client.
- Existing account, option, error, response, and envelope callers compile.
- Activity's promoted contracts and all shared-client tests pass unchanged in behaviour.
