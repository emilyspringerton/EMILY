// api_gpt2.go — POST /api/v1/gpt2/generate
//
// Proxies to the Emily Prime GPT-2 inference server (scripts/serve.py on :8088).
// The Python server loads the fine-tuned checkpoint once at startup and keeps it
// in memory. This endpoint forwards the request and returns the response.
//
// The Python server must be running:
//   python3 /home/fatbaby/gpt2-alpine-c/scripts/serve.py

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const gpt2ServerURL = "http://localhost:8088"

type gpt2GenerateRequest struct {
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
}

type gpt2GenerateResponse struct {
	Text   string `json:"text"`
	Prompt string `json:"prompt"`
	Tokens int    `json:"tokens"`
	Model  string `json:"model"`
	Error  string `json:"error,omitempty"`
}

var gpt2Client = &http.Client{Timeout: 60 * time.Second}

// handleGPT2Generate proxies POST /api/v1/gpt2/generate to the Python inference server.
func (s *Server) handleGPT2Generate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req gpt2GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusBadRequest)
		return
	}
	if req.Prompt == "" {
		req.Prompt = "Emily Prime:"
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = 100
	}
	if req.Temperature == 0 {
		req.Temperature = 0.8
	}

	body, _ := json.Marshal(req)
	upstream, err := http.NewRequestWithContext(r.Context(), http.MethodPost,
		gpt2ServerURL+"/generate", bytes.NewReader(body))
	if err != nil {
		http.Error(w, `{"error":"failed to build upstream request"}`, http.StatusInternalServerError)
		return
	}
	upstream.Header.Set("Content-Type", "application/json")

	resp, err := gpt2Client.Do(upstream)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"gpt2 server unavailable: %v","hint":"run: python3 /home/fatbaby/gpt2-alpine-c/scripts/serve.py"}`, err),
			http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(raw)
}

// handleGPT2Health proxies GET /api/v1/gpt2/health to check if the inference server is up.
func (s *Server) handleGPT2Health(w http.ResponseWriter, r *http.Request) {
	resp, err := gpt2Client.Get(gpt2ServerURL + "/health")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, `{"status":"unavailable","error":%q,"hint":"run: python3 /home/fatbaby/gpt2-alpine-c/scripts/serve.py"}`+"\n", err.Error())
		return
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(raw)
}
