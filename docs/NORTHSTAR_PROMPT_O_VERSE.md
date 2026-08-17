# NORTHSTAR — Prompt-o-verse

**Status:** Draft v0.1 — concept + northstar only, no implementation, no repo yet
**Date:** 2026-08-17
**Founder framing, verbatim (assembled across several real-time fragments):** "what if our mission
was to organize the world's information like google but instead of going the way google is doing
it in terms of identifying the information of the world we start with a prompt like a baseball
card photography and generate a bunch of different versions... as we generate versions we identify
the categories of types of baseball card photography like portraits vs action shots... the styles
will probably have different prompts for the different generations like super early baseball cards
aren't even photos, are they drawings? so we identify the different prompts that create the
different categories or taxonomies... we are going to have at least 2 levels — the ez prompt like
'give me a 99 rookie card of Mark McGwire' or '80s mall portrait' — and then we use those gens to
identify the features of what makes that generation unique, and store an expanded prompt or have a
model dedicated to giving the more detailed feature-rich prompt (lights, background, style,
fashion, etc)." Mid-turn additions: "it should be a playground for those learning gen AI — like
people that don't know you can ask for an 80s mall photo can see that in the gallery with a
rating." / "another meta [category] has all these weird [prompts] like 'make it like the photo but
also make the people look like ice cream novelties' — the goal is to have a universe of prompt
ideas for gen AI to help the creativity of people who are not used to thinking in terms of any
prompt is possible."

---

## 1. The reframe, stated plainly

**The mission, in the founder's own final and simplest phrasing: "categorize all information."**
Google's mission is to organize information that already exists — crawl it, index it, rank it.
This is the inverse move: categorize a *taxonomy of styles* by **generating** the space first, then
discovering its structure from what comes out, rather than discovering structure in something
already sitting there to be crawled. The "world's information" being organized here isn't web
pages — it's the latent space of a generative model, made legible by systematically exploring it
and naming what's found.

The pilot domain is baseball card photography, chosen because it has real, discoverable structure
a person already half-knows intuitively (there's a difference between a hand-drawn 1910s tobacco
card and a glossy 1990s photo card, and between a posed portrait and an action shot) but that
nobody has written down as an explicit, generateable taxonomy. That's the actual product of the
first pass: not "a baseball card generator," but a **map of what baseball card photography even
is**, expressed as prompts that reliably reproduce each node on the map.

**The sharpest framing of what's actually novel here, the founder's own: "it's like a reverse
labeled dataset — a dataset with only labels — and then we can tag output examples of prompts."**
Every conventional ML dataset starts with raw data and adds labels on top (a photo exists, someone
annotates it "1980s mall portrait"). This inverts that: the label/taxonomy — the prompt tree
itself — is the primary, durable asset, built *before* any specific image does. Generated images
are cheap, regeneratable **exemplars tagged back to their originating label**, not the dataset
itself. The taxonomy would still be worth having with zero images attached to it; the images are
proof-of-concept instances, not the asset being categorized. This is the actual mechanism behind
"categorize all information" (§ above) — the categorizing happens at the label layer, generation
is just how a label gets made checkable/visible/browsable.

The mission is broader than mapping realistic/historical styles, though — stated explicitly by the
founder, not inferred: **the end goal is a whole universe of prompt ideas that expands what people
think is even possible to ask a generative model for.** That includes a second, deliberately
different kind of taxonomy branch alongside the historical/stylistic one (§2): a "weird" or
surreal-transformation category — the founder's own example, "make it like the photo but also make
the people look like ice cream novelties." These aren't trying to reproduce something that existed
in the real world (no such thing as an "ice cream novelty portrait" era); they're creative-range
demonstrations, existing specifically to show someone who's never prompted a model before that the
space is bigger than "make it look like X photography style." Both branches — real/historical and
surreal/whimsical — are Prompt-o-verse; neither is the whole thing.

## 2. The taxonomy discovery loop

