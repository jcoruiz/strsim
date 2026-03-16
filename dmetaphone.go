package strsim

import "strings"

// defaultDMMaxLength is the standard Double Metaphone code length.
const defaultDMMaxLength = 4

// DoubleMetaphoneEncoder implements the Double Metaphone algorithm by Lawrence
// Philips. It produces two phonetic codes (primary and alternate) that capture
// different possible pronunciations, accounting for non-English name origins
// including Germanic, Slavic, Celtic, Greek, Italian, Spanish, and Chinese.
type DoubleMetaphoneEncoder struct {
	// MaxLength is the maximum length of each generated code. Default is 4.
	MaxLength int
}

// NewDoubleMetaphone returns a new DoubleMetaphoneEncoder with the default
// code length of 4.
func NewDoubleMetaphone() *DoubleMetaphoneEncoder {
	return &DoubleMetaphoneEncoder{MaxLength: defaultDMMaxLength}
}

// DoubleMetaphone returns the primary and alternate Double Metaphone codes for
// the given string using the default code length of 4.
func DoubleMetaphone(s string) (primary, alternate string) {
	return NewDoubleMetaphone().Encode(s)
}

// DoubleMetaphoneMatch reports whether two strings share at least one common
// Double Metaphone code (primary or alternate).
func DoubleMetaphoneMatch(a, b string) bool {
	return NewDoubleMetaphone().Match(a, b)
}

// Encode returns the primary and alternate Double Metaphone codes for the
// given string. Both codes are empty if the input contains no ASCII letters.
func (e *DoubleMetaphoneEncoder) Encode(s string) (primary, alternate string) {
	maxLen := e.MaxLength
	if maxLen <= 0 {
		maxLen = defaultDMMaxLength
	}

	letters := toASCIIAlpha(s)
	if len(letters) == 0 {
		return "", ""
	}

	word := strings.ToUpper(string(letters))

	// Pad the word to avoid bounds checking.
	padded := word + "     "
	n := len(word)

	var pri, alt strings.Builder
	pri.Grow(maxLen + 2)
	alt.Grow(maxLen + 2)

	slavoGermanic := dmIsSlavoGermanic(word)
	pos := 0

	// Handle initial silent/special pairs.
	if dmStringAt(padded, 0, "GN", "KN", "PN", "AE", "WR") {
		pos++
	}

	// Initial X → S.
	if padded[0] == 'X' {
		pri.WriteByte('S')
		alt.WriteByte('S')
		pos++
	}

	for pos < n && pri.Len() < maxLen && alt.Len() < maxLen {
		ch := padded[pos]

		switch ch {
		case 'A', 'E', 'I', 'O', 'U':
			if pos == 0 {
				pri.WriteByte('A')
				alt.WriteByte('A')
			}
			pos++

		case 'B':
			pri.WriteByte('P')
			alt.WriteByte('P')
			if padded[pos+1] == 'B' {
				pos += 2
			} else {
				pos++
			}

		case 'C':
			pos = dmHandleC(padded, pos, n, &pri, &alt, slavoGermanic)

		case 'D':
			pos = dmHandleD(padded, pos, n, &pri, &alt)

		case 'F':
			pri.WriteByte('F')
			alt.WriteByte('F')
			if padded[pos+1] == 'F' {
				pos += 2
			} else {
				pos++
			}

		case 'G':
			pos = dmHandleG(padded, pos, n, &pri, &alt, slavoGermanic)

		case 'H':
			// H: encode only if before a vowel and not after a vowel.
			if dmIsVowel(padded[pos+1]) && (pos == 0 || !dmIsVowel(padded[pos-1])) {
				pri.WriteByte('H')
				alt.WriteByte('H')
				pos += 2
			} else {
				pos++
			}

		case 'J':
			pos = dmHandleJ(padded, pos, n, word, &pri, &alt, slavoGermanic)

		case 'K':
			pri.WriteByte('K')
			alt.WriteByte('K')
			if padded[pos+1] == 'K' {
				pos += 2
			} else {
				pos++
			}

		case 'L':
			pos = dmHandleL(padded, pos, n, &pri, &alt)

		case 'M':
			pri.WriteByte('M')
			alt.WriteByte('M')
			// Handle UMB at end (e.g., DUMB, THUMB).
			if padded[pos+1] == 'M' {
				pos += 2
			} else {
				pos++
			}

		case 'N':
			pri.WriteByte('N')
			alt.WriteByte('N')
			if padded[pos+1] == 'N' {
				pos += 2
			} else {
				pos++
			}

		case 'P':
			if padded[pos+1] == 'H' {
				pri.WriteByte('F')
				alt.WriteByte('F')
				pos += 2
			} else {
				pri.WriteByte('P')
				alt.WriteByte('P')
				if padded[pos+1] == 'P' {
					pos += 2
				} else {
					pos++
				}
			}

		case 'Q':
			pri.WriteByte('K')
			alt.WriteByte('K')
			if padded[pos+1] == 'Q' {
				pos += 2
			} else {
				pos++
			}

		case 'R':
			// French final R: silent in alternate.
			if pos == n-1 && !slavoGermanic &&
				pos >= 2 && padded[pos-2] == 'I' && padded[pos-1] == 'E' &&
				!dmStringAt(padded, pos-4, "ME", "MA") {
				alt.WriteByte('R')
			} else {
				pri.WriteByte('R')
				alt.WriteByte('R')
			}
			if padded[pos+1] == 'R' {
				pos += 2
			} else {
				pos++
			}

		case 'S':
			pos = dmHandleS(padded, pos, n, &pri, &alt, slavoGermanic)

		case 'T':
			pos = dmHandleT(padded, pos, n, &pri, &alt)

		case 'V':
			pri.WriteByte('F')
			alt.WriteByte('F')
			if padded[pos+1] == 'V' {
				pos += 2
			} else {
				pos++
			}

		case 'W':
			pos = dmHandleW(padded, pos, n, &pri, &alt)

		case 'X':
			// X → KS normally, but not at end for French.
			if !(pos == n-1 &&
				(dmStringAt(padded, pos-3, "IAU", "EAU") ||
					dmStringAt(padded, pos-2, "AU", "OU"))) {
				pri.WriteString("KS")
				alt.WriteString("KS")
			}
			if padded[pos+1] == 'C' || padded[pos+1] == 'X' {
				pos += 2
			} else {
				pos++
			}

		case 'Z':
			pos = dmHandleZ(padded, pos, n, &pri, &alt, slavoGermanic)

		default:
			pos++
		}
	}

	p := pri.String()
	a := alt.String()

	if len(p) > maxLen {
		p = p[:maxLen]
	}

	if len(a) > maxLen {
		a = a[:maxLen]
	}

	return p, a
}

