# Round 10 static self-review

## Scope

Candidate corrections: nonstream lexical cuts in `read.go`, physical envelope
and epilogue proof in `readHintOrigins.go`, and its discovery integration.
Tests cover helper boundaries and public plaintext denial/encrypted controls
for complete epilogues, comment-contained origins, and stream-contained origins.

## Checklist

- Static correctness: no additional defect identified during self-review;
  independent review remains required.
- Security: no authentication bypass, limit loosening, dependency addition,
  private input, or secret introduced. Unknown origins deny hint authority.
- Data preservation: original encoded lengths determine physical boundaries;
  candidate proof does not decrypt or mutate stream bytes.
- Scope: no vendor, lint-policy, CI, public API, or consumer changes.
- Types/style: gofmt, vet, build, and primary LSP checks pass.
- Tests: helper and public regression source exists; none executed in this
  session. Runtime correctness and sensitivity are not established.
- Performance: per-object parsing and sorting add work to eligible discovery;
  object/gap budgets use existing limits. No performance measurement or claim.
- Clarity: proof is separate from ordinary parsing. The xref gap routine is
  longer than the generic style recommendation because it audits one envelope;
  avoid splitting it into speculative parser abstractions before review.

Historical checkpoint: static source hygiene was ready for read-only
feedback, not acceptance or release. Full Preflight, focused/race execution, and the independent acceptance
gate remain incomplete. Churn ranking identified `read.go` as the primary
hotspot; the ranking script emitted its list but exited 141, so its process
status is not claimed as passing validation.

## Current-candidate self-audit checkpoint

The earlier runtime limitations above are superseded by
`plaintext-hint-round10-runtime-verify.yaml`. Current focused, full Preflight,
focused race and full contextutil race runs pass. Owner accepted synthetic
public behavioral evidence. No source changes occurred across verification,
excluding regenerated sample outputs.

- PASS correctness: prior reviewed corrections plus current physical proof
  and finalization inspection reveal no additional defect. Current denial,
  compressed-length, encrypted/plaintext and output-preservation tests pass.
- PASS security: separate scoped scan `result:hint-security10#1` found no
  concrete issue; recorded in `specs/security/REVIEW.md`.
- PASS hygiene: current gofmt listing is empty; whitespace checks pass;
  tracked product diff secret-pattern check found no match.
- PASS scope/supply chain: no dependency, vendor, lint-policy, CI or consumer
  change; authentication and staged-write boundaries retained.
- PASS types/clarity: typed reader-local provenance and existing context
  seams; no new configurable fallback or speculative parser abstraction.
- PASS test coverage: public behavioral tests plus parser/provenance helper
  regressions; synthetic controls and focused race evidence retained.
- Performance limit: bounded inspection is source-reviewed, not benchmarked;
  owner deferred measurements. No performance claim.
- Style limit: existing large reader and long envelope routines exceed the
  generic skill size guideline. Follow existing package shape rather than
  perform an unrelated rewrite. Supplemental internal helper tests do not
  replace public behavior tests.

Self-audit is ready for dual independent acceptance review, not a completed
Phase 5 gate. Completeness tooling is still unresolved: installed wrappers
resolve their package directory, and this active bug has no story matrix.
Do not fabricate a passing verdict. No commit or release authority exists.