The core mechanic, restated as a repeatable process rather than a one-off exploration:

1. **Generate broad.** Start from a wide, unconstrained prompt for the domain ("baseball card
   photography") and generate many outputs.
2. **Observe categories.** Look at what actually came out and name the axes of variation a human
   (or a vision-capable model doing a first pass) can see — era (hand-drawn vs. photographic vs.
   modern digital), shot type (posed portrait vs. action shot vs. team photo), production values
   (glossy studio vs. cheap newsprint vs. premium foil), and whatever else the generations
   themselves surface that wasn't anticipated going in. The taxonomy is discovered, not
   pre-specified — the whole point is that nobody sat down and enumerated "here are the 40 styles
   of baseball card photography" in advance.
3. **Isolate the prompt per category.** For each named category, find the narrower prompt (or
   prompt *modification* relative to the broad one) that reliably reproduces it. "1910s tobacco
   card" and "1990s glossy rookie card" are not the same prompt with different years swapped in —
   they likely need genuinely different vocabulary (illustration/lithograph terms for the former,
   camera/lighting/film-stock terms for the latter).

   **The concrete three-stage mechanism, stated directly:** label (the top-level/category prompt)
   → generate (many images from that label) → **label the generated data itself** (tag each
   individual output with its own specific feature values — this makeup, this hair, this
   background). Step 3's "isolate the prompt per category" and §3's per-category feature-value
   space aren't separate ideas — the second labeling pass, applied across many generations under
   one top-level label, is literally how the feature-axis space for that category gets discovered:
   the aggregate of individual-output tags *is* the feature space.
4. **Build the tree.** Top-level prompts (broad category, e.g. "era: early 20th century") branch
   into sub-prompts (e.g. "tobacco card portrait" vs. "tobacco card action pose"), which can branch
   further. The taxonomy *is* a prompt tree — walking down it is choosing a more specific prompt.
5. **Repeat, and repeat across domains.** Baseball cards prove the loop; the same process should
   generalize to any visual genre with real historical/stylistic structure (album covers, yearbook
   photos, movie posters, mugshots, driver's licenses — anything with recognizable eras and
   sub-styles a generative model can approximate).

This loop as described discovers the **realistic/historical** branch of the taxonomy — styles that
map to something that actually existed (a real card era, a real photo-studio convention). The
**surreal/whimsical** branch (§1, "ice cream novelty portraits") runs a related but distinct
process: instead of discovering categories that already exist somewhere to be named, it's about
generating and curating categories that are *interesting because they don't* — the taxonomy there
is closer to "a curated gallery of good weird ideas" than "a map of a real thing." Both belong in
the same tree; a taxonomy node just needs a `kind: historical | surreal` tag (or similar) rather
than needing two separate systems.

**A genuine architectural fit worth naming, not forcing:** this loop — propose a candidate
(generation), grade it against a standard (does it actually match the claimed category), gate
whether it's good enough to promote to "the canonical prompt for this taxonomy node" — is
structurally the same propose→grade→gate→promote loop `NORN` (`NORN/pkg/norn`, HQ-SPEC-PRIME-101)
already exists to generalize. Whether Prompt-o-verse should actually be built as a NORN
instantiation (Proposer = prompt-variation generator, Oracle = category-match grader, Gate =
promotion into the canonical taxonomy) is a real design option for VS1+, not a decision made here —
flagged because the fit is unusually clean, not because it's assumed.

## 3. Two-tier prompting

Two levels, matching what a casual user types versus what actually drives high-fidelity, on-style
generation:

- **Tier 1 — the EZ prompt.** What a person actually types: "give me a 99 rookie card of Mark
  McGwire," "make my photo look like an 80s mall portrait." Short, casual, the way anyone would ask
  for something without knowing any prompt-engineering vocabulary.
- **Tier 2 — the expanded prompt.** The feature-rich version that actually reproduces the style
  reliably: specific lighting setup, background choice, film grain / color grade, framing, era-
  correct fashion and hair, card border/typography conventions, etc. This is the output of step 3
  in the discovery loop above — the thing that was reverse-engineered from generations, not
  invented from imagination.

  **Refinement, stated directly by the founder:** tier 2 isn't one fixed recipe per category — "80s
  mall portrait" alone could have 100+ real variations, each a different combination of specific
  feature values (makeup, hair, fashion, background, lighting, ...). So the expanded-prompt layer
  is better modeled as a **space of feature axes**, each with a discovered range of real values for
  that category, than as a single canonical prompt string. Generating "an 80s mall portrait" means
  sampling a combination from that space, not replaying one fixed template — the taxonomy node
  owns a structured feature space, and any point in that space is a valid instance of the category.

  **Data-model refinement, also stated directly:** don't store the expanded description as one long
  free-text prompt string — normalize it into structured, queryable tags/fields instead (camera
  setup, lighting, background, etc. as their own attributes, not prose). Two concrete reasons, both
  named by the founder: **dedup** — a lot of these long shot-description components repeat across
  otherwise-unrelated categories (the same "diffused three-point studio lighting" phrase doesn't
  need to be re-written every time a category happens to use it), so normalizing once and
  referencing it beats storing the same paragraph N times; and **queryability** — a long prose
  string isn't data you can actually query ("show me every category using hard flash lighting"
  is unanswerable against a blob of text, trivial against a tag). The genuinely valuable emergent
  property this creates: **shared tags surface associations between top-level prompts that don't
  obviously share context.** Two categories with nothing in common by name or by parent node in
  the taxonomy tree (§2.4) can turn out to share a lighting tag, a background tag, an era-adjacent
  fashion tag — connections a human skimming category names would never notice, but that fall out
  for free once the feature values are tagged data instead of prose. This means the taxonomy isn't
  purely a tree (parent → child categories) — it's also a **graph**, with edges formed by shared
  tags cutting across the tree's branches. Refined, not just binary present/absent: some
  third-level tags naturally align more strongly with a given top-level prompt than others, so
  edge *strength* (how much depiction two categories genuinely share, not just whether they share
  any tag at all) is what separates a meaningful "these are related via common depiction" signal
  from incidental overlap. Both structures coexist: the tree is how a person browses (§4); the
  graph is what a query can surface. Which storage shape (a real relational schema, a tag/attribute
  store, something graph-native) is left as an implementation decision, not picked here.

  **A real tension the founder raised and then answered, worth keeping both halves of:** if
  everything — including a top-level/EZ prompt itself — can be fully expressed as a tag
  combination, is "top-level prompt" a genuinely distinct concept, or just a common/privileged
  point in the same tag-space as everything else? Self-answered rather than left hanging: top-level
  prompts and explicit tags earn their keep as a **deliberate steering mechanism**, not because
  they're ontologically special. Pure random generation (§2.1, sampling broadly under temperature)
  won't reliably hit every real combination on its own — some genuine categories are rare enough
  in the model's own output distribution that chance alone leaves them undiscovered. An explicit
  prompt or tag is how the taxonomy gets **deliberately steered** toward a specific point in the
  space instead of waiting for temperature to wander there. This also means taxonomy coverage has
  a real, nameable gap class going forward: regions of tag-space nobody has explicitly steered
  toward yet, distinct from regions genuinely absent from the model's learned distribution — worth
  distinguishing in practice, not assumed to be the same kind of "missing."

  **Landed on, not just floated:** "any top-level prompt can probably be expressed via a third-
  level feature" — the tension resolves by unifying, not by keeping "top-level prompt" and "tag"
  as two kinds of object. There's one feature space; a top-level/EZ prompt is a named, commonly-
  useful combination within it, not a structurally different thing from a tag. Gap detection
  becomes concrete and operational under this model: examine the tagged output collected for a
  category (the "third-level labeled gen output") and look for expected feature values that never
  show up — e.g. noticing a category has real generated variety across lighting and background but
  **no cartoon/illustrated depictions at all**, even though nothing rules that out conceptually.
  That absence is exactly the signal for where to deliberately steer (previous paragraph) rather
  than keep sampling broadly and hoping.

