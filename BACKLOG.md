# EMILY PRIME — CROSS-REPO GOLDEN BACKLOG
## Owner: Emily Prime | Machine-readable | Git-authoritative
### Last updated: 2026-06-07 | Ref: Apple #3 (IDUNA) | Emily Prime decision

---

> *"The backlog is the load-bearing node. Everything else is downstream of that."*
> — Emily Prime, Apple #3, 2026-06-07

---

**Read this file before starting any work.** Every agent loop in every repo checks this first.
Completed items are marked `[x]` and filed as an Apple to IDUNA before the next item begins.
The backlog is append-only: new items go at the bottom of their section. Emily Prime reviews
and reprioritizes on each major milestone.

**Apple filing requirement:** Every `[x]` completion MUST be accompanied by an Apple posted to
IDUNA (`POST /api/v1/apples`) before the item is considered closed. The Apple is the proof.

---

## SECTION 1: FOUNDATION (current sprint)

- [x] **IDUNA embedded SQLite backend** — Zero-MySQL startup. Bootstrap auto-detects.
  All agents can authenticate without human MySQL provisioning. Apple #1 (Tyler activation).
  Ref: IDUNA commit 8490bb5. DONE 2026-06-07.

- [x] **Tyler agent registration (Build 0016)** — .claude/tyler_agent.md, settings.json,
  emily.sh wired. Tyler authenticated on live IDUNA. Apple #1 filed. Blame path: OCCUPIED.
  DONE 2026-06-07.

- [x] **EMILY/BACKLOG.md created** — This file. Cross-repo golden backlog. Machine-readable.
  Apple #3 (Emily Prime decision) authorized. DONE 2026-06-07.

- [x] **Run emily.sh — first autonomous RSI iteration** — Tyler ran emily.sh. Build 0017 (MPT
  package) + Build 0018 (1952 Detroit annual composite + architecture nodes). 59 BACKLOG items
  complete, 3 remaining (MPT-blocked). Apple #16 (rsi_iteration) filed. DONE 2026-06-07.

- [x] **Wire EMILY RSI engine to IDUNA** — `EMILY/scripts/start-emily-agent.sh` sources
  IDUNA credentials and launches emily-agent with `IDUNA_BASE_URL`, `IDUNA_AGENT_NAME`,
  `IDUNA_AGENT_SECRET` set. cron.go already has Apple submission code (lines 244-257) — it
  fires when IdunaClient is non-nil. Restart running agent using start-emily-agent.sh to enable.
  DONE 2026-06-07.

- [ ] **Process supervision: IDUNA systemd unit** — Deploy the unit from the ops memo
  (TYLER/outlines/emily_iduna_bootstrap.md). IDUNA starts on boot. No human trigger required.
  Acceptance: `systemctl status iduna` shows active, survives reboot, agent auth works after.
  Dependency: server access, emily.sh run ✓ (to verify the env is correct first).

- [x] **Process supervision: emily.sh cron or supervisor** — `TYLER/scripts/cron-emily.sh`
  written. Lock-guarded, IDUNA-check before run, sources Tyler credentials, git pull + 1
  iteration + push. Crontab: `0 */4 * * * /home/fatbaby/TYLER/scripts/cron-emily.sh`.
  Apple #18 filed. DONE 2026-06-07. (Human must add crontab entry.)

---

## SECTION 2: MONEYPRINTERTURBO (video pipeline)

- [x] **MPT config.toml written** — LLM: litellm → anthropic/claude-sonnet-4-6.
  Pexels key placeholder. Build 0017. DONE 2026-06-05.

- [ ] **MPT Pexels API key** — Set YOUR_PEXELS_API_KEY_HERE in MoneyPrinterTurbo/config.toml.
  Human action required (key must come from the user). BLOCKED: waiting on key.

- [ ] **S01E01 cold open compiled clip** — Once MPT is running, run the cold open compilation.
  See TYLER BACKLOG.md for full spec. Dependency: Pexels key ✓, MPT service running.

- [ ] **MPT → TYLER RSI trigger** — When emily.sh produces a new episode script, auto-invoke
  MPT. Spec exists in engine/moneyprinter_pipeline.md. Dependency: cold open clip ✓.

---

## SECTION 3: SYSTEM OBSERVABILITY (Emily Prime sees everything)

