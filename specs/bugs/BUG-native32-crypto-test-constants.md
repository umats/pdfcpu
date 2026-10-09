---
type: bug-fix-plan
context: plaintext hint round13 32-bit verification
bug_id: BUG-native32-crypto-test-constants
status: verified_unreleased
severity: low
scope: tests-only
---

# Native32 crypto-test compile failure

## Reproduction

Round13's full-package 386 compile reports constants at
`pkg/pdfcpu/crypto_test.go:580-581` overflow native int.
`TestNormalizeUnsignedPermission` skips execution on 32-bit, but the compiler
still checks its literal table entries `1<<31` and `1<<32-1`.
Evidence: `.pi-herdsman/hint-fix-round13/red-final.log`.

## Minimal correction and verification

Keep the existing 32-bit runtime skip and all 64-bit permission assertions.
Represent the two unsigned limits as int64 variables, then convert those
runtime variables to int inside the already-skipped test, following the nearby
`TestPermissionBytes` pattern. No production change or new behavior.

Capture failed full-package 386 test compilation and source manifest before
editing. After gofmt, run package-based focused 386 hint/permission regressions,
native permission and hint tests, fresh full Preflight and applicable race.
No source-list workaround is needed once the test package compiles.
Preserve exact sample-output bytes. Record commands, exit codes and final
candidate manifests. No commit, stash, dependency/vendor/CI, release or
consumer action. Keep this discovered test-only fix separable in future
Conventional Commit preparation, which remains unauthorized.

## Resolution

Test-only runtime int64 variables preserve the native32 skip and native64
assertions. Compile RED reproduced; actual package-based private/public386
focused tests, native focused tests, fresh full Preflight and applicable race
pass. Candidate-bound receipts:
`.pi-herdsman/hint-fix-round13-native32/handoff.md`.
Both independent round13 reviewers accepted this correction with no findings.
No full386-suite or ARM-runtime claim. No commit or release performed.
