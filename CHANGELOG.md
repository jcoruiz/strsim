# Changelog

All notable changes to this project are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the project uses
[Semantic Versioning](https://semver.org/). While on v0, a minor version may
change results.

## [Unreleased]

## [0.2.1] - 2026-11-04

### Fixed

- Cosine similarity of a string with itself, or of any two strings with
  proportional n-gram counts, could be 0.9999999999999998 instead of 1
  (`CosineSimilarity("martha", "martha")`), so `FindAboveThreshold` with a
  threshold of 1 missed exact matches. Every other Cosine value is unchanged.

## [0.2.0] - 2026-11-03

Phonetic codes and some Jaro scores change in this release; recompute any
stored codes or scores.

### Fixed

- Double Metaphone follows Philips' reference rules and matches Apache Commons
  Codec; on a 44,000-word corpus of names and English words, v0.1.0 differed
  on about 5% of them. Codes no longer contain spaces (`Horwicz` was `HR K`, now
  `HRTS`), `Y` counts as a vowel, `H` between vowels is kept, initial `PS` is
  silent, and the `VAN `/`SAN `/`MAC ` rules work because word boundaries are
  kept.
- Metaphone applies initial `WH` → `W`, `SCH` → `SK` and `CIA` → `X`, drops the
  `T` of `TCH` and the `G` of `GN`, and keeps `H` before a vowel, matching
  Apache Commons Codec except for an initial `CH` (always `X`).
- Jaro and Jaro-Winkler halve the transposition count with integer division,
  as Winkler's strcmp95, rapidfuzz, jellyfish and Commons Text do.
- Jaro-Winkler with `PrefixSize` above 10 no longer returns values above 1.0.
- `LevenshteinMetric{}` and other all-zero costs returned NaN similarity; they
  now use unit costs. Negative costs count as 0.
- `NgramMetric` with `Size` 0 scored any two non-empty strings 1.0, and a
  negative `Size` panicked; both now use bigrams.
- `FindTopN` and `FindAboveThreshold` could reorder candidates with equal
  scores; ties now keep the candidates' order.
- NYSIIS no longer drops the first letter of a two-letter `AY` key.

### Added

- Phonetic encoders fold Latin diacritics (U+00C0–U+017F) and drop combining
  accents before encoding: `Muñoz` encodes like `Munoz`, `Álvarez` as `A416`
  instead of `L162`, and decomposed (NFD) text like composed text.
- `ErrLengthMismatch`, returned by `Hamming` and used as the panic value of
  `HammingMetric.Distance`.
- `JaroWinklerMetric{}` uses the default threshold and prefix size, so every
  zero-value metric behaves like its constructor (except `NYSIISEncoder`).
- Runnable examples, fuzz tests for the metric invariants and the ASCII fast
  path, and a reference test against Apache Commons Codec.

### Changed

- Jaccard, Dice and Overlap are 2–3x faster and Cosine 1.3x, with a third of
  the allocations; Damerau-Levenshtein allocates one matrix instead of one row
  per character; phonetic encoders allocate less.
- CI tests Go 1.22 and the two supported Go releases, checks formatting, runs
  golangci-lint and short fuzzing.

## [0.1.0] - 2026-03-16

First release: Hamming, Levenshtein, OSA, Damerau-Levenshtein, LCS, Jaro,
Jaro-Winkler, n-gram Cosine/Jaccard/Dice/Overlap, Soundex, Metaphone, Double
Metaphone and NYSIIS, with batch matching and an ASCII fast path.

[Unreleased]: https://github.com/jcoruiz/strsim/compare/v0.2.1...HEAD
[0.2.1]: https://github.com/jcoruiz/strsim/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/jcoruiz/strsim/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/jcoruiz/strsim/releases/tag/v0.1.0
