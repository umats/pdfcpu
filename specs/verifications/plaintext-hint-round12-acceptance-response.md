---
type: review-response
context: BUG-2026-10-08-plaintext-linearization-hints
---

# Round 12 acceptance: FAIL

`result:hint-accept12a#1` and `result:hint-accept12b#1` each score
93.75%, with one distinct must-fix and scoped security FAIL.
The bounded decode correction is accepted by both on source inspection;
overall acceptance fails the dual AND gate.

## Round 13 correction plan

1. A: final recovered `stream` marker may originate inside a comment or
   non-top-level syntax. Audit lexical state at the final selected marker
   before prefix slicing, including ordinary and xref-stream parsing. Mark
   ambiguous origins unreliable while preserving historical ordinary-reader
   behavior. Test the reported nested DecodeParms `/Ignored /stream` plus
   `% >>stream` shape privately and through authenticated public AES reading.
   Plaintext must be denied; encrypted controls must have no deletion
   authority. A legitimate `/stream` name followed by a real marker must
   remain eligible. Include ordinary-origin concealed-reference control.
2. B: unsigned xref fields within int64 can narrow lossy to native int on
   supported 32-bit targets. Audit generation, container number and member
   index fields for lossless native-int round trips before revision skipping.
   Mark unreliable provenance, not a new ordinary parse rejection. Include
   five-byte 4294967296 generation/container/member-index fields, ignored
   older entries and in-range controls. Test private and public plaintext
   denial/encrypted no-deletion behavior on a runnable 32-bit target such as
   GOARCH=386; separately cross-build supported linux/arm. Do not claim an
   ARM runtime test from a build. Preserve encoded-width overflow checks.

One executor owns shared reader source and tests. Use vertical behavioral
RED/GREEN logs with contemporaneous source manifests. Do not invent results
for unavailable cross-target execution. Full fresh Preflight, focused/race,
formatting and whitespace remain required, then fresh dual acceptance.

No general parser rewrite, new API/config/dependency, vendor/CI/lint-policy,
commit/stash, private inputs, release or consumer change. Back up and restore
inherited generated untracked sample PDFs exactly. Lead owns specs and
bug-scoped completeness. Round 13 is within owner-authorized cap of 20.

Phase 5 is incomplete. Production input/Annex F proof remain unavailable;
two historical failures unexplained. Optional layout/performance measurements
and exported-context compatibility caveat remain unchanged.