// Match reports whether two strings share at least one common Double Metaphone
// code.
func (e *DoubleMetaphoneEncoder) Match(a, b string) bool {
	p1, a1 := e.Encode(a)
	p2, a2 := e.Encode(b)

	if p1 == "" && a1 == "" && p2 == "" && a2 == "" {
		return true
	}

	if (p1 == "" && a1 == "") || (p2 == "" && a2 == "") {
		return false
	}

	return p1 == p2 || p1 == a2 || a1 == p2 || a1 == a2
}

// dmHandleC processes the letter C in Double Metaphone.
func dmHandleC(w string, pos, n int, pri, alt *strings.Builder, slavo bool) int {
	// Various Germanic.
	if pos > 1 && !dmIsVowel(w[pos-2]) &&
		w[pos-1] == 'A' && w[pos] == 'C' && w[pos+1] == 'H' &&
		w[pos+2] != 'I' && (w[pos+2] != 'E' || dmStringAt(w, pos-2, "BACHER", "MACHER")) {
		pri.WriteByte('K')
		alt.WriteByte('K')
		return pos + 2
	}

	// Special initial: CAESAR.
	if pos == 0 && dmStringAt(w, pos, "CAESAR") {
		pri.WriteByte('S')
		alt.WriteByte('S')
		return pos + 2
	}

	// Italian: CHIANTI.
	if dmStringAt(w, pos, "CHIA") {
		pri.WriteByte('K')
		alt.WriteByte('K')
		return pos + 2
	}

	// CH
	if w[pos+1] == 'H' {
		// Greek: CHORUS, CHEMISTRY, etc.
		if pos > 0 && dmStringAt(w, pos, "CHAE") {
			pri.WriteByte('K')
			alt.WriteByte('X')
			return pos + 2
		}

		// Greek roots: CHARAC, CHARIS, CHOR, CHEM, CHYM.
		if pos == 0 &&
			(dmStringAt(w, pos+1, "HARAC", "HARIS") ||
				dmStringAt(w, pos+1, "HOR", "HYM", "HIA", "HEM")) &&
			!dmStringAt(w, 0, "CHORE") {
			pri.WriteByte('K')
			alt.WriteByte('K')
			return pos + 2
		}

		// Germanic/Greek: handle CH.
		if dmStringAt(w, 0, "VAN ", "VON ") ||
			dmStringAt(w, 0, "SCH") ||
			dmStringAt(w, pos-2, "ORCHES", "ARCHIT", "ORCHID") ||
			dmStringAt(w, pos+2, "T", "S") ||
			((pos == 0 || dmStringAt(w, pos-1, "A", "O", "U", "E")) &&
				dmStringAt(w, pos+2, "L", "R", "N", "M", "B", "H", "F", "V", "W", " ")) {
			pri.WriteByte('K')
			alt.WriteByte('K')
		} else if pos > 0 {
			if dmStringAt(w, 0, "MC") {
				// McHugh.
				pri.WriteByte('K')
				alt.WriteByte('K')
			} else {
				pri.WriteByte('X')
				alt.WriteByte('K')
			}
		} else {
			pri.WriteByte('X')
			alt.WriteByte('X')
		}

		return pos + 2
	}

	// CZ (Polish) but not WICZ.
	if dmStringAt(w, pos, "CZ") && !dmStringAt(w, pos-2, "WICZ") {
		pri.WriteByte('S')
		alt.WriteByte('X')
		return pos + 2
	}

	// CIA: Italian.
	if dmStringAt(w, pos+1, "IA") {
		pri.WriteByte('X')
		alt.WriteByte('X')
		return pos + 3
	}

	// CC but not McClellan.
	if w[pos+1] == 'C' && !(pos == 1 && w[0] == 'M') {
		// Bellocchio, Bacchus.
		if (w[pos+2] == 'I' || w[pos+2] == 'E' || w[pos+2] == 'H') &&
			!dmStringAt(w, pos+2, "HU") {
			// ACCIO, UCCEE, etc.
			if (pos == 1 && w[0] == 'A') || dmStringAt(w, pos-1, "UCCEE", "UCCES") {
				pri.WriteString("KS")
				alt.WriteString("KS")
			} else {
				pri.WriteByte('X')
				alt.WriteByte('X')
			}
			return pos + 3
		}
		// CC → K (e.g., JACK).
		pri.WriteByte('K')
		alt.WriteByte('K')
		return pos + 2
	}

	if dmStringAt(w, pos, "CK", "CG", "CQ") {
		pri.WriteByte('K')
		alt.WriteByte('K')
		return pos + 2
	}

	if dmStringAt(w, pos, "CI", "CE", "CY") {
		// Italian vs. rest.
		if dmStringAt(w, pos, "CIO", "CIE", "CIA") {
			pri.WriteByte('S')
			alt.WriteByte('X')
		} else {
			pri.WriteByte('S')
			alt.WriteByte('S')
		}
		return pos + 2
	}

	// Default.
	pri.WriteByte('K')
	alt.WriteByte('K')

	// Handle CL, CR, etc. vs CC.
	if dmStringAt(w, pos+1, " C", " Q", " G") {
		return pos + 3
	}

	if w[pos+1] == 'C' || w[pos+1] == 'K' || w[pos+1] == 'Q' {
		if !dmStringAt(w, pos+1, "CE", "CI") {
			return pos + 2
		}
	}

	return pos + 1
}

