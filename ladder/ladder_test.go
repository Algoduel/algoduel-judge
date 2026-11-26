package ladder

import (
	"math"
	"testing"
)

func TestExps(t *testing.T) {
	got := Config{From: 10, To: 22}.Exps()
	if len(got) != 13 || got[0] != 10 || got[12] != 22 {
		t.Fatalf("unexpected rungs: %v", got)
	}
}

func TestScore(t *testing.T) {
	if s := Score(4, 13); math.Abs(s-4.0/13) > 1e-9 {
		t.Fatalf("score = %v", s)
	}
	if Score(1, 0) != 0 {
		t.Fatal("empty ladder must give 0")
	}
}

func TestGrowthExponents(t *testing.T) {
	quadratic := []uint64{1 << 20, 1 << 22, 1 << 24} // ×4 на каждой ступени
	for _, g := range GrowthExponents(quadratic) {
		if math.Abs(g-2) > 1e-9 {
			t.Fatalf("expected 2 for O(n^2), got %v", g)
		}
	}
}
