# NORTHSTAR — Baby ERP

**Classification:** Tier 1 — operational
**Owner:** Emily Springerton / EINHORN_INDUSTRIAL
**Version:** VS0 (bootstrap) | Written 2026-07-19

---

## Mission

One system of record for running both the business and the founder's life, pre-revenue, with
zero external SaaS dependency. Not a toy — a deliberately small, Apple-audited ERP, built on
infrastructure EINHORN_INDUSTRIAL already owns (IDUNA, Apples, emily-agent), shaped by studying
a real enterprise ERP's module boundaries without adopting the enterprise ERP itself.

Founder, verbatim: **"i need baby erp to help run my life and to help run the business as a
business we wont use any other software pre revenue but feel free to use the interface shadowing
pattern and use sap as a model to copy"** — clarified immediately after: **"except the emily
way"**, meaning EINHORN's own existing internal stack (IDUNA, Apples, emily-agent, the whole
operating discipline) is the foundation this is built *on*, not something being replaced.
"No other software" means no third-party SaaS — SAP, NetSuite, Odoo included. We're studying
SAP's module taxonomy as a design reference ("interface shadowing" — see §2), not integrating
with SAP.

## Why now, not later

The trigger was concrete, not abstract: "I want to acquire refurb iPhone 12s, I want to start
tracking stuff like that." That's real procurement (already covered by
`NORTHSTAR_SUPPLY_CHAIN.md`'s VS2 scope — vendor discovery, PO drafting, fulfillment tracking,
all built and working). What's missing is what happens *after* a PO is received: nothing tracks
how many units are actually in hand, where, or what condition. That gap — inventory — is real
ERP territory, and pulling on it surfaces the rest: as STINKIES moves from VS0 (one hoodie SKU)
toward VS3 (Store 0, a real multi-SKU physical location), the founder will need exactly what
ERP systems exist to provide. Better to shape that now, deliberately, than retrofit it later.

## 1. What EINHORN_INDUSTRIAL already has (don't rebuild these)

| SAP-equivalent module | What we call it | Status |
|---|---|---|
| FI (Financial Accounting) + CO (Controlling) | **KAREN** — Controller agent, QuickBooks Online integration, journal entries, bank reconciliation | Specced (`HQ-SPEC-FIN-098`, `HQ-SPEC-FIN-099`), not yet built |
| MM (Materials Management) — procurement half | **Emily Supply Chain AGI** — vendor discovery, PO drafting, fulfillment tracking, reorder intelligence | Built and working (`NORTHSTAR_SUPPLY_CHAIN.md`, SECTION 136) |

Baby ERP does not re-implement either. It's the module taxonomy *around* them — the pieces SAP
has that we don't, plus the one piece SAP was never built for (a person's actual life).

## 2. The interface-shadowing pattern

"Shadowing" SAP means: study how SAP separates concerns at the module boundary — what a
**Material Master** record actually is, what a **Purchase Requisition → Purchase Order → Goods
Receipt** state machine actually looks like, what **Quality Notification** actually captures —
and define our own Go interfaces/data shapes that mirror that *concept*, in our own naming, on
our own infrastructure (IDUNA tables, Apple audit events, eventstore where append-only history
matters). No SAP API calls, no SAP license, no SAP data model verbatim. The value of copying is
that these module boundaries are the result of decades of real enterprise operating experience
about where the seams actually belong — reinventing that from scratch, alone, badly, is the
failure mode this avoids.

Concretely: each new module below gets an `internal/erp/<module>` package (or IDUNA
`internal/<module>` where the module needs to be a trust-authority concern) with a doc comment
naming its SAP-module inspiration and where it deliberately diverges.

## 3. Module map — SAP-inspired (business side)

| Module | SAP name | Scope for Baby ERP VS0 | Depends on |
|---|---|---|---|
| **INV** | MM (materials mgmt, inventory half) | Physical stock: what's on hand, where, condition (new/refurb-grade/damaged), received-against-PO linkage. The actual gap that triggered this doc. | Supply Chain AGI's `supply_orders` (receiving closes a PO into stock) |
| **QM** | QM (Quality Management) | Formalizes what `NORTHSTAR_SUPPLY_CHAIN.md` already sketches informally (QC criteria, pass/fail, Apple-filed) into a real inspection-record schema, reusable beyond STINKIES merch (e.g. grading refurb iPhones: A/B/C cosmetic grade, functional test checklist) | INV (an inspection is *of* a stock item) |
| **SD** | SD (Sales & Distribution) | Order-to-cash, once there's real revenue to run through it — EDIS/WooCommerce orders, Store 0 POS eventually. **Deliberately deferred past VS0** — no real sales channel exists yet to model against. | INV (can't sell what isn't tracked) |
| **PM** | PM (Plant Maintenance) | Physical asset tracking (Store 0 fixtures/equipment once it exists). **Deferred past VS0** — no physical plant yet. | — |

## 4. Module map — not SAP-derived (personal side)

SAP has no concept of a founder's personal life; this half is invented fresh, held to the same
discipline (IDUNA-backed, Apple-audited where a decision/commitment is made) rather than copied
from an enterprise module that doesn't exist.

| Module | Scope for Baby ERP VS0 |
|---|---|
| **PL** (Personal Ledger) | Personal finances as their own ledger, separate from EINHORN's — not mixed into KAREN's business books. |
| **PT** (Personal Tasks) | Founder's own task/commitment tracking — distinct from `EMILY/BACKLOG.md` (engineering work) and HEIMDAL (product sprints); this is personal life admin. |

Scope for these two is intentionally thin at VS0 — the founder's own next input on what "help run
my life" actually needs day to day should shape this further, not a guess made here.

## 5. What VS0 actually ships

Deliberately narrow — one real module (INV), built because it's the concrete gap that started
this conversation, not the whole taxonomy at once:

1. **INV schema in IDUNA**: `stock_items` table — SKU, quantity on hand, location, condition
   grade, linked `supply_orders.order_id` (the PO it was received against). Migration + CRUD
   handler, same shape as `vendors`/`supply_orders` (SECTION 136 precedent).
2. **Receiving closes the loop**: when a `supply_orders` row transitions to `received`, Emily
   Prime creates/updates the corresponding `stock_items` row. Extends the existing fulfillment
   state machine rather than inventing a parallel one.
3. **`emily-agent` tool**: `erp_inventory_query` — read current stock by SKU/category, same
   pattern as `supply_chain_research`.
4. **Apple contract**: `stock_received` (completion), `stock_adjusted` (observation — manual
   correction, e.g. damage found), `stock_depleted` (observation — hits reorder threshold, same
   trigger Supply Chain AGI's reorder intelligence already expects to consume).

QM, SD, PM, PL, PT are named and scoped above so the taxonomy has a place for them, but none of
them ship in VS0 — each becomes its own backlog section when there's a real, concrete trigger for
it, same as INV had one here.

## 6. Human-in-the-loop points

Same discipline as Supply Chain AGI: Emily Prime does NOT autonomously —
1. Adjust recorded stock counts without a source event (a receiving, a manual correction, a sale)
2. Declare inventory QC-failed without human confirmation of physical inspection
3. Create a new module beyond what's specced here without a locked design brief

Vendor discovery, receiving-triggered stock updates, reorder threshold checks, and Apple filing
are automated, matching Supply Chain AGI's existing automation boundary.
