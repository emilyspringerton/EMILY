package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestConversationStoreWritesMarkdownIndexAndSearch(t *testing.T) {
	dir := t.TempDir()
	store, err := NewConversationStore(dir, false, false)
	if err != nil {
		t.Fatalf("NewConversationStore: %v", err)
	}

	turn := Turn{
		ID:        "turn-1",
		Timestamp: time.Date(2026, 5, 30, 14, 45, 0, 0, time.UTC),
		UserInput: "Design the Git workflow for Emily memory",
		Final:     "We will store every exchange as markdown in Git and generate a searchable index.",
		Validated: true,
		Model:     "test-model",
	}
	if err := store.AppendTurn("session-abc", turn); err != nil {
		t.Fatalf("AppendTurn: %v", err)
	}

	jsonl := filepath.Join(dir, "sessions", "session-abc.jsonl")
	if _, err := os.Stat(jsonl); err != nil {
		t.Fatalf("expected JSONL session file: %v", err)
	}

	markdown := filepath.Join(dir, "conversations", "2026", "05", "2026-05-30-14-45-design-the-git-workflow-for-emily-memory.md")
	contentBytes, err := os.ReadFile(markdown)
	if err != nil {
		t.Fatalf("expected markdown conversation: %v", err)
	}
	content := string(contentBytes)
	for _, want := range []string{"# Conversation: 2026-05-30", "## CEO [14:45]", "## Emily [14:45]", "## Decisions Made", "#git"} {
		if !strings.Contains(content, want) {
			t.Fatalf("markdown missing %q:\n%s", want, content)
		}
	}

	indexBytes, err := os.ReadFile(filepath.Join(dir, "conversations", "index.md"))
	if err != nil {
		t.Fatalf("expected index: %v", err)
	}
	if !strings.Contains(string(indexBytes), "Design the Git workflow for Emily memory") {
		t.Fatalf("index missing conversation title:\n%s", string(indexBytes))
	}

	results, err := store.Search("searchable index", 10)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 1 || results[0].Path != filepath.ToSlash(filepath.Join("conversations", "2026", "05", "2026-05-30-14-45-design-the-git-workflow-for-emily-memory.md")) {
		t.Fatalf("unexpected search results: %#v", results)
	}
}

type scriptedLLM struct {
	calls int
}

func (s *scriptedLLM) Complete(ctx context.Context, req LLMRequest) (LLMResponse, error) {
	s.calls++
	if len(req.Messages) > 0 && strings.Contains(req.Messages[0].Content, "objective evaluator") {
		return LLMResponse{
			Content: `{"results":[{"name":"done","passes":false,"value":"artifact lacks done marker","gap":"Add the done marker."}],"all_pass":false,"analysis":"The artifact omits the required done marker.","next_focus":"done: add the done marker so the only criterion can pass."}`,
		}, nil
	}
	return LLMResponse{Content: "draft artifact"}, nil
}

func (s *scriptedLLM) Name() string { return "scripted" }

func TestRSILoopPersistsTaskSnapshots(t *testing.T) {
	dir := t.TempDir()
	loop := NewRSILoopWithStateDir(&Pipeline{Generator: &scriptedLLM{}, GenModel: "test"}, dir)
	task := &ImprovementTask{
		ID:          "demo/task 1",
		Description: "Produce an artifact with a done marker.",
		Criteria:    []AcceptanceCriterion{{Name: "done", Description: "artifact includes a done marker", Target: "done marker present"}},
		MaxIters:    1,
		Status:      "pending",
	}

	if err := loop.Run(context.Background(), task); err != nil {
		t.Fatalf("Run: %v", err)
	}

	path := filepath.Join(dir, "demo-task-1.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected persisted task snapshot: %v", err)
	}
	var persisted ImprovementTask
	if err := json.Unmarshal(data, &persisted); err != nil {
		t.Fatalf("persisted task is not valid JSON: %v\n%s", err, string(data))
	}
	if persisted.ID != task.ID || persisted.Status != "partial" || len(persisted.Iterations) != 1 {
		t.Fatalf("unexpected persisted task: %#v", persisted)
	}
}

