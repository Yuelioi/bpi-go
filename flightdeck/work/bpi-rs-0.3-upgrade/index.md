# bpi-rs 0.3 upgrade

**Status:** Finished

## Goal

Upgrade bpi-go from the locked bpi-rs 0.2.4 baseline to the tagged bpi-rs 0.3.0 release, preserving the idiomatic Go API while matching the promoted contracts, response compatibility fixes, recovery behavior, and offline evidence of the new Rust baseline.

## Current

bpi-go now targets bpi-rs tag `v0.3.0` at `94bcf43e46d6e11b55ec4d36d4848692cecd4213`. The source lock and mapping/contract evidence identify that commit, and an exact detached local clone verified unchanged parity totals: 27 domains, 206 promoted contracts, 432 fixtures, and 642 locked files.

The applicable 0.3.0 response changes are implemented with focused passing tests: semantic API errors outrank incompatible error payloads while explicit compatible `IntoData` remains available; Dynamic detail IDs accept string, number, or null; video FLAC audio is an optional single stream and nullable backup lists remain safe; Bangumi/Cheese use current-quality `durl`; Bangumi totals accept `-1`; and compact `coin`/`play` stats decode correctly. Go's existing zero-value and raw-message models already tolerate the other upstream optional-display and pendant shapes.

## Next

None.

## Progress

- Audited the complete two-commit upstream delta from the former 0.2.4 lock through `v0.3.0`, including source, embedded regressions, release notes, and live-validation findings.
- Added representative offline regressions for every Go-applicable response shape and confirmed they failed before the implementation changes.
- Refreshed source metadata against an exact detached `v0.3.0` clone without changing the sibling Rust checkout.
- Passed focused tests for `client`, `video`, `dynamic`, `bangumi`, and `cheese`.
- Passed exact source generation/check and the offline parity audit with unchanged contract and fixture counts.
- Updated README, changelog, migration, Probe, parity, release, and generated API-index documentation, including both Go-facing breaking field changes.
- Passed `gofmt` drift checking, `go mod tidy`, `go vet ./...`, `go test -count=1 ./...`, API-index checking, parity/privacy audit, and staticcheck v0.7.0.

## Residual verification limits

- `go test -race -count=1 ./...` cannot compile `runtime/cgo` because the installed Scoop GCC 13.2 intrinsic headers contain binary corruption. A second compatible C compiler is not available; this is independent of the repository changes.
- `govulncheck` on local Go 1.25.12 finds five reachable standard-library advisories, all reported fixed in Go 1.25.13. The repository has no third-party runtime dependency, and CI uses the latest `1.25.x` patch.
- Two temporary tag clones remain under the operating-system temp directory because the execution environment rejected recursive cleanup; no temporary source or Probe output exists in the repository.

## References

- [Prior completed port](../bpi-rs-go-port/index.md)
- [Package layout record](../package-layout/index.md)
- [bpi-rs 0.3.0 changelog](../../../../bpi-rs/CHANGELOG.md)
