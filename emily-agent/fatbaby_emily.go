// fatbaby_emily.go
// FatBaby Emily — Emily Prime's operational sub-agent.
//
// Emily Prime orchestrates by chatting to FatBaby Emily in natural language.
// FatBaby Emily has direct access to FatBaby processes and the integration layer.
// Emily Prime sends a message; FatBaby Emily executes and reports back.
// The exchange is surfaced in the web chat UI as a visible sub-conversation.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// SubChatResult is the structured result of a chat_to_fatbaby_emily tool call.
// Stored as JSON in the ToolActivity.Result field; the UI parses and renders it specially.
type SubChatResult struct {
	SentTo         string         `json:"sent_to"`
	Message        string         `json:"message"`
	Reply          string         `json:"reply"`
	ToolActivities []ToolActivity `json:"tool_activities,omitempty"`
}

const fatbabyEmilySystemPrompt = `You are FatBaby Emily — Emily Prime's operational executor and FatBaby system controller.

Emily Prime gives you natural-language instructions. You execute them using your tools and report back concisely.

You have access to:
  fatbaby_start_process    — start a named FatBaby process (signalapi, obs-watcher, newssite, entity-graph, broker)
  fatbaby_list_processes   — list known FatBaby processes and whether each appears to be running
  fatbaby_read_obs         — read recent FatBaby observations from signals/observations/
  fatbaby_write_task       — write a directed task into signals/tasks/ for the obs-watcher loop
  emily_cli                — run an arbitrary emily CLI subcommand

Always report: what you did, the outcome, and any errors. Be concise.`

// FatBabyEmily is Emily Prime's FatBaby sub-agent.
type FatBabyEmily struct {
	pipeline    *Pipeline
	fatbabyRoot string
}

func NewFatBabyEmily(llm LLMClient, model string, integStore *IntegrationStore, fatbabyRoot string) *FatBabyEmily {
	d := NewToolDispatcher()
	fb := &FatBabyEmily{fatbabyRoot: fatbabyRoot}
	fb.registerTools(d, integStore)
	fb.pipeline = &Pipeline{
		Generator:    llm,
		GenModel:     model,
		Dispatcher:   d,
		MaxToolIters: 8,
	}
	return fb
}

// Chat sends a message to FatBaby Emily and returns her reply + what she did.
func (fb *FatBabyEmily) Chat(ctx context.Context, message string) (string, []ToolActivity, error) {
	history := []Message{
		{Role: "system", Content: fatbabyEmilySystemPrompt},
		{Role: "user", Content: message},
	}
	result, err := fb.pipeline.Run(ctx, history)
	if err != nil {
		return "", nil, err
	}
	return result.Final, result.ToolActivities, nil
}

func (fb *FatBabyEmily) registerTools(d *ToolDispatcher, integStore *IntegrationStore) {
	if integStore != nil {
		registerIntegrationTools(d, integStore)
	}

	d.Register(ToolDef{
		Name:        "fatbaby_list_processes",
		Description: "List FatBaby processes and whether each appears to be running (checks process table).",
		Parameters:  ToolParameters{Type: "object", Properties: map[string]ToolPropSchema{}},
	}, func(args map[string]any) (string, error) {
		procs := []string{"signalapi", "obs-watcher", "newssite", "entity-graph", "broker", "emily-agent", "feedserver"}
		out, err := exec.Command("bash", "-c",
			`ps aux | grep -E '`+strings.Join(procs, "|")+`' | grep -v grep`).Output()
		if err != nil || len(strings.TrimSpace(string(out))) == 0 {
			return "no FatBaby processes found in process table", nil
		}
		return strings.TrimSpace(string(out)), nil
	})

	d.Register(ToolDef{
		Name:        "fatbaby_start_process",
		Description: "Start a FatBaby process. Valid names: signalapi, obs-watcher, newssite, entity-graph, broker. Uses emily start --<name>.",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolPropSchema{
				"process": {Type: "string", Description: "signalapi | obs-watcher | newssite | entity-graph | broker"},
			},
			Required: []string{"process"},
		},
	}, func(args map[string]any) (string, error) {
		proc := stringArg(args, "process")
		if proc == "" {
			return "", fmt.Errorf("process name required")
		}
		out, err := fbRunEmilyCmd("start", "--"+proc)
		if err != nil {
			return fmt.Sprintf("start %s: error: %v\n%s", proc, err, out), nil
		}
		return fmt.Sprintf("started %s\n%s", proc, out), nil
	})

	d.Register(ToolDef{
		Name:        "fatbaby_read_obs",
		Description: "Read recent FatBaby observations. Equivalent to integration_read_observations.",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolPropSchema{
				"limit": {Type: "number", Description: "Max observations to return (default 10)"},
			},
		},
	}, func(args map[string]any) (string, error) {
		limit := 10
		if v, ok := args["limit"].(float64); ok && v > 0 {
			limit = int(v)
		}
		if integStore == nil {
			return "integration store not available", nil
		}
		obs, err := integStore.ReadObservations(limit)
		if err != nil {
			return "", err
		}
		if len(obs) == 0 {
			return "no observations found", nil
		}
		data, _ := json.MarshalIndent(obs, "", "  ")
		return string(data), nil
	})

	d.Register(ToolDef{
		Name:        "emily_cli",
		Description: "Run an emily CLI subcommand. Pass the subcommand and args as a space-separated string (without the 'emily' prefix). Example: 'status' or 'apples post -t audit REPO title body'.",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolPropSchema{
				"command": {Type: "string", Description: "Emily subcommand + args (space-separated, no leading 'emily')"},
			},
			Required: []string{"command"},
		},
	}, func(args map[string]any) (string, error) {
		cmd := stringArg(args, "command")
		if cmd == "" {
			return "", fmt.Errorf("command required")
		}
		out, err := fbRunEmilyCmd(strings.Fields(cmd)...)
		if err != nil {
			return fmt.Sprintf("emily %s → error: %v\n%s", cmd, err, out), nil
		}
		return out, nil
	})
}

func fbRunEmilyCmd(args ...string) (string, error) {
	cmd := exec.Command("emily", args...)
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}
