package strsim

import (
	"testing"
)

// TestASCIIOnly verifies that ASCIIOnly mode produces the same results as the
// standard rune-based mode for pure-ASCII input.
func TestASCIIOnly(t *testing.T) {
	pairs := []struct {
		a, b string
	}{
		{"kitten", "sitting"},
		{"", ""},
		{"", "hello"},
		{"a", "b"},
		{"the quick brown fox", "the quikc brown box"},
		{"identical", "identical"},
		{"abc", "cab"},
	}

	t.Run("Levenshtein", func(t *testing.T) {
		for _, p := range pairs {
			std := NewLevenshtein().Distance(p.a, p.b)
			ascii := (&LevenshteinMetric{InsertCost: 1, DeleteCost: 1, ReplaceCost: 1, ASCIIOnly: true}).Distance(p.a, p.b)
			if std != ascii {
				t.Errorf("Levenshtein(%q, %q): std=%d, ascii=%d", p.a, p.b, std, ascii)
			}
		}
	})

	t.Run("LevenshteinSimilarity", func(t *testing.T) {
		for _, p := range pairs {
			std := NewLevenshtein().Similarity(p.a, p.b)
			ascii := (&LevenshteinMetric{InsertCost: 1, DeleteCost: 1, ReplaceCost: 1, ASCIIOnly: true}).Similarity(p.a, p.b)
			if std != ascii {
				t.Errorf("LevenshteinSimilarity(%q, %q): std=%f, ascii=%f", p.a, p.b, std, ascii)
			}
		}
	})

	t.Run("OSA", func(t *testing.T) {
		for _, p := range pairs {
			std := NewOSA().Distance(p.a, p.b)
			ascii := (&OSAMetric{ASCIIOnly: true}).Distance(p.a, p.b)
			if std != ascii {
				t.Errorf("OSA(%q, %q): std=%d, ascii=%d", p.a, p.b, std, ascii)
			}
		}
	})

	t.Run("DamerauLevenshtein", func(t *testing.T) {
		for _, p := range pairs {
			std := NewDamerauLevenshtein().Distance(p.a, p.b)
			ascii := (&DamerauLevenshteinMetric{ASCIIOnly: true}).Distance(p.a, p.b)
			if std != ascii {
				t.Errorf("DamerauLevenshtein(%q, %q): std=%d, ascii=%d", p.a, p.b, std, ascii)
			}
		}
	})

	t.Run("LCS", func(t *testing.T) {
		for _, p := range pairs {
			std := NewLCS().Distance(p.a, p.b)
			ascii := (&LCSMetric{ASCIIOnly: true}).Distance(p.a, p.b)
			if std != ascii {
				t.Errorf("LCS(%q, %q): std=%d, ascii=%d", p.a, p.b, std, ascii)
			}
		}
	})

	t.Run("LCSSimilarity", func(t *testing.T) {
		for _, p := range pairs {
			std := NewLCS().Similarity(p.a, p.b)
			ascii := (&LCSMetric{ASCIIOnly: true}).Similarity(p.a, p.b)
			if std != ascii {
				t.Errorf("LCSSimilarity(%q, %q): std=%f, ascii=%f", p.a, p.b, std, ascii)
			}
		}
	})

	t.Run("JaroWinkler", func(t *testing.T) {
		for _, p := range pairs {
			std := NewJaroWinkler().Similarity(p.a, p.b)
			ascii := (&JaroWinklerMetric{BoostThreshold: 0.7, PrefixSize: 4, ASCIIOnly: true}).Similarity(p.a, p.b)
			if std != ascii {
				t.Errorf("JaroWinkler(%q, %q): std=%f, ascii=%f", p.a, p.b, std, ascii)
			}
		}
	})

	t.Run("Jaro", func(t *testing.T) {
		for _, p := range pairs {
			std := NewJaro().Similarity(p.a, p.b)
			ascii := (&JaroMetric{ASCIIOnly: true}).Similarity(p.a, p.b)
			if std != ascii {
				t.Errorf("Jaro(%q, %q): std=%f, ascii=%f", p.a, p.b, std, ascii)
			}
		}
	})

	equalPairs := []struct {
		a, b string
	}{
		{"kitten", "sittin"},
		{"abcdef", "abcdef"},
	}

	t.Run("Hamming", func(t *testing.T) {
		for _, p := range equalPairs {
			std := NewHamming().Distance(p.a, p.b)
			ascii := (&HammingMetric{ASCIIOnly: true}).Distance(p.a, p.b)
			if std != ascii {
				t.Errorf("Hamming(%q, %q): std=%d, ascii=%d", p.a, p.b, std, ascii)
			}
		}
	})

	t.Run("HammingSimilarity", func(t *testing.T) {
		for _, p := range equalPairs {
			std := NewHamming().Similarity(p.a, p.b)
			ascii := (&HammingMetric{ASCIIOnly: true}).Similarity(p.a, p.b)
			if std != ascii {
				t.Errorf("HammingSimilarity(%q, %q): std=%f, ascii=%f", p.a, p.b, std, ascii)
			}
		}
	})

	t.Run("HammingASCII_empty", func(t *testing.T) {
		m := &HammingMetric{ASCIIOnly: true}
		if s := m.Similarity("", ""); s != 1.0 {
			t.Errorf("Similarity(\"\", \"\") = %f, want 1.0", s)
		}
		if s := m.Similarity("a", "ab"); s != 0.0 {
			t.Errorf("Similarity(\"a\", \"ab\") = %f, want 0.0", s)
		}
	})
}

