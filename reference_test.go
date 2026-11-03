package strsim

import (
	"os"
	"strings"
	"testing"
)

// TestPhoneticReference compares the encoders with Apache Commons Codec on
// the most common US surnames. The documented differences are:
//   - Metaphone encodes an initial CH as X, where Commons Codec gives K before
//     a vowel ("Church" → XRX, not KRX).
//   - Double Metaphone codes never contain spaces; Commons Codec appends one
//     to the alternate code for a final J ("Raj" → R, not "R ").
//   - NYSIIS keeps a key of a single letter ("Ash" → A), where Commons Codec
//     strips the trailing A and returns "".
func TestPhoneticReference(t *testing.T) {
	data, err := os.ReadFile("testdata/phonetic_reference.tsv")
	if err != nil {
		t.Fatal(err)
	}

	unlimited := &NYSIISEncoder{MaxLength: 0}
	rows := 0

	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if strings.HasPrefix(line, "#") {
			continue
		}

		cols := strings.Split(line, "\t")
		if len(cols) != 6 {
			t.Fatalf("malformed line %q", line)
		}

		name, soundex, metaphone, dmPri, dmAlt, nysiis := cols[0], cols[1], cols[2], cols[3], cols[4], cols[5]
		rows++

		if got := Soundex(name); got != soundex {
			t.Errorf("Soundex(%q) = %q, want %q", name, got, soundex)
		}

		if got := Metaphone(name); got != metaphone && !strings.HasPrefix(strings.ToUpper(name), "CH") {
			t.Errorf("Metaphone(%q) = %q, want %q", name, got, metaphone)
		}

		if pri, alt := DoubleMetaphone(name); pri != dmPri || alt != strings.TrimRight(dmAlt, " ") {
			t.Errorf("DoubleMetaphone(%q) = (%q, %q), want (%q, %q)", name, pri, alt, dmPri, dmAlt)
		}

		if got := unlimited.Encode(name); got != nysiis && (nysiis != "" || len(got) != 1) {
			t.Errorf("NYSIIS(%q) = %q, want %q", name, got, nysiis)
		}
	}

	if rows != 3000 {
		t.Errorf("read %d reference rows, want 3000", rows)
	}
}
