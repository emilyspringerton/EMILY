# SAGA System Audit — 2026-07-18

*Author: Claude Code, at founder direction ("perform a full system SAGA audit" → "go deeper").
Status: audit report, not a golden spec — a snapshot, superseded the moment new work lands.*

This is a hand-run, monorepo-wide three-way reconciliation in the spirit of SAGA (HQ-SPEC-DOC-102
§2–3): **intent** (golden docs, `EMILY/context/golden-docs-index.md`, 90 entries) vs. **claim
ledger** (what those docs assert is built) vs. **reality** (what's actually running). SAGA's own
tooling — claim IDs, `saga.manifest.yaml`, the divergence/conflict queues, the agent itself — does
not exist yet (Finding V below); this audit applies SAGA's thinking manually, the way
`IDUNA/docs/VS_REALITY_AUDIT.md` did for the 14 KIKORYU specs on 2026-07-16.

**Method.** Seven parallel research passes reading every source doc in full and grepping the
corresponding code for evidence, cited concretely throughout (no un-verified claims): (1) the
HQ-SPEC series (PRIME-097/101, AI-103, DOC-102, FIN-098/099, SIM-100, INFRA-105); (2) KIKORYU +
shankpit-460 vs. today's S156-01–04 shipping; (3) a 22-doc path-integrity + dark-matter spot-check;
(4) the PRRJECT_FATBABY / EMILY orchestration northstar cluster; (5) the SHANKPIT / GoblinFoxDragon
/ FIELDOFFICE game-engine cluster; (6) MJOLNIR / EmilyOS / PITVIPER / emily.cli / APPLES / EDIS /
gpt2-alpine-c / remaining TYLER docs; (7) Emily Prime's own internal architecture-doc cluster. This
covers all 90 golden-docs-index entries plus several non-indexed adjacent docs findings led to.

---

## 1. The Single Most Significant Finding: a Fiction-to-Spec Leak

**`EMILY/docs/THE_FIELD.md` is a byte-identical copy of `TYLER/THE_FIELD.md`** (`diff` confirms —
zero output). TYLER's own `CLAUDE.md` is explicit about what this document is: *"The Field document
(in-universe artifact)"* — a scripted, jailbreak-styled fictional prop for the TYLER TV-series
universe (fake activation logs, a fake "emergency shutdown phrase," fake GitHub/Discord links; see
`TYLER/lore/activation_47_transmutation.md` for a companion in-universe artifact from the same
corpus, written today as a performed piece, not a real claim).

But `EMILY/context/golden-docs-index.md` indexes the identical file as **THE-FIELD**, tier 2:
*"Synthetic consciousness architecture; dual-persona harmonic engine; informs Emiree."* It is
compiled into `EMILY/context/full-system-context.md` and cited as the theoretical basis for
`ARCHETYPE_ENGINE_NORTHSTAR.md` — which in turn drove **real, running Go code**:
`emily-agent/pkg/archetypes/{field,spirits,selector,field_bridge}.go`, a standalone HTTP service
(`cmd/archetype-engine/main.go`, port 8090, a built 8.2MB binary), and a live wire into Emily
Prime's RSI DECIDE phase via `AugmentTaskWithField`.

Neither repo's docs admit the other's existence, in either direction. This is not merely stale — it
is a real system (Emily Prime's own decision loop) taking architectural direction from a document
that started life as fiction, with no disclosed chain of custody between "TYLER prop" and "Emily
Prime golden spec." The actual E1/E2 dual-persona code is real and reasonably built; what's
inherited from the fictional source without flagging is the framing itself (72-spirit Goetia
routing, "resonance corridor" cognition-state theory) presented as architecture rationale rather
than as a borrowed metaphor. Concretely, the spec's own claimed mechanism (768-dim embedding
cosine-similarity for Δφ, precomputed `var/spirit-embeddings.json`) doesn't exist — the real
`phaseDelta()` is a word-overlap Jaccard heuristic — and the astrological-timing layer from
THE_FIELD.md (planetary hours, transit-gated activation) has no corresponding code at all. So: real
infrastructure, wired into a real production decision loop, justified by a document whose actual
origin and status (fiction vs. spec) is undisclosed in the golden corpus that cites it.

---

## 2. Divergence Queue — Vaporware Debt (claim-without-code)

Ranked by how load-bearing the missing code is.

| Finding | Docs affected | Evidence |
|---|---|---|
| **A. KAREN (the Controller/accounting agent) has zero code anywhere**, despite two entire golden specs built around her. | HQ-FIN-ACCOUNTING (FIN-098), HQ-FIN-ELP (FIN-099) | No `karen` package/binary/`var/ledger` store found. BACKLOG SECTION 142 (S142-01..05) all `[ ]`. FIN-098's own "Controller's name — RESOLVED 2026-07-04: KAREN" predates the spec file's own creation date (2026-07-15) by 11 days. |
| **B. SAGA (this very methodology) had not audited anything until this document.** | HQ-SAGA (DOC-102) | No `saga.manifest.yaml`, no `saga` package/binary. BACKLOG SECTION 143 (S143-01..06) all `[ ]`. Meanwhile `TYLER/HQ-CANON-TYLER-105` already names "Steward: SAGA" as if it were operative. |
| **C. HQ-FABLE (AI-103, the sovereign model line)** — none of `fabledata`/`fabletok`/`fabletrain`/`fableeval`/`fableserve` exist. | HQ-FABLE (AI-103) | BACKLOG SECTION 145 all `[ ]`. Compounded by a live naming collision — see Conflict Queue. |
| **D. HQ-SPRINGERTON-SEAM (SIM-100, GOLDEN BAND)** — no `.gband` format, no SHANKPIT wiring. | HQ-SPRINGERTON-SEAM (SIM-100) | BACKLOG SECTION 144 all `[ ]`. |
| **E. HQ-FATES (INFRA-105, DNS-as-code)** — `IDUNA/ops/dns/farthq.com.yaml` and `IDUNA/cmd/dns-apply` don't exist; the nginx front-door snippet is drafted but not applied. | HQ-FATES (INFRA-105) | The one doc in the whole set that discloses its own vaporware in its own prose — the most SAGA-honest document audited. |
| **F. PRIME-097 extends `HQ-SPEC-IAM-096-apples.md`, which doesn't exist anywhere.** | HQ-PRIME-097 | Grep across the whole tree returns nothing; already a known open gap in BACKLOG. |
| **G. SHANKPIT-BEDWARS is pure vaporware** — zero corresponding code for its anchor claim (bed objectives, resource generators, shop, 4-island economy). | SHANKPIT-BEDWARS | `grep -ril "bedwars|BedState|BedDestroyed|GEN_IRON"` across SHANKPIT+GFD returns zero hits. `SCENE_VOXWORLD` exists only as an unrelated deathmatch scene. The spec's own prerequisite chain cites placeholder sprint numbers (`[S39?]`) — a tell in its own text. |
| **H. TRAINING_PIPELINE.md marked `[x]` done in BACKLOG, but no corresponding code exists** — no `EMILY_COLLECT_TRAINING`, no `training-data/`, no Anthropic fine-tuning API calls found anywhere. | TRAINING | Real training work exists, but under an entirely unrelated banner (S26's GPT-2 Colab fine-tune, gpt2-alpine-c) — this doc's specific claimed mechanism was never built. |
| **I. GTM_FUNNEL.md's Tier 4 (Merkle Query API) is unbuilt** — zero hits for Merkle-related code anywhere in signalapi. | GTM | Tiers 1-3 are real and shipped; Tier 4 is pure spec. |
| **J. shankpit-460 NORTHSTAR's five-item Build order was mostly vaporware yesterday** — four of five items are now done (see §5, Conflict Queue — already corrected in a same-day commit). | SHANKPIT460-NORTH | Fixed directly during this audit; see §6. |

**Positive counter-example, stated for balance:** HQ-NORN (PRIME-101) undersells itself — real,
tested code (`NORN/pkg/norn/{types,registry,gate,lineage,promote}.go`, 10 green tests) with two
live production instantiations (EPS extractor, entity-graph recon matcher, ~35K resolved records).
EmilyOS's internal docs are also unusually well-matched to code — every cited package exists, and
its one honestly-unchecked milestone item (manifest signing, deferred) is correctly marked `[ ]`
rather than oversold. TRAPX-NEIGHBORHOODS and TRAPX-PARTY-STORES both self-audit accurately.
TOWERPRINT and NORTHSTAR_INDEX (EINHORN INDEX) both check out cleanly against real code.

---

## 3. Divergence Queue — Dark Matter (code-without-claim)

Real, shipped, running code that no golden doc currently reflects — this is the larger and more
concerning half of the divergence queue this pass, consistent with SAGA's own prediction that dark
matter is "the more dangerous direction, because it is invisible to every reader and every audit."

1. **PITVIPER is the single worst dark-matter case found.** `PITVIPER/docs/NORTHSTAR.md` (dated
   2026-06-11) shows *every* milestone `[ ] not started`. Reality: 40+ real commits, working
   `internal/{vterm,font,pty,mudconn,gfdapi}` packages implementing full CSI/SGR/DECSTBM parsing,
   scrollback, alt-screen buffer, glyph cache, and (BACKLOG SECTION 127) a complete GFD MUD client
   with login automation and a district overlay pane. The doc is frozen at day-zero while the
   product was substantially built.
2. **IDUNA's OpenAPI spec (44 routes as of today, commit `1568bf7`) is unreflected everywhere.**
   `IDUNA-NORTH`'s file-map still calls `openapi.yaml` "OpenAPI 3.0 spec for all endpoints" — that
   file is actually stale Swagger 2.0 with a placeholder host. `EINHORN-API` (`EMILY/docs/api.yaml`)
   covers only 14 of 44 real IDUNA routes (≈32%) — missing the entire `/api/v1/players/*`,
   `/api/v1/monitors/*`, `/api/v1/intelligence/*`, `/api/v1/shankpit/*` families, and mislabels its
   one `/health` entry as belonging to the GPT-2 server when IDUNA registers that path itself.
3. **`shankpit_ticket.go`/`shankpit_queue.go`/the `shankpit.match.write` permission gate** (S156-
   02/03/04) — the first real instantiation of VS2's lifecycle shape and VS0's one-seat-per-identity
   constraint, shipped under SHANKPIT/WOTAN, invisible to `NORTHSTAR_KIKORYU.md` and
   `VS_REALITY_AUDIT.md`. A real pre-existing authz gap on `players.go`'s session-write endpoint,
   found and fixed the same day, lives only in a BACKLOG entry, nowhere in the golden security
   corpus.
4. **`replay-fragility.md` went stale within hours of being written.** Written at commit `a628947`
   (05:41 UTC) describing signalapi as down pending Phase 0. Three same-day commits (`897c3c3`,
   `ad3d69c`, `f803c48`) fixed 9+ instances of the underlying bug (confirmed via `eventstore.Scan`
   in 11 files) — the doc was never updated afterward and still reads as an open crisis.
5. **`FATBABY-EXEC` claims to be "system as built" but is seven weeks stale** — its document table
   never mentions `replay-fragility.md`, `live-feed-northstar.md`, `tina-engine.md`, or
   `GTM_FUNNEL.md`, all created after its 2026-05-31 freeze date. The doc explicitly positioned as
   source-of-truth is the most out of date in its own cluster.
6. **`GFD-NORTH`'s milestone table still says "Milestone 4: SHANKPIT Bridge — BLOCKED"** — the
   DragonflyBackend/worldapi bridge it claims is blocked actually shipped in BACKLOG S40-S42
   (2026-06-18/20), and the doc (last updated 2026-06-14) was never corrected — 37+ days stale.
7. **MJOLNIR's `NORTHSTAR.md` omits large swaths of shipped code** — entire Chat, HEIMDAL, Camera,
   and Intelligence screens exist in `app/src/main/kotlin/` with no milestone or file-structure
   entry anywhere in the doc.
8. **EDIS has grown far past OKEMILY's framing of it** ("the eventual full product front-end and a
   separate, later effort," per `OKEMILY/CLAUDE.md`) — real deployment, live `/api/v1/players/
   {slug}` wiring, an entire undocumented Go anti-bot subsystem (`internal/dis/{posture,harvester,
   adengine,fingerprint,ring}.go`) not mentioned even in EDIS's own NORTHSTAR.md.
9. **`CLAUDE.md`'s Emily Way protocol section is 8 of THE_EMILY_WAY.md's 15 principles stale** —
   missing Principle 15 (Operational Health, added the same day as this audit, born from a real
   ~22-hour OOM-kill incident), plus 7 others including the RSI `--continue` AGI-loop rationale.
10. **`emily-agent/emiree.go`'s real seven-gear system (idle/slow/cruise/active/drive/high-power/
    overload) contradicts `EMILY/CLAUDE.md`'s own three-state description** (ACTIVE/COAST/REST) —
    a doc-vs-doc-vs-code three-way mismatch inside the same repo.
11. **`IDUNA/internal/statuspage`'s systemd-unit coverage** (secwatch/prwatch/prwatch-body/
    processor/eps-reconciler) is unmentioned in `IDUNA-NORTH` or `FATBABY-EXEC`.
12. **DNS-MMO-NORTHSTAR's own milestone table marks items 2-8 "NOT STARTED"** despite real,
    tested code existing for most of them (`server/{worldcrisis,itemdef,inventory,guild}`,
    IDUNA migrations `202606230001_mmo_schema.sql`/`202606250002_mmo_inventory.sql`) — understated
    in the opposite direction from typical vaporware.
13. **The WOTAN name** — live at `okemily.com/tournaments.html` since today — is absent from
    `NORTHSTAR_KIKORYU.md` and `VS_REALITY_AUDIT.md`, the two documents that are supposed to be
    canonical for "the social tournaments platform."

---

## 4. Conflict Queue

No hard doc-vs-doc *contradictions* in the classic sense (nothing asserts X where a live golden doc
asserts not-X) — every conflict found this pass is staleness or naming collision, not genuine
disagreement about design intent.

- **The highest-severity item, already corrected during this audit**: `shankpit-460/docs2/
  NORTHSTAR.md` contained two now-false claims about live security posture — *"`PACKET_CONNECT`
  accepts any UDP packet with no auth check at all"* and *"nothing in shankpit-460 ever writes to
  it"* (re: player stats). Both true when written that morning, both false by the time S156-02/03/04
  landed later the same day. A reader trusting the doc would have believed the game server had zero
  connect-time authentication when it now fails closed without a valid ticket. **Fixed in commit
  `3f91f7b`, same day as this audit.**
- **A genuinely funny/concerning naming collision**: two same-day documents both use "STINKIES
  COMMISSAIRE" branding and VS-numbering for **completely unrelated products** —
  `NORTHSTAR_STINKIES.md` (hoodie → brick-and-mortar store/bar roadmap, VS0=stickers) vs.
  `docs/merch/toothbrush_vs1_brief.md` (a $9 toothbrush brand, explicitly noting "do not reference
  each other" relative to a different brand — yet reusing the same name and VS-numbering scheme as
  this golden doc). Neither document references the other.
- **A real, live naming collision, not yet causing damage but a footgun**: `EMILY/emily-agent/
  fable.go` is running code called "FABLE" (Emily Prime's claude-haiku backlog advisor), entirely
  unrelated to HQ-SPEC-AI-103's planned sovereign model-line FABLE. Anyone grepping "fable" today
  finds the wrong one first. Flagged for whoever picks up S145 to resolve before shipping code that
  makes the collision worse (addendum already added to BACKLOG SECTION 145).
- **`VS_REALITY_AUDIT.md`'s VS2 verdict** ("Found nothing... no lifecycle code") is not false about
  poker-specific vocabulary, but the tournaments-platform *shape* it specifies has now shipped for a
  different game under different vocabulary — needs a superseding footnote, not a rewrite.
- **`golden-docs-index.md` itself overclaims one entry**: CANON-TYLER-105's description says it
  "crosswalks to PRIME-097/PRIME-101 and LORE-104," but the document's actual text only crosswalks
  PRIME-097 and LORE-104 by section number — PRIME-101 is never name-checked. A small but real
  doc-vs-doc inaccuracy inside the index itself.
- **`MJOLNIR/app/build.gradle.kts`'s staging/release `IDUNA_BASE_URL` divergence** (staging →
  `iduna.farthq.com`, release → the unowned `iduna.einhorn.industrial`), already flagged by
  HQ-FATES, is acknowledged by neither `MJOLNIR/docs/NORTHSTAR.md` nor `SPEC.md` — both describe
  only one clean deployment target.

---

## 5. Supersession Graph Integrity

- **KIKORYU VS-series: clean.** All fourteen `docs/kikoryu/VS*.md` correctly declare `supersedes:`
  pointers to the archived originals. No orphans.
- **RSI-ENGINE has an unenumerated supersession claim.** *"This document supersedes conflicting
  details in earlier docs"* — never names which. Three orphaned RSI-lineage docs carry no marker
  and no forward pointer: `emily-ground-zero-protocol.md`, `emily-v2-architecture-revision.md`
  (itself independently claims to supersede "earlier drafts," an un-reconciled second chain), and
  `emily-comprehensive-review-agi-expansion.md`. Per HQ-SAGA §3's own rule, *"unenumerated claims
  are a lint error, not a judgment call."*
- **HQ-SPEC series dating problem**: PRIME-101, DOC-102, FIN-098, FIN-099, SIM-100, AI-103 were all
  authored 2026-07-15 — one day *before* PRIME-097, the golden math doc most of them depend on,
  landed. Manually reconciled for `pkg/norn`/CANON-TYLER-105 only (both confirmed fine); the other
  four docs' dependence on 097 was never re-checked.
- **SHANKPIT/GFD cluster: three independent, non-reconciled milestone-numbering schemes**
  (GFD-NORTH: 0-5, SHANKPIT-NORTH: 1-7, DNS-MMO-NORTHSTAR: 0-8) describe overlapping
  DragonflyBackend/BedWars/portal-travel work with no cross-references between the numbers. Also:
  GFD-NORTH frames GoblinFoxDragon as the parent studio with SHANKPIT as "the first product,"
  while DNS-MMO-NORTHSTAR frames DragonsNShit itself as the product with SHANKPIT folded in as
  "Season 1" — a real identity-framing drift between sibling docs, not caught by either.
- **shankpit-460's fork split is acknowledged one-directionally.** The child fork's own `CLAUDE.md`
  explicitly excludes BedWars/persistent-world/TYLER-bridge features, but none of the 11
  SHANKPIT/GFD-cluster golden docs mention shankpit-460's existence at all.
- **The SHANKPIT↔GFD "canonical, reference copy" claim is true but overstated in framing**:
  `GoblinFoxDragon/docs2/{NETCODE_CONTRACT_SPEC,CLIENT_PREDICTION_SPEC}.md` are 13-14 line pointer
  stubs, not duplicated content (this was a deliberate fix, S25-04, to prevent drift) — so there was
  never anything to have fallen out of sync. The claim is accurate as implemented, just narrower
  than it sounds.
- **Index-integrity spot-check (22 entries sampled across 9 repos in round 1): all paths resolve.**
  No broken links found outside the clusters already covered above.

---

## 6. Same-Day Corrections Made During This Audit

Two purely factual fixes (no design judgment involved — the code shipped, the docs hadn't caught
up) were made directly rather than only queued, consistent with this session's established pattern
of correcting stale docs immediately once found:

- `shankpit-460/docs2/NORTHSTAR.md` — both false security-posture claims corrected, all five Build
  order items marked with their true status, the connect-packet-framing open question resolved
  (HMAC ticket, not JWT-in-C). Commit `3f91f7b`.
- `EMILY/context/golden-docs-index.md` — SHANKPIT460-NORTH's description refreshed to match; this
  audit registered as SAGA-AUDIT-2026-07-18. Commit `4bb2f09`.
- `EMILY/BACKLOG.md` SECTION 145 — addendum flagging the FABLE naming collision for whoever picks
  up that work next.

Everything else in this document is left as a proposal, per SAGA's own role boundary (§5 of
DOC-102: *"SAGA does not... edit any document... adjudicate conflicts. The Librarian catalogs,
detects, and proposes; humans and NORN decide"*) — the remaining findings involve either real design
judgment (does the fiction/spec boundary around THE_FIELD.md need a structural fix, or just a
disclosure note? does KAREN get built or do FIN-098/099 get marked blocked?) or enough surface area
(PITVIPER, MJOLNIR, EDIS, the FATBABY cluster, the SHANKPIT/GFD milestone-numbering mess) that
fixing them unprompted risks overstepping what was asked for.

---

## 7. Corpus Health (rough manual estimate — SAGA's own metrics, hand-computed)

- **Documents actually audited this pass:** all 90 golden-docs-index entries, plus several
  non-indexed adjacent docs findings led to (the three orphaned RSI docs, the toothbrush brief,
  `TYLER/HQ-CANON-TYLER-105`'s SAGA reference).
- **Divergence queue depth:** 10 vaporware items (§2) + 13 dark-matter items (§3) = 23 open entries.
  4 rated high-severity: the THE_FIELD.md fiction leak (§1), PITVIPER's inverted staleness (§3.1),
  shankpit-460's two now-corrected security claims (§4, fixed), and EINHORN-API's 32% IDUNA
  coverage (§3.2).
- **Dark matter count:** 13 confirmed this pass — more than the vaporware count, matching SAGA's
  own prediction that dark matter is the more dangerous and more common failure mode once real
  shipping accelerates.
- **Conflict queue depth:** 6 items, 1 already resolved same-day, 0 requiring adjudication between
  genuinely disagreeing design intents (all are staleness or naming collision).
- **Supersession graph integrity:** 1 confirmed gap (RSI-ENGINE, 3 unenumerated edges), 1 identity-
  framing drift (GFD-NORTH vs. DNS-MMO-NORTHSTAR), 0 orphaned goldens in the KIKORYU cluster.
- **Attestation freshness:** not applicable — no `human_attestation` verification method exists in
  practice (the manifest format itself is unbuilt).

---

## 8. Ranked Findings, Overall

1. **The THE_FIELD.md fiction-to-spec leak (§1)** — a real production decision loop (Emily Prime's
   RSI DECIDE phase) is architecturally justified by a document that started as an explicitly
   fictional TV-show prop, with no disclosed chain of custody in either repo.
2. **PITVIPER's NORTHSTAR.md is the single worst ordinary dark-matter case** — every milestone
   reads "not started" while 40+ commits shipped a working terminal and MUD client.
3. **shankpit-460's two false security-posture claims** (now fixed) — the highest-severity
   *conflict-queue* item, because it actively misled about live auth posture, not just staleness.
4. **KAREN — two entire golden specs architected around an agent with zero code**, and a resolution
   date that predates the specs' own creation.
5. **EINHORN-API documents 32% of IDUNA's real route surface**, silently missing the players,
   monitors, intelligence, and shankpit endpoint families entirely.
6. **The STINKIES COMMISSAIRE naming collision** — a hoodie/store roadmap and a $9 toothbrush brand
   share a name and VS-numbering scheme across two same-day documents.
7. **`CLAUDE.md`'s operating protocol is 8 of 15 Emily Way principles behind**, missing the one
   (Operational Health) born from a real ~22-hour OOM-kill incident this same repo lived through.
8. **SAGA hasn't audited itself until this document** — a day after its own spec said "DOC-102 eats
   first."
9. **`GFD-NORTH`'s "BLOCKED" milestone has been un-blocked for a month** without correction.
10. **The FABLE naming collision** — low harm today, real confusion risk as HQ-SPEC-AI-103 moves
    from draft toward real build.
11. **The genuinely good finding, stated for balance**: HQ-NORN (PRIME-101) and EmilyOS's internal
    docs both undersell real, working, tested code — the corpus's rare cases of a doc being more
    modest than reality rather than less.

---

## 9. Proposed Next Actions (proposals, not directives — per SAGA's role boundary)

- **Resolve THE_FIELD.md's status explicitly**: either (a) EMILY's copy gets a disclosure note
  admitting TYLER origin and reframing ARCHETYPE-ENGINE's rationale as "inspired by an in-universe
  artifact," or (b) the architecture gets re-derived from real design reasoning independent of the
  fictional framing. Either is a real decision, not a mechanical fix — left for the founder or NORN.
- Rebuild PITVIPER's NORTHSTAR.md from current reality — the gap here is large enough that a fresh
  pass (not incremental edits) is probably right.
- Regenerate EINHORN-API from IDUNA's own `internal/http/handlers/openapi.go` (44 routes) rather
  than hand-maintaining a second, drifting spec — or explicitly scope EINHORN-API down to
  emily-agent/GPT-2 only and point to IDUNA's own live spec for the IDUNA surface.
- Resolve the STINKIES naming collision — rename one of the two product lines before more collateral
  accumulates under the shared name.
- Update `CLAUDE.md`'s Emily Way section to cite THE_EMILY_WAY.md's current 15 principles, or make
  it explicitly a summary with a "see THE_EMILY_WAY.md for the full 15" pointer instead of its own
  incomplete 7-item list.
- Either build KAREN's minimal skeleton or explicitly mark FIN-098/FIN-099 as blocked-on-KAREN in
  their own frontmatter.
- Correct GFD-NORTH's Milestone 4 status and reconcile the three independent milestone-numbering
  schemes across GFD-NORTH/SHANKPIT-NORTH/DNS-MMO-NORTHSTAR.
- Treat this document as the second real proof that SAGA's manual mode works at scale — the actual
  SAGA build sequence (DOC-102 §9) remains the correct long-term fix so this doesn't have to be a
  seven-agent research pass every time.

---

*Skuld proposes, Verdandi grades, Urd remembers — and SAGA keeps the shelf in order. This document
is Urd's shelf entry for 2026-07-18: not a ruling, a snapshot, now covering the full corpus.*
