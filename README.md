# strsim

[![CI](https://github.com/jcoruiz/strsim/actions/workflows/ci.yml/badge.svg)](https://github.com/jcoruiz/strsim/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jcoruiz/strsim)](https://goreportcard.com/report/github.com/jcoruiz/strsim)
[![codecov](https://codecov.io/gh/jcoruiz/strsim/branch/main/graph/badge.svg)](https://codecov.io/gh/jcoruiz/strsim)
[![Go Reference](https://pkg.go.dev/badge/github.com/jcoruiz/strsim.svg)](https://pkg.go.dev/github.com/jcoruiz/strsim)

Comprehensive string similarity metrics for Go. Edit distance, token-based similarity, and phonetic encoding — all in one package with a unified API.

## Features

- **15+ algorithms** in a single zero-dependency package
- **Unified interfaces** — every metric implements `Metric` (similarity) or `DistanceMetric` (distance + similarity)
- **Correct Unicode handling** — all algorithms operate on runes, not bytes
- **Safe edge case handling** — empty strings return 1.0 similarity (never NaN), mismatched lengths return 0.0
- **Batch operations** — `FindBestMatch`, `FindTopN`, `FindAboveThreshold`
- **Zero dependencies** — stdlib only

## Quick Start

```bash
go get github.com/jcoruiz/strsim
```

```go
package main

import (
    "fmt"
    "github.com/jcoruiz/strsim"
)

func main() {
    // Edit distance
    fmt.Println(strsim.Levenshtein("kitten", "sitting"))  // 3

    // Normalized similarity [0, 1]
    fmt.Println(strsim.JaroWinklerSimilarity("martha", "marhta"))  // ~0.961

    // Phonetic matching
    fmt.Println(strsim.SoundexMatch("Robert", "Rupert"))  // true

    // Find best match from candidates
    candidates := []string{"golang", "gold", "golem", "gopher", "python"}
    match := strsim.FindBestMatch("go", candidates, strsim.NewJaroWinkler())
    fmt.Printf("%s (%.2f)\n", match.Value, match.Similarity)
}
```

## Algorithms

### Edit Distance

| Algorithm | Function | Returns |
|---|---|---|
| Hamming | `Hamming(a, b)` | `int, error` (equal-length only) |
| Levenshtein | `Levenshtein(a, b)` | `int` |
| Optimal String Alignment | `OSA(a, b)` | `int` |
| Damerau-Levenshtein | `DamerauLevenshtein(a, b)` | `int` |
| Longest Common Subsequence | `LCSDistance(a, b)` | `int` |

All edit distance functions have a corresponding `*Similarity` variant returning a normalized `float64` in `[0, 1]`.

### Token-based / Set-based Similarity

| Algorithm | Function | Notes |
|---|---|---|
| Jaro | `JaroSimilarity(a, b)` | Matching characters + transpositions |
| Jaro-Winkler | `JaroWinklerSimilarity(a, b)` | Jaro + prefix bonus |
| Cosine (n-gram) | `CosineSimilarity(a, b)` | TF-weighted n-gram vectors |
| Jaccard (n-gram) | `JaccardSimilarity(a, b)` | Set intersection / union |
| Sørensen-Dice (n-gram) | `DiceSimilarity(a, b)` | 2·intersection / (|A| + |B|) |
| Overlap Coefficient | `OverlapSimilarity(a, b)` | Intersection / min(|A|, |B|) |

### Phonetic Encoding

| Algorithm | Function | Returns |
|---|---|---|
| American Soundex | `Soundex(s)` | `string` (4-char code) |
| Metaphone | `Metaphone(s)` | `string` |
| Double Metaphone | `DoubleMetaphone(s)` | `string, string` (primary, alternate) |
| NYSIIS | `NYSIIS(s)` | `string` |

All phonetic functions have a corresponding `*Match(a, b)` variant returning `bool`.

## Interfaces

Every metric implements at least one interface, making them interchangeable:

```go
// Similarity metric — returns [0, 1] where 1.0 = identical.
type Metric interface {
    Similarity(a, b string) float64
}

// Distance metric — also returns raw edit distance.
type DistanceMetric interface {
    Metric
    Distance(a, b string) int
}

// Phonetic encoder.
type Encoder interface {
    Encode(s string) string
}
```

### Using interfaces

```go
// Use any Metric interchangeably
func findSimilar(query string, items []string, m strsim.Metric) {
    for _, item := range items {
        if m.Similarity(query, item) > 0.8 {
            fmt.Println(item)
        }
    }
}

// Works with any metric
findSimilar("golang", items, strsim.NewJaroWinkler())
findSimilar("golang", items, strsim.NewLevenshtein())
findSimilar("golang", items, strsim.NewDamerauLevenshtein())
```

## Batch Operations

```go
candidates := []string{"golang", "gold", "golem", "gopher", "python", "ruby"}
m := strsim.NewJaroWinkler()

// Best single match
best := strsim.FindBestMatch("go", candidates, m)

// Top N matches, sorted by similarity descending
top3 := strsim.FindTopN("go", candidates, 3, m)

// All matches above threshold
matches := strsim.FindAboveThreshold("go", candidates, 0.7, m)
```

## Configurable Metrics

Most algorithms have configurable variants:

```go
// Custom Levenshtein costs
lev := &strsim.LevenshteinMetric{
    InsertCost:  1,
    DeleteCost:  1,
    ReplaceCost: 2,
}
dist := lev.Distance("kitten", "sitting")

// Custom Jaro-Winkler parameters
jw := &strsim.JaroWinklerMetric{
    BoostThreshold: 0.7,
    PrefixSize:     4,
}
sim := jw.Similarity("martha", "marhta")

// Custom n-gram size
ng := &strsim.NgramMetric{Size: 3}  // trigrams instead of bigrams
sim = ng.Cosine("night", "nacht")

// Custom phonetic code length
enc := &strsim.SoundexEncoder{MaxLength: 6}
code := enc.Encode("Washington")
```

## ASCII Fast Path

All metrics support an `ASCIIOnly` mode that skips rune conversion for faster processing of pure-ASCII input. This produces **incorrect results for multi-byte UTF-8 strings** — use only when you know your input is ASCII (identifiers, codes, URLs, English text).

```go
// Up to 11x faster for ASCII-only input
m := &strsim.LevenshteinMetric{
    InsertCost: 1, DeleteCost: 1, ReplaceCost: 1,
    ASCIIOnly: true,
}
dist := m.Distance("kitten", "sitting")

// Works with batch operations too
jw := &strsim.JaroWinklerMetric{
    BoostThreshold: 0.7, PrefixSize: 4,
    ASCIIOnly: true,
}
match := strsim.FindBestMatch("query", candidates, jw)
```

| Algorithm | Rune (ns/op) | ASCII (ns/op) | Speedup |
|---|---|---|---|
| Hamming | 117 | 11 | **11x** |
| Damerau-Levenshtein | 17,477 | 6,804 | **2.6x** |
| Jaro-Winkler | 493 | 234 | **2.1x** |
| Levenshtein | 3,144 | 2,787 | 1.1x |

Inputs: 43-char ASCII strings.

## Benchmarks

Measured on AMD Ryzen, Go 1.22 (`go test -bench=. -benchmem`):

| Algorithm | ns/op | B/op | allocs/op |
|---|---|---|---|
| Soundex | 30 | 21 | 3 |
| Hamming | 39 | 0 | 0 |
| Jaro-Winkler | 48 | 16 | 2 |
| Jaro | 54 | 32 | 2 |
| Metaphone | 56 | 32 | 3 |
| NYSIIS | 86 | 64 | 4 |
| Double Metaphone | 134 | 32 | 4 |
| LCS | 534 | 448 | 2 |
| Cosine (bigram) | 850 | 507 | 18 |
| Levenshtein | 1017 | 448 | 2 |
| OSA | 1030 | 672 | 3 |
| Jaccard (bigram) | 1172 | 1084 | 22 |

Inputs: "kitten"/"sitting" for edit distance, "Schneider" for phonetic, "night"/"nacht" for n-gram.

## Why strsim?

The Go ecosystem for string similarity is fragmented across 8+ libraries, each with a subset of algorithms, inconsistent APIs, and known bugs:

- **go-edlib** — NaN on empty strings, float32 precision, no phonetic
- **strutil** — NaN on empty strings, negative similarity with custom costs
- **smetrics** — operates on bytes not runes (broken Unicode), Soundex panics on empty input
- **matchr** — most complete phonetic support but GPLv3 licensed

strsim consolidates everything into one MIT-licensed package with correct Unicode handling, consistent APIs, and no edge-case surprises.

## License

MIT
