---
type: verification-evidence
context: BUG-2026-10-08-plaintext-linearization-hints
---

# Round 11 correction verification

Executor evidence: `result:hint-fix11#1` and
`.pi-herdsman/hint-fix-round11/handoff.md`.
All four round10 findings are implemented, not yet independently accepted.

## Current-candidate terminal receipts

Each command has its own contiguous stdout/stderr and exit code 0:

- Focused hint/linearization tests: `focused-final.log`.
- Full `go test -count=1 ./... && go vet ./... && go build ./...`:
  `preflight.log`.
- Focused reader/API/model race: `race.log`.
- Full contextutil race: `context-race.log`.
- All 28 changed/untracked Go files: empty `gofmt.log`.
- Whitespace: `whitespace.log`.

Logs and manifests are under `.pi-herdsman/hint-fix-round11/`.
`logs.sha256` binds execution logs. Non-output candidate manifest:
`product-final.sha256`, SHA256
`218e78703b33af2b2e3e05c46ccfe4c88af3dbb8af5ea97f74fac98fa7e244a0`.
Executor verified stability across final runs; Lead independently ran
`sha256sum -c` against current files, exit 0 (`lead-manifest-check.log`).
Lead whitespace check also passes. Active LSP on five correction paths reports
no primary errors/warnings, with one optional range-loop modernization hint.
Auxiliary nil-map warning is false: both control branches initialize claims
before access. GORM warnings describe PDF loops with no database operations;
PDF terminology typos do not indicate source errors.

## Self-audit

PASS hygiene, scope, types and correction inspection: exact live generation
is checked only for physical authority; ordinary reader tolerance is unchanged.
One cancellable claims index replaces per-container table scans. Callback
errors propagate. Normal-cipher and precommit cancellation checks precede
repair/content/accounting and verified identity commits. No public/configuration,
dependency, vendor, CI, lint-policy, consumer or release change.

Regression evidence covers generation mutation, valid/encrypted controls,
deterministic claim traversal and cancellation, decoder integration sensitivity,
and strict verified/unverified finalization cancellation. Current tests pass;
security and independent dual acceptance remain required. Formatting changes
are restricted to changed/untracked files.

## Evidence limits

Initial length/cancellation RED runs did not capture candidate manifests.
One cancellation attempt and initial claims RED failed compilation, not behavior;
they are not behavioral proof. A corrected cancellation RED and length RED
show behavioral failure; decoder integration sensitivity has an exact temporary
candidate manifest. No retrospective hash or two-commit claim is made.

Full Preflight regenerated three retained untracked sample PDFs. Their original
bytes were not backed up; preservation is not claimed. All 446 initially clean
tracked sample outputs were restored. No private inputs or broad cleanup.

Owner accepted synthetic behavioral evidence and selected bug-scoped
completeness checks. Original production input remains unavailable; fixtures
are not fully Annex F conforming. Two historical consumer failures remain
unexplained; performance measurements and optional early-xref-stream scenario
are deferred. Exported ReadContext unkeyed-literal caveat remains.

Phase 5 remains incomplete. No commit, tag, PR, merge, publish or rollout.
