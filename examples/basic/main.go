// Package main demonstrates basic usage of the strsim library.
package main

import (
	"fmt"

	"github.com/jcoruiz/strsim"
)

func main() {
	// Edit distance.
	fmt.Println("=== Edit Distance ===")
	fmt.Printf("Levenshtein(\"kitten\", \"sitting\") = %d\n", strsim.Levenshtein("kitten", "sitting"))
	fmt.Printf("DamerauLevenshtein(\"CA\", \"ABC\") = %d\n", strsim.DamerauLevenshtein("CA", "ABC"))

	// Normalized similarity [0, 1].
	fmt.Println("\n=== Similarity ===")
	fmt.Printf("LevenshteinSimilarity(\"kitten\", \"sitting\") = %.4f\n", strsim.LevenshteinSimilarity("kitten", "sitting"))
	fmt.Printf("JaroWinklerSimilarity(\"martha\", \"marhta\") = %.4f\n", strsim.JaroWinklerSimilarity("martha", "marhta"))
	fmt.Printf("CosineSimilarity(\"night\", \"nacht\") = %.4f\n", strsim.CosineSimilarity("night", "nacht"))
	fmt.Printf("DiceSimilarity(\"night\", \"nacht\") = %.4f\n", strsim.DiceSimilarity("night", "nacht"))

	// Phonetic encoding.
	fmt.Println("\n=== Phonetic ===")
	fmt.Printf("Soundex(\"Robert\") = %s\n", strsim.Soundex("Robert"))
	fmt.Printf("Soundex(\"Rupert\") = %s\n", strsim.Soundex("Rupert"))
	fmt.Printf("SoundexMatch(\"Robert\", \"Rupert\") = %v\n", strsim.SoundexMatch("Robert", "Rupert"))
	fmt.Printf("Metaphone(\"Smith\") = %s\n", strsim.Metaphone("Smith"))
	pri, alt := strsim.DoubleMetaphone("Schmidt")
	fmt.Printf("DoubleMetaphone(\"Schmidt\") = (%s, %s)\n", pri, alt)

	// Batch operations.
	fmt.Println("\n=== Batch ===")
	candidates := []string{"golang", "gold", "golem", "gopher", "python", "ruby"}
	match := strsim.FindBestMatch("go", candidates, strsim.NewJaroWinkler())
	fmt.Printf("FindBestMatch(\"go\", ...) = %q (similarity: %.4f)\n", match.Value, match.Similarity)

	top := strsim.FindTopN("go", candidates, 3, strsim.NewJaroWinkler())
	fmt.Println("FindTopN(\"go\", ..., 3):")
	for _, m := range top {
		fmt.Printf("  %q → %.4f\n", m.Value, m.Similarity)
	}
}
