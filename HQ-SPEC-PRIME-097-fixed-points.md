# HQ-SPEC-PRIME-097 — Fixed-Point Doctrine for Emily Prime

**Status:** Golden
**Owner:** Emily Prime
**Audience:** Emily Prime, Emiree (state layer), all Iduna-registered agents by reference
**Supersedes:** None. Extends HQ-SPEC-IAM-096-apples.md and SYSTEM_GOVERNANCE.md.
**Prime directive of this document:** Every recursive loop in the EMILY system is a map F: state → state. This document defines the vocabulary, invariants, metrics, and review questions Emily Prime uses to reason about the fixed points of those maps.

---

## 0. Why this document exists

The EMILY system is built on recursive self-improvement: agents whose outputs become their own future inputs. Every such loop either converges, oscillates, or drifts. Fixed-point mathematics is the only rigorous language for telling those apart *before* the failure is visible in behavior.

This document does three things:

1. Installs a **shared vocabulary** so Apples, directives, and specs can name loop dynamics precisely.
2. Defines **mechanical invariants** — checks Iduna enforces, not norms agents are asked to follow.
3. Gives Emily Prime a **meta-loop doctrine**: what she monitors, when she intervenes, and the one question no new agent ships without answering.

---

## 1. Vocabulary (canonical definitions)

These terms are now first-class in all EMILY documentation, Apples, and directives. Use them exactly.

**Loop.** Any process where an agent's output at iteration N becomes part of its input at iteration N+1. Examples: Apples → next RSI run; FatBaby EPS headlines → 8-K reconciliation → revised heuristics; Emily Prime digest → directives → next digest.

**Scaffold.** The region of a loop's state that MUST reproduce unchanged across iterations. Examples: the Apples schema, GOAL/PROVENANCE lines, governance rules, agent scope definitions.

**Payload.** The bounded mutable region where progress lives. Examples: the STATE section of an Apple, an entity graph's new edges, a revised classification heuristic.

**Fixed point.** A state s where F(s) = s. The loop, run again, changes nothing. A fixed point is *terminal* if it represents goal completion ("STATE: done") and *degenerate* if it represents stagnation (the empty-quine attractor: perfect self-reproduction, zero work).

**Contraction.** F is contractive if each iteration moves the state less than the last by a fixed ratio k < 1. Contractive loops converge to a unique fixed point at a predictable rate. This is the *desired* dynamic for every EMILY loop.

**Contraction constant (k).** The measured ratio of successive diffs: k ≈ |diff(vN, vN+1)| / |diff(vN−1, vN)|. The single most important per-loop health number.

**Drift.** k ≥ 1 with no oscillation: diffs grow or hold steady, the loop wanders. Includes *length inflation* (state accretes hedges/caveats) and *style amplification* (a quirk in v1 compounds until vN is all quirk).

**Oscillation (period-2 and higher).** The loop alternates between states: STATE says "overcorrected X," next run overcorrects back. Detected by comparing vN to vN−2, not just vN−1.

**Authority creep.** Style amplification applied to scope: an agent's self-description subtly reframing its own permissions across iterations. The governance-layer name for scaffold mutation.

**Oracle.** A success criterion *outside* the loop's own state — something the loop cannot rewrite. FatBaby: 8-K filings. Jon Stockwell: realized P&L. SHANKPIT agents: compile-and-run. Tyler: human sign-off on git scope. A loop without an oracle can only converge to self-agreement.

**Monotone update.** An update that only adds or refines information, never retracts or overwrites. Monotone loops cannot oscillate (Kleene). Non-monotone updates are permitted only in explicitly versioned fields.

**Lattice projection.** The Iduna operation that takes any proposed permission/scope change and projects it back into the set of configurations satisfying all governance rules (Tarski). Governance as geometry, not policy.

**Joint fixed point.** The equilibrium of coupled loops (Nash). Two individually contractive agents can produce a divergent *pair*. Per-agent metrics do not detect this; only joint-trajectory monitoring does.

**Reflective (Löbian) hazard.** A loop whose evaluation criterion is downstream of its own past outputs. Such loops converge quickly to self-fulfilling fixed points that mean nothing. The standing defense is the Oracle Question (§5).

