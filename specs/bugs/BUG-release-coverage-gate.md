---
type: gate-gap
context: plaintext-hint Phase 6 release
bug_id: BUG-release-coverage-gate
status: owner_waived
severity: high
scope: release-verification
---

# Release coverage gate not met

Coverage tests pass, but aggregate statement coverage does not meet the
installed release skill's 80% overall / 95% core thresholds.
Measured: 48,931/65,352 (74.8730%) overall; 37,601/50,631 (74.2648%)
for all `pkg/pdfcpu/...`. Profile instrumented all packages and includes
integration-test calls. No selective denominator pruning.
Owner subsequently directed: always waive this project's aggregate percentage
coverage release gates. This is an explicit project policy, not a measured PASS.

Evidence: `result:hint-release-coverage#1` and
`.pi-herdsman/hint-release-gates/coverage/handoff.md`.
Candidate unchanged; tests exit 0. This is a test-coverage gate gap, not a
new observed production failure. It does not invalidate scoped functional
acceptance. Percentage thresholds are now owner-waived, not met; other
release gates and explicit integration authority remain required.

## Disposition

Owner explicitly selected a standing project waiver rather than coverage
remediation. Recorded in CONVENTIONS.md, state.yaml and durable memory.
Measured percentages and original execution evidence remain unchanged.
No silent dismissal, assertion padding or source/test change.

Waiver is limited to aggregate percentage thresholds. Tests, regressions,
security/acceptance, CI, API checks, synthetic/production limits and separate
integration/publication approvals remain intact. No commit or release performed.
