# pdfcpu — Claude Code

Read `CONVENTIONS.md` before Git or GitHub operations.

<!-- BEGIN bigpowers:context-routing -->
## Context Routing

| Files | Context |
|-------|---------|
| `cmd/pdfcpu/**` | Read CLI commands and dispatch in `pkg/cli/`. |
| `pkg/api/**` | Read the corresponding PDF operation in `pkg/pdfcpu/`. |
| `pkg/pdfcpu/**` | Read `specs/tech-architecture/tech-stack.md` and nearby tests. |
<!-- END bigpowers:context-routing -->

<!-- BEGIN bigpowers:learned-preferences -->
## Learned User Preferences

- Use Conventional Commits.
- Use the team-pr workflow.
- Target `main` after the default-branch migration; the current branch is `master`.

## Workspace Facts

- The codebase is an established Go project.
- Read `specs/tech-architecture/tech-stack.md` for the observed architecture.
<!-- END bigpowers:learned-preferences -->

<!-- BEGIN bigpowers:project -->
## Project

pdfcpu is a Go library and command-line tool for validating and manipulating PDF files.
Stack: Go 1.26, Cobra CLI, and a public Go API.

## Commands

| Action | Command |
|--------|---------|
| Run | `go run ./cmd/pdfcpu` |
| Test | `go test ./...` |
| Build | `go build ./...` |
| Lint | `go vet ./...` |
| Preflight | `go test ./... && go vet ./... && go build ./...` |
| CI | `gh pr checks` when a PR exists |

## Architecture

`cmd/pdfcpu` parses Cobra commands; `pkg/cli` dispatches typed operations.
`pkg/api` owns reusable operation and I/O boundaries; `pkg/pdfcpu` implements PDF processing.

## Conventions

- DO keep changes in the existing package layers.
- DO pass `context.Context` through operations and wrap errors with context.
- DO preserve staged file writes for mutating operations.
- MUST use Conventional Commits: `<type>(<scope>): <description>`.

## Never

- NEVER edit `vendor/`; regenerate vendored files from dependencies instead.
- NEVER land changes directly on `master` or `main`; use a feature branch and team-pr workflow.
- NEVER dismiss reproducible Preflight or CI failures without the fix-or-log process.

## Agent Rules

- MUST read `specs/` and `CONVENTIONS.md` before writing code.
- MUST use relevant bigpowers skills for structured feature and bug work.
- MUST write planning output into `specs/`.
- MUST keep changes minimal and test changed behavior.
- MUST run Preflight before merging and check CI for an open PR.
<!-- END bigpowers:project -->
