# HQ-SPEC-FIN-099 — EINHORN Ledger Platform (ELP)

**Status:** DRAFT v0 — pending Emily Prime review
**Custody:** Emily OS
**Extends:** HQ-SPEC-FIN-098 (Accounting Northstar). 098's phasing stands; this spec defines the full system that Phases 1–3 grow into.
**Northstar:** Functional parity with the three bridge partners — Modern Treasury (payment ops + ledger-as-a-service), Increase (programmatic bank accounts + rails), Column (developer-infrastructure bank) — built in-house, agent-native, Iduna-governed. The bridge partners are scaffolding; ELP is the building.

---

## 0. What Parity Actually Means

The three partners occupy three different layers. Parity is layer-by-layer, not a single finish line:

| Layer | Reference | Parity target | Hard constraint |
|---|---|---|---|
| **Ledger & payment ops** | Modern Treasury | Full parity, in-house software | None. This is pure software. Buildable today. |
| **Programmatic accounts & rails** | Increase | Parity via sponsor bank first, direct later | Rails access requires a bank in the loop until/unless EINHORN pursues charter or acquires one |
| **Being the bank** | Column | Long-horizon option, not a commitment | Charter, capital requirements, regulators, counsel. A board-level decision, out of spec scope |

The strategy: build the Modern Treasury layer completely (it's the brain), rent the Increase layer (rails via sponsor/BaaS) while designing every interface as if the rails were ours, and keep the Column question open with clean books as the down payment.

## 1. System Overview

ELP is five engines over one event store:

```
                    ┌─────────────────────────────┐
                    │        Emily Prime           │  directives, human interface
                    └──────────┬──────────────────┘
                               │
        ┌──────────────────────┼──────────────────────┐
        │                KAREN  ·  Controller           │
        │   proposes · reconciles · reports · never   │
        │            moves money on its own           │
        └──┬────────┬────────┬────────┬────────┬─────┘
           │        │        │        │        │
      ┌────▼──┐ ┌──▼────┐ ┌─▼─────┐ ┌▼──────┐ ┌▼────────┐
      │ Core  │ │Account│ │Payment│ │ Recon │ │ Policy &│
      │Ledger │ │System │ │Engine │ │Engine │ │Approvals│
      └────┬──┘ └──┬────┘ └─┬─────┘ └┬──────┘ └┬────────┘
           │       │        │        │         │
      ┌────▼───────▼────────▼────────▼─────────▼────┐
      │   Append-only event store (var/ledger)      │  NDJSON, same doctrine as FATBABY
      │   + Postgres projections (pgx/v5 pattern)   │
      └──────────────────┬───────────────────────────┘
                         │
              ┌──────────▼──────────┐
              │    Rails adapters   │  qbo · bank-read · baas · (future: direct)
              └─────────────────────┘
```

Everything below KAREN is deterministic Go services. KAREN is the only LLM-driven component; it consumes projections and emits proposal events. Iduna scopes every arrow.

## 2. Core Ledger

The heart. Double-entry, immutable, multi-ledger.

**Entities:**
- **Ledger** — one per (legal entity × currency). E.g., `einhorn-us-usd`. Consolidation is a projection across ledgers, never a merged ledger.
- **LedgerAccount** — typed (`asset`, `liability`, `equity`, `revenue`, `expense`), hierarchical (parent/child paths like `assets:cash:operating:mercury-checking`), with `normal_balance` (debit/credit).
- **LedgerTransaction** — atomic set of ≥2 LedgerEntries that must balance to zero per ledger per currency. Carries `effective_at` (economic time) distinct from `recorded_at` (system time) — backdating is legal, mutation is not.
- **LedgerEntry** — (account, direction, amount, currency). Amounts are integer minor units. No floats anywhere in the money path. Ever.

**Semantics:**
- **Immutability:** transactions are never edited or deleted. Corrections are reversing transactions linked via `reverses_id`. (Invariant 4 from FIN-098, now structural.)
- **Transaction states:** `pending` → `posted` | `archived`. Pending transactions support the authorization/hold pattern (card auths, payment orders awaiting settlement): balances track `pending_debits`/`pending_credits` separately from `posted`.
- **Balance model per account:** `posted_balance`, `pending_balance`, `available_balance` (posted minus pending debits plus configurable rules). Balances are projections, recomputable from entries — the entries are truth, balances are cache.
- **Balance locks:** optimistic concurrency via account version numbers; a transaction can assert `available_balance >= X` at post time and fail cleanly if violated. This is how overdraft-prevention works without a global lock.
- **Idempotency:** every write API takes an idempotency key, stored with request hash; replays return the original result. Non-negotiable for payment ops.

**Event vocabulary:** `ledger_transaction_created`, `ledger_transaction_posted`, `ledger_transaction_reversed`, `ledger_account_created`, `balance_lock_violated`.

## 3. Account System

Two distinct meanings of "account," kept rigorously separate:

**3a. Internal accounts (ours):**
- **BankAccount** — a real account at a real institution (or BaaS partner). Holds routing/account numbers, capabilities (`ach_origination`, `wire`, `rtp`, `book`), and a mapped LedgerAccount. Cash position = sum of bank-account-mapped ledger accounts, confirmed against `bank_txn_observed`.
- **VirtualAccount** — a sub-designation over a real account (FBO-style): unique virtual account numbers that route inbound funds to a specific LedgerAccount without opening real accounts. This is the Increase/Column primitive that makes one operating account behave like thousands. Phase: available via BaaS partner APIs first; the ELP data model treats partner-issued and (future) self-issued virtual accounts identically.
- **AccountHolder** — the legal entity owning internal accounts; hooks for the multi-entity future flagged in FIN-098 §7.

**3b. External accounts (theirs):**
- **Counterparty** — a registry entry: legal name, entity type, verification status, risk notes.
- **CounterpartyAccount** — payment coordinates (routing/account, wire instructions, RTP alias) with verification state (`unverified` → `microdeposit_pending` → `verified`) and provenance events. No payment order may target an unverified account above a policy threshold.

**Chart of accounts** ships as versioned seed data in the repo (reviewed like code), mapped to QBO's chart during the shadow period so both systems speak the same structure.

## 4. Payment Engine

Payment orders as explicit state machines. One model, many rails.

**PaymentOrder:** `{direction: credit|debit, amount, currency, originating_account, receiving_counterparty_account, rail, effective_date, remittance_info, metadata}`

**State machine:**
```
proposed → approved → queued → submitted → confirmed → settled
    │          │                    │           │
 rejected   expired             returned     failed
```
- `proposed` — created by KAREN (or Emily Prime directive). Emits `payment_proposed`.
- `approved` — human biometric approval via companion app (FIN-098 Invariant 2). Approval binds to an exact `(amount, counterparty_account, rail)` tuple hash — any mutation voids it.
- `submitted` — rails adapter has handed it to the bank/partner; carries external trace IDs (ACH trace, IMAD/OMAD for wires).
- `settled` — confirmed against an observed bank transaction by the Recon Engine. Not before.
- `returned` — ACH returns (R01–R85) modeled first-class, each return emits events and an automatic reversing ledger transaction plus an exception-queue item.

**Ledger choreography per order:** on `submitted`, post a pending transaction (cash → in-transit clearing account); on `settled`, post the clearing → final legs. In-transit money is always visible as a balance, never implicit.

**Rails adapters** (uniform interface `SubmitPayment`, `GetStatus`, `CancelIfPossible`):
1. `qbo` — Phase 0: "payment" = recorded transaction only, no money movement.
2. `bank-read` — Phase 1: observation only, feeds recon.
3. `baas` — Phase 2: Modern Treasury / Increase / Column API-backed execution. **All three implemented behind the same interface** so the partner is a config choice, not an architecture choice. This is the parity forcing-function: our internal API must already express everything theirs can.
4. `direct` — Phase 3+: NACHA file generation / FedWire via sponsor arrangements. The adapter interface doesn't change; only what's behind it.

**Scheduled & recurring orders:** cron-driven (existing autonomy pattern), but each occurrence still individually passes policy — recurrence is a proposal generator, not an approval bypass.

## 5. Reconciliation Engine

Continuous three-way match: **intent** (payment orders / expected receipts) × **books** (ledger transactions) × **reality** (`bank_txn_observed`).

- **ExpectedPayment** objects declare inbound funds we're waiting for (client invoice, trade settlement) with match criteria (amount range, counterparty, date window, remittance regex).
- Matcher runs on every new bank observation: exact match → auto-reconcile event; fuzzy match → suggestion queued for KAREN proposal → human confirm above threshold; no match → exception queue with aging (FIN-098 Invariant 3 blocks close on aged exceptions).
- Match confidence scoring is a learnable component: confirmed/rejected suggestions accumulate as labeled data — the same oracle-graded improvement loop as the EPS headline grader, with the same no-regression gate before any matcher update ships.
- **Statements:** period-end statement generation per ledger account (running balance, entries, opening/closing) — both a human artifact and the substrate for continuous close.

## 6. Policy & Approvals Engine

Where Iduna meets money. Declarative policy, evaluated on every state transition:

```yaml
- match: {type: payment_order, rail: wire}
  require: [biometric_approval, second_approver]
- match: {type: payment_order, amount_gte: 10000_00}
  require: [biometric_approval, counterparty_verified]
- match: {type: payment_order, counterparty: unverified}
  require: [deny]
- match: {type: ledger_transaction, source: karen}
  require: [emily_prime_ack]   # Phase 0 training wheels; relaxes per no-regression evidence
```

- Policies are versioned files in the repo under Emily OS custody (locked from Tyler's scopes, same as `CLAUDE.md` et al.).
- Every evaluation emits a `policy_evaluated` event with the matched rules — the audit answer to "why was this allowed" is a query, not an investigation.
- Iduna scopes remain the outer wall: `ledger:read`, `ledger:propose`, `ledger:execute`; execute held only by the post-approval execution path, never standing.
- **Compliance hooks (honest placeholder):** counterparty screening (sanctions/OFAC) and KYB slots are interface stubs in v0, fulfilled by partner APIs in Phase 2. Anything beyond that (BSA/AML program, money transmission analysis) is counsel territory — the system reserves the seams; lawyers decide what fills them.

**Amendment (append-only), 2026-07-23:** the policy engine above predates `pkg/norn`
(`HQ-SPEC-PRIME-101-norn-loop-kernel.md`), the now-canonical propose→grade→gate→promote kernel.
This section is not rewritten — the policy rules and evaluation semantics above remain
authoritative — but it should be read going forward as this domain's NORN instantiation rather
than a bespoke approval loop: PRIME-101 §6's instantiation table already carries "KAREN journal
proposals" (proposer: KAREN, oracle: Emily Prime accept/reject history, tier `prime_ack`, reality
root: human decision) — the exact shape of the `emily_prime_ack` training-wheels rule above (line
137), including the "relaxes per no-regression evidence" note, which is literally NORN's tier-
relaxation-as-gated-promotion concept (see PRIME-101's own text on `prime_ack`). As ELP-policy
tooling is built (Build Sequence step 4, Policy Engine), new policy rules should register as NORN
instantiations rather than accumulate as one-off `require:` entries with no shared grading/lineage
mechanism. Per PRIME-101 §8 build-sequence item 5.

## 7. API Surface

REST, versioned (`/v1/`), JSON, cursor pagination, idempotency keys on all POSTs, webhook/event subscriptions. Resource set mirrors the parity target deliberately:

`/ledgers` · `/ledger_accounts` · `/ledger_transactions` · `/bank_accounts` · `/virtual_accounts` · `/counterparties` · `/counterparty_accounts` · `/payment_orders` · `/expected_payments` · `/reconciliation_exceptions` · `/statements` · `/policies` · `/events`

Consumers: KAREN (primary), Back Office UI, the Android companion app (approvals + read), Jon Stockwell (read-only cash/settlement signals per SYSTEM_GOVERNANCE), external counterparties eventually (scoped read tokens for their own statements — the due-diligence-as-API move from FIN-098 Phase 3).

**Streaming:** ledger and payment events publish to the existing feedserver/tcpstreamsdk pattern; Back Office subscribes live. CAST-style separation applies: streaming failures must be invisible to the core money path.

## 8. Build Sequence (extends FIN-098 §6)

1. **Core Ledger service** — accounts, balanced transactions, balances, idempotency, reversals. Pure Go + Postgres, no external deps. This is the piece that must be perfect; everything else forgives bugs better than the ledger does. Property-based tests: entries always balance, balances always recompute, replays always idempotent.
2. **Shadow the books:** QBO stays system of record (FIN-098 Phase 0); ELP ledger runs in parallel, populated by the same intent events. Divergence between ELP and QBO is itself a recon exception. Promotion criterion: N consecutive clean periods (propose N=3) of zero unexplained divergence — then ELP becomes system of record and QBO demotes to tax/accountant mirror.
3. **Counterparty registry + ExpectedPayments + Recon Engine** over Phase 1 bank-read rails.
4. **Payment Engine + Policy Engine + companion-app approval flow**, executing through one BaaS adapter (pick the first partner on rails coverage + API quality for our actual bank).
5. **Second and third rails adapters** — the parity proof: if our interface can't express a partner's capability, the interface grows until it can.
6. **Virtual accounts, statements-as-API, counterparty read tokens.**
7. **Direct-rails design study** (NACHA generation, sponsor bank requirements) — gated on volume justifying it.

## 9. Back Office Metrics (extends FIN-098 §5)

- Ledger↔QBO divergence count (shadow period) — the promotion gate.
- Payment order cycle time by state (proposed→approved, submitted→settled) per rail.
- Return rate by rail and counterparty; exception queue depth/age (carried from 098).
- Matcher precision/recall on the confirmed-suggestion eval set.
- Approval-policy hit distribution (which rules fire, which deny).
- Available balance accuracy: predicted vs. bank-confirmed at settlement.

## 10. Open Questions

- **Controller's name — RESOLVED 2026-07-04:** **KAREN**. Five engines, one agent, zero tolerance for unmatched transactions.
- **First BaaS partner** — decision inputs: which banks we actually hold accounts at, rail coverage, API quality, exit cost. Adapter architecture makes this reversible, but the first integration sets the tempo.
- **Multi-currency timeline** — the ledger model supports it structurally (ledger per entity×currency); do we ship USD-only v0 or pay the FX-representation cost now?
- **Ledger service isolation** — same process space as FATBABY services, or its own hardened deployment with a stricter release cadence? Instinct says the money path earns its own blast radius.
- **The Column question** — revisit only when settlement volume, capital, and counsel all say the same thing. Until then it stays a northstar, not a roadmap item.

---

*CLEAN BUILDS FIRST.*
