# CLAUDE.md - Obfuscator Project Guidelines

## Build Commands
- Build & run: `go run . -in=test -out=test_obfuscated -key=123 -header=header.txt`
- Install: `go install`
- Full test suite: `go test ./...`
- Single package test: `go test ./[package]` (e.g., `go test ./lexer`)
- Single test: `go test ./[package] -run TestName` (e.g., `go test ./lexer -run TestDeclare`)

## Code Style Guidelines
- **Imports**: Standard library first, then third-party packages
- **Formatting**: Follow standard Go formatting (`go fmt`)
- **Naming**: 
  - Packages: lowercase, single word (lexer, parser, obfuscator)
  - Functions/Methods: CamelCase
  - Variables: camelCase
- **Error Handling**: Use direct error checking with appropriate logging
- **Project Structure**: Organized by component packages (lexer, parser, transpiler, etc.)

## R Code Requirements
- if statements MUST have curly braces (not: `if(TRUE) 1 else 0` but: `if(TRUE) {1} else {0}`)
- Functions MUST have curly braces (not: `\(x) x + 1` but: `\(x) {x + 1}`)
- Expressions in curly braces outside of normal block contexts are not supported
  - Rewrite `tryCatch({x + 1})` as `fn <- \(){x + 1};tryCatch(fn())`