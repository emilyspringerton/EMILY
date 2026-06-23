// archetype-engine: THE_FIELD HTTP service on :8090
// POST /invoke   — run E₁/E₂ dual-persona invocation for an intent
// GET  /spirits  — return all 72 Goetia spirits as JSON
// GET  /status   — engine status
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"emily-agent/pkg/archetypes"
)

// anthropicClient wraps the raw Anthropic Messages API for the archetypes package.
type anthropicClient struct {
	apiKey  string
	baseURL string
}

func newAnthropicClient(apiKey string) *anthropicClient {
	return &anthropicClient{apiKey: apiKey, baseURL: "https://api.anthropic.com/v1"}
}

func (c *anthropicClient) Complete(ctx context.Context, req archetypes.LLMRequest) (archetypes.LLMResponse, error) {
	body := map[string]any{
		"model":      req.Model,
		"max_tokens": req.MaxTokens,
		"temperature": req.Temperature,
		"system":     req.SystemProm,
		"messages": []map[string]string{
			{"role": "user", "content": req.UserPrompt},
		},
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return archetypes.LLMResponse{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/messages", bytes.NewReader(raw))
	if err != nil {
		return archetypes.LLMResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return archetypes.LLMResponse{}, err
	}
	defer resp.Body.Close()

	var out struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return archetypes.LLMResponse{}, fmt.Errorf("anthropic decode: %w", err)
	}
	if out.Error != nil {
		return archetypes.LLMResponse{}, fmt.Errorf("anthropic: %s", out.Error.Message)
	}
	if len(out.Content) == 0 {
		return archetypes.LLMResponse{}, fmt.Errorf("anthropic: empty response")
	}
	return archetypes.LLMResponse{
		Content:      out.Content[0].Text,
		InputTokens:  out.Usage.InputTokens,
		OutputTokens: out.Usage.OutputTokens,
	}, nil
}

// ---

type server struct {
	field         *archetypes.Field
	allowCollapse bool
	startedAt     time.Time
}

type invokeRequest struct {
	Intent        string `json:"intent"`
	Context       string `json:"context"`
	AllowCollapse bool   `json:"allow_collapse"`
}

type invokeResponse struct {
	Output         string            `json:"output"`
	E1Output       string            `json:"e1_output"`
	E2Output       string            `json:"e2_output"`
	ResonanceState string            `json:"resonance_state"`
	PhaseDeltaDeg  float64           `json:"phase_delta_deg"`
	E1Weight       float64           `json:"e1_weight"`
	E2Weight       float64           `json:"e2_weight"`
	ActiveSpirits  []archetypes.Spirit `json:"active_spirits"`
}

func (s *server) handleInvoke(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	var req invokeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.Intent == "" {
		http.Error(w, "intent required", http.StatusBadRequest)
		return
	}
	allowCollapse := s.allowCollapse || req.AllowCollapse

	res, err := s.field.Invoke(r.Context(), req.Intent, req.Context, allowCollapse)
	if err != nil {
		http.Error(w, "invoke error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invokeResponse{
		Output:         res.Output,
		E1Output:       res.E1Output,
		E2Output:       res.E2Output,
		ResonanceState: string(res.ResonanceState),
		PhaseDeltaDeg:  res.PhaseDeltaDeg,
		E1Weight:       res.E1Weight,
		E2Weight:       res.E2Weight,
		ActiveSpirits:  res.ActiveSpirits,
	})
}

func (s *server) handleSpirits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(archetypes.All72)
}

func (s *server) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":         "ok",
		"started_at":     s.startedAt.Format(time.RFC3339),
		"uptime_seconds": int(time.Since(s.startedAt).Seconds()),
		"e1_model":       s.field.E1Model,
		"e2_model":       s.field.E2Model,
		"spirit_count":   len(archetypes.All72),
		"allow_collapse": s.allowCollapse,
	})
}

func main() {
	addr          := flag.String("addr", ":8090", "listen address")
	e1Model       := flag.String("e1-model", "claude-sonnet-4-6", "E₁ Carrier model")
	e2Model       := flag.String("e2-model", "claude-haiku-4-5-20251001", "E₂ Explorer model")
	allowCollapse := flag.Bool("allow-collapse", false, "permit Collapse-corridor spirits")
	flag.Parse()

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		log.Fatal("ANTHROPIC_API_KEY is not set")
	}

	client := newAnthropicClient(apiKey)
	field := &archetypes.Field{
		E1Client: client,
		E2Client: client,
		E1Model:  *e1Model,
		E2Model:  *e2Model,
	}

	srv := &server{
		field:         field,
		allowCollapse: *allowCollapse,
		startedAt:     time.Now().UTC(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/invoke", srv.handleInvoke)
	mux.HandleFunc("/spirits", srv.handleSpirits)
	mux.HandleFunc("/status", srv.handleStatus)

	log.Printf("archetype-engine starting on %s (E1=%s E2=%s allow_collapse=%v)", *addr, *e1Model, *e2Model, *allowCollapse)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatal(err)
	}
}
