# PARENGINE — formalizing a real, shared PARENA game-engine layer (scoping, Principle 19)

Founder real-time, this session: "for the engine can we take a look at all of the games that we
have - SHANKPIT, REDGARDEN, BRAWLPIT ETC - and start to abstract the libraries that make up our
engine - like there must be things that we are doing every time that we can extract into PARENA to
start to more formalize a shape of a game engine instead of just always calling SHANKPIT a game
engine which it isnt really yet". Routed via `emily observe` first (obs `2026-09-11T19-28-04Z`,
Apple #19018) per `EMILY/docs/THE_EMILY_WAY.md` Principle 18, then investigated across every real
game repo in this monorepo before proposing anything — same discipline `EOSUI_NORTHSTAR.md`
already set for the UI half of this exact question.

## The frame-break, stated plainly

"SHANKPIT" has been informally shorthand for "the engine" in conversation, but it was never
structurally true, and that mismatch is worth naming directly: SHANKPIT is a Go UDP-FPS server
with its own C client, in the same lineage as `GoblinFoxDragon` (`packages/protocol`,
`packages/rts`, `server/`, `services/` — a DragonsNShit-MMO-shaped architecture). It shares almost
no code, and no architectural lineage, with `REDGARDEN`/`ECOWAR`/`BRAWLPIT`/`PAPERCRAFT`/
`WEAKNIGHT_BEDROCK_RACERS` — a completely separate lineage built around `packages/simulation`'s
own `arena_game.c`-shaped server-authoritative sim + an SDL2/OpenGL client. Calling either lineage
"the engine" the rest should build on was never going to work, because **neither one is actually
a reusable engine today — each is one specific game that happens to be old enough that newer games
copied pieces of it by hand.** The real ask underneath the founder's own words (Frame-Break
Reframing, `REDGARDEN/NORTHSTAR.md` §28) is: stop informally treating whichever game shipped
first as the engine, and start formalizing what's actually common across *all* of them into
something that deserves the name.

## Real findings — concrete, checked directly, not assumed

**`packages/common`/`packages/simulation` is a real, load-bearing directory *shape* shared across
five repos — but the shape is copied, not the code.** `REDGARDEN`, `ECOWAR`, `BRAWLPIT`,
`PAPERCRAFT`, and `WEAKNIGHT_BEDROCK_RACERS` all independently converged on the identical
top-level layout (`packages/common/`, `packages/simulation/`, `apps/server`, `apps/client`) —
real evidence the underlying problem shape (server-authoritative sim + a matching client) is
genuinely stable across this whole game portfolio. But every one of these is a **separate git
history**, not a shared dependency, so identical filenames inside `packages/common/` diverge
silently the moment anyone touches one without the others:

| File | REDGARDEN | ECOWAR | PAPERCRAFT | WEAKNIGHT_BEDROCK_RACERS |
|---|---|---|---|---|
| `mat4.h` | `05d2b908…` | `05d2b908…` (identical) | *(absent)* | `05d2b908…` (identical — **not even a fork of REDGARDEN**, hand-copied verbatim) |
| `hmac_sha256.h` | `a8e819f7…` | `a8e819f7…` (fork, identical) | `6fe9edd1…` (diverged) | `cd9acfe0…` (diverged) |
| `http_client.h` | `45e793c9…` | `45e793c9…` (fork, identical) | `f8de163d…` (diverged) | `a5f9e04d…` (diverged) |
| `physics.h` | 24 lines (stub) | 24 lines (stub) | *(own file)* | *(absent)* |
| `physics.h` (BRAWLPIT) | — | — | — | 1428 lines, a real, substantial, completely independent implementation under the identical filename |

Three real, distinct failure modes are visible in that one table: (1) `mat4.h` copied byte-for-
byte into a repo that isn't even a fork of the source (pure hand-copy-paste), (2) `hmac_sha256.h`/
`http_client.h` copied once, then each repo's own copy silently drifted apart — meaning a real
bug or security fix in one does **not** reach the other three, (3) `physics.h` shows the same
filename can also mean two totally unrelated pieces of code once forks diverge far enough,
which matters for anyone assuming "same path = same contract" while working across these repos.

