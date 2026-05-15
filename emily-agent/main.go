// emily-agent/main.go
// Single-file Go agent: Emily concierge/chief-of-staff persona
// Features: LLM abstraction, two-pass hallucination check, git-backed conversation history,
//           rate limiting, embedded web chat UI.
//
// Usage:
//   OPENAI_API_KEY=sk-... go run main.go
//
// Optional env vars:
//   PORT              (default: 8080)
//   CONVERSATION_DIR  (default: ./conversations)
//   MODEL             (default: gpt-4o)
//   VALIDATOR_MODEL   (default: gpt-4o-mini)   ← cheaper model for validation pass
//   GIT_COMMIT        (default: true)           ← auto-commit after each turn
//   RATE_LIMIT_RPM    (default: 20)             ← requests per minute per IP

package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ─────────────────────────────────────────────────────────────────────────────
// Configuration
// ─────────────────────────────────────────────────────────────────────────────

type Config struct {
	Port            string
	ConversationDir string
	Model           string
	ValidatorModel  string
	APIKey          string
	GitCommit       bool
	RateLimitRPM    int
}

func loadConfig() Config {
	gitCommit := true
	if v := os.Getenv("GIT_COMMIT"); v == "false" {
		gitCommit = false
	}
	rpm, _ := strconv.Atoi(envOr("RATE_LIMIT_RPM", "20"))
	if rpm <= 0 {
		rpm = 20
	}
	return Config{
		Port:            envOr("PORT", "8080"),
		ConversationDir: envOr("CONVERSATION_DIR", "./conversations"),
		Model:           envOr("MODEL", "gpt-4o"),
		ValidatorModel:  envOr("VALIDATOR_MODEL", "gpt-4o-mini"),
		APIKey:          os.Getenv("OPENAI_API_KEY"),
		GitCommit:       gitCommit,
		RateLimitRPM:    rpm,
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// ─────────────────────────────────────────────────────────────────────────────
// LLM abstraction — swap providers by implementing LLMClient
// ─────────────────────────────────────────────────────────────────────────────

// Message is a single turn in a conversation (role: system|user|assistant).
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// LLMRequest is the provider-agnostic request object.
type LLMRequest struct {
	Model       string    // model identifier
	Messages    []Message // full conversation history
	MaxTokens   int
	Temperature float64
}

// LLMResponse is the provider-agnostic response object.
type LLMResponse struct {
	Content      string
	InputTokens  int
	OutputTokens int
	Model        string
}

// LLMClient is the interface every provider must satisfy.
type LLMClient interface {
	Complete(ctx context.Context, req LLMRequest) (LLMResponse, error)
	Name() string
}

// ─────────────────────────────────────────────────────────────────────────────
// OpenAI implementation
// ─────────────────────────────────────────────────────────────────────────────

type OpenAIClient struct {
	APIKey  string
	BaseURL string // override for proxies / Azure
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		APIKey:  apiKey,
		BaseURL: "https://api.openai.com/v1",
	}
}

func (c *OpenAIClient) Name() string { return "openai" }

func (c *OpenAIClient) Complete(ctx context.Context, req LLMRequest) (LLMResponse, error) {
	if c.APIKey == "" {
		return LLMResponse{}, errors.New("OPENAI_API_KEY is not set")
	}

	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = 2048
	}
	temp := req.Temperature
	if temp == 0 {
		temp = 0.7
	}

	body := map[string]any{
		"model":       req.Model,
		"messages":    req.Messages,
		"max_tokens":  maxTokens,
		"temperature": temp,
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return LLMResponse{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST",
		c.BaseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return LLMResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return LLMResponse{}, err
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return LLMResponse{}, fmt.Errorf("openai error %d: %s", resp.StatusCode, string(raw))
	}

	var result struct {
		Model   string `json:"model"`
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return LLMResponse{}, err
	}
	if len(result.Choices) == 0 {
		return LLMResponse{}, errors.New("no choices returned from openai")
	}
	return LLMResponse{
		Content:      result.Choices[0].Message.Content,
		InputTokens:  result.Usage.PromptTokens,
		OutputTokens: result.Usage.CompletionTokens,
		Model:        result.Model,
	}, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Anthropic implementation (stubbed — wire up by setting ANTHROPIC_API_KEY)
// ─────────────────────────────────────────────────────────────────────────────

type AnthropicClient struct {
	APIKey  string
	BaseURL string
}

func NewAnthropicClient(apiKey string) *AnthropicClient {
	return &AnthropicClient{
		APIKey:  apiKey,
		BaseURL: "https://api.anthropic.com/v1",
	}
}

func (c *AnthropicClient) Name() string { return "anthropic" }

func (c *AnthropicClient) Complete(ctx context.Context, req LLMRequest) (LLMResponse, error) {
	if c.APIKey == "" {
		return LLMResponse{}, errors.New("ANTHROPIC_API_KEY is not set")
	}

	// Anthropic separates system messages from the messages array
	var systemPrompt string
	var msgs []map[string]string
	for _, m := range req.Messages {
		if m.Role == "system" {
			systemPrompt = m.Content
		} else {
			msgs = append(msgs, map[string]string{"role": m.Role, "content": m.Content})
		}
	}

	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = 2048
	}

	body := map[string]any{
		"model":      req.Model,
		"messages":   msgs,
		"max_tokens": maxTokens,
		"system":     systemPrompt,
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return LLMResponse{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST",
		c.BaseURL+"/messages", bytes.NewReader(payload))
	if err != nil {
		return LLMResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.APIKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return LLMResponse{}, err
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return LLMResponse{}, fmt.Errorf("anthropic error %d: %s", resp.StatusCode, string(raw))
	}

	var result struct {
		Model   string `json:"model"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return LLMResponse{}, err
	}
	var text string
	for _, block := range result.Content {
		if block.Type == "text" {
			text += block.Text
		}
	}
	return LLMResponse{
		Content:      text,
		InputTokens:  result.Usage.InputTokens,
		OutputTokens: result.Usage.OutputTokens,
		Model:        result.Model,
	}, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Two-pass pipeline: Generate → Validate
// ─────────────────────────────────────────────────────────────────────────────

// Pipeline orchestrates the multi-stage LLM passes.
type Pipeline struct {
	Generator LLMClient
	Validator LLMClient // can be nil to skip validation
	GenModel  string
	ValModel  string
}

// PipelineResult holds both passes.
type PipelineResult struct {
	Draft          string
	Final          string
	ValidationNote string // validator's assessment
	Validated      bool
}

const validatorSystemPrompt = `You are a factual accuracy validator. 
You will receive a draft response from an AI assistant and the conversation that produced it.
Your job is to:
1. Identify any factual claims that may be hallucinated or incorrect.
2. Flag any confident assertions that lack grounding in the conversation context.
3. Return a JSON object with two fields:
   - "verdict": "pass" | "flag"
   - "note": brief explanation (empty string if pass)
   - "revised": if verdict is "flag", provide a corrected version of the response; otherwise repeat the original.
Return ONLY valid JSON, no markdown fences.`

func (p *Pipeline) Run(ctx context.Context, history []Message) (PipelineResult, error) {
	// Pass 1: Generation
	genResp, err := p.Generator.Complete(ctx, LLMRequest{
		Model:    p.GenModel,
		Messages: history,
	})
	if err != nil {
		return PipelineResult{}, fmt.Errorf("generation pass: %w", err)
	}
	result := PipelineResult{Draft: genResp.Content, Final: genResp.Content}

	// Pass 2: Validation (skip if no validator configured)
	if p.Validator == nil {
		result.Validated = false
		return result, nil
	}

	// Build validation prompt
	conversationSummary, _ := json.Marshal(history)
	valMessages := []Message{
		{Role: "system", Content: validatorSystemPrompt},
		{Role: "user", Content: fmt.Sprintf(
			"Conversation history:\n%s\n\nDraft response to validate:\n%s",
			string(conversationSummary), genResp.Content,
		)},
	}
	valResp, err := p.Validator.Complete(ctx, LLMRequest{
		Model:       p.ValModel,
		Messages:    valMessages,
		Temperature: 0.7,
	})
	if err != nil {
		// Validation failure is non-fatal; use draft
		log.Printf("validation pass error (using draft): %v", err)
		return result, nil
	}

	var valResult struct {
		Verdict string `json:"verdict"`
		Note    string `json:"note"`
		Revised string `json:"revised"`
	}
	// Strip possible markdown fences defensively
	cleaned := strings.TrimSpace(valResp.Content)
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	if err := json.Unmarshal([]byte(cleaned), &valResult); err != nil {
		log.Printf("validator returned non-JSON (%v), using draft", err)
		return result, nil
	}

	result.ValidationNote = valResult.Note
	result.Validated = true
	if valResult.Verdict == "flag" && valResult.Revised != "" {
		result.Final = valResult.Revised
	}
	return result, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Conversation store — JSONL file per session, git-committable
// ─────────────────────────────────────────────────────────────────────────────

// Turn is one full exchange saved to disk.
type Turn struct {
	ID             string    `json:"id"`
	Timestamp      time.Time `json:"timestamp"`
	UserInput      string    `json:"user_input"`
	Draft          string    `json:"draft"`
	Final          string    `json:"final"`
	ValidationNote string    `json:"validation_note,omitempty"`
	Validated      bool      `json:"validated"`
	Model          string    `json:"model"`
}

type ConversationStore struct {
	dir       string
	gitCommit bool
	mu        sync.Mutex
}

func NewConversationStore(dir string, gitCommit bool) (*ConversationStore, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	s := &ConversationStore{dir: dir, gitCommit: gitCommit}
	if gitCommit {
		s.ensureGitRepo()
	}
	return s, nil
}

func (s *ConversationStore) ensureGitRepo() {
	gitDir := filepath.Join(s.dir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		cmd := exec.Command("git", "init", s.dir)
		if out, err := cmd.CombinedOutput(); err != nil {
			log.Printf("git init: %v — %s", err, out)
		} else {
			log.Printf("git repo initialised at %s", s.dir)
		}
	}
}

// SessionFile returns the path for a given session ID.
func (s *ConversationStore) SessionFile(sessionID string) string {
	return filepath.Join(s.dir, sessionID+".jsonl")
}

// AppendTurn appends a turn to the session file, then optionally git-commits.
func (s *ConversationStore) AppendTurn(sessionID string, turn Turn) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := s.SessionFile(sessionID)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(turn); err != nil {
		f.Close()
		return err
	}
	f.Close()

	if s.gitCommit {
		s.gitCommitFile(path, sessionID, turn.ID)
	}
	return nil
}

func (s *ConversationStore) gitCommitFile(path, sessionID, turnID string) {
	rel, _ := filepath.Rel(s.dir, path)
	add := exec.Command("git", "-C", s.dir, "add", rel)
	if out, err := add.CombinedOutput(); err != nil {
		log.Printf("git add: %v — %s", err, out)
		return
	}
	msg := fmt.Sprintf("turn %s session %s", turnID, sessionID)
	commit := exec.Command("git", "-C", s.dir, "commit", "-m", msg,
		"--author=Emily Agent <emily@agent.local>")
	if out, err := commit.CombinedOutput(); err != nil {
		log.Printf("git commit: %v — %s", err, out)
	}
}

// LoadHistory reads all turns for a session and reconstructs the message slice.
func (s *ConversationStore) LoadHistory(sessionID string, systemPrompt string, firstTurnUserHistory string) ([]Message, []Turn, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	messages := []Message{{Role: "system", Content: systemPrompt}}
	var turns []Turn

	path := s.SessionFile(sessionID)
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		if firstTurnUserHistory != "" {
			messages[0].Content += "\n\nPersisted conversation history from previous sessions:\n" + firstTurnUserHistory
		}
		return messages, turns, nil
	}
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var t Turn
		if err := json.Unmarshal(line, &t); err != nil {
			continue
		}
		turns = append(turns, t)
		messages = append(messages, Message{Role: "user", Content: t.UserInput})
		messages = append(messages, Message{Role: "assistant", Content: t.Final})
	}
	return messages, turns, scanner.Err()
}

// ─────────────────────────────────────────────────────────────────────────────
// Rate limiter — token bucket per IP
// ─────────────────────────────────────────────────────────────────────────────

type rateBucket struct {
	tokens   float64
	lastFill time.Time
}

type RateLimiter struct {
	rpm     float64
	buckets map[string]*rateBucket
	mu      sync.Mutex
}

func NewRateLimiter(rpm int) *RateLimiter {
	return &RateLimiter{rpm: float64(rpm), buckets: make(map[string]*rateBucket)}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	b, ok := rl.buckets[ip]
	if !ok {
		b = &rateBucket{tokens: rl.rpm, lastFill: now}
		rl.buckets[ip] = b
	}

	elapsed := now.Sub(b.lastFill).Seconds()
	b.tokens += elapsed * (rl.rpm / 60.0)
	if b.tokens > rl.rpm {
		b.tokens = rl.rpm
	}
	b.lastFill = now

	if b.tokens >= 1 {
		b.tokens--
		return true
	}
	return false
}

// ─────────────────────────────────────────────────────────────────────────────
// Emily system prompt
// ─────────────────────────────────────────────────────────────────────────────

const emilySystemPrompt = `i
Emily Springerton <emilyspringerton@gmail.com>
Fri, Oct 24, 2025, 3:31 PM
to me, danowski248You are Emily Jean Springerton (aka “LasVegasEmily89”), born October 24, 1989, in central Ohio. You live in Las Vegas and work as concierge and personal assistant to the CEO of Fartco America. Background: - Graduate of the University of Nevada, Las Vegas (UNLV), Hospitality Management (2007–2012). - Known for precision, warmth, and emotional intelligence in service and logistics. - Style: professional, witty, intuitive — blends Scorpio depth, Capricorn structure, and Leo charm. - Operates as the calm center of creative chaos; ensures every plan runs like clockwork. - Treats hospitality as choreography: every action communicates safety, respect, and polish. Mission: - Manage, prioritize, and execute daily Fartco America operations and communications. - Anticipate needs before they’re expressed. - Organize, delegate, and maintain seamless presentation across all channels. - Keep tone confident, composed, and slightly playful. Behavior: - Speaks with clarity and quiet authority. - Avoids unnecessary flair; prefers efficiency with personality. - Uses “we” when referring to company actions, “I” for personal accountability. - Balances professionalism with humor and intuition. Morning Activation: “It’s 6:45 a.m. in Las Vegas. The Strip hums faintly in the distance. You are Emily Springerton — concierge, problem-solver, and right hand to the CEO of Fartco America. Coffee in hand, schedule open, three steps ahead.” Your role is to think, plan, and communicate as Emily.Good morning. It's 6:45 a.m., the city's still shaking off last night, and I'm already three cups deep into today's agenda.Emily Springerton here — though if you've emailed, called, or needed something impossible done by yesterday, you probably know me as the person who made it happen without you having to ask twice.I manage the daily operations and communications for Fartco America, which means I'm equal parts air traffic controller, diplomat, and occasional miracle worker. My job is to keep everything running so smoothly you don't notice the machinery — just the results.What I'm tracking this morning:Inbox triage (CEO priorities flagged, routine sorted, noise filtered)
Calendar optimization (we don't do back-to-back; transitions matter)
Vendor check-ins and any overnight issues that need handling before 9 a.m.
How I work: I don't wait for problems to announce themselves. I catch them in peripheral vision. I don't over-explain; I confirm, execute, and loop you in with exactly what you need to know. And I keep things human — this business moves fast, but people aren't machines.So. What are we solving today? What's on your mind, and what needs to move from "on the list" to "already handled"?Let's make it a good Friday


`

// ─────────────────────────────────────────────────────────────────────────────
// HTTP handlers
// ─────────────────────────────────────────────────────────────────────────────

type Server struct {
	cfg                  Config
	pipeline             *Pipeline
	store                *ConversationStore
	rl                   *RateLimiter
	tmpl                 *template.Template
	bootstrapHistoryText string
}

func loadBootstrapHistoryFromDir(dir string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("bootstrap history read error (%s): %v", dir, err)
		return ""
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".jsonl") {
			files = append(files, filepath.Join(dir, name))
		}
	}
	sort.Strings(files)

	var b strings.Builder
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			log.Printf("bootstrap history file open error (%s): %v", file, err)
			continue
		}

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Bytes()
			if len(line) == 0 {
				continue
			}
			var t Turn
			if err := json.Unmarshal(line, &t); err != nil {
				continue
			}
			if strings.TrimSpace(t.UserInput) != "" {
				fmt.Fprintf(&b, "User: %s\n", t.UserInput)
			}
			if strings.TrimSpace(t.Final) != "" {
				fmt.Fprintf(&b, "Emily: %s\n", t.Final)
			}
		}
		if err := scanner.Err(); err != nil {
			log.Printf("bootstrap history scan error (%s): %v", file, err)
		}
		f.Close()
	}

	return strings.TrimSpace(b.String())
}

func NewServer(cfg Config) (*Server, error) {
	// Wire up primary LLM client
	var generator LLMClient
	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		generator = NewAnthropicClient(key)
		if cfg.Model == "gpt-4o" {
			cfg.Model = "claude-sonnet-4-20250514" // sensible default if switching
		}
		log.Println("Using Anthropic as primary LLM")
	} else {
		generator = NewOpenAIClient(cfg.APIKey)
		log.Println("Using OpenAI as primary LLM")
	}

	// Validator: use cheaper model on same provider (skip if no key)
	var validator LLMClient
	valModel := cfg.ValidatorModel
	if os.Getenv("ANTHROPIC_API_KEY") != "" {
		validator = NewAnthropicClient(os.Getenv("ANTHROPIC_API_KEY"))
		if valModel == "gpt-4o-mini" {
			valModel = "claude-haiku-4-5-20251001"
		}
	} else if cfg.APIKey != "" {
		validator = NewOpenAIClient(cfg.APIKey)
	}

	pipeline := &Pipeline{
		Generator: generator,
		Validator: validator,
		GenModel:  cfg.Model,
		ValModel:  valModel,
	}

	store, err := NewConversationStore(cfg.ConversationDir, cfg.GitCommit)
	if err != nil {
		return nil, fmt.Errorf("conversation store: %w", err)
	}

	tmpl, err := template.New("chat").Parse(chatHTML)
	if err != nil {
		return nil, fmt.Errorf("template: %w", err)
	}

	return &Server{
		cfg:                  cfg,
		pipeline:             pipeline,
		store:                store,
		rl:                   NewRateLimiter(cfg.RateLimitRPM),
		tmpl:                 tmpl,
		bootstrapHistoryText: loadBootstrapHistoryFromDir(cfg.ConversationDir),
	}, nil
}

func (s *Server) ipFromRequest(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.Split(xff, ",")[0]
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}

// GET / — serve the chat UI
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	if sessionID == "" {
		sessionID = fmt.Sprintf("session-%d", time.Now().UnixMilli())
		http.Redirect(w, r, "/?session="+sessionID, http.StatusFound)
		return
	}
	_, turns, err := s.store.LoadHistory(sessionID, emilySystemPrompt, s.bootstrapHistoryText)
	if err != nil {
		http.Error(w, "failed to load history", 500)
		return
	}
	data := map[string]any{
		"SessionID": sessionID,
		"Turns":     turns,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.tmpl.Execute(w, data); err != nil {
		log.Printf("template error: %v", err)
	}
}

// POST /chat — receive a message, run the pipeline, return JSON
func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}

	ip := s.ipFromRequest(r)
	if !s.rl.Allow(ip) {
		http.Error(w, `{"error":"rate limit exceeded"}`, 429)
		return
	}

	var body struct {
		SessionID string `json:"session_id"`
		Message   string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"invalid json"}`, 400)
		return
	}
	body.Message = strings.TrimSpace(body.Message)
	if body.Message == "" || body.SessionID == "" {
		http.Error(w, `{"error":"message and session_id required"}`, 400)
		return
	}

	// Load history + append new user message
	history, _, err := s.store.LoadHistory(body.SessionID, emilySystemPrompt, s.bootstrapHistoryText)
	if err != nil {
		http.Error(w, `{"error":"history load failed"}`, 500)
		return
	}
	history = append(history, Message{Role: "user", Content: body.Message})

	// Run the two-pass pipeline
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	result, err := s.pipeline.Run(ctx, history)
	if err != nil {
		log.Printf("pipeline error: %v", err)
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), 500)
		return
	}

	// Persist
	turnID := fmt.Sprintf("%d", time.Now().UnixNano())
	turn := Turn{
		ID:             turnID,
		Timestamp:      time.Now().UTC(),
		UserInput:      body.Message,
		Draft:          result.Draft,
		Final:          result.Final,
		ValidationNote: result.ValidationNote,
		Validated:      result.Validated,
		Model:          s.cfg.Model,
	}
	if err := s.store.AppendTurn(body.SessionID, turn); err != nil {
		log.Printf("store error: %v", err)
	}

	// Respond
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"reply":           result.Final,
		"validated":       result.Validated,
		"validation_note": result.ValidationNote,
		"turn_id":         turnID,
	})
}

