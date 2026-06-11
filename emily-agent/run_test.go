package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSplitArgs(t *testing.T) {
	cases := []struct {
		input string
		want  []string
	}{
		{"observe hello", []string{"observe", "hello"}},
		{"observe 'hello world'", []string{"observe", "hello world"}},
		{"backlog promote --limit=10", []string{"backlog", "promote", "--limit=10"}},
		{"", nil},
		{"  ", nil},
	}
	for _, c := range cases {
		got := splitArgs(c.input)
		if len(got) != len(c.want) {
			t.Errorf("splitArgs(%q) = %v, want %v", c.input, got, c.want)
			continue
		}
		for i := range got {
			if got[i] != c.want[i] {
				t.Errorf("splitArgs(%q)[%d] = %q, want %q", c.input, i, got[i], c.want[i])
			}
		}
	}
}

func TestHandleRun_BlockedSubcommand(t *testing.T) {
	srv := &Server{}
	body, _ := json.Marshal(map[string]string{"command": "install /tmp/evil"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/emily/run", bytes.NewReader(body))
	w := httptest.NewRecorder()
	srv.handleRun(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestHandleRun_MethodNotAllowed(t *testing.T) {
	srv := &Server{}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/emily/run", nil)
	w := httptest.NewRecorder()
	srv.handleRun(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestHandleRun_EmptyCommand(t *testing.T) {
	srv := &Server{}
	body, _ := json.Marshal(map[string]string{"command": ""})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/emily/run", bytes.NewReader(body))
	w := httptest.NewRecorder()
	srv.handleRun(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
