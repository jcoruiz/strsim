package strsim

// defaultSoundexLength is the standard Soundex code length.
const defaultSoundexLength = 4

// soundexTable maps uppercase ASCII letters to their Soundex digit.
// A,E,I,O,U,H,W,Y → '0'; B,F,P,V → '1'; C,G,J,K,Q,S,X,Z → '2';
// D,T → '3'; L → '4'; M,N → '5'; R → '6'.
var soundexTable = [26]byte{
	'0', // A
	'1', // B
	'2', // C
	'3', // D
	'0', // E
	'1', // F
	'2', // G
	'0', // H
	'0', // I
	'2', // J
	'2', // K
	'4', // L
	'5', // M
	'5', // N
	'0', // O
	'1', // P
	'2', // Q
	'6', // R
	'2', // S
	'3', // T
	'0', // U
	'1', // V
	'0', // W
	'2', // X
	'0', // Y
	'2', // Z
}

// SoundexEncoder implements the American Soundex phonetic algorithm.
// It maps a name to a four-character code consisting of one letter followed
// by three digits, enabling approximate matching of names that sound alike.
type SoundexEncoder struct {
	// MaxLength is the length of the generated Soundex code. Default is 4.
	MaxLength int
}

// NewSoundex returns a new SoundexEncoder with the default code length of 4.
func NewSoundex() *SoundexEncoder {
	return &SoundexEncoder{MaxLength: defaultSoundexLength}
}

// Soundex returns the American Soundex code for the given string using the
// default code length of 4.
func Soundex(s string) string {
	return NewSoundex().Encode(s)
}

// SoundexMatch reports whether two strings have the same Soundex code using
// the default code length.
func SoundexMatch(a, b string) bool {
	return NewSoundex().Match(a, b)
}

// Encode returns the Soundex code for the given string. It returns an empty
// string if the input contains no ASCII letters.
func (e *SoundexEncoder) Encode(s string) string {
	length := e.MaxLength
	if length <= 0 {
		length = defaultSoundexLength
	}

	// Extract ASCII letters only.
	letters := toASCIIAlpha(s)
	if len(letters) == 0 {
		return ""
	}

	// Build the code: first character (uppercase) + digits.
	buf := make([]byte, 0, length)
	first := toUpper(letters[0])
	buf = append(buf, first)

	prevCode := soundexTable[first-'A']

	for i := 1; i < len(letters) && len(buf) < length; i++ {
		ch := toUpper(letters[i])
		code := soundexTable[ch-'A']

		if code == '0' {
			// In standard Soundex, H and W are transparent: two consonants
			// with the same code separated by H or W are coded as one digit.
			// Only vowels (A,E,I,O,U,Y) reset the previous code so that two
			// consonants separated by a vowel are encoded separately.
			if ch != 'H' && ch != 'W' {
				prevCode = '0'
			}
			continue
		}

		if code != prevCode {
			buf = append(buf, code)
			prevCode = code
		}
	}

	// Pad with zeros.
	for len(buf) < length {
		buf = append(buf, '0')
	}

	return string(buf)
}

// Match reports whether two strings produce the same Soundex code.
func (e *SoundexEncoder) Match(a, b string) bool {
	return e.Encode(a) == e.Encode(b)
}

// toASCIIAlpha returns only the ASCII alphabetic characters from s as a
// byte slice, each character in its original case.
func toASCIIAlpha(s string) []byte {
	out := make([]byte, 0, len(s))

	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			out = append(out, byte(r))
		} else if r >= 'a' && r <= 'z' {
			out = append(out, byte(r))
		}
		// Skip non-ASCII letters: phonetic encoding only operates on
		// ASCII alphabetic characters.
	}

	return out
}

// toUpper converts an ASCII letter to uppercase.
func toUpper(b byte) byte {
	if b >= 'a' && b <= 'z' {
		return b - 32
	}

	return b
}