// dmHandleD processes the letter D in Double Metaphone.
func dmHandleD(w string, pos, _ int, pri, alt *strings.Builder) int {
	if w[pos+1] == 'G' {
		if w[pos+2] == 'I' || w[pos+2] == 'E' || w[pos+2] == 'Y' {
			// DGI, DGE, DGY → J.
			pri.WriteByte('J')
			alt.WriteByte('J')
			return pos + 3
		}
		// DG → TK.
		pri.WriteString("TK")
		alt.WriteString("TK")
		return pos + 2
	}

	if w[pos+1] == 'T' || w[pos+1] == 'D' {
		pri.WriteByte('T')
		alt.WriteByte('T')
		return pos + 2
	}

	pri.WriteByte('T')
	alt.WriteByte('T')

	return pos + 1
}

// dmHandleG processes the letter G in Double Metaphone.
func dmHandleG(w string, pos, n int, pri, alt *strings.Builder, slavo bool) int {
	if w[pos+1] == 'H' {
		return dmHandleGH(w, pos, n, pri, alt)
	}

	if w[pos+1] == 'N' {
		if pos == 1 && dmIsVowel(w[0]) && !slavo {
			pri.WriteString("KN")
			alt.WriteString("N")
		} else if !dmStringAt(w, pos+2, "EY") && w[pos+1] != 'Y' && !slavo {
			pri.WriteString("N")
			alt.WriteString("KN")
		} else {
			pri.WriteString("KN")
			alt.WriteString("KN")
		}
		return pos + 2
	}

	// Italian: GLIA.
	if dmStringAt(w, pos+1, "LI") && !slavo {
		pri.WriteString("KL")
		alt.WriteString("L")
		return pos + 2
	}

	// Initial GY, GE, GI, GEI, GER, GES.
	if pos == 0 {
		if w[pos+1] == 'Y' {
			pri.WriteByte('K')
			alt.WriteByte('J')
			return pos + 2
		}

		if w[pos+1] == 'E' || w[pos+1] == 'I' {
			pri.WriteByte('K')
			alt.WriteByte('J')
			return pos + 2
		}
	}

	// Parker's rule for -GER, -GEL at end.
	if (dmStringAt(w, pos+1, "ER") || w[pos+1] == 'Y') &&
		!dmStringAt(w, 0, "DANGER", "RANGER", "MANGER") &&
		!dmStringAt(w, pos-1, "E", "I") &&
		!dmStringAt(w, pos-1, "RGY", "OGY") {
		pri.WriteByte('K')
		alt.WriteByte('J')
		return pos + 2
	}

	// Italian.
	if w[pos+1] == 'E' || w[pos+1] == 'I' || w[pos+1] == 'Y' ||
		dmStringAt(w, pos-1, "AGGI", "OGGI") {
		// Germanic.
		if dmStringAt(w, 0, "VAN ", "VON ") ||
			dmStringAt(w, 0, "SCH") ||
			dmStringAt(w, pos+1, "ET") {
			pri.WriteByte('K')
			alt.WriteByte('K')
		} else if dmStringAt(w, pos+1, "IER") {
			pri.WriteByte('J')
			alt.WriteByte('J')
		} else {
			pri.WriteByte('J')
			alt.WriteByte('K')
		}
		return pos + 2
	}

	if w[pos+1] == 'G' {
		pos += 2
	} else {
		pos++
	}

	pri.WriteByte('K')
	alt.WriteByte('K')

	return pos
}

