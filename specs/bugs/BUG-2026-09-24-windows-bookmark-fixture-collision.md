# Windows CI: bookmark test output collides with shared fixture

PR #1, [Test run 36028223392](https://github.com/umats/pdfcpu/actions/runs/36028223392), failed only on Windows Go 1.26: `TestAddSimpleBookmarks` could not replace `pkg/samples/bookmarks/bookmarkSimple.pdf` (`Access is denied`). The test writes to a tracked fixture that other package tests read. `go test ./...` runs packages concurrently; Windows can reject replacement while the file is open. The staged-write production path remains correct and must not be weakened.

Fix: write this test's output under `t.TempDir()` rather than the shared fixture, retaining PDF validation. Focused `TestAddSimpleBookmarks`, full `go test ./... && go vet ./... && go build ./...`, and specs validation pass locally. Windows CI verification remains pending; do not merge while CI is red.
