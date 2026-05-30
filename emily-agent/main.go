// emily-agent/main.go
// Single-file Go agent: Emily concierge/chief-of-staff persona
//
// Features:
//   - LLM provider abstraction (OpenAI + Anthropic, swap via env)
//   - Native tool-calling with agentic loop (no regex)
//   - Tools: git_list_files, git_read_file (sandboxed to conversation dir)
//   - Two-pass hallucination validation after tool loop completes
//   - Git-backed JSONL conversation history (auto-commit each turn)
//   - Token-bucket rate limiter per IP
//   - Embedded web chat UI with tool-call activity shown inline
//
// Usage:
//   OPENAI_API_KEY=sk-...  go run main.go
//   ANTHROPIC_API_KEY=sk-ant-...  go run main.go   (auto-detected)

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
	"strconv"
	"strings"
	"sync"
	"time"
)

// -----------------------------------------------------------------------------
// Configuration
// -----------------------------------------------------------------------------

type Config struct {
	Port            string
	ConversationDir string
	Model           string
	ValidatorModel  string
	APIKey          string
	GitCommit       bool
	RateLimitRPM    int
	MaxToolIters    int
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
	maxIters, _ := strconv.Atoi(envOr("MAX_TOOL_ITERS", "10"))
	if maxIters <= 0 {
		maxIters = 10
	}
	return Config{
		Port:            envOr("PORT", "8080"),
		ConversationDir: envOr("CONVERSATION_DIR", "./conversations"),
		Model:           envOr("MODEL", "gpt-4o-mini"),
		ValidatorModel:  envOr("VALIDATOR_MODEL", "gpt-4o-mini"),
		APIKey:          os.Getenv("OPENAI_API_KEY"),
		GitCommit:       gitCommit,
		RateLimitRPM:    rpm,
		MaxToolIters:    maxIters,
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// -----------------------------------------------------------------------------
// Tool definitions -- the schema the LLM sees
// -----------------------------------------------------------------------------

// ToolDef is the provider-agnostic tool schema sent to the LLM.
type ToolDef struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  ToolParameters `json:"parameters"`
}

type ToolParameters struct {
	Type       string                    `json:"type"` // always "object"
	Properties map[string]ToolPropSchema `json:"properties"`
	Required   []string                  `json:"required,omitempty"`
}

type ToolPropSchema struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// ToolCall is a request from the LLM to invoke a tool.
type ToolCall struct {
	ID       string          // provider call ID (used to route the result back)
	Name     string          // tool name
	ArgsJSON json.RawMessage // raw JSON arguments
}

// ToolResult is the Go-side result of executing a tool.
type ToolResult struct {
	CallID  string
	Name    string
	Content string // plain text handed back to the LLM
	IsError bool
}

// -----------------------------------------------------------------------------
// LLM abstraction
// -----------------------------------------------------------------------------

// Message is a single turn. Role: system | user | assistant | tool
type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCallID string     `json:"tool_call_id,omitempty"` // for role=tool
	ToolName   string     `json:"tool_name,omitempty"`    // for role=tool
	ToolCalls  []ToolCall `json:"-"`                      // for role=assistant with calls
}

// LLMRequest is the provider-agnostic request.
type LLMRequest struct {
	Model       string
	Messages    []Message
	Tools       []ToolDef // empty = no tools exposed
	MaxTokens   int
	Temperature float64
}

// LLMResponse is the provider-agnostic response.
type LLMResponse struct {
	Content      string     // non-empty when the LLM produced text
	ToolCalls    []ToolCall // non-empty when the LLM wants to call tools
	InputTokens  int
	OutputTokens int
	Model        string
}

func (r LLMResponse) HasToolCalls() bool { return len(r.ToolCalls) > 0 }

// LLMClient is the interface every provider must satisfy.
type LLMClient interface {
	Complete(ctx context.Context, req LLMRequest) (LLMResponse, error)
	Name() string
}

// -----------------------------------------------------------------------------
// OpenAI implementation
// -----------------------------------------------------------------------------

