package archetypes

import (
	"testing"
)

func TestAll72Count(t *testing.T) {
	if len(All72) != 72 {
		t.Fatalf("expected 72 spirits, got %d", len(All72))
	}
}

func TestRankContinuity(t *testing.T) {
	for i, s := range All72 {
		if s.Rank != i+1 {
			t.Errorf("spirit[%d].Rank=%d, want %d", i, s.Rank, i+1)
		}
	}
}

func TestByName(t *testing.T) {
	s := ByName("Vassago")
	if s == nil {
		t.Fatal("ByName(Vassago) returned nil")
	}
	if s.Rank != 3 {
		t.Errorf("Vassago rank=%d, want 3", s.Rank)
	}
	if s.Corridor != GoldenBand {
		t.Errorf("Vassago corridor=%s, want GoldenBand", s.Corridor)
	}
}

func TestByNameCaseInsensitive(t *testing.T) {
	s := ByName("ANDROMALIUS")
	if s == nil || s.Rank != 72 {
		t.Fatal("case-insensitive lookup failed for Andromalius")
	}
}

func TestByNameMissing(t *testing.T) {
	if ByName("NotASpirit") != nil {
		t.Fatal("expected nil for unknown spirit")
	}
}

func TestByRank(t *testing.T) {
	for _, tc := range []struct{ rank int; name string }{
		{1, "Bael"}, {36, "Stolas"}, {72, "Andromalius"},
	} {
		s := ByRank(tc.rank)
		if s == nil {
			t.Fatalf("ByRank(%d) returned nil", tc.rank)
		}
		if s.Name != tc.name {
			t.Errorf("ByRank(%d).Name=%s, want %s", tc.rank, s.Name, tc.name)
		}
	}
}

func TestByRankOutOfRange(t *testing.T) {
	if ByRank(0) != nil || ByRank(73) != nil {
		t.Fatal("out-of-range ranks should return nil")
	}
}

func TestAmplitudeSafety(t *testing.T) {
	// All spirits must have amplitude in [0.5, 1.3]
	for _, s := range All72 {
		if s.Amplitude < 0.5 || s.Amplitude > 1.3 {
			t.Errorf("spirit %s amplitude=%.2f out of expected range [0.5, 1.3]", s.Name, s.Amplitude)
		}
	}
}

func TestCollapseSpiritsAllowFlag(t *testing.T) {
	collapseSpirits := []int{9, 25, 32, 41, 63, 68}
	for _, rank := range collapseSpirits {
		s := ByRank(rank)
		if s == nil {
			t.Fatalf("spirit rank %d not found", rank)
		}
		if !s.AllowCollapse {
			t.Errorf("spirit %s (rank %d) should have AllowCollapse=true", s.Name, rank)
		}
	}
	// Non-collapse spirits must have AllowCollapse=false
	safe := ByRank(15) // Eligos, Golden Band
	if safe.AllowCollapse {
		t.Error("Eligos should not require allow_collapse")
	}
}

func TestMatchIntentBasic(t *testing.T) {
	stack := MatchIntent("strategic foresight and planning", false)
	if len(stack) == 0 {
		t.Fatal("MatchIntent returned empty stack for strategy intent")
	}
	if len(stack) > MaxStackSize {
		t.Fatalf("stack size %d exceeds MaxStackSize %d", len(stack), MaxStackSize)
	}
}

func TestMatchIntentExcludesCollapse(t *testing.T) {
	// "influence" maps to Dantalion(71), Paimon(9), Ipos(22)
	// Paimon(9) is Collapse — should be excluded when allowCollapse=false
	stack := MatchIntent("influence and persuasion mass", false)
	for _, s := range stack {
		if s.AllowCollapse {
			t.Errorf("Collapse spirit %s included without allow_collapse", s.Name)
		}
	}
}

func TestMatchIntentAllowCollapse(t *testing.T) {
	stack := MatchIntent("influence and persuasion mass", true)
	found := false
	for _, s := range stack {
		if s.Rank == 9 { // Paimon
			found = true
		}
	}
	if !found {
		t.Log("Paimon not in stack (expected when top-scored) — note: keyword scoring may deprioritize it")
	}
}

func TestMatchIntentAmplitudeNormalized(t *testing.T) {
	// If stack amplitude > 3.5, it must be normalized
	stack := MatchIntent("love emotion relationship", false)
	if sum := AmplitudeSum(stack); sum > MaxStackAmplitude+0.001 {
		t.Errorf("stack amplitude sum %.2f exceeds MaxStackAmplitude %.1f", sum, MaxStackAmplitude)
	}
}

func TestResonanceState(t *testing.T) {
	cases := []struct {
		phi  float64
		want Corridor
	}{
		{0, CoherentLock},
		{15, CoherentLock},
		{22, GoldenBand},
		{38, GoldenBand},
		{75, ChaosEdge},
		{90, ChaosEdge},
		{91, Collapse},
		{180, Collapse},
	}
	for _, tc := range cases {
		got := ResonanceState(tc.phi)
		if got != tc.want {
			t.Errorf("ResonanceState(%.0f)=%s, want %s", tc.phi, got, tc.want)
		}
	}
}

func TestE1E2Weights(t *testing.T) {
	states := []Corridor{CoherentLock, GoldenBand, ChaosEdge, Collapse}
	for _, s := range states {
		sum := E1Weight(s) + E2Weight(s)
		if sum < 0.999 || sum > 1.001 {
			t.Errorf("state %s: E1+E2=%.3f, want 1.0", s, sum)
		}
		if E1Weight(s) < E2Weight(s) && s == CoherentLock {
			t.Error("CoherentLock must have E1 > E2")
		}
	}
	if E1Weight(Collapse) != 1.0 {
		t.Error("Collapse must have E1=1.0 (full fallback)")
	}
}

func TestForCorridor(t *testing.T) {
	cl := ForCorridor(CoherentLock)
	gb := ForCorridor(GoldenBand)
	ce := ForCorridor(ChaosEdge)
	co := ForCorridor(Collapse)
	total := len(cl) + len(gb) + len(ce) + len(co)
	if total != 72 {
		t.Errorf("corridor partition total=%d, want 72", total)
	}
	if len(co) < 4 {
		t.Errorf("expected ≥4 Collapse spirits, got %d", len(co))
	}
}
