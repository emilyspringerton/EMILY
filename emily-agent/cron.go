// emily-agent/cron.go
// Cron-based autonomous execution cycle.
//
// Every 5 minutes (configurable) Emily wakes up, observes her state,
// picks the highest-priority task from her roadmap, runs one RSI iteration,
// updates state, and sleeps.
//
// The lock file prevents concurrent runs. The state file persists across
// cycles so Emily remembers what she was working on.
//
// Usage: start with go run main.go rsi.go cron.go -- --cron
// Or via crontab: */5 * * * * /opt/emily/emily-agent --cron

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CronConfig controls the autonomous cycle.
type CronConfig struct {
	StateDir      string        // directory for state files
	LockFile      string        // prevents concurrent runs
	CycleDuration time.Duration // max time per cycle (default 280s)
	Interval      time.Duration // how often to run (default 5m, used in daemon mode)
}

func defaultCronConfig() CronConfig {
	stateDir := envOr("EMILY_STATE_DIR", "./emily-state")
	return CronConfig{
		StateDir:      stateDir,
		LockFile:      filepath.Join(stateDir, "emily.lock"),
		CycleDuration: 280 * time.Second,
		Interval:      5 * time.Minute,
	}
}

// CyclePhase labels the four phases of each cron cycle.
type CyclePhase string

const (
	PhaseObserve CyclePhase = "observe"
	PhaseDecide  CyclePhase = "decide"
	PhaseAct     CyclePhase = "act"
	PhasePlan    CyclePhase = "plan"
)

// RoadmapItem is one initiative on Emily's evolution roadmap.
type RoadmapItem struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Priority    int                  `json:"priority"` // lower = higher priority
	Status      string               `json:"status"`   // queued|in_progress|complete|blocked
	Criteria    []AcceptanceCriterion `json:"criteria"`
	MaxIters    int                  `json:"max_iters"`
	Notes       string               `json:"notes,omitempty"`
}

// CycleState is persisted between cycles so Emily remembers her context.
type CycleState struct {
	Version          int           `json:"version"`
	CycleNumber      int           `json:"cycle_number"`
	LastCycleAt      time.Time     `json:"last_cycle_at"`
	ActiveTaskID     string        `json:"active_task_id,omitempty"`
	ActiveTask       *ImprovementTask `json:"active_task,omitempty"`
	Roadmap          []RoadmapItem `json:"roadmap"`
	CompletedTasks   []ImprovementTask `json:"completed_tasks,omitempty"`
	Metrics          CycleMetrics  `json:"metrics"`
	NextCyclePlan    string        `json:"next_cycle_plan,omitempty"`
}

// CycleMetrics tracks health across cycles.
type CycleMetrics struct {
	TotalCycles      int `json:"total_cycles"`
	SuccessfulCycles int `json:"successful_cycles"`
	TasksCompleted   int `json:"tasks_completed"`
	ItersRun         int `json:"iters_run"`
	ConsecFailures   int `json:"consecutive_failures"`
}

// CycleRecord is the log entry for one cycle.
type CycleRecord struct {
	Number    int        `json:"number"`
	StartedAt time.Time  `json:"started_at"`
	EndedAt   time.Time  `json:"ended_at"`
	Phase     CyclePhase `json:"last_phase"`
	Action    string     `json:"action"`
	Outcome   string     `json:"outcome"`
	Error     string     `json:"error,omitempty"`
}

// AutonomousCycle runs one 5-minute execution cycle.
type AutonomousCycle struct {
	cfg         CronConfig
	pipeline    *Pipeline
	state       *CycleState
	integration *IntegrationStore // optional; enables prime triage each cycle
	emiree      *EmireeAgent      // witch engine governing RSI operational state
}

// NewAutonomousCycle creates a cycle runner.
func NewAutonomousCycle(cfg CronConfig, p *Pipeline) *AutonomousCycle {
	return &AutonomousCycle{
		cfg:         cfg,
		pipeline:    p,
		integration: buildIntegrationStore(),
		emiree:      NewEmireeAgent(cfg.StateDir),
	}
}

