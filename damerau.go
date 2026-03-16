package strsim

// OSAMetric computes the Optimal String Alignment (restricted edit) distance
// between two strings. OSA extends the Levenshtein distance by also counting
// transpositions of two adjacent characters as a single operation, with the
// restriction that no substring is edited more than once.
//
// Set ASCIIOnly to true for faster processing of pure-ASCII strings.
type OSAMetric struct {
	// ASCIIOnly skips rune conversion for faster ASCII processing.
	// Produces incorrect results for multi-byte UTF-8 input.
	ASCIIOnly bool
}

// NewOSA returns a new OSAMetric instance.
func NewOSA() *OSAMetric {
	return &OSAMetric{}
}

// OSA returns the Optimal String Alignment distance between a and b.
func OSA(a, b string) int {
	return NewOSA().Distance(a, b)
}

// OSASimilarity returns the normalized similarity between a and b using the
// OSA distance as a value in [0, 1].
func OSASimilarity(a, b string) float64 {
	return NewOSA().Similarity(a, b)
}

// Distance returns the OSA distance between a and b. It uses a 3-row
// optimization for O(3 * min(m, n)) space.
func (m *OSAMetric) Distance(a, b string) int {
	if m.ASCIIOnly {
		return osaASCII(a, b)
	}

	ra := []rune(a)
	rb := []rune(b)

	lenA := len(ra)
	lenB := len(rb)

	if lenA == 0 {
		return lenB
	}
	if lenB == 0 {
		return lenA
	}

	// Ensure shorter string is in the inner loop.
	if lenA < lenB {
		ra, rb = rb, ra
		lenA, lenB = lenB, lenA
	}

	// Three rows: prevPrev (i-2), prev (i-1), curr (i).
	prevPrev := make([]int, lenB+1)
	prev := make([]int, lenB+1)
	curr := make([]int, lenB+1)

	for j := 0; j <= lenB; j++ {
		prev[j] = j
	}

	for i := 1; i <= lenA; i++ {
		curr[0] = i

		for j := 1; j <= lenB; j++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}

			curr[j] = min(
				prev[j]+1,      // deletion
				curr[j-1]+1,    // insertion
				prev[j-1]+cost, // substitution
			)

			// Transposition.
			if i > 1 && j > 1 && ra[i-1] == rb[j-2] && ra[i-2] == rb[j-1] {
				curr[j] = min(curr[j], prevPrev[j-2]+cost)
			}
		}

		prevPrev, prev, curr = prev, curr, prevPrev
	}

	return prev[lenB]
}

// Similarity returns the normalized similarity between a and b as a value
// in [0, 1]. Returns 1.0 when both strings are empty.
func (m *OSAMetric) Similarity(a, b string) float64 {
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

	return 1.0 - float64(d)/float64(maxLen)
}

// DamerauLevenshteinMetric computes the true Damerau-Levenshtein distance
// between two strings. Unlike OSA, it allows unrestricted transpositions,
// meaning substrings may be edited more than once.
//
// Set ASCIIOnly to true for faster processing of pure-ASCII strings.
type DamerauLevenshteinMetric struct {
	// ASCIIOnly skips rune conversion for faster ASCII processing.
	// Produces incorrect results for multi-byte UTF-8 input.
	ASCIIOnly bool
}

// NewDamerauLevenshtein returns a new DamerauLevenshteinMetric instance.
func NewDamerauLevenshtein() *DamerauLevenshteinMetric {
	return &DamerauLevenshteinMetric{}
}

// DamerauLevenshtein returns the true Damerau-Levenshtein distance between
// a and b.
func DamerauLevenshtein(a, b string) int {
	return NewDamerauLevenshtein().Distance(a, b)
}

// DamerauLevenshteinSimilarity returns the normalized similarity between a
// and b using the true Damerau-Levenshtein distance as a value in [0, 1].
func DamerauLevenshteinSimilarity(a, b string) float64 {
	return NewDamerauLevenshtein().Similarity(a, b)
}

// Distance returns the true Damerau-Levenshtein distance between a and b
// using the algorithm with a DA (last-row-seen) map. Time and space are
// O(m * n).
func (m *DamerauLevenshteinMetric) Distance(a, b string) int {
	if m.ASCIIOnly {
		return damerauLevenshteinASCII(a, b)
	}

	ra := []rune(a)
	rb := []rune(b)

	lenA := len(ra)
	lenB := len(rb)

	if lenA == 0 {
		return lenB
	}
	if lenB == 0 {
		return lenA
	}

	// da maps each rune to the last row where it was seen.
	da := make(map[rune]int)

	// d is (lenA+2) x (lenB+2) matrix. Indices are shifted by 1 so that
	// d[-1][-1] maps to d[0][0].
	d := make([][]int, lenA+2)
	for i := range d {
		d[i] = make([]int, lenB+2)
	}

	maxDist := lenA + lenB
	d[0][0] = maxDist

	for i := 0; i <= lenA; i++ {
		d[i+1][0] = maxDist
		d[i+1][1] = i
	}
	for j := 0; j <= lenB; j++ {
		d[0][j+1] = maxDist
		d[1][j+1] = j
	}

	for i := 1; i <= lenA; i++ {
		db := 0 // last column where ra[i-1] matched in rb

		for j := 1; j <= lenB; j++ {
			i1 := da[rb[j-1]] // last row where rb[j-1] appeared in ra
			j1 := db

			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
				db = j
			}

			d[i+1][j+1] = min(
				d[i][j]+cost, // substitution
				d[i+1][j]+1,  // insertion
				d[i][j+1]+1,  // deletion
				d[i1][j1]+(i-i1-1)+1+(j-j1-1), // transposition
			)
		}

		da[ra[i-1]] = i
	}

	return d[lenA+1][lenB+1]
}

// Similarity returns the normalized similarity between a and b as a value
// in [0, 1]. Returns 1.0 when both strings are empty.
func (m *DamerauLevenshteinMetric) Similarity(a, b string) float64 {
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

	return 1.0 - float64(d)/float64(maxLen)
}
