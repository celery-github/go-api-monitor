# Go API Monitor

A lightweight API health monitor written in Go.

## Features
- YAML config of endpoints
- Periodic checks on an interval
- Bounded concurrency
- JSON output (stdout) or append-to-file
- Graceful shutdown (Ctrl+C)
- Tests with httptest
- GitHub Actions CI
- Dockerfile

## Quickstart
```bash
go mod tidy
make once
make run
