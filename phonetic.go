package strsim

// upperLetters is sliced to turn an ASCII letter into a string without
// allocating.
const upperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// appendLetters appends the ASCII letters of s to out in uppercase.
// Everything else is dropped, except that with keepSpaces a
// run of non-letters between two words becomes a single space (Double
// Metaphone has rules for "VAN ", "SAN " and "MAC GREGOR"). Apostrophes never
// separate words, so "O'Brien" stays one word.
//
// Callers pass a stack buffer as out to avoid allocating for short names.
func appendLetters(out []byte, s string, keepSpaces bool) []byte {
	gap := false

	for _, r := range s {
		var fold string

		switch {
		case r >= 'A' && r <= 'Z':
			fold = upperLetters[r-'A' : r-'A'+1]
		case r >= 'a' && r <= 'z':
			fold = upperLetters[r-'a' : r-'a'+1]
		case r == '\'' || r == '’':
			continue
		}

		if fold == "" {
			gap = keepSpaces
			continue
		}

		if gap && len(out) > 0 {
			out = append(out, ' ')
		}
		gap = false
		out = append(out, fold...)
	}

	return out
}
