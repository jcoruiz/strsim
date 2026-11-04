package strsim

import (
	"math"
	"math/bits"
	"unicode/utf8"
)

// NgramMetric computes string similarity using n-gram based coefficients.
// It supports Cosine, Jaccard, Sorensen-Dice, and Overlap similarity.
//
// Strings shorter than Size form a single n-gram, so short strings are only
// similar when they are equal.
type NgramMetric struct {
	// Size is the n-gram size, in runes. Values <= 0 use the default of 2
	// (bigrams).
	Size int
}

// NewNgram returns a new NgramMetric with Size = 2 (bigrams).
func NewNgram() *NgramMetric {
	return &NgramMetric{Size: 2}
}

// size returns the n-gram size to use.
func (m *NgramMetric) size() int {
	if m.Size <= 0 {
		return 2
	}

	return m.Size
}

// ngramFrequency returns the term-frequency map of the rune n-grams of s. The
// keys are substrings of s, so building the map does not copy the grams. If s
// has fewer than n runes, the whole string is a single n-gram.
func ngramFrequency(s string, n int) map[string]int {
	// Compare invalid UTF-8 the way []rune decodes it: each bad byte is
	// U+FFFD.
	if !utf8.ValidString(s) {
		s = string([]rune(s))
	}

	// starts holds the byte offset of every rune, plus len(s).
	starts := make([]int, 0, len(s)+1)
	for i := range s {
		starts = append(starts, i)
	}

	runes := len(starts)
	starts = append(starts, len(s))

	if runes <= n {
		return map[string]int{s: 1}
	}

	freq := make(map[string]int, runes-n+1)
	for i := 0; i+n <= runes; i++ {
		freq[s[starts[i]:starts[i+n]]]++
	}

	return freq
}

// ngramIntersection returns the number of distinct n-grams that a and b share.
func ngramIntersection(a, b map[string]int) int {
	if len(a) > len(b) {
		a, b = b, a
	}

	shared := 0

	for gram := range a {
		if _, ok := b[gram]; ok {
			shared++
		}
	}

	return shared
}

// Cosine returns the cosine similarity between a and b using term-frequency
// n-gram vectors: dot(A, B) / (||A|| * ||B||). Returns 1.0 when both strings
// are empty, 0.0 when only one is empty.
func (m *NgramMetric) Cosine(a, b string) float64 {
	if a == "" || b == "" {
		return emptySimilarity(a, b)
	}

	freqA := ngramFrequency(a, m.size())
	freqB := ngramFrequency(b, m.size())

	// Integer sums, exact, so the proportional case below can be detected.
	var dot, normA, normB uint64

	for gram, countA := range freqA {
		if countB, ok := freqB[gram]; ok {
			dot += uint64(countA) * uint64(countB)
		}
	}

	for _, c := range freqA {
		normA += uint64(c) * uint64(c)
	}

	for _, c := range freqB {
		normB += uint64(c) * uint64(c)
	}

	// Proportional count vectors, such as those of two equal strings, have a
	// cosine of exactly 1, but the division below can round it to
	// 0.9999999999999998. dot² = |A|²·|B|² is that case, compared in 128 bits.
	dotHi, dotLo := bits.Mul64(dot, dot)
	normHi, normLo := bits.Mul64(normA, normB)
	if dotHi == normHi && dotLo == normLo {
		return 1.0
	}

	// Rounding can push nearly proportional vectors a hair above 1.
	return min(float64(dot)/(math.Sqrt(float64(normA))*math.Sqrt(float64(normB))), 1.0)
}

// Jaccard returns the Jaccard similarity between a and b using set semantics:
// |intersection| / |union|. Returns 1.0 when both strings are empty.
func (m *NgramMetric) Jaccard(a, b string) float64 {
	if a == "" || b == "" {
		return emptySimilarity(a, b)
	}

	setA := ngramFrequency(a, m.size())
	setB := ngramFrequency(b, m.size())
	intersection := ngramIntersection(setA, setB)

	union := len(setA) + len(setB) - intersection

	return float64(intersection) / float64(union)
}

// Dice returns the Sorensen-Dice similarity between a and b using set
// semantics: 2 * |intersection| / (|A| + |B|). Returns 1.0 when both strings
// are empty.
func (m *NgramMetric) Dice(a, b string) float64 {
	if a == "" || b == "" {
		return emptySimilarity(a, b)
	}

	setA := ngramFrequency(a, m.size())
	setB := ngramFrequency(b, m.size())
	intersection := ngramIntersection(setA, setB)

	return 2.0 * float64(intersection) / float64(len(setA)+len(setB))
}

// Overlap returns the overlap coefficient between a and b:
// |intersection| / min(|A|, |B|). Returns 1.0 when both strings are empty.
func (m *NgramMetric) Overlap(a, b string) float64 {
	if a == "" || b == "" {
		return emptySimilarity(a, b)
	}

	setA := ngramFrequency(a, m.size())
	setB := ngramFrequency(b, m.size())
	intersection := ngramIntersection(setA, setB)

	return float64(intersection) / float64(min(len(setA), len(setB)))
}

// Similarity returns the Sorensen-Dice similarity, which is the default
// similarity coefficient for n-gram based comparison.
func (m *NgramMetric) Similarity(a, b string) float64 {
	return m.Dice(a, b)
}

// CosineSimilarity returns the cosine similarity between a and b using bigrams.
func CosineSimilarity(a, b string) float64 {
	return NewNgram().Cosine(a, b)
}

// JaccardSimilarity returns the Jaccard similarity between a and b using
// bigrams.
func JaccardSimilarity(a, b string) float64 {
	return NewNgram().Jaccard(a, b)
}

// DiceSimilarity returns the Sorensen-Dice similarity between a and b using
// bigrams.
func DiceSimilarity(a, b string) float64 {
	return NewNgram().Dice(a, b)
}

// OverlapSimilarity returns the overlap coefficient between a and b using
// bigrams.
func OverlapSimilarity(a, b string) float64 {
	return NewNgram().Overlap(a, b)
}

// emptySimilarity is the similarity when at least one string is empty: 1.0 if
// both are, 0.0 otherwise.
func emptySimilarity(a, b string) float64 {
	if a == b {
		return 1.0
	}

	return 0.0
}
