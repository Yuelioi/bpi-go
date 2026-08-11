# Package layout reorganization

**Status:** Finished

## Goal

Make the repository navigable by collocating every domain's parameters, models, client implementation, and tests in that domain directory, while preserving the established root `bpi` interface and all 206 promoted-contract behaviours.

## Current

The package-layout reorganization is complete. The root contains only the public facade and root-level examples/parity checks; shared behavior lives in `client/`, and all 27 domain implementations and their contract tests are collocated with their models and parameters.

## Execution

All plan items and the activity proof Slice are complete. The public root call shape and error/type identities were preserved through aliases and wrappers.

## Next

None. Use [context.md](context.md) and [plan.md](plan.md) as the architectural record for future package changes.

## Progress

- Audited root files, exported symbols, domain imports, private request helpers, and test-helper coupling.
- Rejected a simple file move because Go packages cannot span directories and it would create a root/domain import cycle.
- Selected a root facade over a reusable client module so callers retain `client.Activity().Info(...)` while maintainers get domain locality.
- Extracted the shared client foundation and preserved root construction, options, envelope helpers, error identity, and accessors.
- Moved every domain client and all promoted contract tests into the corresponding domain package.
- Consolidated common fixtures and assertions under `internal/contracttest` and removed obsolete root helpers.
- Verified formatting, package compilation, generated parity metadata, static analysis, race safety, vulnerabilities, and privacy checks.

## References

- [Completed parity Work](../bpi-rs-go-port/index.md)
- [Generated API index](../../../docs/api-index.md)
