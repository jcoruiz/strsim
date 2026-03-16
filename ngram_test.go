package strsim

import (
	"math"
	"testing"
)

const ngramEpsilon = 1e-6

func ngramAlmostEqual(a, b, eps float64) bool {
	return math.Abs(a-b) < eps
}

func TestNgramMetric_Cosine(t *testing.T) {
	m := NewNgram()

	tests := []struct {
		name string
		a    string
		b    string
		want float64
	}{
		{name: "both empty", a: "", b: "", want: 1.0},
		{name: "first empty", a: "", b: "abc", want: 0.0},
		{name: "second empty", a: "abc", b: "", want: 0.0},
		{name: "identical", a: "night", b: "night", want: 1.0},
		{name: "night/nacht", a: "night", b: "nacht", want: 0.25},
		{name: "single char a", a: "a", b: "a", want: 1.0},
		{name: "single char different", a: "a", b: "b", want: 0.0},
		{name: "repeated bigrams", a: "aaa", b: "aaa", want: 1.0},
		{name: "unicode japanese", a: "\u6771\u4eac", b: "\u6771\u4eac\u90fd", want: 0.707107},
		{name: "emoji identical", a: "\U0001F600\U0001F601", b: "\U0001F600\U0001F601", want: 1.0},
		{name: "completely different", a: "abc", b: "xyz", want: 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := m.Cosine(tt.a, tt.b)
			if !ngramAlmostEqual(got, tt.want, ngramEpsilon) {
				t.Errorf("Cosine(%q, %q) = %f, want %f", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestNgramMetric_Jaccard(t *testing.T) {
	m := NewNgram()

	tests := []struct {
		name string
		a    string
		b    string
		want float64
	}{
		{name: "both empty", a: "", b: "", want: 1.0},
		{name: "first empty", a: "", b: "abc", want: 0.0},
		{name: "second empty", a: "abc", b: "", want: 0.0},
		{name: "identical", a: "night", b: "night", want: 1.0},
		{name: "night/nacht", a: "night", b: "nacht", want: 0.142857},
		{name: "single char match", a: "a", b: "a", want: 1.0},
		{name: "single char different", a: "a", b: "b", want: 0.0},
		{name: "unicode japanese", a: "\u6771\u4eac", b: "\u6771\u4eac\u90fd", want: 0.5},
		{name: "emoji identical", a: "\U0001F600\U0001F601", b: "\U0001F600\U0001F601", want: 1.0},
		{name: "completely different", a: "abc", b: "xyz", want: 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := m.Jaccard(tt.a, tt.b)
			if !ngramAlmostEqual(got, tt.want, ngramEpsilon) {
				t.Errorf("Jaccard(%q, %q) = %f, want %f", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestNgramMetric_Dice(t *testing.T) {
	m := NewNgram()

	tests := []struct {
		name string
		a    string
		b    string
		want float64
	}{
		{name: "both empty", a: "", b: "", want: 1.0},
		{name: "first empty", a: "", b: "abc", want: 0.0},
		{name: "second empty", a: "abc", b: "", want: 0.0},
		{name: "identical", a: "night", b: "night", want: 1.0},
		{name: "night/nacht", a: "night", b: "nacht", want: 0.25},
		{name: "single char match", a: "a", b: "a", want: 1.0},
		{name: "single char different", a: "a", b: "b", want: 0.0},
		{name: "unicode japanese", a: "\u6771\u4eac", b: "\u6771\u4eac\u90fd", want: 0.666667},
		{name: "emoji identical", a: "\U0001F600\U0001F601", b: "\U0001F600\U0001F601", want: 1.0},
		{name: "completely different", a: "abc", b: "xyz", want: 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := m.Dice(tt.a, tt.b)
			if !ngramAlmostEqual(got, tt.want, ngramEpsilon) {
				t.Errorf("Dice(%q, %q) = %f, want %f", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestNgramMetric_Overlap(t *testing.T) {
	m := NewNgram()

	tests := []struct {
		name string
		a    string
		b    string
		want float64
	}{
		{name: "both empty", a: "", b: "", want: 1.0},
		{name: "first empty", a: "", b: "abc", want: 0.0},
		{name: "second empty", a: "abc", b: "", want: 0.0},
		{name: "identical", a: "night", b: "night", want: 1.0},
		{name: "night/nacht", a: "night", b: "nacht", want: 0.25},
		{name: "single char match", a: "a", b: "a", want: 1.0},
		{name: "single char different", a: "a", b: "b", want: 0.0},
		{name: "subset", a: "abc", b: "abcdef", want: 1.0},
		{name: "unicode japanese", a: "\u6771\u4eac", b: "\u6771\u4eac\u90fd", want: 1.0},
		{name: "emoji identical", a: "\U0001F600\U0001F601", b: "\U0001F600\U0001F601", want: 1.0},
		{name: "completely different", a: "abc", b: "xyz", want: 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := m.Overlap(tt.a, tt.b)
			if !ngramAlmostEqual(got, tt.want, ngramEpsilon) {
				t.Errorf("Overlap(%q, %q) = %f, want %f", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestNgramMetric_Similarity_UsesDice(t *testing.T) {
	m := NewNgram()

	a := "night"
	b := "nacht"

	got := m.Similarity(a, b)
	want := m.Dice(a, b)

	if got != want {
		t.Errorf("Similarity() = %f, should equal Dice() = %f", got, want)
	}
}

func TestNgramMetric_ImplementsMetric(t *testing.T) {
	var _ Metric = NewNgram()
}

func TestNgramMetric_CustomSize(t *testing.T) {
	m := &NgramMetric{Size: 3}

	// Trigrams of "abcde" = {"abc", "bcd", "cde"}
	// Trigrams of "abcfg" = {"abc", "bcf", "cfg"}
	// Intersection = {"abc"}, |A|=3, |B|=3
	// Dice = 2*1/6 = 0.333333
	got := m.Dice("abcde", "abcfg")
	want := 2.0 / 6.0

	if !ngramAlmostEqual(got, want, ngramEpsilon) {
		t.Errorf("Dice with trigrams = %f, want %f", got, want)
	}
}

func TestNgramMetric_NgramSizeGreaterThanString(t *testing.T) {
	// When string is shorter than ngram size, the entire string is one ngram.
	m := &NgramMetric{Size: 5}

	// "ab" < size 5 → treated as single ngram "ab".
	// Identical strings → cosine = 1.0, jaccard = 1.0, dice = 1.0, overlap = 1.0
	if got := m.Cosine("ab", "ab"); !ngramAlmostEqual(got, 1.0, ngramEpsilon) {
		t.Errorf("Cosine short identical = %f, want 1.0", got)
	}

	if got := m.Jaccard("ab", "ab"); !ngramAlmostEqual(got, 1.0, ngramEpsilon) {
		t.Errorf("Jaccard short identical = %f, want 1.0", got)
	}

	if got := m.Dice("ab", "ab"); !ngramAlmostEqual(got, 1.0, ngramEpsilon) {
		t.Errorf("Dice short identical = %f, want 1.0", got)
	}

	if got := m.Overlap("ab", "ab"); !ngramAlmostEqual(got, 1.0, ngramEpsilon) {
		t.Errorf("Overlap short identical = %f, want 1.0", got)
	}

	// Different short strings → 0.0
	if got := m.Cosine("ab", "cd"); !ngramAlmostEqual(got, 0.0, ngramEpsilon) {
		t.Errorf("Cosine short different = %f, want 0.0", got)
	}

	if got := m.Jaccard("ab", "cd"); !ngramAlmostEqual(got, 0.0, ngramEpsilon) {
		t.Errorf("Jaccard short different = %f, want 0.0", got)
	}

	if got := m.Dice("ab", "cd"); !ngramAlmostEqual(got, 0.0, ngramEpsilon) {
		t.Errorf("Dice short different = %f, want 0.0", got)
	}

	if got := m.Overlap("ab", "cd"); !ngramAlmostEqual(got, 0.0, ngramEpsilon) {
		t.Errorf("Overlap short different = %f, want 0.0", got)
	}
}

func TestCosineSimilarity_PackageLevel(t *testing.T) {
	got := CosineSimilarity("night", "nacht")
	want := NewNgram().Cosine("night", "nacht")

	if got != want {
		t.Errorf("CosineSimilarity() = %f, want %f", got, want)
	}
}

func TestJaccardSimilarity_PackageLevel(t *testing.T) {
	got := JaccardSimilarity("night", "nacht")
	want := NewNgram().Jaccard("night", "nacht")

	if got != want {
		t.Errorf("JaccardSimilarity() = %f, want %f", got, want)
	}
}

func TestDiceSimilarity_PackageLevel(t *testing.T) {
	got := DiceSimilarity("night", "nacht")
	want := NewNgram().Dice("night", "nacht")

	if got != want {
		t.Errorf("DiceSimilarity() = %f, want %f", got, want)
	}
}

func TestOverlapSimilarity_PackageLevel(t *testing.T) {
	got := OverlapSimilarity("night", "nacht")
	want := NewNgram().Overlap("night", "nacht")

	if got != want {
		t.Errorf("OverlapSimilarity() = %f, want %f", got, want)
	}
}

func TestCosine_TFVectors(t *testing.T) {
	// Verify that cosine uses term frequency, not binary.
	// "aaa" has bigrams: {"aa": 2}
	// "aa"  has bigrams: {"aa": 1}
	// dot = 2*1 = 2
	// normA = sqrt(4) = 2, normB = sqrt(1) = 1
	// cosine = 2/(2*1) = 1.0
	m := NewNgram()
	got := m.Cosine("aaa", "aa")

	if !ngramAlmostEqual(got, 1.0, ngramEpsilon) {
		t.Errorf("Cosine with TF vectors: got %f, want 1.0", got)
	}

	// "aaab" bigrams: {"aa": 2, "ab": 1}
	// "aab"  bigrams: {"aa": 1, "ab": 1}
	// dot = 2*1 + 1*1 = 3
	// normA = sqrt(4+1) = sqrt(5), normB = sqrt(1+1) = sqrt(2)
	// cosine = 3/sqrt(10)
	got2 := m.Cosine("aaab", "aab")
	want2 := 3.0 / math.Sqrt(10.0)

	if !ngramAlmostEqual(got2, want2, ngramEpsilon) {
		t.Errorf("Cosine TF: got %f, want %f", got2, want2)
	}
}

func BenchmarkCosineSimilarity(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CosineSimilarity("jellyfish", "smellyfish")
	}
}

func BenchmarkJaccardSimilarity(b *testing.B) {
	for i := 0; i < b.N; i++ {
		JaccardSimilarity("jellyfish", "smellyfish")
	}
}

func BenchmarkDiceSimilarity(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DiceSimilarity("jellyfish", "smellyfish")
	}
}

func BenchmarkOverlapSimilarity(b *testing.B) {
	for i := 0; i < b.N; i++ {
		OverlapSimilarity("jellyfish", "smellyfish")
	}
}
