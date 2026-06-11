// emily-agent/integration.go
// Emily Prime ↔ FatBaby-Emily integration layer.
//
// Git is the message bus. Emily Prime reads observations from signals/observations/
// and writes directed tasks to signals/tasks/. Both sides poll their directory.
// Every message is a committed JSON file — fully auditable, diffable, reversible.
//
// Configured via environment:
//   EMILY_INTEGRATION_DIR  — path to the signals/ directory (default: ./signals)
//   EMILY_GIT_ROOT         — git repo root containing signals/ (default: ..)
//
// Routes wired in main.go:
//   GET  /integration/observations  — recent observations from FatBaby
//   POST /integration/task          — issue a directed task to FatBaby
//   GET  /integration/triage        — triage latest observations against signal priorities

package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------
// Observation — FatBaby-Emily → Emily Prime
// Superset of the fatbaby_write_observation output and the spec schema.
// ---------------------------------------------------------------------------

type Observation struct {
	Timestamp            string         `json:"timestamp"`
	Source               string         `json:"source"`              // "fatbaby-emily" | "entity-graph"
	ObservationType      string         `json:"observation_type"`    // anomaly|improvement|escalation|status
	Severity             string         `json:"severity"`            // critical|high|normal|low|info|warn|error
	Summary              string         `json:"summary"`
	Detail               string         `json:"detail,omitempty"`
	Findings             string         `json:"findings,omitempty"`
	SuggestedFix         string         `json:"suggested_fix,omitempty"`
	AffectedTickers      []string       `json:"affected_tickers,omitempty"`
	SignalIDs            []string       `json:"signal_ids,omitempty"`
	SignalsByType        map[string]int `json:"signals_by_type,omitempty"`
	RecommendedAction    string         `json:"recommended_action,omitempty"`
	RequiresCEOVisibility bool          `json:"requires_ceo_visibility"`
	RequestForClaude     string         `json:"request_for_claude,omitempty"`
	// Internal: filename this was loaded from (not serialized)
	filename string `json:"-"`
}

// ---------------------------------------------------------------------------
// DirectedTask — Emily Prime → FatBaby-Emily
// ---------------------------------------------------------------------------

type DirectedTask struct {
	Timestamp          string   `json:"timestamp"`
	From               string   `json:"from"`               // "emily-prime"
	To                 string   `json:"to"`                 // "fatbaby-emily"
	TaskID             string   `json:"task_id"`
	TaskType           string   `json:"task_type"`          // expand_coverage|improve_signal|fix_anomaly|config_change
	Priority           string   `json:"priority"`           // critical|high|normal|low
	Description        string   `json:"description"`
	AcceptanceCriteria []string `json:"acceptance_criteria"`
	Context            string   `json:"context"`            // strategic rationale
	Deadline           string   `json:"deadline,omitempty"`
}

// ---------------------------------------------------------------------------
// SignalPriorities — loaded from context/signal-priorities.json
// ---------------------------------------------------------------------------

type SignalPriorities struct {
	Version                 int                `json:"version"`
	UpdatedAt               string             `json:"updated_at"`
	SignalWeights           map[string]float64 `json:"signal_weights"`
	CEOVisibilityThreshold  float64            `json:"ceo_visibility_threshold"`
	EscalationCriteria      struct {
		DirectorDepartureAlways  bool    `json:"director_departure_always"`
		AuditorChangeAlways      bool    `json:"auditor_change_always"`
		EPSMissThresholdPct      float64 `json:"eps_miss_threshold_pct"`
		GovernanceHealthIndexMin float64 `json:"governance_health_index_min"`
		ActivistRiskAlways       bool    `json:"activist_risk_always"`
	} `json:"escalation_criteria"`
}

// ---------------------------------------------------------------------------
// TriageResult — output of Emily Prime's strategic relevance assessment
// ---------------------------------------------------------------------------

type TriageResult struct {
	ObservationTimestamp  string   `json:"observation_timestamp"`
	StrategicRelevance    float64  `json:"strategic_relevance"`  // 0.0–1.0
	RequiresCEOVisibility bool     `json:"requires_ceo_visibility"`
	TriggeredSignals      []string `json:"triggered_signals,omitempty"`
	Rationale             string   `json:"rationale"`
	SuggestedTaskType     string   `json:"suggested_task_type,omitempty"`
}

// ---------------------------------------------------------------------------
// IntegrationStore — Git-backed message bus for the integration layer
// ---------------------------------------------------------------------------

