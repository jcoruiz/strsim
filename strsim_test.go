package strsim

import "testing"

// TestZeroValues checks the package promise that zero-value structs behave
// like their constructors.
func TestZeroValues(t *testing.T) {
	pairs := [][2]string{{"kitten", "sitting"}, {"martha", "marhta"}, {"", ""}, {"night", "nacht"}, {"日本語", "日本"}}

	metrics := []struct {
		name      string
		zero, def Metric
	}{
		{"Levenshtein", &LevenshteinMetric{}, NewLevenshtein()},
		{"OSA", &OSAMetric{}, NewOSA()},
		{"DamerauLevenshtein", &DamerauLevenshteinMetric{}, NewDamerauLevenshtein()},
		{"LCS", &LCSMetric{}, NewLCS()},
		{"Hamming", &HammingMetric{}, NewHamming()},
		{"Jaro", &JaroMetric{}, NewJaro()},
		{"JaroWinkler", &JaroWinklerMetric{}, NewJaroWinkler()},
		{"Ngram", &NgramMetric{}, NewNgram()},
	}

	for _, m := range metrics {
		for _, p := range pairs {
			if got, want := m.zero.Similarity(p[0], p[1]), m.def.Similarity(p[0], p[1]); got != want {
				t.Errorf("zero %s.Similarity(%q, %q) = %v, want %v", m.name, p[0], p[1], got, want)
			}
		}
	}

	encoders := []struct {
		name      string
		zero, def Encoder
	}{
		{"Soundex", &SoundexEncoder{}, NewSoundex()},
		{"Metaphone", &MetaphoneEncoder{}, NewMetaphone()},
	}

	for _, e := range encoders {
		for _, s := range []string{"Washington", "Schmidt", "Thompson"} {
			if got, want := e.zero.Encode(s), e.def.Encode(s); got != want {
				t.Errorf("zero %s.Encode(%q) = %q, want %q", e.name, s, got, want)
			}
		}
	}

	zeroDM, defDM := &DoubleMetaphoneEncoder{}, NewDoubleMetaphone()
	for _, s := range []string{"Washington", "Schmidt", "Thompson"} {
		p1, a1 := zeroDM.Encode(s)
		if p2, a2 := defDM.Encode(s); p1 != p2 || a1 != a2 {
			t.Errorf("zero DoubleMetaphone.Encode(%q) = (%q, %q), want (%q, %q)", s, p1, a1, p2, a2)
		}
	}
}