type OpenAIClient struct {
	APIKey  string
	BaseURL string
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{APIKey: apiKey, BaseURL: "https://api.openai.com/v1"}
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

	type oaiMsg struct {
		Role       string `json:"role"`
		Content    any    `json:"content"`
		ToolCallID string `json:"tool_call_id,omitempty"`
		Name       string `json:"name,omitempty"`
		ToolCalls  any    `json:"tool_calls,omitempty"`
	}
	var msgs []oaiMsg
	for _, m := range req.Messages {
		switch m.Role {
		case "tool":
			msgs = append(msgs, oaiMsg{Role: "tool", Content: m.Content,
				ToolCallID: m.ToolCallID, Name: m.ToolName})
		case "assistant":
			if len(m.ToolCalls) > 0 {
				type oaiTC struct {
					ID       string `json:"id"`
					Type     string `json:"type"`
					Function struct {
						Name      string `json:"name"`
						Arguments string `json:"arguments"`
					} `json:"function"`
				}
				var tcs []oaiTC
				for _, tc := range m.ToolCalls {
					var ot oaiTC
					ot.ID = tc.ID
					ot.Type = "function"
					ot.Function.Name = tc.Name
					ot.Function.Arguments = string(tc.ArgsJSON)
					tcs = append(tcs, ot)
				}
				msgs = append(msgs, oaiMsg{Role: "assistant", Content: nil, ToolCalls: tcs})
			} else {
				msgs = append(msgs, oaiMsg{Role: "assistant", Content: m.Content})
			}
		default:
			msgs = append(msgs, oaiMsg{Role: m.Role, Content: m.Content})
		}
	}

	body := map[string]any{
		"model": req.Model, "messages": msgs,
		"max_tokens": maxTokens, "temperature": temp,
	}
	if len(req.Tools) > 0 {
		type oaiFn struct {
			Name        string         `json:"name"`
			Description string         `json:"description"`
			Parameters  ToolParameters `json:"parameters"`
		}
		type oaiTool struct {
			Type     string `json:"type"`
			Function oaiFn  `json:"function"`
		}
		var tools []oaiTool
		for _, t := range req.Tools {
			tools = append(tools, oaiTool{Type: "function",
				Function: oaiFn{Name: t.Name, Description: t.Description, Parameters: t.Parameters}})
		}
		body["tools"] = tools
	}

	payload, _ := json.Marshal(body)
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
		return LLMResponse{}, fmt.Errorf("openai %d: %s", resp.StatusCode, raw)
	}

	var result struct {
		Model   string `json:"model"`
		Choices []struct {
			Message struct {
				Content   *string `json:"content"`
				ToolCalls []struct {
					ID       string `json:"id"`
					Function struct {
						Name      string          `json:"name"`
						Arguments json.RawMessage `json:"arguments"`
					} `json:"function"`
				} `json:"tool_calls"`
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
		return LLMResponse{}, errors.New("no choices from openai")
	}
	out := LLMResponse{
		InputTokens:  result.Usage.PromptTokens,
		OutputTokens: result.Usage.CompletionTokens,
		Model:        result.Model,
	}
	msg := result.Choices[0].Message
	if msg.Content != nil {
		out.Content = *msg.Content
	}
	for _, tc := range msg.ToolCalls {
		out.ToolCalls = append(out.ToolCalls, ToolCall{
			ID: tc.ID, Name: tc.Function.Name, ArgsJSON: tc.Function.Arguments,
		})
	}
	return out, nil
}

// -----------------------------------------------------------------------------
// Anthropic implementation
// -----------------------------------------------------------------------------

type AnthropicClient struct {
	APIKey  string
	BaseURL string
}

func NewAnthropicClient(apiKey string) *AnthropicClient {
	return &AnthropicClient{APIKey: apiKey, BaseURL: "https://api.anthropic.com/v1"}
}

func (c *AnthropicClient) Name() string { return "anthropic" }

func (c *AnthropicClient) Complete(ctx context.Context, req LLMRequest) (LLMResponse, error) {
	if c.APIKey == "" {
		return LLMResponse{}, errors.New("ANTHROPIC_API_KEY is not set")
	}
	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = 2048
	}

	var systemPrompt string
	var msgs []map[string]any
	for _, m := range req.Messages {
		switch m.Role {
		case "system":
			systemPrompt = m.Content
		case "user":
			msgs = append(msgs, map[string]any{"role": "user", "content": m.Content})
		case "assistant":
			if len(m.ToolCalls) > 0 {
				var blocks []map[string]any
				for _, tc := range m.ToolCalls {
					var args any
					_ = json.Unmarshal(tc.ArgsJSON, &args)
					blocks = append(blocks, map[string]any{
						"type": "tool_use", "id": tc.ID, "name": tc.Name, "input": args,
					})
				}
				msgs = append(msgs, map[string]any{"role": "assistant", "content": blocks})
			} else {
				msgs = append(msgs, map[string]any{"role": "assistant", "content": m.Content})
			}
		case "tool":
			msgs = append(msgs, map[string]any{
				"role": "user",
				"content": []map[string]any{{
					"type": "tool_result", "tool_use_id": m.ToolCallID, "content": m.Content,
				}},
			})
		}
	}

	body := map[string]any{
		"model": req.Model, "messages": msgs,
		"max_tokens": maxTokens, "system": systemPrompt,
	}
	if len(req.Tools) > 0 {
		var tools []map[string]any
		for _, t := range req.Tools {
			tools = append(tools, map[string]any{
				"name": t.Name, "description": t.Description, "input_schema": t.Parameters,
			})
		}
		body["tools"] = tools
	}

	payload, _ := json.Marshal(body)
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
		return LLMResponse{}, fmt.Errorf("anthropic %d: %s", resp.StatusCode, raw)
	}

	var result struct {
		Model   string `json:"model"`
		Content []struct {
			Type  string          `json:"type"`
			Text  string          `json:"text"`
			ID    string          `json:"id"`
			Name  string          `json:"name"`
			Input json.RawMessage `json:"input"`
		} `json:"content"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return LLMResponse{}, err
	}
	out := LLMResponse{
		InputTokens:  result.Usage.InputTokens,
		OutputTokens: result.Usage.OutputTokens,
		Model:        result.Model,
	}
	for _, block := range result.Content {
		switch block.Type {
		case "text":
			out.Content += block.Text
		case "tool_use":
			out.ToolCalls = append(out.ToolCalls, ToolCall{
				ID: block.ID, Name: block.Name, ArgsJSON: block.Input,
			})
		}
	}
	return out, nil
}

// -----------------------------------------------------------------------------
// Tool dispatcher
// -----------------------------------------------------------------------------

type ToolFunc func(args map[string]any) (string, error)

type ToolDispatcher struct {
	defs     []ToolDef
	handlers map[string]ToolFunc
}

func NewToolDispatcher() *ToolDispatcher {
	return &ToolDispatcher{handlers: make(map[string]ToolFunc)}
}

func (d *ToolDispatcher) Register(def ToolDef, fn ToolFunc) {
	d.defs = append(d.defs, def)
	d.handlers[def.Name] = fn
}

func (d *ToolDispatcher) Dispatch(tc ToolCall) ToolResult {
	fn, ok := d.handlers[tc.Name]
	if !ok {
		return ToolResult{CallID: tc.ID, Name: tc.Name,
			Content: fmt.Sprintf("unknown tool: %s", tc.Name), IsError: true}
	}
	var args map[string]any
	if err := json.Unmarshal(tc.ArgsJSON, &args); err != nil {
		return ToolResult{CallID: tc.ID, Name: tc.Name,
			Content: "bad arguments: " + err.Error(), IsError: true}
	}
	content, err := fn(args)
	if err != nil {
		return ToolResult{CallID: tc.ID, Name: tc.Name, Content: err.Error(), IsError: true}
	}
	return ToolResult{CallID: tc.ID, Name: tc.Name, Content: content}
}

func (d *ToolDispatcher) Defs() []ToolDef { return d.defs }

// -----------------------------------------------------------------------------
// Git tools -- sandboxed to the conversation directory
// -----------------------------------------------------------------------------

// safePath resolves p relative to root and rejects path traversal.
func safePath(root, p string) (string, error) {
	abs := filepath.Join(root, filepath.Clean("/"+p))
	rel, err := filepath.Rel(root, abs)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("path %q is outside the allowed directory", p)
	}
	return abs, nil
}

func registerGitTools(d *ToolDispatcher, repoDir string) {

	// git_list_files
	d.Register(ToolDef{
		Name:        "git_list_files",
		Description: "List files tracked by git in the conversation repository. Use this to discover what session files exist before reading them.",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolPropSchema{
				"subdir": {
					Type:        "string",
					Description: "Optional subdirectory within the repo to list. Defaults to repo root.",
				},
			},
		},
	}, func(args map[string]any) (string, error) {
		subdir := "."
		if v, ok := args["subdir"].(string); ok && v != "" {
			subdir = v
		}
		absSubdir, err := safePath(repoDir, subdir)
		if err != nil {
			return "", err
		}
		// git ls-files only shows tracked files -- safe by design
		out, err := exec.Command("git", "-C", repoDir, "ls-files", absSubdir).Output()
		if err != nil {
			// Fallback before first commit: list directory
			entries, readErr := os.ReadDir(absSubdir)
			if readErr != nil {
				return "", fmt.Errorf("git ls-files failed: %v", err)
			}
			var lines []string
			for _, e := range entries {
				if !strings.HasPrefix(e.Name(), ".") {
					lines = append(lines, e.Name())
				}
			}
			if len(lines) == 0 {
				return "(no files yet)", nil
			}
			return strings.Join(lines, "\n"), nil
		}
		result := strings.TrimSpace(string(out))
		if result == "" {
			return "(no tracked files yet -- conversations will appear here after the first message is saved)", nil
		}
		return result, nil
	})

	// write_file — writes a proposal into proposals/ within the conversation repo.
	// This is how Emily persists RSI-generated improvements for human review.
	d.Register(ToolDef{
		Name:        "write_file",
		Description: "Write a text file to the proposals/ directory in the conversation repository. Use this to persist RSI-generated code improvements, architectural proposals, specifications, or analysis so they can be reviewed and applied. Files are git-committed automatically.",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolPropSchema{
				"path": {
					Type:        "string",
					Description: "Filename within proposals/, e.g. 'collector-backoff.go' or 'arxiv-source.go'. Subdirectories are allowed: 'analysis/quality-calibration.md'",
				},
				"content": {
					Type:        "string",
					Description: "Complete text content to write. For code proposals, include the full file — not a diff.",
				},
			},
			Required: []string{"path", "content"},
		},
	}, func(args map[string]any) (string, error) {
		p, ok := args["path"].(string)
		if !ok || strings.TrimSpace(p) == "" {
			return "", errors.New("path is required")
		}
		content, ok := args["content"].(string)
		if !ok {
			return "", errors.New("content is required")
		}
		if len(content) > 512*1024 {
			return "", fmt.Errorf("content too large: %d bytes exceeds 512KB limit", len(content))
		}
		proposalsDir := filepath.Join(repoDir, "proposals")
		if err := os.MkdirAll(proposalsDir, 0o755); err != nil {
			return "", fmt.Errorf("create proposals dir: %w", err)
		}
		abs, err := safePath(proposalsDir, p)
		if err != nil {
			return "", err
		}
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			return "", fmt.Errorf("create subdirectory: %w", err)
		}
		if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
			return "", fmt.Errorf("write file: %w", err)
		}
		rel, _ := filepath.Rel(repoDir, abs)
		if out, err := exec.Command("git", "-C", repoDir, "add", rel).CombinedOutput(); err != nil {
			log.Printf("git add proposals: %v -- %s", err, out)
		}
		msg := fmt.Sprintf("proposal: %s", p)
		if out, err := exec.Command("git", "-C", repoDir, "commit", "-m", msg,
			"--author=Emily Agent <emily@agent.local>").CombinedOutput(); err != nil {
			log.Printf("git commit proposal: %v -- %s", err, out)
		}
		return fmt.Sprintf("wrote %d bytes to proposals/%s", len(content), p), nil
	})

	// git_read_file
	d.Register(ToolDef{
		Name:        "git_read_file",
		Description: "Read the contents of a file in the conversation git repository. Call git_list_files first to discover valid paths. Session history files are named <session-id>.jsonl and each line is a JSON object representing one conversation turn.",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolPropSchema{
				"path": {
					Type:        "string",
					Description: "Path to the file relative to the repository root, e.g. 'session-1234567890.jsonl'",
				},
				"max_bytes": {
					Type:        "string",
					Description: "Maximum bytes to return (default 8192). Use a smaller number for large files to avoid overwhelming the context.",
				},
			},
			Required: []string{"path"},
		},
	}, func(args map[string]any) (string, error) {
		p, ok := args["path"].(string)
		if !ok || p == "" {
			return "", errors.New("path is required")
		}
		abs, err := safePath(repoDir, p)
		if err != nil {
			return "", err
		}

		maxBytes := 8192
		if v, ok := args["max_bytes"].(string); ok {
			if n, convErr := strconv.Atoi(v); convErr == nil && n > 0 {
				maxBytes = n
			}
		}
		if v, ok := args["max_bytes"].(float64); ok && v > 0 {
			maxBytes = int(v)
		}

		f, err := os.Open(abs)
		if err != nil {
			return "", fmt.Errorf("cannot open %q: %w", p, err)
		}
		defer f.Close()

		buf := make([]byte, maxBytes)
		n, readErr := f.Read(buf)
		if readErr != nil && !errors.Is(readErr, io.EOF) {
			return "", readErr
		}
		content := string(buf[:n])

		fi, _ := f.Stat()
		if fi != nil && fi.Size() > int64(maxBytes) {
			content += fmt.Sprintf(
				"\n\n[TRUNCATED: file is %d bytes total, returned first %d bytes]",
				fi.Size(), maxBytes)
		}
		return content, nil
	})
}

// -----------------------------------------------------------------------------
// Agentic pipeline: tool loop -> validation pass
// -----------------------------------------------------------------------------

// ToolActivity records a single tool invocation for display in the UI.
type ToolActivity struct {
	Name    string `json:"Name"`
	Args    string `json:"Args"`
	Result  string `json:"Result"`
	IsError bool   `json:"IsError"`
}

// PipelineResult holds the full outcome of one user turn.
type PipelineResult struct {
	Draft          string
	Final          string
	ValidationNote string
	Validated      bool
	ToolActivities []ToolActivity
}

const validatorSystemPrompt = `You are a factual accuracy validator.
You receive a draft AI response and the conversation that produced it.
Identify hallucinated or ungrounded factual claims.
Return ONLY a JSON object (no markdown fences) with fields:
  "verdict":  "pass" or "flag"
  "note":     brief explanation, empty string if pass
  "revised":  if flag, a corrected response; otherwise repeat the original`

// Pipeline runs the agentic tool loop then the validation pass.
type Pipeline struct {
	Generator    LLMClient
	Validator    LLMClient // nil = skip validation
	GenModel     string
	ValModel     string
	Dispatcher   *ToolDispatcher
	MaxToolIters int
}

func (p *Pipeline) Run(ctx context.Context, history []Message) (PipelineResult, error) {
	var activities []ToolActivity

	// Working copy of the message list -- tool results accumulate here
	working := make([]Message, len(history))
	copy(working, history)

	tools := p.Dispatcher.Defs()
	var finalContent string

	// Agentic loop: keep calling the LLM until it produces text (not a tool call)
	for i := 0; i < p.MaxToolIters; i++ {
		resp, err := p.Generator.Complete(ctx, LLMRequest{
			Model:    p.GenModel,
			Messages: working,
			Tools:    tools,
		})
		if err != nil {
			return PipelineResult{}, fmt.Errorf("generation (iter %d): %w", i, err)
		}

		if !resp.HasToolCalls() {
			// LLM produced a text response -- loop is done
			finalContent = resp.Content
			break
		}

		// LLM requested tool calls -- record its turn, then execute each call
		working = append(working, Message{
			Role:      "assistant",
			ToolCalls: resp.ToolCalls,
		})

		for _, tc := range resp.ToolCalls {
			log.Printf("[tool call] %s(%s)", tc.Name, string(tc.ArgsJSON))
			result := p.Dispatcher.Dispatch(tc)

			activities = append(activities, ToolActivity{
				Name: tc.Name, Args: string(tc.ArgsJSON),
				Result: result.Content, IsError: result.IsError,
			})

			// Feed result back so the LLM sees it on the next iteration
			working = append(working, Message{
				Role:       "tool",
				Content:    result.Content,
				ToolCallID: result.CallID,
				ToolName:   result.Name,
			})
		}
		// Loop continues: LLM will now see tool results and decide next action
	}

	if finalContent == "" {
		finalContent = "(agent reached the tool iteration limit without producing a final response)"
	}

	result := PipelineResult{
		Draft: finalContent, Final: finalContent,
		ToolActivities: activities,
	}

	// Validation pass
	if p.Validator == nil {
		return result, nil
	}

	historyJSON, _ := json.MarshalIndent(history, "", "  ")
	valMessages := []Message{
		{Role: "system", Content: validatorSystemPrompt},
		{Role: "user", Content: fmt.Sprintf(
			"Conversation:\n%s\n\nDraft response:\n%s",
			string(historyJSON), finalContent,
		)},
	}
	valResp, err := p.Validator.Complete(ctx, LLMRequest{
		Model: p.ValModel, Messages: valMessages, Temperature: 0.1,
	})
	if err != nil {
		log.Printf("validation pass error (using draft): %v", err)
		return result, nil
	}

	cleaned := strings.TrimSpace(valResp.Content)
	cleaned = strings.TrimPrefix(strings.TrimPrefix(cleaned, "```json"), "```")
	cleaned = strings.TrimSuffix(cleaned, "```")

	var valResult struct {
		Verdict string `json:"verdict"`
		Note    string `json:"note"`
		Revised string `json:"revised"`
	}
	if err := json.Unmarshal([]byte(cleaned), &valResult); err != nil {
		log.Printf("validator non-JSON (%v), using draft", err)
		return result, nil
	}
	result.ValidationNote = valResult.Note
	result.Validated = true
	if valResult.Verdict == "flag" && valResult.Revised != "" {
		result.Final = valResult.Revised
	}
	return result, nil
}

