package strsim

import (
	"math"
	"testing"
)

const jaroEpsilon = 1e-6

func almostEqual(a, b, eps float64) bool {
	return math.Abs(a-b) < eps
}

func TestJaroSimilarity(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want float64
	}{
		{name: "both empty", a: "", b: "", want: 1.0},
		{name: "first empty", a: "", b: "abc", want: 0.0},
		{name: "second empty", a: "abc", b: "", want: 0.0},
		{name: "identical", a: "abc", b: "abc", want: 1.0},
		{name: "martha/marhta", a: "martha", b: "marhta", want: 0.944444},
		{name: "dwayne/duane", a: "DWAYNE", b: "DUANE", want: 0.822222},
		{name: "dixon/dicksonx", a: "DIXON", b: "DICKSONX", want: 0.766667},
		{name: "jellyfish/smellyfish", a: "jellyfish", b: "smellyfish", want: 0.896296},
		{name: "single char match", a: "a", b: "a", want: 1.0},
		{name: "single char no match", a: "a", b: "b", want: 0.0},
		{name: "unicode japanese", a: "\u6771\u4eac", b: "\u6771\u4eac\u90fd", want: 0.888889},
		{name: "unicode identical", a: "\u3053\u3093\u306b\u3061\u306f", b: "\u3053\u3093\u306b\u3061\u306f", want: 1.0},
		{name: "emoji identical", a: "\U0001F600\U0001F601", b: "\U0001F600\U0001F601", want: 1.0},
		{name: "emoji different", a: "\U0001F600", b: "\U0001F601", want: 0.0},
		{name: "completely different", a: "abc", b: "xyz", want: 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := JaroSimilarity(tt.a, tt.b)
			if !almostEqual(got, tt.want, jaroEpsilon) {
				t.Errorf("JaroSimilarity(%q, %q) = %f, want %f", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestJaroMetric_Similarity(t *testing.T) {
	m := NewJaro()

	got := m.Similarity("martha", "marhta")
	want := JaroSimilarity("martha", "marhta")

	if got != want {
		t.Errorf("JaroMetric.Similarity() = %f, want %f", got, want)
	}
}

func TestJaroMetric_ImplementsMetric(t *testing.T) {
	var _ Metric = NewJaro()
}

func TestJaroWinklerSimilarity(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want float64
	}{
		{name: "both empty", a: "", b: "", want: 1.0},
		{name: "first empty", a: "", b: "abc", want: 0.0},
		{name: "second empty", a: "abc", b: "", want: 0.0},
		{name: "identical", a: "abc", b: "abc", want: 1.0},
		{name: "martha/marhta", a: "martha", b: "marhta", want: 0.961111},
		{name: "dwayne/duane", a: "DWAYNE", b: "DUANE", want: 0.840000},
		{name: "dixon/dicksonx", a: "DIXON", b: "DICKSONX", want: 0.813333},
		{name: "single char match", a: "a", b: "a", want: 1.0},
		{name: "single char no match", a: "a", b: "b", want: 0.0},
		{name: "unicode japanese", a: "\u6771\u4eac", b: "\u6771\u4eac\u90fd", want: 0.911111},
		{name: "emoji identical", a: "\U0001F600\U0001F601", b: "\U0001F600\U0001F601", want: 1.0},
		{name: "completely different", a: "abc", b: "xyz", want: 0.0},
		{name: "common prefix boost", a: "prefix1", b: "prefix2", want: 0.942857},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := JaroWinklerSimilarity(tt.a, tt.b)
			if !almostEqual(got, tt.want, jaroEpsilon) {
				t.Errorf("JaroWinklerSimilarity(%q, %q) = %f, want %f", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestJaroWinklerMetric_Similarity(t *testing.T) {
	m := NewJaroWinkler()

	got := m.Similarity("martha", "marhta")
	want := JaroWinklerSimilarity("martha", "marhta")

	if got != want {
		t.Errorf("JaroWinklerMetric.Similarity() = %f, want %f", got, want)
	}
}

func TestJaroWinklerMetric_ImplementsMetric(t *testing.T) {
	var _ Metric = NewJaroWinkler()
}

func TestJaroWinklerMetric_CustomSettings(t *testing.T) {
	m := &JaroWinklerMetric{
		BoostThreshold: 0.0,
		PrefixSize:     2,
	}

	// With threshold=0 every score gets boosted; PrefixSize=2 limits prefix.
	got := m.Similarity("abc", "axc")
	if got <= 0 {
		t.Errorf("expected positive similarity, got %f", got)
	}
}

func TestJaroWinklerSimilarity_BelowThreshold(t *testing.T) {
	// When Jaro score is below BoostThreshold, Jaro-Winkler should equal Jaro.
	a := "abc"
	b := "xyz"

	jaro := JaroSimilarity(a, b)
	jw := JaroWinklerSimilarity(a, b)

	if !almostEqual(jaro, jw, jaroEpsilon) {
		t.Errorf("below threshold: Jaro-Winkler (%f) should equal Jaro (%f)", jw, jaro)
	}
}

func BenchmarkJaroSimilarity(b *testing.B) {
	for i := 0; i < b.N; i++ {
		JaroSimilarity("jellyfish", "smellyfish")
	}
}

func BenchmarkJaroWinklerSimilarity(b *testing.B) {
	for i := 0; i < b.N; i++ {
		JaroWinklerSimilarity("martha", "marhta")
	}
}
