package strsim

import (
	"math"
	"testing"
)

func TestHamming(t *testing.T) {
	tests := []struct {
		name    string
		a, b    string
		want    int
		wantErr bool
	}{
		{name: "identical", a: "karolin", b: "karolin", want: 0},
		{name: "one_diff", a: "karolin", b: "kathrin", want: 3},
		{name: "all_diff", a: "abc", b: "xyz", want: 3},
		{name: "both_empty", a: "", b: "", want: 0},
		{name: "single_same", a: "a", b: "a", want: 0},
		{name: "single_diff", a: "a", b: "b", want: 1},
		{name: "different_lengths", a: "abc", b: "ab", wantErr: true},
		{name: "unicode_japanese", a: "\u3053\u3093\u306b\u3061\u306f", b: "\u3053\u3093\u306b\u3061\u306f", want: 0},
		{name: "unicode_japanese_diff", a: "\u3053\u3093\u306b\u3061\u306f", b: "\u3053\u3093\u3070\u3093\u306f", want: 2},
		{name: "emoji_same", a: "\U0001f600\U0001f601\U0001f602", b: "\U0001f600\U0001f601\U0001f602", want: 0},
		{name: "emoji_diff", a: "\U0001f600\U0001f601\U0001f602", b: "\U0001f600\U0001f603\U0001f604", want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Hamming(tt.a, tt.b)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Hamming(%q, %q) expected error, got nil", tt.a, tt.b)
				}
				return
			}
			if err != nil {
				t.Errorf("Hamming(%q, %q) unexpected error: %v", tt.a, tt.b, err)
				return
			}
			if got != tt.want {
				t.Errorf("Hamming(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestHammingMetric_Distance(t *testing.T) {
	m := NewHamming()

	t.Run("normal", func(t *testing.T) {
		got := m.Distance("karolin", "kathrin")
		if got != 3 {
			t.Errorf("Distance = %d, want 3", got)
		}
	})

	t.Run("panics_on_different_lengths", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Distance did not panic for different length strings")
			}
		}()
		m.Distance("abc", "ab")
	})
}

func TestHammingMetric_Similarity(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want float64
	}{
		{name: "identical", a: "karolin", b: "karolin", want: 1.0},
		{name: "one_diff_of_three", a: "abc", b: "axc", want: 2.0 / 3.0},
		{name: "all_diff", a: "abc", b: "xyz", want: 0.0},
		{name: "both_empty", a: "", b: "", want: 1.0},
		{name: "different_lengths", a: "abc", b: "ab", want: 0.0},
		{name: "unicode_japanese", a: "\u3053\u3093\u306b\u3061\u306f", b: "\u3053\u3093\u306b\u3061\u306f", want: 1.0},
		{name: "emoji", a: "\U0001f600\U0001f601", b: "\U0001f600\U0001f602", want: 0.5},
	}

	m := NewHamming()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := m.Similarity(tt.a, tt.b)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Similarity(%q, %q) = %f, want %f", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestHammingMetric_Interfaces(t *testing.T) {
	var _ Metric = NewHamming()
	var _ DistanceMetric = NewHamming()
}

func BenchmarkHamming(b *testing.B) {
	m := NewHamming()
	a := "abcdefghijklmnopqrstuvwxyz"
	s := "zyxwvutsrqponmlkjihgfedcba"

	b.ResetTimer()
	for range b.N {
		m.Distance(a, s)
	}
}

func BenchmarkHammingSimilarity(b *testing.B) {
	m := NewHamming()
	a := "abcdefghijklmnopqrstuvwxyz"
	s := "zyxwvutsrqponmlkjihgfedcba"

	b.ResetTimer()
	for range b.N {
		m.Similarity(a, s)
	}
}
