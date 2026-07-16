# Queued Fable prompt — Act III synthesis (reconcile Act II + KIKORYU audit + Sections 147-150)

**Status:** queued, not dispatched. Written 2026-07-16 after two dispatch attempts both failed on
`529 Overloaded` (server-side, transient — not a session limit, not a content problem). Per founder
instruction, **not being auto-retried** — the intent is documented here so it can be picked up
deliberately later, by a human or a future session, rather than looping on retries.

**Dependency:** none outstanding — everything this prompt reads already landed as of 2026-07-16
(Act II sections 141-150 in `EMILY/BACKLOG.md`, the KIKORYU `VS_REALITY_AUDIT.md` + `docs/kikoryu/*.md`,
`gpt2-alpine-c/docs/TOWERPRINT.md` + `docs/reference/vector_cache.md`). Ready to dispatch as-is
whenever someone decides to.

**To dispatch:** `Agent({ model: "fable", subagent_type: "general-purpose", prompt: <body below> })`

---

## Prompt body

You're doing a cross-document synthesis and prioritization pass for EINHORN_INDUSTRIAL — real backlog-writing work with a real deliverable (new BACKLOG.md sections, properly sequenced), not a memo about the plan.

### Context: a lot landed in one session, none of it reconciled against the rest

Today (2026-07-16), in order:
1. Six `HQ-SPEC-*.md` architecture drafts (NORN loop kernel PRIME-101, FABLE model line AI-103, SAGA doc curation DOC-102, accounting FIN-098, ledger platform FIN-099, Springerton Seam/GOLDEN BAND SIM-100) were registered and turned into **"Act II"** — `EMILY/BACKLOG.md` SECTIONS 141-146, a NORN-first dependency-ordered build plan.
2. Fourteen founder-written 2020-era IDUNA specs (`vs0.md`-`vs13.md`, the original KIKORYU MMO vision) got a full SAGA-style reality-vs-documentation audit: `IDUNA/docs/VS_REALITY_AUDIT.md` + fourteen superseding docs in `IDUNA/docs/kikoryu/`. Verdict: VS0's identity gate is live-but-undocumented, several VS's were reincarnated elsewhere under different names (VS7≈IDUNA's real M2M agent model, VS13≈DragonsNShit's `provenance_chain`), some are superseded-by-different-reality, and **VS2 (tournaments) + VS9 (reputation) + VS10 (scoreboards) are the clean, prerequisite-satisfied, near-term-buildable core of a new "social tournaments platform" direction** — this never got its OWN backlog section, only a "suggested follow-ups" list inside the audit doc (search `VS_REALITY_AUDIT.md` for "Suggested backlog follow-ups" — six items, none filed as backlog yet).
3. `EMILY/BACKLOG.md` SECTION 147 — Apple enrichment (GPT-2 "fingerprint" via `gpt2-alpine-c/pkg/towerprint`, ported from a recovered 2020 personal repo; model-fingerprint provenance; astrology/transit data, source undecided). S147-01/03/05 are DONE (Apples #9915, #9910); S147-02 (the actual emily-agent worker that generates + PATCHes) is still open; S147-04 (astrology data source) is explicitly open.
4. SECTION 148 — a bigger, more sensitive arc: wiring the existing `/chat` endpoint to actually log to the IDUNA Apple ledger (currently gitignored local files, not audit trail — this specific item, S148-00, is small and independently buildable), then a personal writing corpus, a GPT-2 model fine-tuned to predict the founder's own typing, and a firehose ingestion pipeline explicitly including historical Slack data from her administering *other people's* projects/teams. Gated on an undispatched Fable design pass (S148-01) that has to take a real position on the "Slack history includes teammates' words, not just hers" governance problem before anything else in the section is buildable.
5. SECTION 149 — Emily Prime's email operational fabric: AM/PM status digests, founder-directive intake by email, Q&A round-trip by email, MJOLNIR push receipts by email. Infrastructure (`emily-agent/gmail.go`'s `GmailClient`) already exists and is fully coded but dormant — no `GMAIL_*` credentials configured (S149-01, Human Unblock Queue). Explicit, non-negotiable design boundary: email never satisfies biometric human-in-the-loop approval gates.
6. SECTION 150 — towerprint training-data integration (S150-01, augment the corpus with squish/tower/gematria instruction pairs) and a real vector-cache component (S150-02, port target is a founder-supplied reference implementation at `gpt2-alpine-c/docs/reference/vector_cache.md` — FAISS + local sentence-transformer embeddings + Merkle-style hashing; this answers "local vs frontier embeddings" in favor of local). S150-01..04 are all still open — none started.
7. S146-01..04 (the EPS-headline snapshot pipeline feeding FABLE's E0) are now DONE (Apple #9913, gpt2-alpine-c commit `458756c`) — verified live against the real (currently empty-of-confirmed-cases) EPS oracle store.
8. Separately: a data-resale-readiness audit of PRRJECT_FATBABY found its core ingestion (`secwatch`/`processor`/`prwatch`/`prwatch-body`) had been down 29 days; it's been restarted (Apple #9912) and is actively catching up on a month of backlog. Not itself part of this synthesis, but relevant context — FABLE's flagship EPS-headline task depends on this pipeline actually running.

None of items 2-6 have been checked against each other or against Act II for conflicts, redundant infrastructure, or missing sequencing. That reconciliation — applying SAGA's actual *thinking* (three-way match: intent vs. claim-ledger vs. reality; conflict detection between documents that assert incompatible things; nothing gets silently dropped, only explicitly superseded) and NORN's actual *thinking* (what's actually promotable to active-backlog now vs. still draft; what's the dependency order; is there a common kernel these should all instantiate against rather than reinventing) — is your job. Neither NORN nor SAGA are built yet (they're `SECTION 141`/`143` backlog items themselves, still software-shaped intentions) — you're applying their design principles by hand, the same way the VS0-13 audit applied SAGA's thinking before SAGA existed as running code.

Don't go looking for `HQ-SPEC-PRIME-097` — if the founder has found and provided it, it'll be registered in `EMILY/context/golden-docs-index.md` already (check); if not, don't chase it, use PRIME-101 §2's own summary of the math it depends on.

### Required reading, in full, before writing anything

1. `EMILY/BACKLOG.md` — search "STRATEGIC PRIORITY ORDER" (Act I, untouchable), the "ACT II — HQ-SPEC ARCHITECTURE ARC" preamble and SECTIONS 141-146, then 147, 148, 149, 150 in full. Also `EMILY/docs/THE_EMILY_WAY.md` Principle 14 (Continuity Report) and the reports already in `EMILY/continuity/` for how this session's work is already being tracked — match that level of honesty about what's done vs. open.
2. All six `EMILY/docs/hq-specs/HQ-SPEC-*.md` files (you may skim ones you already have strong context on from their BACKLOG sections, but read PRIME-101 and DOC-102 in full regardless — you're applying their thinking directly).
3. `IDUNA/docs/VS_REALITY_AUDIT.md` in full, plus every file in `IDUNA/docs/kikoryu/` (14 docs — the full-rewrite ones especially: VS0, VS1, VS2, VS7, VS9, VS10).
4. `gpt2-alpine-c/docs/TOWERPRINT.md` and `gpt2-alpine-c/docs/reference/vector_cache.md`.
5. `EMILY/docs/fable-prompts/fable-next-backlog.md` (queued/dispatched Fable work already in flight — don't duplicate or contradict what's already queued there: the IDUNA front-door-funnel prompt is ready to dispatch, TYLER S00E-2 retry is ready pending a Fable reset, S148-01's governance design pass is not yet written as a standalone prompt).

### What to actually find and fix (the SAGA-style reconciliation)

Look hard for these, they're real possibilities given how fast this landed:
- **Redundant infrastructure asks.** Does S147's "model fingerprint" field duplicate provenance work FABLE's `fabledata` (§4a) already specifies? Does S148's firehose ingestion reinvent anything `fabledata`'s snapshot/provenance pattern (or S146's newly-built snapshot manifest v1, which is now real, working code, not just a spec) already solves generically? If either is true, don't just flag it — decide which is authoritative and say so, the way the mmo.go/VS13 finding got resolved (reincarnated-elsewhere, cite both, move on).
- **Missing NORN instantiation rows.** PRIME-101 §6 is a living table of "who instantiates the kernel." Does the tournaments engine's NLHE lifecycle state machine want to be one (proposer=game logic, oracle=settled hand outcomes)? Does VS9's reputation/credence signal capture? Does S150-02's vector cache (S150-03 already flags this as a to-do, not yet done — you can do it now since S150-02's reference design exists)? Write the rows if you conclude they belong.
- **Sequencing conflicts.** Act II is explicitly NORN-first (S141 before anything that cites PRIME-101). Does the tournaments-engine work being promoted here also need to wait on S141, or can VS2's core state machine genuinely ship NORN-independent (a card game doesn't obviously need an improvement-loop kernel to exist) with only its *later* self-improvement hooks waiting on NORN? Take a real position. Similarly: S146 is now DONE and produces real snapshots — does that change S145's own E0 sequencing (was E0 waiting on S146, and if S146's done, is E0 now unblocked)?
- **Governance consistency.** S148's not-yet-dispatched design pass has to solve the Slack/teammates problem; does FABLE's `fabledata` license-class taxonomy (`own-exhaust`, "anything else enters through an explicit, logged decision" per AI-103 §4a) — which S146 already implements in real code now — give S148-01 a ready-made framework to inherit rather than invent from scratch? If so, say that explicitly (you may refine S148's existing items if this genuinely changes their shape).

### What to produce

1. **"ACT III" preamble + new section(s)**, appended to `EMILY/BACKLOG.md` following the exact
   existing conventions (see the "ACT II" preamble block for the format: relationship-to-Act-I
   statement, sequencing logic, known-gap acknowledgment). Act III's scope is: the tournaments-
   platform build-out (turning the audit's six "suggested follow-ups" into real, numbered,
   dependency-sequenced backlog items — tournament registry + NLHE engine + lifecycle state
   machine, scoreboard projections + read-only API, conduct/reliability signal capture, VS1's
   session-revocation gap, the stale `app.js` endpoint repair, the NORTHSTAR.md supersession pass)
   PLUS the cross-cutting reconciliation work above (any redundancy fixes, new NORN table rows,
   sequencing decisions, governance inheritance for S148). Number sections starting at whatever's
   next after 150 (`grep -n "^## SECTION" EMILY/BACKLOG.md | tail -3` to confirm — don't assume).
2. **Refine, don't duplicate, existing sections 147/148/149/150** if your reconciliation work changes
   their shape (e.g. adding a note to S148-01 that fabledata's license taxonomy is the framework
   to inherit, or noting S146's completion unblocks something in S145). Small, surgical edits to
   existing item text — not rewrites of what's already there.
3. If you find and resolve a genuine redundancy (not just a superficial overlap), add a short note
   to whichever spec/section loses the duplication saying so and pointing to the authoritative one
   — the same "supersedes, cite it, don't silently drop it" discipline the VS0-13 audit used.
4. A brief closing summary section (or just clear prose in your final report) of what you found:
   what got reconciled, what NORN table rows you added, what sequencing calls you made.

### What NOT to do

- Don't touch Act I's STRATEGIC PRIORITY ORDER or reorder/renumber any existing section.
- Don't chase or fabricate HQ-SPEC-PRIME-097.
- Don't write code — this is a synthesis/planning task.
- Don't dispatch the already-queued Fable prompts yourself or the undispatched S148-01 design pass — reference them, don't send them.
- Don't touch `IDUNA/docs/kikoryu/*.md` content itself (that's the audit's territory, already
  landed) — only `EMILY/BACKLOG.md`, and only add a pointer/note to a kikoryu doc if genuinely
  necessary (e.g. a NORN-table cross-reference), never rewrite its substance.
- **Do file an IDUNA Apple this time** (unlike prior dispatches this session) — this is substantial
  planning work and the session's own audit pass just found and corrected several places where
  Fable-dispatched work landed without one. `emily apples post -t completion -repo EMILY "<title>" "<body>"`
  once you're done, citing what you built.

### When done

Commit `EMILY/BACKLOG.md` (and `EMILY/docs/fable-prompts/fable-next-backlog.md` only if you
touched it, which you shouldn't need to) with a descriptive message, push. File the Apple. Report back: the new
section numbers/titles, every genuine redundancy or conflict you found and how you resolved it,
any new NORN §6 table rows you added (and where), the Apple ID, and confirmation everything's committed and pushed.
