---
doc_id: DOC-102
authority: draft
supersedes: []
amends: []
claims:
  - id: DOC-102.POL-1
    type: POL
    reality_binding: verified
  - id: DOC-102.POL-2
    type: POL
    reality_binding: verified
  - id: DOC-102.POL-3
    type: POL
    reality_binding: verified
  - id: DOC-102.INV-1
    type: INV
    reality_binding: specified
  - id: DOC-102.BEH-1
    type: BEH
    reality_binding: verified
  - id: DOC-102.INV-2
    type: INV
    reality_binding: specified
  - id: DOC-102.NAR-1
    type: NAR
    reality_binding: specified
  - id: DOC-102.BEH-2
    type: BEH
    reality_binding: verified
---

# HQ-SPEC-DOC-102 — SAGA: The Documentation Curation Lifecycle

**Status:** DRAFT v0 — pending Emily Prime review
**Custody:** Emily OS
**Governs:** The documentation corpus as a reconciled knowledge system: claim binding, conflict resolution, supersession, and the divergence lifecycle between what the specs say and what the software does
**Agent:** SAGA — the Librarian. Norse goddess of history; keeper of the story so far. Registered in Iduna as `saga`.
**Prime principle (agile, stated without apology):** working software is truth. Divergence between docs and reality never blocks shipping; it blocks *golden status*. The corpus converges on reality, not the other way around.

---

## 1. The Problem, Precisely

The spec corpus is append-only intent. As it grows, three failure modes compound:

1. **Conflict:** two golden documents assert incompatible things and nothing detects it.
2. **Vaporware debt:** claims with no running counterpart, indistinguishable at read time from claims that are live.
3. **Dark matter:** running behavior with no documented claim — the more dangerous direction, because it is invisible to every reader and every audit.

The fix is the house pattern applied to knowledge: a three-way match between **intent** (specs), **books** (the claim ledger), and **reality** (running software), with exceptions queued, aged, and resolved — KAREN's reconciliation loop, for words.

## 2. The Claim Ledger

Documents are too coarse to reconcile. The unit of reconciliation is the **claim**: an addressable assertion with a stable ID.

- **ID format:** `<DOC>.<TYPE>-<N>` — e.g., `FIN-098.INV-2` (invariant), `SIM-100.BEH-4` (behavior), `PRIME-101.IFC-1` (interface contract). IDs are permanent; a superseding doc that restates a claim *cites* the old ID and issues its own. *(DOC-102.POL-1 — verified: `emily saga lint` enforces this exact format and checks doc-of-origin ownership.)*
- **Claim types:** `INV` (invariant — must always hold), `BEH` (behavior — does X when Y), `IFC` (interface contract), `MET` (metric definition), `POL` (policy), `NAR` (narrative/rationale — explicitly *unverifiable by design*, exempt from reality-binding but still subject to authority lifecycle).
- **Binding manifest:** a versioned file per repo (`saga.manifest.yaml`) mapping claim IDs to their verification anchors:

```yaml
- claim: FIN-098.INV-2
  verified_by:
    - test: ledger/payment_gate_test.go::TestNoUnapprovedMovement
    - telemetry: back_office.metric.unapproved_payment_attempts == 0
- claim: PRIME-101.IFC-1
  verified_by:
    - test: pkg/norn/registry_test.go::TestAppendOnly
- claim: SIM-100.BEH-4
  verified_by: []          # unbound — this is vaporware debt, visible on purpose
```

- **Verification methods:** automated test, telemetry probe, or `human_attestation` (dated, named, expiring — attestations rot and must be renewed). Every method is machine-checkable for *presence*; tests and probes are machine-checkable for *truth*.

## 3. Two Independent Status Axes

Every document carries structured frontmatter with both:

