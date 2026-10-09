---
type: verification-evidence
context: BUG-2026-10-08-plaintext-linearization-hints
---

# Round 13 current-candidate verification

Evidence: `result:hint-fix13#1` and `result:hint-native32#1`.
Handoffs: `.pi-herdsman/hint-fix-round13/handoff.md` and
`.pi-herdsman/hint-fix-round13-native32/handoff.md`.

## Corrections and self-audit

Final recovered stream markers are lexically audited before prefix slicing in
shared ordinary/xref stream buffering. Ambiguous comment/name/string/container
origins mark provenance unreliable without changing historical parsing.
Unsigned xref generation/container/member-index narrowing is audited for
native-int round-trip loss before revision skipping, including older records.
Encoded-width checks remain. Lead inspected both corrections and found no
additional issue within scope; no public/configuration/dependency change.

A discovered test-only native32 compile defect was separately planned in
`specs/bugs/BUG-native32-crypto-test-constants.md`. Oversized unsigned test
constants now use int64 variables and runtime conversions after the unchanged
32-bit skip, following the adjacent test pattern. Native64 assertions remain.
This is a separate test-only fix for future authorized commit preparation.

## Behavioral sensitivity

Round13 `red-complete.sha256` and `red-complete.log` bind repeated behavioral
RED with only the two corrections removed. Native and actual private/public
386 commands fail on provenance/repair/deletion assertions, not compilation.
Reader RED hash matches round12. `green-final.sha256` and log prove GREEN.
Supplemental boundary/cancellation controls were added after correction and
are not claimed RED-sensitive. Earlier invalid test assumptions and compile
failures remain recorded, not substituted for canonical behavioral evidence.

Native32 compile bug has its own `source-before.sha256` and `red.log`, exit 1,
then package-based `focused386.log`, exit 0. No source-list workaround is needed
for current tests. Skip of native32 unsigned-value assertions is intentional;
native64 focused tests execute both assertions successfully.

## Fresh terminal receipts for current candidate

Under `.pi-herdsman/hint-fix-round13-native32/`, each contiguous, exit 0:

- `focused386.log`: actual package-based private/public 386 hint/permission tests.
- `focused-native.log`: native hint/permission controls.
- `preflight.log`:
  `go test -count=1 ./... && go vet ./... && go build ./...`.
- `race.log` and full supplement `context-race.log`.
- `gofmt.log`: all changed/untracked Go paths, empty formatting listing.
- `whitespace.log`: whitespace check.

Final candidate manifest SHA256 (`candidate-final.sha256`):
`05433ec198e8268294248c2c5ca7ecc55405c9fdc11c10d3661436929ae6fbc6`.
Includes crypto test repair; inherited round13 production hashes unchanged.
Stability receipts pass. Lead checked manifest against current files, exit 0
(`lead-manifest-check.log`), and whitespace passes.

Round13 `arm-build.log` proves linux/arm cross-build, not ARM execution.
No full386-suite claim. Active LSP confirms no errors in three tested paths;
`read.go` exceeds the 5000-line LSP cap and is not claimed checked. Native full
Preflight includes compiler/vet coverage of it; independent source review
remains required.

## Preservation and remaining gates

Sample outputs backed up/restored byte-for-byte, including all three inherited
untracked PDFs. No broad cleanup, private input or release actions.
Independent dual acceptance/security and bug-scoped completeness remain pending.
Owner accepted synthetic UAT, not production proof. Fixtures are not fully
Annex F conforming; two historical consumer failures remain unexplained.
Optional early-xref-stream and performance measurements are deferred; exported
ReadContext unkeyed-literal caveat remains. Phase 5 is incomplete.
