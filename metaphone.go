package strsim

import "strings"

// defaultMetaphoneLength is the standard Metaphone code length.
const defaultMetaphoneLength = 4

// MetaphoneEncoder implements the original Metaphone phonetic algorithm by
// Lawrence Philips. It transforms an English word into a phonetic key that
// represents its approximate pronunciation.
//
// The rules follow Apache Commons Codec's Metaphone, except that an initial
// CH is always encoded as X ("Church" → XRX), as in Philips' original rules.
type MetaphoneEncoder struct {
	// MaxLength is the maximum length of the generated code. Values <= 0 use
	// the default of 4.
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

// Encode returns the Metaphone code for the given string. Latin diacritics are
// folded (é → E) and other non-letters are ignored. It returns an empty string
// if the input contains no letters; a single letter is its own code.
func (e *MetaphoneEncoder) Encode(s string) string {
	maxLen := e.MaxLength
	if maxLen <= 0 {
		maxLen = defaultMetaphoneLength
	}

	var lettersBuf [32]byte
	w := string(appendLetters(lettersBuf[:0], s, false))
	if len(w) <= 1 {
		return w
	}

	// Initial exceptions: silent first letter of AE, GN, KN, PN and WR;
	// WH → W; X → S.
	switch {
	case w[1] == 'N' && (w[0] == 'G' || w[0] == 'K' || w[0] == 'P'),
		w[0] == 'A' && w[1] == 'E',
		w[0] == 'W' && w[1] == 'R':
		w = w[1:]
	case w[0] == 'W' && w[1] == 'H':
		w = "W" + w[2:]
	case w[0] == 'X':
		w = "S" + w[1:]
	}

	n := len(w)
	var codeBuf [8]byte
	code := codeBuf[:0]

	for i := 0; i < n && len(code) < maxLen; i++ {
		ch := w[i]

		// Skip duplicate adjacent letters, except C.
		if ch != 'C' && i > 0 && w[i-1] == ch {
			continue
		}

		prev := byteAt(w, i-1)
		next := byteAt(w, i+1)

		switch ch {
		case 'A', 'E', 'I', 'O', 'U':
			// Vowels are only kept as the first letter.
			if i == 0 {
				code = append(code, ch)
			}

		case 'B':
			// Silent in a final MB ("dumb").
			if i < n-1 || prev != 'M' {
				code = append(code, 'B')
			}

		case 'C':
			switch {
			case prev == 'S' && isFrontVowel(next):
				// SCI, SCE, SCY: silent.
			case hasAt(w, i, "CIA"):
				code = append(code, 'X')
			case isFrontVowel(next):
				code = append(code, 'S')
			case prev == 'S' && next == 'H':
				code = append(code, 'K') // SCH
			case next == 'H':
				code = append(code, 'X')
			default:
				code = append(code, 'K')
			}

		case 'D':
			if next == 'G' && isFrontVowel(byteAt(w, i+2)) {
				code = append(code, 'J') // DGE, DGI, DGY
				i += 2
			} else {
				code = append(code, 'T')
			}

		case 'G':
			switch {
			case next == 'H' && !metaphoneIsVowel(byteAt(w, i+2)):
				// GH at the end or before a consonant: silent.
			case i > 0 && next == 'N':
				// GN, GNED: silent.
			case isFrontVowel(next):
				code = append(code, 'J')
			default:
				code = append(code, 'K')
			}

		case 'H':
			// Kept only before a vowel and when not part of CH, SH, PH, TH, GH.
			if metaphoneIsVowel(next) && strings.IndexByte("CSPTG", prev) < 0 {
				code = append(code, 'H')
			}

		case 'F', 'J', 'L', 'M', 'N', 'R':
			code = append(code, ch)

		case 'K':
			if prev != 'C' {
				code = append(code, 'K')
			}

		case 'P':
			if next == 'H' {
				code = append(code, 'F')
			} else {
				code = append(code, 'P')
			}

		case 'Q':
			code = append(code, 'K')

		case 'S':
			if next == 'H' || hasAt(w, i, "SIO", "SIA") {
				code = append(code, 'X')
			} else {
				code = append(code, 'S')
			}

		case 'T':
			switch {
			case hasAt(w, i, "TIA", "TIO"):
				code = append(code, 'X')
			case hasAt(w, i, "TCH"):
				// Silent before CH.
			case next == 'H':
				code = append(code, '0')
			default:
				code = append(code, 'T')
			}

		case 'V':
			code = append(code, 'F')

		case 'W', 'Y':
			// Kept only before a vowel.
			if metaphoneIsVowel(next) {
				code = append(code, ch)
			}

		case 'X':
			code = append(code, 'K', 'S')

		case 'Z':
			code = append(code, 'S')
		}
	}

	if len(code) > maxLen {
		code = code[:maxLen]
	}

	return string(code)
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

// isFrontVowel reports whether ch is E, I or Y, which soften C and G.
func isFrontVowel(ch byte) bool {
	return ch == 'E' || ch == 'I' || ch == 'Y'
}

// byteAt returns s[i], or 0 when i is out of range.
func byteAt(s string, i int) byte {
	if i < 0 || i >= len(s) {
		return 0
	}

	return s[i]
}

// hasAt reports whether any of subs occurs in s at position pos.
func hasAt(s string, pos int, subs ...string) bool {
	if pos < 0 {
		return false
	}

	for _, sub := range subs {
		if pos+len(sub) <= len(s) && s[pos:pos+len(sub)] == sub {
			return true
		}
	}

	return false
}