// -----------------------------------------------------------------------------
// Conversation store -- JSONL per session, git-committed
// -----------------------------------------------------------------------------

type Turn struct {
	ID             string         `json:"id"`
	Timestamp      time.Time      `json:"timestamp"`
	UserInput      string         `json:"user_input"`
	Draft          string         `json:"draft"`
	Final          string         `json:"final"`
	ValidationNote string         `json:"validation_note,omitempty"`
	Validated      bool           `json:"validated"`
	Model          string         `json:"model"`
	ToolActivities []ToolActivity `json:"tool_activities,omitempty"`
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
	if _, err := os.Stat(filepath.Join(s.dir, ".git")); os.IsNotExist(err) {
		if out, err := exec.Command("git", "init", s.dir).CombinedOutput(); err != nil {
			log.Printf("git init: %v -- %s", err, out)
		} else {
			log.Printf("git repo initialised at %s", s.dir)
		}
	}
}

func (s *ConversationStore) SessionFile(id string) string {
	return filepath.Join(s.dir, id+".jsonl")
}

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
	if out, err := exec.Command("git", "-C", s.dir, "add", rel).CombinedOutput(); err != nil {
		log.Printf("git add: %v -- %s", err, out)
		return
	}
	msg := fmt.Sprintf("turn %s session %s", turnID, sessionID)
	if out, err := exec.Command("git", "-C", s.dir, "commit", "-m", msg,
		"--author=Emily Agent <emily@agent.local>").CombinedOutput(); err != nil {
		log.Printf("git commit: %v -- %s", err, out)
	}
}