**Live, self-caused, this exact session: `apps/matchmaker` is the strongest-looking reuse story in
this monorepo, and it just proved the same problem anyway.** `REDGARDEN/apps/matchmaker/src/
main.c` (14063 bytes) and `ECOWAR/apps/matchmaker/src/main.c` (12253 bytes) are **already
different sizes** — because this same session's own S370 work (the Mandelbrot procedural-jungle
seed plumbing) added real `--seed` generation to ECOWAR's copy and never touched REDGARDEN's
identical, still-seedless copy. That wasn't an oversight to fix retroactively (S370 was correctly
scoped to ECOWAR only, per the founder's own explicit choice) — it's a clean, live, easy-to-point-
at example of exactly the cost this whole ask is naming: **the same generically-useful fix now has
to be re-applied by hand to every other copy that wants it, forever, because there was never one
canonical copy to fix in the first place.**

**Every SDL2 client independently reinvented its own render primitives, with zero naming
convergence between lineages.** `REDGARDEN`/`ECOWAR` share `upload_mesh`/`compile_shader`/
`draw_mesh`/`mat4_*` (fork lineage — real, if fragile, reuse). `BRAWLPIT` (a 2D game) built its
own `draw_rect`/`draw_circle`/`draw_hat`/`draw_menu_button` immediate-mode layer from scratch.
`PAPERCRAFT` and `WEAKNIGHT_BEDROCK_RACERS` (both real 3D games, both needing a window + GL
context + camera + login screen, the *closest* possible match to REDGARDEN's own actual needs)
**used none of REDGARDEN's naming or structure** and independently built their own. Three separate
teams solved "create an SDL2 window and draw 3D boxes in it" three separate times in this
monorepo's own recent history, with three incompatible internal APIs, none of which future games
can build against with any confidence they'll still work next time.

**The starkest single piece of evidence: `draw_login_screen` is a near-line-for-line copy between
`PAPERCRAFT` and `WEAKNIGHT_BEDROCK_RACERS`.** Identical function signature
(`SDL_Window *win, int win_w, int win_h, const LoginScreenState *st`), identical layout math,
identical control flow (focused/unfocused field styling, "LOGGING IN..." state, error text) — the
**only** differences are the color palette, the title string, and a `pc_`/`rc_` function-name
prefix. This is a hand-copy-paste of IDUNA's own login-screen UI from one game into another,
proving the exact pattern the founder is pointing at: a real, generically-needed piece of UI (every
networked game in this monorepo needs a login screen talking to IDUNA) that gets rewritten by hand
per game instead of existing once.

## What's already been tried, and its real, current limit

This isn't the first time this monorepo has reached for shared infrastructure — worth crediting,
and worth being honest about why it hasn't actually solved the problem yet:

- **`packages/goldenband`** (the `.gband` animation format + sampler) is the closest thing to a
  real "shared engine package" precedent — used identically by `REDGARDEN` and `ECOWAR`. Checked
  directly: it is **not** a submodule, symlink, or build-system dependency — `gband.c` is
  byte-identical (`5e2241e6…`) only because it hasn't been touched since the fork. It will diverge
  the exact same way `hmac_sha256.h` already has the moment either repo needs to change it.
- **PARENA's own `stdlib/sdl2.prn`** (482 lines, real and already built, `2026-08-26`) is real,
  genuine progress in the right direction — its own header comment says it exists specifically
  because "every call below is a real, grepped SDL2 call from BRAWLPIT/SHANKPIT/REDGARDEN/
  PITVIPER/GoblinFoxDragon" — the PARENA stdlib team already recognized SDL2 access itself as a
  shared, built-in capability worth stopping duplication on. But it's deliberately thin: raw
  window/renderer/event primitives only, 1:1 with the SDL2 C API. It does **not** yet cover the
  layer every game actually rebuilds on top of it — mesh/shader helpers, a camera/view-projection
  convention, a login-screen widget, or a network-snapshot framing convention. `stdlib/sdl2.prn`
  is real, working foundation to build on, not the finished answer.
- **PARENA's `stdlib/<game>/` per-game namespace convention** (`redgarden/`, `ecowar/`,
  `papercraft/`, `racer/`, `gfd/`, `gta7/`, `pitviper/`, `mishri/`, `mixforge/`...) is real and
  works well for what it's for — each game's own *specific* mod logic (`abraham_fireball_mod.prn`,
  etc.). There is currently **no equivalent shared namespace for engine-generic code** — nowhere a
  matrix-math function, a login-screen widget, or a network-snapshot-splitting helper belongs that
  isn't "pick one game's folder and hope everyone else copies it by hand," which is the exact
  status quo this whole investigation is responding to.

