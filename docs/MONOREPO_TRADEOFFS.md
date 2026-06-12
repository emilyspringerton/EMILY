# Monorepo Consideration — Trade-off Analysis

*Written: 2026-06-12 | Status: analysis only — no action yet*

---

## Context

EINHORN_INDUSTRIAL currently uses 9+ repos: EMILY, PRRJECT_FATBABY, IDUNA, emily.cli, SHANKPIT,
GoblinFoxDragon, MJOLNIR, APPLES, TYLER, EmilyOS, PITVIPER. The primary pain points:

- **Context window sprawl**: Claude Code sessions require reading CLAUDE.md + northstar docs across
  multiple repos before meaningful work can start. Each hop costs tokens.
- **Cross-repo changes are hard to coordinate**: A change in emily.cli that requires a change in
  EMILY requires two separate commits, two separate CI runs, two separate PRs.
- **`go.mod` module boundary duplication**: SHANKPIT and GoblinFoxDragon share `module dragonsnshit`
  because they were the same codebase — the split is artificial.

---

## Option A: Full Monorepo

Merge all repos into one `emilyspringerton/einhorn` monorepo.

**Pros:**
- One `git clone`, one CI pipeline, atomic cross-repo commits.
- CLAUDE.md and context files can be consolidated; agents get the full picture in one read.
- Dependency management simplifies (emily.cli doesn't need to import emily-agent via HTTP).

**Cons:**
- Git history merge is complex and lossy. `git merge --allow-unrelated-histories` produces a
  flat history; git grafts / filter-repo preserves it but is error-prone.
- A single large repo amplifies context window pressure: agents reading `git diff HEAD~1` or
  `git log --oneline -50` see unrelated changes from all subprojects at once.
- Android build (MJOLNIR, Gradle) in a Go monorepo creates toolchain friction.
- `module dragonsnshit` (SHANKPIT/GoblinFoxDragon) and `module github.com/emilyspringerton/emily-cli`
  are already different Go module paths — merging requires renaming imports everywhere.

**Verdict: Too aggressive for current stage.** Full monorepo makes sense when cross-repo churn
is high and the team is large enough to justify the migration cost. Right now the bottleneck is
agent context quality, not commit atomicity.

---

## Option B: Partial Consolidation — SHANKPIT + GoblinFoxDragon

These two repos share `module dragonsnshit` and represent two halves of the same system. The
natural merge point is the `WorldBackend` interface (Milestone 3, done 2026-06-12).

**Pros:**
- Eliminates the artificial split. One `go build`, one test suite, one CHANGELOG.
- DragonflyBackend (Milestone 4) can be implemented in the same module without import tricks.
- Reduces the repo count by 1 and removes a source of agent context confusion.

**Cons:**
- Requires a git history merge (same options as Option A, scoped to two repos).
- GoblinFoxDragon's entity/scripting systems and SHANKPIT's FPS systems use the same package
  paths (`server/`, `apps/`, `apps2/`) — package name collisions must be resolved.
- CI for both repos must move to one pipeline.

**Verdict: Worthwhile in Milestone 4.** When DragonflyBackend implementation starts, do the
SHANKPIT + GoblinFoxDragon merge first. Use `filter-repo --to-subdirectory-filter` to preserve
history under `dragonfly/` and `shankpit/` subdirs respectively. Target completion: before
Milestone 4 code lands.

---

## Option C: Go Workspace (go.work)

Use `go.work` to create a multi-module workspace that covers EMILY, emily.cli, PRRJECT_FATBABY,
IDUNA without merging their Git histories.

**Pros:**
- Zero Git migration. Add a `go.work` file at `/home/fatbaby/go.work`.
- Cross-module development (edit emily.cli while running emily-agent) works without replace
  directives or local symlinks.
- Agents can see all `go.work` members in one `go build ./...` pass.

**Cons:**
- `go.work` is a developer tool — it is not checked in to any individual repo, so CI pipelines
  for individual repos are unchanged and don't benefit.
- Does not help with agent context or CLAUDE.md sprawl.
- Doesn't eliminate the separate repos or their separate commit histories.

**Verdict: Low-friction local dev improvement. Do it now.**

---

## Option D: Git Submodules

One parent repo (`emilyos-workspace`) that includes all repos as submodules.

**Pros:**
- Single clone, pinned versions, atomic super-repo commits.

**Cons:**
- Submodule UX is notoriously painful (detached HEAD, update commands).
- CI complexity: every repo still has its own pipeline, plus the parent.
- Does not help with context window sprawl — each submodule is still a separate context.

**Verdict: No. Submodules add friction without solving the agent-context problem.**

---

## Recommendation

| Action | When | Effort |
|--------|------|--------|
| **Add `go.work`** at `/home/fatbaby/go.work` covering EMILY, emily.cli, IDUNA, PRRJECT_FATBABY | Now | 30 min |
| **Merge SHANKPIT + GoblinFoxDragon** using `git filter-repo` | Before Milestone 4 | 1 day |
| **Full monorepo** | Not now — revisit at scale | Multi-week |

The `go.work` change is no-cost and improves cross-module development immediately. The
SHANKPIT+GoblinFoxDragon merge is the natural next step once Milestone 4 starts.

The agent-context problem is better solved by better GOLDEN.md compression and northstar doc
quality than by monorepo consolidation. See the Ops Docs Token Efficiency backlog item for that.

---

## Next Steps

1. Create `/home/fatbaby/go.work` with `use` directives for Go-based repos. *(See below.)*
2. Add `go.work` to `.gitignore` in each repo (it is workspace-local, not repo-local).
3. Schedule SHANKPIT+GoblinFoxDragon merge for Milestone 4 kickoff.

### go.work (draft)

```
go 1.21

use (
    ./EMILY/emily-agent  // or wherever the Go module root is
    ./emily.cli
    ./PRRJECT_FATBABY
    ./IDUNA
)
```

*Note: SHANKPIT and GoblinFoxDragon both use `module dragonsnshit` — they cannot both be `use`d
in a workspace until the merge is complete.*