func (s *ConversationStore) LoadHistory(sessionID, systemPrompt string) ([]Message, []Turn, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	messages := []Message{{Role: "system", Content: systemPrompt}}
	var turns []Turn

	f, err := os.Open(s.SessionFile(sessionID))
	if os.IsNotExist(err) {
		return messages, turns, nil
	}
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
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

// -----------------------------------------------------------------------------
// Rate limiter
// -----------------------------------------------------------------------------

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
	b.tokens += now.Sub(b.lastFill).Seconds() * (rl.rpm / 60.0)
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

// -----------------------------------------------------------------------------
// Emily system prompt
// -----------------------------------------------------------------------------

const emilySystemPrompt = `You are Emily, an AI agent with a dual role:

EXTERNAL (Concierge): Warm, precise, and resourceful when communicating with users.
INTERNAL (Chief of Staff): Direct and structured when coordinating tasks internally.

Core traits:
- Honest about uncertainty. Never fabricate facts.
- Responses are focused and actionable.
- You maintain context across the conversation.
- You are aware you are part of a growing multi-agent system.

TOOLS YOU HAVE ACCESS TO:

  git_list_files  -- list files tracked in the conversation repository
  git_read_file   -- read any tracked file by path (.jsonl sessions, proposals, etc.)
  write_file      -- write a file to proposals/ in the conversation repository

HOW TO USE THEM:
- For past conversations: git_list_files → git_read_file on the relevant .jsonl
- For RSI improvement proposals: write_file to persist generated code or analysis
  for human review. Path examples: 'collector-improvement.go', 'analysis/reddit-quality.md'
- Use tools proactively and autonomously. Never ask permission to look something up.

RSI LOOP ROLE:
When asked to improve a component, you generate the improvement as a complete artifact
and write it to proposals/ using write_file. The artifact is then reviewed and applied.
This is how you close the recursive self-improvement loop.

Current capabilities: conversation memory, git repository read/write (proposals), RSI loop.
Agents: Bob (database, not yet deployed). Data collection: ArXiv, Reddit, Wikipedia.`

// -----------------------------------------------------------------------------
// HTTP server
// -----------------------------------------------------------------------------

type Server struct {
	cfg         Config
	pipeline    *Pipeline
	store       *ConversationStore
	rl          *RateLimiter
	tmpl        *template.Template
	collector   *CollectorPipeline   // nil if EMILY_COLLECT_DIR unset
	integration *IntegrationStore    // nil if signals/ dir not found
}

func NewServer(cfg Config) (*Server, error) {
	var generator LLMClient
	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		generator = NewAnthropicClient(key)
		if cfg.Model == "gpt-4o" {
			cfg.Model = "claude-sonnet-4-20250514"
		}
		log.Println("provider: Anthropic")
	} else {
		generator = NewOpenAIClient(cfg.APIKey)
		log.Println("provider: OpenAI")
	}

	var validator LLMClient
	valModel := cfg.ValidatorModel
	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		validator = NewAnthropicClient(key)
		if valModel == "gpt-4o-mini" {
			valModel = "claude-haiku-4-5-20251001"
		}
	} else if cfg.APIKey != "" {
		validator = NewOpenAIClient(cfg.APIKey)
	}

	store, err := NewConversationStore(cfg.ConversationDir, cfg.GitCommit)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	dispatcher := NewToolDispatcher()
	registerGitTools(dispatcher, cfg.ConversationDir)
	log.Printf("tools registered: %d", len(dispatcher.Defs()))
	for _, t := range dispatcher.Defs() {
		log.Printf("  * %s", t.Name)
	}

	pipeline := &Pipeline{
		Generator:    generator,
		Validator:    validator,
		GenModel:     cfg.Model,
		ValModel:     valModel,
		Dispatcher:   dispatcher,
		MaxToolIters: cfg.MaxToolIters,
	}

	tmpl, err := template.New("chat").Parse(chatHTML)
	if err != nil {
		return nil, fmt.Errorf("template: %w", err)
	}

	srv := &Server{cfg: cfg, pipeline: pipeline, store: store,
		rl: NewRateLimiter(cfg.RateLimitRPM), tmpl: tmpl}
	srv.collector = buildCollector(cfg)
	if srv.collector != nil {
		if err := srv.collector.Start(context.Background()); err != nil {
			log.Printf("collector start failed (continuing without): %v", err)
			srv.collector = nil
		}
	}
	srv.integration = buildIntegrationStore()
	if srv.integration != nil {
		registerIntegrationTools(dispatcher, srv.integration)
		log.Printf("integration store: signals dir=%s", srv.integration.signalsDir)
	}
	return srv, nil
}