// dmHandleGH processes GH combinations in Double Metaphone.
func dmHandleGH(w string, pos, n int, pri, alt *strings.Builder) int {
	if pos > 0 && !dmIsVowel(w[pos-1]) {
		pri.WriteByte('K')
		alt.WriteByte('K')
		return pos + 2
	}

	// GH at beginning.
	if pos == 0 {
		if w[pos+2] == 'I' {
			pri.WriteByte('J')
			alt.WriteByte('J')
		} else {
			pri.WriteByte('K')
			alt.WriteByte('K')
		}
		return pos + 2
	}

	// Parker's rule: GH after a back vowel (except B).
	if (pos > 1 && dmStringAt(w, pos-2, "B", "H", "D")) ||
		(pos > 2 && dmStringAt(w, pos-3, "B", "H", "D")) ||
		(pos > 3 && dmStringAt(w, pos-4, "B", "H")) {
		return pos + 2
	}

	// Italian/Germanic: GH not silent after another pattern.
	if pos > 2 && w[pos-1] == 'U' &&
		dmStringAt(w, pos-3, "C", "G", "L", "R", "T") {
		pri.WriteByte('F')
		alt.WriteByte('F')
	} else if pos > 0 && w[pos-1] != 'I' {
		pri.WriteByte('K')
		alt.WriteByte('K')
	}

	return pos + 2
}