// RunOnce executes exactly one cycle. Designed to be called by cron.
func (ac *AutonomousCycle) RunOnce() error {
	if err := os.MkdirAll(ac.cfg.StateDir, 0o755); err != nil {
		return fmt.Errorf("state dir: %w", err)
	}

	// Acquire lock — prevents concurrent runs
	if err := ac.acquireLock(); err != nil {
		return fmt.Errorf("lock: %w", err)
	}
	defer ac.releaseLock()

	ctx, cancel := context.WithTimeout(context.Background(), ac.cfg.CycleDuration)
	defer cancel()

	rec := CycleRecord{StartedAt: time.Now()}

	// PHASE 1: OBSERVE — load state, check health
	rec.Phase = PhaseObserve
	state, err := ac.loadState()
	if err != nil {
		return fmt.Errorf("observe: %w", err)
	}
	ac.state = state
	state.CycleNumber++
	state.Metrics.TotalCycles++
	log.Printf("[cycle %d] observe: active_task=%s roadmap_items=%d",
		state.CycleNumber, state.ActiveTaskID, len(state.Roadmap))

	// PHASE 2: DECIDE — pick what to work on this cycle (gear-aware)
	rec.Phase = PhaseDecide
	influence := ac.emiree.State.Influence()
	log.Printf("[cycle %d] emiree: %s", state.CycleNumber, ac.emiree.State.Summary())
	task, reason := ac.pickTask(state)
	// Apply gear influence to task max_iters
	if task != nil && influence.MaxIters > 0 {
		task.MaxIters = influence.MaxIters
	}
	rec.Action = reason
	log.Printf("[cycle %d] decide: %s", state.CycleNumber, reason)

	// PHASE 3: ACT — run one RSI iteration
	rec.Phase = PhaseAct
	if task != nil {
		result, err := ac.runIteration(ctx, state, task)
		if err != nil {
			state.Metrics.ConsecFailures++
			rec.Error = err.Error()
			rec.Outcome = "error"
			log.Printf("[cycle %d] act: error: %v", state.CycleNumber, err)
		} else {
			state.Metrics.ConsecFailures = 0
			state.Metrics.ItersRun++
			rec.Outcome = result
			log.Printf("[cycle %d] act: %s", state.CycleNumber, result)
			if task.Status == "success" {
				state.Metrics.TasksCompleted++
				state.CompletedTasks = append(state.CompletedTasks, *task)
				state.ActiveTaskID = ""
				state.ActiveTask = nil
			}
		}
	} else {
		rec.Outcome = "idle — no queued tasks"
		log.Printf("[cycle %d] act: idle", state.CycleNumber)
	}

	// PHASE 4: PLAN — update state, run Prime triage if integration is wired
	rec.Phase = PhasePlan
	state.LastCycleAt = time.Now()
	if task != nil && task.Status == "running" {
		state.ActiveTask = task
	}
	state.Metrics.SuccessfulCycles++
	rec.EndedAt = time.Now()

	triageFindings := 0
	if ac.integration != nil {
		if triageResult, triageErr := ac.runPrimeTriageCycle(ctx, ac.integration); triageErr != nil {
			log.Printf("[cycle %d] triage warn: %v", state.CycleNumber, triageErr)
		} else {
			log.Printf("[cycle %d] triage: %s", state.CycleNumber, triageResult)
			// Count how many tasks were issued to feed back into Emiree
			if strings.Contains(triageResult, "tasks_written=") {
				fmt.Sscanf(triageResult[strings.Index(triageResult, "tasks_written="):], "tasks_written=%d", &triageFindings)
			}
		}
	}

	// Feed outcome back into Emiree; it updates state, saves, returns next gear.
	rsiOutcome := buildRSIOutcome(task, triageFindings)
	nextInfluence := ac.emiree.Tick(rsiOutcome)
	log.Printf("[cycle %d] emiree after: %s | next: max_iters=%d pace=%ds",
		state.CycleNumber, ac.emiree.State.Summary(), nextInfluence.MaxIters, nextInfluence.PaceSeconds)

	if err := ac.saveState(state); err != nil {
		return fmt.Errorf("save state: %w", err)
	}
	if err := ac.appendCycleRecord(rec); err != nil {
		log.Printf("[cycle %d] warn: could not save cycle record: %v", state.CycleNumber, err)
	}
	ac.updateDashboard(state, rec)

	log.Printf("[cycle %d] complete in %s", state.CycleNumber, time.Since(rec.StartedAt).Round(time.Second))
	return nil
}

