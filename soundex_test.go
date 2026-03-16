package strsim

import "testing"

func TestSoundex(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "Robert", input: "Robert", want: "R163"},
		{name: "Rupert", input: "Rupert", want: "R163"},
		{name: "Ashcraft", input: "Ashcraft", want: "A261"},
		{name: "Tymczak", input: "Tymczak", want: "T522"},
		{name: "Pfister", input: "Pfister", want: "P236"},
		{name: "Jackson", input: "Jackson", want: "J250"},
		{name: "Washington", input: "Washington", want: "W252"},
		{name: "Lee", input: "Lee", want: "L000"},
		{name: "Gutierrez", input: "Gutierrez", want: "G362"},
		{name: "Smith", input: "Smith", want: "S530"},
		{name: "Smythe", input: "Smythe", want: "S530"},
		{name: "single_letter", input: "A", want: "A000"},
		{name: "lowercase", input: "robert", want: "R163"},
		{name: "mixed_case", input: "rObErT", want: "R163"},
		{name: "empty", input: "", want: ""},
		{name: "non_alpha", input: "123!@#", want: ""},
		{name: "spaces_and_hyphens", input: "Van De Graaff", want: "V532"},
		{name: "unicode_stripped", input: "Rôbért", want: "R163"},
		{name: "only_unicode", input: "\u00e9\u00e0\u00fc", want: ""},
		{name: "Honeyman", input: "Honeyman", want: "H555"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Soundex(tt.input)
			if got != tt.want {
				t.Errorf("Soundex(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSoundexEncoder_Encode(t *testing.T) {
	e := NewSoundex()

	t.Run("default_length", func(t *testing.T) {
		if e.MaxLength != 4 {
			t.Errorf("MaxLength = %d, want 4", e.MaxLength)
		}
	})

	t.Run("custom_length", func(t *testing.T) {
		enc := &SoundexEncoder{MaxLength: 6}
		got := enc.Encode("Washington")
		if len(got) != 6 {
			t.Errorf("len(Encode) = %d, want 6", len(got))
		}
	})
}

func TestSoundexMatch(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want bool
	}{
		{name: "Robert_Rupert", a: "Robert", b: "Rupert", want: true},
		{name: "Smith_Smythe", a: "Smith", b: "Smythe", want: true},
		{name: "Robert_Smith", a: "Robert", b: "Smith", want: false},
		{name: "empty_empty", a: "", b: "", want: true},
		{name: "empty_nonempty", a: "", b: "Robert", want: false},
		{name: "nonempty_empty", a: "Robert", b: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SoundexMatch(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("SoundexMatch(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestSoundexEncoder_Interface(t *testing.T) {
	var _ Encoder = NewSoundex()
}

func BenchmarkSoundex(b *testing.B) {
	e := NewSoundex()

	b.ResetTimer()
	for range b.N {
		e.Encode("Washington")
	}
}

func BenchmarkSoundexMatch(b *testing.B) {
	e := NewSoundex()

	b.ResetTimer()
	for range b.N {
		e.Match("Robert", "Rupert")
	}
}
