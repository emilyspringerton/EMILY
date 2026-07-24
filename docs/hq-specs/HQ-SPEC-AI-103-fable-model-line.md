---
doc_id: AI-103
authority: draft
supersedes: []
amends: []
claims:
  - id: AI-103.NAR-1
    type: NAR
    reality_binding: specified
  - id: AI-103.POL-1
    type: POL
    reality_binding: building
  - id: AI-103.BEH-1
    type: BEH
    reality_binding: verified
  - id: AI-103.IFC-1
    type: IFC
    reality_binding: specified
---

# HQ-SPEC-AI-103 — FABLE: The In-House Model Line

**Status:** DRAFT v0 — pending Emily Prime review
**Custody:** Emily OS
**Governs:** EINHORN's sovereign language-model program — data engine, training, evaluation, serving, and the evolution loop
**Implements:** S16 backlog (AI layer FABLE); NORN (PRIME-101) as the governing kernel; long-term doctrine of autonomous, auditable infrastructure over API dependency
**Honest premise, stated first:** a GPT-2-class model will not out-reason a frontier model, and FABLE does not pretend otherwise. FABLE wins on a different axis: fully auditable weights, provenance-tagged training data, deterministic serving, zero external dependency, and near-zero marginal cost on high-volume scoped tasks. Sovereignty is a ladder, and this spec climbs it honestly.

---

## 1. Northstar Statement

Every high-volume, well-scoped inference task in EINHORN runs on weights EINHORN trained, from data EINHORN can produce the provenance of, gated by evals rooted in reality. Frontier APIs remain in the stack — as the ceiling FABLE is graded against and the fallback FABLE routes to — but the dependency arrow reverses over time: frontier becomes the exception handler, FABLE the workhorse.

The strategic asset is not the model. It is the **oracle-graded data exhaust** the system already produces: EPS headlines graded by filed 8-Ks, recon matches confirmed by settled transactions, KAREN proposals accepted/rejected by Emily Prime, conflict calls adjudicated for SAGA, playtest metrics from SHANKPIT. Every NORN loop is a label factory whose labels are rooted in reality. That corpus is unpurchasable, and it is exactly the shape small models get strong on. *(AI-103.NAR-1 — specified: narrative framing.)*

## 2. Why GPT-2 Class Is the Right Opening

- **Trainable end-to-end in-house.** The 124M-parameter reproduction is a solved, cheap problem (llm.c / nanoGPT lineage — order of tens of dollars in rented GPU-hours; 1.5B in the hundreds). The entire training run fits in one auditable repo: data snapshot hash in, checkpoint hash out. *(AI-103.POL-1 — building: local CPU LoRA rounds have run in `gpt2-alpine-c`; the actual Colab T4 GPU run is still HITL-blocked, S26-04/HITL-06.)*
- **Small models + narrow tasks + good labels beats big models + no ownership** for the task portfolio in §3. Classification, extraction, tagging, templated generation — these saturate well below frontier scale when the fine-tuning data is clean and the output schema is enforced (Iduna-validated JSON, the house pattern).
- **Auditability is a feature of scale.** A 124M model's training set can be enumerated, its eval history replayed, its behavior diffed checkpoint-to-checkpoint. Try that with an API. *(AI-103.BEH-1 — verified: `gpt2-alpine-c`'s entropy source is functional and measured — base H_mean=4.4877, local 300-step fine-tune H_mean=4.6602, S26-05.)*
- **The ladder exists.** Nothing in this spec caps out at GPT-2. The pipeline (data engine, eval harness, NORN gates, serving) is size-agnostic; parameters are a config value. E0 through E4 in §6. *(AI-103.IFC-1 — specified: E0 (S145-03) not started.)*

## 3. Task Portfolio — Scoped Honestly

**FABLE-eligible (v0/v1):**
- FATBABY signal extraction & classification (governance signal types, entity tagging, `DiscoveryIdentity` assists)
- EPS headline generation — the flagship, because its oracle (the filed 8-K) is the strongest in the house
- PR body triage, dedup, and source classification on the unified BW/PRN feed
- SHANKPIT NPC dialogue and edu_vm hint text (charm proving ground, low stakes)
- Recon-exception summarization and KAREN's journal-entry drafting (behind existing approval tiers)

**Explicitly NOT FABLE (stays frontier, via the existing `LLMClient` interface):**
- Emily Prime's chief-of-staff reasoning, Jon Stockwell's strategy synthesis, Tyler's code work, SAGA's semantic conflict detection, anything legal/compliance-adjacent, and all spec drafting. Routing a task down-model is a *gated promotion with eval evidence*, never an aspiration.

## 4. Full Stack

**4a. Data Engine (`fabledata`)** — Go, house patterns throughout:
- Readers over the NDJSON stores (FATBABY events, NORN registries, STAGEDOOR takes when SIM-101 lands) → normalized examples with **per-record provenance**: source event hash, oracle that labeled it, label date, license class.
- Snapshotting: a dataset is a content-addressed manifest of record hashes — immutable, diffable, citable. Training never reads live stores; it reads snapshots.
- Filters as versioned config: dedup (minhash), PII scrub, license gate (own-exhaust and public-domain first; anything else enters through an explicit, logged decision), contamination guard (eval records are tombstoned out of every training snapshot by hash — mechanical, not procedural).

**4b. Tokenizer (`fabletok`)** — BPE trained on the domain corpus, versioned like everything else. Financial filings and tickers tokenize badly on general-web vocabularies; a domain tokenizer is cheap alpha for small models.

