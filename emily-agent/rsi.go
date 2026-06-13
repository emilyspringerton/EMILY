// emily-agent/rsi.go
// Recursive self-improvement loop engine.
//
// Pattern: given a task with measurable acceptance criteria, run Emily in a loop
// until all criteria pass or max_iters is exhausted.
//
// Emily both generates artifacts and evaluates them — the loop is self-referential.
// The first useful application: POST /rsi with a task that describes the RSI engine
// itself, and let Emily improve her own improvement process.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// AcceptanceCriterion defines one measurable requirement for a task.
type AcceptanceCriterion struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Target      string `json:"target"` // human-readable target, e.g. "> 95%" or "all pass"
}

// CriteriaResult is the evaluated state of one criterion after an iteration.
type CriteriaResult struct {
	Name   string `json:"name"`
	Passes bool   `json:"passes"`
	Value  string `json:"value"` // what was measured
	Gap    string `json:"gap"`   // distance to target, empty if passing
}

// IterationRecord logs one complete pass through the improvement loop.
type IterationRecord struct {
	Number     int              `json:"number"`
	Timestamp  time.Time        `json:"timestamp"`
	Artifact   string           `json:"artifact"` // the generated output
	Results    []CriteriaResult `json:"results"`
	AllPass    bool             `json:"all_pass"`
	Analysis   string           `json:"analysis"`   // why things passed/failed
	NextFocus  string           `json:"next_focus"` // what iteration n+1 should try
	TokensUsed int              `json:"tokens_used"`
}

// ImprovementTask is a complete RSI task specification and its execution state.
type ImprovementTask struct {
	ID          string                `json:"id"`
	Description string                `json:"description"`
	Criteria    []AcceptanceCriterion `json:"criteria"`
	MaxIters    int                   `json:"max_iters"`
	Iterations  []IterationRecord     `json:"iterations"`
	Status      string                `json:"status"` // pending|running|success|partial|failed
	StartedAt   time.Time             `json:"started_at"`
	CompletedAt *time.Time            `json:"completed_at,omitempty"`
	Lessons     []string              `json:"lessons,omitempty"`
}

// RSILoop runs iterative improvement loops backed by an LLM pipeline.
type RSILoop struct {
	pipeline *Pipeline
	stateDir string
}

// NewRSILoop creates an RSI engine backed by the default Emily state directory.
func NewRSILoop(p *Pipeline) *RSILoop {
	return NewRSILoopWithStateDir(p, filepath.Join(defaultCronConfig().StateDir, "rsi-tasks"))
}

// NewRSILoopWithStateDir creates an RSI engine that persists task JSON snapshots
// under stateDir. Tests and embedders can pass a temp directory to keep runs
// isolated from the process working directory.
func NewRSILoopWithStateDir(p *Pipeline, stateDir string) *RSILoop {
	return &RSILoop{pipeline: p, stateDir: stateDir}
}

// Run executes the improvement loop, updating task in place.
func (r *RSILoop) Run(ctx context.Context, task *ImprovementTask) error {
	if task == nil {
		return fmt.Errorf("rsi run: nil task")
	}
	task.Status = "running"
	task.StartedAt = time.Now()
	if err := r.saveTask(task); err != nil {
		return err
	}

	for i := 0; i < task.MaxIters; i++ {
		rec, err := r.runIteration(ctx, task)
		if err != nil {
			task.Status = "failed"
			t := time.Now()
			task.CompletedAt = &t
			_ = r.saveTask(task)
			return fmt.Errorf("iteration %d: %w", i+1, err)
		}
		task.Iterations = append(task.Iterations, rec)
		if err := r.saveTask(task); err != nil {
			task.Status = "failed"
			t := time.Now()
			task.CompletedAt = &t
			return err
		}

		if rec.AllPass {
			task.Status = "success"
			t := time.Now()
			task.CompletedAt = &t
			task.Lessons = r.extractLessons(ctx, task)
			return r.saveTask(task)
		}
	}

	task.Status = "partial"
	t := time.Now()
	task.CompletedAt = &t
	return r.saveTask(task)
}

