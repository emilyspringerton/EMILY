# SAGA System Audit — 2026-07-18

*Author: Claude Code, at founder direction ("perform a full system SAGA audit"). Status: audit
report, not a golden spec — a snapshot, superseded the moment new work lands.*

This is a hand-run, monorepo-wide three-way reconciliation in the spirit of SAGA (HQ-SPEC-DOC-102
§2–3): **intent** (golden docs, `EMILY/context/golden-docs-index.md`, 90 entries) vs. **claim
ledger** (what those docs assert is built) vs. **reality** (what's actually running). SAGA's own
tooling — claim IDs, `saga.manifest.yaml`, the divergence/conflict queues, the agent itself — does
not exist yet (see Finding G below); this audit applies SAGA's thinking manually, the same way
`IDUNA/docs/VS_REALITY_AUDIT.md` did for the 14 KIKORYU specs on 2026-07-16. That audit is one of
the clusters re-checked here, since real work shipped against it since.

**Method.** Three parallel research passes, each reading source docs in full and grepping the
corresponding code for evidence, cited concretely below (no un-verified claims): (1) the 8-document
HQ-SPEC series (PRIME-097/101, AI-103, DOC-102, FIN-098/099, SIM-100, INFRA-105); (2) the KIKORYU
VS-series + `shankpit-460/docs2/NORTHSTAR.md`, checked against today's S156-01–04 work; (3) a
sample of the remaining ~70 golden-docs-index entries for path integrity, recent dark matter, and a
supersession spot-check on RSI-ENGINE. Full untrimmed findings are preserved in the session
transcript this audit was produced in; this document is the durable summary.

---

## 1. Divergence Queue — Vaporware Debt (claim-without-code)

Ranked by how load-bearing the missing code is.

| Finding | Docs affected | Evidence |
|---|---|---|
| **A. KAREN (the Controller/accounting agent) has zero code anywhere in the repo**, despite two entire golden specs architected around her, plus references in PRIME-101's instantiation table and AI-103's task portfolio. | HQ-FIN-ACCOUNTING (FIN-098), HQ-FIN-ELP (FIN-099) | No `karen` package, binary, or `var/ledger` store found. BACKLOG SECTION 142 (S142-01..05) all `[ ]`. FIN-098's own "Controller's name — RESOLVED 2026-07-04: KAREN" predates the spec file's own creation date (2026-07-15) by 11 days — a retrofit-dated resolution for an agent that still doesn't exist 3 days later. |
| **B. HQ-SAGA (this very methodology) has not audited anything until this document** — its own Build Sequence step 1 is "DOC-102 eats first." | HQ-SAGA (DOC-102) | No `saga.manifest.yaml`, no `saga` package/binary anywhere. BACKLOG SECTION 143 (S143-01..06) all `[ ]`. Meanwhile `TYLER/HQ-CANON-TYLER-105-EPOCH-EXTINCTION.md` already names "Steward: SAGA" in its frontmatter, as if SAGA were operative — it isn't. |
| **C. HQ-FABLE (AI-103, the sovereign model-line)** — none of `fabledata`/`fabletok`/`fabletrain`/`fableeval`/`fableserve` exist. | HQ-FABLE (AI-103) | BACKLOG SECTION 145 (S145-01..06) all `[ ]`. See also Finding F below — a real naming collision, not just an empty roadmap. |
| **D. HQ-SPRINGERTON-SEAM (SIM-100, GOLDEN BAND)** — no `.gband` format, no `gb_init`/`gb_sample`, no SHANKPIT wiring. | HQ-SPRINGERTON-SEAM (SIM-100) | BACKLOG SECTION 144 (S144-01..05) all `[ ]`. |
| **E. HQ-FATES (INFRA-105, DNS-as-code)** — `IDUNA/ops/dns/farthq.com.yaml` and `IDUNA/cmd/dns-apply` do not exist; the drafted nginx front-door snippet exists but isn't applied to any live server block. | HQ-FATES (INFRA-105) | `IDUNA/ops/dns/` directory absent entirely. `IDUNA/cmd/` has only `bob-agent`, `bootstrap`, `mailing-list-unlock`. `/etc/nginx/sites-enabled/` has `okemily`, `news-okemily`, `edis` only. **This is the one doc in the set that already discloses its own vaporware in its own prose** ("Divergence, found writing this spec") — the most SAGA-honest document audited. |
| **F. PRIME-097 extends a document that doesn't exist**: `HQ-SPEC-IAM-096-apples.md` is cited as extended but is not present anywhere in the repo. | HQ-PRIME-097 | Grep across the whole tree for `IAM-096` and the filename both return nothing. BACKLOG already flags this gap, unresolved. |
| **G. shankpit-460 NORTHSTAR's own five-item Build order was mostly vaporware when written this morning — four of five are now done, but the doc still reads as a five-item backlog.** | SHANKPIT460-NORTH | See §2 below — this is the cluster with the most same-day movement. |

**One genuine positive surprise, stated for balance (SAGA's own doctrine: report honestly in both
directions):** HQ-NORN (PRIME-101) undersells itself. `NORN/pkg/norn/{types,registry,gate,lineage,
promote}.go` is real, tested (10 green tests), and has two live production instantiations — the EPS
headline extractor (`PRRJECT_FATBABY/internal/eps/norngate`) and the entity-graph recon matcher
(`PRRJECT_FATBABY/internal/entitygraph/norngate`, ~35K resolved records, real NDJSON promotion
events at `var/norn/entity_graph_rules.ndjson`). Of the eight HQ-SPEC docs, this is the only one
where "DRAFT v0" is more modest than reality, not less.

---

## 2. Divergence Queue — Dark Matter (code-without-claim)

Real, shipped, running code that no golden doc currently reflects.

1. **IDUNA's OpenAPI spec went from 15 to 44 documented routes today** (`internal/http/handlers/
   openapi.go`, commit `1568bf7`) — covering the new shankpit ticket/queue endpoints, blog,
   mailing-list, status, monitors, subscriptions, push-tokens, intelligence. **IDUNA-NORTH's own
   file-map still describes `openapi.yaml` as "OpenAPI 3.0 spec for all endpoints"** — that file is
   actually stale Swagger 2.0 with a placeholder `api.example.com` host, a *different* file from
   the live 44-route spec, which IDUNA-NORTH doesn't mention at all. `EINHORN-API`
   (`EMILY/docs/api.yaml`) only documents the emily-agent/IDUNA broker surface, not this route set
   either. **Neither golden doc reflects the real API surface as of today.**
2. **`shankpit_ticket.go`, `shankpit_queue.go`, and the `shankpit.match.write` permission gate on
   `players.go`'s `handleSessionEnd`** (S156-02/03/04, commits `f2a3c69`/`e31db51`/`43343e8`) — the
   first real instantiation of VS2's `QUEUING→STARTING→IN_PROGRESS→COMPLETE` shape and VS0's
   one-seat-per-identity constraint, shipped under SHANKPIT/WOTAN rather than under KIKORYU-VS2's
   own poker vocabulary. `NORTHSTAR_KIKORYU.md` and `VS_REALITY_AUDIT.md` know nothing of it — see
   §3.
3. **A real, pre-existing authz gap on `players.go`'s session-write endpoint, found and fixed
   today, lives only in `EMILY/BACKLOG.md` SECTION 156** — any player's own JWT could previously
   inflate their (or anyone else's) match stats via a direct POST; fixed by gating behind
   `shankpit.match.write`, granted only to the new `SHANKPIT460-SERVER` M2M agent. This is exactly
   SAGA's "dark matter is the more dangerous direction, because it is invisible to every reader and
   every audit" case: a real vulnerability existed and was fixed, and no VS1 (IAM/moderation)-level
   security-surface doc records either the gap or the fix.
4. **`IDUNA/internal/statuspage`'s systemd-unit health coverage** (secwatch/prwatch/prwatch-body/
   processor/eps-reconciler) isn't mentioned in `IDUNA-NORTH` or `FATBABY-EXEC`.
5. **The WOTAN name itself** — founder named the tournaments platform WOTAN today (`EMILY/BACKLOG.md`
   SECTION 156 header), and it's live at `okemily.com/tournaments.html` — but `NORTHSTAR_KIKORYU.md`
   and `VS_REALITY_AUDIT.md`, the two documents that are supposed to be canonical for "the social
   tournaments platform," still describe it only as an unnamed thrust. A naming/supersession gap at
   the top of the doc tree.

---

## 3. Conflict Queue

No hard doc-vs-doc *contradictions* were found (nothing asserts X where another live golden doc
asserts not-X) — the corpus's conflicts this pass are all staleness, not disagreement. The closest
thing to a conflict:

- **`shankpit-460/docs2/NORTHSTAR.md` currently contains at least two now-false factual claims**,
  both about live security posture, not design intent:
  - *"the actual game server (`apps/server/src/main.c`) never validates this token —
    `PACKET_CONNECT` accepts any UDP packet with no auth check at all."* **False as of today.**
    `verify_connect_ticket` (main.c) now runs before every `PACKET_CONNECT` slot allocation and
    fails closed if `SHANKPIT_TICKET_SECRET` is unset. This was closed twice — S156-02 found and
    fixed a second real bypass live-testing the first fix (`PACKET_USERCMD` was auto-welcoming any
    unrecognized sender, completely sidestepping the new ticket check).
  - *"nothing in shankpit-460 ever writes to it"* (re: the players/stats endpoints). **False as of
    today** — `report_match_results()` in `main.c` now POSTs kills/deaths to IDUNA on every match
    completion.
  - This is the single most actively misleading finding in the whole audit: **a reader consulting
    this doc right now would believe shankpit-460's game server has zero connect-time
    authentication, when it now fails closed without a valid ticket.** Not a contradiction between
    two golden docs, but a golden doc contradicting the ground truth — the more dangerous SAGA
    failure mode.
- **`VS_REALITY_AUDIT.md`'s VS2 verdict** — *"Found nothing. ... No tables, no handlers, no
  lifecycle code"* — is not false about poker-specific vocabulary (still zero hits for
  tournament/poker/holdem grep terms), but the *shape* VS2 specifies (the four-state lifecycle,
  one-seat-per-identity) has now shipped, just under different vocabulary, for a different game.
  Needs a superseding footnote, not a rewrite.
- **A real naming collision, not a documentation conflict but worth flagging as a live footgun**:
  `EMILY/emily-agent/fable.go` is running code called "FABLE" (Emily Prime's haiku backlog
  advisor), entirely unrelated to HQ-SPEC-AI-103's "FABLE" (the planned sovereign model line).
  Anyone grepping "fable" in the codebase today finds the wrong one first.

---

## 4. Supersession Graph Integrity

- **KIKORYU VS-series: clean.** All fourteen `docs/kikoryu/VS*.md` files correctly declare
  `supersedes: vsN.md` pointing at the archived originals (`docs/archive/kikoryu-vs-original/`),
  verified present. No orphans here.
- **RSI-ENGINE has an unenumerated supersession claim.** `EMILY/docs/emily-rsi-engine-spec.md`
  states *"This document supersedes conflicting details in earlier docs"* but never names which
  docs. Three RSI-lineage documents in `EMILY/docs/` carry no superseded marker and no forward
  pointer to RSI-ENGINE: `emily-ground-zero-protocol.md` (same-day RSI bootstrap doc),
  `emily-v2-architecture-revision.md` (itself independently claims to supersede "earlier drafts"
  — a second, un-reconciled supersession claim chaining backward but not forward), and
  `emily-comprehensive-review-agi-expansion.md` (no marker at all). Per HQ-SAGA §3's own rule,
  *"Unenumerated claims are a lint error, not a judgment call"* — this is exactly that lint error,
  manually found since the linter doesn't exist yet.
- **HQ-SPEC series dating problem, not strictly a supersession break but adjacent:** PRIME-101,
  DOC-102, FIN-098, FIN-099, SIM-100, and AI-103 were all authored 2026-07-15, one full day *before*
  PRIME-097 — the golden math doc all six either implement or math-depend on — landed on
  2026-07-16. BACKLOG's S141 note discloses this was manually reconciled for `pkg/norn` only; the
  other five docs' dependence on 097 was never re-checked after 097 landed.
- **Index-integrity spot-check (22 entries sampled across EMILY/IDUNA/SHANKPIT/PITVIPER/EmilyOS/
  gpt2-alpine-c/TYLER/PRRJECT_FATBABY/MJOLNIR): all paths resolve.** No broken links found in the
  sample outside the two clusters covered above.

---

## 5. Corpus Health (rough manual estimate — SAGA's own metrics, hand-computed)

SAGA's proposed metrics (§8 of DOC-102) require the automated tooling this audit substitutes for
by hand, so these are order-of-magnitude estimates, not machine-verified numbers:

- **Verification coverage:** not computable — no claim ledger exists yet to compute it against.
- **Divergence queue depth (this audit):** 7 vaporware items (§1) + 5 dark-matter items (§2) = 12
  open entries, 2 of them (shankpit-460's two false security-posture claims, §3) high severity.
- **Dark matter count:** at least 5 confirmed this pass, concentrated in the last 24 hours of real
  shipping — consistent with SAGA's own prediction that dark matter accumulates fastest right after
  a burst of real work, precisely because doc updates lag code by construction.
- **Conflict half-life:** not computable (no prior confirmed-conflict corpus to measure against).
- **Supersession graph integrity:** 1 confirmed gap (RSI-ENGINE's 3 unenumerated edges); 0 orphaned
  goldens in the KIKORYU cluster (clean).
- **Attestation freshness:** not applicable — no `human_attestation` verification methods exist in
  practice yet (the manifest format itself is unbuilt, Finding B).

---

## 6. Ranked Findings, Overall

1. **shankpit-460 NORTHSTAR's two false security-posture claims** (§3) — the highest-severity
   finding in this audit. A doc describing a live product as having "no auth check at all" when it
   now fails closed is the kind of staleness that actively misleads rather than just under-informs.
2. **KAREN — two entire golden specs (FIN-098, FIN-099) architected around an agent with zero
   code**, and a resolution date that predates the specs' own creation.
3. **SAGA hasn't audited itself yet** — ironic given this document is the first real attempt, a day
   after DOC-102 was drafted.
4. **The pre-S156-04 authz gap on `players.go`** — a real vulnerability, found and fixed, that
   exists nowhere in the golden corpus, only in a BACKLOG entry. Exactly SAGA's "dark matter is more
   dangerous because it's invisible to every audit" case, caught only because this session happened
   to touch that code directly.
5. **The FABLE naming collision** (`emily-agent/fable.go` vs. HQ-SPEC-AI-103's FABLE) — low harm
   today, real confusion risk as AI-103 moves from draft toward real build.
6. **IDUNA-NORTH's stale `openapi.yaml` file-map row** — small, mechanical, easy fix.
7. **RSI-ENGINE's unenumerated supersession edges** — low severity, but exactly the kind of gap
   SAGA's linter exists to catch automatically once built.
8. **The one genuinely good finding:** HQ-NORN is real and under-claims itself — the corpus's only
   case of a doc being *more* modest than reality rather than less.

---

## 7. Proposed Next Actions (proposals, not directives — per SAGA's own role boundary, §5 of DOC-102)

These are queued observations for a human or NORN to act on, not unilateral edits made by this
audit:

- Correct shankpit-460 NORTHSTAR's two false claims (§3) — cheapest, highest-value fix available;
  purely factual, no design decision involved.
- Correct IDUNA-NORTH's `openapi.yaml` file-map row and add a pointer to the real 44-route spec.
- Either build KAREN's minimal skeleton (S142-01/S144... whichever lands first) or explicitly mark
  FIN-098/FIN-099 as blocked-on-KAREN in their own frontmatter, so the gap is declared rather than
  silently inherited by every doc that cites her.
- Add a `supersedes:`-equivalent enumeration to RSI-ENGINE naming its three orphaned predecessor
  docs, or explicitly retire them.
- Register WOTAN by name in `NORTHSTAR_KIKORYU.md` and `VS_REALITY_AUDIT.md` with a pointer to
  `shankpit-460/docs2/NORTHSTAR.md`, closing the naming gap at the top of the doc tree.
- Treat this document itself as the first real proof-of-concept that SAGA's manual mode works — the
  actual SAGA build sequence (DOC-102 §9) remains the correct long-term fix so this doesn't have to
  be a full-context multi-agent research pass every time.

---

*Skuld proposes, Verdandi grades, Urd remembers — and SAGA keeps the shelf in order. This document
is Urd's shelf entry for 2026-07-18: not a ruling, a snapshot.*
