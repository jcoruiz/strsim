package strsim

import "testing"

func TestNYSIIS(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "BROOKS", input: "BROOKS", want: "BRAC"},
		{name: "LAKE", input: "LAKE", want: "LAC"},
		{name: "LINCOLN", input: "LINCOLN", want: "LANCAL"},
		{name: "LYNCH", input: "LYNCH", want: "LYNC"},
		{name: "MACINTOSH", input: "MACINTOSH", want: "MCANT"},
		{name: "KNIGHT", input: "KNIGHT", want: "NAGT"},
		{name: "PHONE", input: "PHONE", want: "FAN"},
		{name: "SCHMIDT", input: "SCHMIDT", want: "SNAD"},
		{name: "empty", input: "", want: ""},
		{name: "non_alpha", input: "123", want: ""},
		{name: "single_letter", input: "A", want: "A"},
		{name: "lowercase", input: "brooks", want: "BRAC"},
		{name: "mixed_case", input: "Brooks", want: "BRAC"},
		{name: "unicode_stripped", input: "Brööks", want: "BRC"},
		{name: "non_latin_ignored", input: "Brooks東京", want: "BRAC"},

		// Prefix: MAC → MCC
		{name: "MacArthur", input: "MacArthur", want: "MCARTA"},
		// Prefix: PF → FF
		{name: "Pfeiffer", input: "Pfeiffer", want: "FAFAR"},
		// Prefix: SCH → SSS
		{name: "Schiller", input: "Schiller", want: "SALAR"},
		{name: "Schultz", input: "Schultz", want: "SALT"},
		{name: "Schwartz", input: "Schwartz", want: "SWART"},
		// Prefix: KN → NN
		{name: "Knuth", input: "Knuth", want: "NAT"},
		// Suffix: EE/IE → Y
		{name: "Lee", input: "Lee", want: "LY"},
		// Suffix: DT/RT/RD/NT/ND
		{name: "Dorothy", input: "Dorothy", want: "DARATY"},
		{name: "Barnard", input: "Barnard", want: "BARNAD"},
		{name: "Grant", input: "Grant", want: "GRAD"},
		{name: "Howard", input: "Howard", want: "HAD"},
		{name: "Barnhardt", input: "Barnhardt", want: "BARNAR"},
		{name: "Linda", input: "Linda", want: "LAND"},
		// EV → AF rule
		{name: "Evelyn", input: "Evelyn", want: "EVALYN"},
		{name: "Steve", input: "Steve", want: "STAF"},
		{name: "Devlin", input: "Devlin", want: "DAFLAN"},
		// SCH in body
		{name: "Fischer", input: "Fischer", want: "FASAR"},
		// PH in body
		{name: "Stephen", input: "Stephen", want: "STAFAN"},
		{name: "Phelps", input: "Phelps", want: "FALP"},
		// H rule: H not between two vowels → use previous char
		{name: "Ahab", input: "Ahab", want: "AHAB"},
		// W after vowel → use vowel
		{name: "Hawaii", input: "Hawaii", want: "H"},
		{name: "Sawyer", input: "Sawyer", want: "SAYAR"},
		// AY ending → Y
		{name: "Hayley", input: "Hayley", want: "HAYLY"},
		{name: "Day", input: "Day", want: "DY"},
		{name: "Ray", input: "Ray", want: "RY"},
		{name: "Murray", input: "Murray", want: "MARY"},
		// Trailing S removal
		{name: "Adams", input: "Adams", want: "ADAN"},
		{name: "Harris", input: "Harris", want: "HAR"},
		// Trailing A removal
		{name: "Santa", input: "Santa", want: "SANT"},
		{name: "Sahara", input: "Sahara", want: "SAHAR"},
		// H between vowels → H
		{name: "Mahoney", input: "Mahoney", want: "MAHANY"},
		// K path
		{name: "Schmitt", input: "Schmitt", want: "SNAT"},
		// Prefix: K (not KN) → C
		{name: "Karl", input: "Karl", want: "CARL"},
		{name: "Kelly", input: "Kelly", want: "CALY"},
		// Prefix: PH → FF
		{name: "Phan", input: "Phan", want: "FAN"},
		// Q → G, Z → S
		{name: "Quiz", input: "Quiz", want: "Q"},
		{name: "Buzz", input: "Buzz", want: "B"},
		// Negative MaxLength → default
		// W not after vowel → W
		{name: "Gifford", input: "Gifford", want: "GAFAD"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NYSIIS(tt.input)
			if got != tt.want {
				t.Errorf("NYSIIS(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNYSIISEncoder_MaxLength(t *testing.T) {
	t.Run("default_length", func(t *testing.T) {
		e := NewNYSIIS()
		if e.MaxLength != 6 {
			t.Errorf("MaxLength = %d, want 6", e.MaxLength)
		}
	})

	t.Run("unlimited", func(t *testing.T) {
		e := &NYSIISEncoder{MaxLength: 0}
		got := e.Encode("WASHINGTON")
		// With unlimited, no truncation.
		if got == "" {
			t.Errorf("Encode() returned empty for valid input")
		}
	})

	t.Run("negative_length", func(t *testing.T) {
		e := &NYSIISEncoder{MaxLength: -1}
		got := e.Encode("WASHINGTON")
		// Negative MaxLength → defaults to 6.
		if len(got) > 6 {
			t.Errorf("Encode() with negative MaxLength: len = %d, want <= 6", len(got))
		}
	})

	t.Run("custom_length", func(t *testing.T) {
		e := &NYSIISEncoder{MaxLength: 4}
		got := e.Encode("LINCOLN")
		if len(got) > 4 {
			t.Errorf("Encode() length = %d, want <= 4", len(got))
		}
	})
}

func TestNYSIISMatch(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want bool
	}{
		{name: "BROOKS_BROCKS", a: "BROOKS", b: "BROCKS", want: true},
		{name: "SMITH_JONES", a: "SMITH", b: "JONES", want: false},
		{name: "empty_empty", a: "", b: "", want: true},
		{name: "empty_nonempty", a: "", b: "SMITH", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NYSIISMatch(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("NYSIISMatch(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestNYSIISEncoder_Interface(t *testing.T) {
	var _ Encoder = NewNYSIIS()
}

func BenchmarkNYSIIS(b *testing.B) {
	e := NewNYSIIS()

	b.ResetTimer()
	for range b.N {
		e.Encode("Washington")
	}
}

func BenchmarkNYSIISMatch(b *testing.B) {
	e := NewNYSIIS()

	b.ResetTimer()
	for range b.N {
		e.Match("Brooks", "Brocks")
	}
}