func TestPickTaskStartsHighestPriorityInProgressItem(t *testing.T) {
	cycle := &AutonomousCycle{}
	state := &CycleState{Roadmap: []RoadmapItem{
		{ID: "prime-triage", Name: "Prime triage", Priority: 0, Status: "in_progress", MaxIters: 3},
		{ID: "bob-agent", Name: "Bob", Priority: 6, Status: "queued", MaxIters: 6},
	}}

	task, reason := cycle.pickTask(state)
	if task == nil {
		t.Fatalf("expected task, got nil (%s)", reason)
	}
	if task.ID != "prime-triage" {
		t.Fatalf("expected in-progress prime-triage to be selected first, got %q", task.ID)
	}
	if state.ActiveTaskID != "prime-triage" || state.ActiveTask == nil {
		t.Fatalf("expected active task to be materialized, got id=%q task=%#v", state.ActiveTaskID, state.ActiveTask)
	}
}

func TestPickTaskClearsTerminalActiveTask(t *testing.T) {
	cycle := &AutonomousCycle{}
	state := &CycleState{
		ActiveTaskID: "old",
		ActiveTask:   &ImprovementTask{ID: "old", Status: "partial"},
		Roadmap: []RoadmapItem{
			{ID: "next", Name: "Next", Priority: 1, Status: "queued", MaxIters: 2},
		},
	}

	task, _ := cycle.pickTask(state)
	if task == nil || task.ID != "next" {
		t.Fatalf("expected next queued task after terminal active task, got %#v", task)
	}
	if state.ActiveTaskID != "next" {
		t.Fatalf("expected active task id to move to next, got %q", state.ActiveTaskID)
	}
}

func TestRunIterationMarksMaxIterPartialAsBlocked(t *testing.T) {
	cycle := &AutonomousCycle{pipeline: &Pipeline{Generator: &scriptedLLM{}, GenModel: "test"}}
	state := &CycleState{Roadmap: []RoadmapItem{
		{ID: "demo", Name: "Demo", Priority: 1, Status: "in_progress", MaxIters: 1},
	}}
	task := &ImprovementTask{
		ID:          "demo",
		Description: "Produce an artifact with a done marker.",
		Criteria:    []AcceptanceCriterion{{Name: "done", Description: "artifact includes a done marker", Target: "done marker present"}},
		MaxIters:    1,
		Status:      "pending",
	}

	result, err := cycle.runIteration(context.Background(), state, task)
	if err != nil {
		t.Fatalf("runIteration: %v", err)
	}
	if !strings.Contains(result, "status=partial") {
		t.Fatalf("expected partial result, got %q", result)
	}
	if task.Status != "partial" || task.CompletedAt == nil {
		t.Fatalf("expected task partial with completion time, got status=%q completed=%v", task.Status, task.CompletedAt)
	}
	if state.Roadmap[0].Status != "blocked" {
		t.Fatalf("expected roadmap item blocked, got %q", state.Roadmap[0].Status)
	}
	if !strings.Contains(state.Roadmap[0].Notes, "Reached max_iters (1)") {
		t.Fatalf("expected review note, got %q", state.Roadmap[0].Notes)
	}
}

// ---------------------------------------------------------------------------
// Wikipedia stub detection and section counting
// ---------------------------------------------------------------------------

func TestIsWikipediaStub_DetectsStubs(t *testing.T) {
	stubs := []string{
		"This article is a stub. You can help Wikipedia by expanding it.",
		"this is a stub — please expand",
		"This biography is a stub.",
	}
	for _, s := range stubs {
		if !isWikipediaStub(s) {
			t.Errorf("expected stub detection for: %q", s)
		}
	}
}