// TestASCIIOnlyInterfaces verifies that ASCIIOnly metrics satisfy interfaces.
func TestASCIIOnlyInterfaces(t *testing.T) {
	var _ DistanceMetric = &LevenshteinMetric{ASCIIOnly: true}
	var _ DistanceMetric = &OSAMetric{ASCIIOnly: true}
	var _ DistanceMetric = &DamerauLevenshteinMetric{ASCIIOnly: true}
	var _ DistanceMetric = &LCSMetric{ASCIIOnly: true}
	var _ DistanceMetric = &HammingMetric{ASCIIOnly: true}
	var _ Metric = &JaroMetric{ASCIIOnly: true}
	var _ Metric = &JaroWinklerMetric{BoostThreshold: 0.7, PrefixSize: 4, ASCIIOnly: true}
}

// Benchmarks comparing rune-based vs ASCII-only on medium-length ASCII input.
var benchA = "the quick brown fox jumps over the lazy dog"
var benchB = "the quikc brown box jumped over the layz dog"

func BenchmarkLevenshteinRune(b *testing.B) {
	m := NewLevenshtein()
	for range b.N {
		m.Distance(benchA, benchB)
	}
}

func BenchmarkLevenshteinASCII(b *testing.B) {
	m := &LevenshteinMetric{InsertCost: 1, DeleteCost: 1, ReplaceCost: 1, ASCIIOnly: true}
	for range b.N {
		m.Distance(benchA, benchB)
	}
}

func BenchmarkJaroWinklerRune(b *testing.B) {
	m := NewJaroWinkler()
	for range b.N {
		m.Similarity(benchA, benchB)
	}
}

func BenchmarkJaroWinklerASCII(b *testing.B) {
	m := &JaroWinklerMetric{BoostThreshold: 0.7, PrefixSize: 4, ASCIIOnly: true}
	for range b.N {
		m.Similarity(benchA, benchB)
	}
}

func BenchmarkOSARune(b *testing.B) {
	m := NewOSA()
	for range b.N {
		m.Distance(benchA, benchB)
	}
}

func BenchmarkOSAASCII(b *testing.B) {
	m := &OSAMetric{ASCIIOnly: true}
	for range b.N {
		m.Distance(benchA, benchB)
	}
}

func BenchmarkDamerauLevenshteinRune(b *testing.B) {
	m := NewDamerauLevenshtein()
	for range b.N {
		m.Distance(benchA, benchB)
	}
}

func BenchmarkDamerauLevenshteinASCII(b *testing.B) {
	m := &DamerauLevenshteinMetric{ASCIIOnly: true}
	for range b.N {
		m.Distance(benchA, benchB)
	}
}

func BenchmarkLCSRune(b *testing.B) {
	m := NewLCS()
	for range b.N {
		m.Distance(benchA, benchB)
	}
}

func BenchmarkLCSASCII(b *testing.B) {
	m := &LCSMetric{ASCIIOnly: true}
	for range b.N {
		m.Distance(benchA, benchB)
	}
}

func BenchmarkHammingRune(b *testing.B) {
	a := "the quick brown fox jumps over the lazy"
	bStr := "the quikc brown box jumpe over the lazy"
	m := NewHamming()
	for range b.N {
		m.Distance(a, bStr)
	}
}

func BenchmarkHammingASCII(b *testing.B) {
	a := "the quick brown fox jumps over the lazy"
	bStr := "the quikc brown box jumpe over the lazy"
	m := &HammingMetric{ASCIIOnly: true}
	for range b.N {
		m.Distance(a, bStr)
	}
}