**Axis 1 — Authority** (what the corpus says about the doc):
`draft → golden → amended → superseded` *(DOC-102.POL-2 — verified: `emily saga lint` enforces this exact enum on every document's `authority` field.)*
- Append-only doctrine holds absolutely: supersession is a `supersedes: [HQ-SPEC-X]` pointer in a *new* document, never an edit to an old one. Amendment is a new doc with `amends:` scoped to named claims. *(DOC-102.INV-1 — specified: not yet mechanically enforced — that would need a git-history check, out of scope for the lint tool built so far.)*
- The supersession graph must be a DAG with no orphaned goldens: if doc B supersedes doc A only partially, B must enumerate which of A's claims it inherits, retires, or replaces. Unenumerated claims are a lint error, not a judgment call. *(DOC-102.BEH-1 — verified: `emily saga lint` rejects any `amends` entry with an empty `claims` list.)*

**Axis 2 — Reality-binding** (what the running system says about each claim):
`specified → building → running → verified → diverged` *(DOC-102.POL-3 — verified: `emily saga lint` enforces this exact enum on every claim's `reality_binding` field.)*
- Computed per-claim, rolled up per-doc. `verified` requires a passing bound verification; `diverged` fires automatically when a bound verification fails or an attestation expires.
- A document may be **golden and diverged simultaneously.** That is not a contradiction; it is the honest state of an agile system, and hiding it is the actual failure. Golden ⁠+ diverged means "this was the agreed intent, and reality has moved" — which is precisely the signal that triggers curation. *(DOC-102.NAR-1 — narrative, exempt from reality-binding by design.)*

## 4. The Divergence & Conflict Queues

Two exception queues, Back Office-surfaced, aged like reconciliation exceptions (FIN-098 Invariant 3 pattern — items older than N days block the corpus health gate):

**Divergence queue** (docs vs. reality):
- *Claim-without-code* (vaporware debt): bound verification absent or failing.
- *Code-without-claim* (dark matter): detected by coverage heuristics — exported APIs, event types, config surfaces, and Iduna scopes that appear in no manifest. Dark matter items require either a new claim (documenting reality) or a deprecation (changing reality). Silence is not an option the tooling permits.

**Conflict queue** (docs vs. docs):
- *Mechanical detection:* the same claim ID redefined by two live docs; two golden docs binding contradictory values to the same policy/metric/interface anchor.
- *Semantic detection:* SAGA (LLM-driven) reads the corpus and *proposes* conflicts it suspects — always as proposals into the queue, never as unilateral rulings. Human confirmation required; confirmed/rejected conflict proposals become SAGA's labeled improvement corpus, per the standard oracle-graded loop.

**Resolution is always a new document** citing every party to the conflict, entering at `draft`, promoted through the normal lifecycle. The queues empty only through the corpus, never through quiet edits.

## 5. SAGA — Role and Boundaries

Per SYSTEM_GOVERNANCE: SAGA publishes observations and proposals, never directives.

**SAGA does:**
- Maintain the claim index and supersession graph; run lint (orphan goldens, unenumerated inheritance, ID collisions).
- Run the divergence and conflict detectors; curate both queues; propose resolution drafts, supersession pointers, and manifest bindings.
- Answer corpus queries as a first-class tool surface for every other agent: `saga which-doc-governs <topic>`, `saga status <claim>`, `saga conflicts`, `saga gaps <repo>`. When Tyler or KAREN needs to know what the current truth is, they ask SAGA — agents stop grepping the corpus and start querying it.
- Publish corpus health to Back Office and golden-log every authority transition to Apples.

**SAGA does not:**
- Promote documents (that is a NORN gate decision), edit any document, hold write access to code, or adjudicate conflicts. The Librarian catalogs, detects, and proposes; humans and NORN decide. *(DOC-102.INV-2 — specified: the SAGA agent itself is not built yet, build-sequence step 4.)*
- Overlap with Apples: **Apples is the log of what was promoted; SAGA is the curator of what is currently true and where it conflicts.** Apples remembers events; SAGA maintains the living index over them. SAGA writes to Apples, never the reverse.

## 6. NORN Instantiation

Documentation is an artifact kind in the kernel — this spec adds a row to PRIME-101 §6:

| Instantiation | Proposer (Skuld) | Oracle (Verdandi) | Tier | Reality root |
|---|---|---|---|---|
| Corpus documents | SAGA drafts + human authors | Claim verification (CI) **+ Emily Prime adjudication** | prime_ack | Running software; human decision |
| SAGA's conflict detector | SAGA model/prompt variants | Confirmed/rejected conflict-proposal corpus | autonomous | Human adjudication history |

Promotion to golden requires: lint clean, no unresolved conflicts naming the doc, and every non-`NAR` claim at `building` or better with a manifest entry (binding may be a declared plan, but it must *exist*). Reality-binding thereafter runs continuously; golden status is held, not archived.

## 7. The Curation Lifecycle, End to End

1. **Intake:** new doc lands as `draft` with frontmatter; SAGA extracts claims, assigns IDs, lints the supersession graph.
2. **Binding:** author (or SAGA proposal) populates manifest entries; unbound claims are visible debt from day one.
3. **Promotion:** NORN gate, `prime_ack` tier. `ApplePublished` on golden.
4. **Monitoring:** CI runs claim verification per commit; telemetry probes run continuously; attestations age.
5. **Exception:** divergence or conflict fires → queue → Back Office.
6. **Resolution:** new doc drafted (by human or SAGA proposal), cycle restarts. The corpus's history of resolutions *is* the audit narrative of how the system's understanding of itself improved — Urd's shelf, in SAGA's keeping.

## 8. Corpus Health Metrics (Back Office)

- **Verification coverage:** % of non-`NAR` claims at `verified`.
- **Divergence queue depth & max age** — the headline number; it is to knowledge what reconciliation lag is to cash.
- **Dark matter count** per repo (target: monotone decreasing; every new API without a claim increments it in CI).
- **Conflict half-life:** median days from detection to resolving doc golden.
- **Attestation freshness:** expired human attestations outstanding.
- **Supersession graph integrity:** orphan goldens (must be zero), longest unresolved partial-supersession chain. *(DOC-102.BEH-2 — verified: `emily saga lint` checks this exact condition — a fully-superseded document still marked golden.)*

## 9. Build Sequence

1. Frontmatter schema + claim-ID convention + linter (`saga lint`), retrofit onto the live spec series (097–102). This document ships with its own claims tagged — DOC-102 eats first.
2. `saga.manifest.yaml` format + CI gaps report (claim-without-code and dark-matter detection) in one repo — FATBABY, since it has the most running software to reconcile against.
3. Divergence + conflict queues in Back Office; aging rules wired to the corpus health gate.
4. SAGA agent v0: index, query tools, lint — deterministic parts first.
5. SAGA semantic conflict detection as NORN-looped proposals.
6. Corpus-wide rollout; supersession-graph enforcement turns from warning to gate.

## 10. Open Questions

- **Retrofit depth:** claims extracted from all eleven original architecture docs, or only the live HQ-SPEC series forward? (Instinct: series forward; the originals get `historical` authority status — a fifth authority state worth adding if so.)
- **Attestation TTL defaults** per claim type — propose 90 days for `POL`, 180 for `BEH`, indefinite for `NAR`.
- **Cross-repo claim references** — one global namespace (current design) vs. per-repo; global keeps `FIN-098.INV-2` citable everywhere and is the draft answer.
- **Fiction-universe documents** — they are corpus members with `NAR` claims and full authority lifecycle (canon conflicts are real conflicts), but exempt from reality-binding. Confirm this treatment; the Subscriber's identity remaining sovereign is, structurally, a sealed claim.

---

*Skuld proposes, Verdandi grades, Urd remembers — and SAGA keeps the shelf in order. CLEAN BUILDS FIRST.*
