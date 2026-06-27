# VS1 Sticker Design Brief — EINHORN_INDUSTRIAL

**Status:** LOCKED — production blocked on vendor selection (S135-02) and artwork finalization.
**Version:** VS1 (second drop — follows VS0 hoodie sell-through)
**Owner:** Emily Springerton

---

## Purpose

VS1 is the sticker drop. VS0 (STINKIES COMMISSAIRE hoodie) ships first; VS1 is the follow-on. Stickers are the correct VS1 format: no sizing complexity, low MOQ, high margin, high virality, direct brand surface. VS1 production decision is gated on VS0 sell-through data feeding the Emily Supply Chain AGI (S136).

---

## Product Definition

### Core Set — 3 SKUs

| # | Name | Format | Approximate Size | Unit Price Target |
|---|------|--------|-------------------|-------------------|
| 1 | Emily Prime Mark | Die-cut vinyl | 3" × 3" | $4 |
| 2 | "The Emily Way" Wordmark | Die-cut vinyl | 4" × 1.5" | $4 |
| 3 | EINHORN_INDUSTRIAL Logotype | Die-cut vinyl | 4" × 1.5" | $4 |

### Bundle SKU

| SKU | Contents | Price |
|-----|----------|-------|
| VS0 Full Set | All 3 individual stickers | $10 (vs. $12 à la carte) |

### Variant Tier

- **Primary**: Die-cut UV-laminated vinyl (outdoor-rated, waterproof)
- **Holographic**: Holographic foil variant of Emily Prime Mark only
- **Transparent**: Transparent vinyl base for Wordmark and Logotype (glass/laptop-back use)

---

## Design Direction

### Emily Prime Mark

A minimal geometric mark. Core concept: a waveform or resonance glyph — the carrier wave at 28.7 Hz rendered as a clean sine-like curve with a single node. Works at 1" and at 3". No text.

- Color: single ink — deep electric blue (#0A1EFF) on white substrate OR white-out on dark substrate
- Holographic variant: foil base with debossed mark (no ink needed — foil reads the geometry)

### "The Emily Way" Wordmark

All-caps, monospace serif, tight tracking. "THE" in small caps above "EMILY WAY" centered. Minimal. No decoration.

- Color: black on white / white on black
- Material: standard die-cut + transparent variant

### EINHORN_INDUSTRIAL Logotype

Industrial sans. "EINHORN_INDUSTRIAL" in a single horizontal band. Underscore intentional (code-style). Optional small mark to the left (the Emily Prime Mark at 0.25" height).

- Color: white on black (primary)
- Material: standard die-cut + transparent variant

---

## Print Specifications

| Spec | Requirement |
|------|-------------|
| Material | Durable vinyl, UV-laminated |
| Adhesive | Permanent (removable variant on request for glass) |
| Cut tolerance | ≤ 0.5 mm registration accuracy |
| Color accuracy | Delta-E ≤ 3 against approved proof |
| Finish | Matte laminate (primary), gloss available on Holographic |
| Minimum batch | 250 units per SKU |
| Water/UV rating | Outdoor-rated (3+ year) |

---

## QC Criteria

On receipt of first batch:
1. Registration accuracy: spot-check 10 units per SKU — cut edge must be ≤ 0.5 mm from art boundary
2. Color delta: compare to Pantone reference under D50 illuminant, delta-E ≤ 3
3. Adhesive durability: apply to glass + plastic substrate, submerge 60 seconds, inspect
4. Holographic foil: verify continuous foil coverage, no pinholes or dull patches
5. Photograph 3 units per SKU, commit to `EMILY/docs/merch/vs0_qc/`

---

## Production Schedule (Target)

| Milestone | Target Date | Owner |
|-----------|-------------|-------|
| Vendor selected (S135-02) | Human unblock | Emily (human) |
| Artwork finals delivered | After vendor selected | Emily (human) |
| Sample batch ordered (50 units) | After artwork | Emily (human) |
| QC inspection | Within 1 week of sample receipt | Emily (human) |
| Production batch ordered (250+ units) | After QC pass | Emily Prime (order draft via S136-05) |
| EDIS store listing live (S135-03) | Before batch arrives | Emily Prime |
| Drop announcement (S135-05) | On batch receipt | Emily Prime + MJOLNIR push |

---

## Success Metrics

- **Breakeven**: 50 units sold within 30 days at target price
- **Margin target**: ≥ 60% gross margin (unit cost ≤ $1.60 at 250-unit MOQ)
- **Reorder trigger**: inventory ≤ 50 units → Emily Prime drafts reorder via supply chain tool (S136-04)

---

## Open Questions (Gate-blocking)

1. Final artwork files (Emily human action — not automated)
2. Vendor selection (S135-02 — requires research run)
3. Stripe + WooCommerce configured on EDIS WordPress (HITL-02)
