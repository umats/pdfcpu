---
type: release-preparation
context: verified plaintext-linearization-hint compatibility fix
---

# Release preparation only

Owner confirmed completion of Phase 5 and entry into Phase 6 preparation.
Owner subsequently authorized reviewed separate commits, feature push and
team PR in `umats/pdfcpu` against `master`. Merge, tag, publish, branch migration,
consumer update and rollout remain unauthorized. Initial uncommitted baseline:
`fix/plaintext-linearization-hints`, HEAD
`870fbbf069d1c9501bb27b019d01bd80d1cb1d23`.

## Verified candidate

Current source binding and tests:
`plaintext-hint-round13-verification.md`.
Acceptance: `plaintext-hint-round13-acceptance.md` (A/B 100%, zero findings).
Completeness: `plaintext-hint-completeness.md`.
Synthetic behavioral UAT accepted; scoped security PASS.

## Proposed commit groups (drafts, not executed)

1. Reader fix and all directly related regression tests:
   `fix(reader): handle verified plaintext linearization hints`.
   Includes reader/writer changes, model provenance, internal contextutil and
   untracked hint tests. Do not omit untracked sources when staging later.
   Body: narrow relaxed compatibility for verified unreferenced hints;
   preserve authentication, normal cipher handling, limits and output content.
2. Separate discovered test-only fix:
   `fix(test): make unsigned permission checks compile on native32`.
   Scope: `pkg/pdfcpu/crypto_test.go`, its separate bug record and evidence.
   Preserve native32 skip and native64 assertions.
3. Approved workflow/evidence metadata:
   `docs(specs): record hint compatibility verification and release gates`.
   Review `.gitignore`, `CONVENTIONS.md` and specs changes explicitly;
   never use broad staging. Owner requested ignoring runtime evidence and
   forbidding provider-rejection/no-retry metadata.

Formatting-only changes in existing draft tests belong with their tested fix.
Do not invent retrospective RED commits: recorded source-bound RED/GREEN logs
are the owner-approved substitute for the generic two-commit workflow.
No attribution trailers proposed.

Default Conventional Commit analyzer mapping would treat fix as patch and docs
as no release. This is not a version approval or a compatibility guarantee.
ReadContext gained exported fields, which can break external unkeyed literals.
Before selecting a public release version/title, determine the intended v0
compatibility policy and perform the mechanical API comparison. If classified
as a breaking release, use `fix(reader)!` plus an explicit BREAKING CHANGE footer.
Do not reuse the historical fork-readiness version as an approved fix version.

## PR body draft

### Summary
<!-- bigpowers-provenance: agent-generated -->

- Add narrowly scoped relaxed-mode handling for verified plaintext hints.
- Require bounded physical origins, exact xref identities and complete incoming
  reference coverage before compatibility or full-write deletion authority.
- Preserve cipher-first processing, authentication, strict mode, limits,
  ordinary content and normal encryption for incremental replacements.
- Include separate native32 test compilation repair.

### Verification

- Native full tests/vet/build: PASS, fresh contiguous current-candidate run.
- Focused race and contextutil race: PASS.
- Package-based386 hint/permission regressions: PASS, actual execution.
- Linux/ARM cross-build: PASS, no ARM runtime claim.
- Formatting/whitespace, dual acceptance/security, synthetic UAT: PASS.
- Coverage: 74.8730% overall / 74.2648% core; percentage gates OWNER-WAIVED.
- Baseline mechanical API comparison and module hygiene: PASS, struct caveat.
- Published fork v0.15.1 mechanical API comparison: PASS, struct caveat.
- Tag-version decision and PR CI: still pending.

### Compatibility limits

Internal compatibility exception, not normative exemption or Annex F compliance.
Production sample absent; two historical consumer failures unexplained.
Optional early-xref-stream and performance coverage deferred.
No full386-suite claim. Exported-context shape caveat remains.

## Gates before any integration

- Commits, push and team PR explicitly authorized. Merge/tag/publication
  and rollout still require separate approval.
- Explicit `gh repo view umats/pdfcpu` confirms master as fork default;
  initial unqualified gh selected the parent. Use `--repo umats/pdfcpu` for
  PR/CI commands. Authorized origin fetch confirms unchanged master baseline.
  No migration, rename or remote edit; target main only after migration.
- Owner standing policy waives aggregate percentage release coverage gates
  (generic 80% overall / 95% core). Measurements remain honest; no test,
  security, CI, API or integration-authority waiver. Evidence:
  `plaintext-hint-release-gates.md` and the owner-waived gate gap.
- Owner-authorized pinned apidiff tooling download completed privately;
  recorded-baseline module comparison reports only three compatible fields.
  Fork v0.15.1 comparison also passes. Struct-literal caveat and tag-version
  policy remain; no version tag is authorized.
  Offline `go mod tidy -diff` passes; project module files unchanged.
- Preserve separate discovered-fix commit intent; inspect staged paths and all
  untracked sources. Exclude generated sample PDFs and `.pi-herdsman/` outputs.
- Refresh Preflight for any changed candidate and require green PR CI before
  merge. Only remote default-branch metadata was queried; no PR or CI
  status was queried or inferred.
- Complete applicable release traceability gate using the approved bug scope;
  no story-matrix script PASS is claimed. Take snapshot before a real cut.

Three retained generated untracked sample PDFs remain present, restored exactly.
They are not product fixtures for staging or release; no deletion authorized.
Phase 6 is authorized through team PR and CI checking only;
not complete, merged or released.
