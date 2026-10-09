---
bug_id: BUG-2026-10-08-plaintext-linearization-hints
status: verified_unreleased
severity: high
scope: reader-encryption
title: Plaintext linearization hint stream reaches AES decryption
---

# Plaintext linearization hint streams

## Authorization and evidence limits

### Owner-approved compatibility policy (supersedes standards gate)

Owner explicitly waived standards compliance for this internal production fix
and selected the existing relaxed mode for narrowly scoped plaintext-hint repair.
Strict mode, password authentication and ordinary encrypted streams remain
unchanged. No new configuration setting or consumer change is requested.
The verified standard remains documented; the repair must be labeled a
compatibility exception, never a normative encryption exemption.

Structural hint identification must precede any plaintext decoding probe.
An unchecked H offset, AES alignment, successful inflate or dictionary name
alone cannot classify ordinary content as an exempt hint. Preserve normal
ciphertext handling first where feasible and bound any hint-only repair probe
by existing limits. Primary/overflow identification, stale linearization,
references from ordinary document objects and writer behavior require review.
If sufficient safe identification needs a materially broader parser change,
return the design trade-off to the owner before implementation.

### Broader reader support approved after synthetic RED

Owner selected broader reader support rather than classic-xref-only repair.
Implementation is authorized with complete incoming-reference coverage,
including compressed object origins. Do not waive identity/authentication,
resource/cancellation limits or writer content preservation. No general fallback.

Synthetic RCA completed (`result:hint-reproduction#1`): public PDFInfo,
ReadContext and Decrypt fail on primary hint object 8 with the expected AES
alignment error. The hint is 142 encoded bytes, not 153; the fixture uses real
xref/discovery paths but is not fully Annex F conforming. Controls and focused
prior encryption tests pass. Logs/design: `.pi-herdsman/hint-repro/findings.md`.

Approved implementation plan:

1. RED preserved in the existing new API regression file; identify candidates
   from the physical-first linearization dictionary before any hint decryption.
   Validate direct fields and H ranges, live unique object-start identity,
   object envelope and hint dictionary shape. Stale/repaired/ambiguous claims
   are ineligible. Do not change unrelated strict-mode validation.
2. Defer candidate stream finalization until reference coverage is complete;
   use existing object-stream parsing/materialization under current limits,
   inspect all structural references plus trailer roots, with cancellation and
   recursion checks. Referenced content/image/metadata targets cannot be repaired.
3. For verified unreferenced hints only, preserve Raw while trying normal
   effective AES handling/decode first. Relaxed mode may decode original
   plaintext after normal processing fails. Aligned plaintext must be covered;
   cancellation/limits never trigger fallback. Preserve encrypted-hint behavior.
4. Full writer deletion must use verified hint identity, never unchecked H
   equality that could delete ordinary objects. Keep strict/auth/metadata/xref,
   permissions and incremental-write contracts unchanged.
5. Add vertical regression slices for overflow/aligned/encrypted hints,
   ordinary malformed AES, forged referenced targets including compressed
   origins, malformed H/stale claims, corrupt compression, limits and output
   validation/content/no-encryption. Explicit crypt-filter behavior unchanged.
6. Run focused compatibility tests, full fresh Preflight, applicable race
   checks, gofmt and whitespace checks. Independent read-only review is required
   before acceptance. Preserve exact RED/GREEN and known fixture limitations.

Executor owns product code and regression tests only; parent owns planning,
state, acceptance and independent review. Source scope: existing reader/writer
and model context tracking as needed, related synthetic/internal tests. No new
public API, dependencies, vendor/CI/lint-policy edits or arbitrary refactor.
No commits/stash: owner explicitly prohibits commits; preserve RED/GREEN logs
rather than following the generic skill's commit checkpoints.

