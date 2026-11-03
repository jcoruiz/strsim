package strsim_test

import (
	"errors"
	"fmt"

	"github.com/jcoruiz/strsim"
)

func ExampleLevenshtein() {
	fmt.Println(strsim.Levenshtein("kitten", "sitting"))
	fmt.Printf("%.3f\n", strsim.LevenshteinSimilarity("kitten", "sitting"))
	// Output:
	// 3
	// 0.571
}

func ExampleLevenshteinMetric() {
	// Substitutions cost as much as a deletion plus an insertion.
	m := &strsim.LevenshteinMetric{InsertCost: 1, DeleteCost: 1, ReplaceCost: 2}
	fmt.Println(m.Distance("kitten", "sitting"))
	// Output:
	// 5
}

func ExampleOSA() {
	// OSA counts the swap of adjacent characters as one edit; true
	// Damerau-Levenshtein also allows editing the swapped pair again.
	fmt.Println(strsim.Levenshtein("CA", "ABC"), strsim.OSA("CA", "ABC"), strsim.DamerauLevenshtein("CA", "ABC"))
	// Output:
	// 3 3 2
}

func ExampleHamming() {
	d, err := strsim.Hamming("karolin", "kathrin")
	fmt.Println(d, err)

	_, err = strsim.Hamming("abc", "ab")
	fmt.Println(errors.Is(err, strsim.ErrLengthMismatch))
	// Output:
	// 3 <nil>
	// true
}

func ExampleJaroWinklerSimilarity() {
	fmt.Printf("%.4f\n", strsim.JaroSimilarity("martha", "marhta"))
	fmt.Printf("%.4f\n", strsim.JaroWinklerSimilarity("martha", "marhta"))
	// Output:
	// 0.9444
	// 0.9611
}

func ExampleNgramMetric() {
	trigrams := &strsim.NgramMetric{Size: 3}
	fmt.Printf("%.4f %.4f\n", strsim.DiceSimilarity("night", "nacht"), trigrams.Dice("night", "nacht"))
	// Output:
	// 0.2500 0.0000
}

func ExampleSoundex() {
	fmt.Println(strsim.Soundex("Robert"), strsim.Soundex("Rupert"), strsim.Soundex("Álvarez"))
	// Output:
	// R163 R163 A416
}

func ExampleDoubleMetaphone() {
	fmt.Println(strsim.DoubleMetaphone("Schmidt"))
	fmt.Println(strsim.DoubleMetaphoneMatch("Smith", "Schmidt"))
	// Output:
	// XMT SMT
	// true
}

func ExampleFindBestMatch() {
	candidates := []string{"golang", "gold", "golem", "gopher", "python"}
	m := strsim.FindBestMatch("goland", candidates, strsim.NewJaroWinkler())
	fmt.Printf("%s %d %.3f\n", m.Value, m.Index, m.Similarity)
	// Output:
	// golang 0 0.933
}

func ExampleFindTopN() {
	candidates := []string{"apple", "apply", "ample", "maple", "angle"}
	for _, m := range strsim.FindTopN("appel", candidates, 3, strsim.NewDamerauLevenshtein()) {
		fmt.Printf("%s %.1f\n", m.Value, m.Similarity)
	}
	// Output:
	// apple 0.8
	// apply 0.6
	// ample 0.6
}
