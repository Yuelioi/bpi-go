# Plan

- [x] 0. Audit the package graph, exported root surface, domain helpers, and test coupling.
- [x] 1. [Extract the shared client module and migrate `activity`](slices/01-core-and-activity.md).
- [x] 2. Migrate every remaining domain implementation and test into its domain module.
- [x] 3. Consolidate shared contract-test helpers and remove obsolete root files.
- [x] 4. Refresh documentation and deterministic generated artifacts if symbol ownership changes their presentation.
- [x] 5. Run full compatibility, parity, race, static, vulnerability, and privacy verification.

## Definition of complete

- The package graph matches `context.md` without import cycles or exposed test-only seams.
- Root `bpi` callers retain the documented constructors, options, account/error types, custom request helpers, and all domain accessors.
- All 27 domain folders contain their client implementation and relevant contract tests.
- No contract, fixture, mapping, or risk count changes.
- All offline repository checks pass.
