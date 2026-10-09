---
type: team-pr-handoff
context: plaintext-hint fix
---

# Authorized team PR handoff

Owner approved separate reviewed commits, feature-branch push and a PR in
`umats/pdfcpu` against current default `master`. No merge, tag, package
publication, migration, consumer integration or rollout authority.

## Commits and binding

- `05de1219`: reader fix and all 34 directly related source/regression paths.
- `6e083295`: separate native32 unsigned-permission test compilation fix.
- Approved specs/conventions/ignore metadata follows in a separate docs commit.

Accepted source hashes match round13 candidate; final original whole-manifest
check differs only in `CONVENTIONS.md` because of the explicit standing owner
coverage-percentage waiver. This known metadata delta does not change Go/module
bytes or invalidate code-bound runtime/security/acceptance evidence.
Receipt: `.pi-herdsman/hint-publication/manifest-delta.txt`.

Stage only enumerated reviewed paths. Exclude the three generated untracked
PDFs, all `.pi-herdsman/` runtime evidence and private input. No broad staging
or sample cleanup authorized. No attribution trailers added.

## Gates

Native Preflight, focused races, synthetic UAT, focused386 runtime and ARM
cross-build evidence retained. Aggregate coverage tests exit 0, with measured
74.8730% overall / 74.2648% core; percentage thresholds are explicitly
owner-waived as a standing project policy, not passed.

Mechanical module API comparison against both fork baseline HEAD and published
fork v0.15.1 passes, reporting only three ReadContext field additions. Unkeyed
struct literals remain a compatibility caveat; no tag/version selection implied.
Module hygiene passes without project dependency changes.

Use explicit `--repo umats/pdfcpu` for gh commands: implicit gh selected the
parent. Origin fork master was fetched and matches the original baseline.
Fork latest release is v0.15.1, not the parent's v0.16.1. No existing feature PR
was found before creation. Push/PR/CI receipts will be captured privately;
CI is not yet claimed green. No merge or release until separately authorized.

Production input remains absent; no Annex F, production, full386-suite or
performance claim. Two historical consumer failures remain unexplained.
