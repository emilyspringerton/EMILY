// api_memory.go — POST /api/v1/emily/memory/update
//
// Allows Emily Prime (or any authorized caller) to patch a named section in
// emily-memory/world-state.md. Sections are delimited by "### SECTION_NAME" headers.
// This makes Emily Prime self-updating — resolves AGI Gap 4 (world-state.md is static).

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// handleMemoryUpdate patches a named section in emily-memory/world-state.md.
// POST /api/v1/emily/memory/update
// Body: { "section": "SHANKPIT", "content": "- new state line\n- another line" }
func (s *Server) handleMemoryUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Section string `json:"section"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusBadRequest)
		return
	}
	req.Section = strings.TrimSpace(req.Section)
	if req.Section == "" || req.Content == "" {
		http.Error(w, `{"error":"section and content required"}`, http.StatusBadRequest)
		return
	}

	path := filepath.Join(s.cfg.EmilyRoot, "emily-memory", "world-state.md")
	data, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"world-state.md not found: %v"}`, err), http.StatusInternalServerError)
		return
	}

	content := string(data)
	header := "### " + req.Section
	idx := strings.Index(content, header)
	if idx < 0 {
		http.Error(w, fmt.Sprintf(`{"error":"section not found: %q"}`, req.Section), http.StatusNotFound)
		return
	}

	// Find the end of this section: next ### or ## heading, or EOF.
	afterHeader := idx + len(header)
	rest := content[afterHeader:]
	endIdx := len(rest)
	for _, marker := range []string{"\n### ", "\n## "} {
		if i := strings.Index(rest, marker); i >= 0 && i < endIdx {
			endIdx = i
		}
	}

	newContent := content[:idx] +
		header + "\n" +
		strings.TrimRight(req.Content, "\n") + "\n" +
		rest[endIdx:]

	if err := os.WriteFile(path, []byte(newContent), 0o644); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"write failed: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"updated","section":%q}`+"\n", req.Section)
}