Original authorization (the standards prerequisite is now superseded):
Owner authorized investigation and the smallest standards-grounded synthetic fix
on a feature branch in this existing checkout. No worktree, commit, push, publish,
tag, merge, deploy, SMTP, private sample access, or pdfrelay changes authorized.
Owner approved carrying survey metadata/runtime evidence onto the branch without
commit/stash and requested ignoring `/.pi-herdsman/`.

Reported consuming dependency: `github.com/umats/pdfcpu v0.15.1`.
One of three reported pdfrelay v0.7.1 failures correlates with the supplied log;
build identity was not logged. Two other historical failures remain unexplained.
No production attachment, log, filename, author, content, or password supplied to
this checkout. Do not request or log actual password values.

## Reproduce

Owner reports authenticated AESV2 linearized PDF reading fails with `ciphertext
not a multiple of block size`. The primary `/H` object contains a plaintext
FlateDecode stream (153 encoded bytes, no `/Crypt`) that inflates successfully.
Metadata inspection, decryption and property removal fail before SMTP delivery.
An offset-preserving disposable empty-hint experiment isolated the reported
stream, but stream removal is NOT an authorized repair.

This is transferred evidence, not yet a local synthetic reproduction.

Baseline `870fbbf069d1c9501bb27b019d01bd80d1cb1d23` passed:

```sh
go test -count=1 ./... && go vet ./... && go build ./...
```

Exit 0; `.pi-herdsman/hint-baseline-preflight.log`. Known regenerated sample
outputs were restored from an initially clean samples subtree.
Branch: `fix/plaintext-linearization-hints`.

## Isolate

Read-only reconnaissance (`result:hint-reader-map#1`) finds that linearization
tracking records direct `/H` offsets but does not validate signs, bounds,
uniqueness, lengths or matching stream identity. Stream loading applies ordinary
document encryption before optional decoding. Post-load hint dictionary checks
cannot establish trustworthy identity before decryption. Full writes remove
tracked linearization objects/hints; incremental writes need separate scrutiny.

Existing AESV2 fixtures and public read/decrypt builders can be reused, but no
complete linearized fixture builder was found. The prior missing crypt-filter
Length compatibility fix is related encryption work, not proven recurrence.

## Hypothesize

1. A normative implicit hint exemption is missing: falsify/confirm against the
   complete authoritative encryption and Annex F rules before code changes.
2. Explicit crypt-filter selection is wrong: inspect `/Crypt` and `/Identity`
   with synthetic input; do not conflate with an implicit exemption.
3. Plaintext input is malformed: if normative rules require encryption, stop
   and report the contradiction. Any new compatibility policy needs owner
   approval and independent review, not an offset-only bypass.

Security impact: HIGH for a proposed unchecked `/H` exemption. A hostile
linearization dictionary could designate ordinary content/image/metadata
streams as hints and bypass document encryption. This is prospective repair
risk, not a confirmed exploit in current code. Authentication, permissions,
metadata behavior, xref-stream rules and resource limits must remain intact.

## Verify: implicit exemption rejected for PDF 1.7

Complete-source follow-up (`result:normative-pdf#1`) fetched both authoritative
PDF 1.7 texts and read their complete encryption and hint-stream clauses.
This supersedes the indexed-only evidence below for PDF 1.7/ISO 32000-1.
PDF 2.0 remains unverified: the supplied final-PDF URL returned HTML.

Verified citations:

- ISO 32000-1:2008, Table 20 StmF, printed page 56 (PDF page 64):
  "All streams in the document, except for cross-reference streams
  (see 7.5.8, 'Cross-Reference Streams') or streams that have a Crypt entry
  in their Filter array (see Table 6), shall be decrypted by the security
  handler, using this crypt filter." Default StmF is Identity.
- Section 7.6.2, printed pages 57-58 (PDF pages 65-66):
  "Stream data shall be encrypted after applying all stream encoding filters
  and shall be decrypted before applying any stream decoding filters."
- Section 7.6.5, printed page 67 (PDF page 75), permits an explicit Crypt
  filter to override the default; Identity passes data unchanged. Table 25
  identifies AESV2 as PDF 1.6 AES-CBC with a 16-byte block and prefixed IV.
