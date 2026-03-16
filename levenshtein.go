package strsim

// LevenshteinMetric computes the Levenshtein edit distance between two strings.
// Levenshtein distance counts the minimum number of single-character insertions,
// deletions, and substitutions needed to transform one string into the other.
//
// The metric supports configurable operation costs via InsertCost, DeleteCost,
// and ReplaceCost fields.
//
// Set ASCIIOnly to true for faster processing of pure-ASCII strings. This skips
// the rune conversion but produces incorrect results for multi-byte UTF-8 input.
type LevenshteinMetric struct {
	// InsertCost is the cost of inserting a character. Default: 1.
	InsertCost int
	// DeleteCost is the cost of deleting a character. Default: 1.
	DeleteCost int
	// ReplaceCost is the cost of replacing a character. Default: 1.
	ReplaceCost int
	// ASCIIOnly skips rune conversion for faster ASCII processing.
	// Produces incorrect results for multi-byte UTF-8 input.
	ASCIIOnly bool
}

// NewLevenshtein returns a new LevenshteinMetric with all operation costs set
// to 1 (standard Levenshtein distance).
func NewLevenshtein() *LevenshteinMetric {
	return &LevenshteinMetric{
		InsertCost:  1,
		DeleteCost:  1,
		ReplaceCost: 1,
	}
}

// Levenshtein returns the standard Levenshtein distance between a and b using
// unit costs for all operations.
func Levenshtein(a, b string) int {
	return NewLevenshtein().Distance(a, b)
}

// LevenshteinSimilarity returns the normalized similarity between a and b as
// a value in [0, 1] using the standard Levenshtein distance.
func LevenshteinSimilarity(a, b string) float64 {
	return NewLevenshtein().Similarity(a, b)
}

// Distance returns the Levenshtein distance between a and b using the
// configured operation costs. It uses O(min(m, n)) space.
func (m *LevenshteinMetric) Distance(a, b string) int {
	if m.ASCIIOnly {
		return levenshteinASCII(a, b, m.InsertCost, m.DeleteCost, m.ReplaceCost)
	}

	ra := []rune(a)
	rb := []rune(b)

	lenA := len(ra)
	lenB := len(rb)

	// When costs are symmetric we can swap for space optimization.
	// When costs differ, insert/delete semantics depend on direction, so we
	// must track the effective costs after a potential swap.
	insCost := m.InsertCost
	delCost := m.DeleteCost

	if lenA < lenB {
		ra, rb = rb, ra
		lenA, lenB = lenB, lenA
		insCost, delCost = delCost, insCost
	}

	// Handle trivial cases.
	if lenB == 0 {
		return lenA * delCost
	}

	// Use two rows of size lenB+1.
	prev := make([]int, lenB+1)
	curr := make([]int, lenB+1)

	for j := 0; j <= lenB; j++ {
		prev[j] = j * insCost
	}

	for i := 1; i <= lenA; i++ {
		curr[0] = i * delCost

		for j := 1; j <= lenB; j++ {
			cost := m.ReplaceCost
			if ra[i-1] == rb[j-1] {
				cost = 0
			}

			del := prev[j] + delCost
			ins := curr[j-1] + insCost
			rep := prev[j-1] + cost

			curr[j] = min(del, min(ins, rep))
		}

		prev, curr = curr, prev
	}

	return prev[lenB]
}

// Similarity returns the normalized similarity between a and b as a value
// in [0, 1]. It normalizes by dividing the distance by the maximum possible
// distance (max rune length * max operation cost). Returns 1.0 when both
// strings are empty.
func (m *LevenshteinMetric) Similarity(a, b string) float64 {
	var lenA, lenB int
	if m.ASCIIOnly {
		lenA = len(a)
		lenB = len(b)
	} else {
		lenA = runeLen(a)
		lenB = runeLen(b)
	}

	maxLen := max(lenA, lenB)
	if maxLen == 0 {
		return 1.0
	}

	d := m.Distance(a, b)
	maxCost := max(m.InsertCost, max(m.DeleteCost, m.ReplaceCost))

	return 1.0 - float64(d)/float64(maxLen*maxCost)
}
