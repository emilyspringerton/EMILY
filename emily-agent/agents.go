// emily-agent/agents.go
// Specialist agent communication — Emily-side dispatcher.
//
// Implements the message envelope protocol from emily-agent-protocol.md.
// Agents are HTTP services. Emily sends tasks and receives results.

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// --- Message envelope ---

type MsgType string

const (
	MsgTask       MsgType = "task"
	MsgResult     MsgType = "result"
	MsgEscalation MsgType = "escalation"
	MsgHeartbeat  MsgType = "heartbeat"
)

type Priority string

const (
	PriorityCritical Priority = "critical"
	PriorityHigh     Priority = "high"
	PriorityNormal   Priority = "normal"
	PriorityLow      Priority = "low"
)

type Envelope struct {
	V         int             `json:"v"`
	ID        string          `json:"id"`
	From      string          `json:"from"`
	To        string          `json:"to"`
	Type      MsgType         `json:"type"`
	Priority  Priority        `json:"priority"`
	SentAt    time.Time       `json:"sent_at"`
	ExpiresAt time.Time       `json:"expires_at"`
	Body      json.RawMessage `json:"body"`
	TraceID   string          `json:"trace_id"`
	ReplyTo   string          `json:"reply_to,omitempty"`
}

// --- Body types ---

type TaskBody struct {
	TaskID           string               `json:"task_id"`
	Description      string               `json:"description"`
	Inputs           map[string]any       `json:"inputs"`
	SuccessCriteria  []AcceptanceCriterion `json:"success_criteria"`
	TimeoutSeconds   int                  `json:"timeout_seconds"`
	RequiresConfirm  bool                 `json:"requires_confirm"`
}

type ResultBody struct {
	TaskID          string            `json:"task_id"`
	Status          string            `json:"status"` // success|partial|failed|awaiting_confirm
	Outputs         map[string]any    `json:"outputs"`
	CriteriaResults []CriteriaResult  `json:"criteria_results"`
	AllPass         bool              `json:"all_pass"`
	Error           string            `json:"error"`
	DurationMS      int64             `json:"duration_ms"`
	Artifacts       []string          `json:"artifacts,omitempty"`
}

type EscalationBody struct {
	TaskID      string            `json:"task_id"`
	Category    string            `json:"category"`
	Title       string            `json:"title"`
	Summary     string            `json:"summary"`
	Context     map[string]string `json:"context"`
	Options     []EscalationOption `json:"options"`
	Recommended string            `json:"recommended"`
	Deadline    *time.Time        `json:"deadline,omitempty"`
}

type EscalationOption struct {
	ID             string `json:"id"`
	Label          string `json:"label"`
	RequiresConfirm bool  `json:"requires_confirm"`
}

type HeartbeatBody struct {
	Status         string `json:"status"` // healthy|degraded|down
	UptimeSeconds  int64  `json:"uptime_seconds"`
	TasksPending   int    `json:"tasks_pending"`
	TasksComplete  int    `json:"tasks_complete"`
	LastError      string `json:"last_error"`
}

// --- Agent registry ---

type AgentConfig struct {
	Name    string `json:"name"`
	BaseURL string `json:"base_url"`
	Timeout int    `json:"timeout_seconds"` // default task timeout
}

type AgentRegistry struct {
	agents map[string]AgentConfig
	mu     sync.RWMutex
}

func NewAgentRegistry() *AgentRegistry {
	r := &AgentRegistry{agents: make(map[string]AgentConfig)}
	r.loadFromEnv()
	return r
}

// loadFromEnv reads agent configs from environment variables.
// Format: EMILY_AGENT_BOB=http://localhost:9001
func (r *AgentRegistry) loadFromEnv() {
	prefix := "EMILY_AGENT_"
	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, prefix) {
			continue
		}
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		name := strings.ToLower(strings.TrimPrefix(parts[0], prefix))
		r.Register(AgentConfig{Name: name, BaseURL: parts[1], Timeout: 120})
	}
}

func (r *AgentRegistry) Register(cfg AgentConfig) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.agents[cfg.Name] = cfg
}

func (r *AgentRegistry) Get(name string) (AgentConfig, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cfg, ok := r.agents[name]
	return cfg, ok
}

func (r *AgentRegistry) All() []AgentConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]AgentConfig, 0, len(r.agents))
	for _, cfg := range r.agents {
		out = append(out, cfg)
	}
	return out
}

// --- Emily's agent dispatcher ---

type AgentDispatcher struct {
	registry *AgentRegistry
	auditLog *AuditLog
	traceID  string
}

func NewAgentDispatcher(registry *AgentRegistry, stateDir string) *AgentDispatcher {
	return &AgentDispatcher{
		registry: registry,
		auditLog: NewAuditLog(filepath.Join(stateDir, "agent-audit.jsonl")),
	}
}

func (d *AgentDispatcher) WithTrace(traceID string) *AgentDispatcher {
	return &AgentDispatcher{
		registry: d.registry,
		auditLog: d.auditLog,
		traceID:  traceID,
	}
}

// SendTask sends a task to a named agent and waits for the result.
func (d *AgentDispatcher) SendTask(ctx context.Context, agentName string, task TaskBody) (*ResultBody, error) {
	cfg, ok := d.registry.Get(agentName)
	if !ok {
		return nil, fmt.Errorf("agent %q not registered", agentName)
	}

	taskTimeout := task.TimeoutSeconds
	if taskTimeout <= 0 {
		taskTimeout = cfg.Timeout
	}
	if taskTimeout <= 0 {
		taskTimeout = 120
	}

	bodyJSON, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	env := Envelope{
		V:         1,
		ID:        msgID(),
		From:      "emily",
		To:        agentName,
		Type:      MsgTask,
		Priority:  PriorityNormal,
		SentAt:    time.Now(),
		ExpiresAt: time.Now().Add(time.Duration(taskTimeout) * time.Second),
		Body:      bodyJSON,
		TraceID:   d.traceID,
	}

	d.auditLog.Log("sent", env.TraceID, env.ID, "emily", agentName, string(MsgTask), task.TaskID, "")

	result, err := d.postAndWait(ctx, cfg, env, taskTimeout)
	if err != nil {
		d.auditLog.Log("error", env.TraceID, env.ID, "emily", agentName, string(MsgTask), task.TaskID, err.Error())
		return nil, err
	}

	d.auditLog.Log("received", env.TraceID, env.ID, agentName, "emily", string(MsgResult), task.TaskID, result.Status)
	return result, nil
}

