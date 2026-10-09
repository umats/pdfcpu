# Plaintext hint compatibility — round 4 validation

Round 3 failed the dual-review gate: reviewer A 75%, reviewer B 70%.
Both found xref ambiguity/repair gaps, repairing header identity, and wrapper
recursion-depth inflation. All three have local corrections and regressions;
independent acceptance is pending.

## Corrections and checks

- Same-section classic/stream duplicate assignments now deny authorization.
  Legitimate older-revision precedence remains. Object-zero and free-list
  normalization mark repair provenance without changing ordinary acceptance.
- Hint authorization checks bounded literal headers at the claimed byte start.
  Ordinary repairing header parsing is unchanged.
- Structural wrappers add no PDF nesting depth. Compressed members are
  independent origins. Strict/relaxed encrypted/unencrypted boundary controls
  preserve accepted depth limits.
- Public xref, header, and depth regressions each recorded behavioral RED exit 1
  and GREEN exit 0. Logs and exact commands:
  `.pi-herdsman/hint-fix-round4/handoff.md`.
- Baseline and final full fresh Preflight passed:
  `go test -count=1 ./... && go vet ./... && go build ./...`.
  Intermediate nil-Read unit-context panic was corrected, not dismissed;
  failing log and focused GREEN are retained.
- Focused compatibility tests and focused race checks passed. gofmt,
  `git diff --check`, and YAML validation passed. Generated samples restored;
  final sample subtree clean, no lock removed, owner changes preserved.
- Parent inspected correction diff and public regressions. Active LSP checked
  all eleven Go paths; no primary errors/warnings. Five auxiliary pattern
  warnings remain false positives: guarded initialized map and non-GORM PDF
  loops. Explicit gopls check exit 0; `parent-gopls.log` is empty.
- Scope audit: no dependencies/vendor/CI/lint changes, no source secrets, no
  unrelated redesign. Tests exercise public behavior and preserve controls.
  Additional free-entry snapshot is bounded by the existing xref inventory;
  no performance improvement claimed or benchmark added.

## Acceptance and limits

Next: fresh independent dual review round 4. Owner authorized up to ten total
rounds if needed; zero must-fix and both scores at least 94% still required.
No commit, push, tag, merge, package publication, or consumer rollout authorized.
Synthetic fixtures are not Annex F-proven; production sample unavailable;
two historical consumer failures unexplained. No standards-compliance claim.
Three exported ReadContext tracking fields retain the unkeyed-literal caveat.
