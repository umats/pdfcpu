# Plaintext hint compatibility — round 6 validation

Round 5 failed: A81.25%, B80%. All five combined provenance gaps are locally
corrected. Independent acceptance remains pending.

- Invalid hex omission marks ambiguity without changing ordinary acceptance.
- Trailer prefix scanning audits discarded syntax, including skipped lines.
- Object-stream omitted prolog/member prefixes and index syntax/identity are
  audited under existing decode limits. Comments-only controls stay accepted.
- Single-element trailer ID repair records provenance before authorization.
- Encoded-input length reconstruction records repair before normalization;
  ordinary post-decryption length adjustment is not marked as input repair.

Public behavioral RED exit1 and GREEN exit0, exact commands and complete
finding dispositions: `.pi-herdsman/hint-fix-round6/handoff.md`.
Baseline and final fresh full Preflight, focused compatibility and focused
race passed. Samples safely restored; no lock removed; subtree clean.
Parent inspected prolog/index and encoded-length handling plus public tests.
Scope/security/test audit retains existing gates; no dependency/vendor/CI/lint
changes or unrelated API redesign. No performance improvement claimed.

Parent final checks: whitespace, formatting and specs YAML pass. Active LSP
checked thirteen Go paths: no primary errors/warnings; five auxiliary warnings
remain false positives for a guarded initialized map and non-GORM PDF loops.
Explicit gopls exit0; one advisory about fixture string concatenation remains.
No claim that its output is empty.

Next gate: fresh dual independent round6/10 review. No product edits during
review; no commit, push, merge, tag, publication or consumer rollout authorized.
Production sample unavailable, synthetic tables not Annex F-proven, two
historical failures unexplained, public ReadContext shape caveat unchanged.
