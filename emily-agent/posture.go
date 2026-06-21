// posture.go — lightweight EmilyOS posture reader for emily-agent.
//
// Emily-agent does not import emilyos/ as a module dependency.
// Instead, it reads var/posture.json directly (same on-disk format EmilyOS writes).
// When EMILY_POSTURE_PATH is set and posture=SIEGE or EXITED, LLM calls are skipped.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// postureState is the on-disk EmilyOS posture format (matches emilyos/internal/posture).
type postureState struct {
	Posture string `json:"posture"`
}

// PostureNormal, PostureSiege, PostureExited are the relevant posture states.
const (
	PostureNormal  = "NORMAL"
	PostureSiege   = "SIEGE"
	PostureMercy   = "MERCY"
	PostureGame    = "GAME"
	PostureExited  = "EXITED"
	PostureUnknown = ""
)

// readPosture reads the current posture from EMILY_POSTURE_PATH.
// Returns PostureNormal if the file doesn't exist or the env var is unset.
func readPosture() string {
	path := os.Getenv("EMILY_POSTURE_PATH")
	if path == "" {
		return PostureNormal
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return PostureNormal
	}
	if err != nil {
		log.Printf("[posture] read %s: %v — assuming NORMAL", path, err)
		return PostureNormal
	}
	var ps postureState
	if err := json.Unmarshal(data, &ps); err != nil {
		log.Printf("[posture] parse: %v — assuming NORMAL", err)
		return PostureNormal
	}
	return ps.Posture
}

// postureBlocksLLM returns true (and logs) if the current posture disallows LLM calls.
// SIEGE: no external network calls including LLM.
// EXITED: system has exited; refuse to act.
func postureBlocksLLM() bool {
	p := readPosture()
	switch p {
	case PostureSiege:
		log.Printf("[posture] SIEGE — LLM calls disabled (cap.net ForceOff); skipping cycle")
		return true
	case PostureExited:
		log.Printf("[posture] EXITED — refusing to run; graceful shutdown expected")
		return true
	}
	return false
}

// assertPostureNotExited returns an error if posture is EXITED.
func assertPostureNotExited() error {
	if readPosture() == PostureExited {
		return fmt.Errorf("EmilyOS posture is EXITED — agent startup refused")
	}
	return nil
}
