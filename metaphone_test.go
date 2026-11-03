package strsim

import "testing"

func TestMetaphone(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "Smith", input: "Smith", want: "SM0"},
		{name: "Schmidt", input: "Schmidt", want: "SKMT"},
		{name: "John", input: "John", want: "JN"},
		{name: "Phone", input: "Phone", want: "FN"},
		{name: "Thomas", input: "Thomas", want: "0MS"},
		{name: "Knight", input: "Knight", want: "NT"},
		{name: "Pneumonia", input: "Pneumonia", want: "NMN"},
		{name: "Write", input: "Write", want: "RT"},
		{name: "Aero", input: "Aero", want: "ER"},
		{name: "Gnome", input: "Gnome", want: "NM"},
		{name: "Dumb", input: "Dumb", want: "TM"},
		{name: "Edge", input: "Edge", want: "EJ"},
		{name: "Ghost", input: "Ghost", want: "KST"},
		{name: "Cent", input: "Cent", want: "SNT"},
		{name: "Cyst", input: "Cyst", want: "SST"},
		{name: "Ship", input: "Ship", want: "XP"},
		{name: "Vision", input: "Vision", want: "FXN"},
		{name: "Nation", input: "Nation", want: "NXN"},
		{name: "empty", input: "", want: ""},
		{name: "non_alpha", input: "123", want: ""},
		{name: "single_vowel", input: "A", want: "A"},
		{name: "single_consonant", input: "B", want: "B"},
		{name: "Xerox", input: "Xerox", want: "SRKS"},
		{name: "Why", input: "Why", want: ""},
		{name: "Whole", input: "Whole", want: "WL"},
		{name: "Thumb", input: "Thumb", want: "0M"},

		// SC before I/E/Y → skip C
		{name: "Science", input: "Science", want: "SNS"},
		{name: "Scenic", input: "Scenic", want: "SNK"},
		// DG before I/E/Y → J
		{name: "Judge", input: "Judge", want: "JJ"},
		{name: "Badge", input: "Badge", want: "BJ"},
		{name: "Madge", input: "Madge", want: "MJ"},
		{name: "Dodge", input: "Dodge", want: "TJ"},
		// DG before non-I/E/Y → T
		{name: "Edgar", input: "Edgar", want: "ETKR"},
		// G before H: GH at end → silent
		{name: "Rough", input: "Rough", want: "R"},
		// G before H: GH + non-vowel → both silent
		{name: "Night", input: "Night", want: "NT"},
		// G before H: GH + vowel → K
		{name: "Ghee", input: "Ghee", want: "K"},
		// GN at end → G silent
		{name: "Sign", input: "Sign", want: "SN"},
		// G before I/E/Y after another G → skip (silent second G)
		{name: "Eggs", input: "Eggs", want: "EKS"},
		// G before I/E/Y not after G → J
		{name: "Huge", input: "Huge", want: "HJ"},
		// K after C → silent
		{name: "Acknowledge", input: "Acknowledge", want: "AKNL"},
		{name: "Lacking", input: "Lacking", want: "LKNK"},
		// K not after C
		{name: "King", input: "King", want: "KNK"},
		// W before vowel → W
		{name: "Wow", input: "Wow", want: "W"},
		// W not before vowel → silent
		// Y before vowel → Y
		{name: "Yawn", input: "Yawn", want: "YN"},
		// Y not before vowel → silent
		{name: "By", input: "By", want: "B"},
		// PH → F
		{name: "Phi", input: "Phi", want: "F"},
		// H before vowel after non-vowel → H
		// H between vowels → silent
		{name: "Ohio", input: "Ohio", want: "OH"},
		// CH → X
		{name: "Ache", input: "Ache", want: "AX"},
		// Accept → CC handling
		{name: "Accept", input: "Accept", want: "AKSP"},
		// G before non-I/E/Y → K
		{name: "Gap", input: "Gap", want: "KP"},
		{name: "Gun", input: "Gun", want: "KN"},
		// W not before vowel → silent
		{name: "Brew", input: "Brew", want: "BR"},
		{name: "Stew", input: "Stew", want: "ST"},
		// Duplicate adjacent letters (not C) → skip
		{name: "Hippo", input: "Hippo", want: "HP"},
		// Reference values, checked against Apache Commons Codec.
		{name: "White", input: "White", want: "WT"},
		{name: "Watch", input: "Watch", want: "WX"},
		{name: "School", input: "School", want: "SKL"},
		{name: "Acacia", input: "Acacia", want: "AKX"},
		{name: "Signed", input: "Signed", want: "SNT"},
		{name: "Cohen", input: "Cohen", want: "KHN"},
		// Initial CH is X, as in Philips' rules (Commons Codec gives KRX).
		{name: "Church", input: "Church", want: "XRX"},
		{name: "single_h", input: "H", want: "H"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Metaphone(tt.input)
			if got != tt.want {
				t.Errorf("Metaphone(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestMetaphoneEncoder_MaxLength(t *testing.T) {
	t.Run("custom_length", func(t *testing.T) {
		e := &MetaphoneEncoder{MaxLength: 3}
		got := e.Encode("Washington")
		if len(got) > 3 {
			t.Errorf("Encode() length = %d, want <= 3", len(got))
		}
	})

	t.Run("default_encoder", func(t *testing.T) {
		e := NewMetaphone()
		if e.MaxLength != 4 {
			t.Errorf("MaxLength = %d, want 4", e.MaxLength)
		}
	})
}

func TestMetaphoneMatch(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want bool
	}{
		{name: "Smith_Smyth", a: "Smith", b: "Smyth", want: true},
		{name: "Smith_Jones", a: "Smith", b: "Jones", want: false},
		{name: "empty_empty", a: "", b: "", want: true},
		{name: "empty_nonempty", a: "", b: "Smith", want: false},
		{name: "Knight_Night", a: "Knight", b: "Night", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MetaphoneMatch(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("MetaphoneMatch(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestMetaphoneEncoder_Interface(t *testing.T) {
	var _ Encoder = NewMetaphone()
}

func BenchmarkMetaphone(b *testing.B) {
	e := NewMetaphone()

	b.ResetTimer()
	for range b.N {
		e.Encode("Washington")
	}
}

func BenchmarkMetaphoneMatch(b *testing.B) {
	e := NewMetaphone()

	b.ResetTimer()
	for range b.N {
		e.Match("Smith", "Schmidt")
	}
}