// dmHandleJ processes the letter J in Double Metaphone.
func dmHandleJ(w string, pos, n int, word string, pri, alt *strings.Builder, slavo bool) int {
	// Spanish: SAN JOSE.
	if dmStringAt(w, pos, "JOSE") || dmStringAt(w, 0, "SAN ") {
		if (pos == 0 && w[pos+4] == ' ') || dmStringAt(w, 0, "SAN ") {
			pri.WriteByte('H')
			alt.WriteByte('H')
		} else {
			pri.WriteByte('J')
			alt.WriteByte('H')
		}
		return pos + 1
	}

	if pos == 0 {
		// YAN → J.
		pri.WriteByte('J')
		alt.WriteByte('A')
	} else if dmIsVowel(w[pos-1]) && !slavo && (w[pos+1] == 'A' || w[pos+1] == 'O') {
		pri.WriteByte('J')
		alt.WriteByte('H')
	} else if pos == n-1 {
		pri.WriteByte('J')
		alt.WriteByte(' ')
	} else if !dmStringAt(w, pos+1, "L", "T", "K", "S", "N", "M", "B", "Z") &&
		!dmStringAt(w, pos-1, "S", "K", "L") {
		pri.WriteByte('J')
		alt.WriteByte('J')
	}

	if w[pos+1] == 'J' {
		return pos + 2
	}

	return pos + 1
}

// dmHandleL processes the letter L in Double Metaphone.
func dmHandleL(w string, pos, n int, pri, alt *strings.Builder) int {
	if w[pos+1] == 'L' {
		// Spanish: CABRILLO, GALLEGOS.
		if (pos == n-3 && dmStringAt(w, pos-1, "ILLO", "ILLA", "ALLE")) ||
			((dmStringAt(w, n-2, "AS", "OS") || dmStringAt(w, n-1, "A", "O")) &&
				dmStringAt(w, pos-1, "ALLE")) {
			pri.WriteByte('L')
			alt.WriteByte(' ')
			return pos + 2
		}
		pos += 2
	} else {
		pos++
	}

	pri.WriteByte('L')
	alt.WriteByte('L')

	return pos
}

// dmHandleS processes the letter S in Double Metaphone.
func dmHandleS(w string, pos, n int, pri, alt *strings.Builder, slavo bool) int {
	// Special: ISLAND, ISLE, CARLISLE, CARLYSLE.
	if dmStringAt(w, pos-1, "ISL", "YSL") {
		return pos + 1
	}

	// Special: SUGAR.
	if pos == 0 && dmStringAt(w, pos, "SUGAR") {
		pri.WriteByte('X')
		alt.WriteByte('S')
		return pos + 1
	}

	if dmStringAt(w, pos, "SH") {
		// Germanic: SH → X in most cases.
		if dmStringAt(w, pos+1, "HEIM", "HOEK", "HOLM", "HOLZ") {
			pri.WriteByte('S')
			alt.WriteByte('S')
		} else {
			pri.WriteByte('X')
			alt.WriteByte('X')
		}
		return pos + 2
	}

	// Italian and Armenian: SIO, SIA → S/X.
	if dmStringAt(w, pos, "SIO", "SIA") || dmStringAt(w, pos, "SIAN") {
		if !slavo {
			pri.WriteByte('S')
			alt.WriteByte('X')
		} else {
			pri.WriteByte('S')
			alt.WriteByte('S')
		}
		return pos + 3
	}

	// German & English: initial SC+E,I → S, SC+H → SK.
	if pos == 0 && dmStringAt(w, pos+1, "M", "N", "L", "W") {
		pri.WriteByte('S')
		alt.WriteByte('X')
		return pos + 1
	}

	if dmStringAt(w, pos, "SCH") {
		// Schoenberg.
		if dmStringAt(w, pos+3, "OO", "ER", "EN", "UY", "ED", "EM") {
			// Dutch: SCHOOL, SCHOONER.
			if dmStringAt(w, pos+3, "ER", "EN") {
				pri.WriteString("X")
				alt.WriteString("SK")
			} else {
				pri.WriteString("SK")
				alt.WriteString("SK")
			}
			return pos + 3
		}

		if pos == 0 && !dmIsVowel(w[3]) && w[3] != 'W' {
			pri.WriteByte('X')
			alt.WriteByte('S')
		} else {
			pri.WriteByte('X')
			alt.WriteByte('X')
		}
		return pos + 3
	}

	if dmStringAt(w, pos, "SC") {
		if w[pos+2] == 'I' || w[pos+2] == 'E' || w[pos+2] == 'Y' {
			pri.WriteByte('S')
			alt.WriteByte('S')
		} else {
			pri.WriteString("SK")
			alt.WriteString("SK")
		}
		return pos + 3
	}

	// French: final S silent.
	if pos == n-1 && dmStringAt(w, pos-2, "AI", "OI") {
		alt.WriteByte('S')
	} else {
		pri.WriteByte('S')
		alt.WriteByte('S')
	}

	if w[pos+1] == 'S' || w[pos+1] == 'Z' {
		return pos + 2
	}

	return pos + 1
}

