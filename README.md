# strsim

[![CI](https://github.com/jcoruiz/strsim/actions/workflows/ci.yml/badge.svg)](https://github.com/jcoruiz/strsim/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jcoruiz/strsim)](https://goreportcard.com/report/github.com/jcoruiz/strsim)
[![codecov](https://codecov.io/gh/jcoruiz/strsim/branch/main/graph/badge.svg)](https://codecov.io/gh/jcoruiz/strsim)
[![Go Reference](https://pkg.go.dev/badge/github.com/jcoruiz/strsim.svg)](https://pkg.go.dev/github.com/jcoruiz/strsim)

String similarity metrics for Go: edit distance, token-based similarity and phonetic encoding, all in one package with a unified API.

**Try it in your browser at [strsim.than.cl](https://strsim.than.cl).** The playground runs strsim compiled to WebAssembly. It compares all 15 algorithms on your own text, shows why each one scores what it does, ranks a list and finds the threshold that fits your data.

## Features

- **15 algorithms** in a single zero-dependency package
- **Unified interfaces**: every metric implements `Metric` (similarity) or `DistanceMetric` (distance + similarity)
- **Unicode-aware**: metrics compare runes, not bytes; phonetic encoders fold Latin diacritics (`Muñoz` encodes like `Munoz`)
- **Safe edge cases**: similarities are always in `[0, 1]` and never NaN; a string compared with itself scores exactly 1.0, two empty strings included; zero-value structs work
- **Checked against reference implementations**: distances and Jaro-Winkler match rapidfuzz and jellyfish; phonetic codes match Apache Commons Codec (see [Phonetic Encoding](#phonetic-encoding))
- **Batch operations**: `FindBestMatch`, `FindTopN`, `FindAboveThreshold`, with deterministic tie order

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
    match := strsim.FindBestMatch("goland", candidates, strsim.NewJaroWinkler())
    fmt.Printf("%s (%.2f)\n", match.Value, match.Similarity)  // golang (0.93)
}
```

## Algorithms

### Edit Distance

| Algorithm | Function | Returns |
|---|---|---|
| Hamming | `Hamming(a, b)` | `int, error` (equal-length only, `ErrLengthMismatch` otherwise) |
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
| Sørensen-Dice (n-gram) | `DiceSimilarity(a, b)` | 2·intersection / (\|A\| + \|B\|) |
| Overlap Coefficient | `OverlapSimilarity(a, b)` | Intersection / min(\|A\|, \|B\|) |

### Phonetic Encoding

| Algorithm | Function | Returns |
|---|---|---|
| American Soundex | `Soundex(s)` | `string` (4-char code) |
| Metaphone | `Metaphone(s)` | `string` |
| Double Metaphone | `DoubleMetaphone(s)` | `string, string` (primary, alternate) |
| NYSIIS | `NYSIIS(s)` | `string` |

All phonetic functions have a corresponding `*Match(a, b)` variant returning `bool`.

Encoders read Latin letters: diacritics are folded first (`é` → `E`, `ñ` → `N`, `ß` → `SS`, in composed or decomposed form) and other characters are ignored. Double Metaphone keeps word boundaries, so `"San Jose"` and `"Mac Gregor"` get their special rules.

Codes match [Apache Commons Codec](https://commons.apache.org/proper/commons-codec/) on the 3,000 most common US surnames (`testdata/phonetic_reference.tsv`), with three deliberate differences:

- Metaphone encodes an initial `CH` as `X`, as in Philips' rules (`Church` → `XRX`; Commons Codec gives `KRX`).
- Double Metaphone codes never contain spaces (`Raj` → `RJ` / `R`; Commons Codec gives `"R "`).
- NYSIIS keeps a one-letter key (`Ash` → `A`; Commons Codec gives an empty string).

## Interfaces

Every metric implements at least one interface, making them interchangeable:

```go
// Metric returns a similarity in [0, 1], where 1.0 means identical.
type Metric interface {
    Similarity(a, b string) float64
}

// DistanceMetric also returns the raw edit distance.
type DistanceMetric interface {
    Metric
    Distance(a, b string) int
}

// Phonetic encoder.
type Encoder interface {
    Encode(s string) string
}

// Phonetic encoder with primary and alternate codes (Double Metaphone).
type DualEncoder interface {
    Encode(s string) (primary, alternate string)
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

Ties keep the order of `candidates`, so results are deterministic.

## Configurable Metrics

Metrics and encoders are plain structs. The zero value of each one behaves like its `New*` constructor, except `NYSIISEncoder`, where `MaxLength: 0` means unlimited. Metrics hold no state and are safe for concurrent use.

```go
// Custom Levenshtein costs
lev := &strsim.LevenshteinMetric{
    InsertCost:  1,
    DeleteCost:  1,
    ReplaceCost: 2,
}
dist := lev.Distance("kitten", "sitting")  // 5

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
code := enc.Encode("Washington")  // W25235
```

## ASCII Fast Path

The edit-distance metrics, Jaro, and Jaro-Winkler have an `ASCIIOnly` mode that compares bytes instead of runes. On ASCII input it returns exactly the same results as the default mode; on multi-byte UTF-8 input it returns **incorrect results**. Use it only when you know your input is ASCII (identifiers, codes, URLs, English text).

```go
m := &strsim.LevenshteinMetric{ASCIIOnly: true}
dist := m.Distance("kitten", "sitting")

// Works with batch operations too
jw := &strsim.JaroWinklerMetric{ASCIIOnly: true}
match := strsim.FindBestMatch("query", candidates, jw)
```

| Algorithm | Rune (ns/op) | ASCII (ns/op) | Speedup |
|---|---|---|---|
| Hamming | 105 | 10 | **10.5x** |
| Damerau-Levenshtein | 17,248 | 6,243 | **2.8x** |
| Levenshtein | 2,752 | 2,574 | 1.1x |
| Jaro-Winkler | 315 | 332 | 0.9x |

Inputs: 43-char ASCII strings; median of five runs. On these inputs the rune path of Jaro-Winkler is as fast as the byte path, so `ASCIIOnly` only saves it allocations (96 B instead of 448 B).

## Benchmarks

Measured on an AMD Ryzen 9 9950X3D with Go 1.22 and strsim v0.2.1, median of five runs (`go test -bench=. -benchmem -count=5`):

| Algorithm | ns/op | B/op | allocs/op |
|---|---|---|---|
| Soundex | 33 | 4 | 1 |
| Jaro-Winkler | 35 | 16 | 2 |
| Hamming | 37 | 0 | 0 |
| Metaphone | 47 | 16 | 2 |
| Jaro | 51 | 32 | 2 |
| NYSIIS | 72 | 32 | 2 |
| Double Metaphone | 130 | 21 | 3 |
| Jaccard (bigram) | 510 | 896 | 6 |
| LCS | 512 | 448 | 2 |
| Sørensen-Dice (bigram) | 521 | 896 | 6 |
| Overlap (bigram) | 548 | 896 | 6 |
| Cosine (bigram) | 622 | 896 | 6 |
| Levenshtein | 947 | 448 | 2 |
| OSA | 1,134 | 672 | 3 |
| Damerau-Levenshtein | 6,134 | 7,358 | 6 |

Inputs differ per algorithm (26-letter strings for edit distances, names for phonetic encoders); see the `Benchmark*` functions in the test files.

## Why strsim?

The Go ecosystem for string similarity is fragmented across 8+ libraries, each with a subset of algorithms, inconsistent APIs, and known bugs:

- **go-edlib**: NaN on empty strings, float32 precision, no phonetic
- **strutil**: NaN on empty strings, negative similarity with custom costs
- **smetrics**: operates on bytes not runes (broken Unicode), Soundex panics on empty input
- **matchr**: most complete phonetic support but GPLv3 licensed

strsim consolidates everything into one MIT-licensed package with correct Unicode handling, consistent APIs, and no edge-case surprises.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Changes are listed in [CHANGELOG.md](CHANGELOG.md).

## License

MIT