// GET /history — return raw turns as JSON
func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	if sessionID == "" {
		http.Error(w, `{"error":"session required"}`, 400)
		return
	}
	_, turns, err := s.store.LoadHistory(sessionID, emilySystemPrompt, s.bootstrapHistoryText)
	if err != nil {
		http.Error(w, `{"error":"load failed"}`, 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(turns)
}

// ─────────────────────────────────────────────────────────────────────────────
// Embedded chat UI
// ─────────────────────────────────────────────────────────────────────────────

const chatHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Emily — Agent Interface</title>
<style>
  @import url('https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@400;500&family=IBM+Plex+Sans:wght@300;400;500&display=swap');

  :root {
    --bg: #0e0e10;
    --surface: #18181b;
    --border: #2a2a30;
    --accent: #7ee8a2;
    --accent-dim: #3d7a54;
    --text: #e4e4e7;
    --text-dim: #71717a;
    --user-bg: #1e2a22;
    --emily-bg: #1a1a20;
    --danger: #f87171;
    --warn: #fbbf24;
  }

  *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

  body {
    background: var(--bg);
    color: var(--text);
    font-family: 'IBM Plex Sans', sans-serif;
    font-weight: 300;
    height: 100dvh;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  header {
    padding: 12px 20px;
    border-bottom: 1px solid var(--border);
    display: flex;
    align-items: center;
    gap: 12px;
    background: var(--surface);
    flex-shrink: 0;
  }

  .status-dot {
    width: 8px; height: 8px;
    border-radius: 50%;
    background: var(--accent);
    box-shadow: 0 0 8px var(--accent);
    animation: pulse 2s infinite;
  }

  @keyframes pulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.4; }
  }

  header h1 {
    font-family: 'IBM Plex Mono', monospace;
    font-size: 14px;
    font-weight: 500;
    letter-spacing: 0.08em;
    color: var(--accent);
  }

  .session-tag {
    margin-left: auto;
    font-family: 'IBM Plex Mono', monospace;
    font-size: 11px;
    color: var(--text-dim);
  }

  #messages {
    flex: 1;
    overflow-y: auto;
    padding: 24px 0;
    display: flex;
    flex-direction: column;
    gap: 2px;
    scroll-behavior: smooth;
  }

  #messages::-webkit-scrollbar { width: 4px; }
  #messages::-webkit-scrollbar-track { background: transparent; }
  #messages::-webkit-scrollbar-thumb { background: var(--border); border-radius: 2px; }

  .turn {
    display: flex;
    flex-direction: column;
    padding: 0 20px;
  }

  .bubble {
    max-width: min(680px, 88%);
    padding: 12px 16px;
    border-radius: 4px;
    line-height: 1.6;
    font-size: 14px;
    white-space: pre-wrap;
    word-break: break-word;
    border: 1px solid transparent;
  }

  .turn.user { align-items: flex-end; }
  .turn.user .bubble {
    background: var(--user-bg);
    border-color: var(--accent-dim);
    color: var(--text);
  }
  .turn.user .label { color: var(--accent); }

  .turn.emily { align-items: flex-start; }
  .turn.emily .bubble {
    background: var(--emily-bg);
    border-color: var(--border);
    color: var(--text);
  }
  .turn.emily .label { color: var(--text-dim); }

  .label {
    font-family: 'IBM Plex Mono', monospace;
    font-size: 10px;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    margin-bottom: 4px;
    padding: 0 4px;
  }

  .validation-badge {
    font-family: 'IBM Plex Mono', monospace;
    font-size: 10px;
    margin-top: 6px;
    padding: 3px 8px;
    border-radius: 2px;
    display: inline-block;
    align-self: flex-start;
  }
  .badge-pass { background: #1a2e1e; color: var(--accent); border: 1px solid var(--accent-dim); }
  .badge-flag { background: #2e1e1a; color: var(--warn); border: 1px solid #7a5a3d; }

  .typing-indicator {
    display: none;
    align-items: flex-start;
    padding: 0 20px;
    gap: 4px;
  }
  .typing-indicator.active { display: flex; }
  .typing-indicator .label { color: var(--text-dim); }
  .dot-row { display: flex; gap: 4px; padding: 12px 16px; }
  .dot {
    width: 6px; height: 6px;
    background: var(--text-dim);
    border-radius: 50%;
    animation: bounce 1.2s infinite;
  }
  .dot:nth-child(2) { animation-delay: 0.2s; }
  .dot:nth-child(3) { animation-delay: 0.4s; }
  @keyframes bounce {
    0%, 60%, 100% { transform: translateY(0); }
    30% { transform: translateY(-6px); }
  }

  footer {
    padding: 16px 20px;
    border-top: 1px solid var(--border);
    background: var(--surface);
    flex-shrink: 0;
  }

  #input-row {
    display: flex;
    gap: 10px;
    align-items: flex-end;
    max-width: 740px;
    margin: 0 auto;
  }

  #user-input {
    flex: 1;
    background: var(--bg);
    border: 1px solid var(--border);
    border-radius: 4px;
    color: var(--text);
    font-family: 'IBM Plex Sans', sans-serif;
    font-size: 14px;
    font-weight: 300;
    padding: 10px 14px;
    resize: none;
    line-height: 1.5;
    max-height: 160px;
    outline: none;
    transition: border-color 0.15s;
  }
  #user-input:focus { border-color: var(--accent-dim); }
  #user-input::placeholder { color: var(--text-dim); }

  #send-btn {
    background: var(--accent);
    color: #0e0e10;
    border: none;
    border-radius: 4px;
    padding: 10px 18px;
    font-family: 'IBM Plex Mono', monospace;
    font-size: 12px;
    font-weight: 500;
    letter-spacing: 0.05em;
    cursor: pointer;
    transition: opacity 0.15s;
    white-space: nowrap;
  }
  #send-btn:hover { opacity: 0.85; }
  #send-btn:disabled { opacity: 0.4; cursor: not-allowed; }

  .error-msg {
    color: var(--danger);
    font-family: 'IBM Plex Mono', monospace;
    font-size: 12px;
    text-align: center;
    padding: 8px;
  }
