---
type: acceptance-report
context: BUG-2026-10-08-plaintext-linearization-hints
---

# Round 13 dual acceptance: PASS

| Independent reviewer | Items | Must-fix | Should-fix | Score | Acceptance | Security |
| --- | --- | --- | --- | --- | --- | --- |
| A: `result:hint-accept13a#1` | 16 | 0 | 0 | 100% | PASS | PASS |
| B: `result:hint-accept13b#1` | 16 | 0 | 0 | 100% | PASS | PASS |

Both fresh read-only reviewers inspected current source, connected reader,
parser, crypto, API and writer flows, and regression implementations.
They reviewed supplied execution receipts; neither executed tests or edits.
Both verified the round13 provenance corrections and test-only native32 repair.
No findings remain. Both scoped security verdicts have confidence 8/10;
no concrete high-confidence vulnerability identified. This is bounded review,
not a universal security guarantee.

Current candidate and terminal receipts are bound in
`plaintext-hint-round13-verification.md`. Lead manifest check and unchanged
candidate establish applicability. No source edits after reviewer dispatch.
The owner-authorized 94% zero-must-fix AND gate is satisfied at round 13/20.

Owner accepted synthetic behavioral UAT and bug-scoped completeness checks.
Completeness evidence: `plaintext-hint-completeness.md`.
Phase transition requires owner confirmation; release authority is separate.

## Limits retained

Approved relaxed-mode compatibility exception, not normative exemption or
standards compliance. Original production input unavailable; fixtures not fully
Annex F conforming; two historical consumer failures remain unexplained.
ARM is build-only; actual 386 evidence is focused package execution, not a
full386-suite claim. Performance and optional early-xref-stream coverage remain
deferred. Exported ReadContext additions retain unkeyed-literal compatibility
risk. Historical incomplete RED hashes and regenerated outputs are disclosed
in earlier receipts, not retroactively repaired. No commit, release or rollout.
