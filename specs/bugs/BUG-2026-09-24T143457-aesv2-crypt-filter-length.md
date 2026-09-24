# BUG-2026-09-24T143457: AESV2 crypt filter without local Length

## Problem
A synthetic Standard V=4/R=4 encrypted PDF with explicit top-level 128-bit Length and a referenced AESV2 crypt filter lacking local Length is rejected by public read/decrypt. Strict rejection must remain; relaxed mode needs a narrowly gated interoperability exception.

## Root Cause Analysis
The encryption validator rejects missing local crypt-filter Length before it considers the Standard handler's revision or the explicitly declared top-level key length. All referenced stream/string/embedded-file filters use that validator. The reader checks encryption before it identifies Catalog Version, so effective-version handling needs a separate review. Security impact: **MEDIUM** if a broad exception were to admit unrelated malformed dictionaries; the proposed guard keeps authentication mandatory. Risk: **High** for shared encryption parsing.

## Stage 1: observed RED
- `go test -count=1 ./pkg/api -run '^TestMissingAESV2CryptFilterLength$' -v`: fails on `encrypt dict entry "StmF": malformed encryption: crypt filter missing entry "Length"`.
- `go test -count=1 ./pkg/pdfcpu -run '^TestSupportedEncryptionMissingAESV2Length$' -v`: fails on the same missing-Length error.
- Synthetic source is constructed in `pkg/api/encrypted_xref_reconstruction_test.go`. Public `Encrypt` uses test-only passwords and AES-128; the in-memory output's PDF header and local Length are changed in fixed-width positions without changing xref offsets. No production fixture or password is used.

## Stage 2 plan
1. **GREEN:** In relaxed mode only, permit missing local Length on specifically referenced AESV2 under Standard V=4/R=4 with explicit top-level Length=128. Keep strict mode rejection. Verify both RED tests become green.
2. **RED/GREEN:** Add a public API rejection matrix for wrong password, missing/wrong top-level Length, missing/unknown CFM, other V/R, public-key handler, and malformed/conflicting/wrong-size explicit local Length; check both modes and decrypted output content/validation.
3. **RED/GREEN:** Assert Catalog `/Version` override behavior, rather than relying only on `%PDF-1.4` header. Ensure parser errors remain errors and no arbitrary bytes are rewritten.
4. Run focused tests, full Preflight, and cross-platform builds; keep module rename separately reviewable.

## Acceptance Criteria
- [x] Relaxed public read/decrypt succeeds; strict continues to reject.
- [x] Wrong password and malformed/inconsistent dictionaries remain errors.
- [x] Decrypted output can be read and validated with synthetic credentials.
- [x] Full tests, vet, builds, Windows cross-builds pass after the module rename.

## Resolution
Stage 2 focused tests and full Go preflight pass before and after the module rename. Windows amd64/arm64 builds and Docker image smoke also pass. The user accepted UAT on 2026-09-24. A reviewer found a pre-authentication Catalog stream recursion and xref-repair risk; bootstrap now parses Catalog syntax without resolving stream filters or caching undecrypted objects. Unit tests cover cyclic stream-filter syntax, indirect Catalog Version, and corrupt Catalog offset deferral. No production fixtures are used. Independent review of the module rename remains a separate integration gate.