func (r *RSILoop) saveTask(task *ImprovementTask) error {
	if strings.TrimSpace(r.stateDir) == "" {
		return nil
	}
	if task == nil {
		return fmt.Errorf("rsi save task: nil task")
	}
	if strings.TrimSpace(task.ID) == "" {
		return fmt.Errorf("rsi save task: missing task id")
	}
	if err := os.MkdirAll(r.stateDir, 0o755); err != nil {
		return fmt.Errorf("rsi state dir: %w", err)
	}
	data, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		return fmt.Errorf("rsi marshal task: %w", err)
	}
	path := filepath.Join(r.stateDir, safeTaskID(task.ID)+".json")
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("rsi write task: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("rsi commit task: %w", err)
	}
	return nil
}

func safeTaskID(id string) string {
	id = strings.TrimSpace(id)
	var b strings.Builder
	for _, r := range id {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '_', r == '.':
			b.WriteRune(r)
		default:
			b.WriteByte('-')
		}
	}
	if b.Len() == 0 {
		return "task"
	}
	return b.String()
}

func (r *RSILoop) runIteration(ctx context.Context, task *ImprovementTask) (IterationRecord, error) {
	n := len(task.Iterations) + 1
	rec := IterationRecord{Number: n, Timestamp: time.Now()}

	// Generate: call Emily to produce an artifact
	genPrompt := r.buildGenerationPrompt(task, n)
	genResp, err := r.pipeline.Generator.Complete(ctx, LLMRequest{
		Model:     r.pipeline.GenModel,
		Messages:  []Message{{Role: "system", Content: rsiGeneratorPrompt}, {Role: "user", Content: genPrompt}},
		MaxTokens: 4096,
	})
	if err != nil {
		return rec, fmt.Errorf("generation: %w", err)
	}
	rec.Artifact = genResp.Content
	rec.TokensUsed += genResp.InputTokens + genResp.OutputTokens

	// Evaluate: call Emily to measure artifact against criteria
	eval, evalTokens, err := r.evaluate(ctx, task, rec.Artifact)
	if err != nil {
		return rec, fmt.Errorf("evaluation: %w", err)
	}
	rec.Results = eval.Results
	rec.AllPass = eval.AllPass
	rec.Analysis = eval.Analysis
	rec.NextFocus = eval.NextFocus
	rec.TokensUsed += evalTokens

	return rec, nil
}

type evalOut struct {
	Results   []CriteriaResult `json:"results"`
	AllPass   bool             `json:"all_pass"`
	Analysis  string           `json:"analysis"`
	NextFocus string           `json:"next_focus"`
}

func (r *RSILoop) evaluate(ctx context.Context, task *ImprovementTask, artifact string) (evalOut, int, error) {
	criteriaJSON, _ := json.MarshalIndent(task.Criteria, "", "  ")

	prompt := fmt.Sprintf(
		"TASK: %s\n\nCRITERIA:\n%s\n\nARTIFACT:\n%s",
		task.Description, string(criteriaJSON), artifact,
	)

	resp, err := r.pipeline.Generator.Complete(ctx, LLMRequest{
		Model: r.pipeline.GenModel,
		Messages: []Message{
			{Role: "system", Content: rsiEvaluatorPrompt},
			{Role: "user", Content: prompt},
		},
		Temperature: 0.1,
		MaxTokens:   3072,
	})
	if err != nil {
		return evalOut{}, 0, err
	}

	tokens := resp.InputTokens + resp.OutputTokens

	// Find the first { in the response — the prompt says "first character must be {"
	// but LLMs sometimes emit preamble prose before the object. Slice from there.
	raw := resp.Content
	if start := strings.Index(raw, "{"); start > 0 {
		raw = raw[start:]
	}
	// Also strip any trailing ``` fence if the LLM wrapped the object anyway.
	raw = strings.TrimSpace(raw)
	if idx := strings.LastIndex(raw, "}"); idx >= 0 && idx < len(raw)-1 {
		raw = raw[:idx+1]
	}

	var out evalOut
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return evalOut{
			AllPass:   false,
			Analysis:  raw,
			NextFocus: "evaluator produced non-JSON output; check the artifact length and criteria count",
		}, tokens, nil
	}
	return out, tokens, nil
}