func TestIsWikipediaStub_PassesNonStubs(t *testing.T) {
	articles := []string{
		"Python is a high-level programming language. It was created by Guido van Rossum...",
		"Transformer models use self-attention to process sequences in parallel.",
		"The article about stubs in software engineering covers many cases.",
	}
	for _, s := range articles {
		if isWikipediaStub(s) {
			t.Errorf("false-positive stub detection for: %q", s)
		}
	}
}

func TestCountWikiSections_MultiSection(t *testing.T) {
	text := "Introduction paragraph.\n\nSection one content here.\n\nSection two content."
	n := countWikiSections(text)
	if n < 2 {
		t.Errorf("expected >= 2 sections, got %d", n)
	}
}

func TestCountWikiSections_Empty(t *testing.T) {
	if n := countWikiSections(""); n != 0 {
		t.Errorf("expected 0 for empty, got %d", n)
	}
}

// ---------------------------------------------------------------------------
// QualityScorer — technical content bonus
// ---------------------------------------------------------------------------

func TestQualityScorer_TechBonusForCodeBlock(t *testing.T) {
	qs := &QualityScorer{}
	withCode := CollectedDoc{
		ID:   "test",
		Body: "Here is an example:\n\n```python\ndef train(model, data):\n    optimizer.step()\n    return loss\n```\n\nThis shows a basic training loop using PyTorch with gradient descent and backpropagation through the neural network layers.",
	}
	without := CollectedDoc{
		ID:   "test2",
		Body: "Here is an example of how to do something with data. You can process it step by step in order to achieve the desired result.",
	}
	scoreWith, _ := qs.Score(withCode)
	scoreWithout, _ := qs.Score(without)
	if scoreWith <= scoreWithout {
		t.Errorf("expected code+tech-keywords to score higher: with=%.2f without=%.2f", scoreWith, scoreWithout)
	}
}

func TestQualityScorer_TechBonusForEquation(t *testing.T) {
	qs := &QualityScorer{}
	// Compare score of same-length body with and without equations.
	plain := CollectedDoc{
		ID:   "eq-plain",
		Body: "The loss function is defined. Gradient descent minimises it. This approach is fundamental. The model converges after training.",
	}
	withEq := CollectedDoc{
		ID:   "eq-test",
		Body: "The loss function is defined as $L = -\\sum y \\log p$. Gradient descent minimises $L$ using $\\theta \\leftarrow \\theta - \\alpha \\nabla L$. This approach is fundamental to neural network training.",
	}
	scorePlain, _ := qs.Score(plain)
	scoreEq, _ := qs.Score(withEq)
	// Equation-dense content should score higher than equivalent plain prose.
	if scoreEq <= scorePlain {
		t.Errorf("equation content (%.2f) should outscore plain text (%.2f)", scoreEq, scorePlain)
	}
}

func TestQualityScorer_NoFalsePositiveOnNormalText(t *testing.T) {
	qs := &QualityScorer{}
	// Plain prose without tech signals should not get the bonus.
	doc := CollectedDoc{
		ID:   "plain",
		Body: "The company announced strong quarterly results today. Revenue grew by fifteen percent year over year. The board approved a dividend increase for shareholders.",
	}
	score, _ := qs.Score(doc)
	// Should be normal score, not inflated by spurious tech detection.
	techBonus := score - func() float64 {
		// Re-score without any tech signals — approximate by checking score is < 6.5
		return 0
	}()
	_ = techBonus
	// Main check: clean prose doesn't get gold from tech bonus alone
	if score >= 6.5 {
		t.Errorf("plain prose scored %.2f (gold tier) unexpectedly; tech bonus may be too generous", score)
	}
}

// ---------------------------------------------------------------------------
// ArXiv abstract LaTeX cleanup
// ---------------------------------------------------------------------------

