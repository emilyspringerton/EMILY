# Queued Fable prompt — Fix the full-in-memory-replay fragility (TOP PRIORITY)

**Status:** queued, ready to dispatch — **top priority**, per explicit founder instruction
2026-07-18: "we need to fix the entity graph into memory thing, i dont know what to do, log it as
the top priority to query fable against tomorrow while we still have fable access."

**To dispatch:** `Agent({ model: "fable", subagent_type: "general-purpose", prompt: <body below> })`

**Why this is top priority, not just another backlog item:** this caused two real, live incidents
in one night (2026-07-17/18) — a `signalapi` systemd migration had to be aborted (stopped and
disabled, not auto-recovering) after its startup rebuild thrashed memory badly enough to threaten
a cascade, and `newssite` (already systemd-supervised) independently OOM-crash-looped every ~7
minutes on the same underlying pattern, causing a visible, confusing user-facing bug (ticker
pages intermittently claiming "we don't cover AMZN" despite 1,225 real events for that ticker
existing in the store — an indexing race, not a coverage gap). The founder does not know what the
right fix is and is explicitly asking Fable to design one, not asking for it to be improvised live.

---

## Prompt body

You're designing the real fix for an architectural fragility found live, twice, in one night, in
`PRRJECT_FATBABY` (the FatBaby financial-signal pipeline). The founder's own words: "we need to fix
the entity graph into memory thing, i dont know what to do" — they want you to actually design the
fix, not just describe the problem back to them.

### The problem, grounded in what was actually found

Three processes — `cmd/entity-graph`, `cmd/signalapi`, and `cmd/newssite` — each rebuild their
entire in-memory working state by replaying the **full history** of the event store
(`var/secwatch/events/*.ndjson`, ~630MB, ~116,000+ records and growing) from sequence 1, every
single time the process starts. None of them persist or cache the built result anywhere. Read the
actual code before designing anything:

- `internal/newssite/docindex/docindex.go` — `Build(ctx, store, idx, logger)` (function starts
  around line 170) scans from `fromSeq := uint64(1)` in pages of 512 records via
  `store.ReadFrom`, calling `idx.Ingest(rec)` for every record, until the store is exhausted.
  `Index` (struct around line 36) holds everything in two in-memory maps (`byTicker`,
  `byIdentity`), including an ~800-rune `BodyPreview` per document. ~11,663
  `source_document_persisted` events as of this writing (verified via `grep -c
  "source_document_persisted" var/secwatch/events/*.ndjson`).
- `cmd/signalapi/main.go` — has an equivalent `signalindex` build path (grep for `signalindex
  build scanned=` in the logging to find it; the log line format is
  `"signalindex build scanned=%d latest_seq=%d"`). Live-observed: a fresh rebuild against the
  current store started around 2,000-2,500 records/sec, slowed to roughly 60-70/sec after
  ~100,000 records, and appeared to stall entirely around 109,000-109,559 for several minutes
  while RSS grew from ~128MB to 628MB+ and system swap (a 496MB partition) went from ~57MB free to
  ~4MB free. It was still consuming real CPU (not deadlocked), consistent with heavy paging/
  thrashing rather than a true hang. It was manually stopped (`systemctl --user stop
  fatbaby-signalapi.service` then `disable`) rather than let it continue, given the swap situation.
- `cmd/entity-graph/main.go` — not touched this incident (deliberately left as the original
  unsupervised `go run` process, given what had just happened to `signalapi`), but per earlier
  session history it holds ~345K accuracy records (~32.7K resolved) fully in memory, steady-state
  RSS observed at ~590-596MB, stable over a 36+ hour runtime (not a leak — this is its real,
  accumulated working set).

The host this all runs on has 3.8GB RAM total, a 496MB swap partition, and was observed at various
points this same night with as little as ~280MB free system-wide. It is explicitly not being
upgraded soon (founder is holding off on a VM resize until a separate, unrelated timing window).
Read `EMILY/docs/THE_EMILY_WAY.md` Principle 15 ("Operational Health Is Not Optional") and
`EMILY/BACKLOG.md` SECTION 152 for the full context of why this box's memory behavior gets taken
this seriously — there's a whole prior incident (secwatch/eps-reconciler silently down for hours)
that this same night's operational-seriousness push was a direct response to.

