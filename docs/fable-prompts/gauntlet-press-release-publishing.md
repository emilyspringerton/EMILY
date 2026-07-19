# Queued Fable prompt — "The Gauntlet": northstar for FatBaby content management, licensing + publishing

**Status:** queued, not yet dispatched. Written 2026-07-19, assembled from five founder messages
across one session — see "What 'the gauntlet' means" below for the exact synthesis. Scope
broadened twice during the same session (press-release publishing → full northstar umbrella) —
this file reflects the final, broadest framing; don't dispatch against an earlier mental model.
**Dependency:** none — ready to dispatch now.
**Deliverable type:** **a NORTHSTAR.md, not a one-off report.** Founder's own words: "as a north
star... continuing evolving concept and process." Write it the way `SHANKPIT/docs2/NORTHSTAR.md`
is written — "not a roadmap with dates... a statement of direction that makes priorities obvious
even when the shape of the work isn't yet fully visible," meant to be kept alive alongside the
codebase, not a single point-in-time report. No code in this dispatch.

**To dispatch:** `Agent({ model: "fable", subagent_type: "general-purpose", prompt: <body below> })`

---

## What "the gauntlet" means — pieced together from the founder's own messages, 2026-07-19

Original ask was cryptic ("plan the formal gauntlet content engine wordpress as a backend for
editorial tools") and a grep for "gauntlet" found nothing in the codebase — it's a new concept,
not an existing one. Founder clarified across five follow-up messages, in order:

1. "gauntlet is the metadata enrichment queue and editorial pipeline/calendar"
2. "it will enable us to pivot into providing press release publishing services to satisfy capital
   market requirements"
3. "gated queues content types editorial review (agent and human depending on content type) etc"
4. "as a north star we are a hybrid ai and also we plan to be a traditional human powered ai
   augmented news room too — we need the platform to manage all the content and the backend
   systems and queues etc to help pull the actual content out of our data as a north star and
   continuing evolving concept and process"
5. "gauntlet may be several systems but is the general concept of the system that manages fatbaby
   financial content for licensing and publishing to the public internet"

**Final synthesis (this is the scope to plan against, supersedes the narrower press-release-only
framing from messages 1-3):** the gauntlet is not one system — it's the umbrella concept for
**everything FatBaby's raw data becomes on its way out the door**, across three distinct exit
paths that all share the same underlying pipeline (extraction → enrichment → gated review →
release):

- **Publishing** — the press-release-publishing-as-a-service angle (message 2), *and* newssite's
  existing aggregated wire-service output — both are "publishing to the public internet," just two
  different content types moving through the same conceptual gauntlet.
- **Licensing** — selling/licensing FatBaby's structured financial content/data to third parties.
  Not previously scoped anywhere in this backlog — new.
- **The newsroom itself** — "a hybrid AI and also... traditional human powered AI augmented news
  room" (message 4). This is an operating-model statement, not just a technical one: some content
  is produced/reviewed by agents alone, some needs a human in a genuinely editorial (not just
  compliance-gate) role — this is broader than message 3's "agent vs human review by content type"
  framing, which was about *approval gates*; message 4 is about who's actually doing editorial
  *work*, which may be a different axis entirely (worth the report distinguishing the two).

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
6. `PRRJECT_FATBABY/docs/headlines/tina.md` (golden-docs-index entry `FATBABY-TINA`) — founder
   raised TINA (Trading Idea, Not Advice engine) as related to the gauntlet, then **immediately
   self-corrected the direction**: first said "TINA would be a consumer of gauntlet," then "or a
   producer for gauntlet rather." **This directionality is genuinely unresolved, not settled —
   don't silently pick one.** Read TINA's own design doc and figure out which is actually true (or
   whether it's both, at different pipeline stages), and say so explicitly with reasoning, rather
   than asserting a direction the founder themselves wasn't sure of. This is real evidence the
   gauntlet's exact position in the architecture is still being worked out live — treat "name every
   known relationship, flag every unresolved one" as a requirement for this doc, not just TINA.

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

Write to `PRRJECT_FATBABY/docs/northstar/gauntlet.md` (new file — not the narrower
`-press-release-publishing` name, given the scope broadened to the full umbrella concept). This
should be golden-doc registered — given "continuing evolving concept and process" is the founder's
own framing, this is exactly the kind of living document `EMILY/context/golden-docs-index.md`
exists for. Append it there (`GAUNTLET | PRRJECT_FATBABY/docs/northstar/gauntlet.md | 1 | <budget> |
<one-line>`), commit + push EMILY per the standard golden-doc registration protocol. Add concrete
`[ ]` backlog items to `EMILY/BACKLOG.md` under a new section for whatever this plan concludes
should happen first. No code in this dispatch.

---

## Prompt body ends here