## The real architectural question this ask raises, named honestly rather than assumed away

Every existing PARENA-mod integration in this monorepo (`REDGARDEN`'s tree-passive/bloodflower/
build-template mods, `ECOWAR`'s card-effect mod, `SPIDERBEETLE`'s battery-UI functions) follows
one pattern: **PARENA compiles a small, specific logic function; a hand-written C (or Java) host
engine calls into it.** "PARENA mod is the trigger, host does the real work" — this repo's own
words for it, everywhere it's used.

Pulling `mat4`/mesh-upload/shader-compile/window-management **into** PARENA is a different, bigger
kind of ask than that pattern — it's asking PARENA to own pieces of the *host* itself, not just
supply logic the host calls into. `stdlib/sdl2.prn`'s own real existence proves this direction is
plausible (PARENA's C emitter is mature enough to wrap real SDL2 calls today) — but "wrap
individual SDL2 calls" and "own mesh upload / shader compilation / a camera convention that every
game's own render loop depends on for every frame" are not the same size of commitment, and
conflating them would repeat exactly the "we said PARENA could do X, then found out X meant
something much narrower" mistake `DEADWEIGHT/NORTHSTAR.md`'s own PARENA-Java-emitter audit already
caught once this same week. Named here so it isn't assumed away: **this needs its own real design
pass** (does a game's C `main.c` call INTO PARENA-emitted render helpers, or does PARENA-emitted
code become the actual render loop with the C host reduced to a thin OS/window shell?) before any
of the mesh/shader layer gets built — not blocking the lower-risk extractions below, which don't
raise this question at all.

## Proposed shared namespace: `PARENA/stdlib/engine/`

A new namespace, sibling to the existing per-game ones, for code that belongs to **no** specific
game because every game needs it. Same "thin wrapper around a real, grepped API, not a redesigned
one" discipline `stdlib/sdl2.prn` already set. Grounded in the actual function signatures found
in this investigation, not invented:

- **`stdlib/engine/mat4.prn`** — the lowest-risk, highest-confidence extraction in this whole
  list: `mat4.h` is *already* byte-identical across three independent repos today. Port
  `mat4_identity`/`mat4_multiply`/`mat4_translate`/`mat4_scale`/`mat4_rotate_y`/
  `mat4_perspective`/`mat4_orbit_view` directly — there is no design question here, only a port.
- **`stdlib/engine/hmac.prn`** — `hmac_sha256.h`'s own real HMAC-SHA256 implementation
  (IDUNA connect-ticket verification, used identically in spirit everywhere it appears even
  though the actual bytes have already drifted) — one real, audited implementation instead of N
  silently-diverging ones is a genuine security win, not just a maintenance one.
- **`stdlib/engine/http_client.prn`** — `http_client.h`'s own real minimal HTTP client (talking
  to IDUNA) — same reasoning as `hmac.prn`.
- **`stdlib/engine/net_snapshot.prn`** — the `NetHeader` + typed-payload-packet + MTU-aware
  multi-packet-split convention `ECOWAR`'s own `PACKET_ARENA_SNAPSHOT_HEROES`/
  `PACKET_ARENA_SNAPSHOT_OBSTACLES` split already proved out for real, twice, in this monorepo's
  own recent history (S170-193, S370-04) — currently a pattern every new game re-derives by hand
  from scratch (`DEADWEIGHT/docs/PHASE_D2_SERVER_AND_ACCOUNTS.md` already has to say "copied
  directly, not redesigned" about this exact thing, because there's nowhere to actually import it
  from yet).