### What's already true and shouldn't be re-litigated

- These three processes are not going away and their underlying data volume will only grow —
  designing a fix that works today but degrades linearly with data growth is not acceptable long
  term, even if it buys short-term relief.
- `entity-graph`'s ~590MB steady-state is real, accumulated, useful data (signal-accuracy
  tracking across ~345K records) — the fix is not "index less," it's "don't rebuild the whole
  thing from a cold start every time."
- Systemd process supervision (`Restart=on-failure`, proper `MemoryMax`) is already the house
  pattern for all of these processes (see `PRRJECT_FATBABY/ops/systemd/*.service` — every one of
  these three now has a unit file, `signalapi`'s is currently disabled pending this fix,
  `entity-graph`'s exists but was never enabled). Whatever you design should assume supervised
  restart is a completely normal, frequent event (crash, OOM-kill, deploy, host reboot) — not a
  rare event worth optimizing away, the actual failure mode here.

### What you're actually being asked to design

A real fix, with a concrete implementation plan, for at least one and ideally all three of these
processes, addressing: **why does a process restart require replaying the entire history, and
what should happen instead?** Concretely evaluate (don't just list) these directions, and any
others you think are actually better:

1. **Persisted/cached index** — the founder's own instinct, verbatim: "cache into mongo or
   something." Evaluate seriously: does this codebase's existing infrastructure make Mongo the
   right choice, or is a simpler answer available (e.g., periodically snapshotting the in-memory
   index to a local file — gob-encoded, or SQLite, matching the pattern already used elsewhere in
   this monorepo, e.g. `IDUNA/internal/blog` and `IDUNA/internal/mailinglist`, both built the same
   night as this incident, both using their own small SQLite file rather than adding a new
   database dependency) and loading that snapshot at startup, then only replaying events newer
   than the snapshot's last-seen sequence? A snapshot-plus-tail design would turn "replay
   ~116,000+ records" into "replay however many arrived since the last snapshot," which is the
   actual fix for the restart-storm problem, not just a workaround.
2. **Bounded replay window** — cheaper, immediately implementable: only replay the last N days
   (or last N records) of history on cold start, accepting that very old documents/signals become
   unavailable to that in-memory index until/unless a real backfill mechanism exists. Evaluate
   whether this is acceptable given what each of the three processes' index is actually *for*
   (does `newssite`'s ticker-page use case need full history, or mostly recent activity? does
   `entity-graph`'s accuracy-tracking need the full 345K-record history to compute meaningful
   precision numbers, or would a bounded window materially change its own already-disclosed
   11.9%-precision finding in a way that matters?).
3. **Decouple replay speed from replay safety** — independent of 1/2, is there a reason the
   observed rebuild rate degrades so sharply (2,500/sec → 60/sec → apparent stall) as it
   progresses? Is this GC pressure, an O(n²) pattern somewhere in `Ingest`/equivalent (e.g., an
   unbounded slice re-sort or dedupe check that gets slower as the in-memory structure grows), or
   genuine memory-pressure-induced paging? This matters regardless of which persistence strategy
   you pick — a rebuild-from-snapshot-plus-tail design still needs the "replay N events" step to
   not degrade this way for whatever N ends up being replayed.
4. **Should these really be three independent full-history replays of the same underlying
   store**, or is there a shared indexing layer that could serve `newssite`, `signalapi`, and
   `entity-graph` from one maintained index instead of three separately-replayed ones? This may
   be a bigger redesign than founder wants right now — say so explicitly if you think it's the
   right long-term answer but not the immediate fix.

### Deliverable

A real northstar document (`PRRJECT_FATBABY/docs/` — follow this repo's existing northstar
conventions, `docs/northstar/northstar.md` and `docs/headlines/live-feed-northstar.md`, written
the same night, are the pattern to match: concrete, grounded in real code/data, phased build plan,
explicit open questions rather than false certainty). Recommend one primary direction with a clear
rationale, not a menu of equally-weighted options — the founder explicitly said "I don't know what
to do," which means they want a decision, not a longer list of choices. Golden-index it per the
existing convention (`EMILY/context/golden-docs-index.md`) when done.