func TestCleanAbstract_StripsInlineCommands(t *testing.T) {
	raw := `We propose \textbf{FLAME}, a novel framework for \emph{efficient} training.`
	out := cleanAbstract(raw)
	if strings.Contains(out, `\textbf`) || strings.Contains(out, `\emph`) {
		t.Errorf("LaTeX commands not stripped: %s", out)
	}
	if !strings.Contains(out, "FLAME") || !strings.Contains(out, "efficient") {
		t.Errorf("content was removed: %s", out)
	}
}

func TestCleanAbstract_StripsMathAndEnvironments(t *testing.T) {
	raw := `The loss $\mathcal{L} = \sum_{i} \ell_i$ is minimised via SGD. See \begin{equation} x^2 + y^2 = r^2 \end{equation} for details.`
	out := cleanAbstract(raw)
	if strings.Contains(out, `\begin`) || strings.Contains(out, `\end`) {
		t.Errorf("LaTeX environment not removed: %s", out)
	}
	if strings.Contains(out, `\mathcal`) || strings.Contains(out, `\sum`) {
		t.Errorf("math block not cleaned: %s", out)
	}
	if !strings.Contains(out, "minimised via SGD") {
		t.Errorf("prose was removed: %s", out)
	}
}

func TestCleanAbstract_PassesCleanText(t *testing.T) {
	clean := "We present a new approach to large language model fine-tuning. Our method achieves state-of-the-art results."
	out := cleanAbstract(clean)
	if out != clean {
		t.Errorf("clean text was modified: got %q, want %q", out, clean)
	}
}

func TestCleanAbstract_StripsBareCommands(t *testing.T) {
	raw := `We achieve \alpha = 0.01 and \beta accuracy improvements.`
	out := cleanAbstract(raw)
	if strings.Contains(out, `\alpha`) || strings.Contains(out, `\beta`) {
		t.Errorf("bare LaTeX commands not stripped: %s", out)
	}
}

// ---------------------------------------------------------------------------
// Per-source stats
// ---------------------------------------------------------------------------

func TestCollectorPipelinePerSourceStats(t *testing.T) {
	dir := t.TempDir()
	store, err := NewDocStore(dir)
	if err != nil {
		t.Fatalf("NewDocStore: %v", err)
	}
	pipe := NewCollectorPipeline(CollectorConfig{
		Store:   store,
		Scorer:  &QualityScorer{},
		Workers: 1,
		MinTier: "bronze",
	})
	// Manually call processOne to increment per-source counter.
	pipe.seen = newDocSeenSet(nil)
	pipe.processOne(CollectedDoc{
		ID: "test-1", SourceName: "reddit",
		Body: strings.Repeat("word sentence. ", 100), // enough to be bronze+
	})
	pipe.processOne(CollectedDoc{
		ID: "test-2", SourceName: "arxiv",
		Body: strings.Repeat("abstract sentence. ", 100),
	})
	pipe.processOne(CollectedDoc{
		ID: "test-3", SourceName: "reddit",
		Body: strings.Repeat("another sentence. ", 100),
	})
	stats := pipe.Stats()
	if stats.PerSource["reddit"] != 2 {
		t.Errorf("reddit count = %d, want 2", stats.PerSource["reddit"])
	}
	if stats.PerSource["arxiv"] != 1 {
		t.Errorf("arxiv count = %d, want 1", stats.PerSource["arxiv"])
	}
}

// ---------------------------------------------------------------------------
// stripHTML
// ---------------------------------------------------------------------------

func TestStripHTML_RemovesTags(t *testing.T) {
	input := `<html><body><article><h1>Title</h1><p>First paragraph.</p><p>Second paragraph with <strong>emphasis</strong>.</p></article></body></html>`
	out := stripHTML(input)
	if strings.Contains(out, "<") || strings.Contains(out, ">") {
		t.Errorf("stripHTML left HTML tags in output: %s", out)
	}
	if !strings.Contains(out, "First paragraph") {
		t.Errorf("stripHTML removed text content: %s", out)
	}
	if !strings.Contains(out, "emphasis") {
		t.Errorf("stripHTML removed emphasised text: %s", out)
	}
}

