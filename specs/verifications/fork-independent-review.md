# Fork independent review

- Round 1: two standalone, read-only Pi reviewers independently failed (88/100 each). Both identified the indirect Catalog Version/adjacent-stream boundary bug and misleading upstream-only operational links.
- Round 2: after the regression and fixes, the same independent setup passed (94/100 and 95/100; no must-fix findings). Reviewers did not edit files or independently run the Go suite.
- RED: `go test -count=1 ./pkg/pdfcpu -run '^TestBootstrapCatalogVersionIndirectName$'` failed with effective version 1.4, expected 2.0.
- GREEN: targeted bootstrap and public AESV2 tests passed. `go test ./... && go vet ./... && go build ./...`, `validate-specs-yaml.sh "$PWD/specs"`, and `git diff --check` passed.
- Remaining: public-API malformed-encryption matrix and end-to-end indirect-version case could strengthen coverage; upstream CI/coverage badges remain unlabeled. Fork-specific private vulnerability reporting is not yet available and is explicitly disclosed in `SECURITY.md`.
- No commit, PR, merge, tag, or release was performed.