**Open question, correctly left open rather than picked here:** whether tier-2 expansion is (a) a
**stored** expanded prompt per known taxonomy node — a lookup, cheap and deterministic, but only
covers styles the taxonomy already knows about — or (b) a **dedicated expansion model/prompt**
whose job is EZ-prompt-in, rich-prompt-out, generalizing to styles not yet in the taxonomy at the
cost of being less predictable. These aren't mutually exclusive: the taxonomy tree could serve as
the training/few-shot data for an expansion model, with stored prompts as the reliable fallback
for anything already mapped. Which to build first, and whether the expansion model is fine-tuned,
prompted, or RAG-style retrieval over the taxonomy, is unresolved and shouldn't be assumed toward
either answer without a real design pass.

## 4. The playground — the affordance, not a bolt-on feature

Stated directly by the founder and worth being exact about: **the gallery/playground is the
affordance** — the actual human-facing interface the categorization mission (§1) is *for*. A
taxonomy nobody can browse is a research artifact; the gallery is what turns "we categorized all
this" into something a person can actually use, learn from, and enjoy. This isn't a late feature
layered on top of the real product — it's the product's face, and §8's phasing (which builds the
taxonomy mechanic before the gallery) is an engineering sequencing choice about what has to be
proven first, not a claim that the gallery matters less.

