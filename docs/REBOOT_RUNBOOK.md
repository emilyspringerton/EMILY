# VM Reboot Runbook

Written 2026-07-15 ahead of a planned reboot for pending OS updates
(`/var/run/reboot-required` was set; uptime was 27 days). Updated
2026-07-19 after the next reboot surfaced 8 new systemd units (added
2026-07-17/18, after this doc was first written) that weren't reflected
here yet — see "Update 2026-07-19" below. If you are a fresh Claude Code
session picking this up with no memory of the prior conversation, this
file is self-contained — start here.

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
auto-start now. Only **entity-graph** and **eps-processor** still have no
unit and genuinely need `emily start --all` (or `--entity-graph`
/`--eps-processor`) after a reboot. Run `systemctl --user list-unit-files
--state=enabled` first to see what's already covered before reaching for
`emily start --all`, to avoid the eps-reconciler double-start described
above.

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

Unrelated to the IDUNA/FatBaby chain but also died on reboot — restart if
wanted:
```bash
cd /home/fatbaby/GoblinFoxDragon && go run ./apps2/mud &
```

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
