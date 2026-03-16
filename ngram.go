package strsim

import (
	"math"
)

// NgramMetric computes string similarity using n-gram based coefficients.
// It supports Cosine, Jaccard, Sorensen-Dice, and Overlap similarity.
type NgramMetric struct {
	// Size is the n-gram size. Default: 2 (bigrams).
	Size int
}

// NewNgram returns a new NgramMetric with Size = 2 (bigrams).
func NewNgram() *NgramMetric {
	return &NgramMetric{Size: 2}
}

// ngramFrequency returns the term-frequency map of n-grams for the given rune
// slice. If the slice is shorter than n, the entire string is treated as a
// single n-gram.
func ngramFrequency(runes []rune, n int) map[string]int {
	freq := make(map[string]int)

	if len(runes) == 0 {
		return freq
	}

	if len(runes) < n {
		freq[string(runes)] = 1
		return freq
	}

	for i := 0; i <= len(runes)-n; i++ {
		gram := string(runes[i : i+n])
		freq[gram]++
	}

	return freq
}

// ngramSet returns the set of distinct n-grams for the given rune slice.
func ngramSet(runes []rune, n int) map[string]bool {
	freq := ngramFrequency(runes, n)
	set := make(map[string]bool, len(freq))

	for k := range freq {
		set[k] = true
	}

	return set
}

// Cosine returns the cosine similarity between a and b using term-frequency
// n-gram vectors: dot(A, B) / (||A|| * ||B||). Returns 1.0 when both strings
// are empty, 0.0 when only one is empty.
func (m *NgramMetric) Cosine(a, b string) float64 {
	ra := []rune(a)
	rb := []rune(b)

	if len(ra) == 0 && len(rb) == 0 {
		return 1.0
	}

	if len(ra) == 0 || len(rb) == 0 {
		return 0.0
	}

	freqA := ngramFrequency(ra, m.Size)
	freqB := ngramFrequency(rb, m.Size)

	var dot float64

	for gram, countA := range freqA {
		if countB, ok := freqB[gram]; ok {
			dot += float64(countA) * float64(countB)
		}
	}

	var normA, normB float64

	for _, c := range freqA {
		normA += float64(c) * float64(c)
	}

	for _, c := range freqB {
		normB += float64(c) * float64(c)
	}

	denom := math.Sqrt(normA) * math.Sqrt(normB)
	if denom == 0 {
		return 0.0
	}

	return dot / denom
}

// Jaccard returns the Jaccard similarity between a and b using set semantics:
// |intersection| / |union|. Returns 1.0 when both strings are empty.
func (m *NgramMetric) Jaccard(a, b string) float64 {
	ra := []rune(a)
	rb := []rune(b)

	if len(ra) == 0 && len(rb) == 0 {
		return 1.0
	}

	if len(ra) == 0 || len(rb) == 0 {
		return 0.0
	}

	setA := ngramSet(ra, m.Size)
	setB := ngramSet(rb, m.Size)

	intersection := 0

	for gram := range setA {
		if setB[gram] {
			intersection++
		}
	}

	union := len(setA) + len(setB) - intersection

	if union == 0 {
		return 0.0
	}

	return float64(intersection) / float64(union)
}

// Dice returns the Sorensen-Dice similarity between a and b using set
// semantics: 2 * |intersection| / (|A| + |B|). Returns 1.0 when both strings
// are empty.
func (m *NgramMetric) Dice(a, b string) float64 {
	ra := []rune(a)
	rb := []rune(b)

	if len(ra) == 0 && len(rb) == 0 {
		return 1.0
	}

	if len(ra) == 0 || len(rb) == 0 {
		return 0.0
	}

	setA := ngramSet(ra, m.Size)
	setB := ngramSet(rb, m.Size)

	intersection := 0

	for gram := range setA {
		if setB[gram] {
			intersection++
		}
	}

	total := len(setA) + len(setB)
	if total == 0 {
		return 0.0
	}

	return 2.0 * float64(intersection) / float64(total)
}

// Overlap returns the overlap coefficient between a and b:
// |intersection| / min(|A|, |B|). Returns 1.0 when both strings are empty.
func (m *NgramMetric) Overlap(a, b string) float64 {
	ra := []rune(a)
	rb := []rune(b)

	if len(ra) == 0 && len(rb) == 0 {
		return 1.0
	}

	if len(ra) == 0 || len(rb) == 0 {
		return 0.0
	}

	setA := ngramSet(ra, m.Size)
	setB := ngramSet(rb, m.Size)

	intersection := 0

	for gram := range setA {
		if setB[gram] {
			intersection++
		}
	}

	minSize := min(len(setA), len(setB))
	if minSize == 0 {
		return 0.0
	}

	return float64(intersection) / float64(minSize)
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