func TestStripHTML_DecodesEntities(t *testing.T) {
	out := stripHTML("&amp; &lt;tag&gt; &quot;quoted&quot;")
	if strings.Contains(out, "&amp;") || strings.Contains(out, "&lt;") {
		t.Errorf("HTML entities not decoded: %s", out)
	}
	if !strings.Contains(out, "&") || !strings.Contains(out, "<") {
		t.Errorf("entity decoding incorrect: %s", out)
	}
}

func TestStripHTML_RemovesScriptAndStyle(t *testing.T) {
	input := `<p>Visible text</p><script>alert("xss")</script><style>.hidden{display:none}</style>`
	out := stripHTML(input)
	if strings.Contains(out, "alert") || strings.Contains(out, "display:none") {
		t.Errorf("script/style not stripped: %s", out)
	}
	if !strings.Contains(out, "Visible text") {
		t.Errorf("visible text was removed: %s", out)
	}
}

// ---------------------------------------------------------------------------
// IdunaClient — nil-safe when env vars are absent
// ---------------------------------------------------------------------------

func TestNewIdunaClientFromEnv_NilWhenUnconfigured(t *testing.T) {
	// Ensure env vars are not set in test environment.
	for _, k := range []string{"IDUNA_BASE_URL", "IDUNA_AGENT_NAME", "IDUNA_AGENT_SECRET"} {
		t.Setenv(k, "")
	}
	client := NewIdunaClientFromEnv()
	if client != nil {
		t.Error("expected nil IdunaClient when env vars are absent, got non-nil")
	}
}

func TestNewIdunaClientFromEnv_NonNilWhenConfigured(t *testing.T) {
	t.Setenv("IDUNA_BASE_URL", "http://localhost:8080")
	t.Setenv("IDUNA_AGENT_NAME", "TEST_AGENT")
	t.Setenv("IDUNA_AGENT_SECRET", "sk-test-secret")
	client := NewIdunaClientFromEnv()
	if client == nil {
		t.Error("expected non-nil IdunaClient when all env vars are set, got nil")
	}
}

// ---------------------------------------------------------------------------
// buildCycleApple
// ---------------------------------------------------------------------------

func TestBuildCycleApple_SuccessTask(t *testing.T) {
	now := time.Now().UTC()
	state := &CycleState{CycleNumber: 7, Metrics: CycleMetrics{TasksCompleted: 3}}
	task := &ImprovementTask{ID: "reddit-collector", Status: "success", MaxIters: 5,
		Iterations: []IterationRecord{{AllPass: true, Analysis: "All criteria met."}}}
	rec := CycleRecord{StartedAt: now, EndedAt: now.Add(45 * time.Second), Outcome: "converged"}
	payload := buildCycleApple(state, task, rec, 0)
	if payload.AppleType != "improvement" {
		t.Errorf("type = %q, want improvement for successful task", payload.AppleType)
	}
	if !strings.Contains(payload.Title, "cycle-7") && !strings.Contains(payload.Title, "Cycle 7") {
		t.Errorf("title %q should mention cycle number", payload.Title)
	}
	if payload.SourceRepo != "emily" {
		t.Errorf("source_repo = %q, want emily", payload.SourceRepo)
	}
	if !strings.Contains(payload.RunID, "emily-cycle-7") {
		t.Errorf("run_id = %q, want emily-cycle-7 prefix", payload.RunID)
	}
}

func TestBuildCycleApple_IdleCycle(t *testing.T) {
	state := &CycleState{CycleNumber: 1}
	rec := CycleRecord{StartedAt: time.Now(), EndedAt: time.Now(), Outcome: "idle"}
	payload := buildCycleApple(state, nil, rec, 0)
	if payload.AppleType != "audit" {
		t.Errorf("type = %q, want audit for idle cycle", payload.AppleType)
	}
}