type IntegrationStore struct {
	signalsDir  string // path to signals/
	gitRoot     string // git repo root (EMILY)
	fatBabyRoot string // PRRJECT_FATBABY root (for backlog curation)
	mu          sync.Mutex
}

func NewIntegrationStore(signalsDir, gitRoot string) (*IntegrationStore, error) {
	for _, sub := range []string{"observations", "tasks"} {
		if err := os.MkdirAll(filepath.Join(signalsDir, sub), 0o755); err != nil {
			return nil, fmt.Errorf("integration_store mkdir %s: %w", sub, err)
		}
	}
	fatBabyRoot := os.Getenv("FATBABY_ROOT")
	if fatBabyRoot == "" {
		fatBabyRoot = "/home/fatbaby/PRRJECT_FATBABY"
	}
	return &IntegrationStore{signalsDir: signalsDir, gitRoot: gitRoot, fatBabyRoot: fatBabyRoot}, nil
}

func (s *IntegrationStore) WriteObservation(obs Observation) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if obs.Timestamp == "" {
		obs.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}
	if obs.Source == "" {
		obs.Source = "fatbaby-emily"
	}
	fname := fmt.Sprintf("%s-%s.json",
		strings.ReplaceAll(obs.Timestamp, ":", ""),
		slugify(obs.Summary, 32))
	path := filepath.Join(s.signalsDir, "observations", fname)
	data, err := json.MarshalIndent(obs, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return err
	}
	s.gitCommit(path, "observation: "+obs.Summary)
	return nil
}

func (s *IntegrationStore) WriteTask(task DirectedTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if task.Timestamp == "" {
		task.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}
	if task.From == "" {
		task.From = "emily-prime"
	}
	if task.To == "" {
		task.To = "fatbaby-emily"
	}
	if task.TaskID == "" {
		task.TaskID = fmt.Sprintf("task-%d", time.Now().UnixNano())
	}

	// Deduplication: skip if an identical task (same type + description) was written
	// in the last 4 hours. The triage cycle runs every 5 minutes and re-scores the
	// same observations on each pass; without this guard, a single high-relevance
	// observation generates a new task file every cycle, flooding the dispatch queue
	// and causing duplicate Claude Code sessions (~100K tokens each).
	tasksDir := filepath.Join(s.signalsDir, "tasks")
	if s.recentDuplicateExists(tasksDir, task.TaskType, task.Description) {
		log.Printf("integration: skipping duplicate task type=%s (identical task pending within 4h)", task.TaskType)
		return nil
	}

	fname := fmt.Sprintf("%s-%s.json",
		strings.ReplaceAll(task.Timestamp, ":", ""),
		task.TaskID)
	path := filepath.Join(tasksDir, fname)
	data, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return err
	}
	s.gitCommit(path, "task: "+task.Description)
	return nil
}

// recentDuplicateExists scans the tasks directory for a task with matching type
// and description written within the last 4 hours. Caller must hold s.mu.
func (s *IntegrationStore) recentDuplicateExists(dir, taskType, description string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	// Cutoff in filename format: RFC3339 with colons stripped, e.g. "2026-06-10T102648Z".
	cutoff := strings.ReplaceAll(
		time.Now().UTC().Add(-4*time.Hour).Format(time.RFC3339), ":", "")
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		if e.Name() < cutoff {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		var t DirectedTask
		if err := json.Unmarshal(data, &t); err != nil {
			continue
		}
		if t.TaskType == taskType && t.Description == description {
			return true
		}
	}
	return false
}

// ReadObservations returns the N most recent observations (sorted descending).
func (s *IntegrationStore) ReadObservations(limit int) ([]Observation, error) {
	if limit <= 0 {
		limit = 20
	}
	dir := filepath.Join(s.signalsDir, "observations")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read observations dir: %w", err)
	}

	// Sort descending (newest first) — filenames are timestamp-prefixed.
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() > entries[j].Name()
	})

	var out []Observation
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		var obs Observation
		if err := json.Unmarshal(data, &obs); err != nil {
			continue
		}
		obs.filename = e.Name()
		out = append(out, obs)
		if len(out) >= limit {
			break
		}
	}
	return out, nil
}

