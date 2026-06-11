// emily-agent/emilytools.go
// Emily Prime repo-access tools: emily_list_files, emily_read_file, emily_write_file.
//
// Read tools are sandboxed to EMILY, PRRJECT_FATBABY, and IDUNA repos.
// Write tool is sandboxed to EMILY only and auto-commits every mutation.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const maxReadBytes = 200 * 1024 // 200 KB

// registerEmilyPrimeTools adds emily_list_files, emily_read_file, emily_write_file
// to the dispatcher. All root arguments must be absolute paths.
func registerEmilyPrimeTools(d *ToolDispatcher, emilyRoot, fatbabyRoot, idunaRoot string) {
	roots := map[string]string{
		"EMILY":           emilyRoot,
		"PRRJECT_FATBABY": fatbabyRoot,
		"IDUNA":           idunaRoot,
	}

	resolveRoot := func(args map[string]any) (string, error) {
		repo, _ := args["repo"].(string)
		root, ok := roots[repo]
		if !ok {
			return "", fmt.Errorf("unknown repo %q — use EMILY, PRRJECT_FATBABY, or IDUNA", repo)
		}
		return root, nil
	}

	// ── emily_list_files ──────────────────────────────────────────────────────
	d.Register(ToolDef{
		Name:        "emily_list_files",
		Description: "List files in a directory of one of the permitted repos: EMILY, PRRJECT_FATBABY, or IDUNA. Directories are shown with a trailing slash.",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolPropSchema{
				"repo": {
					Type:        "string",
					Description: "One of: EMILY, PRRJECT_FATBABY, IDUNA",
				},
				"path": {
					Type:        "string",
					Description: "Relative path within the repo (default: '.')",
				},
				"depth": {
					Type:        "integer",
					Description: "Max depth (default 1, max 3)",
				},
			},
			Required: []string{"repo"},
		},
	}, func(args map[string]any) (string, error) {
		root, err := resolveRoot(args)
		if err != nil {
			return "", err
		}
		relPath, _ := args["path"].(string)
		if relPath == "" {
			relPath = "."
		}
		depth := 1
		if d, ok := args["depth"].(float64); ok && d > 0 {
			depth = int(d)
		}
		if depth > 3 {
			depth = 3
		}
		safe, err := safePath(root, relPath)
		if err != nil {
			return "", err
		}
		lines, err := emilyListDir(safe, depth, 0)
		if err != nil {
			return "", err
		}
		return strings.Join(lines, "\n"), nil
	})

	// ── emily_read_file ───────────────────────────────────────────────────────
	d.Register(ToolDef{
		Name:        "emily_read_file",
		Description: "Read the contents of a file in one of the permitted repos: EMILY, PRRJECT_FATBABY, or IDUNA. Limit: 200KB.",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolPropSchema{
				"repo": {
					Type:        "string",
					Description: "One of: EMILY, PRRJECT_FATBABY, IDUNA",
				},
				"path": {
					Type:        "string",
					Description: "Relative file path within the repo",
				},
			},
			Required: []string{"repo", "path"},
		},
	}, func(args map[string]any) (string, error) {
		root, err := resolveRoot(args)
		if err != nil {
			return "", err
		}
		relPath, _ := args["path"].(string)
		if relPath == "" {
			return "", fmt.Errorf("path is required")
		}
		safe, err := safePath(root, relPath)
		if err != nil {
			return "", err
		}
		info, err := os.Stat(safe)
		if err != nil {
			return "", fmt.Errorf("file not found: %s", relPath)
		}
		if info.IsDir() {
			return "", fmt.Errorf("%s is a directory — use emily_list_files", relPath)
		}
		if info.Size() > maxReadBytes {
			return "", fmt.Errorf("file too large (%dKB > 200KB): %s", info.Size()/1024, relPath)
		}
		data, err := os.ReadFile(safe)
		if err != nil {
			return "", fmt.Errorf("read error: %w", err)
		}
		return string(data), nil
	})

	// ── emily_write_file ──────────────────────────────────────────────────────
	d.Register(ToolDef{
		Name:        "emily_write_file",
		Description: "Write or append to a file in the EMILY repo. Auto-commits after writing. Use mode='append' to add to an existing file (e.g. BACKLOG.md). Every write gets a git commit.",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolPropSchema{
				"path": {
					Type:        "string",
					Description: "Relative path within the EMILY repo",
				},
				"content": {
					Type:        "string",
					Description: "Content to write or append",
				},
				"mode": {
					Type:        "string",
					Description: "write | append (default: write)",
				},
				"commit_message": {
					Type:        "string",
					Description: "Git commit message (auto-generated if omitted)",
				},
			},
			Required: []string{"path", "content"},
		},
	}, func(args map[string]any) (string, error) {
		relPath, _ := args["path"].(string)
		content, _ := args["content"].(string)
		mode, _ := args["mode"].(string)
		commitMsg, _ := args["commit_message"].(string)

		if relPath == "" {
			return "", fmt.Errorf("path is required")
		}
		if mode == "" {
			mode = "write"
		}

		safe, err := safePath(emilyRoot, relPath)
		if err != nil {
			return "", err
		}
		if err := os.MkdirAll(filepath.Dir(safe), 0755); err != nil {
			return "", fmt.Errorf("mkdir: %w", err)
		}

		flag := os.O_CREATE | os.O_WRONLY
		if mode == "append" {
			flag |= os.O_APPEND
		} else {
			flag |= os.O_TRUNC
		}
		f, err := os.OpenFile(safe, flag, 0644)
		if err != nil {
			return "", fmt.Errorf("open: %w", err)
		}
		_, writeErr := f.WriteString(content)
		f.Close()
		if writeErr != nil {
			return "", fmt.Errorf("write: %w", writeErr)
		}

		if commitMsg == "" {
			op := "write"
			if mode == "append" {
				op = "append"
			}
			commitMsg = fmt.Sprintf("emily-prime: %s %s (%s)", op, relPath, time.Now().UTC().Format("2006-01-02T15:04Z"))
		}
		if err := emilyGitAddCommit(emilyRoot, relPath, commitMsg); err != nil {
			return fmt.Sprintf("wrote %s but git commit failed: %v", relPath, err), nil
		}
		return fmt.Sprintf("ok: %s %s — committed", mode, relPath), nil
	})
}

// emilyListDir recursively lists files/dirs up to maxDepth.
func emilyListDir(path string, maxDepth, curDepth int) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir()
		}
		return entries[i].Name() < entries[j].Name()
	})
	var lines []string
	prefix := strings.Repeat("  ", curDepth)
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}
		if e.IsDir() {
			lines = append(lines, prefix+e.Name()+"/")
			if curDepth+1 < maxDepth {
				sub, _ := emilyListDir(filepath.Join(path, e.Name()), maxDepth, curDepth+1)
				lines = append(lines, sub...)
			}
		} else {
			lines = append(lines, prefix+e.Name())
		}
	}
	return lines, nil
}

// emilyGitAddCommit stages relPath and creates a commit in repoDir.
func emilyGitAddCommit(repoDir, relPath, msg string) error {
	add := exec.Command("git", "-C", repoDir, "add", relPath)
	if out, err := add.CombinedOutput(); err != nil {
		return fmt.Errorf("git add: %w: %s", err, strings.TrimSpace(string(out)))
	}
	commit := exec.Command("git", "-C", repoDir, "commit", "-m", msg)
	if out, err := commit.CombinedOutput(); err != nil {
		if strings.Contains(string(out), "nothing to commit") {
			return nil
		}
		return fmt.Errorf("git commit: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}
