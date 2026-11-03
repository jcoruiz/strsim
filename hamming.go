package strsim

import (
	"errors"
	"unicode/utf8"
)

// ErrLengthMismatch is returned by Hamming when the strings differ in length.
var ErrLengthMismatch = errors.New("strsim: hamming distance requires strings of equal length")

// HammingMetric computes the Hamming distance between two strings of equal
// length. Hamming distance is the number of positions at which the
// corresponding characters differ.
//
// Set ASCIIOnly to true for faster processing of pure-ASCII strings. In this
// mode, length is measured in bytes instead of runes.
type HammingMetric struct {
	// ASCIIOnly skips rune conversion for faster ASCII processing.
	// In this mode, length is measured in bytes instead of runes.
	ASCIIOnly bool
}

// NewHamming returns a new HammingMetric instance.
func NewHamming() *HammingMetric {
	return &HammingMetric{}
}

// Hamming returns the Hamming distance between two strings. It returns
// ErrLengthMismatch if the strings have different rune lengths, since Hamming
// distance is only defined for strings of equal length.
func Hamming(a, b string) (int, error) {
	ra := []rune(a)
	rb := []rune(b)

	if len(ra) != len(rb) {
		return 0, ErrLengthMismatch
	}

	return hammingRunes(ra, rb), nil
}

// hammingRunes counts the positions where two equal-length rune slices differ.
func hammingRunes(ra, rb []rune) int {
	d := 0
	for i := range ra {
		if ra[i] != rb[i] {
			d++
		}
	}

	return d
}

// HammingSimilarity returns the normalized Hamming similarity between a and b
// as a value in [0, 1]. Returns 1.0 when both strings are empty. Returns 0.0
// when the strings have different rune lengths.
func HammingSimilarity(a, b string) float64 {
	return NewHamming().Similarity(a, b)
}

// Distance returns the Hamming distance between a and b. It panics if the
// strings have different lengths (rune length, or byte length if ASCIIOnly);
// use Hamming to get an error instead.
func (m *HammingMetric) Distance(a, b string) int {
	if m.ASCIIOnly {
		if len(a) != len(b) {
			panic(ErrLengthMismatch)
		}
		return hammingASCII(a, b)
	}

	d, err := Hamming(a, b)
	if err != nil {
		panic(err)
	}

	return d
}

// Similarity returns the normalized similarity between a and b as a value
// in [0, 1]. Returns 1.0 when both strings are empty. Returns 0.0 when the
// strings have different lengths.
func (m *HammingMetric) Similarity(a, b string) float64 {
	if m.ASCIIOnly {
		if len(a) != len(b) {
			return 0.0
		}

		n := len(a)
		if n == 0 {
			return 1.0
		}

		d := hammingASCII(a, b)

		return 1.0 - float64(d)/float64(n)
	}

	ra := []rune(a)
	rb := []rune(b)

	if len(ra) != len(rb) {
		return 0.0
	}

	n := len(ra)
	if n == 0 {
		return 1.0
	}

	return 1.0 - float64(hammingRunes(ra, rb))/float64(n)
}

// runeLen returns the number of runes in a string. This is a package-level
// helper to avoid repeated utf8.RuneCountInString calls.
func runeLen(s string) int {
	return utf8.RuneCountInString(s)
}