// ReadPendingTasks returns unprocessed tasks from signals/tasks/.
func (s *IntegrationStore) ReadPendingTasks(limit int) ([]DirectedTask, error) {
	if limit <= 0 {
		limit = 20
	}
	dir := filepath.Join(s.signalsDir, "tasks")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read tasks dir: %w", err)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() > entries[j].Name()
	})
	var out []DirectedTask
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		var task DirectedTask
		if err := json.Unmarshal(data, &task); err != nil {
			continue
		}
		out = append(out, task)
		if len(out) >= limit {
			break
		}
	}
	return out, nil
}

// LoadSignalPriorities reads context/signal-priorities.json relative to gitRoot.
func (s *IntegrationStore) LoadSignalPriorities() (SignalPriorities, error) {
	path := filepath.Join(s.gitRoot, "context", "signal-priorities.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return SignalPriorities{}, fmt.Errorf("load signal_priorities: %w", err)
	}
	var p SignalPriorities
	if err := json.Unmarshal(data, &p); err != nil {
		return SignalPriorities{}, fmt.Errorf("parse signal_priorities: %w", err)
	}
	return p, nil
}

// gitCommit adds and commits a single file.
func (s *IntegrationStore) gitCommit(absPath, msg string) {
	if s.gitRoot == "" {
		return
	}
	rel, err := filepath.Rel(s.gitRoot, absPath)
	if err != nil {
		return
	}
	if out, err := exec.Command("git", "-C", s.gitRoot, "add", rel).CombinedOutput(); err != nil {
		_ = out // best-effort
		return
	}
	exec.Command("git", "-C", s.gitRoot, "commit", "-m", msg,
		"--author=Emily Prime <emily-prime@agent.local>").Run()
}

// ---------------------------------------------------------------------------
// Triage — Emily Prime's strategic relevance assessment
// ---------------------------------------------------------------------------

func Triage(obs Observation, priorities SignalPriorities) TriageResult {
	result := TriageResult{
		ObservationTimestamp: obs.Timestamp,
	}

	var maxWeight float64
	var triggered []string

	// Score based on signal types present in this observation.
	for sigType, count := range obs.SignalsByType {
		if count <= 0 {
			continue
		}
		if w, ok := priorities.SignalWeights[sigType]; ok && w > maxWeight {
			maxWeight = w
		}
		if priorities.SignalWeights[sigType] >= 0.5 {
			triggered = append(triggered, sigType)
		}
	}

	// Severity bonus.
	switch obs.Severity {
	case "critical", "error":
		if maxWeight < 0.9 {
			maxWeight = 0.9
		}
	case "high", "warn":
		if maxWeight < 0.6 {
			maxWeight = 0.6
		}
	}

	// Observation-type bonus.
	if obs.ObservationType == "escalation" && maxWeight < 0.85 {
		maxWeight = 0.85
	}

	result.StrategicRelevance = maxWeight
	result.TriggeredSignals = triggered

	// CEO visibility criteria.
	ec := priorities.EscalationCriteria
	result.RequiresCEOVisibility = obs.RequiresCEOVisibility ||
		maxWeight >= priorities.CEOVisibilityThreshold ||
		(ec.ActivistRiskAlways && obs.SignalsByType["activist_risk"] > 0) ||
		(ec.AuditorChangeAlways && obs.SignalsByType["auditor_change"] > 0) ||
		(ec.DirectorDepartureAlways && obs.SignalsByType["nomination_rejection"] > 0)

	// Suggest task type.
	switch {
	case obs.ObservationType == "escalation" || obs.Severity == "critical" || obs.Severity == "error":
		result.SuggestedTaskType = "fix_anomaly"
	case obs.SignalsByType["activist_risk"] > 0 || obs.SignalsByType["governance_entrenchment"] > 0:
		result.SuggestedTaskType = "improve_signal"
	case obs.ObservationType == "improvement":
		result.SuggestedTaskType = "improve_signal"
	default:
		result.SuggestedTaskType = "expand_coverage"
	}

	// Rationale.
	var parts []string
	if len(triggered) > 0 {
		parts = append(parts, fmt.Sprintf("high-weight signals: %s", strings.Join(triggered, ", ")))
	}
	if result.RequiresCEOVisibility {
		parts = append(parts, "CEO visibility required")
	}
	if len(obs.AffectedTickers) > 0 {
		parts = append(parts, fmt.Sprintf("tickers: %s", strings.Join(obs.AffectedTickers, ", ")))
	}
	if len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("severity=%s relevance=%.2f", obs.Severity, maxWeight))
	}
	result.Rationale = strings.Join(parts, "; ")

	return result
}

// ---------------------------------------------------------------------------
// PrimeTasks — automatic directed task generation from a triage result
// ---------------------------------------------------------------------------

