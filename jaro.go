package strsim

// JaroMetric computes the Jaro similarity between two strings. Jaro similarity
// is based on the number of matching characters and transpositions.
//
// Set ASCIIOnly to true for faster processing of pure-ASCII strings.
type JaroMetric struct {
	// ASCIIOnly skips rune conversion for faster ASCII processing.
	// Produces incorrect results for multi-byte UTF-8 input.
	ASCIIOnly bool
}

// NewJaro returns a new JaroMetric instance.
func NewJaro() *JaroMetric {
	return &JaroMetric{}
}

// JaroSimilarity returns the Jaro similarity between a and b as a value in
// [0, 1], where 1.0 means identical. The algorithm considers characters as
// matching if they appear within a window of max(len(a), len(b))/2 - 1
// positions and counts transpositions among matched characters.
func JaroSimilarity(a, b string) float64 {
	ra := []rune(a)
	rb := []rune(b)

	la := len(ra)
	lb := len(rb)

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

	// Find matching characters.
	for i := 0; i < la; i++ {
		lo := max(0, i-window)
		hi := min(lb-1, i+window)

		for j := lo; j <= hi; j++ {
			if matchedB[j] || ra[i] != rb[j] {
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

	// Count transpositions.
	k := 0
	for i := 0; i < la; i++ {
		if !matchedA[i] {
			continue
		}

		for !matchedB[k] {
			k++
		}

		if ra[i] != rb[k] {
			transpositions++
		}

		k++
	}

	m := float64(matches)
	t := float64(transpositions) / 2.0

	return (m/float64(la) + m/float64(lb) + (m-t)/m) / 3.0
}

// Similarity returns the Jaro similarity between a and b.
func (m *JaroMetric) Similarity(a, b string) float64 {
	if m.ASCIIOnly {
		return jaroASCII(a, b)
	}

	return JaroSimilarity(a, b)
}

// JaroWinklerMetric computes the Jaro-Winkler similarity between two strings.
// Jaro-Winkler extends Jaro with a prefix bonus that increases the score when
// the strings share a common prefix.
//
// Set ASCIIOnly to true for faster processing of pure-ASCII strings.
type JaroWinklerMetric struct {
	// BoostThreshold is the minimum Jaro score required to apply the prefix
	// bonus. Default: 0.7.
	BoostThreshold float64

	// PrefixSize is the maximum number of prefix characters considered for
	// the bonus. Default: 4.
	PrefixSize int

	// ASCIIOnly skips rune conversion for faster ASCII processing.
	// Produces incorrect results for multi-byte UTF-8 input.
	ASCIIOnly bool
}

// NewJaroWinkler returns a new JaroWinklerMetric with default settings
// (BoostThreshold = 0.7, PrefixSize = 4).
func NewJaroWinkler() *JaroWinklerMetric {
	return &JaroWinklerMetric{
		BoostThreshold: 0.7,
		PrefixSize:     4,
	}
}

// JaroWinklerSimilarity returns the Jaro-Winkler similarity between a and b
// using default settings (BoostThreshold = 0.7, PrefixSize = 4).
func JaroWinklerSimilarity(a, b string) float64 {
	return NewJaroWinkler().Similarity(a, b)
}

// Similarity returns the Jaro-Winkler similarity between a and b as a value
// in [0, 1], where 1.0 means identical.
func (m *JaroWinklerMetric) Similarity(a, b string) float64 {
	if m.ASCIIOnly {
		la := len(a)
		lb := len(b)

		if la == 0 && lb == 0 {
			return 1.0
		}

		jaro := jaroASCII(a, b)

		if jaro < m.BoostThreshold {
			return jaro
		}

		prefixLen := min(min(la, lb), m.PrefixSize)
		prefix := 0

		for i := 0; i < prefixLen; i++ {
			if a[i] != b[i] {
				break
			}

			prefix++
		}

		return jaro + float64(prefix)*0.1*(1.0-jaro)
	}

	ra := []rune(a)
	rb := []rune(b)

	if len(ra) == 0 && len(rb) == 0 {
		return 1.0
	}

	jaro := JaroSimilarity(a, b)

	if jaro < m.BoostThreshold {
		return jaro
	}

	prefixLen := min(min(len(ra), len(rb)), m.PrefixSize)
	prefix := 0

	for i := 0; i < prefixLen; i++ {
		if ra[i] != rb[i] {
			break
		}

		prefix++
	}

	return jaro + float64(prefix)*0.1*(1.0-jaro)
}
