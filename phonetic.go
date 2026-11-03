package strsim

import "unicode"

// Phonetic encoders work on uppercase ASCII letters. Latin letters with
// diacritics (U+00C0–U+017F) are folded to their base letters first, so that
// "Muñoz" encodes like "Munoz" and "Ångström" like "Angstrom"; combining
// accents of decomposed text are dropped. Other runes are ignored.

// upperLetters is sliced to turn an ASCII letter into a string without
// allocating.
const upperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// latinFoldFirst is the first rune covered by latinFold.
const latinFoldFirst = 0xC0

// latinFold maps U+00C0–U+017F (Latin-1 Supplement letters and Latin
// Extended-A) to uppercase ASCII. An empty entry means the rune is not a
// letter (× and ÷).
var latinFold = [...]string{
	"A", "A", "A", "A", "A", "A", "AE", "C", // U+00C0 ÀÁÂÃÄÅÆÇ
	"E", "E", "E", "E", "I", "I", "I", "I", // U+00C8 ÈÉÊËÌÍÎÏ
	"D", "N", "O", "O", "O", "O", "O", "", // U+00D0 ÐÑÒÓÔÕÖ×
	"O", "U", "U", "U", "U", "Y", "TH", "SS", // U+00D8 ØÙÚÛÜÝÞß
	"A", "A", "A", "A", "A", "A", "AE", "C", // U+00E0 àáâãäåæç
	"E", "E", "E", "E", "I", "I", "I", "I", // U+00E8 èéêëìíîï
	"D", "N", "O", "O", "O", "O", "O", "", // U+00F0 ðñòóôõö÷
	"O", "U", "U", "U", "U", "Y", "TH", "Y", // U+00F8 øùúûüýþÿ
	"A", "A", "A", "A", "A", "A", "C", "C", // U+0100 ĀāĂăĄąĆć
	"C", "C", "C", "C", "C", "C", "D", "D", // U+0108 ĈĉĊċČčĎď
	"D", "D", "E", "E", "E", "E", "E", "E", // U+0110 ĐđĒēĔĕĖė
	"E", "E", "E", "E", "G", "G", "G", "G", // U+0118 ĘęĚěĜĝĞğ
	"G", "G", "G", "G", "H", "H", "H", "H", // U+0120 ĠġĢģĤĥĦħ
	"I", "I", "I", "I", "I", "I", "I", "I", // U+0128 ĨĩĪīĬĭĮį
	"I", "I", "IJ", "IJ", "J", "J", "K", "K", // U+0130 İıĲĳĴĵĶķ
	"K", "L", "L", "L", "L", "L", "L", "L", // U+0138 ĸĹĺĻļĽľĿ
	"L", "L", "L", "N", "N", "N", "N", "N", // U+0140 ŀŁłŃńŅņŇ
	"N", "N", "N", "N", "O", "O", "O", "O", // U+0148 ňŉŊŋŌōŎŏ
	"O", "O", "OE", "OE", "R", "R", "R", "R", // U+0150 ŐőŒœŔŕŖŗ
	"R", "R", "S", "S", "S", "S", "S", "S", // U+0158 ŘřŚśŜŝŞş
	"S", "S", "T", "T", "T", "T", "T", "T", // U+0160 ŠšŢţŤťŦŧ
	"U", "U", "U", "U", "U", "U", "U", "U", // U+0168 ŨũŪūŬŭŮů
	"U", "U", "U", "U", "W", "W", "Y", "Y", // U+0170 ŰűŲųŴŵŶŷ
	"Y", "Z", "Z", "Z", "Z", "Z", "Z", "S", // U+0178 ŸŹźŻżŽžſ
}

// appendLetters appends the letters of s to out as uppercase ASCII, folding
// Latin diacritics. Everything else is dropped, except that with keepSpaces a
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
		case r >= latinFoldFirst && int(r-latinFoldFirst) < len(latinFold):
			fold = latinFold[r-latinFoldFirst]
		case r == '\'' || r == '’':
			continue
		case r >= 0x300 && unicode.Is(unicode.Mn, r):
			// Combining accents in decomposed text ("e\u0301") belong to
			// the previous letter.
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
