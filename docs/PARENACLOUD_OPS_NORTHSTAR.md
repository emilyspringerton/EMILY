# PARENACLOUD — a standard ops-playbook format for backend services

Founder real-time, this session: "lets also start to create a standard format for the backend
services in terms of like systemd? like for example if someone wanted to start their own fork of
REDGARDEN on their own server there should be part of the command line interface that we
discussed that pulls the repo reads the ops playbook or whatever and actually sets up those
backend services (via systemd or whatever) it would be good to have a standard format - mention
but defer work on additional dependent services setup via ansible and terraform." Direct
continuation of `PARENACLOUD_NORTHSTAR.md`'s own PC2 ("build the `GameConfig` schema for real")
and the CLI sketch's `parena cloud deploy` entry, now made concrete.

## Real, existing state — checked directly across every real deploy, not assumed

Every game repo in the C/SDL2 arena lineage (`REDGARDEN`, `ECOWAR`, `BRAWLPIT`, `PAPERCRAFT`,
`WEAKNIGHT_BEDROCK_RACERS`) already has its own `ops/systemd/*.service` files — real, live,
hand-written unit files, not a hypothetical. Reading them side by side surfaces exactly
`PARENGINE_NORTHSTAR.md`'s own "copy-pasted by hand" finding playing out a third time, now on the
deploy side: `ECOWAR/ops/systemd/` still contains files literally named
`redgarden-matchmaker-bots-3v3.service`, `redgarden-bot-pool.service`, etc. — REDGARDEN's own
units, hand-copied wholesale into ECOWAR's own directory and never renamed, sitting alongside the
two that WERE correctly adapted (`ecowar-matchmaker.service`, `ecowar-bot-pool.service`).

But the hand-written files that ARE correct are also **real, valuable prior art** — every one of
them independently converged on the same real shape:

- User-level systemd throughout (`systemctl --user`, no sudo) — a real, consistent, deliberate
  convention across every single example, not a coincidence.
- `[Service]`: `Type=simple`, a `WorkingDirectory`, an `ExecStart` (binary + flags), `Restart=
  on-failure` + `RestartSec=5s`, `StandardOutput`/`StandardError` appended to a `var/logs/*.log`
  file, a `MemoryMax`.
- An optional `EnvironmentFile` for injecting per-install secrets (an IDUNA M2M agent credential,
  almost always) — real, found-the-hard-way significance: `ECOWAR/ops/systemd/
  ecowar-matchmaker.service`'s own header comment tells the real story of a required (no leading
  `-`) `EnvironmentFile=` on a missing file taking down the ENTIRE unit's activation for 5 real
  days, not just the feature the file was actually for (WOTAN match-result reporting).
- A dependency chain between a persistent matchmaker and its own bot-pool (`After=`/`Wants=`
  pointing at the matchmaker unit, an `ExecStartPre=/bin/sleep 2` before the pool starts).
- A deploy runbook hand-typed as a comment block at the top of every single file (`mkdir -p
  ~/.config/systemd/user && cp ... && systemctl --user daemon-reload && enable --now`) — this is
  already, functionally, a hand-written "ops playbook" per file; it's just never been factored
  into one real, machine-readable source those comments could instead be generated FROM.

**A real, live, already-existing risk found by this exact comparison, not fixed here**:
`REDGARDEN/ops/systemd/redgarden-matchmaker-bots.service` still has a **required** (no leading
`-`) `EnvironmentFile=` line — the identical landmine `ECOWAR`'s own incident already taught this
org to avoid, sitting live in REDGARDEN's own current production config, undiscovered until this
side-by-side read surfaced it. Named here, not silently patched — REDGARDEN's own live matchmaker
is a shared, currently-running service (see this session's own standing memory on never touching
one of those without asking first); flagging it is this doc's job, fixing it is a real, separate,
explicitly-confirmed action.