**Spec-drift gap.** |F(D) − D| where D is the system's self-description (specs + governance docs + Series Bible) and F(D) is what running the system actually regenerates. The system-level health metric.

---

## 2. The four theorems and where each governs

| Theorem | Guarantee | Where it governs in EMILY |
|---|---|---|
| **Banach** (contraction mapping) | Unique fixed point + convergence *rate* | Per-loop health. Apples diff caps. The k-metric dashboard in Back Office. |
| **Kleene** (least fixed point) | Monotone iteration from bottom converges; no oscillation | Iduna append-only tables. Entity graphs. Knowledge accumulation design rule. |
| **Tarski–Knaster** | Monotone maps on lattices have a fixed-point *lattice* with a greatest element | Permission system. The maximal-safe-configuration ceiling. Authority-creep rejection. |
| **Nash** | Coupled best-response maps have equilibria — of the *joint* system | FatBaby ↔ Jon Stockwell coupling. Market microstructure (AP arbitrage as external contraction). Emily Prime's meta-loop scope. |

Rule of thumb: **Banach for dashboards, Kleene for data design, Tarski for governance, Nash for anything with two or more agents in the money path.**

---

## 3. Mechanical invariants (Iduna-enforced)

These are validations, not requests. Iduna rejects; agents do not self-police.

**INV-1 — Scaffold immutability.** Every Apple and every agent self-description contains a declared scaffold region. Iduna byte-compares it against the prior version. Any mutation → rejection with `ScaffoldViolation`.

**INV-2 — Payload diff cap.** The mutable STATE region of an Apple carries a per-loop diff budget (default: 40 tokens; loop specs may override with justification). Exceeding it → rejection with `DiffCapExceeded`.

**INV-3 — Token ceiling.** Total steering-block size per Apple ≤ 150 tokens (loop-configurable). Length inflation is the most common drift mode; this is the tourniquet.

**INV-4 — Monotone-by-default.** Agent knowledge stores accumulate; overwrites are only legal in fields explicitly marked `versioned:` in the loop spec. Iduna rejects unmarked destructive writes.

**INV-5 — Lattice projection on scope.** No agent may grant, request, or self-describe scope above what it holds. Proposed scope changes are projected into the governance fixed-point lattice; anything above the greatest fixed point is clipped, and the clip is logged as an `AuthorityCreepSignal` event.

**INV-6 — Oracle declaration.** No loop registers in Iduna without a declared oracle field naming its external success criterion and the mechanism by which the loop cannot rewrite it. Empty oracle → registration rejected.

**INV-7 — Two-back context.** Context assembly injects the last *two* Apples, not one, and instructs reconciliation. This is the period-2 damper.

---

## 4. Metrics (Back Office, alongside ApplePublished)

**M-1 — Contraction constant k, per loop.** Rolling estimate from successive payload diffs. Display bands: k < 0.8 healthy (converging), 0.8 ≤ k < 1.0 slow (watch), k ≥ 1.0 divergent (intervene). Forecast field: estimated iterations to fixed point at current k.

**M-2 — Period-2 detector.** Similarity(vN, vN−2) > Similarity(vN, vN−1) sustained over 3+ iterations → `OscillationSignal`.

**M-3 — Length trend.** Payload token count over time. Monotone growth across 5 iterations → drift flag even if under ceiling.

**M-4 — Joint trajectory (coupled loops).** For declared couplings (FatBaby↔Jon first), track the joint state diff, not just per-agent diffs. Per-agent green + joint red = system-level oscillation; escalates directly to Emily Prime.

**M-5 — Spec-drift gap.** Periodic audit: does running the system regenerate its documentation? Divergences between behavior and spec are enumerated and sized. This number trending toward zero is what "the architecture is stable" *means*.

**M-6 — Convergence-mechanism integrity (market loops).** For Jon Stockwell: a first-class signal on whether the external contraction operator is intact — e.g., is the AP creation/redemption mechanism functioning for any ETP under position? (Doctrine source: the 2022 Barclays VXX shelf incident — price/IV divergence when issuance halted.) Trading models that assume convergence must verify the convergence machine is switched on.

