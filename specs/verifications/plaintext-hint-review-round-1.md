# Plaintext hint review — round 1

Independent read-only reports: `result:hint-review-a#1` (70/100),
`result:hint-review-b#1` (72/100). Both changes-requested; AND gate FAIL.
No reviewer executed tests; findings are source-path verified and must be
runtime-reproduced before claiming each correction proved.

## Consolidated disposition: all fix before acceptance

1. MUST: invalidate deferred candidate state when xref reconstruction replaces
   entries already decrypted by parseAndLoad. Encrypted primary/overflow plus
   unrelated repairable corrupt offset must reopen with preserved content;
   no compatibility/deletion authorization after rebuilding.
2. MUST: strict policy is not sufficient provenance. Duplicate candidate
   H/Length/S fields, parser fallback, and duplicate overwritten references in
   ordinary/trailer/compressed origins must disqualify hint authorization.
   Preserve ordinary parsing behavior; no unrelated global rejection tightening.
3. MUST (A; SHOULD B): repaired/unreliable xref marker must cover Size/Index
   corrections and skipped in-use classic entries, not only bypass rebuild.
   Missing origins cannot authorize plaintext or deletion. Test plaintext
   rejection and ordinary encrypted successful controls.
4. MUST (B): incremental mutation after repairing hints must not produce an
   unreadable append. Track actual plaintext repairs and serialize normally
   encrypted replacement hint objects in incremental output, so stale original
   linearization claims do not require plaintext bypass on the next read.
   Test public annotation increment primary/overflow, reopening/validation and
   content preservation. Do not relax stale-L eligibility or silently reject
   operations without owner decision if a safe replacement is not possible.

Executor: existing implementer continuation owns all source/tests/fix validation.
Parent owns metadata, acceptance and re-review. No commits/publication, private
sample use, standards-compliance claim or consumer rollout authorized.

Validation: preserve RED/GREEN for each finding; focused suite, fresh complete
Preflight, focused race, gopls and diff checks; fresh dual-blind review round 2.
