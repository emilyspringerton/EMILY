# VM Reboot Runbook

Written 2026-07-15 ahead of a planned reboot for pending OS updates
(`/var/run/reboot-required` was set; uptime was 27 days). If you are a
fresh Claude Code session picking this up with no memory of the prior
conversation, this file is self-contained — start here.

## What's now systemd-supervised (auto-starts on boot, no login needed)

Linger is enabled for `fatbaby` (`sudo loginctl enable-linger fatbaby`,
confirmed via `/var/lib/systemd/linger/fatbaby`), and these two user units
are enabled:

- `~/.config/systemd/user/iduna.service` — runs `~/.local/bin/iduna` (:8080).
  Env from `~/.config/iduna/env` (JWT_SECRET, JWT_ISSUER, SERVER_PORT,
  APPLES_GIT_DIR — already populated).
- `~/.config/systemd/user/emily-system.service` — oneshot, runs
  `emily start --iduna` (After/Wants iduna.service). Brings up
  `observation-watcher` and the `emily-agent` RSI daemon.

These three (IDUNA, observation-watcher, emily-agent) should come back on
their own after reboot. Verify with:
```bash
export XDG_RUNTIME_DIR=/run/user/1000 DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/1000/bus
systemctl --user status iduna.service emily-system.service
curl -s localhost:8080/health
```

## What does NOT auto-start — restart manually

These were running as bare `go run`/binaries (PPID 1, launched ad hoc, no
unit, no `emily start` flag covers them all) as of the last snapshot:

```bash
cd /home/fatbaby/PRRJECT_FATBABY

# newssite (:8082) + signalapi (:9091) — covered by `emily start --all`,
# but confirm/run explicitly if that doesn't pick them up:
go run ./cmd/newssite -store var/secwatch -graph-dir var/entity-graph \
  -eps-dir var/eps -commentary-dir var/commentary -guidance-dir var/guidance \
  -earnings-cal-dir var/earnings-calendar &

go run ./cmd/signalapi -addr :9091 -store var/secwatch &

# entity-graph, eps-reconciler, eps-processor — no unit, no emily start flag.
# Always started by hand; do so again:
go run ./cmd/entity-graph -store ./var/secwatch &
go run ./cmd/eps-reconciler -store ./var/secwatch -eps-dir ./var/eps &
go run ./cmd/eps-processor -body-store ./var/prwatch-body \
  -discovery-store ./var/prwatch -eps-dir ./var/eps &
```

Not running before the reboot (per PRRJECT_FATBABY/CLAUDE.md's fuller
process table) and presumed intentionally stopped — confirm with the user
before starting: `secwatch`, `prwatch`, `prwatch-body`, `dashboard`,
`feedserver`, `broker`, `guidance-watcher`, `jon-agent`, `form4-watcher`,
`dividend-watcher`, `buyback-watcher`, `nt-watcher`, SHANKPIT's
`shank_go_server` + emily-bots.

Unrelated to the IDUNA/FatBaby chain but also died on reboot — restart if
wanted:
```bash
cd /home/fatbaby/GoblinFoxDragon && go run ./apps2/mud &
```

## Startup order

1. Confirm boot, confirm `systemctl --user` reachable for fatbaby (proves
   linger survived the update/reboot — re-run `enable-linger` if not).
2. `iduna.service` should already be active — `curl localhost:8080/health`.
3. `emily-system.service` (oneshot) should already have run — check
   `var/logs/observation-watcher.log` and `var/logs/emily-agent.log` in
   `EMILY/` for clean startup.
4. `emily start --iduna --all` — idempotent, picks up newssite/signalapi if
   step 3 didn't already.
5. Manually start entity-graph / eps-reconciler / eps-processor per above.
6. Verify ports: `ss -tlnp | grep -E ':8080|:8082|:9091'`.
7. Tail each `var/logs/*.log` for startup errors.
8. `emily status --fatbaby` and `emily apples list --limit 5` — confirm the
   Apple round-trip works end-to-end through the freshly booted IDUNA.
9. Wait one 5-min emily-agent RSI cycle, confirm a new Apple appears.
10. Resume normal Emily Way operation — read `EMILY/BACKLOG.md`.

## Known pre-existing issue (not reboot-related)

`PITVIPER`'s `origin` remote (`github.com/emilyspringerton/PITVIPER.git`)
returns "Repository not found" on fetch, and its `main` branch has no
upstream configured. No local changes are at risk (working tree was clean
at last check), but it isn't backed up to GitHub. Needs a human decision
(recreate the repo? wrong URL? access issue?) — not fixed as part of this
runbook.
