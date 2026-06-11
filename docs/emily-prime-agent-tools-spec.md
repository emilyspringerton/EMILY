# Emily Prime Agent Tools — Golden Specification
## Owner: Emily Prime | 2026-06-11
## Status: Implementation Ready

---

## What This Is

Emily Prime needs agency over her own repo. Currently she can read conversations
and list files in the conversation store, but she cannot read the broader codebase,
cannot write to BACKLOG.md, and has no visibility into the state of FatBaby,
IDUNA, or TYLER beyond what she sees in Apple feed summaries.

This spec defines four tools that close that gap:

| Tool | Scope | Operation |
|---|---|---|
| `emily_list_files` | EMILY, PRRJECT_FATBABY, IDUNA | Read-only: list files in a repo directory |
| `emily_read_file` | EMILY, PRRJECT_FATBABY, IDUNA | Read-only: read file contents |
| `emily_write_file` | EMILY only | Write/append a file + auto-git-commit |
| `emily_run_shell` | EMILY only — allow-listed commands | Run safe read-only shell ops (git log, go test) |

---

## Safety Model

- **Read tools** are sandboxed to three allowed roots: EMILY, PRRJECT_FATBABY, IDUNA.
  Any path traversal (`..`) is rejected. Symlinks are not followed beyond the root.
- **Write tool** is sandboxed to EMILY only. No `..` traversal. Auto-commits every write
  so every Emily Prime mutation has a git audit trail.
- **Shell tool** (future, not in this spec) uses an explicit allow-list of commands.
  Not implemented here.

---

## Tool Definitions

### `emily_list_files`

```json
{
  "name": "emily_list_files",
  "description": "List files in a directory of one of the permitted repos: EMILY, PRRJECT_FATBABY, or IDUNA. Use repo='EMILY' and path='.' to see the top-level EMILY repo. Subdirectories are listed with a trailing slash.",
  "parameters": {
    "type": "object",
    "properties": {
      "repo":  { "type": "string", "description": "One of: EMILY, PRRJECT_FATBABY, IDUNA" },
      "path":  { "type": "string", "description": "Relative path within the repo (default: '.')" },
      "depth": { "type": "integer", "description": "Max depth (default 1, max 3)" }
    },
    "required": ["repo"]
  }
}
```

---

### `emily_read_file`

```json
{
  "name": "emily_read_file",
  "description": "Read the full contents of a file in one of the permitted repos: EMILY, PRRJECT_FATBABY, or IDUNA. Returns an error if the file does not exist or if the path escapes the repo root. Limit: 200KB.",
  "parameters": {
    "type": "object",
    "properties": {
      "repo": { "type": "string", "description": "One of: EMILY, PRRJECT_FATBABY, IDUNA" },
      "path": { "type": "string", "description": "Relative file path within the repo" }
    },
    "required": ["repo", "path"]
  }
}
```

---

### `emily_write_file`

```json
{
  "name": "emily_write_file",
  "description": "Write or append to a file in the EMILY repo. Creates parent directories if needed. Automatically runs `git add <path> && git commit` after writing. The commit message is auto-generated from the path and operation unless commit_message is provided. Use mode='append' to add to an existing file (e.g. BACKLOG.md). Use mode='write' to overwrite.",
  "parameters": {
    "type": "object",
    "properties": {
      "path":           { "type": "string", "description": "Relative path within EMILY repo" },
      "content":        { "type": "string", "description": "Content to write or append" },
      "mode":           { "type": "string", "description": "write | append (default: write)" },
      "commit_message": { "type": "string", "description": "Git commit message (auto-generated if omitted)" }
    },
    "required": ["path", "content"]
  }
}
```

---

## Implementation

New file: `emily-agent/emilytools.go`
- `registerEmilyPrimeTools(d *ToolDispatcher, emilyRoot, fatbabyRoot, idunaRoot string)`
- Roots resolved from env: EMILY_ROOT, FATBABY_ROOT, IDUNA_ROOT (with /home/fatbaby/ defaults)
- Called from `NewServer()` after `registerWebAuditTools(dispatcher)`

Config additions:
- `EmilyRoot string` — env EMILY_ROOT, default /home/fatbaby/EMILY
- `FatBabyRoot string` — env FATBABY_ROOT, default /home/fatbaby/PRRJECT_FATBABY
- `IDUNARoot string` — env IDUNA_ROOT, default /home/fatbaby/IDUNA

---

## How Emily Prime Uses These

1. **Curating the backlog**: Emily Prime reads BACKLOG.md via emily_read_file, identifies
   a new observation, appends a `[ ]` item via emily_write_file, which auto-commits.

2. **Auditing FatBaby config**: Emily Prime reads entity-graph-rules.json or watchlist.json
   to understand current signal configuration before writing directed tasks.

3. **Cross-repo review**: Emily Prime reads IDUNA golden.md or TYLER BACKLOG.md to understand
   the system state across repos before planning the next sprint.

4. **Writing golden docs**: Emily Prime appends new observations, decisions, or specs to
   EMILY/docs/ via emily_write_file. Every write is auto-committed.

---

## Roadmap

- `emily_run_shell` — allow-listed: `git log`, `git diff`, `go test ./...` in any permitted repo
- `emily_post_apple` — direct Apple submission without routing through the CLI
- `emily_prime_task` — direct task dispatch (already exists via CLI; add as agent tool for symmetry)