- Complete F.3.6, printed pages 681-683 (PDF pages 689-691), covers primary
  and overflow hint streams with no encryption exception. Logical exclusion
  permits regeneration/removal, not plaintext treatment on the read path.
- Table F.1 H, printed page 680 (PDF page 688):
  "This is the beginning of the stream object, not the beginning of the
  stream data." Offset/length pairs include object overhead.
- Adobe PDF Reference, sixth edition, PDF 1.7 (November 2006), Table 3.18
  StmF page 117 and complete F.2.5 pages 1032-1034 independently agree.

Complete-source URLs:

- https://www.reportlab.com/ftp/PDF32000_2008.pdf
- https://openn.library.upenn.edu/Data/0014/ArchimedesPalimpsest/Documents/External/pdf_reference_1-7.pdf

Fetched PDFs, extracted text, page evidence, hashes and extraction commands:
`.pi-herdsman/hint-standards/`; receipts in `receipts.json`, details in
`findings.md`. Installed Ghostscript 10.08.0 txtwrite was used, no installation.

Decision: normal effective crypt-filter selection governs both hint streams.
AESV2-selected plaintext hints are inconsistent input. No Crypt declaration
alone does not prove this production file malformed: default StmF can be
Identity. Authentication does not establish filter selection. No sample or
local synthetic RCA establishes the actual production selection.

### Historical indexed-only attempt (superseded for PDF 1.7)

Normative research completed (`result:hint-standards#1`). The researcher reported
no implicit hint-stream exemption. However, complete standard PDF retrieval
failed (size rejection or 404); the result relied on indexed passages and errata.
The parent therefore does NOT accept its claimed gate pass as full normative
verification. Exact complete authoritative scope for both hints is still missing.
No offset-based exemption is approved.

Researcher-reported passages, not independently verified complete text:

- ISO 32000-1:2008 section 7.6.1, General:
  "Encryption applies to all strings and streams in the document's PDF file,
  with the following exceptions:" followed by trailer ID values, strings in
  the encryption dictionary, and strings contained within encrypted streams.
  URL: https://www.reportlab.com/ftp/PDF32000_2008.pdf
- ISO 32000-2:2020 section 7.6.2, printed page 71, reportedly retains this rule
  and adds signature Contents hexadecimal strings. Section 7.6.3 reportedly
  states: "Stream data shall be encrypted after applying all stream encoding
  filters and shall be decrypted before applying any stream decoding filters."
  URL: https://previewnorm.com/iso/ISO%2032000-2-2020%20PDF.pdf
- Explicit Crypt/Identity selection differs from an implicit exception;
  reported locator ISO 32000-2 section 7.6.6. Table F.1 locates the beginning
  of the hint-stream object, not its stream data; F.3.6 covers both hints.

If these complete rules confirm no implicit exemption and the effective stream
filter is AESV2, plaintext hints without an effective Identity selection are
inconsistent input, not a standards-authorized bypass. Actual sample selection
is unavailable and synthetic crypt-filter RCA has not been performed. Treat
input conformance as unresolved, not a verified malformed-production conclusion.

Gate requirements remain:
Require public URL, edition, section/page/table and exact relevant passage for
primary AND overflow hints, supported PDF versions and AESV2. Establish whether
`/H` locates object starts versus stream data. Successful inflate, block alignment
and error messages do NOT prove an exemption.

Research leads (not normative proof):

- https://pdf-issues.pdfa.org/32000-2-2020/clauseAnnexF.html
  Table F.1, F.3.1, F.3.6.
- https://pdf-issues.pdfa.org/32000-2-2020/clause07.html, section 7.6.2.
- https://www.loc.gov/preservation/digital/formats/fdd/fdd000030.shtml

Historical standards prerequisite superseded by owner compatibility policy.
Local synthetic RCA and bounded repair design are still required.