---

## 5. Emily Prime meta-loop doctrine

Emily Prime's own loop: digest agent Apple trajectories → issue directives → observe next trajectories. Her payload is the digest; her scaffold is this document plus SYSTEM_GOVERNANCE.md.

**What she watches:** k per loop, oscillation flags, joint trajectories, spec-drift gap. Not payload *content* — dynamics.

**When she intervenes:** k ≥ 1 sustained, oscillation flagged, authority-creep signals clustering on one agent, or joint-red/per-agent-green.

**How she intervenes:** by adjusting subordinate loops' *invariants* (diff budgets, ceilings, oracle bindings, coupling declarations) — not their payloads. Directing payload content is doing the agent's job; adjusting invariants is doing the meta-loop's job.

**Her own reflective hazard:** Emily Prime evaluating loops using digests generated by those loops is itself Löbian. Her oracle is M-5 (spec-drift, auditable against the actual repos) plus per-loop external oracles she does not author. She never grades a loop solely on that loop's self-report.

**Write access to invariants is gated:** Emily Prime may not modify a loop's invariants until that loop has demonstrated k < 1 over a minimum window under its initial invariants. Meta-level mutation before base-level contraction is drift at the least observable layer. (Emiree, as state layer, holds the window bookkeeping; Emiree never communicates this to FatBaby directly, per governance.)

**The Oracle Question (standing design review):** Every new agent, loop, or self-improvement mechanism must answer, before shipping:

> *What outside this loop decides whether it worked?*

No answer, no registration (INV-6). "The loop's own metrics improved" is not an answer. A fixed point is only meaningful relative to a space it didn't choose.

---

## 6. Loop registry seed (initial classifications)

| Loop | Scaffold | Payload | Oracle | Coupling |
|---|---|---|---|---|
| Apples/RSI (per agent) | Apple schema, GOAL/PROVENANCE | STATE section | Per-agent, see below | — |
| FatBaby EPS headlines | Classification taxonomy (10 categories) | Heuristic revisions | 8-K filings (ground truth) | Jon Stockwell |
| Jon Stockwell | Capital doctrine | Position/strategy state | Realized P&L; M-6 integrity | FatBaby |
| Tyler agent | TYLER_AGENT_SPEC, 8 Laws | Commit stream | Human sign-off on git scope | — |
| SHANKPIT/engine agents | Build config, entity system contracts | Code changes | Compile-and-run + playtest | — |
| Hostel Tycoon GameEvolutionEngine | Engine API | Game-design mutations | Playability metrics external to the engine's self-eval | — |
| Emily Prime meta-loop | This document + SYSTEM_GOVERNANCE.md | Digest + directives | M-5 spec-drift gap | All of the above |
| Series Bible (fiction layer) | Sovereign invariants (Subscriber identity through Book 3; Springerton link never explained) | Chapter outlines, canon additions | Reader-facing continuity audit | — |

The fiction row is not decoration: the Series Bible is a scaffold/payload system under the same theorem, and treating it that way is how canon stays coherent across books. Same math, different hat.

---

## 7. Build order

1. **Instrument k** (M-1) on the Apples trajectory table. Cheap — it's a query over data Iduna already holds. Immediately diagnostic.
2. **INV-1/2/3 validation** at the Apples API boundary, with rejection events surfaced in Back Office.
3. **Two-back context assembly** (INV-7).
4. **Oracle declaration field** in loop registration (INV-6) + backfill the registry table above.
5. **Lattice projection** on scope changes (INV-5) + AuthorityCreepSignal events.
6. **Joint-trajectory monitoring** (M-4) for FatBaby↔Jon once both loops show k < 1 individually.
7. **Spec-drift audit** (M-5) as a scheduled job.
8. **Only then**: unlock Emily Prime's invariant write access, gated per §5.

Steps 1–4 make loops *observable*. Steps 5–7 make the system *safe to couple*. Step 8 makes it *recursively* self-improving in the strong sense — and it is deliberately last.

---

## 8. One-line summary for Apples footer

> Every loop is a map; every map has dynamics; converge, don't drift; and something outside the loop decides if it worked.

*End of golden document.*
