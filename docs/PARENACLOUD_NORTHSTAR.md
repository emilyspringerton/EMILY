# PARENACLOUD — a real multi-tenant game-platform-as-a-service, scoped (Principle 19)

Founder real-time, this session: "lets start to think about multi tennant - what if we wanted to
build DEADWEIGHT as a different organization on top of emily for business or whatever the gaming
equivalent is? we can offer services such as engine and backend just like ROBLOX does... what do
the backend APIs look like? a dev focused product with a command line tool to scaffold projects
and to also do cloud things like the gcloud cli - like we can define somewhere what the
multiplayer paradigm is and the different matchmaking queues (casual vs ranked vs bots etc)...
now we have copy and pasted the same thing over and over its time to start paving the cow paths
and abstracting stuff that would be actually useful and that wouldnt give us the problem of the
shared library making it hard to edit code that would be really nice to edit so the library API
or engine API or whatever it needs to be overloadable in terms of like batteries included but
extensible mods first as a montra". Routed via `emily observe` first (obs `2026-09-11T19-35-55Z`,
Apple #19021) per Principle 18. Working name for this product: **PARENACLOUD** — a real, if
unlocked, name (same "pick something reasonable, let the founder correct it" convention `SAND`'s
own IRONCLAD → JEWEL → SAND history already used) chosen to reuse PARENA's own existing brand
rather than invent a new one, and to echo the `gcloud`/Roblox framing directly in the name.

This doc is the direct continuation of last turn's `EMILY/docs/PARENGINE_NORTHSTAR.md` (the
shared *engine* half of this question) — this one is the *platform/backend/multi-tenancy* half.
Read that one first; several findings below build on it without repeating them.

## Real, already-existing building blocks — checked directly, not assumed

This ask isn't starting from zero. Three real, separately-built pieces of this monorepo already
answer parts of it, and the honest first job here is finding and connecting them, not inventing:

1. **A real, live, shipped tenant-provisioning control plane already exists**: `IDUNA_PRO`'s
   `internal/tenantprovision` (`EMILY/BACKLOG.md` SECTION 267, `IDUNA-TENANTS-1`, live-verified
   2026-09-07). It allocates a free port, generates a per-tenant JWT secret, writes a per-tenant
   systemd user unit parameterized from `IDUNA_PRO`'s own real service template, starts it, and
   polls `/health` — reusing **one shared binary** (`~/.local/bin/idunapro`) config-driven per
   tenant via env vars, zero per-tenant rebuild. This is real, working infrastructure for "spin up
   an isolated instance of something for a new organization" — exactly the primitive a game
   platform needs, just not yet extended past its current scope of one long-lived web-server
   process per tenant (see the real gap below).
