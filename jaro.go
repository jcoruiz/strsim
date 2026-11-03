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
	return jaroRunes([]rune(a), []rune(b))
}

// jaroRunes computes Jaro similarity over rune slices.
func jaroRunes(ra, rb []rune) float64 {
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
	// Half the out-of-order matches, rounded down as in Winkler's strcmp95.
	t := float64(transpositions / 2)

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
//
// The zero value uses the defaults of NewJaroWinkler.
type JaroWinklerMetric struct {
	// BoostThreshold is the minimum Jaro score required to apply the prefix
	// bonus. Default: 0.7.
	BoostThreshold float64

	// PrefixSize is the maximum number of prefix characters considered for
	// the bonus. Each shared prefix character adds 10% of the remaining
	// distance to 1.0, so a prefix of 10 or more characters scores 1.0.
	// Default: 4.
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
	threshold, prefixSize := m.BoostThreshold, m.PrefixSize
	if threshold == 0 && prefixSize == 0 {
		threshold, prefixSize = 0.7, 4
	}

	if m.ASCIIOnly {
		jaro := jaroASCII(a, b)
		if jaro < threshold {
			return jaro
		}

		prefix := 0
		for prefix < min(len(a), len(b), prefixSize) && a[prefix] == b[prefix] {
			prefix++
		}

		return winklerBoost(jaro, prefix)
	}

	ra := []rune(a)
	rb := []rune(b)

	jaro := jaroRunes(ra, rb)
	if jaro < threshold {
		return jaro
	}

	prefix := 0
	for prefix < min(len(ra), len(rb), prefixSize) && ra[prefix] == rb[prefix] {
		prefix++
	}

	return winklerBoost(jaro, prefix)
}

// winklerBoost adds the Winkler prefix bonus to a Jaro score, capped so the
// result never exceeds 1.0.
func winklerBoost(jaro float64, prefix int) float64 {
	return jaro + min(float64(prefix)*0.1, 1.0)*(1.0-jaro)
}