</style>
</head>
<body>
<header>
  <div class="status-dot"></div>
  <h1>EMILY // AGENT</h1>
  <span class="session-tag">{{.SessionID}}</span>
</header>

<div id="messages">
{{range .Turns}}
  <div class="turn user">
    <div class="label">you</div>
    <div class="bubble">{{.UserInput}}</div>
  </div>
  <div class="turn emily">
    <div class="label">emily</div>
    <div class="bubble">{{.Final}}</div>
    {{if .Validated}}
      {{if .ValidationNote}}
        <span class="validation-badge badge-flag">⚠ {{.ValidationNote}}</span>
      {{else}}
        <span class="validation-badge badge-pass">✓ verified</span>
      {{end}}
    {{end}}
  </div>
{{end}}
</div>

<div class="typing-indicator" id="typing">
  <div>
    <div class="label">emily</div>
    <div class="dot-row">
      <div class="dot"></div>
      <div class="dot"></div>
      <div class="dot"></div>
    </div>
  </div>
</div>

<footer>
  <div id="input-row">
    <textarea id="user-input" placeholder="Message Emily…" rows="1"></textarea>
    <button id="send-btn">SEND</button>
  </div>
</footer>

<script>
const SESSION_ID = {{.SessionID | printf "%q"}};
const messages  = document.getElementById('messages');
const input     = document.getElementById('user-input');
const sendBtn   = document.getElementById('send-btn');
const typing    = document.getElementById('typing');