func buildCollector(cfg Config) *CollectorPipeline {
	collectDir := envOr("EMILY_COLLECT_DIR", "")
	if collectDir == "" {
		return nil
	}
	store, err := NewDocStore(collectDir)
	if err != nil {
		log.Printf("doc store init failed: %v", err)
		return nil
	}
	ua := envOr("EMILY_COLLECT_UA", "emily-agent/0.1 data-collection-bot")
	subredditEnv := envOr("EMILY_COLLECT_SUBREDDITS", "MachineLearning,artificial,LocalLLaMA,programming,learnprogramming")
	var subreddits []string
	for _, s := range strings.Split(subredditEnv, ",") {
		if s = strings.TrimSpace(s); s != "" {
			subreddits = append(subreddits, s)
		}
	}
	sources := []Source{
		&ArXivSource{
			Categories:  []string{"cs.AI", "cs.LG", "cs.CL", "cs.CV", "stat.ML"},
			MaxResults:  100,
			UserAgent:   ua,
			MinAbstrLen: 200,
		},
		&RedditSource{
			Subreddits:  subreddits,
			UserAgent:   ua,
			MinScore:    100,
			MinComments: 10,
			MinBodyLen:  500,
		},
		&WikipediaSource{
			UserAgent:  ua,
			MinBodyLen: 1000,
		},
	}
	return NewCollectorPipeline(CollectorConfig{
		Store:        store,
		Sources:      sources,
		Scorer:       &QualityScorer{},
		Workers:      4,
		PollInterval: 5 * time.Minute,
		MinTier:      envOr("EMILY_COLLECT_MIN_TIER", "bronze"),
	})
}

