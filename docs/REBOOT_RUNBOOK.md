# VM Reboot Runbook

Written 2026-07-15 ahead of a planned reboot for pending OS updates
(`/var/run/reboot-required` was set; uptime was 27 days). Updated
2026-07-19 after the next reboot surfaced 8 new systemd units (added
2026-07-17/18, after this doc was first written) that weren't reflected
here yet — see "Update 2026-07-19" below. If you are a fresh Claude Code
session picking this up with no memory of the prior conversation, this
file is self-contained — start here.

**Update 2026-08-10, ahead of a planned box upgrade (founder: "we will
upgrade the box... but just in case document everything meticulously
before you go down"):** re-verified against live `systemctl --user
list-unit-files --state=enabled` rather than trusting the 2026-07-19
snapshot below, which had drifted significantly — 21 more units exist now
that this doc never listed (REDGARDEN's R&D+stable split, GoblinFoxDragon,
EINHORN_SURVIVAL, SHANKPIT-460, WEAKNIGHT_RACERS, and several more
`fatbaby-*` watchers), and two of the doc's own claims were stale
(`fatbaby-entity-graph.service` is enabled now, not disabled;
`fatbaby-eps-processor.service` has a real unit now, contradicting "no
unit and genuinely needs `emily start --all`"). See the fully rewritten
"What's now systemd-supervised" list immediately below — treat it as the
current source of truth, not the history above it.

## What's now systemd-supervised (auto-starts on boot, no login needed)

Linger is enabled for `fatbaby` (`sudo loginctl enable-linger fatbaby`,
confirmed via `/var/lib/systemd/linger/fatbaby`).

Core chain (documented since 2026-07-15):
- `~/.config/systemd/user/iduna.service` — runs `~/.local/bin/iduna` (:8080).
  Env from `~/.config/iduna/env` (JWT_SECRET, JWT_ISSUER, SERVER_PORT,
  APPLES_GIT_DIR — already populated).
- `~/.config/systemd/user/emily-system.service` — oneshot, runs
  `emily start --iduna` (After/Wants iduna.service). Brings up
  `observation-watcher` and the `emily-agent` RSI daemon.

**Update 2026-07-19:** BACKLOG SECTION 152 (S152-03, added 2026-07-17)
gave 7 more FatBaby processes their own supervised units, after secwatch
and eps-reconciler were both silently OOM-killed and stayed down for
hours with zero auto-recovery. All are `Type=simple`, `Restart=on-failure`,
`WantedBy=default.target` (so they auto-start same as iduna/emily-system):
`fatbaby-secwatch.service`, `fatbaby-prwatch.service`,
`fatbaby-prwatch-body.service`, `fatbaby-processor.service` (`cmd/processor`
— distinct binary from `cmd/eps-processor`, which has no unit of its own),
`fatbaby-eps-reconciler.service`, `fatbaby-newssite.service`,
`fatbaby-signalapi.service`. A matching `fatbaby-entity-graph.service` file
exists too but is currently **disabled** (not yet cut over) — entity-graph
still needs to be started the old way (`emily start --all` / `go run`).

Also added 2026-07-18: `shankpit460-server.service` — first-ever supervised
run of the SHANKPIT-460 game server (UDP :6969, `Restart=on-failure`). This
auto-starts on every boot now, unconditionally — see the "confirm before
starting" caveat below, which this unit now bypasses.

**Conflict this creates:** step 4 below (`emily start --iduna --all`) still
tries to launch `eps-reconciler` (and would try newssite/signalapi/secwatch
too if their pgrep guards ever mismatch the systemd-launched command line).
Confirmed 2026-07-19: it silently spawned a second, unsupervised `go run
./cmd/eps-reconciler` process alongside the systemd-managed one — both
polling/writing the same event store. **Before running `emily start --all`
after a reboot, check `systemctl --user status fatbaby-*.service` first**
and kill any resulting duplicate (compare command lines — the systemd one
runs the prebuilt `bin/<name>` binary with the full flag set; a duplicate
from `emily start` runs `go run ./cmd/<name>` with a shorter flag set).

These should all come back on their own after reboot. Verify with:
```bash
export XDG_RUNTIME_DIR=/run/user/1000 DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/1000/bus
systemctl --user status iduna.service emily-system.service
systemctl --user list-unit-files --state=enabled   # full current picture, not just the core two
curl -s localhost:8080/health
```

**Full current enabled-unit list (2026-08-10, re-verified live — this supersedes every partial
list above in this doc, which are kept only as history of how we got here):**

| Unit | What it is |
|---|---|
| `iduna.service` | IDUNA IAM (:8080), `Restart=on-failure`, has an `ExecStartPost` health check |
| `emily-system.service` | oneshot — `emily start --iduna --agi`, brings up observation-watcher + the emily-agent RSI daemon (see "Known gaps" below: NOT individually restart-supervised after that) |
| `fatbaby-secwatch.service`, `fatbaby-prwatch.service`, `fatbaby-prwatch-body.service`, `fatbaby-processor.service`, `fatbaby-eps-reconciler.service`, `fatbaby-eps-processor.service`, `fatbaby-newssite.service`, `fatbaby-signalapi.service`, `fatbaby-entity-graph.service`, `fatbaby-form4-watcher.service`, `fatbaby-guidance-watcher.service`, `fatbaby-dividend-watcher.service`, `fatbaby-buyback-watcher.service`, `fatbaby-nt-watcher.service`, `fatbaby-schd13-watcher.service`, `fatbaby-market-data-watcher.service`, `fatbaby-earnings-calendar.service`, `fatbaby-pr-reaction-watcher.service`, `fatbaby-pr-indexer.service` (added 2026-08-13, closes the real "press releases never wrote source_document_persisted" gap — see PRRJECT_FATBABY CHANGELOG same date) | every PRRJECT_FATBABY watcher/processor now has its own unit — `Type=simple`, `Restart=on-failure` |
| `fatbaby-bond-watcher.timer`, `fatbaby-movers-watcher.timer` | timer-triggered FatBaby jobs |
| `einhorn-survival.service` | the real community Minecraft server (Paper, :25565) |
| `gfd-mud.service` | GoblinFoxDragon's DragonsNShit MUD (`:2323`, `-api-port 7171`) — **the old "died on reboot, restart with `go run ./apps2/mud &`" instruction below this table is stale, left only as history; this is a real supervised unit now** |
| `gfd-server-go.service` | GoblinFoxDragon's real Go UDP backend (`-udp-port 6970 -worldapi-port 7070 -trapx-port 7071`) |
| `dragonfly-debug.service` | Dragonfly/Bedrock fork debug instance |
| `redgarden-matchmaker-bots.service` (:7778), `redgarden-matchmaker-players.service` (:7779), `redgarden-bot-pool.service` | REDGARDEN's **R&D** deployment (2026-08-10 split — see `REDGARDEN/CLAUDE.md`'s own "Deployments" table) |
| `redgarden-stable-matchmaker-bots.service` (:8778), `redgarden-stable-bot-pool.service` | REDGARDEN's **stable** deployment, dedicated to GoblinFoxDragon's Battlegrounds — separate checkout (`/home/fatbaby/redgarden-stable`), NOT touched by `redgarden-auto-deploy.timer` |
| `redgarden-auto-deploy.timer` | polls CI, redeploys the **R&D** REDGARDEN units only when a new green build lands — never the stable ones |
| `shankpit460-server.service`, `shankpit460-emily-bot.service` | SHANKPIT-460 game server (UDP :6969) + its bot |
| `weaknight-racers-server.service` | WEAKNIGHT_BEDROCK_RACERS server |
| `session-migration.service` | (check its own unit file if unfamiliar — not otherwise documented here) |
| `gpt2-serve.service` | gpt2-alpine-c inference server (:8088, emily-ft checkpoint) — added 2026-08-13, see "processify" note below, was manual-restart-only before this |
| `fatbaby-broker.service` | tenant-aware proxy fronting gpt2-serve.service (:8679) — added 2026-08-13, source is `PRRJECT_FATBABY/cmd/broker` despite the `gpt2-alpine-c`-looking name history below |

Confirm linger + this full list survived the upgrade with one command:
```bash
loginctl show-user fatbaby | grep -i linger   # expect: Linger=yes
systemctl --user list-unit-files --state=enabled
```

## What does NOT auto-start — restart manually

**Update 2026-07-15:** `emily start`'s `--all` flag now covers newssite,
signalapi, entity-graph, eps-reconciler, and eps-processor (previously
only newssite/signalapi; entity-graph/eps-reconciler/eps-processor had no
flag at all and had to be started by hand with no idempotency guard —
fixed in emily.cli commit history, see CHANGELOG). `--all` deliberately
does **not** include `shank_go_server` (SHANKPIT) any more — it used to,
which meant a plain `emily start --iduna --all` would silently launch a
live game server + fill bots. That's now gated only by `--shankpit`.

**Update 2026-07-19:** most of what this section used to call "manual" is
now systemd-supervised instead (see S152-03 above) — newssite, signalapi,
secwatch, prwatch, prwatch-body, processor, and eps-reconciler all
auto-start now.

**Update 2026-08-10, correcting the above — re-verified live, both claims in the previous
paragraph are now stale:** `fatbaby-entity-graph.service` and `fatbaby-eps-processor.service`
are BOTH enabled, real units (`systemctl --user is-enabled` confirms both) — this section's own
2026-07-19 claim that they "still have no unit and genuinely need `emily start --all`" is wrong
as of today. The real remaining gaps are `gpt2-alpine-c`'s `serve.py` and `cmd/broker`, covered
in their own section below — not entity-graph/eps-processor. Always run `systemctl --user
list-unit-files --state=enabled` first to see current real coverage before reaching for `emily
start --all`, same "don't trust a doc, check live" discipline this whole update applies
throughout.

```bash
emily start --iduna --all
```

This is idempotent (pgrep-guarded per process) — safe to re-run — *except*
where a systemd-launched process's command line doesn't match the pgrep
pattern `emily start` expects, which is exactly what happened with
eps-reconciler on 2026-07-19. Verify with `pgrep -af` after running, not
just by trusting "idempotent."

Also now auto-starting via `shankpit460-server.service` (2026-07-18) and
**no longer gated by a confirm-first check**: SHANKPIT's game server is up
on every boot unconditionally. If that's not wanted, `systemctl --user
disable --now shankpit460-server.service`.

Still not running before the reboot (per PRRJECT_FATBABY/CLAUDE.md's fuller
process table) and presumed intentionally stopped — confirm with the user
before starting: `dashboard`,
`feedserver`, `broker`, `guidance-watcher`, `jon-agent`, `form4-watcher`,
`dividend-watcher`, `buyback-watcher`, `nt-watcher`. (SHANKPIT's fill-bots
are still manual — `emily start --shankpit` — even though the game server
itself now auto-starts via `shankpit460-server.service`.)

**Update 2026-08-13 — both of the below are now real, enabled systemd units. Processified during
post-upgrade reboot recovery (founder real-time: "we need to processify that") — this whole
subsection is now history, not an active gap.** `gpt2-serve.service` (source:
`gpt2-alpine-c/ops/systemd/gpt2-serve.service`) and `fatbaby-broker.service` (source:
`PRRJECT_FATBABY/ops/systemd/fatbaby-broker.service`) — both `Type=simple`, `Restart=on-failure`,
`WantedBy=default.target`, deployed the same way as every other unit in this doc (`cp` to
`~/.config/systemd/user/`, `daemon-reload`, `enable --now`). Add both to the "What's now
systemd-supervised" table above on the next full re-verification pass.

One correction while fixing this: this doc previously said `cmd/broker` lives in
`gpt2-alpine-c` — wrong. The actual Go source is `PRRJECT_FATBABY/cmd/broker`;
`gpt2-alpine-c/config/broker-routes.json` is only the routing config (tenant_id "emily-prime",
upstream `http://localhost:8088`). `fatbaby-broker.service` builds/runs from PRRJECT_FATBABY and
points `--routes` at the gpt2-alpine-c config file.

**Known issue carried forward, not fixed by processifying:** `cmd/broker/main.go` hardcodes a 30s
`ResponseHeaderTimeout`/context timeout on the upstream call. `serve.py`'s real `/generate`
latency has been observed at ~8 minutes on a cold request (stdlib `http.server`, single-threaded,
not `ThreadingHTTPServer`) — so the broker path 502s ("upstream unavailable") on any real
generation call today. Confirmed via direct smoke test 2026-08-13, Apple #13254. Both the
timeout mismatch and the single-threaded server are real, separate follow-up work — not touched
here since fixing them wasn't in scope for reboot recovery / processifying.

Original 2026-08-10 manual-restart note, kept as history:
- **`gpt2-alpine-c/scripts/serve.py` (GPT-2 inference server, :8088)** — no systemd unit exists.
  Restart manually: `cd /home/fatbaby/gpt2-alpine-c && python3 scripts/serve.py --model ft --port 8088 &`
  (or nohup/disown it properly — don't leave it attached to a login shell).
- **`cmd/broker` (gpt2-alpine-c's routing broker, :8679)** — no systemd unit exists. Restart:
  `cd /home/fatbaby/gpt2-alpine-c && go run ./cmd/broker --routes config/broker-routes.json --addr :8679 --store /home/fatbaby/EMILY/var/broker-events &`

The GFD MUD instruction that used to be here (`cd GoblinFoxDragon && go run ./apps2/mud &`) is
now WRONG — see `gfd-mud.service` in the table above, it's a real supervised unit as of
2026-08-06 and needs no manual restart.

## Reboot-as-deploy caution

Every auto-started and manually-restarted process here runs via `go run`,
which compiles HEAD at launch time — this was already true day-to-day, but
a reboot means several weeks of accumulated, never-yet-executed commits
all go live simultaneously for the first time (as of 2026-07-15: EMILY and
PRRJECT_FATBABY both had a large drift between what the running binaries
were built from and current HEAD). Watch the logs closely for the first
RSI cycle or two after a reboot rather than assuming a clean process start
means the code is behaving as expected.

Separately: `~/.local/bin/iduna` is a **prebuilt** binary (not `go run`),
so it does NOT pick up IDUNA HEAD automatically — it keeps running
whatever was last built with `go build -o ~/.local/bin/iduna .`. If IDUNA
source has moved on since that build (check `git -C /home/fatbaby/IDUNA
log -1` vs the binary's mtime), decide deliberately whether to rebuild +
restart it — don't do this bundled with the reboot itself, since IDUNA
runs DB migrations at startup and you want to isolate any migration
failure from "did the reboot work" as a separate question.

## Startup order

1. Confirm boot, confirm `systemctl --user` reachable for fatbaby (proves
   linger survived the update/reboot — re-run `enable-linger` if not).
2. `iduna.service` should already be active and *healthy*, not just
   started — `curl -sf localhost:8080/health` (the unit now has an
   `ExecStartPost` health-check loop, so `systemctl --user is-active`
   alone means the health check already passed; if the unit is stuck
   "activating" for >30s the health check is failing — check
   `journalctl --user -u iduna.service`).
3. `emily-system.service` (oneshot) should already have run. **Don't just
   trust `systemctl --user status` being "active (exited)"** — `emily
   start` can fail to launch a child and still exit 0 in older builds; as
   of this fix it now exits 1 if anything failed, but verify directly:
   ```bash
   pgrep -af 'cmd/observation-watcher --root'
   pgrep -af 'emily-agent.*--daemon'
   ```
   and check `var/logs/observation-watcher.log` / `var/logs/emily-agent.log`
   in `EMILY/` for clean startup.
4. `emily start --iduna --all` — idempotent, picks up newssite / signalapi /
   entity-graph / eps-reconciler / eps-processor if step 3 didn't already.
   Safe to run even if everything's already up.
5. Verify ports: `ss -tlnp | grep -E ':8080|:8082|:9091'`. (Nothing on
   :8086 is expected — `emily-agent --daemon` returns before ever starting
   its HTTP server; that port is only live in the non-daemon interactive
   mode, which `emily start` doesn't use.)
6. Tail each `var/logs/*.log` for startup errors.
7. `emily status --fatbaby` and `emily apples list --limit 5` — confirm the
   Apple round-trip works end-to-end through the freshly booted IDUNA.
8. Wait one 15-min emily-agent RSI cycle, confirm a new Apple appears.
9. Check the APPLES git mirror actually received it, not just IDUNA's db:
   ```bash
   git -C /home/fatbaby/APPLES log -1 --format='%ci %s'
   ```
   (APPLES git push failures are currently log-only in IDUNA — a stale
   mirror wouldn't otherwise be obvious.)
10. Check for truncated NDJSON writes per the state-hydration note below.
11. Resume normal Emily Way operation — read `EMILY/BACKLOG.md`.

Don't run `start.sh` and `run.sh` concurrently against a cold boot — both
launch `emily start`-adjacent work and, while now pgrep-guarded, there's
no reason to race them.

## State hydration note

No separate data/state hydration step is needed. Everything meaningful —
git repos, IDUNA's `var/iduna.db`, PRRJECT_FATBABY's NDJSON event stores,
memory files, CLAUDE.md, BACKLOG.md — is disk-resident and survives a
plain reboot untouched. Only in-memory process state is lost (covered
above) and any live Claude Code session's context (covered by `start.sh`,
which hands a fresh session this runbook instead of expecting it to
already know what happened).

One thing worth checking since the event-store processes get hard-killed
by the reboot rather than shut down cleanly: confirm the append-only
NDJSON stores don't have a truncated trailing line from a write in
progress at kill time.
```bash
for f in $(find /home/fatbaby/PRRJECT_FATBABY/var/secwatch -name '*.ndjson' -newer /home/fatbaby/PRRJECT_FATBABY/go.mod 2>/dev/null); do
  tail -c1 "$f" | od -c | grep -q '\\n' || echo "possible truncated write: $f"
done
```
If one turns up, drop the partial last line — the event store's
sequence numbering is monotonic and expects whole records only.

## Resolved issues from the earlier audit (2026-07-15)

- `PITVIPER`'s remote was broken ("Repository not found") — the user fixed
  it by reconnecting the remote; `git fetch` now succeeds and `main` tracks
  `origin/main` cleanly. A stray empty nested clone at
  `PITVIPER/PITVIPER/.git` (0 objects, pure debris) was also removed.
- Deep audit (Fable) surfaced and fixed: `--all` accidentally starting
  SHANKPIT's game server, `emily start` always exiting 0 regardless of
  child failures, an overly-broad pgrep pattern for observation-watcher
  that could match a `tail -f` of its own log file, a bare `emily` exec in
  obs-watcher's failure-escalation path that could silently no-op under
  systemd's restricted PATH, IDUNA's world-readable JWT secret file
  (0664 → 0600), and start.sh/run.sh being group-writable (0775 → 0755).
  entity-graph/eps-reconciler/eps-processor gained proper `emily start`
  integration (pgrep idempotency, `--all` coverage) instead of being fully
  manual with no guard against double-starting and double-writing to the
  append-only event stores.

## Update 2026-09-05 — full home-config loss, not just a normal reboot

A reboot this day was preceded by `~/.ssh`, `~/.gitconfig`, `~/.config/systemd/user`,
`~/.local/bin`, and `~/go.work` all turning up **empty/missing** — not a normal reboot recovery,
since none of the "should auto-start" units above had any unit file left to load. All *data*
(git repos, `IDUNA/var/*.db`, `IDUNA/iduna-key.json` — the actual ES256 signing keypair) survived
untouched; only home-scoped config/dotfiles were gone. Checked for intrusion (users, authorized_keys,
cron, listening ports) — all clean; root logins all traced to one consistent residential/mobile IP.
Read as an infra/home-reset event, not a compromise, but not root-caused with certainty from inside
the box alone.

Real new gotchas found doing full recovery from this state, not covered above:
- **`~/go.work` is NOT git-tracked** (confirmed via `git log --all -- go.work`) — it's a
  machine-local convenience file, gone for good on a home reset with no git backup. Recreate with
  `go work init ./EMILY/emily-agent ./EDIS ./IDUNA ./PRRJECT_FATBABY ./emily.cli ./gpt2-alpine-c`
  — note `gpt2-alpine-c` is a 6th real member, needed because `EMILY/emily-agent/enrichworker.go`
  imports `gpt2-alpine-c/pkg/towerprint`; CLAUDE.md's own "5 modules" table predates this and is
  stale. Without it, `emily-agent`'s daemon fails to `go run` at all with a confusing "not in std"
  error, which reads like a code bug but is purely a missing workspace file.
- **The REDGARDEN matchmaker/bot-pool unit files exist in 5 places** (`REDGARDEN/`, `ECOWAR/`,
  `redgarden-deploy/`, `redgarden-stable/`, `fix/REDGARDEN/`) — only two are real: R&D units run
  from `/home/fatbaby/REDGARDEN` itself, stable units from `/home/fatbaby/redgarden-stable` (check
  each unit's own `WorkingDirectory=` to confirm, don't guess). `redgarden-deploy` is only
  `auto_deploy.sh`'s own private CI-polling checkout (see its `DEPLOY_DIR` var) and never serves
  live traffic; `ECOWAR`'s copies are fork leftovers (it has its own distinct
  `ecowar-matchmaker`/`ecowar-bot-pool` units instead); `fix/REDGARDEN` is a scratch checkout.
- **`gpt2-serve.service` needs Python `transformers`/`torch`** installed somewhere on the box (no
  `requirements.txt` exists to pin versions) — if `~/.local` is what's wiped, these likely lived in
  user site-packages there. Not reinstalled blind during this recovery (unclear size/version
  pinning) — a deliberate, separate decision, not bundled into "just get services back."
- **`~/.config/iduna/env` and `~/.config/idunapro/env` contents are NOT recoverable if lost** —
  `JWT_SECRET`/`JWT_ISSUER`/`BASE_URL` have no backup anywhere found. Safe to regenerate
  `JWT_SECRET` fresh (only used for the narrow device-auth flow) as long as the real ES256 keypair
  (`IDUNA/iduna-key.json` / IDUNA_PRO's own equivalent) is confirmed intact first — that's the
  actual trust root, not the env file.
- **Git push breaks silently and system-wide if `~/.ssh` is wiped** — no private key, no
  `known_hosts`. `ssh-keyscan -t ed25519,rsa github.com >> ~/.ssh/known_hosts` fixes host
  verification, but a lost private key needs a **new keypair generated and its public half added
  to GitHub by a human** — not something recoverable from inside the box. This also silently broke
  IDUNA's own kanban-git auto-sync and `redgarden-auto-deploy.timer` (both retry-loop and fail
  until the new key is registered).

## Known gaps, deliberately not auto-fixed

- **Alerting env vars aren't wired into the systemd path.** `EMILY/var/emily-secrets.env`
  (the only file `emily start`'s children read via `config.Resolve()`)
  contains only `ANTHROPIC_API_KEY`. The S131 Slack/check-in alerting reads
  `SLACK_WEBHOOK_URL`, `SLACK_DEFAULT_CHANNEL`, `GMAIL_*`, `FCM_*` from env,
  none of which are supplied to `emily-agent` when launched via
  `emily-system.service`. If those alerts matter, add them to
  `emily-secrets.env` or an `EnvironmentFile=` on the unit — not done here
  since the actual values weren't available to set blind.
- **No restart-on-crash for observation-watcher/emily-agent.** Only
  `iduna.service` has `Restart=on-failure`; the RSI daemon and obs-watcher
  are launched once by `emily-system.service`'s oneshot `ExecStart` and
  then live outside systemd's supervision — if either crashes at 3am,
  nothing brings it back until someone runs `emily start` again. Converting
  them to their own `Type=simple`/`Restart=on-failure` units would close
  this, but it's a bigger change to the current "one oneshot launcher"
  model than warranted for this reboot — flagging as a follow-up.
- **IDUNA binary/HEAD drift** — see "Reboot-as-deploy caution" above;
  left as a deliberate human decision, not bundled into reboot recovery.