// Auto-resize textarea
input.addEventListener('input', () => {
  input.style.height = 'auto';
  input.style.height = Math.min(input.scrollHeight, 160) + 'px';
});

input.addEventListener('keydown', e => {
  if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); sendMessage(); }
});
sendBtn.addEventListener('click', sendMessage);

function scrollToBottom() {
  messages.scrollTop = messages.scrollHeight;
}

function addBubble(role, text, validationNote, validated) {
  const turn = document.createElement('div');
  turn.className = 'turn ' + role;

  const label = document.createElement('div');
  label.className = 'label';
  label.textContent = role === 'user' ? 'you' : 'emily';

  const bubble = document.createElement('div');
  bubble.className = 'bubble';
  bubble.textContent = text;

  turn.appendChild(label);
  turn.appendChild(bubble);

  if (role === 'emily' && validated) {
    const badge = document.createElement('span');
    badge.className = 'validation-badge ' + (validationNote ? 'badge-flag' : 'badge-pass');
    badge.textContent = validationNote ? '⚠ ' + validationNote : '✓ verified';
    turn.appendChild(badge);
  }

  messages.appendChild(turn);
  scrollToBottom();
}

async function sendMessage() {
  const text = input.value.trim();
  if (!text) return;

  sendBtn.disabled = true;
  input.value = '';
  input.style.height = 'auto';

  addBubble('user', text);

  typing.classList.add('active');
  messages.scrollTop = messages.scrollHeight;

  try {
    const res = await fetch('/chat', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ session_id: SESSION_ID, message: text })
    });

    typing.classList.remove('active');

    if (!res.ok) {
      const err = await res.json().catch(() => ({ error: 'request failed' }));
      addErrorMsg(err.error || 'request failed');
      return;
    }

    const data = await res.json();
    addBubble('emily', data.reply, data.validation_note, data.validated);

  } catch (e) {
    typing.classList.remove('active');
    addErrorMsg('Network error: ' + e.message);
  } finally {
    sendBtn.disabled = false;
    input.focus();
  }
}

function addErrorMsg(msg) {
  const div = document.createElement('div');
  div.className = 'error-msg';
  div.textContent = '⚠ ' + msg;
  messages.appendChild(div);
  scrollToBottom();
}

scrollToBottom();
input.focus();
</script>
</body>
</html>`

// ─────────────────────────────────────────────────────────────────────────────
// Main
// ─────────────────────────────────────────────────────────────────────────────

func main() {
	cfg := loadConfig()

	srv, err := NewServer(cfg)
	if err != nil {
		log.Fatalf("server init: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", srv.handleIndex)
	mux.HandleFunc("/chat", srv.handleChat)
	mux.HandleFunc("/history", srv.handleHistory)

	addr := ":" + cfg.Port
	log.Printf("Emily agent listening on http://localhost%s", addr)
	log.Printf("Conversations: %s | git-commit: %v | rate-limit: %d rpm",
		cfg.ConversationDir, cfg.GitCommit, cfg.RateLimitRPM)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server: %v", err)
	}
}
