# Fork checkpoint

- Checkout base before these changes: `4d0e9ff2e547383e02b6b385ea9e6b48c4cc9f65`, reporting `v0.16.0-rc.1` and module `github.com/pdfcpu/pdfcpu`. An authoritative upstream-base comparison has **not** yet been established; do not call this commit an upstream base without verification.
- Fork modification at this checkpoint: relaxed Standard V=4/R=4 AESV2 reader compatibility when a referenced crypt filter omits local Length and top-level Length is explicitly 128 bits. Strict mode and negative cases retain rejection. Early Catalog Version parsing is syntax-only so effective version can be used before encryption validation.
- Synthetic regression is generated from the test-only PDF constructor and test-only passwords. The in-memory encrypted PDF is edited at fixed-width positions for the PDF-1.4 header, omitted crypt-filter Length, and Catalog Version override; no xref offsets are rewritten or production data included.
- The requested module rename to `github.com/umats/pdfcpu` is **pending**. Do not integrate into pdfrelay until reviewed independently. Apache-2.0 upstream license and copyrights remain untouched; no NOTICE file existed in this checkout.
- Before transfer: `go test ./... && go vet ./... && go build ./...` passed. Windows cross-builds and final fork-wide validation remain pending after module rename. No tag, PR, merge, or release was made.
