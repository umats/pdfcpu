# Round 10 corrections — incomplete checkpoint

Round 9 dual review failed: A87.5%, B90%. Three findings remain unresolved:
comment-contained nonstream terminators, complete counterfeit epilogues, and
manufactured hint header origins. Final independent review 10 has NOT run;
the review-round cap has not been exhausted.

## Recovery evidence

First round10 executor retired without edits after context pressure. Fresh
executor then passed baseline full Preflight (`baseline.log`, exit0) and wrote
only partial regression code:

- `pkg/api/plaintext_linearization_hint_test.go`: review10 option and auxiliary
  comment-contained terminator fixture.
- `pkg/api/plaintext_hint_round10_test.go`: draft public denial regression.
- `pkg/pdfcpu/read_hint_round10_test.go`: draft internal value/provenance test.

No regression execution result or RED log exists for the draft round10 tests.

Local recovery checked prior round9 source hashes: product source unchanged;
only the existing API fixture changed, plus two draft tests. No product fix or
final verification exists for round10. Preserve these partial changes; do not
claim accepted or fixed.

Baseline root status proved samples initially clean. Parent restored exactly
446 known test-generated tracked sample outputs and removed the three known
new sample PDFs. No lock removed; product/owner/test changes preserved; final
sample status clean. Evidence under `.pi-herdsman/hint-fix-round10/`:
`baseline-status.log`, `baseline.log`, `recovery-cleanup.log`.

## Current source corrections (unverified)

The shared `object` reader now checks `lexicalCutState` at a nonstream
`endobj` cut before trusting it for hint authorization. A marker inside a
comment, string, token, or open container marks provenance unreliable.
Ordinary object parsing and returned values remain on the existing path.
Cancellation from the bounded audit is propagated.

Discovery inventories candidates without granting repair or deletion.
Finalization calls `proveHintOrigins` in `readHintOrigins.go` after ordinary
object-stream materialization. It orders physical xref claims, anchors the
first claim to the physical first object, and proves each predecessor envelope
using original encoded stream lengths. Accepted xref origins must also have
been visited by ordinary reader traversal; an unvisited table cannot authorize
an epilogue or escape the retained-trailer reference audit.
Overlapping claims, clipped comments, unknown intervening syntax, and hidden
values after an epilogue deny compatibility. Gaps are bounded by the existing
object limit; stream postludes are bounded to 1025 bytes. Complete epilogues
must reference a proven physical xref origin (or the early linearized zero
sentinel after a proven classic xref), and EOF comment lines cannot be clipped
into manufactured object origins. Cancellation and terminal limits propagate.

Added `readHintOrigins_test.go` covering gap framing, hidden structural
continuation, clipped EOF/comments, stream-contained headers, valid chains,
and cancellation. Public regressions also cover unvisited counterfeit xrefs
and compressed indirect stream lengths. Updated the discovery unit fixture
to include a real classic xref/trailer/epilogue. Existing round10 public/internal
drafts are preserved. Review findings and source responses are recorded in
`plaintext-hint-round10-review-response.md`.
No regression tests were executed in this session.

Formatting and source whitespace checks pass. Active LSP checked five paths
and reports no primary errors/warnings; three auxiliary GORM-pattern warnings
concern non-database PDF loops. Fresh `go vet ./... && go build ./...` completed
with exit 0 after all source corrections. This is static verification only.
The changes are not accepted: fresh tests, Preflight, race, and independent
review remain outstanding. Unsupported physical layouts intentionally deny
hint compatibility and deletion rather than trusting declared offsets.

## Approved design remains pending

Conservative physical-origin proof must be anchored to physical header/first
object or proven predecessor envelope chains, not unchecked xref or H offsets.
Use original encoded stream lengths and actual terminators; bind epilogues to
actual revision boundaries. Only bounded whitespace/comment gaps permitted;
unknown layouts deny compatibility/deletion, ordinary reads unchanged.
No unbounded scan, general parser rewrite, or limit loosening approved.

Remaining work: complete round10 regressions and corrections, then run fresh
validation and independent review before acceptance.
No commit/push/merge/tag/publication/consumer rollout. Production sample absent,
Annex F synthetic proof absent, two historical failures unexplained, public
ReadContext shape caveat unchanged.
