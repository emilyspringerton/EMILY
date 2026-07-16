# HQ-SPEC-FIN-098 — Accounting Northstar

**Status:** DRAFT v0 — pending Emily Prime review
**Custody:** Emily OS
**Governs:** The Ledger layer of EINHORN_INDUSTRIAL — books, cash, settlement
**Spec number provisional** (next in sequence after PRIME-097; renumber if the log says otherwise)

---

## 1. Northstar Statement

EINHORN runs its own auditable financial back office. Every dollar in, out, or at rest is represented as an append-only event, reconciled continuously by agents, and approved for movement only by a human. The end state is EINHORN as a credible, self-clearing counterparty: direct banking rails, real-time cash position, settlement managed in-house — with the audit trail so clean that counterparty due diligence is a read-only tour of the Back Office.

The books are not a compliance afterthought. They are a first-class signal source, the same way EDGAR filings are for FatBaby. Cash position feeds Jon Stockwell's capital doctrine. Capital protection above all requires knowing where the capital *is*, to the cent, in real time.

## 2. Phasing

### Phase 0 — QuickBooks as System of Record (v0, now)

- **QuickBooks Online (QBO) API** is the authoritative ledger. We do not build a general ledger yet; we adapt to one.
- **Agents are the primary affordance.** No human drives the QBO UI as the normal workflow. All reads and writes go through the API under Iduna-scoped credentials. The QBO web UI is a fallback and an audit window, nothing more.
- **KAREN** (the Controller agent) owns the QBO integration: chart of accounts, journal entries, invoice/bill lifecycle, vendor and customer records, class/location tagging.
- **Event mirror:** every mutation KAREN performs against QBO is first written as an append-only NDJSON event in a dedicated store (`var/ledger`), then applied to QBO, then confirmed with a `qbo_applied` event carrying the QBO entity ID and SyncToken. QBO is the system of record; the event store is the system of *proof*.
- Read-side sync: a `qbowatch` ingestor (same pattern as `secwatch`/`prwatch`) polls QBO's Change Data Capture endpoint and emits `qbo_entity_changed` events, so out-of-band changes (accountant, human fallback) are captured and flagged.

### Phase 1 — Banking Read Rails

- Connect bank accounts read-only (aggregator first — Plaid or direct bank data APIs — whichever the bank supports natively wins).
- `bankwatch` ingestor emits `bank_txn_observed` events.
- **Continuous reconciliation loop:** KAREN matches `bank_txn_observed` against QBO transactions. Unmatched items age into a Back Office queue. Reconciliation is no longer a monthly ritual; it is a standing invariant with a lag metric.
- Real-time cash position becomes a queryable signal, published to Emily Prime (and readable by Jon Stockwell as *observation only*, per SYSTEM_GOVERNANCE — KAREN publishes, it does not take trading direction).

### Phase 2 — Payment Initiation

- Move from reading money to moving it: ACH, wire, RTP via banking APIs (direct bank API where available; a BaaS/payment-ops layer like Modern Treasury, Increase, or Column as the pragmatic bridge).
- **Hard invariant:** no agent initiates external money movement autonomously. Every payment is a *proposal* event (`payment_proposed`) requiring human approval through the Android companion app's biometric gate before a `payment_approved` event unlocks execution. This is exactly the escalation-control surface the companion app was scoped for.
- Dual representation: every executed payment exists as (a) a ledger event chain, (b) a QBO transaction, (c) a bank-confirmed transaction. Three-way match or it goes to the exception queue.

### Phase 3 — Counterparty Standing

- Goal: EINHORN as a counterparty other institutions treat seriously — settlement discipline, clean books on demand, operational maturity.
- **Honest constraint, stated plainly:** direct Fedwire/FedACH access requires being a chartered institution. Prime brokerage and ISDA relationships require capital, track record, and legal entity maturity. The realistic path runs through sponsor banks, a payments partner, and possibly licensing (money transmitter or beyond) — each a legal/regulatory decision made with actual counsel, not a spec document. This spec defines the *operational readiness* to deserve those relationships; it does not conjure them.
- The Back Office audit trail (Apples-logged, Iduna-enforced) is the asset here. Counterparty due diligence should be answerable in hours, not weeks.

## 3. Governance & Agent Boundaries

Consistent with SYSTEM_GOVERNANCE.md:

- **KAREN** publishes ledger observations and reconciliation state; proposes journal entries and payments; never commands other agents.
- **Emily Prime** issues directives to KAREN and is the sole path for human-originated financial instructions.
- **Jon Stockwell** may *read* cash position and settlement status as published signals. He may not instruct KAREN. Capital doctrine consumes the ledger; it does not drive it.
- **Iduna** holds all credentials (QBO OAuth tokens, bank API keys) and scopes them per-agent: `ledger:read`, `ledger:propose`, `ledger:execute`. `ledger:execute` is granted only to the execution path behind the human approval gate — no agent holds it as a standing permission.
- **Apples** receives a golden-log entry for every approved payment, every period close, and every reconciliation exception resolution. `ApplePublished` events surface in Back Office.

## 4. Invariants (Iduna-enforced)

1. Every QBO mutation has a preceding intent event in `var/ledger`. No silent writes.
2. No external money movement without a biometric-approved `payment_approved` event. No exceptions, no batch waivers.
3. Reconciliation exceptions older than N days (default 5 business days) block period close.
4. The event store is append-only. Corrections are new events (reversing entries), never edits — same doctrine as FATBABY's store.
5. Out-of-band QBO changes (detected by `qbowatch`, not originated by KAREN) are flagged and require acknowledgment, not silently absorbed.
6. Credentials never live in agent code or prompts. Iduna issues short-lived scoped tokens. (Bill & Ted convention applies: configuration arrives via environment; the time travel already happened.)

## 5. Back Office Metrics

- **Reconciliation lag:** median hours between `bank_txn_observed` and matched QBO transaction.
- **Exception queue depth & age:** unmatched items, oldest item age.
- **Proposal→approval funnel:** agent-proposed journal entries and payments vs. human-approved, with rejection reasons captured as events (this is training signal for KAREN's improvement loop, same oracle pattern as the EPS headline grader).
- **Close time:** days from period end to books closed. Northstar: continuous close — the metric trends toward zero.
- **Cash position freshness:** age of the newest confirmed bank observation.

## 6. v0 Build Order

1. `qbowatch` ingestor + `var/ledger` store (Go, same skeleton as `prwatch`).
2. KAREN agent skeleton with read-only QBO tools (`qbo_get_accounts`, `qbo_get_txns`, `qbo_get_reports`) registered in Iduna under `ledger:read`.
3. Intent-event → QBO write path for the two highest-volume entity types (likely Invoice and Bill), gated at `ledger:propose` with Emily Prime approval closing the loop.
4. Cash-position signal published to Back Office.
5. Only then: bank read rails (Phase 1).

## 7. Open Questions

- **Controller's name — RESOLVED 2026-07-04.** The Controller is **KAREN**. Registered in Iduna as `karen`; speaks to the manager (Emily Prime) and no one else.
- **Legal entity structure.** Which entity holds the QBO file and bank accounts? Phase 2+ decisions (licensing, sponsor bank) depend on this and belong with counsel.
- **QBO tier & API limits.** QBO's API throttles (per-minute request caps) may constrain the continuous-reconciliation cadence; needs a rate-budget the same way the Go agent core does token buckets.
- **Multi-entity future.** If EINHORN spawns entities, does each get its own QBO file with a consolidating layer, or do we accelerate toward an in-house GL sooner than planned?

---

*CLEAN BUILDS FIRST.*
