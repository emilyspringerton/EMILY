# Queued Fable prompt — SHANKPIT-460 map creation pipeline (report, not code)

**Status:** queued, not yet dispatched. Written 2026-07-19.
**Dependency:** none — ready to dispatch now.
**Deliverable type:** a written report/recommendation doc. Founder explicitly asked for a report
here, not an implementation — do not write map-format code or a level editor in this dispatch.

**To dispatch:** `Agent({ model: "fable", subagent_type: "general-purpose", prompt: <body below> })`

---

## Reality check (2026-07-19, re-verify if dispatching later)

- `shankpit-460/packages/map/map.h` + `map.c` (142 lines) is the entire current map system: a
  `GameMap` struct is a fixed array of up to 100 axis-aligned boxes (`Wall{id, x,y,z, sx,sy,sz,
  r,g,b, friction}`), `map_init()` hardcodes one layout directly in C, plus a raycast and AABB
  collision resolver. **No file format, no loader, no editor, no multi-map support exists at
  all.** This is a real, load-bearing gap, not a nice-to-have — competitive esports play (this
  fork's whole mission, per `docs2/NORTHSTAR.md`) needs a map rotation, not one hardcoded box.
- `docs2/NORTHSTAR.md`'s mission constraints apply directly to this question: **low system specs
  required** (evaluate any pipeline against low-end hardware cost, not just dev-machine
  convenience) and **competitive-first** (this isn't SHANKPIT's persistent-world/DragonsNShit
  ambition — no voxel worlds, no buildable terrain, that's explicitly out of scope for this fork).
- The parent `SHANKPIT` repo has real, more developed world code worth knowing about before
  recommending anything — `packages/world/terrain.c`/`.h` and `packages/common/block_map.h` (plus
  a Go port at `packages2/common/block_map.go`) — check whether any of that is reusable/adaptable
  for shankpit-460's much narrower needs, or whether it's tied to the persistent-world ambitions
  this fork deliberately isn't carrying forward.

## Prompt body (copy from here down)

Write a report — not code — answering: **how should shankpit-460 actually create and ship
competitive maps?**

### Required reading first

1. `shankpit-460/packages/map/map.h` and `map.c` — the entire current reality, read in full.
2. `shankpit-460/docs2/NORTHSTAR.md` — full document. Sections 1 (WOTAN/IDUNA platform reuse
   doctrine — matchmaking assumes discrete match instances; think about whether maps are
   per-match-format or global), 5 (explicit non-goals), and 7 (spatial audio backlog, added same
   day as this dispatch — not required for maps, but shows the current documentation register/
   depth this report should match).
3. `SHANKPIT/packages/world/terrain.c`/`.h` and `SHANKPIT/packages/common/block_map.h` +
   `SHANKPIT/packages2/common/block_map.go` — evaluate for reuse, don't assume either way.
4. `shankpit-460/CLAUDE.md` — the "stripped-down racecar" mission framing and the two hard
   constraints (low-spec hardware, competitive-first not persistent-world-first).

### What the report needs to cover

- **Authoring**: hand-placed geometry (like a simple level format — think Quake `.map`/Source
  `.vmf` in spirit, radically simpler in practice given the `Wall`-box primitive already exists),
  a minimal in-repo editor tool, or something else entirely. Weigh against the actual primitive
  available today (axis-aligned boxes) — recommend whether that primitive is sufficient for
  competitive FPS maps or whether it needs to grow (e.g. simple BSP-ish geometry, ramps) before
  map variety is even possible, and if so, how much that work is worth doing now vs. deferring.
- **File format + loading**: a real, versioned, parseable map file format so maps aren't compiled
  into the binary — this is the concrete unblock. Keep it low-spec-friendly (plain text/JSON-ish,
  not a heavy asset pipeline) and consistent with this fork's "stripped-down racecar" ethos.
- **Map rotation for competitive play**: how matchmaking (already live, `IDUNA/internal/http/
  handlers/shankpit_queue.go`) should pick a map — fixed rotation, random, veto/pick like real
  esports titles, or out of scope for v0.
- **Content pipeline** for actually producing several maps (who/what builds them — is this
  something Claude Code can generate reasonably from a text description, or does it need a human
  with a 3D tool, and if so which, given the low-spec/no-heavyweight-pipeline constraint).
- **Sequencing**: what's the smallest real unblock (e.g. "just get map.c reading one external file
  instead of hardcoding it, ship 2 maps") vs. what's a later refinement.

### Where this lands

Write the report to `shankpit-460/docs2/maps-report.md`. Do not implement anything. End the report
with a short, concrete "if I were doing v0 next" recommendation — this is a report meant to
unblock a real decision, not an open-ended survey.

---

## Prompt body ends here
