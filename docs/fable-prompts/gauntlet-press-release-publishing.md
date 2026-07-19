# Queued Fable prompt — "The Gauntlet": press-release publishing service (planning doc)

**Status:** queued, not yet dispatched. Written 2026-07-19, assembled from several founder
messages across one session — see "What 'the gauntlet' means" below for the exact synthesis.
**Dependency:** none — ready to dispatch now.
**Deliverable type:** a written plan/design doc, not code. This is a new product direction
question, not an implementation task yet.

**To dispatch:** `Agent({ model: "fable", subagent_type: "general-purpose", prompt: <body below> })`

---

## What "the gauntlet" means — pieced together from the founder's own messages, 2026-07-19

Original ask was cryptic ("plan the formal gauntlet content engine wordpress as a backend for
editorial tools") and a grep for "gauntlet" found nothing in the codebase — it's a new concept,
not an existing one. Founder clarified across three follow-up messages, in order:

1. "gauntlet is the metadata enrichment queue and editorial pipeline/calendar"
2. "it will enable us to pivot into providing press release publishing services to satisfy capital
   market requirements"
3. "gated queues content types editorial review (agent and human depending on content type) etc"

Synthesis: **the gauntlet is the editorial pipeline for a new product — EINHORN_INDUSTRIAL
publishing press releases on behalf of companies, not just discovering/aggregating other
companies' releases the way FatBaby does today.** Companies have real regulatory reasons to need
"widely disseminated" press release distribution (Reg FD and similar disclosure requirements are
why services like PR Newswire/Business Wire/GlobeNewswire exist as a category) — this is a
plausible, real pivot, not a toy feature. The gauntlet itself is the multi-stage gated queue a
submitted release passes through: metadata enrichment, an editorial calendar, and content-type-
dependent review — some content types reviewable by an agent alone, others requiring a human,
per the founder's own words.

## Required reading first

1. `PRRJECT_FATBABY/docs/northstar/newssite.md` — **read line ~335 specifically**: "WordPress / any
   CMS — explicitly excluded. This is a generated publication over an event store, not an editorial
   tool." This is newssite's own explicit reasoning for *why it isn't* an editorial CMS — it reads
   and republishes discovered press releases, it doesn't originate them. **The gauntlet is not
   newssite** — don't let this dispatch accidentally propose rebuilding newssite as a CMS; that
   line explains why that would be wrong for newssite specifically. The gauntlet is a new, separate
   product: outbound publishing, not inbound aggregation. Read the rest of `newssite.md` too, for
   what FatBaby's wire-service reading already looks like — the new product's *output* should
   probably be indistinguishable from a normal FatBaby-discovered release once published, so the
   two systems need to interoperate even though they're separately purposed.
2. `EMILY/docs/hq-specs/HQ-SPEC-PRIME-101-norn-loop-kernel.md` — §4's gate tiers
   (`autonomous`/`prime_ack`/`biometric`) are close to exactly what "agent and human depending on
   content type" describes. Don't invent a bespoke review-tier system from scratch; evaluate
   whether NORN's existing tiers instantiate directly onto content types (e.g. routine
   earnings-release boilerplate = autonomous, anything with forward-looking financial claims or
   M&A language = prime_ack/human), same instantiate-don't-reinvent framing used elsewhere in this
   backlog tonight (see the SHANKPIT-460 maps dispatch, entry #11, for the same pattern applied to
   generated-map fuzz-eval).
3. `EDIS/NORTHSTAR.md`, `EDIS/CLAUDE.md`, `EDIS/docs/architecture.md` — EDIS is the actual
   WordPress product surface in this company (edis-core, edis-signals, edis-ask-emily plugins).
   If WordPress belongs anywhere for this, it's as a new EDIS plugin (`edis-gauntlet`?), reusing
   `edis-core`'s existing API-client/shortcode pattern, not a standalone WordPress install. Confirm
   or reject that placement explicitly rather than leaving it implicit.
4. `PRRJECT_FATBABY/docs/GTM_FUNNEL.md` — the existing Ask Emily revenue funnel design, for
   consistency of tone/structure with how this company already scopes a product-facing funnel.
5. Skim `EMILY/BACKLOG.md` SECTION 160 (added same session) — the tickerization/regex-extraction
   work already underway on the *inbound* side (FatBaby discovering other companies' releases).
   The gauntlet's metadata-enrichment stage should reuse `internal/prwatch/tickers.go`'s
   `ExtractTickers`/`ExtractFromHTML` rather than reimplementing ticker extraction a second time —
   name this explicitly as a reuse point.

## What the report needs to cover

- **What "publishing" actually requires to satisfy capital-markets disclosure norms** — real wire
  services provide timestamped, widely-disseminated, non-editable-after-publish distribution.
  What's the minimum credible version of that this company could stand behind? Don't hand-wave
  the compliance angle; name what's actually required vs. what's marketing.
- **The gated queue, concretely**: stages (submission → metadata enrichment → editorial
  calendar/scheduling → review gate → publish), what "metadata enrichment" actually enriches
  (ticker/company identity via the existing regex extractor, plus whatever else a real release
  needs — dateline, category, distribution list), and where NORN's gate tiers map onto content
  types.
- **WordPress placement** — argue it, using EDIS as the default per required reading #3, or argue
  against it explicitly if the editorial workflow genuinely doesn't fit WordPress's model.
- **Interop with existing FatBaby ingestion** — once EINHORN_INDUSTRIAL is a publisher, does its
  own FatBaby pipeline discover its own published releases the same way it discovers everyone
  else's? Should it? Name the loop explicitly, don't leave it implicit.
- **Sequencing** — smallest credible v0 vs. later bets, same discipline as every other planning
  doc tonight (no wishlist, a real recommendation).

## Where this lands

Write to `PRRJECT_FATBABY/docs/northstar/gauntlet-press-release-publishing.md` (new file). Register
in `EMILY/context/golden-docs-index.md` if you judge it load-bearing enough — your call, one-line
justification either way. Add concrete `[ ]` backlog items to `EMILY/BACKLOG.md` under a new
section for whatever this plan concludes should happen first. No code in this dispatch.

---

## Prompt body ends here