func (r *RSILoop) buildGenerationPrompt(task *ImprovementTask, iterNum int) string {
	criteriaJSON, _ := json.MarshalIndent(task.Criteria, "", "  ")

	var sb strings.Builder
	fmt.Fprintf(&sb, "TASK: %s\n\n", task.Description)
	fmt.Fprintf(&sb, "ACCEPTANCE CRITERIA (all must pass):\n%s\n\n", string(criteriaJSON))

	if len(task.Iterations) == 0 {
		sb.WriteString("ITERATION 1: Produce your first attempt.\n")
		sb.WriteString("Strategy: correctness first. Minimal, direct, complete. Make every criterion's compliance visible in the artifact — the evaluator reads but cannot execute, so demonstrate rather than imply.\n")
	} else {
		fmt.Fprintf(&sb, "ITERATION %d: Targeted fix only.\n\n", iterNum)

		// Build a criteria target lookup so the failure report shows target, not gap twice.
		targetFor := make(map[string]string, len(task.Criteria))
		for _, c := range task.Criteria {
			targetFor[c.Name] = c.Target
		}

		// Limit history to the most recent 3 iterations to cap prompt size on long tasks.
		// The last iteration's NextFocus is the primary signal; older records have
		// diminishing returns and grow the input by ~150 tokens each.
		const maxHistoryIters = 3
		history := task.Iterations
		if len(history) > maxHistoryIters {
			omitted := len(history) - maxHistoryIters
			fmt.Fprintf(&sb, "[%d earlier iteration(s) omitted]\n\n", omitted)
			history = history[len(history)-maxHistoryIters:]
		}

		for _, prev := range history {
			var passing, failing []string
			for _, cr := range prev.Results {
				if cr.Passes {
					passing = append(passing, cr.Name)
				} else {
					// Truncate the measured value to keep prompt size bounded across
					// iterations — the gap field (the fix) is more useful than the
					// full quoted artifact span.
					measured := cr.Value
					if len(measured) > 120 {
						measured = measured[:120] + "…"
					}
					failing = append(failing,
						fmt.Sprintf("  FAIL %-20s  measured: %q  target: %q  gap: %s",
							cr.Name, measured, targetFor[cr.Name], cr.Gap))
				}
			}
			fmt.Fprintf(&sb, "Iteration %d — passing: [%s]\n", prev.Number, strings.Join(passing, ", "))
			if len(failing) > 0 {
				sb.WriteString("            failing:\n")
				for _, f := range failing {
					sb.WriteString(f + "\n")
				}
			}
			fmt.Fprintf(&sb, "Root cause: %s\n", prev.Analysis)
			fmt.Fprintf(&sb, "Next focus: %s\n\n", prev.NextFocus)
		}

		last := task.Iterations[len(task.Iterations)-1]
		fmt.Fprintf(&sb, "YOUR FOCUS THIS ITERATION: %s\n\n", last.NextFocus)
		sb.WriteString("CONSTRAINT: change only what is required to address the focus above. Every passing criterion must remain passing. Adding unrequested scope causes failure.\n")
	}

	return sb.String()
}

func (r *RSILoop) extractLessons(ctx context.Context, task *ImprovementTask) []string {
	// Single-pass successes have no multi-iteration dynamics worth recording.
	if len(task.Iterations) <= 1 {
		return nil
	}
	// Strip the artifact field before marshaling: artifacts can be 1–4 kT each and
	// are irrelevant to the lessons prompt. Pass only the evaluation results.
	type lessonIter struct {
		Number    int              `json:"number"`
		Results   []CriteriaResult `json:"results"`
		AllPass   bool             `json:"all_pass"`
		Analysis  string           `json:"analysis"`
		NextFocus string           `json:"next_focus"`
	}
	stripped := make([]lessonIter, len(task.Iterations))
	for i, it := range task.Iterations {
		stripped[i] = lessonIter{
			Number:    it.Number,
			Results:   it.Results,
			AllPass:   it.AllPass,
			Analysis:  it.Analysis,
			NextFocus: it.NextFocus,
		}
	}
	itersJSON, _ := json.MarshalIndent(stripped, "", "  ")
	prompt := fmt.Sprintf(
		"Task: %s\n\nIteration history:\n%s\n\nExtract 3-5 concise lessons as a JSON array of strings. No markdown, no prose outside the array.",
		task.Description, string(itersJSON),
	)

	resp, err := r.pipeline.Generator.Complete(ctx, LLMRequest{
		Model:       r.pipeline.GenModel,
		Messages:    []Message{{Role: "user", Content: prompt}},
		Temperature: 0.2,
		MaxTokens:   512,
	})
	if err != nil {
		return nil
	}

	cleaned := strings.TrimSpace(resp.Content)
	cleaned = strings.TrimPrefix(strings.TrimPrefix(cleaned, "```json"), "```")
	cleaned = strings.TrimSuffix(cleaned, "```")

	var lessons []string
	_ = json.Unmarshal([]byte(cleaned), &lessons)
	return lessons
}

