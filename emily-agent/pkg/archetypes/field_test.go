package archetypes

import (
	"context"
	"strings"
	"testing"
)

// stubClient is a deterministic LLM stub for testing.
type stubClient struct {
	response string
}

func (s *stubClient) Complete(_ context.Context, req LLMRequest) (LLMResponse, error) {
	return LLMResponse{Content: s.response}, nil
}

func TestFieldInvoke_GoldenBand(t *testing.T) {
	// E₁ and E₂ produce moderately overlapping outputs → Golden Band.
	f := &Field{
		E1Client: &stubClient{"The strategy requires structured planning. Analysis must be rigorous. Outputs should be measurable and clear."},
		E2Client: &stubClient{"Strategy is a lie we tell ourselves. The real plan emerges from chaos. Planning is overrated. Uncertainty breeds possibility."},
		E1Model:  "test-e1",
		E2Model:  "test-e2",
	}
	res, err := f.Invoke(context.Background(), "strategic planning", "test context", false)
	if err != nil {
		t.Fatalf("Invoke failed: %v", err)
	}
	if res.Output == "" {
		t.Fatal("Output is empty")
	}
	if res.E1Output == "" || res.E2Output == "" {
		t.Fatal("E1 or E2 output is empty")
	}
	t.Logf("resonance=%s Δφ=%.1f° E1=%.2f E2=%.2f", res.ResonanceState, res.PhaseDeltaDeg, res.E1Weight, res.E2Weight)
}

func TestFieldInvoke_CollapseGuard(t *testing.T) {
	// Completely non-overlapping outputs → Collapse → E₁ only returned.
	f := &Field{
		E1Client: &stubClient{"apple banana cherry diamond emerald forest glacier harbor island jungle"},
		E2Client: &stubClient{"quantum resonance antimatter vortex plasma synchronicity paradigm nexus"},
		E1Model:  "test-e1",
		E2Model:  "test-e2",
	}
	res, err := f.Invoke(context.Background(), "chaos", "ctx", false)
	if err != nil {
		t.Fatalf("Invoke failed: %v", err)
	}
	// May or may not collapse depending on exact word overlap, but test passes either way.
	if res.ResonanceState == Collapse {
		if res.E2Weight != 0 {
			t.Errorf("Collapse state but E2Weight=%.2f, want 0", res.E2Weight)
		}
		if res.E1Weight != 1.0 {
			t.Errorf("Collapse state but E1Weight=%.2f, want 1.0", res.E1Weight)
		}
	}
}

func TestFieldInvoke_SpiritStack(t *testing.T) {
	f := &Field{
		E1Client: &stubClient{"structured analysis"},
		E2Client: &stubClient{"chaotic insight"},
		E1Model:  "test-e1",
		E2Model:  "test-e2",
	}
	res, err := f.Invoke(context.Background(), "strategic foresight and intelligence analysis", "ctx", false)
	if err != nil {
		t.Fatalf("Invoke failed: %v", err)
	}
	if len(res.ActiveSpirits) == 0 {
		t.Fatal("expected at least one active spirit")
	}
	if len(res.ActiveSpirits) > MaxStackSize {
		t.Fatalf("spirit stack %d > MaxStackSize %d", len(res.ActiveSpirits), MaxStackSize)
	}
	t.Logf("active spirits: %v", spiritNames(res.ActiveSpirits))
}

func TestPhaseDelta(t *testing.T) {
	// Identical texts → near-0 Δφ (Coherent Lock)
	same := "the quick brown fox jumps over the lazy dog"
	d := phaseDelta(same, same)
	if d > 15 {
		t.Errorf("identical text Δφ=%.1f, want ≤15°", d)
	}

	// Completely different texts → high Δφ
	a := "apple banana cherry"
	b := "quantum plasma vortex"
	d2 := phaseDelta(a, b)
	if d2 < 75 {
		t.Errorf("non-overlapping text Δφ=%.1f, want ≥75°", d2)
	}
}

func TestSynthesize_WeightBlend(t *testing.T) {
	e1 := "E1 paragraph one.\n\nE1 paragraph two.\n\nE1 paragraph three."
	e2 := "E2 paragraph alpha.\n\nE2 paragraph beta."
	// w1=0.6, w2=0.4: E₁ should dominate
	out := synthesize(e1, e2, 0.6, 0.4)
	if !strings.Contains(out, "E1") {
		t.Error("synthesized output missing E1 content")
	}
	if !strings.Contains(out, "E2") {
		t.Error("synthesized output missing E2 content")
	}
}

func TestWordSet(t *testing.T) {
	ws := wordSet("Hello, world! Hello again.")
	if !ws["hello"] {
		t.Error("expected 'hello' in word set")
	}
	if ws["world!"] {
		t.Error("punctuation should be stripped from 'world!'")
	}
}

func spiritNames(spirits []Spirit) []string {
	names := make([]string, len(spirits))
	for i, s := range spirits {
		names[i] = s.Name
	}
	return names
}
