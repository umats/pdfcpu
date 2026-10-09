---
type: review-response
context: BUG-2026-10-08-plaintext-linearization-hints
---

# Round 11 acceptance: FAIL

Evidence: `result:hint-accept11a#1` (100%, no findings) and
`result:hint-accept11b#1` (93.75%, one must-fix, security FAIL).
The dual AND gate fails despite prior terminal checks passing.

## Correction plan for round 12

Must-fix: object-stream prefix audit uses a positive-length Flate decode
request which does not enforce the configured decode limit. An attacker can
supply a large first-member gap of compressed whitespace, causing eager
expansion even for an ordinary PDF with no hint candidate.

Bound all private object-stream prolog/prefix/lookahead requests before
allocation against the effective decode budget. Use checked arithmetic for
relative-plus-First offsets and lookahead additions. Preserve zero/default
and explicitly unlimited budget semantics; do not change public decoder APIs
or silently impose a new ordinary-reader policy. Budget exhaustion must return
an error wrapping `filter.ErrDecodeLimitExceeded`, not merely mark provenance
ambiguous or authorize compatibility fallback.

Add one vertical RED/GREEN regression slice: small Flate object stream with
small First and a whitespace gap beyond a small configured budget. Assert the
terminal sentinel through the private decoder and public reader, with no-hint
input and in-budget controls. Cover lookahead/offset arithmetic boundaries.
Sweep the three positive-length calls in object-stream auditing so a related
prolog/lookahead request cannot retain the bypass. No giant allocation or
new performance feature required.

Executor owns source/tests and candidate-bound execution logs. Capture RED
and GREEN manifests contemporaneously. Fresh focused/full Preflight/race,
gofmt and whitespace checks remain required. Lead owns acceptance/specs and
bug-scoped completeness. No commit, stash, production inputs, dependency,
vendor, CI, release, consumer or rollout action. Preserve existing generated
untracked outputs with a pre-run backup when running the full suite.

Round 12 remains within the owner-authorized review cap of 20.
Phase 5 is incomplete. Production evidence and Annex F conformance remain
unavailable; two historical consumer failures are unexplained. Deferred
performance/layout coverage and exported-context compatibility caveat remain.
