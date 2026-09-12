# Monorepo Product Audit — 2026-09-12

## Why this exists

Founder real-time, this session: *"Format that'll get you the most useful triage back: product,
what a stranger can do in the first 60 seconds, install friction, and what's actually broken.
Onboarding time alone won't rank them — a 10-second demo that needs a download loses to a
60-second one that doesn't. GO THROUGH EVERY REPO CURRENTLY ON THIS BOX IN THE MONO REPO AUDIT
THE STATE OF THE PRODUCT..."*

This is that audit. Every real git repo under `/home/fatbaby` (56 of them — `for d in */; do [ -d
"$d/.git" ] && echo "$d"; done`, checked directly rather than assumed) got the same four
questions:

1. **Product** — what it actually is, one or two sentences, no marketing.
2. **First 60 seconds, zero install** — what can a stranger with just a browser (or an already-
   installed terminal/SSH/Minecraft client) actually *do*, right now, with nothing to download.
   This is scored as a hard category, not folded into a speed number, per the founder's own
   explicit instruction: **a 10-second demo that needs a download loses to a 60-second one that
   doesn't.**
3. **Install friction** — for everything that fails #2, what does getting to a working demo
   actually cost (build toolchain, an app store CI that's broken, a client binary that doesn't
   exist yet)?
4. **What's actually broken** — named plainly, not hedged.

Data sources, all checked directly in this session rather than trusted from memory: each repo's
own README/NORTHSTAR (head of file), `git log -1` + commit count for recency/maturity, `ss -tlnp`
for what's actually listening on this box right now, and the live nginx config
(`/etc/nginx/sites-enabled/*`) for what's actually reachable from the public internet and at what
hostname/path.

**Live, listening on this box right now** (confirmed via `ss -tlnp`, not assumed from a README):
`:8080` (IDUNA), `:2323`+`:2222` (GoblinFoxDragon telnet+SSH), `:25565` (EINHORN_SURVIVAL
Minecraft/Paper), `:8890` (JEWEL Jupyter, proxied). **Live, reachable from the public internet**
(confirmed via nginx `server_name`/`proxy_pass` blocks): `okemily.com` + `www.okemily.com` (→
IDUNA :8080 for `/api`, `/admin`, `/portal`; static marketing elsewhere; `/gfd-ws/*` → GFD
websocket ports), `news.okemily.com` (→ :8082), `carepyre.org` + `www.carepyre.org` (static site +
`/api` → :8080), `iduna.farthq.com` (WordPress — this is EDIS, not IDUNA itself, despite the
hostname).

---

## Tier 1 — Zero-install, works in a browser or a tool everyone already has

These are the only real candidates for "a stranger does something in 60 seconds with nothing to
download." Ranked by how close to frictionless the actual first action is.