// TaskFromTriage generates a DirectedTask from a triage result when the
// strategic relevance exceeds a threshold.
func TaskFromTriage(obs Observation, triage TriageResult) *DirectedTask {
	if triage.StrategicRelevance < 0.5 {
		return nil
	}
	priority := "normal"
	switch {
	case triage.StrategicRelevance >= 0.9:
		priority = "critical"
	case triage.StrategicRelevance >= 0.75:
		priority = "high"
	}

	criteria := []string{
		"Signal extraction covers the identified gap",
		"No regression in existing signal counts",
		"Parse errors reduced or eliminated",
	}
	if obs.RequestForClaude != "" {
		criteria = append(criteria, obs.RequestForClaude)
	}

	ctx := fmt.Sprintf(
		"Emily Prime triage: relevance=%.2f, triggered=%s. "+
			"Strategic context: these signals inform M&A risk, activist targeting, and board stability assessment.",
		triage.StrategicRelevance, strings.Join(triage.TriggeredSignals, "+"),
	)

	return &DirectedTask{
		From:               "emily-prime",
		To:                 "fatbaby-emily",
		TaskType:           triage.SuggestedTaskType,
		Priority:           priority,
		Description:        obs.Summary,
		AcceptanceCriteria: criteria,
		Context:            ctx,
	}
}

// ---------------------------------------------------------------------------
// Register integration tools on Emily Prime's dispatcher
// ---------------------------------------------------------------------------

func registerIntegrationTools(d *ToolDispatcher, store *IntegrationStore) {
	d.Register(ToolDef{
		Name:        "integration_read_observations",
		Description: "Read recent observations published by FatBaby-Emily to the integration layer. Returns the N most recent observations from signals/observations/.",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolPropSchema{
				"limit": {Type: "number", Description: "Number of observations to return (default 10)"},
			},
		},
	}, func(args map[string]any) (string, error) {
		limit := 10
		if v, ok := args["limit"].(float64); ok && v > 0 {
			limit = int(v)
		}
		obs, err := store.ReadObservations(limit)
		if err != nil {
			return "", err
		}
		if len(obs) == 0 {
			return "no observations yet — FatBaby-Emily has not published to signals/observations/", nil
		}
		data, _ := json.MarshalIndent(obs, "", "  ")
		return string(data), nil
	})

	d.Register(ToolDef{
		Name:        "integration_write_task",
		Description: "Issue a directed task to FatBaby-Emily. The task is written to signals/tasks/ and committed. FatBaby's observation-watcher will pick it up on next poll.",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolPropSchema{
				"task_type":   {Type: "string", Description: "expand_coverage|improve_signal|fix_anomaly|config_change"},
				"priority":    {Type: "string", Description: "critical|high|normal|low"},
				"description": {Type: "string", Description: "What FatBaby should do"},
				"context":     {Type: "string", Description: "Strategic rationale FatBaby needs to act well"},
				"criteria":    {Type: "string", Description: "Comma-separated acceptance criteria"},
				"deadline":    {Type: "string", Description: "ISO8601 deadline (optional)"},
			},
			Required: []string{"task_type", "description"},
		},
	}, func(args map[string]any) (string, error) {
		task := DirectedTask{
			TaskType:    stringArg(args, "task_type"),
			Priority:    stringArg(args, "priority"),
			Description: stringArg(args, "description"),
			Context:     stringArg(args, "context"),
			Deadline:    stringArg(args, "deadline"),
		}
		if task.Priority == "" {
			task.Priority = "normal"
		}
		if raw := stringArg(args, "criteria"); raw != "" {
			for _, c := range strings.Split(raw, ",") {
				if c = strings.TrimSpace(c); c != "" {
					task.AcceptanceCriteria = append(task.AcceptanceCriteria, c)
				}
			}
		}
		if err := store.WriteTask(task); err != nil {
			return "", err
		}
		return fmt.Sprintf("task written: %s priority=%s", task.TaskID, task.Priority), nil
	})

	d.Register(ToolDef{
		Name:        "integration_triage",
		Description: "Triage the N most recent FatBaby observations against Emily Prime's signal priorities. Returns strategic relevance scores and CEO visibility flags. Auto-generates directed tasks for high-relevance observations.",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolPropSchema{
				"limit":        {Type: "number", Description: "Observations to triage (default 5)"},
				"write_tasks":  {Type: "string", Description: "Set to 'true' to auto-write tasks for high-relevance observations (relevance >= 0.75)"},
			},
		},
	}, func(args map[string]any) (string, error) {
		limit := 5
		if v, ok := args["limit"].(float64); ok && v > 0 {
			limit = int(v)
		}
		writeTasks := stringArg(args, "write_tasks") == "true"

		priorities, err := store.LoadSignalPriorities()
		if err != nil {
			return "", err
		}
		obs, err := store.ReadObservations(limit)
		if err != nil {
			return "", err
		}
		if len(obs) == 0 {
			return "no observations to triage", nil
		}

		type triageEntry struct {
			Observation Observation  `json:"observation"`
			Triage      TriageResult `json:"triage"`
			TaskWritten bool         `json:"task_written,omitempty"`
		}
		var results []triageEntry
		for _, o := range obs {
			t := Triage(o, priorities)
			entry := triageEntry{Observation: o, Triage: t}
			if writeTasks {
				if task := TaskFromTriage(o, t); task != nil {
					if writeErr := store.WriteTask(*task); writeErr == nil {
						entry.TaskWritten = true
					}
				}
			}
			results = append(results, entry)
		}
		data, _ := json.MarshalIndent(results, "", "  ")
		return string(data), nil
	})

	d.Register(ToolDef{
		Name:        "integration_write_observation",
		Description: "Publish an observation to the integration layer from Emily Prime. Useful for publishing Prime's own strategic findings that FatBaby should act on.",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolPropSchema{
				"summary":          {Type: "string", Description: "One-line headline"},
				"detail":           {Type: "string", Description: "Full detail"},
				"severity":         {Type: "string", Description: "critical|high|normal|low"},
				"observation_type": {Type: "string", Description: "anomaly|improvement|escalation|status"},
				"tickers":          {Type: "string", Description: "Comma-separated affected tickers"},
			},
			Required: []string{"summary", "severity"},
		},
	}, func(args map[string]any) (string, error) {
		obs := Observation{
			Source:          "emily-prime",
			Summary:         stringArg(args, "summary"),
			Detail:          stringArg(args, "detail"),
			Severity:        stringArg(args, "severity"),
			ObservationType: stringArg(args, "observation_type"),
		}
		if raw := stringArg(args, "tickers"); raw != "" {
			for _, t := range strings.Split(raw, ",") {
				if t = strings.TrimSpace(t); t != "" {
					obs.AffectedTickers = append(obs.AffectedTickers, t)
				}
			}
		}
		if err := store.WriteObservation(obs); err != nil {
			return "", err
		}
		return "observation written to signals/observations/", nil
	})
}

