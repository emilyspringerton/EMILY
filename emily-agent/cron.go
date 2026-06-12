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
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"emily-agent/pkg/fcm"
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
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Priority    int                   `json:"priority"` // lower = higher priority
	Status      string                `json:"status"`   // queued|in_progress|complete|blocked
	Criteria    []AcceptanceCriterion `json:"criteria"`
	MaxIters    int                   `json:"max_iters"`
	Notes       string                `json:"notes,omitempty"`
}

// CycleState is persisted between cycles so Emily remembers her context.
type CycleState struct {
	Version        int               `json:"version"`
	CycleNumber    int               `json:"cycle_number"`
	LastCycleAt    time.Time         `json:"last_cycle_at"`
	ActiveTaskID   string            `json:"active_task_id,omitempty"`
	ActiveTask     *ImprovementTask  `json:"active_task,omitempty"`
	Roadmap        []RoadmapItem     `json:"roadmap"`
	CompletedTasks []ImprovementTask `json:"completed_tasks,omitempty"`
	Metrics        CycleMetrics      `json:"metrics"`
	NextCyclePlan  string            `json:"next_cycle_plan,omitempty"`
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
	gmail       *GmailClient      // optional; enables CEO escalation alerts from triage
	iduna       *IdunaClient      // optional; submits Apples to IDUNA after each cycle
	fcmSender   *fcm.Sender       // optional; fires push notifications to MJOLNIR on CEO-visible escalations
}

