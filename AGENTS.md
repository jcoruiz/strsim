# strsim

Go library of string similarity metrics, edit distances and phonetic encoders. Standard library only; the minimum Go version is the one in `go.mod`.

## Invariants

- Every `Similarity` returns a value in [0, 1], never NaN, and a string compared with itself scores exactly 1.0, two empty strings included (`FuzzMetrics`, which also covers the `NgramMetric` methods).
- Zero-value structs behave like their `New*` constructor (`TestZeroValues`), except `NYSIISEncoder`, where `MaxLength: 0` means unlimited. Keep it that way when adding fields.
- The `ASCIIOnly` byte paths in `ascii.go` return exactly what the rune paths return on ASCII input (`FuzzASCIIOnly`); change both together.
- Batch functions keep the candidates' order on ties (stable sort).
- Metric values hold no state, so they are safe for concurrent use.

## Reference implementations

When changing an algorithm, compare against these, not against the current output:

- Phonetic encoders match Apache Commons Codec on `testdata/phonetic_reference.tsv`, except for the three deviations documented in `TestPhoneticReference` (`reference_test.go`). Regenerate the file with Commons Codec when the reference version changes.
- Distances, Jaro and Jaro-Winkler match rapidfuzz and jellyfish, including the integer halving of transpositions from Winkler's strcmp95.
- Phonetic input is folded with the table in `phonetic.go` (Latin-1 Supplement and Latin Extended-A); other scripts are dropped.

## Checks

```bash
go vet ./...
go test -race ./...
golangci-lint run   # CI pins the version in .github/workflows/ci.yml
go test -run '^$' -fuzz '^FuzzMetrics$' -fuzztime 30s .   # also FuzzASCIIOnly, FuzzPhonetic
```

## Releases

Tags are `vX.Y.Z`. While on v0, a change in results bumps the minor version. Move the `Unreleased` notes in `CHANGELOG.md` under the new version, tag it, and publish the GitHub release with that section as notes.

# Git conventions

All git text is written in **English**: commit messages, PR titles and descriptions, branch names.

## Commits: Conventional Commits

Format: `<type>(<scope>): <subject>`

- **type** (required): `feat`, `fix`, `refactor`, `perf`, `test`, `docs`, `style`, `build`, `ci`, `chore`, `revert`.
- **scope** (optional): the area touched, in lowercase, such as `levenshtein`, `damerau`, `hamming`, `jaro`, `lcs`, `ngram`, `phonetic`, `batch`, `ascii`, `examples`.
- **subject**: imperative mood, lowercase start, no trailing period, ≤ 72 chars for the whole header.
- **body** (optional): blank line after the header, then 1–3 short lines or bullets explaining *what/why*. No multi-paragraph prose.
- **Breaking changes**: `feat(scope)!: ...` plus a `BREAKING CHANGE: ...` footer.
- One logical change per commit; don't mix refactors with features.
- No `Co-Authored-By` or AI attribution lines.
- Enforced locally by `scripts/hooks/commit-msg` (enable once per clone with `git config core.hooksPath scripts/hooks`) and in CI by the `commits` job, which runs the same hook on every pushed commit and on pull request titles. (A GitHub ruleset can't do it: this account doesn't support commit message rules.)

Examples:

```
feat(phonetic): fold latin diacritics before encoding
fix(jaro): cap the winkler prefix bonus at 1.0
perf(ngram): key n-grams by substring to cut allocations
docs(readme): update benchmark numbers from a fresh run
```

## Pull requests

- Title follows the same format: it becomes the squash commit, and the `commits` job checks it.
- Description: short summary of the change and how it was verified. No AI attribution.

## Branches

`<type>/<short-kebab-description>`, such as `feat/ngram-ascii-path`, `fix/jaro-transpositions`.