// ---------------------------------------------------------------------------
// HTTP handlers
// ---------------------------------------------------------------------------

func (s *Server) handleIntegrationObservations(w http.ResponseWriter, r *http.Request) {
	if s.integration == nil {
		http.Error(w, `{"error":"integration not configured"}`, 503)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		obs, err := s.integration.ReadObservations(20)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), 500)
			return
		}
		json.NewEncoder(w).Encode(obs)
	case http.MethodPost:
		var obs Observation
		if err := json.NewDecoder(r.Body).Decode(&obs); err != nil {
			http.Error(w, `{"error":"invalid json"}`, 400)
			return
		}
		if err := s.integration.WriteObservation(obs); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "written"})
	default:
		http.Error(w, "method not allowed", 405)
	}
}

func (s *Server) handleIntegrationTask(w http.ResponseWriter, r *http.Request) {
	if s.integration == nil {
		http.Error(w, `{"error":"integration not configured"}`, 503)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	var task DirectedTask
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, `{"error":"invalid json"}`, 400)
		return
	}
	if strings.TrimSpace(task.Description) == "" {
		http.Error(w, `{"error":"description required"}`, 400)
		return
	}
	if err := s.integration.WriteTask(task); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), 500)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "written", "task_id": task.TaskID})
}

func (s *Server) handleIntegrationTriage(w http.ResponseWriter, r *http.Request) {
	if s.integration == nil {
		http.Error(w, `{"error":"integration not configured"}`, 503)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", 405)
		return
	}
	priorities, err := s.integration.LoadSignalPriorities()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), 500)
		return
	}
	obs, err := s.integration.ReadObservations(10)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), 500)
		return
	}
	type entry struct {
		O Observation  `json:"observation"`
		T TriageResult `json:"triage"`
	}
	var results []entry
	for _, o := range obs {
		results = append(results, entry{O: o, T: Triage(o, priorities)})
	}
	json.NewEncoder(w).Encode(results)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func slugify(s string, maxLen int) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
		} else if b.Len() > 0 {
			b.WriteRune('-')
		}
		if b.Len() >= maxLen {
			break
		}
	}
	return strings.TrimRight(b.String(), "-")
}

