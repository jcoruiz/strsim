package strsim

import (
	"math"
	"testing"
)

func TestLCS(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want int
	}{
		{name: "identical", a: "hello", b: "hello", want: 5},
		{name: "both_empty", a: "", b: "", want: 0},
		{name: "one_empty", a: "abc", b: "", want: 0},
		{name: "other_empty", a: "", b: "abc", want: 0},
		{name: "single_char_same", a: "a", b: "a", want: 1},
		{name: "single_char_diff", a: "a", b: "b", want: 0},
		{name: "abcde_ace", a: "abcde", b: "ace", want: 3},
		{name: "abc_def", a: "abc", b: "def", want: 0},
		{name: "subsequence", a: "AGGTAB", b: "GXTXAYB", want: 4},
		{name: "unicode_japanese", a: "\u3053\u3093\u306b\u3061\u306f", b: "\u3053\u3093\u3070\u3093\u306f", want: 3},
		{name: "emoji", a: "\U0001f600\U0001f601\U0001f602", b: "\U0001f600\U0001f603\U0001f602", want: 2},
		{name: "long_same_prefix", a: "abcdefghij", b: "abcdefghik", want: 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LCS(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("LCS(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestLCSDistance(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want int
	}{
		{name: "identical", a: "hello", b: "hello", want: 0},
		{name: "both_empty", a: "", b: "", want: 0},
		{name: "one_empty", a: "abc", b: "", want: 3},
		{name: "other_empty", a: "", b: "abc", want: 3},
		{name: "abcde_ace", a: "abcde", b: "ace", want: 2}, // 5 + 3 - 2*3 = 2
		{name: "abc_def", a: "abc", b: "def", want: 6},     // 3 + 3 - 0 = 6
		{name: "unicode_japanese", a: "\u3053\u3093\u306b\u3061\u306f", b: "\u3053\u3093\u3070\u3093\u306f", want: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LCSDistance(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("LCSDistance(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestLCSSimilarity(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want float64
	}{
		{name: "identical", a: "hello", b: "hello", want: 1.0},
		{name: "both_empty", a: "", b: "", want: 1.0},
		{name: "one_empty", a: "abc", b: "", want: 0.0},
		{name: "completely_different", a: "abc", b: "def", want: 0.0},
		{name: "abcde_ace", a: "abcde", b: "ace", want: 3.0 / 5.0},
		{name: "emoji", a: "\U0001f600\U0001f601\U0001f602", b: "\U0001f600\U0001f603\U0001f602", want: 2.0 / 3.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LCSSimilarity(tt.a, tt.b)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("LCSSimilarity(%q, %q) = %f, want %f", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestLCSMetric_Distance(t *testing.T) {
	m := NewLCS()

	got := m.Distance("abcde", "ace")
	if got != 2 {
		t.Errorf("Distance = %d, want 2", got)
	}
}

func TestLCSMetric_Similarity(t *testing.T) {
	m := NewLCS()

	got := m.Similarity("abcde", "ace")
	want := 3.0 / 5.0
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("Similarity = %f, want %f", got, want)
	}
}

func TestLCSMetric_Interfaces(t *testing.T) {
	var _ Metric = NewLCS()
	var _ DistanceMetric = NewLCS()
}

func BenchmarkLCS(b *testing.B) {
	a := "abcdefghijklmnopqrstuvwxyz"
	s := "zyxwvutsrqponmlkjihgfedcba"

	b.ResetTimer()
	for range b.N {
		LCS(a, s)
	}
}

func BenchmarkLCSDistance(b *testing.B) {
	m := NewLCS()
	a := "abcdefghijklmnopqrstuvwxyz"
	s := "zyxwvutsrqponmlkjihgfedcba"

	b.ResetTimer()
	for range b.N {
		m.Distance(a, s)
	}
}
