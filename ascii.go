package strsim

// This file contains byte-level implementations of edit distance and similarity
// algorithms. These are used when ASCIIOnly is set to true on a metric struct.
//
// Byte-level operations skip the []rune conversion, which is significantly
// faster for ASCII-only text. However, they produce INCORRECT results for
// strings containing multi-byte UTF-8 characters (any non-ASCII text).
//
// Users should only enable ASCIIOnly when they are certain their input is
// pure ASCII (e.g., English identifiers, codes, URLs).

// levenshteinASCII computes Levenshtein distance operating on bytes.
func levenshteinASCII(a, b string, insCost, delCost, repCost int) int {
	lenA := len(a)
	lenB := len(b)

	ic := insCost
	dc := delCost

	if lenA < lenB {
		a, b = b, a
		lenA, lenB = lenB, lenA
		ic, dc = dc, ic
	}

	if lenB == 0 {
		return lenA * dc
	}

	prev := make([]int, lenB+1)
	curr := make([]int, lenB+1)

	for j := 0; j <= lenB; j++ {
		prev[j] = j * ic
	}

	for i := 1; i <= lenA; i++ {
		curr[0] = i * dc

		for j := 1; j <= lenB; j++ {
			cost := repCost
			if a[i-1] == b[j-1] {
				cost = 0
			}

			del := prev[j] + dc
			ins := curr[j-1] + ic
			rep := prev[j-1] + cost

			curr[j] = min(del, min(ins, rep))
		}

		prev, curr = curr, prev
	}

	return prev[lenB]
}

// osaASCII computes OSA distance operating on bytes.
func osaASCII(a, b string) int {
	lenA := len(a)
	lenB := len(b)

	if lenA == 0 {
		return lenB
	}
	if lenB == 0 {
		return lenA
	}

	if lenA < lenB {
		a, b = b, a
		lenA, lenB = lenB, lenA
	}

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
			if a[i-1] == b[j-1] {
				cost = 0
			}

			curr[j] = min(
				prev[j]+1,
				curr[j-1]+1,
				prev[j-1]+cost,
			)

			if i > 1 && j > 1 && a[i-1] == b[j-2] && a[i-2] == b[j-1] {
				curr[j] = min(curr[j], prevPrev[j-2]+cost)
			}
		}

		prevPrev, prev, curr = prev, curr, prevPrev
	}

	return prev[lenB]
}

// damerauLevenshteinASCII computes true Damerau-Levenshtein distance on bytes.
func damerauLevenshteinASCII(a, b string) int {
	lenA := len(a)
	lenB := len(b)

	if lenA == 0 {
		return lenB
	}
	if lenB == 0 {
		return lenA
	}

	// da maps each byte to the last row where it was seen.
	// Using a fixed 256-byte array instead of a map for speed.
	var da [256]int

	// The (lenA+2) x (lenB+2) matrix lives in one slice; at(i, j) is its
	// index.
	cols := lenB + 2
	d := make([]int, (lenA+2)*cols)
	at := func(i, j int) int { return i*cols + j }

	maxDist := lenA + lenB
	d[0] = maxDist

	for i := 0; i <= lenA; i++ {
		d[at(i+1, 0)] = maxDist
		d[at(i+1, 1)] = i
	}
	for j := 0; j <= lenB; j++ {
		d[at(0, j+1)] = maxDist
		d[at(1, j+1)] = j
	}

	for i := 1; i <= lenA; i++ {
		db := 0

		for j := 1; j <= lenB; j++ {
			i1 := da[b[j-1]]
			j1 := db

			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
				db = j
			}

			d[at(i+1, j+1)] = min(
				d[at(i, j)]+cost,
				d[at(i+1, j)]+1,
				d[at(i, j+1)]+1,
				d[at(i1, j1)]+(i-i1-1)+1+(j-j1-1),
			)
		}

		da[a[i-1]] = i
	}

	return d[at(lenA+1, lenB+1)]
}

// lcsASCII computes LCS length operating on bytes.
func lcsASCII(a, b string) int {
	lenA := len(a)
	lenB := len(b)

	if lenA < lenB {
		a, b = b, a
		lenA, lenB = lenB, lenA
	}

	if lenB == 0 {
		return 0
	}

	prev := make([]int, lenB+1)
	curr := make([]int, lenB+1)

	for i := 1; i <= lenA; i++ {
		for j := 1; j <= lenB; j++ {
			if a[i-1] == b[j-1] {
				curr[j] = prev[j-1] + 1
			} else {
				curr[j] = max(prev[j], curr[j-1])
			}
		}

		prev, curr = curr, prev

		for j := range curr {
			curr[j] = 0
		}
	}

	return prev[lenB]
}

// hammingASCII computes Hamming distance operating on bytes.
func hammingASCII(a, b string) int {
	d := 0
	for i := range len(a) {
		if a[i] != b[i] {
			d++
		}
	}

	return d
}

// jaroASCII computes Jaro similarity operating on bytes.
func jaroASCII(a, b string) float64 {
	la := len(a)
	lb := len(b)

	if la == 0 && lb == 0 {
		return 1.0
	}

	if la == 0 || lb == 0 {
		return 0.0
	}

	window := max(la, lb)/2 - 1
	if window < 0 {
		window = 0
	}

	matchedA := make([]bool, la)
	matchedB := make([]bool, lb)

	matches := 0
	transpositions := 0

	for i := 0; i < la; i++ {
		lo := max(0, i-window)
		hi := min(lb-1, i+window)

		for j := lo; j <= hi; j++ {
			if matchedB[j] || a[i] != b[j] {
				continue
			}

			matchedA[i] = true
			matchedB[j] = true
			matches++

			break
		}
	}

	if matches == 0 {
		return 0.0
	}

	k := 0
	for i := 0; i < la; i++ {
		if !matchedA[i] {
			continue
		}

		for !matchedB[k] {
			k++
		}

		if a[i] != b[k] {
			transpositions++
		}

		k++
	}

	m := float64(matches)
	// Half the out-of-order matches, rounded down as in Winkler's strcmp95.
	t := float64(transpositions / 2)

	return (m/float64(la) + m/float64(lb) + (m-t)/m) / 3.0
}
