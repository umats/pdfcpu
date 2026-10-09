# Paused checkpoint — plaintext linearization hints

Paused at owner request on 2026-10-08T17:42:51Z.
Branch: fix/plaintext-linearization-hints.
No active managed agents; no implementation or review resumes automatically.
No commits, push, publication, version, private sample or delivery actions.

## Latest product state

Owner explicitly approved internal nonstandards compatibility in existing
relaxed mode, with broader object-stream support. Strict/authentication,
ordinary encrypted content and resource limits must remain unchanged.

Original primary public metadata/read/decrypt RED and first implementation
recorded in `.pi-herdsman/hint-repro/` and `.pi-herdsman/hint-fix/`.
Round1 independent reviews both failed; four findings were corrected and public
regressions proved GREEN. Round2 independent reviews also failed:

- result:hint-round2-a#1: 84/100; null duplicates, remaining xref repair paths,
  repeated structural inspection notices.
- result:hint-round2-b#1: 82/100; null duplicates and classic trailer Size repair.

All round2 findings now have local corrections and regression evidence:

1. Model dictionary parser tracks encountered names independently of retained
   non-null values for provenance, without changing normal parse acceptance or
   notices. Null-first/null-last H/Length/S and ordinary/trailer/compressed
   cases reject plaintext compatibility and retain encrypted controls.
2. Effective trailer Size mismatch marked before candidate discovery;
   prefixed/signed noncanonical classic fields and surplus xref-stream entries
   disqualify compatibility/deletion without changing normal acceptance.
3. Structural lazy-reference inspection is notice-free while retaining parser
   provenance. Candidate Flate normalization runs on a dictionary copy without
   emitting an extra repair notice; ordinary loading still emits its notice.

Latest changed source/test scope:
internal/contextutil/parse.go; pkg/pdfcpu/model/{parse.go,context.go};
pkg/pdfcpu/{read.go,readHints.go,write.go,readHints_test.go};
pkg/api/{plaintext_linearization_hint_test.go,plaintext_hint_review_test.go}.
Metadata: .gitignore (owner-requested .pi-herdsman ignore), bug/registry/state,
verification documents. No dependencies/vendor/CI/lint-policy edits.

## Exact latest validation (all exit 0)

Logs `.pi-herdsman/hint-fix-round3/`:

- null-red.log / null-green.log: public TestHintReviewNullDuplicates RED1/GREEN0.
- xref-red.log / xref-green.log: public TestHintReviewUnreliableXRef RED1/GREEN0.
- notices-red.log / notices-green.log: public TestHintReviewInspectionNotices
  RED1/GREEN0.
- `go test -count=1 ./pkg/api ./pkg/pdfcpu ./pkg/pdfcpu/model ./pkg/cli
  ./internal/contextutil -run
  'Hint|OrdinaryAES|Linearization|AESV2CryptFilter|EncryptedXRef|Metadata|Crypt|Step24a|ReturnedNotice'`
  focused.log exit0.
- `go test -count=1 ./... && go vet ./... && go build ./...`
  preflight.log exit0, PREFLIGHT_EXIT=0.
- Same focused command with `-race`: race.log exit0.
- gofmt on latest changed paths completed before validation.

Final round3 gopls/whitespace and fresh independent dual review NOT performed
before pause. Round2 gopls/whitespace were clean before the latest edits.
No current independent pass; do NOT declare accepted/fixed or publish.

## Housekeeping blocker

Latest Preflight regenerated tracked pkg/samples outputs and the three known
new sample PDFs (basic/FormDemo.pdf, bookmarks/bookmarkTree2Levels.pdf,
bookmarks/bookmarkTreeImported.pdf). Initial sample subtree status was clean
(`hint-fix-round3/samples-before.status`). Scoped restoration attempted during
pause housekeeping FAILED because `.git/index.lock` exists. No lock was deleted,
no fallback file overwrite attempted, no generated sample cleanup completed.
`samples-after.status` records remaining changes. Lock ownership/staleness must
be checked safely on resume; do not delete an active Git lock. Then restore
only proven test-generated sample output paths, preserving all product/owner
changes. This is a Git housekeeping blocker, not a failed test gate.

## Resume only after owner requests

1. Read this checkpoint, specs/state.yaml and the active bug document. Review
   current worktree without resetting/stashing/committing anything.
2. Resolve lock housekeeping safely and restore only known generated samples.
3. Run final gopls/whitespace/YAML checks on current paths; preserve results.
4. Dispatch fresh independent dual reviewers for round3 on current source and
   these exact validation logs; neither sees the other's current report.
5. Address all findings; max5 total review rounds. Current completed rounds=2,
   next round=3. No commit/merge/release/consumer rollout without approval.

Limits: production attachment unavailable; synthetic fixtures not proven
AnnexF-conforming. Unsupported hint shapes fail closed. Complete PDF2.0 text
unverified; owner waived standards compliance. Two historical pdfrelay failures
remain unexplained. New exported ReadContext tracking fields may affect unkeyed
external literals; no exported function/interface/config setting changed.