## Conditional synthetic TDD acceptance

- Primary plaintext compressed hint: valid encrypted linearized synthetic PDF,
  nonaligned hint length (153 if valid), ordinary encrypted streams intact.
  Capture RED public metadata failure, then GREEN read/decrypt and validated
  output, page count/content preservation, no encryption and existing permissions.
- Overflow: valid four-element `/H`, prove both primary/overflow public paths.
- Security/parser: malformed ordinary AES remains terminal; negative,
  out-of-range, duplicate, indirect, stale and wrong-type `/H` cannot exempt
  ordinary ciphertext; corrupt consumed compressed hint fails without panic.
  Preserve documented strict/relaxed behavior rather than unrelated tightening.
- Compatibility: unencrypted linearized and encrypted non-linearized files;
  metadata, Identity crypt-filter, wrong password, xref exclusion, affected
  RC4/AES variants and aligned plaintext hints.
- Preserve encoded bytes, filters, decode/size/cancellation limits. No generic
  plaintext fallback, AES error swallowing, deletion repair or dependencies.
- Validate full-write output without stale linearization; do not promise output
  linearization beyond existing API contract.
- Reuse synthetic builders; external fixture generation must record version and
  commands and cannot become an uninstalled routine-test dependency.
- Run gofmt, narrow new regressions, reader/crypto/API/metadata/linearization
  compatibility tests, full Preflight, configured lint if present, applicable
  race checks and diff whitespace checks. Independent read-only review focuses
  on hostile hint identification and accidental ordinary-stream exemptions.

## Required handoff

Return normative scope, verified RCA/fix, paths/diff, exact RED/GREEN/validation,
independent review/resolution, skipped checks and the two unexplained failures.
No fixed commit until separately authorized; no invented or published version.
A later pdfrelay dependency/vendor update is a separate authorized task and must
preserve source retention, Author selection, metadata, permissions and atomic
email delivery.

## Resolution

Approved synthetic/internal compatibility fix verified, not released.
Candidate-bound current execution, dual acceptance and completeness:

- `specs/verifications/plaintext-hint-round13-verification.md`
- `specs/verifications/plaintext-hint-round13-acceptance.md`
- `specs/verifications/plaintext-hint-completeness.md`

Fresh native full Preflight and applicable race checks pass. Actual private/public
package-based386 regressions pass; ARM cross-build only. Independent A/B each
score 100%, with no findings and scoped security PASS. Owner accepted synthetic
behavioral evidence. Unknown origins deny repair/deletion; ordinary cipher,
authentication, permissions, limits, staged writes and content are retained.

No production reproduction or Annex F compliance claim. Two historical consumer
failures remain unexplained; optional layout/performance coverage and exported
ReadContext shape caveat remain. Historical incomplete RED hashes are disclosed,
not fabricated. No commit, publication, consumer update or release approval.
Owner confirmed Phase 5 completion and Phase 6 preparation only.
Commit/PR/integration/publication authority remains unapproved; see
`specs/verifications/plaintext-hint-release-preparation.md`.

### Historical resolution checkpoint (superseded)

Owner resumed as relaxed-mode internal compatibility work. Synthetic RCA and
safe identity design were next; production changes were pending design review.

Previous checkpoint (superseded): stopped before implementation because the complete PDF 1.7 normative rules
do not support the proposed implicit exemption. PDF 2.0 retrieval remains
incomplete; no claim about its complete normative scope is made. No RED/GREEN fixture, product
fix, independent patch review, commit identifier, or patch version exists.
Fresh baseline Preflight passed; YAML validation, metadata LSP diagnostics and
`git diff --check` passed. No additional lint configuration was found at the
repository root or one directory below it; repository lint command is go vet.

At that historical checkpoint, resume required owner direction for a separate
crypt-filter investigation or a new compatibility policy. The owner-approved
policy and current verification above supersede that pause. No offset-only
exemption or normative encryption exemption is authorized.