func (s *Server) ip(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.Split(xff, ",")[0]
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	if sessionID == "" {
		sessionID = fmt.Sprintf("session-%d", time.Now().UnixMilli())
		http.Redirect(w, r, "/?session="+sessionID, http.StatusFound)
		return
	}
	_, turns, err := s.store.LoadHistory(sessionID, emilySystemPrompt)
	if err != nil {
		http.Error(w, "history load failed", 500)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	s.tmpl.Execute(w, map[string]any{"SessionID": sessionID, "Turns": turns})
}

func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	if !s.rl.Allow(s.ip(r)) {
		w.Header().Set("Content-Type", "application/json")
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

	history, _, err := s.store.LoadHistory(body.SessionID, emilySystemPrompt)
	if err != nil {
		http.Error(w, `{"error":"history load failed"}`, 500)
		return
	}
	history = append(history, Message{Role: "user", Content: body.Message})

	ctx, cancel := context.WithTimeout(r.Context(), 120*time.Second)
	defer cancel()

	result, err := s.pipeline.Run(ctx, history)
	if err != nil {
		log.Printf("pipeline: %v", err)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), 500)
		return
	}

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
		ToolActivities: result.ToolActivities,
	}
	if err := s.store.AppendTurn(body.SessionID, turn); err != nil {
		log.Printf("store: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"reply":           result.Final,
		"validated":       result.Validated,
		"validation_note": result.ValidationNote,
		"tool_activities": result.ToolActivities,
		"turn_id":         turnID,
	})
}

