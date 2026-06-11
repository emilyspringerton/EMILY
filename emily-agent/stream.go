// stream.go — append-only NDJSON event stream for Emily OS system observability.
// All system inputs (cycles, apples, act errors, triage) feed into
// var/emily-stream.ndjson. The file is git-committed by `emily sync --stream`.

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// StreamEntry is one line in the emily-stream.ndjson log.
type StreamEntry struct {
	TS     string         `json:"ts"`
	Source string         `json:"source"`
	Type   string         `json:"type"`
	Cycle  int            `json:"cycle,omitempty"`
	Msg    string         `json:"msg,omitempty"`
	Data   map[string]any `json:"data,omitempty"`
}

var streamMu sync.Mutex

// streamPath returns the path to emily-stream.ndjson, derived from stateDir.
// The stream lives at <EMILY_ROOT>/var/emily-stream.ndjson where EMILY_ROOT is
// one level above the stateDir (typically ./emily-state → ./).
func streamPath(stateDir string) string {
	root := filepath.Dir(stateDir)
	return filepath.Join(root, "var", "emily-stream.ndjson")
}

// appendStream writes one NDJSON line to the stream file. Thread-safe, best-effort.
func appendStream(path string, entry StreamEntry) {
	if entry.TS == "" {
		entry.TS = time.Now().UTC().Format(time.RFC3339)
	}
	if entry.Source == "" {
		entry.Source = "emily-agent"
	}

	streamMu.Lock()
	defer streamMu.Unlock()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()

	line, err := json.Marshal(entry)
	if err != nil {
		return
	}
	_, _ = f.Write(append(line, '\n'))
}