func TestBuildCycleApple_TriageObservation(t *testing.T) {
	state := &CycleState{CycleNumber: 2}
	rec := CycleRecord{StartedAt: time.Now(), EndedAt: time.Now(), Outcome: "triage"}
	payload := buildCycleApple(state, nil, rec, 3)
	if payload.AppleType != "observation" {
		t.Errorf("type = %q, want observation when triage_findings > 0", payload.AppleType)
	}
}

// ---------------------------------------------------------------------------
// buildChatTurnApple (S148-00: /chat wired to the Apple ledger)
// ---------------------------------------------------------------------------

func TestBuildChatTurnApple_Fields(t *testing.T) {
	turn := Turn{
		ID:        "turn-1",
		UserInput: "what should I work on next?",
		Final:     "S148-00 looks like the smallest buildable item.",
		Model:     "claude-sonnet-5",
		Validated: true,
	}
	payload := buildChatTurnApple("session-abc", turn)

	if payload.SourceRepo != "EMILY" {
		t.Errorf("SourceRepo = %q, want EMILY", payload.SourceRepo)
	}
	if payload.AppleType != "conversation" {
		t.Errorf("AppleType = %q, want conversation", payload.AppleType)
	}
	if payload.RunID != "chat-session-abc-turn-1" {
		t.Errorf("RunID = %q, want chat-session-abc-turn-1", payload.RunID)
	}
	if !strings.Contains(payload.Title, "what should I work on next") {
		t.Errorf("Title = %q, should contain the user's message", payload.Title)
	}
	if !strings.Contains(payload.Body, turn.UserInput) || !strings.Contains(payload.Body, turn.Final) {
		t.Errorf("Body should contain both the user input and the final reply, got: %q", payload.Body)
	}

	var meta map[string]any
	if err := json.Unmarshal(payload.Metadata, &meta); err != nil {
		t.Fatalf("Metadata should be valid JSON: %v", err)
	}
	if meta["session_id"] != "session-abc" || meta["turn_id"] != "turn-1" {
		t.Errorf("Metadata missing session_id/turn_id: %v", meta)
	}
	if meta["validated"] != true {
		t.Errorf("Metadata.validated = %v, want true", meta["validated"])
	}
}

func TestBuildChatTurnApple_TitleTruncatesLongMessages(t *testing.T) {
	longMsg := strings.Repeat("a", 200)
	turn := Turn{ID: "t2", UserInput: longMsg, Final: "ok"}
	payload := buildChatTurnApple("s1", turn)
	if len(payload.Title) > 73 { // truncate's own "..." allowance
		t.Errorf("Title should be truncated, got length %d", len(payload.Title))
	}
}

func TestBuildChatTurnApple_IncludesValidationNoteWhenPresent(t *testing.T) {
	turn := Turn{ID: "t3", UserInput: "q", Final: "a", ValidationNote: "flagged: unverified claim"}
	payload := buildChatTurnApple("s1", turn)
	if !strings.Contains(payload.Body, "flagged: unverified claim") {
		t.Errorf("Body should include the validation note when present, got: %q", payload.Body)
	}
}

func TestBuildChatTurnApple_OmitsValidationNoteWhenAbsent(t *testing.T) {
	turn := Turn{ID: "t4", UserInput: "q", Final: "a"}
	payload := buildChatTurnApple("s1", turn)
	if strings.Contains(payload.Body, "validation note") {
		t.Errorf("Body should not mention a validation note when there isn't one, got: %q", payload.Body)
	}
}

func TestPostChatTurnApple_NoopWithoutIdunaEnv(t *testing.T) {
	for _, k := range []string{"IDUNA_BASE_URL", "IDUNA_AGENT_NAME", "IDUNA_AGENT_SECRET"} {
		t.Setenv(k, "")
	}
	s := &Server{}
	// Must not panic even with a zero-value Server -- postChatTurnApple
	// should short-circuit entirely on the nil IdunaClient before touching
	// anything else on s.
	s.postChatTurnApple("session-x", Turn{ID: "t5", UserInput: "q", Final: "a"})
}
