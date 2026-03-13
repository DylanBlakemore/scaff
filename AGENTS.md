# AGENTS.md

This file tells LLM-based coding agents how to work effectively in this repository.

## Non-negotiables (repo workflow)

- **Format**: run `make format` after code changes.
- **Lint/style**: run `make lint` and fix issues you introduced.
- **Tests**: run `make test` and ensure it’s green.
- **Security/deps**: run `make security` for vulnerability + dependency checks.
- **CI parity**: prefer `make ci` before finalizing a change (it runs all of the above).

If you cannot run commands (or they fail for infra reasons), explain what you would run and why, and keep changes minimal.

## Go language guidelines

### Code style and ergonomics

- **Keep it simple and idiomatic**: prefer the standard library; avoid clever abstractions.
- **Minimize comments**: comment when necessary to explain WHY, not WHAT.
- **Prefer early returns**: reduce indentation; use guard clauses to keep happy paths flat.
- **Naming**: choose meaningful names that are as short as possible while still clear.
- **Testability**: keep “edges” (filesystem/network/process/env) at the edges; make core logic stateless/pure where practical.
- **Dependency injection (lightweight)**: inject dependencies when it improves testability; avoid over-engineering.
- **Errors**:
  - Return errors, don’t panic, except for truly unrecoverable programmer errors.
  - Wrap errors with context using `fmt.Errorf("...: %w", err)` at boundaries.
  - Don’t double-wrap; add context once per layer.
  - Prefer sentinel errors only when callers must branch; otherwise use typed/contextual errors.
- **Context**:
  - Accept `context.Context` as the first parameter for calls that do I/O, RPC, or may block.
  - Never store contexts in structs; don’t pass `nil` contexts.
- **Logging**:
  - Don’t log and return the same error at the same layer (avoid duplicate logs).
  - Keep logs structured/consistent; avoid noisy debug output in normal operation.
- **Concurrency**:
  - Avoid goroutines unless needed; always reason about cancellation and shutdown.
  - Avoid data races; prefer ownership, channels, or `sync` primitives appropriately.
- **Performance**:
  - Optimize only when there’s evidence; prefer clarity.
  - Avoid premature micro-optimizations (e.g., manual inlining, unsafe).

### Project layout assumptions

This is a Go CLI using Cobra (`github.com/spf13/cobra`) with the typical layout:

- `main.go` calls into `cmd.Execute()`.
- Commands live under `cmd/`.
- **`cmd/` is the user interface layer**: it should only deal with user input/output (args, flags, env, prompts, printing, exit codes).
- **All business logic must live under `internal/`** (and be callable from the CLI layer).

When adding new functionality, prefer putting user-facing behavior behind a Cobra command/subcommand rather than adding logic directly in `main.go`. Keep `cmd/` thin by delegating to `internal/` packages.

### Public API hygiene

- Keep exported identifiers minimal.
- Write doc comments for exported types/functions where the intent is not obvious.
- Prefer smaller packages with clear responsibilities; avoid cyclic dependencies.

### Configuration and environment

- Prefer explicit flags/env vars over hidden global state.
- Treat environment variables as inputs: validate early and fail with actionable messages.
- Don’t bake secrets into code, tests, or fixtures.

## CLI-specific guidelines (Cobra-based)

### UX and command design

- **Consistency**:
  - Use consistent flag names across commands (e.g., `--output`, `--format`, `--json`).
  - Keep help text short and scannable; include examples for non-trivial commands.
- **Exit codes**:
  - `0` for success.
  - Non-zero for failures; distinguish user error vs internal error when feasible.
- **Output discipline**:
  - Write primary command output to **stdout**.
  - Write diagnostics/errors to **stderr**.
  - Make output stable for scripting; avoid extra chatter unless `--verbose` is set.
- **Error messages**:
  - Prefer actionable errors: what failed, why it matters, how to fix.
  - Avoid printing Go internals (stack traces) by default.

### Flag parsing and validation

- Validate flags/args early in `PreRunE`/`Args` validation.
- Prefer `RunE` (returning an error) over `Run` so failures propagate uniformly.
- Don’t use global mutable variables for flags unless that’s the established Cobra pattern here; prefer binding flags to a command-specific options struct.

### Command implementation structure

Recommended pattern for new commands:

- Define an `options` struct for flags.
- Build the command in `NewXCommand()` (or similar) and bind flags there.
- Keep business logic in `internal/` (separate package/functions) so it’s testable without CLI parsing.

### Testing CLI behavior

- Prefer unit tests for the underlying logic plus focused tests for CLI wiring.
- In CLI tests:
  - Set args via Cobra (`cmd.SetArgs(...)`) and capture output buffers.
  - Assert on stdout/stderr content and exit behavior (error returned).
  - Avoid golden files unless output is large and stable.

## Change discipline for agents

- Keep PRs/changesets focused: one feature/fix per change where possible.
- Avoid broad refactors unless explicitly requested.
- Update docs/help text when behavior changes.
- If you introduce a new dependency, justify it and ensure `make security` remains clean.

