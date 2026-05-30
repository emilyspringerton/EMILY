package main

import (
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
