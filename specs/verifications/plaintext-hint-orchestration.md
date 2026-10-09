---
type: verification-plan
context: plaintext-linearization-hint bug
---

# Standard orchestration resume

## Approved phase map

Owner selected **Verify gap loop**, then confirmed Phase 5 completion after
current-candidate verification, dual acceptance and completeness passed.
Phase 6 is now release preparation only; the existing bug loop is verified
through step 4. Earlier phase markers alone were not verification evidence. Do not start a new epic or fabricate stories from
skill metadata. The active scope/design is the existing plaintext-hint BUG file.

Mode: standard. Branch and recorded baseline match:
`fix/plaintext-linearization-hints`,
`870fbbf069d1c9501bb27b019d01bd80d1cb1d23`.
The candidate includes uncommitted and untracked source/tests. HEAD alone
cannot identify the source exercised by future verification.

## Gate order

| Gate | Current evidence | Required before advancement |
| --- | --- | --- |
| Nested bug correction | Round13 corrections and test-only native32 repair verified | PASS within approved synthetic scope |
| Runtime validation | Fresh native full Preflight/race and focused package386 pass | PASS, candidate-bound receipts |
| User-observable verification | Public flow evidence; owner accepts synthetic UAT | PASS, no production claim |
| Security and acceptance | Round13 A/B 100%, zero findings, scoped security PASS | PASS; bug-scoped completeness also passes |
| Phase transition | Technical gates pass; owner confirmed Phase 5 | PASS; Phase 6 preparation only |
| Release | Not authorized; version proposal is historical | Separate explicit commit/PR/release authorization; team-pr workflow |

Default verify depth is P1 because no active story risk field is supplied;
this does not waive the security scan for shared PDF/encryption handling.
The existing standards waiver is a compatibility-design choice, not a waiver
of tests, authentication, content preservation, limits, or acceptance gates.

## Evidence rules

- Bind each execution result to the exact candidate, including untracked source
  and relevant fixture changes; preserve a source manifest before/after a run.
- Record stdout/stderr and terminal exit code from one contiguous run. Never
  combine old baseline output with a newer candidate or synthesize a verdict
  from excerpts of separate runs.
- Empty diagnostics caches, builds, static reviews, and source inspection do
  not establish a passing regression suite.
- Failed gates reopen this bug; preserve reproduction evidence and use its
  correction loop. Do not create duplicate bug/story entries or mark done.
- Preserve owner changes and sample evidence. Any generated-file cleanup must
  be scoped to proven test outputs; no index-lock removal or broad reset.
- No production attachment/password access, private processing, SMTP, pdfrelay
  changes, publication, or rollout is authorized by this orchestration resume.

## Deferred work and limitations

Owner deferred optional early-xref-stream coverage and performance measurements.
No universal layout, performance, Annex F compliance, or production-fix claim.
Original production sample is unavailable; two historical consumer failures
remain unexplained; public ReadContext shape compatibility remains a caveat.

Product/planning YAML placeholders and the fork version index are historical,
not a new release specification. Referenced installed gates/checkpoints docs
are absent; the skill's available embedded gate definitions were consumed.
Do not edit installed tooling or restart product planning as unrelated work.

## Handoff

Current evidence: round13 verification, acceptance and bug-scoped completeness
reports under `specs/verifications/`. Owner confirmed Phase 5 advancement.
Release preparation is saved in `plaintext-hint-release-preparation.md`;
commit/PR/merge, tag/publication and rollout authority remain separate.
No release action performed.
