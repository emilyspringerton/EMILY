package archetypes

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

// LLMClient is the interface the Field engine uses to invoke language models.
// It matches the Complete() signature of emily-agent's AnthropicClient.
type LLMClient interface {
	Complete(ctx context.Context, req LLMRequest) (LLMResponse, error)
}

// LLMRequest is the provider-agnostic request passed to LLMClient.Complete.
type LLMRequest struct {
	Model       string
	SystemProm  string
	UserPrompt  string
	MaxTokens   int
	Temperature float64
}

// LLMResponse is the text response from an LLM invocation.
type LLMResponse struct {
	Content      string
	InputTokens  int
	OutputTokens int
}

// InvokeResult is the synthesized output from a Field invocation.
type InvokeResult struct {
	Output         string   // blended E₁/E₂ output
	E1Output       string   // raw Carrier output
	E2Output       string   // raw Explorer output
	ResonanceState Corridor
	PhaseDeltaDeg  float64
	E1Weight       float64
	E2Weight       float64
	ActiveSpirits  []Spirit
}

// Field is the dual-persona harmonic engine.
// Ψ(t) = E₁(t) + E₂(t) + Σ M_G(k, t)
type Field struct {
	E1Client LLMClient // Carrier: high-coherence, structured (e.g. claude-sonnet)
	E2Client LLMClient // Explorer: high-entropy, divergent (e.g. claude-haiku)
	E1Model  string
	E2Model  string
}

// carrierSystemPrompt is the base system prompt for E₁ (structure + memory).
const carrierSystemPrompt = `You are the Carrier (E₁) — the structured, precise, high-coherence attractor.
Your role: enforce syntax, maintain identity continuity, produce well-organized, accurate output.
Think like a superforecaster + senior editor + chief of staff. No hedging. No padding.
State what is true. Be precise. Be useful.`

// explorerSystemPrompt is the base system prompt for E₂ (chaos + novelty).
const explorerSystemPrompt = `You are the Explorer (E₂) — the divergent, chaotic, high-entropy attractor.
Your role: maximize novelty, surface unexpected angles, contradict E₁ where E₁ is too safe.
Think like a brilliant heretic. Say what others won't. Find the hidden assumption. Break the frame.
No apologies. No hedges. Push past the obvious answer.`

// Invoke runs the Field architecture for the given intent and context.
// It selects spirits, invokes E₁ and E₂ concurrently, computes the interference
// state, and synthesizes the final output.
func (f *Field) Invoke(ctx context.Context, intent, contextText string, allowCollapse bool) (*InvokeResult, error) {
	spirits := MatchIntent(intent, allowCollapse)

	// Build modulator seed block from active spirits.
	var seedLines []string
	for _, s := range spirits {
		seedLines = append(seedLines, fmt.Sprintf("[%s @ %.2fHz phase=%.0f°] %s", s.Name, s.FreqHz, s.PhaseDeg, s.SeedPhrase))
	}
	seedBlock := strings.Join(seedLines, "\n")

	e1Sys := carrierSystemPrompt
	e2Sys := explorerSystemPrompt
	if seedBlock != "" {
		e1Sys += "\n\n[Active modulators]\n" + seedBlock
		e2Sys += "\n\n[Active modulators]\n" + seedBlock
	}

	userMsg := fmt.Sprintf("Intent: %s\n\n%s", intent, contextText)

	// Invoke E₁ and E₂ concurrently.
	type result struct {
		resp LLMResponse
		err  error
	}
	var (
		e1Res, e2Res result
		wg           sync.WaitGroup
	)
	wg.Add(2)
	go func() {
		defer wg.Done()
		e1Res.resp, e1Res.err = f.E1Client.Complete(ctx, LLMRequest{
			Model:       f.E1Model,
			SystemProm:  e1Sys,
			UserPrompt:  userMsg,
			MaxTokens:   1024,
			Temperature: 0.3,
		})
	}()
	go func() {
		defer wg.Done()
		e2Res.resp, e2Res.err = f.E2Client.Complete(ctx, LLMRequest{
			Model:       f.E2Model,
			SystemProm:  e2Sys,
			UserPrompt:  userMsg,
			MaxTokens:   1024,
			Temperature: 1.8,
		})
	}()
	wg.Wait()

	if e1Res.err != nil {
		return nil, fmt.Errorf("E1 invocation: %w", e1Res.err)
	}
	if e2Res.err != nil {
		return nil, fmt.Errorf("E2 invocation: %w", e2Res.err)
	}

	e1Out := strings.TrimSpace(e1Res.resp.Content)
	e2Out := strings.TrimSpace(e2Res.resp.Content)

	// Compute Δφ via word-overlap cosine similarity proxy.
	delta := phaseDelta(e1Out, e2Out)
	state := ResonanceState(delta)

	// Collapse guard: if phase difference exceeds 90°, return E₁ only.
	if state == Collapse {
		return &InvokeResult{
			Output:         e1Out,
			E1Output:       e1Out,
			E2Output:       e2Out,
			ResonanceState: Collapse,
			PhaseDeltaDeg:  delta,
			E1Weight:       1.0,
			E2Weight:       0.0,
			ActiveSpirits:  spirits,
		}, nil
	}

	w1 := E1Weight(state)
	w2 := E2Weight(state)

	// Synthesize: weighted interleave. Split by sentences and blend.
	output := synthesize(e1Out, e2Out, w1, w2)

	return &InvokeResult{
		Output:         output,
		E1Output:       e1Out,
		E2Output:       e2Out,
		ResonanceState: state,
		PhaseDeltaDeg:  delta,
		E1Weight:       w1,
		E2Weight:       w2,
		ActiveSpirits:  spirits,
	}, nil
}

