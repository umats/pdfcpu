# Plaintext hint compatibility self-audit

- Scope: owner-approved relaxed-mode internal compatibility, broader object
  streams included; no normative exemption claimed.
- PASS conventions: feature branch; no commits/publication/dependency/vendor/CI
  or lint-policy edits; planning artifacts in specs; generated samples restored.
- PASS hygiene: source/test bodies and tracked reader/writer/model diffs read;
  original RED preserved; implementation full fresh Preflight/race/gopls logs
  report exit 0; no source secrets/private payloads.
- PASS initial correctness/security audit: independent copied cipher probes;
  structurally pinned candidate identities; reference walk includes lazy objects
  and trailers; ordinary/strict paths preserve encryption; writer deletion no
  longer trusts offset equality. Final adversarial acceptance requires review.
- PASS initial test audit: public read/metadata/decrypt/output variants, aligned,
  overflow/encrypted/object-stream, strict/wrongpassword, hostile claims and
  writer-preservation coverage. Internal tests cover reference provenance,
  dictionaries, limits, unreadable origins and byte preservation.
- PASS dependency/performance hygiene: standard library/existing parser only;
  per-read candidate inventory; no global helper state or external test tools.
  Candidate discovery probes first 1KiB; reference auditing limited to files
  containing supported candidates. No benchmark improvement claimed.
- Compatibility note: three fields were added to exported model.ReadContext for
  internal tracking. Unkeyed external literals/exact shape may be affected;
  no exported function/interface/configuration signature changed. No release
  compatibility tooling downloaded or version proposed.
- House-style exception: new source is 452 lines, public tests 547, internal
  tests 340; generic skill size/function heuristics are not project conventions.
  Do not mechanically split or add abstractions solely for line limits.
- Churn priority: read.go first (25 commits/90days). Churn helper printed ranking
  then exited 141 (bounded head pipeline), not a product Preflight failure.
- Limits: fixture layout and hint table are not proven Annex F-conforming;
  private production sample unavailable. Unsupported shapes remain fail-closed.
  Complete PDF2.0 standard not retrieved; owner waived standards requirement.

## Resume audit — 2026-10-09

- Owner explicitly resumed verification, not publication.
- PASS cleanup: sample status matched the saved checkpoint exactly. Restored
  446 tracked generated outputs and removed the three known generated PDFs.
  No Git lock existed or was removed; sample subtree is clean.
- PASS latest source inspection: reviewed reader/writer/model diff, hint helper,
  parser provenance helper and public round2 regressions. Null names are tracked
  before omission; repaired xrefs deny compatibility; inspection stays quiet.
- PASS final checks: gopls check exit 0, empty gofmt output, git diff --check
  exit 0, YAML validation OK, fresh focused tests exit 0. Logs are in
  `.pi-herdsman/hint-fix-round3/resume-{gopls,focused,cleanup}.log`.
- PASS diagnostics triage: no LSP errors or warnings. The auxiliary nil-map
  warning is false: seen is initialized and written under the same provenance
  guard. GORM warnings match PDF structural loops, not database queries.
  Spelling and modernization advisories are not acceptance failures.
- PASS prior evidence retained: unchanged product source has full fresh
  Preflight and focused race exit 0 before pause. Neither suite was repeated
  on resume; only focused tests and final diagnostics were repeated.
- Scope/security/performance/test checks retain the initial audit above;
  final adversarial acceptance remains the independent dual-review gate.
- Churn ranking again printed then exited 141 from its bounded head pipeline;
  no product gate failed. No heuristic-only splitting or new dependency added.

Next gate: fresh independent round3 read-only correctness/security reviews.
No final fix acceptance, package version, commit or consumer rollout authorized.
