package main

import (
	"strings"
	"testing"
	"time"
)

func TestBuildBriefingMessageEmpty(t *testing.T) {
	title, body := buildBriefingMessage(nil)
	if !strings.Contains(title, "Good morning") {
		t.Errorf("empty: title = %q, want 'Good morning'", title)
	}
	if !strings.Contains(body, "quiet") {
		t.Errorf("empty: body = %q, want 'quiet'", body)
	}
}

func TestBuildBriefingMessageCounts(t *testing.T) {
	now := time.Now()
	apples := []AppleListItem{
		{AppleType: "completion", Title: "S57 done", RecordedAt: now},
		{AppleType: "completion", Title: "S58 done", RecordedAt: now.Add(-1 * time.Hour)},
		{AppleType: "observation", Title: "obs note", RecordedAt: now.Add(-2 * time.Hour)},
	}
	title, body := buildBriefingMessage(apples)
	if !strings.Contains(title, "3 apples") {
		t.Errorf("title should mention 3 apples: %q", title)
	}
	if !strings.Contains(body, "2") {
		t.Errorf("body should show 2 completions: %q", body)
	}
}

func TestBuildBriefingMessageTitleCap(t *testing.T) {
	// Many distinct types → title falls back to "N apples (M types)"
	var apples []AppleListItem
	types := []string{"completion", "observation", "audit", "escalation", "improvement"}
	for i, tp := range types {
		for j := 0; j <= i; j++ {
			apples = append(apples, AppleListItem{AppleType: tp, Title: "t"})
		}
	}
	title, _ := buildBriefingMessage(apples)
	// Title should be ≤90 chars or fallback form used.
	if len(title) > 100 {
		t.Errorf("title too long: %d chars: %q", len(title), title)
	}
}

func TestBuildBriefingMessageNotable(t *testing.T) {
	now := time.Now()
	// 4 completions — only 3 most recent should appear as bullets.
	apples := []AppleListItem{
		{AppleType: "completion", Title: "Newest", RecordedAt: now},
		{AppleType: "completion", Title: "Second", RecordedAt: now.Add(-1 * time.Hour)},
		{AppleType: "completion", Title: "Third", RecordedAt: now.Add(-2 * time.Hour)},
		{AppleType: "completion", Title: "Oldest (should be excluded)", RecordedAt: now.Add(-3 * time.Hour)},
	}
	_, body := buildBriefingMessage(apples)
	if strings.Contains(body, "Oldest") {
		t.Errorf("4th-oldest notable should be excluded, but found in body: %q", body)
	}
	if !strings.Contains(body, "Newest") {
		t.Errorf("most recent should appear in body: %q", body)
	}
}

func TestTypeLabel(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"completion", "✓"},
		{"escalation", "!"},
		{"improvement", "↑"},
		{"signal_observation", "~"},
		{"rsi_iteration", "⟳"},
		{"status", "·"},
		{"unknown-type", "unknown-type"}, // default returns input
	}
	for _, tc := range cases {
		got := typeLabel(tc.input)
		if got != tc.want {
			t.Errorf("typeLabel(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
