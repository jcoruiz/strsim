package strsim

// LCSMetric computes the Longest Common Subsequence (LCS) between two strings.
// The LCS distance is defined as len(a) + len(b) - 2*LCS(a, b), representing
// the minimum number of characters that must be deleted from both strings to
// make them equal.
//
// Set ASCIIOnly to true for faster processing of pure-ASCII strings.
type LCSMetric struct {
	// ASCIIOnly skips rune conversion for faster ASCII processing.
	// Produces incorrect results for multi-byte UTF-8 input.
	ASCIIOnly bool
}

// NewLCS returns a new LCSMetric instance.
func NewLCS() *LCSMetric {
	return &LCSMetric{}
}

// LCS returns the length of the longest common subsequence of a and b. It
// uses O(min(m, n)) space with a two-row optimization.
func LCS(a, b string) int {
	ra := []rune(a)
	rb := []rune(b)

	lenA := len(ra)
	lenB := len(rb)

	// Ensure shorter string is in the inner loop.
	if lenA < lenB {
		ra, rb = rb, ra
		lenA, lenB = lenB, lenA
	}

	if lenB == 0 {
		return 0
	}

	prev := make([]int, lenB+1)
	curr := make([]int, lenB+1)

	for i := 1; i <= lenA; i++ {
		for j := 1; j <= lenB; j++ {
			if ra[i-1] == rb[j-1] {
				curr[j] = prev[j-1] + 1
			} else {
				curr[j] = max(prev[j], curr[j-1])
			}
		}

		prev, curr = curr, prev

		// Clear curr for the next row.
		for j := range curr {
			curr[j] = 0
		}
	}

	return prev[lenB]
}

// LCSDistance returns the LCS distance between a and b, defined as
// runeLen(a) + runeLen(b) - 2*LCS(a, b).
func LCSDistance(a, b string) int {
	lenA := runeLen(a)
	lenB := runeLen(b)
	lcsLen := LCS(a, b)

	return lenA + lenB - 2*lcsLen
}

// LCSSimilarity returns the normalized similarity between a and b based on
// the LCS length as a value in [0, 1]. Returns 1.0 when both strings are
// empty.
func LCSSimilarity(a, b string) float64 {
	lenA := runeLen(a)
	lenB := runeLen(b)

	maxLen := max(lenA, lenB)
	if maxLen == 0 {
		return 1.0
	}

	lcsLen := LCS(a, b)

	return float64(lcsLen) / float64(maxLen)
}

// Distance returns the LCS distance between a and b.
func (m *LCSMetric) Distance(a, b string) int {
	if m.ASCIIOnly {
		lcsLen := lcsASCII(a, b)
		return len(a) + len(b) - 2*lcsLen
	}

	return LCSDistance(a, b)
}

// Similarity returns the normalized similarity between a and b as a value
// in [0, 1]. Returns 1.0 when both strings are empty.
func (m *LCSMetric) Similarity(a, b string) float64 {
	if m.ASCIIOnly {
		maxLen := max(len(a), len(b))
		if maxLen == 0 {
			return 1.0
		}
		lcsLen := lcsASCII(a, b)
		return float64(lcsLen) / float64(maxLen)
	}

	return LCSSimilarity(a, b)
}
