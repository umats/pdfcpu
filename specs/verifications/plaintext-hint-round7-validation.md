# Plaintext hint compatibility — round 7 validation

Round 6 failed. Its five combined findings are locally corrected with executable
RED exit1 and GREEN exit0. Evidence and exact commands:
`.pi-herdsman/hint-fix-round7/handoff.md`.

- Shared encoded loader audits bounded actual stream postludes before
  decryption. Noncomment syntax or unreadable/overlong tails deny authorization
  without changing ordinary acceptance; cancellation remains terminal.
- Hint envelopes check a real byte beyond the claimed span; clipped physical
  terminator tokens no longer establish identity.
- Lossy xref field conversions are audited before revision precedence; effective
  public fields and ignored older records have sensitive regressions.
- Missing/noncanonical stream EOL handling records provenance in shared parsing.
- Consumed lexical whitespace and case normalization record repairs, with
  legitimate string/comment payload controls retained.

The previous executor retired safely without corrections; fresh executor reused
its unchanged passing baseline and completed the work. Final fresh full
Preflight, focused compatibility, focused race passed. Final formatting-only
edit retested. Samples restored safely; subtree clean; no lock removed.
Parent inspected shared postlude, loader, literal envelope and public tests.
Self-audit retains correctness/security/limits/writer gates for independent
review. No dependencies/vendor/CI/lint changes, unrelated API redesign or
performance claim.

Parent active LSP checked fifteen Go paths: no primary errors/warnings; same
five auxiliary false positives for guarded initialized map and non-GORM PDF
loops. Explicit gopls exit0; fixture concatenation advisory remains. Whitespace
and YAML pass; formatting evidence is in the implementer handoff.

Next: fresh dual independent round7/10 review. No product changes during review.
No commit, push, merge, tag, publication or consumer rollout authorized.
Production sample unavailable; synthetic tables not Annex F-proven; two
historical failures unexplained; public ReadContext shape caveat unchanged.
