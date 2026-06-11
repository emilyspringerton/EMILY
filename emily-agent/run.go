// emily-agent/run.go
// POST /api/v1/emily/run — invoke emily CLI commands over HTTP.
//
// Allows external orchestration (HEIMDAL webhooks, MJOLNIR, cron) to drive
// Emily Prime without a human at the terminal.
//
// Routes wired in main.go:
//   POST /api/v1/emily/run  — execute an emily CLI subcommand; returns stdout/stderr/exit_code

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

// allowedRunCommands is the whitelist of emily subcommands that may be invoked
// via the API. This prevents arbitrary shell execution.
var allowedRunCommands = map[string]bool{
	"observe":       true,
	"obs":           true,
	"apples":        true,
	"status":        true,
	"sync":          true,
	"backlog":       true,
	"primetask":     true,
	"agents":        true,
	"watch":         true,
}

type runRequest struct {
	Command string `json:"command"` // e.g. "backlog promote" or "observe 'something happened'"
	Timeout int    `json:"timeout"` // seconds; default 60, max 300
}

type runResponse struct {
	ExitCode int    `json:"exit_code"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	Duration string `json:"duration"`
}

func (s *Server) handleRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	var req runRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	args := splitArgs(req.Command)
	if len(args) == 0 {
		http.Error(w, `"command" is required`, http.StatusBadRequest)
		return
	}

	sub := args[0]
	if !allowedRunCommands[sub] {
		http.Error(w, "subcommand not permitted: "+sub, http.StatusForbidden)
		return
	}

	timeout := req.Timeout
	if timeout <= 0 {
		timeout = 60
	}
	if timeout > 300 {
		timeout = 300
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeout)*time.Second)
	defer cancel()

	start := time.Now()
	cmd := exec.CommandContext(ctx, "emily", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	exitCode := 0
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			exitCode = ee.ExitCode()
		} else {
			exitCode = -1
		}
	}

	elapsed := time.Since(start)
	log.Printf("run: emily %s → exit %d (%s)", req.Command, exitCode, elapsed.Round(time.Millisecond))

	resp := runResponse{
		ExitCode: exitCode,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Duration: elapsed.Round(time.Millisecond).String(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// splitArgs splits a command string into argv, respecting single-quoted strings.
// Does not handle double-quotes or backslash escapes — sufficient for trusted internal use.
func splitArgs(s string) []string {
	var args []string
	var current strings.Builder
	inQuote := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '\'' && !inQuote:
			inQuote = true
		case c == '\'' && inQuote:
			inQuote = false
		case c == ' ' && !inQuote:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteByte(c)
		}
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}
