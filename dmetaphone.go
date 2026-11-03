package strsim

import "strings"

// defaultDMMaxLength is the standard Double Metaphone code length.
const defaultDMMaxLength = 4

// DoubleMetaphoneEncoder implements the Double Metaphone algorithm by Lawrence
// Philips. It produces two phonetic codes (primary and alternate) that capture
// different possible pronunciations, accounting for non-English name origins
// including Germanic, Slavic, Celtic, Greek, Italian, Spanish, and Chinese.
//
// The rules follow Philips' reference implementation as ported by Apache
// Commons Codec. Unlike that port, codes never contain spaces: where the
// reference adds nothing to one of the codes ("Raj" → RJ / R), so does this
// encoder.
type DoubleMetaphoneEncoder struct {
	// MaxLength is the maximum length of each generated code. Values <= 0 use
	// the default of 4.
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
// given string. Latin diacritics are folded (é → E), and words are kept apart
// because some rules look at word boundaries ("San Jose", "Mac Gregor"). Both
// codes are empty if the input contains no letters.
func (e *DoubleMetaphoneEncoder) Encode(s string) (primary, alternate string) {
	maxLen := e.MaxLength
	if maxLen <= 0 {
		maxLen = defaultDMMaxLength
	}

	var lettersBuf [32]byte
	w := string(appendLetters(lettersBuf[:0], s, true))
	if w == "" {
		return "", ""
	}

	d := dmState{w: w, slavo: dmIsSlavoGermanic(w), max: maxLen}
	d.pri.Grow(maxLen)
	d.alt.Grow(maxLen)
	d.encode()

	return d.pri.String(), d.alt.String()
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

// dmState holds the word being encoded and both codes built so far.
type dmState struct {
	w     string
	slavo bool
	max   int
	pri   strings.Builder
	alt   strings.Builder
}

// at returns the byte at i, or 0 when i is out of range.
func (d *dmState) at(i int) byte {
	return byteAt(d.w, i)
}

// has reports whether any of subs occurs at position pos.
func (d *dmState) has(pos int, subs ...string) bool {
	return hasAt(d.w, pos, subs...)
}

// add appends p to the primary code and a to the alternate code, each
// truncated to the maximum length.
func (d *dmState) add(p, a string) {
	writeUpTo(&d.pri, p, d.max)
	writeUpTo(&d.alt, a, d.max)
}

// both appends s to both codes.
func (d *dmState) both(s string) {
	d.add(s, s)
}

// writeUpTo writes s to code without exceeding limit bytes.
func writeUpTo(code *strings.Builder, s string, limit int) {
	if room := limit - code.Len(); len(s) > room {
		s = s[:max(room, 0)]
	}

	code.WriteString(s)
}

// skip returns the next position after i, skipping one repeat of ch.
func (d *dmState) skip(i int, ch byte) int {
	if d.at(i+1) == ch {
		return i + 2
	}

	return i + 1
}

func (d *dmState) encode() {
	i := 0
	if d.has(0, "GN", "KN", "PN", "WR", "PS") {
		i = 1
	}

	for i < len(d.w) && (d.pri.Len() < d.max || d.alt.Len() < d.max) {
		switch d.w[i] {
		case 'A', 'E', 'I', 'O', 'U', 'Y':
			if i == 0 {
				d.both("A")
			}
			i++
		case 'B':
			d.both("P")
			i = d.skip(i, 'B')
		case 'C':
			i = d.c(i)
		case 'D':
			i = d.d(i)
		case 'F':
			d.both("F")
			i = d.skip(i, 'F')
		case 'G':
			i = d.g(i)
		case 'H':
			i = d.h(i)
		case 'J':
			i = d.j(i)
		case 'K':
			d.both("K")
			i = d.skip(i, 'K')
		case 'L':
			i = d.l(i)
		case 'M':
			d.both("M")
			if d.at(i+1) == 'M' || (d.has(i-1, "UMB") && (i+1 == len(d.w)-1 || d.has(i+2, "ER"))) {
				i += 2 // "dumb", "thumb"
			} else {
				i++
			}
		case 'N':
			d.both("N")
			i = d.skip(i, 'N')
		case 'P':
			i = d.p(i)
		case 'Q':
			d.both("K")
			i = d.skip(i, 'Q')
		case 'R':
			i = d.r(i)
		case 'S':
			i = d.s(i)
		case 'T':
			i = d.t(i)
		case 'V':
			d.both("F")
			i = d.skip(i, 'V')
		case 'W':
			i = d.wch(i)
		case 'X':
			i = d.x(i)
		case 'Z':
			i = d.z(i)
		default:
			i++
		}
	}
}

func (d *dmState) c(i int) int {
	switch {
	case d.isGermanicCH(i):
		d.both("K")
		return i + 2
	case i == 0 && d.has(i, "CAESAR"):
		d.both("S")
		return i + 2
	case d.has(i, "CH"):
		return d.ch(i)
	case d.has(i, "CZ") && !d.has(i-2, "WICZ"):
		d.add("S", "X") // "Czerny"
		return i + 2
	case d.has(i+1, "CIA"):
		d.both("X") // "focaccia"
		return i + 3
	case d.has(i, "CC") && (i != 1 || d.at(0) != 'M'):
		// Double C, but not "McClellan".
		if d.has(i+2, "I", "E", "H") && !d.has(i+2, "HU") {
			// "bellocchio" but not "bacchus".
			if (i == 1 && d.at(0) == 'A') || d.has(i-1, "UCCEE", "UCCES") {
				d.both("KS") // "accident", "accede", "succeed"
			} else {
				d.both("X") // "bacci", "bertucci"
			}
			return i + 3
		}
		d.both("K") // Pierce's rule
		return i + 2
	case d.has(i, "CK", "CG", "CQ"):
		d.both("K")
		return i + 2
	case d.has(i, "CI", "CE", "CY"):
		if d.has(i, "CIO", "CIE", "CIA") {
			d.add("S", "X") // Italian
		} else {
			d.both("S")
		}
		return i + 2
	}

	d.both("K")

	switch {
	case d.has(i+1, " C", " Q", " G"):
		return i + 3 // "Mac Caffrey", "Mac Gregor"
	case d.has(i+1, "C", "K", "Q") && !d.has(i+1, "CE", "CI"):
		return i + 2
	}

	return i + 1
}

// isGermanicCH reports the Italian "chianti" and the Germanic -ACH- cases,
// where CH sounds like K.
func (d *dmState) isGermanicCH(i int) bool {
	if d.has(i, "CHIA") {
		return true
	}

	if i <= 1 || dmIsVowel(d.at(i-2)) || !d.has(i-1, "ACH") {
		return false
	}

	c := d.at(i + 2)

	return (c != 'I' && c != 'E') || d.has(i-2, "BACHER", "MACHER")
}

func (d *dmState) ch(i int) int {
	switch {
	case i > 0 && d.has(i, "CHAE"):
		d.add("K", "X") // "Michael"
	case i == 0 && (d.has(i+1, "HARAC", "HARIS") || d.has(i+1, "HOR", "HYM", "HIA", "HEM")) &&
		!d.has(0, "CHORE"):
		d.both("K") // Greek roots: "chemistry", "chorus"
	case d.has(0, "VAN ", "VON ", "SCH") ||
		d.has(i-2, "ORCHES", "ARCHIT", "ORCHID") ||
		d.has(i+2, "T", "S") ||
		((i == 0 || d.has(i-1, "A", "O", "U", "E")) &&
			(d.has(i+2, "L", "R", "N", "M", "B", "H", "F", "V", "W", " ") || i+1 == len(d.w)-1)):
		d.both("K") // Germanic, Greek, or otherwise CH for the KH sound
	case i > 0 && d.has(0, "MC"):
		d.both("K") // "McHugh"
	case i > 0:
		d.add("X", "K")
	default:
		d.both("X")
	}

	return i + 2
}

func (d *dmState) d(i int) int {
	switch {
	case d.has(i, "DG"):
		if d.has(i+2, "I", "E", "Y") {
			d.both("J") // "edge"
			return i + 3
		}
		d.both("TK") // "Edgar"
		return i + 2
	case d.has(i, "DT", "DD"):
		d.both("T")
		return i + 2
	}

	d.both("T")

	return i + 1
}

func (d *dmState) g(i int) int {
	next := d.at(i + 1)

	switch {
	case next == 'H':
		return d.gh(i)
	case next == 'N':
		switch {
		case i == 1 && dmIsVowel(d.at(0)) && !d.slavo:
			d.add("KN", "N")
		case !d.has(i+2, "EY") && !d.slavo:
			d.add("N", "KN")
		default:
			d.both("KN")
		}
		return i + 2
	case d.has(i+1, "LI") && !d.slavo:
		d.add("KL", "L") // "tagliaro"
		return i + 2
	case i == 0 && (next == 'Y' || d.has(i+1, "ES", "EP", "EB", "EL", "EY", "IB", "IL", "IN", "IE", "EI", "ER")):
		d.add("K", "J") // -ges-, -gep-, -gel-, -gie- at the beginning
		return i + 2
	case (d.has(i+1, "ER") || next == 'Y') &&
		!d.has(0, "DANGER", "RANGER", "MANGER") &&
		!d.has(i-1, "E", "I") &&
		!d.has(i-1, "RGY", "OGY"):
		d.add("K", "J") // -ger-, -gy-
		return i + 2
	case d.has(i+1, "E", "I", "Y") || d.has(i-1, "AGGI", "OGGI"):
		switch {
		case d.has(0, "VAN ", "VON ", "SCH") || d.has(i+1, "ET"):
			d.both("K") // obviously Germanic
		case d.has(i+1, "IER"):
			d.both("J")
		default:
			d.add("J", "K") // Italian "biaggi"
		}
		return i + 2
	}

	d.both("K")

	return d.skip(i, 'G')
}

func (d *dmState) gh(i int) int {
	switch {
	case i > 0 && !dmIsVowel(d.at(i-1)):
		d.both("K")
	case i == 0:
		if d.at(i+2) == 'I' {
			d.both("J") // "ghislane"
		} else {
			d.both("K") // "ghost"
		}
	case (i > 1 && d.has(i-2, "B", "H", "D")) ||
		(i > 2 && d.has(i-3, "B", "H", "D")) ||
		(i > 3 && d.has(i-4, "B", "H")):
		// Parker's rule: silent in "hugh", "bough", "broughton".
	case i > 2 && d.at(i-1) == 'U' && d.has(i-3, "C", "G", "L", "R", "T"):
		d.both("F") // "laugh", "McLaughlin", "cough", "rough", "tough"
	case d.at(i-1) != 'I':
		d.both("K")
	}

	return i + 2
}

func (d *dmState) h(i int) int {
	// Only kept as the first letter or between two vowels.
	if (i == 0 || dmIsVowel(d.at(i-1))) && dmIsVowel(d.at(i+1)) {
		d.both("H")
		return i + 2 // also skips HH
	}

	return i + 1
}

func (d *dmState) j(i int) int {
	if d.has(i, "JOSE") || d.has(0, "SAN ") {
		// Obviously Spanish: "Jose", "San Jacinto".
		if (i == 0 && d.at(i+4) == ' ') || len(d.w) == 4 || d.has(0, "SAN ") {
			d.both("H")
		} else {
			d.add("J", "H")
		}
		return i + 1
	}

	switch {
	case i == 0:
		d.add("J", "A") // "Yankelovich" / "Jankelowicz"
	case dmIsVowel(d.at(i-1)) && !d.slavo && (d.at(i+1) == 'A' || d.at(i+1) == 'O'):
		d.add("J", "H") // Spanish pronunciation of "bajador"
	case i == len(d.w)-1:
		d.add("J", "")
	case !d.has(i+1, "L", "T", "K", "S", "N", "M", "B", "Z") && !d.has(i-1, "S", "K", "L"):
		d.both("J")
	}

	return d.skip(i, 'J')
}

func (d *dmState) l(i int) int {
	if d.at(i+1) != 'L' {
		d.both("L")
		return i + 1
	}

	n := len(d.w)
	if (i == n-3 && d.has(i-1, "ILLO", "ILLA", "ALLE")) ||
		((d.has(n-2, "AS", "OS") || d.has(n-1, "A", "O")) && d.has(i-1, "ALLE")) {
		d.add("L", "") // Spanish: "cabrillo", "gallegos"
	} else {
		d.both("L")
	}

	return i + 2
}

func (d *dmState) p(i int) int {
	if d.at(i+1) == 'H' {
		d.both("F")
		return i + 2
	}

	d.both("P")

	if d.has(i+1, "P", "B") {
		return i + 2 // "campbell"
	}

	return i + 1
}

func (d *dmState) r(i int) int {
	// French final R, silent in the primary code: "Rogier".
	if i == len(d.w)-1 && !d.slavo && d.has(i-2, "IE") && !d.has(i-4, "ME", "MA") {
		d.add("", "R")
	} else {
		d.both("R")
	}

	return d.skip(i, 'R')
}

func (d *dmState) s(i int) int {
	switch {
	case d.has(i-1, "ISL", "YSL"):
		return i + 1 // "island", "isle", "carlisle", "carlysle"
	case i == 0 && d.has(i, "SUGAR"):
		d.add("X", "S")
		return i + 1
	case d.has(i, "SH"):
		if d.has(i+1, "HEIM", "HOEK", "HOLM", "HOLZ") {
			d.both("S") // Germanic
		} else {
			d.both("X")
		}
		return i + 2
	case d.has(i, "SIO", "SIA", "SIAN"):
		// Italian and Armenian.
		if d.slavo {
			d.both("S")
		} else {
			d.add("S", "X")
		}
		return i + 3
	case (i == 0 && d.has(i+1, "M", "N", "L", "W")) || d.has(i+1, "Z"):
		// German and anglicisations ("Smith" matches "Schmidt", "Snider"
		// matches "Schneider"); -SZ- in Slavic languages.
		d.add("S", "X")
		return d.skip(i, 'Z')
	case d.has(i, "SC"):
		return d.sc(i)
	}

	if i == len(d.w)-1 && d.has(i-2, "AI", "OI") {
		d.add("", "S") // French: "Resnais", "Artois"
	} else {
		d.both("S")
	}

	if d.has(i+1, "S", "Z") {
		return i + 2
	}

	return i + 1
}

func (d *dmState) sc(i int) int {
	switch {
	case d.at(i+2) == 'H':
		// Schlesinger's rule.
		switch {
		case d.has(i+3, "ER", "EN"):
			d.add("X", "SK") // "schermerhorn", "schenker"
		case d.has(i+3, "OO", "UY", "ED", "EM"):
			d.both("SK") // Dutch: "school", "schooner"
		case i == 0 && !dmIsVowel(d.at(3)) && d.at(3) != 'W':
			d.add("X", "S")
		default:
			d.both("X")
		}
	case d.has(i+2, "I", "E", "Y"):
		d.both("S")
	default:
		d.both("SK")
	}

	return i + 3
}

func (d *dmState) t(i int) int {
	switch {
	case d.has(i, "TION", "TIA", "TCH"):
		d.both("X")
		return i + 3
	case d.has(i, "TH", "TTH"):
		if d.has(i+2, "OM", "AM") || d.has(0, "VAN ", "VON ", "SCH") {
			d.both("T") // "thomas", "thames", or Germanic
		} else {
			d.add("0", "T")
		}
		return i + 2
	}

	d.both("T")

	if d.has(i+1, "T", "D") {
		return i + 2
	}

	return i + 1
}

// wch handles W; the name w is taken by the word field.
func (d *dmState) wch(i int) int {
	if d.has(i, "WR") {
		d.both("R") // can also be in the middle of a word
		return i + 2
	}

	switch {
	case i == 0 && dmIsVowel(d.at(i+1)):
		d.add("A", "F") // "Wasserman" matches "Vasserman"
	case i == 0 && d.at(i+1) == 'H':
		d.both("A") // "Uomo" matches "Womo"
	case (i == len(d.w)-1 && dmIsVowel(d.at(i-1))) ||
		d.has(i-1, "EWSKI", "EWSKY", "OWSKI", "OWSKY") ||
		d.has(0, "SCH"):
		d.add("", "F") // "Arnow" matches "Arnoff"
	case d.has(i, "WICZ", "WITZ"):
		d.add("TS", "FX") // Polish: "Filipowicz"
		return i + 4
	}

	return i + 1
}

func (d *dmState) x(i int) int {
	if i == 0 {
		d.both("S")
		return i + 1
	}

	// Silent in French final -IAUX, -EAUX, -AUX, -OUX: "Breaux".
	frenchFinal := i == len(d.w)-1 && (d.has(i-3, "IAU", "EAU") || d.has(i-2, "AU", "OU"))
	if !frenchFinal {
		d.both("KS")
	}

	if d.has(i+1, "C", "X") {
		return i + 2
	}

	return i + 1
}

func (d *dmState) z(i int) int {
	if d.at(i+1) == 'H' {
		d.both("J") // Chinese pinyin: "Zhao"
		return i + 2
	}

	if d.has(i+1, "ZO", "ZI", "ZA") || (d.slavo && i > 0 && d.at(i-1) != 'T') {
		d.add("S", "TS")
	} else {
		d.both("S")
	}

	return d.skip(i, 'Z')
}

// dmIsVowel reports whether ch is a vowel for Double Metaphone, which counts Y.
func dmIsVowel(ch byte) bool {
	switch ch {
	case 'A', 'E', 'I', 'O', 'U', 'Y':
		return true
	}

	return false
}

// dmIsSlavoGermanic reports whether the word has Slavo-Germanic characteristics.
func dmIsSlavoGermanic(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == 'W' || s[i] == 'K' || (s[i] == 'C' && byteAt(s, i+1) == 'Z') {
			return true
		}
	}

	return false
}