2. **The extensibility-mod question was already asked and answered for `IDUNA_PRO` itself**,
   2026-09-03 (`IDUNA/docs/EMILY_FOR_BUSINESS_NORTHSTAR.md`'s own "Extensibility — PARENA mods"
   section, founder's own words preserved there: "how do we abstract the primitives... to let
   people build on top... PARENA mods built in I guess"). The real, decisive, already-made
   finding: `IDUNA_PRO`'s host is **Go**, so its mods compile via **`BURROW`'s native Go emission
   target** (`burrow build ... -o file.go`, a real, ordinary Go import — no cgo/FFI boundary),
   **not** PARENA's C target, which is what every game-engine mod (`REDGARDEN`/`ECOWAR`/
   `PAPERCRAFT`) uses instead because *their* host is C. This is a real, load-bearing, two-track
   split already decided, and this doc inherits it rather than re-deciding it: **platform/backend
   logic extensibility → BURROW-Go mods; game-engine/client logic extensibility → PARENA-C mods**
   (per `PARENGINE_NORTHSTAR.md`'s own `stdlib/engine/` proposal).
3. **The exact duplication this ask is naming already exists, concretely, in IDUNA today.**
   Checked directly: `IDUNA/internal/http/handlers/` has real, separate, hand-written files —
   `redgarden_ticket.go`, `redgarden_player_ticket.go`, `redgarden_self_ticket.go`,
   `redgarden_stats.go`, `racer_ticket.go`, `papercraft_ticket.go`, `shankpit_auth.go`,
   `shankpit_queue.go`, `shankpit_ticket.go` — one hand-coded Go file family **per game** solving
   the identical problem (verify a connect ticket, track a match queue, record stats/leaderboard)
   with zero shared implementation. This is `PARENGINE_NORTHSTAR.md`'s own "copy-pasted-by-hand"
   finding, playing out again on the backend/IDUNA side instead of the client/engine side — the
   exact same root problem the founder is naming, found in a second place independently, which is
   real corroborating evidence rather than a coincidence.

## Real gap: game-backend tenancy is a different *shape* than IDUNA_PRO's current tenancy

`tenantprovision`'s real, current model is **one long-lived web-server process per tenant** — the
right shape for `IDUNA_PRO`'s own customers (CarePyre etc.), who need one HTTP service that's
always up. A game backend's real operational shape is different and more complex: a **long-lived
matchmaker daemon** that itself **spawns short-lived, per-match child server processes on
demand** (`apps/matchmaker`'s own real `fork()`/`execl()` pattern, already live in every arena
game), plus optionally a **persistent bot-pool process**. "Provision a tenant" for a game org
means standing up that whole small process tree, not one binary — a real, materially different
lifecycle `tenantprovision` doesn't handle today. **Named here as a real, undecided fork, not
resolved**: (a) extend `internal/tenantprovision` to understand a multi-process "service group"
shape generically, so both IDUNA_PRO's own single-web-server tenants and a game tenant's
matchmaker+bot-pool tree provision through the same code path, or (b) build a separate, game-
specific provisioning control plane that reuses `tenantprovision`'s own proven primitives
(port allocation, secret generation, systemd unit templating, health polling) without sharing its
actual code path. (a) is less duplication if it's real; (b) is safer if game-tenant lifecycle
turns out to need enough that's genuinely different (killing/reaping ephemeral match-server
children, for instance, which IDUNA_PRO's own single-process tenants never need) that forcing it
through the same abstraction bends it out of shape. Not decided here — a real design pass, the
same "real fork named explicitly" discipline every scoping doc in this monorepo already follows.

## What "the backend APIs" concretely look like — grounded in real, existing code, not invented

The founder's own question — "what do the backend APIs look like?" — has a real, close-at-hand
answer once the per-game ticket/queue files above and `apps/matchmaker`'s own real flags
(`--server-bin`, `--listen-port`, `--first-game-port`, `--lobby-size`) are looked at side by side:
**the matchmaker is already almost config-driven — it just externalizes its config as CLI flags
set once at process launch, instead of a real, declarative, per-tenant document.** The real,
concrete gap between what exists and what "casual vs ranked vs bots" needs is small and nameable:

- **Today**: one matchmaker process = one queue = one `--lobby-size`. No ranked/casual
  distinction exists anywhere in this codebase; "bots" means a separately-run, always-on
  `arena_bot` pool filling whichever single queue exists, not a selectable queue property.
- **Needed**: a declarative `GameConfig` a tenant defines once — real, concrete schema, grounded
  in what already exists rather than invented from nothing:

  ```
  # deadweight.parenacloud.toml (illustrative — real format TBD in implementation)
  [server]
  binary = "apps/deadweight_server"
  port_range = "9700-9899"

  [[queue]]
  name = "casual_1v1"
  lobby_size = 2
  ranked = false
  bot_fill = "always"        # matches today's arena_bot-pool-always-on convention
  bot_fill_after_seconds = 0

  [[queue]]
  name = "ranked_1v1"
  lobby_size = 2
  ranked = true
  bot_fill = "after_timeout"  # never in the shipped design today -- new, named work
  bot_fill_after_seconds = 30
  ```

  A `GameConfig`-aware matchmaker either (a) runs one queue per matchmaker process, same as
  today, just with `--lobby-size 2` etc. replaced by reading this file (a small, low-risk change
  to the existing generic `apps/matchmaker` binary), or (b) becomes genuinely multi-queue-aware
  within one process (players pick a queue, `try_match()` runs per-queue instead of globally) — a
  real, bigger change to `apps/matchmaker`'s own internals. Not decided here; (a) is the safe
  first step (multiple matchmaker processes, one per queue, behind one client-facing "which queue
  do you want" selector) and doesn't touch `apps/matchmaker`'s own already-proven internals at
  all — the same "smallest rule that resolves the gap" discipline `DEADWEIGHT/docs/SPEC_REVIEW.md`
  already used.
- **Backend API surface** (real, minimal, HTTP-over-IDUNA_PRO-or-a-successor, not designed in
  full here): `POST /api/v1/games/{tenant}/queues/{queue}/find-match` (replaces `PACKET_FIND_MATCH`
  UDP for any web/mobile-lobby front end that wants one — the UDP path stays the real, low-latency
  one for the actual in-match client, this is only for pre-match lobby/queue UX), `GET
  /api/v1/games/{tenant}/leaderboard` (replaces N hand-written `*_stats.go` files with one generic
  ladder-read endpoint parameterized by tenant+queue), and the guest-account provider
  `DEADWEIGHT/docs/PHASE_D2_SERVER_AND_ACCOUNTS.md` already scoped — real, direct evidence that
  a **generic, per-tenant game-account/queue/leaderboard API** is exactly what DEADWEIGHT's own
  D2 phase needs anyway, making it the natural first real consumer of this too (same role it
  already plays for `PARENGINE_NORTHSTAR.md`'s own `net_snapshot.prn`).

## The CLI — "gcloud for games," grounded in what already exists

**Real recommendation: extend the existing `parena` CLI with a `cloud` subcommand group, rather
than inventing a new binary.** `parena build` is already a real, live, known command in this
monorepo; `gcloud`'s own UX convention is itself "one binary, many resource-scoped subcommand
groups" (`gcloud compute`, `gcloud sql`, ...) — the same shape, not a coincidence to route around.
Real, concrete subcommand sketch, each grounded in a real, already-scoped capability rather than
invented:

- `parena cloud init <name>` — scaffold a new game project. Real, concrete scaffolding target:
  `DEADWEIGHT`'s own real repo shape (`CLAUDE.md`, `NORTHSTAR.md`, `docs/PHASE_*.md`,
  `packages/common`/`packages/simulation`/`apps/{server,client,matchmaker}` directory shape) is
  the literal template — **deliberately not designed until DEADWEIGHT itself has actually reached
  D2** (a real, working server) — scaffolding a shape that hasn't been proven real yet risks
  baking in guesses instead of a working pattern, the same "prove it with one real consumer
  first" discipline `PARENGINE_NORTHSTAR.md`'s own P1/P2 ordering already uses.
- `parena cloud deploy` — wraps whatever `tenantprovision`'s own real, generalized multi-process
  provisioning ends up looking like (see the real, undecided fork above) — not designed further
  until that fork is resolved.
- `parena cloud queues list/create/edit` — a real, thin CLI front end over the `GameConfig`
  schema above, once it exists.
- `parena cloud tenants` — admin-only, wraps `IDUNA_PRO`'s own real `GET/POST /api/v1/tenants`
  (already shipped, `iduna.admin`-gated) directly, no new backend work needed for this one.

## The extensibility mantra, made concrete: "batteries included, mods-first, overridable"

This is the real technical answer to the founder's own stated worry — "that wouldnt give us the
problem of the shared library making it hard to edit code that would be really nice to edit."
The rule, stated once so every future `stdlib/engine/*.prn` and `GameConfig` field can be checked
against it: **every shared default is a swappable mod slot, not a vendored dependency.** Concretely:

- A `stdlib/engine/login_ui.prn` (per `PARENGINE_NORTHSTAR.md`) ships a real, working default
  implementation — but the host (the C client, or PARENGINE's own eventual render layer) calls it
  through a named hook, the same "host calls into swappable PARENA logic" shape every existing
  gameplay mod already uses. A tenant who wants their own login screen writes their own
  `login_ui.prn`-shaped mod and registers it in their own `GameConfig` — **they never touch, fork,
  or vendor the shared file itself.** This is the literal difference between "a shared library
  that's hard to edit" (what the founder is naming, correctly, as a real risk) and "a batteries-
  included default that's designed to be replaced," and it's the same principle `PARENA`'s own
  "mod is the trigger, host does the real work" idiom has been quietly proving out for game logic
  this whole time — this doc's real contribution is naming it as the platform's own explicit,
  written contract instead of an implicit convention only visible by reading enough source code.
- Same two-track split as the mod-compilation target itself: engine/client hooks are PARENA-C
  overridable, platform/backend hooks (auth extensions, custom leaderboard scoring, whatever a
  tenant's own business logic needs) are BURROW-Go overridable, per the already-decided split
  above.
- **Real, honest constraint, not oversold**: `BURROW`'s own real current emission capability is
  narrow — `LO`'s own capability audit (this monorepo's own words) already found it "scalar+flat-
  struct only, no `let`/construction/`defenum`/`match`/`loop`/`Vec`" — the same class of "v0,
  don't oversell it" caution `DEADWEIGHT/NORTHSTAR.md`'s own PARENA-Java-emitter audit already
  applied this same week, now applying to the *other* half of this platform's own mod story.
  "Batteries included, mods-first" is the real, correct target architecture; whether a tenant's
  own BURROW-Go mod can express everything they'd want on day one is a real, separate,
  already-known-narrow constraint, not a solved problem.

## Phased plan (not yet built) — explicitly gated on the two prerequisite docs finishing first

This doc does not start its own phase list from zero — it's gated on real, already-scoped work
landing first, and says so rather than pretending otherwise:

- [ ] **Gate 1**: `PARENGINE_NORTHSTAR.md`'s own P1-P2 (the zero-controversy `stdlib/engine/`
  extractions, `net_snapshot.prn`) — a multi-tenant game *platform* is premature to design in
  detail before there's even one real, shared game *engine* module to provision tenants around.
- [ ] **Gate 2**: `DEADWEIGHT`'s own D1-D2 (`docs/PHASE_D1_CORE_LOOP.md`,
  `docs/PHASE_D2_SERVER_AND_ACCOUNTS.md`) — a real, working server + guest-account flow is the
  concrete thing this doc's own `GameConfig`/backend-API sketch needs to be checked against, not
  designed further in the abstract.
- [ ] **PC1**: resolve the real, undecided tenancy-shape fork above (extend `tenantprovision`
  generically vs. a separate game-specific control plane) — a founder-level architecture decision,
  same as `MULTI-TENANCY-NORTH`'s own still-open DB-per-tenant-vs-row-level fork, which this whole
  question inherits rather than re-litigates.
- [ ] **PC2**: build the `GameConfig` schema for real against DEADWEIGHT's own actual queue needs
  (at minimum: casual 1v1, ranked 1v1 — DEADWEIGHT's own V0 cut per `NORTHSTAR.md` already
  excludes tournaments/2v2, so this doesn't need to solve queue types it doesn't have yet either).
- [ ] **PC3**: `parena cloud init` scaffolding, once DEADWEIGHT's own real repo shape (proven
  through D1-D2, not guessed at) is the actual template being generalized.
- [ ] **PC4**: the mods-first override contract, formalized as real PARENA/BURROW tooling
  (a registration mechanism, not just a written convention) — gated on `PARENGINE_NORTHSTAR.md`'s
  own P3 (the mesh/shader architecture pass) landing first, since the override mechanism needs a
  real hook system to attach to.

## Non-goals

Not a commitment to DB-per-tenant vs. row-level isolation (inherits `MULTI-TENANCY-NORTH`'s own
open fork, doesn't resolve it here). Not a ranking/MMR algorithm design — orthogonal, real, and
deferred. Not a public self-serve signup flow for third-party game studios — this doc scopes
DEADWEIGHT as the first real internal tenant proving the pattern, the same role CarePyre already
plays for `IDUNA_PRO`'s own general multi-tenancy work; opening this to outside customers is a
real, separate, later business decision, not assumed here.
