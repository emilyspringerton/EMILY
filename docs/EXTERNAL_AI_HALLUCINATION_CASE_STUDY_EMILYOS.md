# Case study: an external AI hallucinating about "EmilyOS" (2026-08-25)

**Status: research note, not a design doc.** Nothing in this file describes anything real about
this monorepo's own systems. It documents an external model's fabricated output and what
actually explains it, for future reference if something similar happens again.

## What happened

Founder found a LinkedIn Pulse post by Martin Robson (Founder & CEO, Robson Communications Inc.,
`ca.linkedin.com/in/robsonai`, tagline "🇨🇦 Canadian Data Sovereignty - Emily OS"), dated
2026-03-09, describing his own, entirely separate project also named "Emily OS" — a "sovereign
cognitive architecture" with acronyms `EMEB`/`EARL`/`ECCR`/`ECGL`/`AOPO`/`APC`. This repo's own
`EmilyOS` (first commit 2026-01-31) predates that post by ~5 weeks and shares nothing with it —
confirmed by grepping this entire monorepo for every one of those acronyms and "Martin Robson":
zero matches anywhere.

Separately, the founder asked "Google Chat AI" (model uncertain — possibly Gemini, possibly
Grok) something that returned a multi-part, confident description of an "EINHORN_INDUSTRIAL"
programming paradigm: a fictional language called "Smilium" with `@shell`/`@game` dual execution
modes, "Sovereign Code Management," "Deterministic Object Routing," and an AI that "makes
optimizations on its own community servers."

## What's confirmed real vs. fabricated

| Claim | Verdict | Evidence |
|---|---|---|
| Repo names `SHANKPIT`, `REDGARDEN`, `EmilyOS`, `PARENA`, `EMILY`, `emily.cli`, `GoblinFoxDragon` exist | **Real** | All return HTTP 200 from `api.github.com/repos/emilyspringerton/<name>` — genuinely public repos |
| `EmilyOS`'s real public content is a "Game-UI Filesystem" GUI design spec (tiling panes, double-click semantics, keyboard nav) | **Real** | Fetched directly from the live GitHub repo page |
| A language called "Smilium" | **Fabricated** | Grepped this entire monorepo, zero matches, no such thing exists |
| "Sovereign," "data sovereignty," "sovereign cognitive architecture" framing | **Borrowed from Martin Robson's unrelated post**, not this codebase | His own LinkedIn tagline uses this exact language; not present anywhere in this monorepo before this incident |
| Erlang/BEAM fault-tolerance philosophy mentioned in the hallucinated text | **Real, but repurposed** | `GoblinFoxDragon/docs2/MOD_SURFACE_NORTHSTAR.md` §3a genuinely discusses an Erlang/BEAM idea — added and pushed to that public repo earlier this same session |
| "Object-capability routing," "hijacking the AI's memory spaces," "own community servers" | **Fabricated, generic filler** | Generic secure-computing buzzwords, not distinctive to this codebase, not found anywhere in it |

## Working theory

A search-grounded model was asked about "EmilyOS"/"Emily"/"EINHORN_INDUSTRIAL," found two real,
unrelated things sharing the name "Emily OS" (this monorepo's public GitHub org, and Martin
Robson's separate LinkedIn-published project), and produced a fabricated synthesis: real
fragments (repo names, the real Erlang idea, the real GUI-filesystem concept) stitched together
with borrowed language from the unrelated post and generic AI-architecture filler invented to
bridge the gaps. This is an ordinary LLM hallucination failure mode — confident confabulation
around partial real grounding — not evidence of any actual system behavior. Nothing in this
monorepo has any mechanism to publish externally, self-modify in the way described, or run
outside this machine.

## Why this is worth keeping

Real, useful signal for the future: this monorepo's public repos (`EmilyOS`, `EMILY`, `PARENA`,
`emily.cli`, `SHANKPIT`, `REDGARDEN`, `GoblinFoxDragon`, ...) are genuinely crawlable and get
picked up fast — the Erlang §3a addition was live and apparently already reachable within the
same session it was pushed. Worth remembering next time an external AI's output about this
project seems oddly specific: check whether the specific-sounding parts trace to something
actually public, and treat anything that doesn't as invented, not as investigated fact.

## Not done here

No attempt was made to reproduce the original Google Chat AI/Grok query — the exact prompt that
produced this output isn't known, and neither service was queried directly as part of this
research pass.
