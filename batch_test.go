package strsim

import (
	"math"
	"testing"
)

func TestFindBestMatch(t *testing.T) {
	candidates := []string{"golang", "gold", "golem", "gopher", "python"}
	m := NewJaroWinkler()

	t.Run("normal", func(t *testing.T) {
		got := FindBestMatch("go", candidates, m)
		if got.Index < 0 || got.Value == "" {
			t.Fatal("expected a match")
		}
		if got.Similarity <= 0 {
			t.Errorf("expected positive similarity, got %f", got.Similarity)
		}
	})

	t.Run("exact_match", func(t *testing.T) {
		got := FindBestMatch("golang", candidates, m)
		if got.Value != "golang" {
			t.Errorf("value = %q, want %q", got.Value, "golang")
		}
		if got.Similarity != 1.0 {
			t.Errorf("similarity = %f, want 1.0", got.Similarity)
		}
	})

	t.Run("empty_candidates", func(t *testing.T) {
		got := FindBestMatch("test", nil, m)
		if got.Index != -1 {
			t.Errorf("index = %d, want -1", got.Index)
		}
	})

	t.Run("single_candidate", func(t *testing.T) {
		got := FindBestMatch("test", []string{"testing"}, m)
		if got.Index != 0 || got.Value != "testing" {
			t.Errorf("unexpected match: %+v", got)
		}
	})
}

func TestFindTopN(t *testing.T) {
	candidates := []string{"golang", "gold", "golem", "gopher", "python"}
	m := NewJaroWinkler()

	t.Run("top_3", func(t *testing.T) {
		got := FindTopN("go", candidates, 3, m)
		if len(got) != 3 {
			t.Fatalf("len = %d, want 3", len(got))
		}
		// Verify descending order.
		for i := 1; i < len(got); i++ {
			if got[i].Similarity > got[i-1].Similarity {
				t.Errorf("not sorted: [%d].Similarity=%f > [%d].Similarity=%f",
					i, got[i].Similarity, i-1, got[i-1].Similarity)
			}
		}
	})

	t.Run("n_exceeds_candidates", func(t *testing.T) {
		got := FindTopN("go", candidates, 100, m)
		if len(got) != len(candidates) {
			t.Errorf("len = %d, want %d", len(got), len(candidates))
		}
	})

	t.Run("empty_candidates", func(t *testing.T) {
		got := FindTopN("test", nil, 3, m)
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})

	t.Run("zero_n", func(t *testing.T) {
		got := FindTopN("test", candidates, 0, m)
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
}

func TestFindAboveThreshold(t *testing.T) {
	candidates := []string{"golang", "gold", "golem", "gopher", "python"}
	m := NewJaroWinkler()

	t.Run("threshold_0.8", func(t *testing.T) {
		got := FindAboveThreshold("golang", candidates, 0.8, m)
		if len(got) == 0 {
			t.Fatal("expected at least one match")
		}
		for _, match := range got {
			if match.Similarity < 0.8 {
				t.Errorf("match %q has similarity %f < 0.8", match.Value, match.Similarity)
			}
		}
		// Verify descending order.
		for i := 1; i < len(got); i++ {
			if got[i].Similarity > got[i-1].Similarity {
				t.Errorf("not sorted: [%d].Similarity=%f > [%d].Similarity=%f",
					i, got[i].Similarity, i-1, got[i-1].Similarity)
			}
		}
	})

	t.Run("threshold_1.0", func(t *testing.T) {
		got := FindAboveThreshold("golang", candidates, 1.0, m)
		if len(got) != 1 || got[0].Value != "golang" {
			t.Errorf("expected only exact match, got %v", got)
		}
	})

	t.Run("threshold_too_high", func(t *testing.T) {
		got := FindAboveThreshold("xyz", candidates, 0.99, m)
		if len(got) != 0 {
			t.Errorf("expected no matches, got %v", got)
		}
	})

	t.Run("empty_candidates", func(t *testing.T) {
		got := FindAboveThreshold("test", nil, 0.5, m)
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
}

func TestBatchWithDistanceMetric(t *testing.T) {
	candidates := []string{"kitten", "sitting", "kit", "kitchen"}
	m := NewLevenshtein()

	got := FindBestMatch("kitten", candidates, m)
	if got.Value != "kitten" {
		t.Errorf("value = %q, want %q", got.Value, "kitten")
	}

	if math.Abs(got.Similarity-1.0) > 1e-10 {
		t.Errorf("similarity = %f, want 1.0", got.Similarity)
	}
}

func BenchmarkFindBestMatch(b *testing.B) {
	candidates := make([]string, 1000)
	for i := range candidates {
		candidates[i] = "candidate_" + string(rune('A'+i%26))
	}

	m := NewJaroWinkler()

	b.ResetTimer()
	for range b.N {
		FindBestMatch("golang", candidates, m)
	}
}

func BenchmarkFindTopN(b *testing.B) {
	candidates := make([]string, 1000)
	for i := range candidates {
		candidates[i] = "candidate_" + string(rune('A'+i%26))
	}

	m := NewJaroWinkler()

	b.ResetTimer()
	for range b.N {
		FindTopN("golang", candidates, 10, m)
	}
}