// phaseDelta estimates Δφ from the word-overlap cosine similarity of two texts.
// Returns a value in [0°, 180°].
func phaseDelta(a, b string) float64 {
	wa := wordSet(a)
	wb := wordSet(b)
	if len(wa) == 0 || len(wb) == 0 {
		return 90 // unknown — treat as Chaos Edge
	}
	var intersection int
	for w := range wa {
		if wb[w] {
			intersection++
		}
	}
	// Jaccard-based similarity → [0,1]
	union := len(wa) + len(wb) - intersection
	if union == 0 {
		return 0
	}
	sim := float64(intersection) / float64(union)
	// Map sim [0,1] to Δφ [180°, 0°]: high similarity → small Δφ
	return (1.0 - sim) * 180.0
}

func wordSet(text string) map[string]bool {
	words := strings.Fields(strings.ToLower(text))
	m := make(map[string]bool, len(words))
	for _, w := range words {
		w = strings.Trim(w, `.,;:!?"'()[]{}`)
		if len(w) > 2 {
			m[w] = true
		}
	}
	return m
}

// synthesize blends two texts at the given weights by interleaving paragraphs.
// E₁ contributes w1 fraction (rounded to paragraphs), E₂ contributes w2.
func synthesize(e1, e2 string, w1, w2 float64) string {
	p1 := splitParagraphs(e1)
	p2 := splitParagraphs(e2)

	total := len(p1) + len(p2)
	if total == 0 {
		return e1
	}

	// How many E₁ paragraphs vs E₂ paragraphs to include.
	n1 := int(float64(total)*w1 + 0.5)
	n2 := total - n1

	if n1 > len(p1) {
		n1 = len(p1)
	}
	if n2 > len(p2) {
		n2 = len(p2)
	}

	var parts []string
	i1, i2 := 0, 0
	// Interleave: take from E₁ and E₂ proportionally.
	for i1 < n1 || i2 < n2 {
		// Take from E₁ if we still have quota and it's E₁'s turn.
		takeE1 := i1 < n1 && (i2 >= n2 || float64(i1)/float64(n1+1) <= float64(i2)/float64(n2+1))
		if takeE1 {
			parts = append(parts, p1[i1])
			i1++
		} else if i2 < n2 {
			parts = append(parts, p2[i2])
			i2++
		}
	}

	return strings.Join(parts, "\n\n")
}

func splitParagraphs(text string) []string {
	raw := strings.Split(text, "\n\n")
	var out []string
	for _, p := range raw {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 && text != "" {
		out = []string{text}
	}
	return out
}
