# Code Review Reports

This directory contains individual code review reports for each source file in the go-chrono project.

## Review Criteria

Each report evaluates:
- **Code Style (20%)**: Naming conventions, comments, formatting, readability
- **Error Handling (25%)**: Error wrapping, boundary checks, panic prevention
- **Concurrency Safety (20%)**: Lock usage, race conditions, resource leaks
- **Architecture Design (25%)**: Interface design, coupling, extensibility
- **Performance (10%)**: Algorithm efficiency, memory usage

## Priority Levels

- **P0 - Critical**: Security vulnerabilities, data races, panic risks, resource leaks
- **P1 - Important**: Missing error handling, performance bottlenecks, design flaws
- **P2 - Normal**: Code duplication, naming issues, missing comments
- **P3 - Suggestions**: Optimization opportunities, style improvements

## Project-Specific Checks

Per `CLAUDE.md`:
1. Bilingual comments (English + Chinese) completeness
2. `sync.Mutex` protection for shared state
3. Error wrapping with `fmt.Errorf`
4. Function comments start with function name (not struct name)

## Reports

<!-- Add links as reports are generated -->
