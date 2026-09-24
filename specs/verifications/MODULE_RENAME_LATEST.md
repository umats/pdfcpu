# Fork module rename verification

## Scope
- Base: upstream [`4d0e9ff2e547383e02b6b385ea9e6b48c4cc9f65`](https://github.com/pdfcpu/pdfcpu/commit/4d0e9ff2e547383e02b6b385ea9e6b48c4cc9f65), independently matched to upstream `HEAD` and `master` by `git ls-remote` on 2026-09-24.
- Current branch: `fix/aesv2-missing-cf-length`; prior AESV2 fix and docs are commits `21c66543` and `3048b4ee`.
- Main module, 506 tracked Go files' first-party import string literals, two nested module directives, release linker path, coverage scripts, local Docker build, and fork documentation now use `github.com/umats/pdfcpu`. Historical upstream attribution URLs, LICENSE.txt, and copyright remain unchanged; no NOTICE exists in the checkout. No pdfrelay or vendor edits.

## Commands and results
- Baseline on Go 1.27.1/darwin/arm64: `go test ./... && go vet ./... && go build ./...` — PASS before rename.
- `go mod tidy` — PASS; only unused old checksums removed from go.sum.
- `go list ./...` — PASS, 32 fork packages; `go list -m` in root and both nested modules reported fork paths.
- `go test ./... && go vet ./... && go build ./...` — PASS after rename.
- `GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o /tmp/pdfcpu-fork-windows-amd64.exe ./cmd/pdfcpu` — PASS.
- `GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -o /tmp/pdfcpu-fork-windows-arm64.exe ./cmd/pdfcpu` — PASS.
- Windows amd64 and arm64 `CGO_ENABLED=0 go build ./pkg/... ./internal/...` — PASS.
- `goreleaser check` — PASS. `git diff --check` — PASS. Changed Go files are gofmt-clean.
- Exact import audit: all 506 modified Go files matched the former content with only `"github.com/pdfcpu/pdfcpu/` replaced by `"github.com/umats/pdfcpu/`, except 13 Cobra CLI files where gofmt reordered imports. No old first-party Go imports or maintained tooling module paths remain. The full changed-path list is `specs/verifications/CHANGED_FILES_LATEST.txt`; use `git diff 3048b4ee..HEAD` after the documentation checkpoint to inspect the complete diff.
- Docker image build — PASS on Colima/Linux arm64: `docker build -t pdfcpu-fork-uat .` built the fork from local source. `docker run --rm pdfcpu-fork-uat --help`, `version`, and invalid-flag rejection passed. The legacy Docker builder lacks `--progress`; retry without it passed. Audit found the initial build context was 755.4 MB; `.dockerignore` now excludes `.git`, PDFs, and sample/testdata modules. Rebuild and smoke passed with a 27.74 MB context; the final image contains only `/root/pdfcpu`. No image was pushed.
- Full test run produced sample output PDFs; restored only files known clean at baseline and removed newly generated sample artifacts. No sample fixtures are part of the change.

## Remaining review/risk
- Integration into pdfrelay is separate and untouched. The fork is not merged, tagged, or released. The malformed crypt-filter omission remains a deliberately narrow relaxed-mode exception; accepting nonconforming PDFs can expose parser paths, so review the AESV2 and Catalog bootstrap diff before shipping.
- No upstream issue or PR was opened without approval. Existing README upstream badges and historical links remain as attribution; fork package docs and examples point to the new path.
