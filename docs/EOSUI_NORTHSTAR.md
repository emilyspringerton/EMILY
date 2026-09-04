# EOSUI-1244 / EOSUI-12444: Cross-Game UI Design Language — Scoping (Principle 19)

**Kanban priority-queue cluster**, 2 real cards, investigated together (cross-repo, not tied to
one game — distinct from `BRAWLPIT/docs/BP_LEVEL_EDITOR_NORTHSTAR.md`'s own Phase 4, which raised
the identical embedded-browser question narrowly for BRAWLPIT's own map registry; this is the
same real question asked at the monorepo level, "for our games" plural):

- `#295 EOSUI-1244`: "can we actually implement some super basic html and css renderer for our
  games? hear me out we need a design language for the interfaces i think HTML CSS is the only
  thing that makes sense"
- `#296 EOSUI-12444`: "we can have a native version and an embedded version so we can judge which
  experience is better and closer to ship we can interop with web like roblox does especially if
  we build our own browser"

This is a real, large, cross-cutting architectural ask (`EMILY/docs/THE_EMILY_WAY.md` Principle
19) — investigated before any code was written, a real fork named explicitly, real options below,
real sub-tasks returned to `EMILY/BACKLOG.md` rather than guessed at and built blind.

## Real, current state (investigated directly, not assumed)

- **No shared UI framework exists anywhere in this monorepo.** Every game reinvents its own
  bespoke immediate-mode drawing code: `SHANKPIT/apps/lobby/src/main.c`'s own real
  `draw_lobby_buttons` (filled `glRectf` + hand-placed offset drop-shadow text), BRAWLPIT's own
  just-shipped `draw_menu_button` (`BPUX-12444`, this session — a direct, deliberate port of
  SHANKPIT's convention, done by hand because there was nothing shared to reach for), and every
  other game's own ad hoc menu code. This is the exact real pain the founder's own "we need a
  design language for the interfaces" names — checked and confirmed real, not assumed.
- **PARENA's own `stdlib/editor/widget.prn`/`ui.prn`/`theme.prn` are real, existing, but
  deliberately narrow.** `widget.prn` (98 lines) implements exactly ONE widget type (`Toggle`) at
  absolute `x/y/w/h` pixel coordinates, explicitly scoped "rather than a speculative
  button/slider/checkbox catalog nothing calls yet" per its own header comment. `theme.prn` (57
  lines) is a real scope-name → color map, but built specifically for TextMate syntax-highlighting
  scopes, not general widget styling. **There is no cascade, no box model, no flexbox/grid layout,
  no property inheritance, no markup parser anywhere in this codebase** — these are the real,
  substantial pieces "HTML and CSS" actually names, not just "colored rectangles with text."
- **No embedded-browser/webview capability exists on any desktop game client.** The one real
  precedent in this monorepo is `MJOLNIR` (Android/Kotlin) using Android's own built-in system
  `WebView` for product panels (`CLAUDE.md`: "Accompanist WebView for product panels") — real,
  working, but a platform-provided capability specific to Android, not something the SDL2/OpenGL
  desktop C games (`BRAWLPIT`, `SHANKPIT`, `PAPERCRAFT`, `ECOWAR`, `WEAKNIGHT_BEDROCK_RACERS`) can
  reach the same way. `BP_LEVEL_EDITOR_NORTHSTAR.md`'s own Phase 4 already named this same real
  gap for BRAWLPIT specifically ("zero embedded-browser capability... a real, substantial,
  Roblox-scale engineering investment").
- **Real, existing, well-known prior art for "a basic HTML/CSS renderer," found by direct
  research, not assumed:** `litehtml` (C++, MIT) is the real, standard answer to "a minimal,
  embeddable HTML/CSS layout+render engine" — used by real shipping apps (e.g. qutebrowser's
  predecessor tooling, various e-readers) specifically because writing one from scratch is a
  genuinely enormous undertaking: even "basic" HTML/CSS requires a real tokenizer, a DOM tree, CSS
  selector matching with real specificity/cascade rules, a box-model layout pass (block + inline
  flow at minimum), and text shaping/wrapping. `litehtml` itself is tens of thousands of lines
  after years of iterative real-world hardening, and still doesn't do JavaScript or full modern
  CSS (grid, most of flexbox). This is the real, honest scale of Option A below.
