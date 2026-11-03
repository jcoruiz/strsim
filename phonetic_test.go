package strsim

import "testing"

func TestPhonetic_ComposedAndDecomposedAccentsMatch(t *testing.T) {
	pairs := [][2]string{
		{"García", "García"},
		{"Muñoz", "Muñoz"},
		{"Álvarez", "Álvarez"},
		{"José María", "José María"},
	}

	dm := NewDoubleMetaphone()

	for _, p := range pairs {
		for name, e := range map[string]Encoder{"Soundex": NewSoundex(), "Metaphone": NewMetaphone(), "NYSIIS": NewNYSIIS()} {
			if c1, c2 := e.Encode(p[0]), e.Encode(p[1]); c1 != c2 {
				t.Errorf("%s: %q → %q, decomposed %q → %q", name, p[0], c1, p[1], c2)
			}
		}

		p1, a1 := dm.Encode(p[0])
		if p2, a2 := dm.Encode(p[1]); p1 != p2 || a1 != a2 {
			t.Errorf("DoubleMetaphone: %q → (%q, %q), decomposed %q → (%q, %q)", p[0], p1, a1, p[1], p2, a2)
		}
	}
}
