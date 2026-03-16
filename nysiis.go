package strsim

import "strings"

// defaultNYSIISLength is the standard NYSIIS code length.
const defaultNYSIISLength = 6

// NYSIISEncoder implements the New York State Identification and Intelligence
// System (NYSIIS) phonetic algorithm. It produces a code that groups
// similar-sounding names together. The algorithm handles common English name
// patterns and is particularly effective for American names.
type NYSIISEncoder struct {
	// MaxLength is the maximum length of the generated code. Default is 6.
	// Set to 0 for unlimited length.
	MaxLength int
}

// NewNYSIIS returns a new NYSIISEncoder with the default code length of 6.
func NewNYSIIS() *NYSIISEncoder {
	return &NYSIISEncoder{MaxLength: defaultNYSIISLength}
}

// NYSIIS returns the NYSIIS code for the given string using the default code
// length of 6.
func NYSIIS(s string) string {
	return NewNYSIIS().Encode(s)
}

// NYSIISMatch reports whether two strings produce the same NYSIIS code.
func NYSIISMatch(a, b string) bool {
	return NewNYSIIS().Match(a, b)
}

// Encode returns the NYSIIS code for the given string. It returns an empty
// string if the input contains no ASCII letters.
func (e *NYSIISEncoder) Encode(s string) string {
	maxLen := e.MaxLength
	if maxLen < 0 {
		maxLen = defaultNYSIISLength
	}

	// Extract ASCII alpha, convert to uppercase.
	letters := toASCIIAlpha(s)
	if len(letters) == 0 {
		return ""
	}

	word := strings.ToUpper(string(letters))

	// Step 1: Handle initial prefixes.
	word = nysiisHandlePrefixes(word)

	// Step 2: Handle final suffixes.
	word = nysiisHandleSuffixes(word)

	if len(word) == 0 {
		return ""
	}

	// Step 3: First character of key = first character of name.
	key := make([]byte, 0, len(word))
	key = append(key, word[0])

	// Step 4: Process remaining characters with inline deduplication.
	// Each translated character is appended only if different from the last
	// character in the key.
	n := len(word)
	i := 1

	for i < n {
		ch := word[i]
		var translated byte

		// EV → AF.
		if ch == 'E' && i+1 < n && word[i+1] == 'V' {
			nysiisAppendUniq(&key, 'A')
			nysiisAppendUniq(&key, 'F')
			i += 2
			continue
		}

		// Vowels → A.
		if nysiisIsVowel(ch) {
			translated = 'A'
			nysiisAppendUniq(&key, translated)
			i++
			continue
		}

		switch ch {
		case 'Q':
			translated = 'G'
		case 'Z':
			translated = 'S'
		case 'M':
			translated = 'N'
		case 'K':
			if i+1 < n && word[i+1] == 'N' {
				translated = 'N'
			} else {
				translated = 'C'
			}
		case 'S':
			if i+1 < n && word[i+1] == 'C' && i+2 < n && word[i+2] == 'H' {
				nysiisAppendUniq(&key, 'S')
				nysiisAppendUniq(&key, 'S')
				nysiisAppendUniq(&key, 'S')
				i += 3
				continue
			}
			translated = 'S'
		case 'P':
			if i+1 < n && word[i+1] == 'H' {
				nysiisAppendUniq(&key, 'F')
				nysiisAppendUniq(&key, 'F')
				i += 2
				continue
			}
			translated = 'P'
		case 'H':
			// H: if previous or next char is not a vowel, use the previous
			// key character instead of H.
			prevCh := key[len(key)-1]
			nextIsVowel := i+1 < n && nysiisIsVowel(word[i+1])
			prevIsVowel := nysiisIsVowel(prevCh)

			if !prevIsVowel || !nextIsVowel {
				translated = prevCh
			} else {
				translated = 'H'
			}
		case 'W':
			// W: if previous key char is a vowel, use it instead of W.
			prevCh := key[len(key)-1]
			if nysiisIsVowel(prevCh) {
				translated = prevCh
			} else {
				translated = 'W'
			}
		default:
			translated = ch
		}

		nysiisAppendUniq(&key, translated)
		i++
	}

	result := string(key)

	// Step 5: Remove trailing S (unless it's the only character).
	if len(result) > 1 && result[len(result)-1] == 'S' {
		result = result[:len(result)-1]
	}

	// Step 6: Translate AY at end → Y.
	if len(result) >= 2 && result[len(result)-2:] == "AY" {
		result = result[:len(result)-2] + "Y"
	}

	// Step 7: Remove trailing A.
	if len(result) > 1 && result[len(result)-1] == 'A' {
		result = result[:len(result)-1]
	}

	// Step 8: Truncate to MaxLength.
	if maxLen > 0 && len(result) > maxLen {
		result = result[:maxLen]
	}

	return result
}

// Match reports whether two strings produce the same NYSIIS code.
func (e *NYSIISEncoder) Match(a, b string) bool {
	return e.Encode(a) == e.Encode(b)
}

// nysiisAppendUniq appends ch to key only if it differs from the last character.
func nysiisAppendUniq(key *[]byte, ch byte) {
	if len(*key) == 0 || (*key)[len(*key)-1] != ch {
		*key = append(*key, ch)
	}
}

// nysiisHandlePrefixes translates initial prefixes.
func nysiisHandlePrefixes(s string) string {
	switch {
	case strings.HasPrefix(s, "MAC"):
		return "MCC" + s[3:]
	case strings.HasPrefix(s, "KN"):
		return "NN" + s[2:]
	case strings.HasPrefix(s, "K"):
		return "C" + s[1:]
	case strings.HasPrefix(s, "PH"):
		return "FF" + s[2:]
	case strings.HasPrefix(s, "PF"):
		return "FF" + s[2:]
	case strings.HasPrefix(s, "SCH"):
		return "SSS" + s[3:]
	}

	return s
}

// nysiisHandleSuffixes translates final suffixes.
func nysiisHandleSuffixes(s string) string {
	switch {
	case strings.HasSuffix(s, "EE") || strings.HasSuffix(s, "IE"):
		return s[:len(s)-2] + "Y"
	case strings.HasSuffix(s, "DT") || strings.HasSuffix(s, "RT") ||
		strings.HasSuffix(s, "RD") || strings.HasSuffix(s, "NT") ||
		strings.HasSuffix(s, "ND"):
		return s[:len(s)-2] + "D"
	}

	return s
}

// nysiisIsVowel reports whether ch is a vowel in NYSIIS context.
func nysiisIsVowel(ch byte) bool {
	switch ch {
	case 'A', 'E', 'I', 'O', 'U':
		return true
	}

	return false
}
