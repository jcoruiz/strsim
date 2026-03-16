package strsim

import "strings"

// defaultMetaphoneLength is the standard Metaphone code length.
const defaultMetaphoneLength = 4

// MetaphoneEncoder implements the original Metaphone phonetic algorithm by
// Lawrence Philips. It transforms an English word into a phonetic key that
// represents its approximate pronunciation.
type MetaphoneEncoder struct {
	// MaxLength is the maximum length of the generated code. Default is 4.
	MaxLength int
}

// NewMetaphone returns a new MetaphoneEncoder with the default code length of 4.
func NewMetaphone() *MetaphoneEncoder {
	return &MetaphoneEncoder{MaxLength: defaultMetaphoneLength}
}

// Metaphone returns the Metaphone code for the given string using the default
// code length of 4.
func Metaphone(s string) string {
	return NewMetaphone().Encode(s)
}

// MetaphoneMatch reports whether two strings have the same Metaphone code.
func MetaphoneMatch(a, b string) bool {
	return NewMetaphone().Match(a, b)
}

// Encode returns the Metaphone code for the given string. It returns an empty
// string if the input contains no ASCII letters.
func (e *MetaphoneEncoder) Encode(s string) string {
	maxLen := e.MaxLength
	if maxLen <= 0 {
		maxLen = defaultMetaphoneLength
	}

	// Normalize: extract ASCII alpha, convert to uppercase.
	letters := toASCIIAlpha(s)
	if len(letters) == 0 {
		return ""
	}

	word := strings.ToUpper(string(letters))
	n := len(word)

	// Handle initial silent letter pairs.
	switch {
	case n >= 2 && word[:2] == "AE":
		word = word[1:]
		n--
	case n >= 2 && word[:2] == "GN":
		word = word[1:]
		n--
	case n >= 2 && word[:2] == "KN":
		word = word[1:]
		n--
	case n >= 2 && word[:2] == "PN":
		word = word[1:]
		n--
	case n >= 2 && word[:2] == "WR":
		word = word[1:]
		n--
	case n >= 1 && word[0] == 'X':
		word = "S" + word[1:]
	}

	var buf strings.Builder
	buf.Grow(maxLen)

	i := 0
	for i < n && buf.Len() < maxLen {
		ch := word[i]

		// Skip duplicate adjacent letters, except C.
		if i > 0 && ch == word[i-1] && ch != 'C' {
			i++
			continue
		}

		switch ch {
		case 'A', 'E', 'I', 'O', 'U':
			// Vowels are only kept if they are the first character.
			if i == 0 {
				buf.WriteByte(ch)
			}
			i++

		case 'B':
			// B is silent if after M at end of word.
			if !(i > 0 && word[i-1] == 'M' && i == n-1) {
				buf.WriteByte('B')
			}
			i++

		case 'C':
			if i+1 < n && word[i+1] == 'H' {
				buf.WriteByte('X')
				i += 2
			} else if i+1 < n && (word[i+1] == 'I' || word[i+1] == 'E' || word[i+1] == 'Y') {
				if i > 0 && word[i-1] == 'S' {
					// SC before I, E, Y → skip the C.
					i++
					continue
				}
				buf.WriteByte('S')
				i++
			} else {
				buf.WriteByte('K')
				i++
			}

		case 'D':
			if i+1 < n && word[i+1] == 'G' {
				if i+2 < n && (word[i+2] == 'I' || word[i+2] == 'E' || word[i+2] == 'Y') {
					buf.WriteByte('J')
					i += 3
				} else {
					buf.WriteByte('T')
					i++
				}
			} else {
				buf.WriteByte('T')
				i++
			}

		case 'F':
			buf.WriteByte('F')
			i++

		case 'G':
			silent := false

			// G before H: silent if H is not at end and char after H is not a vowel,
			// or GH at end of word.
			if i+1 < n && word[i+1] == 'H' {
				if i+2 >= n {
					// GH at end of word: G is silent, skip both.
					silent = true
					i += 2
				} else if !metaphoneIsVowel(word[i+2]) {
					// GH followed by non-vowel: both silent.
					silent = true
					i += 2
				} else {
					// GH followed by vowel: G → K.
					buf.WriteByte('K')
					i += 2
					silent = true
				}
			}

			// GN at end: G is silent.
			if !silent && i+1 < n && word[i+1] == 'N' && i+2 >= n {
				silent = true
				i += 2
			}

			if silent {
				continue
			}

			if i+1 < n && (word[i+1] == 'I' || word[i+1] == 'E' || word[i+1] == 'Y') {
				if i > 0 && word[i-1] == 'G' {
					i++
					continue
				}
				buf.WriteByte('J')
			} else {
				buf.WriteByte('K')
			}
			i++

		case 'H':
			// H is encoded only before a vowel and if the previous char is not a vowel.
			if i+1 < n && metaphoneIsVowel(word[i+1]) {
				if i == 0 || !metaphoneIsVowel(word[i-1]) {
					buf.WriteByte('H')
				}
			}
			i++

		case 'J':
			buf.WriteByte('J')
			i++

		case 'K':
			if i > 0 && word[i-1] == 'C' {
				i++
				continue
			}
			buf.WriteByte('K')
			i++

		case 'L':
			buf.WriteByte('L')
			i++

		case 'M':
			buf.WriteByte('M')
			i++

		case 'N':
			buf.WriteByte('N')
			i++

		case 'P':
			if i+1 < n && word[i+1] == 'H' {
				buf.WriteByte('F')
				i += 2
			} else {
				buf.WriteByte('P')
				i++
			}

		case 'Q':
			buf.WriteByte('K')
			i++

		case 'R':
			buf.WriteByte('R')
			i++

		case 'S':
			if i+1 < n && word[i+1] == 'H' {
				buf.WriteByte('X')
				i += 2
			} else if i+1 < n && word[i+1] == 'I' && i+2 < n && (word[i+2] == 'O' || word[i+2] == 'A') {
				buf.WriteByte('X')
				i += 3
			} else {
				buf.WriteByte('S')
				i++
			}

		case 'T':
			if i+1 < n && word[i+1] == 'H' {
				buf.WriteByte('0')
				i += 2
			} else if i+1 < n && word[i+1] == 'I' && i+2 < n && (word[i+2] == 'A' || word[i+2] == 'O') {
				buf.WriteByte('X')
				i += 3
			} else {
				buf.WriteByte('T')
				i++
			}

		case 'V':
			buf.WriteByte('F')
			i++

		case 'W':
			// W is kept only before a vowel.
			if i+1 < n && metaphoneIsVowel(word[i+1]) {
				buf.WriteByte('W')
			}
			i++

		case 'X':
			buf.WriteByte('K')
			if buf.Len() < maxLen {
				buf.WriteByte('S')
			}
			i++

		case 'Y':
			// Y is kept only before a vowel.
			if i+1 < n && metaphoneIsVowel(word[i+1]) {
				buf.WriteByte('Y')
			}
			i++

		case 'Z':
			buf.WriteByte('S')
			i++

		default:
			i++
		}
	}

	return buf.String()
}

// Match reports whether two strings produce the same Metaphone code.
func (e *MetaphoneEncoder) Match(a, b string) bool {
	return e.Encode(a) == e.Encode(b)
}

// metaphoneIsVowel reports whether ch is an ASCII vowel.
func metaphoneIsVowel(ch byte) bool {
	switch ch {
	case 'A', 'E', 'I', 'O', 'U':
		return true
	}

	return false
}
