# SAGA Frontmatter Schema + Claim-ID Convention

`status: draft — build sequence step 1 of HQ-SPEC-DOC-102`
`role: defines the machine-parseable frontmatter every doc in EMILY/docs/hq-specs/ carries,
and the claim-ID convention SAGA's tooling (emily saga lint) checks against`

*Per HQ-SPEC-DOC-102 §9 Build Sequence step 1: "Frontmatter schema + claim-ID convention +
linter (saga lint), retrofit onto the live spec series. This document ships with its own
claims tagged — DOC-102 eats first." This document is the schema itself; DOC-102's own
retrofit (and the other six HQ-SPEC docs') is the worked example.*

## 1. Where it lives

A YAML-ish frontmatter block (`---` ... `---`) at the very top of the file, before the `#
Title` line. Deliberately a *restricted subset* of YAML — hand-parsed by `emily saga lint`
(no new third-party dependency; this codebase is stdlib-first throughout), not a full YAML
parser. The restriction: scalar `key: value` lines, `key: []` for an explicitly-empty list,
and one nested list shape (`claims:`) with a fixed set of per-item keys. No multi-line
strings, no anchors, no arbitrary nesting.

## 2. Fields

```yaml
---
doc_id: DOC-102
authority: draft
supersedes: []
amends: []
claims:
  - id: DOC-102.POL-1
    type: POL
    reality_binding: specified
  - id: DOC-102.INV-1
    type: INV
    reality_binding: specified
---
```

- **`doc_id`** (required) — the spec's own ID, e.g. `DOC-102`, `FIN-098`. Must match the
  `HQ-SPEC-<doc_id>-...md` filename.
- **`authority`** (required) — Axis 1 from HQ-SPEC-DOC-102 §3: one of
  `draft | golden | amended | superseded`.
- **`supersedes`** (optional, default `[]`) — list of `doc_id`s this document fully
  supersedes. Every entry must resolve to a real `doc_id` elsewhere in the corpus (a
  dangling reference is a lint error — DOC-102 §3: "unenumerated claims are a lint error,
  not a judgment call").
- **`amends`** (optional, default `[]`) — list of `{doc: <doc_id>, claims: [<claim-id>, ...]}`
  for *partial* supersession. The `claims` list must be non-empty — an `amends` entry with
  no enumerated claims is exactly the "unenumerated inheritance" case §3 calls out as a lint
  error, not a judgment call.
- **`claims`** (required, may be empty for a doc with no claims yet — rare) — list of:
  - **`id`** — `<doc_id>.<TYPE>-<N>` (§2's claim-ID format). `N` is sequential per doc,
    starting at 1. The `<doc_id>` portion must match this document's own `doc_id` (a claim
    can only be *declared* in its own document; other documents *cite* it by ID, they never
    redeclare it under a different doc's frontmatter).
  - **`type`** — one of `INV | BEH | IFC | MET | POL | NAR` (DOC-102 §2).
  - **`reality_binding`** — Axis 2 from §3: one of
    `specified | building | running | verified | diverged`. `NAR` claims are exempt from
    reality-binding by design (§2, §10) but still carry the field for schema uniformity —
    lint does not enforce progression on `NAR` claims, only presence.

## 3. What `emily saga lint` checks

Per HQ-SPEC-DOC-102 §5 ("run lint: orphan goldens, unenumerated inheritance, ID collisions")
and §3 ("unenumerated claims are a lint error, not a judgment call"):

1. **Parse validity** — every `.md` file's frontmatter block parses under the restricted
   grammar above; missing `doc_id`/`authority`/`claims` keys are errors.
2. **Enum validity** — `authority` and every claim's `type`/`reality_binding` must be one of
   the values listed above.
3. **Claim-ID format + ownership** — every claim ID matches `<doc_id>.<TYPE>-<N>` and its
   `<doc_id>` prefix matches the declaring document's own `doc_id`.
4. **ID collisions** — no claim ID appears twice anywhere in the corpus.
5. **Dangling references** — every `doc_id` named in any `supersedes`/`amends` list must
   exist as some document's own `doc_id`.
6. **Unenumerated inheritance** — every `amends` entry must name at least one claim.
7. **Orphan goldens** — if document B's `supersedes` list names document A, and A is not
   itself in B's `amends` list (i.e., B claims *full* supersession of A), then A's own
   `authority` must not still be `golden` — a fully-superseded document that still claims
   golden status is the orphan-goldens case DOC-102 §8 tracks as a corpus health metric.

Lint does **not** attempt semantic conflict detection (DOC-102 §5: that is SAGA's own,
separate, LLM-driven proposal surface, explicitly out of scope for the deterministic linter
built in this step) or reality-binding verification against actual test/telemetry results
(§9 step 2, `saga.manifest.yaml`, a later build-sequence step) — this tool is the mechanical,
deterministic half only, exactly as step 1 scopes it.

## 4. Non-goals of this pass (explicitly deferred, per DOC-102 §9)

- `saga.manifest.yaml` binding claims to real tests/telemetry (step 2).
- Divergence/conflict queues in Back Office (step 3).
- The SAGA agent itself as a query tool surface (step 4).
- Semantic conflict detection (step 5).
- Corpus-wide enforcement turning warning → gate (step 6).

This step is exactly what DOC-102 §9 step 1 names: schema, convention, and a lint tool —
nothing more, retrofit onto the docs that exist today.