**4c. Training (`fabletrain`)** — thin orchestration around the llm.c/nanoGPT lineage. A **training bundle** = hash over {dataset snapshot, tokenizer, code rev, config, seed}. Bundle in, checkpoint out, both logged. Fine-tunes and LoRA adapters are bundles too. Hardware: local 24GB-class GPU for fine-tunes; rented A100/H100 hours for pretrains; budget is a config value with a KAREN-visible cost event per run.

**4d. Eval Harness (`fableeval`)** — frozen suites per task, Iduna-schema'd outputs, exact/graded scoring. Every suite declares its **reality root** (which oracle produced its labels). Suites are NORN oracles and live under oracle governance (PRIME-101 §5): frozen, versioned, rotated only against held-out reality.

**4e. Serving (`fableserve`)** — checkpoints exported to a quantized standard format, served behind the existing Go `LLMClient` interface — meaning every agent in the fleet can be pointed at FABLE **without code changes**. The router (see §6, E4) also implements `LLMClient`: callers never know which brain answered, but the event log always does — every inference records model hash, prompt hash, and route decision.

## 5. "Let It Evolve" — NORN Instantiations

New rows for the PRIME-101 §6 table:

| Instantiation | Proposer (Skuld) | Oracle (Verdandi) | Tier | Reality root |
|---|---|---|---|---|
| FABLE checkpoints (per task) | fabletrain bundles | Frozen task suites | autonomous | Task's oracle (8-Ks, settlements, human decisions) |
| Dataset snapshots | fabledata pipelines | Contamination/provenance audit + downstream eval delta | prime_ack | Source event stores |
| Tokenizers | fabletok candidates | Compression + downstream eval delta | prime_ack | Corpus |
| Prompts & adapters | Agents + humans | Frozen task suites | autonomous | Task oracle |
| Router policies | Routing proposals | Cost/quality frontier on logged traffic | prime_ack | Production outcomes |

**Löbian rules, model edition (the load-bearing paragraph):**
1. No FABLE model's outputs may enter a training snapshot unless graded by a reality-rooted oracle first. Self-generated data is admissible only as *graded* data.
2. No FABLE model may author or amend an eval suite that any FABLE model is graded by. Suites descend from reality and humans only.
3. Frontier-model distillation, if pursued, is a licensing/ToS decision made with counsel per provider — flagged here, not assumed. The architecture doesn't require it; the own-exhaust corpus is the primary diet.
4. Checkpoint lineage is tracked like all NORN lineage; a model whose training data includes another model's outputs carries that ancestry in its hash chain, forever queryable.

## 6. The Evolution Ladder

- **E0 — Fine-tune published GPT-2 weights** on one task (EPS headlines). Purpose: stand up the whole stack end-to-end and lose to the API on quality while winning on cost — measured, not assumed.
- **E1 — Pretrain FABLE-124M from scratch** on the curated corpus + own exhaust, domain tokenizer. First fully sovereign checkpoint: every byte of provenance answerable.
- **E2 — Scale to 350M/774M per-domain models** where E1 evals show headroom; task-specific fine-tunes branch from the pretrained trunk.
- **E3 — 1.5B trunk** when (and only when) E2 saturates on tasks that matter; budget gate through KAREN.
- **E4 — The router** (`LLMClient` over FABLE-first, frontier-fallback): each task's traffic migrates down-model only by gated promotion. Router logs become the next training set — production is the final oracle.
- **Quiescence is a valid outcome per rung.** If a task's eval says the 124M is at parity, scaling it further is budget vandalism; NORN declares the fixed point and the money goes elsewhere.

## 7. Back Office Metrics

- Win-rate vs. incumbent (frontier or prior checkpoint) per task suite — the promotion currency.
- Cost per 1k tasks, FABLE vs. API, with router offload % over time — the sovereignty curve, the number that makes this program legible to a CFO.
- Regression catches at the gate; contamination audit results (must be zero findings).
- Data freshness: age of newest oracle-graded records per task; label acquisition rate (the factories' output).
- Serving determinism drift: same prompt hash + model hash ⇒ same output hash at temperature 0, monitored.

## 8. Build Order

1. `fabledata` snapshotting over the EPS headline corpus (the richest oracle-graded set in the house) + contamination tombstoning.
2. `fableeval` with the EPS suite frozen as oracle v1 — *before any training*, so day-one numbers are honest.
3. E0 fine-tune; wire NORN; first gated promotion (or honest rejection) of a FABLE checkpoint.
4. `fableserve` behind `LLMClient`; shadow-route EPS traffic (FABLE answers logged, frontier answers served) — the ELP shadow-books pattern applied to inference.
5. Promote to live routing on eval + shadow evidence. One task fully sovereign, end to end.
6. `fabletok` + E1 pretrain. Then breadth per §3.

## 9. Open Questions

- **Corpus beyond own-exhaust for E1 pretraining** — which public-domain/open corpora, at what mix; provenance standards apply to acquisitions.
- **GPU posture** — owned 24GB-class box for the daily loop vs. purely rented burst; KAREN gets the cost model either way.
- **Distillation stance** (per Löbian rule 3) — counsel question, parked with the FIN-099 licensing set.
- **Naming the checkpoints** — FABLE-124M is a size, not a character. The fiction universe has opinions about what stories are called; first promoted sovereign checkpoint deserves a name from it.

---

*Small model, long memory, honest grades. CLEAN BUILDS FIRST.*
