# HQ-SPEC-PRIME-101 — NORN: The Loop Kernel

**Status:** DRAFT v0 — pending Emily Prime review
**Custody:** Emily OS
**Governs:** The single, canonical implementation of the propose→grade→gate→promote loop, instantiated by every improving system in EINHORN_INDUSTRIAL
**Implements:** HQ-SPEC-PRIME-097 (Fixed Points). Where 097 is the mathematics, 101 is the machine.
**Working name:** NORN. Skuld proposes (what shall be), Verdandi evaluates (what is becoming), Urd records (what has become). Three phases, one loop, fate as infrastructure.

---

## 1. Why This Spec Exists

The same loop has now been independently derived at least six times: the EPS headline grader (8-K as oracle), the GameEvolutionEngine (structured improvement proposals), the recon matcher (confirmed suggestions as labels), the reward compiler (frozen sim evals), KAREN's proposal funnel (rejection reasons as training signal), and sim-fidelity improvement (reality deviation as grade). Each re-derivation re-argues the Löbian hazard from scratch and re-invents the no-regression gate slightly differently.

NORN extracts the loop once. Domains stop owning loops; they own **instantiations** — a proposer, an oracle binding, and a gate policy. The hazard analysis is done once, at the kernel, and inherited everywhere.

**Cardinal rule:** NORN decides *promotion*, never *execution*. It answers "is this artifact now golden?" It never moves money, never actuates hardware, never deploys code. Execution lives in domain systems behind their own gates (KAREN's payment approval, the Springerton Seam's actuation ladder). A kernel that could act would be a kernel worth attacking; this one can only bless.

## 2. The Loop, Formally (per PRIME-097)

- Artifacts live in a **quality lattice** per domain; the promotion operator must be **monotone**: a promoted artifact is never worse than the incumbent on the frozen eval set. No-regression is not a policy preference — it is the monotonicity condition under which the fixed-point theorems apply at all.
- **Frozen evals are the fixed metric.** Banach-style convergence arguments require the distance function to hold still. An eval set that drifts under the loop it grades destroys the contraction; therefore eval mutation is a *separate, slower loop* with its own oracle (§5).
- The **fixed point** we converge toward per domain is "no proposable improvement passes the gate" — quiescence, which is a reportable state, not a failure.
- **Löbian hazard, resolved structurally:** no artifact may be graded by an oracle whose lineage includes that artifact's proposer. Lineage is tracked by content hash; the check is mechanical, not procedural. A system may never prove its own soundness; it may only be graded by something rooted outside itself. The root of every oracle chain must terminate in **reality**: a filed 8-K, a bank-confirmed transaction, physical telemetry, a human decision. Reality is the root oracle; everything else is a cache of it.

## 3. Core Abstractions (Go)

Five interfaces, one event store. Bare metal doctrine applies: NORN is a library plus a small daemon, not a platform.

```go
// A content-addressed candidate for promotion.
type Artifact struct {
    Hash     string            // sha256 over canonical bytes
    Kind     string            // "eps_extractor", "reward_profile", "matcher_model", ...
    Lineage  []string          // hashes of proposer + inputs that produced it
    Payload  []byte            // or a content-store reference
    Meta     map[string]string
}

// Skuld — emits candidates. LLM-driven or procedural; NORN doesn't care.
type Proposer interface {
    Propose(ctx context.Context, incumbent Artifact, budget Budget) ([]Artifact, error)
}

// Verdandi — grades a candidate against a frozen, versioned eval set.
type Oracle interface {
    Version() OracleVersion              // content hash of the frozen eval set + grader code
    Grade(ctx context.Context, a Artifact) (Report, error)
}

// The no-regression gate + domain policy (approval tier, extra invariants).
type Gate interface {
    Decide(incumbentReport, candidateReport Report, policy GatePolicy) Decision
}

// Urd — append-only. The only writer of promotion truth.
type Registry interface {
    Golden(kind string) (Artifact, error)          // current fixed point
    Record(e Event) error                          // append-only, NDJSON + pg projection
    History(kind string) ([]Event, error)
}
```

**Event vocabulary** (append-only, house NDJSON store, `pgx/v5` projection pattern):
`proposal_submitted` · `evaluation_completed` · `gate_evaluated` · `artifact_promoted` · `artifact_rejected` · `regression_detected` · `oracle_frozen` · `oracle_rotated` · `quiescence_declared` · `lineage_violation_blocked`

Every `artifact_promoted` also emits an `ApplePublished` golden-log entry. Promotions are civic events.

## 4. Gate Policy & Approval Tiers

Gate policy is declarative config (versioned in repo, Emily OS custody), per instantiation:

```yaml
kind: reward_profile
oracle_min_version: <hash>
no_regression_on: [pose_error, fall_rate, energy]     # monotonicity dimensions
improvement_required_on_any: true
approval_tier: autonomous          # autonomous | prime_ack | biometric
quiesce_after_rejections: 50       # declare fixed point, stop burning budget
proposal_budget: {per_day: 200, concurrent: 8}
```