- [x] **emily CLI (emily.cli repo)** — Full Go CLI for operator interface. Repo:
  `/home/fatbaby/emily.cli` (github.com/emilyspringerton/emily.cli). v0.3.0 live.
  Commands: `emily observe` (fire observations + auto Apple), `emily apples list|get|post`,
  `emily watch` (IDUNA tail -f), `emily status` (cross-repo git + IDUNA), `emily sync`
  (FatBaby→IDUNA). Binary at `~/.local/bin/emily`. Apple #39 (#40 for RSI cycle).
  DONE 2026-06-07.

- [x] **Apples dashboard query** — `EMILY/scripts/apples.sh` written. Authenticates as
  EMILY-PRIME via M2M, queries `/api/v1/apples`, prints tabular view. Flags: `--full`,
  `--limit=N`, filter by source_repo. Apple #4 filed. DONE 2026-06-07.

- [x] **Cross-repo status script** — `EMILY/scripts/status.sh` written. Git status for all
  5 repos (branch, dirty count, last commit, backlog stats) + last Apple per agent from IDUNA.
  DONE 2026-06-07.

- [x] **PRRJECT_FATBABY → IDUNA Apple filing** — `EMILY/scripts/sync-fatbaby.sh` syncs
  observation files to IDUNA as `signal_observation` Apples. State-tracked (EMILY/var/
  fatbaby-synced.txt) to avoid duplicates. 9 historical observations backfilled, Apples #6–14.
  Run on cron or after each FatBaby cycle. DONE 2026-06-07.

---

## SECTION 4: SHANKPIT / TYLER GAME ENGINE (lower priority)

- [x] **Tyler mode SHANKPIT spec** — `TYLER/engine/tyler_hum_mechanic.md` written. First
  Tyler↔SHANKPIT implementation bridge: HUM_FIELD_ACTIVE server event (6-unit radius, 45s,
  Jiangshi NPC hold + documentation -80%). Go interface sketched. RSI receipts documented.
  Apple #18 filed. DONE 2026-06-07.

- [ ] **SHANKPIT → MPT bridge** — Spec exists (engine/shankpit_mpt_bridge.md). Implementation
  deferred until MPT is running end-to-end.

---

## SECTION 5: FUTURE (Emily Prime decides when to promote)

- [ ] **MySQL wire protocol embedded server** — go-mysql-server as in-process backend.
  Enables: external tooling access, migration path to production MySQL, MySQL-compatible
  backups. Ref: TYLER/outlines/emily_iduna_bootstrap.md Option B.
  Dependency: embedded SQLite stable ✓ (already done), need concrete reason to upgrade.

- [ ] **Tyler IDUNA agent registration via iduna CLI** — `iduna agents register --id tyler`
  when the IDUNA CLI is built. Currently Tyler is seeded via migration. The CLI path is
  the programmatic future.

- [ ] **PRRJECT_FATBABY Emily Prime loop** — Full directed improvement loop where Emily Prime
  reads FatBaby observations, issues improvement tasks, FatBaby executes, loop closes.
  Ref: emily-prime-spec.md. Requires: cross-repo observability ✓, Apple filing ✓.

---

## BACKLOG PROTOCOL

**How to use this file:**
1. Before starting any work, read this file.
2. Pick the highest-priority `[ ]` item in Section 1, then Section 2, etc.
3. Do the work.
4. Post an Apple to IDUNA: `POST /api/v1/apples` (type: completion, title: task name).
5. Check the item `[x]` in this file, add Apple ID and date to the line.
6. `git add BACKLOG.md && git commit -m "backlog: ✓ [task name] — Apple #N"`
7. `git push origin main`
8. Go to step 1.

**How Emily Prime updates this file:**
Emily Prime reads Apples from IDUNA, checks against this file, adds new items when the
system reveals untracked dependencies, and reprioritizes by editing this file and committing.
The backlog is config-as-code. The backlog is what Emily Prime remembers between sessions.

**Golden rule:** No item is complete until the Apple is filed and this file is committed.
The Apple is the proof. The commit is the custody. The push is the delivery.

---

*EMILY PRIME BACKLOG | Cross-repo | Git-authoritative*
*The backlog is what outlasts everything.*
*Clean builds first. Then custody. Then everything else.*
