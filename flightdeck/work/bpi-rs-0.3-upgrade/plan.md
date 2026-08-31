# Plan

1. **Inventory bpi-rs 0.3.0 delta** — Complete
   - Compare source, contracts, fixtures, generated metadata, documentation, and tests from the locked 0.2.4 commit through tag `v0.3.0`.
   - Separate changes already represented in Go from required code, evidence, tooling, and documentation updates.
2. **Refresh locked parity evidence and tooling** — Complete
   - Update source locks, promoted contracts, fixtures, generated metadata, and source-generator assumptions against the tagged commit.
3. **Implement Go behavior and model compatibility** — Complete
   - Port every applicable 0.3.0 response-model and envelope-decoding change with focused regression tests.
4. **Update public documentation and release notes** — Complete
   - Refresh API/parity docs, migration guidance, README claims, and changelog for the new baseline and any Go-facing changes.
5. **Verify and finish the Work** — Complete
   - Run formatting, unit and contract tests, parity audits, source-generation checks, static analysis, race checks, vulnerability checks, and privacy-sensitive repository checks as applicable.
   - Formatting, tidy, vet, full tests, generated-document checks, parity/source checks, staticcheck, and privacy-sensitive audits pass.
   - Race execution is unavailable because the installed Scoop GCC headers are corrupt; govulncheck reports only Go 1.25.12 standard-library findings fixed by Go 1.25.13.