func (d *AgentDispatcher) postAndWait(ctx context.Context, cfg AgentConfig, env Envelope, timeoutSec int) (*ResultBody, error) {
	payload, err := json.Marshal(env)
	if err != nil {
		return nil, err
	}

	deadline := time.Now().Add(time.Duration(timeoutSec) * time.Second)
	ctx, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", cfg.BaseURL+"/task", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send to %s: %w", cfg.Name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusAccepted {
		// Async: poll GET /task/{id}
		return d.pollResult(ctx, cfg, env.Body, timeoutSec)
	}

	// Synchronous response
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s returned %d: %s", cfg.Name, resp.StatusCode, raw)
	}

	var resultEnv Envelope
	if err := json.Unmarshal(raw, &resultEnv); err != nil {
		return nil, fmt.Errorf("parse result from %s: %w", cfg.Name, err)
	}
	var rb ResultBody
	if err := json.Unmarshal(resultEnv.Body, &rb); err != nil {
		return nil, fmt.Errorf("parse result body from %s: %w", cfg.Name, err)
	}
	return &rb, nil
}

func (d *AgentDispatcher) pollResult(ctx context.Context, cfg AgentConfig, taskBodyRaw json.RawMessage, timeoutSec int) (*ResultBody, error) {
	var tb TaskBody
	_ = json.Unmarshal(taskBodyRaw, &tb)
	taskID := tb.TaskID

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout waiting for task %s from %s", taskID, cfg.Name)
		case <-ticker.C:
			url := fmt.Sprintf("%s/task/%s", cfg.BaseURL, taskID)
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				continue
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				continue
			}
			raw, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			var resultEnv Envelope
			if err := json.Unmarshal(raw, &resultEnv); err != nil {
				continue
			}
			if resultEnv.Type != MsgResult {
				continue // still running
			}
			var rb ResultBody
			if err := json.Unmarshal(resultEnv.Body, &rb); err != nil {
				return nil, err
			}
			return &rb, nil
		}
	}
}

// CheckHealth returns the health of a named agent.
func (d *AgentDispatcher) CheckHealth(ctx context.Context, agentName string) (*HeartbeatBody, error) {
	cfg, ok := d.registry.Get(agentName)
	if !ok {
		return nil, fmt.Errorf("agent %q not registered", agentName)
	}
	req, err := http.NewRequestWithContext(ctx, "GET", cfg.BaseURL+"/health", nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &HeartbeatBody{Status: "down", LastError: err.Error()}, nil
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var hb HeartbeatBody
	if err := json.Unmarshal(raw, &hb); err != nil {
		return &HeartbeatBody{Status: "degraded", LastError: "parse error"}, nil
	}
	return &hb, nil
}

// --- Audit log ---

type AuditLog struct {
	path string
	mu   sync.Mutex
}

func NewAuditLog(path string) *AuditLog {
	return &AuditLog{path: path}
}

type auditEntry struct {
	Event   string `json:"event"`
	TraceID string `json:"trace_id"`
	MsgID   string `json:"msg_id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Type    string `json:"type"`
	TaskID  string `json:"task_id"`
	Status  string `json:"status,omitempty"`
	At      string `json:"at"`
}

func (a *AuditLog) Log(event, traceID, msgID, from, to, msgType, taskID, status string) {
	entry := auditEntry{
		Event: event, TraceID: traceID, MsgID: msgID,
		From: from, To: to, Type: msgType, TaskID: taskID,
		Status: status, At: time.Now().UTC().Format(time.RFC3339),
	}
	data, _ := json.Marshal(entry)
	a.mu.Lock()
	defer a.mu.Unlock()
	f, err := os.OpenFile(a.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Printf("audit log error: %v", err)
		return
	}
	fmt.Fprintf(f, "%s\n", data)
	f.Close()
}

// --- HTTP handler: Emily receives results from agents ---

func (s *Server) handleAgentResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	var env Envelope
	if err := json.NewDecoder(r.Body).Decode(&env); err != nil {
		http.Error(w, `{"error":"invalid envelope"}`, 400)
		return
	}
	if env.Type != MsgResult && env.Type != MsgEscalation {
		http.Error(w, `{"error":"expected type result or escalation"}`, 400)
		return
	}
	// Store for retrieval — in production this would notify the waiting cycle
	stateDir := envOr("EMILY_STATE_DIR", "./emily-state")
	if err := os.MkdirAll(stateDir, 0o755); err == nil {
		path := filepath.Join(stateDir, "agent-results", string(env.ID)+".json")
		_ = os.MkdirAll(filepath.Dir(path), 0o755)
		data, _ := json.MarshalIndent(env, "", "  ")
		_ = os.WriteFile(path, data, 0o644)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "accepted", "id": env.ID})
}

// --- Helpers ---

func msgID() string {
	return fmt.Sprintf("msg_%d", time.Now().UnixNano())
}

// ResultBody.Status helper
func (rb *ResultBody) IsComplete() bool {
	return rb.Status == "success" || rb.Status == "partial" || rb.Status == "failed"
}
