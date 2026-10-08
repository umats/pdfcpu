# Parallel test font restoration failure

Status: fixed in working tree; fresh parallel gates green, pending owner review and separate fix commit.

## Evidence

On baseline `57dc6e3f`, `go test ./...` failed in TestTextRenderingUsesConfigurationFontRepository restoring missing installed UnifontUpperMedium.gob. A focused rerun passed.

After resolving upstream/master `2f7a89a0` on `merge/upstream-master`, fresh `go test -count=1 ./...` failed in TestPrepBytesPropagatesFontErrors at pkg/pdfcpu/model/text_test.go:73 restoring missing installed UnifontMedium.gob. Log: `.pi-herdsman/merge-fresh-full.log`.

`GOFLAGS=-p=1 go test -count=1 ./... && go vet ./... && go build ./...` passed. Cached standard Preflight also passed, but does not invalidate the fresh parallel failure.

## Approved bounded fix plan

Trace font installation and restoration; add deterministic assertions that each font-installing integration test binary uses its own temporary configuration/font store, including legacy default-configuration callers. Run the assertions red against current setup, then isolate API and CLI TestMain configuration under their existing temporary output roots. Keep all edits test-only and unstaged, preserve the merge index, run focused tests and fresh parallel Preflight. No production/CI serialization change. Red/green commits and workflow-state transitions are deferred to the owner because this assignment prohibits commits and owns only this defect.

## Hypothesis and next steps

Initial hypothesis: cross-package tests may interfere through the shared installed font/configuration directory. The investigation and test isolation below establish the causal path and remove these shared test writers. Keep the fix in a separate Conventional Commit; merge integration remains staged and the defect fix remains unstaged.

## Root cause and fix (2026-10-08)

- API and CLI integration `TestMain` both called `api.LoadConfiguration` with no explicit root, selected the same OS user store, assigned `font.UserFontDir`, and installed all fixtures from `pkg/testdata/fonts`. The fixture inputs include `unifont-13.0.03.ttf` and `unifont_upper-13.0.03.ttf`; the new passing regressions verify their installed `UnifontMedium.gob` and `UnifontUpperMedium.gob` files in each private store.
- `pkg/api/font.go`, `commitStagedFontsWithOperations`, renames an existing installed target into `.pdfcpu-font-backup-*`, syncs directories, then renames the staged replacement into place. During that interval the target is absent. Both integration binaries performed these replacements against the store read by other package tests.
- `pkg/font/metrics.go`, `loadUserFontMetrics`, first lists `.gob` files with `os.ReadDir`, then opens each file separately. `pkg/pdfcpu/model/text_test.go` cleanup restores the original directory and calls `ReloadUserFonts`, so a listed file can disappear before it is opened. This explains the exact missing-installed-font restoration failures preserved above. The failure is test resource ownership, not a reason to weaken font errors or serialize CI.
- Minimal fix: in `pkg/api/test/api_test.go` and `pkg/cli/test/cli_test.go`, create the already-existing suite temporary output root before configuration initialization; initialize legacy defaults via `model.EnsureDefaultConfigAt(outDir, false)` and load suite configuration with `Root: outDir`. Legacy `model.NewDefaultConfiguration()` callers, explicit suite configurations, and global font installation now all use this binary's private store. Existing end-of-suite cleanup also removes this configuration; setup failure paths remove it too.
- No production, vendor, lint configuration, CI or resource-limit changes. No modifications to the staged merge index.

## TDD and validation

- Added `TestMainFontStoreIsPrivate` to each integration package. Assertions cover both suite and legacy default store selection, global installer destination and presence of both implicated font files.
- API RED: `go test -count=1 ./pkg/api/test -run '^TestMainFontStoreIsPrivate$'` failed (exit 1) because the suite store was `/home/umats/.config/pdfcpu/fonts`, not `/tmp/pdfcpu_api_tests1365268954/pdfcpu/fonts`. Log: `.pi-herdsman/font-isolation-api-red.log`. Same command after API setup isolation passed (exit 0): `.pi-herdsman/font-isolation-api-green.log`.
- CLI RED: `go test -count=1 ./pkg/cli/test -run '^TestMainFontStoreIsPrivate$'` failed (exit 1) because the suite store was `/home/umats/.config/pdfcpu/fonts`, not `/tmp/pdfcpu_cli_tests3978843524/pdfcpu/fonts`. Log: `.pi-herdsman/font-isolation-cli-red.log`. Same command after CLI setup isolation passed (exit 0): `.pi-herdsman/font-isolation-cli-green.log`.
- `go test -count=20 ./pkg/pdfcpu/model -run '^(TestPrepBytesPropagatesFontErrors|TestTextRenderingUsesConfigurationFontRepository)$'` and `go test -count=3 ./pkg/api/test ./pkg/cli/test -run '^TestMainFontStoreIsPrivate$'` passed (exit 0). Log: `.pi-herdsman/font-restoration-focused.log`.
- Behavioral check: `go test -count=1 ./pkg/api/test ./pkg/cli/test -run '^TestMainFontStoreIsPrivate$'` passed while the shared installed `.gob` filenames and SHA256 hashes before/after were identical (`cmp` exit 0). Log: `.pi-herdsman/font-isolation-behavior.log`; manifests: `.pi-herdsman/font-store-{before,after}.sha256`.
- Fresh default-parallel `go test -count=1 ./... && go vet ./... && go build ./...` passed twice (exit 0), including a final run after setup cleanup edits. Logs: `.pi-herdsman/font-fix-preflight.log`, `.pi-herdsman/font-fix-preflight-final.log`. `go env GOFLAGS` was empty; no serialization override. Both original failures remain recorded, not erased by these successes.
- Defect-class sweep: `rg -n 'api.InstallFonts\(context.Background\(\), fonts\)' --glob '*test.go'` found exactly 2 suite startup writers, both patched. Reviewed the other 2 `TestMain` implementations: command tests default to stateless configuration or supplied test roots; core image tests do not bulk-install fixture fonts. Other isolated per-test font helpers do not perform these shared startup replacements.
- Ran gofmt and `git diff --check`; inspected the final diff. Restored only known tracked `pkg/samples/` test outputs and removed exactly `basic/FormDemo.pdf`, `bookmarks/bookmarkTree2Levels.pdf`, and `bookmarks/bookmarkTreeImported.pdf` below that directory. Upstream TIFF fixture changes remain staged.
- Binary staged-diff snapshots before/after compared equal (`cmp` exit 0): `.pi-herdsman/font-fix-{initial,final}-index.diff`. Both changed Go files were previously unstaged/unchanged, so no staged-file separation is needed. Bug and merge verification notes remain untracked/unstaged. No commits, push, branch change or merge abort performed.

## Remaining work

Owner review, merge-integration commit followed by a separate Conventional fix commit (suggested `fix(test): isolate integration font stores`), then PR/CI and merge only when green. Red/green evidence is persisted, but no red/green commits were made because this execution assignment expressly prohibited commits. Platform validation here was Linux only; CI must confirm the existing supported platform matrix.

See `specs/verifications/UPSTREAM_MERGE_2026-10-08.md` for integration evidence.