// RunDaemon runs cycles on the configured interval. Use for development;
// in production prefer cron + RunOnce.
func (ac *AutonomousCycle) RunDaemon() {
	log.Printf("Emily daemon starting — cycle interval %s", ac.cfg.Interval)
	for {
		start := time.Now()
		if err := ac.RunOnce(); err != nil {
			log.Printf("cycle error: %v", err)
		}
		elapsed := time.Since(start)
		sleep := ac.cfg.Interval - elapsed
		if sleep > 0 {
			time.Sleep(sleep)
		}
	}
}

// pickTask selects what to work on this cycle.
// Priority: resume active task > pick highest-priority queued item.
func (ac *AutonomousCycle) pickTask(state *CycleState) (*ImprovementTask, string) {
	// Resume in-progress task
	if state.ActiveTask != nil && state.ActiveTask.Status == "running" {
		return state.ActiveTask, fmt.Sprintf("resume active task: %s", state.ActiveTask.ID)
	}

	// Find highest-priority queued roadmap item
	var best *RoadmapItem
	for i := range state.Roadmap {
		item := &state.Roadmap[i]
		if item.Status != "queued" {
			continue
		}
		if best == nil || item.Priority < best.Priority {
			best = item
		}
	}
	if best == nil {
		return nil, "no queued tasks"
	}

	// Promote roadmap item to active task
	best.Status = "in_progress"
	maxIters := best.MaxIters
	if maxIters <= 0 {
		maxIters = 10
	}
	task := &ImprovementTask{
		ID:          best.ID,
		Description: best.Description,
		Criteria:    best.Criteria,
		MaxIters:    maxIters,
		Status:      "pending",
	}
	state.ActiveTaskID = task.ID
	state.ActiveTask = task
	return task, fmt.Sprintf("start new task: %s (%s)", best.ID, best.Name)
}

// runIteration runs one RSI iteration on the active task.
func (ac *AutonomousCycle) runIteration(ctx context.Context, state *CycleState, task *ImprovementTask) (string, error) {
	loop := NewRSILoop(ac.pipeline)

	// Run a single iteration (not the full loop — one cycle = one iteration)
	task.Status = "running"
	rec, err := loop.runIteration(ctx, task)
	if err != nil {
		return "", err
	}
	task.Iterations = append(task.Iterations, rec)

	if rec.AllPass {
		task.Status = "success"
		t := time.Now()
		task.CompletedAt = &t
		task.Lessons = loop.extractLessons(ctx, task)
		// Mark roadmap item complete
		for i := range state.Roadmap {
			if state.Roadmap[i].ID == task.ID {
				state.Roadmap[i].Status = "complete"
			}
		}
		return fmt.Sprintf("task %s COMPLETE in %d iterations", task.ID, len(task.Iterations)), nil
	}

	if len(task.Iterations) >= task.MaxIters {
		task.Status = "partial"
		t := time.Now()
		task.CompletedAt = &t
		return fmt.Sprintf("task %s reached max_iters (%d), status=partial", task.ID, task.MaxIters), nil
	}

	return fmt.Sprintf("task %s iteration %d: %d/%d criteria pass",
		task.ID, rec.Number, passCount(rec.Results), len(rec.Results)), nil
}

func passCount(results []CriteriaResult) int {
	n := 0
	for _, r := range results {
		if r.Passes {
			n++
		}
	}
	return n
}

// --- State persistence ---

func (ac *AutonomousCycle) loadState() (*CycleState, error) {
	path := filepath.Join(ac.cfg.StateDir, "current-evolution.json")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return ac.defaultState(), nil
	}
	if err != nil {
		return nil, err
	}
	var state CycleState
	if err := json.Unmarshal(data, &state); err != nil {
		log.Printf("warn: corrupt state file, resetting: %v", err)
		return ac.defaultState(), nil
	}
	return &state, nil
}

func (ac *AutonomousCycle) defaultState() *CycleState {
	return &CycleState{
		Version: 1,
		Roadmap: defaultRoadmap(),
	}
}

