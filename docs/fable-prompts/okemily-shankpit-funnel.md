# Queued Fable prompt — SHANKPIT funnel page on okemily.com (Phase 1)

**Status:** queued, not yet dispatched. Written 2026-07-19.
**Dependency:** none — ready to dispatch now. Phase 1 only (see phased plan below); phases 2/3 are
explicitly future work, not part of this dispatch.

**To dispatch:** `Agent({ model: "fable", subagent_type: "general-purpose", prompt: <body below> })`

---

## The phased plan (founder-approved 2026-07-19, this dispatch is Phase 1 only)

1. **Phase 1 (this dispatch):** build a long-form "wow style" funnel page for SHANKPIT on
   okemily.com. CTA buttons sprinkled throughout, all linking to
   `https://github.com/emilyspringerton/SHANKPIT/releases/tag/460` (verified live 2026-07-19 —
   this is farthq.com's current, actual CTA target, confirmed by fetching the page, not guessed).
   `farthq.com` stays up unchanged, coexisting.
2. **Phase 2 (later, not this dispatch):** swap the CTA target once a real download page, signup
   flow, or Steam wishlist link (S19-05) exists. Should be a one-line change if Phase 1 is built
   with the CTA href as a single, consistently-used value rather than hardcoded N times.
3. **Phase 3 (later, not this dispatch, explicitly undecided):** whether SHANKPIT's long-term
   canonical home is okemily.com, farthq.com, or something else. Founder's own words: "farthq
   isn't the best url for shankpit anyways, not sure what is right now." Don't resolve this in
   Phase 1 — don't retire, redirect, or de-emphasize farthq.com.

## Reality check (2026-07-19, re-verify if dispatching later)

- `farthq.com` today is a **short, single-viewport teaser**, not a funnel: neon/glitch cyberpunk
  aesthetic ("ENTER THE VOID"), a 2x2 stat-card grid (Language, Engine Bloat, Physics Chaos, Vibe
  Check — playful, not literal), one CTA button (`>./INSTALL_SHANKPIT`) linking to the GitHub
  release tag above, footer text "v460 // Pre-Alpha // Use at own risk." Full HTML fetched and
  read in full this session — use it as the tone/aesthetic seed, not something to just re-host.
  Note it also says "Language: C++/OpenGL," which is stale/jokey — the real stack per
  `SHANKPIT/CLAUDE.md` is C/SDL2 client + Go server; don't copy that line literally into the new
  page as a factual claim.
- Real, current status (verify against `EMILY/BACKLOG.md` SECTION 19 before writing marketing
  copy, don't assume this stays accurate): standalone EA build exists and works (`S19-04`, done —
  `dist/ea/` via `make ea` / `make ea-windows`, 4-player LAN/internet sessions, see
  `SHANKPIT/docs/EA_BUILD.md`). Steam Early Access is the goal but **not live** — `S19-03`
  (Steamworks account) and `S19-05` (actual Steam launch) are both still unchecked, blocked on a
  human creating a Steamworks developer account. **Do not write copy that implies Steam
  availability today.** The Dragonfly/DragonsNShit persistent-world layer (BedWars, buildable
  voxel worlds, season lineage) is real, designed, and worth describing as vision/roadmap — but is
  Layer 3 in `SHANKPIT/docs2/NORTHSTAR.md`, not yet live. Keep the "what's real today" vs "what's
  coming" distinction sharp; see `EMILY/docs/THE_EMILY_WAY.md`'s "real data over demos" value —
  the same discipline that's already keeping okemily.com's main landing page honest should govern
  this page too.
- okemily.com is static HTML/CSS/JS, no build step, no framework — see `OKEMILY/CLAUDE.md`. Keep
  this page the same way. Live deploy needs a human running `sudo rsync` + nginx reload — you are
  not expected to deploy this yourself, just build and commit it.

## Prompt body (copy from here down)

Build Phase 1 of a SHANKPIT funnel page for okemily.com: `OKEMILY/shankpit.html` (or `OKEMILY/
shankpit/index.html` if you want a clean `/shankpit/` URL — your call, match whatever pattern
`OKEMILY/blog/` already uses for its own clean URLs). This is a **long-form, scrolling landing
page** — meaningfully longer and more developed than farthq.com's single-viewport teaser, with
multiple sections and a CTA button placed in each one, not just at the top and bottom.

### Required reading first

1. Fetch `https://farthq.com/` yourself and read the full HTML/CSS — it's the aesthetic and tone
   seed (neon cyan / damage red / crit yellow palette, monospace/terminal voice, glitch text,
   grid-floor background, health-bar-styled stat cards, `>./COMMAND_STYLE` button labels). Evolve
   it into something longer, don't discard it.
2. `SHANKPIT/docs2/NORTHSTAR.md` — the three-sentence version, Layer 1 (what's real today: the FPS
   core, movement feel, combat), Layer 3 (DragonsNShit vision, described as vision not fact).
3. `EMILY/BACKLOG.md` SECTION 19 — real current status. Confirm what's checked vs not before
   writing a single status claim.
4. `SHANKPIT/docs/EA_BUILD.md` if it exists (check) — for accurate "how do I actually play this
   today" content, since the funnel should end with something real and specific, not just a
   GitHub link with no context on what the visitor is about to get.
5. `OKEMILY/index.html` and `OKEMILY/CLAUDE.md` — match the existing site's dark/light theme
   variable system (`--bg`, `--fg`, `--accent`, etc.) so this page feels like part of the same
   site, not a bolted-on skin. It's fine — encouraged — for this page to have its own bolder
   visual identity (SHANKPIT's neon aesthetic vs the main page's clean minimalism), the way a
   product page can look different from a company homepage, but it should still feel related, not
   jarring. Link back to `/` somewhere (header wordmark, like the blog posts already do).

### What the page needs

- Hero section: SHANKPIT identity, the tagline energy of farthq's page, one CTA.
- A "what it feels like to play" section — grounded in NORTHSTAR's Layer 1 description of the
  movement/combat feel, not generic FPS marketing copy. Another CTA.
- A "where it's going" section — DragonsNShit/persistent-world vision, clearly framed as roadmap
  (season lineage, buildable/destructible voxel worlds, BedWars-style modes), not present-tense.
  Another CTA.
- An honest "get it now" section — standalone EA build, what platforms, what a first session looks
  like (4-player LAN/internet, from EA_BUILD.md), Pre-Alpha framing carried over from farthq's own
  footer disclaimer. Final CTA.
- Every CTA button: same href, `https://github.com/emilyspringerton/SHANKPIT/releases/tag/460` —
  use one CSS class / one templated link so Phase 2's swap is a find-and-replace, not a rewrite.

### Where this lands

Commit to `OKEMILY/`. Update `OKEMILY/CHANGELOG.md`. Link the new page from `OKEMILY/index.html`
somewhere sensible (the footer's existing link list is the obvious place, next to the blog links).
In your final summary: confirm the CTA href count and that they're all identical, and give the
exact `sudo rsync` deploy command from `OKEMILY/CLAUDE.md` for the founder to run — don't attempt
the live deploy yourself.

---

## Prompt body ends here