func stringArg(args map[string]any, key string) string {
	v, _ := args[key].(string)
	return v
}

// buildIntegrationStore constructs the integration store from env, or returns nil.
func buildIntegrationStore() *IntegrationStore {
	// Default: signals/ dir relative to the process working directory parent.
	// When emily-agent runs from emily-agent/, the EMILY repo root is ..
	gitRoot := os.Getenv("EMILY_GIT_ROOT")
	if gitRoot == "" {
		// Try to auto-detect: look for the signals/ dir one level up.
		if _, err := os.Stat("../signals"); err == nil {
			gitRoot = ".."
		} else {
			gitRoot = "."
		}
	}
	signalsDir := os.Getenv("EMILY_INTEGRATION_DIR")
	if signalsDir == "" {
		signalsDir = filepath.Join(gitRoot, "signals")
	}
	store, err := NewIntegrationStore(signalsDir, gitRoot)
	if err != nil {
		return nil
	}
	return store
}

// PrimeTriageContext holds context for the autonomous triage cycle.
// Called from cron.go when the cycle picks up triage work.
// runPrimeTriage reads FatBaby observations, scores strategic relevance, writes
// directed tasks for high-relevance findings, and sends Gmail alerts for CEO-visible
// observations that haven't been escalated before (deduped via a cursor file).
// gmail may be nil — in that case escalations are counted but no email is sent.
// PushFunc is a best-effort push notification callback. title, body, and data
// are forwarded to FCM. If nil, no push is attempted.
type PushFunc func(title, body string, data map[string]string)

// runBacklogCuration reads uncurated FatBaby observations and appends non-trivial
// ones to EMILY/BACKLOG.md as INTAKE QUEUE items. Idempotent via state file.
// Called from runPrimeTriage each cycle — max 5 items per cycle to stay bounded.
func (s *IntegrationStore) runBacklogCuration() (int, error) {
	if s.fatBabyRoot == "" || s.gitRoot == "" {
		return 0, nil
	}
	obsDir := filepath.Join(s.fatBabyRoot, "var", "emily-observations")
	backlogPath := filepath.Join(s.gitRoot, "BACKLOG.md")
	stateFile := filepath.Join(s.gitRoot, "var", "backlog-curated.txt")

	if err := os.MkdirAll(filepath.Dir(stateFile), 0755); err != nil {
		return 0, fmt.Errorf("mkdir var: %w", err)
	}

	curated := loadCuratedSet(stateFile)
	entries, err := os.ReadDir(obsDir)
	if err != nil {
		return 0, nil // obs dir may not exist yet
	}

	// Sort by mtime descending (newest first).
	type obsEntry struct {
		name  string
		mtime int64
	}
	var oes []obsEntry
	for _, e := range entries {
		n := e.Name()
		if e.IsDir() || strings.HasPrefix(n, ".") || n == "latest.json" || !strings.HasSuffix(n, ".json") {
			continue
		}
		info, _ := e.Info()
		if info != nil {
			oes = append(oes, obsEntry{name: n, mtime: info.ModTime().UnixNano()})
		}
	}
	sort.Slice(oes, func(i, j int) bool { return oes[i].mtime > oes[j].mtime })

	const maxPerCycle = 5
	today := time.Now().UTC().Format("2006-01-02")
	var lines []string
	var newKeys []string

	for _, oe := range oes {
		if len(lines) >= maxPerCycle {
			break
		}
		key := strings.TrimSuffix(oe.name, ".json")
		if _, seen := curated[key]; seen {
			continue
		}
		data, err := os.ReadFile(filepath.Join(obsDir, oe.name))
		if err != nil {
			newKeys = append(newKeys, key) // mark as seen even if unreadable
			continue
		}
		var obs struct {
			Timestamp string `json:"timestamp"`
			Summary   string `json:"summary"`
		}
		if err := json.Unmarshal(data, &obs); err != nil || obs.Summary == "" {
			newKeys = append(newKeys, key)
			continue
		}
		if isTrivialBacklogObs(obs.Summary) {
			newKeys = append(newKeys, key) // mark as seen but don't append
			continue
		}
		summary := obs.Summary
		if len(summary) > 120 {
			summary = summary[:119] + "…"
		}
		// Escape markdown bold markers.
		summary = strings.ReplaceAll(summary, "**", "*")
		ts := obs.Timestamp
		if ts == "" {
			ts = key
		}
		lines = append(lines, fmt.Sprintf("- [ ] **%s** — obs `%s`. CURATED: %s.", summary, ts, today))
		newKeys = append(newKeys, key)
	}

	if len(lines) > 0 {
		if err := appendToBacklogMd(backlogPath, lines); err != nil {
			log.Printf("backlog curation: append failed: %v", err)
		} else {
			msg := fmt.Sprintf("emily-prime: backlog curate %d obs -> INTAKE QUEUE (%s)", len(lines), today)
			emilyGitAddCommit(s.gitRoot, "BACKLOG.md", msg)
		}
	}

	// Update state file with all processed keys (including trivials and empties).
	if len(newKeys) > 0 {
		f, err := os.OpenFile(stateFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			w := bufio.NewWriter(f)
			for _, k := range newKeys {
				fmt.Fprintln(w, k)
			}
			w.Flush()
			f.Close()
		}
	}

	return len(lines), nil
}

