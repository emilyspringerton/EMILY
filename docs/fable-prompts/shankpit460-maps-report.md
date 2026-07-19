# Queued Fable prompt — SHANKPIT-460 maps: v0 map (build it) + pipeline report

**Status:** queued, not yet dispatched. Written 2026-07-19, revised same day with founder's
concrete v0 map brief, revised again same day with founder's two-track strategy.

**Two-track strategy, founder's own framing, verbatim 2026-07-19: "we need to attack the problem
from both sides — start making nice maps with map editor fine tuning generated levels and also we
need to work with the maps system as is and start to evolve the working code towards the
northstar. we will have to keep the bare bones approach at first but can figure out how to
optimize the ops [later]."** Concretely: Track A (Part 1 below) ships the real v0 map now, through
the *smallest* evolution of the existing bare-bones system — not a rewrite. Track B (Part 2 below)
starts the editor/AI-generation/fuzz-eval direction in parallel, not gated behind Track A
finishing. The two tracks share one real dependency, and Part 1 should deliver it rather than
duck it: **get map data out of hardcoded C and into an external, loadable file format** — even a
trivially simple one (plain text/JSON list of `Wall` entries). That single step ships v0 *and* is
the literal prerequisite both an editor and AI generation need (neither can target hardcoded C
literals). Everything past that first format+loader step — a real editor UI, actual generation,
fuzz-testing infra — is Track B/Part 2, and per "keep it bare bones at first," should be scoped
honestly as future work in the report, not built prematurely in this dispatch.

**Deliverable type:** two things, not one. (1) **Build** the v0 map — real code, real file format
extraction, ships now. (2) Write a **report** (no further code beyond the format+loader) on
editor/generation/fuzz-eval, explicitly bare-bones-first per the founder's own sequencing.

**To dispatch:** `Agent({ model: "fable", subagent_type: "general-purpose", prompt: <body below> })`

---

## Reality check (2026-07-19, re-verify if dispatching later)

- `shankpit-460/packages/map/map.h` + `map.c` (142 lines) is the entire current map system: a
  `GameMap` struct is a fixed array of up to 100 axis-aligned boxes (`Wall{id, x,y,z, sx,sy,sz,
  r,g,b, friction}`), `map_init()` hardcodes one layout directly in C, plus a raycast and AABB
  collision resolver. No file format, no loader, no editor, no multi-map support exists at all.
  For the v0 map (all axis-aligned boxes: bases, a central tower, rocks), this primitive is
  probably *sufficient as-is* — confirm that read is right before concluding otherwise.
