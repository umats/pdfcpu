# Conventions

## Conventional Commits

- MUST format commits as `<type>(<scope>): <description>`.
- MUST use `feat` for features and `fix` for bug fixes.
- MUST mark breaking changes with `!` or a `BREAKING CHANGE:` footer.
- MUST keep discovered fixes in separate commits within the same PR.

## Git Workflow

- ALWAYS work on a feature branch and use the `team-pr` workflow.
- NEVER land changes directly on `master` or `main`.
- MUST target `main` after the default-branch migration; the current default branch is `master`.
- MUST NOT rename branches or change remote settings as part of unrelated work.

## Project Structure

- DO put command parsing in `cmd/pdfcpu/` and dispatch in `pkg/cli/`.
- DO put reusable file and stream operations in `pkg/api/`.
- DO put PDF parsing, transformations, validation, and writing in `pkg/pdfcpu/`.
- MUST write planning and verification artifacts in `specs/`.
- MUST preserve `specs/tech-architecture/tech-stack.md` as the observed architecture map.
- NEVER modify `vendor/` directly; `go mod vendor` fully manages its contents.
- DO notify the user when a vendor fix is necessary; running `go mod vendor` is allowed.

## Bigpowers Scripts

- ALWAYS use the full `~/.pi/agent/npm/node_modules/bigpowers/scripts/<script-name>` path for bigpowers scripts.
- DO expand `~` from the shell; NEVER hard-code a user's home directory.

## Go Changes

- NEVER modify `.golangci.yml`.
- ALWAYS fix lint errors in code; DO ignore a lint error only when absolutely necessary.
- MUST use `gofmt` on changed Go files.
- DO pass `context.Context` through cancellable operations.
- DO wrap returned errors with operation context and `%w`.
- DO preserve staged file writes and cleanup on failure.
- DO add focused tests for changed behavior and bug fixes.
- DO reuse existing package seams before adding dependencies or abstractions.

## Always Green / Shift Left

- MUST run Preflight before forward implementation or merge: `go test ./... && go vet ./... && go build ./...`.
- MUST verify passing CI checks before merging an open PR.
- MUST stop forward work when Preflight or CI is red.
- DO fix defects early: a development fix costs less than a production fix.

## Discovered Defects

- MUST run the fix-or-log ladder for reproducible gate failures.
- DO use `quick-fix` for trivial data-only repairs within its guardrails.
- DO use `fix-bug` when investigation or a regression test is needed.
- MUST log blocked reproductions in `specs/bugs/` and stop forward work until triaged.
- MUST ship discovered fixes in separate Conventional Commits within the same PR.
- MUST NOT dismiss red gates as pre-existing, unrelated, not introduced here, or out of scope.

## Defensive Code

No additional defensive-code category was selected during seeding.
DO preserve existing PDF input limits, timeouts, and network protections when changing affected code.
