# AESV2 crypt-filter Length compatibility: impact and plan

## Checkout and scope
- Fork HEAD `4d0e9ff2e547383e02b6b385ea9e6b48c4cc9f65` (`fix Windows configuration test portability`), on `docs/seed-conventions-team-pr`; `model.VersionStr` says `v0.16.0-rc.1`, **not v0.15.0**. No local tags or upstream remote establish an upstream base yet.
- Stage 1: add a synthetic test-only PDF generated from an existing synthetic source and AES-128 encryption with test-only passwords; omit the generated crypt filter's `/Length` without changing byte lengths or xref offsets. Prove `ReadContext`/`Decrypt` reject it specifically for missing crypt-filter Length before the fix.
- Stage 2: accept the omission only in relaxed mode for Standard V=4/R=4, explicitly declared top-level 128-bit Length, and a referenced AESV2 stream/string crypt filter. Keep strict mode and all other invalid cases rejecting it. Confirm authentication and output readability/validation. No parser-error success or general repair.
- Separately reviewable module rename: change `go.mod`, maintained imports, two nested module requirements, examples and operational tooling to `github.com/umats/pdfcpu`; preserve upstream attribution and Apache-2.0 license. Do not touch `vendor/` or pdfrelay, publish, tag, push, or release.

## Dependents and risks
- `supportedEncryption` validates `Filter`, `V`, top-level `Length`, then crypt filters, then `R`. The exception needs `R` before crypt-filter validation without broadening other cases. `validateCryptFilters` routes `StmF`, `StrF`, and `EFF` through `locateCFEntry` to `validateCryptFilter`; only referenced filters are checked. Effective PDF version is `XRefTable.Version()` (Catalog `/Version` overrides header) and `ctx.PDF20()` reflects that.
- Shared parser and writer encryption paths affect public API read, decrypt, validate, and CLI operations. Risk: **High** (security-sensitive shared code). Tests must reject missing/wrong top-level key length, unknown/missing CFM, other V/R, public-key handler, malformed/conflicting/wrong explicit crypt-filter lengths, and wrong passwords in strict/relaxed modes.
- Existing tests: `pkg/pdfcpu/crypto_test.go` checks crypt-filter classification and AES lengths; `pkg/api/encrypted_xref_reconstruction_test.go` builds an encrypted synthetic PDF and reads/decrypts it. Gap: omission acceptance and rejection matrix across public API and effective Catalog version.
- Baseline `go test ./... && go vet ./... && go build ./...` passed. Tests generated changes in tracked sample PDFs; restored only these baseline-clean generated files and three new generated sample PDFs.

## Gate and next steps
- Discover complete; await confirmation to proceed to investigation/plan and separate Stage 1 implementation. Feature branch/worktree must not absorb the pre-existing untracked `CLAUDE.md`, `CONVENTIONS.md`, or `specs/` accidentally. Verify upstream base against an authoritative upstream commit before documenting it as fact.