// NewAutonomousCycle creates a cycle runner.
func NewAutonomousCycle(cfg CronConfig, p *Pipeline, gmail *GmailClient) *AutonomousCycle {
	ac := &AutonomousCycle{
		cfg:         cfg,
		pipeline:    p,
		integration: buildIntegrationStore(),
		emiree:      NewEmireeAgent(cfg.StateDir),
		gmail:       gmail,
		iduna:       NewIdunaClientFromEnv(),
	}
	if fcm.IsConfigured() {
		sender, err := fcm.NewFromEnv()
		if err != nil {
			log.Printf("fcm: init failed (push notifications disabled): %v", err)
		} else {
			ac.fcmSender = sender
			log.Printf("fcm: sender initialized (project=%s)", os.Getenv("FCM_PROJECT_ID"))
		}
	}
	return ac
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
	sPath := streamPath(ac.cfg.StateDir)
	appendStream(sPath, StreamEntry{
		Type:  "cycle_start",
		Cycle: state.CycleNumber,
		Msg:   fmt.Sprintf("active_task=%s roadmap_items=%d", state.ActiveTaskID, len(state.Roadmap)),
	})

	// PHASE 2: DECIDE — pick what to work on this cycle (gear-aware)
	rec.Phase = PhaseDecide
	influence := ac.emiree.State.Influence()
	log.Printf("[cycle %d] emiree: %s", state.CycleNumber, ac.emiree.State.Summary())
	task, reason := ac.pickTask(state)
	// Apply gear influence to task max_iters without shortening an
	// already-running task below the work it has accumulated. Emiree may shift
	// into a more conservative gear between cycles, but that should not force a
	// task into partial completion simply because it already used more iterations
	// than the new gear would have allowed from scratch.
	if task != nil && influence.MaxIters > 0 {
		minAllowed := len(task.Iterations) + 1
		if influence.MaxIters > minAllowed {
			task.MaxIters = influence.MaxIters
		} else if task.MaxIters < minAllowed {
			task.MaxIters = minAllowed
		}
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
			if ac.iduna != nil {
				failCtx, fcancel := context.WithTimeout(context.Background(), 8*time.Second)
				defer fcancel()
				_, _ = ac.iduna.PostApple(failCtx, ApplePayload{
					AppleType:  "escalation",
					SourceRepo: "EMILY",
					Title:      fmt.Sprintf("cycle %d act error: %v", state.CycleNumber, err),
					Body:       fmt.Sprintf("task_id: %s\nerror: %v\nconsec_failures: %d", task.ID, err, state.Metrics.ConsecFailures),
				})
			}
		} else {
			state.Metrics.ConsecFailures = 0
			state.Metrics.ItersRun++
			rec.Outcome = result
			log.Printf("[cycle %d] act: %s", state.CycleNumber, result)
			switch task.Status {
			case "success":
				state.Metrics.TasksCompleted++
				state.CompletedTasks = append(state.CompletedTasks, *task)
				state.ActiveTaskID = ""
				state.ActiveTask = nil
				if strings.HasPrefix(task.ID, "heimdal-") {
					go ac.notifyHeimdalStatus(task, "complete")
				}
			case "partial", "failed":
				// Terminal-but-not-successful tasks must not pin the active slot forever;
				// the roadmap item is marked blocked in runIteration for operator review.
				state.ActiveTaskID = ""
				state.ActiveTask = nil
				if strings.HasPrefix(task.ID, "heimdal-") {
					go ac.notifyHeimdalStatus(task, "blocked")
				}
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
	if rec.Error == "" {
		state.Metrics.SuccessfulCycles++
	}
	rec.EndedAt = time.Now()

	triageFindings := 0
	if ac.integration != nil {
		if triageResult, triageErr := ac.runPrimeTriageCycle(ctx, ac.integration); triageErr != nil {
			log.Printf("[cycle %d] triage warn: %v", state.CycleNumber, triageErr)
			if ac.iduna != nil {
				warnCtx, warnCancel := context.WithTimeout(context.Background(), 8*time.Second)
				defer warnCancel()
				_, _ = ac.iduna.PostApple(warnCtx, ApplePayload{
					AppleType:  "observation",
					SourceRepo: "EMILY",
					Title:      fmt.Sprintf("cycle %d prime-triage warning", state.CycleNumber),
					Body:       fmt.Sprintf("error: %v", triageErr),
				})
			}
		} else {
			log.Printf("[cycle %d] triage: %s", state.CycleNumber, triageResult)
			// Count how many tasks were issued to feed back into Emiree
			if strings.Contains(triageResult, "tasks_written=") {
				fmt.Sscanf(triageResult[strings.Index(triageResult, "tasks_written="):], "tasks_written=%d", &triageFindings)
			}
		}
	}

	// Vision cycle: drain pending camera_observations from IDUNA.
	if ac.iduna != nil {
		visionCtx, visionCancel := context.WithTimeout(context.Background(), 60*time.Second)
		go func() {
			defer visionCancel()
			result := runVisionCycle(visionCtx, ac.iduna)
			log.Printf("[cycle %d] %s", state.CycleNumber, result)
		}()
	}

	// HEIMDAL sprint cycle: translate pending sprint requirements into RSI roadmap items.
	if ac.iduna != nil {
		heimdalCtx, heimdalCancel := context.WithTimeout(context.Background(), 60*time.Second)
		go func() {
			defer heimdalCancel()
			var push PushFunc
			if ac.fcmSender != nil && ac.iduna != nil {
				sender := ac.fcmSender
				iduna := ac.iduna
				push = func(title, body string, data map[string]string) {
					pushCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					defer cancel()
					deviceToken, err := iduna.GetPushToken(pushCtx, "mjolnir-emily")
					if err != nil || deviceToken == "" {
						return
					}
					_ = sender.Send(pushCtx, deviceToken, fcm.Message{
						Title:    title,
						Body:     body,
						Priority: "high",
						Data:     data,
					})
				}
			}
			result := ac.runHeimdalCycle(heimdalCtx, push)
			log.Printf("[cycle %d] %s", state.CycleNumber, result)
		}()
	}

	// Morning briefing: 09:00 UTC ±30 min, once per calendar day.
	if briefingDue(ac.cfg.StateDir) {
		var push PushFunc
		if ac.fcmSender != nil && ac.iduna != nil {
			sender := ac.fcmSender
			iduna := ac.iduna
			push = func(title, body string, data map[string]string) {
				pushCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				deviceToken, err := iduna.GetPushToken(pushCtx, "mjolnir-emily")
				if err != nil || deviceToken == "" {
					log.Printf("briefing: no device token (push skipped): %v", err)
					return
				}
				if err := sender.Send(pushCtx, deviceToken, fcm.Message{
					Title:    title,
					Body:     body,
					Priority: "normal",
					Data:     data,
				}); err != nil {
					log.Printf("briefing: send failed: %v", err)
					// File an escalation Apple so the failure is visible in the audit trail.
					if iduna != nil {
						failCtx, fcancel := context.WithTimeout(context.Background(), 8*time.Second)
						defer fcancel()
						_, _ = iduna.PostApple(failCtx, ApplePayload{
							AppleType:  "escalation",
							SourceRepo: "EMILY",
							Title:      "morning briefing FCM send failed",
							Body:       fmt.Sprintf("error: %v\ntitle: %s", err, title),
						})
					}
				}
			}
		}
		runMorningBriefing(ctx, ac.iduna, push, ac.cfg.StateDir)
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

	// Capture task completion state before goroutine (task pointer may change).
	taskCompleted := task != nil && task.Status == "success"
	var taskTitle string
	if task != nil {
		taskTitle = task.Description
		if taskTitle == "" {
			taskTitle = task.ID
		}
		if len(taskTitle) > 120 {
			taskTitle = taskTitle[:120] + "…"
		}
	}

	// Submit Apple to IDUNA — best-effort, never blocks or fails the cycle.
	if ac.iduna != nil {
		appleCtx, appleCancel := context.WithTimeout(context.Background(), 15*time.Second)
		go func() {
			defer appleCancel()
			payload := buildCycleApple(state, task, rec, triageFindings)
			id, err := ac.iduna.PostApple(appleCtx, payload)
			if err != nil {
				log.Printf("[cycle %d] apple: submit failed (non-fatal): %v", state.CycleNumber, err)
			} else {
				log.Printf("[cycle %d] apple: submitted id=%d type=%s", state.CycleNumber, id, payload.AppleType)
				appendStream(sPath, StreamEntry{
					Type:  "apple_filed",
					Cycle: state.CycleNumber,
					Msg:   payload.Title,
					Data:  map[string]any{"id": id, "apple_type": payload.AppleType},
				})
				// RSI cycle completion push: fire on task success only (not idle/audit cycles).
				if taskCompleted && ac.fcmSender != nil {
					pushCtx, pushCancel := context.WithTimeout(context.Background(), 8*time.Second)
					defer pushCancel()
					deviceToken, tokenErr := ac.iduna.GetPushToken(pushCtx, "mjolnir-emily")
					if tokenErr == nil && deviceToken != "" {
						_ = ac.fcmSender.Send(pushCtx, deviceToken, fcm.Message{
							Title:    fmt.Sprintf("RSI cycle %d complete", state.CycleNumber),
							Body:     taskTitle,
							Priority: "normal",
							Data: map[string]string{
								"apple_id":   fmt.Sprintf("%d", id),
								"cycle":      fmt.Sprintf("%d", state.CycleNumber),
								"tasks_done": fmt.Sprintf("%d", state.Metrics.TasksCompleted),
							},
						})
					}
				}
			}
		}()
	}

	duration := time.Since(rec.StartedAt).Round(time.Second)
	log.Printf("[cycle %d] complete in %s", state.CycleNumber, duration)
	appendStream(sPath, StreamEntry{
		Type:  "cycle_end",
		Cycle: state.CycleNumber,
		Msg:   rec.Outcome,
		Data: map[string]any{
			"duration_s":   duration.Seconds(),
			"phase":        string(rec.Phase),
			"error":        rec.Error,
			"total_cycles": state.Metrics.TotalCycles,
			"tasks_done":   state.Metrics.TasksCompleted,
		},
	})
	return nil
}

// RunDaemon runs cycles at approximately the configured interval.
// A ±60s jitter is added to avoid synchronising with external schedules
// and to reduce token waste from perfectly-timed back-to-back runs.
func (ac *AutonomousCycle) RunDaemon() {
	log.Printf("Emily daemon starting — base interval %s (±60s jitter)", ac.cfg.Interval)
	for {
		start := time.Now()
		if err := ac.RunOnce(); err != nil {
			log.Printf("cycle error: %v", err)
		}
		elapsed := time.Since(start)
		jitter := time.Duration(rand.Int63n(120)-60) * time.Second // ±60s
		sleep := ac.cfg.Interval + jitter - elapsed
		if sleep < 30*time.Second {
			sleep = 30 * time.Second // floor: never spin faster than 30s
		}
		log.Printf("daemon sleeping %s before next cycle", sleep.Round(time.Second))
		time.Sleep(sleep)
	}
}

// pickTask selects what to work on this cycle.
// Priority: resume active task > highest-priority executable roadmap item.
func (ac *AutonomousCycle) pickTask(state *CycleState) (*ImprovementTask, string) {
	// Resume any non-terminal active task. Older state files may persist an
	// ActiveTask with status pending before the first iteration, while current
	// cycles persist it as running between iterations.
	if state.ActiveTask != nil {
		switch state.ActiveTask.Status {
		case "pending", "running":
			return state.ActiveTask, fmt.Sprintf("resume active task: %s", state.ActiveTask.ID)
		case "success", "partial", "failed":
			state.ActiveTaskID = ""
			state.ActiveTask = nil
		}
	}

	// Find the highest-priority executable roadmap item. The canonical roadmap
	// ships some always-on initiatives as in_progress even before an ActiveTask is
	// materialized; treating those as executable lets a fresh state begin with the
	// actual highest-priority work instead of skipping to later queued items.
	var best *RoadmapItem
	for i := range state.Roadmap {
		item := &state.Roadmap[i]
		if !isExecutableRoadmapStatus(item.Status) {
			continue
		}
		if best == nil || item.Priority < best.Priority {
			best = item
		}
	}
	if best == nil {
		return nil, "no executable tasks"
	}

	// Promote roadmap item to active task.
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

func isExecutableRoadmapStatus(status string) bool {
	switch status {
	case "", "queued", "in_progress":
		return true
	default:
		return false
	}
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
		for i := range state.Roadmap {
			if state.Roadmap[i].ID == task.ID {
				state.Roadmap[i].Status = "blocked"
				state.Roadmap[i].Notes = fmt.Sprintf("Reached max_iters (%d) without satisfying all criteria; last pass count %d/%d. Requires review or a revised task.", task.MaxIters, passCount(rec.Results), len(rec.Results))
			}
		}
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
// roadmap item is active. It reads FatBaby observations, triages them, issues
// directed tasks for high-relevance findings, and sends Gmail alerts for
// CEO-visible observations (when gmail credentials are configured).
func (ac *AutonomousCycle) runPrimeTriageCycle(ctx context.Context, store *IntegrationStore) (string, error) {
	var push PushFunc
	if ac.fcmSender != nil && ac.iduna != nil {
		sender := ac.fcmSender
		iduna := ac.iduna
		const mjolnirAgent = "mjolnir-emily"
		push = func(title, body string, data map[string]string) {
			pushCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			deviceToken, err := iduna.GetPushToken(pushCtx, mjolnirAgent)
			if err != nil || deviceToken == "" {
				log.Printf("fcm: no device token for %s (push skipped): %v", mjolnirAgent, err)
				return
			}
			if err := sender.Send(pushCtx, deviceToken, fcm.Message{
				Title:    title,
				Body:     body,
				Priority: "high",
				Data:     data,
			}); err != nil {
				log.Printf("fcm: send failed (non-fatal): %v", err)
				failCtx, fcancel := context.WithTimeout(context.Background(), 8*time.Second)
				defer fcancel()
				_, _ = iduna.PostApple(failCtx, ApplePayload{
					AppleType:  "escalation",
					SourceRepo: "EMILY",
					Title:      "prime-triage FCM send failed",
					Body:       fmt.Sprintf("error: %v\ntoken: %s\ntitle: %s", err, deviceToken, title),
				})
			} else {
				log.Printf("fcm: push sent to %s", mjolnirAgent)
			}
		}
	}
	return runPrimeTriage(ctx, store, ac.pipeline, ac.gmail, push)
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
			Status:      "complete",
			Criteria: []AcceptanceCriterion{
				{Name: "throughput", Description: "collection rate from configured subreddits", Target: ">= 500 posts/hour sustained"},
				{Name: "quality_filter", Description: "posts below quality bar are excluded before write", Target: "score/comments/body filters applied at source"},
				{Name: "dedup", Description: "same post is never written twice across restarts", Target: "seenSet loaded from corpus at startup"},
				{Name: "resilience", Description: "rate limits and transient errors handled gracefully", Target: "exponential backoff, zero panics, fail-open"},
			},
			MaxIters: 5,
			Notes:    "Completed: httpGetWithRetry provides exponential backoff on 429/503; seenSet O(1) dedup; quality filters applied at fetchSubreddit; per-source stats added to CollectionStats.",
		},
		{
			ID:          "wikipedia-collector",
			Name:        "Wikipedia collection pipeline — tune and expand",
			Description: "The Go WikipediaSource in collector.go fetches recent changes + plain-text extracts. Improve it: add category filtering, expand to 100 pages/batch, add article quality signals (length, section count, citation count from wikitext), ensure stub articles are filtered.",
			Priority:    3,
			Status:      "complete",
			Criteria: []AcceptanceCriterion{
				{Name: "extraction", Description: "article text is clean plain text, no wikitext artifacts", Target: "extract field contains prose only"},
				{Name: "stub_filter", Description: "stub articles are excluded", Target: "length < 1000 chars filtered at source"},
				{Name: "throughput", Description: "articles collected per hour", Target: ">= 500 articles/hour"},
				{Name: "schema", Description: "CollectedDoc fields populated", Target: "id, title, body, body_bytes, quality_score, quality_tier all present"},
			},
			MaxIters: 5,
			Notes:    "Completed: explaintext gives clean prose; isWikipediaStub filter added; 100-page batch at 500ms inter-batch delay; countWikiSections adds sections:N tag; all schema fields present.",
		},
		{
			ID:          "quality-scorer",
			Name:        "Quality scorer — calibrate and extend",
			Description: "The Go QualityScorer in collector.go is live (deterministic, no ML). Calibrate thresholds against real corpus samples: verify gold/silver/bronze distribution, tune spam regex patterns, add domain-specific signals for ML content (code blocks, equations, citations), verify false-positive rate < 5% on known-good samples.",
			Priority:    4,
			Status:      "complete",
			Criteria: []AcceptanceCriterion{
				{Name: "scoring", Description: "every record gets a numeric quality score 0-10", Target: "0.0-10.0 float, mean score > 5.0 for reddit/wiki corpus"},
				{Name: "tier_distribution", Description: "gold/silver/bronze distribution is reasonable", Target: "gold 10-20%, silver 30-50%, bronze remainder"},
				{Name: "spam_detection", Description: "obvious spam/ads are scored low", Target: "score < 4.0 for samples containing 3+ spam phrases"},
				{Name: "false_positive_rate", Description: "legitimate technical content is not penalized", Target: "< 5% of gold-quality human-labeled docs score < 5.0"},
			},
			MaxIters: 5,
			Notes:    "Completed: 0-10 float score; gold>=6.5/silver>=4.5; spam regex with 1.5x penalty; citation richness signal; technical bonus for code/equations/ML keywords (max score ~9.0).",
		},
		{
			ID:          "arxiv-collector",
			Name:        "ArXiv collection pipeline — live and collecting",
			Description: "ArXivSource is live in collector.go fetching cs.AI/cs.LG/cs.CL/cs.CV/stat.ML abstracts. Improve: add full-paper HTML extraction via arxiv HTML endpoint, add citation count heuristic to quality scorer, tune category list against corpus quality metrics.",
			Priority:    4,
			Status:      "complete",
			Criteria: []AcceptanceCriterion{
				{Name: "throughput", Description: "papers collected per hour", Target: ">= 200 papers/hour from default categories"},
				{Name: "body_quality", Description: "body field contains clean prose, no LaTeX artifacts", Target: "< 5% of collected docs contain raw LaTeX macros"},
				{Name: "dedup", Description: "same paper not written twice across restarts", Target: "seenSet deduplication verified on restart"},
			},
			MaxIters: 4,
			Notes:    "Completed: HTML full-paper extraction (MaxFullText=10 per batch); cleanAbstract strips raw LaTeX macros (\\cmd, \\begin/\\end environments, inline math) from abstracts; seenSet dedup verified; throughput met via batched Atom feed.",
		},
		{
			ID:          "bob-agent",
			Name:        "Bob database admin agent",
			Description: "Design and implement Bob: a minimal Go HTTP service that accepts structured task requests from Emily, executes MySQL operations, returns structured results, and maintains an audit log.",
			Priority:    6,
			Status:      "complete",
			Criteria: []AcceptanceCriterion{
				{Name: "api", Description: "Emily can POST a task to Bob and GET the result", Target: "POST /task, GET /task/{id}"},
				{Name: "isolation", Description: "Bob only executes database operations", Target: "no filesystem or network access beyond DB"},
				{Name: "audit", Description: "every operation is logged with timestamp, operation, result", Target: "100% operations logged"},
				{Name: "safety", Description: "destructive operations require explicit confirmation flag", Target: "drop/delete require {confirm: true}"},
			},
			MaxIters: 6,
			Notes:    "Completed in IDUNA repo (cmd/bob-agent). MySQL schema admin agent with safety gates for destructive operations.",
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

// buildCycleApple constructs the Apple payload for a completed RSI cycle.
// apple_type is "improvement" when a task succeeded, "observation" when triage
// found new tasks, "audit" for idle/monitoring cycles.
func buildCycleApple(state *CycleState, task *ImprovementTask, rec CycleRecord, triageFindings int) ApplePayload {
	appleType := "audit"
	title := fmt.Sprintf("Cycle %d — idle (no queued tasks)", state.CycleNumber)
	var bodyLines []string

	if task != nil {
		switch task.Status {
		case "success":
			appleType = "improvement"
			title = fmt.Sprintf("Cycle %d — %s: converged", state.CycleNumber, task.ID)
		case "running":
			appleType = "improvement"
			title = fmt.Sprintf("Cycle %d — %s: iteration %d", state.CycleNumber, task.ID, len(task.Iterations))
		default:
			appleType = "observation"
			title = fmt.Sprintf("Cycle %d — %s: %s", state.CycleNumber, task.ID, task.Status)
		}
		bodyLines = append(bodyLines,
			fmt.Sprintf("## Task: %s", task.ID),
			fmt.Sprintf("**Status:** %s  **Iterations:** %d/%d",
				task.Status, len(task.Iterations), task.MaxIters),
			"",
			"**Description:** "+task.Description,
		)
		if len(task.Iterations) > 0 {
			last := task.Iterations[len(task.Iterations)-1]
			bodyLines = append(bodyLines, "",
				"### Last iteration",
				fmt.Sprintf("- All criteria pass: %v", last.AllPass),
				"- Analysis: "+last.Analysis,
			)
		}
	} else if triageFindings > 0 {
		appleType = "observation"
		title = fmt.Sprintf("Cycle %d — triage: %d new tasks", state.CycleNumber, triageFindings)
	}

	bodyLines = append(bodyLines, "",
		"## Cycle summary",
		fmt.Sprintf("- Cycle: %d", state.CycleNumber),
		fmt.Sprintf("- Duration: %s", rec.EndedAt.Sub(rec.StartedAt).Round(time.Second)),
		fmt.Sprintf("- Outcome: %s", rec.Outcome),
		fmt.Sprintf("- Total completed tasks: %d", state.Metrics.TasksCompleted),
		fmt.Sprintf("- Triage findings this cycle: %d", triageFindings),
	)

	var cycleTokens int
	if task != nil {
		for _, iter := range task.Iterations {
			cycleTokens += iter.TokensUsed
		}
	}
	if cycleTokens > 0 {
		bodyLines = append(bodyLines, fmt.Sprintf("- Tokens this cycle: %d", cycleTokens))
	}

	meta, _ := json.Marshal(map[string]any{
		"cycle_number":         state.CycleNumber,
		"task_id":              func() string { if task != nil { return task.ID }; return "" }(),
		"triage_findings":      triageFindings,
		"consecutive_failures": state.Metrics.ConsecFailures,
		"duration_s":           int(rec.EndedAt.Sub(rec.StartedAt).Seconds()),
		"tokens_used":          cycleTokens,
	})

	return ApplePayload{
		SourceRepo: "emily",
		RunID:      fmt.Sprintf("emily-cycle-%d-%s", state.CycleNumber, rec.StartedAt.UTC().Format("20060102T150405Z")),
		AppleType:  appleType,
		Title:      title,
		Body:       strings.Join(bodyLines, "\n"),
		Metadata:   json.RawMessage(meta),
	}
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
	cycle := NewAutonomousCycle(cronCfg, s.pipeline, s.gmail)

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
	cycle := NewAutonomousCycle(cronCfg, s.pipeline, s.gmail)

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
