---
type: verification-evidence
context: BUG-2026-10-08-plaintext-linearization-hints
---

# Round 12 bounded decode verification

Executor: `result:hint-fix12#1`; detailed handoff and logs:
`.pi-herdsman/hint-fix-round12/handoff.md`.

All three private object-stream prolog/cut/prefix requests now use a checked
budget guard before decoding. Checked native-int offset/lookahead arithmetic
prevents overflow. Zero/default and negative/unlimited semantics are preserved;
excess returns a wrapped terminal decode-limit sentinel. No public decoder API,
filter implementation, dependency, configuration or consumer change.

## Candidate-bound execution

Logs under `.pi-herdsman/hint-fix-round12/`, each contiguous with exit code:

- `red.log`: exit 1, behavioral private/public no-hint gap failures and
  arithmetic sentinel failures; in-budget controls passed.
- `green.log`: exit 0, corrective regression suite including boundaries.
- `focused.log`: exit 0, hint/linearization regressions.
- `preflight.log`: exit 0,
  `go test -count=1 ./... && go vet ./... && go build ./...`.
- `race.log` and full supplement `context-race.log`: exit 0.
- `gofmt.log`: empty formatting listing; `whitespace.log`: exit 0.

RED and GREEN manifests were captured contemporaneously. Supplemental boundary
controls added after correction are not claimed as RED. Final product manifest
`product-final.sha256` SHA256:
`ab7c6048bfcdc34417707eaf92412a705152f5fc8f300d4b98a2523986d8f664`.
Executor stability receipt verifies no changes across validation; Lead ran
`sha256sum -c` on current source, exit 0 (`lead-manifest-check.log`).
`logs.sha256` binds execution logs. Lead whitespace check passes.

## Self-audit and limits

Correction inspection found no additional source issue: budget guard runs
before the underlying positive-length Flate request; all three calls use it;
terminal routing cannot turn limit exhaustion into provenance-only denial.
Regression coverage includes a small compressed-whitespace gap, public no-hint
reader denial, in-budget/default/unlimited controls and arithmetic boundaries.
Source hygiene and approved scope pass; dual acceptance/security remain pending.

Three inherited untracked sample PDFs were backed up and restored byte-for-byte.
Only proven initially clean tracked test outputs were restored. Round11's
historical output-regeneration and missing initial RED hashes are not erased;
round12 has its own complete candidate receipts.

Owner accepts synthetic UAT and selected bug-scoped completeness checks.
The no-hint fixture has zero pages and tests ReadContext, not conformance.
Production input and full Annex F proof remain unavailable. Two historical
consumer failures remain unexplained. Optional early-xref-stream coverage and
performance measurements are deferred; exported context shape caveat remains.
Phase 5 is not passed; no release actions or authorization.