// dmHandleT processes the letter T in Double Metaphone.
func dmHandleT(w string, pos, n int, pri, alt *strings.Builder) int {
	if dmStringAt(w, pos, "TION") || dmStringAt(w, pos, "TIA", "TCH") {
		pri.WriteByte('X')
		alt.WriteByte('X')
		return pos + 3
	}

	if dmStringAt(w, pos, "TH") || dmStringAt(w, pos, "TTH") {
		// Greek origin.
		if dmStringAt(w, pos+2, "OM", "AM") || dmStringAt(w, 0, "VAN ", "VON ") || dmStringAt(w, 0, "SCH") {
			pri.WriteByte('T')
			alt.WriteByte('T')
		} else {
			pri.WriteByte('0')
			alt.WriteByte('T')
		}
		return pos + 2
	}

	if w[pos+1] == 'T' || w[pos+1] == 'D' {
		pos += 2
	} else {
		pos++
	}

	pri.WriteByte('T')
	alt.WriteByte('T')

	return pos
}

// dmHandleW processes the letter W in Double Metaphone.
func dmHandleW(w string, pos, n int, pri, alt *strings.Builder) int {
	// WR → R.
	if w[pos+1] == 'R' {
		pri.WriteByte('R')
		alt.WriteByte('R')
		return pos + 2
	}

	if pos == 0 {
		if dmIsVowel(w[pos+1]) {
			// WH → A in alternate.
			pri.WriteByte('A')
			alt.WriteByte('F')
		} else if w[pos+1] == 'H' {
			// WH → A.
			pri.WriteByte('A')
			alt.WriteByte('A')
		}
	}

	// WICZ, WITZ: Polish.
	if (pos == n-1 && dmIsVowel(w[pos-1])) ||
		dmStringAt(w, pos, "WICZ", "WITZ") ||
		dmStringAt(w, pos-1, "EWSKI", "EWSKY", "OWSKI", "OWSKY") {
		pri.WriteByte(' ')
		alt.WriteString("F")
		return pos + 1
	}

	return pos + 1
}

// dmHandleZ processes the letter Z in Double Metaphone.
func dmHandleZ(w string, pos, n int, pri, alt *strings.Builder, slavo bool) int {
	if w[pos+1] == 'H' {
		// Chinese: ZH → J.
		pri.WriteByte('J')
		alt.WriteByte('J')
		return pos + 2
	}

	if dmStringAt(w, pos+1, "ZO", "ZI", "ZA") || (slavo && pos > 0 && w[pos-1] != 'T') {
		pri.WriteByte('S')
		alt.WriteString("TS")
	} else {
		pri.WriteByte('S')
		alt.WriteByte('S')
	}

	if w[pos+1] == 'Z' {
		return pos + 2
	}

	return pos + 1
}

// dmIsVowel reports whether ch is an ASCII vowel.
func dmIsVowel(ch byte) bool {
	switch ch {
	case 'A', 'E', 'I', 'O', 'U':
		return true
	}

	return false
}

// dmIsSlavoGermanic reports whether the word has Slavo-Germanic characteristics.
func dmIsSlavoGermanic(s string) bool {
	return strings.Contains(s, "W") ||
		strings.Contains(s, "K") ||
		strings.Contains(s, "CZ") ||
		strings.Contains(s, "WITZ")
}

// dmStringAt reports whether any of the given substrings appear at position pos
// in the string s.
func dmStringAt(s string, pos int, subs ...string) bool {
	if pos < 0 {
		return false
	}

	for _, sub := range subs {
		end := pos + len(sub)
		if end <= len(s) && s[pos:end] == sub {
			return true
		}
	}

	return false
}
