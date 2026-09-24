# Fork change security review

Scope: `fix/aesv2-missing-cf-length` versus upstream base `4d0e9ff2`.

- The committed AESV2 exception was reviewed independently before the module rename. Reviewer findings on pre-auth Catalog stream recursion, corrupt-offset repair, and indirect Catalog Version were addressed with syntax-only bootstrap parsing and focused tests.
- The locally committed rename changes first-party import literals in 506 Go files (plus gofmt import ordering in 13 CLI files), module names, build/linker paths, coverage scripts, Docker build source, and docs. A mechanical content comparison found no unexpected Go edits. No new dependency or authentication/cryptography logic is introduced by the rename.
- `go list -deps ./...` and `go list -m all` contain no upstream pdfcpu identity; the public API authentication regression (wrong password) remains covered in `pkg/api/crypt_filter_length_test.go` and passes under the renamed module.
- Self-audit found that Docker's broad `COPY . .` sent `.git` and PDF fixtures to the local build daemon (755.4 MB context). A native `.dockerignore` now excludes `.git`, all PDFs, `pkg/samples`, and `pkg/testdata`. Rebuilding passed with a 27.74 MB context; the final image contains the CLI binary only.
- No new HIGH-confidence exploit path was found in the rename diff. Residual risk remains from processing intentionally nonconforming encrypted PDFs and pre-auth Catalog syntax. The rename received a self-audit, **not an independent reviewer audit**; obtain one before integration. A code review cannot establish safety for all malformed PDFs.
- No production PDF, password, customer identity, or metadata was used for the regression. No pdfrelay/vendor changes or release/publish actions occurred.
