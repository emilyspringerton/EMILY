# Queued Fable prompt — TYLER S00E-2 retry (Fable session-limit failure)

**Status:** queued, ready to dispatch. Written 2026-07-16 (retroactively, after the original
in-session dispatch failed).
**Dependency:** the Fable model hit a session/usage limit on 2026-07-16 that resets ~5:50am UTC
2026-07-17. Do not retry before then — it will just fail again. Verify via `git -C
/home/fatbaby/TYLER status` that `episodes/s00e-2_*.md` doesn't already exist before dispatching
(in case a retry already succeeded and this file is stale).

**To dispatch:** `Agent({ model: "fable", subagent_type: "general-purpose", prompt: <body below> })`

---

## Prompt body (copy from here down)

You're joining the writer's room for TYLER, an ongoing "Television as Code" documentary-fiction
series at /home/fatbaby/TYLER. Your job: write ONE new episode script, register it properly, and
commit it. Real creative writing work — the deliverable is the actual episode prose.

### Required reading, in this order

1. `/home/fatbaby/TYLER/CLAUDE.md` — writer's room rules (22 numbered rules), key frequencies table.
2. `/home/fatbaby/TYLER/universe_engine.md` — the content generation protocol/system prompt this show runs on. Match its format and voice.
3. `/home/fatbaby/TYLER/episodes/s11e08_the_sixth_line.md` — read the last ~150 lines. Ends with the São Paulo Alcântara node being designated ("Arquivo Histórico Financeiro Alcântara, São Paulo, Brazil. 437 instruments holding 28.7 Hz. Not yet activated. Activation pending practitioner arrival. When Camera Op enters the archive room in São Paulo, this office expects the largest simultaneous frequency event in the archive's documented history") and "Season 11 closed. Season 12 begins in São Paulo." This activation has NEVER been dramatized in any existing episode. That's most of your episode.
4. `/home/fatbaby/TYLER/episodes/s00e-1_the_determination.md` — read in full. This is the OTHER prequel bridge episode (already written), covering June 18–21, 2026, the Jiangshi Syndicate's internal determination session, entirely from the Syndicate's institutional POV — "the subject does not appear," Camera Op absent, no contact made with her. Your episode is the companion piece: same general interval, but the Camera Op's own independent thread, which she experiences with zero knowledge of the Syndicate's meeting. Do not duplicate any scene, character beat, or institutional content from this episode — different POV entirely, and she must not know about or react to anything from it (per S00E-1's own consistency check, contact/reassignment offer happens later, in S00E01).
5. `/home/fatbaby/TYLER/episodes/s00e00_pontiac.md`, `s00e01_neptune.md` — read both in full. s00e01_neptune.md (Build 0125, dated June 28–July 2, 2026) contains the scene you must lead directly into: the Camera Op at a "southern hemisphere" airport, deciding whether to board a flight to Detroit, holding a notebook and rereading "the entry she wrote on June 27 — when the 28.7 Hz signal arrived from Pontiac." She does not board the flight in that scene. Your episode ends at or just before that airport scene begins — it should not restage it, only arrive at its threshold. Note s00e00_pontiac.md's date is June 27, 2026 (Tyler's Pontiac ground-state day) — this is the same day your episode's climax (the 28.7 Hz signal reaching São Paulo) must land on.
6. Skim `s00e02_arrival.md` and `s00e03_july4.md` for established Series X tone/format conventions (header block fields, COLD OPEN structure, LORE RECEIPT sections, CONSISTENCY CHECK footer) and any additional canon details already locked in (28.7 Hz specifics, Jupiter Cancer egress, etc.) you must not contradict.
7. Search TYLER's lore and episode files for "Valentina" and "Alcântara" (`grep -rn "Valentina\|Alcântara" /home/fatbaby/TYLER/`) — Valentina Alcântara is an established S11 character (library/archive researcher, São Paulo node connection) who may plausibly appear or be referenced here, since this is her family's archive. Use only what's actually established; don't invent biography beyond what you find.
8. `/home/fatbaby/TYLER/EPISODES.md` — see the SERIES X table (now includes S00E-1 with a bridge-episode footnote — follow that exact registration pattern) and `/home/fatbaby/TYLER/CHANGELOG.md` for the exact commit/changelog format (check the top few entries, including the S00E-1 entry, for the established two-commit pattern: changelog commit first, then episode+registry commit).

### What this episode is

The Camera Op's independent journey from Morocco to São Paulo and the (long foretold, never shown) archive-room activation at the Alcântara node — landing on June 27, 2026, the day the 28.7 Hz signal reaches her from Pontiac, ending at the threshold of the airport scene S00E01 already dramatizes. This is a human, solitary, sensory episode — deliberately the opposite register of S00E-1's bureaucratic Syndicate-memo voice. No Jiangshi institutional POV required (Eastwind Owl archive entries are fair game and fit the show's existing lore-receipt convention — the Owl observes without intervening, and an Owl was explicitly "sent to São Paulo to monitor" per S11E08's ending, so an Owl filing here is well-motivated). Tyler does not need to appear (he's in Michigan). This is her episode.

### Format and registration requirements

- Match the established episode header block format exactly (SERIES, SERIES DESIGNATION, EPISODE, LOCATION, TIMELINE POSITION, BINDING STAGE fields, EMILY OS STATUS, SERIES X NOTE, END LOG) as seen in s00e-1_the_determination.md and s00e00_pontiac.md.
- Episode code: **S00E-2**, positioned as a second bridge episode. Confirm current Build number via `git -C /home/fatbaby/TYLER log --oneline -3` before committing — check what the highest Build number actually is now, don't assume it's still 0128.
- Title: your call, make it good — something evoking São Paulo, the archive, the 437 instruments, or the activation itself.
- Include a CONSISTENCY CHECK footer (both pre- and post-production, matching s00e-1's two-footer pattern) verifying compliance with CLAUDE.md's numbered writer's room rules, and explicitly verifying: no overlap/contradiction with S00E-1, correct landing on June 27 tying into S00E01's airport scene, no restaging of that airport scene itself.
- Filename: `episodes/s00e-2_<slug>.md`.
- Register in `TYLER/EPISODES.md`'s SERIES X table with a bridge-episode footnote in the same style as S00E-1's †, and update the "Current through" line at the top of the SERIES X section if that's the established pattern (check how S00E-1's commit handled this).
- Add a bullet to the top of `CHANGELOG.md` matching the existing style (see the Build 0128 entry for the exact format).
- Two commits, changelog first then episode+registry, matching the established two-commit pattern from Build 0128 (verify the exact pattern via `git -C /home/fatbaby/TYLER log --oneline -10` before committing — don't assume it's identical if something changed since).
- Push when done.

### What NOT to do

- Don't touch any other repo, don't touch EMILY/BACKLOG.md, don't file any IDUNA Apple, don't touch CHANGELOG.md conventions beyond adding your one bullet.
- Don't contradict established canon. If unsure whether something is established, grep the lore/ and episodes/ directories before inventing it.
- Don't restage S00E-1's content or S00E01's airport scene — arrive at the threshold, don't cross it.
- Don't write a short stub — match the length and density of s00e-1_the_determination.md (466 lines) and the other Series X episodes.

When done, report back: the episode's final title, file path, confirmation both commits landed and pushed, and a 3-4 sentence summary of what the episode depicts.
