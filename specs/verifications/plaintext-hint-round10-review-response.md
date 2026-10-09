# Round 10 static review response

Evidence: `result:hint-static-a#1` and `result:hint-static-b#1`.
Both independent source reviews identified the same three defects. Neither
review executed tests or provided acceptance.

## Source corrections applied; runtime validation pending

1. **Must-fix: unvisited counterfeit xrefs.** Reader-local provenance now
   records xref origins from ordinary classic/stream traversal. Physical gap
   proof rejects classic xref offsets absent from that traversal; stream xref
   authority also requires a recorded origin. Accepted ordinary trailers remain
   in `Read.TrailerDicts` for the existing incoming-reference audit.
   Added an unvisited-table helper regression and a public counterfeit-xref
   fixture between the hint and ordinary objects. The latter contains a
   `/Target` reference that ordinary revision traversal does not visit.
2. **Must-fix: compressed indirect lengths.** Discovery now inventories
   candidates without granting repair/deletion. Physical envelope proof runs
   during finalization, after object streams and ordinary origins are loaded.
   The original file dictionary is reparsed; compressed integer lengths can
   then resolve through the already-materialized object stream. Added public
   plaintext/encrypted hint controls with an ordinary stream's length stored
   as a compressed integer. Ciphertext bytes are preserved by replacing only
   the first matched dictionary length entry.
3. **Should-fix: missing-length notice leak.** Physical inspection rejects a
   missing direct/indirect length quietly before the ordinary helper can emit
   its skip notice. Ordinary loading retains its existing reporting behavior.
   Added a missing-length envelope-denial helper case.

The existing byte-preservation unit test now targets the extracted content
finalizer with an explicitly verified candidate. Full physical proof has
separate draft tests; no production bypass flag or configuration was introduced.

## Follow-up static review

Evidence: `result:hint-static-a2#1` and `result:hint-static-b2#1`.
A2 identified a further source-level false denial; B2 reported no findings.
This was not a passing independent acceptance gate.

The proof previously required `Read.XRefStreams`, which the reader does not
populate. It now recognizes the cached `types.XRefStreamDict` produced by
successful ordinary traversal while still requiring the recorded physical
offset. Added a focused helper regression for parsed/visited, unvisited, and
name-only stream claims. Existing public xref-object-stream and compressed
indirect-length controls remain required. Fresh vet/build and primary LSP
checks pass after this correction; runtime results remain unverified.

## Optional items

Owner selected **Defer both**: positive early-xref-stream coverage and added
performance measurements are deferred. No performance or universal layout
support claim is made.

## Validation limits

Fresh vet/build checks pass; source checks do not establish runtime behavior.
All new regressions remain unexecuted in this session. Pending checks include
`TestHintPhysicalOrigins`, `TestHintProbePreservesOriginalAndEncryptedBytes`,
`TestPrimaryHintVariants`, `TestHintCompressedIndirectLength`,
`TestHintDiscoveryIdentityGuards`, and `TestUnencryptedPrimaryHint`.

Evidence reconciliation found no corresponding failure output in available
`.pi-herdsman` logs, cached error diagnostics, or recalled session entries.
Earlier notes described these names as observed superseded failures without
preserving supporting run output; that description is not a verified result.
Do not infer either passing or failing outcomes from it. The available round10
`baseline.log` records full Preflight exit 0 before the current corrections,
not verification of these changes. The source defects identified by reviewers
remain real review evidence; their corrections still require runtime proof.
Full Preflight, focused/race checks, and independent acceptance remain required.
No commit, merge, publication, or rollout is authorized.