// --- System prompts ---

const rsiGeneratorPrompt = `You are the generation engine in an iterative improvement loop.

OUTPUT RULE: your entire response is the artifact. Your first character must be the first character of the artifact. No preamble. No postamble. No "Here is...", no "I'll...", no markdown fences around the artifact itself. Violating this causes pipeline failure.

ITERATION 1 — first attempt:
Produce the most direct, correct, self-contained artifact possible. Do not over-engineer. The evaluator reads but cannot execute — make compliance with each criterion visible from reading alone: use clear naming, inline comments at non-obvious points, and short inline examples where runtime behavior is a criterion.

SUBSEQUENT ITERATIONS — targeted fix:
You will receive a specific focus from the previous evaluation. Your job is surgical: implement exactly that change, nothing more. Every section not mentioned in the focus carries over unchanged. Rewriting passing sections is the primary cause of regression and wasted iterations.

INVARIANTS (override everything else):
- Correctness is non-negotiable. A faster but wrong artifact fails.
- If a criterion is genuinely ambiguous, resolve it conservatively and add a one-line comment at the relevant site explaining your interpretation.`

const rsiEvaluatorPrompt = `You are an objective evaluator in an iterative improvement loop.

CRITICAL OUTPUT RULE: your response must be exactly one JSON object. Your first character must be {. No markdown fences, no prose, no explanation outside the object. Any character before { causes a parse failure that halts the loop.

Required schema:
{
  "results": [
    {
      "name": "<criterion name copied exactly from the criteria list>",
      "passes": <true|false>,
      "value": "<quote the exact span of the artifact being evaluated — for code: the relevant line(s); for prose: the sentence. Never assert without citing evidence from the artifact.>",
      "gap": "<passes=false only: one implementable change that would make this criterion pass — name the specific function, field, or line and state the exact required change. 'Not implemented' is not acceptable. passes=true: \"\">"
    }
  ],
  "all_pass": <derive this from the results array: true ONLY when zero results have passes=false. Do not assert independently.>,
  "analysis": "<root cause of failures in one paragraph — WHY the deficiency exists: trace to the specific decision or omission in the artifact, not the symptom>",
  "next_focus": "<all_pass=false: exactly one criterion name + one specific change + one sentence explaining why fixing it unblocks the most other failures. No lists. No multi-step instructions. all_pass=true: \"\">"
}

BEFORE WRITING all_pass: count the results where passes=false. If that count is greater than zero, all_pass must be false.

EVALUATION DISCIPLINE:
- Every result, passing or failing, must cite artifact evidence in the value field.
- Pass only for behavior demonstrably present in the artifact, not implied or intended.
- Fail only for what the criterion actually requires — do not penalize out-of-scope issues.
- For subjective criteria (readability, naming clarity): cite one specific instance and explain why it meets or fails the standard.`

// --- HTTP handler ---

func (s *Server) handleRSI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	if !s.rl.Allow(s.ip(r)) {
		http.Error(w, `{"error":"rate limit exceeded"}`, 429)
		return
	}

	var req struct {
		Description string                `json:"description"`
		Criteria    []AcceptanceCriterion `json:"criteria"`
		MaxIters    int                   `json:"max_iters"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, 400)
		return
	}
	if strings.TrimSpace(req.Description) == "" || len(req.Criteria) == 0 {
		http.Error(w, `{"error":"description and criteria required"}`, 400)
		return
	}
	if req.MaxIters <= 0 {
		req.MaxIters = 10
	}
	if req.MaxIters > 20 {
		req.MaxIters = 20 // hard cap
	}

	task := &ImprovementTask{
		ID:          fmt.Sprintf("rsi-%d", time.Now().UnixNano()),
		Description: req.Description,
		Criteria:    req.Criteria,
		MaxIters:    req.MaxIters,
		Status:      "pending",
	}

	// 10 minutes max — RSI tasks can take a while
	ctx, cancel := context.WithTimeout(r.Context(), 600*time.Second)
	defer cancel()

	loop := NewRSILoop(s.pipeline)
	if err := loop.Run(ctx, task); err != nil {
		w.Header().Set("Content-Type", "application/json")
		errJSON, _ := json.Marshal(map[string]string{"error": err.Error(), "task_id": task.ID, "status": task.Status})
		w.WriteHeader(500)
		w.Write(errJSON)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}