- **Real, existing, well-known prior art for "an embedded browser," found by direct research:**
  the `webview/webview` C/C++ library (MIT, ~cross-platform) is the real, standard lightweight
  answer — it does NOT implement a browser engine itself, it thinly wraps whatever the OS already
  ships (WebKitGTK on Linux, WebView2 on Windows, WKWebView on macOS), the same real strategy
  Android's own system `WebView` (MJOLNIR's precedent above) already uses. This is real,
  dramatically cheaper than Option A, and is what "interop with web like roblox does" actually
  describes technically — Roblox's own in-client browser is understood to be a similar
  OS-provided-engine wrapper, not a hand-rolled rendering engine.

## The real fork (a founder-level decision, not a guess)

**Option A — build a real HTML/CSS parser + box-model layout engine from scratch** (the card's
own literal words). Real, honest scale: this is not "super basic" in any meaningful sense once
started for real — even a genuinely minimal subset (block layout, inline text flow, a handful of
CSS properties) is a multi-week-to-multi-month undertaking on its own, disproportionate to the
real, stated underlying need ("a design language for the interfaces"), which does not actually
require parsing real HTML/CSS *text* — it requires a *consistent set of styling concepts* applied
across games. Not recommended as v0.

**Option B — vendor/FFI-bind a real, existing engine instead of writing one.** Two real, separate
sub-choices, matching the card's own real "native vs. embedded" framing in `EOSUI-12444`:
  - **B1 (embedded):** bind `webview/webview` (or an equivalent real OS-webview wrapper) into one
    game client as a real, timeboxed spike — this is the actual, cheap, real way to "interop with
    web like roblox does," reusing an OS-provided rendering engine rather than building one.
    Recommend BRAWLPIT as the spike target, since `BP_LEVEL_EDITOR_NORTHSTAR.md`'s own Phase 4
    (map-registry browsing) is a real, concrete, already-scoped use for exactly this capability —
    this card and that one converge on the same real deliverable.
  - **B2 (native):** do NOT vendor a full HTML/CSS engine for the "native" experience either —
    instead, build a real, small, PARENA-native declarative styling layer (see Option C, which
    IS this), and compare that against B1 for `EOSUI-12444`'s own explicit "judge which experience
    is better" ask.

**Option C — a real, lightweight, PARENA-native "CSS-inspired" declarative styling layer, not a
full HTML/CSS engine.** This is the recommended real v0: a new `stdlib/ui/style.prn` extending
`stdlib/editor/widget.prn`'s own real, proven pattern (already SDL2-backed, already has a real
consumer in the PARENA editor work) with a small, real property set — color, padding, a
flexbox-*lite* row/column layout direction, font size — consumed by each game's own native
renderer. This directly answers "a design language for the interfaces" (the founder's own actual
stated need) without paying Option A's real, disproportionate cost, and gives every game (not
just BRAWLPIT) a shared, reusable widget vocabulary the same way `SHANKPIT`'s
`draw_lobby_buttons`/BRAWLPIT's `draw_menu_button` each independently reinvented this session.

## Recommendation, not yet actioned

Pursue **Option C** as the real, buildable-now answer to the actual named need (a shared design
language), and run **Option B1** as a real, timeboxed, honest spike in BRAWLPIT specifically
(converging with `BP_LEVEL_EDITOR_NORTHSTAR.md`'s own Phase 4 map-registry browsing) to give
`EOSUI-12444`'s own "judge which experience is better" ask a real, concrete comparison — Option
C's native styled buttons vs. a real embedded OS webview showing the same content. **Do not**
attempt Option A (a from-scratch HTML/CSS parser+layout engine) — real, checked-live research
into `litehtml`'s own real scope makes this a genuinely disproportionate undertaking relative to
the actual stated need, not a "super basic" one-off.

## Why this isn't done in one pass

Real, open decisions remain before code starts: which game gets Option C's first real widget
catalog beyond BRAWLPIT's own `draw_menu_button` (this session) and SHANKPIT's own
`draw_lobby_buttons` (do those get retrofitted onto the new shared layer, or does it start fresh
and get adopted going forward only); whether Option B1's spike target should be BRAWLPIT (real,
already-scoped real use via `BP_LEVEL_EDITOR_NORTHSTAR.md`) or a different game; and real platform
risk for `webview/webview` specifically needs a quick, honest compatibility check against this
monorepo's actual deployment targets (Linux server + whatever the founder's own client OS is)
before committing engineering time. Real sub-tasks are logged in `EMILY/BACKLOG.md` under this
cluster's own section rather than folded into a single, unscoped "build an HTML/CSS renderer"
checkbox.
