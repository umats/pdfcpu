# Plaintext hint compatibility — review corrections and validation

## Recovery and ownership

The implementer continuation ended without a final handoff.
The parent inspected resolved partial work and completed ordinary local PDF
reader/writer correctness regressions. Product work is now resolved pending
fresh independent review; no source editing while those reviews are live.

## Round 1 findings addressed (pending independent acceptance)

1. Xref rebuild: pending hint inventory invalidated before bypass/retry so
   normally loaded replacement entries are not decrypted twice. Public primary
   and overflow encrypted controls retained; existing issue1 RED/GREEN logs.
2. Parser ambiguity: read-local atomic provenance detects duplicate dictionary
   keys and fallback parsing, including lazy compressed origins. Notice-free
   candidate probes use separate provenance; ordinary parser acceptance and
   notice behavior remain. Public ordinary/trailer/compressed/candidate cases
   include duplicate overwritten refs and parser fallback. Issue2 RED/GREEN.
3. Unreliable xref: Size/Index correction, skipped live offsets and malformed
   classic entry fields/extra entries disqualify compatibility and deletion;
   ordinary encrypted controls continue. Public issue3 RED/GREEN captured.
4. Incremental output: track only actually plaintext-repaired hints in
   ReadContext.RepairedHints; include them in modified objects for ordinary
   encrypted replacement writing. Original file prefix remains untouched.
   Public annotation increment primary/overflow reopens and validates with
   content intact; strict read succeeds without plaintext retry. Issue4 RED/GREEN.

## Changed source/test paths

- internal/contextutil/parse.go (new internal provenance)
- pkg/pdfcpu/model/parse.go (record duplicate/fallback without global rejection)
- pkg/pdfcpu/model/context.go (reader tracking fields)
- pkg/pdfcpu/read.go
- pkg/pdfcpu/readHints.go (new)
- pkg/pdfcpu/write.go
- pkg/api/plaintext_linearization_hint_test.go (new)
- pkg/api/plaintext_hint_review_test.go (new)
- pkg/pdfcpu/readHints_test.go (new)

Parent metadata and owner-requested gitignore changes are separate from product
logic. No dependencies, vendor, lint policy, CI or consumer project changed.
Public ReadContext fields are additive but unkeyed literals/exact shape can be
incompatible; no exported function/interface/configuration signature changed.

## Exact validation

Logs: `.pi-herdsman/hint-fix-round2/`.

- `go test -count=1 -v ./pkg/api -run '^TestHintReviewUnreliableXRef$'`
  RED exit 1 (`issue3-red.log`), GREEN exit 0 (`issue3-green.log`).
- `go test -count=1 -v ./pkg/api -run '^TestHintReviewIncrementReopens$'`
  RED exit 1 (`issue4-red.log`), GREEN exit 0 (`issue4-green.log`).
- Prior continuation `issue1-{red,green}.log`, `issue2-{red,green}.log` retained;
  combined focus also covers these tests after local changes.
- `go test -count=1 ./pkg/api ./pkg/pdfcpu ./pkg/pdfcpu/model ./pkg/cli
  ./internal/contextutil -run
  'Hint|OrdinaryAES|Linearization|AESV2CryptFilter|EncryptedXRef|Metadata|Crypt|Step24a|ReturnedNotice'`
  exit 0 (`focused.log`).
- Initial full run exposed a nil ReadContext in an existing parser-limit unit
  fixture, because new optional xref tracking dereferenced ReadContext early.
  Fixed by supplying tracking only when reader state exists. Failure retained
  (`preflight.log`); `TestParseXRefSectionUsesConfiguredLimits` green
  (`nil-read-green.log`). No test weakened.
- `go test -count=1 ./... && go vet ./... && go build ./...`
  final exit 0 (`preflight-final.log`, PREFLIGHT_EXIT=0).
- `go test -race -count=1 ./pkg/api ./pkg/pdfcpu ./pkg/pdfcpu/model ./pkg/cli
  ./internal/contextutil -run
  'Hint|OrdinaryAES|Linearization|AESV2CryptFilter|EncryptedXRef|Metadata|Crypt|Step24a|ReturnedNotice'`
  exit 0 (`race.log`).
- gofmt all nine changed Go files; gopls check all nine paths exit 0
  (`gopls.log`); active lens diagnostics no type errors. Advisory typos/style
  hints and GORM N+1 pattern warnings are not product findings: there is no
  database/GORM query in these byte-decoding/structural loops.
- `git diff --check` exit 0. Samples before/after status empty; 446 regenerated
  tracked sample outputs restored and exactly three known new outputs removed.

## Remaining limits

Original production attachment unavailable; no direct production proof. Fixtures
are synthetic reader-layout reproductions, not proven Annex F-conforming hint
tables. Owner approved nonstandards relaxed compatibility; unsupported hint
shapes fail closed. PDF 2.0 complete normative verification not obtained and no
compliance claim made. Two other historical consumer failures unexplained.
No commit, version, publication, private processing, delivery or rollout.

Next gate: fresh dual independent review round 2, AND pass required.