// isTrivialBacklogObs returns true for RSI completion pings and other status noise.
func isTrivialBacklogObs(summary string) bool {
	lower := strings.ToLower(summary)
	trivial := []string{
		"rsi loop iteration", "tic-toc cycle done", "fatbaby+tyler tic-toc",
		"prime task complete", "rsi cycle complete", "rsi session complete",
		"rsi session ", "emily.cli v0.", "emily.cli rsi", "emily.cli section",
		"emily.cli first real", "emily v0.2", "emily v0.3", "rsi: rsi loop",
	}
	for _, s := range trivial {
		if strings.Contains(lower, s) {
			return true
		}
	}
	// Single-word typos / test observations.
	for _, t := range []string{"tuo", "test", "ok"} {
		if strings.EqualFold(strings.TrimSpace(summary), t) {
			return true
		}
	}
	return false
}

// appendToBacklogMd adds lines to the INTAKE QUEUE section of BACKLOG.md,
// creating the section if absent.
func appendToBacklogMd(path string, lines []string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(data)
	block := strings.Join(lines, "\n") + "\n"

	const intakeMarker = "## INTAKE QUEUE"
	if idx := strings.Index(content, intakeMarker); idx >= 0 {
		searchFrom := idx + len(intakeMarker)
		insertPos := len(content)
		if next := strings.Index(content[searchFrom:], "\n## "); next >= 0 {
			insertPos = searchFrom + next + 1
		}
		content = content[:insertPos] + block + content[insertPos:]
	} else {
		const protocol = "## BACKLOG PROTOCOL"
		const header = "\n## INTAKE QUEUE (curated by emily backlog curate)\n\nItems below have been auto-curated from FatBaby observations. Emily Prime reviews\nand promotes them into the appropriate section when she plans the next sprint.\n\n"
		if idx2 := strings.Index(content, protocol); idx2 >= 0 {
			content = content[:idx2] + header + block + "\n---\n\n" + content[idx2:]
		} else {
			content += header + block
		}
	}
	return os.WriteFile(path, []byte(content), 0644)
}

// loadCuratedSet reads EMILY/var/backlog-curated.txt and returns the set of file keys.
func loadCuratedSet(stateFile string) map[string]struct{} {
	m := make(map[string]struct{})
	f, err := os.Open(stateFile)
	if err != nil {
		return m
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		if k := strings.TrimSpace(s.Text()); k != "" {
			m[k] = struct{}{}
		}
	}
	return m
}

