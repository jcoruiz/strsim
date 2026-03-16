package strsim

import (
	"math"
	"testing"
)

func TestOSA(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want int
	}{
		{name: "identical", a: "hello", b: "hello", want: 0},
		{name: "both_empty", a: "", b: "", want: 0},
		{name: "one_empty", a: "abc", b: "", want: 3},
		{name: "other_empty", a: "", b: "abc", want: 3},
		{name: "single_char_same", a: "a", b: "a", want: 0},
		{name: "single_char_diff", a: "a", b: "b", want: 1},
		{name: "transposition", a: "ab", b: "ba", want: 1},
		{name: "ca_abc", a: "CA", b: "ABC", want: 3},
		{name: "kitten_sitting", a: "kitten", b: "sitting", want: 3},
		{name: "adjacent_transposition", a: "abcd", b: "abdc", want: 1},
		{name: "unicode_japanese", a: "\u3053\u3093\u306b\u3061\u306f", b: "\u3053\u3093\u306b\u3061\u306f", want: 0},
		{name: "unicode_japanese_transposition", a: "\u3053\u3093", b: "\u3093\u3053", want: 1},
		{name: "emoji", a: "\U0001f600\U0001f601", b: "\U0001f601\U0001f600", want: 1},
		{name: "long_strings", a: "abcdefghij", b: "abdcefghij", want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := OSA(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("OSA(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestOSASimilarity(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want float64
	}{
		{name: "identical", a: "hello", b: "hello", want: 1.0},
		{name: "both_empty", a: "", b: "", want: 1.0},
		{name: "transposition", a: "ab", b: "ba", want: 0.5},
		{name: "one_empty", a: "abc", b: "", want: 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := OSASimilarity(tt.a, tt.b)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("OSASimilarity(%q, %q) = %f, want %f", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestDamerauLevenshtein(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want int
	}{
		{name: "identical", a: "hello", b: "hello", want: 0},
		{name: "both_empty", a: "", b: "", want: 0},
		{name: "one_empty", a: "abc", b: "", want: 3},
		{name: "other_empty", a: "", b: "abc", want: 3},
		{name: "single_char_same", a: "a", b: "a", want: 0},
		{name: "single_char_diff", a: "a", b: "b", want: 1},
		{name: "transposition", a: "ab", b: "ba", want: 1},
		{name: "ca_abc", a: "CA", b: "ABC", want: 2},
		{name: "kitten_sitting", a: "kitten", b: "sitting", want: 3},
		{name: "adjacent_transposition", a: "abcd", b: "abdc", want: 1},
		// The key difference between OSA and true DL:
		// OSA("CA", "ABC") = 3, but DL("CA", "ABC") = 2
		{name: "unicode_japanese", a: "\u3053\u3093\u306b\u3061\u306f", b: "\u3053\u3093\u306b\u3061\u306f", want: 0},
		{name: "unicode_japanese_transposition", a: "\u3053\u3093", b: "\u3093\u3053", want: 1},
		{name: "emoji", a: "\U0001f600\U0001f601", b: "\U0001f601\U0001f600", want: 1},
		{name: "long_strings", a: "abcdefghij", b: "abdcefghij", want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DamerauLevenshtein(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("DamerauLevenshtein(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestDamerauLevenshteinSimilarity(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want float64
	}{
		{name: "identical", a: "hello", b: "hello", want: 1.0},
		{name: "both_empty", a: "", b: "", want: 1.0},
		{name: "transposition", a: "ab", b: "ba", want: 0.5},
		{name: "one_empty", a: "abc", b: "", want: 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DamerauLevenshteinSimilarity(tt.a, tt.b)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("DamerauLevenshteinSimilarity(%q, %q) = %f, want %f", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestOSA_vs_DamerauLevenshtein(t *testing.T) {
	// Demonstrate the difference between OSA and true Damerau-Levenshtein.
	// For "CA" -> "ABC":
	// OSA = 3 (restricted: each substring edited at most once)
	// DL  = 2 (unrestricted transpositions)
	t.Run("ca_abc_difference", func(t *testing.T) {
		osa := OSA("CA", "ABC")
		dl := DamerauLevenshtein("CA", "ABC")

		if osa != 3 {
			t.Errorf("OSA(CA, ABC) = %d, want 3", osa)
		}
		if dl != 2 {
			t.Errorf("DamerauLevenshtein(CA, ABC) = %d, want 2", dl)
		}
		if dl > osa {
			t.Error("DL distance should be <= OSA distance")
		}
	})
}

func TestOSAMetric_Interfaces(t *testing.T) {
	var _ Metric = NewOSA()
	var _ DistanceMetric = NewOSA()
}

func TestDamerauLevenshteinMetric_Interfaces(t *testing.T) {
	var _ Metric = NewDamerauLevenshtein()
	var _ DistanceMetric = NewDamerauLevenshtein()
}

func BenchmarkOSA(b *testing.B) {
	m := NewOSA()
	a := "abcdefghijklmnopqrstuvwxyz"
	s := "zyxwvutsrqponmlkjihgfedcba"

	b.ResetTimer()
	for range b.N {
		m.Distance(a, s)
	}
}

func BenchmarkDamerauLevenshtein(b *testing.B) {
	m := NewDamerauLevenshtein()
	a := "abcdefghijklmnopqrstuvwxyz"
	s := "zyxwvutsrqponmlkjihgfedcba"

	b.ResetTimer()
	for range b.N {
		m.Distance(a, s)
	}
}
