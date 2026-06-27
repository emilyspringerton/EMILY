# NORTHSTAR — Emily Supply Chain AGI

**Classification:** Tier 1 — operational  
**Owner:** Emily Springerton / EINHORN_INDUSTRIAL  
**Version:** VS0 (bootstrap) | Updated: 2026-06-27

---

## Mission

Emily manages a transparent, Apple-audited physical product supply chain. Every vendor interaction, purchase order, and fulfillment event is traceable to a source Apple. The goal is to eliminate the need for Emily (human) to manually discover, evaluate, negotiate with, or manage vendors. Emily Prime does it.

---

## Scope

| Version | What it covers |
|---------|---------------|
| VS0 | Sticker drop — vinyl print-on-demand (S135) |
| VS1 | Broader merch (shirts, pins, patches, packaging) |
| VS2 | Arbitrary physical product sourcing (any SKU) |

This document governs VS0–VS2.

---

## Core Components

### 1. Vendor Discovery

Emily Prime uses the Research tool (S137-01) to find and evaluate vendors for a given product category and specification. Research queries are logged as `research_log` Apples with full source provenance.

Input: `SupplySpec{product, category, quantity, quality_tier, budget_cents_per_unit}`  
Output: `[]VendorOption{vendor_id, name, url, moq, unit_cost_cents, lead_days, quality_tier, confidence}`

Discovery sources (priority order):
1. IDUNA vendor registry (known vendors from prior runs)
2. Research engine (S137-01) — web fetch from supplier directories, subreddit recommendations, industry directories
3. Manual addition via Emily human override

### 2. Quality Specification

Each supply order references a `SupplySpec` with measurable QC criteria:
- Print: delta-E, cut tolerance, adhesive spec
- Apparel: thread count, shrinkage, color accuracy
- General: dimensional tolerance, material certification

The spec is the contract. All QC outcomes are filed as Apples against the spec ID.

### 3. Quote Negotiation Loop

Emily Prime drafts purchase inquiry emails (or API calls for vendors with programmatic interfaces). The draft is queued in HEIMDAL as a sprint for human approval. Once approved, Emily files the order in IDUNA `supply_orders` and marks the sprint complete.

Cycle: `Research → Draft PO → HEIMDAL approval sprint → Order placed → Track fulfillment`

### 4. Fulfillment Tracking

Status machine: `pending → ordered → shipped → received → qc_pass | qc_fail`

Each transition: Emily Prime calls IDUNA to update the order status and files a state-change Apple. On `qc_fail`, Emily Prime files a `research_log` Apple with failure details and routes a new vendor comparison.

### 5. Reorder Intelligence

Reorder threshold is stored per SKU in IDUNA `supply_orders`. Emily Prime checks inventory counts (from EDIS WooCommerce sales data) each cycle against the threshold. Below threshold: draft reorder via `ResearchSupplyChain` → HEIMDAL approval sprint.

---

## Data Schema

### IDUNA: `vendors` table (S136-02)

| Column | Type | Description |
|--------|------|-------------|
| vendor_id | UUID | Primary key |
| name | VARCHAR(120) | Vendor display name |
| category | VARCHAR(64) | Product category (print, apparel, packaging) |
| url | TEXT | Homepage or order portal |
| moq | INTEGER | Minimum order quantity |
| unit_cost_cents | INTEGER | Cost per unit at MOQ in USD cents |
| lead_days | INTEGER | Production + shipping days |
| quality_tier | VARCHAR(16) | budget / standard / premium |
| last_evaluated_at | DATETIME | When Emily last scored this vendor |
| notes | TEXT | Free-form notes from research runs |
| status | VARCHAR(16) | active / inactive / disqualified |

### IDUNA: `supply_orders` table (S136-03)

| Column | Type | Description |
|--------|------|-------------|
| order_id | UUID | Primary key |
| vendor_id | UUID | FK to vendors |
| product | VARCHAR(120) | Product description |
| quantity | INTEGER | Units ordered |
| unit_cost_cents | INTEGER | Locked price at time of order |
| total_cost_cents | INTEGER | quantity × unit_cost_cents + shipping |
| status | VARCHAR(16) | pending / ordered / shipped / received / qc_pass / qc_fail |
| ordered_at | DATETIME | Null until human approves |
| received_at | DATETIME | Null until delivery confirmed |
| notes | TEXT | QC notes, failure details |

---

## Emily Prime Tool Surface (S136-04/05)

```go
// ResearchSupplyChain discovers and scores vendors for a given spec.
func ResearchSupplyChain(product string, spec SupplySpec) []VendorOption

// DraftPurchaseOrder generates a PO draft and queues a HEIMDAL approval sprint.
func DraftPurchaseOrder(vendor VendorOption, quantity int, spec SupplySpec) error
```

Both functions file Apples: `ResearchSupplyChain` files `research_log`; `DraftPurchaseOrder` files `observation` (action pending human approval).

---

## Apple Audit Contract

Every supply chain event produces an Apple:

| Event | Apple Type | Title Pattern |
|-------|-----------|---------------|
| Vendor research completed | research_log | `supply-chain: vendor research for {product}` |
| Order draft queued | observation | `supply-chain: PO draft pending approval — {vendor}` |
| Order placed | completion | `supply-chain: order placed — {vendor} × {qty} {product}` |
| Shipped | completion | `supply-chain: shipped — tracking {tracking_id}` |
| QC pass | completion | `supply-chain: QC pass — {product} batch {order_id}` |
| QC fail | escalation | `supply-chain: QC FAIL — {product} batch {order_id}: {reason}` |
| Reorder triggered | observation | `supply-chain: reorder triggered for {sku} (inventory below threshold)` |

---

## Revenue Integration

Supply chain cost basis is logged per order. IDUNA `supply_orders.total_cost_cents` feeds into the financial model. Margin reporting: Emily Prime compares unit cost to EDIS/WooCommerce revenue per SKU each weekly briefing.

Target margin by tier:
- VS0 (stickers): ≥ 60% gross
- VS1 (apparel): ≥ 45% gross
- VS2 (arbitrary): ≥ 30% gross

---

## Human-in-the-Loop Points

Emily Prime does NOT autonomously:
1. Submit payment (all orders require human approval in HEIMDAL)
2. QC physical goods (Emily human inspects and confirms)
3. Create new product categories without a locked design brief (per VS0 process)

All other steps — vendor discovery, quote drafting, order tracking, reorder triggering, Apple filing — are fully automated.
