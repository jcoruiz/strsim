# Contributing to strsim

## Setup

```bash
git clone https://github.com/jcoruiz/strsim.git
cd strsim
git config core.hooksPath scripts/hooks   # checks commit messages
go test ./...
```

Go 1.22 or later. No external dependencies.

## Before opening a pull request

- `go vet ./...` and `go test -race ./...` pass.
- `gofmt -s -l .` prints nothing and `golangci-lint run` is clean.
- Coverage stays above 95% (`go test -cover .`).
- Exported symbols have godoc comments; behaviour that users can see has a runnable `Example` or a test.
- Commit messages and the PR title follow the conventions in [AGENTS.md](AGENTS.md#git-conventions); user-visible changes get a line under `Unreleased` in [CHANGELOG.md](CHANGELOG.md).

Style follows the [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md).

## Adding an algorithm

1. Implement it in `algorithm.go` with a `New*` constructor, package-level helpers, and the `Metric`, `DistanceMetric` or `Encoder` interface (add it to the compile-time checks in `strsim.go`).
2. Make the zero value work like the constructor and add it to `TestZeroValues`.
3. Write table-driven tests with values from a reference implementation, covering empty strings, identical strings, single characters and Unicode; add it to the fuzz targets in `fuzz_test.go` and a `Benchmark*`.
4. Document it in README.md.

The invariants the package promises are listed in [AGENTS.md](AGENTS.md).