- **`stdlib/engine/mesh.prn`/`shader.prn`** — the `upload_mesh`/`compile_shader`/`draw_mesh`
  layer, once the architectural question above is actually resolved — the highest-value single
  extraction (it's the one every 3D client has independently rebuilt) and correctly the one that
  needs its own real design pass first, not the one to start with.
- **`stdlib/engine/login_ui.prn`** — the `draw_login_screen` pattern, most naturally built as a
  real widget on `EOSUI-NORTH`'s own Option C styling layer (`EMILY/docs/EOSUI_NORTHSTAR.md`) once
  that exists, rather than as a bespoke immediate-mode function ported as-is — this is the one
  extraction target that two separate scoping docs (this one and EOSUI's) both converge on
  independently, worth noting as real, corroborating signal rather than coincidence.

## Phased plan (not yet built)

- [ ] **P1 — the zero-controversy extractions.** `mat4.prn`, `hmac.prn`, `http_client.prn`: small,
  self-contained, already proven near-identical across real repos, no architecture question to
  resolve first. Real acceptance bar: at least one *existing* game (candidate: `WEAKNIGHT_
  BEDROCK_RACERS`, since its own `mat4.h` is a pure hand-copy with nothing else depending on the
  old file's exact identity) retrofit onto the shared version, proving it's a real drop-in, not
  just a port that compiles in isolation.
- [ ] **P2 — `net_snapshot.prn`.** Codify the MTU-split snapshot-packet convention for real, with
  `DEADWEIGHT`'s own D2 phase (`docs/PHASE_D2_SERVER_AND_ACCOUNTS.md`) as the first real NEW
  consumer — it needs this pattern anyway and is being built from scratch this exact week, the
  ideal low-risk proving ground before asking any *existing*, already-shipped game to adopt it.
- [ ] **P3 — the mesh/shader/render architecture design pass.** Resolve the open question named
  above (PARENA-emits-helpers vs. PARENA-owns-the-render-loop) before writing any code — the same
  "real fork named explicitly, not guessed at" discipline this whole doc has followed throughout.
- [ ] **P4 — `login_ui.prn`, jointly with `EOSUI-NORTH`'s own Option C.** Not this repo's sole
  deliverable (same "gated on a separate cross-repo doc" relationship `DEADWEIGHT`'s own D5 phase
  already has with EOSUI) — flagged here so whoever picks up EOSUI's Option C work knows a real,
  concrete, already-duplicated second use case (login screens, not just shop panels) is waiting.
- [ ] **P5 — opportunistic retrofit of existing games**, explicitly **not required** for any of
  the above to ship: once `stdlib/engine/` exists and DEADWEIGHT proves it out, `REDGARDEN`/
  `ECOWAR`/`BRAWLPIT`/`PAPERCRAFT`/`WEAKNIGHT_BEDROCK_RACERS` can each adopt pieces of it
  whenever it's convenient (a bugfix touching `hmac_sha256.h` is a natural moment to swap that one
  file for the shared version instead of re-patching the local copy) — never a mandated migration
  that risks destabilizing a shipped, working game for the sake of tidiness.

## Non-goals

Not a rewrite of any existing game's engine. Not a requirement that `SHANKPIT`/`GoblinFoxDragon`
(a genuinely different, Go-server lineage) converge onto the same `stdlib/engine/` as the C/SDL2
arena-lineage games — that would be forcing two structurally different problems into one answer,
the same mistake this whole doc opened by naming in "SHANKPIT isn't the engine" itself. If SHANKPIT/
GFD's own Go-side duplication (across `packages/protocol`, `packages/rts`, `server/`, `services/`)
turns out to be real once someone looks, that's a real, separate investigation for that lineage,
not assumed to be solved by this one.
