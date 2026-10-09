# Plaintext hint compatibility — round 8 validation

Round 7 failed: A85%, B87.5%. All three combined findings locally corrected.
Executable RED/GREEN, exact commands and dispositions:
`.pi-herdsman/hint-fix-round8/handoff.md`.

- Actual nonstream terminators are audited before truncation, including physical
  boundaries, non-PDF whitespace, junk recovery and name-marker confusion.
- Raw xref framing and consumed epilogue bytes/separators are audited before
  normalization and revision skips. Ordinary parsing remains unchanged.
- Keyword/reference prefixes record missing physical token boundaries.
  Compressed-member artificial EOF equivalents check actual decoded endpoints.
  Valid separated tokens, names, strings and comments retain acceptance.

Fresh baseline/final Preflight, focused compatibility, focused race, formatting
and whitespace passed. Samples safely restored and clean; no lock removed.
Parent inspected boundary predicates/public regressions and reused executor's
connected defect sweep and self-audit. No public interface/dependency/vendor/
CI/lint changes, unrelated redesign or performance claim.

Active LSP checked eighteen Go paths: no primary errors/warnings; same five
auxiliary false positives for guarded initialized map and non-GORM PDF loops.
Explicit gopls exit0 with fixture concatenation advisory. YAML/whitespace pass.

Next: fresh dual independent round8/10. No product edits during review; no
commit/push/tag/merge/publication or consumer rollout authorized. Production
sample absent; synthetic Annex F proof absent; two historical failures and
public ReadContext shape caveat remain unchanged.