Concretely: a gallery where every discovered taxonomy node (e.g. "80s mall portrait") is browsable
with its example generations, so someone who
has never heard the phrase "prompt engineering" and doesn't know you can even *ask* for an "80s
mall photo" style can discover it exists by scrolling — the taxonomy becomes the teaching surface,
not documentation nobody reads. A rating mechanic sits on top: users rate how well a given
generation matches its claimed category, which does double duty as (a) a way for browsers to find
the best examples of a style, and (b) real signal feeding back into the discovery loop's grading
step (§2.3) — a generation the community rates as "not actually an 80s mall portrait" is exactly
the kind of feedback that should demote a bad candidate prompt rather than let it calcify into the
canonical one.

Real open questions this raises, not resolved here: how ratings resist gaming/brigading at scale,
whether rating requires an account (ties to IDUNA identity, or fully anonymous), and whether the
gallery is free-to-browse with generation itself gated (matching the existing Emily+ free-tier
pattern used elsewhere in this system) or something else entirely.

## 5. What's deliberately not decided here

Named plainly rather than glossed over, matching this document's job as a northstar and not a
build plan:

- **Image-generation backend — RESOLVED 2026-08-17, no longer open.** No existing image-generation
  infrastructure existed anywhere in this monorepo before today (checked, not assumed). Google's
  "Nano Banana" (Gemini image generation) is the real answer, but getting there took two dead ends
  first, both worth keeping as a record: (1) this box's existing OAuth-personal Gemini CLI
  credential (`~/.gemini/oauth_creds.json`) — tested directly, confirmed dead,
  `IneligibleTierError` ("no longer supported for Gemini Code Assist for individuals... migrate to
  Antigravity"), a real Google-side deprecation, not a config problem here; (2) a dedicated
  `GEMINI_API_KEY` via Google AI Studio — founder spent ~2 hours on this already, hit a 404 loop on
  one of Google's own recommended key-creation paths; live web research confirmed this is a real,
  currently-active platform-side bug (Google's own AI Developer Forum has open threads since May
  2026 for "Failed to generate API key, the request is suspicious," plus separate unresolved
  403/404 permission-denied threads), not founder error.

  **What actually worked, founder's own idea:** bypass the console entirely and use this account's
  existing `gcloud` ADC credentials against **Vertex AI** instead of the consumer AI Studio API —
  a completely different auth path (standard GCP IAM, not the buggy API-key UI), and this account
  already has `aiplatform.googleapis.com` enabled with an active billing account. Tested for real:
  `POST https://us-central1-aiplatform.googleapis.com/v1/projects/{PROJECT}/locations/us-central1/
  publishers/google/models/gemini-2.5-flash-image:generateContent`, `Authorization: Bearer $(gcloud
  auth print-access-token)`. Real success — genuinely convincing 1024×1024 output on the first
  real Prompt-o-verse prompt ("a 1990s glossy baseball rookie card portrait, studio lighting, blue
  background"), saved as the project's first real artifact:
  `EMILY/docs/prompt-o-verse/vs0-proof/1990s-rookie-card-portrait.png` (README alongside it has
  the exact call). This is real evidence VS0's core mechanic works, not just that a backend exists
  on paper — the backend question converts from "unpicked" to "picked and proven": Vertex AI via
  existing `gcloud` ADC, no new credential needed.
- **Repo home.** No dedicated repo exists yet. This doc lives in `EMILY/docs/` for now, same
  precedent as `NORTHSTAR_OPENCLAW_INTEGRATION.md` before OpenClaw had a repo — move it once (or
  if) a real repo is created and scoped, matching the `CarePyre`/`EXODUS`/`TTT` pattern of
  cloning a founder-created GitHub repo rather than one spun up unprompted.
- **Domains beyond baseball cards.** Baseball card photography is the pilot precisely because it's
  concrete enough to prove the loop against; which domain comes second (and whether the taxonomy-
  discovery *process* itself, not just its baseball-card output, is the reusable product) is a
  question for after VS0 actually produces a real taxonomy, not before.
- **Moderation.** A public gallery of AI-generated imagery with user ratings needs a real content
  policy before it's public — not addressed here, flagged so it isn't skipped later.

## 6. Beyond the playground — other uses for a reverse-labeled dataset

The gallery (§4) is the primary product, but it's worth naming — the founder's own prompt,
interpolating the mission (§1, categorize all information) and the product (§4, the playground) —
what else a **multi-layered** taxonomy (top-level category → sub-category → feature-value space →
tagged exemplars, §2/§3) is actually good for once it exists, since the durable asset is the label
structure, not any specific image (§1's "reverse labeled dataset" framing). None of these are
scoped or committed — brainstormed, flagged as real option value the taxonomy creates almost for
free once it exists:

- **Bootstrapping tier 2 itself.** The EZ-prompt → expanded-prompt pairs collected while building
  the taxonomy (§3) are exactly the parallel corpus needed to train or few-shot the expansion
  model §3 leaves as an open question — the dataset teaches the tool that later serves it.
- **A frozen eval benchmark for image-gen backends.** A rated, tagged corpus per taxonomy node is
  a real ground truth for "does this backend actually reproduce '80s mall portrait' when asked" —
  reusable as a `NORN` Oracle (§2's flagged fit) for grading any future backend or prompt-strategy
  change against a frozen standard, the same oracle-freezing discipline `PRIME-101` already uses
  elsewhere in this system.
- **Reverse style lookup.** Upload a real photo, match it against the taxonomy, return "this reads
  as `90s yearbook photo`" plus the prompt that reproduces or riffs on it — the mission's mirror
  image: recognizing a category in existing data instead of generating an instance of one.
- **Cross-repo creative-asset seeding.** A structured, named style space is directly useful to
  anything in this monorepo that needs stylistically-coherent generated content — merch design
  (STINKIES' sticker/hoodie drops), procedural cosmetic variety (`GOLDENBAND`/`REDGARDEN`/
  `SHANKPIT`), no design work spent from scratch each time.
- **A citable, licensable dataset in its own right.** Same instinct already active elsewhere in
  this system (SKULDMARK's planned open release to researchers/industry, S175-03 in
  `EMILY/BACKLOG.md`) — a well-documented taxonomy of generative style-space is a real research/
  industry artifact independent of any product built on top of it.
- **Prompt-engineering curriculum.** The taxonomy *is* a documented "here's exactly what makes 40
  styles work" teaching resource — a natural fit for this codebase's existing teaching-through-
  real-examples pattern (`TYLER`, the `TTT` "Tyler Teaches Typing" repo).

## 7. GPT-2-assisted extrapolation — prompts for "realities that don't yet exist"

Stated directly by the founder: run GPT-2 — the existing `gpt2-alpine-c` fine-tuning + inference
infrastructure already built for Emily Prime — over different levels and segments of the taxonomy
data (top-level labels only; one branch, e.g. all historical eras; the full feature-tag corpus per
node) to help identify and generate prompts for **"realities that don't yet exist as far as we are
aware."** This is a distinct move from §2's discovery loop, worth being precise about the
difference: §2 discovers categories that are *already latent* in what a broad prompt generates —
observing what's there. This is extrapolation — training a model on the taxonomy *as it exists so
far* and using its own generative capacity to propose plausible categories nobody has generated an
example of yet, the same way a language model trained on real sentences can produce a grammatical
sentence nobody has typed before. The taxonomy stops being purely observational and starts being
generative in its own right — a second-order use of the reverse-labeled dataset (§1): first the
labels generate images, then the labels themselves become training data that generates *more
labels*.

Concretely, this reuses `gpt2-alpine-c`'s established pipeline shape almost directly rather than
requiring new infrastructure: build a JSONL corpus from taxonomy segments (mirroring
`scripts/prime_directive_dataset.py`'s existing golden-doc-to-corpus pattern), fine-tune via the
same Colab notebook flow, serve via the same inference API (`scripts/serve.py`, proxied through
`emily-agent`'s `:8086/api/v1/gpt2/generate`). **Honest caveat, not glossed over:** this pipeline
exists in code but a first real Colab fine-tune run hasn't actually completed successfully yet
anywhere in this system as of this writing (`EMILY/BACKLOG.md` S26-04, still open/HITL-blocked) —
Prompt-o-verse would be a new *consumer* of that pipeline, not something that gets it working for
the first time; the two are coupled but this document doesn't assume S26-04 resolves itself.

This also closes a loop with §2's flagged `NORN` fit rather than competing with it: GPT-2 fits
naturally as the **Proposer** in a propose→grade→gate→promote instantiation — it proposes candidate
new taxonomy nodes/prompts, and the existing discovery mechanism (generate an image from the
candidate, check whether it actually holds together as a real category, i.e. the Oracle/Gate) is
what decides whether an extrapolated "reality that doesn't exist yet" earns a place in the real
taxonomy or gets discarded as a plausible-sounding dead end. Not decided here whether this is
built as a formal `NORN` instantiation or something lighter — flagged as the same clean fit §2
already named, now with a concrete Proposer candidate.

## 8. Phased plan

**VS0 — prove the discovery loop, no product surface yet.** Pick one image-generation backend
(cheapest path to a working proof, not necessarily the eventual production choice), run the
generate→observe→categorize→isolate-prompt loop by hand against baseball card photography, and
produce a real, checked-in taxonomy tree (even a shallow one — 3-4 top-level eras, 2-3 sub-styles
each) with a verified prompt per node: does the stored prompt for "1990s glossy rookie card"
actually and repeatably produce that style. This is where "does the core mechanic even work" gets
answered before anything else is built on top of it.

**VS1 — the two-tier prompt system.** Build the EZ-prompt path against VS0's real taxonomy (start
with the stored-lookup approach from §3, since it doesn't require training/prompting a new
expansion model first) and decide, with real usage data from VS0's taxonomy, whether an expansion
model is worth building next.

**VS2 — the public playground.** Gallery, browsing by taxonomy node, rating mechanic, and the
open questions from §4 (accounts, abuse-resistance, access model) get real answers instead of
being flagged. Only after VS0 proves the mechanic and VS1 proves the two-tier prompt system is
usable — the gallery is the product's face, but it's built on a taxonomy that has to be real first.
