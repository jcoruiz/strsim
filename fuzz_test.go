package strsim

import (
	"math"
	"strings"
	"testing"
	"unicode/utf8"
)

// similarity turns a method such as NgramMetric.Cosine into a Metric.
type similarity func(a, b string) float64

func (f similarity) Similarity(a, b string) float64 { return f(a, b) }

func addFuzzPairs(f *testing.F) {
	for _, p := range [][2]string{
		{"", ""},
		{"", "abc"},
		{"kitten", "sitting"},
		{"martha", "marhta"},
		{"cabb", "abacbcbba"},
		{"日本語", "日本"},
		{"😀a", "a😀"},
		{"Schmidt", "Smith"},
		{"\xff", "\xfe"},
	} {
		f.Add(p[0], p[1])
	}
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			return false
		}
	}

	return true
}

// FuzzMetrics checks the invariants every similarity metric promises.
func FuzzMetrics(f *testing.F) {
	addFuzzPairs(f)

	symmetric := map[string]Metric{
		"Levenshtein":        NewLevenshtein(),
		"OSA":                NewOSA(),
		"DamerauLevenshtein": NewDamerauLevenshtein(),
		"LCS":                NewLCS(),
		"Hamming":            NewHamming(),
		"Ngram":              NewNgram(),
		"Trigram":            &NgramMetric{Size: 3},
		"Cosine":             similarity(NewNgram().Cosine),
		"TrigramCosine":      similarity((&NgramMetric{Size: 3}).Cosine),
		"Jaccard":            similarity(NewNgram().Jaccard),
		"Overlap":            similarity(NewNgram().Overlap),
	}
	other := map[string]Metric{
		"LevenshteinCosts": &LevenshteinMetric{InsertCost: 2, DeleteCost: 3, ReplaceCost: 5},
		"Jaro":             NewJaro(),
		"JaroWinkler":      NewJaroWinkler(),
		"JaroWinklerLong":  &JaroWinklerMetric{BoostThreshold: 0.7, PrefixSize: 20},
	}

	f.Fuzz(func(t *testing.T, a, b string) {
		check := func(name string, m Metric) {
			s := m.Similarity(a, b)
			if math.IsNaN(s) || s < 0 || s > 1 {
				t.Errorf("%s.Similarity(%q, %q) = %v, want a value in [0, 1]", name, a, b, s)
			}

			if utf8.ValidString(a) {
				if self := m.Similarity(a, a); self != 1 {
					t.Errorf("%s.Similarity(%q, %q) = %v, want 1", name, a, a, self)
				}
			}
		}

		for name, m := range symmetric {
			check(name, m)

			if ab, ba := m.Similarity(a, b), m.Similarity(b, a); math.Abs(ab-ba) > 1e-12 {
				t.Errorf("%s is not symmetric for (%q, %q): %v vs %v", name, a, b, ab, ba)
			}
		}

		for name, m := range other {
			check(name, m)
		}

		if utf8.ValidString(a) && utf8.ValidString(b) {
			maxLen := max(utf8.RuneCountInString(a), utf8.RuneCountInString(b))
			for _, m := range []DistanceMetric{NewLevenshtein(), NewOSA(), NewDamerauLevenshtein(), NewLCS()} {
				d := m.Distance(a, b)
				if (d == 0) != (a == b) {
					t.Errorf("%T.Distance(%q, %q) = %d", m, a, b, d)
				}
				if _, isLCS := m.(*LCSMetric); !isLCS && d > maxLen {
					t.Errorf("%T.Distance(%q, %q) = %d, exceeds the longer length %d", m, a, b, d, maxLen)
				}
			}
		}
	})
}

// FuzzASCIIOnly checks that the byte-level fast paths agree with the rune
// paths on ASCII input.
func FuzzASCIIOnly(f *testing.F) {
	addFuzzPairs(f)

	f.Fuzz(func(t *testing.T, a, b string) {
		if !isASCII(a) || !isASCII(b) {
			return
		}

		pairs := []struct {
			name      string
			std, fast Metric
		}{
			{"Levenshtein", &LevenshteinMetric{InsertCost: 2, DeleteCost: 3, ReplaceCost: 5},
				&LevenshteinMetric{InsertCost: 2, DeleteCost: 3, ReplaceCost: 5, ASCIIOnly: true}},
			{"OSA", NewOSA(), &OSAMetric{ASCIIOnly: true}},
			{"DamerauLevenshtein", NewDamerauLevenshtein(), &DamerauLevenshteinMetric{ASCIIOnly: true}},
			{"LCS", NewLCS(), &LCSMetric{ASCIIOnly: true}},
			{"Hamming", NewHamming(), &HammingMetric{ASCIIOnly: true}},
			{"Jaro", NewJaro(), &JaroMetric{ASCIIOnly: true}},
			{"JaroWinkler", NewJaroWinkler(), &JaroWinklerMetric{ASCIIOnly: true}},
		}

		for _, p := range pairs {
			if std, fast := p.std.Similarity(a, b), p.fast.Similarity(a, b); std != fast {
				t.Errorf("%s.Similarity(%q, %q): runes %v, ASCII %v", p.name, a, b, std, fast)
			}

			dm, ok := p.std.(DistanceMetric)
			if !ok || (p.name == "Hamming" && len(a) != len(b)) {
				continue
			}

			if std, fast := dm.Distance(a, b), p.fast.(DistanceMetric).Distance(a, b); std != fast {
				t.Errorf("%s.Distance(%q, %q): runes %d, ASCII %d", p.name, a, b, std, fast)
			}
		}
	})
}

// FuzzPhonetic checks that encoders never panic and only emit well-formed
// codes.
func FuzzPhonetic(f *testing.F) {
	for _, s := range []string{"", "Schmidt", "O'Brien", "San Jose", "Muñoz", "Garci\u0301a", "Strauß", "東京", "Mac Gregor", "\xff"} {
		f.Add(s)
	}

	wellFormed := func(code, alphabet string, maxLen int) bool {
		if len(code) > maxLen {
			return false
		}

		for i := 0; i < len(code); i++ {
			if strings.IndexByte(alphabet, code[i]) < 0 {
				return false
			}
		}

		return true
	}

	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	f.Fuzz(func(t *testing.T, s string) {
		if c := Soundex(s); c != "" && (len(c) != 4 || !wellFormed(c[:1], letters, 1) || !wellFormed(c[1:], "0123456", 3)) {
			t.Errorf("Soundex(%q) = %q", s, c)
		}

		if c := Metaphone(s); !wellFormed(c, letters+"0", 4) {
			t.Errorf("Metaphone(%q) = %q", s, c)
		}

		if p, a := DoubleMetaphone(s); !wellFormed(p, letters+"0", 4) || !wellFormed(a, letters+"0", 4) {
			t.Errorf("DoubleMetaphone(%q) = (%q, %q)", s, p, a)
		}

		if c := NYSIIS(s); !wellFormed(c, letters, 6) {
			t.Errorf("NYSIIS(%q) = %q", s, c)
		}

		if isASCII(s) {
			upper, lower := strings.ToUpper(s), strings.ToLower(s)
			if Soundex(upper) != Soundex(lower) || Metaphone(upper) != Metaphone(lower) || NYSIIS(upper) != NYSIIS(lower) {
				t.Errorf("encoding of %q depends on case", s)
			}
			p1, a1 := DoubleMetaphone(upper)
			if p2, a2 := DoubleMetaphone(lower); p1 != p2 || a1 != a2 {
				t.Errorf("DoubleMetaphone of %q depends on case", s)
			}
		}
	})
}
