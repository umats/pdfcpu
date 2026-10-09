---
type: review-response
context: BUG-2026-10-08-plaintext-linearization-hints
---

# Round 10 acceptance: FAIL

Independent results: `result:hint-accept10a#1` and
`result:hint-accept10b#1`. Both score 81.25% (16 items, two must-fix
items and one should-fix item each). Neither meets the 94% AND gate.
Scoped security scan PASS does not supersede these concrete findings.

## Approved correction plan

1. Must-fix (A): physical ordinary-stream indirect `/Length` resolution
   loses generation identity. Before resolving, require a live target with
   matching generation, including generation-zero compressed targets.
   Deny hint authority on mismatch without tightening ordinary loading.
   Add public compressed-length generation-mutation plaintext denial and
   encrypted readability/no-deletion controls.
2. Must-fix (A+B): compressed-member auditing scans the full xref table
   once per container, even without hints. Build claims by container once
   under cancellation checks, then audit only each container's claims.
   Keep read-local state private; no exported/configuration seam.
   Add deterministic scaling/traversal and cancellation regressions plus
   multi-container identity controls. This is a concrete complexity defect,
   not the owner's deferred performance-measurement item.
3. Must-fix (B): deferred normal cipher processing can miss cancellation
   when repair is false. Check cancellation after normal processing and
   before committing content or verified deletion authority.
   Add strict encrypted and unverified-candidate cancellation regressions.
4. Should-fix (A+B): gofmt the two draft round10 files and check every
   changed/untracked Go source, not only the earlier subset.

Use vertical regression slices with RED/GREEN logs, not commits or stash.
Executor owns source/tests and runtime execution receipts; Lead owns specs,
acceptance and integration. Preserve old logs as candidate-bound evidence.
Run fresh focused tests, full Preflight and applicable race checks after fixes.
Then self-audit/security review and dual acceptance at round 11, under the
existing owner-authorized cap of 20. No production or release action.

## Completeness decision

Owner selected bug-scoped checks instead of creating story artifacts.
Map approved bug criteria to current tests and receipts. Story-only installed
wrappers target their own package root and are inapplicable here; no substantive
verification, security, acceptance or phase-confirmation gate is waived.

Phase 5 and bug acceptance remain incomplete. Production input is unavailable;
fixtures are synthetic, not fully Annex F conforming. Optional early-xref-stream
coverage and performance measurements remain deferred. Two historical consumer
failures remain unexplained; exported ReadContext shape caveat remains.