| # | Repo | Product | What a stranger does, right now, no install | Real friction that remains |
|---|---|---|---|---|
| 1 | **OKEMILY** | Marketing/company site + blog + tournaments page for EINHORN_INDUSTRIAL | Visit `okemily.com` in any browser. Reads, no account needed. | None for reading. Anything beyond reading (account creation) hands off to IDUNA below. |
| 2 | **CarePyre** | Community trust landing page + Contact form + a 6-pillar "C.H.A.N.G.E. Initiative" page | Visit `carepyre.org`, read the pitch, submit the real Contact form (posts to IDUNA's `/api/v1/carepyre/contact`, not a dead form). | None to read/contact. No actual product to *use* yet — this is pre-launch messaging + intake, honestly. |
| 3 | **EDIS** | WordPress site (`news.okemily.com`) surfacing FatBaby's governance/insider-trading signal pipeline as readable articles | Visit the site, read a signal-driven article. | None to read. The interesting part (the pipeline) is invisible to a visitor — they just see a WordPress blog. |
| 4 | **IDUNA** | Platform IAM: account creation (`okemily.com/` → Google OAuth → gamertag), dev portal (`/portal`), admin/kanban UI | Click "sign in with Google" at `okemily.com`, get a real account + gamertag in one flow — genuinely fast. | The portal/admin/kanban surfaces behind it are gated (devportal access, admin role) — a stranger's 60 seconds ends at "I have an account," not at anything that looks like a finished product. |
| 5 | **GoblinFoxDragon (GFD)** | Persistent MUD ("DragonsNShit") — real jobs, combat, spells | `telnet <host> 2323` or `ssh <host> 2222` — both clients ship preinstalled on macOS/Linux; Windows 10+ ships `ssh` but not `telnet`. Guest login exists (no account needed) per this session's SSH-identity work. Genuinely playable inside 60 seconds *if* the visitor already has a terminal open. | Not literally zero-install for a Windows user with no terminal habit — real but small (`ssh` is one Settings toggle away on modern Windows, still cheaper than any client download). The README itself is written assuming the reader is comfortable pasting a `telnet`/`ssh` command — that's a real, self-selecting audience narrower than "any stranger." |
| 6 | **JEWEL** | Real Jupyter kernel/backend for PARENA, served at `okemily.com/jewel/` | Open the URL. | Gated behind HTTP Basic Auth (a deliberate interim measure pending Google OAuth) — a stranger cannot actually get in without a shared credential today. So this is zero-install but **not** zero-friction; effectively not a public demo right now. |

**Everything else on this list is either not deployed publicly, requires a client the visitor
doesn't already have, or isn't consumer-facing at all** (dev tooling, compilers, internal
infra) — see Tiers 2-4.

---

## Tier 2 — Real, playable/usable products that require a real install

These are legitimate products with real depth, but the founder's own ranking rule demotes them
below Tier 1 regardless of polish: the first thing a stranger has to do is download or build
something.

| Repo | Product | What exists | Install cost | What's actually broken |
|---|---|---|---|---|
| **EINHORN_SURVIVAL** | Real, live community Minecraft survival server | Live on `mc.okemily.com:25565` (Java) / `:19133` (Bedrock) | Requires owning Minecraft (Java or Bedrock) — a real, non-trivial purchase for a brand-new stranger, but "already own Minecraft" describes a huge chunk of the internet. Cheapest real install cost on this whole list once that's true. | Nothing named broken; this is one of the more solid, actually-live things in the monorepo. |
| **GTA7** | TRAPX doctrine plugin running *inside* EINHORN_SURVIVAL | Same connection as above, `/flow` and other in-game commands | Same as EINHORN_SURVIVAL (no separate install — it's a plugin, not a separate server) | VS0 not yet fully built per its own NORTHSTAR; scope is real but shallow today. |
| **SHANKPIT** | Server-authoritative UDP FPS, C/SDL2 client, multiple modes (arena, CTF, story, helicopter physics) | CI-built Windows client exists | Requires downloading a built client (CI artifact) or building from source — no browser path | CI was fixed this cycle (cutscene/audio linking), but "download a game client from CI artifacts" is real, meaningful friction for a stranger vs. a store listing. |
| **shankpit-460** | SHANKPIT fork — bot-pool deathmatch, 3-button lobby | Real bot-pool deathmatch live via `shankpit460-emily-bot.service` | Same C/SDL2 client-download friction as SHANKPIT | Root-caused connectivity bug (missing `SHANKPIT_TICKET_SECRET` in `PLAY.bat`) already fixed per its own README — genuinely in better shape than the friction score alone suggests. |
| **REDGARDEN** / **redgarden-deploy** | "Knights of the Void" arena MOBA, `apps/arena` | Deep, documented keybind contract, real heroes/shop system (124 heroes per TYLER lore tie-in) | Native C client — no browser path, no store listing found | Nothing named broken in the README head; this looks like the most feature-complete game in the portfolio, entirely gated behind a client download. |
| **redgarden-stable** | Frozen/stable branch of REDGARDEN | Same game, last commit 2026-08-20 (three weeks stale vs. `REDGARDEN`'s same-day commits) | Same as REDGARDEN | Its own staleness relative to the actively-developed `REDGARDEN`/`redgarden-deploy` pair is the notable thing — worth confirming which of the three is actually meant to be "the" build strangers get pointed at, since right now there are three repos with the same game and no obvious single canonical download link. |
| **ECOWAR** | Hard fork of REDGARDEN — separate 1v1 matchmaker/bot-pool, its own 16-card system | Own matchmaker staged (`:9779`) | Same native-client friction as REDGARDEN | Most actively committed repo in the whole monorepo (commits same day as this audit) — real, fast-moving, but the card/mechanic layer is the newest, least-battle-tested part. |
| **BRAWLPIT** | 2.5D platform fighter built on SHANKPIT (Build 178); TIPJAR bartender mode lives inside it (press T) | Steps 1-3 shipped (core shift loop, per-player sim, real 2-player split-screen) | Same native-client friction | Steps 4-7 (competitive/co-op, content, polish) explicitly not done yet per its own README — an honestly-scoped partial build, not a finished game. |
| **WEAKNIGHT_BEDROCK_RACERS** | Voxel-terrain racing game, real physics, Xbox controller support, bot opponents | Phase 0 shipped + live-verified, two selectable tracks | Native client, no browser path | Nothing named broken; realistically the most "done for its own current phase" of the racing/sandbox cluster. |
| **PAPERCRAFT** | Non-voxel destructible-geometry open-world sandbox ("Minecraft meets GTA3 plus Skate2"), mods-first via PARENA | NORTHSTAR + one verified PARENA mod only | No server/client host code exists yet — literally nothing to install or play | This is pre-alpha scoping, not a product yet — listed here for completeness, not because a stranger can do anything with it. |
| **MJOLNIR** | Android intelligence terminal (push notifications, Apple feed, front door to the products) | Real Kotlin/Compose app exists in source | **Cannot be installed at all right now** — every CI run (13/13) fails at "Validate google-services.json," so there is no APK to download, period. | This is the single worst "broken" finding in the whole audit: a real, apparently-otherwise-working app that a stranger cannot obtain in any form, blocked on one human-only GCP Console step the README says nobody else can do. |
| **PITVIPER** | SDL2 terminal emulator + GFD game client, Emily Prime integration hooks, new SSH client mode | Real, shipping steadily (SSH client mode landed 2026-09-06) | Build-from-source only, no distributed binary found | Nothing named broken; it's simply a developer-facing terminal, not something a general stranger would seek out regardless of install cost. |

---

## Tier 3 — Real infrastructure, compilers, and internal tooling (not consumer products)

The "stranger in 60 seconds" frame doesn't really apply to these — they're B2B/dev-tool/infra
plumbing, not something a random visitor would be handed a link to. Judged instead on: is it
real, is it live, is it moving.

| Repo | What it is | Real/shipped? | Notable gap |
|---|---|---|---|
| **PARENA** | Core in-house language/compiler (S-expr, region-safety, multi-target C/JVM/TS/Wasm) | Very real, 609 commits, actively developed (last commit same day as this audit) | This is the load-bearing piece under half the portfolio (KARAMBIT, SPIDERBEETLE, LO, JEWEL, MIXFORGE, ECOWAR's card logic, PAPERCRAFT mods) — its health *is* the portfolio's health for anything "PARENA-native." |
| **BURROW** | Full Go+PARENA reimplementation of the PARENA compiler itself | Real, 33 commits, lexer/parser parity verified against the C reference's own 111-file corpus | Emission-target work (needed by DUNG, IDUNA_PRO's customer-mod story) is the current bottleneck. |
| **emily.cli** | Operator CLI for the whole agent system | Real, mature (216 commits), zero-LLM, zero-dependency by design | None named — this is one of the most solid, boring-in-a-good-way pieces of infra here. |
| **EMILY** (this repo) | Meta-orchestration layer — Emily Prime, RSI loop, cross-repo backlog | Real, by far the largest repo (5,963 commits) | N/A — this is the process itself, not a product being audited for consumer fit. |
| **IDUNA_PRO** | Standalone extracted core of IDUNA for the "Emily for Business" product | Real, 79 commits, live-verified boot + real self-serve register/login/JWT flow | Found-but-not-fixed: `/api/v1/identities/me` doesn't resolve a local-auth JWT subject — matters more here since this product leans on self-serve local auth more than IDUNA itself does. |
| **PRRJECT_FATBABY** | SEC/PR financial signal intelligence pipeline (630MB+ accumulated data) | Real, mature (519 commits), the actual value driver behind EDIS | No consumer-facing surface of its own — its value only reaches a stranger secondhand, via EDIS. |
| **EmilyOS** | Go policy kernel (posture-gated sessions, RBAC, tamper-evident audit log) | Real, 58 commits, has its own CI/release pipeline now | Not a consumer product by design — SOC2-readiness infra. |
| **GOLDENBAND** | `.gband` motion-asset format + C sampler + Blender pipeline | Real, working rigs/animation pipeline | Purely a content-pipeline tool for GFD, not player-facing. |
| **gpt2-alpine-c** | Pure-C GPT-2 inference + Emily fine-tune pipeline | Real, actively developed (68 commits, this week) | Dev/research tool, not a product surface. |
| **KARAMBIT** | Android network toolkit (LAN SSH-port scanner), PARENA-compiled decision logic | Real, small, working v0 | No installable APK path documented (no Android SDK in this sandbox to build one) — same class of gap as MJOLNIR, smaller stakes. |
| **SPIDERBEETLE** | Proof that PARENA's Java emitter can power a real Android app | Real generated Java, verified against a real JDK | No actual Android project (Gradle/Manifest/Activity) exists yet — this sandbox has no Android SDK at all, so even the existing MJOLNIR app can't be built here, let alone this one. |
| **LO** | Hyper-minimalist esolang over a PARENA backend | Real critique + phased plan, 39 commits of real progress since | Still frontend/compiler work, nothing runnable by an outsider yet. |
| **JEWEL** (compiler angle) | Jupyter kernel for PARENA | Real, live (see Tier 1) | Covered above — the same repo is genuinely both an infra tool and a Tier-1 zero-install demo. |
| **MIXFORGE** | PARENA-native DJ app (YouTube import, BPM/key detection, crossfading) | Real CLI host just landed (this week), otherwise NORTHSTAR-stage | Key/BPM detection has no PARENA story yet at all — a genuinely open technical gap, not just unimplemented. |
| **ladybug** | BDD testing framework for PARENA (Ginkgo/Gomega-shaped) | Real, small, active | Dev-tool, not consumer-facing by design. |
| **NORN** | Read-only CLI over NDJSON registries | Real, working, narrow | Narrow internal tool. |
| **SKULDMARK** | 25-char instrument identifier format | Real, small, done, public-domain | Genuinely complete for its narrow scope — nothing to fix. |
| **longma** | CSV file splitter | Small, working, generic OSS-style utility | Not really part of this company's product line — reads as an orphaned utility repo. |
| **GOLDENBAND**, **SKULDMARK**, **longma** | — | — | (grouped above/here to avoid re-listing) |

---

## Tier 4 — Stubs, dormant, forked/vendored, or non-product repos

| Repo | Status |
|---|---|
| **TIPJAR** | Empty repository — one initial commit, no content. The real, live TIPJAR feature actually lives *inside* `BRAWLPIT` (press T from the lobby) — this standalone repo appears vestigial. |
| **TIPJAR-wiki** | 4 commits, wiki-only, presumably for the same feature that actually lives in BRAWLPIT. |
| **TTT** | Explicit stub — "scope not yet decided," the real VS0 (Tyler Teaches Typing) lives in SHANKPIT instead. |
| **XCVBNM-OR** | Two files (`MISHRI_CONSTRUCT.txt`, a bare `README`), 2 commits — no real content, unclear purpose, reads as accidental/scratch. |
| **QUEENSALLYONLINEBOOKOFMAGIFICATIONANDUNICOR** | Unusual, personal/creative-writing-styled content, not a software product in the normal sense — 47 commits but no discernible product shape from the README. |
| **REDGARDEN-wiki** | Wiki repo, 4 commits, last touched 2026-02-08 — stale relative to REDGARDEN's own active development. |
| **EMILY_FOR_BUSINESS** | Licensing-only repo (a real custom source-available license + trademark policy) — explicitly no product code here; the actual scoping doc lives in IDUNA's docs. |
| **APPLES** | Not a product — the git-authoritative audit-trail data store for the whole agent system. Real and load-bearing, just not something a stranger is ever handed. |
| **EDIS** (infra angle) | Already covered in Tier 1 as the visible WordPress site; the plugin code itself is PHP infra, not separately consumer-facing. |
| **dragonfly** | Vendored fork of an external Minecraft Bedrock server library (2,777 commits, but last commit 2026-07-26 and clearly upstream-authored, not this org's own work) — not wired up as this org's live Minecraft server (that's the Paper-based `EINHORN_SURVIVAL`, a different stack entirely). Appears to be an unused/reference dependency sitting in the monorepo rather than a deployed product. |
| **emsdk** | Vendored Emscripten SDK (1,217 commits, all upstream) — a build dependency, not a product. |
| **image-builder-rpi** | Forked from `hypriot/image-builder-rpi`, explicitly not yet touched per the founder's own sequencing (FLASH comes first) — real fork, zero local work done. |
| **FLASH** | Forked from `hypriot/flash`, real local work added (Bazel build, Windows/Git-Bash support) — but the tool itself is Mac/Linux-only and explicitly refuses to run under WSL, and the founder's own real blocker (they're on Windows) is still unresolved. A CLI tool for flashing SD cards, not a consumer product regardless. |
| **MISHRI** | Forked Minecraft bot (mineflayer/Node) — real Bazel build added on top, but licensing is genuinely ambiguous (package.json says ISC, no LICENSE file present) — worth resolving before this is used or shared any further. |
| **MoneyPrinterTurbo** | Forked video-compilation tool supporting TYLER's own video pipeline — internal tooling, not itself a product. |
| **WOTAN** | Esports/stats hub being split out of `OKEMILY/tournaments.html` — explicitly **not deployed**; `index.html` is a placeholder, DNS/nginx not live yet. Until that lands, `OKEMILY/tournaments.html` remains the actual live page — so right now there are two copies of this idea, one live-but-old, one new-but-fake. |
| **DEADWEIGHT** | Real critical-review NORTHSTAR only (5 commits) — no code. |
| **JEWEL** | (see Tier 1/3 — real and live, not dormant, listed here only for cross-reference) |
| **skateboard** | Historical/ancestor repo, superseded by PAPERCRAFT — kept intentionally as an audit trail, not an active product. |
| **CarePyre** (infra angle) | Covered in Tier 1; noted here only that beyond the landing page + contact form, no actual "product" (the 4-layer stack it pitches) is built yet. |

---

## Cross-cutting findings

1. **The single biggest "broken" finding across the whole monorepo is MJOLNIR's CI.** It is the
   one app most directly aimed at being "the front door to the products" (per its own README) and
   it has been completely un-installable for 13 straight CI runs. Every other broken thing found
   in this audit is either partial-feature (BRAWLPIT steps 4-7, PAPERCRAFT's missing host code) or
   friction (native-client-only games). MJOLNIR is the only *zero* — no APK exists in any form.

2. **REDGARDEN exists as three near-identical repos** (`REDGARDEN`, `redgarden-deploy`,
   `redgarden-stable`) with real drift between them (redgarden-stable is 3+ weeks behind the other
   two). Worth a deliberate decision about which one is the canonical "download link," since none
   of the three READMEs say so themselves.

3. **The founder's own ranking rule (no-download beats fast-but-download) genuinely inverts the
   naive "most polished" ranking.** REDGARDEN and SHANKPIT are almost certainly the most
   feature-complete, actively-developed *games* in this portfolio — and they rank below IDUNA's
   one-click Google sign-up and even below GFD's telnet/SSH MUD, purely because a stranger has
   nothing to click that just works in a browser or a preinstalled terminal. If a public "try our
   stuff" link is ever assembled, GoblinFoxDragon (telnet/ssh) and OKEMILY/IDUNA (browser) are the
   only real Tier-1 candidates today; everything else either isn't public, is gated behind a
   credential (JEWEL), or requires a client download.

4. **A meaningful fraction of the 56 repos are not products at all** — vendored dependencies
   (`dragonfly`, `emsdk`), forks not yet built on (`image-builder-rpi`), pure infra/data stores
   (`APPLES`, `EmilyOS`), licensing-only shells (`EMILY_FOR_BUSINESS`), and stubs with no content
   (`TIPJAR`, `XCVBNM-OR`, `TTT`). None of that is a problem by itself, but it means the real
   "portfolio" for an outside-facing question like this is closer to 15-20 repos than 56.

5. **PARENA is the highest-leverage single repo in the whole tree that isn't itself consumer-
   facing.** Its health gates real progress in at least seven other repos (KARAMBIT, SPIDERBEETLE,
   LO, JEWEL, MIXFORGE, BURROW/DUNG, and PAPERCRAFT's mod layer, plus ECOWAR's card logic). It is
   also, encouragingly, one of the most actively-committed repos in the monorepo as of this audit.

---

*Compiled 2026-09-12, session `sess-20260905-0720-ec33e7c5`. Methodology: direct `git log`/`ss -
tlnp`/nginx-config inspection plus each repo's own README/NORTHSTAR head, not memory or
assumption. A repo's absence of a "what's broken" line above means none was found in this pass —
not a guarantee nothing is broken, only that nothing surfaced in a README-plus-recency-plus-
liveness check at this level of depth.*
