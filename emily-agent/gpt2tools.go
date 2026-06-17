// emily-agent/gpt2tools.go
// Emily Prime GPT-2 tools: gpt2_generate, gpt2_health, gpt2_start.
//
// These let Emily Prime call the GPT-2 inference stack directly from a
// conversation — generate text, check service health, or start the server
// and proxy if they are not running.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"time"
)

const (
	gpt2ServeURL = "http://localhost:8088"
	gpt2ProxyURL = "http://localhost:8679"
)

var gpt2HealthClient = &http.Client{Timeout: 5 * time.Second}

// registerGPT2Tools adds gpt2_generate, gpt2_health, and gpt2_start to the dispatcher.
func registerGPT2Tools(d *ToolDispatcher) {

	// ── gpt2_generate ─────────────────────────────────────────────────────────
	d.Register(ToolDef{
		Name:        "gpt2_generate",
		Description: "Generate text using the Emily Prime GPT-2 fine-tuned inference server running on :8088. Returns the generated continuation and token count. The server must be running (use gpt2_start or emily gpt2 start).",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolPropSchema{
				"prompt":      {Type: "string", Description: "The text prompt to continue"},
				"max_tokens":  {Type: "number", Description: "Maximum tokens to generate (default 100)"},
				"temperature": {Type: "number", Description: "Sampling temperature 0.1–2.0 (default 0.8)"},
			},
			Required: []string{"prompt"},
		},
	}, func(args map[string]any) (string, error) {
		prompt := stringArg(args, "prompt")
		if prompt == "" {
			return "", fmt.Errorf("prompt is required")
		}
		maxTok := 100
		if v, ok := args["max_tokens"].(float64); ok && v > 0 {
			maxTok = int(v)
		}
		temp := 0.8
		if v, ok := args["temperature"].(float64); ok && v > 0 {
			temp = v
		}

		reqBody, _ := json.Marshal(map[string]any{
			"prompt":      prompt,
			"max_tokens":  maxTok,
			"temperature": temp,
		})
		req, _ := http.NewRequest(http.MethodPost, gpt2ServeURL+"/generate", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 60 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return "", fmt.Errorf("gpt2 server unreachable (run: emily gpt2 start): %w", err)
		}
		defer resp.Body.Close()

		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("gpt2 server returned %d: %s", resp.StatusCode, raw)
		}

		var out struct {
			Text   string `json:"text"`
			Tokens int    `json:"tokens"`
			Model  string `json:"model"`
			Error  string `json:"error"`
		}
		if err := json.Unmarshal(raw, &out); err != nil {
			return string(raw), nil
		}
		if out.Error != "" {
			return "", fmt.Errorf("gpt2 generate error: %s", out.Error)
		}

		result := map[string]any{
			"prompt":         prompt,
			"generated_text": out.Text,
			"full_text":      prompt + out.Text,
			"tokens":         out.Tokens,
			"model":          out.Model,
		}
		data, _ := json.Marshal(result)
		return string(data), nil
	})

	// ── gpt2_health ───────────────────────────────────────────────────────────
	d.Register(ToolDef{
		Name:        "gpt2_health",
		Description: "Check the health of the GPT-2 inference stack: the Python serve.py on :8088 and the broker proxy on :8679. Returns status for each component. Use this before gpt2_generate to ensure the server is ready.",
		Parameters: ToolParameters{
			Type:       "object",
			Properties: map[string]ToolPropSchema{},
		},
	}, func(args map[string]any) (string, error) {
		type componentStatus struct {
			Name      string `json:"name"`
			Port      string `json:"port"`
			Status    string `json:"status"`
			Model     string `json:"model,omitempty"`
			ErrorMsg  string `json:"error,omitempty"`
		}

		var components []componentStatus

		// serve.py :8088
		serveStatus := componentStatus{Name: "gpt2-serve", Port: ":8088"}
		resp, err := gpt2HealthClient.Get(gpt2ServeURL + "/health")
		if err != nil {
			serveStatus.Status = "unavailable"
			serveStatus.ErrorMsg = err.Error()
		} else {
			raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				serveStatus.Status = "ok"
				var body map[string]any
				json.Unmarshal(raw, &body)
				if m, ok := body["model"].(string); ok {
					serveStatus.Model = m
				}
			} else {
				serveStatus.Status = fmt.Sprintf("error_%d", resp.StatusCode)
				serveStatus.ErrorMsg = string(raw)
			}
		}
		components = append(components, serveStatus)

		// broker proxy :8679
		proxyStatus := componentStatus{Name: "gpt2-proxy", Port: ":8679"}
		resp2, err := gpt2HealthClient.Get(gpt2ProxyURL + "/health")
		if err != nil {
			proxyStatus.Status = "unavailable"
			proxyStatus.ErrorMsg = "run: emily gpt2 proxy"
		} else {
			resp2.Body.Close()
			if resp2.StatusCode < 500 {
				proxyStatus.Status = "ok"
			} else {
				proxyStatus.Status = fmt.Sprintf("error_%d", resp2.StatusCode)
			}
		}
		components = append(components, proxyStatus)

		result := map[string]any{
			"components": components,
			"inference_ready": serveStatus.Status == "ok",
			"proxy_ready":     proxyStatus.Status == "ok",
		}
		data, _ := json.Marshal(result)
		return string(data), nil
	})

	// ── gpt2_start ────────────────────────────────────────────────────────────
	d.Register(ToolDef{
		Name:        "gpt2_start",
		Description: "Start the GPT-2 inference server and/or broker proxy if they are not already running. Runs 'emily gpt2 start' and/or 'emily gpt2 proxy' as background processes. After starting, wait a few seconds then call gpt2_health to confirm readiness.",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolPropSchema{
				"component": {
					Type:        "string",
					Description: "Which component to start: 'server' (serve.py :8088), 'proxy' (broker :8679), or 'all' (default: all)",
				},
				"model": {
					Type:        "string",
					Description: "Model to load for the server: 'ft' (fine-tuned, default) or 'base'",
				},
			},
		},
	}, func(args map[string]any) (string, error) {
		component := stringArg(args, "component")
		if component == "" {
			component = "all"
		}
		model := stringArg(args, "model")
		if model == "" {
			model = "ft"
		}

		type startResult struct {
			Component string `json:"component"`
			Action    string `json:"action"`
			Output    string `json:"output,omitempty"`
			Error     string `json:"error,omitempty"`
		}
		var results []startResult

		startCmd := func(name string, cmdArgs []string) startResult {
			r := startResult{Component: name, Action: "started"}
			cmd := exec.Command("emily", cmdArgs...)
			out, err := cmd.CombinedOutput()
			r.Output = string(out)
			if err != nil {
				// emily gpt2 start returns 0 when already running; nonzero is a real error
				r.Action = "error"
				r.Error = err.Error()
			}
			return r
		}

		if component == "server" || component == "all" {
			results = append(results, startCmd("gpt2-serve", []string{"gpt2", "start", "--model", model}))
		}
		if component == "proxy" || component == "all" {
			results = append(results, startCmd("gpt2-proxy", []string{"gpt2", "proxy"}))
		}

		data, _ := json.Marshal(map[string]any{
			"results": results,
			"hint":    "call gpt2_health after a few seconds to confirm readiness",
		})
		return string(data), nil
	})
}
