---
type: release-gate-report
context: verified plaintext-linearization-hint candidate
---

# Release gates: team PR authorized, merge/release pending

Phase 5 acceptance remains valid. Owner authorized separate commits, feature
push and team PR in `umats/pdfcpu` against `master`. No merge/tag/publication.
No product source changed during these checks.

## Coverage: OWNER-WAIVED percentage gate

Execution: `result:hint-release-coverage#1`; detailed receipt:
`.pi-herdsman/hint-release-gates/coverage/handoff.md`.
One aggregate `go test -count=1 -coverpkg=./... -coverprofile=... ./...`
run passes, exit 0. Unique statement blocks are weighted/unioned across
integration-package calls, not averaged or selectively excluded.

| Scope | Covered / statements | Measured | Gate |
| --- | --- | --- | --- |
| Entire aggregate profile | 48,931 / 65,352 | 74.8730% | 80%: OWNER-WAIVED |
| All `pkg/pdfcpu/...` core | 37,601 / 50,631 | 74.2648% | 95%: OWNER-WAIVED |

New-file coverage: contextutil parse 93.33%, physical origins 83.33%,
hint handling 84.57%. These subsets do not substitute for the release gates.
Source manifests unchanged; all 619 backed-up PDF/JSON hashes restored exactly.
This is native statement coverage, not branch/changed-line/platform coverage.
Owner explicitly directed that aggregate percentage release thresholds
are always waived for this project. CONVENTIONS.md and state record that
standing exception. Original measurement is not relabeled PASS; tests,
security/acceptance, CI, API and separate release approvals remain intact.
Disposition: `specs/bugs/BUG-release-coverage-gate.md` (owner-waived).
No release-readiness or publication claim.

## Mechanical API comparison: baseline PASS with caveat

Owner approved downloading pinned tooling into a private workspace.
Tool: `golang.org/x/exp/cmd/apidiff` at
`v0.0.0-20250620022241-b7579e27df2b`, installed under
`.pi-herdsman/hint-release-gates/api/tools/`. Project modules unchanged.
Compare module export data for recorded baseline HEAD
`870fbbf069d1c9501bb27b019d01bd80d1cb1d23` versus the current candidate.
Export/load/comparison logs under `.pi-herdsman/hint-release-gates/api/`.
Both exports and comparison exit 0. Tool reports three compatible additions:

- `model.ReadContext.RepairedHints`
- `model.ReadContext.RepairedXRef`
- `model.ReadContext.TrailerDicts`

No incompatible exported API change reported against this baseline.
This does not prove every external caller compiles: unkeyed struct literals
and exact struct-type assignments can break after field additions.
Tag-version policy remains pending. Published fork v0.15.1 comparison
also passes (details below); no release version inferred.
Earlier offline tool attempts failed dependency lookup; no API PASS was claimed
until the approved pinned download and successful comparison.

Initial unqualified gh queries selected the parent `pdfcpu/pdfcpu`, not
origin fork `umats/pdfcpu`. Those upstream v0.16.1 results are not fork release
evidence. Explicit `--repo umats/pdfcpu` reports latest fork v0.15.1; origin
remote tag resolves to `57dc6e3f231abdb7e1ee07e9c1fd9aa0098ee95d`.

Exported API comparison of that published fork commit versus current candidate
exits 0, with the same three compatible field additions and no incompatible
report. Logs: `published-export.log`, `published-apidiff.log` in the API receipt
directory. Struct-literal caveat remains. PR can proceed without a version tag;
publication still requires separately approved version policy.

## Module hygiene: PASS

Offline `go mod tidy -diff` exits 0 with no diff. `go.mod`/`go.sum` unchanged.
Final accepted candidate manifest check passes after API/hygiene checks.
No dependency version, vendor or toolchain requirement change.

## PR base: confirmed, no migration

Explicit `gh repo view umats/pdfcpu` confirms fork default `master`.
Origin is `umats/pdfcpu`; upstream is `pdfcpu/pdfcpu`. Use explicit gh repository
selection for every PR/CI action. Origin master was fetched after publication
authorization and still equals the recorded baseline. No branch rename,
remote edit or migration. No existing PR on this feature branch was found.

## Remaining decisions

- Aggregate coverage percentage gate explicitly owner-waived by standing
  project policy; retain measurements and all non-percentage gates.
- Decide v0 compatibility/release version policy with exported-struct caveat.
- Commit/push/team-PR authority explicitly granted; separate merge/tag,
  publication and rollout authority remains unapproved.
- Refresh tests for changed candidates and verify green PR CI before merging.

Reader source and native32 test repair committed separately as 05de1219 and
6e083295 after authorization; docs commit/push/PR follow. No merge/tag/publish.
Production sample, Annex F, optional layout/performance and historical failure
limits unchanged. Current progress: `plaintext-hint-publication.md`.