func (ac *AutonomousCycle) saveState(state *CycleState) error {
	path := filepath.Join(ac.cfg.StateDir, "current-evolution.json")
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func (ac *AutonomousCycle) appendCycleRecord(rec CycleRecord) error {
	path := filepath.Join(ac.cfg.StateDir, "cycle-log.jsonl")
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%s\n", data)
	return err
}

// updateDashboard writes a human-readable status file.
func (ac *AutonomousCycle) updateDashboard(state *CycleState, rec CycleRecord) {
	path := filepath.Join(ac.cfg.StateDir, "dashboard.txt")
	var sb strings.Builder
	fmt.Fprintf(&sb, "EMILY STATUS DASHBOARD\n")
	fmt.Fprintf(&sb, "Updated: %s\n\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(&sb, "Cycle:          #%d\n", state.CycleNumber)
	fmt.Fprintf(&sb, "Last cycle:     %s\n", state.LastCycleAt.Format(time.RFC3339))
	fmt.Fprintf(&sb, "Total cycles:   %d\n", state.Metrics.TotalCycles)
	fmt.Fprintf(&sb, "Tasks complete: %d\n", state.Metrics.TasksCompleted)
	fmt.Fprintf(&sb, "Iters run:      %d\n\n", state.Metrics.ItersRun)

	if state.ActiveTask != nil {
		t := state.ActiveTask
		fmt.Fprintf(&sb, "ACTIVE TASK: %s\n", t.ID)
		fmt.Fprintf(&sb, "Status:      %s\n", t.Status)
		fmt.Fprintf(&sb, "Iterations:  %d / %d\n", len(t.Iterations), t.MaxIters)
		if len(t.Iterations) > 0 {
			last := t.Iterations[len(t.Iterations)-1]
			fmt.Fprintf(&sb, "Last iter:   %d/%d criteria pass\n\n", passCount(last.Results), len(last.Results))
		}
	} else {
		fmt.Fprintf(&sb, "ACTIVE TASK: none\n\n")
	}

	fmt.Fprintf(&sb, "ROADMAP:\n")
	for _, item := range state.Roadmap {
		fmt.Fprintf(&sb, "  [%s] p%d %s — %s\n", item.Status, item.Priority, item.ID, item.Name)
	}

	fmt.Fprintf(&sb, "\nLAST CYCLE:\n")
	fmt.Fprintf(&sb, "  Action:  %s\n", rec.Action)
	fmt.Fprintf(&sb, "  Outcome: %s\n", rec.Outcome)
	if rec.Error != "" {
		fmt.Fprintf(&sb, "  Error:   %s\n", rec.Error)
	}

	_ = os.WriteFile(path, []byte(sb.String()), 0o644)
}

// --- Lock file ---

func (ac *AutonomousCycle) acquireLock() error {
	path := ac.cfg.LockFile
	// Check for stale lock (older than cycle duration + buffer)
	if info, err := os.Stat(path); err == nil {
		age := time.Since(info.ModTime())
		if age > ac.cfg.CycleDuration+30*time.Second {
			log.Printf("removing stale lock (age %s)", age.Round(time.Second))
			_ = os.Remove(path)
		} else {
			return fmt.Errorf("already running (lock age %s)", age.Round(time.Second))
		}
	}
	return os.WriteFile(path, []byte(fmt.Sprintf("%d", os.Getpid())), 0o644)
}

func (ac *AutonomousCycle) releaseLock() {
	_ = os.Remove(ac.cfg.LockFile)
}

// --- Default roadmap: the bootstrap sequence from emily-ground-zero-protocol.md ---

// runPrimeTriageCycle is called from the autonomous cycle when the prime-triage
// roadmap item is active. It reads FatBaby observations, triages them, and
// issues directed tasks for high-relevance findings.
func (ac *AutonomousCycle) runPrimeTriageCycle(ctx context.Context, store *IntegrationStore) (string, error) {
	return runPrimeTriage(ctx, store, ac.pipeline)
}

func defaultRoadmap() []RoadmapItem {
	return []RoadmapItem{
		{
			ID:          "prime-triage",
			Name:        "Emily Prime triage loop — active each cycle",
			Description: "Each cron cycle: read FatBaby observations from signals/observations/, score strategic relevance against signal-priorities.json, issue DirectedTasks to signals/tasks/ for high-relevance findings, flag CEO-visibility observations for Gmail escalation.",
			Priority:    0,
			Status:      "in_progress",
			Criteria: []AcceptanceCriterion{
				{Name: "reads_observations", Description: "Emily Prime reads FatBaby observations from signals/observations/", Target: "observation count increases after fatbaby_commit_to_prime"},
				{Name: "issues_tasks", Description: "directed tasks written to signals/tasks/ for relevance >= 0.75", Target: "task file present with correct schema"},
				{Name: "ceo_visibility", Description: "high-relevance observations trigger Gmail alert", Target: "alert sent for activist_risk and auditor_change observations"},
			},
			MaxIters: 3,
		},
		{
			ID:          "rsi-self-improve",
			Name:        "RSI engine self-improvement",
			Description: "Improve the RSI engine's own system prompts and iteration logic.",
			Priority:    1,
			Status:      "complete",
			Criteria: []AcceptanceCriterion{
				{Name: "convergence", Description: "test tasks converge in fewer iterations", Target: "< 4 avg iterations"},
				{Name: "quality", Description: "artifacts satisfy criteria on first evaluation", Target: "> 70% first-pass rate"},
			},
			MaxIters: 5,
			Notes:    "Completed: tightened generator/evaluator prompts, fixed buildGenerationPrompt target/gap bug, added structural JSON guard, improved evaluate() JSON extraction robustness.",
		},
		{
			ID:          "reddit-collector",
			Name:        "Reddit data collection pipeline — tune and harden",
			Description: "The Go RedditSource in collector.go is live. Improve it: tune quality thresholds (score/comments/body_len) against real collection runs, add exponential backoff on 429s, add per-subreddit collection metrics, and ensure zero data loss on transient errors.",
			Priority:    2,
			Status:      "in_progress",
			Criteria: []AcceptanceCriterion{
				{Name: "throughput", Description: "collection rate from configured subreddits", Target: ">= 500 posts/hour sustained"},
				{Name: "quality_filter", Description: "posts below quality bar are excluded before write", Target: "score/comments/body filters applied at source"},
				{Name: "dedup", Description: "same post is never written twice across restarts", Target: "seenSet loaded from corpus at startup"},
				{Name: "resilience", Description: "rate limits and transient errors handled gracefully", Target: "exponential backoff, zero panics, fail-open"},
			},
			MaxIters: 5,
		},
		{
			ID:          "wikipedia-collector",
			Name:        "Wikipedia collection pipeline — tune and expand",
			Description: "The Go WikipediaSource in collector.go fetches recent changes + plain-text extracts. Improve it: add category filtering, expand to 100 pages/batch, add article quality signals (length, section count, citation count from wikitext), ensure stub articles are filtered.",
			Priority:    3,
			Status:      "in_progress",
			Criteria: []AcceptanceCriterion{
				{Name: "extraction", Description: "article text is clean plain text, no wikitext artifacts", Target: "extract field contains prose only"},
				{Name: "stub_filter", Description: "stub articles are excluded", Target: "length < 1000 chars filtered at source"},
				{Name: "throughput", Description: "articles collected per hour", Target: ">= 500 articles/hour"},
				{Name: "schema", Description: "CollectedDoc fields populated", Target: "id, title, body, body_bytes, quality_score, quality_tier all present"},
			},
			MaxIters: 5,
		},
		{
			ID:          "quality-scorer",
			Name:        "Quality scorer — calibrate and extend",
			Description: "The Go QualityScorer in collector.go is live (deterministic, no ML). Calibrate thresholds against real corpus samples: verify gold/silver/bronze distribution, tune spam regex patterns, add domain-specific signals for ML content (code blocks, equations, citations), verify false-positive rate < 5% on known-good samples.",
			Priority:    4,
			Status:      "in_progress",
			Criteria: []AcceptanceCriterion{
				{Name: "scoring", Description: "every record gets a numeric quality score 0-10", Target: "0.0-10.0 float, mean score > 5.0 for reddit/wiki corpus"},
				{Name: "tier_distribution", Description: "gold/silver/bronze distribution is reasonable", Target: "gold 10-20%, silver 30-50%, bronze remainder"},
				{Name: "spam_detection", Description: "obvious spam/ads are scored low", Target: "score < 4.0 for samples containing 3+ spam phrases"},
				{Name: "false_positive_rate", Description: "legitimate technical content is not penalized", Target: "< 5% of gold-quality human-labeled docs score < 5.0"},
			},
			MaxIters: 5,
		},
		{
			ID:          "arxiv-collector",
			Name:        "ArXiv collection pipeline — live and collecting",
			Description: "ArXivSource is live in collector.go fetching cs.AI/cs.LG/cs.CL/cs.CV/stat.ML abstracts. Improve: add full-paper HTML extraction via arxiv HTML endpoint, add citation count heuristic to quality scorer, tune category list against corpus quality metrics.",
			Priority:    4,
			Status:      "in_progress",
			Criteria: []AcceptanceCriterion{
				{Name: "throughput", Description: "papers collected per hour", Target: ">= 200 papers/hour from default categories"},
				{Name: "body_quality", Description: "body field contains clean prose, no LaTeX artifacts", Target: "< 5% of collected docs contain raw LaTeX macros"},
				{Name: "dedup", Description: "same paper not written twice across restarts", Target: "seenSet deduplication verified on restart"},
			},
			MaxIters: 4,
		},
		{
			ID:          "bob-agent",
			Name:        "Bob database admin agent",
			Description: "Design and implement Bob: a minimal Go HTTP service that accepts structured task requests from Emily, executes PostgreSQL operations, returns structured results, and maintains an audit log.",
			Priority:    6,
			Status:      "queued",
			Criteria: []AcceptanceCriterion{
				{Name: "api", Description: "Emily can POST a task to Bob and GET the result", Target: "POST /task, GET /task/{id}"},
				{Name: "isolation", Description: "Bob only executes database operations", Target: "no filesystem or network access beyond DB"},
				{Name: "audit", Description: "every operation is logged with timestamp, operation, result", Target: "100% operations logged"},
				{Name: "safety", Description: "destructive operations require explicit confirmation flag", Target: "drop/delete require {confirm: true}"},
			},
			MaxIters: 6,
		},
	}
}

// buildRSIOutcome maps a completed task (or nil) + triage count to RSIOutcome.
func buildRSIOutcome(task *ImprovementTask, triageFindings int) RSIOutcome {
	out := RSIOutcome{TriageFindings: triageFindings}
	if task == nil {
		return out
	}
	out.TaskID = task.ID
	out.MaxIterations = task.MaxIters
	out.Iterations = len(task.Iterations)
	out.Converged = task.Status == "success"

	// First-pass rate: fraction of criteria passing on iteration 1
	if len(task.Iterations) > 0 {
		first := task.Iterations[0]
		passing := 0
		for _, r := range first.Results {
			if r.Passes {
				passing++
			}
		}
		if len(first.Results) > 0 {
			out.FirstPassRate = float64(passing) / float64(len(first.Results))
		}
	}
	return out
}

// AddRoadmapItem adds an item to the roadmap from the HTTP API.
// Called when Emily receives a new task from humans or from her own planning.
func (ac *AutonomousCycle) AddRoadmapItem(item RoadmapItem) error {
	state, err := ac.loadState()
	if err != nil {
		return err
	}
	// Assign next priority if not set
	if item.Priority == 0 {
		max := 0
		for _, existing := range state.Roadmap {
			if existing.Priority > max {
				max = existing.Priority
			}
		}
		item.Priority = max + 1
	}
	if item.Status == "" {
		item.Status = "queued"
	}
	state.Roadmap = append(state.Roadmap, item)
	return ac.saveState(state)
}

// --- HTTP handlers on Server ---

// handleRoadmap: GET returns current roadmap; POST adds an item.
func (s *Server) handleRoadmap(w http.ResponseWriter, r *http.Request) {
	cronCfg := defaultCronConfig()
	cycle := NewAutonomousCycle(cronCfg, s.pipeline)

	switch r.Method {
	case http.MethodGet:
		state, err := cycle.loadState()
		if err != nil {
			http.Error(w, `{"error":"state load failed"}`, 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(state.Roadmap)

	case http.MethodPost:
		var item RoadmapItem
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, `{"error":"invalid json"}`, 400)
			return
		}
		if strings.TrimSpace(item.ID) == "" || strings.TrimSpace(item.Name) == "" {
			http.Error(w, `{"error":"id and name required"}`, 400)
			return
		}
		if err := cycle.AddRoadmapItem(item); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "added", "id": item.ID})

	default:
		http.Error(w, "method not allowed", 405)
	}
}

// handleStatus: GET returns the current dashboard / cycle state.
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", 405)
		return
	}
	cronCfg := defaultCronConfig()
	cycle := NewAutonomousCycle(cronCfg, s.pipeline)

	// Try dashboard.txt first (human-readable)
	if r.Header.Get("Accept") == "text/plain" {
		data, err := os.ReadFile(filepath.Join(cronCfg.StateDir, "dashboard.txt"))
		if err != nil {
			http.Error(w, "no status yet", 404)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write(data)
		return
	}

	// Default: JSON state
	state, err := cycle.loadState()
	if err != nil {
		http.Error(w, `{"error":"state load failed"}`, 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}