func runPrimeTriage(ctx context.Context, store *IntegrationStore, pipeline *Pipeline, gmail *GmailClient, push PushFunc) (string, error) {
	priorities, err := store.LoadSignalPriorities()
	if err != nil {
		return "", fmt.Errorf("load priorities: %w", err)
	}
	obs, err := store.ReadObservations(10)
	if err != nil {
		return "", fmt.Errorf("read observations: %w", err)
	}
	if len(obs) == 0 {
		return "no observations to triage", nil
	}

	// Load escalation cursor so we don't re-send alerts for observations already escalated.
	escalationCursorPath := filepath.Join(store.signalsDir, "observations", ".escalation-cursor")
	lastEscalated := ""
	if b, err := os.ReadFile(escalationCursorPath); err == nil {
		lastEscalated = strings.TrimSpace(string(b))
	}

	tasksWritten := 0
	escalations := 0
	alertsSent := 0
	newLastEscalated := lastEscalated

	// obs is sorted newest-first; iterate newest-first so the cursor advances to latest.
	for _, o := range obs {
		t := Triage(o, priorities)
		if task := TaskFromTriage(o, t); task != nil {
			if writeErr := store.WriteTask(*task); writeErr == nil {
				tasksWritten++
			}
		}
		if t.RequiresCEOVisibility {
			escalations++
			// Only alert for observations newer than the escalation cursor.
			if o.Timestamp > lastEscalated {
				if gmail != nil {
					alert := FormatSignalAlert(o, t)
					if sendErr := gmail.SendAlert(ctx, alert); sendErr != nil {
						log.Printf("triage: gmail send failed: %v", sendErr)
					} else {
						alertsSent++
						if o.Timestamp > newLastEscalated {
							newLastEscalated = o.Timestamp
						}
					}
				}
				// FCM push to MJOLNIR — best-effort, never blocks triage.
				if push != nil {
					title := fmt.Sprintf("[%s] %s", strings.ToUpper(o.Severity), o.Source)
					body := o.Summary
					if len(body) > 140 {
						body = body[:139] + "…"
					}
					go push(title, body, map[string]string{
						"severity":   o.Severity,
						"source":     o.Source,
						"obs_type":   o.ObservationType,
						"timestamp":  o.Timestamp,
					})
				}
			}
		}
	}

	// Persist cursor so the next cycle doesn't re-send the same alerts.
	if newLastEscalated != lastEscalated {
		_ = os.WriteFile(escalationCursorPath, []byte(newLastEscalated), 0o644)
	}

	// Autonomous backlog curation: append uncurated FatBaby obs to INTAKE QUEUE.
	curated := 0
	if store.fatBabyRoot != "" {
		if n, err := store.runBacklogCuration(); err != nil {
			log.Printf("triage: backlog curation error (non-fatal): %v", err)
		} else {
			curated = n
		}
	}

	// Autonomous backlog promotion: classify INTAKE QUEUE items into numbered sections.
	// Throttled internally to at most once per 30 min; no-ops if emily CLI unavailable.
	promoted := 0
	if store.gitRoot != "" {
		if n, err := store.runBacklogPromotion(ctx); err != nil {
			log.Printf("triage: backlog promotion error (non-fatal): %v", err)
		} else if n > 0 {
			promoted = n
			log.Printf("triage: promoted %d items from INTAKE QUEUE", n)
		}
	}

	return fmt.Sprintf("triage complete observations=%d tasks_written=%d escalations=%d alerts_sent=%d backlog_curated=%d backlog_promoted=%d",
		len(obs), tasksWritten, escalations, alertsSent, curated, promoted), nil
}

// runBacklogPromotion calls `emily backlog promote` to classify INTAKE QUEUE items.
// Throttled to at most once per 30 minutes via a cursor file so haiku spend stays bounded.
func (s *IntegrationStore) runBacklogPromotion(ctx context.Context) (int, error) {
	// Throttle: skip if promoted within the last 30 minutes.
	cursorPath := filepath.Join(s.gitRoot, "var", "backlog-promote-cursor")
	if data, err := os.ReadFile(cursorPath); err == nil {
		if t, err := time.Parse(time.RFC3339, strings.TrimSpace(string(data))); err == nil {
			if time.Since(t) < 30*time.Minute {
				return 0, nil
			}
		}
	}

	emilyBin, err := exec.LookPath("emily")
	if err != nil {
		return 0, nil // emily CLI not installed; skip silently
	}

	cmd := exec.CommandContext(ctx, emilyBin, "backlog", "promote",
		"--limit=15", "--batch=10", "--no-apple", "--json")
	cmd.Env = os.Environ() // inherit ANTHROPIC_API_KEY, EMILY_ROOT, etc.

	out, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("emily backlog promote: %w", err)
	}

	// Update cursor on success.
	_ = os.WriteFile(cursorPath, []byte(time.Now().UTC().Format(time.RFC3339)), 0644)

	var result struct {
		Promoted int `json:"promoted"`
		Removed  int `json:"removed"`
	}
	if jsonErr := json.Unmarshal(out, &result); jsonErr != nil {
		return 0, fmt.Errorf("parse promote output: %w", jsonErr)
	}
	return result.Promoted + result.Removed, nil
}
