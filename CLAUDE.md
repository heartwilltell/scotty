# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Scotty is a zero-dependency Go library for building simple CLI applications. It wraps Go's standard `flag.FlagSet` with subcommand support, struct tag-based config binding, and environment variable integration. All source lives in a single `scotty` package at the root — no subdirectories.

## Build / Test / Lint Commands

```bash
go build                                    # Build
go test ./...                               # Run all tests
go test -v -race -cover ./...               # Run tests (as CI does)
go test -run TestFuncName ./...             # Run a single test
golangci-lint run --timeout=3m              # Lint (30 linters, strict config in .golangci.yml)
```

## Architecture

**Core types and their roles:**

- **`Command`** (`command.go`) — Represents a CLI command. Holds a `Run` function, flags, subcommands (stored in a map), and parent pointer for hierarchy traversal. `Exec()` is the main entry point: it parses args, discovers subcommands, binds config, validates, then calls `Run`.

- **`FlagSet`** (`flagset.go`) — Wraps `flag.FlagSet`. Adds `*VarE` methods (e.g. `StringVarE`, `BoolVarE`) that bind a flag to both a CLI flag name and an environment variable. Also tracks bound config and required fields.

- **Config binding** (`config.go`) — Reflection-based system that reads struct tags (`flag`, `env`, `default`, `usage`, `required`) to auto-register flags. `BindConfig()` attaches a struct; `MustConfig[T]()` / `GetConfig[T]()` provide type-safe retrieval. Structs can implement `ConfigValidator` for custom validation.

- **`Error` / `RequiredFieldError`** (`error.go`) — Custom error types. `RequiredFieldError` reports missing required config fields with flag/env context.

- **Usage generation** (`usage.go`) — Builds formatted help text showing command chain, subcommand list, and flag documentation.

**Execution flow:** Create `Command` → add flags via `SetFlags` callback or `Flags()` → optionally `BindConfig()` → `Exec()` parses args → subcommand dispatch → config validation → `Run()`.

## Code Conventions

- **Zero external dependencies** — only stdlib. This is intentional; do not add third-party deps.
- **Table-driven tests** using `map[string]tcase` pattern with descriptive keys.
- **Test helpers** prefixed with `helper` (e.g. `helperDisableStdout`, `helperCatchPanic`).
- **`sync.Once`** for lazy `FlagSet` initialization on `Command`.
- **Generics** used for type-safe config access (`MustConfig[T]`, `GetConfig[T]`) and the `tern` helper.
- **Comments must end with a period** (enforced by `godot` linter).
- **Revive limits:** cognitive complexity ≤30, cyclomatic ≤15, function length ≤30 statements, line length ≤150, max 5 params, max 2 return values.
- Supported config field types: `string`, `bool`, `int`, `int64`, `uint`, `uint64`, `float64`, `time.Duration`.