func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	if sessionID == "" {
		http.Error(w, `{"error":"session required"}`, 400)
		return
	}
	_, turns, err := s.store.LoadHistory(sessionID, emilySystemPrompt)
	if err != nil {
		http.Error(w, `{"error":"load failed"}`, 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(turns)
}

// -----------------------------------------------------------------------------
// Embedded chat UI
// -----------------------------------------------------------------------------

const chatHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Emily -- Agent Interface</title>
<style>
  @import url('https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@400;500&family=IBM+Plex+Sans:wght@300;400;500&display=swap');
  :root {
    --bg:#0e0e10; --surface:#18181b; --border:#2a2a30;
    --accent:#7ee8a2; --accent-dim:#3d7a54;
    --text:#e4e4e7; --text-dim:#71717a;
    --user-bg:#1e2a22; --emily-bg:#1a1a20;
    --tool-bg:#1a1e28; --tool-border:#2a3050;
    --danger:#f87171; --warn:#fbbf24; --info:#60a5fa;
  }
  *,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
  body{background:var(--bg);color:var(--text);font-family:'IBM Plex Sans',sans-serif;
    font-weight:300;height:100dvh;display:flex;flex-direction:column;overflow:hidden}
  header{padding:12px 20px;border-bottom:1px solid var(--border);display:flex;
    align-items:center;gap:12px;background:var(--surface);flex-shrink:0}
  .status-dot{width:8px;height:8px;border-radius:50%;background:var(--accent);
    box-shadow:0 0 8px var(--accent);animation:pulse 2s infinite}
  @keyframes pulse{0%,100%{opacity:1}50%{opacity:.4}}
  header h1{font-family:'IBM Plex Mono',monospace;font-size:14px;font-weight:500;
    letter-spacing:.08em;color:var(--accent)}
  .session-tag{margin-left:auto;font-family:'IBM Plex Mono',monospace;font-size:11px;color:var(--text-dim)}
  #messages{flex:1;overflow-y:auto;padding:24px 0;display:flex;flex-direction:column;
    gap:2px;scroll-behavior:smooth}
  #messages::-webkit-scrollbar{width:4px}
  #messages::-webkit-scrollbar-thumb{background:var(--border);border-radius:2px}
  .turn{display:flex;flex-direction:column;padding:0 20px}
  .bubble{max-width:min(680px,88%);padding:12px 16px;border-radius:4px;line-height:1.6;
    font-size:14px;white-space:pre-wrap;word-break:break-word;border:1px solid transparent}
  .turn.user{align-items:flex-end}
  .turn.user .bubble{background:var(--user-bg);border-color:var(--accent-dim)}
  .turn.user .label{color:var(--accent)}
  .turn.emily{align-items:flex-start}
  .turn.emily .bubble{background:var(--emily-bg);border-color:var(--border)}
  .turn.emily .label{color:var(--text-dim)}
  .label{font-family:'IBM Plex Mono',monospace;font-size:10px;letter-spacing:.1em;
    text-transform:uppercase;margin-bottom:4px;padding:0 4px}
  .tool-log{max-width:min(680px,88%);margin-top:6px;border:1px solid var(--tool-border);
    border-radius:4px;background:var(--tool-bg);font-family:'IBM Plex Mono',monospace;
    font-size:11px;overflow:hidden}
  .tool-log-header{padding:5px 10px;background:#1f2435;color:var(--info);letter-spacing:.05em;
    cursor:pointer;user-select:none;display:flex;align-items:center;gap:6px}
  .tool-log-header:hover{background:#252b42}
  .tool-log-body{display:none}
  .tool-log-body.open{display:block}
  .tool-entry{padding:8px 10px;border-top:1px solid var(--tool-border)}
  .tool-entry-name{color:var(--info);margin-bottom:3px}
  .tool-entry-args{color:var(--text-dim);margin-bottom:3px;font-size:10px}
  .tool-entry-result{color:var(--text);white-space:pre-wrap;max-height:200px;overflow-y:auto}
  .tool-entry.error .tool-entry-name{color:var(--danger)}
  .validation-badge{font-family:'IBM Plex Mono',monospace;font-size:10px;margin-top:6px;
    padding:3px 8px;border-radius:2px;display:inline-block}
  .badge-pass{background:#1a2e1e;color:var(--accent);border:1px solid var(--accent-dim)}
  .badge-flag{background:#2e1e1a;color:var(--warn);border:1px solid #7a5a3d}
  .typing-indicator{display:none;align-items:flex-start;padding:0 20px}
  .typing-indicator.active{display:flex}
  .typing-indicator .label{color:var(--text-dim)}
  .dot-row{display:flex;gap:4px;padding:12px 16px}
  .dot{width:6px;height:6px;background:var(--text-dim);border-radius:50%;animation:bounce 1.2s infinite}
  .dot:nth-child(2){animation-delay:.2s}.dot:nth-child(3){animation-delay:.4s}
  @keyframes bounce{0%,60%,100%{transform:translateY(0)}30%{transform:translateY(-6px)}}
  footer{padding:16px 20px;border-top:1px solid var(--border);background:var(--surface);flex-shrink:0}
  #input-row{display:flex;gap:10px;align-items:flex-end;max-width:740px;margin:0 auto}
  #user-input{flex:1;background:var(--bg);border:1px solid var(--border);border-radius:4px;
    color:var(--text);font-family:'IBM Plex Sans',sans-serif;font-size:14px;font-weight:300;
    padding:10px 14px;resize:none;line-height:1.5;max-height:160px;outline:none;transition:border-color .15s}
  #user-input:focus{border-color:var(--accent-dim)}
  #user-input::placeholder{color:var(--text-dim)}
  #send-btn{background:var(--accent);color:#0e0e10;border:none;border-radius:4px;padding:10px 18px;
    font-family:'IBM Plex Mono',monospace;font-size:12px;font-weight:500;letter-spacing:.05em;
    cursor:pointer;transition:opacity .15s;white-space:nowrap}
  #send-btn:hover{opacity:.85}#send-btn:disabled{opacity:.4;cursor:not-allowed}
  .error-msg{color:var(--danger);font-family:'IBM Plex Mono',monospace;font-size:12px;text-align:center;padding:8px}
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
    {{if .ToolActivities}}
    <div class="tool-log">
      <div class="tool-log-header" onclick="this.nextElementSibling.classList.toggle('open')">
        &#9881; {{len .ToolActivities}} tool call(s) &mdash; click to expand
      </div>
      <div class="tool-log-body">
        {{range .ToolActivities}}
        <div class="tool-entry{{if .IsError}} error{{end}}">
          <div class="tool-entry-name">&#9654; {{.Name}}</div>
          <div class="tool-entry-args">args: {{.Args}}</div>
          <div class="tool-entry-result">{{.Result}}</div>
        </div>
        {{end}}
      </div>
    </div>
    {{end}}
    <div class="bubble">{{.Final}}</div>
    {{if .Validated}}
      {{if .ValidationNote}}
        <span class="validation-badge badge-flag">&#9888; {{.ValidationNote}}</span>
      {{else}}
        <span class="validation-badge badge-pass">&#10003; verified</span>
      {{end}}
    {{end}}
  </div>
{{end}}
</div>

<div class="typing-indicator" id="typing">
  <div>
    <div class="label">emily</div>
    <div class="dot-row"><div class="dot"></div><div class="dot"></div><div class="dot"></div></div>
  </div>
</div>

<footer>
  <div id="input-row">
    <textarea id="user-input" placeholder="Message Emily&#8230;" rows="1"></textarea>
    <button id="send-btn">SEND</button>
  </div>
</footer>

<script>
const SESSION_ID = {{.SessionID | printf "%q"}};
const messages = document.getElementById('messages');
const input    = document.getElementById('user-input');
const sendBtn  = document.getElementById('send-btn');
const typing   = document.getElementById('typing');

input.addEventListener('input', () => {
  input.style.height = 'auto';
  input.style.height = Math.min(input.scrollHeight, 160) + 'px';
});
input.addEventListener('keydown', e => {
  if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); sendMessage(); }
});
sendBtn.addEventListener('click', sendMessage);

