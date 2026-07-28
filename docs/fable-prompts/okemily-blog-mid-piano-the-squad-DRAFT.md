# Pre-pass draft — "Mid-Piano Presents: The Squad" (OKEMILY blog)

**Status:** pre-pass complete, ready to queue for Fable final pass. Written 2026-07-28 by Claude
Code, same session as the REDGARDEN squad/node-capture fix it's grounded in.
**Voice/format reference:** the existing "Mid-Piano Presents" episodes —
`https://okemily.com/blog/mid-piano-presents-the-mark/` and
`https://okemily.com/blog/mid-piano-presents-the-new-guys/` (also queryable via
`GET /api/v1/blog/posts/{slug}`). Cold-open banter → Duck's "for legal reasons" bit → Tyler
walks the topic → hero dialogue grounds the real mechanic in-character → closing
`**TYLER:** *(to camera, or whatever the garage has instead of one)* We're Mid-Piano. <line>.`
This draft follows that shape exactly, dialogue in full below, not a summary to expand.
**Origin material:** `TYLER/just_a_duck.md` — the real "Jack's Factory" short transcript the
whole roster (and the show's own name — "the ghost is playing the piano") is built on. Re-read
before any voice edit; Duck's telekinesis/secret-agent double claim and Unicorn's deadpan are
sourced dialogue, not invented flavor.
**Publish path:** IDUNA's `blog.write` endpoint (`IDUNA/internal/blog`), same as every prior
Mid-Piano episode — SQLite-backed, renders to static HTML in `/var/www/okemily/blog/<slug>/` on
publish.

---

## What this is not

Not a request to write the episode from scratch — it's written below, in full. Fable's job on
the final pass: line-edit for voice, verify every mechanical/lore claim below against current
repo state (this session shipped fast — check the commits before repeating anything as fact),
and actually publish it via the blog API.

## Facts this draft relies on — verify before publishing

1. The actual bug and fix, REDGARDEN `ed59bc1` (S170-201/202), `apps/arena_bot/src/main.c`:
   the old node-capture "anchor" rule picked one bot per node via
   `owner_index mod ARENA_SNAPSHOT_NODE_COUNT == node_index` — coincidental, not tied to which
   node that bot was actually heading toward. If both of a node's coincidental slot-owners were
   dead/engaged/elsewhere, nobody ever anchored it, and the flock's own separation force pushed
   everyone else outside capture radius forever. Fixed by picking the anchor from whichever
   living teammates are ACTUALLY assigned to that node (lowest owner index among them).
2. The second, related fix: bots previously all independently computed "my own nearest uncapped
   node," which converged the whole team onto the same single node. Fixed with squad-scoped
   flocking (tight cohesion within a small squad, same Reynolds boid math as before, just
   scoped down) plus squads claiming DIFFERENT nodes via a deterministic greedy pass — a
   "fractal" application of the same grouping/spreading idea at two scales, per the founder's
   own framing ("add like fractal boids so we naturally split more into squads").
3. Gary's kit, `docs/HEROES_VS0.md`: the one hero on the roster with no dash/blink/gap-closer at
   all — "Gary doesn't chase, he watches from where he's standing." His R line "Slow Down, This
   Isn't a Track Meet" is real, current kit text, not invented for this episode.
4. MnM's kit, `docs/HEROES_VS0.md`: R ("Absorbing Hits Meant For Somebody Else") is a
   self-root + guaranteed-survival window — "the shell takes the hit instead of the crab
   underneath it," per the hero's own lore entry (`TYLER/multiverse_heroes.md` #114, co-written
   "by Tyler and Mid-Piano," per the founder's own real-time direction when the hero was added).
5. Duck's double claim (telekinesis + government agent) and Unicorn's TEXT CARD convention are
   both sourced directly from `TYLER/just_a_duck.md`, the original "Jack's Factory" transcript
   every hero in this roster is grounded in.
6. Loki's "I was never missing" bit referenced in passing is established in the prior episode
   ("Mid-Piano Presents: The New Guys"), not new to this draft — consistent callback, not a new
   claim needing separate sourcing.

---

## Draft post body (copy from here down, then edit)

**Title:** Mid-Piano Presents: The Squad

**By:** EINHORN_MEDIA

---

*RED GARDEN — LOADING — "MID-PIANO PRESENTS." Transcript, unedited, as recorded on whatever was
running in the garage that day.*

**TYLER:** Different kind of episode today. Not lore. Behavior.

**DUCK:** For legal reasons, before we start: whatever "behavior" means, I did not do it.

**TYLER:** Nobody's accused you of anything.

**DUCK:** I'm getting ahead of it.

**TYLER:** Okay. So. For a while now, whenever you all get put on a team and told to go take a
contested node, something's been going wrong. Not losing-the-fight wrong. Standing-right-there-
and-somehow-not-capturing-it wrong.

**UNICORN (TEXT CARD):** I noticed this. I stood on the flag for what felt like several
geological eras and nothing happened.

**TYLER:** Right. So here's what was actually going on, and I want Gary to weigh in on this one
specifically, because it's kind of his whole deal.

**GARY:** I'm listening.

**TYLER:** There was supposed to be one of you — an anchor — who'd ignore everybody else milling
around and just walk to the exact spot and stand on it. Whoever got picked for that job was
based on... which numbered slot you happened to be standing in. Not which node you were actually
heading to. Just a number, checked against another number.

**GARY:** So the anchor could be assigned to a node nobody assigned to that job was ever near.

**TYLER:** Correct.

**GARY:** And everyone else just orbited it forever.

**TYLER:** Also correct.

**GARY:** That's the least surprising thing I've heard all week. I don't chase. I've never
chased. My entire kit is "stand somewhere and don't leave," and I get told constantly that this
is unusual, borderline concerning behavior for a teammate to have. Turns out it's not concerning.
It's just correct, and the rest of you have been doing it wrong by committee.

**DUCK:** That's extremely rude and also probably true.

**TYLER:** It's fixed now. Whoever's actually closest to a node — actually walking there, not just
assigned a number that happens to match it — that's who anchors it. Real reason, not a
coincidence.

**GARY:** Good. "This isn't a track meet" was never supposed to be a punchline. It was a
strategy.

**TYLER:** There was a second problem, though, and this one's less "one guy stands still" and
more "all of you, personally." Once the anchor thing got fixed, the whole team just... walked to
the same node. Every single time. The nearest one. Together. As one guy.

**DUCK:** I resent "as one guy."

**TYLER:** You're a duck.

**DUCK:** I'm aware of what I am.

**TYLER:** Ducks flock. That's not an insult, that's a documented fact about your species.

**DUCK:** I have telekinesis and I work for the government. I do not "flock." Flocking is what
regular ducks do. I am, and I want this in the transcript exactly like this, *not a regular
duck.*

**UNICORN (TEXT CARD):** You quite literally just moved toward a node because everyone else did.

**DUCK:** That was strategy.

**UNICORN (TEXT CARD):** It was ten of you standing on top of each other while four other flags
sat completely unguarded.

**DUCK:** Coordinated strategy.

**TYLER:** That's the actual fix, though — so let's bring in the guy whose whole kit is basically
about this. MnM.

**MNM:** Yeah?

**TYLER:** Explain your R to the room. Not the rap version. The short version.

**MNM:** Somebody's gotta take the hit so the shell underneath doesn't have to. That's it. That's
the whole bit, mechanically and otherwise.

**TYLER:** Right, so — new system splits everybody into smaller squads instead of one big pile,
and each squad picks its *own* node instead of racing the other nine guys to the same one. And
within a squad, same anchor logic as before: one of you actually stands on the spot, everyone
else holds the ground around it.

**MNM:** So basically what I already do. Somebody takes the hit, somebody else gets to keep
moving.

**TYLER:** Basically exactly what you already do.

**MNM:** Then good. Glad the rest of the roster caught up to the shell.

**LOKI:** Can I say something.

**TYLER:** Go.

**LOKI:** You're describing this like it's new information — "stick with your small group,
spread out, don't all chase the same shiny thing." That's just what a smart crowd does
naturally. I've been doing the small-group thing since before there was a group to be small in.

**DUCK:** He wasn't in the group. There wasn't a group. He's saying this now.

**LOKI:** I'm saying it *retroactively.* That's different from saying it late.

**TYLER:** *(to camera, or whatever the garage has instead of one)* We're Mid-Piano. Spread out.
Still hold the flag.
