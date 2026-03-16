# Contributing to strsim

Thank you for your interest in contributing to strsim!

## Development Setup

```bash
git clone https://github.com/jcoruiz/strsim.git
cd strsim
go test ./...
```

### Requirements

- Go 1.22 or later
- No external dependencies needed

## Running Tests

```bash
# All tests
go test ./... -v

# With coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Benchmarks
go test -bench=. -benchmem ./...

# Race detector
go test -race ./...
```

## Code Standards

- Follow the [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- `go vet ./...` must pass with no warnings
- `gofmt -s` must produce no changes
- Test coverage must stay above 95%
- All exported symbols must have godoc comments
- Zero external dependencies — stdlib only

## Adding a New Algorithm

1. Create `algorithm.go` with the implementation
2. Create `algorithm_test.go` with table-driven tests covering:
   - Normal cases with known expected values
   - Empty strings (both empty, one empty)
   - Identical strings
   - Unicode (Japanese, emoji, accented characters)
   - Single character strings
3. Add benchmark functions (`Benchmark*`)
4. If it's a similarity metric, implement the `Metric` interface
5. If it's a distance metric, implement the `DistanceMetric` interface
6. If it's a phonetic encoder, implement the `Encoder` interface
7. Add convenience top-level functions
8. Update README.md

## Pull Request Process

1. Fork the repository
2. Create a feature branch from `main`
3. Make your changes
4. Run `go test -race -count=1 ./...`
5. Run `go vet ./...`
6. Submit a PR with a clear description

## Release Process

Releases follow [Semantic Versioning](https://semver.org/):

```bash
git tag v0.x.0
git push origin v0.x.0
gh release create v0.x.0 --title "v0.x.0" --notes "Release notes"
```
