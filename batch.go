package strsim

import (
	"cmp"
	"slices"
)

// FindBestMatch returns the candidate with the highest similarity to query
// according to the given Metric. Ties go to the earliest candidate. If
// candidates is empty, it returns a Match with Index -1 and Similarity 0.
func FindBestMatch(query string, candidates []string, m Metric) Match {
	if len(candidates) == 0 {
		return Match{Index: -1}
	}

	best := Match{Index: 0, Value: candidates[0], Similarity: m.Similarity(query, candidates[0])}

	for i := 1; i < len(candidates); i++ {
		s := m.Similarity(query, candidates[i])
		if s > best.Similarity {
			best = Match{Index: i, Value: candidates[i], Similarity: s}
		}
	}

	return best
}

// FindTopN returns the top n candidates with the highest similarity to query,
// sorted by similarity descending; ties keep the order of candidates. If n
// exceeds the number of candidates, all candidates are returned.
func FindTopN(query string, candidates []string, n int, m Metric) []Match {
	if len(candidates) == 0 || n <= 0 {
		return nil
	}

	matches := make([]Match, len(candidates))
	for i, c := range candidates {
		matches[i] = Match{Index: i, Value: c, Similarity: m.Similarity(query, c)}
	}

	slices.SortStableFunc(matches, compareMatchDesc)

	if n > len(matches) {
		n = len(matches)
	}

	return matches[:n]
}

// FindAboveThreshold returns all candidates whose similarity to query meets or
// exceeds the given threshold, sorted by similarity descending; ties keep the
// order of candidates.
func FindAboveThreshold(query string, candidates []string, threshold float64, m Metric) []Match {
	if len(candidates) == 0 {
		return nil
	}

	var matches []Match

	for i, c := range candidates {
		s := m.Similarity(query, c)
		if s >= threshold {
			matches = append(matches, Match{Index: i, Value: c, Similarity: s})
		}
	}

	slices.SortStableFunc(matches, compareMatchDesc)

	return matches
}

// compareMatchDesc sorts matches by similarity descending.
func compareMatchDesc(a, b Match) int {
	return cmp.Compare(b.Similarity, a.Similarity)
}