**Tiers:**
- `autonomous` — gate decision is final (pure-software artifacts with reality-rooted oracles: matcher models, extractors).
- `prime_ack` — Emily Prime acknowledges before promotion (early-life domains; KAREN's journal proposals start here per FIN-099 training wheels, and *tier relaxation is itself a gated promotion* with rejection-rate evidence as the eval).
- `biometric` — companion-app approval; reserved for artifacts whose promotion feeds a physical or financial execution path (deployment bundles at the Springerton Seam, policy files in ELP). NORN promotion at this tier still doesn't execute anything — it makes the artifact *eligible* for the domain's own execution gate. Two locks, deliberately.

## 5. Oracle Governance — The Slow Loop

Evals rot; graders have bugs; reality drifts. Oracle updates are therefore their own NORN instantiation, with three hard rules:

1. **Frozen means frozen.** An `oracle_frozen` event fixes the eval-set hash. Grading always cites the oracle version; reports against different versions are incomparable and the kernel refuses to compare them.
2. **Oracle candidates are graded only by reality or by humans** — held-out ground truth (later 8-Ks, settled bank transactions, physical telemetry, human adjudication) — never by any loop the oracle will grade. This is the lineage check applied to the graders themselves.
3. **Rotation is versioned, never in-place.** New oracle → incumbent artifacts are re-graded → the quality lattice is re-based explicitly, as an event. History remains interpretable forever because every report names its ruler.

## 6. Instantiation Table (initial)

| Instantiation | Proposer (Skuld) | Oracle (Verdandi) | Tier | Reality root |
|---|---|---|---|---|
| EPS headline extractor | LLM prompt/code variants | Frozen labeled set graded by filed 8-Ks | autonomous | EDGAR filings |
| Recon matcher | Matcher model/threshold variants | Confirmed/rejected suggestion corpus | autonomous | Bank-settled transactions |
| KAREN journal proposals | KAREN | Emily Prime accept/reject history | prime_ack | Human decision |
| GameEvolutionEngine | Structured JSON improvement proposals | Frozen playtest/metric suites | autonomous | Player telemetry |
| Reward & DR profiles | GEE reward-weight/DR variants | Frozen sim eval suites | autonomous | — (sim) |
| Deployment bundles (Seam) | Training pipeline | Frozen sim evals **+ prior-rung telemetry** | biometric | Physical telemetry |
| Sim fidelity | Sim-parameter proposals | Real-vs-sim deviation on returned telemetry | prime_ack | Physical telemetry |
| Gate policies themselves | Emily Prime directives only | Human review | biometric | Human decision |

Note the last row: NORN's own configuration is inside NORN, at the strictest tier, with humans as the oracle. The kernel governs itself the same way it governs everything — which is the only honest answer to "who watches the watcher" that doesn't regress infinitely: the chain terminates in a person with a thumbprint.

## 7. Operational Surface

- **Library:** `pkg/norn` — the five interfaces + reference implementations (NDJSON registry, hash/lineage utilities, budget/token-bucket reuse from the agent core).
- **Daemon:** `nornd` — runs instantiation schedules (cron-based autonomy pattern), enforces budgets, publishes events to feedserver for Back Office live view.
- **CLI:** `norn propose <kind>` · `norn grade <hash>` · `norn status <kind>` (incumbent, pending, quiescence state) · `norn history <kind>` · `norn freeze-oracle <kind>`.
- **Back Office metrics:** promotion rate, rejection reasons distribution, regression catch count (gate saves), oracle staleness (age since last reality-rooted validation), time-in-gate by tier, quiescent domains list, lineage violations blocked (should be zero; each one is an attempted Löbian shortcut, and each is an incident).

## 8. Migration & Build Sequence

1. `pkg/norn` interfaces + NDJSON registry + lineage checker + property tests (monotone promotion, replay determinism, lineage-violation rejection, oracle-version comparability refusal).
2. **First migration: EPS headline extractor** — the loop that invented the pattern proves the kernel. Success = identical promotion decisions replayed through NORN from historical events.
3. GameEvolutionEngine second (most structurally similar), recon matcher third (lands with KAREN Phase 0/1 per FIN-098 build order).
4. `nornd` + Back Office view + Apples hookup.
5. Seam and ELP-policy instantiations land with their host specs (SIM-100 §6, FIN-099 §6) — those specs now *reference* NORN rather than restating the loop. Amendment notes to both, append-only.

## 9. Anti-Goals

- NORN never executes side effects. No payments, no actuation, no deploys, no writes outside its registry.
- NORN never trains models; it grades and promotes artifacts that training produces.
- NORN is not a workflow engine, not a scheduler beyond its own cadence, and not a place for domain logic to hide. If an instantiation needs kernel changes, the kernel is probably right and the instantiation is probably confused — burden of proof on the domain.
- No fast path around the gate. There is no `--force`. There will never be a `--force`.

## 10. Open Questions

- **NORN as the name** — the three-phase mapping is proposed above; confirm or overrule.
- **Eval corpus storage** — content-addressed blob store shape: flat files + manifest vs. dedicated store; decide at build step 1.
- **Proposal budgeting strategy** — flat budgets v0; bandit-style allocation across proposers is a later NORN instantiation (naturally: budget policies graded by realized improvement).
- **Cross-domain lineage** — when one domain's golden artifact is another's input (reward profile → deployment bundle), promotion cascades need ordering rules; draft answer is "downstream re-grades on upstream promotion, never auto-promotes," to be pressure-tested in step 3.

---

*Skuld proposes. Verdandi grades. Urd remembers. CLEAN BUILDS FIRST.*
