// Package strsim provides string similarity metrics, edit distances, and
// phonetic encoding algorithms with a unified API and no dependencies.
//
// Every metric returns a similarity in [0, 1]; a string compared with itself
// scores exactly 1.0, two empty strings included, and none returns NaN.
// Edit-distance metrics also
// return the raw distance. All metrics compare runes, not bytes, unless their
// ASCIIOnly field is set.
//
// # Quick Start
//
//	// Edit distance
//	d := strsim.Levenshtein("kitten", "sitting")  // 3
//
//	// Normalized similarity [0, 1]
//	s := strsim.JaroWinklerSimilarity("martha", "marhta")  // ~0.961
//
//	// Phonetic encoding
//	code := strsim.Soundex("Robert")  // "R163"
//
//	// Find best match from a list
//	m := strsim.FindBestMatch("golang", candidates, strsim.NewJaroWinkler())
//
// # Interfaces
//
// All similarity metrics implement the [Metric] interface, allowing them to be
// used interchangeably. Distance metrics additionally implement [DistanceMetric].
// Phonetic encoders implement the [Encoder] interface, and Double Metaphone the
// [DualEncoder] interface.
//
// # Configuration
//
// Metrics and encoders are structs with exported fields. The zero value of
// every struct works and, except for [NYSIISEncoder] (where MaxLength 0 means
// unlimited), behaves like the matching New constructor. Metrics hold no
// state, so one value can be shared across goroutines.
//
// # Phonetic encoding
//
// Phonetic encoders read Latin letters only. Diacritics are folded first, in
// composed or decomposed form ("Muñoz" encodes like "Munoz", "Strauß" like
// "Strauss"); other characters are ignored.
package strsim

// Metric computes similarity between two strings.
// Implementations return values in [0, 1] where 1.0 means identical.
type Metric interface {
	Similarity(a, b string) float64
}

// DistanceMetric computes both raw edit distance and normalized similarity.
type DistanceMetric interface {
	Metric
	Distance(a, b string) int
}

// Encoder produces a phonetic encoding of a string.
type Encoder interface {
	Encode(s string) string
}

// DualEncoder produces primary and alternate phonetic encodings.
type DualEncoder interface {
	Encode(s string) (primary, alternate string)
}

// Match represents a candidate string and its similarity score.
type Match struct {
	// Value is the matched candidate string.
	Value string
	// Similarity is the similarity score in [0, 1].
	Similarity float64
	// Index is the position of this candidate in the original slice.
	Index int
}

// Compile-time interface compliance checks.
var (
	_ DistanceMetric = (*LevenshteinMetric)(nil)
	_ DistanceMetric = (*OSAMetric)(nil)
	_ DistanceMetric = (*DamerauLevenshteinMetric)(nil)
	_ DistanceMetric = (*LCSMetric)(nil)
	_ DistanceMetric = (*HammingMetric)(nil)
	_ Metric         = (*JaroMetric)(nil)
	_ Metric         = (*JaroWinklerMetric)(nil)
	_ Metric         = (*NgramMetric)(nil)
	_ Encoder        = (*SoundexEncoder)(nil)
	_ Encoder        = (*MetaphoneEncoder)(nil)
	_ Encoder        = (*NYSIISEncoder)(nil)
	_ DualEncoder    = (*DoubleMetaphoneEncoder)(nil)
)
