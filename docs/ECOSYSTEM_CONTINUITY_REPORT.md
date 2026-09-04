# EINHORN_INDUSTRIAL Ecosystem — Full Continuity Report

**EMILY-CR-12414** ("the emily ecosystem is a comprehensive thing at this point containing
multitudes we need an actual list of all of our repos and the purpose and the status full full
full continuity report"). Real, checked directly against `/home/fatbaby/CLAUDE.md`'s own
per-repo table (the actual, maintained source of truth every session reads first) as of
**2026-09-04** — organized here as a real report, not a restatement of the operational
reference. Keep both in sync going forward: this doc groups and synthesizes; `CLAUDE.md` stays
the authoritative, line-by-line record.

## 1. What this actually is

Not one project — a real, sprawling monorepo (~50 repos) spanning: the RSI agent loop that runs
this whole operation, a live IAM/audit backbone, several real games and game engines, a
from-scratch programming language with two independent compilers, a live Minecraft community
server with its own game-design layer, a media/blog pipeline, and a growing set of small
product/tooling spinoffs. Held together by one real, enforced operating discipline (**The Emily
Way** — backlog-first, Apple-before-mark-done, session-tagged commits) rather than by shared
code.

## 2. The operational spine — how work actually happens here

| Repo | Real role |
|---|---|
| `EMILY` | Emily Prime agent: RSI cron loop (OBSERVE→DECIDE→ACT→PLAN), the golden `BACKLOG.md` every session reads first, HEIMDAL sprint intake from MJOLNIR. |
| `IDUNA` | The real, central trust authority — Google/local/M2M auth, ES256 JWTs, RBAC, the Apples audit ledger, HEIMDAL sprint storage, the kanban↔BACKLOG.md bridge every "priority"/"cruise" cycle actually runs through, and a growing set of real admin GUI pages for GFD content. |
| `emily.cli` | The `emily` operator CLI — apples/observe/kanban/changelog/session, the actual command surface this whole report was assembled through. |
| `PRRJECT_FATBABY` | SEC/PR financial signal pipeline; also hosts the real observation-file queue `emily observe` writes to. |
| `APPLES` | Git-authoritative backup of every Apple ever filed — pure data, no code. |
| `EmilyOS` | A separate, real policy kernel (posture-gated sessions, RBAC, SOC 2 audit log) — governance infrastructure, not the agent itself. |
| `PITVIPER` | SDL2 terminal emulator with Emily Prime integration hooks — the founder's own stated "this VPS is my IDE" vehicle, real precedent `DUNG` is rewriting forward. |

## 3. Games & interactive worlds — real, live, or actively developed

| Repo | Real status |
|---|---|
| `GoblinFoxDragon` | **Most actively developed this session.** Real PvE MMO (`apps2/mud`) + PvP arena client (`apps2/battlegrounds_gui`), dungeons, Notorious Monsters, a new 23rd job, GFD's first two real PARENA/BURROW mods. See its own `docs2/STACK_CONTINUITY_REPORT.md` for the full picture — real, decisive finding this session: its Battlegrounds PvP mode has no native server, parasitic on REDGARDEN's own live, independently-evolving binaries, and has drifted (`docs2/BATTLEGROUNDS_MIGRATION_NORTHSTAR.md`). |
| `REDGARDEN` | The real, live, independently-evolving MOBA/arena source `GoblinFoxDragon`'s own client was forked from — still being actively developed on its own, which is exactly what's causing GFD Battlegrounds' own drift. |
| `ECOWAR` | Hard fork of REDGARDEN — first repo to ship real PARENA decision-logic mods (`card_effect_mod.prn`), not just triggers; real 16-card content shipped. |
| `SHANKPIT` | Server-authoritative UDP FPS; the deep ancestor of both `GoblinFoxDragon` and `skateboard`'s own C/SDL2 lineage. |
| `BRAWLPIT` | Server-authoritative platform fighter; this session shipped real matchmaking (4 phases), a sandbox training mode, and two new fully-kitted characters (Sunlit Draw, Sequel Duck). |
| `PAPERCRAFT` | "Minecraft meets GTA3 plus Skate2" sandbox, mods-first from day one — real PARENA mods for XP/leveling/items already compiled and tested; no server/client host code yet. |
| `WEAKNIGHT_BEDROCK_RACERS` | Real-Minecraft-Bedrock-protocol racer game, consumes `GoblinFoxDragon/apps2/server-go`'s own `/heightmap` API. |
| `GTA7` | TRAPX doctrine as a real Paper plugin on `EINHORN_SURVIVAL` — scoped, VS0 not yet built. |
| `EINHORN_SURVIVAL` | The real, live, actual Minecraft community server (`:25565`) all the TRAPX/GTA7 game-design work targets. |
| `skateboard` | Historical repo where the Papercraft concept was first scoped — kept as the real audit trail, not the live build. |
| `TYLER` + `MoneyPrinterTurbo` | TV series bible/scripts + the flat-stream video compiler that turns them into real episodes — the media-production half of this ecosystem. |
| `MISHRI` | Forked external Minecraft bot (human-like movement/chat) aimed at `EINHORN_SURVIVAL`, for testing. |

## 4. The PARENA language ecosystem — real, growing, increasingly self-hosting

| Repo | Real status |
|---|---|
| `PARENA` | The language itself (S-expression syntax, region-based memory safety). VS0 (the `parena-c` compiler, C/JVM/TS/Wasm targets) is real and shipping — the C target is mature (defstruct/defenum/match/loop/Vec/FFI), other targets narrower. Self-hosting effort (`selfhost/*.prn`) in progress. |
| `BURROW` | A full, parallel Go+PARENA rewrite of the PARENA compiler. Phases 0-6 shipped, including a real native Go emission target — **first real host adoption happened this session** (`GoblinFoxDragon/server/modevent`'s `nm_bonus_mod`), closing out a flag this repo's own docs carried since the target shipped. |
| `JEWEL` | Real Jupyter kernel for PARENA, live behind an nginx+basic-auth proxy on okemily.com. |
| `SAND` | Named-only stub — a future PARENA-native code editor, distinct from `DUNG`. |
| `DUNG` | Scoped, not built — the real BURROW-compiled rewrite of `PITVIPER` + PARENA's own editor stdlib, gated on BURROW's own struct-emission work landing first. |
| `LO` | Esolang spec under real critical review (grammar/tokenization gaps found and named) — compiles down to `.prn` for `parena`/`burrow` to finish the job; no compiler code yet. |
| `SPIDERBEETLE` | Real, verified proof that PARENA's Java emitter works end-to-end (a real battery-UI decision function, tested with an actual JDK) — found and fixed a real PARENA compiler bug (`(not x)` unhandled) along the way. |
| `MIXFORGE` | Planned DJ-app product, explicitly committing to build in PARENA end-to-end rather than the more conventional C++/JUCE stack a sourced design doc proposed. |

## 5. Business/product layer — extracting a product out of internal infrastructure

| Repo | Real status |
|---|---|
| `IDUNA_PRO` | The real, standalone extraction of IDUNA's generic IAM core (OAuth/JWT/RBAC/Apples/logging) for `EMILY_FOR_BUSINESS` — zero cross-imports from IDUNA's own monorepo-custom code, confirmed directly before extraction. A real, ongoing BURROW/PARENA CLI effort (kanban card 9988) is ahead of it. |
| `EMILY_FOR_BUSINESS` | Licensing-only so far (a real custom source-available license + trademark policy) — real product scope lives in `IDUNA/docs/EMILY_FOR_BUSINESS_NORTHSTAR.md`, naming a real, unresolved tension against IDUNA's own "not a product, the backbone" self-framing. |
| `WOTAN` | Splitting EINHORN_INDUSTRIAL's esports/stats hub out of the marketing site into its own repo/subdomain — a hat-store NORTHSTAR exists; not deployed yet. |
| `CONSOLE` (console.okemily.com) | Real subdomain live, real signup/status front-end scoped into SECTION 244 subtasks (tenant DB, provisioning, TLS, a real emily.cli Cloudflare-token command) — none built yet, all real, concrete next steps. |

## 6. Infra / dev tooling

| Repo | Real status |
|---|---|
| `gpt2-alpine-c` | GPT-2 inference in C + the Emily Prime fine-tune pipeline. |
| `SKULDMARK` | A real, small, complete, public-domain instrument-identifier format spec — done, stable. |
| `FLASH` + `image-builder-rpi` | Forked HypriotOS tooling for a future EmilyOS-Pi/PARENA distro (S213) — `FLASH` has a real Bazel build and two real, found Bazel gotchas fixed; blocked on Windows/WSL support the founder's own current machine needs. `image-builder-rpi` untouched, sequenced after `FLASH` proves out. |
| `EDIS` | WordPress plugin suite fronting FatBaby's signal API and Emily Prime for `okemily.com`. |

## 7. New, small, or stub repos (real, but early)

`CarePyre`, `TTT` (Tyler Teaches Typing — real VS0 actually lives in `SHANKPIT`), `EXODUS`
(future home for REDGARDEN's own PFSP autocurriculum research), `XCVBNM-OR`,
`QUEENSALLYONLINEBOOKOFMAGIFICATIONANDUNICOR`, `NORN` — named and tracked, no real scope decided
or code written yet for most of these; not a gap, a real, honest reflection of where founder
attention hasn't landed.

## 8. Real, current cross-cutting themes

- **Mods-first is now a working discipline, not just a stated goal.** `ECOWAR`, `PAPERCRAFT`,
  and — as of this session — `GoblinFoxDragon` (on both its C and Go hosts) all have real,
  shipped, tested PARENA/BURROW mods doing genuine decision logic, not placeholder triggers.
- **BURROW's Go target went from "shipped, unused" to "shipped, proven"** in a single session —
  real evidence the PARENA/BURROW dual-compiler bet is paying off across more than one host
  language, not just the original C target.
- **The Emily Way's own bottleneck got found and fixed this session too**: the kanban↔BACKLOG.md
  bridge every one of these repos' own priority-queue work runs through was re-parsing a
  30,000-line file on every call (~250ms); now cached (~29µs). A real, ecosystem-wide throughput
  fix, not scoped to any one repo.
- **Forking is the default growth pattern, and it has a real, recurring cost.** `GoblinFoxDragon`
  ← `REDGARDEN`, `ECOWAR` ← `REDGARDEN`, `SAND`/`DUNG` ← `PITVIPER`, `SHANKPIT` ← `skateboard` —
  every major game/tool in this ecosystem is a fork of something else in it. `GoblinFoxDragon`'s
  own Battlegrounds drift (§3) is the real, concrete cost of that pattern when a fork's source
  keeps moving and nothing tracks the two staying in sync.
- **Persistence keeps being the real, deferred piece.** GFD's inventory, GFD's per-job levels,
  and (this session) IDUNA's own schema limits all hit the same real wall — an in-memory or
  single-column model that a real feature ask outgrows — and all three got the same honest
  "scope it, don't build blind" treatment rather than a rushed schema change.

## Related

- `/home/fatbaby/CLAUDE.md` — the real, authoritative, line-by-line per-repo table this report groups and synthesizes. Update both when a repo's status materially changes.
- `GoblinFoxDragon/docs2/STACK_CONTINUITY_REPORT.md` — the same real report pattern, one level deeper, for GFD specifically.