**Two real, documented incidents this exact standardization would have structurally prevented**,
both already in this monorepo's own git history: `BRAWLPIT/ops/systemd/brawlpit-server.service`'s
own header comment describes its server binary hardcoding port 6969 — the exact port SHANKPIT's
own live process already permanently holds on the same host — with no port-collision check at
all, silently swallowed for an unknown period before a real investigation found it.
`PAPERCRAFT/ops/systemd/papercraft-server.service` names the same class of incident ("the real
CI-released client works, but no server was ever actually running for it to connect to"). A
format that makes every service's own port a real, structured, checkable field — not a number
buried inside a free-text `ExecStart=` line — is what turns "silently swallowed until someone
investigates" into "the deploy tool refuses before it happens."

## The real format: `ops/playbook.toml`

TOML, living at `ops/playbook.toml` in each game repo — the same `ops/` directory every one of
these repos already uses for its hand-written units, so a maintainer looking for deploy config
finds it exactly where they'd already look. Two top-level sections plus a repeatable service
block:

```toml
# ops/playbook.toml — <GAME> backend services. Generated units are written by the deploy tool,
# never hand-edited (same "do not hand-edit, a re-provision will overwrite this" discipline
# IDUNA_PRO's own internal/tenantprovision.writeSystemdUnit already establishes).

[game]
name = "ecowar"        # unit-name prefix and default log-path namespace
repo = "ECOWAR"         # informational only

# One [[service]] block per real, independent process. `kind` is the one field every reader
# should check first:
#   "persistent" -- a real, always-on systemd unit (a matchmaker, a bot pool, a single-process
#                    game server like BRAWLPIT's).
#   "spawned"    -- NOT a systemd unit at all. Documents a binary a persistent service spawns
#                    itself at runtime (arena_server, spawned per-match by the matchmaker's own
#                    --server-bin flag) -- real, deliberate: giving this its own systemd unit
#                    would be structurally wrong (systemd can't supervise "however many of these
#                    exist at once, however long each one lives"), but leaving it out of the
#                    playbook entirely would make the playbook an incomplete map of what this
#                    game's backend actually consists of.
[[service]]
name = "matchmaker"
kind = "persistent"
description = "ECOWAR arena matchmaker — 1v1 queue with bots"
binary = "build/red_garden_matchmaker"
args = ["--listen-port", "{{listen_port}}", "--lobby-size", "{{lobby_size}}",
        "--server-bin", "{{server_bin}}", "--first-game-port", "{{first_game_port}}"]
port = "{{listen_port}}"          # real, structured -- checkable before deploy, not buried in args
environment = { REDGARDEN_TICKET_SECRET = "test-secret-for-vs0-vs1-validation" }
environment_file = { path = "var/{{env_file_name}}", required = false }  # the ECOWAR fix, as the default
restart = "on-failure"
restart_sec = 5
memory_max = "256M"
after = ["network-online.target"]

[[service]]
name = "bot-pool"
kind = "persistent"
description = "ECOWAR arena bot pool — fills the queue's open human slot(s)"
binary = "scripts/run_bot_pool.sh"
args = ["{{bot_count}}", "{{listen_port}}"]
environment = { REDGARDEN_TICKET_SECRET = "test-secret-for-vs0-vs1-validation" }
after = ["network-online.target", "ecowar-matchmaker.service"]
wants = ["ecowar-matchmaker.service"]
start_delay_sec = 2               # renders as ExecStartPre=/bin/sleep 2
restart = "on-failure"
restart_sec = 5
memory_max = "128M"

[[service]]
name = "arena-server"
kind = "spawned"
description = "Per-match arena server, spawned on demand by matchmaker's own --server-bin"
binary = "build/red_garden_arena_server"

[values]
# Real, per-install defaults for every {{placeholder}} above. A second install on the same box
# (or a genuinely different fork/tenant) overrides only what's different via a small
# `--values-file` passed to the deploy tool, never by hand-editing the playbook itself.
listen_port = 9779
lobby_size = 2
server_bin = "build/red_garden_arena_server"
first_game_port = 9600
bot_count = 1
env_file_name = "ecowar-iduna-agent.env"
```

Real fields not shown above, with real defaults: `user_level` (bool, default `true` — every real
example in this monorepo is user-level, so that's the default rather than something every
playbook has to state); `log_path` (default `var/logs/<service-name>.log`, matching every real
example's own convention).

**Two real, concrete wins this format gets for free, both grounded directly in incidents already
in this monorepo's own history**: (1) a real port-collision check becomes possible before deploy
(`port` is a structured field, not free text inside `args`) — the actual, structural fix for the
BRAWLPIT/PAPERCRAFT incident class named above; (2) REDGARDEN's own real 3v3/"stable"-branch/
players-only unit variants (`redgarden-matchmaker-bots-3v3.service`,
`redgarden-stable-matchmaker-bots.service`, `redgarden-matchmaker-players.service` — real, live,
currently hand-maintained as separate near-duplicate files) collapse to ONE playbook plus small,
named `--values-file` overrides (`lobby_size = 6` for 3v3, a different `listen_port`/working
directory for the stable branch) instead of N independently-drifting copies of the whole unit
file.

## Real, working implementation: the renderer

Built and tested this session, not just specified: `emily.cli`'s new `internal/opsplaybook`
package (`ParsePlaybook`, real TOML parsing via `github.com/BurntSushi/toml`) + `RenderUnit`
(playbook service + values → real systemd unit file text, matching the exact `[Unit]`/`[Service]`/
`[Install]` shape every hand-written example above already uses) + a new `emily ops render-unit`
subcommand. Real, live-verified: rendering `ECOWAR/ops/playbook.toml`'s own `matchmaker` service
against its own `[values]` produces a unit functionally equivalent to the real, currently-deployed
`ecowar-matchmaker.service` (same `ExecStart` command line, same `EnvironmentFile=-` optional-file
convention, same `Restart=`/`MemoryMax=`/log-path shape) — checked field by field in
`internal/opsplaybook/opsplaybook_test.go`, not eyeballed.

**Real, deliberately NOT built in this pass**: the "pulls the repo" half of the founder's own
ask — cloning/fetching a fork, running its build script, actually copying rendered units into
`~/.config/systemd/user/` and running `daemon-reload`/`enable --now` against a real host. That's
real, live infrastructure-mutation work (the same category of action this session has held back
from without explicit confirmation throughout — see the standing "never restart a shared
matchmaker without asking" precedent) and belongs behind the actual `parena cloud deploy` /
whichever CLI surface `PARENACLOUD_NORTHSTAR.md`'s own still-open PC1 architecture question
resolves to, not bolted onto a rendering library as an afterthought. The renderer is real,
tested, and ready for that CLI to call into once it exists.

**Real, honest architecture note found while placing this code**: `PARENACLOUD_NORTHSTAR.md`'s
own CLI sketch recommended extending the `parena` binary itself with a `cloud` subcommand group —
but `parena` is a C program (`src/main.c`). TOML parsing, systemd interaction, and future HTTP
orchestration against IDUNA_PRO's own real tenant-provisioning API are all substantially easier in
Go, and `emily.cli` is already this monorepo's own real, working, Go ops CLI with exactly this
kind of task as its established job. This doc's own renderer lives there for that reason. Whether
the eventual customer-facing tool is `parena cloud`, a new `emily.cli` subcommand tree, or its own
binary that imports this same `internal/opsplaybook` package is still a real, open decision — not
resolved here, and not blocking this format or its renderer from being real and usable today.

## Explicitly mentioned, explicitly deferred: Ansible and Terraform

Per the founder's own direct instruction. A `playbook.toml` describes **one game's own backend
process tree on one box it's already deployed to** — it says nothing about provisioning the box
itself, a database, a DNS record, a TLS certificate, or a load balancer in front of multiple
boxes. Those are real, genuinely different problems (infrastructure provisioning vs. process
supervision) with their own real, standard tools already named by the founder — a future
`ops/infra.tf` (Terraform, box/DNS/cert provisioning) and/or `ops/site.yml` (Ansible, OS-level
package/config management across a fleet) are real, deliberate extension points this format
leaves room for, not solved, not started, and not design-committed to in any way that would need
undoing later — `playbook.toml` only ever describes what runs once a box already exists and this
repo's own build has already produced real binaries on it.

## Phased plan (renderer done; the rest not yet built)

- [x] **OPS1**: the real format, specified and grounded in every existing hand-written unit file
  across the C/SDL2 arena lineage.
- [x] **OPS2**: real example `ops/playbook.toml` for `ECOWAR` (both real services) and
  `REDGARDEN` (its core `:7778` matchmaker+bot-pool pair — the 3v3/stable/players-only variants
  are named as real, future `--values-file` overrides of the same playbook, not modeled
  individually in this first pass).
- [x] **OPS3**: `emily.cli`'s real `internal/opsplaybook` renderer + `emily ops render-unit`,
  live-verified against ECOWAR's own real, currently-deployed unit shape.
- [ ] **OPS4**: the port-collision check this format makes possible (`port` is real+structured
  now) — not implemented yet, a real, concrete near-term win once someone needs it for a real
  second deploy on a shared box.
- [ ] **OPS5**: the actual "pull the repo, build it, install+enable the rendered units" apply
  step — gated on `PARENACLOUD_NORTHSTAR.md`'s own still-open CLI-surface decision, real live-
  infra-mutation work needing explicit confirmation the same way every other action of this kind
  in this session has.
- [ ] **OPS6** (real, honest, separately actionable): REDGARDEN's own live `redgarden-matchmaker-
  bots.service` still has the required (non-optional) `EnvironmentFile=` ECOWAR's own incident
  already proved is a real risk — named here, not fixed without being asked.
