package strsim

import (
	"math"
	"testing"
)

func TestLevenshtein(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want int
	}{
		{name: "kitten_sitting", a: "kitten", b: "sitting", want: 3},
		{name: "identical", a: "hello", b: "hello", want: 0},
		{name: "both_empty", a: "", b: "", want: 0},
		{name: "one_empty", a: "abc", b: "", want: 3},
		{name: "other_empty", a: "", b: "abc", want: 3},
		{name: "single_char_same", a: "a", b: "a", want: 0},
		{name: "single_char_diff", a: "a", b: "b", want: 1},
		{name: "insertion", a: "abc", b: "abcd", want: 1},
		{name: "deletion", a: "abcd", b: "abc", want: 1},
		{name: "substitution", a: "abc", b: "axc", want: 1},
		{name: "completely_different", a: "abc", b: "xyz", want: 3},
		{name: "unicode_japanese", a: "\u3053\u3093\u306b\u3061\u306f", b: "\u3053\u3093\u3070\u3093\u306f", want: 2},
		{name: "emoji", a: "\U0001f600\U0001f601\U0001f602", b: "\U0001f600\U0001f603", want: 2},
		{name: "long_strings", a: "abcdefghijklmnopqrstuvwxyz", b: "abcdefghijklmnopqrstuvwxy", want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Levenshtein(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Levenshtein(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestLevenshteinSimilarity(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want float64
	}{
		{name: "identical", a: "hello", b: "hello", want: 1.0},
		{name: "both_empty", a: "", b: "", want: 1.0},
		{name: "completely_different", a: "abc", b: "xyz", want: 0.0},
		{name: "one_empty", a: "abc", b: "", want: 0.0},
		{name: "kitten_sitting", a: "kitten", b: "sitting", want: 1.0 - 3.0/7.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LevenshteinSimilarity(tt.a, tt.b)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("LevenshteinSimilarity(%q, %q) = %f, want %f", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestLevenshteinMetric_CustomCosts(t *testing.T) {
	t.Run("insert_cost_2", func(t *testing.T) {
		m := &LevenshteinMetric{InsertCost: 2, DeleteCost: 1, ReplaceCost: 1}
		// "abc" -> "abcd" requires 1 insert, cost = 2
		got := m.Distance("abc", "abcd")
		if got != 2 {
			t.Errorf("Distance = %d, want 2", got)
		}
	})

	t.Run("delete_cost_2", func(t *testing.T) {
		m := &LevenshteinMetric{InsertCost: 1, DeleteCost: 2, ReplaceCost: 1}
		// "abcd" -> "abc" requires 1 delete, cost = 2
		got := m.Distance("abcd", "abc")
		if got != 2 {
			t.Errorf("Distance = %d, want 2", got)
		}
	})

	t.Run("replace_cost_2", func(t *testing.T) {
		m := &LevenshteinMetric{InsertCost: 1, DeleteCost: 1, ReplaceCost: 2}
		// "abc" -> "axc" requires 1 replace at cost 2, but insert+delete = 2 as well.
		got := m.Distance("abc", "axc")
		if got != 2 {
			t.Errorf("Distance = %d, want 2", got)
		}
	})

	t.Run("similarity_custom_costs_normalized", func(t *testing.T) {
		m := &LevenshteinMetric{InsertCost: 2, DeleteCost: 2, ReplaceCost: 2}
		// Same as standard but doubled costs. Similarity should be the same
		// after normalization since maxCost also doubles.
		got := m.Similarity("kitten", "sitting")
		want := 1.0 - float64(6)/float64(7*2)
		if math.Abs(got-want) > 1e-9 {
			t.Errorf("Similarity = %f, want %f", got, want)
		}
	})
}

func TestLevenshteinMetric_Interfaces(t *testing.T) {
	var _ Metric = NewLevenshtein()
	var _ DistanceMetric = NewLevenshtein()
}

func BenchmarkLevenshtein(b *testing.B) {
	m := NewLevenshtein()
	a := "abcdefghijklmnopqrstuvwxyz"
	s := "zyxwvutsrqponmlkjihgfedcba"

	b.ResetTimer()
	for range b.N {
		m.Distance(a, s)
	}
}

func BenchmarkLevenshteinSimilarity(b *testing.B) {
	m := NewLevenshtein()
	a := "abcdefghijklmnopqrstuvwxyz"
	s := "zyxwvutsrqponmlkjihgfedcba"

	b.ResetTimer()
	for range b.N {
		m.Similarity(a, s)
	}
}
