package main

import (
	"context"
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