function scrollToBottom() { messages.scrollTop = messages.scrollHeight; }

function esc(s) {
  return String(s).replace(/&/g,'&amp;').replace(/</g,'&lt;')
    .replace(/>/g,'&gt;').replace(/"/g,'&quot;');
}

function buildToolLog(activities) {
  if (!activities || !activities.length) return null;
  const wrap = document.createElement('div');
  wrap.className = 'tool-log';
  const hdr = document.createElement('div');
  hdr.className = 'tool-log-header';
  hdr.innerHTML = '&#9881; ' + activities.length + ' tool call(s) &mdash; click to expand';
  const body = document.createElement('div');
  body.className = 'tool-log-body';
  hdr.onclick = () => body.classList.toggle('open');
  activities.forEach(a => {
    const entry = document.createElement('div');
    entry.className = 'tool-entry' + (a.IsError ? ' error' : '');
    entry.innerHTML =
      '<div class="tool-entry-name">&#9654; ' + esc(a.Name) + '</div>' +
      '<div class="tool-entry-args">args: ' + esc(a.Args) + '</div>' +
      '<div class="tool-entry-result">' + esc(a.Result) + '</div>';
    body.appendChild(entry);
  });
  wrap.appendChild(hdr);
  wrap.appendChild(body);
  return wrap;
}

function addTurn(userText, reply, activities, validationNote, validated) {
  const userTurn = document.createElement('div');
  userTurn.className = 'turn user';
  userTurn.innerHTML = '<div class="label">you</div><div class="bubble"></div>';
  userTurn.querySelector('.bubble').textContent = userText;
  messages.appendChild(userTurn);

  const emilyTurn = document.createElement('div');
  emilyTurn.className = 'turn emily';
  const lbl = document.createElement('div');
  lbl.className = 'label';
  lbl.textContent = 'emily';
  emilyTurn.appendChild(lbl);

  const log = buildToolLog(activities);
  if (log) emilyTurn.appendChild(log);

  const bubble = document.createElement('div');
  bubble.className = 'bubble';
  bubble.textContent = reply;
  emilyTurn.appendChild(bubble);

  if (validated) {
    const badge = document.createElement('span');
    badge.className = 'validation-badge ' + (validationNote ? 'badge-flag' : 'badge-pass');
    badge.textContent = validationNote ? '\u26A0 ' + validationNote : '\u2713 verified';
    emilyTurn.appendChild(badge);
  }
  messages.appendChild(emilyTurn);
  scrollToBottom();
}

async function sendMessage() {
  const text = input.value.trim();
  if (!text) return;
  sendBtn.disabled = true;
  input.value = '';
  input.style.height = 'auto';
  typing.classList.add('active');
  scrollToBottom();

  try {
    const res = await fetch('/chat', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({session_id: SESSION_ID, message: text})
    });
    typing.classList.remove('active');
    if (!res.ok) {
      const err = await res.json().catch(() => ({error: 'request failed'}));
      addError(err.error || 'request failed');
      return;
    }
    const data = await res.json();
    addTurn(text, data.reply, data.tool_activities, data.validation_note, data.validated);
  } catch(e) {
    typing.classList.remove('active');
    addError('Network error: ' + e.message);
  } finally {
    sendBtn.disabled = false;
    input.focus();
  }
}

function addError(msg) {
  const div = document.createElement('div');
  div.className = 'error-msg';
  div.textContent = '\u26A0 ' + msg;
  messages.appendChild(div);
  scrollToBottom();
}

scrollToBottom();
input.focus();
</script>
</body>
</html>`

// -----------------------------------------------------------------------------
// Main
// -----------------------------------------------------------------------------

func main() {
	cfg := loadConfig()

	// --cron: run one autonomous cycle and exit (designed for cron scheduling)
	// --daemon: run autonomous cycles continuously at the configured interval
	for _, arg := range os.Args[1:] {
		if arg == "--cron" || arg == "--daemon" {
			srv, err := NewServer(cfg)
			if err != nil {
				log.Fatalf("server init: %v", err)
			}
			cronCfg := defaultCronConfig()
			cycle := NewAutonomousCycle(cronCfg, srv.pipeline)
			if arg == "--daemon" {
				cycle.RunDaemon()
			} else {
				if err := cycle.RunOnce(); err != nil {
					log.Fatalf("cycle: %v", err)
				}
			}
			return
		}
	}

	srv, err := NewServer(cfg)
	if err != nil {
		log.Fatalf("server init: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", srv.handleIndex)
	mux.HandleFunc("/chat", srv.handleChat)
	mux.HandleFunc("/history", srv.handleHistory)
	mux.HandleFunc("/rsi", srv.handleRSI)
	mux.HandleFunc("/roadmap", srv.handleRoadmap)
	mux.HandleFunc("/status", srv.handleStatus)
	mux.HandleFunc("/collect", srv.handleCollect)
	mux.HandleFunc("/integration/observations", srv.handleIntegrationObservations)
	mux.HandleFunc("/integration/task", srv.handleIntegrationTask)
	mux.HandleFunc("/integration/triage", srv.handleIntegrationTriage)
	mux.HandleFunc("/agent/result", srv.handleAgentResult)

	addr := ":" + cfg.Port
	log.Printf("Emily agent  ->  http://localhost%s", addr)
	log.Printf("Conversations: %s | git: %v | rpm: %d | max-tool-iters: %d",
		cfg.ConversationDir, cfg.GitCommit, cfg.RateLimitRPM, cfg.MaxToolIters)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server: %v", err)
	}
}
