# Queued Fable prompt — evolve the TYLER easter egg into an ARG (YOLO, one-shot, just ship it)

**Status:** queued, not yet dispatched. Written 2026-07-19.
**Dependency:** none — ready to dispatch now.
**Founder's own framing, verbatim:** "we want to evolve TYLER into an ARG via the easter egg —
call out to fable to YOLO ONE SHOT IT (use field activation)." This is deliberately not a big
planning dispatch like the others tonight — build and ship a real first version in one pass. Use
your judgment on scope; don't ask for another round of confirmation, just make something real and
land it.

**To dispatch:** `Agent({ model: "fable", subagent_type: "general-purpose", prompt: <body below> })`

---

## Prompt body (copy from here down)

`okemily.com`'s footer already has a tiny TYLER easter egg: triple-click the copyright year and one
italic quote appears ("television as code..."). Evolve this into a real ARG hook using TYLER's
existing Field Activation mechanic — don't invent a new fictional system, use the one that's
already there and already has a public data point.

### Required reading first

1. `EMILY/docs/THE_FIELD.md` — the Astrological Goetia Activation Protocol. 72 spirits, each with
   a frequency band, amplitude, function, planetary ruler, and seed phrase. Activations are
   numbered, timestamped against real transits (e.g. "Activation #47 – 2025-10-31 03:17 UTC,
   Saturn stations direct in Pisces, spirit Zagan, seed 'Base metal screams when it remembers'").
2. `/var/www/okemily/blog/activation-114/index.html` (or fetch `https://okemily.com/blog/activation-114/`)
   — an existing guest post that already treats Activation #114 as real in-universe content,
   written in the same register you should extend. This is your tone anchor.
3. `OKEMILY/index.html` — the current easter egg implementation (search `yr-msg`, the triple-click
   handler). This is what you're evolving, not replacing wholesale — keep the site's existing
   voice/minimalism.

### What "ARG" should mean here, concretely

Pick a real, shippable mechanic — options, not a mandate, use judgment:
- The easter egg reveals an Activation entry computed against *today's actual date/time* (not
  static) — a seed phrase, frequency, spirit — something that changes if you come back tomorrow,
  rewarding people who check repeatedly.
- A breadcrumb: the revealed content points somewhere else (another hidden page, a TYLER repo
  path, a numbered Activation log page) rather than being a dead end.
- Consider whether IDUNA's blog system could serve a growing "Activation log" as its own
  discoverable path, or whether this stays purely front-end/static — your call, argue it briefly
  in your summary.

Whatever you build: it should feel like a real ARG breadcrumb (rewards curiosity, has *some*
persistence or logic beyond a static string), not just a second static quote bolted onto the
first. But keep it small enough to actually ship in one pass — this is explicitly not the moment
for a whole new subsystem.

### Where this lands

Commit to `OKEMILY` (and `TYLER` if you add canon there). Update `OKEMILY/CHANGELOG.md`. State
plainly in your summary: what you built, why you scoped it the way you did, and the exact
`sudo rsync` deploy command from `OKEMILY/CLAUDE.md` for the founder to run — don't attempt the
live deploy yourself.

---

## Prompt body ends here