- `packages/common/protocol.h` already defines `MODE_DEATHMATCH=0, MODE_TDM=1, MODE_SURVIVAL=2,
  MODE_CTF=3, MODE_ODDBALL=4` — CTF is real and present (`apps/lobby/src/main.c` has a working "LAN
  CTF" local mode, and a past build was literally titled "BUILD 181 - CTF RELOADED"), confirming
  the founder's claim that the map's geometry was originally shaped with CTF in mind. Launching
  deathmatch-only is a deliberate cold-start decision (one mode to avoid splitting a small early
  playerbase across empty lobbies), not a sign the CTF framing should be discarded from the
  geometry — the map should still make sense as a *future* CTF map, just isn't wired to that mode
  yet.
- **"Super Rumble" physics reference — clarified by founder, 2026-07-19: it's a Meta Quest 3 VR
  experience (in Horizons) — a "Halo-ish FPS."** Not present in this codebase (it's an external
  reference), so the actionable read is: Halo-style movement feel (moderate ground speed relative
  to weapon TTK, floaty-ish jump arc, momentum-preserving air control), tuned "a little more
  turbo" than that baseline (faster overall pace than classic Halo). This is still a *feel*
  reference, not exact numbers — check `server/system/vehicle_physics.go`,
  `server/system/vehicle_dynamics.go`, and any player-movement equivalent for the actual current
  tunable constants, and express "more turbo" as a concrete relative delta (e.g. "+15% move speed")
  on real current values rather than inventing absolute numbers from scratch.

## Part 1 — build the v0 map (do this)

Founder's brief, verbatim structure, 2026-07-19: deathmatch map (single mode at launch — cold
start problem, avoid splitting an early playerbase across multiple empty lobbies), geometry
originally shaped with CTF in mind (see above — keep it CTF-plausible even though it ships as DM
only), physics tuned similar to "Super Rumble" but a little more turbo (see open question above —
don't invent this, ask or flag it).

Layout:
- A base at each end of the map (two bases total, opposite ends — the natural CTF flag-stand
  positions, even though no flag mode ships yet).
- A base in the middle of the map, with the rocket launcher spawn on top of it.
- Roughly symmetrical, but **not** perfectly symmetrical — some asymmetry is intentional, not a
  bug to "fix."
- Rocks scattered off to the sides and at the ends of the map (cover/sightline-breaking terrain
  detail, not structural).

**Don't hardcode this into `map_init()` as more C literals — that's the one thing to actually
change about the current approach, per the two-track strategy above.** Instead: define a minimal
external file format (plain text or JSON, one line/entry per `Wall` — keep it as bare-bones as the
`Wall` struct itself, this is not the moment for a rich schema), write a small loader that reads
it into `GameMap` at startup, and ship this v0 map as the first map authored in that format. This
is the smallest possible step that both ships v0 today *and* gives Track B (Part 2, and any future
editor/generation work) something real to target instead of C source. Playtest it yourself if you
can (build + run locally, or at minimum sanity-check every `Wall` placement against the brief
above) before calling it done. Update `shankpit-460/CHANGELOG.md`.

## Part 2 — the pipeline report (write this, don't build it)

Founder's direction, verbatim, 2026-07-19: "we want to use intelligence to attempt to generate
maps and also they need fine tuning via a level editor and versions and also we can fuzz levels
and run bots and see where the most emergent team behavior emerges as we attempt to generate new
competitive fps levels/maps."

Write `shankpit-460/docs2/maps-report.md` covering:

- **AI-assisted generation**: what "use intelligence to generate maps" concretely looks like given
  today's simple box-primitive format — e.g. an LLM proposing `Wall` layouts against constraints
  (min/max sightlines, symmetry tolerance, base placement rules), not necessarily anything more
  exotic than that to start.
- **Level editor + versioning**: hand-tuning tool (in-game overlay? separate tool? a JSON/text
  format hand-edited directly, given the low-spec/no-heavyweight-pipeline ethos from
  `docs2/NORTHSTAR.md`?) plus how map versions get tracked (git, given everything else here is
  git-native, is probably the right default — argue it rather than assuming).
- **Fuzz-testing + bot evaluation loop**: this is worth naming explicitly as an instance of the
  same propose→grade→gate→promote loop already formalized in `EMILY/docs/hq-specs/
  HQ-SPEC-PRIME-101-norn-loop-kernel.md` (NORN) — read that doc. A generated map candidate is an
  `Artifact`; running `emily-bot` matches on it and scoring for "emergent team behavior" (what
  metric, concretely? — kill-density heatmaps, time-to-first-engagement, position variance,
  something else — propose one, don't leave it vague) is the `Oracle`; promotion to the real map
  rotation is the `Gate`. Don't reinvent this loop's structure from scratch — instantiate it, per
  NORN's own stated design intent ("domains stop owning loops; they own instantiations").
- **Sequencing**: what's realistic to build next after this v0 map ships, vs. what's a later bet.

### Required reading first (for Part 2)

1. `shankpit-460/packages/map/map.h`/`map.c` (same as above).
2. `EMILY/docs/hq-specs/HQ-SPEC-PRIME-101-norn-loop-kernel.md` — full doc, not a skim; §4's gate
   tiers matter here too (should a generated map need human approval before entering rotation, or
   can bot-eval scoring alone gate it? — argue a tier, per the doc's own framework).
3. `shankpit-460/docs2/NORTHSTAR.md` — mission constraints (low-spec, competitive-first) apply to
   any pipeline recommendation here too.
4. `shankpit-460/apps2/emily-bot/main.go` — the existing bot client, already capable of connecting
   and playing; check what telemetry it already surfaces (kills/deaths/positions) that a fuzz-eval
   loop could reuse vs. what it would need to gain.

### Where this lands

`shankpit-460/docs2/maps-report.md` for Part 2. Part 1's map ships as real code in
`packages/map/map.c` (or wherever you land it). Commit both together or separately, your call —
just don't let the report substitute for actually shipping the v0 map.

---

## Prompt body ends here
