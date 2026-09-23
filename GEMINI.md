# GEMINI.md — Guidance for Gemini / Antigravity in EMILY

## What This Is

`EMILY` is the home of Emily Prime — the meta-orchestrator and chief of staff for EINHORN_INDUSTRIAL. She runs as a Go service on port 8086, executes RSI cycles every 15 minutes, triages observations, manages the golden backlog, files Apples to IDUNA, and issues tasks.

Key references:
- `BACKLOG.md`: Canonical cross-repo backlog.
- `GOLDEN.md`: Compressed backlog context for prompt loading.
- `docs/THE_EMILY_WAY.md`: Core operating principles.
- `context/golden-docs-index.md`: Registry of system specs and North Star docs.

## Operating Protocols (The Emily Way)

1. **Founder Direction**:
   Always post through `emily observe -s info "Founder real-time: <summary>"` first.
2. **Backlog Workflow**:
   Pick unchecked items from `BACKLOG.md`. Mark done with Apple ID.
3. **Apples**:
   Post completion Apples via `emily apples post -t completion "<title>"`.
4. **Session Tag**:
   Auto-tagged via `emily session current` on commits and Apples.
5. **Commit Protocol**:
   Commit and push immediately when verified, with the `session: <tag>` trailer.
